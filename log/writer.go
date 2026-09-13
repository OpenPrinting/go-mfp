// MFP - Miulti-Function Printers and scanners toolkit
// Logging facilities
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// io.Writer that writes to Logger

package log

type writer struct {
	lgr    *Logger // Target logger
	level  Level   // Log level used to write logs
	prefix string  // Log prefix
}

// Write writes text message into the target logger.
// It implements io.Writer interface.
func (w writer) Write(data []byte) (int, error) {
	w.lgr.Begin(w.prefix).text(w.level, 0, data).Commit()
	return len(data), nil
}
