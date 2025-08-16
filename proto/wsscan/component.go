// MFP - Miulti-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Component element

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Component identifies the device component that a DeviceCondition or
// ConditionHistoryEntry describes.
//
// Values are defined by WS-Scan spec: ADF, Film, MediaPath, Platen.
type Component int

// Known component values
const (
	UnknownComponent Component = iota
	ADFComponent
	FilmComponent
	MediaPathComponent
	PlatenComponent
)

// decodeComponent decodes [Component] from the XML tree.
func decodeComponent(root xmldoc.Element) (c Component, err error) {
	return decodeEnum(root, DecodeComponent)
}

// toXML generates XML tree for the [Component].
func (c Component) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: c.String(),
	}
}

// String returns a string representation of the [Component].
func (c Component) String() string {
	switch c {
	case ADFComponent:
		return "ADF"
	case FilmComponent:
		return "Film"
	case MediaPathComponent:
		return "MediaPath"
	case PlatenComponent:
		return "Platen"
	}

	return "Unknown"
}

// DecodeComponent decodes [Component] out of its XML string representation.
func DecodeComponent(s string) Component {
	switch s {
	case "ADF":
		return ADFComponent
	case "Film":
		return FilmComponent
	case "MediaPath":
		return MediaPathComponent
	case "Platen":
		return PlatenComponent
	}

	return UnknownComponent
}
