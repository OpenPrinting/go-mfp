// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "list-printers" command.

package cups

import (
	"context"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
	"github.com/OpenPrinting/go-mfp/internal/env"
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
	dest := optCUPSURL(inv)

	sel := &cups.GetPPDsSelection{
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
	clnt := cups.NewClient(dest, nil)
	clnt.SetDecoderOptions(&ipp.DecoderOptions{KeepTrying: true})

	ppds, err := clnt.CUPSGetPPDs(ctx, sel, nil)
	if err != nil {
		return err
	}

	// Format output
	pager := env.NewPager()

	pager.Printf("CUPS: %s", dest)
	for _, prn := range ppds {
		pager.Printf("")
		ppdAttrsFormat(pager, prn)
	}

	return pager.Display()
}
