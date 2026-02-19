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

// CreateScanJobRequest prepares a scan device to scan (host-initiated mode).
// Only host-initiated scan is supported: the client starts the application and
// acquires an image. This mode uses only the ScanTicket parameter.
// Device-initiated scan (user pushes button on device; requires ScanIdentifier
// and DestinationToken) is not supported.
type CreateScanJobRequest struct {
	ScanTicket ScanTicket
}

// toXML generates XML tree for the CreateScanJobRequest.
func (csjr CreateScanJobRequest) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name:     name,
		Children: []xmldoc.Element{csjr.ScanTicket.toXML(NsWSCN + ":ScanTicket")},
	}
}

// decodeCreateScanJobRequest decodes CreateScanJobRequest from the XML tree.
func decodeCreateScanJobRequest(root xmldoc.Element) (
	csjr CreateScanJobRequest,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	scanTicket := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanTicket",
		Required: true,
	}

	if missed := root.Lookup(&scanTicket); missed != nil && missed.Required {
		return csjr, xmldoc.XMLErrMissed(missed.Name)
	}

	if csjr.ScanTicket, err = decodeScanTicket(scanTicket.Elem); err != nil {
		return csjr, fmt.Errorf("ScanTicket: %w", err)
	}

	return csjr, nil
}
