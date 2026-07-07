// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP response

package ipp

import (
	"io"

	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// Response is the IPP response interface.
type Response interface {
	// Header() returns *ResponseHeader.
	//
	// Each concrete Response implementation inherits it by
	// embedding this structure
	Header() *ResponseHeader

	// The following methods each concrete Response implementation
	// must define by itself:
	//   - Encode encodes Response into goipp.Message.
	//   - Decode decodes Response from goipp.Message.
	Encode() *goipp.Message
	Decode(*goipp.Message, *DecoderOptions) error
}

// ResponseHeader is the common [Response] header. It contains common
// fields and implements common interfaces.
//
// It should be embedded at the beginning of every structure that
// implements the [Response] interface.
type ResponseHeader struct {
	// IPP version, RequestID, IPP Status code.
	Version   goipp.Version
	RequestID uint32
	Status    goipp.Status

	// Common Operation attributes.
	AttributesCharset         string               `ipp:"attributes-charset"`
	AttributesNaturalLanguage string               `ipp:"attributes-natural-language"`
	StatusMessage             string               `ipp:"status-message"`
	DetailedStatusMessage     optional.Val[string] `ipp:"detailed-status-message"`

	// Raw IPP message.
	//
	// This field is filled when Response is received as result
	// of Client.Do or Client.DoWithBody and contains the raw
	// IPP message.
	IPPMessage *goipp.Message

	// Response Body.
	//
	// If Response is received as result of Client.DoWithBody,
	// the Body is never nil and always valid, even if empty.
	// Caller of Client.DoWithBody MUST Close the response
	// to avoid leak of resources.
	//
	// As IPP requests rarely returns non-empty Body, the
	// convenience wrapper, Client.Do, is provided. It
	// closes the Body and sets it to nil.
	Body io.ReadCloser
}

// Header returns [ResponseHeader], which gives uniform
// access to the header of any [Response]
func (rsph *ResponseHeader) Header() *ResponseHeader {
	return rsph
}
