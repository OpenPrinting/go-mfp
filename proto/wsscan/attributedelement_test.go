// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for AttributedElement

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestAttributedElement_RoundTrip(t *testing.T) {
	orig := AttributedElement[RotationValue]{
		Value:       Rotation90,
		MustHonor:   optional.New(BooleanElement("true")),
		Override:    optional.New(BooleanElement("false")),
		UsedDefault: optional.New(BooleanElement("true")),
	}

	elm := orig.toXML(NsWSCN+":Rotation", func(rv RotationValue) string {
		return rv.String()
	})

	if elm.Name != NsWSCN+":Rotation" {
		t.Errorf("expected element name '%s', got '%s'", NsWSCN+":Rotation", elm.Name)
	}
	if elm.Text != "90" {
		t.Errorf("expected text '90', got '%s'", elm.Text)
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
	if attrsMap[NsWSCN+":UsedDefault"] != "true" {
		t.Errorf("expected UsedDefault='true', got '%s'", attrsMap[NsWSCN+":UsedDefault"])
	}

	// Decode back
	decoded, err := decodeAttributedElement(elm, func(s string) (RotationValue, error) {
		return DecodeRotationValue(s), nil
	})
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

func TestAttributedElement_NoAttributes(t *testing.T) {
	orig := AttributedElement[RotationValue]{
		Value: Rotation180,
	}

	elm := orig.toXML(NsWSCN+":Rotation", func(rv RotationValue) string {
		return rv.String()
	})

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}

	decoded, err := decodeAttributedElement(elm, func(s string) (RotationValue, error) {
		return DecodeRotationValue(s), nil
	})
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestAttributedElement_StringValue(t *testing.T) {
	orig := AttributedElement[string]{
		Value:     "some-value",
		MustHonor: optional.New(BooleanElement("1")),
	}

	elm := orig.toXML(NsWSCN+":SomeElement", func(s string) string {
		return s
	})

	if elm.Text != "some-value" {
		t.Errorf("expected text 'some-value', got '%s'", elm.Text)
	}
	if len(elm.Attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(elm.Attrs))
	}

	decoded, err := decodeAttributedElement(elm, func(s string) (string, error) {
		return s, nil
	})
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestAttributedElement_FromXML(t *testing.T) {
	// Create XML element manually
	root := xmldoc.Element{
		Name: NsWSCN + ":Rotation",
		Text: "270",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "false"},
			{Name: NsWSCN + ":Override", Value: "true"},
		},
	}

	decoded, err := decodeAttributedElement(root, func(s string) (RotationValue, error) {
		val := DecodeRotationValue(s)
		if val == UnknownRotationValue {
			return val, xmldoc.XMLErrWrap(root, nil)
		}
		return val, nil
	})
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Value != Rotation270 {
		t.Errorf("expected Rotation270, got %v", decoded.Value)
	}
	if mustHonor := optional.Get(decoded.MustHonor); string(mustHonor) != "false" {
		t.Errorf("expected MustHonor='false', got '%s'", mustHonor)
	}
	if override := optional.Get(decoded.Override); string(override) != "true" {
		t.Errorf("expected Override='true', got '%s'", override)
	}
	if usedDefault := optional.Get(decoded.UsedDefault); usedDefault != "" {
		t.Errorf("expected empty UsedDefault, got '%s'", usedDefault)
	}
}

func TestAttributedElement_InvalidBooleanAttributes(t *testing.T) {
	tests := []struct {
		name    string
		attr    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid true",
			attr:    "MustHonor",
			value:   "true",
			wantErr: false,
		},
		{
			name:    "valid false",
			attr:    "MustHonor",
			value:   "false",
			wantErr: false,
		},
		{
			name:    "valid 1",
			attr:    "MustHonor",
			value:   "1",
			wantErr: false,
		},
		{
			name:    "valid 0",
			attr:    "MustHonor",
			value:   "0",
			wantErr: false,
		},
		{
			name:    "invalid value",
			attr:    "MustHonor",
			value:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid empty",
			attr:    "MustHonor",
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid Override",
			attr:    "Override",
			value:   "yes",
			wantErr: true,
		},
		{
			name:    "invalid UsedDefault",
			attr:    "UsedDefault",
			value:   "maybe",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := xmldoc.Element{
				Name: NsWSCN + ":Rotation",
				Text: "90",
				Attrs: []xmldoc.Attr{
					{Name: NsWSCN + ":" + tt.attr, Value: tt.value},
				},
			}

			_, err := decodeAttributedElement(root, func(s string) (RotationValue, error) {
				val := DecodeRotationValue(s)
				if val == UnknownRotationValue {
					return val, xmldoc.XMLErrWrap(root, nil)
				}
				return val, nil
			})

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %s='%s', got nil", tt.attr, tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s='%s': %v", tt.attr, tt.value, err)
				}
			}
		})
	}
}
