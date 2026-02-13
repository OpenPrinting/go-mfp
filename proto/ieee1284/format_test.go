// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Document format detection tests

package ieee1284

import (
	"testing"
)

func TestDetectFormatByMagic(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		expect DocFormat
	}{
		{"PostScript", []byte("%!PS-Adobe-3.0\n"), DocFormatPostScript},
		{"PDF", []byte("%PDF-1.4\n"), DocFormatPDF},
		{"PCL-XL", []byte(") HP-PCL XL;2;0\n"), DocFormatPCLXL},
		{"PCL5 ESC-E", []byte{0x1b, 'E'}, DocFormatPCL5},
		{"PCL5 ESC-command", []byte{0x1b, '&', 'l', '1', 'O'}, DocFormatPCL5},
		{"Unknown", []byte("Hello, World!"), DocFormatUnknown},
		{"Empty", []byte{}, DocFormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectFormatByMagic(tt.data)
			if got != tt.expect {
				t.Errorf("detectFormatByMagic(%q) = %v, want %v",
					tt.data, got, tt.expect)
			}
		})
	}
}

func TestDetectFormatByLanguage(t *testing.T) {
	tests := []struct {
		lang   string
		expect DocFormat
	}{
		{"POSTSCRIPT", DocFormatPostScript},
		{"postscript", DocFormatPostScript},
		{"PostScript", DocFormatPostScript},
		{"PDF", DocFormatPDF},
		{"pdf", DocFormatPDF},
		{"PCL", DocFormatPCL5},
		{"pcl", DocFormatPCL5},
		{"PCLXL", DocFormatPCLXL},
		{"pclxl", DocFormatPCLXL},
		{"UNKNOWN", DocFormatUnknown},
		{"", DocFormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := detectFormatByLanguage(tt.lang)
			if got != tt.expect {
				t.Errorf("detectFormatByLanguage(%q) = %v, want %v",
					tt.lang, got, tt.expect)
			}
		})
	}
}

// TestDocFormatString tests the String() method of DocFormat.
func TestDocFormatString(t *testing.T) {
	tests := []struct {
		format DocFormat
		expect string
	}{
		{DocFormatUnknown, "Unknown"},
		{DocFormatPostScript, "PostScript"},
		{DocFormatPDF, "PDF"},
		{DocFormatPCL5, "PCL"},
		{DocFormatPCLXL, "PCL-XL"},
	}

	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			got := tt.format.String()
			if got != tt.expect {
				t.Errorf("DocFormat(%d).String() = %q, want %q",
					tt.format, got, tt.expect)
			}
		})
	}
}
