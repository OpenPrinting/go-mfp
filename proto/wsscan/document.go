// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Document: details about documents scanned during a scan job

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Document contains the details about one of the documents that are scanned
// during a scan job.
type Document struct {
	DocumentDescription DocumentDescription
}

// toXML generates XML tree for the Document.
func (d Document) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			d.DocumentDescription.toXML(NsWSCN + ":DocumentDescription"),
		},
	}
}

// decodeDocument decodes Document from the XML tree.
func decodeDocument(root xmldoc.Element) (
	d Document,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	documentDescription := xmldoc.Lookup{
		Name:     NsWSCN + ":DocumentDescription",
		Required: true,
	}

	missed := root.Lookup(&documentDescription)
	if missed != nil && missed.Required {
		return d, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode DocumentDescription
	if d.DocumentDescription, err = decodeDocumentDescription(
		documentDescription.Elem,
	); err != nil {
		return d, fmt.Errorf("DocumentDescription: %w", err)
	}

	return d, nil
}
