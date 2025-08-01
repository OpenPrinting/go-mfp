// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestRangeElement_XMLRoundTrip(t *testing.T) {
	cases := []struct {
		min, max int
	}{
		{1, 100},
		{50, 75},
		{100, 100},
		{0, 1000},
	}

	for _, c := range cases {
		re := Range{MinValue: c.min, MaxValue: c.max}
		elm := xmldoc.Element{
			Name:     NsWSCN + ":TestRange",
			Children: re.toXML(),
		}
		parsed, err := decodeRange(elm)
		if err != nil {
			t.Errorf("decodeRangeElement: input %+v, unexpected error: %v",
				re, err)
		}
		if !reflect.DeepEqual(parsed, re) {
			t.Errorf("XML round-trip: expected %+v, got %+v", re, parsed)
		}
	}
}

func TestRangeElement_DecodeErrors(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() xmldoc.Element
		errContains string
	}{
		{
			name: "missing MinValue",
			setup: func() xmldoc.Element {
				// Only MaxValue present
				return xmldoc.Element{
					Name: NsWSCN + ":TestRange",
					Children: []xmldoc.Element{{
						Name: NsWSCN + ":MaxValue",
						Text: "100",
					}},
				}
			},
			errContains: "MinValue",
		},
		{
			name: "missing MaxValue",
			setup: func() xmldoc.Element {
				// Only MinValue present
				return xmldoc.Element{
					Name: NsWSCN + ":TestRange",
					Children: []xmldoc.Element{{
						Name: NsWSCN + ":MinValue",
						Text: "1",
					}},
				}
			},
			errContains: "MaxValue",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			elm := tc.setup()
			_, err := decodeRange(elm)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !containsSubstring(err.Error(), tc.errContains) {
				t.Errorf("expected error to contain %q, got %q",
					tc.errContains, err)
			}
		})
	}
}
