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

func TestADFFeederSideElement_RoundTrip(t *testing.T) {
	orig := ADFFeederSideElement{
		ADFColor: []ColorEntry{BlackAndWhite1, RGB24},
		ADFMaximumSize: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "297"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "210"}},
		},
		ADFMinimumSize: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "100"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "50"}},
		},
		ADFOpticalResolution: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "600"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "600"}},
		},
		ADFResolutions: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "300"}, {Text: "600"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "300"}, {Text: "600"}},
		},
	}
	elm := orig.toXML("wscn:ADFBack")
	if elm.Name != "wscn:ADFBack" {
		t.Errorf("expected element name 'wscn:ADFBack', got '%s'", elm.Name)
	}

	parsed, err := decodeADFFeederSideElement(elm)
	if err != nil {
		t.Fatalf("decodeADFFeederSideElement returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, parsed) {
		t.Errorf("expected %+v, got %+v", orig, parsed)
	}
}

func TestADFFeederSideElement_MissingRequired(t *testing.T) {
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
	_, err := decodeADFFeederSideElement(elm)
	if err == nil {
		t.Errorf("expected error for missing required elements, got nil")
	}
}

func TestADFFeederSideElement_InvalidColorEntry(t *testing.T) {
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
	_, err := decodeADFFeederSideElement(elm)
	if err == nil {
		t.Errorf("expected error for invalid color entry, got nil")
	}
}
