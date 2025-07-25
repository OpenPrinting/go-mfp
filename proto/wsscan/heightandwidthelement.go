// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// HeightAndWidthElement: reusable element with Height and Width children

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// HeightAndWidthElement holds Heights and Widths child elements,
// each with text and optional wscn:Override and wscn:UsedDefault attributes.
type HeightAndWidthElement struct {
	Heights []TextWithOverrideAndDefault
	Widths  []TextWithOverrideAndDefault
}

// toXML creates an XML element for HeightAndWidthElement.
func (hwe HeightAndWidthElement) toXML(name string) xmldoc.Element {
	children := make([]xmldoc.Element, 0, 2)
	// Heights
	if len(hwe.Heights) == 1 {
		children = append(children, hwe.Heights[0].toXML(NsWSCN+":Height"))
	} else if len(hwe.Heights) > 1 {
		heightChildren := make([]xmldoc.Element, 0, len(hwe.Heights))
		for _, h := range hwe.Heights {
			heightChildren = append(heightChildren, h.toXML(NsWSCN+":Height"))
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":Heights",
			Children: heightChildren,
		})
	}
	// Widths
	if len(hwe.Widths) == 1 {
		children = append(children, hwe.Widths[0].toXML(NsWSCN+":Width"))
	} else if len(hwe.Widths) > 1 {
		widthChildren := make([]xmldoc.Element, 0, len(hwe.Widths))
		for _, w := range hwe.Widths {
			widthChildren = append(widthChildren, w.toXML(NsWSCN+":Width"))
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":Widths",
			Children: widthChildren,
		})
	}
	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeHeightAndWidthElement decodes a HeightAndWidthElement
// from an XML element. Accepts both direct and wrapped forms.
func decodeHeightAndWidthElement(root xmldoc.Element) (
	HeightAndWidthElement, error) {
	var hwe HeightAndWidthElement
	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Height":
			h, err := new(TextWithOverrideAndDefault).decodeTextWithOverrideAndDefault(child)
			if err != nil {
				return hwe, fmt.Errorf("height: %w", err)
			}
			hwe.Heights = append(hwe.Heights, h)
		case NsWSCN + ":Heights":
			for _, hchild := range child.Children {
				if hchild.Name == NsWSCN+":Height" {
					h, err := new(TextWithOverrideAndDefault).decodeTextWithOverrideAndDefault(hchild)
					if err != nil {
						return hwe, fmt.Errorf("height: %w", err)
					}
					hwe.Heights = append(hwe.Heights, h)
				}
			}
		case NsWSCN + ":Width":
			w, err := new(TextWithOverrideAndDefault).decodeTextWithOverrideAndDefault(child)
			if err != nil {
				return hwe, fmt.Errorf("width: %w", err)
			}
			hwe.Widths = append(hwe.Widths, w)
		case NsWSCN + ":Widths":
			for _, wchild := range child.Children {
				if wchild.Name == NsWSCN+":Width" {
					w, err := new(TextWithOverrideAndDefault).decodeTextWithOverrideAndDefault(wchild)
					if err != nil {
						return hwe, fmt.Errorf("width: %w", err)
					}
					hwe.Widths = append(hwe.Widths, w)
				}
			}
		}
	}
	if len(hwe.Heights) == 0 || len(hwe.Widths) == 0 {
		return hwe, fmt.Errorf("missing Height or Width element")
	}
	return hwe, nil
}
