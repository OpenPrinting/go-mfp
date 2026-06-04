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
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/optional"
)

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

// TestDocumentReceiverCalled verifies that the DocumentReceiver callback
// is invoked with the correct jobID, format, and data.
func TestDocumentReceiverCalled(t *testing.T) {
	printer := testNewCaptPrinter(t)

	var (
		gotJobID  int
		gotFormat string
		gotData   []byte
		called    bool
	)

	printer.SetDocumentReceiver(func(jobID int, format string, data []byte) {
		called = true
		gotJobID = jobID
		gotFormat = format
		gotData = append([]byte(nil), data...)
	})

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

	wantJobID := createRsp.Job.JobID

	// Step 2: Send-Document
	wantData := []byte("Hello, virtual printer!")
	wantFormat := "application/pdf"

	sendRq := &SendDocumentRequest{
		RequestHeader:  DefaultRequestHeader,
		PrinterURI:     optional.New(ippURI),
		JobID:          optional.New(wantJobID),
		DocumentFormat: optional.New(wantFormat),
		LastDocument:   true,
		Job:            &JobAttributes{},
	}
	sendRq.Body = bytes.NewReader(wantData)

	sendRsp := &SendDocumentResponse{}
	if err := client.Do(ctx, sendRq, sendRsp); err != nil {
		t.Fatalf("Send-Document: %v", err)
	}

	// Verify receiver was called
	if !called {
		t.Fatal("DocumentReceiver was not called")
	}
	if gotJobID != wantJobID {
		t.Errorf("jobID: got %d, want %d", gotJobID, wantJobID)
	}
	if gotFormat != wantFormat {
		t.Errorf("format: got %q, want %q", gotFormat, wantFormat)
	}
	if !bytes.Equal(gotData, wantData) {
		t.Errorf("data: got %q, want %q", gotData, wantData)
	}
}

// TestDocumentReceiverNilNoPanic verifies that a nil receiver
// does not panic when a document arrives.
func TestDocumentReceiverNilNoPanic(t *testing.T) {
	printer := testNewCaptPrinter(t)
	// No SetDocumentReceiver — receiver stays nil

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
		t.Fatalf("Send-Document with nil receiver: %v", err)
	}
	// reaching here without panic = success
}

// TestDocumentReceiverLargeDoc verifies that large documents
// are captured fully without truncation.
func TestDocumentReceiverLargeDoc(t *testing.T) {
	printer := testNewCaptPrinter(t)

	var gotData []byte
	printer.SetDocumentReceiver(func(_ int, _ string, data []byte) {
		gotData = append([]byte(nil), data...)
	})

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

	if len(gotData) != len(wantData) {
		t.Errorf("data length: got %d, want %d", len(gotData), len(wantData))
	}
	if !bytes.Equal(gotData, wantData) {
		t.Error("large document data mismatch")
	}
}
