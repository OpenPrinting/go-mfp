// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "default-printer" command.

package cups

import (
	"context"
	"fmt"
	"os"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// cmdDefaultPrinter defines the "default-printer" sub-command
var cmdDefaultPrinter = argv.Command{
	Name:    "default-printer",
	Help:    "Get default printer",
	Handler: cmdDefaultPrinterHandler,
	Options: []argv.Option{
		optAttrs,
		argv.HelpOption,
	},
}

// cmdDefaultPrinterHandler is the "default-printer" command handler
func cmdDefaultPrinterHandler(ctx context.Context, inv *argv.Invocation) error {
	attrList := optAttrsGet(inv)
	attrList = append(attrList, prnAttrsRequested...)

	// Perform the query
	clnt := cups.NewClient(optCUPSURL, nil)
	clnt.SetDecoderOptions(&ipp.DecoderOptions{KeepTrying: true})

	prn, err := clnt.CUPSGetDefault(ctx, attrList)
	if err != nil {
		return err
	}

	// Format output
	fmt.Printf("CUPS: %s\n", optCUPSURL)
	prnAttrsFormat(os.Stdout, prn)

	return nil
}
