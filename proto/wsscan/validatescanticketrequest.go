// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ValidateScanTicketRequest: validates scan ticket settings

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ValidateScanTicketRequest enables a client to determine if the settings
// for future scan operations are valid.
type ValidateScanTicketRequest struct {
	ScanTicket ScanTicket
}

// toXML generates XML tree for the ValidateScanTicketRequest.
func (vstr ValidateScanTicketRequest) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			vstr.ScanTicket.toXML(NsWSCN + ":ScanTicket"),
		},
	}
}

// decodeValidateScanTicketRequest decodes ValidateScanTicketRequest from the XML tree.
func decodeValidateScanTicketRequest(root xmldoc.Element) (
	vstr ValidateScanTicketRequest,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	scanTicket := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanTicket",
		Required: true,
	}

	missed := root.Lookup(&scanTicket)
	if missed != nil && missed.Required {
		return vstr, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode ScanTicket (required)
	if vstr.ScanTicket, err = decodeScanTicket(
		scanTicket.Elem,
	); err != nil {
		return vstr, fmt.Errorf("ScanTicket: %w", err)
	}

	return vstr, nil
}
