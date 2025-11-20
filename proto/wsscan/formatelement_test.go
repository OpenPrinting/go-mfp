// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for FormatElement

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestFormatElement_RoundTrip(t *testing.T) {
	orig := FormatElement{
		Value:       "png",
		Override:    optional.New(BooleanElement("false")),
		UsedDefault: optional.New(BooleanElement("1")),
	}

	elm := toXMLFormatElement(orig, NsWSCN+":Format")

	if elm.Name != NsWSCN+":Format" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":Format", elm.Name)
	}
	if elm.Text != "png" {
		t.Errorf("expected text 'png', got '%s'", elm.Text)
	}
	if len(elm.Attrs) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(elm.Attrs))
	}

	// Check attributes
	attrsMap := make(map[string]string)
	for _, attr := range elm.Attrs {
		attrsMap[attr.Name] = attr.Value
	}
	if attrsMap[NsWSCN+":Override"] != "false" {
		t.Errorf("expected Override='false', got '%s'", attrsMap[NsWSCN+":Override"])
	}
	if attrsMap[NsWSCN+":UsedDefault"] != "1" {
		t.Errorf("expected UsedDefault='1', got '%s'", attrsMap[NsWSCN+":UsedDefault"])
	}
	// MustHonor should not be present
	if _, found := attrsMap[NsWSCN+":MustHonor"]; found {
		t.Errorf("expected MustHonor to not be present, but it was found")
	}

	// Decode back
	decoded, err := decodeFormatElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
	if !reflect.DeepEqual(orig.Override, decoded.Override) {
		t.Errorf("expected Override %+v, got %+v", orig.Override, decoded.Override)
	}
	if !reflect.DeepEqual(orig.UsedDefault, decoded.UsedDefault) {
		t.Errorf("expected UsedDefault %+v, got %+v", orig.UsedDefault, decoded.UsedDefault)
	}
}

func TestFormatElement_NoAttributes(t *testing.T) {
	orig := FormatElement{
		Value: "jpeg2k",
	}

	elm := toXMLFormatElement(orig, NsWSCN+":Format")

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}
	if elm.Text != "jpeg2k" {
		t.Errorf("expected text 'jpeg2k', got '%s'", elm.Text)
	}

	decoded, err := decodeFormatElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestFormatElement_StandardValues(t *testing.T) {
	standardValues := []string{
		"dib",
		"exif",
		"jbig",
		"jfif",
		"jpeg2k",
		"pdf-a",
		"png",
		"tiff-single-uncompressed",
		"tiff-single-g4",
		"tiff-single-g3mh",
		"tiff-single-jpeg-tn2",
		"tiff-multi-uncompressed",
		"tiff-multi-g4",
		"tiff-multi-g3mh",
		"tiff-multi-jpeg-tn2",
		"xps",
	}

	for _, val := range standardValues {
		t.Run(val, func(t *testing.T) {
			orig := FormatElement{
				Value: val,
			}

			elm := toXMLFormatElement(orig, NsWSCN+":Format")
			if elm.Text != val {
				t.Errorf("expected text '%s', got '%s'", val, elm.Text)
			}

			decoded, err := decodeFormatElement(elm)
			if err != nil {
				t.Fatalf("decode returned error: %v", err)
			}
			if decoded.Value != val {
				t.Errorf("expected value %s, got %s", val, decoded.Value)
			}
		})
	}
}

func TestFormatElement_VendorDefinedValues(t *testing.T) {
	// Test that vendor-defined values are accepted
	vendorValues := []string{"vendor-format-1", "custom-format", "extended-value"}

	for _, val := range vendorValues {
		t.Run(val, func(t *testing.T) {
			orig := FormatElement{
				Value: val,
			}

			elm := toXMLFormatElement(orig, NsWSCN+":Format")
			decoded, err := decodeFormatElement(elm)
			if err != nil {
				t.Fatalf("decode returned error for vendor-defined value '%s': %v", val, err)
			}
			if decoded.Value != val {
				t.Errorf("expected value %s, got %s", val, decoded.Value)
			}
		})
	}
}

func TestFormatElement_FromXML(t *testing.T) {
	// Create XML element manually with Override and UsedDefault (no MustHonor)
	root := xmldoc.Element{
		Name: NsWSCN + ":Format",
		Text: "tiff-multi-g4",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":Override", Value: "0"},
			{Name: NsWSCN + ":UsedDefault", Value: "true"},
		},
	}

	decoded, err := decodeFormatElement(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Value != "tiff-multi-g4" {
		t.Errorf("expected value 'tiff-multi-g4', got '%s'", decoded.Value)
	}
	if override := optional.Get(decoded.Override); string(override) != "0" {
		t.Errorf("expected Override='0', got '%s'", override)
	}
	if usedDefault := optional.Get(decoded.UsedDefault); string(usedDefault) != "true" {
		t.Errorf("expected UsedDefault='true', got '%s'", usedDefault)
	}
	// MustHonor should not be set
	if decoded.MustHonor != nil {
		t.Errorf("expected MustHonor to be nil, got %+v", decoded.MustHonor)
	}
}

func TestFormatElement_InvalidBooleanAttributes(t *testing.T) {
	// Test that invalid boolean values in attributes are rejected
	root := xmldoc.Element{
		Name: NsWSCN + ":Format",
		Text: "png",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":Override", Value: "invalid"},
		},
	}

	_, err := decodeFormatElement(root)
	if err == nil {
		t.Errorf("expected error for invalid Override value 'invalid', got nil")
	}
}

func TestFormatElement_OnlyOverride(t *testing.T) {
	// Test with only Override attribute
	orig := FormatElement{
		Value:    "pdf-a",
		Override: optional.New(BooleanElement("1")),
	}

	elm := toXMLFormatElement(orig, NsWSCN+":Format")

	if len(elm.Attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(elm.Attrs))
	}
	if elm.Attrs[0].Name != NsWSCN+":Override" {
		t.Errorf("expected Override attribute, got %s", elm.Attrs[0].Name)
	}

	decoded, err := decodeFormatElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != "pdf-a" {
		t.Errorf("expected value 'pdf-a', got '%s'", decoded.Value)
	}
}

func TestFormatElement_OnlyUsedDefault(t *testing.T) {
	// Test with only UsedDefault attribute
	orig := FormatElement{
		Value:       "exif",
		UsedDefault: optional.New(BooleanElement("false")),
	}

	elm := toXMLFormatElement(orig, NsWSCN+":Format")

	if len(elm.Attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(elm.Attrs))
	}
	if elm.Attrs[0].Name != NsWSCN+":UsedDefault" {
		t.Errorf("expected UsedDefault attribute, got %s", elm.Attrs[0].Name)
	}

	decoded, err := decodeFormatElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != "exif" {
		t.Errorf("expected value 'exif', got '%s'", decoded.Value)
	}
}
