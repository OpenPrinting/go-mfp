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

// CapturedDoc holds a single captured print document
// with its negotiated job parameters and raw bytes.
type CapturedDoc struct {
	Params abstract.PrinterRequest
	Data   []byte
}

// DocumentCapture implements abstract.Printer and collects
// all incoming print documents for later inspection.
// It is safe for concurrent use.
type DocumentCapture struct {
	mu   sync.Mutex
	docs []CapturedDoc
	done chan struct{}
}

// NewDocumentCapture creates a new DocumentCapture.
func NewDocumentCapture() *DocumentCapture {
	return &DocumentCapture{
		done: make(chan struct{}),
	}
}

// PrintDocument implements abstract.Printer.
// It reads the full document body and stores it along with params.
func (dc *DocumentCapture) PrintDocument(
	params abstract.PrinterRequest, body io.Reader) error {

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	dc.mu.Lock()
	dc.docs = append(dc.docs, CapturedDoc{
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

// OnDocument returns a channel that is closed when the first
// document is received. Useful for waiting without polling.
func (dc *DocumentCapture) OnDocument() <-chan struct{} {
	return dc.done
}

// Wait blocks until at least one document has been captured.
func (dc *DocumentCapture) Wait() {
	<-dc.done
}

// Docs returns a snapshot of all captured documents so far.
func (dc *DocumentCapture) Docs() []CapturedDoc {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	out := make([]CapturedDoc, len(dc.docs))
	copy(out, dc.docs)
	return out
}
