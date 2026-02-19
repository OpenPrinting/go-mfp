// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ScanTicket tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestScanTicketRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		st   ScanTicket
	}{
		{
			name: "minimal ScanTicket with JobDescription only",
			st: ScanTicket{
				JobDescription: JobDescription{
					JobName:                "TestJob",
					JobOriginatingUserName: "testuser",
				},
			},
		},
		{
			name: "ScanTicket with DocumentParameters",
			st: ScanTicket{
				DocumentParameters: optional.New(DocumentParameters{
					Format: optional.New(Format(ValWithOptions[FormatValue]{
						Text: JFIF,
					})),
					ImagesToTransfer: optional.New(ImagesToTransfer(ValWithOptions[int]{
						Text: 10,
					})),
				}),
				JobDescription: JobDescription{
					JobName:                "ScanJob",
					JobOriginatingUserName: "admin",
					JobInformation:         optional.New("Test scan job"),
				},
			},
		},
		{
			name: "ScanTicket with full DocumentParameters",
			st: ScanTicket{
				DocumentParameters: optional.New(DocumentParameters{
					Format: optional.New(Format(
						ValWithOptions[FormatValue]{
							Text: PNG,
						})),
					ImagesToTransfer: optional.New(ImagesToTransfer(
						ValWithOptions[int]{
							Text: 5,
						})),
					InputSource: optional.New(InputSource(
						ValWithOptions[InputSourceValue]{
							Text: InputSourcePlaten,
						})),
					CompressionQualityFactor: optional.New(
						CompressionQualityFactor(ValWithOptions[int]{
							Text: 85,
						})),
				}),
				JobDescription: JobDescription{
					JobName:                "FullScanJob",
					JobOriginatingUserName: "poweruser",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.st.toXML(NsWSCN + ":ScanTicket")

			// Decode back
			decoded, err := decodeScanTicket(xml)
			if err != nil {
				t.Fatalf("decodeScanTicket() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.st) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.st, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestScanTicketDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing required JobDescription",
			xml: xmldoc.Element{
				Name: NsWSCN + ":ScanTicket",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DocumentParameters",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":Format",
								Text: "jfif",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid minimal ScanTicket",
			xml: xmldoc.Element{
				Name: NsWSCN + ":ScanTicket",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":JobDescription",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":JobName",
								Text: "TestJob",
							},
							{
								Name: NsWSCN + ":JobOriginatingUserName",
								Text: "testuser",
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
			_, err := decodeScanTicket(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeScanTicket() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestScanTicketToXML(t *testing.T) {
	st := ScanTicket{
		DocumentParameters: optional.New(DocumentParameters{
			Format: optional.New(Format(ValWithOptions[FormatValue]{
				Text: JFIF,
			})),
		}),
		JobDescription: JobDescription{
			JobName:                "TestJob",
			JobOriginatingUserName: "testuser",
		},
	}

	xml := st.toXML(NsWSCN + ":ScanTicket")

	// Verify element name
	if xml.Name != NsWSCN+":ScanTicket" {
		t.Errorf("Expected name %s, got %s", NsWSCN+":ScanTicket", xml.Name)
	}

	// Verify children exist (DocumentParameters + JobDescription)
	if len(xml.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(xml.Children))
	}

	// Verify child element names
	foundDocParams := false
	foundJobDesc := false
	for _, child := range xml.Children {
		if child.Name == NsWSCN+":DocumentParameters" {
			foundDocParams = true
		}
		if child.Name == NsWSCN+":JobDescription" {
			foundJobDesc = true
		}
	}

	if !foundDocParams {
		t.Error("DocumentParameters child element not found")
	}
	if !foundJobDesc {
		t.Error("JobDescription child element not found")
	}
}

func TestScanTicketWithoutDocumentParameters(t *testing.T) {
	st := ScanTicket{
		JobDescription: JobDescription{
			JobName:                "MinimalJob",
			JobOriginatingUserName: "user",
		},
	}

	xml := st.toXML(NsWSCN + ":ScanTicket")

	// Should only have JobDescription child
	if len(xml.Children) != 1 {
		t.Errorf("Expected 1 child (JobDescription only), got %d",
			len(xml.Children))
	}

	if xml.Children[0].Name != NsWSCN+":JobDescription" {
		t.Errorf("Expected JobDescription child, got %s",
			xml.Children[0].Name)
	}

	// Round trip test
	decoded, err := decodeScanTicket(xml)
	if err != nil {
		t.Fatalf("decodeScanTicket() error = %v", err)
	}

	if !reflect.DeepEqual(decoded, st) {
		t.Errorf("Round trip failed:\nOriginal: %+v\nDecoded:  %+v",
			st, decoded)
	}
}
