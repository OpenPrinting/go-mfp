// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for InputSize

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestInputSize_RoundTrip(t *testing.T) {
	orig := InputSize{
		MustHonor:              optional.New(BooleanElement("true")),
		DocumentSizeAutoDetect: optional.New(BooleanElement("1")),
		InputMediaSize: InputMediaSize{
			Width: AttributedElement[int]{
				Value: 8500,
			},
			Height: AttributedElement[int]{
				Value: 11000,
			},
		},
	}

	elm := orig.toXML(NsWSCN + ":InputSize")

	if elm.Name != NsWSCN+":InputSize" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":InputSize", elm.Name)
	}
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}
	if len(elm.Attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(elm.Attrs))
	}

	// Check MustHonor attribute
	if elm.Attrs[0].Name != NsWSCN+":MustHonor" {
		t.Errorf("expected attribute name '%s', got '%s'",
			NsWSCN+":MustHonor", elm.Attrs[0].Name)
	}
	if elm.Attrs[0].Value != "true" {
		t.Errorf("expected MustHonor='true', got '%s'", elm.Attrs[0].Value)
	}

	// Decode back
	decoded, err := decodeInputSize(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig.MustHonor, decoded.MustHonor) {
		t.Errorf("expected MustHonor %+v, got %+v", orig.MustHonor, decoded.MustHonor)
	}
	if !reflect.DeepEqual(orig.DocumentSizeAutoDetect, decoded.DocumentSizeAutoDetect) {
		t.Errorf("expected DocumentSizeAutoDetect %+v, got %+v",
			orig.DocumentSizeAutoDetect, decoded.DocumentSizeAutoDetect)
	}
	if !reflect.DeepEqual(orig.InputMediaSize, decoded.InputMediaSize) {
		t.Errorf("expected InputMediaSize %+v, got %+v",
			orig.InputMediaSize, decoded.InputMediaSize)
	}
	if decoded.InputMediaSize.Width.Value != 8500 {
		t.Errorf("expected InputMediaSize.Width=8500, got %d", decoded.InputMediaSize.Width.Value)
	}
	if decoded.InputMediaSize.Height.Value != 11000 {
		t.Errorf("expected InputMediaSize.Height=11000, got %d", decoded.InputMediaSize.Height.Value)
	}
}

func TestInputSize_DocumentSizeAutoDetectOnly(t *testing.T) {
	orig := InputSize{
		DocumentSizeAutoDetect: optional.New(BooleanElement("true")),
	}

	elm := orig.toXML(NsWSCN + ":InputSize")

	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}
	if elm.Children[0].Name != NsWSCN+":DocumentSizeAutoDetect" {
		t.Errorf("expected child name '%s', got '%s'",
			NsWSCN+":DocumentSizeAutoDetect", elm.Children[0].Name)
	}
	if elm.Children[0].Text != "true" {
		t.Errorf("expected DocumentSizeAutoDetect text 'true', got '%s'", elm.Children[0].Text)
	}

	decoded, err := decodeInputSize(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.DocumentSizeAutoDetect == nil {
		t.Errorf("expected DocumentSizeAutoDetect to be set")
	}
	if optional.Get(decoded.DocumentSizeAutoDetect) != BooleanElement("true") {
		t.Errorf("expected DocumentSizeAutoDetect='true', got '%v'",
			optional.Get(decoded.DocumentSizeAutoDetect))
	}
}

func TestInputSize_InputMediaSizeOnly(t *testing.T) {
	orig := InputSize{
		InputMediaSize: InputMediaSize{
			Width: AttributedElement[int]{
				Value: 8500,
			},
			Height: AttributedElement[int]{
				Value: 11000,
			},
		},
	}

	elm := orig.toXML(NsWSCN + ":InputSize")

	if len(elm.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(elm.Children))
	}
	if elm.Children[0].Name != NsWSCN+":InputMediaSize" {
		t.Errorf("expected child name '%s', got '%s'",
			NsWSCN+":InputMediaSize", elm.Children[0].Name)
	}

	decoded, err := decodeInputSize(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.InputMediaSize.Width.Value != 8500 {
		t.Errorf("expected Width=8500, got %d", decoded.InputMediaSize.Width.Value)
	}
	if decoded.InputMediaSize.Height.Value != 11000 {
		t.Errorf("expected Height=11000, got %d", decoded.InputMediaSize.Height.Value)
	}
}

func TestInputSize_FromXML(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputSize",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "1"},
		},
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":DocumentSizeAutoDetect",
				Text: "true",
			},
			{
				Name: NsWSCN + ":InputMediaSize",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":Width",
						Text: "8500",
					},
					{
						Name: NsWSCN + ":Height",
						Text: "11000",
					},
				},
			},
		},
	}

	decoded, err := decodeInputSize(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "1" {
		t.Errorf("expected MustHonor='1', got '%s'", mustHonor)
	}
	if autoDetect := optional.Get(decoded.DocumentSizeAutoDetect); autoDetect != BooleanElement("true") {
		t.Errorf("expected DocumentSizeAutoDetect='true', got '%v'", autoDetect)
	}
	if decoded.InputMediaSize.Width.Value != 8500 {
		t.Errorf("expected Width=8500, got %d", decoded.InputMediaSize.Width.Value)
	}
	if decoded.InputMediaSize.Height.Value != 11000 {
		t.Errorf("expected Height=11000, got %d", decoded.InputMediaSize.Height.Value)
	}
}

func TestInputSize_InvalidBooleanAttribute(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputSize",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
	}

	_, err := decodeInputSize(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}

func TestInputSize_InvalidDocumentSizeAutoDetect(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":DocumentSizeAutoDetect",
				Text: "invalid",
			},
			{
				Name: NsWSCN + ":InputMediaSize",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":Width",
						Text: "8500",
					},
					{
						Name: NsWSCN + ":Height",
						Text: "11000",
					},
				},
			},
		},
	}

	_, err := decodeInputSize(root)
	if err == nil {
		t.Errorf("expected error for invalid DocumentSizeAutoDetect value 'invalid', got nil")
	}
}

func TestInputSize_MissingInputMediaSize(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":InputSize",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":DocumentSizeAutoDetect",
				Text: "true",
			},
		},
	}

	_, err := decodeInputSize(root)
	if err == nil {
		t.Errorf("expected error for missing InputMediaSize, got nil")
	}
}
