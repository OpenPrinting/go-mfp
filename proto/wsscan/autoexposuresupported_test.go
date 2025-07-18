// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla
// See LICENSE for license terms and conditions
//
// Test for auto exposure supported

package wsscan

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestAutoExposureSupported(t *testing.T) {
	cases := []struct {
		xmlValue      string
		expectedBool  bool
		expectedValid bool
	}{
		{"1", true, true},
		{"0", false, true},
		{"true", true, true},
		{"false", false, true},
		{" TRUE ", true, true},
		{" False ", false, true},
		{"maybe", false, false},
		{"", false, false},
	}

	for _, c := range cases {
		aes := AutoExposureSupported(c.xmlValue)
		if aes.IsValid() != c.expectedValid {
			t.Errorf(
				"IsValid: input %q, expected %v, got %v",
				c.xmlValue, c.expectedValid, aes.IsValid(),
			)
		}
		b, err := aes.Bool()
		if c.expectedValid {
			if err != nil {
				t.Errorf(
					"Bool: input %q, unexpected error: %v",
					c.xmlValue, err,
				)
			}
			if b != c.expectedBool {
				t.Errorf(
					"Bool: input %q, expected %v, got %v",
					c.xmlValue, c.expectedBool, b,
				)
			}
		} else {
			if err == nil {
				t.Errorf(
					"Bool: input %q, expected error, got nil",
					c.xmlValue,
				)
			}
		}
	}

	// Test XML round-trip for valid values
	validCases := []string{
		"1", "0", "true", "false", " TRUE ", " False ",
	}
	for _, val := range validCases {
		aes := AutoExposureSupported(val)
		elm := aes.toXML("wscn:AutoExposureSupported")
		aes2, err := decodeAutoExposureSupported(elm)
		if err != nil {
			t.Errorf(
				"decodeAutoExposureSupported: input %q, unexpected error: %v",
				val, err,
			)
		}
		if !aes2.IsValid() {
			t.Errorf(
				"decodeAutoExposureSupported: input %q, result not valid",
				val,
			)
		}
	}

	// Test decodeAutoExposureSupported with invalid value
	elm := xmldoc.Element{
		Name: "wscn:AutoExposureSupported", Text: "maybe",
	}
	_, err := decodeAutoExposureSupported(elm)
	if err == nil {
		t.Errorf(
			"decodeAutoExposureSupported: " +
				"expected error for invalid value, got nil",
		)
	}
}
