// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Get-Jobs request and response

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// GetJobsRequest operation returns a list of jobs.
type GetJobsRequest struct {
	ObjectRawAttrs
	RequestHeader
	OperationGroup

	PrinterURI          string                    `ipp:"printer-uri"`
	RequestingUserName  optional.Val[string]      `ipp:"requesting-user-name"`
	Limit               optional.Val[int]         `ipp:"limit"`
	RequestedAttributes []KwRequestedAttribute    `ipp:"requested-attributes"`
	WhichJobs           optional.Val[KwWhichJobs] `ipp:"which-jobs"`
	MyJobs              optional.Val[bool]        `ipp:"my-jobs"`
}

// GetJobsResponse is the Get-Jobs response.
type GetJobsResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	UnsupportedAttributes goipp.Attributes
	Jobs                  []JobGroupEntry
}

// GetOp returns GetJobsRequest IPP Operation code.
func (rq *GetJobsRequest) GetOp() goipp.Op {
	return goipp.OpGetJobs
}

// Encode encodes GetJobsRequest into the goipp.Message.
func (rq *GetJobsRequest) Encode() *goipp.Message {
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

// Decode decodes GetJobsRequest from goipp.Message.
func (rq *GetJobsRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	return dec.Decode(rq, msg.Operation)
}

// Encode encodes GetJobsResponse into the goipp.Message.
func (rsp *GetJobsResponse) Encode() *goipp.Message {
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

	for i := range rsp.Jobs {
		groups = append(groups, goipp.Group{
			Tag:   goipp.TagJobGroup,
			Attrs: enc.Encode(&rsp.Jobs[i]),
		})
	}

	return goipp.NewMessageWithGroups(
		rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups,
	)
}

// Decode decodes GetJobsResponse from goipp.Message.
func (rsp *GetJobsResponse) Decode(
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

	groups := msg.AttrGroups()

	for _, grp := range groups {
		if grp.Tag != goipp.TagJobGroup {
			continue
		}

		entry, err := DecodeJobGroupEntry(grp.Attrs, opt)
		if err != nil {
			return err
		}
		rsp.Jobs = append(rsp.Jobs, entry)
	}

	return nil
}

// Apply applies the Get-Jobs request to jobs and returns the encoded response.
func (rq *GetJobsRequest) Apply(jobs []*job) *goipp.Message {
	whichJobs := optional.Get(rq.WhichJobs)
	if whichJobs == "" {
		whichJobs = KwWhichJobsNotCompleted
	}

	if whichJobs != KwWhichJobsCompleted && whichJobs != KwWhichJobsNotCompleted {
		rsp := &GetJobsResponse{
			ResponseHeader: rq.ResponseHeader(goipp.StatusErrorAttributesOrValues),
			UnsupportedAttributes: goipp.Attributes{
				goipp.MakeAttribute("which-jobs",
					goipp.TagKeyword, goipp.String(string(whichJobs))),
			},
		}
		return rsp.Encode()
	}

	requestedAttributes := rq.RequestedAttributes

	if len(requestedAttributes) == 0 {
		requestedAttributes = []KwRequestedAttribute{
			KwRequestedAttributeJobURI,
			KwRequestedAttributeJobID,
		}
	}

	requestedNames := make([]string, len(requestedAttributes))

	for i, attr := range requestedAttributes {
		requestedNames[i] = string(attr)
	}

	unsupported := unsupportedJobAttributes(requestedNames)

	myJobs := optional.Get(rq.MyJobs)
	user := optional.Get(rq.RequestingUserName)

	matched := make([]*job, 0, len(jobs))
	for _, j := range jobs {
		var match bool
		switch whichJobs {
		case KwWhichJobsCompleted:
			match = j.JobState == EnJobStateCompleted ||
				j.JobState == EnJobStateCanceled ||
				j.JobState == EnJobStateAborted
		case KwWhichJobsNotCompleted:
			match = j.JobState == EnJobStatePending ||
				j.JobState == EnJobStatePendingHeld ||
				j.JobState == EnJobStateProcessing ||
				j.JobState == EnJobStateProcessingStopped
		}
		if !match {
			continue
		}
		if myJobs && (user == "" || j.JobOriginatingUserName == nil ||
			optional.Get(j.JobOriginatingUserName) != user) {
			continue
		}
		matched = append(matched, j)
	}

	limit := optional.Get(rq.Limit)
	if limit > 0 && len(matched) > limit {
		matched = matched[:limit]
	}

	enc := ippEncoder{}
	jobGroups := make(goipp.Groups, 0, len(matched))
	for _, j := range matched {
		j.Lock()
		encoded := enc.Encode(&JobGroupEntry{
			JobDescriptionAttrs: j.JobDescriptionAttrs,
			JobStatusAttrs:      j.JobStatusAttrs,
			JobTemplateAttrs:    j.JobTemplateAttrs,
		})
		j.Unlock()

		filtered, _ := jobAttributesFilter.Apply(requestedNames, encoded)
		jobGroups = append(jobGroups, goipp.Group{
			Tag:   goipp.TagJobGroup,
			Attrs: filtered,
		})
	}

	status := goipp.StatusOk
	if len(unsupported) > 0 {
		status = goipp.StatusOkIgnoredOrSubstituted
	}

	rsp := &GetJobsResponse{
		ResponseHeader: rq.ResponseHeader(status),
		UnsupportedAttributes: unsupportedAttrsFromNames(
			"requested-attributes", unsupported),
	}

	return rsp.encodeWithJobGroups(jobGroups)
}

func (rsp *GetJobsResponse) encodeWithJobGroups(jobGroups goipp.Groups) *goipp.Message {
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

	groups = append(groups, jobGroups...)

	return goipp.NewMessageWithGroups(
		rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups,
	)
}

func unsupportedJobAttributes(requested []string) []string {
	all := jobAttrGroups["all"]
	seen := generic.NewSet[string]()
	var unsupported []string

	for _, name := range requested {
		if _, ok := jobAttrGroups[name]; ok {
			continue
		}
		if all.Contains(name) {
			continue
		}
		if seen.TestAndAdd(name) {
			unsupported = append(unsupported, name)
		}
	}

	return unsupported
}

func unsupportedAttrsFromNames(attr string, names []string) goipp.Attributes {
	if len(names) == 0 {
		return nil
	}

	values := make(goipp.Values, 0, len(names))
	for _, name := range names {
		values.Add(goipp.TagKeyword, goipp.String(name))
	}

	return goipp.Attributes{
		{Name: attr, Values: values},
	}
}
