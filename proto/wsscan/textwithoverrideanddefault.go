// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// TextWithOverrideDefaultElement: reusable type for elements with
// text and optional wscn:Override and wscn:UsedDefault attributes

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TextWithOverrideAndDefault holds a text value and
// optional wscn:Override and wscn:UsedDefault boolean attributes.
type TextWithOverrideAndDefault struct {
	Text        string
	Override    optional.Val[BooleanElement]
	UsedDefault optional.Val[BooleanElement]
}

// decodeTextWithOverrideAndDefault fills the struct from an XML element.
func (t *TextWithOverrideAndDefault) decodeTextWithOverrideAndDefault(
	root xmldoc.Element) (TextWithOverrideAndDefault, error) {
	t.Text = root.Text
	if attr, found := root.AttrByName(NsWSCN + ":Override"); found {
		override, err := decodeBooleanElement(xmldoc.Element{Text: attr.Value})
		if err != nil {
			return *t, err
		}
		t.Override = optional.New(override)
	}
	if attr, found := root.AttrByName(NsWSCN + ":UsedDefault"); found {
		usedDefault, err := decodeBooleanElement(
			xmldoc.Element{Text: attr.Value})
		if err != nil {
			return *t, err
		}
		t.UsedDefault = optional.New(usedDefault)
	}
	return *t, nil
}

// toXML creates an XML element from the struct.
func (t TextWithOverrideAndDefault) toXML(name string) xmldoc.Element {
	elm := xmldoc.Element{Name: name, Text: t.Text}
	if t.Override != nil {
		elm.Attrs = append(elm.Attrs, xmldoc.Attr{
			Name:  NsWSCN + ":Override",
			Value: string(optional.Get(t.Override)),
		})
	}
	if t.UsedDefault != nil {
		elm.Attrs = append(elm.Attrs, xmldoc.Attr{
			Name:  NsWSCN + ":UsedDefault",
			Value: string(optional.Get(t.UsedDefault)),
		})
	}
	return elm
}
