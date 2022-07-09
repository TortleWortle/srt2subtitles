package srt

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Segment struct {
	ID    int64
	Start float64
	End   float64
	Text  string
}

var timestampRegex = regexp.MustCompile("([0-9][0-9]):([0-9][0-9]):([0-9][0-9],[0-9][0-9][0-9]) --> ([0-9][0-9]):([0-9][0-9]):([0-9][0-9],[0-9][0-9][0-9])")

func segmentStartEndTime(segmentStr string) (float64, float64, error) {
	matches := [][]string{timestampRegex.FindStringSubmatch(segmentStr)}

	startHours, err := strconv.ParseInt(matches[0][1], 10, 0)
	if err != nil {
		return 0, 0, fmt.Errorf("starthours: %v", err)
	}

	startMinutes, err := strconv.ParseInt(matches[0][2], 10, 0)
	if err != nil {
		return 0, 0, fmt.Errorf("startminutes: %v", err)
	}

	startSeconds, err := strconv.ParseFloat(strings.Replace(matches[0][3], ",", ".", -1), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("startseconds: %v", err)
	}

	endHours, err := strconv.ParseFloat(matches[0][4], 0)
	if err != nil {
		return 0, 0, fmt.Errorf("endhours: %v", err)
	}

	endMinutes, err := strconv.ParseFloat(matches[0][5], 0)
	if err != nil {
		return 0, 0, fmt.Errorf("endminutes: %v", err)
	}

	endSeconds, err := strconv.ParseFloat(strings.Replace(matches[0][6], ",", ".", -1), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("endseconds: %v", err)
	}

	start := float64(startHours*60*60+startMinutes*60) + startSeconds
	end := float64(endHours*60*60+endMinutes*60) + endSeconds

	return start, end, nil
}

func ReadSRT(file io.Reader) ([]Segment, error) {
	reader := bufio.NewReader(file)

	// skip first 3 bytes if it matches the premiere srt export
	b, err := reader.Peek(3)
	if err != nil {
		return nil, fmt.Errorf("peek err: %v", err)
	}
	if b[0] == 239 && b[1] == 187 && b[2] == 191 {
		buf := make([]byte, 3)
		reader.Read(buf)
	}

	scanner := bufio.NewScanner(reader)

	var segments []Segment

	for scanner.Scan() {
		idIn := scanner.Text()
		if idIn == "" {
			break
		}
		id, err := strconv.ParseInt(idIn, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("id \"%v\" not int", []byte(idIn))
		}
		scanner.Scan()
		timestamps := scanner.Text()
		var text []string
		for scanner.Scan() {
			t := scanner.Text()
			if t == "" {
				break
			}
			text = append(text, t)
		}

		start, end, err := segmentStartEndTime(timestamps)
		if err != nil {
			return nil, fmt.Errorf("startendtime: %v", err)
		}
		seg := Segment{
			ID:    id,
			Start: start,
			End:   end,
			Text:  strings.Join(text, "\n"),
		}
		segments = append(segments, seg)
	}

	return segments, nil
}
