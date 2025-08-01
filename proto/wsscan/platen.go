// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Platen: describes the capabilities of the flatbed platen

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Platen describes the capabilities of the flatbed platen available on the scanner.
type Platen struct {
	PlatenColor             []ColorEntry
	PlatenMaximumSize       Dimension
	PlatenMinimumSize       Dimension
	PlatenOpticalResolution Dimension
	PlatenResolutions       Resolutions
}

// toXML creates an XML element for Platen.
func (p Platen) toXML(name string) xmldoc.Element {
	elm := xmldoc.Element{Name: name}
	// PlatenColor as parent with ColorEntry children

	colorChildren := make([]xmldoc.Element, len(p.PlatenColor))
	for i, color := range p.PlatenColor {
		colorChildren[i] = color.toXML(NsWSCN + ":ColorEntry")
	}
	elm.Children = append(elm.Children, xmldoc.Element{
		Name:     NsWSCN + ":PlatenColor",
		Children: colorChildren,
	})

	elm.Children = append(elm.Children, p.PlatenMaximumSize.toXML(NsWSCN+
		":PlatenMaximumSize"))
	elm.Children = append(elm.Children, p.PlatenMinimumSize.toXML(NsWSCN+
		":PlatenMinimumSize"))
	elm.Children = append(elm.Children, p.PlatenOpticalResolution.toXML(NsWSCN+":PlatenOpticalResolution"))
	elm.Children = append(elm.Children, p.PlatenResolutions.toXML(NsWSCN+
		":PlatenResolutions"))
	return elm
}

// decodePlaten decodes a Platen from an XML element using the lookup pattern.
func decodePlaten(root xmldoc.Element) (Platen, error) {
	var p Platen
	// Setup lookups for all possible child elements
	platenColor := xmldoc.Lookup{
		Name:     NsWSCN + ":PlatenColor",
		Required: true,
	}
	platenMaximumSize := xmldoc.Lookup{
		Name:     NsWSCN + ":PlatenMaximumSize",
		Required: true,
	}
	platenMinimumSize := xmldoc.Lookup{
		Name:     NsWSCN + ":PlatenMinimumSize",
		Required: true,
	}
	platenOpticalResolution := xmldoc.Lookup{
		Name:     NsWSCN + ":PlatenOpticalResolution",
		Required: true,
	}
	platenResolutions := xmldoc.Lookup{
		Name:     NsWSCN + ":PlatenResolutions",
		Required: true,
	}

	missed := root.Lookup(
		&platenColor,
		&platenMaximumSize,
		&platenMinimumSize,
		&platenOpticalResolution,
		&platenResolutions,
	)
	if missed != nil {
		return p, xmldoc.XMLErrMissed(missed.Name)
	}

	// PlatenColor
	for _, child := range platenColor.Elem.Children {
		val, err := decodeColorEntry(child)
		if err != nil {
			return p, fmt.Errorf("PlatenColor: %w", err)
		}
		p.PlatenColor = append(p.PlatenColor, val)

	}

	// PlatenMaximumSize
	max, err := decodeDimension(platenMaximumSize.Elem)
	if err != nil {
		return p, fmt.Errorf("PlatenMaximumSize: %w", err)
	}
	p.PlatenMaximumSize = max

	// PlatenMinimumSize
	min, err := decodeDimension(platenMinimumSize.Elem)
	if err != nil {
		return p, fmt.Errorf("PlatenMinimumSize: %w", err)
	}
	p.PlatenMinimumSize = min

	// PlatenOpticalResolution
	opt, err := decodeDimension(platenOpticalResolution.Elem)
	if err != nil {
		return p, fmt.Errorf("PlatenOpticalResolution: %w", err)
	}
	p.PlatenOpticalResolution = opt

	// PlatenResolutions
	res, err := decodeResolutions(platenResolutions.Elem)
	if err != nil {
		return p, fmt.Errorf("PlatenResolutions: %w", err)
	}
	p.PlatenResolutions = res

	return p, nil
}
