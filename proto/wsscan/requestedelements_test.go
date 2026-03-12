// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// RequestedElements tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TestRequestedElementsRoundTrip checks that encoding RequestedElements to
// XML and decoding it back produces an identical struct.
func TestRequestedElementsRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		re   RequestedElements
	}{
		{
			// Minimum valid case: exactly one ScannerRequestedElementNames value.
			name: "single name",
			re: RequestedElements{
				Names: []ScannerRequestedElementNames{
					ScannerRequestedElementDescription,
				},
			},
		},
		{
			// All three standard QName values to verify order is
			// preserved and all entries survive the round trip.
			name: "multiple names",
			re: RequestedElements{
				Names: []ScannerRequestedElementNames{
					ScannerRequestedElementDescription,
					ScannerRequestedElementConfiguration,
					ScannerRequestedElementStatus,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML then decode back.
			xml := tt.re.toXML(NsWSCN + ":RequestedElements")

			decoded, err := decodeRequestedElements(xml)
			if err != nil {
				t.Fatalf("decodeRequestedElements() error = %v", err)
			}

			// Decoded value must be identical to the original.
			if !reflect.DeepEqual(decoded, tt.re) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.re, decoded, xml.EncodeString(nil))
			}
		})
	}
}

// TestRequestedElementsDecodeError checks that decode returns an error when
// no Name children are present, and succeeds for a valid element.
func TestRequestedElementsDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			// At least one Name is required; an empty element must fail.
			name: "no Name children",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":RequestedElements",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			// A single valid Name child must decode without error.
			name: "valid single Name",
			xml: xmldoc.Element{
				Name: NsWSCN + ":RequestedElements",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":Name",
						Text: NsWSCN + ":ScannerDescription",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeRequestedElements(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeRequestedElements() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}
