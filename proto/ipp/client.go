// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP client

package ipp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/user"
	"strings"
	"sync/atomic"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// Client implements Client-side IPP Printer object.
type Client struct {
	URL        *url.URL          // Destination URL (ipp://...)
	HTTPClient *transport.Client // HTTP Client
	RequestID  uint32            // RequestID of the next request
	decoderOpt *DecoderOptions   // Options for message decoder
}

// NewClient creates a new IPP client.
//
// If tr is nil, [transport.NewTransport] will be used to create
// a new transport.
func NewClient(u *url.URL, tr *transport.Transport) *Client {
	c := &Client{
		URL:        u,
		HTTPClient: transport.NewClient(tr),
	}

	return c
}

// SetDecoderOptions updates the [DecoderrOptions] that affect decoding
// of the received IPP messages
func (c *Client) SetDecoderOptions(opt *DecoderOptions) {
	c.decoderOpt = opt
}

// requestid generates a next RequestID
func (c *Client) requestid() uint32 {
	// IPP doesn't allow RequestID to be zero, so roll
	// until first non-zero value
	var id uint32
	for id == 0 {
		id = atomic.AddUint32(&c.RequestID, 1)
	}

	return id
}

// Do sends the [Request] and waits for [Response].
//
// The following Request fields are filled automatically:
//   - Version, if zero, will be set to goipp.DefaultVersion
//   - RequestID will be set to next Client's RequestID in sequence
//
// It automatically closes Response Body. This is convenient
// for most IPP requests, as body is rarely returned by IPP.
//
// For requests with returned body, use [Client.DoWithBody] instead.
func (c *Client) Do(ctx context.Context, rq Request, rsp Response) error {
	err := c.DoWithBody(ctx, rq, rsp)
	if err == nil {
		if body := rsp.Header().Body; body != nil {
			body.Close()
			rsp.Header().Body = nil
		}
	}
	return err
}

// DoWithBody sends the Request and waits for Response.
//
// The following Request fields are filled automatically:
//   - Version, if zero, will be set to goipp.DefaultVersion
//   - RequestID will be set to next Client's RequestID in sequence
//
// On success, caller MUST close Response body after use.
func (c *Client) DoWithBody(ctx context.Context,
	rq Request, rsp Response) error {

	// Encode IPP message
	buf := &bytes.Buffer{}
	msg := rq.Encode()

	if msg.Version == 0 {
		msg.Version = goipp.DefaultVersion
	}

	if msg.RequestID == 0 {
		msg.RequestID = c.requestid()
	}

	msg.Encode(buf)

	// Log the IPP request
	f := goipp.NewFormatter()
	f.SetIndent(2)
	f.FmtRequest(msg)
	log.Debug(ctx, "IPP request:\n%s", f.Bytes())

	// Attach Request body, if any
	body := rq.Header().Body
	if body == nil {
		body = buf
	} else {
		body = io.MultiReader(buf, body)
	}

	// Create HTTP request
	httpRq, err := transport.NewRequest(ctx, "POST", c.URL, body)
	if err != nil {
		return err
	}

	httpRq.Header.Set("Content-Type", "application/ipp")

	// If we are on local socket, set "PeerCred username" as
	// authentication information...
	if strings.ToLower(httpRq.URL.Scheme) == "unix" {
		usr, err := user.Current()
		if err != nil {
			return err
		}

		auth := fmt.Sprintf("PeerCred %s", usr.Username)
		httpRq.Header.Set("Authorization", auth)
	}

	// Call server
	httpRsp, err := c.HTTPClient.Do(httpRq)
	if err != nil {
		return err
	}

	if httpRsp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP: %s", httpRsp.Status)
		goto ERROR
	}

	// Decode IPP message
	msg.Reset()
	err = msg.Decode(httpRsp.Body)
	if err != nil {
		goto ERROR
	}

	// Log the IPP response
	f.Reset()
	f.SetIndent(2)
	f.FmtResponse(msg)
	log.Debug(ctx, "IPP response:\n%s", f.Bytes())

	// Decode Response
	err = rsp.Decode(msg, c.decoderOpt)
	if err != nil {
		goto ERROR
	}

	// Save IPPMessage, remainder of body and return
	rsp.Header().IPPMessage = msg
	rsp.Header().Body = httpRsp.Body

	return nil

