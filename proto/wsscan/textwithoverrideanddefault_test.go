// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for TextWithOverrideAndDefault

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestTextWithOverrideAndDefault_RoundTrip(t *testing.T) {
	orig := TextWithOverrideAndDefault{
		Text:        "210",
		Override:    optional.New(BooleanElement("true")),
		UsedDefault: optional.New(BooleanElement("false")),
	}
	elm := orig.toXML("wscn:Width")
	if elm.Name != "wscn:Width" {
		t.Errorf("expected element name 'wscn:Width', got '%s'", elm.Name)
	}
	if elm.Text != orig.Text {
		t.Errorf("expected text '%s', got '%s'", orig.Text, elm.Text)
	}
	if len(elm.Attrs) != 2 {
		t.Errorf("expected 2 attributes, got %+v", elm.Attrs)
	}

	decoded, err := orig.decodeTextWithOverrideAndDefault(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestTextWithOverrideAndDefault_EmptyAttrs(t *testing.T) {
	orig := TextWithOverrideAndDefault{
		Text: "297",
	}
	elm := orig.toXML("wscn:Width")
	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}

	decoded, err := orig.decodeTextWithOverrideAndDefault(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestTextWithOverrideAndDefault_InvalidBool(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:Width",
		Text: "210",
		Attrs: []xmldoc.Attr{{
			Name:  "wscn:Override",
			Value: "maybe",
		}},
	}
	var twod TextWithOverrideAndDefault
	_, err := twod.decodeTextWithOverrideAndDefault(elm)
	if err == nil {
		t.Errorf("expected error for invalid boolean attribute, got nil")
	}
}
