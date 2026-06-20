// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP request handler

package ipp

import (
	"context"
	"io"

	"github.com/OpenPrinting/goipp"
)

// Handler is the IPP request handler. It implements http.Handler interface.
type Handler struct {
	Op       goipp.Op
	callback func(context.Context, *goipp.Message, io.Reader) (
		*goipp.Message, io.ReadCloser, error)
}

// NewHandler creates a new IPP handler from the function that
// consumes [Request] and returns the [goipp.Message] response:
//
//	func DoCUPSGetDefaultRequest(rq *CUPSGetDefaultRequest) *goipp.Message {
//	. . .
//	}
//
//	handler := NewHandler(DoCUPSGetDefaultRequest)
func NewHandler[RQT any,
	RQ interface {
		*RQT
		Request
	}](f func(ctx context.Context, rq RQ) (*goipp.Message, io.ReadCloser, error)) *Handler {

	callback := func(ctx context.Context,
		rqMsg *goipp.Message, body io.Reader) (

		*goipp.Message, io.ReadCloser, error) {

		rq := RQ(new(RQT))
		rq.Header().setBody(body)

		err := rq.Decode(rqMsg, nil)
		if err != nil {
			return nil, nil, err
		}

		return f(ctx, rq)
	}

	return &Handler{
		Op:       RQ.GetOp(nil),
		callback: callback,
	}
}

// handle handles the received request.
func (h *Handler) handle(ctx context.Context, rq *goipp.Message, body io.Reader) (
	*goipp.Message, io.ReadCloser, error) {
	return h.callback(ctx, rq, body)
}
