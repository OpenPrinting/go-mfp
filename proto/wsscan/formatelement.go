// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Format element (not to be confused with FormatValue enum type)

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// FormatElement represents the optional <wscn:Format> element
// that indicates a single file format and compression type supported by the scanner.
//
// Standard values include: "dib", "exif", "jbig", "jfif", "jpeg2k", "pdf-a", "png",
// "tiff-single-uncompressed", "tiff-single-g4", "tiff-single-g3mh",
// "tiff-single-jpeg-tn2", "tiff-multi-uncompressed", "tiff-multi-g4",
// "tiff-multi-g3mh", "tiff-multi-jpeg-tn2", "xps".
// Vendor-defined values are also allowed and will decode as UnknownFormatValue.
//
// It includes optional wscn:Override and wscn:UsedDefault attributes
// (all xs:string, but should be boolean values: 0, false, 1, or true).
// Note: This element does NOT have a MustHonor attribute.
type FormatElement = AttributedElement[FormatValue]

// decodeFormatElement decodes [FormatElement] from the XML tree.
func decodeFormatElement(root xmldoc.Element) (FormatElement, error) {
	return decodeAttributedElement(root, func(s string) (FormatValue, error) {
		return DecodeFormatValue(s), nil
	})
}

// toXMLFormatElement generates XML tree for the [FormatElement].
func toXMLFormatElement(f FormatElement, name string) xmldoc.Element {
	return f.toXML(name, FormatValue.String)
}
