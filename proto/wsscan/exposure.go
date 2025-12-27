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
// The element contains child elements: AutoExposure (required)
// and ExposureSettings (required).
type Exposure struct {
	MustHonor        optional.Val[BooleanElement]
	AutoExposure     BooleanElement
	ExposureSettings ExposureSettings
}

// toXML generates XML tree for the [Exposure].
func (exp Exposure) toXML(name string) xmldoc.Element {
	children := []xmldoc.Element{
		exp.AutoExposure.toXML(NsWSCN + ":AutoExposure"),
		exp.ExposureSettings.toXML(NsWSCN + ":ExposureSettings"),
	}

	elm := xmldoc.Element{
		Name:     name,
		Children: children,
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

	// Decode required child elements
	var autoExposureFound, exposureSettingsFound bool
	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":AutoExposure":
			autoExp, err := decodeBooleanElement(child)
			if err != nil {
				return exp, fmt.Errorf("autoExposure: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			exp.AutoExposure = autoExp
			autoExposureFound = true
		case NsWSCN + ":ExposureSettings":
			expSettings, err := decodeExposureSettings(child)
			if err != nil {
				return exp, fmt.Errorf("exposureSettings: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			exp.ExposureSettings = expSettings
			exposureSettingsFound = true
		}
	}

	if !autoExposureFound {
		return exp, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:AutoExposure", NsWSCN))
	}
	if !exposureSettingsFound {
		return exp, xmldoc.XMLErrWrap(root,
			fmt.Errorf("missing required element: %s:ExposureSettings", NsWSCN))
	}

	return exp, nil
}
