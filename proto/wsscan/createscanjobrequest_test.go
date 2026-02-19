// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// CreateScanJobRequest tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestCreateScanJobRequestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		csjr CreateScanJobRequest
	}{
		{
			name: "minimal CreateScanJobRequest",
			csjr: CreateScanJobRequest{
				ScanTicket: ScanTicket{
					JobDescription: JobDescription{
						JobName:                "TestScan",
						JobOriginatingUserName: "user@example.com",
					},
				},
			},
		},
		{
			name: "CreateScanJobRequest with DocumentParameters",
			csjr: CreateScanJobRequest{
				ScanTicket: ScanTicket{
					DocumentParameters: optional.New(DocumentParameters{
						Format: optional.New(Format(ValWithOptions[FormatValue]{
							Text: JFIF,
						})),
						ImagesToTransfer: optional.New(
							ImagesToTransfer(ValWithOptions[int]{
								Text: 5,
							}),
						),
					}),
					JobDescription: JobDescription{
						JobName:                "DocumentScan",
						JobOriginatingUserName: "admin",
						JobInformation:         optional.New("Scanning documents"),
					},
				},
			},
		},
		{
			name: "CreateScanJobRequest with full JobDescription",
			csjr: CreateScanJobRequest{
				ScanTicket: ScanTicket{
					JobDescription: JobDescription{
						JobName:                "FullScan",
						JobOriginatingUserName: "system",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.csjr.toXML(NsWSCN + ":CreateScanJobRequest")

			// Decode back
			decoded, err := decodeCreateScanJobRequest(xml)
			if err != nil {
				t.Fatalf("decodeCreateScanJobRequest() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.csjr) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.csjr, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestCreateScanJobRequestDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing ScanTicket",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":CreateScanJobRequest",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			name: "valid CreateScanJobRequest",
			xml: xmldoc.Element{
				Name: NsWSCN + ":CreateScanJobRequest",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":ScanTicket",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":JobDescription",
								Children: []xmldoc.Element{
									{
										Name: NsWSCN + ":JobName",
										Text: "ValidJob",
									},
									{
										Name: NsWSCN + ":JobOriginatingUserName",
										Text: "validuser",
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
			_, err := decodeCreateScanJobRequest(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCreateScanJobRequest() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestCreateScanJobRequestToXML(t *testing.T) {
	csjr := CreateScanJobRequest{
		ScanTicket: ScanTicket{
			JobDescription: JobDescription{
				JobName:                "MyJob",
				JobOriginatingUserName: "testuser",
			},
		},
	}

	xml := csjr.toXML(NsWSCN + ":CreateScanJobRequest")

	// Verify element name
	if xml.Name != NsWSCN+":CreateScanJobRequest" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":CreateScanJobRequest", xml.Name)
	}

	// Host-initiated mode: only ScanTicket child
	if len(xml.Children) != 1 {
		t.Errorf("Expected 1 child (ScanTicket), got %d", len(xml.Children))
	}
	if len(xml.Children) > 0 && xml.Children[0].Name != NsWSCN+":ScanTicket" {
		t.Errorf("Expected child %s, got %s",
			NsWSCN+":ScanTicket", xml.Children[0].Name)
	}
}
