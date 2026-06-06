// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Tests for document capture

package ipp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// testBackend is a test implementation of abstract.PrintBackend.
type testBackend struct {
	called bool
	params abstract.PrintJobParams
	data   []byte
	err    error
}

func (b *testBackend) PrintDocument(
	params abstract.PrintJobParams, body io.Reader) error {
	b.called = true
	b.params = params
	b.data, b.err = io.ReadAll(body)
	return b.err
}

// testNewCaptPrinter creates a minimal Printer for capture testing.
func testNewCaptPrinter(t *testing.T) *Printer {
	t.Helper()
	return NewPrinter(&PrinterAttributes{}, PrinterOptions{})
}

// testCaptPrinterURL returns the HTTP URL and IPP URI for the test server.
func testCaptPrinterURL(srv *httptest.Server) (*url.URL, string) {
	httpURL, _ := url.Parse(srv.URL + "/ipp/print")
	ippURI := fmt.Sprintf("ipp://%s/ipp/print", srv.Listener.Addr())
	return httpURL, ippURI
}

// TestDocumentReceiverCalled verifies that the PrintBackend is invoked
// with the correct format and data.
func TestDocumentReceiverCalled(t *testing.T) {
	printer := testNewCaptPrinter(t)
	backend := &testBackend{}
	printer.SetPrintBackend(backend)

	srv := httptest.NewServer(printer)
	defer srv.Close()

	httpURL, ippURI := testCaptPrinterURL(srv)
	client := NewClient(httpURL, nil)
	ctx := context.Background()

	// Step 1: Create-Job
	createRq := &CreateJobRequest{
		RequestHeader: DefaultRequestHeader,
		JobCreateOperation: JobCreateOperation{
			PrinterURI: ippURI,
		},
		Job: &JobAttributes{},
	}
	createRsp := &CreateJobResponse{}
	if err := client.Do(ctx, createRq, createRsp); err != nil {
		t.Fatalf("Create-Job: %v", err)
	}

	// Step 2: Send-Document
	wantData := []byte("Hello, virtual printer!")
	wantFormat := "application/pdf"

	sendRq := &SendDocumentRequest{
		RequestHeader:  DefaultRequestHeader,
		PrinterURI:     optional.New(ippURI),
		JobID:          optional.New(createRsp.Job.JobID),
		DocumentFormat: optional.New(wantFormat),
		LastDocument:   true,
		Job:            &JobAttributes{},
	}
	sendRq.Body = bytes.NewReader(wantData)

	sendRsp := &SendDocumentResponse{}
	if err := client.Do(ctx, sendRq, sendRsp); err != nil {
		t.Fatalf("Send-Document: %v", err)
	}

	if !backend.called {
		t.Fatal("PrintBackend was not called")
	}
	if backend.params.Format != wantFormat {
		t.Errorf("format: got %q, want %q",
			backend.params.Format, wantFormat)
	}
	if !bytes.Equal(backend.data, wantData) {
		t.Errorf("data: got %q, want %q", backend.data, wantData)
	}
}

// TestDocumentReceiverNilNoPanic verifies that a nil backend
// does not panic when a document arrives.
func TestDocumentReceiverNilNoPanic(t *testing.T) {
	printer := testNewCaptPrinter(t)
	// No SetPrintBackend — backend stays nil

	srv := httptest.NewServer(printer)
	defer srv.Close()

	httpURL, ippURI := testCaptPrinterURL(srv)
	client := NewClient(httpURL, nil)
	ctx := context.Background()

	createRq := &CreateJobRequest{
		RequestHeader: DefaultRequestHeader,
		JobCreateOperation: JobCreateOperation{
			PrinterURI: ippURI,
		},
		Job: &JobAttributes{},
	}
	createRsp := &CreateJobResponse{}
	if err := client.Do(ctx, createRq, createRsp); err != nil {
		t.Fatalf("Create-Job: %v", err)
	}

	sendRq := &SendDocumentRequest{
		RequestHeader: DefaultRequestHeader,
		PrinterURI:    optional.New(ippURI),
		JobID:         optional.New(createRsp.Job.JobID),
		LastDocument:  true,
		Job:           &JobAttributes{},
	}
	sendRq.Body = bytes.NewReader([]byte("test data"))

	sendRsp := &SendDocumentResponse{}
	if err := client.Do(ctx, sendRq, sendRsp); err != nil {
		t.Fatalf("Send-Document with nil backend: %v", err)
	}
	// reaching here without panic = success
}

// TestDocumentReceiverLargeDoc verifies that large documents
// are streamed fully without truncation.
func TestDocumentReceiverLargeDoc(t *testing.T) {
	printer := testNewCaptPrinter(t)
	backend := &testBackend{}
	printer.SetPrintBackend(backend)

	srv := httptest.NewServer(printer)
	defer srv.Close()

	httpURL, ippURI := testCaptPrinterURL(srv)
	client := NewClient(httpURL, nil)
	ctx := context.Background()

	createRq := &CreateJobRequest{
		RequestHeader: DefaultRequestHeader,
		JobCreateOperation: JobCreateOperation{
			PrinterURI: ippURI,
		},
		Job: &JobAttributes{},
	}
	createRsp := &CreateJobResponse{}
	if err := client.Do(ctx, createRq, createRsp); err != nil {
		t.Fatalf("Create-Job: %v", err)
	}

	// 1 MB document
	wantData := bytes.Repeat([]byte("A"), 1024*1024)

	sendRq := &SendDocumentRequest{
		RequestHeader: DefaultRequestHeader,
		PrinterURI:    optional.New(ippURI),
		JobID:         optional.New(createRsp.Job.JobID),
		LastDocument:  true,
		Job:           &JobAttributes{},
	}
	sendRq.Body = bytes.NewReader(wantData)

	sendRsp := &SendDocumentResponse{}
	if err := client.Do(ctx, sendRq, sendRsp); err != nil {
		t.Fatalf("Send-Document: %v", err)
	}

	if len(backend.data) != len(wantData) {
		t.Errorf("data length: got %d, want %d",
			len(backend.data), len(wantData))
	}
	if !bytes.Equal(backend.data, wantData) {
		t.Error("large document data mismatch")
	}
}
