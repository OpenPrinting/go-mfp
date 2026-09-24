// MFP - Miulti-Function Printers and scanners toolkit
// CUPS Client and Server
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// CUPS Client

package cups

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// Client represents the CUPS client.
type Client struct {
	IPPClient *ipp.Client // Underlying IPP client
}

// NewClient creates a new CUPS client.
//
// If tr is nil, [transport.NewTransport] will be used to create
// a new transport.
func NewClient(u *url.URL, tr *transport.Transport) *Client {
	return &Client{
		IPPClient: ipp.NewClient(u, tr),
	}
}

// SetDecoderOptions updates the [DecoderrOptions] that affect decoding
// of the received IPP messages
func (c *Client) SetDecoderOptions(opt *ipp.DecoderOptions) {
	c.IPPClient.SetDecoderOptions(opt)
}

// CUPSGetDefault returns information on default printer.
// The attrs attribute allows to specify list of requested attributes.
func (c *Client) CUPSGetDefault(ctx context.Context,
	attrs []string) (*ipp.PrinterAttributes, error) {

	rq := &ipp.CUPSGetDefaultRequest{
		RequestHeader:       ipp.DefaultRequestHeader,
		RequestedAttributes: attrs,
	}

	rsp := &ipp.CUPSGetDefaultResponse{}

	err := c.IPPClient.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Printer, nil
}

// CUPSGetPrinters returns printer attributes for printers known
// to the system.
//
// If [GetPrintersSelection] argument is not nil, it allows to
// specify a subset of printers to be returned.
//
// The attrs attribute allows to specify list of requested attributes.
func (c *Client) CUPSGetPrinters(ctx context.Context,
	sel *GetPrintersSelection, attrs []string) (
	[]*ipp.PrinterAttributes, error) {

	if sel == nil {
		sel = DefaultGetPrintersSelection
	}

	rq := &ipp.CUPSGetPrintersRequest{
		RequestHeader:       ipp.DefaultRequestHeader,
		FirstPrinterName:    optional.NotZero(sel.FirstPrinterName),
		Limit:               optional.NotZero(sel.Limit),
		PrinterID:           optional.NotZero(sel.PrinterID),
		PrinterLocation:     optional.NotZero(sel.PrinterLocation),
		PrinterType:         optional.NotZero(sel.PrinterType),
		PrinterTypeMask:     optional.NotZero(sel.PrinterTypeMask),
		RequestedUserName:   optional.NotZero(sel.User),
		RequestedAttributes: attrs,
	}

	rsp := &ipp.CUPSGetPrintersResponse{}

	err := c.IPPClient.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Printer, nil
}

// CUPSGetDevices performs search for available devices and returns
// found devices.
//
// If [GetDevicesSelection] argument is not nil, it allows to
// specify a subset of devices to be returned.
//
// The attrs attribute allows to specify list of requested attributes.
func (c *Client) CUPSGetDevices(ctx context.Context,
	sel *GetDevicesSelection, attrs []string) (
	[]*ipp.DeviceAttributes, error) {

	if sel == nil {
		sel = DefaultGetDevicesSelection
	}

	tm := 0
	if sel.Timeout != 0 {
		tm = int((sel.Timeout + time.Second - 1) / time.Second)
	}

	rq := &ipp.CUPSGetDevicesRequest{
		RequestHeader:       ipp.DefaultRequestHeader,
		ExcludeSchemes:      sel.ExcludeSchemes,
		IncludeSchemes:      sel.IncludeSchemes,
		Limit:               optional.NotZero(sel.Limit),
		Timeout:             optional.NotZero(tm),
		RequestedAttributes: attrs,
	}

	rsp := &ipp.CUPSGetDevicesResponse{}

	err := c.IPPClient.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Printer, nil
}

// CUPSGetPPD requests PPD file by printer URI or the PPD file name.
//
// It returns one of the following:
//   - non-nil body where requested PPD file can be read from
//   - nil body and non-empty seeOtherURI string, that specify
//     the printer URI that can serve the request
//   - nil body, empty seeOtherURI and non-nil err in a case of error.
//
// If non-nil body returned, caller MUST close it after use.
func (c *Client) CUPSGetPPD(ctx context.Context,
	printerURI, ppdName string) (
	body io.ReadCloser, seeOtherURI string, err error) {

	rq := &ipp.CUPSGetPPDRequest{
		RequestHeader: ipp.DefaultRequestHeader,
		PrinterURI:    optional.NotZero(printerURI),
		PPDName:       optional.NotZero(ppdName),
	}

	rsp := &ipp.CUPSGetPPDResponse{}

	err = c.IPPClient.DoWithBody(ctx, rq, rsp)
	if err != nil {
		return
	}

	if rsp.Status == goipp.StatusOk {
		return rsp.Body, "", nil
	}

	rsp.Body.Close()
	if rsp.Status == goipp.StatusCupsSeeOther {
		return nil, optional.Get(rsp.PrinterURI), nil
	}

	return nil, "", fmt.Errorf("IPP: %s", rsp.Status)
}

// CUPSGetPPDs requests information about PPD files available at the server.
//
// If GetPPDsSelection is nil, all PPDs will be returned (the response could be
// really large at this case).
//
// The attrs attribute allows to specify list of requested attributes
// ("all" if set to nil).
func (c *Client) CUPSGetPPDs(ctx context.Context,
	sel *GetPPDsSelection, attrs []string) ([]*ipp.PPDAttributes, error) {

	if sel == nil {
		sel = &GetPPDsSelection{}
	}

	if attrs == nil {
		attrs = []string{"all"}
	}

	rq := &ipp.CUPSGetPPDsRequest{
		RequestHeader:       ipp.DefaultRequestHeader,
		ExcludeSchemes:      sel.ExcludeSchemes,
		IncludeSchemes:      sel.IncludeSchemes,
		Limit:               optional.NotZero(sel.Limit),
		PPDMake:             optional.NotZero(sel.PPDMake),
		PPDMakeAndModel:     optional.NotZero(sel.PPDMakeAndModel),
		PPDModelNumber:      optional.NotZero(sel.PPDModelNumber),
		PPDNaturalLanguage:  optional.NotZero(sel.PPDNaturalLanguage),
		PPDProduct:          optional.NotZero(sel.PPDProduct),
		PPDPsVersion:        optional.NotZero(sel.PPDPsVersion),
		PPDType:             optional.NotZero(sel.PPDType),
		RequestedAttributes: attrs,
	}

	rsp := &ipp.CUPSGetPPDsResponse{}

	err := c.IPPClient.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	if rsp.Status == goipp.StatusOk {
		return rsp.PPDs, err
	}

	return nil, fmt.Errorf("IPP: %s", rsp.Status)
}
