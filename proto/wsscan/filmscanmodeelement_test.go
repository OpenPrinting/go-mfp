// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for FilmScanModeElement

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestFilmScanModeElement_RoundTrip(t *testing.T) {
	orig := FilmScanModeElement{
		Value:       "ColorSlideFilm",
		MustHonor:   optional.New(BooleanElement("true")),
		Override:    optional.New(BooleanElement("false")),
		UsedDefault: optional.New(BooleanElement("1")),
	}

	elm := toXMLFilmScanModeElement(orig, NsWSCN+":FilmScanMode")

	if elm.Name != NsWSCN+":FilmScanMode" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":FilmScanMode", elm.Name)
	}
	if elm.Text != "ColorSlideFilm" {
		t.Errorf("expected text 'ColorSlideFilm', got '%s'", elm.Text)
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
	decoded, err := decodeFilmScanModeElement(elm)
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

func TestFilmScanModeElement_NoAttributes(t *testing.T) {
	orig := FilmScanModeElement{
		Value: "NotApplicable",
	}

	elm := toXMLFilmScanModeElement(orig, NsWSCN+":FilmScanMode")

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}
	if elm.Text != "NotApplicable" {
		t.Errorf("expected text 'NotApplicable', got '%s'", elm.Text)
	}

	decoded, err := decodeFilmScanModeElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestFilmScanModeElement_StandardValues(t *testing.T) {
	standardValues := []string{
		"NotApplicable",
		"ColorSlideFilm",
		"ColorNegativeFilm",
		"BlackandWhiteNegativeFilm",
	}

	for _, val := range standardValues {
		t.Run(val, func(t *testing.T) {
			orig := FilmScanModeElement{
				Value: val,
			}

			elm := toXMLFilmScanModeElement(orig, NsWSCN+":FilmScanMode")
			if elm.Text != val {
				t.Errorf("expected text '%s', got '%s'", val, elm.Text)
			}

			decoded, err := decodeFilmScanModeElement(elm)
			if err != nil {
				t.Fatalf("decode returned error: %v", err)
			}
			if decoded.Value != val {
				t.Errorf("expected value %s, got %s", val, decoded.Value)
			}
		})
	}
}

func TestFilmScanModeElement_ExtendedValues(t *testing.T) {
	// Test that extended values are accepted (as per spec: "You can both extend and subset values")
	extendedValues := []string{"CustomFilmType", "AnotherType", "ExtendedValue"}

	for _, val := range extendedValues {
		t.Run(val, func(t *testing.T) {
			orig := FilmScanModeElement{
				Value: val,
			}

			elm := toXMLFilmScanModeElement(orig, NsWSCN+":FilmScanMode")
			decoded, err := decodeFilmScanModeElement(elm)
			if err != nil {
				t.Fatalf("decode returned error for extended value '%s': %v", val, err)
			}
			if decoded.Value != val {
				t.Errorf("expected value %s, got %s", val, decoded.Value)
			}
		})
	}
}

func TestFilmScanModeElement_FromXML(t *testing.T) {
	// Create XML element manually with all attributes
	root := xmldoc.Element{
		Name: NsWSCN + ":FilmScanMode",
		Text: "ColorNegativeFilm",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "0"},
			{Name: NsWSCN + ":Override", Value: "1"},
			{Name: NsWSCN + ":UsedDefault", Value: "false"},
		},
	}

	decoded, err := decodeFilmScanModeElement(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Value != "ColorNegativeFilm" {
		t.Errorf("expected value 'ColorNegativeFilm', got '%s'", decoded.Value)
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

func TestFilmScanModeElement_InvalidBooleanAttributes(t *testing.T) {
	// Test that invalid boolean values in attributes are rejected
	root := xmldoc.Element{
		Name: NsWSCN + ":FilmScanMode",
		Text: "ColorSlideFilm",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
	}

	_, err := decodeFilmScanModeElement(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}
