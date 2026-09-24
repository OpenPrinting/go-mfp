// MFP - Miulti-Function Printers and scanners toolkit
// CUPS Client and Server
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Common types

package cups

import (
	"time"
)

// Default values for common types:
var (
	DefaultGetPrintersSelection = &GetPrintersSelection{}
	DefaultGetDevicesSelection  = &GetDevicesSelection{
		Timeout: DefaultGetDevicesTimeout,
	}
	DefaultGetDevicesTimeout = 5 * time.Second
)

// GetPrintersSelection configures a selection of printers returned
// by [Client.CUPSGetPrinters].
type GetPrintersSelection struct {
	// Printer name (also, queue name) is the unique name, under
	// which printer is registered in the CUPS system.
	//
	// If this parameter is not empty, it specifies the first printer
	// name to be returned.
	FirstPrinterName string

	// If not zero, specifies maximum number of printers to be returned.
	Limit int

	// Each printer in the CUPS system has its own PrinterID,
	// which is system-unique positive number, assigned to the
	// printer when it is added to the system.
	//
	// If this parameter is not zero, only the printer with specified
	// ID (and matching other criteria) will be returned.
	PrinterID int

	// PrinterLocation is the arbitrary string which can be
	// configured on the printer device by owner or system
	// administrator. Its purpose is to simplify selection of
	// particular devices in a big network. It corresponds to the
	// "printer-location" IPP attribute and may sound like
	// "Printers on 1st floor" or "Printers at reception".
	//
	// If this parameter is not empty, only devices with the specified
	// locatiom will be returned.
	PrinterLocation string

	// TODO
	PrinterType     int
	PrinterTypeMask int

	// If not empty, only printers accessible to that user will
	// be returned. User name is the user's **login** name,
	User string
}

// GetDevicesSelection configures a selection of devices returned by
// [Client.CUPSGetDevices].
type GetDevicesSelection struct {
	// ExcludeSchemes and IncludeSchemes specifies which schemes
	// to include into output.
	//
	// CUPS supports the following schemes:
	//	- http, https
	//	- ipp, ipps
	//	- lpd
	//	- smb
	//	- socket
	//
	// This list may change and actually depends on installed backends.
	//
	// By default, CUPS returns all available schemes. If ExcludeSchemes
	// is not empty, schemes listed there are excluded. If IncludeSchemes
	// is not empty, schemes not listed there are not included.
	ExcludeSchemes, IncludeSchemes []string

	// If not zero, specifies maximum number of devices to be returned.
	Limit int

	// Timeout, if set, specifies a maximum search time.
	//
	// Note, Timeout is rounded up to a whole number of seconds.
	// If not set, CUPS server default is used.
	Timeout time.Duration
}

// GetPPDsSelection configures a selection of devices returned by
// [Client.CUPSGetPPDs].
type GetPPDsSelection struct {
	// ExcludeSchemes and IncludeSchemes specifies which schemes
	// to include into output.
	ExcludeSchemes []string
	IncludeSchemes []string

	// If not zero, specifies maximum number of PPDs to be returned.
	Limit int

	// These parameters, allows to filter returned  PPDs by
	// make (manufacturer), make and model, model number,
	// language, PostSript product string, PostScript version
	// string and PPD type, respectively.
	//
	// Use "" for strings and 0 for int to bypass the particular
	// filtering.
	PpdMake            string
	PpdMakeAndModel    string
	PpdModelNumber     int
	PpdNaturalLanguage string
	PpdProduct         string
	PpdPsVersion       string
	PpdType            string
}
