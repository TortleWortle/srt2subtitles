package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"

	"github.com/google/uuid"
	"github.com/sqweek/dialog"
)

//go:embed template.goxml
var xmltemplate string

const name = "OverStimulationRoomCCStylized.srt_BQOUGHCLWH4iqz3iSbYz_1655388165"
const framerate = 24

type TemplateVars struct {
	UUID      string
	Name      string
	Duration  int
	Framerate int
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

var timestampRegex = regexp.MustCompile("([0-9][0-9]):([0-9][0-9]):([0-9][0-9],[0-9][0-9][0-9]) --> ([0-9][0-9]):([0-9][0-9]):([0-9][0-9],[0-9][0-9][0-9])")

var tmpl = template.Must(template.New("xml").Parse(xmltemplate))

func main() {
	filename, err := dialog.File().Filter("Srt file", "srt").Load()
	if err != nil {
		log.Fatalf("opening dialog: %v", err)
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}
	defer file.Close()

	segments, err := ReadSRT(file)
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
	id := uuid.New()
	vars := TemplateVars{
		UUID:      id.String(),
		Name:      name,
		Framerate: framerate,
		Duration:  int(framerate * segments[len(segments)-1].End),
		Segments:  tmplSegments,
		Width:     1920,
		Height:    1080,
	}
	filename, err = dialog.File().Filter("xml file", "xml").Title("Export to xml").Save()
	if err != nil {
		log.Fatalf("opening output dialog: %v", err)
		return
	}
	outfile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("outfile: %v", err)
	}
	defer outfile.Close()

	outfile.Truncate(0)

	fmt.Println(tmpl.Execute(outfile, vars))
}
