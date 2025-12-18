// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for InputSourceValue

package wsscan

import "testing"

func TestInputSourceValue_String(t *testing.T) {
	tests := []struct {
		value    InputSourceValue
		expected string
	}{
		{InputSourceADF, "ADF"},
		{InputSourceADFDuplex, "ADFDuplex"},
		{InputSourceFilm, "Film"},
		{InputSourcePlaten, "Platen"},
		{UnknownInputSourceValue, "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.value.String(); got != tt.expected {
			t.Errorf("InputSourceValue(%d).String() = %q, want %q", tt.value, got, tt.expected)
		}
	}
}

func TestDecodeInputSourceValue(t *testing.T) {
	tests := []struct {
		input    string
		expected InputSourceValue
	}{
		{"ADF", InputSourceADF},
		{"ADFDuplex", InputSourceADFDuplex},
		{"Film", InputSourceFilm},
		{"Platen", InputSourcePlaten},
		{"Unknown", UnknownInputSourceValue},
		{"invalid", UnknownInputSourceValue},
		{"", UnknownInputSourceValue},
	}

	for _, tt := range tests {
		if got := DecodeInputSourceValue(tt.input); got != tt.expected {
			t.Errorf("DecodeInputSourceValue(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
