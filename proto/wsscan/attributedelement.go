// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// AttributedElement: reusable type for elements with
// text value and optional wscn:MustHonor, wscn:Override, wscn:UsedDefault attributes

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// AttributedElement holds a value and optional wscn:MustHonor, wscn:Override,
// and wscn:UsedDefault attributes.
//
// The attributes are xs:string but must be boolean values: "0", "1", "false", or "true"
// (case-insensitive, whitespace ignored).
//
// This type is generic and can be used for elements like <wscn:Rotation>
// that have these attributes along with text content.
type AttributedElement[T any] struct {
	Value       T
	MustHonor   optional.Val[BooleanElement]
	Override    optional.Val[BooleanElement]
	UsedDefault optional.Val[BooleanElement]
}

// decodeAttributedElement fills the struct from an XML element.
//
// decodeValue is a function that decodes the value type T from a string.
func decodeAttributedElement[T any](
	root xmldoc.Element,
	decodeValue func(string) (T, error),
) (AttributedElement[T], error) {
	var elem AttributedElement[T]

	// Decode the value from text content
	var err error
	elem.Value, err = decodeValue(root.Text)
	if err != nil {
		return elem, err
	}

	// Decode optional attributes with validation
	if attr, found := root.AttrByName(NsWSCN + ":MustHonor"); found {
		mustHonor := BooleanElement(attr.Value)
		if err := mustHonor.Validate(); err != nil {
			return elem, xmldoc.XMLErrWrap(root, fmt.Errorf("mustHonor: %w", err))
		}
		elem.MustHonor = optional.New(mustHonor)
	}
	if attr, found := root.AttrByName(NsWSCN + ":Override"); found {
		override := BooleanElement(attr.Value)
		if err := override.Validate(); err != nil {
			return elem, xmldoc.XMLErrWrap(root, fmt.Errorf("override: %w", err))
		}
		elem.Override = optional.New(override)
	}
	if attr, found := root.AttrByName(NsWSCN + ":UsedDefault"); found {
		usedDefault := BooleanElement(attr.Value)
		if err := usedDefault.Validate(); err != nil {
			return elem, xmldoc.XMLErrWrap(root, fmt.Errorf("usedDefault: %w", err))
		}
		elem.UsedDefault = optional.New(usedDefault)
	}

	return elem, nil
}

// toXML creates an XML element from the struct.
//
// name is the XML element name (e.g., "wscn:Rotation").
// valueToString converts the value type T to its string representation.
func (a AttributedElement[T]) toXML(
	name string,
	valueToString func(T) string,
) xmldoc.Element {
	elm := xmldoc.Element{
		Name: name,
		Text: valueToString(a.Value),
	}

	// Add optional attributes if present
	attrs := make([]xmldoc.Attr, 0, 3)
	if mustHonor := optional.Get(a.MustHonor); mustHonor != "" {
		attrs = append(attrs, xmldoc.Attr{
			Name:  NsWSCN + ":MustHonor",
			Value: string(mustHonor),
		})
	}
	if override := optional.Get(a.Override); override != "" {
		attrs = append(attrs, xmldoc.Attr{
			Name:  NsWSCN + ":Override",
			Value: string(override),
		})
	}
	if usedDefault := optional.Get(a.UsedDefault); usedDefault != "" {
		attrs = append(attrs, xmldoc.Attr{
			Name:  NsWSCN + ":UsedDefault",
			Value: string(usedDefault),
		})
	}

	if len(attrs) > 0 {
		elm.Attrs = attrs
	}

	return elm
}
