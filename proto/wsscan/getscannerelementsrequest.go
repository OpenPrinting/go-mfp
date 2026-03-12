// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// GetScannerElementsRequest: requests scanner element information

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// GetScannerElementsRequest enables a client to request information
// about the scanner from the WSD Scan Service.
type GetScannerElementsRequest struct {
	RequestedElements RequestedElements
}

// toXML generates XML tree for the GetScannerElementsRequest.
func (gser GetScannerElementsRequest) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			gser.RequestedElements.toXML(NsWSCN + ":RequestedElements"),
		},
	}
}

// decodeGetScannerElementsRequest
// decodes GetScannerElementsRequest from the XML tree.
func decodeGetScannerElementsRequest(root xmldoc.Element) (
	gser GetScannerElementsRequest,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	requestedElements := xmldoc.Lookup{
		Name:     NsWSCN + ":RequestedElements",
		Required: true,
	}

	missed := root.Lookup(&requestedElements)
	if missed != nil && missed.Required {
		return gser, xmldoc.XMLErrMissed(missed.Name)
	}

	if gser.RequestedElements, err = decodeRequestedElements(
		requestedElements.Elem,
	); err != nil {
		return gser, err
	}

	return gser, nil
}
