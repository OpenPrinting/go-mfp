// MFP          - Miulti-Function Printers and scanners toolkit
// cmd/mfp-cups - CUPS client
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import "github.com/OpenPrinting/go-mfp/cmd/mfp-virtual/virtual"

// main function for the mfp-cups command
func main() {
	virtual.Command.Main(nil)
}
