// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Scaling

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestScaling_RoundTrip_Complete(t *testing.T) {
	orig := Scaling{
		MustHonor: optional.New(BooleanElement("true")),
		ScalingWidth: AttributedElement[int]{
			Value:       500,
			Override:    optional.New(BooleanElement("false")),
			UsedDefault: optional.New(BooleanElement("true")),
		},
		ScalingHeight: AttributedElement[int]{
			Value:    600,
			Override: optional.New(BooleanElement("1")),
		},
	}

	elm := orig.toXML(NsWSCN + ":Scaling")

	if elm.Name != NsWSCN+":Scaling" {
		t.Errorf("expected element name '%s', got '%s'", NsWSCN+":Scaling", elm.Name)
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
	decoded, err := decodeScaling(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig.ScalingWidth, decoded.ScalingWidth) {
		t.Errorf("expected ScalingWidth %+v, got %+v", orig.ScalingWidth, decoded.ScalingWidth)
	}
	if !reflect.DeepEqual(orig.ScalingHeight, decoded.ScalingHeight) {
		t.Errorf("expected ScalingHeight %+v, got %+v", orig.ScalingHeight, decoded.ScalingHeight)
	}
	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "true" {
		t.Errorf("expected MustHonor='true', got '%s'", mustHonor)
	}
}

func TestScaling_RoundTrip_NoAttributes(t *testing.T) {
	orig := Scaling{
		ScalingWidth: AttributedElement[int]{
			Value: 100,
		},
		ScalingHeight: AttributedElement[int]{
			Value: 200,
		},
	}

	elm := orig.toXML(NsWSCN + ":Scaling")

	if elm.Name != NsWSCN+":Scaling" {
		t.Errorf("expected element name '%s', got '%s'", NsWSCN+":Scaling", elm.Name)
	}
	if len(elm.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(elm.Children))
	}
	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %d", len(elm.Attrs))
	}

	// Decode back
	decoded, err := decodeScaling(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.ScalingWidth.Value != orig.ScalingWidth.Value {
		t.Errorf("expected ScalingWidth.Value %d, got %d", orig.ScalingWidth.Value, decoded.ScalingWidth.Value)
	}
	if decoded.ScalingHeight.Value != orig.ScalingHeight.Value {
		t.Errorf("expected ScalingHeight.Value %d, got %d", orig.ScalingHeight.Value, decoded.ScalingHeight.Value)
	}
}

func TestScaling_FromXML_Complete(t *testing.T) {
	// Create XML element manually with all attributes
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "0"},
		},
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "300",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":Override", Value: "1"},
					{Name: NsWSCN + ":UsedDefault", Value: "false"},
				},
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "400",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":UsedDefault", Value: "true"},
				},
			},
		},
	}

	decoded, err := decodeScaling(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.ScalingWidth.Value != 300 {
		t.Errorf("expected ScalingWidth.Value 300, got %d", decoded.ScalingWidth.Value)
	}
	if decoded.ScalingHeight.Value != 400 {
		t.Errorf("expected ScalingHeight.Value 400, got %d", decoded.ScalingHeight.Value)
	}
	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "0" {
		t.Errorf("expected MustHonor='0', got '%s'", mustHonor)
	}
	if override := optional.Get(decoded.ScalingWidth.Override); string(override) != "1" {
		t.Errorf("expected ScalingWidth.Override='1', got '%s'", override)
	}
	if usedDefault := optional.Get(decoded.ScalingWidth.UsedDefault); string(usedDefault) != "false" {
		t.Errorf("expected ScalingWidth.UsedDefault='false', got '%s'", usedDefault)
	}
}

func TestScaling_MissingScalingWidth(t *testing.T) {
	// Create XML element without required ScalingWidth
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "500",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for missing required ScalingWidth element, got nil")
	}
}

func TestScaling_MissingScalingHeight(t *testing.T) {
	// Create XML element without required ScalingHeight
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "500",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for missing required ScalingHeight element, got nil")
	}
}

func TestScaling_InvalidScalingWidthValue_TooLow(t *testing.T) {
	// Create XML element with ScalingWidth value below minimum (0)
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "0",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "500",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for ScalingWidth value 0 (below minimum 1), got nil")
	}
}

