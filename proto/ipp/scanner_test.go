// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Tests for IPP Scan Service

package ipp

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

func testScannerCapabilities() *abstract.ScannerCapabilities {
	profiles := []abstract.SettingsProfile{
		{
			ColorModes: generic.MakeBitset(abstract.ColorModeColor),
			Depths:     generic.MakeBitset(abstract.ColorDepth8),
			Resolutions: []abstract.Resolution{
				{XResolution: 300, YResolution: 300},
			},
		},
	}

	platen := &abstract.InputCapabilities{
		MinWidth:  abstract.A4Width,
		MinHeight: abstract.A4Height,
		MaxWidth:  abstract.A4Width,
		MaxHeight: abstract.A4Height,
		Profiles:  profiles,
	}

	return &abstract.ScannerCapabilities{
		DocumentFormats:   []string{"image/jpeg", "application/pdf"},
		BrightnessRange:   abstract.Range{Min: 0, Max: 100, Normal: 50},
		ContrastRange:     abstract.Range{Min: -127, Max: 127, Normal: 0},
		SharpenRange:      abstract.Range{Min: 0, Max: 100, Normal: 0},
		NoiseRemovalRange: abstract.Range{Min: 0, Max: 100, Normal: 0},
		CompressionRange:  abstract.Range{Min: 0, Max: 100, Normal: 0},
		Platen:            platen,
	}
}

func testNewScanner(t *testing.T, backend abstract.Scanner) *Scanner {
	t.Helper()
	return NewScanner(&PrinterAttributes{}, ScannerOptions{Scanner: backend})
}

func testScannerURL(srv *httptest.Server) (*url.URL, string) {
	httpURL, _ := url.Parse(srv.URL + "/ipp/scan")
	ippURI := fmt.Sprintf("ipp://%s/ipp/scan", srv.Listener.Addr())
	return httpURL, ippURI
}

