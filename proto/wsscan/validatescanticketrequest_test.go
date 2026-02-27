// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ValidateScanTicketRequest tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestValidateScanTicketRequestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		vstr ValidateScanTicketRequest
	}{
		{
			name: "minimal ValidateScanTicketRequest",
			vstr: ValidateScanTicketRequest{
				ScanTicket: ScanTicket{
					JobDescription: JobDescription{
						JobName:                "TestValidation",
						JobOriginatingUserName: "user@example.com",
					},
				},
			},
		},
		{
			name: "ValidateScanTicketRequest with DocumentParameters",
			vstr: ValidateScanTicketRequest{
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
						InputSource: optional.New(
							InputSource(ValWithOptions[InputSourceValue]{
								Text: InputSourcePlaten,
							}),
						),
					}),
					JobDescription: JobDescription{
						JobName:                "ComplexValidation",
						JobOriginatingUserName: "admin",
						JobInformation:         optional.New("Validating scan settings"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.vstr.toXML(NsWSCN + ":ValidateScanTicketRequest")

			// Decode back
			decoded, err := decodeValidateScanTicketRequest(xml)
			if err != nil {
				t.Fatalf("decodeValidateScanTicketRequest() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.vstr) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.vstr, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestValidateScanTicketRequestDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing ScanTicket",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":ValidateScanTicketRequest",
				Children: []xmldoc.Element{},
			},
			wantErr: true,
		},
		{
			name: "invalid ScanTicket - missing JobDescription",
			xml: xmldoc.Element{
				Name: NsWSCN + ":ValidateScanTicketRequest",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":ScanTicket",
						Children: []xmldoc.Element{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid ValidateScanTicketRequest",
			xml: xmldoc.Element{
				Name: NsWSCN + ":ValidateScanTicketRequest",
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
			_, err := decodeValidateScanTicketRequest(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeValidateScanTicketRequest() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateScanTicketRequestToXML(t *testing.T) {
	vstr := ValidateScanTicketRequest{
		ScanTicket: ScanTicket{
			JobDescription: JobDescription{
				JobName:                "MyValidation",
				JobOriginatingUserName: "testuser",
			},
		},
	}

	xml := vstr.toXML(NsWSCN + ":ValidateScanTicketRequest")

	// Verify element name
	if xml.Name != NsWSCN+":ValidateScanTicketRequest" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":ValidateScanTicketRequest", xml.Name)
	}

	// Verify one child (ScanTicket)
	if len(xml.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(xml.Children))
	}

	// Verify child is ScanTicket
	if xml.Children[0].Name != NsWSCN+":ScanTicket" {
		t.Errorf("Expected ScanTicket child, got %s", xml.Children[0].Name)
	}
}
