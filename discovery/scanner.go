// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Scanner information

package discovery

import (
	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// ScannerParameters represents the discoverable information about the scanner.
type ScannerParameters struct {
	// Scanner capabilities
	Duplex optional.Val[bool]                 // Duplex mode supported
	Inputs generic.Bitset[abstract.Input]     // Supported inputs
	Colors generic.Bitset[abstract.ColorMode] // Supported color modes
	PDL    []string                           // Supported MIME types
}
