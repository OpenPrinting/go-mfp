// MFP - Multi-Function Printers and scanners toolkit
// The "mfp-test" command
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Document capture

package test

import (
	"io"
	"sync"

	"github.com/OpenPrinting/go-mfp/abstract"
)

// capturedDoc holds a single captured print document
// with its negotiated job parameters and raw bytes.
type capturedDoc struct {
	Params abstract.PrinterRequest
	Data   []byte
}

// documentCapture implements abstract.Printer and collects
// all incoming print documents for later inspection.
// It is safe for concurrent use.
type documentCapture struct {
	mu       sync.Mutex
	captured []capturedDoc
	done     chan struct{}
}

// newDocumentCapture creates a new documentCapture.
func newDocumentCapture() *documentCapture {
	return &documentCapture{
		done: make(chan struct{}),
	}
}

// PrintDocument implements abstract.Printer.
// It reads the full document body and stores it along with params.
func (dc *documentCapture) PrintDocument(
	params abstract.PrinterRequest, body io.Reader) error {

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	dc.mu.Lock()
	dc.captured = append(dc.captured, capturedDoc{
		Params: params,
		Data:   data,
	})
	dc.mu.Unlock()

	// Signal that at least one document has arrived.
	select {
	case <-dc.done:
	default:
		close(dc.done)
	}

	return nil
}

// onDocument returns a channel that is closed when the first
// document is received. Useful for waiting without polling.
func (dc *documentCapture) onDocument() <-chan struct{} {
	return dc.done
}

// docs returns a snapshot of all captured documents so far.
func (dc *documentCapture) docs() []capturedDoc {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	out := make([]capturedDoc, len(dc.captured))
	copy(out, dc.captured)
	return out
}

// reset clears all captured documents and resets the onDocument
// signal so the capture can be reused for the next test run.
func (dc *documentCapture) reset() {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dc.captured = nil
	dc.done = make(chan struct{})
}
