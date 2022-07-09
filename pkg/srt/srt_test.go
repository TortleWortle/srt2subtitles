package srt

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed test/premieresample.srt
var premiereSample []byte

//go:embed test/manualsample.srt
var manualsample []byte

var manualSegments = []Segment{
	{
		ID:    1,
		Start: 0.498,
		End:   2.827,
		Text:  "- Here's what I love most\nabout food and diet.",
	},
	{
		ID:    2,
		Start: 2.827,
		End:   6.383,
		Text:  "We all eat several times a day,\nand we're totally in charge",
	},
	{
		ID:    3,
		Start: 6.383,
		End:   9.427,
		Text:  "of what goes on our plate\nand what stays off.",
	},
}

var premierSegments = []Segment{
	{
		ID:    1,
		Start: 0,
		End:   2.961,
		Text:  "hii",
	},
	{
		ID:    2,
		Start: 18.476,
		End:   21.438,
		Text:  "how are",
	},
	{
		ID:    3,
		Start: 41.416,
		End:   44.377,
		Text:  "you\ndoing",
	},
}

func TestReadSRTManual(t *testing.T) {
	segments, err := ReadSRT(bytes.NewBuffer(manualsample))
	if err != nil {
		t.Error(err)
	}
	if len(segments) != 3 {
		t.Errorf("insufficient segments expected: %v got: %v", 3, len(segments))
	}
	for i, seg := range segments {
		if seg.ID != manualSegments[i].ID {
			t.Errorf("ID doesn't match for segment: %v", manualSegments[i].ID)
		}
		if seg.Start != manualSegments[i].Start {
			t.Errorf("Start doesn't match for segment: %v", manualSegments[i].ID)
		}
		if seg.End != manualSegments[i].End {
			t.Errorf("End doesn't match for segment: %v", manualSegments[i].ID)
		}
		if seg.Text != manualSegments[i].Text {
			t.Errorf("Text doesn't match for segment: %v", manualSegments[i].ID)
		}
	}
}

func TestReadSRTPremiere(t *testing.T) {
	segments, err := ReadSRT(bytes.NewBuffer(premiereSample))
	if err != nil {
		t.Error(err)
	}
	if len(segments) != 3 {
		t.Errorf("insufficient segments expected: %v got: %v", 3, len(segments))
	}

	for i, seg := range segments {
		if seg.ID != premierSegments[i].ID {
			t.Errorf("ID doesn't match for segment: %v", premierSegments[i].ID)
		}
		if seg.Start != premierSegments[i].Start {
			t.Errorf("Start doesn't match for segment: %v", premierSegments[i].ID)
		}
		if seg.End != premierSegments[i].End {
			t.Errorf("End doesn't match for segment: %v", premierSegments[i].ID)
		}
		if seg.Text != premierSegments[i].Text {
			t.Errorf("Text doesn't match for segment: %v", premierSegments[i].ID)
		}
	}
}
