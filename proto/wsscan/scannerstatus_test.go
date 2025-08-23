// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// scanner status tests

package wsscan

import (
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestScannerStatus_toXML(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	clearTime, _ := time.Parse(time.RFC3339, "2024-01-01T13:00:00Z")

	ss := ScannerStatus{
		ActiveConditions: []DeviceCondition{
			{
				Component: PlatenComponent,
				Name:      CoverOpen,
				Severity:  Warning,
				Time:      DateTime(testTime),
			},
		},
		ConditionHistory: []ConditionHistoryEntry{
			{
				ClearTime: DateTime(clearTime),
				Component: PlatenComponent,
				Name:      CoverOpen,
				Severity:  Warning,
				Time:      DateTime(testTime),
			},
		},
		ScannerCurrentTime:  DateTime(testTime),
		ScannerState:        Idle,
		ScannerStateReasons: []ScannerStateReason{StateNone},
	}

	expected := xmldoc.Element{
		Name: "wscn:ScannerStatus",
		Children: []xmldoc.Element{
			{
				Name: "wscn:ActiveConditions",
				Children: []xmldoc.Element{
					{
						Name: "wscn:DeviceCondition",
						Children: []xmldoc.Element{
							{
								Name: "wscn:Component",
								Text: "Platen",
							},
							{
								Name: "wscn:Name",
								Text: "CoverOpen",
							},
							{
								Name: "wscn:Severity",
								Text: "Warning",
							},
							{
								Name: "wscn:Time",
								Text: "2024-01-01T12:00:00Z",
							},
						},
					},
				},
			},
			{
				Name: "wscn:ConditionHistory",
				Children: []xmldoc.Element{
					{
						Name: "wscn:ConditionHistoryEntry",
						Children: []xmldoc.Element{
							{
								Name: "wscn:ClearTime",
								Text: "2024-01-01T13:00:00Z",
							},
							{
								Name: "wscn:Component",
								Text: "Platen",
							},
							{
								Name: "wscn:Name",
								Text: "CoverOpen",
							},
							{
								Name: "wscn:Severity",
								Text: "Warning",
							},
							{
								Name: "wscn:Time",
								Text: "2024-01-01T12:00:00Z",
							},
						},
					},
				},
			},
			{
				Name: "wscn:ScannerCurrentTime",
				Text: "2024-01-01T12:00:00Z",
			},
			{
				Name: "wscn:ScannerState",
				Text: "Idle",
			},
			{
				Name: "wscn:ScannerStateReasons",
				Children: []xmldoc.Element{
					{
						Name: "wscn:ScannerStateReason",
						Text: "None",
					},
				},
			},
		},
	}

	got := ss.toXML("wscn:ScannerStatus")

	if got.Name != expected.Name {
		t.Errorf("toXML().Name = %v, want %v", got.Name, expected.Name)
	}

	if len(got.Children) != len(expected.Children) {
		t.Errorf("toXML().Children length = %v, want %v", len(got.Children),
			len(expected.Children))
		return
	}

	for i, child := range got.Children {
		if child.Name != expected.Children[i].Name {
			t.Errorf("toXML().Children[%d].Name = %v, want %v",
				i, child.Name, expected.Children[i].Name)
		}
	}
}

func Test_decodeScannerStatus(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")

	tests := []struct {
		name     string
		element  xmldoc.Element
		expected ScannerStatus
		wantErr  bool
	}{
		{
			name: "Valid ScannerStatus with all fields",
			element: xmldoc.Element{
				Name: "wscn:ScannerStatus",
				Children: []xmldoc.Element{
					{
						Name: "wscn:ActiveConditions",
						Children: []xmldoc.Element{
							{
								Name: "wscn:DeviceCondition",
								Children: []xmldoc.Element{
									{
										Name: "wscn:Component",
										Text: "Platen",
									},
									{
										Name: "wscn:Name",
										Text: "CoverOpen",
									},
									{
										Name: "wscn:Severity",
										Text: "Warning",
									},
									{
										Name: "wscn:Time",
										Text: "2024-01-01T12:00:00Z",
									},
								},
							},
						},
					},
					{
						Name: "wscn:ScannerCurrentTime",
						Text: "2024-01-01T12:00:00Z",
					},
					{
						Name: "wscn:ScannerState",
						Text: "Idle",
					},
					{
						Name: "wscn:ScannerStateReasons",
						Children: []xmldoc.Element{
							{
								Name: "wscn:ScannerStateReason",
								Text: "None",
							},
						},
					},
				},
			},
			expected: ScannerStatus{
				ActiveConditions: []DeviceCondition{
					{
						Component: PlatenComponent,
						Name:      CoverOpen,
						Severity:  Warning,
						Time:      DateTime(testTime),
					},
				},
				ScannerCurrentTime: DateTime(testTime),
				ScannerState:       Idle,
			},
			wantErr: false,
		},
		{
			name: "Missing ScannerCurrentTime",
			element: xmldoc.Element{
				Name: "wscn:ScannerStatus",
				Children: []xmldoc.Element{
					{
						Name: "wscn:ScannerState",
						Text: "Idle",
					},
					{
						Name: "wscn:ScannerStateReasons",
						Children: []xmldoc.Element{
							{
								Name: "wscn:ScannerStateReason",
								Text: "None",
							},
						},
					},
				},
			},
			expected: ScannerStatus{},
			wantErr:  true,
		},
		{
			name: "Missing ScannerState",
			element: xmldoc.Element{
				Name: "wscn:ScannerStatus",
				Children: []xmldoc.Element{
					{
						Name: "wscn:ScannerCurrentTime",
						Text: "2024-01-01T12:00:00Z",
					},
					{
						Name: "wscn:ScannerStateReasons",
						Children: []xmldoc.Element{
							{
								Name: "wscn:ScannerStateReason",
								Text: "None",
							},
						},
					},
				},
			},
			expected: ScannerStatus{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeScannerStatus(tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeScannerStatus() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got.ActiveConditions) != len(tt.expected.ActiveConditions) {
					t.Errorf("decodeScannerStatus().ActiveConditions length = %v, want %v",
						len(got.ActiveConditions), len(tt.expected.ActiveConditions))
				}
				if got.ScannerCurrentTime != tt.expected.ScannerCurrentTime {
					t.Errorf("decodeScannerStatus().ScannerCurrentTime = %v, want %v",
						got.ScannerCurrentTime, tt.expected.ScannerCurrentTime)
				}
				if got.ScannerState != tt.expected.ScannerState {
					t.Errorf("decodeScannerStatus().ScannerState = %v, want %v",
						got.ScannerState, tt.expected.ScannerState)
				}
			}
		})
	}
}
