// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// condition history entry tests

package wsscan

import (
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestConditionHistoryEntry_toXML(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	clearTime, _ := time.Parse(time.RFC3339, "2024-01-01T13:00:00Z")
	che := ConditionHistoryEntry{
		ClearTime: DateTime(clearTime),
		Component: PlatenComponent,
		Name:      CoverOpen,
		Severity:  Warning,
		Time:      DateTime(testTime),
	}

	expected := xmldoc.Element{
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
	}

	got := che.toXML("wscn:ConditionHistoryEntry")

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
		if child.Text != expected.Children[i].Text {
			t.Errorf("toXML().Children[%d].Text = %v, want %v",
				i, child.Text, expected.Children[i].Text)
		}
	}
}

func Test_decodeConditionHistoryEntry(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	clearTime, _ := time.Parse(time.RFC3339, "2024-01-01T13:00:00Z")
	tests := []struct {
		name     string
		element  xmldoc.Element
		expected ConditionHistoryEntry
		wantErr  bool
	}{
		{
			name: "Valid ConditionHistoryEntry",
			element: xmldoc.Element{
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
			expected: ConditionHistoryEntry{
				ClearTime: DateTime(clearTime),
				Component: PlatenComponent,
				Name:      CoverOpen,
				Severity:  Warning,
				Time:      DateTime(testTime),
			},
			wantErr: false,
		},
		{
			name: "Missing ClearTime",
			element: xmldoc.Element{
				Name: "wscn:ConditionHistoryEntry",
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
			expected: ConditionHistoryEntry{},
			wantErr:  true,
		},
		{
			name: "Missing Component",
			element: xmldoc.Element{
				Name: "wscn:ConditionHistoryEntry",
				Children: []xmldoc.Element{
					{
						Name: "wscn:ClearTime",
						Text: "2024-01-01T13:00:00Z",
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
			expected: ConditionHistoryEntry{},
			wantErr:  true,
		},
		{
			name: "Missing Name",
			element: xmldoc.Element{
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
						Name: "wscn:Severity",
						Text: "Warning",
					},
					{
						Name: "wscn:Time",
						Text: "2024-01-01T12:00:00Z",
					},
				},
			},
			expected: ConditionHistoryEntry{},
			wantErr:  true,
		},
		{
			name: "Missing Severity",
			element: xmldoc.Element{
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
						Name: "wscn:Time",
						Text: "2024-01-01T12:00:00Z",
					},
				},
			},
			expected: ConditionHistoryEntry{},
			wantErr:  true,
		},
		{
			name: "Missing Time",
			element: xmldoc.Element{
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
				},
			},
			expected: ConditionHistoryEntry{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeConditionHistoryEntry(tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeConditionHistoryEntry() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ClearTime != tt.expected.ClearTime {
					t.Errorf("decodeConditionHistoryEntry().ClearTime = %v, want %v",
						got.ClearTime, tt.expected.ClearTime)
				}
				if got.Component != tt.expected.Component {
					t.Errorf("decodeConditionHistoryEntry().Component = %v, want %v",
						got.Component, tt.expected.Component)
				}
				if got.Name != tt.expected.Name {
					t.Errorf("decodeConditionHistoryEntry().Name = %v, want %v",
						got.Name, tt.expected.Name)
				}
				if got.Severity != tt.expected.Severity {
					t.Errorf("decodeConditionHistoryEntry().Severity = %v, want %v",
						got.Severity, tt.expected.Severity)
				}
				if got.Time != tt.expected.Time {
					t.Errorf("decodeConditionHistoryEntry().Time = %v, want %v",
						got.Time, tt.expected.Time)
				}
			}
		})
	}
}
