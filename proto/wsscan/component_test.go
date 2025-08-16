// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for Component element

package wsscan

import "testing"

var testComponent = testEnum[Component]{
	decodeStr: DecodeComponent,
	decodeXML: decodeComponent,
	dataset: []testEnumData[Component]{
		{ADFComponent, "ADF"},
		{FilmComponent, "Film"},
		{MediaPathComponent, "MediaPath"},
		{PlatenComponent, "Platen"},
	},
}

// TestComponent tests [Component] common methods and functions.
func TestComponent(t *testing.T) {
	testComponent.run(t)
}
