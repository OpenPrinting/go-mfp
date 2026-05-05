// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// HTTP Query context

package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// nextQueryID maintains is the next Query ID
var nextQueryID atomic.Uint64

// ServerQuery maintains the HTTP [Server] query processing,
// allowing per-request centralized logging and hooking.
//
// It's interface rules is stronger that of the [http.ResponseWriter]:
//   - [ServerQuery.WriteHeader] must be explicitly called before
//     the body can be written.
//   - [ServerQuery.WriteHeader] cannot be called twice.
//
// It keeps the reference to the original [http.Request] and wraps
// the corresponding [http.ResponseWriter], passed to the
// [http.Handler.ServeHTTP].
type ServerQuery struct {
	log          *log.Record          // Log record for the query
	id           uint64               // Query ID
	logprefix    string               // Log prefix
	rq           *http.Request        // Incoming request
	w            http.ResponseWriter  // Underlying http.ResponseWriter
	status       atomic.Int32         // HTTP status, 0 if not known yet
	finished     atomic.Bool          // True if query already finished
	oncompletion []func(*ServerQuery) // Called on query completion
}

// NewServerQuery returns the new [ServerQuery].
func NewServerQuery(w http.ResponseWriter, rq *http.Request) *ServerQuery {
	ctx := rq.Context()
	query := &ServerQuery{
		log:       log.Begin(ctx),
		id:        nextQueryID.Add(1),
		logprefix: "HTTP-SRVR",
		rq:        rq,
		w:         w,
	}

	return query
}

// ID returns the unique identifier assigned to each [ServerQuery].
//
// IDs are intended to correlate query activities for debugging and
// diagnostics.
//
// The only guarantee is that the ID remains unique throughout the
// process lifetime.
func (query *ServerQuery) ID() uint64 {
	return query.id
}

// SetLogPrefix sets log prefix used for log messages generated
// during the [ServerQuery] processing. Default is "HTTP-SRVR".
func (query *ServerQuery) SetLogPrefix(prefix string) {
	query.logprefix = prefix
}

// OnCompletion registers a callback to be invoked when the [ServerQuery] completes.
//
// Multiple callbacks can be registered and are executed in reverse order of
// their registration (last registered, first called).
//
// The callback must not modify the query.
func (query *ServerQuery) OnCompletion(cb func(*ServerQuery)) {
	query.oncompletion = append(query.oncompletion, cb)
}

// Request returns the underlying Request.
// Caller should not modify the request obtained this way.
func (query *ServerQuery) Request() *http.Request {
	return query.rq
}

// RequestContext returns the Request context.
func (query *ServerQuery) RequestContext() context.Context {
	return query.rq.Context()
}

// RequestURL returns the Request URL.
//
// The returned URL is taken unmodified from the underlying [http.Request],
// so in most cases it will not contain scheme and post parts.
func (query *ServerQuery) RequestURL() *url.URL {
	return query.rq.URL
}

// RequestFullURL returns the Request full URL, with the reconstructed
// Scheme and Host parts,
func (query *ServerQuery) RequestFullURL() *url.URL {
	u := URLClone(query.rq.URL)
	u.Scheme = query.RequestScheme()
	u.Host = query.RequestHost()

	return u
}

// RequestMethod returns the Request Method.
func (query *ServerQuery) RequestMethod() string {
	return query.rq.Method
}

// RequestScheme returns the Request Scheme.
func (query *ServerQuery) RequestScheme() string {
	if query.rq.TLS != nil {
		return "https"
	}
	return "http"
}

// RequestHost returns the Request Host.
func (query *ServerQuery) RequestHost() string {
	return query.rq.Host
}

// RequestContentType returns the Request content type,
// normalized to the lowercase.
func (query *ServerQuery) RequestContentType() string {
	return strings.ToLower(query.rq.Header.Get("Content-Type"))
}

// RequestContentLength returns the Request content Length.
func (query *ServerQuery) RequestContentLength() int64 {
	return query.rq.ContentLength
}

// RequestHeader returns http.Header of the request
func (query *ServerQuery) RequestHeader() http.Header {
	return query.rq.Header
}

