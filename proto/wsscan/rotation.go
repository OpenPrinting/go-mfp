// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Rotation element

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Rotation represents the optional <wscn:Rotation> element
// that specifies the amount to rotate each image of the scanned document.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (xs:string, but should be boolean values: 0, false, 1, or true).
//
// The element contains a required text value that must be one of: 0, 90, 180, or 270.
type Rotation = AttributedElement[RotationValue]

// decodeRotation decodes [Rotation] from the XML tree.
func decodeRotation(root xmldoc.Element) (Rotation, error) {
	return decodeAttributedElement(root, func(s string) (RotationValue, error) {
		val := DecodeRotationValue(s)
		if val == UnknownRotationValue {
			return val, xmldoc.XMLErrWrap(root,
				fmt.Errorf("rotation value must be one of 0, 90, 180, or 270, got %q", s))
		}
		return val, nil
	})
}

// toXMLRotation generates XML tree for the [Rotation].
func toXMLRotation(r Rotation, name string) xmldoc.Element {
	return r.toXML(name, RotationValue.String)
}
