// MFP - Miulti-Function Printers and scanners toolkit
// Logging facilities
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Console backend

package log

import (
	"os"
	"sync"
	"syscall"

	"golang.org/x/term"
)

// backendConsole is the Backend that writes logs to console
type backendConsole struct {
	mutex sync.Mutex // Send lock
}

// isTerminal keeps information if program's stdout was connected
// to the terminal.
//
// It is initialized early, to allow further os.Stdout redirection
// to the pager program.
var isTerminal = term.IsTerminal(int(syscall.Stdout))

// Line implements [Backend.Send] method
func (bk *backendConsole) Send(levels []Level, lines [][]byte) {
	// Build the entire message in the buffer
	buf := bufAlloc()
	defer bufFree(buf)

	for i := range lines {
		level := levels[i]
		line := lines[i]

		var color string
		var restore string
		if isTerminal {
			switch level {
			case LevelTrace:
				color = "\033[2m" // Dark Gray
			case LevelVerbose:
				color = "\033[36m" // Cyan
			case LevelDebug:
				color = "\033[37m" // Gray
			case LevelInfo:
				color = "\033[32m" // Green
			case LevelWarning:
				color = "\033[33m" // Yellow
			case LevelError, LevelFatal:
				color = "\033[1;37;41m" // White on Red
			}
			restore = "\033[0m"
		}

		buf.Write([]byte(color))
		buf.Write(line)
		buf.Write([]byte(restore + "\n"))
	}

	// Now send buffer to the os.Stdout
	bk.mutex.Lock()
	buf.WriteTo(os.Stdout)
	bk.mutex.Unlock()
}
