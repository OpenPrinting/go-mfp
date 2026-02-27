// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ActiveJobs: list of all currently active scan jobs

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ActiveJobs contains a list of all currently active scan jobs.
type ActiveJobs struct {
	Job        []Job
	JobSummary []JobSummary
}

// toXML generates XML tree for the ActiveJobs.
func (aj ActiveJobs) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{}

	// Add all Job elements
	for _, job := range aj.Job {
		children = append(children, job.toXML(NsWSCN+":Job"))
	}

	// Add all JobSummary elements
	for _, js := range aj.JobSummary {
		children = append(children, js.toXML(NsWSCN+":JobSummary"))
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeActiveJobs decodes ActiveJobs from the XML tree.
func decodeActiveJobs(root xmldoc.Element) (
	aj ActiveJobs,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Job":
			var job Job
			if job, err = decodeJob(child); err != nil {
				return aj, err
			}
			aj.Job = append(aj.Job, job)
		case NsWSCN + ":JobSummary":
			var js JobSummary
			if js, err = decodeJobSummary(child); err != nil {
				return aj, err
			}
			aj.JobSummary = append(aj.JobSummary, js)
		}
	}

	return aj, nil
}
