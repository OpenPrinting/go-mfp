// MFP          - Miulti-Function Printers and scanners toolkit
// cmd/mfp-cups - CUPS client
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import cups "github.com/OpenPrinting/go-mfp/cmd/mfp-cups/command"

// main function for the mfp-cups command
func main() {
	cups.Command.Main(nil)
}
