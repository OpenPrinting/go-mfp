// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// IPP Scan Service implementation (PWG5100.17).

package ipp

import (
	"context"
	"net/http"
	"sync"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/goipp"
)

// Scanner implements the IPP Scan Service as defined in PWG5100.17.
type Scanner struct {
	options ScannerOptions
	server  *Server
	attrs   *PrinterAttributes
	q       *queue

	activeDoc abstract.Document
	activeJob int
	lock      sync.Mutex
}

// ScannerOptions extends [ServerOptions] with scanner-specific parameters.
type ScannerOptions struct {
	ServerOptions

	// Scanner is the underlying abstract scanner. Required.
	Scanner abstract.Scanner

	// UseRawPrinterAttributes, if set, instructs [Scanner] to return
	// attributes based on PrinterAttributes.RawAttrs instead of the
	// PrinterAttributes.Encode result. See [PrinterOptions] for details.
	UseRawPrinterAttributes bool
}

// NewScanner creates a new [Scanner], whose facilities and behavior
// are defined by the supplied [PrinterAttributes] and the underlying
// [abstract.Scanner].
func NewScanner(attrs *PrinterAttributes, options ScannerOptions) *Scanner {
	// Populate ScannerDescription from the underlying abstract scanner.
	attrs.ScannerDescription =
		fromAbstractScannerDescription(options.Scanner.Capabilities())

	server := NewServer(options.ServerOptions)
	scanner := &Scanner{
		options: options,
		server:  server,
		attrs:   attrs,
		q:       newQueue(),
	}

	// Install scan-service handlers.
	server.RegisterHandler(NewHandler(scanner.handleGetPrinterAttributes))
	server.RegisterHandler(NewHandler(scanner.handleCreateScanJob))

	return scanner
}

// ServeHTTP handles incoming HTTP requests. It implements
// [http.Handler] interface.
func (scanner *Scanner) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	scanner.server.ServeHTTP(w, rq)
}

// handleGetPrinterAttributes handles Get-Printer-Attributes request.
func (scanner *Scanner) handleGetPrinterAttributes(
	ctx context.Context,
	rq *GetPrinterAttributesRequest) (*goipp.Message, error) {

	return rq.Apply(scanner.attrs, scanner.options.UseRawPrinterAttributes), nil
}

// handleCreateScanJob handles Create-Job request on the Scan Service
func (scanner *Scanner) handleCreateScanJob(
	ctx context.Context,
	rq *CreateJobRequest) (*goipp.Message, error) {

	// input-attributes MUST be supplied by the client.
	if rq.InputAttributes == nil {
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missing required input-attributes")
	}

	// Validate scan parameters against scanner capabilities.
	caps := scanner.options.Scanner.Capabilities()
	req := rq.JobCreateOperation.ToAbstract()
	filled, err := caps.FillRequest(&req)
	if err != nil {
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorAttributesOrValues,
			"invalid scan parameters: %s", err)
	}

	// Single-document model: reject if another scan is already active.
	scanner.lock.Lock()
	if scanner.activeDoc != nil {
		scanner.lock.Unlock()
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBusy,
			"scanner is busy with another job")
	}

	doc, err := scanner.options.Scanner.Scan(ctx, *filled)
	if err != nil {
		scanner.lock.Unlock()
		return nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorDevice,
			"scan failed: %s", err)
	}

	j := newJob(&rq.JobCreateOperation, rq.Job)
	scanner.q.Push(j)

	scanner.activeDoc = doc
	scanner.activeJob = j.JobID
	scanner.lock.Unlock()

	j.Lock()
	j.JobState = EnJobStateProcessing
	j.JobStateReasons = []KwJobStateReasons{KwJobStateReasonsNone}
	j.Unlock()

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
