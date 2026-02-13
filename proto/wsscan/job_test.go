// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Job tests

package wsscan

import (
	"reflect"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestJobRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		job  Job
	}{
		{
			name: "minimal Job",
			job: Job{
				Documents: Documents{
					DocumentFinalParameters: DocumentParameters{
						Format: optional.New(Format(
							ValWithOptions[FormatValue]{
								Text: JFIF,
							})),
					},
				},
				JobStatus: JobStatus{
					JobID:          12345,
					JobState:       JobStatePending,
					ScansCompleted: 0,
				},
				ScanTicket: ScanTicket{
					JobDescription: JobDescription{
						JobName:                "TestJob",
						JobOriginatingUserName: "user@example.com",
					},
				},
			},
		},
		{
			name: "complete Job with Documents",
			job: Job{
				Documents: Documents{
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
								DocumentName: "page1.png",
							},
						},
						{
							DocumentDescription: DocumentDescription{
								DocumentName: "page2.png",
							},
						},
					},
				},
				JobStatus: JobStatus{
					JobID:          67890,
					JobState:       JobStateCompleted,
					ScansCompleted: 2,
					JobCreatedTime: optional.New(
						time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)),
					JobCompletedTime: optional.New(
						time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)),
				},
				ScanTicket: ScanTicket{
					DocumentParameters: optional.New(DocumentParameters{
						Format: optional.New(Format(
							ValWithOptions[FormatValue]{
								Text: PNG,
							})),
					}),
					JobDescription: JobDescription{
						JobName:                "CompleteJob",
						JobOriginatingUserName: "admin",
						JobInformation:         optional.New("Full scan job"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.job.toXML(NsWSCN + ":Job")

			// Decode back
			decoded, err := decodeJob(xml)
			if err != nil {
				t.Fatalf("decodeJob() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.job) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.job, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestJobDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing Documents",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Job",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":JobStatus",
						Children: []xmldoc.Element{
							{Name: NsWSCN + ":JobId", Text: "1"},
							{Name: NsWSCN + ":JobState", Text: "Pending"},
							{Name: NsWSCN + ":ScansCompleted", Text: "0"},
						},
					},
					{
						Name: NsWSCN + ":ScanTicket",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":JobDescription",
								Children: []xmldoc.Element{
									{Name: NsWSCN + ":JobName", Text: "Test"},
									{Name: NsWSCN + ":JobOriginatingUserName",
										Text: "user"},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing JobStatus",
			xml: xmldoc.Element{
				Name: NsWSCN + ":Job",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":Documents",
						Children: []xmldoc.Element{
							{Name: NsWSCN + ":DocumentFinalParameters"},
						},
					},
					{
						Name: NsWSCN + ":ScanTicket",
						Children: []xmldoc.Element{
							{
								Name: NsWSCN + ":JobDescription",
								Children: []xmldoc.Element{
									{Name: NsWSCN + ":JobName", Text: "Test"},
									{Name: NsWSCN + ":JobOriginatingUserName",
										Text: "user"},
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
				Name: NsWSCN + ":Job",
				Children: []xmldoc.Element{
					{
						Name: NsWSCN + ":Documents",
						Children: []xmldoc.Element{
							{Name: NsWSCN + ":DocumentFinalParameters"},
						},
					},
					{
						Name: NsWSCN + ":JobStatus",
						Children: []xmldoc.Element{
							{Name: NsWSCN + ":JobId", Text: "1"},
							{Name: NsWSCN + ":JobState", Text: "Pending"},
							{Name: NsWSCN + ":ScansCompleted", Text: "0"},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeJob(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeJob() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestJobToXML(t *testing.T) {
	job := Job{
		Documents: Documents{
			DocumentFinalParameters: DocumentParameters{
				Format: optional.New(Format(
					ValWithOptions[FormatValue]{
						Text: JFIF,
					})),
			},
		},
		JobStatus: JobStatus{
			JobID:          100,
			JobState:       JobStateProcessing,
			ScansCompleted: 1,
		},
		ScanTicket: ScanTicket{
			JobDescription: JobDescription{
				JobName:                "XMLTest",
				JobOriginatingUserName: "testuser",
			},
		},
	}

	xml := job.toXML(NsWSCN + ":Job")

	// Verify element name
	if xml.Name != NsWSCN+":Job" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":Job", xml.Name)
	}

	// Verify all three children exist
	if len(xml.Children) != 3 {
		t.Errorf("Expected 3 children, got %d", len(xml.Children))
	}

	// Verify child element names and order
	expectedChildren := []string{
		NsWSCN + ":Documents",
		NsWSCN + ":JobStatus",
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
}
