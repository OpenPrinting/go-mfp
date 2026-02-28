// MFP - Miulti-Function Printers and scanners toolkit
// XML mini library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Element.Expand test

package xmldoc

import (
	"testing"
)

// TestExpand tests Element.Expand function
func TestExpand(t *testing.T) {
	in := Element{
		Name: "root",
		Text: "expanded $ROOT",
		Attrs: []Attr{
			{"attr1", "expanded $ATTR1"},
			{"attr2", "expanded ${ATTR2}"},
		},
		Children: []Element{
			{
				Name: "child-1",
				Text: "expanded $CHILD1",
			},
			{
				Name: "child-2",
				Text: "expanded $CHILD2",
			},
		},
	}

	expected := Element{
		Name: "root",
		Text: "expanded root value",
		Attrs: []Attr{
			{"attr1", "expanded attr1 value"},
			{"attr2", "expanded attr2 value"},
		},
		Children: []Element{
			{
				Name: "child-1",
				Text: "expanded child1 value",
			},
			{
				Name: "child-2",
				Text: "expanded child2 value",
			},
		},
	}

	out := in.Expand(
		func(in string) string {
			switch in {
			case "ROOT":
				return "root value"
			case "ATTR1":
				return "attr1 value"
			case "ATTR2":
				return "attr2 value"
			case "CHILD1":
				return "child1 value"
			case "CHILD2":
				return "child2 value"
			}

			return ""
		})

	if !out.Equal(expected) {
		t.Errorf("expected:\n%s\npresent:\n%s\n",
			expected.EncodeIndentString(nil, "  "),
			out.EncodeIndentString(nil, "  "))
	}
}
