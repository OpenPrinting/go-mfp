// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "get-ppd" command.

package cups

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
)

// cmdGetPPD defines the "get-ppd" sub-command.
var cmdGetPPD = argv.Command{
	Name:    "get-ppd",
	Help:    "Get PPD file by printer URI of PPD file name",
	Handler: cmdGetPPDHandler,
	Options: []argv.Option{
		optPrinterURI,
		optPPDName,
		argv.HelpOption,
	},
}

// cmdGetPPDHandler is the "get-ppd" command handler
func cmdGetPPDHandler(ctx context.Context, inv *argv.Invocation) error {
	// Validate options
	printerURI := optPrinterURIGet(inv)
	ppdName := optPPDNameGet(inv)

	switch {
	case printerURI == "" && ppdName == "":
		return fmt.Errorf("either %s or %s option required",
			optPrinterURI.Name, optPPDName.Name)
	case printerURI != "" && ppdName != "":
		return fmt.Errorf("confliction options: %s and %s",
			optPrinterURI.Name, optPPDName.Name)
	}

	// Perform the query
	clnt := cups.NewClient(optCUPSURL, nil)
	body, uri, err := clnt.CUPSGetPPD(ctx, printerURI, ppdName)
	if err != nil {
		return err
	}

	// Format output
	if body == nil {
		fmt.Printf("Use printer-uri: %q\n", uri)
	} else {
		_, err = io.Copy(os.Stdout, body)
		body.Close()
		if err != nil {
			fmt.Printf("\n%s\n", err)
		}
	}

	return nil
}
