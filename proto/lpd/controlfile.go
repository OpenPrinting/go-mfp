// MFP - Multi-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Control file parser

package lpd

import (
	"bytes"
	"errors"
	"fmt"
)

// ControlFile represents the parsed LPD control file
// (RFC 1179, section 7).
//
// Only the lines that carry useful information are retained.
// Others (mail, unlink, indent, width, class, title for pr,
// symlink data and so on) are silently ignored, the same way
// as cups-lpd does.
type ControlFile struct {
	Host      string     // 'H' line: host name
	User      string     // 'P' line: user identification
	JobName   string     // 'J' line: job name for banner page
	Banner    bool       // 'L' line: print banner page
	Documents []Document // Print lines, in order of appearance
}

// Document represents a single print line of the control file.
//
// The same data file may be referenced multiple times. Each
// reference is a separate Document; this is how BSD lpr
// requests multiple copies.
type Document struct {
	Format   byte   // Format letter: 'l', 'f', 'o', ...
	DataFile string // Data file name
	Name     string // 'N' line: source file name, if any
}

// IsRaw reports if the document should be printed as is
// ('l' format: "print file leaving control characters").
func (doc Document) IsRaw() bool {
	return doc.Format == 'l'
}

// Errors returned by ParseControlFile:
var (
	ErrNoDataFileName = errors.New("lpd: print line without data file name")
)

// isPrintFormat reports if c is the format letter of the
// print line.
func isPrintFormat(c byte) bool {
	switch c {
	case 'c', // Plot CIF file
		'd', // Print DVI file
		'f', // Print formatted file
		'g', // Plot file
		'l', // Print file leaving control characters (raw)
		'n', // Print ditroff output file
		'o', // Print PostScript output file
		'p', // Print file with 'pr' format
		'r', // File to print with FORTRAN carriage control
		't', // Print troff output file
		'v': // Print raster file
		return true
	}
	return false
}

// ParseControlFile parses the LPD control file.
//
// Lines are LF-terminated; trailing CR, if any, is ignored.
// Empty lines and unknown commands are skipped.
//
// The 'N' (source file name) line is placed differently by
// different clients: BSD lpr and CUPS write it after the print
// line(s) it describes, while some others write it before. Both
// orders are handled: 'N' names the immediately preceding print
// lines that have no name yet, and if there are no such lines,
// it names all subsequent print lines until the next 'N'.
func ParseControlFile(data []byte) (*ControlFile, error) {
	cf := &ControlFile{}
	pendingName := ""

	for lineno, line := range bytes.Split(data, []byte("\n")) {
		line = bytes.TrimSuffix(line, []byte("\r"))
		if len(line) == 0 {
			continue
		}

		cmd, arg := line[0], string(line[1:])

		switch {
		case cmd == 'H':
			cf.Host = arg

		case cmd == 'P':
			cf.User = arg

		case cmd == 'J':
			cf.JobName = arg

		case cmd == 'L':
			cf.Banner = true

		case cmd == 'N':
			if cf.nameTrailing(arg) {
				pendingName = ""
			} else {
				pendingName = arg
			}

		case isPrintFormat(cmd):
			if arg == "" {
				return nil, fmt.Errorf("line %d: %w",
					lineno+1, ErrNoDataFileName)
			}

			cf.Documents = append(cf.Documents, Document{
				Format:   cmd,
				DataFile: arg,
				Name:     pendingName,
			})
		}
	}

	return cf, nil
}

// nameTrailing assigns name to the trailing run of unnamed
// documents. It returns false if there are no such documents.
func (cf *ControlFile) nameTrailing(name string) bool {
	i := len(cf.Documents)
	for i > 0 && cf.Documents[i-1].Name == "" {
		i--
	}

	if i == len(cf.Documents) {
		return false
	}

	for ; i < len(cf.Documents); i++ {
		cf.Documents[i].Name = name
	}

	return true
}

// Title returns the job title: the 'J' line, if present,
// otherwise the name of the first named document, otherwise
// an empty string.
func (cf *ControlFile) Title() string {
	if cf.JobName != "" {
		return cf.JobName
	}

	for _, doc := range cf.Documents {
		if doc.Name != "" {
			return doc.Name
		}
	}

	return ""
}

// DataFiles returns the names of data files referenced by
// the control file, without duplicates, in order of the first
// reference.
func (cf *ControlFile) DataFiles() []string {
	var names []string
	seen := make(map[string]struct{})

	for _, doc := range cf.Documents {
		if _, found := seen[doc.DataFile]; !found {
			seen[doc.DataFile] = struct{}{}
			names = append(names, doc.DataFile)
		}
	}

	return names
}
