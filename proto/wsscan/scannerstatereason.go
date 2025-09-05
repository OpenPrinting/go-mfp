// MFP - Miulti-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scan state reason

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ScannerStateReason defines the reason why the scanner is in its current state.
type ScannerStateReason int

// known scanner state reasons:
const (
	UnknownScannerStateReason ScannerStateReason = iota
	StateAttentionRequired
	StateCalibrating
	StateCoverOpen
	StateInterlockOpen
	StateInternalStorageFull
	StateLampError
	StateLampWarming
	StateMediaJam
	StateMultipleFeedError
	StateNone
	StatePaused
)

// decodeScannerStateReason decodes [ScannerStateReason] from the XML tree.
func decodeScannerStateReason(root xmldoc.Element) (ssr ScannerStateReason, err error) {
	return decodeEnum(root, DecodeScannerStateReason)
}

// toXML generates XML tree for the [ScannerStateReason].
func (ssr ScannerStateReason) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: ssr.String(),
	}
}

// String returns a string representation of the [ScannerStateReason]
func (ssr ScannerStateReason) String() string {
	switch ssr {
	case StateAttentionRequired:
		return "AttentionRequired"
	case StateCalibrating:
		return "Calibrating"
	case StateCoverOpen:
		return "CoverOpen"
	case StateInterlockOpen:
		return "InterlockOpen"
	case StateInternalStorageFull:
		return "InternalStorageFull"
	case StateLampError:
		return "LampError"
	case StateLampWarming:
		return "LampWarming"
	case StateMediaJam:
		return "MediaJam"
	case StateMultipleFeedError:
		return "MultipleFeedError"
	case StateNone:
		return "None"
	case StatePaused:
		return "Paused"
	}

	return "Unknown"
}

// DecodeScannerStateReason decodes [ScannerStateReason] out of its XML string representation.
func DecodeScannerStateReason(s string) ScannerStateReason {
	switch s {
	case "AttentionRequired":
		return StateAttentionRequired
	case "Calibrating":
		return StateCalibrating
	case "CoverOpen":
		return StateCoverOpen
	case "InterlockOpen":
		return StateInterlockOpen
	case "InternalStorageFull":
		return StateInternalStorageFull
	case "LampError":
		return StateLampError
	case "LampWarming":
		return StateLampWarming
	case "MediaJam":
		return StateMediaJam
	case "MultipleFeedError":
		return StateMultipleFeedError
	case "None":
		return StateNone
	case "Paused":
		return StatePaused
	}

	return UnknownScannerStateReason
}
