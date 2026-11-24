// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions
//
// Disable-Printer request

package ipp

import (
	"github.com/OpenPrinting/goipp"
)

// DisablePrinterRequest operation (0x0023) makes the printer reject
// new print jobs.
type DisablePrinterRequest struct {
	ObjectRawAttrs
	RequestHeader
	OperationGroup

	// Operation attributes
	PrinterURI string `ipp:"printer-uri"`
}

// DisablePrinterResponse is the Disable-Printer response.
type DisablePrinterResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup
}

// GetOp returns DisablePrinterRequest IPP Operation code.
func (rq *DisablePrinterRequest) GetOp() goipp.Op {
	return goipp.OpDisablePrinter
}

// Encode encodes DisablePrinterRequest into the goipp.Message.
func (rq *DisablePrinterRequest) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rq),
		},
	}

	msg := goipp.NewMessageWithGroups(rq.Version, goipp.Code(rq.GetOp()),
		rq.RequestID, groups)

	return msg
}

// Decode decodes DisablePrinterRequest from goipp.Message.
func (rq *DisablePrinterRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(rq, msg.Operation)
	if err != nil {
		return err
	}

	return nil
}

// Encode encodes DisablePrinterResponse into goipp.Message.
func (rsp *DisablePrinterResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes DisablePrinterResponse from goipp.Message.
func (rsp *DisablePrinterResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)

	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(rsp, msg.Operation)
	if err != nil {
		return err
	}

	return nil
}
