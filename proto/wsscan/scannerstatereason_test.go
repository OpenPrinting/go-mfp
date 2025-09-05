// MFP - Miulti-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scan state reason tests

package wsscan

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestScannerStateReason_String(t *testing.T) {
	tests := []struct {
		name     string
		reason   ScannerStateReason
		expected string
	}{
		{"Unknown", UnknownScannerStateReason, "Unknown"},
		{"AttentionRequired", StateAttentionRequired, "AttentionRequired"},
		{"Calibrating", StateCalibrating, "Calibrating"},
		{"CoverOpen", StateCoverOpen, "CoverOpen"},
		{"InterlockOpen", StateInterlockOpen, "InterlockOpen"},
		{"InternalStorageFull", StateInternalStorageFull, "InternalStorageFull"},
		{"LampError", StateLampError, "LampError"},
		{"LampWarming", StateLampWarming, "LampWarming"},
		{"MediaJam", StateMediaJam, "MediaJam"},
		{"MultipleFeedError", StateMultipleFeedError, "MultipleFeedError"},
		{"None", StateNone, "None"},
		{"Paused", StatePaused, "Paused"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reason.String(); got != tt.expected {
				t.Errorf("ScannerStateReason.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDecodeScannerStateReason(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ScannerStateReason
	}{
		{"Unknown", "Unknown", UnknownScannerStateReason},
		{"AttentionRequired", "AttentionRequired", StateAttentionRequired},
		{"Calibrating", "Calibrating", StateCalibrating},
		{"CoverOpen", "CoverOpen", StateCoverOpen},
		{"InterlockOpen", "InterlockOpen", StateInterlockOpen},
		{"InternalStorageFull", "InternalStorageFull", StateInternalStorageFull},
		{"LampError", "LampError", StateLampError},
		{"LampWarming", "LampWarming", StateLampWarming},
		{"MediaJam", "MediaJam", StateMediaJam},
		{"MultipleFeedError", "MultipleFeedError", StateMultipleFeedError},
		{"None", "None", StateNone},
		{"Paused", "Paused", StatePaused},
		{"Empty", "", UnknownScannerStateReason},
		{"Invalid", "InvalidReason", UnknownScannerStateReason},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeScannerStateReason(tt.input); got != tt.expected {
				t.Errorf("DecodeScannerStateReason() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestScannerStateReason_toXML(t *testing.T) {
	tests := []struct {
		name     string
		reason   ScannerStateReason
		xmlName  string
		expected xmldoc.Element
	}{
		{
			name:    "AttentionRequired",
			reason:  StateAttentionRequired,
			xmlName: "ScannerStateReason",
			expected: xmldoc.Element{
				Name: "ScannerStateReason",
				Text: "AttentionRequired",
			},
		},
		{
			name:    "None",
			reason:  StateNone,
			xmlName: "ScannerStateReason",
			expected: xmldoc.Element{
				Name: "ScannerStateReason",
				Text: "None",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.reason.toXML(tt.xmlName)
			if got.Name != tt.expected.Name {
				t.Errorf("toXML().Name = %v, want %v", got.Name, tt.expected.Name)
			}
			if got.Text != tt.expected.Text {
				t.Errorf("toXML().Text = %v, want %v", got.Text, tt.expected.Text)
			}
		})
	}
}

func Test_decodeScannerStateReason(t *testing.T) {
	tests := []struct {
		name     string
		element  xmldoc.Element
		expected ScannerStateReason
		wantErr  bool
	}{
		{
			name: "Valid AttentionRequired",
			element: xmldoc.Element{
				Name: "ScannerStateReason",
				Text: "AttentionRequired",
			},
			expected: StateAttentionRequired,
			wantErr:  false,
		},
		{
			name: "Valid None",
			element: xmldoc.Element{
				Name: "ScannerStateReason",
				Text: "None",
			},
			expected: StateNone,
			wantErr:  false,
		},
		{
			name: "Invalid value",
			element: xmldoc.Element{
				Name: "ScannerStateReason",
				Text: "InvalidReason",
			},
			expected: UnknownScannerStateReason,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeScannerStateReason(tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeScannerStateReason() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("decodeScannerStateReason() = %v, want %v", got, tt.expected)
			}
		})
	}
}
