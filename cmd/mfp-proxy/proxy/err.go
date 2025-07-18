// MFP - Miulti-Function Printers and scanners toolkit
// The "proxy" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Common errors

package proxy

import "errors"

var (
	// ErrShutdown indicates that proxy shutdown is in progress
	ErrShutdown = errors.New("Proxy shutdown")
)
