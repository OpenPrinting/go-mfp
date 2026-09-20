// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP request

package ipp

import (
	"io"

	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/goipp"
)

// Request is the IPP request interface.
type Request interface {
	// Header() returns *RequestHeader.
	//
	// Each concrete Request implementation inherits it by
	// embedding this structure.
	Header() *RequestHeader

	// The following methods each concrete Request implementation
	// must define by itself:
	//   - GetOp returns IPP Operation code of the Request.
	//   - Encode encodes Request into the goipp.Message.
	//   - Decode decodes Request from the goipp.Message.
	GetOp() goipp.Op
	Encode() *goipp.Message
	Decode(*goipp.Message, *DecoderOptions) error
}

// RequestHeader is the common [Request] header. It contains common
// fields and implements common interfaces.
//
// It should be embedded at the beginning of every structure that
// implements the [Request] interface.
type RequestHeader struct {
	// IPP version and RequestID.
	Version   goipp.Version
	RequestID uint32

	// Common Operation attributes
	AttributesCharset         string `ipp:"attributes-charset"`
	AttributesNaturalLanguage string `ipp:"attributes-natural-language"`

	// Request body. Sent after IPP message. May be nil
	Body io.Reader
}

// Header returns [RequestHeader], which gives uniform
// access to the header of any [Request]
func (rqh *RequestHeader) Header() *RequestHeader {
	return rqh
}

// ResponseHeader returns the appropriate [ResponseHeader]
// for the request.
//
// The response version matches the request version (RFC 8011,
// 4.1.8), limited by the [MaxVersion].
func (rqh *RequestHeader) ResponseHeader(status goipp.Status) ResponseHeader {
	ver := generic.Min(rqh.Version, MaxVersion)
	if ver == 0 {
		ver = DefaultVersion
	}

	return ResponseHeader{
		Version:                   ver,
		RequestID:                 rqh.RequestID,
		Status:                    status,
		StatusMessage:             status.String(),
		AttributesCharset:         DefaultCharset,
		AttributesNaturalLanguage: DefaultNaturalLanguage,
	}
}

// setBody sets the RequestHeader.Body
func (rqh *RequestHeader) setBody(body io.Reader) {
	rqh.Body = body
}
