// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ColorProcessing element

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ColorProcessing represents the optional <wscn:ColorProcessing> element
// that specifies the color-processing mode of the input source on the scanner.
//
// Standard values are: "BlackAndWhite1", "Grayscale4", "Grayscale8", "Grayscale16",
// "RGB24", "RGB48", "RGBa32", "RGBa64" (see ColorEntry for details).
// Vendor-defined values are also allowed and will decode as UnknownColorEntry.
//
// It includes optional wscn:MustHonor, wscn:Override, and wscn:UsedDefault
// attributes (all xs:string, but should be boolean values: 0, false, 1, or true).
type ColorProcessing = AttributedElement[ColorEntry]

// decodeColorProcessing decodes [ColorProcessing] from the XML tree.
func decodeColorProcessing(root xmldoc.Element) (ColorProcessing, error) {
	return decodeAttributedElement(root, func(s string) (ColorEntry, error) {
		return DecodeColorEntry(s), nil
	})
}

// toXMLColorProcessing generates XML tree for the [ColorProcessing].
func toXMLColorProcessing(cp ColorProcessing, name string) xmldoc.Element {
	return cp.toXML(name, ColorEntry.String)
}
