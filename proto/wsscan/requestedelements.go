// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// RequestedElements: identifies elements the client wants data for

package wsscan

import (
	"errors"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// RequestedElements holds one or more ScannerRequestedElementNames QNames
// to request from the WSD Scan Service.
type RequestedElements struct {
	Names []ScannerRequestedElementNames // At least one required
}

// toXML generates XML tree for the RequestedElements.
func (re RequestedElements) toXML(name string) xmldoc.Element {
	nameElems := make([]xmldoc.Element, len(re.Names))
	for i, n := range re.Names {
		nameElems[i] = n.toXML(NsWSCN + ":Name")
	}
	return xmldoc.Element{Name: name, Children: nameElems}
}

// decodeRequestedElements decodes RequestedElements from the XML tree.
func decodeRequestedElements(root xmldoc.Element) (
	re RequestedElements, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	for _, child := range root.Children {
		if child.Name == NsWSCN+":Name" {
			sren, decErr := decodeScannerRequestedElementNames(child)
			if decErr != nil {
				return re, decErr
			}
			re.Names = append(re.Names, sren)
		}
	}

	if len(re.Names) == 0 {
		return re, errors.New("at least one Name is required")
	}

	return re, nil
}
