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

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// CreateScanJobRequest prepares a scan device to scan.
// ScanTicket is required (host-initiated scan). DestinationToken and
// ScanIdentifier are optional and used for device-initiated scan (user
// pushes button on device); that mode is not fully supported.
type CreateScanJobRequest struct {
	DestinationToken optional.Val[string]
	ScanIdentifier optional.Val[string]
	ScanTicket      ScanTicket
}

// toXML generates XML tree for the CreateScanJobRequest.
func (csjr CreateScanJobRequest) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{}
	if csjr.DestinationToken != nil {
		children = append(children, xmldoc.Element{
			Name: NsWSCN + ":DestinationToken",
			Text: optional.Get(csjr.DestinationToken),
		})
	}
	if csjr.ScanIdentifier != nil {
		children = append(children, xmldoc.Element{
			Name: NsWSCN + ":ScanIdentifier",
			Text: optional.Get(csjr.ScanIdentifier),
		})
	}
	children = append(children, csjr.ScanTicket.toXML(NsWSCN+":ScanTicket"))
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

	scanTicket := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanTicket",
		Required: true,
	}
	destinationToken := xmldoc.Lookup{
		Name:     NsWSCN + ":DestinationToken",
		Required: false,
	}
	scanIdentifier := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanIdentifier",
		Required: false,
	}

	if missed := root.Lookup(
		&scanTicket, 
		&destinationToken,
		&scanIdentifier,
	); missed != nil {
		return csjr, xmldoc.XMLErrMissed(missed.Name)
	}
	
	if csjr.ScanTicket, err = decodeScanTicket(scanTicket.Elem); err != nil {
		return csjr, fmt.Errorf("ScanTicket: %w", err)
	}
	if destinationToken.Found {
		csjr.DestinationToken = optional.New(destinationToken.Elem.Text)
	}
	if scanIdentifier.Found {
		csjr.ScanIdentifier = optional.New(scanIdentifier.Elem.Text)
	}
	return csjr, nil
}
