// MFP - Miulti-Function Printers and scanners toolkit
// The "mfp-cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "list-printers" command.

package command

import (
	"context"
	"fmt"
	"os"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// cmdListPPDs defines the "list-ppds" sub-command.
var cmdListPPDs = argv.Command{
	Name:    "list-ppds",
	Help:    "Get information on available PPDs",
	Handler: cmdListPPDsHandler,
	Options: []argv.Option{
		optSchemesExclude,
		optSchemesInclude,
		optLimit,
		optMake,
		optMakeModel,
		optModelNumber,
		optLanguage,
		optProduct,
		optPsVersion,
		optPPDType,
		argv.HelpOption,
	},
}

// cmdListPPDsHandler is the "list-ppds" command handler
func cmdListPPDsHandler(ctx context.Context, inv *argv.Invocation) error {
	// Prepare arguments
	sel := &cups.GetPPDsSelection{
		ExcludeSchemes:     optSchemesExcludeGet(inv),
		IncludeSchemes:     optSchemesIncludeGet(inv),
		Limit:              optLimitGet(inv),
		PPDMake:            optMakeGet(inv),
		PPDMakeAndModel:    optMakeModelGet(inv),
		PPDModelNumber:     optoptModelNumberGet(inv),
		PPDNaturalLanguage: optLanguageGet(inv),
		PPDProduct:         optProductGet(inv),
		PPDPsVersion:       optPsVersionGet(inv),
		PPDType:            optPPDTypeGet(inv),
	}

	// Perform the query
	clnt := cups.NewClient(optCUPSURL, nil)
	clnt.SetDecoderOptions(&ipp.DecoderOptions{KeepTrying: true})

	ppds, err := clnt.CUPSGetPPDs(ctx, sel, nil)
	if err != nil {
		return err
	}

	// Format output
	for _, prn := range ppds {
		fmt.Println("")
		ppdAttrsFormat(os.Stdout, prn)
	}

	return nil
}
