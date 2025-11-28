// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for CompressionQualityFactor

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestCompressionQualityFactor_RoundTrip(t *testing.T) {
	orig := CompressionQualityFactor{
		Value:       85,
		MustHonor:   optional.New(BooleanElement("true")),
		Override:    optional.New(BooleanElement("false")),
		UsedDefault: optional.New(BooleanElement("1")),
	}

	elm := toXMLCompressionQualityFactor(orig, NsWSCN+":CompressionQualityFactor")

	if elm.Name != NsWSCN+":CompressionQualityFactor" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":CompressionQualityFactor", elm.Name)
	}
	if elm.Text != "85" {
		t.Errorf("expected text '85', got '%s'", elm.Text)
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
	decoded, err := decodeCompressionQualityFactor(elm)
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

func TestCompressionQualityFactor_NoAttributes(t *testing.T) {
	orig := CompressionQualityFactor{
		Value: 50,
	}

	elm := toXMLCompressionQualityFactor(orig, NsWSCN+":CompressionQualityFactor")

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}
	if elm.Text != "50" {
		t.Errorf("expected text '50', got '%s'", elm.Text)
	}

	decoded, err := decodeCompressionQualityFactor(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestCompressionQualityFactor_Validation(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
		wantVal int
	}{
		{
			name:    "valid minimum",
			text:    "0",
			wantErr: false,
			wantVal: 0,
		},
		{
			name:    "valid maximum",
			text:    "100",
			wantErr: false,
			wantVal: 100,
		},
		{
			name:    "valid middle",
			text:    "50",
			wantErr: false,
			wantVal: 50,
		},
		{
			name:    "invalid negative",
			text:    "-1",
			wantErr: true,
		},
		{
			name:    "invalid too large",
			text:    "101",
			wantErr: true,
		},
		{
			name:    "invalid not a number",
			text:    "abc",
			wantErr: true,
		},
		{
			name:    "invalid empty",
			text:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := xmldoc.Element{
				Name: NsWSCN + ":CompressionQualityFactor",
				Text: tt.text,
			}

			decoded, err := decodeCompressionQualityFactor(root)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if decoded.Value != tt.wantVal {
					t.Errorf("expected value %d, got %d", tt.wantVal, decoded.Value)
				}
			}
		})
	}
}

func TestCompressionQualityFactor_FromXML(t *testing.T) {
	// Create XML element manually with all attributes
	root := xmldoc.Element{
		Name: NsWSCN + ":CompressionQualityFactor",
		Text: "75",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "0"},
			{Name: NsWSCN + ":Override", Value: "1"},
			{Name: NsWSCN + ":UsedDefault", Value: "false"},
		},
	}

	decoded, err := decodeCompressionQualityFactor(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Value != 75 {
		t.Errorf("expected value 75, got %v", decoded.Value)
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
