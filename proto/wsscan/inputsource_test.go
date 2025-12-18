// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for InputSource

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestInputSource_RoundTrip(t *testing.T) {
	orig := InputSource{
		Value:       InputSourceADF,
		MustHonor:   optional.New(BooleanElement("true")),
		Override:    optional.New(BooleanElement("false")),
		UsedDefault: optional.New(BooleanElement("1")),
	}

	elm := toXMLInputSource(orig, NsWSCN+":InputSource")

	if elm.Name != NsWSCN+":InputSource" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":InputSource", elm.Name)
	}
	if elm.Text != "ADF" {
		t.Errorf("expected text 'ADF', got '%s'", elm.Text)
	}
	if len(elm.Attrs) != 3 {
		t.Errorf("expected 3 attributes, got %d", len(elm.Attrs))
	}

	// Check attributes
	attrsMap := make(map[string]string)
	for _, attr := range elm.Attrs {
		attrsMap[attr.Name] = attr.Value
	}
	if attrsMap[NsWSCN+":MustHonor"] != "true" {
		t.Errorf("expected MustHonor='true', got '%s'", attrsMap[NsWSCN+":MustHonor"])
	}
	if attrsMap[NsWSCN+":Override"] != "false" {
		t.Errorf("expected Override='false', got '%s'", attrsMap[NsWSCN+":Override"])
	}
	if attrsMap[NsWSCN+":UsedDefault"] != "1" {
		t.Errorf("expected UsedDefault='1', got '%s'", attrsMap[NsWSCN+":UsedDefault"])
	}

	// Decode back
	decoded, err := decodeInputSource(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
	if !reflect.DeepEqual(orig.MustHonor, decoded.MustHonor) {
		t.Errorf("expected MustHonor %+v, got %+v", orig.MustHonor, decoded.MustHonor)
	}
	if !reflect.DeepEqual(orig.Override, decoded.Override) {
		t.Errorf("expected Override %+v, got %+v", orig.Override, decoded.Override)
	}
	if !reflect.DeepEqual(orig.UsedDefault, decoded.UsedDefault) {
		t.Errorf("expected UsedDefault %+v, got %+v", orig.UsedDefault, decoded.UsedDefault)
	}
}

func TestInputSource_NoAttributes(t *testing.T) {
	orig := InputSource{
		Value: InputSourcePlaten,
	}

	elm := toXMLInputSource(orig, NsWSCN+":InputSource")

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}
	if elm.Text != "Platen" {
		t.Errorf("expected text 'Platen', got '%s'", elm.Text)
	}

	decoded, err := decodeInputSource(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestInputSource_StandardValues(t *testing.T) {
	standardValues := []struct {
		enumValue InputSourceValue
		textValue string
	}{
		{InputSourceADF, "ADF"},
		{InputSourceADFDuplex, "ADFDuplex"},
		{InputSourceFilm, "Film"},
		{InputSourcePlaten, "Platen"},
	}

	for _, tc := range standardValues {
		t.Run(tc.textValue, func(t *testing.T) {
			orig := InputSource{
				Value: tc.enumValue,
			}

			elm := toXMLInputSource(orig, NsWSCN+":InputSource")
			if elm.Text != tc.textValue {
				t.Errorf("expected text '%s', got '%s'", tc.textValue, elm.Text)
			}

			decoded, err := decodeInputSource(elm)
			if err != nil {
				t.Fatalf("decode returned error: %v", err)
			}
			if decoded.Value != tc.enumValue {
				t.Errorf("expected value %v, got %v", tc.enumValue, decoded.Value)
			}
		})
	}
}

func TestInputSource_FromXML(t *testing.T) {
	// Create XML element manually with all attributes
	root := xmldoc.Element{
		Name: NsWSCN + ":InputSource",
		Text: "ADFDuplex",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "0"},
			{Name: NsWSCN + ":Override", Value: "1"},
			{Name: NsWSCN + ":UsedDefault", Value: "false"},
		},
	}

	decoded, err := decodeInputSource(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Value != InputSourceADFDuplex {
		t.Errorf("expected value InputSourceADFDuplex, got %v", decoded.Value)
	}
	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "0" {
		t.Errorf("expected MustHonor='0', got '%s'", mustHonor)
	}
	if override := optional.Get(decoded.Override); string(override) != "1" {
		t.Errorf("expected Override='1', got '%s'", override)
	}
	if usedDefault := optional.Get(decoded.UsedDefault); string(usedDefault) != "false" {
		t.Errorf("expected UsedDefault='false', got '%s'", usedDefault)
	}
}

func TestInputSource_InvalidBooleanAttributes(t *testing.T) {
	// Test that invalid boolean values in attributes are rejected
	root := xmldoc.Element{
		Name: NsWSCN + ":InputSource",
		Text: "ADF",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
	}

	_, err := decodeInputSource(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}

func TestInputSource_AllStandardValuesWithAttributes(t *testing.T) {
	standardValues := []struct {
		enumValue InputSourceValue
		textValue string
	}{
		{InputSourceADF, "ADF"},
		{InputSourceADFDuplex, "ADFDuplex"},
		{InputSourceFilm, "Film"},
		{InputSourcePlaten, "Platen"},
	}

	for _, tc := range standardValues {
		t.Run(tc.textValue, func(t *testing.T) {
			orig := InputSource{
				Value:       tc.enumValue,
				MustHonor:   optional.New(BooleanElement("1")),
				Override:    optional.New(BooleanElement("0")),
				UsedDefault: optional.New(BooleanElement("true")),
			}

			elm := toXMLInputSource(orig, NsWSCN+":InputSource")
			decoded, err := decodeInputSource(elm)
			if err != nil {
				t.Fatalf("decode returned error for value '%s': %v", tc.textValue, err)
			}
			if decoded.Value != tc.enumValue {
				t.Errorf("expected value %v, got %v", tc.enumValue, decoded.Value)
			}
			if len(elm.Attrs) != 3 {
				t.Errorf("expected 3 attributes for value '%s', got %d", tc.textValue, len(elm.Attrs))
			}
		})
	}
}

func TestInputSource_VendorDefinedValues(t *testing.T) {
	// Test that vendor-defined values decode to UnknownInputSourceValue
	vendorValues := []string{"vendor-source-1", "custom-source", "extended-value"}

	for _, val := range vendorValues {
		t.Run(val, func(t *testing.T) {
			root := xmldoc.Element{
				Name: NsWSCN + ":InputSource",
				Text: val,
			}

			decoded, err := decodeInputSource(root)
			if err != nil {
				t.Fatalf("decode returned error for vendor-defined value '%s': %v", val, err)
			}
			if decoded.Value != UnknownInputSourceValue {
				t.Errorf("expected UnknownInputSourceValue, got %v", decoded.Value)
			}
			// When encoding UnknownInputSourceValue, it will return "Unknown"
			elm := toXMLInputSource(decoded, NsWSCN+":InputSource")
			if elm.Text != "Unknown" {
				t.Errorf("expected text 'Unknown' for UnknownInputSourceValue, got '%s'", elm.Text)
			}
		})
	}
}
