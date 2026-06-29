// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// WS-Scan server on a top of abstract.Scanner

package wsscan

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/log/trace"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/uuid"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

var (
	errBusy       = errors.New("scanner busy")
	errInvalidJob = errors.New("invalid job")
)

// AbstractServer implements a WS-Scan server on top of
// [abstract.Scanner].
type AbstractServer struct {
	options   AbstractServerOptions
	caps      *abstract.ScannerCapabilities
	status    ScannerStatus
	document  abstract.Document
	jobs      jobList
	nextJobID int
	lock      sync.Mutex
}

// AbstractServerOptions allows specifying options that can
// modify the [AbstractServer] behavior.
type AbstractServerOptions struct {
	Scanner abstract.Scanner

	// BasePath is required so the server knows how to
	// interpret incoming request paths.
	BasePath string
}

// NewAbstractServer returns a new [AbstractServer].
func NewAbstractServer(options AbstractServerOptions) *AbstractServer {
	srv := &AbstractServer{
		options:   options,
		caps:      options.Scanner.Capabilities(),
		nextJobID: 1,
	}

	srv.status = ScannerStatus{
		ScannerState: Idle,
	}

	return srv
}

// ServeHTTP serves incoming HTTP requests.
// It implements the [http.Handler] interface.
func (srv *AbstractServer) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	query := transport.NewServerQuery(w, rq)
	defer query.Finish()

	// Check path prefix
	if !strings.HasPrefix(query.RequestURL().Path, srv.options.BasePath) {
		query.Reject(http.StatusNotFound, nil)
		return
	}

	// WS-Scan uses POST only
	if query.RequestMethod() != "POST" {
		query.Reject(http.StatusMethodNotAllowed, nil)
		return
	}

	// Read and decode SOAP message
	data, err := io.ReadAll(query.RequestBody())
	if err != nil {
		query.Reject(http.StatusBadRequest, err)
		return
	}

	root, err := xmldoc.Decode(NsMap, bytes.NewReader(data))
	if err != nil {
		query.Reject(http.StatusBadRequest, err)
		return
	}

	msg, err := DecodeMessage(root)
	if err != nil {
		query.Reject(http.StatusBadRequest, err)
		return
	}

	// Notify tracer
	trace.OnRequest(query, traceMessage{msg}, nil)

	// Dispatch by body type
	var rsp Body
	switch body := msg.Body.(type) {
	case *CancelJobRequest:
		rsp, err = srv.handleCancelJobRequest(query, body)

	case *CreateScanJobRequest:
		rsp, err = srv.handleCreateScanJobRequest(query, body)

	case *GetActiveJobsRequest:
		rsp, err = srv.handleGetActiveJobsRequest(query, body)

	case *GetJobElementsRequest:
		rsp, err = srv.handleGetJobElementsRequest(query, body)

	case *GetJobHistoryRequest:
		rsp, err = srv.handleGetJobHistoryRequest(query, body)

	case *GetScannerElementsRequest:
		rsp, err = srv.handleGetScannerElementsRequest(query, body)

	case *RetrieveImageRequest:
		rsp, err = srv.handleRetrieveImageRequest(query, body)

	default:
		query.Reject(http.StatusBadRequest, nil)
		return
	}

	if err == nil {
		srv.sendSOAPResponse(query, msg, rsp)
	}
}

// handleGetScannerElementsRequest handles GetScannerElements requests.
func (srv *AbstractServer) handleGetScannerElementsRequest(
	query *transport.ServerQuery,
	req *GetScannerElementsRequest,
) (Body, error) {

	// Build ScannerElemData for each requested element, skipping duplicates.
	var elements []ScannerElemData
	seen := generic.NewSet[ScannerElemName]()

	for _, re := range req.RequestedElements {
		if !seen.TestAndAdd(re) {
			continue
		}
		switch re {
		case ScannerElemDefaultScanTicket:
			req := srv.caps.DefaultRequest()
			if req != nil {
				ticket := fromAbstractScannerRequest(req)
				elements = append(elements, ScannerElemData{
					Name:              ScannerElemDefaultScanTicket,
					Valid:             BooleanElement("true"),
					DefaultScanTicket: optional.New(ticket),
				})
			}

		case ScannerElemDescription:
			desc := fromAbstractScannerDescription(srv.caps)
			elements = append(elements, ScannerElemData{
				Name:               ScannerElemDescription,
				Valid:              BooleanElement("true"),
				ScannerDescription: optional.New(desc),
			})

		case ScannerElemConfiguration:
			conf := fromAbstractScannerConfiguration(srv.caps)
			elements = append(elements, ScannerElemData{
				Name:                 ScannerElemConfiguration,
				Valid:                BooleanElement("true"),
				ScannerConfiguration: optional.New(conf),
			})

		case ScannerElemStatus:
			srv.lock.Lock()
			status := srv.status
			srv.lock.Unlock()
			status.ScannerCurrentTime = time.Now()
			elements = append(elements, ScannerElemData{
				Name:          ScannerElemStatus,
				Valid:         BooleanElement("true"),
				ScannerStatus: optional.New(status),
			})
		}
	}

	return &GetScannerElementsResponse{
		ScannerElements: elements,
	}, nil
}

