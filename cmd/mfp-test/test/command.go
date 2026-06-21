// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// mfp-test command definition

package test

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/transport"
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

	// Parse port number
	port := DefaultTCPPort
	if portStr, ok := inv.Get("-P"); ok {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid port %q: %w", portStr, err)
		}
		port = p
	}

	// Create document capture backend
	capture := NewDocumentCapture()

	// Create virtual IPP printer from model and hook capture into it
	ippPrinter := model.NewIPPServer()
	if ippPrinter == nil {
		return fmt.Errorf("model has no IPP printer attributes")
	}
	ippPrinter.SetPrintBackend(capture)

	// Register IPP handler on the URL path /ipp/print
	mux := transport.NewPathMux()
	mux.Add("/ipp/print", ippPrinter)

	// Open TCP port and start HTTP server (IPP runs over HTTP)
	addr := fmt.Sprintf("localhost:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srvr := transport.NewServer(ctx, nil, mux)
	log.Info(ctx, "virtual IPP printer at ipp://%s/ipp/print", addr)
	go srvr.Serve(ln)
	defer srvr.Close()

	// Get CUPS queue name
	queueName := DefaultQueueName
	if name, ok := inv.Get("-n"); ok {
		queueName = name
	}

	// Register virtual printer with CUPS
	ippURL := fmt.Sprintf("ipp://localhost:%d/ipp/print", port)
	if err := CreateCUPSQueue(ctx, queueName, ippURL); err != nil {
		return err
	}
	// WithoutCancel preserves logging and values from ctx but
	// prevents cancellation from stopping the cleanup operation.
	defer RemoveCUPSQueue(context.WithoutCancel(ctx), queueName)

	log.Info(ctx, "CUPS queue %q ready at %s", queueName, ippURL)

	// Send a minimal test document through the full pipeline
	log.Info(ctx, "sending test document via lp...")
	lpCmd := exec.CommandContext(ctx, "lp", "-d", queueName)
	lpCmd.Stdin = strings.NewReader("mfp-test sanity check\n")
	if out, err := lpCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("lp -d %s: %w: %s", queueName, err, out)
	}

	// Wait for the document to arrive at the capture backend
	select {
	case <-capture.OnDocument():
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout: no document received after 30s")
	case <-ctx.Done():
		return nil
	}

	// Report what was captured
	docs := capture.Docs()
	for i, d := range docs {
		log.Info(ctx, "captured doc %d: %d bytes, format=%q, job=%q",
			i+1, len(d.Data), d.Params.Format, d.Params.JobName)
	}

	<-ctx.Done()
	return nil
}
