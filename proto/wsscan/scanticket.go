// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ScanTicket: defines all description and processing parameters of a scan job

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ScanTicket defines all of the description and processing parameters
// of the currently identified scan job.
type ScanTicket struct {
	DocumentParameters optional.Val[DocumentParameters]
	JobDescription     JobDescription
}

// toXML generates XML tree for the ScanTicket.
func (st ScanTicket) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{}

	// DocumentParameters is optional
	if st.DocumentParameters != nil {
		children = append(children, optional.Get(
			st.DocumentParameters,
		).toXML(NsWSCN+":DocumentParameters"))
	}

	// JobDescription is required
	children = append(children, st.JobDescription.toXML(
		NsWSCN+":JobDescription"))

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeScanTicket decodes ScanTicket from the XML tree.
func decodeScanTicket(root xmldoc.Element) (st ScanTicket, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	documentParameters := xmldoc.Lookup{
		Name:     NsWSCN + ":DocumentParameters",
		Required: false,
	}
	jobDescription := xmldoc.Lookup{
		Name:     NsWSCN + ":JobDescription",
		Required: true,
	}

	missed := root.Lookup(
		&documentParameters,
		&jobDescription,
	)
	if missed != nil && missed.Required {
		return st, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode DocumentParameters if present
	if documentParameters.Found {
		var dp DocumentParameters
		if dp, err = decodeDocumentParameters(
			documentParameters.Elem,
		); err != nil {
			return st, fmt.Errorf("DocumentParameters: %w", err)
		}
		st.DocumentParameters = optional.New(dp)
	}

	// Decode JobDescription (required)
	if st.JobDescription, err = decodeJobDescription(
		jobDescription.Elem,
	); err != nil {
		return st, fmt.Errorf("JobDescription: %w", err)
	}

	return st, nil
}
