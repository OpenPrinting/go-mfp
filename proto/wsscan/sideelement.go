// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ADFFeederSideElement: reusable element for ADFBack, ADFFront, etc.

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ADFFeederSideElement describes the capabilities of a scanner ADF feeder side (ADFBack, ADFFront, etc.).
type ADFFeederSideElement struct {
	ADFColor             []ColorEntry
	ADFMaximumSize       HeightAndWidthElement
	ADFMinimumSize       HeightAndWidthElement
	ADFOpticalResolution HeightAndWidthElement
	ADFResolutions       HeightAndWidthElement
}

// toXML creates an XML element for ADFFeederSideElement.
func (s ADFFeederSideElement) toXML(name string) xmldoc.Element {
	elm := xmldoc.Element{Name: name}
	if len(s.ADFColor) > 0 {
		colorChildren := make([]xmldoc.Element, len(s.ADFColor))
		for i, color := range s.ADFColor {
			colorChildren[i] = color.toXML(NsWSCN + ":ColorEntry")
		}
		elm.Children = append(elm.Children, xmldoc.Element{
			Name:     NsWSCN + ":ADFColor",
			Children: colorChildren,
		})
	}
	elm.Children = append(elm.Children, s.ADFMaximumSize.toXML(NsWSCN+
		":ADFMaximumSize"))
	elm.Children = append(elm.Children, s.ADFMinimumSize.toXML(NsWSCN+
		":ADFMinimumSize"))
	elm.Children = append(elm.Children, s.ADFOpticalResolution.toXML(NsWSCN+":ADFOpticalResolution"))
	elm.Children = append(elm.Children, s.ADFResolutions.toXML(NsWSCN+
		":ADFResolutions"))
	return elm
}

// decodeADFFeederSideElement decodes an ADFFeederSideElement from an XML element.
func decodeADFFeederSideElement(root xmldoc.Element) (
	ADFFeederSideElement, error) {
	var s ADFFeederSideElement
	adfColor := xmldoc.Lookup{Name: NsWSCN + ":ADFColor"}
	adfMaximumSize := xmldoc.Lookup{Name: NsWSCN + ":ADFMaximumSize"}
	adfMinimumSize := xmldoc.Lookup{Name: NsWSCN + ":ADFMinimumSize"}
	adfOpticalResolution := xmldoc.Lookup{Name: NsWSCN + ":ADFOpticalResolution"}
	adfResolutions := xmldoc.Lookup{Name: NsWSCN + ":ADFResolutions"}

	missed := root.Lookup(
		&adfColor,
		&adfMaximumSize,
		&adfMinimumSize,
		&adfOpticalResolution,
		&adfResolutions,
	)
	if missed != nil {
		return s, xmldoc.XMLErrMissed(missed.Name)
	}

	for _, child := range adfColor.Elem.Children {
		val, err := decodeColorEntry(child)
		if err != nil {
			return s, fmt.Errorf("ADFColor: %w", err)
		}
		s.ADFColor = append(s.ADFColor, val)
	}

	max, err := decodeHeightAndWidthElement(adfMaximumSize.Elem)
	if err != nil {
		return s, fmt.Errorf("ADFMaximumSize: %w", err)
	}
	s.ADFMaximumSize = max

	min, err := decodeHeightAndWidthElement(adfMinimumSize.Elem)
	if err != nil {
		return s, fmt.Errorf("ADFMinimumSize: %w", err)
	}
	s.ADFMinimumSize = min

	opt, err := decodeHeightAndWidthElement(adfOpticalResolution.Elem)
	if err != nil {
		return s, fmt.Errorf("ADFOpticalResolution: %w", err)
	}
	s.ADFOpticalResolution = opt

	res, err := decodeHeightAndWidthElement(adfResolutions.Elem)
	if err != nil {
		return s, fmt.Errorf("ADFResolutions: %w", err)
	}
	s.ADFResolutions = res

	return s, nil
}
