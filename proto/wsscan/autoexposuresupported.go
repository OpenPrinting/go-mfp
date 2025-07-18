// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla
// See LICENSE for license terms and conditions
//
// auto exposure supported

package wsscan

import (
	"errors"
	"strings"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// AutoExposureSupported is a string type representing whether the scan device
// supports automatic exposure adjustment. Allowed values: "0", "1", "false",
// "true" (case-insensitive, whitespace ignored).
type AutoExposureSupported string

// IsValid returns true if the value is a valid AutoExposureSupported value.
func (aes AutoExposureSupported) IsValid() bool {
	switch strings.ToLower(strings.TrimSpace(string(aes))) {
	case "0", "1", "false", "true":
		return true
	default:
		return false
	}
}

// Bool returns the boolean value of AutoExposureSupported.
// Returns true for "1" or "true", false for "0" or "false".
func (aes AutoExposureSupported) Bool() (bool, error) {
	switch strings.ToLower(strings.TrimSpace(string(aes))) {
	case "1", "true":
		return true, nil
	case "0", "false":
		return false, nil
	default:
		return false, errors.New(
			"invalid value for AutoExposureSupported: " +
				"must be 0, 1, false, or true",
		)
	}
}

// decodeAutoExposureSupported decodes an AutoExposureSupported from an XML
// element.
func decodeAutoExposureSupported(root xmldoc.Element) (
	AutoExposureSupported, error,
) {
	val := AutoExposureSupported(strings.TrimSpace(root.Text))
	if !val.IsValid() {
		return "", errors.New(
			"invalid value for AutoExposureSupported: " +
				"must be 0, 1, false, or true",
		)
	}
	return val, nil
}

// toXML converts an AutoExposureSupported to an XML element.
func (aes AutoExposureSupported) toXML(name string) xmldoc.Element {
	return xmldoc.Element{Name: name, Text: string(aes)}
}
