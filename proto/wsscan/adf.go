// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ADF: describes the capabilities of the automatic document feeder (ADF)

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ADF describes the capabilities of the automatic document feeder (ADF)
// attached to the scanner. It is a optional element.
type ADF struct {
	ADFBack           optional.Val[ADFFeederSide]
	ADFFront          optional.Val[ADFFeederSide]
	ADFSupportsDuplex BooleanElement
}

// toXML creates an XML element for ADF.
func (a ADF) toXML(name string) xmldoc.Element {
	elm := xmldoc.Element{Name: name}
	if a.ADFBack != nil {
		elm.Children = append(elm.Children,
			optional.Get(a.ADFBack).toXML(NsWSCN+":ADFBack"))
	}
	if a.ADFFront != nil {
		elm.Children = append(elm.Children,
			optional.Get(a.ADFFront).toXML(NsWSCN+":ADFFront"))
	}

	elm.Children = append(elm.Children,
		a.ADFSupportsDuplex.toXML(NsWSCN+":ADFSupportsDuplex"))

	return elm
}

// decodeADF decodes an ADF from an XML element.
func decodeADF(root xmldoc.Element) (ADF, error) {
	var a ADF
	adfBack := xmldoc.Lookup{
		Name:     NsWSCN + ":ADFBack",
		Required: false,
	}
	adfFront := xmldoc.Lookup{
		Name:     NsWSCN + ":ADFFront",
		Required: false,
	}
	adfSupportsDuplex := xmldoc.Lookup{
		Name:     NsWSCN + ":ADFSupportsDuplex",
		Required: true,
	}

	missed := root.Lookup(&adfBack, &adfFront, &adfSupportsDuplex)
	if missed != nil {
		return a, xmldoc.XMLErrMissed(missed.Name)
	}

	if adfBack.Found {
		back, err := decodeADFFeederSide(adfBack.Elem)
		if err != nil {
			return a, fmt.Errorf("ADFBack: %w", err)
		}
		a.ADFBack = optional.New(back)
	}
	if adfFront.Found {
		front, err := decodeADFFeederSide(adfFront.Elem)
		if err != nil {
			return a, fmt.Errorf("ADFFront: %w", err)
		}
		a.ADFFront = optional.New(front)
	}

	duplex, err := decodeBooleanElement(adfSupportsDuplex.Elem)
	if err != nil {
		return a, fmt.Errorf("ADFSupportsDuplex: %w", err)
	}
	a.ADFSupportsDuplex = duplex

	return a, nil
}
