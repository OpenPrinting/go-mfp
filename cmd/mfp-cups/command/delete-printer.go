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

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
)

// cmdDeletePrinter defines the "delete-printer" sub-command.
var cmdDeletePrinter = argv.Command{
	Name:    "delete-printer",
	Help:    "Delete a print queue",
	Handler: cmdDeletePrinterHandler,
	Options: []argv.Option{
		argv.HelpOption,
	},
	Parameters: []argv.Parameter{
		{
			Name: "printer",
			Help: "Queue name",
		},
	},
}

// cmdDeletePrinterHandler is the "delete-printer" command handler.
func cmdDeletePrinterHandler(ctx context.Context, inv *argv.Invocation) error {
	clnt := cups.NewClient(optCUPSURL, nil)
	return clnt.CUPSDeletePrinter(ctx, inv.GetString("printer"))
}
