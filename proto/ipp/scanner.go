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
	"io"
	"net/http"
	"sync"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// docResult holds the outcome of a doc.Next() call delivered via channel.
type docResult struct {
	file abstract.DocumentFile
	err  error
}

type docBodyCloser struct {
	file    abstract.DocumentFile
	onClose func()
}

func (d *docBodyCloser) Read(p []byte) (int, error) { return d.file.Read(p) }
func (d *docBodyCloser) Close() error {
	d.onClose()
	return nil
}

// Scanner implements the IPP Scan Service as defined in PWG5100.17.
type Scanner struct {
	options ScannerOptions
	server  *Server
	attrs   *PrinterAttributes
	q       *queue

	activeDoc        abstract.Document
	activeJob        int
	activeDocPageNum int
	activeDocCh      chan docResult
	lock             sync.Mutex
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

	// DocumentDataGetInterval is the number of seconds the scanner tells
	// the client to wait before retrying Get-Next-Document-Data when no
	// document data is immediately available.
	DocumentDataGetInterval int
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
		options:     options,
		server:      server,
		attrs:       attrs,
		q:           newQueue(),
		activeDocCh: make(chan docResult, 1),
	}

	// Install scan-service handlers.
	server.RegisterHandler(NewHandler(scanner.handleGetPrinterAttributes))
	server.RegisterHandler(NewHandler(scanner.handleCreateScanJob))
	server.RegisterHandler(NewHandler(scanner.handleGetNextDocumentData))
	server.RegisterHandler(NewHandler(scanner.handleCancelJob))

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
	rq *GetPrinterAttributesRequest) (*goipp.Message, io.ReadCloser, error) {

	return rq.Apply(scanner.attrs, scanner.options.UseRawPrinterAttributes), nil, nil
}

// handleCreateScanJob handles Create-Job request on the Scan Service
func (scanner *Scanner) handleCreateScanJob(
	ctx context.Context,
	rq *CreateJobRequest) (*goipp.Message, io.ReadCloser, error) {

	// input-attributes MUST be supplied by the client.
	if rq.InputAttributes == nil {
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missing required input-attributes")
	}

	// Validate scan parameters against scanner capabilities.
	caps := scanner.options.Scanner.Capabilities()
	req := rq.JobCreateOperation.ToAbstract()
	filled, err := caps.FillRequest(&req)
	if err != nil {
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorAttributesOrValues,
			"invalid scan parameters: %s", err)
	}

	// Single-document model: reject if another scan is already active.
	scanner.lock.Lock()
	if scanner.activeDoc != nil {
		scanner.lock.Unlock()
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBusy,
			"scanner is busy with another job")
	}

	doc, err := scanner.options.Scanner.Scan(ctx, *filled)
	if err != nil {
		scanner.lock.Unlock()
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorDevice,
			"scan failed: %s", err)
	}

	j := newJob(&rq.JobCreateOperation, rq.Job)
	scanner.q.Push(j)

	scanner.activeDoc = doc
	scanner.activeJob = j.JobID
	scanner.activeDocPageNum = 0
	go func() {
		f, err := doc.Next()
		scanner.activeDocCh <- docResult{file: f, err: err}
	}()
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

	return rsp.Encode(), nil, nil
}

// handleGetNextDocumentData handles Get-Next-Document-Data requests
func (scanner *Scanner) handleGetNextDocumentData(
	ctx context.Context,
	rq *GetNextDocumentDataRequest) (*goipp.Message, io.ReadCloser, error) {

	scanner.lock.Lock()
	defer scanner.lock.Unlock()

	jobID := optional.Get(rq.JobID)
	if jobID == 0 {
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missing job-id")
	}

	if jobID != scanner.activeJob {
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotFound,
			"no active scan job")
	}

	doc := scanner.activeDoc

	var result docResult
	select {
	case result = <-scanner.activeDocCh:
	default:
		rsp := GetNextDocumentDataResponse{
			ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
			DocumentDataGetInterval: optional.New(
				scanner.documentDataGetInterval()),
		}
		return rsp.Encode(), nil, nil
	}

	if result.err == io.EOF {
		scanner.cleanupActiveScan()
		rsp := GetNextDocumentDataResponse{
			ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
			LastDocument:   true,
		}
		return rsp.Encode(), nil, nil
	}

	if result.err != nil {
		scanner.cleanupActiveScan()
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorDevice,
			"scan error: %s", result.err)
	}

	scanner.activeDocPageNum++
	docPageNum := scanner.activeDocPageNum

	body := &docBodyCloser{
		file: result.file,
		onClose: func() {
			go func() {
				f, err := doc.Next()
				scanner.activeDocCh <- docResult{file: f, err: err}
			}()
		},
	}

	rsp := GetNextDocumentDataResponse{
		ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
		LastDocument:   false,
		DocumentFormat: optional.New(result.file.Format()),
		Document: &DocumentStatus{
			DocumentNumber: optional.New(docPageNum)},
	}

	return rsp.Encode(), body, nil
}

func (scanner *Scanner) documentDataGetInterval() int {
	if scanner.options.DocumentDataGetInterval > 0 {
		return scanner.options.DocumentDataGetInterval
	}
	return 5
}

// handleCancelJob handles Cancel-Job request on the Scan Service.
func (scanner *Scanner) handleCancelJob(
	ctx context.Context,
	rq *CancelJobRequest) (*goipp.Message, io.ReadCloser, error) {

	var j *job

	switch {
	case rq.PrinterURI != nil && rq.JobID != nil:
		j = scanner.q.JobByID(*rq.JobID)
		if j == nil {
			return nil, nil, NewErrIPPFromRequest(rq,
				goipp.StatusErrorNotFound,
				"job not found (job-id=%d)", *rq.JobID)
		}

	case rq.JobURI != nil:
		j = scanner.q.JobByURI(*rq.JobURI)
		if j == nil {
			return nil, nil, NewErrIPPFromRequest(rq,
				goipp.StatusErrorNotFound,
				"job not found (job-uri=%q)", *rq.JobURI)
		}

	default:
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorBadRequest,
			"missing job-id and job-uri attributes")
	}

	j.Lock()
	if !j.isCancelable() {
		j.Unlock()
		return nil, nil, NewErrIPPFromRequest(rq,
			goipp.StatusErrorNotPossible,
			"job cannot be canceled in %v state", j.JobState)
	}

	j.beginCancel()
	jobID := j.JobID
	j.Unlock()

	scanner.lock.Lock()
	if scanner.activeJob == jobID {
		scanner.cleanupActiveScan()
	}
	scanner.lock.Unlock()

	j.Lock()
	j.finishCancel()
	j.Unlock()

	rsp := CancelJobResponse{
		ResponseHeader: rq.ResponseHeader(goipp.StatusOk),
	}

	return rsp.Encode(), nil, nil
}

func (scanner *Scanner) cleanupActiveScan() {
	if scanner.activeDoc != nil {
		scanner.activeDoc.Close()
		scanner.activeDoc = nil
	}
	scanner.activeJob = 0
	scanner.activeDocPageNum = 0
}
