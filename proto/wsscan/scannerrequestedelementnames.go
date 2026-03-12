// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Name for RequestedElements element

package wsscan

import "github.com/OpenPrinting/go-mfp/util/xmldoc"

// ScannerRequestedElementNames identifies a section of the WSD Scan Service
// schema that a client wants data for in a GetScannerElementsRequest.
//
// For GetScannerElementsRequest, one of the following QName values:
//   - wscn:ScannerDescription
//   - wscn:ScannerConfiguration
//   - wscn:ScannerStatus
//   - xmlns:VendorSection (vendor-defined extension)
type ScannerRequestedElementNames int

// Known ScannerRequestedElementNames values.
const (
	UnknownScannerRequestedElementNames ScannerRequestedElementNames = iota
	ScannerRequestedElementDescription                               // wscn:ScannerDescription
	ScannerRequestedElementConfiguration                             // wscn:ScannerConfiguration
	ScannerRequestedElementStatus                                    // wscn:ScannerStatus
	ScannerRequestedElementVendorSection                             // xmlns:VendorSection
)

// decodeScannerRequestedElementNames decodes [ScannerRequestedElementNames]
// from the XML tree.
func decodeScannerRequestedElementNames(root xmldoc.Element) (
	ScannerRequestedElementNames, error) {
	return decodeEnum(root, DecodeScannerRequestedElementNames)
}

// toXML generates XML tree for the [ScannerRequestedElementNames].
func (sren ScannerRequestedElementNames) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: sren.String(),
	}
}

// String returns the string representation of the
// [ScannerRequestedElementNames].
func (sren ScannerRequestedElementNames) String() string {
	switch sren {
	case ScannerRequestedElementDescription:
		return NsWSCN + ":ScannerDescription"
	case ScannerRequestedElementConfiguration:
		return NsWSCN + ":ScannerConfiguration"
	case ScannerRequestedElementStatus:
		return NsWSCN + ":ScannerStatus"
	case ScannerRequestedElementVendorSection:
		return NsXML + ":VendorSection"
	}

	return "Unknown"
}

// DecodeScannerRequestedElementNames decodes [ScannerRequestedElementNames]
// out of its XML string representation.
func DecodeScannerRequestedElementNames(s string) ScannerRequestedElementNames {
	switch s {
	case NsWSCN + ":ScannerDescription":
		return ScannerRequestedElementDescription
	case NsWSCN + ":ScannerConfiguration":
		return ScannerRequestedElementConfiguration
	case NsWSCN + ":ScannerStatus":
		return ScannerRequestedElementStatus
	case NsXML + ":VendorSection":
		return ScannerRequestedElementVendorSection
	}

	return UnknownScannerRequestedElementNames
}
