// MFP - Miulti-Function Printers and scanners toolkit
// Logging facilities
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package log

import (
	"sync"
)

// Standard loggers:
var (
	// StderrLogger writes logs to the stderr.
	StderrLogger = NewLogger(LevelAll, Stderr)

	// ConsoleLogger writes logs to the console.
	// Logs are colored, if console supports that.
	ConsoleLogger = NewLogger(LevelAll, Console)

	// FatalLogger writes LevelFatal logs to the
	// stderr and discards less important messages.
	FatalLogger = NewLogger(LevelFatal, Stderr)

	// DiscardLogger discards all logs written to.
	DiscardLogger = NewLogger(LevelNone, Discard)

	// DefaultLogger is the default logging destination.
	DefaultLogger = ConsoleLogger
)

// Logger is the logging destination.
// It can be connected to console, to the disk file etc...
type Logger struct {
	out     []loggerDest // Attached destinations
	outLock sync.Mutex   // Destinations modification lock
}

// loggerDest represents logging destination
type loggerDest struct {
	level   Level
	backend Backend
}

// NewLogger returns a new logger, attached to the specified backend
func NewLogger(lvl Level, b Backend) *Logger {
	return &Logger{
		out: []loggerDest{
			{
				level:   lvl,
				backend: b,
			},
		},
	}
}

// Attach adds an additional [Backend] to send logs to.
//
// If this backend already attached to this logger, it
// only updates the log level.
func (lgr *Logger) Attach(lvl Level, b Backend) {
	// Must do under the lock
	lgr.outLock.Lock()
	defer lgr.outLock.Unlock()

	// If Backend already attached just update a Level
	for i := range lgr.out {
		if lgr.out[i].backend == b {
			lgr.out[i].level = lvl
			return
		}
	}

	// Create new attachment
	lgr.out = append(lgr.out, loggerDest{level: lvl, backend: b})
}

// Begin initiates creation of a new multi-line log [Record].
//
// Records are always written atomically. Records written from
// the concurrently running goroutines are never intermixed at
// output. During log rotation, Records are not split between
// different log files.
func (lgr *Logger) Begin(prefix string) *Record {
	return &Record{parent: lgr, prefix: prefix}
}

// Trace writes a Trace-level message to the Logger.
func (lgr *Logger) Trace(prefix, format string, v ...any) *Logger {
	return lgr.Begin(prefix).Trace(format, v...).Commit()
}

// Verbose writes a Verbose-level message to the Logger.
func (lgr *Logger) Verbose(prefix, format string, v ...any) *Logger {
	return lgr.Begin(prefix).Verbose(format, v...).Commit()
}

// Debug writes a Debug-level message to the Logger.
func (lgr *Logger) Debug(prefix, format string, v ...any) *Logger {
	return lgr.Begin(prefix).Debug(format, v...).Commit()
}

// Info writes a Info-level message to the Logger.
func (lgr *Logger) Info(prefix, format string, v ...any) *Logger {
	return lgr.Begin(prefix).Info(format, v...).Commit()
}

// Warning writes a Warning-level message to the Logger.
func (lgr *Logger) Warning(prefix, format string, v ...any) *Logger {
	return lgr.Begin(prefix).Warning(format, v...).Commit()
}

// Error writes a Error-level message to the Logger.
func (lgr *Logger) Error(prefix, format string, v ...any) *Logger {
	return lgr.Begin(prefix).Error(format, v...).Commit()
}

// Fatal writes a Fatal-level message to the Logger.
//
// It calls os.Exit(1) and never returns.
func (lgr *Logger) Fatal(prefix, format string, v ...any) {
	lgr.Begin(prefix).Fatal(format, v...)
}

// Dump writes the hex dump to the Logger.
func (lgr *Logger) Dump(prefix string, level Level, data []byte) {
	lgr.Begin(prefix).Dump(level, data).Commit()
}

// Object writes any object that implements [Marshaler]
// interface to the Logger.
func (lgr *Logger) Object(prefix string, level Level, indent int, obj Marshaler) *Logger {
	return lgr.Begin(prefix).Object(level, indent, obj).Commit()
}

// send writes some lines to the Logger.
func (lgr *Logger) send(prefix string, levels []Level, lines [][]byte) *Logger {
	// Prepend prefix
	if prefix != "" {
		prefixed := make([][]byte, len(lines))
		for i := range lines {
			prefixed[i] = []byte(prefix + ": " + string(lines[i]))
		}
		lines = prefixed
	}

	// Send message to all destinations
	lgr.outLock.Lock()
	out := lgr.out
	lgr.outLock.Unlock()

	for _, dest := range out {
		// Filter lines by level
		filteredLevels := make([]Level, 0, len(lines))
		filteredLines := make([][]byte, 0, len(lines))

		for i := range lines {
			lvl := levels[i]
			if lvl >= dest.level {
				filteredLevels = append(filteredLevels, lvl)
				filteredLines = append(filteredLines,
					trim(lines[i]))
			}
		}

		// Send to destination
		if len(filteredLines) > 0 {
			dest.backend.Send(filteredLevels, filteredLines)
		}
	}

	return lgr
}
