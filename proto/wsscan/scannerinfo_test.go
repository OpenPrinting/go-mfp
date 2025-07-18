// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for scanner info

package wsscan

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

func TestScannerInfo(t *testing.T) {
	dataset := []ScannerInfo{
		{Text: "some info", Lang: optional.New("en-AU, en-GB")},
		{Text: "some other info", Lang: optional.New("en-AU")},
		{Text: "some more info"},
	}

	for _, si := range dataset {
		elm := si.toXML("wscn:ScannerInfo")
		if elm.Name != "wscn:ScannerInfo" {
			t.Errorf(
				"expected element name 'wscn:ScannerInfo', got '%s'",
				elm.Name,
			)
		}
		if elm.Text != si.Text {
			t.Errorf(
				"expected element text '%s', got '%s'",
				si.Text, elm.Text,
			)
		}

		si2, err := decodeScannerInfo(elm)
		assert.NoError(err)
		if !reflect.DeepEqual(si, si2) {
			t.Errorf("expected %v, got %v", si, si2)
		}
	}
}
