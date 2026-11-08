// MFP           - Miulti-Function Printers and scanners toolkit
// cmd/mfp-proxy - IPP/eSCL/WSD proxy
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import "github.com/OpenPrinting/go-mfp/cmd/mfp-proxy/command"

// main function for the mfp-cups command
func main() {
	command.Command.Main(nil)
}
