// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for SideElement (ADFBack/ADFFront)

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestADFFeederSide_RoundTrip(t *testing.T) {
	orig := ADFFeederSide{
		ADFColor: []ColorEntry{BlackAndWhite1, RGB24},
		ADFMaximumSize: Dimension{
			Width:  210,
			Height: 297,
		},
		ADFMinimumSize: Dimension{
			Width:  50,
			Height: 100,
		},
		ADFOpticalResolution: Dimension{
			Width:  600,
			Height: 600,
		},
		ADFResolutions: Dimension{
			Width:  300,
			Height: 300,
		},
	}
	elm := orig.toXML("wscn:ADFBack")
	if elm.Name != "wscn:ADFBack" {
		t.Errorf("expected element name 'wscn:ADFBack', got '%s'", elm.Name)
	}

	parsed, err := decodeADFFeederSide(elm)
	if err != nil {
		t.Fatalf("decodeADFFeederSide returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, parsed) {
		t.Errorf("expected %+v, got %+v", orig, parsed)
	}
}

func TestADFFeederSide_MissingRequired(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:ADFBack",
		Children: []xmldoc.Element{
			// Only ADFColor present
			{
				Name: NsWSCN + ":ADFColor",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":ColorEntry", Text: "BlackAndWhite1"},
				},
			},
		},
	}
	_, err := decodeADFFeederSide(elm)
	if err == nil {
		t.Errorf("expected error for missing required elements, got nil")
	}
}

func TestADFFeederSide_InvalidColorEntry(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:ADFBack",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ADFColor",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":ColorEntry", Text: "NotAColor"},
				},
			},
			{
				Name: NsWSCN + ":ADFMaximumSize",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "297"},
					{Name: NsWSCN + ":Width", Text: "210"},
				},
			},
			{
				Name: NsWSCN + ":ADFMinimumSize",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "100"},
					{Name: NsWSCN + ":Width", Text: "50"},
				},
			},
			{
				Name: NsWSCN + ":ADFOpticalResolution",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "600"},
					{Name: NsWSCN + ":Width", Text: "600"},
				},
			},
			{
				Name: NsWSCN + ":ADFResolutions",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "300"},
					{Name: NsWSCN + ":Width", Text: "300"},
				},
			},
		},
	}
	_, err := decodeADFFeederSide(elm)
	if err == nil {
		t.Errorf("expected error for invalid color entry, got nil")
	}
}
