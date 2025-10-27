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
	return ipp.NewPrinter(attrs)
}

// ippLoad decodes the IPP part of the model. The model file assumed to
// be already loaded into the Model's Python interpreter (model.py).
func (model *Model) ippLoad() error {
	// Load and decode printer capabilities
	obj, err := model.py.Eval("ipp.attrs")
	if err != nil {
		err = fmt.Errorf("ipp.attrs: %s", err)
		return err
	}

	if !obj.IsNone() {
		pa, err := model.pyImportPrinterAppributes(obj)
		if err != nil {
			err = fmt.Errorf("escl.caps: %s", err)
			return err
		}

		model.ippPrinterAttrs = pa
	}

	// Load IPP hooks
	// TODO

	return nil
}
