// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// argv validators tests

package transport

import "testing"

// TestArgvValidatePortRange tests the ArgvValidatePortRange function
func TestArgvValidatePortRange(t *testing.T) {
	type testData struct {
		s  string
		ok bool
	}

	tests := []testData{
		{"1", true},
		{"1-2", true},
		{"2-1", false},
		{"0", false},
		{"65535", true},
		{"65536", false},
		{"0-2", false},
		{"1-65535", true},
		{"1-65536", false},
		{"-2", false},
		{"2-", false},
	}

	for _, test := range tests {
		err := ArgvValidatePortRange(test.s)
		switch {
		case test.ok && err != nil:
			t.Errorf("%s: %s", test.s, err)
		case !test.ok && err == nil:
			t.Errorf("%s: error expected, nil returned", test.s)
		}
	}
}
