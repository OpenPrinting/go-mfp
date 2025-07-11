// MFP - Miulti-Function Printers and scanners toolkit
// WSD core protocol
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device types, for discovery

package wsd

import (
	"strings"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Types represents set of device types, for discovery
type Types []Type

// Type represents a device type.
type Type int

// Known types:
const (
	UnknownType Type = 1 << iota
	Device
	PrinterServiceType
	ScannerServiceType
)

// DecodeTypes decodes [Types] from the XML tree
func DecodeTypes(root xmldoc.Element) (types Types, err error) {
	names := strings.Fields(root.Text)

	for _, n := range names {
		// Note, type names looks as follows: namespace:name
		// (for example, devprof:Device). However, this is very
		// hard to bring here information from the original
		// XMP about namespace prefixes assignments. So as a
		// workaround, we just ignore prefixes here.
		if i := strings.IndexByte(n, ':'); i >= 0 {
			n = n[i+1:]
		}

		switch n {
		case "Device":
			types = append(types, Device)
		case "PrintDeviceType":
			types = append(types, PrinterServiceType)
		case "ScanDeviceType":
			types = append(types, ScannerServiceType)
		}
	}

	return
}

// DecodeMetadataTypes decodes [Types] from the XML tree.
//
// It works like [DecodeTypes] but for types encoded within [Metadata]
// messages.
func DecodeMetadataTypes(root xmldoc.Element) (types Types, err error) {
	names := strings.Fields(root.Text)

	for _, n := range names {
		// Note, type names looks as follows: namespace:name
		// (for example, devprof:Device). However, this is very
		// hard to bring here information from the original
		// XMP about namespace prefixes assignments. So as a
		// workaround, we just ignore prefixes here.
		if i := strings.IndexByte(n, ':'); i >= 0 {
			n = n[i+1:]
		}

		switch n {
		case "PrinterServiceType":
			types = append(types, PrinterServiceType)
		case "ScannerServiceType":
			types = append(types, ScannerServiceType)
		}
	}

	return
}

// Contains reports if type is member of types.
func (types Types) Contains(t Type) bool {
	for _, contained := range types {
		if t == contained {
			return true
		}
	}
	return false
}

// String returns text representation for [Types].
//
// The returned value can be directly used as a text value of Types XML
// element, except for [Metadata] message encoding.
//
// Use for Metadata, you need to use the [Types.MetadataString] function.
func (types Types) String() string {
	names := make([]string, len(types))

	for i := range types {
		names[i] = types[i].String()
	}

	return strings.Join(names, " ")
}

// MetadataString returns the XML text representation for [Types],
// suitable for [Metadata] message encoding.
//
// This is very similar to the [Types.String] but uses slightly
// different spelling of keywords.
func (types Types) MetadataString() string {
	names := make([]string, 0, 3)

	for _, t := range types {
		switch t {
		case PrinterServiceType:
			names = append(names, "print:PrinterServiceType")

		case ScannerServiceType:
			names = append(names, "scan:ScannerServiceType")
		}
	}

	return strings.Join(names, " ")
}

// ToXML generates XML tree for the Types.
//
// For [Metadata] encoding, use [Types.MetadataToXML].
func (types Types) ToXML() xmldoc.Element {
	elm := xmldoc.Element{
		Name: NsDiscovery + ":Types",
		Text: types.String(),
	}

	return elm
}

// MetadataToXML generates XML tree for the [Types].
//
// It is intended for encoding the [Metadata] messages, which use
// slightly different encoding for the Types element.
func (types Types) MetadataToXML() xmldoc.Element {
	elm := xmldoc.Element{
		Name: NsDevprof + ":Types",
		Text: types.MetadataString(),
	}

	return elm
}

// MarkUsedNamespace marks [xmldoc.Namespace] entries used by
// data elements within the message body, if any.
func (types Types) MarkUsedNamespace(ns xmldoc.Namespace) {
	// Note, xmldoc.Namespace may have multiple entries with the
	// same prefix and different URLs. Only the first one should
	// be used for output, while others allow to handle different
	// namespace URLs as equal on input (for example, SOUP 1.1 and
	// 1.2 use different URLs).
	//
	// So it is better to leave Namespace.MarkUsedPrefix to handle
	// all these nuances rather that to duplicate its work, trading
	// simplicity for efficiency.
	for _, t := range types {
		switch t {
		case Device:
			ns.MarkUsedPrefix("devprof")
		case PrinterServiceType:
			ns.MarkUsedPrefix("print")
		case ScannerServiceType:
			ns.MarkUsedPrefix("scan")
		}
	}
}

// String returns text representation for [Type].
func (t Type) String() string {
	switch t {
	case Device:
		return "devprof:Device"
	case PrinterServiceType:
		return "print:PrintDeviceType"
	case ScannerServiceType:
		return "scan:ScanDeviceType"
	}

	return "Unknown"
}
