// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// TextWithLangElement: reusable type for elements with
// text and optional xml:lang

package wsscan

import (
	"reflect"
	"strings"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TextWithLangElement holds a text value and an optional xml:lang attribute.
type TextWithLangElement struct {
	Text string
	Lang optional.Val[string]
}

// HasOptions reports if value really has any options set.
// It implements the [Wrapper] interface.
func (t TextWithLangElement) HasOptions() bool {
	return t.Lang != nil
}

// Unwrap returns the underlying value, if t has no options, or the
// t's value itself otherwise.
//
// It implements the [Wrapper] interface.
func (t TextWithLangElement) Unwrap() any {
	if !t.HasOptions() {
		return t.Text
	}
	return t
}

// Wrap wraps the simple value into the Wrapper
// type and returns the new wrapped value.
//
// In case the provided value cannot be converted
// into the Wrapper's underlying type, this function
// returns nil.
func (t TextWithLangElement) Wrap(v any) any {
	val, ok := v.(string)
	if ok {
		return TextWithLangElement{Text: val}
	}
	println("FAIL", reflect.TypeOf(v).String())
	return nil
}

// decodeTextWithLangElement fills the struct from an XML element.
func (t *TextWithLangElement) decodeTextWithLangElement(root xmldoc.Element) (
	TextWithLangElement, error) {
	t.Text = root.Text
	if attr, found := root.AttrByName("xml:lang"); found {
		t.Lang = optional.New(attr.Value)
	}
	return *t, nil
}

// toXML creates an XML element from the struct.
func (t TextWithLangElement) toXML(name string) xmldoc.Element {
	elm := xmldoc.Element{Name: name, Text: t.Text}
	lang := optional.Get(t.Lang)
	if lang != "" {
		elm.Attrs = []xmldoc.Attr{{
			Name:  "xml:lang",
			Value: lang,
		}}
	}
	return elm
}

// TextWithLangList represents a list of text elements with
// optional language attributes.
type TextWithLangList []TextWithLangElement

// NeutralLang returns the most language-neutral entry from the list.
//
// It uses the following preferences:
//   - entry without Lang is the best match
//   - if not found, search for the "en" version
//   - if not found, search for the "en-US" version
//   - if not found, search for the first entry starting with "en-"
//   - if not found, return the first entry from the list
//
// If list is empty, it returns TextWithLangElement{}.
func (tl TextWithLangList) NeutralLang() TextWithLangElement {
	var en, enUS, enAny TextWithLangElement
	var enFound, enUSFound, enAnyFound bool

	for _, t := range tl {
		lang := strings.ToLower(optional.Get(t.Lang))
		switch lang {
		case "":
			return t // The best match; return immediately
		case "en":
			if !enFound {
				en = t
				enFound = true
			}
		case "en-us":
			if !enUSFound {
				enUS = t
				enUSFound = true
			}
		default:
			if strings.HasPrefix(lang, "en-") && !enAnyFound {
				enAny = t
				enAnyFound = true
			}
		}
	}

	switch {
	case enFound:
		return en
	case enUSFound:
		return enUS
	case enAnyFound:
		return enAny
	case len(tl) != 0:
		return tl[0]
	}

	return TextWithLangElement{}
}
