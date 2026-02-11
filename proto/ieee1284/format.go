// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Document format detection

package ieee1284

import (
	"bytes"
	"strings"
)

// DocFormat represents a detected document format.
type DocFormat int
// Known values for DocFormat.
const (
	DocFormatUnknown    DocFormat = iota // Unknown format
	DocFormatPostScript                  // PostScript
	DocFormatPDF                         // PDF
	DocFormatPCL                         // PCL 5
	DocFormatPCLXL5                       // PCL-XL / PCL 6
)

// String returns a human-readable name for the format.
func (f DocFormat) String() string {
	switch f {
	case DocFormatPostScript:
		return "PostScript"
	case DocFormatPDF:
		return "PDF"
	case DocFormatPCL:
		return "PCL"
	case DocFormatPCLXL5:
		return "PCL-XL"
	default:
		return "Unknown"
	}
}

// magicEntry maps a byte prefix to a document format.
type magicEntry struct {
	prefix []byte
	format DocFormat
}

// magicTable lists document format signatures, checked in order.
// More specific prefixes come before less specific ones.
var magicTable = []magicEntry{
	{[]byte("%!PS"), DocFormatPostScript},
	{[]byte("%PDF-"), DocFormatPDF},
	{[]byte(") HP-PCL XL"), DocFormatPCLXL5},
}

// detectFormatByMagic detects a document format by magic byte prefix.
// If the data doesn't match any known signature, it checks for a bare
// ESC byte which indicates PCL 5.
func detectFormatByMagic(data []byte) DocFormat {
	for _, m := range magicTable {
		if bytes.HasPrefix(data, m.prefix) {
			return m.format
		}
	}

	// PCL 5: starts with ESC followed by a PCL command character.
	// PCL commands use characters in the range 0x21-0x7e after ESC,
	// but not '%' (which would be UEL).
	if len(data) >= 2 && data[0] == 0x1b && data[1] != '%' {
		return DocFormatPCL
	}

	return DocFormatUnknown
}

// detectFormatByLanguage maps a PJL ENTER LANGUAGE value to a
// document format. The comparison is case-insensitive.
func detectFormatByLanguage(lang string) DocFormat {
	switch strings.ToUpper(lang) {
	case "POSTSCRIPT":
		return DocFormatPostScript
	case "PDF":
		return DocFormatPDF
	case "PCL":
		return DocFormatPCL
	case "PCLXL":
		return DocFormatPCLXL5
	default:
		return DocFormatUnknown
	}
}
