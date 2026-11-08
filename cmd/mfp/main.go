// MFP     - Miulti-Function Printers and scanners toolkit
// cmd/mfp - Universal command that implements all other as sub-commands
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenPrinting/go-mfp/argv"
	cups "github.com/OpenPrinting/go-mfp/cmd/mfp-cups/command"
	discover "github.com/OpenPrinting/go-mfp/cmd/mfp-discover/command"
	ippusb "github.com/OpenPrinting/go-mfp/cmd/mfp-ipp-usb/command"
	model "github.com/OpenPrinting/go-mfp/cmd/mfp-model/command"
	ppd "github.com/OpenPrinting/go-mfp/cmd/mfp-ppd/command"
	proxy "github.com/OpenPrinting/go-mfp/cmd/mfp-proxy/command"
	usb "github.com/OpenPrinting/go-mfp/cmd/mfp-usb/command"
	virtual "github.com/OpenPrinting/go-mfp/cmd/mfp-virtual/command"
)

var allCommands = []argv.Command{
	cups.Command,
	discover.Command,
	ippusb.Command,
	model.Command,
	ppd.Command,
	proxy.Command,
	usb.Command,
	virtual.Command,
}

// main function for the mfp command
func main() {
	base := filepath.Base(os.Args[0])
	for i := range allCommands {
		cmd := &allCommands[i]
		if cmd.Name == base {
			cmd.Main(nil)
			return
		}
	}

	fmt.Printf("%s is the universal executable fpr all mfp-xxx tools\n", base)
	fmt.Printf("\n")
	fmt.Printf("You can put it somewhere into your path, create symbolic links\n")
	fmt.Printf("and you have a whole set of `mfp-XXX` commands in a single executable:\n")
	fmt.Printf("\n")

	for i := range allCommands {
		fmt.Printf("    ln -s %s %s\n", base, allCommands[i].Name)
	}

	fmt.Printf("\n")
	fmt.Printf("It saves a lot of installation space, providing the same functionality\n")
}
