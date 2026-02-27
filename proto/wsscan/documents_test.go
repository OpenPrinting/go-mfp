// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Documents tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDocumentsRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		docs Documents
	}{
		{
			name: "Documents with no Document elements",
			docs: Documents{
				DocumentFinalParameters: DocumentParameters{
					Format: optional.New(Format(
						ValWithOptions[FormatValue]{
							Text: JFIF,
						})),
				},
			},
		},
		{
			name: "Documents with one Document",
			docs: Documents{
				DocumentFinalParameters: DocumentParameters{
					Format: optional.New(Format(
						ValWithOptions[FormatValue]{
							Text: PNG,
						})),
					ImagesToTransfer: optional.New(ImagesToTransfer(
						ValWithOptions[int]{
							Text: 5,
						})),
				},
				Document: []Document{
					{
						DocumentDescription: DocumentDescription{
							DocumentName: "scan001.png",
						},
					},
				},
			},
		},
		{
			name: "Documents with multiple Documents",
			docs: Documents{
				DocumentFinalParameters: DocumentParameters{
					Format: optional.New(Format(
						ValWithOptions[FormatValue]{
							Text: JFIF,
						})),
					ImagesToTransfer: optional.New(ImagesToTransfer(
						ValWithOptions[int]{
							Text: 10,
						})),
					InputSource: optional.New(InputSource(
						ValWithOptions[InputSourceValue]{
							Text: InputSourcePlaten,
						})),
				},
				Document: []Document{
					{
						DocumentDescription: DocumentDescription{
							DocumentName: "page1.jpg",
						},
					},
					{
						DocumentDescription: DocumentDescription{
							DocumentName: "page2.jpg",
						},
					},
					{
						DocumentDescription: DocumentDescription{
							DocumentName: "page3.jpg",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.docs.toXML(NsWSCN + ":Documents")

			// Decode back
			decoded, err := decodeDocuments(xml)
			if err != nil {
				t.Fatalf("decodeDocuments() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.docs) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.docs, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestDocumentsDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing DocumentFinalParameters",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":Documents",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			name: "invalid DocumentFinalParameters",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Documents",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DocumentFinalParameters",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":ImagesToTransfer",
								Text: "not-a-number",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid Documents with no Document elements",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Documents",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":DocumentFinalParameters",
						Children: []xmldoc.Element{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Documents with Document elements",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Documents",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":DocumentFinalParameters",
						Children: []xmldoc.Element{},
					},
					{
						Name: NsWSCN + ":Document",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":DocumentDescription",
								Children: []xmldoc.Element{
									{
										Name: NsWSCN + ":DocumentName",
										Text: "doc1.pdf",
									},
								},
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
			_, err := decodeDocuments(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeDocuments() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestDocumentsToXML(t *testing.T) {
	docs := Documents{
		DocumentFinalParameters: DocumentParameters{
			Format: optional.New(Format(
				ValWithOptions[FormatValue]{
					Text: JFIF,
				})),
		},
		Document: []Document{
			{
				DocumentDescription: DocumentDescription{
					DocumentName: "doc1.jpg",
				},
			},
			{
				DocumentDescription: DocumentDescription{
					DocumentName: "doc2.jpg",
				},
			},
		},
	}

	xml := docs.toXML(NsWSCN + ":Documents")

	// Verify element name
	if xml.Name != NsWSCN+":Documents" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":Documents", xml.Name)
	}

	// Verify children (DocumentFinalParameters + 2 Documents)
	if len(xml.Children) != 3 {
		t.Errorf("Expected 3 children, got %d", len(xml.Children))
	}

	// Verify first child is DocumentFinalParameters
	if xml.Children[0].Name != NsWSCN+":DocumentFinalParameters" {
		t.Errorf("Expected DocumentFinalParameters child, got %s",
			xml.Children[0].Name)
	}

	// Verify remaining children are Documents
	for i := 1; i < len(xml.Children); i++ {
		if xml.Children[i].Name != NsWSCN+":Document" {
			t.Errorf("Expected Document child at index %d, got %s",
				i, xml.Children[i].Name)
		}
	}
}