// Finish is called when query processing is finished.
//
// The following simple query completion methods call
// it automatically:
//   - [ServerQuery.Reject]
//   - [ServerQuery.Created]
//   - [ServerQuery.SendXML]
//   - [ServerQuery.SendData]
func (query *ServerQuery) Finish() {
	// Don't allow query to be finished twice.
	//
	// FIXME. The present API is unclear and overcomplicated.
	// It requires Finish to be called on some paths and
	// implies the automatic finish on others.
	if !query.finished.CompareAndSwap(false, true) {
		return
	}

	// Call OnCompletion callbacks
	for i := len(query.oncompletion) - 1; i >= 0; i-- {
		query.oncompletion[i](query)
	}

	// Commit the log message
	query.log.Commit()
}

// DumpRequest creates a request dump, for logging.
func (query *ServerQuery) DumpRequest() []byte {
	dump, _ := httputil.DumpRequest(query.Request(), false)
	return dump
}

// RequestBody returns body of the http.Request
func (query *ServerQuery) RequestBody() io.ReadCloser {
	return query.rq.Body
}

// ResponseHeader returns http.Header of the response
func (query *ServerQuery) ResponseHeader() http.Header {
	return query.w.Header()
}

// ResponseStatus returns the HTTP response status of the query.
// It returns 0, if status is not yet set.
func (query *ServerQuery) ResponseStatus() int {
	return int(query.status.Load())
}

// IsStatusSet returns true if the HTTP response status is set.
func (query *ServerQuery) IsStatusSet() bool {
	return query.ResponseStatus() != 0
}

// assertStatusSet panics if HTTP response status is not set.
func (query *ServerQuery) assertStatusSet() {
	assert.MustMsg(query.IsStatusSet(),
		"ServerQuery.WriteHeader: HTTP status is not set")
}

// assertStatusNotSet panics if HTTP response status is set.
func (query *ServerQuery) assertStatusNotSet() {
	assert.MustMsg(!query.IsStatusSet(),
		"ServerQuery.WriteHeader: HTTP status already set")
}

// Write writes response body bytes.
func (query *ServerQuery) Write(data []byte) (int, error) {
	query.assertStatusSet()
	return query.w.Write(data)
}

// WriteHeader writes HTTP response header.
func (query *ServerQuery) WriteHeader(status int) {
	assert.MustMsg(status != 0,
		"ServerQuery.WriteHeader: invalid HTTP status %d", status)

	if query.status.CompareAndSwap(0, int32(status)) {
		query.w.WriteHeader(status)
		query.log.Debug("%s %s %s -- %d %s",
			query.logprefix,
			query.rq.Method, query.rq.URL,
			status, http.StatusText(status))
		query.log.Flush()
		return
	}

	query.assertStatusNotSet()
}

// NoCache set response headers to disable client-side response cacheing.
func (query *ServerQuery) NoCache() {
	hdr := query.ResponseHeader()
	hdr.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	hdr.Set("Pragma", "no-cache")
	hdr.Set("Expires", "0")
}

// Reject completes request with a error.
//
// The error parameter is sent as the text body of the HTTP
// response. If it is nil, the reasonable default will be provided.
func (query *ServerQuery) Reject(status int, err error) {
	query.ResponseHeader().Set("Content-Type", "text/plain; charset=utf-8")
	query.NoCache()
	query.WriteHeader(status)

	if err == nil {
		err = errors.New(http.StatusText(status))
	}

	s := fmt.Sprintf("%3.3d %s\n", status, err)
	query.Write([]byte(s))
	query.Write([]byte("\n"))

	query.Finish()
}

// Created completes request with the http.StatusCreated
// status and Location: URL
func (query *ServerQuery) Created(location string) {
	scheme := "http"
	if query.rq.TLS != nil {
		scheme = "https"
	}

	location = scheme + "://" + query.rq.Host + location

	query.ResponseHeader().Set("Location", location)
	query.WriteHeader(http.StatusCreated)

	query.Finish()
}

// SendXML sends the XML response and completes the request.
func (query *ServerQuery) SendXML(
	status int, ns xmldoc.Namespace, xml xmldoc.Element) {

	query.ResponseHeader().Set("Content-Type", "text/xml")
	query.WriteHeader(status)
	xml.EncodeIndent(query, ns, "  ")

	query.Finish()
}

// SendData sends the binary data, represented as [io.Reader]
// and completes the request.
func (query *ServerQuery) SendData(
	status int, contentType string, data io.Reader) {

	query.ResponseHeader().Set("Content-Type", contentType)
	query.WriteHeader(status)
	io.Copy(query, data)

	query.Finish()
}
