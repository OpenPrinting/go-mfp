// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Resolutions

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestResolutions_SingleValues(t *testing.T) {
	orig := Resolutions{
		Widths:  []int{210},
		Heights: []int{297},
	}

	elm := orig.toXML("wscn:Resolutions")
	if elm.Name != "wscn:Resolutions" {
		t.Errorf("expected element name 'wscn:Resolutions', got '%s'", elm.Name)
	}

	// Should have 2 children: Widths and Heights containers
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}

	decoded, err := decodeResolutions(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestResolutions_MultipleValues(t *testing.T) {
	orig := Resolutions{
		Widths:  []int{210, 297, 420},
		Heights: []int{297, 420, 594},
	}

	elm := orig.toXML("wscn:Resolutions")
	if elm.Name != "wscn:Resolutions" {
		t.Errorf("expected element name 'wscn:Resolutions', got '%s'", elm.Name)
	}

	// Should have 2 children: Widths and Heights containers
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}

	decoded, err := decodeResolutions(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestResolutions_MixedValues(t *testing.T) {
	orig := Resolutions{
		Widths:  []int{210},
		Heights: []int{297, 420, 594},
	}

	elm := orig.toXML("wscn:Resolutions")
	decoded, err := decodeResolutions(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestResolutions_EmptyArrays(t *testing.T) {
	orig := Resolutions{
		Widths:  []int{},
		Heights: []int{},
	}

	elm := orig.toXML("wscn:Resolutions")
	_, err := decodeResolutions(elm)
	if err == nil {
		t.Errorf("expected error for empty arrays, got nil")
	}
}

func TestResolutions_MissingWidths(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Resolutions",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Heights",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "297"},
				},
			},
		},
	}

	_, err := decodeResolutions(elm)
	if err == nil {
		t.Errorf("expected error for missing Width elements, got nil")
	}
}

func TestResolutions_MissingHeights(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Resolutions",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Widths",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Width", Text: "210"},
				},
			},
		},
	}

	_, err := decodeResolutions(elm)
	if err == nil {
		t.Errorf("expected error for missing Height elements, got nil")
	}
}

func TestResolutions_InvalidWidthValue(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Resolutions",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Widths",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Width", Text: "invalid"},
				},
			},
			{
				Name: NsWSCN + ":Heights",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "297"},
				},
			},
		},
	}

	_, err := decodeResolutions(elm)
	if err == nil {
		t.Errorf("expected error for invalid width value, got nil")
	}
}

func TestResolutions_InvalidHeightValue(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Resolutions",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Widths",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Width", Text: "210"},
				},
			},
			{
				Name: NsWSCN + ":Heights",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "invalid"},
				},
			},
		},
	}

	_, err := decodeResolutions(elm)
	if err == nil {
		t.Errorf("expected error for invalid height value, got nil")
	}
}

func TestResolutions_WrappedForm(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Resolutions",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Widths",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Width", Text: "210"},
					{Name: NsWSCN + ":Width", Text: "297"},
				},
			},
			{
				Name: NsWSCN + ":Heights",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Height", Text: "297"},
					{Name: NsWSCN + ":Height", Text: "420"},
				},
			},
		},
	}

	expected := Resolutions{
		Widths:  []int{210, 297},
		Heights: []int{297, 420},
	}

	decoded, err := decodeResolutions(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Errorf("expected %+v, got %+v", expected, decoded)
	}
}
