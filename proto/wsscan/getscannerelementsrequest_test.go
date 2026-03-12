// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// GetScannerElementsRequest tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TestGetScannerElementsRequestRoundTrip checks that encoding a
// GetScannerElementsRequest to XML and decoding it back produces
// an identical struct.
func TestGetScannerElementsRequestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		gser GetScannerElementsRequest
	}{
		{
			// A request with a single element name: the minimum
			// valid case.
			name: "single name",
			gser: GetScannerElementsRequest{
				RequestedElements: RequestedElements{
					Names: []ScannerRequestedElementNames{
						ScannerRequestedElementDescription,
					},
				},
			},
		},
		{
			// A request with all three standard element names to
			// verify that multiple Names survive the round trip
			// in order.
			name: "multiple names",
			gser: GetScannerElementsRequest{
				RequestedElements: RequestedElements{
					Names: []ScannerRequestedElementNames{
						ScannerRequestedElementDescription,
						ScannerRequestedElementConfiguration,
						ScannerRequestedElementStatus,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML then decode back.
			xml := tt.gser.toXML(NsWSCN + ":GetScannerElementsRequest")

			decoded, err := decodeGetScannerElementsRequest(xml)
			if err != nil {
				t.Fatalf("decodeGetScannerElementsRequest() error = %v", err)
			}

			// Decoded value must be identical to the original.
			if !reflect.DeepEqual(decoded, tt.gser) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.gser, decoded, xml.EncodeString(nil))
			}
		})
	}
}

// TestGetScannerElementsRequestDecodeError checks that decode returns an
// error for invalid XML and succeeds for well-formed input.
func TestGetScannerElementsRequestDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			// RequestedElements is required; an empty body must fail.
			name: "missing RequestedElements",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":GetScannerElementsRequest",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			// RequestedElements must contain at least one Name child;
			// an empty RequestedElements element must fail.
			name: "RequestedElements with no Name children",
			xml: xmldoc.Element{
				Name: NsWSCN + ":GetScannerElementsRequest",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":RequestedElements",
						Children: []xmldoc.Element{},
					},
				},
			},
			wantErr: true,
		},
		{
			// A fully valid request with one Name must decode
			// without error.
			name: "valid single name",
			xml: xmldoc.Element{
				Name: NsWSCN + ":GetScannerElementsRequest",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":RequestedElements",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":Name",
								Text: NsWSCN + ":ScannerDescription",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeGetScannerElementsRequest(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeGetScannerElementsRequest() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

// TestGetScannerElementsRequestToXML checks the exact XML structure produced
// by toXML: element name, single RequestedElements child, and its Name
// children.
func TestGetScannerElementsRequestToXML(t *testing.T) {
	gser := GetScannerElementsRequest{
		RequestedElements: RequestedElements{
			Names: []ScannerRequestedElementNames{
				ScannerRequestedElementDescription,
				ScannerRequestedElementConfiguration,
			},
		},
	}

	xml := gser.toXML(NsWSCN + ":GetScannerElementsRequest")

	// Root element must be wscn:GetScannerElementsRequest.
	if xml.Name != NsWSCN+":GetScannerElementsRequest" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":GetScannerElementsRequest", xml.Name)
	}

	// Must have exactly one child: wscn:RequestedElements.
	if len(xml.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(xml.Children))
	}

	reqElem := xml.Children[0]
	if reqElem.Name != NsWSCN+":RequestedElements" {
		t.Errorf("Expected RequestedElements child, got %s",
			reqElem.Name)
	}

	// RequestedElements must contain exactly two wscn:Name children
	// with the correct QName text values in order.
	if len(reqElem.Children) != 2 {
		t.Fatalf("Expected 2 Name children, got %d",
			len(reqElem.Children))
	}

	expected := []ScannerRequestedElementNames{
		ScannerRequestedElementDescription,
		ScannerRequestedElementConfiguration,
	}
	for i, child := range reqElem.Children {
		if child.Name != NsWSCN+":Name" {
			t.Errorf("Child[%d]: expected element name %s, got %s",
				i, NsWSCN+":Name", child.Name)
		}
		if child.Text != expected[i].String() {
			t.Errorf("Child[%d]: expected text %s, got %s",
				i, expected[i], child.Text)
		}
	}
}
