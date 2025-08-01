// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Dimension: simple struct for width and height values

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Dimension represents a simple width and height pair.
type Dimension struct {
	Width  int
	Height int
}

// toXML creates an XML element for Dimension.
func (d Dimension) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":Width",
				Text: strconv.Itoa(d.Width),
			},
			{
				Name: NsWSCN + ":Height",
				Text: strconv.Itoa(d.Height),
			},
		},
	}
}

// decodeDimension decodes a Dimension from an XML element.
func decodeDimension(root xmldoc.Element) (Dimension, error) {
	var dim Dimension
	var widthFound, heightFound bool

	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Width":
			width, err := strconv.Atoi(child.Text)
			if err != nil {
				return dim, fmt.Errorf("invalid width value: %w", err)
			}
			dim.Width = width
			widthFound = true
		case NsWSCN + ":Height":
			height, err := strconv.Atoi(child.Text)
			if err != nil {
				return dim, fmt.Errorf("invalid height value: %w", err)
			}
			dim.Height = height
			heightFound = true
		}
	}

	if !widthFound || !heightFound {
		return dim, fmt.Errorf("missing Width or Height element")
	}

	return dim, nil
}
