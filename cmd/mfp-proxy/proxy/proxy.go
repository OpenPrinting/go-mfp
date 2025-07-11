// MFP - Miulti-Function Printers and scanners toolkit
// The "proxy" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package proxy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/goipp"
)

// proxy implements an IPP/eSCL/WSD proxy
type proxy struct {
	ctx       context.Context   // Logging/shutdown context
	trace     *traceWriter      // Trace writer (may be nil)
	cancel    func()            // ctx cancel function
	m         mapping           // Local/remote mapping
	l         net.Listener      // TCP listener for incoming connections
	srv       *transport.Server // HTTP server for incoming connections
	clnt      *transport.Client // HTTP client part of proxy
	closeWait sync.WaitGroup    // Wait for proxy.Close completion
	rqnum     atomic.Uint32     // Request number, for logging
}

// newProxy creates a new proxy for the specified mapping.
func newProxy(ctx context.Context, m mapping, trace *traceWriter) (
	*proxy, error) {
	log.Debug(ctx, "proxy started: %d->%s", m.localPort, m.targetURL)

	// Create TCP listener
	l, err := newListener(ctx, m.localPort)
	if err != nil {
		return nil, err
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(ctx)

	// Create proxy structure
	p := &proxy{
		ctx:    ctx,
		trace:  trace,
		cancel: cancel,
		m:      m,
		l:      l,
		clnt:   transport.NewClient(nil),
	}

	// Ensure cancellation propagation
	p.closeWait.Add(1)
	go p.kill()

	// Start HTTP server
	p.srv = transport.NewServer(nil, p)

	p.closeWait.Add(1)
	go func() {
		p.srv.Serve(l)
		p.closeWait.Done()
	}()

	return p, nil
}

// kill closes the proxy and terminates all active session when proxy.ctx
// is canceled.
func (p *proxy) kill() {
	<-p.ctx.Done()

	p.srv.Close()

	p.closeWait.Done()
}

// Shutdown performs proxy shutdown.
func (p *proxy) Shutdown() {
	p.cancel()
	p.closeWait.Wait()

	log.Debug(p.ctx, "proxy finished: %d->%s",
		p.m.localPort, p.m.targetURL)
}

// ServeHTTP handles incoming HTTP requests.
// It implements [http.Handler] interface.
func (p *proxy) ServeHTTP(w http.ResponseWriter, in *http.Request) {
	// Catch panics to log
	defer func() {
		v := recover()
		if v != nil {
			log.Panic(p.ctx, v)
		}
	}()

	// Handle request
	log.Debug(p.ctx, "%s %s", in.Method, in.URL)

	ct := strings.ToLower(in.Header.Get("Content-Type"))

	switch {
	case p.m.proto == protoIPP && in.Method == "POST" &&
		ct == "application/ipp":
		p.doIPP(w, in)

	case in.Method == "GET":
		p.doHTTP(w, in)

	default:
		p.httpReject(w, in,
			http.StatusBadRequest, errors.New("Bad Request"))
	}
}

// outreq creates an outgoing HTTP request based on request
// received by the server side of proxy.
func (p *proxy) outreq(in *http.Request, body io.ReadCloser) *http.Request {
	// Create request
	out, _ := transport.NewRequest(p.ctx, in.Method, in.URL, body)
	out.Header = in.Header.Clone()
	p.httpRemoveHopByHopHeaders(out.Header)

	// Adjust target URL
	prq := httputil.ProxyRequest{
		Out: out,
	}
	prq.SetURL(p.m.targetURL)
	out.Host = out.URL.Host

	return out
}

// msgxlat returns goipp.Message translator that rewrites message
// attributes when message is being forwarded via proxy.
//
// Currently, only URLs embedded into the message are translated.
func (p *proxy) msgxlat(in *http.Request) (*msgXlat, error) {
	s := "http://" + in.Host
	u, err := transport.ParseURL(s)
	if err != nil {
		err = fmt.Errorf("%q: can't parse local URL", s)
		return nil, err
	}

	urlxlat := transport.NewURLXlat(u, p.m.targetURL)
	msgxlat := newMsgXlat(urlxlat)

	return msgxlat, nil
}

// doHTTP implements proxy for the bare HTTP requests
func (p *proxy) doHTTP(w http.ResponseWriter, in *http.Request) {
	// Dump request headers
	p.httpLogRequest("HTTP", in)

	// Prepare outgoing request
	out := p.outreq(in, in.Body)
	out.ContentLength = in.ContentLength

	// Execute outgoing request
	log.Debug(p.ctx, "HTTP: forward request to: %s", out.URL)

	rsp, err := p.clnt.Do(out)
	if err != nil {
		log.Debug(p.ctx, "IPP: %s", err)
		p.httpReject(w, in, http.StatusBadGateway, err)
		return
	}

	// Copy response headers and status to the client
	p.httpRemoveHopByHopHeaders(rsp.Header)
	p.httpCopyHeaders(w.Header(), rsp.Header)

	if rsp.ContentLength >= 0 {
		rsp.Header.Set("Content-Length",
			strconv.FormatInt(rsp.ContentLength, 10))
	}

	w.WriteHeader(rsp.StatusCode)

	// Dump response headers
	p.httpLogResponse("HTTP", rsp)

	// Forward response body
	io.Copy(w, rsp.Body)
	rsp.Body.Close()
}

// doIPP implements proxy for IPP requests
func (p *proxy) doIPP(w http.ResponseWriter, in *http.Request) {
	rqnum := p.rqnum.Add(1)

	// Create goipp.Message translator
	msgxlat, err := p.msgxlat(in)
	if err != nil {
		p.httpReject(w, in, http.StatusBadGateway, err)
		return
	}

	// Prepare outgoing request
	out, ipplen, err := p.doIPPreq(in, msgxlat, rqnum)
	if err != nil {
		err = fmt.Errorf("IPP error: %w", err)
		p.httpReject(w, in, http.StatusBadGateway, err)
		return
	}

	// Shiff outgoing data, if trace is active
	var sniffBuff bytes.Buffer
	if p.trace != nil {
		out.Body = transport.TeeReadCloser(out.Body, &sniffBuff)
	}

	// Execute outgoing request
	log.Debug(p.ctx, "IPP: forward request to: %s", out.URL)

	rsp, err := p.clnt.Do(out)
	if err != nil {
		log.Debug(p.ctx, "IPP: %s", err)
		p.httpReject(w, in, http.StatusBadGateway, err)
		return
	}

	// Save sniffed request data
	if p.trace != nil && sniffBuff.Len() > ipplen {
		data := sniffBuff.Bytes()[ipplen:]
		name := fmt.Sprintf("%8.8d-data.%s", rqnum, magic(data))
		p.trace.Send(name, data)
	}

	// Dump response HTTP headers
	p.httpLogResponse("IPP", rsp)

	// Translate IPP response
	ct := strings.ToLower(rsp.Header.Get("Content-Type"))
	if ct == "application/ipp" {
		err = p.doIPPrsp(rsp, msgxlat, rqnum)
		if err != nil {
			log.Debug(p.ctx, "IPP: %s", err)
			p.httpReject(w, in, http.StatusBadGateway, err)
			return
		}
	}

	// Copy response headers and status to the client
	p.httpRemoveHopByHopHeaders(rsp.Header)
	p.httpCopyHeaders(w.Header(), rsp.Header)

	if rsp.ContentLength >= 0 {
		rsp.Header.Set("Content-Length",
			strconv.FormatInt(rsp.ContentLength, 10))
	}

	w.WriteHeader(rsp.StatusCode)

	// Forward response body
	io.Copy(w, rsp.Body)
	rsp.Body.Close()
}

// doIPPreq performs (client->server) part of the IPP request handling
//
// It returns modified request ready to be send to the server,
// length of the IPP part of that request and error, if any.
func (p *proxy) doIPPreq(in *http.Request,
	msgxlat *msgXlat, rqnum uint32) (*http.Request, int, error) {

	ops := goipp.DecoderOptions{EnableWorkarounds: true}

	// Dump request HTTP headers
	p.httpLogRequest("IPP", in)

	// Fetch IPP Request message
	peeker := transport.NewPeeker(in.Body)
	var msg goipp.Message
	err := msg.DecodeEx(peeker, ops)
	if err != nil {
		return nil, 0, err
	}

	// Write trace
	if p.trace != nil {
		name := fmt.Sprintf("%8.8d-%s.ipp",
			rqnum, goipp.Op(msg.Code))
		p.trace.Send(name, peeker.Bytes())
	}

	// Translate IPP message
	msg2, chg := msgxlat.Forward(&msg)

	// Log the message
	var buf bytes.Buffer
	msg2.Print(&buf, true)
	log.Debug(p.ctx, "IPP: request message:")
	log.Debug(p.ctx, buf.String())

	if !chg.Empty() {
		log.Debug(p.ctx, "IPP: translated attributes:")
		log.Object(p.ctx, log.LevelDebug, 4, chg)
	}

	// Setup outgoing request
	msg2bytes, _ := msg2.EncodeBytes()
	peeker.Replace(msg2bytes)

	out := p.outreq(in, peeker)
	out.ContentLength = in.ContentLength
	if out.ContentLength >= 0 {
		out.ContentLength += int64(len(msg2bytes))
		out.ContentLength -= peeker.Count()
	}

	return out, len(msg2bytes), nil
}

// doIPPreq performs (client->server) part of the IPP request handling
func (p *proxy) doIPPrsp(rsp *http.Response,
	msgxlat *msgXlat, rqnum uint32) error {

	ops := goipp.DecoderOptions{EnableWorkarounds: true}

	// Fetch IPP response message
	peeker := transport.NewPeeker(rsp.Body)
	var msg goipp.Message
	err := msg.DecodeEx(peeker, ops)
	if err != nil {
		peeker.Rewind()
		return err
	}

	// Translate IPP response
	msg2, chg := msgxlat.Reverse(&msg)

	// Log the message
	var buf bytes.Buffer
	msg2.Print(&buf, false)
	log.Debug(p.ctx, "IPP: response message (translated):")
	log.Debug(p.ctx, buf.String())

	if !chg.Empty() {
		log.Debug(p.ctx, "IPP: translated attributes:")
		log.Object(p.ctx, log.LevelDebug, 4, chg)
	}

	// Replace http.Response body
	msg2bytes, _ := msg2.EncodeBytes()
	peeker.Replace(msg2bytes)
	rsp.Body = peeker

	// Adjust rsp.ContentLength
	if rsp.ContentLength >= 0 {
		rsp.ContentLength += int64(len(msg2bytes))
		rsp.ContentLength -= peeker.Count()
	}

	// Write trace
	if p.trace != nil {
		name := fmt.Sprintf("%8.8d-%s.ipp",
			rqnum, goipp.Status(msg.Code))
		p.trace.Send(name, msg2bytes)
	}

	return nil
}

// httpRemoveHopByHopHeaders removes HTTP hop-by-hop headers,
// RFC 7230, section 6.1
func (p *proxy) httpRemoveHopByHopHeaders(hdr http.Header) {
	// Per RFC 7230, section 6.1:
	//
	// Hence, the Connection header field provides a declarative way of
	// distinguishing header fields that are only intended for the immediate
	// recipient ("hop-by-hop") from those fields that are intended for all
	// recipients on the chain ("end-to-end"), enabling the message to be
	// self-descriptive and allowing future connection-specific extensions
	// to be deployed without fear that they will be blindly forwarded by
	// older intermediaries.
	if c := hdr.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				hdr.Del(f)
			}
		}
	}

	// These headers are always considered hop-by-hop.
	for _, c := range []string{"Connection", "Keep-Alive",
		"Proxy-Authenticate", "Proxy-Connection",
		"Proxy-Authorization", "Te", "Trailer", "Transfer-Encoding"} {
		hdr.Del(c)
	}
}

// httpCopyHeaders copies HTTP headers from src to dst
func (p *proxy) httpCopyHeaders(dst, src http.Header) {
	for k, v := range src {
		if strings.ToLower(k) != "content-length" {
			dst[k] = v
		}
	}
}

// httpLogRequest writes http.Request headers the log
func (p *proxy) httpLogRequest(proto string, rq *http.Request) {
	dump, _ := httputil.DumpRequest(rq, false)
	log.Debug(p.ctx, "%s: request received:", proto)
	log.Debug(p.ctx, "%s", dump)
}

// httpLogResponse writes http.Response headers the log
func (p *proxy) httpLogResponse(proto string, rsp *http.Response) {
	dump, _ := httputil.DumpResponse(rsp, false)
	log.Debug(p.ctx, "%s: response received:", proto)
	log.Debug(p.ctx, "%s", dump)
}

// httpReject completes request with a error
func (p *proxy) httpReject(w http.ResponseWriter, in *http.Request,
	status int, err error) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	p.httpNoCache(w)
	w.WriteHeader(status)

	w.Write([]byte(err.Error()))
	w.Write([]byte("\n"))
}

// httpNoCache set response headers to disable client-side
// response cacheing.
func (p *proxy) httpNoCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
