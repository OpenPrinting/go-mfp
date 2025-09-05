// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// device condition tests

package wsscan

import (
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDeviceCondition_toXML(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	dc := DeviceCondition{
		Component: PlatenComponent,
		Name:      CoverOpen,
		Severity:  Warning,
		Time:      DateTime(testTime),
	}

	expected := xmldoc.Element{
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
	}

	got := dc.toXML("wscn:DeviceCondition")

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

func Test_decodeDeviceCondition(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	tests := []struct {
		name     string
		element  xmldoc.Element
		expected DeviceCondition
		wantErr  bool
	}{
		{
			name: "Valid DeviceCondition",
			element: xmldoc.Element{
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
			expected: DeviceCondition{
				Component: PlatenComponent,
				Name:      CoverOpen,
				Severity:  Warning,
				Time:      DateTime(testTime),
			},
			wantErr: false,
		},
		{
			name: "Missing Component",
			element: xmldoc.Element{
				Name: "wscn:DeviceCondition",
				Children: []xmldoc.Element{
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
			expected: DeviceCondition{},
			wantErr:  true,
		},
		{
			name: "Missing Name",
			element: xmldoc.Element{
				Name: "wscn:DeviceCondition",
				Children: []xmldoc.Element{
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
			expected: DeviceCondition{},
			wantErr:  true,
		},
		{
			name: "Missing Severity",
			element: xmldoc.Element{
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
						Name: "wscn:Time",
						Text: "2024-01-01T12:00:00Z",
					},
				},
			},
			expected: DeviceCondition{},
			wantErr:  true,
		},
		{
			name: "Missing Time",
			element: xmldoc.Element{
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
				},
			},
			expected: DeviceCondition{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeDeviceCondition(tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeDeviceCondition() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Component != tt.expected.Component {
					t.Errorf("decodeDeviceCondition().Component = %v, want %v",
						got.Component, tt.expected.Component)
				}
				if got.Name != tt.expected.Name {
					t.Errorf("decodeDeviceCondition().Name = %v, want %v",
						got.Name, tt.expected.Name)
				}
				if got.Severity != tt.expected.Severity {
					t.Errorf("decodeDeviceCondition().Severity = %v, want %v",
						got.Severity, tt.expected.Severity)
				}
				if got.Time != tt.expected.Time {
					t.Errorf("decodeDeviceCondition().Time = %v, want %v",
						got.Time, tt.expected.Time)
				}
			}
		})
	}
}
