// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// JobSummary tests

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestJobSummaryRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		js   JobSummary
	}{
		{
			name: "minimal JobSummary",
			js: JobSummary{
				JobID:                  12345,
				JobName:                "TestJob",
				JobOriginatingUserName: "user@example.com",
				JobState:               JobStatePending,
				ScansCompleted:         0,
			},
		},
		{
			name: "JobSummary with JobStateReasons",
			js: JobSummary{
				JobID:                  67890,
				JobName:                "ComplexJob",
				JobOriginatingUserName: "admin",
				JobState:               JobStateProcessing,
				JobStateReasons: []JobStateReason{
					JobScanning,
					JobTransferring,
				},
				ScansCompleted: 5,
			},
		},
		{
			name: "completed JobSummary",
			js: JobSummary{
				JobID:                  99999,
				JobName:                "CompletedJob",
				JobOriginatingUserName: "system",
				JobState:               JobStateCompleted,
				ScansCompleted:         10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to XML
			xml := tt.js.toXML(NsWSCN + ":JobSummary")

			// Decode back
			decoded, err := decodeJobSummary(xml)
			if err != nil {
				t.Fatalf("decodeJobSummary() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(decoded, tt.js) {
				t.Errorf(
					"Round trip failed:\nOriginal: %+v\nDecoded:  %+v\nXML: %s",
					tt.js, decoded, xml.EncodeString(nil))
			}
		})
	}
}

func TestJobSummaryDecodeError(t *testing.T) {
	tests := []struct {
		name    string
		xml     xmldoc.Element
		wantErr bool
	}{
		{
			name: "missing JobId",
			xml: xmldoc.Element{
				Name: NsWSCN + ":JobSummary",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":JobName", Text: "Test"},
					{Name: NsWSCN + ":JobOriginatingUserName", Text: "user"},
					{Name: NsWSCN + ":JobState", Text: "Pending"},
					{Name: NsWSCN + ":ScansCompleted", Text: "0"},
				},
			},
			wantErr: true,
		},
		{
			name: "missing JobName",
			xml: xmldoc.Element{
				Name: NsWSCN + ":JobSummary",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":JobId", Text: "1"},
					{Name: NsWSCN + ":JobOriginatingUserName", Text: "user"},
					{Name: NsWSCN + ":JobState", Text: "Pending"},
					{Name: NsWSCN + ":ScansCompleted", Text: "0"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty JobName",
			xml: xmldoc.Element{
				Name: NsWSCN + ":JobSummary",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":JobId", Text: "1"},
					{Name: NsWSCN + ":JobName", Text: ""},
					{Name: NsWSCN + ":JobOriginatingUserName", Text: "user"},
					{Name: NsWSCN + ":JobState", Text: "Pending"},
					{Name: NsWSCN + ":ScansCompleted", Text: "0"},
				},
			},
			wantErr: true,
		},
		{
			name: "missing JobOriginatingUserName",
			xml: xmldoc.Element{
				Name: NsWSCN + ":JobSummary",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":JobId", Text: "1"},
					{Name: NsWSCN + ":JobName", Text: "Test"},
					{Name: NsWSCN + ":JobState", Text: "Pending"},
					{Name: NsWSCN + ":ScansCompleted", Text: "0"},
				},
			},
			wantErr: true,
		},
		{
			name: "valid JobSummary",
			xml: xmldoc.Element{
				Name: NsWSCN + ":JobSummary",
				Children: []xmldoc.Element{
					{Name: NsWSCN + ":JobId", Text: "100"},
					{Name: NsWSCN + ":JobName", Text: "ValidJob"},
					{Name: NsWSCN + ":JobOriginatingUserName",
						Text: "validuser"},
					{Name: NsWSCN + ":JobState", Text: "Completed"},
					{Name: NsWSCN + ":ScansCompleted", Text: "5"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeJobSummary(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"decodeJobSummary() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestJobSummaryToXML(t *testing.T) {
	js := JobSummary{
		JobID:                  555,
		JobName:                "XMLTest",
		JobOriginatingUserName: "testuser",
		JobState:               JobStateProcessing,
		JobStateReasons: []JobStateReason{
			JobScanning,
		},
		ScansCompleted: 3,
	}

	xml := js.toXML(NsWSCN + ":JobSummary")

	// Verify element name
	if xml.Name != NsWSCN+":JobSummary" {
		t.Errorf("Expected name %s, got %s",
			NsWSCN+":JobSummary", xml.Name)
	}

	// Verify required children exist
	requiredChildren := map[string]bool{
		NsWSCN + ":JobId":                  false,
		NsWSCN + ":JobName":                false,
		NsWSCN + ":JobOriginatingUserName": false,
		NsWSCN + ":JobState":               false,
		NsWSCN + ":ScansCompleted":         false,
	}

	for _, child := range xml.Children {
		if _, exists := requiredChildren[child.Name]; exists {
			requiredChildren[child.Name] = true
		}
	}

	for name, found := range requiredChildren {
		if !found {
			t.Errorf("Required child %s not found", name)
		}
	}

	// Verify JobStateReasons is present
	foundJobStateReasons := false
	for _, child := range xml.Children {
		if child.Name == NsWSCN+":JobStateReasons" {
			foundJobStateReasons = true
			if len(child.Children) != 1 {
				t.Errorf("Expected 1 JobStateReason, got %d",
					len(child.Children))
			}
		}
	}
	if !foundJobStateReasons {
		t.Error("JobStateReasons not found in XML")
	}
}
