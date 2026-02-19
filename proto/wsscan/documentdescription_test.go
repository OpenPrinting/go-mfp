// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// DocumentDescription tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDocumentDescriptionRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		dd   DocumentDescription
	}{
		{
			name: "simple document name",
			dd: DocumentDescription{
				DocumentName: "MyDocument.pdf",
			},
		},
		{
			name: "document with path",
			dd: DocumentDescription{
				DocumentName: "/path/to/document.jpg",
			},
		},
		{
			name: "document with special characters",
			dd: DocumentDescription{
				DocumentName: "Document (Copy) #2.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.dd.toXML(NsWSCN + ":DocumentDescription")

			// Decode back
			decoded, err := decodeDocumentDescription(xml)
			if err != nil {
				t.Fatalf("decodeDocumentDescription() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.dd) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.dd, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestDocumentDescriptionDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing DocumentName",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":DocumentDescription",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			name: "empty DocumentName",
			xml: xmldoc.Element{
				Name: NsWSCN + ":DocumentDescription",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DocumentName",
						Text: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid DocumentDescription",
			xml: xmldoc.Element{
				Name: NsWSCN + ":DocumentDescription",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DocumentName",
						Text: "ValidDocument.pdf",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeDocumentDescription(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeDocumentDescription() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestDocumentDescriptionToXML(t *testing.T) {
	dd := DocumentDescription{
		DocumentName: "TestDoc.jpg",
	}

	xml := dd.toXML(NsWSCN + ":DocumentDescription")

	// Verify element name
	if xml.Name != NsWSCN+":DocumentDescription" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":DocumentDescription", xml.Name)
	}

	// Verify one child (DocumentName)
	if len(xml.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(xml.Children))
	}

	// Verify child is DocumentName
	if xml.Children[0].Name != NsWSCN+":DocumentName" {
		t.Errorf("Expected DocumentName child, got %s",
			xml.Children[0].Name)
	}

	// Verify text value
	if xml.Children[0].Text != "TestDoc.jpg" {
		t.Errorf("Expected text 'TestDoc.jpg', got %s",
			xml.Children[0].Text)
	}
}
