// MFP - Miulti-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Name element for DeviceCondition and ConditionHistoryEntry

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// NameElement names the current error condition specified in
// DeviceCondition or ConditionHistoryEntry.
//
// Values are defined by the WS-Scan spec.
type NameElement int

// Known NameElement values.
const (
	UnknownNameElement NameElement = iota
	Calibrating
	CoverOpen
	InputTrayEmpty
	InterlockOpen
	InternalStorageFull
	LampError
	LampWarming
	MediaJam
	MultipleFeedError
)

// decodeNameElement decodes [NameElement] from the XML tree.
func decodeNameElement(root xmldoc.Element) (ne NameElement, err error) {
	return decodeEnum(root, DecodeNameElement)
}

// toXML generates XML tree for the [NameElement].
func (ne NameElement) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: ne.String(),
	}
}

// String returns a string representation of the [NameElement].
func (ne NameElement) String() string {
	switch ne {
	case Calibrating:
		return "Calibrating"
	case CoverOpen:
		return "CoverOpen"
	case InputTrayEmpty:
		return "InputTrayEmpty"
	case InterlockOpen:
		return "InterlockOpen"
	case InternalStorageFull:
		return "InternalStorageFull"
	case LampError:
		return "LampError"
	case LampWarming:
		return "LampWarming"
	case MediaJam:
		return "MediaJam"
	case MultipleFeedError:
		return "MultipleFeedError"
	}

	return "Unknown"
}

// DecodeNameElement decodes [NameElement] out of its XML string representation.
func DecodeNameElement(s string) NameElement {
	switch s {
	case "Calibrating":
		return Calibrating
	case "CoverOpen":
		return CoverOpen
	case "InputTrayEmpty":
		return InputTrayEmpty
	case "InterlockOpen":
		return InterlockOpen
	case "InternalStorageFull":
		return InternalStorageFull
	case "LampError":
		return LampError
	case "LampWarming":
		return LampWarming
	case "MediaJam":
		return MediaJam
	case "MultipleFeedError":
		return MultipleFeedError
	}

	return UnknownNameElement
}