ERROR:
	log.Debug(ctx, "HTTP %s %s - %s", httpRq.Method, httpRq.URL, err)
	httpRsp.Body.Close()
	return err
}

// GetPrinterAttributes returns printer attributes.
// The attrs attribute allows to specify list of requested attributes.
//
// Note, certain printer attributes may depend on the format being
// printer, so second argument, if not "", allows to specify the
// desired document format.
//
// According to the RFC8011, only the following attributes may depend
// on the document format:
//
//   - Job Template attributes ("xxx-default", "xxx-supported", and "xxx-ready")
//   - "pdl-override-supported"
//   - "compression-supported"
//   - "job-k-octets-supported"
//   - "job-impressions-supported
//   - "job-media-sheets-supported"
//   - "printer-driver-installer"
//   - "color-supported"
//   - "reference-uri-schemes-supported"
//
// See RFC8011, 4.2.5.1. for details.
func (c *Client) GetPrinterAttributes(ctx context.Context,
	attrs []string, format string) (*PrinterAttributes, error) {

	rq := &GetPrinterAttributesRequest{
		RequestHeader:       DefaultRequestHeader,
		PrinterURI:          c.URL.String(),
		DocumentFormat:      optional.NotZero(format),
		RequestedAttributes: attrs,
	}

	rsp := &GetPrinterAttributesResponse{}

	err := c.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Printer, nil
}

// CreateJob sends a Create-Job request.
//
// For scan jobs (PWG5100.17), op.InputAttributes must be set.
// For print jobs, document data must be supplied separately using
// Send-Document.
func (c *Client) CreateJob(
	ctx context.Context,
	op JobCreateOperation, job *JobTemplate) (
	*CreateJobResponse, error) {

	if job == nil {
		job = &JobTemplate{}
	}

	rq := &CreateJobRequest{
		RequestHeader:      DefaultRequestHeader,
		JobCreateOperation: op,
		JobTemplate:        job,
	}

	rsp := &CreateJobResponse{}

	err := c.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

// GetNextDocumentData sends a Get-Next-Document-Data request.
func (c *Client) GetNextDocumentData(
	ctx context.Context,
	rq *GetNextDocumentDataRequest) (
	*GetNextDocumentDataResponse, error) {

	rsp := &GetNextDocumentDataResponse{}

	err := c.DoWithBody(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

// CancelJob cancels the job identified by jobID on the IPP object at c.URL.
func (c *Client) CancelJob(ctx context.Context, jobID int,
	message string) error {

	rq := &CancelJobRequest{
		RequestHeader: DefaultRequestHeader,
		JobCancelOperation: JobCancelOperation{
			PrinterURI: optional.New(c.URL.String()),
			JobID:      optional.New(jobID),
			Message:    optional.NotZero(message),
		},
	}

	rsp := &CancelJobResponse{}

	return c.Do(ctx, rq, rsp)
}

// GetJobs returns jobs from the IPP object at c.URL.
func (c *Client) GetJobs(ctx context.Context,
	which KwWhichJobs, myJobs bool, limit int,
	attrs []KwRequestedAttribute) ([]JobGroupEntry, error) {

	rq := &GetJobsRequest{
		RequestHeader:       DefaultRequestHeader,
		PrinterURI:          c.URL.String(),
		RequestedAttributes: attrs,
		WhichJobs:           optional.NotZero(which),
		MyJobs:              optional.New(myJobs),
		Limit:               optional.NotZero(limit),
	}

	rsp := &GetJobsResponse{}

	err := c.Do(ctx, rq, rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Jobs, nil
}

// GetJobAttributes returns attributes for jobID on the IPP object at c.URL.
func (c *Client) GetJobAttributes(ctx context.Context, jobID int,
	attrs []KwRequestedAttribute) (JobGroupEntry, error) {

	rq := &GetJobAttributesRequest{
		RequestHeader:       DefaultRequestHeader,
		PrinterURI:          optional.New(c.URL.String()),
		JobID:               optional.New(jobID),
		RequestedAttributes: attrs,
	}

	rsp := &GetJobAttributesResponse{}

	err := c.Do(ctx, rq, rsp)
	if err != nil {
		return JobGroupEntry{}, err
	}

	return rsp.Job, nil
}
