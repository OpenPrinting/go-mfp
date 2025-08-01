// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scaling range supported (ScalingRangeSupported)

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ScalingRangeSupported represents the <wscn:ScalingRangeSupported> element,
// containing ScalingWidth and ScalingHeight.
type ScalingRangeSupported struct {
	ScalingWidth  Range
	ScalingHeight Range
}

// toXML generates XML tree for the [ScalingRangeSupported].
func (srs ScalingRangeSupported) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			{
				Name:     NsWSCN + ":ScalingWidth",
				Children: srs.ScalingWidth.toXML(),
			},
			{
				Name:     NsWSCN + ":ScalingHeight",
				Children: srs.ScalingHeight.toXML(),
			},
		},
	}
}

// decodeScalingRangeSupported decodes [ScalingRangeSupported] from the XML tree.
func decodeScalingRangeSupported(root xmldoc.Element) (
	srs ScalingRangeSupported, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	widthLookup := xmldoc.Lookup{
		Name:     NsWSCN + ":ScalingWidth",
		Required: true,
	}
	heightLookup := xmldoc.Lookup{
		Name:     NsWSCN + ":ScalingHeight",
		Required: true,
	}

	missed := root.Lookup(&widthLookup, &heightLookup)
	if missed != nil {
		return srs, xmldoc.XMLErrMissed(missed.Name)
	}

	width, err := decodeRange(widthLookup.Elem)
	if err != nil {
		return srs, fmt.Errorf("invalid ScalingWidth: %w", err)
	}
	srs.ScalingWidth = width

	height, err := decodeRange(heightLookup.Elem)
	if err != nil {
		return srs, fmt.Errorf("invalid ScalingHeight: %w", err)
	}
	srs.ScalingHeight = height

	return srs, nil
}
