// MFP - Multi-Function Printers and scanners toolkit
// Abstract definition for printer and scanner interfaces
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// The print backend interface

package abstract

import "io"

// PrintJobParams contains protocol-independent parameters of a
// print job, as negotiated between the client and the printer.
//
// Fields are optional; zero value means the parameter was not
// provided by the protocol layer.
type PrintJobParams struct {
	// Format is the MIME type of the document
	// (e.g., "application/pdf", "image/pwg-raster").
	Format string

	// JobName is the name of the print job, if provided.
	JobName string

	// Copies is the requested number of copies.
	// Zero means unset (use printer default).
	Copies int

	// Sides controls duplex printing.
	// Common values: "one-sided", "two-sided-long-edge",
	// "two-sided-short-edge". Empty string means unset.
	Sides string

	// ColorMode is the requested color mode.
	// Common values: "color", "monochrome", "auto".
	// Empty string means unset.
	ColorMode string

	// Media is the requested media size or type
	// (e.g., "iso_a4_210x297mm", "na_letter_8.5x11in").
	// Empty string means unset.
	Media string

	// Variables holds protocol-specific key=value parameters that
	// do not map to the standard fields above.
	// For IEEE 1284/PJL, these are @PJL SET key=value pairs.
	// Keys are stored in uppercase.
	Variables map[string]string
}

// PrintBackend is the protocol-independent interface for receiving
// print jobs from the virtual printer.
//
// It is the printing analogue of [Scanner] for the scanning side.
// Implementations are called by the protocol layer (IPP, IEEE 1284)
// when a print job is ready to be processed.
type PrintBackend interface {
	// PrintDocument is called when a new print document arrives.
	//
	// params contains the negotiated job parameters extracted from
	// the protocol layer (IPP job attributes, PJL commands, etc.).
	//
	// body provides streaming access to the document data.
	// The implementation must fully consume body before returning.
	// body is valid only for the duration of this call.
	PrintDocument(params PrintJobParams, body io.Reader) error
}
