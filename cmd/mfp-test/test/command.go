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

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
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

	// Tasks 12-14 will wire virtual printer, CUPS queue,
	// and test execution here.
	_ = model
	<-ctx.Done()
	return nil
}
