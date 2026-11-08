// MFP           - Miulti-Function Printers and scanners toolkit
// cmd/mfp-ppd   - Utility for PPD files
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import "github.com/OpenPrinting/go-mfp/cmd/mfp-ppd/command"

// main function for the mfp-ppd command
func main() {
	command.Command.Main(nil)
}
