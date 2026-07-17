// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Create-Job request

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// CreateJobRequest operation (0x0005) creates a new print Job.
// This operation requires that the document data is supplied
// by the client separately, using Send-Document or Send-URI
// operations.
type CreateJobRequest struct {
	ObjectRawAttrs
	RequestHeader

	// Operation attributes
	JobCreateOperation

	// Job Template attributes (RFC8011 Group 2)
	JobTemplate *JobTemplate
}

// CreateJobResponse is the Create-Job response.
type CreateJobResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	// PWG5100.17, 7.1.2: compression of the delivered document data.
	Compression optional.Val[KwCompression] `ipp:"compression"`

	// Unsupported attributes, if any
	UnsupportedAttributes goipp.Attributes

	// Job status
	Job *JobDescriptionAndStatus
}

// GetOp returns CreateJobRequest IPP Operation code.
func (rq *CreateJobRequest) GetOp() goipp.Op {
	return goipp.OpCreateJob
}

// Encode encodes CreateJobRequest into the goipp.Message.
func (rq *CreateJobRequest) Encode() *goipp.Message {
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

// Decode decodes CreateJobRequest from goipp.Message.
func (rq *CreateJobRequest) Decode(
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

// Encode encodes CreateJobResponse into goipp.Message.
func (rsp *CreateJobResponse) Encode() *goipp.Message {
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

// Decode decodes CreateJobResponse from goipp.Message.
func (rsp *CreateJobResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)
	rsp.UnsupportedAttributes = msg.Unsupported

	var err error
	rsp.Job, err = DecodeJobDescriptionAndStatus(msg.Job, opt)
	if err != nil {
		return err
	}

	return nil
}
