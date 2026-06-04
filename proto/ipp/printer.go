// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP printer implementation.

package ipp

import (
	"context"
	"io"
	"net/http"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/goipp"
)

// Printer implements the IPP printer.
type Printer struct {
	options  PrinterOptions     // Printer options
	server   *Server            // Underlying IPP server
	attrs    *PrinterAttributes // Printer attributes
	q        *queue             // Job queue
	receiver DocumentReceiver   // Document capture callback
}

// PrinterOptions extends [ServerOptions] with printer-specific
// parameters.
type PrinterOptions struct {
	ServerOptions

	// UseRawPrinterAttributes, if set, instruct [Printer]
	// to return attributes, based on PrinterAttributes.RawAttrs
	// instead of the the PrinterAttributes.Encode.
	//
	// It can be useful when the exact content and ordering of
	// printer attributes needs to be specified, because conversion
	// from the IPP attributes to and from the Go structure
	// is not lossless.
	UseRawPrinterAttributes bool
}

// NewPrinter creates a new [Printer], which facilities and
// behavior is defined by the supplied [PrinterAttributes].
func NewPrinter(attrs *PrinterAttributes, options PrinterOptions) *Printer {
	// Create the Printer structure
	server := NewServer(options.ServerOptions)
	printer := &Printer{
		options: options,
		server:  server,
		attrs:   attrs,
		q:       newQueue(),
	}

	// Install request handlers
	server.RegisterHandler(NewHandler(printer.handleGetPrinterAttributes))
	server.RegisterHandler(NewHandler(printer.handleValidateJob))
	server.RegisterHandler(NewHandler(printer.handleCreateJob))
	server.RegisterHandler(NewHandler(printer.handleSendDocument))

	return printer
}

// ServeHTTP handles incoming HTTP request. It implements
// [http.Handler] interface.
func (printer *Printer) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	printer.server.ServeHTTP(w, rq)
}

// handleGetPrinterAttributes handles Get-Printer-Attributes request.
func (printer *Printer) handleGetPrinterAttributes(
	ctx context.Context,
	rq *GetPrinterAttributesRequest) (*goipp.Message, error) {

	return rq.Apply(printer.attrs, printer.options.UseRawPrinterAttributes), nil
}

// handleValidateJob handles Validate-Job request.
func (printer *Printer) handleValidateJob(
	ctx context.Context,
	rq *ValidateJobRequest) (*goipp.Message, error) {

	rsp := ValidateJobResponse{
		ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
	}

	return rsp.Encode(), nil
}

// handleCreateJob handles Create-Job request.
func (printer *Printer) handleCreateJob(
	ctx context.Context,
	rq *CreateJobRequest) (*goipp.Message, error) {

	// Create new job
	j := newJob(&rq.JobCreateOperation, rq.Job)
	j.Lock()
	defer j.Unlock()

	printer.q.Push(j)

	// Prepare the CreateJobResponse
	rsp := CreateJobResponse{
		ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
		Job: &JobStatus{
			JobID:           j.JobID,
			JobState:        j.JobState,
			JobStateReasons: j.JobStateReasons,
			JobURI:          j.JobURI,
		},
	}

	return rsp.Encode(), nil
}

// handleCreateJob handles Send-Document request.
func (printer *Printer) handleSendDocument(
	ctx context.Context,
	rq *SendDocumentRequest) (*goipp.Message, error) {

	// Lookup the job
	var j *job

	switch {
	case rq.PrinterURI != nil && rq.JobID != nil:
		j = printer.q.JobByID(*rq.JobID)
		if j == nil {
			err := NewErrIPPFromRequest(rq,
				goipp.StatusErrorNotFound,
				"job not found (job-id=%d)", *rq.JobID)
			return nil, err
		}

	case rq.JobURI != nil:
		j = printer.q.JobByURI(*rq.JobURI)
		if j == nil {
			err := NewErrIPPFromRequest(rq,
				goipp.StatusErrorNotFound,
				"job not found (job-uri=%q)", *rq.JobURI)
			return nil, err
		}

	default:
		err := NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missed job-id and job-uri attributes")
		return nil, err
	}

	j.Lock()
	defer j.Unlock()

	// Check job state
	if j.JobState != EnJobStatePendingHeld {
		err := NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotPossible,
			"job is not in pending-held state")
		return nil, err
	}

	if j.SendDocumentActive {
		err := NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotPossible,
			"Send-Document already in progress")
		return nil, err
	}

	// Consume the document body
	j.SendDocumentActive = true
	j.Unlock()

	data, err := io.ReadAll(rq.Body)
	if err == nil {
		log.Debug(ctx, "Send-Document: %d bytes received", len(data))
	} else {
		log.Error(ctx, "Send-Document: %s", err)
	}

	// Invoke document receiver callback if installed
	if err == nil && printer.receiver != nil {
		format := ""
		if rq.DocumentFormat != nil {
			format = *rq.DocumentFormat
		}
		printer.receiver(j.JobID, format, data)
	}

	j.Lock()
	j.SendDocumentActive = false

	// Generate response
	rsp := &SendDocumentResponse{
		Job: &JobStatus{
			JobID:           j.JobID,
			JobState:        j.JobState,
			JobStateReasons: j.JobStateReasons,
			JobURI:          j.JobURI,
		},
	}

	return rsp.Encode(), nil
}
