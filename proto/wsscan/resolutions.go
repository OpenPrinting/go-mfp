// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Resolutions: simple struct for arrays of width and height values

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Resolutions represents arrays of width and height values.
type Resolutions struct {
	Widths  []int
	Heights []int
}

// toXML creates an XML element for Resolutions.
func (r Resolutions) toXML(name string) xmldoc.Element {
	children := make([]xmldoc.Element, 0, 2)

	// Widths
	if len(r.Widths) > 0 {
		widthChildren := make([]xmldoc.Element, 0, len(r.Widths))
		for _, width := range r.Widths {
			widthChildren = append(widthChildren, xmldoc.Element{
				Name: NsWSCN + ":Width",
				Text: strconv.Itoa(width),
			})
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":Widths",
			Children: widthChildren,
		})
	}

	// Heights
	if len(r.Heights) > 0 {
		heightChildren := make([]xmldoc.Element, 0, len(r.Heights))
		for _, height := range r.Heights {
			heightChildren = append(heightChildren, xmldoc.Element{
				Name: NsWSCN + ":Height",
				Text: strconv.Itoa(height),
			})
		}
		children = append(children, xmldoc.Element{
			Name:     NsWSCN + ":Heights",
			Children: heightChildren,
		})
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeResolutions decodes a Resolutions from an XML element.
// Expects wrapped form with Widths and Heights containers.
func decodeResolutions(root xmldoc.Element) (Resolutions, error) {
	var res Resolutions

	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Widths":
			for _, wchild := range child.Children {
				if wchild.Name == NsWSCN+":Width" {
					width, err := strconv.Atoi(wchild.Text)
					if err != nil {
						return res, fmt.Errorf("invalid width value: %w", err)
					}
					res.Widths = append(res.Widths, width)
				}
			}
		case NsWSCN + ":Heights":
			for _, hchild := range child.Children {
				if hchild.Name == NsWSCN+":Height" {
					height, err := strconv.Atoi(hchild.Text)
					if err != nil {
						return res, fmt.Errorf("invalid height value: %w", err)
					}
					res.Heights = append(res.Heights, height)
				}
			}
		}
	}

	if len(res.Widths) == 0 || len(res.Heights) == 0 {
		return res, fmt.Errorf("missing Width or Height elements")
	}

	return res, nil
}
