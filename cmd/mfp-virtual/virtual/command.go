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
	"strings"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/log/trace"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/transport"
)

// DefaultTCPPortMin and DefaultTCPPortMax define the
// default TCP port range for the MFP simulator
const (
	DefaultTCPPortMin = 50000
	DefaultTCPPortMax = 59999
)

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
			Name:    "-P",
			Aliases: []string{"--port"},
			HelpArg: "N or MIN-MAX",
			Help: fmt.Sprintf("TCP port range. Default: %d-%d",
				DefaultTCPPortMin, DefaultTCPPortMax),
			Validate: transport.ArgvValidatePortRange,
		},
		argv.Option{
			Name:      "-U",
			Aliases:   []string{"--usbip"},
			Help:      "USBIP mode",
			Singleton: true,
			Conflicts: []string{"-P"},
		},
		argv.Option{
			Name:      "-m",
			Aliases:   []string{"--model"},
			Help:      "read model from file",
			HelpArg:   "file",
			Required:  true,
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		argv.Option{
			Name:     "-t",
			Aliases:  []string{"--trace"},
			Help:     "write trace to file.log and file.tar",
			HelpArg:  "file",
			Validate: argv.ValidateAny,
			Complete: argv.CompleteOSPath,
		},
		argv.Option{
			Name:    "-d",
			Aliases: []string{"--debug"},
			Help:    "Enable debug output",
		},
		argv.Option{
			Name:    "-v",
			Aliases: []string{"--verbose"},
			Help:    "Verbose logging (-vv for very verbose)",
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
	level := log.LevelInfo
	switch {
	case len(inv.Values("-v")) > 1:
		level = log.LevelTrace
	case len(inv.Values("-v")) > 0:
		level = log.LevelVerbose
	case inv.Flag("-d"):
		level = log.LevelDebug
	}

	logger := log.NewLogger(level, log.Console)
	ctx = log.NewContext(ctx, logger)

	var err error

	// Setup tracer
	if traceName, _ := inv.Get("-t"); traceName != "" {
		tracer, err := trace.NewWriter(ctx, traceName)
		if err != nil {
			return err
		}

		defer tracer.Close()
		ctx = trace.NewContext(ctx, tracer)
	}

	// Create MFP model
	model, err := modeling.NewModel()
	if err != nil {
		return err
	}

	defer model.Close()

	// Load model file
	modelfile, _ := inv.Get("-m")
	err = model.Load(modelfile)
	if err != nil {
		return err
	}

	// Obtain remaining parameters
	portmin := DefaultTCPPortMin
	portmax := DefaultTCPPortMax
	if portrange, ok := inv.Get("-P"); ok {
		if n := strings.IndexByte(portrange, '-'); n < 0 {
			portmin, _ = strconv.Atoi(portrange)
			portmax = portmin
		} else {
			min := portrange[:n]
			max := portrange[n+1:]

			portmin, _ = strconv.Atoi(min)
			portmax, _ = strconv.Atoi(max)
		}
	}

	argv := []string{}
	if command, ok := inv.Get("command"); ok {
		argv = append(argv, command)
		argv = append(argv, inv.Values("args")...)
	}

	// Run the simulator
	usbip := inv.Flag("-U")
	return simulate(ctx, model,
		uint16(portmin), uint16(portmax), usbip, argv)
}
