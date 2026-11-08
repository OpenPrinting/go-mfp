// MFP           - Miulti-Function Printers and scanners toolkit
// cmd/mfp-usb   - Troubleshooting tool for USB-connected MFPs
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import "github.com/OpenPrinting/go-mfp/cmd/mfp-usb/command"

// main function for the mfp-usb command
func main() {
	command.Command.Main(nil)
}
