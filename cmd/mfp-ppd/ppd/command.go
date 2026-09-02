// MFP - Miulti-Function Printers and scanners toolkit
// The "ppd" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Command description

package ppd

import (
	"context"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
)

// Command is the 'ppd' command description
var Command = argv.Command{
	Name: "ppd",
	Help: "Utility for PPD files",
	Options: []argv.Option{
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
	SubCommands: []argv.Command{
		cmdDump,
		argv.HelpCommand,
	},
	Handler: cmdPpdHandler,
}

// cmdPpdHandler is the top-level handler for the 'ppd' command.
func cmdPpdHandler(ctx context.Context, inv *argv.Invocation) error {
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

	// Execute subcommand
	return argv.DefaultHandler(ctx, inv)
}
