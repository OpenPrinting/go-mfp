// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Validate-Job request

package ipp

import (
	"github.com/OpenPrinting/goipp"
)

// ValidateJobRequest operation (0x0004) performs the print job
// validation as if the document was actually printed, bypassing
// the actual fetching and printing of the document data.
type ValidateJobRequest struct {
	ObjectRawAttrs
	RequestHeader

	// Operation attributes
	JobCreateOperation

	// Job Template attributes (RFC8011 Group 2)
	JobTemplate *JobTemplate
}

// ValidateJobResponse is the Validate-Job response.
type ValidateJobResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	// Unsupported attributes, if any
	UnsupportedAttributes goipp.Attributes
}

// GetOp returns ValidateJobRequest IPP Operation code.
func (rq *ValidateJobRequest) GetOp() goipp.Op {
	return goipp.OpValidateJob
}

// Encode encodes ValidateJobRequest into the goipp.Message.
func (rq *ValidateJobRequest) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rq),
		},

		{
			Tag:   goipp.TagJobGroup,
			Attrs: enc.Encode(rq.JobTemplate),
		},
	}

	msg := goipp.NewMessageWithGroups(rq.Version, goipp.Code(rq.GetOp()),
		rq.RequestID, groups)

	return msg
}

// Decode decodes ValidateJobRequest from goipp.Message.
func (rq *ValidateJobRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(rq, msg.Operation)
	if err != nil {
		return err
	}

	rq.JobTemplate, err = DecodeJobTemplate(msg.Job, opt)
	if err != nil {
		return err
	}

	return nil
}

// Encode encodes ValidateJobResponse into goipp.Message.
func (rsp *ValidateJobResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	if len(rsp.UnsupportedAttributes) > 0 {
		groups = append(groups, goipp.Group{
			Tag:   goipp.TagUnsupportedGroup,
			Attrs: rsp.UnsupportedAttributes,
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes ValidateJobResponse from goipp.Message.
func (rsp *ValidateJobResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)
	rsp.UnsupportedAttributes = msg.Unsupported

	return nil
}
