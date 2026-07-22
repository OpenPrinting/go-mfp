// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Version test

package usb

import (
	"fmt"
	"testing"
)

// TestVersion tests Version functions and methods
func TestVersion(t *testing.T) {
	type testData struct {
		v Version // The version
		s string  // Its string representation
	}

	tests := []testData{
		{MakeVersion(1, 0), "1.0"},
		{MakeVersion(2, 1), "2.1"},
		{MakeVersion(3, 12), "3.12"},
	}

	for _, test := range tests {
		s := test.v.String()
		if s != test.s {
			t.Errorf("Version(0x%4.4x).String():\n"+
				"expected: %s\n"+
				"present:  %s\n",
				int(test.v), test.s, s)
		}

		v, err := ParseVersion(test.s)
		if err != nil || v != test.v {
			present := fmt.Sprintf("0x%4.4x", int(v))
			if err != nil {
				present = err.Error()
			}

			t.Errorf("ParseVersion(%q):\n"+
				"expected: 0x%4.4x\n"+
				"present:  %s\n",
				test.s, int(test.v), present)
		}
	}
}
