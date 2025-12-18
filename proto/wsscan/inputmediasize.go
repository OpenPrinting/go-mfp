// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// InputMediaSize element

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// InputMediaSize represents the required <wscn:InputMediaSize> element
// that specifies the size of the media to be scanned for the current job.
//
// It has no attributes and the following required child elements:
//   - <wscn:Width>
//   - <wscn:Height>
//
// Both Width and Height are modeled as [AttributedElement] with
// int value and optional Override / UsedDefault attributes.
// The values must be in the range from 1 through 2147483648 and are
// in units of one-thousandths (1/1000) of an inch.
type InputMediaSize struct {
	Width  AttributedElement[int]
	Height AttributedElement[int]
}

// toXML generates XML tree for the [InputMediaSize].
func (ims InputMediaSize) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Children: []xmldoc.Element{
			ims.Width.toXML(NsWSCN+":Width", strconv.Itoa),
			ims.Height.toXML(NsWSCN+":Height", strconv.Itoa),
		},
	}
}

// decodeInputMediaSize decodes [InputMediaSize] from the XML tree.
func decodeInputMediaSize(root xmldoc.Element) (InputMediaSize, error) {
	var ims InputMediaSize

	decodeValue := func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %w", err)
		}
		return val, nil
	}

	var widthFound, heightFound bool
	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Width":
			width, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return ims, fmt.Errorf("width: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			ims.Width = width
			widthFound = true
		case NsWSCN + ":Height":
			height, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return ims, fmt.Errorf("height: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			ims.Height = height
			heightFound = true
		}
	}

	if !widthFound {
		return ims, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:Width", NsWSCN))
	}
	if !heightFound {
		return ims, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:Height", NsWSCN))
	}

	return ims, nil
}
