// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// InputSize element

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// InputSize represents the optional <wscn:InputSize> element
// that specifies the size of the original scan media.
//
// It includes an optional wscn:MustHonor attribute (xs:string,
// but should be a boolean value: 0, false, 1, or true).
//
// The element contains child elements:
//   - DocumentSizeAutoDetect (optional BooleanElement)
//   - InputMediaSize (required InputMediaSize)
type InputSize struct {
	MustHonor              optional.Val[BooleanElement]
	DocumentSizeAutoDetect optional.Val[BooleanElement]
	InputMediaSize         InputMediaSize
}

// toXML generates XML tree for the [InputSize].
func (is InputSize) toXML(name string) xmldoc.Element {
	children := make([]xmldoc.Element, 0, 2)

	// Add DocumentSizeAutoDetect if present
	if is.DocumentSizeAutoDetect != nil {
		autoDetect := optional.Get(is.DocumentSizeAutoDetect)
		children = append(children,
			autoDetect.toXML(NsWSCN+":DocumentSizeAutoDetect"))
	}

	// Add InputMediaSize (required)
	children = append(children, is.InputMediaSize.toXML(NsWSCN+":InputMediaSize"))

	elm := xmldoc.Element{
		Name:     name,
		Children: children,
	}

	// Add optional MustHonor attribute if present
	if mustHonor := optional.Get(is.MustHonor); mustHonor != "" {
		elm.Attrs = []xmldoc.Attr{
			{
				Name:  NsWSCN + ":MustHonor",
				Value: string(mustHonor),
			},
		}
	}

	return elm
}

// decodeInputSize decodes [InputSize] from the XML tree.
func decodeInputSize(root xmldoc.Element) (InputSize, error) {
	var is InputSize

	// Decode optional MustHonor attribute with validation
	if attr, found := root.AttrByName(NsWSCN + ":MustHonor"); found {
		mustHonor := BooleanElement(attr.Value)
		if err := mustHonor.Validate(); err != nil {
			return is, xmldoc.XMLErrWrap(root, fmt.Errorf("mustHonor: %w", err))
		}
		is.MustHonor = optional.New(mustHonor)
	}

	// Decode child elements
	var inputMediaSizeFound bool
	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":DocumentSizeAutoDetect":
			autoDetect, err := decodeBooleanElement(child)
			if err != nil {
				return is, fmt.Errorf("documentSizeAutoDetect: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			is.DocumentSizeAutoDetect = optional.New(autoDetect)
		case NsWSCN + ":InputMediaSize":
			mediaSize, err := decodeInputMediaSize(child)
			if err != nil {
				return is, fmt.Errorf("inputMediaSize: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			is.InputMediaSize = mediaSize
			inputMediaSizeFound = true
		}
	}

	if !inputMediaSizeFound {
		return is, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:InputMediaSize", NsWSCN))
	}

	return is, nil
}
