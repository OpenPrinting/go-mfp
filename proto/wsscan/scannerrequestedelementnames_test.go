// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Name for RequestedElements element tests

package wsscan

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TestScannerRequestedElementNames_String checks that each known
// ScannerRequestedElementNames produces the correct QName string.
func TestScannerRequestedElementNames_String(t *testing.T) {
	tests := []struct {
		name     string
		sren     ScannerRequestedElementNames
		expected string
	}{
		{"Unknown", UnknownScannerRequestedElementNames, "Unknown"},
		{"Description", ScannerRequestedElementDescription,
			NsWSCN + ":ScannerDescription"},
		{"Configuration", ScannerRequestedElementConfiguration,
			NsWSCN + ":ScannerConfiguration"},
		{"Status", ScannerRequestedElementStatus,
			NsWSCN + ":ScannerStatus"},
		{"VendorSection", ScannerRequestedElementVendorSection,
			NsXML + ":VendorSection"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sren.String(); got != tt.expected {
				t.Errorf("ScannerRequestedElementNames.String() = %v, want %v",
					got, tt.expected)
			}
		})
	}
}

// TestDecodeScannerRequestedElementNames checks that valid QName strings
// decode to the correct constant and invalid ones return Unknown.
func TestDecodeScannerRequestedElementNames(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ScannerRequestedElementNames
	}{
		{"Description", NsWSCN + ":ScannerDescription",
			ScannerRequestedElementDescription},
		{"Configuration", NsWSCN + ":ScannerConfiguration",
			ScannerRequestedElementConfiguration},
		{"Status", NsWSCN + ":ScannerStatus",
			ScannerRequestedElementStatus},
		{"VendorSection", NsXML + ":VendorSection",
			ScannerRequestedElementVendorSection},
		{"Empty", "", UnknownScannerRequestedElementNames},
		{"Invalid", "InvalidName", UnknownScannerRequestedElementNames},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeScannerRequestedElementNames(
				tt.input); got != tt.expected {
				t.Errorf("DecodeScannerRequestedElementNames() = %v, want %v",
					got, tt.expected)
			}
		})
	}
}

// TestScannerRequestedElementNames_toXML checks that toXML produces an
// element with the correct name and QName text value.
func TestScannerRequestedElementNames_toXML(t *testing.T) {
	tests := []struct {
		name     string
		sren     ScannerRequestedElementNames
		xmlName  string
		expected xmldoc.Element
	}{
		{
			// Verify Description produces wscn:ScannerDescription text.
			name:    "Description",
			sren:    ScannerRequestedElementDescription,
			xmlName: NsWSCN + ":Name",
			expected: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsWSCN + ":ScannerDescription",
			},
		},
		{
			// Verify Configuration produces wscn:ScannerConfiguration text.
			name:    "Configuration",
			sren:    ScannerRequestedElementConfiguration,
			xmlName: NsWSCN + ":Name",
			expected: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsWSCN + ":ScannerConfiguration",
			},
		},
		{
			// Verify Status produces wscn:ScannerStatus text.
			name:    "Status",
			sren:    ScannerRequestedElementStatus,
			xmlName: NsWSCN + ":Name",
			expected: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsWSCN + ":ScannerStatus",
			},
		},
		{
			// Verify VendorSection produces xmlns:VendorSection text.
			name:    "VendorSection",
			sren:    ScannerRequestedElementVendorSection,
			xmlName: NsWSCN + ":Name",
			expected: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsXML + ":VendorSection",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sren.toXML(tt.xmlName)
			if got.Name != tt.expected.Name {
				t.Errorf("toXML().Name = %v, want %v",
					got.Name, tt.expected.Name)
			}
			if got.Text != tt.expected.Text {
				t.Errorf("toXML().Text = %v, want %v",
					got.Text, tt.expected.Text)
			}
		})
	}
}

// Test_decodeScannerRequestedElementNames checks that decoding from XML
// elements works for valid values and returns an error for invalid ones.
func Test_decodeScannerRequestedElementNames(t *testing.T) {
	tests := []struct {
		name     string
		element  xmldoc.Element
		expected ScannerRequestedElementNames
		wantErr  bool
	}{
		{
			// Valid wscn:ScannerDescription.
			name: "Valid Description",
			element: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsWSCN + ":ScannerDescription",
			},
			expected: ScannerRequestedElementDescription,
			wantErr:  false,
		},
		{
			// Valid wscn:ScannerConfiguration.
			name: "Valid Configuration",
			element: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsWSCN + ":ScannerConfiguration",
			},
			expected: ScannerRequestedElementConfiguration,
			wantErr:  false,
		},
		{
			// Valid wscn:ScannerStatus.
			name: "Valid Status",
			element: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsWSCN + ":ScannerStatus",
			},
			expected: ScannerRequestedElementStatus,
			wantErr:  false,
		},
		{
			// Valid xmlns:VendorSection.
			name: "Valid VendorSection",
			element: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: NsXML + ":VendorSection",
			},
			expected: ScannerRequestedElementVendorSection,
			wantErr:  false,
		},
		{
			// Invalid value must return an error.
			name: "Invalid value",
			element: xmldoc.Element{
				Name: NsWSCN + ":Name",
				Text: "InvalidName",
			},
			expected: UnknownScannerRequestedElementNames,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeScannerRequestedElementNames(tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeScannerRequestedElementNames() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("decodeScannerRequestedElementNames() = %v, want %v",
					got, tt.expected)
			}
		})
	}
}
