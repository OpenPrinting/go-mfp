// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman(officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Shared test helpers

package ieee1284

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/log"
)

// newTestContext creates a context with logging disabled.
// Use log.LevelAll instead of log.LevelNone to enable verbose
// output when debugging tests.
func newTestContext() context.Context {
	logger := log.NewLogger(log.LevelNone, log.Console)
	return log.NewContext(context.Background(), logger)
}

// docResult captures a single backend invocation.
type docResult struct {
	params abstract.PrinterRequest
	data   []byte
}

// testBackend is a test implementation of abstract.Printer
// that captures all PrintDocument calls.
type testBackend struct {
	results *[]docResult
}

func (b *testBackend) PrintDocument(
	params abstract.PrinterRequest, body io.Reader) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	*b.results = append(*b.results, docResult{params, cp})
	return nil
}

// testHandler returns a testBackend that appends results
// to the provided slice.
func testHandler(results *[]docResult) *testBackend {
	return &testBackend{results: results}
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

// docBytes extracts the raw bytes from a docResult.
func docBytes(r docResult) []byte {
	return r.data
}

// docFormat extracts the DocFormat from a docResult by reverse-mapping
// the MIME type.
func docFormat(r docResult) DocFormat {
	switch r.params.Format {
	case "application/postscript":
		return DocFormatPostScript
	case "application/pdf":
		return DocFormatPDF
	case "application/vnd.hp-pcl":
		return DocFormatPCL
	case "application/vnd.hp-pclxl":
		return DocFormatPCLXL
	case "text/plain":
		return DocFormatPlainText
	default:
		return DocFormatUnknown
	}
}

// mustContain checks that got contains want as a substring.
func mustContain(t *testing.T, got, want []byte) {
	t.Helper()
	if !bytes.Contains(got, want) {
		t.Errorf("expected %q to contain %q", got, want)
	}
}
