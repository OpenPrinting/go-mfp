// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// InputSource element

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// InputSource represents the optional <wscn:InputSource> element
// that specifies the source of the original document on the scanning device.
//
// Standard values are: "ADF", "ADFDuplex", "Film", "Platen".
// Vendor-defined values are also allowed and will decode as UnknownInputSourceValue.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
type InputSource = AttributedElement[InputSourceValue]

// decodeInputSource decodes [InputSource] from the XML tree.
func decodeInputSource(root xmldoc.Element) (InputSource, error) {
	return decodeAttributedElement(root, func(s string) (InputSourceValue, error) {
		return DecodeInputSourceValue(s), nil
	})
}

// toXMLInputSource generates XML tree for the [InputSource].
func toXMLInputSource(is InputSource, name string) xmldoc.Element {
	return is.toXML(name, InputSourceValue.String)
}
