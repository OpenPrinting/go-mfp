// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Format element: indicates a single file format and compression type

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Format indicates a single file format and compression type
// supported by the scanner.
// The text value can be standard values or vendor-defined values.
type Format ValWithOptions[FormatValue]

// decodeFormat decodes a Format from an XML element.
func decodeFormat(root xmldoc.Element) (Format, error) {
	var base ValWithOptions[FormatValue]
	decoded, err := base.decodeValWithOptions(root, formatValueDecoder)
	if err != nil {
		return Format{}, err
	}
	return Format(decoded), nil
}

// toXML converts a Format to an XML element.
func (f Format) toXML(name string) xmldoc.Element {
	return ValWithOptions[FormatValue](f).toXML(name, formatValueEncoder)
}

// formatValueDecoder converts a string to a FormatValue.
func formatValueDecoder(s string) (FormatValue, error) {
	return DecodeFormatValue(s), nil
}

// formatValueEncoder converts a FormatValue to a string.
func formatValueEncoder(fv FormatValue) string {
	return string(fv)
}
