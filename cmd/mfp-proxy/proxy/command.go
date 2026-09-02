// MFP - Miulti-Function Printers and scanners toolkit
// The "proxy" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Command description.

package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/internal/env"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/log/trace"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/transport"
)

// DefaultTCPPort is the default TCP port for the MFP proxy
const DefaultTCPPort = 50000

// description is printed as a command description text
const description = "" +
	"This command runs the IPP/eSCL/WSD proxy\n" +
	"The proxy can be useful for two purposes:\n" +
	"  - to sniff the protocol traffic between the IPP/eSCL/WSD\n" +
	"    client and server\n" +
	"  - to logically bring the device into the different IP address\n" +
	"    or port\n" +
	"\n" +
	"If optional command is specified, the CUPS_SERVER and the\n" +
	"SANE_AIRSCAN_DEVICE environment variables will be set properly\n" +
	"and the command will be executed, The proxy  will exit when\n" +
	"the command finished.\n" +
	"\n" +
	"Without that the proxy  will run until termination signal\n" +
	"is received.\n"

// Command is the 'proxy' command description
var Command = argv.Command{
	Name:                     "proxy",
	Help:                     "IPP/eSCL/WSD masquerading proxy",
	Description:              description,
	NoOptionsAfterParameters: true,
	Options: []argv.Option{
		argv.Option{
			Name:    "-P",
			Aliases: []string{"--port"},
			Help: fmt.Sprintf("TCP port. Default: %d",
				DefaultTCPPort),
			HelpArg:   "port",
			Singleton: true,
			Validate:  argv.ValidateUint16,
		},
		argv.Option{
			Name:      "-U",
			Aliases:   []string{"--usbip"},
			Help:      "USBIP mode",
			Singleton: true,
			Conflicts: []string{"-P"},
		},
		argv.Option{
			Name:     "-E",
			Aliases:  []string{"--escl"},
			Help:     "Forward eSCL requests from local path to url",
			HelpArg:  "path=url",
			Validate: validateMapping,
		},
		argv.Option{
			Name:     "-I",
			Aliases:  []string{"--ipp"},
			Help:     "Forward IPP requests from local path to url",
			HelpArg:  "path=url",
			Validate: validateMapping,
		},
		argv.Option{
			Name:     "-W",
			Aliases:  []string{"--wsd"},
			Help:     "Forward WSD requests from local path to url",
			HelpArg:  "path=url",
			Validate: validateMapping,
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
			Help: "command to run under the proxy",
		},
		{
			Name: "[args...]",
			Help: "the command's arguments",
		},
	},
	Handler: cmdProxyHandler,
}

// cmdProxyHandler is the top-level handler for the 'proxy' command.
func cmdProxyHandler(ctx context.Context, inv *argv.Invocation) error {
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

	// Setup trace
	if traceName, _ := inv.Get("-t"); traceName != "" {
		tracer, err := trace.NewWriter(ctx, traceName)
		if err != nil {
			return err
		}

		defer tracer.Close()
		ctx = trace.NewContext(ctx, tracer)
	}

	// Validate parameters
	portnum := DefaultTCPPort
	if portname, ok := inv.Get("-P"); ok {
		var err error
		portnum, err = strconv.Atoi(portname)
		assert.NoError(err)
	}

	// Parse mappings
	var mappings []mapping

	for _, opt := range inv.Values("--escl") {
		m, err := parseMapping(protoESCL, opt)
		if err != nil {
			return err
		}
		mappings = append(mappings, m)
	}

	for _, opt := range inv.Values("--ipp") {
		m, err := parseMapping(protoIPP, opt)
		if err != nil {
			return err
		}
		mappings = append(mappings, m)
	}

	for _, opt := range inv.Values("--wsd") {
		m, err := parseMapping(protoWSD, opt)
		if err != nil {
			return err
		}
		mappings = append(mappings, m)
	}

	if len(mappings) == 0 {
		err := errors.New("at least one option required: --escl, --ipp or --wsd")
		return err
	}

	// Create and populate the PathMux
	runner := env.Runner{
		ESCLName: "Virtual MFP Scanner",
	}

	mux := transport.NewPathMux()
	for _, m := range mappings {
		if mux.Contains(m.localPath) {
			err := fmt.Errorf("Local path %q used multiple times",
				m.localPath)
			return err
		}

		switch m.proto {
		case protoIPP:
			proxy := ipp.NewProxy(m.localPath, m.targetURL)
			mux.Add(m.localPath, proxy)

			runner.CUPSPort = portnum

		case protoESCL:
			proxy := escl.NewProxy(m.localPath, m.targetURL)
			mux.Add(m.localPath, proxy)

			runner.ESCLPort = portnum
			runner.ESCLPath = m.localPath

		case protoWSD:
			return errors.New("WSD proxy not implemented")
		}
	}

	// Create server for incoming connections.
	if !inv.Flag("-U") {
		l, err := newListener(ctx, portnum)
		if err != nil {
			return err
		}

		srvr := transport.NewServer(ctx, nil, mux)
		log.Info(ctx, "starting MFP proxy at http://localhost:%d",
			portnum)
		go srvr.Serve(l)

		defer srvr.Close()
	} else {
		addr := &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 3240,
		}

		log.Info(ctx, "starting USBIP server at %s", addr)
		newUsbipServer(ctx, addr, mux)
	}

	// Run external program if requested
	if command, ok := inv.Get("command"); ok {
		argv := inv.Values("args")
		return runner.Run(ctx, command, argv...)
	}

	// Wait for termination signal
	<-ctx.Done()
	log.Info(ctx, "Exiting...")

	return nil
}
