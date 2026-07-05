// MFP - Miulti-Function Printers and scanners toolkit
// Utility functions and data BLOBs for testing
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Test data examples for Xerox printers

package testutils

import (
	// Import "embed" for its side effects
	_ "embed"
)

// Xerox contains data samples taken from the Xerox printers
var Xerox struct {
	// M2040dn model
	B235 struct {
		// IPP protocol samples
		IPP struct {
			PrinterAttributes []byte
		}
		// ESCL protocol samples
		ESCL struct {
			ScannerCapabilities []byte
		}
	}
}

func init() {
	Xerox.B235.IPP.PrinterAttributes =
		kyoceraECOSYSM2040dnPrinterAttributes
	Xerox.B235.ESCL.ScannerCapabilities =
		kyoceraECOSYSM2040dnScannerCapabilities
}

//go:embed "data/Xerox-B235-Printer-Attributes.ipp"
var xeroxB235PrinterAttributes []byte

//go:embed "data/Xerox-B235-ScannerCapabilities.xml"
var xeroxB235ScannerCapabilities []byte
