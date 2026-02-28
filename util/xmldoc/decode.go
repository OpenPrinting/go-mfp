// MFP - Miulti-Function Printers and scanners toolkit
// XML mini library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// XML decoder

package xmldoc

import (
	"encoding/xml"
	"io"
	"strings"
)

// Decode parses XML document, and represents it as a tree of
// [Element]s.
//
// Namespace prefixes are rewritten according to the 'ns' map.
// Full namespace URL used as map index, and value that corresponds
// to the index replaced with map value. If URL is not found in the
// map, prefix replaced with "-" string
func Decode(ns Namespace, in io.Reader) (Element, error) {
	var elem Element
	stack := []Element{}
	decoder := xml.NewDecoder(in)

	for {
		token, err := decoder.Token()
		if err != nil {
			return Element{}, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			// Decode name and path.
			// Namespace translation is handled here.
			var name string
			if t.Name.Space != "" {
				var ok bool
				name, ok = ns.ByURL(t.Name.Space)
				if !ok {
					name = "-"
				}
			}

			if name != "" {
				name += ":"
			}
			name += t.Name.Local

			// Create an element
			stack = append(stack, elem)
			line, _ := decoder.InputPos()
			elem = Element{Name: name, Line: line}

			// Decode attributes
			for _, attr := range t.Attr {
				if attr.Name.Space == "xmlns" {
					// Skip xmlns attributes, they
					// are for XML namespace management.
					// On encoding we are insert them
					// automatically, so they are
					// removed on decoding, for symmetry.
					continue
				}

				name = ""
				if attr.Name.Space != "" {
					var ok bool
					name, ok = ns.ByURL(attr.Name.Space)
					if !ok {
						name = "-"
					}
				}

				if name != "" {
					name += ":"
				}
				name += attr.Name.Local

				elem.Attrs = append(elem.Attrs,
					Attr{name, attr.Value})
			}

		case xml.EndElement:
			elem.Text = strings.TrimSpace(elem.Text)

			if len(stack) == 1 {
				return elem, nil
			}

			parent := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			parent.Children = append(parent.Children, elem)
			elem = parent

		case xml.CharData:
			elem.Text += string(t)
		}
	}
}
