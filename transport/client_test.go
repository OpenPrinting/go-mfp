// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// HTTP client wrapper test.

package transport

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
)

// TestNewClient tests NewClient function
func TestNewClient(t *testing.T) {
	// NewClient(nil) must create a new Transport
	clnt := NewClient(nil)
	if clnt.Transport == nil {
		t.Errorf("NewClient(nil): clnt.Transport == nil")
	}

	// NewClient(tr) must use provided Transport
	tr := NewTransport(nil)
	clnt = NewClient(tr)
	if clnt.Transport != tr {
		t.Errorf("NewClient(tr): clnt.Transport != tr")
	}
}

// TestClientRedirect tests that HTTP redirects are properly
// handled by the Client
func TestClientRedirect(t *testing.T) {
	const (
		badhost  = "badhost"
		goodhost = "goodhost"
	)

	// Create loopbacked Client and Server
	tr, l := NewLoopback()

	handler := func(w http.ResponseWriter, rq *http.Request) {
		fmt.Println(">>>")
		rq.Write(os.Stdout)

		switch rq.Host {
		case badhost:
			http.Redirect(w, rq, "http://goodhost/",
				http.StatusFound)
		case goodhost:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	srv := &http.Server{Handler: http.HandlerFunc(handler)}
	clnt := NewClient(tr)

	// Run srv.Serve() on its own goroutine
	var done sync.WaitGroup
	done.Add(1)

	go func() {
		srv.Serve(l)
		done.Done()
	}()

	defer func() {
		srv.Close()
		done.Wait()
	}()

	// Perform HTTP transaction and test results
	resp, err := clnt.Get("http://badhost/")
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		t.Errorf("HTTP error: %s", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("HTTP response status mismatch:\n"+
			"expected: %d\n"+
			"present:  %d",
			http.StatusOK, resp.StatusCode)
		return
	}
}
