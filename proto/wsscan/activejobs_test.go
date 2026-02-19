// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ActiveJobs tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestActiveJobsRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		aj   ActiveJobs
	}{
		{
			name: "empty ActiveJobs",
			aj:   ActiveJobs{},
		},
		{
			name: "ActiveJobs with only JobSummary",
			aj: ActiveJobs{
				JobSummary: []JobSummary{
					{
						JobID:                  100,
						JobName:                "Summary1",
						JobOriginatingUserName: "user1",
						JobState:               JobStatePending,
						ScansCompleted:         0,
					},
					{
						JobID:                  101,
						JobName:                "Summary2",
						JobOriginatingUserName: "user2",
						JobState:               JobStateProcessing,
						ScansCompleted:         3,
					},
				},
			},
		},
		{
			name: "ActiveJobs with only Job",
			aj: ActiveJobs{
				Job: []Job{
					{
						Documents: Documents{
							DocumentFinalParameters: DocumentParameters{
								Format: optional.New(Format(
									ValWithOptions[FormatValue]{
										Text: JFIF,
									})),
							},
						},
						JobStatus: JobStatus{
							JobID:          200,
							JobState:       JobStatePending,
							ScansCompleted: 0,
						},
						ScanTicket: ScanTicket{
							JobDescription: JobDescription{
								JobName:                "FullJob1",
								JobOriginatingUserName: "admin",
							},
						},
					},
				},
			},
		},
		{
			name: "ActiveJobs with both Job and JobSummary",
			aj: ActiveJobs{
				Job: []Job{
					{
						Documents: Documents{
							DocumentFinalParameters: DocumentParameters{
								Format: optional.New(Format(
									ValWithOptions[FormatValue]{
										Text: PNG,
									})),
							},
						},
						JobStatus: JobStatus{
							JobID:          300,
							JobState:       JobStateProcessing,
							ScansCompleted: 2,
						},
						ScanTicket: ScanTicket{
							JobDescription: JobDescription{
								JobName:                "MixedJob",
								JobOriginatingUserName: "system",
							},
						},
					},
				},
				JobSummary: []JobSummary{
					{
						JobID:                  301,
						JobName:                "MixedSummary",
						JobOriginatingUserName: "user",
						JobState:               JobStateCompleted,
						ScansCompleted:         5,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.aj.toXML(NsWSCN + ":ActiveJobs")

			// Decode back
			decoded, err := decodeActiveJobs(xml)
			if err != nil {
				t.Fatalf("decodeActiveJobs() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.aj) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.aj, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestActiveJobsDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "empty ActiveJobs",
			xml: xmldoc.Element{
				Name:     NsWSCN + ":ActiveJobs",
				Children: []xmldoc.Element{},
			},
			wantErr: false,
		},
		{
			name: "invalid Job",
			xml: xmldoc.Element{
				Name: NsWSCN + ":ActiveJobs",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":Job",
						Children: []xmldoc.Element{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid JobSummary",
			xml: xmldoc.Element{
				Name: NsWSCN + ":ActiveJobs",
				Children: []xmldoc.Element{
					{
						Name:     NsWSCN + ":JobSummary",
						Children: []xmldoc.Element{},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeActiveJobs(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeActiveJobs() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestActiveJobsToXML(t *testing.T) {
	aj := ActiveJobs{
		Job: []Job{
			{
				Documents: Documents{
					DocumentFinalParameters: DocumentParameters{
						Format: optional.New(Format(
							ValWithOptions[FormatValue]{
								Text: JFIF,
							})),
					},
				},
				JobStatus: JobStatus{
					JobID:          400,
					JobState:       JobStatePending,
					ScansCompleted: 0,
				},
				ScanTicket: ScanTicket{
					JobDescription: JobDescription{
						JobName:                "XMLTestJob",
						JobOriginatingUserName: "testuser",
					},
				},
			},
		},
		JobSummary: []JobSummary{
			{
				JobID:                  401,
				JobName:                "XMLTestSummary",
				JobOriginatingUserName: "testuser2",
				JobState:               JobStateCompleted,
				ScansCompleted:         10,
			},
		},
	}

	xml := aj.toXML(NsWSCN + ":ActiveJobs")

	// Verify element name
	if xml.Name != NsWSCN+":ActiveJobs" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":ActiveJobs", xml.Name)
	}

	// Verify children (1 Job + 1 JobSummary)
	if len(xml.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(xml.Children))
	}

	// Verify first child is Job
	if xml.Children[0].Name != NsWSCN+":Job" {
		t.Errorf("Expected Job child at index 0, got %s",
			xml.Children[0].Name)
	}

	// Verify second child is JobSummary
	if xml.Children[1].Name != NsWSCN+":JobSummary" {
		t.Errorf("Expected JobSummary child at index 1, got %s",
			xml.Children[1].Name)
	}
}
