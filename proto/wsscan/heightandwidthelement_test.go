// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for HeightAndWidthElement

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestHeightAndWidthElement_RoundTrip(t *testing.T) {
	orig := HeightAndWidthElement{
		Heights: []TextWithOverrideAndDefault{{
			Text:     "297",
			Override: optional.New(BooleanElement("true")),
		}},
		Widths: []TextWithOverrideAndDefault{{
			Text:        "210",
			UsedDefault: optional.New(BooleanElement("false")),
		}},
	}
	elm := orig.toXML("wscn:PlatenMaximumSize")
	if elm.Name != "wscn:PlatenMaximumSize" {
		t.Errorf("expected element name 'wscn:PlatenMaximumSize', got '%s'",
			elm.Name)
	}
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}

	decoded, err := decodeHeightAndWidthElement(elm)
	if err != nil {
		t.Fatalf("decodeHeightAndWidthElement returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestHeightAndWidthElement_MissingChild(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:PlatenMaximumSize",
		Children: []xmldoc.Element{
			{Name: "wscn:Height", Text: "297"},
		},
	}
	_, err := decodeHeightAndWidthElement(elm)
	if err == nil {
		t.Errorf("expected error for missing Width, got nil")
	}
}

func TestHeightAndWidthElement_InvalidBool(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:PlatenMaximumSize",
		Children: []xmldoc.Element{
			{
				Name: "wscn:Height",
				Text: "297",
				Attrs: []xmldoc.Attr{{
					Name:  "wscn:Override",
					Value: "maybe",
				}},
			},
			{
				Name: "wscn:Width",
				Text: "210",
			},
		},
	}
	_, err := decodeHeightAndWidthElement(elm)
	if err == nil {
		t.Errorf("expected error for invalid boolean attribute, got nil")
	}
}

func TestHeightAndWidthElement_MultipleHeightsWidths(t *testing.T) {
	orig := HeightAndWidthElement{
		Heights: []TextWithOverrideAndDefault{{Text: "297"}, {Text: "300"}},
		Widths:  []TextWithOverrideAndDefault{{Text: "210"}, {Text: "215"}},
	}
	elm := orig.toXML("wscn:PlatenMaximumSize")
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children (Heights, Widths), got %d", len(elm.Children))
	}
	if elm.Children[0].Name != NsWSCN+":Heights" || elm.Children[1].Name != NsWSCN+":Widths" {
		t.Errorf("expected wrapper elements <Heights> and <Widths>")
	}
	if len(elm.Children[0].Children) != 2 || len(elm.Children[1].Children) != 2 {
		t.Errorf("expected 2 Height and 2 Width children")
	}

	decoded, err := decodeHeightAndWidthElement(elm)
	if err != nil {
		t.Fatalf("decodeHeightAndWidthElement returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestHeightAndWidthElement_DecodeWrappedAndDirectMix(t *testing.T) {
	// XML with one direct Height, wrapped Widths
	elm := xmldoc.Element{
		Name: "wscn:PlatenMaximumSize",
		Children: []xmldoc.Element{
			{Name: NsWSCN + ":Height", Text: "297"},
			{
				Name: NsWSCN + ":Widths",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":Width", Text: "210"},
					{Name: NsWSCN + ":Width", Text: "215"},
				},
			},
		},
	}
	decoded, err := decodeHeightAndWidthElement(elm)
	if err != nil {
		t.Fatalf("decodeHeightAndWidthElement returned error: %v", err)
	}
	if len(decoded.Heights) != 1 || decoded.Heights[0].Text != "297" {
		t.Errorf("expected 1 Height with text '297'")
	}
	if len(decoded.Widths) != 2 || decoded.Widths[0].Text != "210" || decoded.Widths[1].Text != "215" {
		t.Errorf("expected 2 Widths with text '210' and '215'")
	}
}
