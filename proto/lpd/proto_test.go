// MFP - Miulti-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Tests for protocol constants and low-level wire encoding

package lpd

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

// TestCommandString tests Command.String and Subcommand.String
func TestCommandString(t *testing.T) {
	type testData struct {
		s, expected string
	}

	tests := []testData{
		{CmdPrintWaiting.String(), "print-waiting-jobs"},
		{CmdReceiveJob.String(), "receive-job"},
		{CmdQueueStateShort.String(), "queue-state-short"},
		{CmdQueueStateLong.String(), "queue-state-long"},
		{CmdRemoveJobs.String(), "remove-jobs"},
		{Command(0x42).String(), "unknown(0x42)"},
		{SubcmdAbortJob.String(), "abort-job"},
		{SubcmdControlFile.String(), "control-file"},
		{SubcmdDataFile.String(), "data-file"},
		{Subcommand(0x00).String(), "unknown(0x00)"},
	}

	for _, test := range tests {
		if test.s != test.expected {
			t.Errorf("expected %q, present %q", test.expected, test.s)
		}
	}
}

// TestReadLine tests readLine
func TestReadLine(t *testing.T) {
	type testData struct {
		in     string
		maxlen int
		lines  []string // Expected lines
		err    error    // Expected final error
	}

	tests := []testData{
		{
			in:     "\x02lp\n",
			maxlen: MaxLineLen,
			lines:  []string{"\x02lp"},
			err:    io.EOF,
		},
		{
			in:     "\x02lp\n\x0212 cfA001host\n",
			maxlen: MaxLineLen,
			lines:  []string{"\x02lp", "\x0212 cfA001host"},
			err:    io.EOF,
		},
		{
			in:     "\n",
			maxlen: MaxLineLen,
			lines:  []string{""},
			err:    io.EOF,
		},
		{
			in:     "",
			maxlen: MaxLineLen,
			err:    io.EOF,
		},
		{
			in:     "\x02lp",
			maxlen: MaxLineLen,
			err:    io.ErrUnexpectedEOF,
		},
		{
			in:     "12345\n",
			maxlen: 5,
			lines:  []string{"12345"},
			err:    io.EOF,
		},
		{
			in:     "123456\n",
			maxlen: 5,
			err:    ErrLineTooLong,
		},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.in))
		var lines []string
		var err error

		for {
			var line []byte
			line, err = readLine(r, test.maxlen)
			if err != nil {
				break
			}
			lines = append(lines, string(line))
		}

		if !reflect.DeepEqual(lines, test.lines) {
			t.Errorf("%q: lines mismatch:\n"+
				"expected: %q\n"+
				"present:  %q", test.in, test.lines, lines)
		}

		if err != test.err {
			t.Errorf("%q: error mismatch:\n"+
				"expected: %v\n"+
				"present:  %v", test.in, test.err, err)
		}
	}
}

// TestParseLine tests parseLine
func TestParseLine(t *testing.T) {
	type testData struct {
		in       string
		code     byte
		operands []string
		err      error
	}

	tests := []testData{
		{
			in:   "\x02lp",
			code: 0x02, operands: []string{"lp"},
		},
		{
			in:   "\x05lp root 12 13",
			code: 0x05, operands: []string{"lp", "root", "12", "13"},
		},
		{
			in:   "\x03lp \t 12\t\t13 ",
			code: 0x03, operands: []string{"lp", "12", "13"},
		},
		{
			in:   "\x01",
			code: 0x01,
		},
		{
			in:  "",
			err: ErrEmptyLine,
		},
	}

	for _, test := range tests {
		code, operands, err := parseLine([]byte(test.in))

		if err != test.err {
			t.Errorf("%q: error mismatch:\n"+
				"expected: %v\n"+
				"present:  %v", test.in, test.err, err)
			continue
		}

		if code != test.code {
			t.Errorf("%q: code mismatch: expected 0x%2.2x, present 0x%2.2x",
				test.in, test.code, code)
		}

		if !reflect.DeepEqual(operands, test.operands) {
			t.Errorf("%q: operands mismatch:\n"+
				"expected: %q\n"+
				"present:  %q", test.in, test.operands, operands)
		}
	}
}

// TestParseFileHeader tests parseFileHeader
func TestParseFileHeader(t *testing.T) {
	type testData struct {
		operands []string
		maxsize  int64
		count    int64
		name     string
		err      error
	}

	tests := []testData{
		{
			operands: []string{"123", "dfA001host"},
			maxsize:  1000,
			count:    123,
			name:     "dfA001host",
		},
		{
			operands: []string{"1000", "dfA001host"},
			maxsize:  1000,
			count:    1000,
			name:     "dfA001host",
		},
		{
			operands: []string{"1001", "dfA001host"},
			maxsize:  1000,
			err:      ErrFileTooLarge,
		},
		{
			operands: []string{"0", "dfA001host"},
			maxsize:  1000,
			err:      ErrBadFileHeader,
		},
		{
			operands: []string{"-5", "dfA001host"},
			maxsize:  1000,
			err:      ErrBadFileHeader,
		},
		{
			operands: []string{"abc", "dfA001host"},
			maxsize:  1000,
			err:      ErrBadFileHeader,
		},
		{
			operands: []string{"123"},
			maxsize:  1000,
			err:      ErrBadFileHeader,
		},
		{
			operands: []string{"123", "dfA001host", "extra"},
			maxsize:  1000,
			err:      ErrBadFileHeader,
		},
	}

	for _, test := range tests {
		count, name, err := parseFileHeader(test.operands, test.maxsize)

		if err != test.err {
			t.Errorf("%q: error mismatch:\n"+
				"expected: %v\n"+
				"present:  %v", test.operands, test.err, err)
			continue
		}

		if count != test.count || name != test.name {
			t.Errorf("%q: expected (%d, %q), present (%d, %q)",
				test.operands, test.count, test.name, count, name)
		}
	}
}

// TestReadFileBody tests readFileBody
func TestReadFileBody(t *testing.T) {
	type testData struct {
		in    string
		count int64
		data  string
		rest  string // Remaining unread input
		err   error
	}

	tests := []testData{
		{
			in:    "hello\x00",
			count: 5,
			data:  "hello",
		},
		{
			in:    "hello\x00\x03next",
			count: 5,
			data:  "hello",
			rest:  "\x03next",
		},
		{
			// Embedded zeroes are part of the data
			in:    "a\x00b\x00",
			count: 3,
			data:  "a\x00b",
		},
		{
			in:    "hello!",
			count: 5,
			err:   ErrNoTerminator,
		},
		{
			in:    "hello",
			count: 5,
			err:   io.ErrUnexpectedEOF,
		},
		{
			in:    "hel",
			count: 5,
			err:   io.ErrUnexpectedEOF,
		},
		{
			in:    "",
			count: 5,
			err:   io.ErrUnexpectedEOF,
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.in)
		data, err := readFileBody(r, test.count)

		if err != test.err {
			t.Errorf("%q: error mismatch:\n"+
				"expected: %v\n"+
				"present:  %v", test.in, test.err, err)
			continue
		}

		if string(data) != test.data {
			t.Errorf("%q: data mismatch:\n"+
				"expected: %q\n"+
				"present:  %q", test.in, test.data, data)
		}

		rest, _ := io.ReadAll(r)
		if string(rest) != test.rest {
			t.Errorf("%q: rest mismatch:\n"+
				"expected: %q\n"+
				"present:  %q", test.in, test.rest, rest)
		}
	}
}
