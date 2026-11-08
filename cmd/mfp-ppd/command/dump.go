// MFP - Miulti-Function Printers and scanners toolkit
// The "mfp-ppd" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "dump" command.

package command

import (
	"context"
	"os"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/proto/ppd"
	"github.com/OpenPrinting/goipp"
)

// cmdDump defines the "Dump" sub-command
var cmdDump = argv.Command{
	Name:    "dump",
	Help:    "Dump PPD file as IPP attributes",
	Handler: cmdDumpHandler,
	Parameters: []argv.Parameter{{
		Name:     "file",
		Help:     "PPD file name",
		Complete: argv.CompleteOSPath,
	}},
}

// cmdDumpHandler is the "list" command handler.
func cmdDumpHandler(ctx context.Context, inv *argv.Invocation) error {
	// Read the PPD file
	file, _ := inv.Get("file")
	ppddata, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Convert to IPP attributes
	attrs, err := ppd.ToIPP(ppddata)
	if err != nil {
		return err
	}

	// Format output
	f := goipp.NewFormatter()
	f.FmtGroup(goipp.Group{
		Tag:   goipp.TagPrinterGroup,
		Attrs: attrs,
	})

	_, err = f.WriteTo(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
