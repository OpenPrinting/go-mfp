// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Platen

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestPlaten_RoundTrip(t *testing.T) {
	orig := Platen{
		PlatenColor: []ColorEntry{BlackAndWhite1, RGB24},
		PlatenMaximumSize: Dimension{
			Width:  210,
			Height: 297,
		},
		PlatenMinimumSize: Dimension{
			Width:  50,
			Height: 100,
		},
		PlatenOpticalResolution: Dimension{
			Width:  600,
			Height: 600,
		},
		PlatenResolutions: Resolutions{
			Widths:  []int{300, 600},
			Heights: []int{300, 600},
		},
	}
	elm := orig.toXML("wscn:Platen")
	if elm.Name != "wscn:Platen" {
		t.Errorf("expected element name 'wscn:Platen', got '%s'", elm.Name)
	}

	parsed, err := decodePlaten(elm)
	if err != nil {
		t.Fatalf("decodePlaten returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, parsed) {
		t.Errorf("expected %+v, got %+v", orig, parsed)
	}
}

func TestPlaten_MissingRequired(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Platen",
		Children: []xmldoc.Element{
			// Only PlatenColor present
			{
				Name: NsWSCN + ":PlatenColor",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":ColorEntry", Text: "BlackAndWhite1"},
				},
			},
		},
	}
	_, err := decodePlaten(elm)
	if err == nil {
		t.Errorf("expected error for missing required elements, got nil")
	}
}

func TestPlaten_InvalidColorEntry(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Platen",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":PlatenColor",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":ColorEntry", Text: "NotAColor"},
				},
			},
			{
				Name: NsWSCN + ":PlatenMaximumSize",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "297"},
					{Name: NsWSCN + ":Width", Text: "210"},
				},
			},
			{
				Name: NsWSCN + ":PlatenMinimumSize",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "100"},
					{Name: NsWSCN + ":Width", Text: "50"},
				},
			},
			{
				Name: NsWSCN + ":PlatenOpticalResolution",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "600"},
					{Name: NsWSCN + ":Width", Text: "600"},
				},
			},
			{
				Name: NsWSCN + ":PlatenResolutions",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "300"},
					{Name: NsWSCN + ":Width", Text: "300"},
				},
			},
		},
	}
	_, err := decodePlaten(elm)
	if err == nil {
		t.Errorf("expected error for invalid color entry, got nil")
	}
}
