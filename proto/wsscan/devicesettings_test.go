// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions

package wsscan

import (
	"reflect"
	"strings"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func createValidDeviceSettings() DeviceSettings {
	autoExp := BooleanElement("true")
	brightness := BooleanElement("true")
	contrast := BooleanElement("true")
	docSize := BooleanElement("true")

	return DeviceSettings{
		AutoExposureSupported:             autoExp,
		BrightnessSupported:               brightness,
		CompressionQualityFactorSupported: Range{MinValue: 1, MaxValue: 100},
		ContentTypesSupported:             []ContentTypeValue{Auto},
		ContrastSupported:                 contrast,
		DocumentSizeAutoDetectSupported:   docSize,
		FormatsSupported:                  []FormatValue{PNG},
		RotationsSupported:                []RotationValue{Rotation0},
		ScalingRangeSupported: ScalingRangeSupported{
			ScalingWidth:  Range{MinValue: 1, MaxValue: 1000},
			ScalingHeight: Range{MinValue: 1, MaxValue: 1000},
		},
	}
}

// Test for DeviceSettings XML round-trip
func TestDeviceSettings_XMLRoundTrip(t *testing.T) {
	ds := createValidDeviceSettings()
	elm := ds.toXML(NsWSCN + ":DeviceSettings")
	parsed, err := decodeDeviceSettings(elm)
	if err != nil {
		t.Fatalf("decodeDeviceSettings: input %+v, unexpected error: %v", ds, err)
	}

	// Compare individual fields since BooleanElements are not comparable with DeepEqual
	if parsed.AutoExposureSupported != ds.AutoExposureSupported {
		t.Errorf("AutoExposureSupported: expected %v, got %v",
			ds.AutoExposureSupported, parsed.AutoExposureSupported)
	}
	if parsed.BrightnessSupported != ds.BrightnessSupported {
		t.Errorf("BrightnessSupported: expected %v, got %v",
			ds.BrightnessSupported, parsed.BrightnessSupported)
	}
	if !reflect.DeepEqual(parsed.CompressionQualityFactorSupported, ds.CompressionQualityFactorSupported) {
		t.Errorf("CompressionQualityFactorSupported: expected %v, got %v",
			ds.CompressionQualityFactorSupported, parsed.CompressionQualityFactorSupported)
	}
	if !reflect.DeepEqual(parsed.ContentTypesSupported, ds.ContentTypesSupported) {
		t.Errorf("ContentTypesSupported: expected %v, got %v",
			ds.ContentTypesSupported, parsed.ContentTypesSupported)
	}
	if parsed.ContrastSupported != ds.ContrastSupported {
		t.Errorf("ContrastSupported: expected %v, got %v",
			ds.ContrastSupported, parsed.ContrastSupported)
	}
	if parsed.DocumentSizeAutoDetectSupported != ds.DocumentSizeAutoDetectSupported {
		t.Errorf("DocumentSizeAutoDetectSupported: expected %v, got %v",
			ds.DocumentSizeAutoDetectSupported, parsed.DocumentSizeAutoDetectSupported)
	}
	if !reflect.DeepEqual(parsed.FormatsSupported, ds.FormatsSupported) {
		t.Errorf("FormatsSupported: expected %v, got %v",
			ds.FormatsSupported, parsed.FormatsSupported)
	}
	if !reflect.DeepEqual(parsed.RotationsSupported, ds.RotationsSupported) {
		t.Errorf("RotationsSupported: expected %v, got %v",
			ds.RotationsSupported, parsed.RotationsSupported)
	}
	if !reflect.DeepEqual(parsed.ScalingRangeSupported, ds.ScalingRangeSupported) {
		t.Errorf("ScalingRangeSupported: expected %v, got %v",
			ds.ScalingRangeSupported, parsed.ScalingRangeSupported)
	}
}

func TestDeviceSettings_DecodeErrors(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() xmldoc.Element
		errContains string
	}{
		{
			name: "missing AutoExposureSupported",
			setup: func() xmldoc.Element {
				ds := createValidDeviceSettings()
				elm := ds.toXML(NsWSCN + ":DeviceSettings")
				// Remove AutoExposureSupported
				var children []xmldoc.Element
				for _, child := range elm.Children {
					if child.Name != NsWSCN+":AutoExposureSupported" {
						children = append(children, child)
					}
				}
				elm.Children = children
				return elm
			},
			errContains: "AutoExposureSupported",
		},
		{
			name: "invalid AutoExposureSupported value",
			setup: func() xmldoc.Element {
				ds := createValidDeviceSettings()
				ds.AutoExposureSupported = BooleanElement("maybe")
				return ds.toXML(NsWSCN + ":DeviceSettings")
			},
			errContains: "AutoExposureSupported",
		},
		{
			name: "missing ContentTypesSupported values",
			setup: func() xmldoc.Element {
				ds := createValidDeviceSettings()
				ds.ContentTypesSupported = nil
				return ds.toXML(NsWSCN + ":DeviceSettings")
			},
			errContains: "at least one ContentTypeValue is required",
		},
		{
			name: "missing FormatsSupported values",
			setup: func() xmldoc.Element {
				ds := createValidDeviceSettings()
				ds.FormatsSupported = nil
				return ds.toXML(NsWSCN + ":DeviceSettings")
			},
			errContains: "at least one FormatValue is required",
		},
		{
			name: "missing RotationsSupported values",
			setup: func() xmldoc.Element {
				ds := createValidDeviceSettings()
				ds.RotationsSupported = nil
				return ds.toXML(NsWSCN + ":DeviceSettings")
			},
			errContains: "at least one RotationValue is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			elm := tc.setup()
			_, err := decodeDeviceSettings(elm)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !containsSubstring(err.Error(), tc.errContains) {
				t.Errorf("expected error to contain %q, got %q",
					tc.errContains, err)
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}
