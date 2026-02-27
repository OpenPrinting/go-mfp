// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Job: all elements associated with a scan job

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Job contains all elements that are associated with a scan job.
type Job struct {
	Documents  Documents
	JobStatus  JobStatus
	ScanTicket ScanTicket
}

// toXML generates XML tree for the Job.
func (j Job) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			j.Documents.toXML(NsWSCN + ":Documents"),
			j.JobStatus.toXML(NsWSCN + ":JobStatus"),
			j.ScanTicket.toXML(NsWSCN + ":ScanTicket"),
		},
	}
}

// decodeJob decodes Job from the XML tree.
func decodeJob(root xmldoc.Element) (
	j Job,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	documents := xmldoc.Lookup{
		Name:     NsWSCN + ":Documents",
		Required: true,
	}
	jobStatus := xmldoc.Lookup{
		Name:     NsWSCN + ":JobStatus",
		Required: true,
	}
	scanTicket := xmldoc.Lookup{
		Name:     NsWSCN + ":ScanTicket",
		Required: true,
	}

	missed := root.Lookup(
		&documents,
		&jobStatus,
		&scanTicket,
	)
	if missed != nil && missed.Required {
		return j, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode Documents
	if j.Documents, err = decodeDocuments(
		documents.Elem,
	); err != nil {
		return j, fmt.Errorf("Documents: %w", err)
	}

	// Decode JobStatus
	if j.JobStatus, err = decodeJobStatus(
		jobStatus.Elem,
	); err != nil {
		return j, fmt.Errorf("JobStatus: %w", err)
	}

	// Decode ScanTicket
	if j.ScanTicket, err = decodeScanTicket(
		scanTicket.Elem,
	); err != nil {
		return j, fmt.Errorf("ScanTicket: %w", err)
	}

	return j, nil
}
