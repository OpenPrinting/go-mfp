// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Dimension

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDimension_RoundTrip(t *testing.T) {
	orig := Dimension{
		Width:  210,
		Height: 297,
	}

	elm := orig.toXML("wscn:Dimension")
	if elm.Name != "wscn:Dimension" {
		t.Errorf("expected element name 'wscn:Dimension', got '%s'", elm.Name)
	}

	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}

	decoded, err := decodeDimension(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestDimension_ZeroValues(t *testing.T) {
	orig := Dimension{
		Width:  0,
		Height: 0,
	}

	elm := orig.toXML("wscn:Dimension")
	decoded, err := decodeDimension(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestDimension_MissingWidth(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Dimension",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Height",
				Text: "297",
			},
		},
	}

	_, err := decodeDimension(elm)
	if err == nil {
		t.Errorf("expected error for missing Width element, got nil")
	}
}

func TestDimension_MissingHeight(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Dimension",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "210",
			},
		},
	}

	_, err := decodeDimension(elm)
	if err == nil {
		t.Errorf("expected error for missing Height element, got nil")
	}
}

func TestDimension_InvalidWidth(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Dimension",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "invalid",
			},
			{
				Name: NsWSCN + ":Height",
				Text: "297",
			},
		},
	}

	_, err := decodeDimension(elm)
	if err == nil {
		t.Errorf("expected error for invalid width value, got nil")
	}
}

func TestDimension_InvalidHeight(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Dimension",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "210",
			},
			{
				Name: NsWSCN + ":Height",
				Text: "invalid",
			},
		},
	}

	_, err := decodeDimension(elm)
	if err == nil {
		t.Errorf("expected error for invalid height value, got nil")
	}
}
