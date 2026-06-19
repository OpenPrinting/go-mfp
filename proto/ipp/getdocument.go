// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Request and Response for Get-Next-Document-Data operation

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// GetNextDocumentDataRequest is the Get-Next-Document-Data request.
//
// The Get-Next-Document-Data operation allows a Scan
// Client to retrieve the scan data from an existing Job object, enabling
// pull scanning. The target Job MUST be in the 'processing' or 'completed'
// state.
type GetNextDocumentDataRequest struct {
	ObjectRawAttrs
	RequestHeader

	OperationGroup

	PrinterURI         optional.Val[string] `ipp:"printer-uri"`
	JobID              optional.Val[int]    `ipp:"job-id"`
	RequestingUserName optional.Val[string] `ipp:"requesting-user-name"`
	RequestingUserURI  optional.Val[string] `ipp:"requesting-user-uri"`
	DocumentDataWait   optional.Val[bool]   `ipp:"document-data-wait"`
}

// GetNextDocumentDataResponse is the Get-Next-Document-Data response.
type GetNextDocumentDataResponse struct {
	ObjectRawAttrs
	ResponseHeader

	OperationGroup

	Compression             optional.Val[KwCompression] `ipp:"compression"`
	DocumentFormat          optional.Val[string]        `ipp:"document-format"`
	DocumentDataGetInterval optional.Val[int]           `ipp:"document-data-get-interval"`
	LastDocument            bool                        `ipp:"last-document"`

	Document *DocumentStatus
}

type DocumentStatus struct {
	ObjectRawAttrs
	DocumentStatusGroup

	DocumentNumber optional.Val[int] `ipp:"document-number"`
}

// DecodeDocumentStatusAttributes decodes [DocumentStatus] from
// [goipp.Attributes].
func DecodeDocumentStatusAttributes(attrs goipp.Attributes,
	opt *DecoderOptions) (*DocumentStatus, error) {

	doc := &DocumentStatus{}
	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(doc, attrs)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// GetOp returns GetNextDocumentDataRequest IPP Operation code.
func (rq *GetNextDocumentDataRequest) GetOp() goipp.Op {
	return goipp.OpGetNextDocumentData
}

// Encode encodes GetNextDocumentDataRequest into the goipp.Message.
func (rq *GetNextDocumentDataRequest) Encode() *goipp.Message {
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

// Decode decodes GetNextDocumentDataRequest from goipp.Message.
func (rq *GetNextDocumentDataRequest) Decode(
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

// Encode encodes GetNextDocumentDataResponse into the goipp.Message.
func (rsp *GetNextDocumentDataResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	if rsp.Document != nil {
		groups = append(groups, goipp.Group{
			Tag:   goipp.TagDocumentGroup,
			Attrs: enc.Encode(rsp.Document),
		})
	}

	return goipp.NewMessageWithGroups(
		rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups,
	)
}

// Decode decodes GetNextDocumentDataResponse from goipp.Message.
func (rsp *GetNextDocumentDataResponse) Decode(
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

	rsp.Document, err = DecodeDocumentStatusAttributes(msg.Document, opt)
	if err != nil {
		return err
	}

	return nil
}
