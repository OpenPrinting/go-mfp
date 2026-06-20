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

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/goipp"
)

// Printer implements the IPP printer.
type Printer struct {
	options PrinterOptions       // Printer options
	server  *Server              // Underlying IPP server
	attrs   *PrinterAttributes   // Printer attributes
	q       *queue               // Job queue
	backend abstract.Printer // Print backend
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

// SetPrintBackend installs backend as the handler for incoming
// print documents. Pass nil to clear a previously set backend.
func (printer *Printer) SetPrintBackend(backend abstract.Printer) {
	printer.backend = backend
}

// ServeHTTP handles incoming HTTP request. It implements
// [http.Handler] interface.
func (printer *Printer) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	printer.server.ServeHTTP(w, rq)
}

// handleGetPrinterAttributes handles Get-Printer-Attributes request.
func (printer *Printer) handleGetPrinterAttributes(
	ctx context.Context,
	rq *GetPrinterAttributesRequest) (*goipp.Message, io.ReadCloser, error) {

	return rq.Apply(printer.attrs, printer.options.UseRawPrinterAttributes), nil, nil
}

// handleValidateJob handles Validate-Job request.
func (printer *Printer) handleValidateJob(
	ctx context.Context,
	rq *ValidateJobRequest) (*goipp.Message, io.ReadCloser, error) {

	rsp := ValidateJobResponse{
		ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
	}

	return rsp.Encode(), nil, nil
}

// handleCreateJob handles Create-Job request.
func (printer *Printer) handleCreateJob(
	ctx context.Context,
	rq *CreateJobRequest) (*goipp.Message, io.ReadCloser, error) {

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

	return rsp.Encode(), nil, nil
}

// handleCreateJob handles Send-Document request.
func (printer *Printer) handleSendDocument(
	ctx context.Context,
	rq *SendDocumentRequest) (*goipp.Message, io.ReadCloser, error) {

	// Lookup the job
	var j *job

	switch {
	case rq.PrinterURI != nil && rq.JobID != nil:
		j = printer.q.JobByID(*rq.JobID)
		if j == nil {
			err := NewErrIPPFromRequest(rq,
				goipp.StatusErrorNotFound,
				"job not found (job-id=%d)", *rq.JobID)
			return nil, nil, err
		}

	case rq.JobURI != nil:
		j = printer.q.JobByURI(*rq.JobURI)
		if j == nil {
			err := NewErrIPPFromRequest(rq,
				goipp.StatusErrorNotFound,
				"job not found (job-uri=%q)", *rq.JobURI)
			return nil, nil, err
		}

	default:
		err := NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missed job-id and job-uri attributes")
		return nil, nil, err
	}

	j.Lock()
	defer j.Unlock()

	// Check job state
	if j.JobState != EnJobStatePendingHeld {
		err := NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotPossible,
			"job is not in pending-held state")
		return nil, nil, err
	}

	if j.SendDocumentActive {
		err := NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotPossible,
			"Send-Document already in progress")
		return nil, nil, err
	}

	// Consume the document body
	j.SendDocumentActive = true
	j.Unlock()

	if printer.backend != nil {
		// Build protocol-independent job parameters
		params := abstract.PrinterRequest{}

		if rq.DocumentFormat != nil {
			params.Format = *rq.DocumentFormat
		}
		if rq.DocumentName != nil {
			params.JobName = *rq.DocumentName
		} else if j.JobStatus.JobName != nil {
			params.JobName = *j.JobStatus.JobName
		}
		if rq.Job != nil {
			if rq.Job.Copies != nil {
				params.Copies = *rq.Job.Copies
			}
			if rq.Job.Sides != nil {
				params.Sides = sidesToAbstract(*rq.Job.Sides)
			}
			if rq.Job.PrintColorMode != nil {
				params.ColorMode = colorModeToAbstract(*rq.Job.PrintColorMode)
			}
			if rq.Job.Media != nil {
				params.Media = mediaSizeToAbstract(*rq.Job.Media)
			}
		}

		if err := printer.backend.PrintDocument(params, rq.Body); err != nil {
			log.Error(ctx, "Send-Document: backend error: %s", err)
		}
	} else {
		// No backend — drain the body so the connection stays clean
		n, err := io.Copy(io.Discard, rq.Body)
		if err != nil {
			log.Error(ctx, "Send-Document: %s", err)
		} else {
			log.Debug(ctx, "Send-Document: %d bytes discarded (no backend)", n)
		}
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

	return rsp.Encode(), nil, nil
}
