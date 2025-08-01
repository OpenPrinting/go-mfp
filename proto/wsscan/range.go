// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// range element: reusable min/max range type

package wsscan

import (
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// Range represents a range with minimum and maximum values.
// This is a common pattern used across multiple scanner configuration elements.
type Range struct {
	MinValue int
	MaxValue int
}

// toXML generates XML elements for MinValue and MaxValue.
func (re Range) toXML() []xmldoc.Element {
	return []xmldoc.Element{
		{
			Name: NsWSCN + ":MinValue",
			Text: strconv.Itoa(re.MinValue),
		},
		{
			Name: NsWSCN + ":MaxValue",
			Text: strconv.Itoa(re.MaxValue),
		},
	}
}

// decodeRangeElement decodes a RangeElement from an XML element.
func decodeRange(root xmldoc.Element) (re Range, err error) {
	defer func() { err = xmldoc.XMLErrWrap(root, err) }()

	// Lookup MinValue and MaxValue elements
	minLookup := xmldoc.Lookup{
		Name:     NsWSCN + ":MinValue",
		Required: true,
	}
	maxLookup := xmldoc.Lookup{
		Name:     NsWSCN + ":MaxValue",
		Required: true,
	}

	missed := root.Lookup(&minLookup, &maxLookup)
	if missed != nil {
		return re, xmldoc.XMLErrMissed(missed.Name)
	}

	// Decode MinValue
	minVal, err := strconv.Atoi(minLookup.Elem.Text)
	if err != nil {
		return re, fmt.Errorf("MinValue: invalid integer %q",
			minLookup.Elem.Text)
	}
	re.MinValue = minVal

	// Decode MaxValue
	maxVal, err := strconv.Atoi(maxLookup.Elem.Text)
	if err != nil {
		return re, fmt.Errorf("MaxValue: invalid integer %q",
			maxLookup.Elem.Text)
	}
	re.MaxValue = maxVal

	return re, nil
}
