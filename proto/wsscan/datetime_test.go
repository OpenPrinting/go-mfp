// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Tests for DateTime element

package wsscan

import (
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

func TestDateTime_StringAndXML(t *testing.T) {
	// 2008-10-12T14:10:00Z
	ref, _ := time.Parse(time.RFC3339, "2008-10-12T14:10:00Z")
	dt := DateTime(ref)

	if got, want := dt.String(), "2008-10-12T14:10:00Z"; got != want {
		// ensure UTC RFC3339
		t.Fatalf("String() = %q; want %q", got, want)
	}

	xml := dt.toXML(NsWSCN + ":DateTime")
	if xml.Name != NsWSCN+":DateTime" || xml.Text != "2008-10-12T14:10:00Z" {
		t.Fatalf("toXML() unexpected: %s", xml.EncodeString(nil))
	}
}

func TestDecodeDateTime_OK(t *testing.T) {
	for _, s := range []string{
		"2008-10-12T14:10:00Z",
		"2008-10-12T14:10:00.123Z",
		"2008-10-12T14:10:00.123456Z",
		"2008-10-12T14:10:00.123456789Z",
	} {
		xml := xmldoc.Element{Name: NsWSCN + ":DateTime", Text: s}
		if _, err := decodeDateTime(xml); err != nil {
			t.Fatalf("decodeDateTime(%q) unexpected error: %v", s, err)
		}
	}
}

func TestDecodeDateTime_Errors(t *testing.T) {
	for _, s := range []string{
		"",
		"2008-10-12 14:10:00",
		"not-a-date",
		"2020-13-01T00:00:00Z",
		"2020-02-30T00:00:00Z",
	} {
		xml := xmldoc.Element{Name: NsWSCN + ":DateTime", Text: s}
		if _, err := decodeDateTime(xml); err == nil {
			t.Fatalf("decodeDateTime(%q) expected error, got nil", s)
		}
	}
}
