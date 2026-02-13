// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// CreateScanJobRequest: prepares a scan device to scan

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// CreateScanJobRequest prepares a scan device to scan.
// All child elements are required.
type CreateScanJobRequest struct {
	DestinationToken string
	ScanIdentifier   string
	ScanTicket       ScanTicket
}

// toXML generates XML tree for the CreateScanJobRequest.
func (csjr CreateScanJobRequest) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		{
			Name: NsWSCN + ":DestinationToken",
			Text: csjr.DestinationToken,
		},
		{
			Name: NsWSCN + ":ScanIdentifier",
			Text: csjr.ScanIdentifier,
		},
		csjr.ScanTicket.toXML(NsWSCN + ":ScanTicket"),
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeCreateScanJobRequest decodes CreateScanJobRequest from the XML tree.
func decodeCreateScanJobRequest(root xmldoc.Element) (
	csjr CreateScanJobRequest,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	destinationToken := xmldoc.Lookup{
		Name:     NsWSCN + ":DestinationToken",
		Required: true,
	}
	scanIdentifier := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanIdentifier",
		Required: true,
	}
	scanTicket := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanTicket",
		Required: true,
	}

	missed := root.Lookup(
		&destinationToken,
		&scanIdentifier,
		&scanTicket,
	)
	if missed != nil && missed.Required {
		return csjr, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode DestinationToken (required)
	csjr.DestinationToken = destinationToken.Elem.Text

	// Decode ScanIdentifier (required)
	csjr.ScanIdentifier = scanIdentifier.Elem.Text

	// Decode ScanTicket (required)
	if csjr.ScanTicket, err = decodeScanTicket(
		scanTicket.Elem,
	); err != nil {
		return csjr, fmt.Errorf("ScanTicket: %w", err)
	}

	return csjr, nil
}
