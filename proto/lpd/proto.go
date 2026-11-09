// MFP - Multi-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Protocol constants and low-level wire encoding

package lpd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Command is the LPD top-level command code (RFC 1179, section 5).
type Command byte

const (
	CmdPrintWaiting    Command = 0x01 // Print any waiting jobs
	CmdReceiveJob      Command = 0x02 // Receive a printer job
	CmdQueueStateShort Command = 0x03 // Send queue state (short)
	CmdQueueStateLong  Command = 0x04 // Send queue state (long)
	CmdRemoveJobs      Command = 0x05 // Remove jobs
)

// String returns the Command name, for logging.
func (cmd Command) String() string {
	switch cmd {
	case CmdPrintWaiting:
		return "print-waiting-jobs"
	case CmdReceiveJob:
		return "receive-job"
	case CmdQueueStateShort:
		return "queue-state-short"
	case CmdQueueStateLong:
		return "queue-state-long"
	case CmdRemoveJobs:
		return "remove-jobs"
	}
	return fmt.Sprintf("unknown(0x%2.2x)", byte(cmd))
}

// Subcommand is the code of the Receive-Job subcommand
// (RFC 1179, section 6).
type Subcommand byte

const (
	SubcmdAbortJob    Subcommand = 0x01 // Abort job
	SubcmdControlFile Subcommand = 0x02 // Receive control file
	SubcmdDataFile    Subcommand = 0x03 // Receive data file
)

// String returns the Subcommand name, for logging.
func (sub Subcommand) String() string {
	switch sub {
	case SubcmdAbortJob:
		return "abort-job"
	case SubcmdControlFile:
		return "control-file"
	case SubcmdDataFile:
		return "data-file"
	}
	return fmt.Sprintf("unknown(0x%2.2x)", byte(sub))
}

// Acknowledgement bytes. RFC 1179 defines a single zero octet as
// a positive acknowledgement; any other value is negative.
const (
	ack  byte = 0x00
	nack byte = 0x01
)

// MaxLineLen is the maximum length of the command or subcommand
// line, including the command byte but not including the
// terminating LF.
const MaxLineLen = 4096

var (
	ErrLineTooLong   = errors.New("lpd: command line too long")
	ErrEmptyLine     = errors.New("lpd: empty command line")
	ErrBadFileHeader = errors.New("lpd: malformed file header")
	ErrFileTooLarge  = errors.New("lpd: file too large")
	ErrNoTerminator  = errors.New("lpd: missed file terminating zero octet")
)

// readLine reads a single LF-terminated line from r.
// The LF is not included into the returned line.
//
// It returns io.EOF if the stream ends before any byte is read,
// io.ErrUnexpectedEOF if it ends in the middle of the line,
// and ErrLineTooLong if the line exceeds maxlen bytes.
func readLine(r *bufio.Reader, maxlen int) ([]byte, error) {
	var line []byte

	for {
		c, err := r.ReadByte()
		switch {
		case err == io.EOF && len(line) == 0:
			return nil, io.EOF
		case err == io.EOF:
			return nil, io.ErrUnexpectedEOF
		case err != nil:
			return nil, err
		}

		if c == '\n' {
			return line, nil
		}

		if len(line) >= maxlen {
			return nil, ErrLineTooLong
		}

		line = append(line, c)
	}
}

// parseLine splits the command or subcommand line into the
// leading code byte and the white space separated operands.
//
// RFC 1179 defines white space as SP or HT.
func parseLine(line []byte) (code byte, operands []string, err error) {
	if len(line) == 0 {
		return 0, nil, ErrEmptyLine
	}

	operands = strings.FieldsFunc(string(line[1:]), func(c rune) bool {
		return c == ' ' || c == '\t'
	})

	if len(operands) == 0 {
		operands = nil
	}

	return line[0], operands, nil
}

// parseFileHeader parses operands of the control file or data
// file subcommand, which are "count SP name".
//
// count must be a positive decimal number not exceeding maxsize.
// Zero count ("read until EOF") is not supported.
func parseFileHeader(operands []string, maxsize int64) (
	count int64, name string, err error) {

	if len(operands) != 2 {
		return 0, "", ErrBadFileHeader
	}

	count, err = strconv.ParseInt(operands[0], 10, 64)
	if err != nil || count <= 0 {
		return 0, "", ErrBadFileHeader
	}

	if count > maxsize {
		return 0, "", ErrFileTooLarge
	}

	return count, operands[1], nil
}

// readFileBody reads count bytes of the file contents, followed
// by the terminating zero octet.
//
// The buffer grows as data actually arrives, so a peer that
// announces a huge count but sends nothing doesn't cause
// a huge allocation.
func readFileBody(r io.Reader, count int64) ([]byte, error) {
	var buf bytes.Buffer

	n, err := io.CopyN(&buf, r, count+1)
	if err == io.EOF || (err == nil && n < count+1) {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()
	if data[count] != 0 {
		return nil, ErrNoTerminator
	}

	return data[:count], nil
}
