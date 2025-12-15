// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for ExposureSettings

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestExposureSettings_RoundTrip(t *testing.T) {
	orig := ExposureSettings{
		Brightness: optional.New(AttributedElement[int]{
			Value:    50,
			Override: optional.New(BooleanElement("1")),
		}),
		Contrast: optional.New(AttributedElement[int]{
			Value: 75,
		}),
	}

	elm := orig.toXML(NsWSCN + ":ExposureSettings")

	if elm.Name != NsWSCN+":ExposureSettings" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":ExposureSettings", elm.Name)
	}

	decoded, err := decodeExposureSettings(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if !reflect.DeepEqual(optional.Get(orig.Brightness), optional.Get(decoded.Brightness)) {
		t.Errorf("Brightness mismatch: expected %+v, got %+v",
			optional.Get(orig.Brightness), optional.Get(decoded.Brightness))
	}
	if !reflect.DeepEqual(optional.Get(orig.Contrast), optional.Get(decoded.Contrast)) {
		t.Errorf("Contrast mismatch: expected %+v, got %+v",
			optional.Get(orig.Contrast), optional.Get(decoded.Contrast))
	}
	if decoded.Sharpness != nil {
		t.Errorf("expected Sharpness to be nil, got %+v", optional.Get(decoded.Sharpness))
	}
}

func TestExposureSettings_FromXML(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":ExposureSettings",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Brightness",
				Text: "50",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":Override", Value: "0"},
				},
			},
			{
				Name: NsWSCN + ":Sharpness",
				Text: "25",
			},
		},
	}

	decoded, err := decodeExposureSettings(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if optional.Get(decoded.Brightness).Value != 50 {
		t.Errorf("expected Brightness value 50, got '%v'",
			optional.Get(decoded.Brightness).Value)
	}
	if override := optional.Get(optional.Get(decoded.Brightness).Override); override != "0" {
		t.Errorf("expected Brightness Override='0', got '%s'", override)
	}
	if optional.Get(decoded.Sharpness).Value != 25 {
		t.Errorf("expected Sharpness value 25, got '%v'",
			optional.Get(decoded.Sharpness).Value)
	}
}

func TestExposureSettings_InvalidIntegerValue(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":ExposureSettings",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Brightness",
				Text: "invalid",
			},
		},
	}

	_, err := decodeExposureSettings(root)
	if err == nil {
		t.Errorf("expected error for invalid Brightness value, got nil")
	}
}