func TestScaling_InvalidScalingWidthValue_TooHigh(t *testing.T) {
	// Create XML element with ScalingWidth value above maximum (1001)
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "1001",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "500",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for ScalingWidth value 1001 (above maximum 1000), got nil")
	}
}

func TestScaling_InvalidScalingHeightValue_TooLow(t *testing.T) {
	// Create XML element with ScalingHeight value below minimum (0)
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "500",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "0",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for ScalingHeight value 0 (below minimum 1), got nil")
	}
}

func TestScaling_InvalidScalingHeightValue_TooHigh(t *testing.T) {
	// Create XML element with ScalingHeight value above maximum (1001)
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "500",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "1001",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for ScalingHeight value 1001 (above maximum 1000), got nil")
	}
}

func TestScaling_InvalidScalingWidthValue_NotInteger(t *testing.T) {
	// Create XML element with invalid ScalingWidth value
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "invalid",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "500",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for invalid ScalingWidth value 'invalid', got nil")
	}
}

func TestScaling_InvalidScalingHeightValue_NotInteger(t *testing.T) {
	// Create XML element with invalid ScalingHeight value
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "500",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "invalid",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for invalid ScalingHeight value 'invalid', got nil")
	}
}

func TestScaling_InvalidMustHonor(t *testing.T) {
	// Create XML element with invalid MustHonor attribute
	root := xmldoc.Element{
		Name: NsWSCN + ":Scaling",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ScalingWidth",
				Text: "500",
			},
			{
				Name: NsWSCN + ":ScalingHeight",
				Text: "600",
			},
		},
	}

	_, err := decodeScaling(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}

func TestScaling_BoundaryValues(t *testing.T) {
	// Test minimum value (1)
	orig1 := Scaling{
		ScalingWidth:  AttributedElement[int]{Value: 1},
		ScalingHeight: AttributedElement[int]{Value: 1},
	}
	elm1 := orig1.toXML(NsWSCN + ":Scaling")
	decoded1, err := decodeScaling(elm1)
	if err != nil {
		t.Fatalf("decode with value 1 returned error: %v", err)
	}
	if decoded1.ScalingWidth.Value != 1 || decoded1.ScalingHeight.Value != 1 {
		t.Errorf("expected values 1, got %d, %d", decoded1.ScalingWidth.Value, decoded1.ScalingHeight.Value)
	}

	// Test maximum value (1000)
	orig2 := Scaling{
		ScalingWidth:  AttributedElement[int]{Value: 1000},
		ScalingHeight: AttributedElement[int]{Value: 1000},
	}
	elm2 := orig2.toXML(NsWSCN + ":Scaling")
	decoded2, err := decodeScaling(elm2)
	if err != nil {
		t.Fatalf("decode with value 1000 returned error: %v", err)
	}
	if decoded2.ScalingWidth.Value != 1000 || decoded2.ScalingHeight.Value != 1000 {
		t.Errorf("expected values 1000, got %d, %d", decoded2.ScalingWidth.Value, decoded2.ScalingHeight.Value)
	}
}

func TestScaling_AttributesOnChildElements(t *testing.T) {
	// Test ScalingWidth and ScalingHeight with all attributes
	orig := Scaling{
		ScalingWidth: AttributedElement[int]{
			Value:       300,
			MustHonor:   optional.New(BooleanElement("true")),
			Override:    optional.New(BooleanElement("false")),
			UsedDefault: optional.New(BooleanElement("1")),
		},
		ScalingHeight: AttributedElement[int]{
			Value:       400,
			MustHonor:   optional.New(BooleanElement("0")),
			Override:    optional.New(BooleanElement("true")),
			UsedDefault: optional.New(BooleanElement("false")),
		},
	}

	elm := orig.toXML(NsWSCN + ":Scaling")
	if len(elm.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(elm.Children))
	}

	decoded, err := decodeScaling(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(orig.ScalingWidth, decoded.ScalingWidth) {
		t.Errorf("expected ScalingWidth %+v, got %+v", orig.ScalingWidth, decoded.ScalingWidth)
	}
	if !reflect.DeepEqual(orig.ScalingHeight, decoded.ScalingHeight) {
		t.Errorf("expected ScalingHeight %+v, got %+v", orig.ScalingHeight, decoded.ScalingHeight)
	}
}
