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
	"sync/atomic"

	"golang.org/x/term"
)

// backendConsole is the Backend that writes logs to console
type backendConsole struct {
	color int32      // No: -1, Yes: +1, Unknown: 0
	mutex sync.Mutex // Send lock
}

// Line implements [Backend.Send] method
func (bk *backendConsole) Send(levels []Level, lines [][]byte) {
	// Color auto-detection
	if atomic.LoadInt32(&bk.color) == 0 {
		isatty := term.IsTerminal(int(os.Stdout.Fd()))
		if isatty {
			atomic.CompareAndSwapInt32(&bk.color, 0, +1)
		}
	}

	// Build the entire message in the buffer
	buf := bufAlloc()
	defer bufFree(buf)

	for i := range lines {
		level := levels[i]
		line := lines[i]

		var color string
		if atomic.LoadInt32(&bk.color) > 0 {
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
		}

		buf.Write([]byte(color))
		buf.Write(line)
		buf.Write([]byte("\033[0m" + "\n"))
	}

	// Now send buffer to the os.Stdout
	bk.mutex.Lock()
	buf.WriteTo(os.Stdout)
	bk.mutex.Unlock()
}
