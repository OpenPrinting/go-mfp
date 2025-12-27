// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ImagesToTransfer element

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ImagesToTransfer represents the optional <wscn:ImagesToTransfer> element
// that specifies the number of images to scan for the current job.
//
// The value must be an integer in the range from 0 through 2147483648.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
type ImagesToTransfer = AttributedElement[int]

// decodeImagesToTransfer decodes [ImagesToTransfer] from the XML tree.
func decodeImagesToTransfer(root xmldoc.Element) (
	ImagesToTransfer, error) {
	return decodeAttributedElement(root, func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %q", s)
		}
		if val < 0 || val > 2147483648 {
			return 0, fmt.Errorf("value out of range [0-2147483648]: %d", val)
		}
		return val, nil
	})
}

// toXMLImagesToTransfer generates XML tree for the [ImagesToTransfer].
func toXMLImagesToTransfer(itt ImagesToTransfer, name string) xmldoc.Element {
	return itt.toXML(name, strconv.Itoa)
}
