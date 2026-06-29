// MFP - Miulti-Function Printers and scanners toolkit
// Utility functions and data BLOBs for testing
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Test data examples for Kyocera printers

package testutils

import (
	// Import "embed" for its side effects
	_ "embed"
)

// Kyocera contains data samples taken from the Kyocera printers
var Kyocera struct {
	// ECOSYS series
	ECOSYS struct {
		// M2040dn model
		M2040dn struct {
			// IPP protocol samples
			IPP struct {
				PrinterAttributes []byte
			}
			// ESCL protocol samples
			ESCL struct {
				ScannerCapabilities []byte
				ScannerStatus       []byte
			}
			// WS-Scan protocol samples
			WSD struct {
				GetScannerElementsResponse []byte
			}
		}
	}
}

func init() {
	Kyocera.ECOSYS.M2040dn.IPP.PrinterAttributes =
		kyoceraECOSYSM2040dnPrinterAttributes
	Kyocera.ECOSYS.M2040dn.ESCL.ScannerCapabilities =
		kyoceraECOSYSM2040dnScannerCapabilities
	Kyocera.ECOSYS.M2040dn.ESCL.ScannerStatus =
		kyoceraECOSYSM2040dnScannerStatus
	Kyocera.ECOSYS.M2040dn.WSD.GetScannerElementsResponse =
		kyoceraECOSYSM2040dnWSDGetScannerElementsResponse

}

//go:embed "data/Kyocera-ECOSYS-M2040dn-Printer-Attributes.ipp"
var kyoceraECOSYSM2040dnPrinterAttributes []byte

//go:embed "data/Kyocera-ECOSYS-M2040dn-ScannerCapabilities.xml"
var kyoceraECOSYSM2040dnScannerCapabilities []byte

//go:embed "data/Kyocera-ECOSYS-M2040dn-ScannerStatus.xml"
var kyoceraECOSYSM2040dnScannerStatus []byte

//go:embed "data/Kyocera-ECOSYS-M2040dn-WSD-GetScannerElementsResponse.xml"
var kyoceraECOSYSM2040dnWSDGetScannerElementsResponse []byte