// handleGetJobElementsRequest handles GetJobElements requests.
// It returns the requested job schema elements for the job identified
// by JobID.
func (srv *AbstractServer) handleGetJobElementsRequest(
	query *transport.ServerQuery,
	req *GetJobElementsRequest,
) (Body, error) {

	srv.lock.Lock()
	job := srv.jobs.get(req.JobID)
	if job == nil {
		srv.lock.Unlock()
		query.Reject(http.StatusNotFound, nil)
		return nil, errInvalidJob
	}
	j := *job
	srv.lock.Unlock()

	var elements []JobElemData
	seen := generic.NewSet[JobElemName]()

	for _, re := range req.RequestedElements {
		if !seen.TestAndAdd(re) {
			continue
		}
		switch re {
		case JobElemStatus:
			elements = append(elements, JobElemData{
				Name:      JobElemStatus,
				Valid:     BooleanElement("true"),
				JobStatus: optional.New(jobStatusFrom(j)),
			})

		case JobElemScanTicket:
			elements = append(elements, JobElemData{
				Name:       JobElemScanTicket,
				Valid:      BooleanElement("true"),
				ScanTicket: optional.New(j.scanTicket),
			})

		case JobElemDocuments:
			docs := Documents{}
			if j.scanTicket.DocumentParameters != nil {
				docs.DocumentFinalParameters =
					optional.Get(j.scanTicket.DocumentParameters)
			}
			elements = append(elements, JobElemData{
				Name:      JobElemDocuments,
				Valid:     BooleanElement("true"),
				Documents: optional.New(docs),
			})
		}
	}

	return &GetJobElementsResponse{
		JobElements: elements,
	}, nil
}

// jobStatusFrom builds a [JobStatus] from a [jobInfo].
func jobStatusFrom(j jobInfo) JobStatus {
	js := JobStatus{
		JobID:          j.jobID,
		JobState:       j.state,
		ScansCompleted: j.scansCompleted,
	}
	if !j.createdTime.IsZero() {
		js.JobCreatedTime = optional.New(j.createdTime)
	}
	if !j.completedTime.IsZero() {
		js.JobCompletedTime = optional.New(j.completedTime)
	}
	return js
}

// handleCreateScanJobRequest handles CreateScanJob requests.
func (srv *AbstractServer) handleCreateScanJobRequest(
	query *transport.ServerQuery,
	req *CreateScanJobRequest,
) (Body, error) {

	srv.lock.Lock()
	defer srv.lock.Unlock()

	// Check if previous scan is still in progress
	if srv.document != nil {
		query.Reject(http.StatusServiceUnavailable, nil)
		return nil, errBusy
	}

	// Convert ScanTicket to abstract.ScannerRequest
	absreq := req.ScanTicket.ToAbstract()

	// Fill missing parameters with scanner defaults and validate against
	// capabilities. FillRequest returns an error for unsupported params.
	filled, err := srv.caps.FillRequest(&absreq)
	if err != nil {
		query.Reject(http.StatusConflict, err)
		return nil, err
	}

	// Send filled request to the underlying abstract.Scanner
	ctx := query.RequestContext()
	document, err := srv.options.Scanner.Scan(ctx, *filled)
	if err != nil {
		query.Reject(http.StatusConflict, err)
		return nil, err
	}

	// Store document and update status
	srv.document = document
	srv.status.ScannerState = Processing

	// Convert the filled request back to DocumentParameters so the
	// response reflects the actual parameters used for the scan.
	finalTicket := fromAbstractScannerRequest(filled)
	finalTicket.JobDescription = req.ScanTicket.JobDescription

	// Register job in history
	jobID := srv.nextJobID
	srv.nextJobID++
	jobToken := uuid.Random().URN()

	srv.jobs.put(jobInfo{
		jobID:       jobID,
		jobToken:    jobToken,
		state:       JobStateProcessing,
		scanTicket:  finalTicket,
		createdTime: time.Now(),
	})

	srv.status.ScannerState = Processing

	return &CreateScanJobResponse{
		DocumentFinalParameters: optional.Get(finalTicket.DocumentParameters),
		JobID:                   jobID,
		JobToken:                jobToken,
	}, nil
}

