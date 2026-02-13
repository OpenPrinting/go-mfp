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
				DestinationToken: "dest-token-123",
				ScanIdentifier:   "scan-id-456",
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
				DestinationToken: "dest-uuid-001",
				ScanIdentifier:   "scan-uuid-002",
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
			name: "CreateScanJobRequest with UUID-like identifiers",
			csjr: CreateScanJobRequest{
				DestinationToken: "550e8400-e29b-41d4-a716-446655440000",
				ScanIdentifier:   "660e8400-e29b-41d4-a716-446655440001",
				ScanTicket: ScanTicket{
					JobDescription: JobDescription{
						JobName:                "UUIDScan",
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
			name: "missing DestinationToken",
			xml: xmldoc.Element{
				Name: NsWSCN + ":CreateScanJobRequest",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":ScanIdentifier",
						Text: "scan-123",
					},
					{
						Name: NsWSCN + ":ScanTicket",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":JobDescription",
								Children: []xmldoc.Element{
									{
										Name: NsWSCN + ":JobName",
										Text: "Test",
									},
									{
										Name: NsWSCN + ":JobOriginatingUserName",
										Text: "user",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing ScanIdentifier",
			xml: xmldoc.Element{
				Name: NsWSCN + ":CreateScanJobRequest",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DestinationToken",
						Text: "dest-123",
					},
					{
						Name: NsWSCN + ":ScanTicket",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":JobDescription",
								Children: []xmldoc.Element{
									{
										Name: NsWSCN + ":JobName",
										Text: "Test",
									},
									{
										Name: NsWSCN + ":JobOriginatingUserName",
										Text: "user",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing ScanTicket",
			xml: xmldoc.Element{
				Name: NsWSCN + ":CreateScanJobRequest",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DestinationToken",
						Text: "dest-123",
					},
					{
						Name: NsWSCN + ":ScanIdentifier",
						Text: "scan-123",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid CreateScanJobRequest",
			xml: xmldoc.Element{
				Name: NsWSCN + ":CreateScanJobRequest",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":DestinationToken",
						Text: "dest-valid",
					},
					{
						Name: NsWSCN + ":ScanIdentifier",
						Text: "scan-valid",
					},
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
		DestinationToken: "token-abc",
		ScanIdentifier:   "id-xyz",
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

	// Verify all three children exist
	if len(xml.Children) != 3 {
		t.Errorf("Expected 3 children, got %d", len(xml.Children))
	}

	// Verify child element names and order
	expectedChildren := []string{
		NsWSCN + ":DestinationToken",
		NsWSCN + ":ScanIdentifier",
		NsWSCN + ":ScanTicket",
	}

	for i, expected := range expectedChildren {
		if i >= len(xml.Children) {
			t.Errorf("Missing child at index %d", i)
			continue
		}
		if xml.Children[i].Name != expected {
			t.Errorf("Child %d: expected %s, got %s",
				i, expected, xml.Children[i].Name)
		}
	}

	// Verify text values for string fields
	if xml.Children[0].Text != "token-abc" {
		t.Errorf("DestinationToken: expected 'token-abc', got %s",
			xml.Children[0].Text)
	}
	if xml.Children[1].Text != "id-xyz" {
		t.Errorf("ScanIdentifier: expected 'id-xyz', got %s",
			xml.Children[1].Text)
	}
}
