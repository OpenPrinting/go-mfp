// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// DateTime element

package wsscan

import (
	"fmt"
	"time"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// DateTime specifies the time at which a condition was set.
// It is represented as an XML Schema dateTime value.
// Example: 2008-10-12T14:10:00Z
type DateTime time.Time

// String returns the RFC3339 representation of the time in UTC.
func (dt DateTime) String() string {
	return time.Time(dt).Format(time.RFC3339)
}

// toXML generates XML tree for the [DateTime].
func (dt DateTime) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: dt.String(),
	}
}

// decodeDateTime decodes [DateTime] from the XML tree.
// Accepts valid RFC3339/RFC3339Nano date-time strings.
func decodeDateTime(root xmldoc.Element) (DateTime, error) {
	text := root.Text
	if text == "" {
		return DateTime{}, fmt.Errorf("invalid dateTime: %q", text)
	}

	// Try RFC3339 first, then RFC3339Nano.
	if t, err := time.Parse(time.RFC3339, text); err == nil {
		return DateTime(t), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, text); err == nil {
		return DateTime(t), nil
	}

	return DateTime{}, fmt.Errorf("invalid dateTime: %q", text)
}
