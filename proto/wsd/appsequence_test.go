// MFP - Miulti-Function Printers and scanners toolkit
// WSD core protocol
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// AppSequence test

package wsd

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TestAppSequence tests AppSequence
func TestAppSequence(t *testing.T) {
	type testData struct {
		seq AppSequence
		xml xmldoc.Element
	}

	const urn AnyURI = "urn:uuid:2a443ed7-5ee5-498d-a302-73ff91ea9ea0"

	tests := []testData{
		{
			seq: AppSequence{
				InstanceID:    123456789,
				MessageNumber: 123,
			},
			xml: xmldoc.WithAttrs(NsDiscovery+":"+"AppSequence",
				xmldoc.Attr{Name: "InstanceId", Value: "123456789"},
				xmldoc.Attr{Name: "MessageNumber", Value: "123"},
			),
		},

		{
			seq: AppSequence{
				InstanceID:    987654321,
				MessageNumber: 321,
				SequenceID:    optional.New(urn),
			},
			xml: xmldoc.WithAttrs(NsDiscovery+":"+"AppSequence",
				xmldoc.Attr{Name: "InstanceId", Value: "987654321"},
				xmldoc.Attr{Name: "MessageNumber", Value: "321"},
				xmldoc.Attr{
					Name:  "SequenceId",
					Value: "urn:uuid:2a443ed7-5ee5-498d-a302-73ff91ea9ea0",
				},
			),
		},
	}

	for _, test := range tests {
		xml := test.seq.ToXML()
		if !reflect.DeepEqual(xml, test.xml) {
			t.Errorf("ToXML:\nexpected: %s\npresent: %s\n",
				test.xml.EncodeString(NsMap),
				xml.EncodeString(NsMap))
		}

		seq, err := DecodeAppSequence(xml)
		if err != nil {
			t.Errorf("DecodeAppSequence: %s", err)
			continue
		}

		if !reflect.DeepEqual(seq, test.seq) {
			t.Errorf("DecodeAppSequence:\n"+
				"expected: %#v\npresent:  %#v\n",
				test.seq, seq)
		}
	}
}

// TestAppSequenceDecodeErrors tests AppSequence decode errors
func TestAppSequenceDecodeErrors(t *testing.T) {
	type testData struct {
		xml  xmldoc.Element
		estr string
	}

	tests := []testData{
		{
			xml: xmldoc.WithAttrs(NsDiscovery+":"+"AppSequence",
				xmldoc.Attr{Name: "InstanceId", Value: "123456789"},
				xmldoc.Attr{Name: "MessageNumber", Value: "123"},
			),
			estr: "",
		},

		{
			xml: xmldoc.WithAttrs(NsDiscovery+":"+"AppSequence",
				xmldoc.Attr{Name: "InstanceId", Value: "123456789"},
			),
			estr: "/d:AppSequence/d:AppSequence/@MessageNumber: missed attribyte",
		},

		{
			xml: xmldoc.WithAttrs(NsDiscovery+":"+"AppSequence",
				xmldoc.Attr{Name: "MessageNumber", Value: "123"},
			),
			estr: "/d:AppSequence/d:AppSequence/@InstanceId: missed attribyte",
		},

		{
			xml: xmldoc.WithAttrs(NsDiscovery+":"+"AppSequence",
				xmldoc.Attr{Name: "InstanceId", Value: "ABC"},
				xmldoc.Attr{Name: "MessageNumber", Value: "123"},
			),
			estr: `/d:AppSequence/@InstanceId: invalid uint: "ABC"`,
		},
	}

	for _, test := range tests {
		_, err := DecodeAppSequence(test.xml)
		estr := ""
		if err != nil {
			estr = err.Error()
		}

		if estr != test.estr {
			t.Errorf("%s\nexpected: %s\npresent:  %s",
				test.xml.EncodeString(NsMap),
				test.estr, estr)
		}
	}
}
