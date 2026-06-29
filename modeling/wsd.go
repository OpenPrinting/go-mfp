// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// WSD (WS-Scan) part of Model

package modeling

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
)

// SetWSDScanCaps sets the WS-Scan scanner capabilities.
//
// It requires [wsscan.GetScannerElementsResponse] with filled
// [wsscan.GetScannerElementsResponse.ScannerConfiguration] and
// [wsscan.GetScannerElementsResponse.ScannerDescription] fields.
func (model *Model) SetWSDScanCaps(caps *wsscan.GetScannerElementsResponse) {
	model.wsdScanCaps = caps
}

// GetWSDScanCaps returns the WS-Scan scanner capabilities, previously
// set with [Model.SetWSDScanCaps].
func (model *Model) GetWSDScanCaps() *wsscan.GetScannerElementsResponse {
	return model.wsdScanCaps
}

// NewWSDServer creates a virtual WS-Scan server based on WS-Scan scanner
// capabilities, defined by the model (see also [Model.SetWSDScanCaps]).
//
// The actual scanning facilities provided by the supplied [abstract.Scanner].
//
// It will return nil, if model doesn't have the eSCL scanner capabilities.
func (model *Model) NewWSDServer(
	scanner abstract.Scanner) *wsscan.AbstractServer {

	// Obtain scanner capabilities
	caps := model.GetWSDScanCaps()
	if caps == nil {
		return nil
	}

	// Setup hooks
	// TODO

	// Setup options
	options := wsscan.AbstractServerOptions{
		Scanner:  scanner,
		BasePath: "/WSScan",
	}

	// Create the WS-Scan server
	return wsscan.NewAbstractServer(options)
}

// wsdLoad decodes WS-Scan part of model. The model file assumed to
// be preloaded into the Model's Python interpreter (model.py).
func (model *Model) wsdLoad() error {
	// Load wsscan.GetScannerElementsResponse
	obj := model.py.Eval("wsd.caps")
	if err := obj.Err(); err != nil {
		err = fmt.Errorf("wsd.caps: %w", err)
		return err
	}

	if obj.IsNone() {
		return nil
	}

	// Decode wsscan.GetScannerElementsResponse
	var caps *wsscan.GetScannerElementsResponse
	err := structImport(obj, keywordMapWSD, &caps)
	if err != nil {
		err = fmt.Errorf("wsd.caps: %w", err)
		return err
	}

	model.wsdScanCaps = caps
	return nil
}
