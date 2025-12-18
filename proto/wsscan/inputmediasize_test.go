// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for InputMediaSize

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestInputMediaSize_RoundTrip(t *testing.T) {
	orig := InputMediaSize{
		Width: AttributedElement[int]{
			Value:    8500,
			Override: optional.New(BooleanElement("1")),
		},
		Height: AttributedElement[int]{
			Value:       11000,
			UsedDefault: optional.New(BooleanElement("true")),
		},
	}

	elm := orig.toXML(NsWSCN + ":InputMediaSize")

	if elm.Name != NsWSCN+":InputMediaSize" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":InputMediaSize", elm.Name)
	}
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}

	decoded, err := decodeInputMediaSize(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig.Width, decoded.Width) {
		t.Errorf("Width mismatch: expected %+v, got %+v", orig.Width, decoded.Width)
	}
	if !reflect.DeepEqual(orig.Height, decoded.Height) {
		t.Errorf("Height mismatch: expected %+v, got %+v", orig.Height, decoded.Height)
	}
}

func TestInputMediaSize_NoAttributes(t *testing.T) {
	orig := InputMediaSize{
		Width: AttributedElement[int]{
			Value: 8500,
		},
		Height: AttributedElement[int]{
			Value: 11000,
		},
	}

	elm := orig.toXML(NsWSCN + ":InputMediaSize")

	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}
	if elm.Children[0].Text != "8500" {
		t.Errorf("expected Width text '8500', got '%s'", elm.Children[0].Text)
	}
	if elm.Children[1].Text != "11000" {
		t.Errorf("expected Height text '11000', got '%s'", elm.Children[1].Text)
	}

	decoded, err := decodeInputMediaSize(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Width.Value != 8500 {
		t.Errorf("expected Width=8500, got %d", decoded.Width.Value)
	}
	if decoded.Height.Value != 11000 {
		t.Errorf("expected Height=11000, got %d", decoded.Height.Value)
	}
}

func TestInputMediaSize_FromXML(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputMediaSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "8500",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":Override", Value: "0"},
				},
			},
			{
				Name: NsWSCN + ":Height",
				Text: "11000",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":UsedDefault", Value: "false"},
				},
			},
		},
	}

	decoded, err := decodeInputMediaSize(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Width.Value != 8500 {
		t.Errorf("expected Width=8500, got %d", decoded.Width.Value)
	}
	if override := optional.Get(decoded.Width.Override); override != "0" {
		t.Errorf("expected Width Override='0', got '%s'", override)
	}
	if decoded.Height.Value != 11000 {
		t.Errorf("expected Height=11000, got %d", decoded.Height.Value)
	}
	if usedDefault := optional.Get(decoded.Height.UsedDefault); usedDefault != "false" {
		t.Errorf("expected Height UsedDefault='false', got '%s'", usedDefault)
	}
}

func TestInputMediaSize_MissingWidth(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputMediaSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Height",
				Text: "11000",
			},
		},
	}

	_, err := decodeInputMediaSize(root)
	if err == nil {
		t.Errorf("expected error for missing Width, got nil")
	}
}

func TestInputMediaSize_MissingHeight(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputMediaSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "8500",
			},
		},
	}

	_, err := decodeInputMediaSize(root)
	if err == nil {
		t.Errorf("expected error for missing Height, got nil")
	}
}

func TestInputMediaSize_InvalidIntegerValue(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputMediaSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "invalid",
			},
			{
				Name: NsWSCN + ":Height",
				Text: "11000",
			},
		},
	}

	_, err := decodeInputMediaSize(root)
	if err == nil {
		t.Errorf("expected error for invalid Width value, got nil")
	}
}

func TestInputMediaSize_InvalidBooleanAttributes(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputMediaSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "8500",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":Override", Value: "invalid"},
				},
			},
			{
				Name: NsWSCN + ":Height",
				Text: "11000",
			},
		},
	}

	_, err := decodeInputMediaSize(root)
	if err == nil {
		t.Errorf("expected error for invalid Override value, got nil")
	}
}
