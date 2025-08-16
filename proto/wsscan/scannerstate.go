// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scanner state

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ScannerState defines the current state of the scanning portion of the scan device.
type ScannerState int

// known scanner states:
const (
	UnknownScannerState ScannerState = iota
	Idle
	Processing
	Stopped
)

// decodeScannerState decodes [ScannerState] from the XML tree.
func decodeScannerState(root xmldoc.Element) (ss ScannerState, err error) {
	return decodeEnum(root, DecodeScannerState)
}

// toXML generates XML tree for the [ScannerState].
func (ss ScannerState) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: ss.String(),
	}
}

// String returns a string representation of the [ScannerState]
func (ss ScannerState) String() string {
	switch ss {
	case Idle:
		return "Idle"
	case Processing:
		return "Processing"
	case Stopped:
		return "Stopped"
	}

	return "Unknown"
}

// DecodeScannerState decodes [ScannerState] out of its XML string representation.
func DecodeScannerState(s string) ScannerState {
	switch s {
	case "Idle":
		return Idle
	case "Processing":
		return Processing
	case "Stopped":
		return Stopped
	}

	return UnknownScannerState
}
