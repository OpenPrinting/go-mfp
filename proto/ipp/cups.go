// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// CUPS requests and responses

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

type (
	// CUPSGetDefaultRequest operation (0x4001) returns the default printer URI
	// and attributes.
	CUPSGetDefaultRequest struct {
		ObjectRawAttrs
		RequestHeader
		OperationGroup

		// Operation attributes
		RequestedAttributes []string `ipp:"requested-attributes"`
	}

	// CUPSGetDefaultResponse is the CUPS-Get-Default Response.
	CUPSGetDefaultResponse struct {
		ObjectRawAttrs
		ResponseHeader
		OperationGroup

		// Other attributes.
		Printer *PrinterAttributes
	}

	// CUPSGetPrintersRequest operation (0x4002) returns the printer
	// attributes for every printer known to the system.
	CUPSGetPrintersRequest struct {
		ObjectRawAttrs
		RequestHeader
		OperationGroup

		// Operation attributes
		FirstPrinterName    optional.Val[string] `ipp:"first-printer-name"`
		Limit               optional.Val[int]    `ipp:"limit"`
		PrinterID           optional.Val[int]    `ipp:"printer-id"`
		PrinterLocation     optional.Val[string] `ipp:"printer-location"`
		PrinterType         optional.Val[int]    `ipp:"printer-type"`
		PrinterTypeMask     optional.Val[int]    `ipp:"printer-type-mask"`
		RequestedAttributes []string             `ipp:"requested-attributes"`
		RequestedUserName   optional.Val[string] `ipp:"requested-user-name,name"`
	}

	// CUPSGetPrintersResponse is the CUPS-Get-Printers Response.
	CUPSGetPrintersResponse struct {
		ObjectRawAttrs
		ResponseHeader
		OperationGroup

		// Other attributes.
		Printer []*PrinterAttributes
	}

	// CUPSGetDevicesRequest operation (0x400b) performs search
	// for available printers and returns all of the supported
	// device-uri's
	CUPSGetDevicesRequest struct {
		ObjectRawAttrs
		RequestHeader
		OperationGroup

		// Operational attributes
		ExcludeSchemes      []string          `ipp:"exclude-schemes"`
		IncludeSchemes      []string          `ipp:"include-schemes"`
		Limit               optional.Val[int] `ipp:"limit"`
		RequestedAttributes []string          `ipp:"requested-attributes"`
		Timeout             optional.Val[int] `ipp:"timeout"`
	}

	// CUPSGetDevicesResponse is the CUPS-Get-Devices Response.
	CUPSGetDevicesResponse struct {
		ObjectRawAttrs
		ResponseHeader
		OperationGroup

		// Other attributes.
		Printer []*DeviceAttributes
	}

	// CUPSGetPPDsRequest operation (0x400c) returns list of available PPDs
	CUPSGetPPDsRequest struct {
		ObjectRawAttrs
		RequestHeader

		OperationGroup
		CUPSPPDAttributesGroup

		// Operational attributes
		ExcludeSchemes      []string             `ipp:"exclude-schemes"`
		IncludeSchemes      []string             `ipp:"include-schemes"`
		Limit               optional.Val[int]    `ipp:"limit"`
		PPDMake             optional.Val[string] `ipp:"ppd-make"`
		PPDMakeAndModel     optional.Val[string] `ipp:"ppd-make-and-model"`
		PPDModelNumber      optional.Val[int]    `ipp:"ppd-model-number"`
		PPDNaturalLanguage  optional.Val[string] `ipp:"ppd-natural-language"`
		PPDProduct          optional.Val[string] `ipp:"ppd-product"`
		PPDPsVersion        optional.Val[string] `ipp:"ppd-psversion"`
		PPDType             optional.Val[string] `ipp:"ppd-type,keyword"`
		RequestedAttributes []string             `ipp:"requested-attributes"`
	}

	// CUPSGetPPDsResponse is the CUPS-Get-PPDs Response.
	CUPSGetPPDsResponse struct {
		ObjectRawAttrs
		ResponseHeader
		OperationGroup

		// Other attributes.
		PPDs []*PPDAttributes
	}

	// CUPSGetPPDRequest operation (0x400f) returns PPD file from
	// the server.
	CUPSGetPPDRequest struct {
		ObjectRawAttrs
		RequestHeader

		OperationGroup
		CUPSPPDAttributesGroup

		// Operational attributes
		//
		// Use PrinterURI to specify particular print queue
		// or PPDName to request PPD file by its name.
		PrinterURI optional.Val[string] `ipp:"printer-uri"`
		PPDName    optional.Val[string] `ipp:"ppd-name"`
	}

	// CUPSGetPPDResponse is the CUPS-Get-PPD Response.
	//
	// If the PPD file is found, goipp.StatusOk is returned with the PPD
	// file represented by the ResponseHeader.Body.
	//
	// If the PPD file cannot be served by the local server because the
	// printer-uri attribute points to an external printer, a
	// goipp.StatusCupsSeeOther is returned and PrinterURI contains
	// the correct URI to use.
	//
	// If the PPD file does not exist, goipp.StatusErrorNotFound is
	// returned.
	CUPSGetPPDResponse struct {
		ObjectRawAttrs
		ResponseHeader
		OperationGroup

		// Operational attributes
		PrinterURI optional.Val[string] `ipp:"printer-uri"`
	}
)

// ----- CUPS-Get-Default methods -----

// GetOp returns CUPSGetDefaultRequest IPP Operation code.
func (rq *CUPSGetDefaultRequest) GetOp() goipp.Op {
	return goipp.OpCupsGetDefault
}

