// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Scaling element

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Scaling represents the optional <wscn:Scaling> element
// that specifies the scaling of both the width and height of the scanned document.
//
// It includes an optional wscn:MustHonor attribute (xs:string,
// but should be a boolean value: 0, false, 1, or true).
//
// The element contains child elements:
//   - ScalingWidth (required AttributedElement[int]) - scaling width value in range 1-1000
//     Note: ScalingWidth should only use Override and UsedDefault attributes, not MustHonor
//   - ScalingHeight (required AttributedElement[int]) - scaling height value in range 1-1000
//     Note: ScalingHeight should only use Override and UsedDefault attributes, not MustHonor
type Scaling struct {
	MustHonor     optional.Val[BooleanElement]
	ScalingWidth  AttributedElement[int]
	ScalingHeight AttributedElement[int]
}

// toXML generates XML tree for the [Scaling].
func (sc Scaling) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		sc.ScalingWidth.toXML(NsWSCN+":ScalingWidth", strconv.Itoa),
		sc.ScalingHeight.toXML(NsWSCN+":ScalingHeight", strconv.Itoa),
	}

	elm := xmldoc.Element{
		Name:     name,
		Children: children,
	}

	if mustHonor := optional.Get(sc.MustHonor); mustHonor != "" {
		elm.Attrs = []xmldoc.Attr{
			{
				Name:  NsWSCN + ":MustHonor",
				Value: string(mustHonor),
			},
		}
	}

	return elm
}

// decodeScaling decodes [Scaling] from the XML tree.
func decodeScaling(root xmldoc.Element) (Scaling, error) {
	var sc Scaling

	// Decode optional MustHonor attribute with validation
	if attr, found := root.AttrByName(NsWSCN + ":MustHonor"); found {
		mustHonor := BooleanElement(attr.Value)
		if err := mustHonor.Validate(); err != nil {
			return sc, xmldoc.XMLErrWrap(root, fmt.Errorf("mustHonor: %w", err))
		}
		sc.MustHonor = optional.New(mustHonor)
	}

	decodeValue := func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %w", err)
		}
		return val, nil
	}

	// Decode child elements
	var widthFound, heightFound bool
	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":ScalingWidth":
			width, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return sc, fmt.Errorf("scalingWidth: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			sc.ScalingWidth = width
			widthFound = true
		case NsWSCN + ":ScalingHeight":
			height, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return sc, fmt.Errorf("scalingHeight: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			sc.ScalingHeight = height
			heightFound = true
		}
	}

	if !widthFound {
		return sc, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:ScalingWidth", NsWSCN))
	}
	if !heightFound {
		return sc, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:ScalingHeight", NsWSCN))
	}

	return sc, nil
}
