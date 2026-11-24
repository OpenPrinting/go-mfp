// MFP - Multi-Function Printers and scanners toolkit
// The "mfp-cups" command
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions
//
// The "disable-printer" command.

package command

import (
	"context"
	"fmt"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
)

// cmdDisablePrinter defines the "disable-printer" sub-command.
var cmdDisablePrinter = argv.Command{
	Name:    "disable-printer",
	Help:    "Make a print queue reject jobs",
	Handler: cmdDisablePrinterHandler,
	Options: []argv.Option{
		optPrinterURI,
		argv.HelpOption,
	},
}

// cmdDisablePrinterHandler is the "disable-printer" command handler.
func cmdDisablePrinterHandler(ctx context.Context, inv *argv.Invocation) error {
	printerURI := optPrinterURIGet(inv)
	if printerURI == "" {
		return fmt.Errorf("%s option required", optPrinterURI.Name)
	}

	clnt := cups.NewClient(optCUPSURL, nil)
	return clnt.DisablePrinter(ctx, printerURI)
}
