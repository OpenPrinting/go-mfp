// MFP - Multi-Function Printers and scanners toolkit
// ipp-usb API client
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Helpers for logging

package ippusbapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
)

// jsonLogger formats JSON text for logging.
// It implements the [log.Marshaler] interface.
type jsonLogger []byte

// MarshalLog implements log.Marshaler interface for jsonLogger.
func (jsn jsonLogger) MarshalLog() []byte {
	buf := bytes.Buffer{}
	if err := json.Indent(&buf, jsn, "", "  "); err != nil {
		return jsn
	}

	return buf.Bytes()
}

// httpOutRequestLogger formats outgoing http.Request for logging
type httpOutRequestLogger http.Request

// MarshalLog implements log.Marshaler interface for httpOutRequestLogger.
func (rq *httpOutRequestLogger) MarshalLog() []byte {
	dump, _ := httputil.DumpRequestOut((*http.Request)(rq), false)
	return dump
}

// httpResponseLogger formats outgoing http.Request for logging
type httpResponseLogger http.Response

// MarshalLog implements log.Marshaler interface for httpResponseLogger.
func (rsp *httpResponseLogger) MarshalLog() []byte {
	dump, _ := httputil.DumpResponse((*http.Response)(rsp), false)
	return dump
}
