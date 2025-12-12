// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ContentType element

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ContentType represents the optional <wscn:ContentType> element
// that specifies the document content type.
//
// Standard values are: "Auto", "Text", "Photo", "Halftone", "Mixed".
// Values can be extended or subset, so any string value is accepted.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
type ContentType = AttributedElement[string]

// decodeContentType decodes [ContentType] from the XML tree.
func decodeContentType(root xmldoc.Element) (ContentType, error) {
	return decodeAttributedElement(root, func(s string) (string, error) {
		// Accept any string value as values can be extended/subset
		return s, nil
	})
}

// toXMLContentType generates XML tree for the [ContentType].
func toXMLContentType(ct ContentType, name string) xmldoc.Element {
	return ct.toXML(name, func(s string) string {
		return s
	})
}
