// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ExposureSettings element
//
// The required <wscn:ExposureSettings> element contains individual
// adjustment values that the WSD Scan Service should apply to the
// image data after acquisition.

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ExposureSettings represents the <wscn:ExposureSettings> element.
//
// It has no attributes and the following optional child elements:
//   - <wscn:Brightness>
//   - <wscn:Contrast>
//   - <wscn:Sharpness>
//
// Each child element is modeled as [AttributedElement] with
// int value and optional Override / UsedDefault attributes.
type ExposureSettings struct {
	Brightness optional.Val[AttributedElement[int]]
	Contrast   optional.Val[AttributedElement[int]]
	Sharpness  optional.Val[AttributedElement[int]]
}

// toXML generates XML tree for the [ExposureSettings].
func (es ExposureSettings) toXML(name string) xmldoc.Element {
	children := make([]xmldoc.Element, 0, 3)

	if es.Brightness != nil {
		b := optional.Get(es.Brightness)
		children = append(children, b.toXML(NsWSCN+":Brightness",
			strconv.Itoa))
	}
	if es.Contrast != nil {
		c := optional.Get(es.Contrast)
		children = append(children, c.toXML(NsWSCN+":Contrast",
			strconv.Itoa))
	}
	if es.Sharpness != nil {
		s := optional.Get(es.Sharpness)
		children = append(children, s.toXML(NsWSCN+":Sharpness",
			strconv.Itoa))
	}

	return xmldoc.Element{
		Name:     name,
		Children: children,
	}
}

// decodeExposureSettings decodes [ExposureSettings] from the XML tree.
func decodeExposureSettings(root xmldoc.Element) (ExposureSettings, error) {
	var es ExposureSettings

	decodeValue := func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %w", err)
		}
		return val, nil
	}

	for _, child := range root.Children {
		switch child.Name {
		case NsWSCN + ":Brightness":
			val, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return es, fmt.Errorf("brightness: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			es.Brightness = optional.New(val)
		case NsWSCN + ":Contrast":
			val, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return es, fmt.Errorf("contrast: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			es.Contrast = optional.New(val)
		case NsWSCN + ":Sharpness":
			val, err := decodeAttributedElement(child, decodeValue)
			if err != nil {
				return es, fmt.Errorf("sharpness: %w",
					xmldoc.XMLErrWrap(child, err))
			}
			es.Sharpness = optional.New(val)
		}
	}

	return es, nil
}
