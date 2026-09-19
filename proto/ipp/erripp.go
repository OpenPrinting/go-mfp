// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP errors

package ipp

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/goipp"
)

// Common errors, reported as ErrIPP:
var (
// HTTPErrorMethodNotAllowed = NewHTTPError(http.StatusMethodNotAllowed, "")
)

// ErrIPP represents IPP error that can be returned to the IPP client
// as the IPP error response.
//
// It consist of the IPP status and optional message text.
// Implements [error] interface.
type ErrIPP struct {
	Version       goipp.Version // IPP version
	RequestID     uint32        // IPP Request ID
	Op            goipp.Op      // IPP operation
	Status        goipp.Status  // IPP status
	StatusMessage string        // Optional error message
}

// NewErrIPPFromMessage creates a new IPP error that can be sent as
// response to the [goipp.Message].
func NewErrIPPFromMessage(rq *goipp.Message, code goipp.Status,
	format string, args ...any) *ErrIPP {

	ver := generic.Min(rq.Version, MaxVersion)
	return &ErrIPP{
		Version:       ver,
		RequestID:     rq.RequestID,
		Op:            goipp.Op(rq.Code),
		Status:        code,
		StatusMessage: fmt.Sprintf(format, args...),
	}
}

// NewErrIPPFromRequest creates a new IPP error that can be sent as
// response to the decoded [Request].
func NewErrIPPFromRequest(rq Request, code goipp.Status,
	format string, args ...any) *ErrIPP {

	hdr := rq.Header()

	ver := generic.Min(hdr.Version, MaxVersion)
	return &ErrIPP{
		Version:       ver,
		RequestID:     hdr.RequestID,
		Status:        code,
		StatusMessage: fmt.Sprintf(format, args...),
	}
}

// Error returns an error string. It implements [error] interface.
func (e *ErrIPP) Error() string {
	msg := e.StatusMessage
	if msg == "" {
		msg = e.Status.String()
	}
	return fmt.Sprintf("IPP %s", msg)
}

// Encode encodes ErrIPP into the goipp.Message.
func (e *ErrIPP) Encode() *goipp.Message {
	msg := &goipp.Message{
		Version:   e.Version,
		RequestID: e.RequestID,
		Code:      goipp.Code(e.Status),
	}

	msg.Operation.Add(goipp.MakeAttribute("attributes-charset",
		goipp.TagCharset, goipp.String("utf-8")))
	msg.Operation.Add(goipp.MakeAttribute("attributes-natural-language",
		goipp.TagLanguage, goipp.String(DefaultNaturalLanguage)))

	if e.StatusMessage != "" {
		msg.Operation.Add(goipp.MakeAttribute("status-message",
			goipp.TagText, goipp.String(e.StatusMessage)))
	}

	return msg
}
