// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Resolution element

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Resolution represents the optional <wscn:Resolution> element
// that specifies the resolution of the scanned image.
//
// It includes an optional wscn:MustHonor attribute (xs:string,
// but should be a boolean value: 0, false, 1, or true).
//
// The element contains child elements:
//   - Width (required AttributedElement[int]) - resolution width in pixels per inch
//   - Height (optional AttributedElement[int]) - resolution height in pixels per inch
//     If Height is missing, the Width value should be used, yielding a square resolution
//     (for example, 300 x 300).
type Resolution struct {
	MustHonor optional.Val[BooleanElement]
	Width     AttributedElement[int]
	Height    optional.Val[AttributedElement[int]]
}

// toXML generates XML tree for the [Resolution].
func (res Resolution) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		res.Width.toXML(NsWSCN+":Width", strconv.Itoa),
	}

	// Add Height if present
	if res.Height != nil {
		height := optional.Get(res.Height)
		children = append(children, height.toXML(NsWSCN+":Height", strconv.Itoa))
	}

	elm := xmldoc.Element{
		Name:     name,
		Children: children,
	}

	// Add optional MustHonor attribute if present
	if mustHonor := optional.Get(res.MustHonor); mustHonor != "" {
		elm.Attrs = []xmldoc.Attr{
			{
				Name:  NsWSCN + ":MustHonor",
				Value: string(mustHonor),
			},
		}
	}

	return elm
}

// decodeResolution decodes [Resolution] from the XML tree.
func decodeResolution(root xmldoc.Element) (Resolution, error) {
	var res Resolution

	// Decode optional MustHonor attribute with validation
	if attr, found := root.AttrByName(NsWSCN + ":MustHonor"); found {
		mustHonor := BooleanElement(attr.Value)
		if err := mustHonor.Validate(); err != nil {
			return res, xmldoc.XMLErrWrap(root, fmt.Errorf("mustHonor: %w", err))
		}
		res.MustHonor = optional.New(mustHonor)
	}

	decodeValue := func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %w", err)
		}
		return val, nil
	}

	// Decode child elements
	var widthFound bool
	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Width":
			width, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return res, fmt.Errorf("width: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			res.Width = width
			widthFound = true
		case NsWSCN + ":Height":
			height, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return res, fmt.Errorf("height: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			res.Height = optional.New(height)
		}
	}

	if !widthFound {
		return res, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:Width", NsWSCN))
	}

	return res, nil
}
