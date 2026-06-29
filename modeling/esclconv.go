// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// eSCL type conversions

package modeling

import (
	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/proto/escl"
)

// esclImportScannerCapabilities imports escl.ScannerCapabilities
// from the Python object.
func esclImportScannerCapabilities(obj *cpython.Object) (
	*escl.ScannerCapabilities, error) {

	var caps escl.ScannerCapabilities
	err := structImport(obj, keywordMapESCL, &caps)
	if err != nil {
		return nil, err
	}

	return &caps, nil
}

// esclImportScanSettings imports escl.ScanSettings from the Python object.
func esclImportScanSettings(obj *cpython.Object) (*escl.ScanSettings, error) {
	var ss escl.ScanSettings
	err := structImport(obj, keywordMapESCL, &ss)
	if err != nil {
		return nil, err
	}

	return &ss, nil
}

// esclDecodeJobStateReason decodes escl.JobStateReason from the Python object
func esclDecodeJobStateReason(obj *cpython.Object) (escl.JobStateReason, error) {
	s, err := obj.Str()
	if err != nil {
		return "", err
	}

	return escl.JobStateReason(s), nil
}

// esclDecodeVersion decodes escl.Version from the Python object
func esclDecodeVersion(obj *cpython.Object) (escl.Version, error) {
	s, err := obj.Str()
	if err != nil {
		return 0, err
	}

	ver, err := escl.DecodeVersion(s)
	if err != nil {
		return 0, err
	}

	return ver, nil
}
