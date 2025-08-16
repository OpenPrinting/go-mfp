// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Severity element

package wsscan

import "testing"

var testSeverity = testEnum[Severity]{
	decodeStr: DecodeSeverity,
	decodeXML: decodeSeverity,
	dataset: []testEnumData[Severity]{
		{Informational, "Informational"},
		{Warning, "Warning"},
		{Critical, "Critical"},
	},
}

// TestSeverity tests [Severity] common methods and functions.
func TestSeverity(t *testing.T) {
	testSeverity.run(t)
}

