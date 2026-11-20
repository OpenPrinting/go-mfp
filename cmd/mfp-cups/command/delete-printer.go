// MFP - Multi-Function Printers and scanners toolkit
// The "mfp-cups" command
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions
//
// The "delete-printer" command.

package command

import (
	"context"
	"fmt"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
)

// cmdDeletePrinter defines the "delete-printer" sub-command.
var cmdDeletePrinter = argv.Command{
	Name:    "delete-printer",
	Help:    "Delete a print queue",
	Handler: cmdDeletePrinterHandler,
	Options: []argv.Option{
		optPrinterURI,
		argv.HelpOption,
	},
}

// cmdDeletePrinterHandler is the "delete-printer" command handler.
func cmdDeletePrinterHandler(ctx context.Context, inv *argv.Invocation) error {
	printerURI := optPrinterURIGet(inv)
	if printerURI == "" {
		return fmt.Errorf("%s option required", optPrinterURI.Name)
	}

	clnt := cups.NewClient(optCUPSURL, nil)
	return clnt.CUPSDeletePrinter(ctx, printerURI)
}
