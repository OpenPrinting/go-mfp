// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IEEE 1284 printer

package ieee1284

import (
	"context"
	"io"
	"sync"

	"github.com/OpenPrinting/go-mfp/log"
)

// Printer implements the IEEE-1284 printer.
type Printer struct {
	ctx     context.Context // Logging context
	state   parserState     // Current parser state
	format  DocFormat       // Detected format of current document
	buf     []byte          // Parser input buffer
	lineBuf []byte          // Partial PJL line buffer
	docBuf  []byte          // Accumulated document content
	handler DocumentHandler // Called when document is complete
	model   string          // Printer model name for PJL INFO ID

	mu      sync.Mutex // Protects respBuf and closed
	cond    *sync.Cond // Signaled when respBuf has data or closed
	respBuf []byte     // Response data waiting to be read
	closed  bool       // Set by Close()
}

// NewPrinter creates a new printer.
func NewPrinter(ctx context.Context, handler DocumentHandler) *Printer {
	p := &Printer{
		ctx:     ctx,
		handler: handler,
	}
	p.cond = sync.NewCond(&p.mu)
	return p
}

// SetModel sets the printer model name returned by PJL INFO ID.
func (p *Printer) SetModel(name string) {
	p.model = name
}

// Write consumes data sent to the Printer from the host.
// It implements the [io.Writer] interface.
func (p *Printer) Write(data []byte) (int, error) {
	log.Debug(p.ctx, "ieee1284: Write: %d bytes", len(data))
	log.Dump(p.ctx, log.LevelTrace, data)
	p.feed(data)
	return len(data), nil
}

// Read returns data sent by the Printer to the host.
// It blocks until response data is available or the printer is closed.
// It implements the [io.Reader] interface.
func (p *Printer) Read(buf []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.respBuf) == 0 && !p.closed {
		p.cond.Wait()
	}

	if len(p.respBuf) == 0 {
		return 0, io.EOF
	}

	n := copy(buf, p.respBuf)
	p.respBuf = p.respBuf[n:]
	return n, nil
}

// Close closes the printer.
// It implements the [io.Closer] interface.
func (p *Printer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true
	p.cond.Broadcast()
	return nil
}

// queueResponse appends response data to the response buffer
// and wakes any goroutine blocked in Read().
func (p *Printer) queueResponse(data []byte) {
	if len(data) == 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.respBuf = append(p.respBuf, data...)
	p.cond.Signal()
}
