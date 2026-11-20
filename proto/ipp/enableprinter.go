// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions
//
// Enable-Printer request

package ipp

import (
	"github.com/OpenPrinting/goipp"
)

// EnablePrinterRequest operation (0x0022) makes the printer accept
// new print jobs.
type EnablePrinterRequest struct {
	ObjectRawAttrs
	RequestHeader
	OperationGroup

	// Operation attributes
	PrinterURI string `ipp:"printer-uri"`
}

// EnablePrinterResponse is the Enable-Printer response.
type EnablePrinterResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup
}

// GetOp returns EnablePrinterRequest IPP Operation code.
func (rq *EnablePrinterRequest) GetOp() goipp.Op {
	return goipp.OpEnablePrinter
}

// Encode encodes EnablePrinterRequest into the goipp.Message.
func (rq *EnablePrinterRequest) Encode() *goipp.Message {
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

// Decode decodes EnablePrinterRequest from goipp.Message.
func (rq *EnablePrinterRequest) Decode(
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

// Encode encodes EnablePrinterResponse into goipp.Message.
func (rsp *EnablePrinterResponse) Encode() *goipp.Message {
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

// Decode decodes EnablePrinterResponse from goipp.Message.
func (rsp *EnablePrinterResponse) Decode(
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
