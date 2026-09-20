// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP part of the Model

package modeling

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// SetIPPPrinterAttrs sets the IPP printer attributes.
func (model *Model) SetIPPPrinterAttrs(attrs *ipp.PrinterAttributes) {
	model.ippPrinterAttrs = attrs
}

// GetIPPPrinterAttrs returns the IPP printer attributes.
func (model *Model) GetIPPPrinterAttrs() *ipp.PrinterAttributes {
	return model.ippPrinterAttrs
}

// SetIPPScannerAttrs sets the IPP scanner attributes.
func (model *Model) SetIPPScannerAttrs(attrs *ipp.PrinterAttributes) {
	model.ippScannerAttrs = attrs
}

// GetIPPScannerAttrs returns the IPP scanner attributes.
func (model *Model) GetIPPScannerAttrs() *ipp.PrinterAttributes {
	return model.ippScannerAttrs
}

// NewIPPPrinter creates a virtual IPP printer.
// It will return nil, if model doesn't have the IPP printer attributes.
func (model *Model) NewIPPPrinter() *ipp.Printer {
	// Obtain printer attributes
	attrs := model.GetIPPPrinterAttrs()
	if attrs == nil {
		return nil
	}

	// Create the IPP print server
	options := ipp.PrinterOptions{
		UseRawPrinterAttributes: true,
	}
	return ipp.NewPrinter(attrs, options)
}

// NewIPPScanner creates a virtual IPP scanner.
// It will return nil, if model doesn't have the IPP scanner. attributes.
func (model *Model) NewIPPScanner(scanner abstract.Scanner) *ipp.Scanner {
	// Obtain printer attributes
	attrs := model.GetIPPScannerAttrs()
	if attrs == nil {
		return nil
	}

	// Create the IPP print server
	options := ipp.ScannerOptions{
		Scanner:                 scanner,
		UseRawPrinterAttributes: true,
	}
	return ipp.NewScanner(attrs, options)
}

// ippLoadPrinter decodes the IPP-print part of the model.
func (model *Model) ippLoadPrinter() error {
	// Load and decode printer capabilities
	name := "ipp.printer"
	obj := model.py.Eval(name)

	if err := obj.Err(); err != nil {
		err = fmt.Errorf("%s: %w", name, err)
		return err
	}

	if !obj.IsNone() {
		pa, err := ippImportPrinterAppributes(obj)
		if err != nil {
			err = fmt.Errorf("%s: %w", name, err)
			return err
		}

		model.ippPrinterAttrs = pa
	}

	// Load IPP hooks
	// TODO

	return nil
}

// ippLoadScanner decodes the IPP-scan part of the model.
func (model *Model) ippLoadScanner() error {
	// Load and decode printer capabilities
	name := "ipp.scanner"
	obj := model.py.Eval(name)

	if err := obj.Err(); err != nil {
		err = fmt.Errorf("%s: %w", name, err)
		return err
	}

	if !obj.IsNone() {
		pa, err := ippImportPrinterAppributes(obj)
		if err != nil {
			err = fmt.Errorf("%s: %w", name, err)
			return err
		}

		model.ippScannerAttrs = pa
	}

	// Load IPP hooks
	// TODO

	return nil
}
