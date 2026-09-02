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
	"image"
	"image/color"
	"image/png"
	"net"
	"os"
	"strconv"
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
			Required:  true,
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
			Help:      fmt.Sprintf("minimum similarity score to pass (0.0-1.0, default %.2f)", DefaultThreshold),
			HelpArg:   "score",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name:      "--timeout",
			Help:      fmt.Sprintf("document capture timeout (default %s)", DefaultTimeout),
			HelpArg:   "duration",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name: "--list",
			Help: "list all test configurations (full batch matrix) and exit",
		},
		{
			Name: "--list-quick",
			Help: "list quick test configurations and exit",
		},
		{
			Name:      "--batch",
			Help:      "run all test configurations",
			Singleton: true,
		},
		{
			Name:      "--single",
			Help:      "run a single test configuration by name (sides/color-mode/format)",
			HelpArg:   "name",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name: "--quick",
			Help: "run reduced test matrix (all sides × all color modes, first format only)",
		},
		{
			Name: "--keep",
			Help: "save captured document bytes to current directory",
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

	// Query printer capabilities for test matrix generation.
	caps, err := QueryPrinterCaps(ctx, ippURL)
	if err != nil {
		return fmt.Errorf("query printer capabilities: %w", err)
	}

	// --list: print all configurations (full batch matrix) and exit.
	if inv.Flag("--list") {
		for _, cfg := range BatchMatrix(caps) {
			fmt.Println(cfg.Name)
		}
		return nil
	}

	// --list-quick: print quick test configurations and exit.
	if inv.Flag("--list-quick") {
		for _, cfg := range QuickMatrix(caps) {
			fmt.Println(cfg.Name)
		}
		return nil
	}

	// Determine which test configurations to run.
	var configs []TestConfig
	switch {
	case inv.Flag("--batch"):
		configs = BatchMatrix(caps)
	case inv.Flag("--quick"):
		configs = QuickMatrix(caps)
	default:
		if spec, ok := inv.Get("--single"); ok {
			cfg, err := SingleConfig(spec)
			if err != nil {
				return err
			}
			configs = []TestConfig{*cfg}
		} else {
			// Default: single run with printer defaults.
			configs = QuickMatrix(caps)
		}
	}

	// Parse similarity threshold.
	threshold := DefaultThreshold
	if ts, ok := inv.Get("--threshold"); ok {
		var t float64
		if _, err := fmt.Sscanf(ts, "%f", &t); err != nil {
			return fmt.Errorf("invalid threshold %q: %w", ts, err)
		}
		threshold = t
	}

	// Parse capture timeout.
	timeout := DefaultTimeout
	if ts, ok := inv.Get("--timeout"); ok {
		d, err := time.ParseDuration(ts)
		if err != nil {
			return fmt.Errorf("invalid timeout %q: %w", ts, err)
		}
		timeout = d
	}

	keep := inv.Flag("--keep")
	verbose := inv.Flag("-v")

	// Run each test configuration.
	for _, cfg := range configs {
		log.Info(ctx, "running test: %s", cfg.Name)
		result, err := RunTest(ctx, cfg, queueName, capture, threshold, timeout, keep, verbose)
		if err != nil {
			log.Info(ctx, "FAIL %s: %v", cfg.Name, err)
			continue
		}
		if result.Passed {
			log.Info(ctx, "PASS %s (score=%.4f)", cfg.Name, result.Score)
		} else {
			log.Info(ctx, "FAIL %s (score=%.4f < threshold=%.4f)", cfg.Name, result.Score, threshold)
		}
	}

	return nil
}

// generateTestPNG creates a temporary PNG test image with three
// horizontal colour bands (red, green, blue) and returns its path.
// The caller is responsible for removing the file after use.
func generateTestPNG() (string, error) {
	const size = 300
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			switch {
			case y < size/3:
				img.Set(x, y, color.RGBA{R: 255, A: 255}) // red
			case y < 2*size/3:
				img.Set(x, y, color.RGBA{G: 255, A: 255}) // green
			default:
				img.Set(x, y, color.RGBA{B: 255, A: 255}) // blue
			}
		}
	}

	f, err := os.CreateTemp("", "mfp-test-*.png")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		os.Remove(f.Name())
		return "", err
	}

	return f.Name(), nil
}
