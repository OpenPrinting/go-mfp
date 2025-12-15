// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ContentType element

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ContentType represents the optional <wscn:ContentType> element
// that specifies the document content type.
//
// Supported values are defined by [ContentTypeValue].
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
type ContentType = AttributedElement[ContentTypeValue]

// decodeContentType decodes [ContentType] from the XML tree.
func decodeContentType(root xmldoc.Element) (ContentType, error) {
	return decodeAttributedElement(root, func(s string) (ContentTypeValue, error) {
		val := DecodeContentTypeValue(s)
		if val == UnknownContentTypeValue {
			return val, xmldoc.XMLErrWrap(root, fmt.Errorf("invalid ContentTypeValue: %q", s))
		}
		return val, nil
	})
}

// toXMLContentType generates XML tree for the [ContentType].
func toXMLContentType(ct ContentType, name string) xmldoc.Element {
	return ct.toXML(name, func(v ContentTypeValue) string {
		return v.String()
	})
}
