// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Rotation

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestRotation_RoundTrip(t *testing.T) {
	orig := Rotation{
		Value:       Rotation90,
		MustHonor:   optional.New(BooleanElement("true")),
		Override:    optional.New(BooleanElement("false")),
		UsedDefault: optional.New(BooleanElement("1")),
	}

	elm := toXMLRotation(orig, NsWSCN+":Rotation")

	if elm.Name != NsWSCN+":Rotation" {
		t.Errorf("expected element name '%s', got '%s'",
			NsWSCN+":Rotation", elm.Name)
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
	if attrsMap[NsWSCN+":UsedDefault"] != "1" {
		t.Errorf("expected UsedDefault='1', got '%s'", attrsMap[NsWSCN+":UsedDefault"])
	}

	// Decode back
	decoded, err := decodeRotation(elm)
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

func TestRotation_NoAttributes(t *testing.T) {
	orig := Rotation{
		Value: Rotation180,
	}

	elm := toXMLRotation(orig, NsWSCN+":Rotation")

	if len(elm.Attrs) != 0 {
		t.Errorf("expected no attributes, got %+v", elm.Attrs)
	}
	if elm.Text != "180" {
		t.Errorf("expected text '180', got '%s'", elm.Text)
	}

	decoded, err := decodeRotation(elm)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}
	if decoded.Value != orig.Value {
		t.Errorf("expected value %v, got %v", orig.Value, decoded.Value)
	}
}

func TestRotation_AllValidValues(t *testing.T) {
	validValues := []struct {
		enumValue RotationValue
		textValue string
	}{
		{Rotation0, "0"},
		{Rotation90, "90"},
		{Rotation180, "180"},
		{Rotation270, "270"},
	}

	for _, tc := range validValues {
		t.Run(tc.textValue, func(t *testing.T) {
			orig := Rotation{
				Value: tc.enumValue,
			}

			elm := toXMLRotation(orig, NsWSCN+":Rotation")
			if elm.Text != tc.textValue {
				t.Errorf("expected text '%s', got '%s'", tc.textValue, elm.Text)
			}

			decoded, err := decodeRotation(elm)
			if err != nil {
				t.Fatalf("decode returned error: %v", err)
			}
			if decoded.Value != tc.enumValue {
				t.Errorf("expected value %v, got %v", tc.enumValue, decoded.Value)
			}
		})
	}
}

func TestRotation_FromXML(t *testing.T) {
	// Create XML element manually with all attributes
	root := xmldoc.Element{
		Name: NsWSCN + ":Rotation",
		Text: "270",
		Attrs: []xmldoc.Attr{
			{Name: NsWSCN + ":MustHonor", Value: "0"},
			{Name: NsWSCN + ":Override", Value: "1"},
			{Name: NsWSCN + ":UsedDefault", Value: "false"},
		},
	}

	decoded, err := decodeRotation(root)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Value != Rotation270 {
		t.Errorf("expected value Rotation270, got %v", decoded.Value)
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

func TestRotation_InvalidValues(t *testing.T) {
	invalidValues := []string{
		"45",
		"91",
		"179",
		"271",
		"360",
		"-90",
		"invalid",
		"",
		" 90 ",
	}

	for _, val := range invalidValues {
		t.Run(val, func(t *testing.T) {
			root := xmldoc.Element{
				Name: NsWSCN + ":Rotation",
				Text: val,
			}

			_, err := decodeRotation(root)
			if err == nil {
				t.Errorf("expected error for invalid rotation value '%s', got nil", val)
			}
		})
	}
}

func TestRotation_InvalidBooleanAttributes(t *testing.T) {
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

			_, err := decodeRotation(root)
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

func TestRotation_AllValidValuesWithAttributes(t *testing.T) {
	validValues := []struct {
		enumValue RotationValue
		textValue string
	}{
		{Rotation0, "0"},
		{Rotation90, "90"},
		{Rotation180, "180"},
		{Rotation270, "270"},
	}

	for _, tc := range validValues {
		t.Run(tc.textValue, func(t *testing.T) {
			orig := Rotation{
				Value:       tc.enumValue,
				MustHonor:   optional.New(BooleanElement("1")),
				Override:    optional.New(BooleanElement("0")),
				UsedDefault: optional.New(BooleanElement("true")),
			}

			elm := toXMLRotation(orig, NsWSCN+":Rotation")
			decoded, err := decodeRotation(elm)
			if err != nil {
				t.Fatalf("decode returned error for value '%s': %v", tc.textValue, err)
			}
			if decoded.Value != tc.enumValue {
				t.Errorf("expected value %v, got %v", tc.enumValue, decoded.Value)
			}
			if len(elm.Attrs) != 3 {
				t.Errorf("expected 3 attributes for value '%s', got %d", tc.textValue, len(elm.Attrs))
			}
		})
	}
}

func TestRotation_BoundaryValues(t *testing.T) {
	// Test all valid rotation values
	boundaryTests := []struct {
		name      string
		enumValue RotationValue
		textValue string
	}{
		{"minimum", Rotation0, "0"},
		{"90 degrees", Rotation90, "90"},
		{"180 degrees", Rotation180, "180"},
		{"maximum", Rotation270, "270"},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			orig := Rotation{
				Value: tt.enumValue,
			}

			elm := toXMLRotation(orig, NsWSCN+":Rotation")
			decoded, err := decodeRotation(elm)
			if err != nil {
				t.Fatalf("decode with value %s returned error: %v", tt.textValue, err)
			}
			if decoded.Value != tt.enumValue {
				t.Errorf("expected value %v, got %v", tt.enumValue, decoded.Value)
			}
			if elm.Text != tt.textValue {
				t.Errorf("expected text '%s', got '%s'", tt.textValue, elm.Text)
			}
		})
	}
}

func TestRotation_AttributesOnAllValues(t *testing.T) {
	// Test that all rotation values work correctly with all attribute combinations
	rotationValues := []RotationValue{Rotation0, Rotation90, Rotation180, Rotation270}
	attrCombinations := []struct {
		name        string
		mustHonor   optional.Val[BooleanElement]
		override    optional.Val[BooleanElement]
		usedDefault optional.Val[BooleanElement]
	}{
		{"no attributes", nil, nil, nil},
		{"only MustHonor", optional.New(BooleanElement("true")), nil, nil},
		{"only Override", nil, optional.New(BooleanElement("false")), nil},
		{"only UsedDefault", nil, nil, optional.New(BooleanElement("1"))},
		{"all attributes", optional.New(BooleanElement("0")), optional.New(BooleanElement("1")), optional.New(BooleanElement("true"))},
	}

	for _, rotVal := range rotationValues {
		for _, attrCombo := range attrCombinations {
			t.Run(rotVal.String()+"_"+attrCombo.name, func(t *testing.T) {
				orig := Rotation{
					Value:       rotVal,
					MustHonor:   attrCombo.mustHonor,
					Override:    attrCombo.override,
					UsedDefault: attrCombo.usedDefault,
				}

				elm := toXMLRotation(orig, NsWSCN+":Rotation")
				decoded, err := decodeRotation(elm)
				if err != nil {
					t.Fatalf("decode returned error: %v", err)
				}

				if decoded.Value != rotVal {
					t.Errorf("expected value %v, got %v", rotVal, decoded.Value)
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
			})
		}
	}
}
