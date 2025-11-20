// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Exposure

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestExposure_RoundTrip(t *testing.T) {
	orig := Exposure{
		MustHonor: optional.New(BooleanElement("true")),
		Children: []xmldoc.Element{
			{Name: NsWSCN + ":SomeChild", Text: "value1"},
			{Name: NsWSCN + ":AnotherChild", Text: "value2"},
		},
	}

	elm := toXMLExposure(orig, NsWSCN+":Exposure")

	if elm.Name != NsWSCN+":Exposure" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":Exposure", elm.Name)
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

	// Check children
	if elm.Children[0].Name != NsWSCN+":SomeChild" {
		t.Errorf("expected first child name '%s', got '%s'",
			NsWSCN+":SomeChild", elm.Children[0].Name)
	}
	if elm.Children[1].Name != NsWSCN+":AnotherChild" {
		t.Errorf("expected second child name '%s', got '%s'",
			NsWSCN+":AnotherChild", elm.Children[1].Name)
	}

	// Decode back
	decoded, err := decodeExposure(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig.MustHonor, decoded.MustHonor) {
		t.Errorf("expected MustHonor %+v, got %+v", orig.MustHonor, decoded.MustHonor)
	}
	if len(decoded.Children) != len(orig.Children) {
		t.Errorf("expected %d children, got %d", len(orig.Children), len(decoded.Children))
	}
	for i := range orig.Children {
		if !orig.Children[i].Equal(decoded.Children[i]) {
			t.Errorf("child %d mismatch: expected %+v, got %+v",
				i, orig.Children[i], decoded.Children[i])
		}
	}
}

func TestExposure_NoAttributes(t *testing.T) {
	orig := Exposure{
		Children: []xmldoc.Element{
			{Name: NsWSCN + ":Child", Text: "value"},
		},
	}

	elm := toXMLExposure(orig, NsWSCN+":Exposure")

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}
	if len(elm.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(elm.Children))
	}

	decoded, err := decodeExposure(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.MustHonor != nil {
		t.Errorf("expected empty MustHonor, got %+v", decoded.MustHonor)
	}
	if len(decoded.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(decoded.Children))
	}
}

func TestExposure_NoChildren(t *testing.T) {
	orig := Exposure{
		MustHonor: optional.New(BooleanElement("false")),
		Children:  nil,
	}

	elm := toXMLExposure(orig, NsWSCN+":Exposure")

	if len(elm.Children) != 0 {
		t.Errorf("expected no children, got %d", len(elm.Children))
	}
	if len(elm.Attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(elm.Attrs))
	}

	decoded, err := decodeExposure(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "false" {
		t.Errorf("expected MustHonor='false', got '%s'", mustHonor)
	}
	if len(decoded.Children) != 0 {
		t.Errorf("expected no children, got %d", len(decoded.Children))
	}
}

func TestExposure_FromXML(t *testing.T) {
	// Create XML element manually with MustHonor and children
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "1"},
		},
		Children: []xmldoc.Element{
			{Name: NsWSCN + ":ExposureSettings", Text: "auto"},
			{Name: NsWSCN + ":Brightness", Text: "50"},
		},
	}

	decoded, err := decodeExposure(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "1" {
		t.Errorf("expected MustHonor='1', got '%s'", mustHonor)
	}
	if len(decoded.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(decoded.Children))
	}
	if decoded.Children[0].Name != NsWSCN+":ExposureSettings" {
		t.Errorf("expected first child name '%s', got '%s'",
			NsWSCN+":ExposureSettings", decoded.Children[0].Name)
	}
	if decoded.Children[1].Name != NsWSCN+":Brightness" {
		t.Errorf("expected second child name '%s', got '%s'",
			NsWSCN+":Brightness", decoded.Children[1].Name)
	}
}

func TestExposure_InvalidBooleanAttribute(t *testing.T) {
	// Test that invalid boolean value in MustHonor attribute is rejected
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
		Children: []xmldoc.Element{
			{Name: NsWSCN + ":Child", Text: "value"},
		},
	}

	_, err := decodeExposure(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}

func TestExposure_ValidBooleanValues(t *testing.T) {
	validValues := []string{"0", "1", "false", "true", "False", "True"}

	for _, val := range validValues {
		t.Run(val, func(t *testing.T) {
			root := xmldoc.Element{
				Name: NsWSCN + ":Exposure",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":MustHonor", Value: val},
				},
			}

			decoded, err := decodeExposure(root)
			if err != nil {
				t.Errorf("unexpected error for valid value '%s': %v", val, err)
			}
			if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) == "" {
				t.Errorf("expected MustHonor to be set for value '%s'", val)
			}
		})
	}
}
