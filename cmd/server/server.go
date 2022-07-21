package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	txtTemplate "text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/tortlewortle/srt2subtitles/pkg/srt"
)

//go:embed final_cut_project.goxml
var xmltemplate string

//go:embed upload.gohtml
var htmltemplate string

var tmpl = txtTemplate.Must(txtTemplate.New("xml").Parse(xmltemplate))
var indexTmpl = template.Must(template.New("index").Parse(htmltemplate))

type TemplateVars struct {
	UUID      string
	Name      string
	Duration  int
	Framerate float64
	Segments  []TmplSegment
	Width     int
	Height    int
}

type TmplSegment struct {
	ID         int64
	StartFrame int
	EndFrame   int
	Text       string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexTmpl.Execute(w, nil)
	})

	r.Post("/upload", uploadHandler)
	http.ListenAndServe(":"+port, r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	inheight := r.PostFormValue("height")
	inwidth := r.PostFormValue("width")
	inrate := r.PostFormValue("framerate")

	if inheight == "" {
		inheight = "1080"
	}

	if inwidth == "" {
		inwidth = "1920"
	}

	if inrate == "" {
		inrate = "24"
	}

	var width int64 = 1920
	var height int64 = 1080

	// ignore errors as defaults have been set
	height, _ = strconv.ParseInt(inheight, 10, 0)
	width, _ = strconv.ParseInt(inwidth, 10, 0)

	framerate, err := strconv.ParseFloat(strings.Replace(inrate, ",", ".", -1), 64)
	// just default to sane figure if it fails ig
	if err != nil {
		framerate = 24
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) != ".srt" {
		http.Error(w, "not srt", http.StatusBadRequest)
		return
	}

	segments, err := srt.ReadSRT(file)
	if err != nil {
		log.Fatalf("gimmsegments: %v", err)
	}

	var tmplSegments []TmplSegment

	for _, segment := range segments {
		tmplSegments = append(tmplSegments, TmplSegment{
			ID:         segment.ID,
			StartFrame: int(segment.Start * framerate),
			EndFrame:   int(segment.End * framerate),
			Text:       segment.Text,
		})
	}

	rawName := handler.Filename[:len(handler.Filename)-len(filepath.Ext(handler.Filename))]
	id := uuid.New()
	vars := TemplateVars{
		UUID:      id.String(),
		Name:      rawName,
		Framerate: framerate,
		Duration:  int(framerate * segments[len(segments)-1].End),
		Segments:  tmplSegments,
		Width:     int(width),
		Height:    int(height),
	}

	w.Header().Set("content-type", "application/octet-stream")
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s.xml\"", rawName))
	tmpl.Execute(w, vars)
}
