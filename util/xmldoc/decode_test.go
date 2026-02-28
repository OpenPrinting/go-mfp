// MFP - Miulti-Function Printers and scanners toolkit
// XML mini library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package xmldoc

import (
	"bytes"
	"testing"
)

// TestDecode tests Decode function
func TestDecode(t *testing.T) {
	ns := Namespace{
		{URL: `http://example.com/a`, Prefix: `a`},
		{URL: `http://example.com/b`, Prefix: `b`},
		{URL: `http://example.com/c`, Prefix: `c`},
	}

	in := `` +
		`<?xml version="1.0" ?>` +
		`<env xmlns:ns-a="http://example.com/a" xmlns:ns-b="http://example.com/b" xmlns:ns-c="http://example.com/c">` +
		`  <ns-a:elem-a ns-a:attr="value">body a</ns-a:elem-a>` +
		`  <ns-b:elem-b>body b` +
		`    <ns-b:nested-1>nested body 1</ns-b:nested-1>` +
		`    <ns-b:nested-2>nested body 2` +
		`      <ns-b:nested-2-1>nested body 2-1</ns-b:nested-2-1>` +
		`    </ns-b:nested-2>` +
		`  </ns-b:elem-b>` +
		`  <ns-c:elem-c>body c</ns-c:elem-c>` +
		`  <ns-d:elem-d>body d</ns-d:elem-d>` +
		`</env>` +
		``

	expect := Element{
		Name: "env",
		//Attrs: []Attr{
		//	{Name: "xmlns:ns-a", Value: "http://example.com/a"},
		//	{Name: "xmlns:ns-b", Value: "http://example.com/b"},
		//	{Name: "xmlns:ns-c", Value: "http://example.com/c"},
		//},
		Children: []Element{
			{
				Name: "a:elem-a",
				Text: "body a",
				Attrs: []Attr{
					{Name: "a:attr", Value: "value"},
				},
			},
			{
				Name: "b:elem-b",
				Text: "body b",
				Children: []Element{
					{
						Name: "b:nested-1",
						Text: "nested body 1",
					},
					{
						Name: "b:nested-2",
						Text: "nested body 2",
						Children: []Element{
							{
								Name: "b:nested-2-1",
								Text: "nested body 2-1",
							},
						},
					},
				},
			},
			{
				Name: "c:elem-c",
				Text: "body c",
			},
			{
				Name: "-:elem-d",
				Text: "body d",
			},
		},
	}

	out, err := Decode(ns, bytes.NewReader([]byte(in)))
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if !out.Equal(expect) {
		fmtexp := expect.EncodeIndentString(nil, "  ")
		fmtout := out.EncodeIndentString(nil, "  ")
		t.Errorf("expected:\n%s\npresent:\n%s\n",
			fmtexp, fmtout)
	}
}
