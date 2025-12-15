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
		Value:       PNG,
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
		Value: JPEG2K,
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
	standardValues := []struct {
		formatValue FormatValue
		textValue   string
	}{
		{DIB, "dib"},
		{EXIF, "exif"},
		{JBIG, "jbig"},
		{JFIF, "jfif"},
		{JPEG2K, "jpeg2k"},
		{PDFA, "pdf-a"},
		{PNG, "png"},
		{TIFFSingleUncompressed, "tiff-single-uncompressed"},
		{TIFFSingleG4, "tiff-single-g4"},
		{TIFFSingleG3MH, "tiff-single-g3mh"},
		{TIFFSingleJPEGTN2, "tiff-single-jpeg-tn2"},
		{TIFFMultiUncompressed, "tiff-multi-uncompressed"},
		{TIFFMultiG4, "tiff-multi-g4"},
		{TIFFMultiG3MH, "tiff-multi-g3mh"},
		{TIFFMultiJPEGTN2, "tiff-multi-jpeg-tn2"},
		{XPS, "xps"},
	}

	for _, tc := range standardValues {
		t.Run(tc.textValue, func(t *testing.T) {
			orig := FormatElement{
				Value: tc.formatValue,
			}

			elm := toXMLFormatElement(orig, NsWSCN+":Format")
			if elm.Text != tc.textValue {
				t.Errorf("expected text '%s', got '%s'", tc.textValue, elm.Text)
			}

			decoded, err := decodeFormatElement(elm)
			if err != nil {
				t.Fatalf("decode returned error: %v", err)
			}
			if decoded.Value != tc.formatValue {
				t.Errorf("expected value %v, got %v", tc.formatValue, decoded.Value)
			}
		})
	}
}

func TestFormatElement_VendorDefinedValues(t *testing.T) {
	// Test that vendor-defined values decode to UnknownFormatValue
	vendorValues := []string{"vendor-format-1", "custom-format", "extended-value"}

	for _, val := range vendorValues {
		t.Run(val, func(t *testing.T) {
			root := xmldoc.Element{
				Name: NsWSCN + ":Format",
				Text: val,
			}

			decoded, err := decodeFormatElement(root)
			if err != nil {
				t.Fatalf("decode returned error for vendor-defined value '%s': %v", val, err)
			}
			if decoded.Value != UnknownFormatValue {
				t.Errorf("expected UnknownFormatValue, got %v", decoded.Value)
			}
			// When encoding UnknownFormatValue, it will return "Unknown"
			elm := toXMLFormatElement(decoded, NsWSCN+":Format")
			if elm.Text != "Unknown" {
				t.Errorf("expected text 'Unknown' for UnknownFormatValue, got '%s'", elm.Text)
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

	if decoded.Value != TIFFMultiG4 {
		t.Errorf("expected value TIFFMultiG4, got %v", decoded.Value)
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
		Value:    PDFA,
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
	if decoded.Value != PDFA {
		t.Errorf("expected value PDFA, got %v", decoded.Value)
	}
}

func TestFormatElement_OnlyUsedDefault(t *testing.T) {
	// Test with only UsedDefault attribute
	orig := FormatElement{
		Value:       EXIF,
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
	if decoded.Value != EXIF {
		t.Errorf("expected value EXIF, got %v", decoded.Value)
	}
}
