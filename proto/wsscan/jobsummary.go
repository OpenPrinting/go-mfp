// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// JobSummary: summary about a scan job

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// JobSummary contains a summary about a scan job.
type JobSummary struct {
	JobID                  int
	JobName                string
	JobOriginatingUserName string
	JobState               JobState
	JobStateReasons        []JobStateReason
	ScansCompleted         int
}

// toXML generates XML tree for the JobSummary.
func (js JobSummary) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		{
			Name: NsWSCN + ":JobId",
			Text: strconv.Itoa(js.JobID),
		},
		{
			Name: NsWSCN + ":JobName",
			Text: js.JobName,
		},
		{
			Name: NsWSCN + ":JobOriginatingUserName",
			Text: js.JobOriginatingUserName,
		},
		js.JobState.toXML(NsWSCN + ":JobState"),
	}

	if len(js.JobStateReasons) > 0 {
		jsrChildren := make([]xmldoc.Element, len(js.JobStateReasons))
		for i, v := range js.JobStateReasons {
			jsrChildren[i] = v.toXML(NsWSCN + ":JobStateReason")
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":JobStateReasons",
			Children: jsrChildren,
		})
	}

	children = append(children, xmldoc.Element{
		Name: NsWSCN + ":ScansCompleted",
		Text: strconv.Itoa(js.ScansCompleted),
	})

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeJobSummary decodes JobSummary from the XML tree.
func decodeJobSummary(root xmldoc.Element) (
	js JobSummary,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	jobID := xmldoc.Lookup{
		Name:     NsWSCN + ":JobId",
		Required: true,
	}
	jobName := xmldoc.Lookup{
		Name:     NsWSCN + ":JobName",
		Required: true,
	}
	jobOriginatingUserName := xmldoc.Lookup{
		Name:     NsWSCN + ":JobOriginatingUserName",
		Required: true,
	}
	jobState := xmldoc.Lookup{
		Name:     NsWSCN + ":JobState",
		Required: true,
	}
	scansCompleted := xmldoc.Lookup{
		Name:     NsWSCN + ":ScansCompleted",
		Required: true,
	}
	jobStateReasons := xmldoc.Lookup{
		Name:     NsWSCN + ":JobStateReasons",
		Required: false,
	}

	missed := root.Lookup(
		&jobID,
		&jobName,
		&jobOriginatingUserName,
		&jobState,
		&scansCompleted,
		&jobStateReasons,
	)
	if missed != nil && missed.Required {
		return js, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode JobID
	if js.JobID, err = decodeNonNegativeInt(jobID.Elem); err != nil {
		return js, fmt.Errorf("JobId: %w", err)
	}

	// Decode JobName
	js.JobName = jobName.Elem.Text
	if js.JobName == "" {
		return js, fmt.Errorf("JobName: empty value")
	}

	// Decode JobOriginatingUserName
	js.JobOriginatingUserName = jobOriginatingUserName.Elem.Text
	if js.JobOriginatingUserName == "" {
		return js, fmt.Errorf("JobOriginatingUserName: empty value")
	}

	// Decode JobState
	if js.JobState, err = decodeJobState(jobState.Elem); err != nil {
		return js, fmt.Errorf("JobState: %w", err)
	}

	// Decode ScansCompleted
	if js.ScansCompleted, err = decodeNonNegativeInt(
		scansCompleted.Elem); err != nil {
		return js, fmt.Errorf("ScansCompleted: %w", err)
	}

	// Decode JobStateReasons (optional)
	if jobStateReasons.Found && jobStateReasons.Elem.Children != nil {
		for _, child := range jobStateReasons.Elem.Children {
			if child.Name == NsWSCN+":JobStateReason" {
				val, err := decodeJobStateReason(child)
				if err != nil {
					return js, fmt.Errorf(
						"JobStateReasons: invalid JobStateReason: %w", err)
				}
				js.JobStateReasons = append(js.JobStateReasons, val)
			}
		}
	}

	return js, nil
}
