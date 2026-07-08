// MFP - Miulti-Function Printers and scanners toolkit
// The "model" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Command description.

package model

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/discovery/dnssd"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/transport"
)

// DefaultTCPPort is the default TCP port for the MFP simulator
const DefaultTCPPort = 50000

// description is printed as a command description text
const description = "" +
	"This command generates and validates models for the MFP simulator\n" +
	""

// Command is the 'model' command description
var Command = argv.Command{
	Name:                     "model",
	Help:                     "Model generator for MFP simulator",
	Description:              description,
	NoOptionsAfterParameters: true,
	Options: []argv.Option{
		argv.Option{
			Name:      "-D",
			Aliases:   []string{"--dnssd"},
			Help:      "DNS-SD name of the device",
			HelpArg:   "name",
			Conflicts: []string{"-E", "-I", "-W"},
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  dnssd.ArgvCompleter,
		},
		argv.Option{
			Name:      "-E",
			Aliases:   []string{"--escl"},
			Help:      "eSCL scanner URL",
			HelpArg:   "URL",
			Singleton: true,
			Validate:  transport.ValidateAddr,
		},
		argv.Option{
			Name:      "-I",
			Aliases:   []string{"--ipp"},
			Help:      "IPP printer URL",
			HelpArg:   "URL",
			Singleton: true,
			Validate: func(s string) error {
				err := transport.ValidateAddr(s)
				return err
			},
		},
		argv.Option{
			Name:      "-W",
			Aliases:   []string{"--wsd"},
			Help:      "WSD scanner URL",
			HelpArg:   "URL",
			Singleton: true,
			Validate: func(s string) error {
				err := transport.ValidateAddr(s)
				return err
			},
		},
		argv.Option{
			Name:      "-m",
			Aliases:   []string{"--model"},
			Help:      "write model to file (use - for stdout)",
			HelpArg:   "file",
			Required:  true,
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		argv.Option{
			Name:      "-V",
			Aliases:   []string{"--validate"},
			Help:      "validate existent model",
			Singleton: true,
			Conflicts: []string{
				"--dnssd", "--escl", "--ipp", "--wsd",
			},
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
		argv.HelpOption,
	},
	Handler: cmdModelHandler,
}

// cmdModelHandler is the top-level handler for the 'model' command.
func cmdModelHandler(ctx context.Context, inv *argv.Invocation) error {
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

	_ = ctx

	// Check options
	optDHSSD, haveDNSSD := inv.Get("--dnssd")
	_, validate := inv.Get("--validate")
	optESCL := inv.Values("--escl")
	optIPP := inv.Values("--ipp")
	optWSD := inv.Values("--wsd")

	if !haveDNSSD && !validate &&
		optIPP == nil && optESCL == nil && optWSD == nil {

		err := errors.New("at least one option required: --dnssd, --escl, --ipp, --wsd or --validate")
		return err
	}

	// Handle the --validate option
	if validate {
		model, err := modeling.NewModel()
		if err != nil {
			return err
		}

		defer model.Close()

		file, _ := inv.Get("-m")
		err = model.Load(file)
		if err != nil {
			return err
		}

		caps := model.GetIPPPrinterAttrs()
		if caps != nil {
			errors := caps.Errors()
			if errors != nil {
				log.Warning(ctx, "ipp: printer attributes decoded with warnings:")
				for _, err := range errors {
					log.Warning(ctx, "  %s", err)
				}
			}
		}

		return err
	}

	// Gather endpoints
	if haveDNSSD {
		dev, err := discoverByName(ctx, optDHSSD)
		if err != nil {
			return err
		}

		for _, unit := range dev.PrintUnits {
			if unit.Proto == discovery.ServiceIPP {
				optIPP = append(optIPP, unit.Endpoints...)
			}
		}

		for _, unit := range dev.ScanUnits {
			switch unit.Proto {
			case discovery.ServiceESCL:
				optESCL = append(optESCL, unit.Endpoints...)
			case discovery.ServiceWSD:
				optWSD = append(optWSD, unit.Endpoints...)
			}
		}

		// Check that something was discovered.
		if optIPP == nil && optESCL == nil && optWSD == nil {
			err := errors.New("no eSCL/IPP/WSD endpoints discovered")
			return err
		}
	}

	// Create a model
	model, err := modeling.NewModel()
	if err != nil {
		return err
	}

	// Query printer and scanner capabilities
	if optIPP != nil {
		err = model.DownloadIPPPrinterAttrs(ctx, optIPP)
		if err != nil {
			err = fmt.Errorf(
				"Can't get IPP Printer Attributes: %s", err)
			return err
		}

		attrs := model.GetIPPPrinterAttrs()
		if errors := attrs.Errors(); errors != nil {
			log.Warning(ctx, "ipp: printer attributes decoded with warnings:")
			for _, err := range errors {
				log.Warning(ctx, "  %s", err)
			}
		}

	}

	if optESCL != nil {
		err = model.DownloadESCLScannerCapabilities(ctx, optESCL)
		if err != nil {
			err = fmt.Errorf(
				"Can't get eSCL ScannerCapabilities: %s", err)
			return err
		}
	}

	if optWSD != nil {
		err = model.DownloadWSDScannerCapabilities(ctx, optWSD)
		if err != nil {
			err = fmt.Errorf(
				"Can't get WSD ScannerCapabilities: %s", err)
			return err
		}
	}

	// Save model to file
	file, _ := inv.Get("-m")
	if file == "-" {
		return model.Write(os.Stdout)
	}

	return model.Save(file)
}
