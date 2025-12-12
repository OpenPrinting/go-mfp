// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Exposure element

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Exposure represents the optional <wscn:Exposure> element
// that specifies the exposure settings of the document.
//
// It includes an optional wscn:MustHonor attribute (xs:string,
// but should be a boolean value: 0, false, 1, or true).
//
// The element contains child elements that define the exposure settings.
type Exposure struct {
	MustHonor optional.Val[BooleanElement]
	Children  []xmldoc.Element
}

// decodeExposure decodes [Exposure] from the XML tree.
func decodeExposure(root xmldoc.Element) (Exposure, error) {
	var exp Exposure

	// Decode optional MustHonor attribute with validation
	if attr, found := root.AttrByName(NsWSCN + ":MustHonor"); found {
		mustHonor := BooleanElement(attr.Value)
		if err := mustHonor.Validate(); err != nil {
			return exp, xmldoc.XMLErrWrap(root, fmt.Errorf("mustHonor: %w", err))
		}
		exp.MustHonor = optional.New(mustHonor)
	}

	// Copy child elements
	if root.Children != nil {
		exp.Children = make([]xmldoc.Element, len(root.Children))
		copy(exp.Children, root.Children)
	}

	return exp, nil
}

// toXMLExposure generates XML tree for the [Exposure].
func toXMLExposure(exp Exposure, name string) xmldoc.Element {
	elm := xmldoc.Element{
		Name:     name,
		Children: exp.Children,
	}

	// Add optional MustHonor attribute if present
	if mustHonor := optional.Get(exp.MustHonor); mustHonor != "" {
		elm.Attrs = []xmldoc.Attr{
			{
				Name:  NsWSCN + ":MustHonor",
				Value: string(mustHonor),
			},
		}
	}

	return elm
}