func testGetNextDocumentData(
	ctx context.Context,
	client *Client,
	ippURI string,
	jobID int,
) (*GetNextDocumentDataResponse, error) {

	const maxAttempts = 20

	for attempt := 0; attempt < maxAttempts; attempt++ {
		rq := &GetNextDocumentDataRequest{
			RequestHeader: DefaultRequestHeader,
			PrinterURI:    optional.New(ippURI),
			JobID:         optional.New(jobID),
		}
		rsp := &GetNextDocumentDataResponse{}

		err := client.DoWithBody(ctx, rq, rsp)
		if err != nil {
			return nil, err
		}

		if rsp.LastDocument || rsp.Header().Body != nil {
			return rsp, nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil, fmt.Errorf("timed out waiting for document data")
}

func testScanJobCreateOperation(ippURI string) JobCreateOperation {
	return JobCreateOperation{
		PrinterURI:             ippURI,
		RequestingUserName:     optional.New("alice"),
		RequestingUserURI:      optional.New("mailto:alice@example.com"),
		JobName:                optional.New("scan-1"),
		IppAttributeFidelity:   optional.New(true),
		CompressionAccepted:    []KwCompression{KwCompressionNone, KwCompressionGzip},
		DocumentFormatAccepted: []string{"image/jpeg", "application/pdf"},
		InputAttributes: optional.New(InputAttributes{
			InputAutoExposure:         optional.New(true),
			InputAutoScaling:          optional.New(false),
			InputAutoSkewCorrection:   optional.New(true),
			InputBrightness:           optional.New(50),
			InputColorMode:            optional.New(KwInputColorModeColor),
			InputContentType:          optional.New(KwInputContentTypeTextAndPhoto),
			InputContrast:             optional.New(-10),
			InputFilmScanMode:         optional.New(KwInputFilmScanModeNotApplicable),
			InputImagesToTransfer:     optional.New(1),
			InputMedia:                optional.New("iso_a4_210x297mm"),
			InputOrientationRequested: optional.New(EnInputOrientationPortrait),
			InputQuality:              optional.New(EnInputQualityNormal),
			InputResolution:           optional.New(goipp.Resolution{Xres: 300, Yres: 300, Units: goipp.UnitsDpi}),
			InputScalingHeight:        optional.New(100),
			InputScalingWidth:         optional.New(100),
			InputScanRegions: []InputScanRegion{
				{
					XDimension: optional.New(21000),
					XOrigin:    optional.New(0),
					YDimension: optional.New(29700),
					YOrigin:    optional.New(0),
				},
			},
			InputSharpness: optional.New(0),
			InputSides:     optional.New(KwSidesOneSided),
			InputSource:    optional.New(KwInputSourcePlaten),
		}),
		OutputAttributes: optional.New(OutputAttributes{
			NoiseRemoval:                   optional.New(20),
			OutputCompressionQualityFactor: optional.New(75),
		}),
	}
}

// Test_Scanner_HappyPath verifies the full scan workflow:
// Get-Printer-Attributes → Create-Job → Get-Next-Document-Data.
func Test_Scanner_HappyPath(t *testing.T) {
	vscan := &abstract.VirtualScanner{
		ScanCaps: testScannerCapabilities(),
		Resolution: abstract.Resolution{
			XResolution: 300,
			YResolution: 300,
		},
		PlatenImage: testutils.Images.JPEG100x75rgb8,
	}

	scanner := testNewScanner(t, vscan)
	srv := httptest.NewServer(scanner)
	defer srv.Close()

	httpURL, ippURI := testScannerURL(srv)
	client := NewClient(httpURL, nil)
	ctx := context.Background()

	// Step 1: Get-Printer-Attributes
	gpaRq := &GetPrinterAttributesRequest{
		RequestHeader: DefaultRequestHeader,
		PrinterURI:    ippURI,
		RequestedAttributes: []string{
			"input-color-mode-supported",
			"input-resolution-supported",
			"input-source-supported",
		},
	}
	gpaRsp := &GetPrinterAttributesResponse{}
	if err := client.Do(ctx, gpaRq, gpaRsp); err != nil {
		t.Fatalf("Get-Printer-Attributes: %v", err)
	}
	if gpaRsp.Status != goipp.StatusOk {
		t.Fatalf("Get-Printer-Attributes status: %s", gpaRsp.Status)
	}

	sd := gpaRsp.Printer.ScannerDescription
	if !slices.Contains(sd.InputSourceSupported, KwInputSourcePlaten) {
		t.Errorf("input-source-supported: got %v, want platen",
			sd.InputSourceSupported)
	}
	if !slices.Contains(sd.InputColorModeSupported, KwInputColorModeColor8) {
		t.Errorf("input-color-mode-supported: got %v, want color_8",
			sd.InputColorModeSupported)
	}
	wantRes := goipp.Resolution{Xres: 300, Yres: 300, Units: goipp.UnitsDpi}
	if !slices.Contains(sd.InputResolutionSupported, wantRes) {
		t.Errorf("input-resolution-supported: got %v, want %v",
			sd.InputResolutionSupported, wantRes)
	}

	// Step 2: Create-Job
	createRq := &CreateJobRequest{
		RequestHeader:      DefaultRequestHeader,
		JobCreateOperation: testScanJobCreateOperation(ippURI),
		JobTemplate:        &JobTemplate{},
	}
	createRsp := &CreateJobResponse{}
	if err := client.Do(ctx, createRq, createRsp); err != nil {
		t.Fatalf("Create-Job: %v", err)
	}
	if createRsp.Status != goipp.StatusOk {
		t.Fatalf("Create-Job status: %s", createRsp.Status)
	}
	if createRsp.Job == nil {
		t.Fatal("Create-Job: missing job status")
	}
	if createRsp.Job.JobState != EnJobStateProcessing {
		t.Errorf("job-state: got %v, want processing", createRsp.Job.JobState)
	}
	if createRsp.Job.JobID == 0 {
		t.Fatal("Create-Job: missing job-id")
	}

	jobID := createRsp.Job.JobID

	// Step 3: Get-Next-Document-Data (first page)
	docRsp, err := testGetNextDocumentData(ctx, client, ippURI, jobID)
	if err != nil {
		t.Fatalf("Get-Next-Document-Data: %v", err)
	}
	if docRsp.LastDocument {
		t.Fatal("expected scan page, got last-document")
	}
	if optional.Get(docRsp.DocumentFormat) != "image/jpeg" {
		t.Errorf("document-format: got %q, want image/jpeg",
			optional.Get(docRsp.DocumentFormat))
	}
	if optional.Get(docRsp.Document.DocumentNumber) != 1 {
		t.Errorf("document-number: got %d, want 1",
			optional.Get(docRsp.Document.DocumentNumber))
	}

	body := docRsp.Header().Body
	if body == nil {
		t.Fatal("Get-Next-Document-Data: missing document body")
	}
	gotPage, err := io.ReadAll(body)
	body.Close()
	if err != nil {
		t.Fatalf("read document body: %v", err)
	}
	if len(gotPage) == 0 {
		t.Fatal("document body is empty")
	}

	// Step 4: Get-Next-Document-Data (end of scan)
	endRsp, err := testGetNextDocumentData(ctx, client, ippURI, jobID)
	if err != nil {
		t.Fatalf("Get-Next-Document-Data (end): %v", err)
	}
	if !endRsp.LastDocument {
		t.Fatal("expected last-document after final page")
	}
	if endRsp.Header().Body != nil {
		rest, err := io.ReadAll(endRsp.Header().Body)
		endRsp.Header().Body.Close()
		if err != nil {
			t.Fatalf("read last-document body: %v", err)
		}
		if len(rest) != 0 {
			t.Fatalf("last-document response must not include body, got %d bytes", len(rest))
		}
	}
}
