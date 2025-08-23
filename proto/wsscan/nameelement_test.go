// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for NameElement

package wsscan

import "testing"

var testNameElement = testEnum[NameElement]{
	decodeStr: DecodeNameElement,
	decodeXML: decodeNameElement,
	dataset: []testEnumData[NameElement]{
		{Calibrating, "Calibrating"},
		{CoverOpen, "CoverOpen"},
		{InputTrayEmpty, "InputTrayEmpty"},
		{InterlockOpen, "InterlockOpen"},
		{InternalStorageFull, "InternalStorageFull"},
		{LampError, "LampError"},
		{LampWarming, "LampWarming"},
		{MediaJam, "MediaJam"},
		{MultipleFeedError, "MultipleFeedError"},
	},
}

// TestNameElement tests [NameElement] common methods and functions.
func TestNameElement(t *testing.T) {
	testNameElement.run(t)
}
