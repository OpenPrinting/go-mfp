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

// cmdEnablePrinter defines the "enable-printer" sub-command.
var cmdEnablePrinter = argv.Command{
	Name:    "enable-printer",
	Help:    "Make a print queue accept jobs",
	Handler: cmdEnablePrinterHandler,
	Options: []argv.Option{
		optPrinter,
		argv.HelpOption,
	},
}

// cmdEnablePrinterHandler is the "enable-printer" command handler.
func cmdEnablePrinterHandler(ctx context.Context, inv *argv.Invocation) error {
	printerURI := optPrinterGet(inv)
	if printerURI == "" {
		return fmt.Errorf("%s option required", optPrinter.Name)
	}

	clnt := cups.NewClient(optCUPSURL, nil)
	return clnt.EnablePrinter(ctx, printerURI)
}
