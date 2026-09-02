// MFP - Miulti-Function Printers and scanners toolkit
// Logging facilities
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Log records

package log

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/OpenPrinting/go-mfp/util/generic"
)

// Record allows to build a multi-line log message, which
// will be written atomically, and will not be intermixed
// with logging activity from other goroutines and/or split
// between different files during log rotation.
type Record struct {
	parent *Logger    // Parent logger
	prefix string     // Log prefix
	lines  [][]byte   // Collected lines
	levels []Level    // Corresponding levels
	mutex  sync.Mutex // Access lock
}

// Commit writes Record to the parent [Logger].
func (rec *Record) Commit() *Logger {
	rec.Flush()
	return rec.parent
}

// Rollback drops the Record.
func (rec *Record) Rollback() *Logger {
	return rec.parent
}

// Flush writes Record to the parent [Logger] and resets its buffers.
func (rec *Record) Flush() *Record {
	rec.mutex.Lock()
	lines, levels := rec.lines, rec.levels
	rec.lines = rec.lines[:0]
	rec.levels = rec.levels[:0]
	rec.mutex.Unlock()

	rec.parent.send(rec.prefix, levels, lines)
	return rec
}

// Trace writes a Trace-level message to the Record.
func (rec *Record) Trace(format string, v ...any) *Record {
	return rec.format(LevelTrace, format, v...)
}

// Verbose writes a Verbose-level message to the Record.
func (rec *Record) Verbose(format string, v ...any) *Record {
	return rec.format(LevelVerbose, format, v...)
}

// Debug writes a Debug-level message to the Record.
func (rec *Record) Debug(format string, v ...any) *Record {
	return rec.format(LevelDebug, format, v...)
}

// Info writes a Info-level message to the Record.
func (rec *Record) Info(format string, v ...any) *Record {
	return rec.format(LevelInfo, format, v...)
}

// Warning writes a Warning-level message to the Record.
func (rec *Record) Warning(format string, v ...any) *Record {
	return rec.format(LevelWarning, format, v...)
}

// Error writes a Error-level message to the Record.
func (rec *Record) Error(format string, v ...any) *Record {
	return rec.format(LevelError, format, v...)
}

// Fatal writes a Fatal-level message to the Record.
//
// It calls os.Exit(1) and never returns.
func (rec *Record) Fatal(format string, v ...any) {
	rec.format(LevelFatal, format, v...)
	rec.Commit()
	os.Exit(1)
}

// Dump writes the hex dump to the Record.
func (rec *Record) Dump(level Level, data []byte) *Record {
	buf := bufAlloc()
	defer bufFree(buf)

	for len(data) > 0 {
		wid := generic.Min(len(data), 16)

		for i := 0; i < wid; i++ {
			fmt.Fprintf(buf, "%2.2x", data[i])
			if i%4 == 3 {
				buf.Write([]byte(":"))
			} else {
				buf.Write([]byte(" "))
			}

		}

		if pad := 16 - wid; pad > 0 {
			buf.Write(bytes.Repeat([]byte(" "), pad*3))
		}

		buf.Write([]byte{' ', ' '})

		for i := 0; i < wid; i++ {
			c := data[i]

			if 0x20 <= c && c < 0x80 {
				buf.WriteByte(c)
			} else {
				buf.WriteByte('.')
			}
		}

		buf.Write([]byte{'\n'})
		data = data[wid:]
	}

	return rec.text(level, 0, generic.CopySlice(buf.Bytes()))
}

// Object writes any object that implements [Marshaler]
// interface to the Record.
func (rec *Record) Object(level Level, indent int, obj Marshaler) *Record {
	text := obj.MarshalLog()
	return rec.text(level, indent, text)
}

// format writes a single formatted message to the Record
func (rec *Record) format(level Level, format string, v ...any) *Record {
	buf := bufAlloc()
	defer bufFree(buf)

	fmt.Fprintf(buf, format, v...)
	return rec.text(level, 0, generic.CopySlice(buf.Bytes()))
}

// text writes a text message to the Record
func (rec *Record) text(level Level, indent int, text []byte) *Record {
	if len(text) == 0 {
		return rec
	}

	if text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}

	lines := bytes.Split(text, []byte("\n"))
	levels := make([]Level, len(lines))

	for i := range lines {
		levels[i] = level
	}

	if indent > 0 {
		buf := bufAlloc()
		defer bufFree(buf)

		buf.Write(bytes.Repeat([]byte(" "), indent))

		for i := range lines {
			buf.Truncate(indent)
			buf.Write(lines[i])
			lines[i] = generic.CopySlice(buf.Bytes())
		}
	}

	rec.mutex.Lock()
	rec.lines = append(rec.lines, lines...)
	rec.levels = append(rec.levels, levels...)
	rec.mutex.Unlock()

	return rec
}
