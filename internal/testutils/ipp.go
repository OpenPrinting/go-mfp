// MFP - Miulti-Function Printers and scanners toolkit
// Utility functions and data BLOBs for testing
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP BLOB parser

package testutils

import (
	"github.com/OpenPrinting/goipp"
)

// IPPParse parses IPP binary data into the [goipp.Message].
func IPPParse(blob []byte) (*goipp.Message, error) {
	msg := &goipp.Message{}
	err := msg.DecodeBytesEx(blob,
		goipp.DecoderOptions{EnableWorkarounds: true})

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// IPPMustParse parses IPP binary data into the [goipp.Message].
// It panics if data cannot be parsed.
func IPPMustParse(blob []byte) *goipp.Message {
	msg, err := IPPParse(blob)

	if err != nil {
		panic(err)
	}

	return msg
}

// IPPDiffAttributes returns line diff between two sequences of IPP attributes
func IPPDiffAttributes(oldName string, old goipp.Attributes,
	newName string, new goipp.Attributes) string {
	f := goipp.NewFormatter()
	f.SetIndent(4)
	f.FmtAttributes(old)
	s1 := f.String()
	f.Reset()

	f.SetIndent(4)
	f.FmtAttributes(new)
	s2 := f.String()

	return DiffLines(oldName, s1, newName, s2)
}
