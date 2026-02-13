// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Print job stream parser

package ieee1284

import (
	"bytes"
	"strings"
)

// UEL is the Universal Exit Language escape sequence.
const uel = "\x1b%-12345X"

// DocumentHandler is called when a complete document is extracted
// from the print stream.
type DocumentHandler func(format DocFormat, data []byte)

// parserState represents the current state of the stream parser.
type parserState int

const (
	stateIdle     parserState = iota // Between jobs or at start
	statePJL                         // Inside PJL command block
	stateDocument                    // Inside document content
)

// feed processes incoming data through the parser state machine.
// It is called from Printer.Write() with each chunk of data.
func (p *Printer) feed(data []byte) {
	p.buf = append(p.buf, data...)

	for len(p.buf) > 0 {
		switch p.state {
		case stateIdle:
			if !p.feedIdle() {
				return
			}
		case statePJL:
			if !p.feedPJL() {
				return
			}
		case stateDocument:
			if !p.feedDocument() {
				return
			}
		}
	}
}

// feedIdle handles the idle state.
// Returns true if a transition occurred and parsing should continue.
func (p *Printer) feedIdle() bool {
	// Check for UEL prefix
	if hasOrWillHavePrefix(p.buf, []byte(uel)) {
		if len(p.buf) < len(uel) {
			// Partial UEL — wait for more data
			return false
		}

		// Full UEL found. Check what follows.
		after := p.buf[len(uel):]
		if hasOrWillHavePrefix(after, []byte("@PJL")) {
			if len(after) < 4 {
				// Partial "@PJL" — wait for more data
				return false
			}
			// UEL + @PJL → transition to PJL state
			p.buf = after
			p.state = statePJL
			p.lineBuf = nil
			return true
		}

		// UEL without @PJL — discard the UEL and stay idle
		p.buf = after
		return len(after) > 0
	}

	// No UEL — try magic byte detection for raw documents
	format := detectFormatByMagic(p.buf)
	if format != DocFormatUnknown {
		p.state = stateDocument
		p.format = format
		p.docBuf = nil
		return true
	}

	// Unrecognized data — discard it
	p.buf = nil
	return false
}

// feedPJL handles the PJL state, parsing PJL commands line by line.
// Returns true if a transition occurred and parsing should continue.
func (p *Printer) feedPJL() bool {
	for {
		// Find the next line ending (\r\n or \n)
		idx := bytes.IndexByte(p.buf, '\n')
		if idx < 0 {
			// No complete line yet — wait for more data
			return false
		}

		// Extract line (strip trailing \r\n)
		line := string(p.buf[:idx])
		line = strings.TrimRight(line, "\r")
		p.buf = p.buf[idx+1:]

		// Parse PJL command
		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)

		// Strip "@PJL" prefix
		const pjlPrefix = "@PJL"
		if !strings.HasPrefix(upper, pjlPrefix) {
			continue
		}
		after := trimmed[len(pjlPrefix):]

		// Check for @PJL ENTER LANGUAGE=X
		const enterPrefix = "@PJL ENTER LANGUAGE="
		if strings.HasPrefix(upper, enterPrefix) {
			lang := strings.TrimSpace(line[len(enterPrefix):])
			lang = strings.TrimSpace(lang)
			p.format = detectFormatByLanguage(lang)
			p.state = stateDocument
			p.docBuf = nil
			return true
		}

		// Dispatch other PJL commands
		cmd := parsePJLCommand(after)
		p.handlePJL(cmd)
	}
}

// feedDocument handles the document state, scanning for end markers.
// Returns true if a transition occurred and parsing should continue.
func (p *Printer) feedDocument() bool {
	for i := 0; i < len(p.buf); i++ {
		b := p.buf[i]

		// Check for UEL start (ESC)
		if b == 0x1b {
			remaining := p.buf[i:]
			if hasOrWillHavePrefix(remaining, []byte(uel)) {
				if len(remaining) < len(uel) {
					// Partial UEL at end of buffer.
					// Keep doc data before it, hold
					// the partial UEL.
					p.docBuf = append(p.docBuf,
						p.buf[:i]...)
					p.buf = remaining
					return false
				}

				// Full UEL found — document ends here
				p.docBuf = append(p.docBuf, p.buf[:i]...)
				p.emitDocument()
				p.buf = p.buf[i:]
				p.state = stateIdle
				return true
			}
		}

		// Check for Ctrl-D (PostScript/PDF end marker)
		if b == 0x04 {
			switch p.format {
			case DocFormatPostScript, DocFormatPDF,
				DocFormatUnknown:
				p.docBuf = append(p.docBuf, p.buf[:i]...)
				p.emitDocument()
				p.buf = p.buf[i+1:] // skip Ctrl-D
				p.state = stateIdle
				return true
			}
		}

		// Check for PCL reset (ESC E) as end marker for PCL 5
		if b == 0x1b && p.format == DocFormatPCL5 {
			if i+1 < len(p.buf) {
				if p.buf[i+1] == 'E' {
					p.docBuf = append(p.docBuf,
						p.buf[:i]...)
					p.emitDocument()
					p.buf = p.buf[i+2:] // skip ESC E
					p.state = stateIdle
					return true
				}
			} else {
				// ESC at end of buffer — could be
				// ESC E or start of UEL. Wait for
				// more data.
				p.docBuf = append(p.docBuf, p.buf[:i]...)
				p.buf = p.buf[i:]
				return false
			}
		}
	}

	// No end marker found — accumulate everything and wait
	p.docBuf = append(p.docBuf, p.buf...)
	p.buf = nil
	return false
}

// emitDocument calls the handler with the completed document.
func (p *Printer) emitDocument() {
	if p.handler != nil && len(p.docBuf) > 0 {
		p.handler(p.format, p.docBuf)
	}
	p.docBuf = nil
	p.format = DocFormatUnknown
}

// Flush should be called when the stream ends to emit any
// remaining document data that wasn't terminated by an explicit
// end marker.
func (p *Printer) Flush() {
	if p.state == stateDocument && len(p.docBuf)+len(p.buf) > 0 {
		p.docBuf = append(p.docBuf, p.buf...)
		p.buf = nil
		p.emitDocument()
	}
	p.state = stateIdle
}

// hasOrWillHavePrefix checks whether data starts with a prefix,
// or could start with it once more data arrives. This handles
// the case where a multi-byte sequence (like UEL) is split across
// Write() calls.
func hasOrWillHavePrefix(data, prefix []byte) bool {
	if len(data) >= len(prefix) {
		return bytes.HasPrefix(data, prefix)
	}
	return bytes.HasPrefix(prefix, data)
}
