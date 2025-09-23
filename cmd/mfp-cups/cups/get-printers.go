// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "get-printers" command.

package cups

import (
	"context"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
	"github.com/OpenPrinting/go-mfp/internal/env"
)

// cmdGetPrinters defines the "list-printers" sub-command.
var cmdListPrinters = argv.Command{
	Name:    "list-printers",
	Help:    "Get information on configured printers",
	Handler: cmdListPrintersHandler,
	Options: []argv.Option{
		optAttrs,
		optID,
		optLimit,
		optLocation,
		optUser,
		argv.HelpOption,
	},
}

// cmdListPrintersHandler is the "list-printers" command handler
func cmdListPrintersHandler(ctx context.Context, inv *argv.Invocation) error {
	// Prepare arguments
	dest := optCUPSURL(inv)

	sel := &cups.GetPrintersSelection{
		PrinterID:       optIDGet(inv),
		Limit:           optLimitGet(inv),
		PrinterLocation: optLocationGet(inv),
		User:            optUserGet(inv),
	}

	attrList := optAttrsGet(inv)
	attrList = append(attrList, prnAttrsRequested...)

	// Perform the query
	clnt := cups.NewClient(dest, nil)
	printers, err := clnt.CUPSGetPrinters(ctx, sel, attrList)
	if err != nil {
		return err
	}

	// Format output
	pager := env.NewPager()

	pager.Printf("CUPS: %s", dest)
	for _, prn := range printers {
		pager.Printf("")
		prnAttrsFormat(pager, prn)
	}

	return pager.Display()
}
