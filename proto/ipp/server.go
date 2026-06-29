// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP server

package ipp

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"net/http/httputil"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/log/trace"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/goipp"
)

// Server represents the IPP server.
type Server struct {
	options ServerOptions         // Server options
	ops     map[goipp.Op]*Handler // Installed handlers
}

// ServerOptions allows to specify options that can modify
// the [Server] behavior.
type ServerOptions struct {
	// Hooks defines IPP server hooks. See [ServerHooks]
	// for details.
	Hooks ServerHooks
}

// NewServer returns a new Sever.
func NewServer(options ServerOptions) *Server {
	s := &Server{
		options: options,
		ops:     make(map[goipp.Op]*Handler),
	}
	return s
}

// ServeHTTP handles incoming HTTP request. It implements
// [http.Handler] interface.
//
// Using this interface, the [Server] can work on a top of
// existent [http.Server], [http.ServeMux], [transport.PathMux]
// and so on.
func (s *Server) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	// Setup things
	query := transport.NewServerQuery(w, rq)
	query.SetLogPrefix("IPP")
	ctx := query.RequestContext()

	// Dump request HTTP headers
	dump, _ := httputil.DumpRequest(query.Request(), false)
	log.Debug(ctx, "IPP request received:")
	log.Debug(ctx, "%s", dump)

	// Call the OnHTTPRequest hook
	if s.options.Hooks.OnHTTPRequest != nil {
		s.options.Hooks.OnHTTPRequest(query)
		if query.IsStatusSet() {
			return
		}
	}

	// Check HTTP parameters
	if query.RequestMethod() != "POST" {
		s.httpError(query, ErrHTTPMethodNotAllowed)
		return
	}

	ctype := query.RequestContentType()
	mediatype, _, _ := mime.ParseMediaType(ctype)
	if mediatype != goipp.ContentType {
		err := NewErrHTTP(http.StatusUnsupportedMediaType,
			fmt.Sprintf("unsupported media type: %q", ctype))
		s.httpError(query, err)
		return
	}

	// Decode IPP message
	msg := &goipp.Message{}
	err := msg.Decode(query.RequestBody())
	if err != nil {
		s.httpError(query, err)
		return
	}

	// Notify tracer, if present
	body := query.RequestBody()
	body = trace.OnRequest(query, goippRequest{msg}, body)
	defer body.Close()

	// Call the OnIPPRequest hook
	if s.options.Hooks.OnIPPRequest != nil {
		msg2 := s.options.Hooks.OnIPPRequest(query, msg)
		if query.IsStatusSet() {
			return
		}
		if msg2 != nil {
			msg = msg2
		}
	}

	// Log the IPP request
	var buf bytes.Buffer
	msg.Print(&buf, true)
	log.Debug(ctx, "IPP request message:")
	log.Debug(ctx, buf.String())

	// Check IPP parameters
	if msg.RequestID == 0 {
		err := NewErrIPPFromMessage(msg,
			goipp.StatusErrorBadRequest,
			"bad request ID %d", msg.RequestID)

		s.httpError(query, err)
		return
	}

	if msg.Version < goipp.MakeVersion(1, 0) ||
		msg.Version > goipp.DefaultVersion {
		err := NewErrIPPFromMessage(msg,
			goipp.StatusErrorVersionNotSupported,
			"bad request version %s", msg.Version)

		s.httpError(query, err)
		return
	}

	handler := s.ops[goipp.Op(msg.Code)]
	if handler == nil {
		op := goipp.Op(msg.Code)
		err := NewErrIPPFromMessage(msg,
			goipp.StatusErrorOperationNotSupported,
			"unsupported operation %s", op)

		s.httpError(query, err)
		return
	}

	// Handle the message
	rsp, err := handler.handle(ctx, msg, body)
	if err != nil {
		s.httpError(query, err)
		return
	}

	// Close the body. It will notify tracer that request is
	// fully consumed, so tracer can finish writing it.
	body.Close()
	body = nil

	// Call the OnIPPResponse hook
	if s.options.Hooks.OnIPPResponse != nil {
		rsp2 := s.options.Hooks.OnIPPResponse(query, rsp)
		if query.IsStatusSet() {
			return
		}
		if rsp2 != nil {
			rsp = rsp2
		}
	}

	// Log the response
	buf.Reset()
	rsp.Print(&buf, false)
	log.Debug(ctx, "IPP response message:")
	log.Debug(ctx, buf.String())

	// Send response
	query.ResponseHeader().Set("Content-Type", "application/ipp")
	query.WriteHeader(http.StatusOK) // At HTTP level everything OK.

	// Notify tracer, if present (must be after WriteHeader so
	// DumpResponse can read the correct response status).
	trace.OnResponse(query, goippResponse{rsp}, nil)

	err = rsp.Encode(query)
	if err != nil {
		log.Error(ctx, "IPP error sending response: %s", err)
	}

	query.Finish()
}

// RegisterHandler adds the request [Handler].
func (s *Server) RegisterHandler(handler *Handler) {
	s.ops[handler.Op] = handler
}

// httpError finishes HTTP request with an error.
func (s *Server) httpError(query *transport.ServerQuery, err error) {
	switch err := err.(type) {
	case *ErrHTTP:
		query.Reject(err.Status, err)

	case *ErrIPP:
		// Create the IPP response
		rsp := err.Encode()

		// Notify tracer if present
		trace.OnResponse(query, goippResponse{rsp}, nil)

		// Call OnIPPResponse hook
		if s.options.Hooks.OnIPPResponse != nil {
			rsp2 := s.options.Hooks.OnIPPResponse(query, rsp)
			if query.IsStatusSet() {
				return
			}
			if rsp2 != nil {
				rsp = rsp2
			}
		}

		// Log the response
		var buf bytes.Buffer
		ctx := query.RequestContext()

		rsp.Print(&buf, false)
		log.Debug(ctx, "IPP response message:")
		log.Debug(ctx, buf.String())

		// Finish the HTTP query
		query.ResponseHeader().Set("Content-Type", "application/ipp")
		query.WriteHeader(http.StatusOK) // At HTTP level everything OK.

		rsp.Encode(query)

	default:
		query.Reject(http.StatusInternalServerError, err)
	}
}
