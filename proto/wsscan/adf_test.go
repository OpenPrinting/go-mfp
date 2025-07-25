// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for ADF

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestADF_RoundTrip(t *testing.T) {
	adfBack := ADFFeederSideElement{
		ADFColor: []ColorEntry{BlackAndWhite1},
		ADFMaximumSize: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "297"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "210"}},
		},
		ADFMinimumSize: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "100"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "50"}},
		},
		ADFOpticalResolution: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "600"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "600"}},
		},
		ADFResolutions: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "300"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "300"}},
		},
	}
	adfFront := ADFFeederSideElement{
		ADFColor: []ColorEntry{RGB24},
		ADFMaximumSize: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "297"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "210"}},
		},
		ADFMinimumSize: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "100"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "50"}},
		},
		ADFOpticalResolution: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "600"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "600"}},
		},
		ADFResolutions: HeightAndWidthElement{
			Heights: []TextWithOverrideAndDefault{{Text: "300"}},
			Widths:  []TextWithOverrideAndDefault{{Text: "300"}},
		},
	}
	orig := ADF{
		ADFBack:           optional.New(adfBack),
		ADFFront:          optional.New(adfFront),
		ADFSupportsDuplex: BooleanElement("true"),
	}
	elm := orig.toXML("wscn:ADF")
	if elm.Name != "wscn:ADF" {
		t.Errorf("expected element name 'wscn:ADF', got '%s'", elm.Name)
	}

	parsed, err := decodeADF(elm)
	if err != nil {
		t.Fatalf("decodeADF returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, parsed) {
		t.Errorf("expected %+v, got %+v", orig, parsed)
	}
}

func TestADF_OptionalFields(t *testing.T) {
	orig := ADF{
		ADFSupportsDuplex: BooleanElement("false"),
	}
	elm := orig.toXML("wscn:ADF")
	parsed, err := decodeADF(elm)
	if err != nil {
		t.Fatalf("decodeADF returned error: %v", err)
	}
	if !reflect.DeepEqual(orig, parsed) {
		t.Errorf("expected %+v, got %+v", orig, parsed)
	}
}

func TestADF_InvalidBoolean(t *testing.T) {
	elm := xmldoc.Element{
		Name: "wscn:ADF",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":ADFSupportsDuplex",
				Text: "maybe",
			},
		},
	}
	_, err := decodeADF(elm)
	if err == nil {
		t.Errorf("expected error for invalid boolean value, got nil")
	}
}
