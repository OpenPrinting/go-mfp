// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// USB Endpoint emulation.

package usbip

import (
	"context"
	"io"
	"sync"
	"syscall"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// Endpoint is the virtual USB endpoint. Effectively, it implements
// the uni-direction data queue.
type Endpoint struct {
	ty       usb.EndpointType        // Endpoint type
	attrs    usb.EndpointAttributes  // Endpoint attributes
	pktsize  int                     // The packet size
	inqueue  []*protoIOSubmitRequest // Queue of input submit requests
	outqueue []*protoIOSubmitRequest // Queue of output submit requests
	inchan   chan struct{}           // Signaled when inqueue is pushed
	outchan  chan struct{}           // Signaled when outqueue is pushed
	lock     sync.Mutex              // Access lock
}

// NewEndpoint creates a new endpoint.
func NewEndpoint(desc usb.EndpointDescriptor) *Endpoint {
	ep := &Endpoint{
		ty:      desc.Type,
		attrs:   desc.BMAttributes,
		pktsize: int(desc.WMaxPacketSize),
		inchan:  make(chan struct{}, 1),
		outchan: make(chan struct{}, 1),
	}

	return ep
}

// Type returns the endpoint type.
func (ep *Endpoint) Type() usb.EndpointType {
	return ep.ty
}

// Attrs returns endpoint attributes.
func (ep *Endpoint) Attrs() usb.EndpointAttributes {
	return ep.attrs
}

// PktSize returns max packet size for the endpoint.
func (ep *Endpoint) PktSize() int {
	return ep.pktsize
}

// submit routes protoIOSubmitRequest to the endpoint.
// It returns 0 on success and some error code otherwise.
func (ep *Endpoint) submit(rq *protoIOSubmitRequest) syscall.Errno {
	ep.lock.Lock()
	if rq.Input {
		rq.Buffer = make([]byte, rq.Length)
		ep.inqueue = append(ep.inqueue, rq)
		select {
		case ep.inchan <- struct{}{}:
		default:
		}
	} else {
		ep.outqueue = append(ep.outqueue, rq)
		select {
		case ep.outchan <- struct{}{}:
		default:
		}
	}
	ep.lock.Unlock()

	return 0
}

// unlink routes protoIOUnlinkRequest to the endpoint
// It returns syscall.ECONNRESET on success and 0 otherwise.
func (ep *Endpoint) unlink(rq *protoIOUnlinkRequest) syscall.Errno {
	ep.lock.Lock()
	defer ep.lock.Unlock()

	if rq.Input {
		for i := range ep.inqueue {
			if ep.inqueue[i].Seqnum == rq.UnlinkSeqnum {
				copy(ep.inqueue[i:], ep.inqueue[i+1:])
				ep.inqueue = ep.inqueue[:len(ep.inqueue)-1]
				return syscall.ECONNRESET
			}
		}
	} else {
		for i := range ep.outqueue {
			if ep.outqueue[i].Seqnum == rq.UnlinkSeqnum {
				copy(ep.outqueue[i:], ep.outqueue[i+1:])
				ep.outqueue = ep.outqueue[:len(ep.outqueue)-1]
				return syscall.ECONNRESET
			}
		}
	}

	return 0
}

// shutdown cancels all pending protoIOSubmitRequest.
func (ep *Endpoint) shutdown() {
	ep.lock.Lock()
	ep.inqueue = ep.inqueue[:0]
	ep.outqueue = ep.outqueue[:0]
	ep.lock.Unlock()
}

// Read returns data that was sent to the [Endpoint] from the USB side.
func (ep *Endpoint) Read(buf []byte) (int, error) {
	return ep.ReadContext(context.Background(), buf)
}

// ReadContext is the [context.Context]-aware version of the [Endpoint.Read].
func (ep *Endpoint) ReadContext(ctx context.Context, buf []byte) (int, error) {
	// Check that Endpoint direction allows reading
	if ep.ty == usb.EndpointIn {
		return 0, io.ErrClosedPipe
	}

	// If context already canceled, return immediately
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	// Wait for the data to arrive
	ep.lock.Lock()
	defer ep.lock.Unlock()

	for len(ep.outqueue) == 0 {
		ep.lock.Unlock()
		select {
		case <-ep.outchan:
		case <-ctx.Done():
			return 0, ctx.Err()
		}
		ep.lock.Lock()
	}

	rq := ep.outqueue[0]
	n := copy(buf, rq.Buffer[rq.actualLength:])
	rq.actualLength += n

	if rq.actualLength == int(rq.Length) {
		l := copy(ep.outqueue, ep.outqueue[1:])
		ep.outqueue = ep.outqueue[:l]

		if rq.completion != nil {
			rq.completion(rq, 0)
		}
	}

	return n, nil
}

// Write writes data that will arrive at the [Endpoint] at the USB side.
func (ep *Endpoint) Write(buf []byte) (int, error) {
	return ep.WriteContext(context.Background(), buf)
}

// WriteContext is the [context.Context]-aware version of the [Endpoint.Write].
func (ep *Endpoint) WriteContext(ctx context.Context, buf []byte) (int, error) {
	// Check that Endpoint direction allows writing
	if ep.ty == usb.EndpointOut {
		return 0, io.ErrClosedPipe
	}

	// If context already canceled, return immediately
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	// Wait for the opportunity to write the data
	ep.lock.Lock()
	defer ep.lock.Unlock()

	for len(ep.inqueue) == 0 {
		ep.lock.Unlock()
		select {
		case <-ep.inchan:
		case <-ctx.Done():
			return 0, ctx.Err()
		}
		ep.lock.Lock()
	}

	rq := ep.inqueue[0]
	n := copy(rq.Buffer[rq.actualLength:], buf)
	rq.actualLength += n

	if true || rq.actualLength == int(rq.Length) {
		l := copy(ep.inqueue, ep.inqueue[1:])
		ep.inqueue = ep.inqueue[:l]
		if rq.completion != nil {
			rq.completion(rq, 0)
		}
	}

	return n, nil
}

// EndpointsInOut accept a group of endpoints (assuming they all
// belongs to the same interface/alt setting) and returns the first
// input (EndpointIn) and the first output (EndpointOut) endpoint
// in the group.
func EndpointsInOut(endpoints []*Endpoint) (in, out *Endpoint) {
	for _, ep := range endpoints {
		switch {
		case in == nil && ep.Type() == usb.EndpointIn:
			in = ep
		case out == nil && ep.Type() == usb.EndpointOut:
			out = ep
		}
	}

	return
}