// handleRetrieveImageRequest handles RetrieveImage requests.
func (srv *AbstractServer) handleRetrieveImageRequest(
	query *transport.ServerQuery,
	req *RetrieveImageRequest,
) (Body, error) {

	srv.lock.Lock()

	job := srv.jobs.get(req.JobID)
	if srv.document == nil || job == nil || req.JobToken != job.jobToken {
		srv.lock.Unlock()
		query.Reject(http.StatusNotFound, nil)
		return nil, errInvalidJob
	}

	// Get next document file
	file, err := srv.document.Next()
	srv.lock.Unlock()

	switch {
	case err == io.EOF:
		srv.finish(req.JobID, JobStateCompleted)
		query.Reject(http.StatusNotFound, nil)
		return nil, err
	case err != nil:
		srv.finish(req.JobID, JobStateAborted)
		query.Reject(http.StatusServiceUnavailable, err)
		return nil, err
	}

	// Increment scansCompleted for this job
	srv.lock.Lock()
	if j := srv.jobs.get(req.JobID); j != nil {
		j.scansCompleted++
	}
	srv.lock.Unlock()

	return &RetrieveImageResponse{
		ScanData:    ScanData{ContentID: uuid.Random().String()},
		Image:       io.NopCloser(file),
		ContentType: file.Format(),
	}, nil
}

// handleCancelJobRequest handles CancelJob requests.
func (srv *AbstractServer) handleCancelJobRequest(
	query *transport.ServerQuery,
	req *CancelJobRequest,
) (Body, error) {

	srv.lock.Lock()
	job := srv.jobs.get(req.JobID)
	active := job.state == JobStateProcessing
	srv.lock.Unlock()

	if !active {
		query.Reject(http.StatusNotFound, nil)
		return nil, errInvalidJob
	}

	srv.finish(req.JobID, JobStateCanceled)
	return &CancelJobResponse{}, nil
}

// handleGetActiveJobsRequest handles GetActiveJobs requests.
func (srv *AbstractServer) handleGetActiveJobsRequest(
	query *transport.ServerQuery,
	req *GetActiveJobsRequest,
) (Body, error) {

	srv.lock.Lock()
	defer srv.lock.Unlock()

	var summaries []JobSummary
	for _, j := range srv.jobs {
		if j.state == JobStateProcessing {
			summaries = append(summaries, jobSummaryFrom(j))
		}
	}

	return &GetActiveJobsResponse{
		ActiveJobs: ActiveJobs{JobSummary: summaries},
	}, nil
}

// handleGetJobHistoryRequest handles GetJobHistory requests.
func (srv *AbstractServer) handleGetJobHistoryRequest(
	query *transport.ServerQuery,
	req *GetJobHistoryRequest,
) (Body, error) {

	srv.lock.Lock()
	defer srv.lock.Unlock()

	var history []JobSummary
	for _, j := range srv.jobs {
		if j.state != JobStateProcessing {
			history = append(history, jobSummaryFrom(j))
		}
	}

	return &GetJobHistoryResponse{JobHistory: history}, nil
}

// jobSummaryFrom builds a [JobSummary] from a [jobInfo].
func jobSummaryFrom(j jobInfo) JobSummary {
	return JobSummary{
		JobID:                  j.jobID,
		JobName:                j.scanTicket.JobDescription.JobName,
		JobOriginatingUserName: j.scanTicket.JobDescription.JobOriginatingUserName,
		JobState:               j.state,
		ScansCompleted:         j.scansCompleted,
	}
}

// finish closes the current document, updates the job state, and
// resets the server to idle.
func (srv *AbstractServer) finish(jobID int, state JobState) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if srv.document != nil {
		srv.document.Close()
		srv.document = nil
	}

	srv.status.ScannerState = Idle

	if j := srv.jobs.get(jobID); j != nil {
		j.state = state
		j.completedTime = time.Now()
	}
}

// sendSOAPResponse wraps a response body in a SOAP envelope and sends it.
// If the body is a [RetrieveImageResponse], it sends an MTOM/XOP
// multipart message with the image as a binary attachment.
func (srv *AbstractServer) sendSOAPResponse(
	query *transport.ServerQuery,
	req Message,
	body Body) {

	rsp := Message{
		Header: Header{
			Action:    body.Action(),
			MessageID: AnyURI(uuid.Random().URN()),
			To:        optional.New(AnyURI(AddrAnonymous)),
			RelatesTo: optional.New(req.Header.MessageID),
		},
		Body: body,
	}

	// RetrieveImageResponse requires MTOM/XOP multipart encoding
	if _, ok := body.(*RetrieveImageResponse); ok {
		rsp.writeMTOM(query)
		return
	}

	query.ResponseHeader().Set("Content-Type", "application/soap+xml")
	query.SendXML(http.StatusOK, NsMap, rsp.toXML())

	// Notify tracer
	trace.OnResponse(query, traceMessage{rsp}, nil)
}
