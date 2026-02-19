// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Document tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDocumentRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		doc  Document
	}{
		{
			name: "simple document",
			doc: Document{
				DocumentDescription: DocumentDescription{
					DocumentName: "MyDocument.pdf",
				},
			},
		},
		{
			name: "document with path",
			doc: Document{
				DocumentDescription: DocumentDescription{
					DocumentName: "/home/user/scans/scan001.jpg",
				},
			},
		},
		{
			name: "document with special characters",
			doc: Document{
				DocumentDescription: DocumentDescription{
					DocumentName: "Report (Final) - 2024.png",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.doc.toXML(NsWSCN + ":Document")

			// Decode back
			decoded, err := decodeDocument(xml)
			if err != nil {
				t.Fatalf("decodeDocument() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.doc) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.doc, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestDocumentDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing DocumentDescription",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":Document",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			name: "invalid DocumentDescription - missing DocumentName",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Document",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":DocumentDescription",
						Children: []xmldoc.Element{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DocumentDescription - empty DocumentName",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Document",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DocumentDescription",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":DocumentName",
								Text: "",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid Document",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Document",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DocumentDescription",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":DocumentName",
								Text: "ValidDocument.pdf",
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
			_, err := decodeDocument(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeDocument() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestDocumentToXML(t *testing.T) {
	doc := Document{
		DocumentDescription: DocumentDescription{
			DocumentName: "TestDoc.jpg",
		},
	}

	xml := doc.toXML(NsWSCN + ":Document")

	// Verify element name
	if xml.Name != NsWSCN+":Document" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":Document", xml.Name)
	}

	// Verify one child (DocumentDescription)
	if len(xml.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(xml.Children))
	}

	// Verify child is DocumentDescription
	if xml.Children[0].Name != NsWSCN+":DocumentDescription" {
		t.Errorf("Expected DocumentDescription child, got %s",
			xml.Children[0].Name)
	}

	// Verify DocumentDescription has DocumentName child
	if len(xml.Children[0].Children) != 1 {
		t.Errorf("Expected DocumentDescription to have 1 child, got %d",
			len(xml.Children[0].Children))
	}

	if xml.Children[0].Children[0].Name != NsWSCN+":DocumentName" {
		t.Errorf("Expected DocumentName child, got %s",
			xml.Children[0].Children[0].Name)
	}

	if xml.Children[0].Children[0].Text != "TestDoc.jpg" {
		t.Errorf("Expected text 'TestDoc.jpg', got %s",
			xml.Children[0].Children[0].Text)
	}
}
