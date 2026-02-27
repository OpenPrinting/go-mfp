// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// DocumentDescription: description attributes for document creation

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// DocumentDescription defines all of the description attributes that pertain
// to the basic creation information of the
// currently identified Document element.
type DocumentDescription struct {
	DocumentName string
}

// toXML generates XML tree for the DocumentDescription.
func (dd DocumentDescription) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			{
				Name: NsWSCN + ":DocumentName",
				Text: dd.DocumentName,
			},
		},
	}
}

// decodeDocumentDescription decodes DocumentDescription from the XML tree.
func decodeDocumentDescription(root xmldoc.Element) (
	dd DocumentDescription,
	err error,
) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	documentName := xmldoc.Lookup{
		Name:     NsWSCN + ":DocumentName",
		Required: true,
	}

	missed := root.Lookup(&documentName)
	if missed != nil && missed.Required {
		return dd, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode DocumentName (required)
	dd.DocumentName = documentName.Elem.Text
	if dd.DocumentName == "" {
		return dd, fmt.Errorf("DocumentName: empty value")
	}

	return dd, nil
}
