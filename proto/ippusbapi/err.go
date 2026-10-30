// MFP - Multi-Function Printers and scanners toolkit
// ipp-usb API client
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Errors

package ippusbapi

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
)

var (
	// ErrNoIppUsb is returned when ipp-usb daemon is not running
	ErrNoIppUsb = errors.New("ipp-usb daemon not running")

	// ErrAccess indicates an access error
	ErrAccess = errors.New("access denied")

	// ErrNotSupported indicates that operation is not
	// supported by the running ipp-usb daemon
	ErrNotSupported = errors.New("operation not supported")
)

// ErrHTTP represents an error, caused by non-successful HTTP status
type ErrHTTP int

// Error returns error string. It implements the [error] interface.
func (err ErrHTTP) Error() string {
	return fmt.Sprintf("%3.3d %s", int(err), http.StatusText(int(err)))
}

// errDecode decodes known errors.
func errDecode(err error) error {
	switch {
	case errors.Is(err, syscall.ECONNREFUSED):
		err = ErrNoIppUsb
	case errors.Is(err, syscall.ENOENT):
		err = ErrNoIppUsb
	case errors.Is(err, syscall.EACCES):
		err = ErrAccess
	case errors.Is(err, syscall.EPERM):
		err = ErrAccess
	case errors.Is(err, ErrHTTP(http.StatusNotFound)):
		err = ErrNotSupported
	}

	return err
}
