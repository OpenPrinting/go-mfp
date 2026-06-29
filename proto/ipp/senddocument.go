// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Send-Document request and response

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// SendDocumentRequest operation (0x0006) adds document to
// the job previously created by the Create-Job request.
type SendDocumentRequest struct {
	ObjectRawAttrs
	RequestHeader

	OperationGroup

	// Operation attributes
	PrinterURI         optional.Val[string] `ipp:"printer-uri"`
	JobID              optional.Val[int]    `ipp:"job-id"`
	JobURI             optional.Val[string] `ipp:"job-uri"`
	RequestingUserName optional.Val[string] `ipp:"requesting-user-name"`

	Compression             optional.Val[KwCompression] `ipp:"compression"`
	DocumentFormat          optional.Val[string]        `ipp:"document-format"`
	DocumentName            optional.Val[string]        `ipp:"document-name"`
	DocumentNaturalLanguage optional.Val[string]        `ipp:"document-natural-language"`
	LastDocument            bool                        `ipp:"last-document"`

	// Job attributes
	Job *JobAttributes
}

// SendDocumentResponse is the Create-Job response.
type SendDocumentResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	// Unsupported attributes, if any
	UnsupportedAttributes goipp.Attributes

	// Job status
	Job *JobStatus
}

// GetOp returns SendDocumentRequest IPP Operation code.
func (rq *SendDocumentRequest) GetOp() goipp.Op {
	return goipp.OpSendDocument
}

// Encode encodes SendDocumentRequest into the goipp.Message.
func (rq *SendDocumentRequest) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rq),
		},

		{
			Tag:   goipp.TagJobGroup,
			Attrs: enc.Encode(rq.Job),
		},
	}

	msg := goipp.NewMessageWithGroups(rq.Version, goipp.Code(rq.GetOp()),
		rq.RequestID, groups)

	return msg
}

// Decode decodes SendDocumentRequest from goipp.Message.
func (rq *SendDocumentRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(rq, msg.Operation)
	if err != nil {
		return err
	}

	rq.Job, err = DecodeJobAttributes(msg.Job, opt)
	if err != nil {
		return err
	}

	return nil
}

// Encode encodes SendDocumentResponse into goipp.Message.
func (rsp *SendDocumentResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	if rsp.Job != nil {
		groups = append(groups, goipp.Group{
			Tag:   goipp.TagJobGroup,
			Attrs: enc.Encode(rsp.Job),
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes SendDocumentResponse from goipp.Message.
func (rsp *SendDocumentResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)
	rsp.UnsupportedAttributes = msg.Unsupported

	var err error
	rsp.Job, err = DecodeJobStatusAttributes(msg.Job, opt)
	if err != nil {
		return err
	}

	return nil
}
