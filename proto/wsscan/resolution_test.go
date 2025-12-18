// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Resolution

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestResolution_RoundTrip_WithHeight(t *testing.T) {
	orig := Resolution{
		MustHonor: optional.New(BooleanElement("true")),
		Width: AttributedElement[int]{
			Value:       300,
			Override:    optional.New(BooleanElement("false")),
			UsedDefault: optional.New(BooleanElement("true")),
		},
		Height: optional.New(AttributedElement[int]{
			Value:    600,
			Override: optional.New(BooleanElement("1")),
		}),
	}

	elm := orig.toXML(NsWSCN + ":Resolution")

	if elm.Name != NsWSCN+":Resolution" {
		t.Errorf("expected element name '%s', got '%s'", NsWSCN+":Resolution", elm.Name)
	}
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}
	if len(elm.Attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(elm.Attrs))
	}

	// Check MustHonor attribute
	if elm.Attrs[0].Name != NsWSCN+":MustHonor" || elm.Attrs[0].Value != "true" {
		t.Errorf("expected MustHonor='true', got %+v", elm.Attrs[0])
	}

	// Decode back
	decoded, err := decodeResolution(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Width.Value != orig.Width.Value {
		t.Errorf("expected Width.Value %d, got %d", orig.Width.Value, decoded.Width.Value)
	}
	if decoded.Height == nil {
		t.Errorf("expected Height to be present, got nil")
	} else {
		decodedHeight := optional.Get(decoded.Height)
		origHeight := optional.Get(orig.Height)
		if decodedHeight.Value != origHeight.Value {
			t.Errorf("expected Height.Value %d, got %d", origHeight.Value, decodedHeight.Value)
		}
	}
}

func TestResolution_RoundTrip_WithoutHeight(t *testing.T) {
	orig := Resolution{
		Width: AttributedElement[int]{
			Value: 300,
		},
		// Height is nil (optional)
	}

	elm := orig.toXML(NsWSCN + ":Resolution")

	if elm.Name != NsWSCN+":Resolution" {
		t.Errorf("expected element name '%s', got '%s'", NsWSCN+":Resolution", elm.Name)
	}
	if len(elm.Children) != 1 {
		t.Errorf("expected 1 child (only Width), got %d", len(elm.Children))
	}
	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %d", len(elm.Attrs))
	}

	// Check that only Width is present
	if elm.Children[0].Name != NsWSCN+":Width" {
		t.Errorf("expected first child to be Width, got %s", elm.Children[0].Name)
	}
	if elm.Children[0].Text != "300" {
		t.Errorf("expected Width text '300', got '%s'", elm.Children[0].Text)
	}

	// Decode back
	decoded, err := decodeResolution(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Width.Value != orig.Width.Value {
		t.Errorf("expected Width.Value %d, got %d", orig.Width.Value, decoded.Width.Value)
	}
	if decoded.Height != nil {
		t.Errorf("expected Height to be nil, got %+v", decoded.Height)
	}
}

func TestResolution_FromXML_WithHeight(t *testing.T) {
	// Create XML element manually with Width and Height
	root := xmldoc.Element{
		Name: NsWSCN + ":Resolution",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "0"},
		},
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "300",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":Override", Value: "1"},
				},
			},
			{
				Name: NsWSCN + ":Height",
				Text: "600",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":UsedDefault", Value: "true"},
				},
			},
		},
	}

	decoded, err := decodeResolution(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Width.Value != 300 {
		t.Errorf("expected Width.Value 300, got %d", decoded.Width.Value)
	}
	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "0" {
		t.Errorf("expected MustHonor='0', got '%s'", mustHonor)
	}
	if override := optional.Get(decoded.Width.Override); string(override) != "1" {
		t.Errorf("expected Width.Override='1', got '%s'", override)
	}
	if decoded.Height == nil {
		t.Fatalf("expected Height to be present, got nil")
	}
	height := optional.Get(decoded.Height)
	if height.Value != 600 {
		t.Errorf("expected Height.Value 600, got %d", height.Value)
	}
	if usedDefault := optional.Get(height.UsedDefault); string(usedDefault) != "true" {
		t.Errorf("expected Height.UsedDefault='true', got '%s'", usedDefault)
	}
}

func TestResolution_FromXML_WithoutHeight(t *testing.T) {
	// Create XML element manually with only Width (no Height)
	root := xmldoc.Element{
		Name: NsWSCN + ":Resolution",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "300",
			},
		},
	}

	decoded, err := decodeResolution(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Width.Value != 300 {
		t.Errorf("expected Width.Value 300, got %d", decoded.Width.Value)
	}
	if decoded.Height != nil {
		t.Errorf("expected Height to be nil when not present in XML, got %+v", decoded.Height)
	}
}

func TestResolution_MissingWidth(t *testing.T) {
	// Create XML element without required Width
	root := xmldoc.Element{
		Name: NsWSCN + ":Resolution",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Height",
				Text: "600",
			},
		},
	}

	_, err := decodeResolution(root)
	if err == nil {
		t.Errorf("expected error for missing required Width element, got nil")
	}
}

func TestResolution_InvalidWidthValue(t *testing.T) {
	// Create XML element with invalid Width value
	root := xmldoc.Element{
		Name: NsWSCN + ":Resolution",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "invalid",
			},
		},
	}

	_, err := decodeResolution(root)
	if err == nil {
		t.Errorf("expected error for invalid Width value 'invalid', got nil")
	}
}

func TestResolution_InvalidHeightValue(t *testing.T) {
	// Create XML element with invalid Height value
	root := xmldoc.Element{
		Name: NsWSCN + ":Resolution",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "300",
			},
			{
				Name: NsWSCN + ":Height",
				Text: "invalid",
			},
		},
	}

	_, err := decodeResolution(root)
	if err == nil {
		t.Errorf("expected error for invalid Height value 'invalid', got nil")
	}
}

func TestResolution_InvalidMustHonor(t *testing.T) {
	// Create XML element with invalid MustHonor attribute
	root := xmldoc.Element{
		Name: NsWSCN + ":Resolution",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: "300",
			},
		},
	}

	_, err := decodeResolution(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}

func TestResolution_WidthAttributes(t *testing.T) {
	// Test Width with all attributes
	orig := Resolution{
		Width: AttributedElement[int]{
			Value:       300,
			MustHonor:   optional.New(BooleanElement("true")),
			Override:    optional.New(BooleanElement("false")),
			UsedDefault: optional.New(BooleanElement("1")),
		},
	}

	elm := orig.toXML(NsWSCN + ":Resolution")
	if len(elm.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(elm.Children))
	}

	widthElm := elm.Children[0]
	if widthElm.Name != NsWSCN+":Width" {
		t.Fatalf("expected Width element, got %s", widthElm.Name)
	}
	if len(widthElm.Attrs) != 3 {
		t.Errorf("expected 3 attributes on Width, got %d", len(widthElm.Attrs))
	}

	decoded, err := decodeResolution(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig.Width, decoded.Width) {
		t.Errorf("expected Width %+v, got %+v", orig.Width, decoded.Width)
	}
}
