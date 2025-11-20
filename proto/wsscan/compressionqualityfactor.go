// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// CompressionQualityFactor element

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// CompressionQualityFactor represents the optional <wscn:CompressionQualityFactor>
// element that specifies an idealized integer amount of image quality,
// on a scale from 0 through 100.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
type CompressionQualityFactor = AttributedElement[int]

// decodeCompressionQualityFactor decodes [CompressionQualityFactor] from the XML tree.
func decodeCompressionQualityFactor(root xmldoc.Element) (
	CompressionQualityFactor, error) {
	return decodeAttributedElement(root, func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %q", s)
		}
		if val < 0 || val > 100 {
			return 0, fmt.Errorf("value out of range [0-100]: %d", val)
		}
		return val, nil
	})
}

// toXMLCompressionQualityFactor generates XML tree for the [CompressionQualityFactor].
func toXMLCompressionQualityFactor(cqf CompressionQualityFactor, name string) xmldoc.Element {
	return cqf.toXML(name, strconv.Itoa)
}