// Encode encodes CUPSGetDefaultRequest into the goipp.Message.
func (rq *CUPSGetDefaultRequest) Encode() *goipp.Message {
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

// Decode decodes CUPSGetDefaultRequest from goipp.Message.
func (rq *CUPSGetDefaultRequest) Decode(
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

// Encode encodes CUPSGetDefaultResponse into goipp.Message.
func (rsp *CUPSGetDefaultResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	if rsp.Printer != nil {
		groups.Add(goipp.Group{
			Tag:   goipp.TagPrinterGroup,
			Attrs: enc.Encode(rsp.Printer),
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes CUPSGetDefaultResponse from goipp.Message.
func (rsp *CUPSGetDefaultResponse) Decode(
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

	if len(msg.Printer) != 0 {
		rsp.Printer, err = DecodePrinterAttributes(msg.Printer, opt)
		if err != nil {
			return err
		}
	}

	return nil
}

// ----- CUPS-Get-Printers methods -----

// GetOp returns CUPSGetPrintersRequest IPP Operation code.
func (rq *CUPSGetPrintersRequest) GetOp() goipp.Op {
	return goipp.OpCupsGetPrinters
}

// Encode encodes CUPSGetPrintersRequest into the goipp.Message.
func (rq *CUPSGetPrintersRequest) Encode() *goipp.Message {
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

// Decode decodes CUPSGetPrintersRequest from goipp.Message.
func (rq *CUPSGetPrintersRequest) Decode(
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

// Encode encodes CUPSGetPrintersResponse into goipp.Message.
func (rsp *CUPSGetPrintersResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	for _, prn := range rsp.Printer {
		groups.Add(goipp.Group{
			Tag:   goipp.TagPrinterGroup,
			Attrs: enc.Encode(prn),
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes CUPSGetPrintersResponse from goipp.Message.
func (rsp *CUPSGetPrintersResponse) Decode(
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

	for _, grp := range msg.Groups {
		if grp.Tag == goipp.TagPrinterGroup && len(grp.Attrs) > 0 {
			prn, err := DecodePrinterAttributes(grp.Attrs, opt)
			if err != nil {
				return err
			}

			rsp.Printer = append(rsp.Printer, prn)
		}
	}

	return nil
}

// ----- CUPS-Get-Devices methods -----

// GetOp returns CUPSGetDevicesRequest IPP Operation code.
func (rq *CUPSGetDevicesRequest) GetOp() goipp.Op {
	return goipp.OpCupsGetDevices
}

// Encode encodes CUPSGetDevicesRequest into the goipp.Message.
func (rq *CUPSGetDevicesRequest) Encode() *goipp.Message {
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

// Decode decodes CUPSGetDevicesRequest from goipp.Message.
func (rq *CUPSGetDevicesRequest) Decode(
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

// Encode encodes CUPSGetDevicesResponse into goipp.Message.
func (rsp *CUPSGetDevicesResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	for _, prn := range rsp.Printer {
		groups.Add(goipp.Group{
			Tag:   goipp.TagPrinterGroup,
			Attrs: enc.Encode(prn),
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes CUPSGetDevicesResponse from goipp.Message.
func (rsp *CUPSGetDevicesResponse) Decode(
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

	for _, grp := range msg.Groups {
		if grp.Tag == goipp.TagPrinterGroup && len(grp.Attrs) > 0 {
			dev := &DeviceAttributes{}
			err = dec.Decode(dev, grp.Attrs)
			if err != nil {
				return err
			}

			rsp.Printer = append(rsp.Printer, dev)
		}
	}

	return nil
}

// ----- CUPS-Get-PPDs methods -----

// GetOp returns CUPSGetPPDsRequest IPP Operation code.
func (rq *CUPSGetPPDsRequest) GetOp() goipp.Op {
	return goipp.OpCupsGetPpds
}

// Encode encodes CUPSGetPPDsRequest into the goipp.Message.
func (rq *CUPSGetPPDsRequest) Encode() *goipp.Message {
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

// Decode decodes CUPSGetPPDsRequest from goipp.Message.
func (rq *CUPSGetPPDsRequest) Decode(
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

// Encode encodes CUPSGetPPDsResponse into goipp.Message.
func (rsp *CUPSGetPPDsResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	for _, ppd := range rsp.PPDs {
		groups.Add(goipp.Group{
			Tag:   goipp.TagPrinterGroup,
			Attrs: enc.Encode(ppd),
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// Decode decodes CUPSGetPPDsResponse from goipp.Message.
func (rsp *CUPSGetPPDsResponse) Decode(
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

	for _, grp := range msg.Groups {
		if grp.Tag == goipp.TagPrinterGroup && len(grp.Attrs) > 0 {
			ppd := &PPDAttributes{}
			err = dec.Decode(ppd, grp.Attrs)
			if err != nil {
				return err
			}

			rsp.PPDs = append(rsp.PPDs, ppd)
		}
	}

	return nil
}

// ----- CUPS-Get-PPD methods -----

// GetOp returns CUPSGetPPDRequest IPP Operation code.
func (rq *CUPSGetPPDRequest) GetOp() goipp.Op {
	return goipp.OpCupsGetPpd
}

// Encode encodes CUPSGetPPDRequest into the goipp.Message.
func (rq *CUPSGetPPDRequest) Encode() *goipp.Message {
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

// Decode decodes CUPSGetPPDRequest from goipp.Message.
func (rq *CUPSGetPPDRequest) Decode(
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

// Encode encodes CUPSGetPPDResponse into goipp.Message.
func (rsp *CUPSGetPPDResponse) Encode() *goipp.Message {
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

// Decode decodes CUPSGetPPDResponse from goipp.Message.
func (rsp *CUPSGetPPDResponse) Decode(
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
