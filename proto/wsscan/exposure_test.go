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
		MustHonor:    optional.New(BooleanElement("true")),
		AutoExposure: BooleanElement("1"),
		ExposureSettings: ExposureSettings{
			Brightness: optional.New(AttributedElement[int]{
				Value: 50,
			}),
		},
	}

	elm := orig.toXML(NsWSCN + ":Exposure")

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

	// Decode back
	decoded, err := decodeExposure(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if !reflect.DeepEqual(orig.MustHonor, decoded.MustHonor) {
		t.Errorf("expected MustHonor %+v, got %+v", orig.MustHonor, decoded.MustHonor)
	}
	if orig.AutoExposure != decoded.AutoExposure {
		t.Errorf("expected AutoExposure %+v, got %+v", orig.AutoExposure, decoded.AutoExposure)
	}
	if !reflect.DeepEqual(orig.ExposureSettings, decoded.ExposureSettings) {
		t.Errorf("expected ExposureSettings %+v, got %+v", orig.ExposureSettings, decoded.ExposureSettings)
	}
}

func TestExposure_MissingAutoExposure(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Children: []xmldoc.Element{
			{
				Name:     NsWSCN + ":ExposureSettings",
				Children: []xmldoc.Element{},
			},
		},
	}

	_, err := decodeExposure(root)
	if err == nil {
		t.Errorf("expected error for missing AutoExposure, got nil")
	}
}

func TestExposure_MissingExposureSettings(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":AutoExposure",
				Text: "true",
			},
		},
	}

	_, err := decodeExposure(root)
	if err == nil {
		t.Errorf("expected error for missing ExposureSettings, got nil")
	}
}

func TestExposure_FromXML(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "1"},
		},
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":AutoExposure",
				Text: "true",
			},
			{
				Name: NsWSCN + ":ExposureSettings",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":Brightness",
						Text: "50",
					},
				},
			},
		},
	}

	decoded, err := decodeExposure(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "1" {
		t.Errorf("expected MustHonor='1', got '%s'", mustHonor)
	}
	if decoded.AutoExposure != BooleanElement("true") {
		t.Errorf("expected AutoExposure='true', got '%v'", decoded.AutoExposure)
	}
	expSettings := decoded.ExposureSettings
	if brightness := optional.Get(expSettings.Brightness); brightness.Value != 50 {
		t.Errorf("expected Brightness=50, got %d", brightness.Value)
	}
}

func TestExposure_InvalidBooleanAttribute(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "invalid"},
		},
	}

	_, err := decodeExposure(root)
	if err == nil {
		t.Errorf("expected error for invalid MustHonor value 'invalid', got nil")
	}
}

func TestExposure_InvalidAutoExposure(t *testing.T) {
	root := xmldoc.Element{
		Name: NsWSCN + ":Exposure",
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":AutoExposure",
				Text: "invalid",
			},
		},
	}

	_, err := decodeExposure(root)
	if err == nil {
		t.Errorf("expected error for invalid AutoExposure value 'invalid', got nil")
	}
}
