// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// mfp-test command definition

package test

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// DefaultTCPPort is the default IPP server TCP port.
const DefaultTCPPort = 60000

// DefaultQueueName is the default CUPS queue name.
const DefaultQueueName = "mfp-test"

// Command is the mfp-test command description.
var Command = argv.Command{
	Name: "mfp-test",
	Help: "Print system testing pipeline",
	Options: []argv.Option{
		{
			Name:      "-m",
			Aliases:   []string{"--model"},
			Help:      "printer model file",
			HelpArg:   "file",
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		{
			Name:      "-P",
			Aliases:   []string{"--port"},
			Help:      fmt.Sprintf("IPP server TCP port (default %d)", DefaultTCPPort),
			HelpArg:   "port",
			Singleton: true,
			Validate:  argv.ValidateUint16,
		},
		{
			Name:      "-n",
			Aliases:   []string{"--name"},
			Help:      fmt.Sprintf("CUPS queue name (default %q)", DefaultQueueName),
			HelpArg:   "name",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name:      "-o",
			Aliases:   []string{"--output"},
			Help:      "write JSON report to file",
			HelpArg:   "file",
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		{
			Name:      "--threshold",
			Help:      "minimum similarity score to pass (0.0-1.0, default 0.95)",
			HelpArg:   "score",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name: "--list",
			Help: "list all test configurations and exit",
		},
		{
			Name:      "--batch",
			Help:      "run all test configurations",
			Singleton: true,
		},
		{
			Name:      "--single",
			Help:      "run a single test configuration by name",
			HelpArg:   "name",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name:    "-v",
			Aliases: []string{"--verbose"},
			Help:    "enable verbose output",
		},
		argv.HelpOption,
	},
	Handler: cmdTestHandler,
}

// cmdTestHandler is the top-level handler for the mfp-test command.
func cmdTestHandler(ctx context.Context, inv *argv.Invocation) error {
	level := log.LevelInfo
	if inv.Flag("-v") {
		level = log.LevelTrace
	}
	logger := log.NewLogger(level, log.Console)
	ctx = log.NewContext(ctx, logger)

	// Model file is required: without it, NewIPPServer() returns nil.
	modelfile, ok := inv.Get("-m")
	if !ok {
		return fmt.Errorf("model file required: use -m <file>")
	}

	model, err := modeling.NewModel()
	if err != nil {
		return err
	}
	defer model.Close()

	if err := model.Load(modelfile); err != nil {
		return fmt.Errorf("load model %q: %w", modelfile, err)
	}

	// Create document capture backend
	capture := NewDocumentCapture()

	// Create virtual IPP printer from model and hook capture into it
	ippPrinter := model.NewIPPServer()
	if ippPrinter == nil {
		return fmt.Errorf("model has no IPP printer attributes")
	}
	ippPrinter.SetPrintBackend(capture)

	// Create in-process loopback — no real TCP socket needed
	tr, loopback := transport.NewLoopback()

	mux := transport.NewPathMux()
	mux.Add("/ipp/print", ippPrinter)

	srvr := transport.NewServer(ctx, nil, mux)
	log.Info(ctx, "virtual IPP printer started (in-process loopback)")
	go srvr.Serve(loopback)
	defer srvr.Close()

	// Send test document via in-process IPP client
	log.Info(ctx, "sending test document...")

	printerURL, _ := url.Parse("ipp://loopback/ipp/print")
	ippURI := "ipp://loopback/ipp/print"
	client := ipp.NewClient(printerURL, tr)

	// Step 1: Create-Job
	createRq := &ipp.CreateJobRequest{
		RequestHeader: ipp.DefaultRequestHeader,
		JobCreateOperation: ipp.JobCreateOperation{
			PrinterURI: ippURI,
		},
		Job: &ipp.JobAttributes{},
	}
	createRsp := &ipp.CreateJobResponse{}
	if err := client.Do(ctx, createRq, createRsp); err != nil {
		return fmt.Errorf("Create-Job: %w", err)
	}

	// Step 2: Send-Document
	sendRq := &ipp.SendDocumentRequest{
		RequestHeader:  ipp.DefaultRequestHeader,
		PrinterURI:     optional.New(ippURI),
		JobID:          optional.New(createRsp.Job.JobID),
		DocumentFormat: optional.New("text/plain"),
		LastDocument:   true,
		Job:            &ipp.JobAttributes{},
	}
	sendRq.Body = bytes.NewReader([]byte("mfp-test sanity check\n"))

	sendRsp := &ipp.SendDocumentResponse{}
	if err := client.Do(ctx, sendRq, sendRsp); err != nil {
		return fmt.Errorf("Send-Document: %w", err)
	}

	// Wait for the document to arrive at capture backend
	select {
	case <-capture.OnDocument():
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout: no document received after 30s")
	case <-ctx.Done():
		return nil
	}

	// Report captured result
	docs := capture.Docs()
	for i, d := range docs {
		log.Info(ctx, "captured doc %d: %d bytes, format=%q, job=%q",
			i+1, len(d.Data), d.Params.Format, d.Params.JobName)
	}

	<-ctx.Done()
	return nil
}
