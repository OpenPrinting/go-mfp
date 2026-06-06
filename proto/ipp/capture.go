// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// IPP print backend hook

package ipp

import "github.com/OpenPrinting/go-mfp/abstract"

// SetPrintBackend installs backend as the handler for incoming
// print documents. It is called for each Send-Document request
// with the negotiated job parameters and a streaming reader for
// the document body. Pass nil to clear a previously set backend.
func (printer *Printer) SetPrintBackend(backend abstract.PrintBackend) {
	printer.backend = backend
}
