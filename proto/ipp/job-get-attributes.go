// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Get-Job-Attributes request and response

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// GetJobAttributesRequest operation (0x0009) returns Job attributes.
type GetJobAttributesRequest struct {
	ObjectRawAttrs
	RequestHeader
	OperationGroup

	PrinterURI          optional.Val[string]   `ipp:"printer-uri"`
	JobID               optional.Val[int]      `ipp:"job-id"`
	JobURI              optional.Val[string]   `ipp:"job-uri"`
	RequestingUserName  optional.Val[string]   `ipp:"requesting-user-name"`
	RequestedAttributes []KwRequestedAttribute `ipp:"requested-attributes"`
}

// GetJobAttributesResponse is the Get-Job-Attributes response.
type GetJobAttributesResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	UnsupportedAttributes goipp.Attributes
	Job                   JobGroupEntry
}

// GetOp returns GetJobAttributesRequest IPP Operation code.
func (rq *GetJobAttributesRequest) GetOp() goipp.Op {
	return goipp.OpGetJobAttributes
}

// Encode encodes GetJobAttributesRequest into the goipp.Message.
func (rq *GetJobAttributesRequest) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rq),
		},
	}

	return goipp.NewMessageWithGroups(
		rq.Version, goipp.Code(rq.GetOp()),
		rq.RequestID, groups,
	)
}

// Decode decodes GetJobAttributesRequest from goipp.Message.
func (rq *GetJobAttributesRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	return dec.Decode(rq, msg.Operation)
}

// Encode encodes GetJobAttributesResponse into the goipp.Message.
func (rsp *GetJobAttributesResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	return rsp.EncodeRaw(enc.Encode(&rsp.Job))
}

// EncodeRaw encodes GetJobAttributesResponse using pre-built Job attributes.
func (rsp *GetJobAttributesResponse) EncodeRaw(
	jobAttrs goipp.Attributes) *goipp.Message {

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

	if len(jobAttrs) > 0 {
		groups = append(groups, goipp.Group{
			Tag:   goipp.TagJobGroup,
			Attrs: jobAttrs,
		})
	}

	return goipp.NewMessageWithGroups(
		rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups,
	)
}

// Decode decodes GetJobAttributesResponse from goipp.Message.
func (rsp *GetJobAttributesResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)
	rsp.UnsupportedAttributes = msg.Unsupported

	dec := NewDecoder(opt)
	defer dec.Free()

	if err := dec.Decode(rsp, msg.Operation); err != nil {
		return err
	}

	if len(msg.Job) == 0 {
		return nil
	}

	var err error
	rsp.Job, err = DecodeJobGroupEntry(msg.Job, opt)
	return err
}

// Apply applies the Get-Job-Attributes request to jobs and returns the response.
func (rq *GetJobAttributesRequest) Apply(jobs []*job) (*goipp.Message, error) {
	j, err := rq.lookupJob(jobs)
	if err != nil {
		return nil, err
	}

	requested := rq.effectiveRequestedAttributes()
	requestedNames := make([]string, len(requested))
	for i, attr := range requested {
		requestedNames[i] = string(attr)
	}
	unsupported := unsupportedJobAttributes(requestedNames)

	j.Lock()
	enc := ippEncoder{}
	encoded := enc.Encode(&JobGroupEntry{
		JobDescriptionAttrs: j.JobDescriptionAttrs,
		JobStatusAttrs:      j.JobStatusAttrs,
		JobTemplateAttrs:    j.JobTemplateAttrs,
	})
	j.Unlock()

	filtered, _ := jobAttributesFilter.Apply(requestedNames, encoded)

	status := goipp.StatusOk
	if len(unsupported) > 0 {
		status = goipp.StatusOkIgnoredOrSubstituted
	}

	rsp := &GetJobAttributesResponse{
		ResponseHeader: rq.ResponseHeader(status),
		UnsupportedAttributes: unsupportedAttrsFromNames(
			"requested-attributes", unsupported),
	}

	return rsp.EncodeRaw(filtered), nil
}

func (rq *GetJobAttributesRequest) lookupJob(jobs []*job) (*job, error) {
	switch {
	case rq.PrinterURI != nil && rq.JobID != nil:
		for _, j := range jobs {
			if j.JobID == *rq.JobID {
				return j, nil
			}
		}
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotFound,
			"job not found (job-id=%d)", *rq.JobID)

	case rq.JobURI != nil:
		for _, j := range jobs {
			if j.JobURI == *rq.JobURI {
				return j, nil
			}
		}
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotFound,
			"job not found (job-uri=%q)", *rq.JobURI)

	default:
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missing job-id and job-uri attributes")
	}
}

func (rq *GetJobAttributesRequest) effectiveRequestedAttributes() []KwRequestedAttribute {
	if len(rq.RequestedAttributes) > 0 {
		return rq.RequestedAttributes
	}
	return []KwRequestedAttribute{KwRequestedAttributeAll}
}
