// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "detect-printers" command.

package cups

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// cmdDetectPrinters defines the "detect-printers" sub-command.
var cmdDetectPrinters = argv.Command{
	Name:    "detect-printers",
	Help:    "Search for available devices",
	Handler: cmdDetectPrintersHandler,
	Options: []argv.Option{
		optSchemesExclude,
		optSchemesInclude,
		optLimit,
		optTimeout,
		argv.HelpOption,
	},
}

// cmdDetectPrintersHandler is the "detect-printers" command handler
func cmdDetectPrintersHandler(ctx context.Context, inv *argv.Invocation) error {
	// Prepare arguments
	sel := &cups.GetDevicesSelection{
		Limit:          optLimitGet(inv),
		Timeout:        optTimeoutGet(inv),
		ExcludeSchemes: optSchemesExcludeGet(inv),
		IncludeSchemes: optSchemesIncludeGet(inv),
	}

	// Perform the query
	clnt := cups.NewClient(optCUPSURL, nil)
	clnt.SetDecoderOptions(&ipp.DecoderOptions{KeepTrying: true})

	devices, err := clnt.CUPSGetDevices(ctx, sel, []string{"all"})
	if err != nil {
		return err
	}

	// Format output
	for _, dev := range devices {
		switch {
		// Skip devices with missed or incomplete URIs
		case dev.DeviceURI == nil:
		case !strings.Contains(*dev.DeviceURI, "://"):

		// Dump others
		default:
			fmt.Println("")
			devAttrsFormat(os.Stdout, dev)
		}
	}

	return nil
}
