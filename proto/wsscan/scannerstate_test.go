// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scanner state tests

package wsscan

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestScannerState_String(t *testing.T) {
	tests := []struct {
		name     string
		state    ScannerState
		expected string
	}{
		{"Unknown", UnknownScannerState, "Unknown"},
		{"Idle", Idle, "Idle"},
		{"Processing", Processing, "Processing"},
		{"Stopped", Stopped, "Stopped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("ScannerState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDecodeScannerState(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ScannerState
	}{
		{"Unknown", "Unknown", UnknownScannerState},
		{"Idle", "Idle", Idle},
		{"Processing", "Processing", Processing},
		{"Stopped", "Stopped", Stopped},
		{"Empty", "", UnknownScannerState},
		{"Invalid", "InvalidState", UnknownScannerState},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeScannerState(tt.input); got != tt.expected {
				t.Errorf("DecodeScannerState() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestScannerState_toXML(t *testing.T) {
	tests := []struct {
		name     string
		state    ScannerState
		xmlName  string
		expected xmldoc.Element
	}{
		{
			name:    "Idle",
			state:   Idle,
			xmlName: "ScannerState",
			expected: xmldoc.Element{
				Name: "ScannerState",
				Text: "Idle",
			},
		},
		{
			name:    "Processing",
			state:   Processing,
			xmlName: "ScannerState",
			expected: xmldoc.Element{
				Name: "ScannerState",
				Text: "Processing",
			},
		},
		{
			name:    "Stopped",
			state:   Stopped,
			xmlName: "ScannerState",
			expected: xmldoc.Element{
				Name: "ScannerState",
				Text: "Stopped",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.state.toXML(tt.xmlName)
			if got.Name != tt.expected.Name {
				t.Errorf("toXML().Name = %v, want %v", got.Name, tt.expected.Name)
			}
			if got.Text != tt.expected.Text {
				t.Errorf("toXML().Text = %v, want %v", got.Text, tt.expected.Text)
			}
		})
	}
}

func Test_decodeScannerState(t *testing.T) {
	tests := []struct {
		name     string
		element  xmldoc.Element
		expected ScannerState
		wantErr  bool
	}{
		{
			name: "Valid Idle",
			element: xmldoc.Element{
				Name: "ScannerState",
				Text: "Idle",
			},
			expected: Idle,
			wantErr:  false,
		},
		{
			name: "Valid Processing",
			element: xmldoc.Element{
				Name: "ScannerState",
				Text: "Processing",
			},
			expected: Processing,
			wantErr:  false,
		},
		{
			name: "Valid Stopped",
			element: xmldoc.Element{
				Name: "ScannerState",
				Text: "Stopped",
			},
			expected: Stopped,
			wantErr:  false,
		},
		{
			name: "Invalid value",
			element: xmldoc.Element{
				Name: "ScannerState",
				Text: "InvalidState",
			},
			expected: UnknownScannerState,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeScannerState(tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeScannerState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("decodeScannerState() = %v, want %v", got, tt.expected)
			}
		})
	}
}
