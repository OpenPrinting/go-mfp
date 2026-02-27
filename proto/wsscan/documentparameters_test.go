// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// DocumentParameters tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDocumentParametersRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		dp   DocumentParameters
	}{
		{
			name: "empty DocumentParameters",
			dp:   DocumentParameters{},
		},
		{
			name: "with Format only",
			dp: DocumentParameters{
				Format: optional.New(Format(ValWithOptions[FormatValue]{
					Text: JFIF,
				})),
			},
		},
		{
			name: "with ImagesToTransfer only",
			dp: DocumentParameters{
				ImagesToTransfer: optional.New(ImagesToTransfer(ValWithOptions[int]{
					Text: 10,
				})),
			},
		},
		{
			name: "with InputSource only",
			dp: DocumentParameters{
				InputSource: optional.New(InputSource(ValWithOptions[InputSourceValue]{
					Text: InputSourcePlaten,
				})),
			},
		},
		{
			name: "with multiple elements",
			dp: DocumentParameters{
				Format: optional.New(Format(ValWithOptions[FormatValue]{
					Text: PNG,
				})),
				ImagesToTransfer: optional.New(ImagesToTransfer(ValWithOptions[int]{
					Text: 5,
				})),
				CompressionQualityFactor: optional.New(CompressionQualityFactor(ValWithOptions[int]{
					Text: 75,
				})),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.dp.toXML(NsWSCN + ":DocumentParameters")

			// Decode back
			decoded, err := decodeDocumentParameters(xml)
			if err != nil {
				t.Fatalf("decodeDocumentParameters() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.dp) {
				t.Errorf("Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.dp, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestDocumentParametersDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "invalid ImagesToTransfer",
			xml: xmldoc.Element{
				Name: NsWSCN + ":DocumentParameters",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":ImagesToTransfer",
						Text: "not-a-number",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeDocumentParameters(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeDocumentParameters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocumentParametersToXML(t *testing.T) {
	dp := DocumentParameters{
		Format: optional.New(Format(ValWithOptions[FormatValue]{
			Text: JFIF,
		})),
		ImagesToTransfer: optional.New(ImagesToTransfer(ValWithOptions[int]{
			Text: 10,
		})),
	}

	xml := dp.toXML(NsWSCN + ":DocumentParameters")

	// Verify element name
	if xml.Name != NsWSCN+":DocumentParameters" {
		t.Errorf("Expected name %s, got %s", NsWSCN+":DocumentParameters", xml.Name)
	}

	// Verify children exist
	if len(xml.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(xml.Children))
	}

	// Verify Format child
	foundFormat := false
	foundImagesToTransfer := false
	for _, child := range xml.Children {
		if child.Name == NsWSCN+":Format" {
			foundFormat = true
		}
		if child.Name == NsWSCN+":ImagesToTransfer" {
			foundImagesToTransfer = true
		}
	}

	if !foundFormat {
		t.Error("Format child element not found")
	}
	if !foundImagesToTransfer {
		t.Error("ImagesToTransfer child element not found")
	}
}
