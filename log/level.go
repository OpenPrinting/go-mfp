// MFP - Miulti-Function Printers and scanners toolkit
// Logging facilities
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Log levels

package log

// Level enumerates possible log levels
type Level int

// Log levels:
const (
	LevelTrace   Level = iota // Protocol trace (very verbose)
	LevelVerbose              // Verbose messages
	LevelDebug                // Debug messages
	LevelInfo                 // Informational messages
	LevelWarning              // Warning messages
	LevelError                // Error messages
	LevelFatal                // Fatal errors

	// These constants are useful for filtering
	LevelAll  = LevelTrace     // Allow all levels
	LevelNone = LevelFatal + 1 // Allow no levels
)
