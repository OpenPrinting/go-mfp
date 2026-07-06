// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Cancel-Job request and response

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// JobCancelOperation contains operation attributes common for
// job cancellation requests.
type JobCancelOperation struct {
	OperationGroup

	PrinterURI         optional.Val[string] `ipp:"printer-uri"`
	JobID              optional.Val[int]    `ipp:"job-id"`
	JobURI             optional.Val[string] `ipp:"job-uri"`
	RequestingUserName optional.Val[string] `ipp:"requesting-user-name"`
	Message            optional.Val[string] `ipp:"message"`
}

// CancelJobRequest operation (0x0008) cancels a Job.
type CancelJobRequest struct {
	ObjectRawAttrs
	RequestHeader

	JobCancelOperation
}

// CancelJobResponse is the Cancel-Job response.
type CancelJobResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	// Unsupported attributes, if any
	UnsupportedAttributes goipp.Attributes
}

// GetOp returns CancelJobRequest IPP Operation code.
func (rq *CancelJobRequest) GetOp() goipp.Op {
	return goipp.OpCancelJob
}

// Encode encodes CancelJobRequest into the goipp.Message.
func (rq *CancelJobRequest) Encode() *goipp.Message {
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

// Decode decodes CancelJobRequest from goipp.Message.
func (rq *CancelJobRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	return dec.Decode(rq, msg.Operation)
}

// Encode encodes CancelJobResponse into the goipp.Message.
func (rsp *CancelJobResponse) Encode() *goipp.Message {
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

	return goipp.NewMessageWithGroups(
		rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups,
	)
}

// Decode decodes CancelJobResponse from goipp.Message.
func (rsp *CancelJobResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)
	rsp.UnsupportedAttributes = msg.Unsupported

	dec := NewDecoder(opt)
	defer dec.Free()

	return dec.Decode(rsp, msg.Operation)
}

// beginCancel accepts Cancel-Job for this job.
//
// If the job is still processing, it remains in that state with
// processing-to-stop-point in job-state-reasons until finishCancel
// is called. Otherwise the job transitions immediately to canceled.
func (j *job) beginCancel() {
	switch {
	case j.JobState == EnJobStateProcessing,
		j.JobState == EnJobStateProcessingStopped,
		j.SendDocumentActive:
		j.cancelPending = true
		j.JobStateReasons = append(
			j.JobStateReasons, KwJobStateReasonsProcessingToStopPoint)
	default:
		j.cancelPending = false
		j.JobState = EnJobStateCanceled
		j.JobStateReasons = []KwJobStateReasons{
			KwJobStateReasonsJobCanceledByUser,
		}
	}
}

// finishCancel transitions a cancel-pending job to the canceled state.
func (j *job) finishCancel() {
	if !j.cancelPending {
		return
	}
	j.cancelPending = false
	j.JobState = EnJobStateCanceled
	j.JobStateReasons = []KwJobStateReasons{
		KwJobStateReasonsJobCanceledByUser,
	}
}

// isCancelable reports whether the job can be canceled.
func (j *job) isCancelable() bool {
	switch j.JobState {
	case EnJobStatePending, EnJobStatePendingHeld,
		EnJobStateProcessing, EnJobStateProcessingStopped:
		return true
	default:
		return false
	}
}
