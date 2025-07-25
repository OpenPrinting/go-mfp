// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for TextWithLangElement

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
)

func TestTextWithLangElement_RoundTrip(t *testing.T) {
	orig := TextWithLangElement{
		Text: "Test Scanner",
		Lang: optional.New("en-US"),
	}
	elm := orig.toXML("wscn:ScannerName")
	if elm.Name != "wscn:ScannerName" {
		t.Errorf("expected element name 'wscn:ScannerName', got '%s'", elm.Name)
	}
	if elm.Text != orig.Text {
		t.Errorf("expected text '%s', got '%s'", orig.Text, elm.Text)
	}
	if len(elm.Attrs) != 1 || elm.Attrs[0].Name != "xml:lang" || elm.Attrs[0].Value != "en-US" {
		t.Errorf("expected xml:lang attribute 'en-US', got %+v", elm.Attrs)
	}

	decoded, err := orig.decodeTextWithLangElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}

func TestTextWithLangElement_EmptyLang(t *testing.T) {
	orig := TextWithLangElement{
		Text: "No Lang",
	}
	elm := orig.toXML("wscn:ScannerInfo")
	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}

	decoded, err := orig.decodeTextWithLangElement(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, decoded) {
		t.Errorf("expected %+v, got %+v", orig, decoded)
	}
}
