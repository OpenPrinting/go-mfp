// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Shared test helpers

package ieee1284

import (
	"context"
	"testing"

	"github.com/OpenPrinting/go-mfp/log"
)

// newTestContext creates a context with logging disabled.
// Use log.LevelAll instead of log.LevelNone to enable verbose
// output when debugging tests.
func newTestContext() context.Context {
	logger := log.NewLogger(log.LevelNone, log.Console)
	return log.NewContext(context.Background(), logger)
}

// docResult captures a single handler invocation.
type docResult struct {
	format DocFormat
	data   []byte
}

// testHandler returns a DocumentHandler that appends results
// to the provided slice.
func testHandler(results *[]docResult) DocumentHandler {
	return func(format DocFormat, data []byte) {
		cp := make([]byte, len(data))
		copy(cp, data)
		*results = append(*results, docResult{format, cp})
	}
}

// writeInChunks writes data to the printer in fixed-size chunks,
// simulating USB transfer behavior (USB bulk transfers are
// typically 512 bytes).
func writeInChunks(t *testing.T, p *Printer, data []byte, chunkSize int) {
	t.Helper()
	for len(data) > 0 {
		chunk := chunkSize
		if chunk > len(data) {
			chunk = len(data)
		}
		n, err := p.Write(data[:chunk])
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		if n != chunk {
			t.Fatalf("Write: got %d, want %d", n, chunk)
		}
		data = data[n:]
	}
}
