// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// BooleanElement: reusable boolean XML element type

package wsscan

import (
	"errors"
	"strings"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// BooleanElement is a string type representing a boolean XML value.
// Allowed values: "0", "1", "false", "true"
// (case-insensitive, whitespace ignored).
type BooleanElement string

// Validate checks that the value is a valid BooleanElement value.
func (b BooleanElement) Validate() error {
	switch strings.ToLower(strings.TrimSpace(string(b))) {
	case "0", "1", "false", "true":
		return nil
	default:
		return errors.New(
			"BooleanElement: must be 0, 1, false, or true",
		)
	}
}

// Bool returns the boolean value of BooleanElement.
// Returns true for "1" or "true", false for "0" or "false" defaulting to false.
func (b BooleanElement) Bool() bool {
	switch strings.ToLower(strings.TrimSpace(string(b))) {
	case "1", "true":
		return true
	case "0", "false":
		return false
	default:
		return false
	}
}

// decodeBooleanElement decodes a BooleanElement from an XML element.
func decodeBooleanElement(root xmldoc.Element) (BooleanElement, error) {
	val := BooleanElement(strings.TrimSpace(root.Text))
	return val, val.Validate()
}

// toXML converts a BooleanElement to an XML element.
func (b BooleanElement) toXML(name string) xmldoc.Element {
	return xmldoc.Element{Name: name, Text: string(b)}
}
