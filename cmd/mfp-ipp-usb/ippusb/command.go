// MFP - Miulti-Function Printers and scanners toolkit
// The "mfp-ipp-usb" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Command description.

package ippusb

import (
	"context"
	"fmt"
	"os"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/ippusbapi"
)

// Command is the 'mfp-ipp-usb' command description
var Command = argv.Command{
	Name: "mfp-ipp-usb",
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
	Description: "" +
		"This is the query tool for ipp-usb daemon API\n" +
		"See ipp-usb-api(7) for details",
	Handler: cmdIPPUSBHandler,

	SubCommands: []argv.Command{
		argv.Command{
			Name:    "status",
			Help:    "Print ipp-usb status",
			Handler: cmdStatusHandler,
		},
		argv.Command{
			Name:    "version",
			Help:    "Print ipp-usb API version",
			Handler: cmdVersionHandler,
		},
	},
}

// cmdIPPUSBHandler is the top-level handler for the 'mfp-ipp-usb' command.
func cmdIPPUSBHandler(ctx context.Context, inv *argv.Invocation) error {
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

// cmdStatusHandler is the "status" sub-command handler
func cmdStatusHandler(ctx context.Context, inv *argv.Invocation) error {
	clnt := ippusbapi.NewClient(ctx)
	status, err := clnt.Status()
	if err != nil {
		return err
	}

	os.Stdout.Write(status)
	return nil
}

// cmdVersionHandler is the "version" sub-command handler
func cmdVersionHandler(ctx context.Context, inv *argv.Invocation) error {
	clnt := ippusbapi.NewClient(ctx)
	version, err := clnt.Version()
	if err != nil {
		return err
	}

	fmt.Println(version)
	return nil
}
