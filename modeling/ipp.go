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

	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// SetIPPPrinterAttrs sets the [escl.ScannerCapabilities].
func (model *Model) SetIPPPrinterAttrs(attrs *ipp.PrinterAttributes) {
	model.ippPrinterAttrs = attrs
}

// GetIPPPrinterAttrs returns the [escl.ScannerCapabilities].
func (model *Model) GetIPPPrinterAttrs() *ipp.PrinterAttributes {
	return model.ippPrinterAttrs
}

// NewIPPServer creates a virtual IPP server.
// It will return nil, if model doesn't have the IPP printer attributes.
func (model *Model) NewIPPServer() *ipp.Printer {
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

// ippLoad decodes the IPP part of the model. The model file assumed to
// be already loaded into the Model's Python interpreter (model.py).
func (model *Model) ippLoad() error {
	// Load and decode printer capabilities
	name := "ipp.printer"
	obj := model.py.Eval(name)
	if obj.Err() != nil {
		name = "ipp.attrs"
		obj = model.py.Eval(name)
	}

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
