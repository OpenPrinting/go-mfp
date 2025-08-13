// MFP - Miulti-Function Printers and scanners toolkit
// The "virtual" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Command description.

package virtual

import (
	"context"
	"fmt"
	"strconv"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/proto/escl"
)

// DefaultTCPPort is the default TCP port for the MFP simulator
const DefaultTCPPort = 50000

// description is printed as a command description text
const description = "" +
	"This command runs the MFP simulator\n" +
	"\n" +
	"If optional command is specified, the CUPS_SERVER and the\n" +
	"SANE_AIRSCAN_DEVICE environment variables will be set properly\n" +
	"and the command will be executed, The simulator will exit when\n" +
	"the command finished.\n" +
	"\n" +
	"Without that the simulator will run until termination signal\n" +
	"is received.\n"

// Command is the 'virtual' command description
var Command = argv.Command{
	Name:                     "virtual",
	Help:                     "Virtual MFP simulator",
	Description:              description,
	NoOptionsAfterParameters: true,
	Options: []argv.Option{
		argv.Option{
			Name:      "-m",
			Aliases:   []string{"--model"},
			Help:      "read model from file",
			HelpArg:   "file",
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		argv.Option{
			Name:    "-d",
			Aliases: []string{"--debug"},
			Help:    "Enable debug output",
		},
		argv.Option{
			Name:    "-v",
			Aliases: []string{"--verbose"},
			Help:    "Enable verbose debug output",
		},
		argv.Option{
			Name:    "-p",
			Aliases: []string{"--port"},
			HelpArg: "port",
			Help: fmt.Sprintf("TCP port. Default: %d",
				DefaultTCPPort),
			Validate: argv.ValidateUint16,
		},
		argv.HelpOption,
	},
	Parameters: []argv.Parameter{
		{
			Name: "[command]",
			Help: "command to run under the simulator",
		},
		{
			Name: "[args...]",
			Help: "the command's arguments",
		},
	},
	Handler: cmdVirtualHandler,
}

// cmdVirtualHandler is the top-level handler for the 'cups' command.
func cmdVirtualHandler(ctx context.Context, inv *argv.Invocation) error {
	// Setup logging
	_, dbg := inv.Get("-d")
	_, vrb := inv.Get("-v")

	level := log.LevelInfo
	if dbg {
		level = log.LevelDebug
	}
	if vrb {
		level = log.LevelTrace
	}

	logger := log.NewLogger(level, log.Console)
	ctx = log.NewContext(ctx, logger)

	var err error

	// Create MFP model
	model, err := modeling.NewModel()
	if err != nil {
		return err
	}

	defer model.Close()

	// Load model file, or use defaults
	if modelfile, ok := inv.Get("-m"); ok {
		err = model.Load(modelfile)
		if err != nil {
			return err
		}
	} else {
		caps := modeling.DefaultScannerCapabilities()
		esclcaps := escl.FromAbstractScannerCapabilities(
			escl.DefaultVersion, caps)
		model.SetESCLScanCaps(esclcaps)
	}

	// Obtain remaining parameters
	port := DefaultTCPPort
	if portname, ok := inv.Get("-p"); ok {
		port, err = strconv.Atoi(portname)
		if err != nil {
			return err
		}
	}

	argv := []string{}
	if command, ok := inv.Get("command"); ok {
		argv = append(argv, command)
		argv = append(argv, inv.Values("args")...)
	}

	// Run the simulator
	return simulate(ctx, model, port, argv)
}
