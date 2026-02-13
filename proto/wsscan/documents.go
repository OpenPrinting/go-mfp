// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Documents: scan characteristics and document collection

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Documents contains the actual scan characteristics that are used during
// image acquisition, plus a collection of all Document elements that the
// scan job contains.
type Documents struct {
	DocumentFinalParameters DocumentParameters
	Document                []Document
}

// toXML generates XML tree for the Documents.
func (ds Documents) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		ds.DocumentFinalParameters.toXML(NsWSCN + ":DocumentFinalParameters"),
	}

	// Add all Document elements
	for _, doc := range ds.Document {
		children = append(children, doc.toXML(NsWSCN+":Document"))
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeDocuments decodes Documents from the XML tree.
func decodeDocuments(root xmldoc.Element) (
	ds Documents,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	documentFinalParameters := xmldoc.Lookup{
		Name:     NsWSCN + ":DocumentFinalParameters",
		Required: true,
	}

	missed := root.Lookup(&documentFinalParameters)
	if missed != nil && missed.Required {
		return ds, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode DocumentFinalParameters
	if ds.DocumentFinalParameters, err = decodeDocumentParameters(
		documentFinalParameters.Elem,
	); err != nil {
		return ds, fmt.Errorf("DocumentFinalParameters: %w", err)
	}

	// Decode all Document elements (optional, can be zero or more)
	for _, child := range root.Children {
		if child.Name == NsWSCN+":Document" {
			var doc Document
			if doc, err = decodeDocument(child); err != nil {
				return ds, fmt.Errorf("Document: %w", err)
			}
			ds.Document = append(ds.Document, doc)
		}
	}

	return ds, nil
}
