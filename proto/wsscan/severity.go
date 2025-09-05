// MFP - Miulti-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Severity element for DeviceCondition and ConditionHistoryEntry

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Severity specifies the severity level of the current condition.
// Allowed values: Informational, Warning, Critical.
type Severity int

// Known Severity values.
const (
	UnknownSeverity Severity = iota
	Informational
	Warning
	Critical
)

// decodeSeverity decodes [Severity] from the XML tree.
func decodeSeverity(root xmldoc.Element) (s Severity, err error) {
	return decodeEnum(root, DecodeSeverity)
}

// toXML generates XML tree for the [Severity].
func (s Severity) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: s.String(),
	}
}

// String returns a string representation of the [Severity].
func (s Severity) String() string {
	switch s {
	case Informational:
		return "Informational"
	case Warning:
		return "Warning"
	case Critical:
		return "Critical"
	}

	return "Unknown"
}

// DecodeSeverity decodes [Severity] out of its XML string representation.
func DecodeSeverity(str string) Severity {
	switch str {
	case "Informational":
		return Informational
	case "Warning":
		return Warning
	case "Critical":
		return Critical
	}

	return UnknownSeverity
}

