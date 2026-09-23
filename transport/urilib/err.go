// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URI errors

package urilib

import "errors"

// URI errors:
var (
	ErrURIInvalid       = errors.New(`URI: syntax error`)
	ErrURISchemeMissed  = errors.New(`URI: missed scheme`)
	ErrURISchemeInvalid = errors.New(`URI: invalid scheme`)
	ErrURIHostMissed    = errors.New(`URI: missed host`)
	ErrURIUNIXHost      = errors.New(`URI: host must be "localhost" or empty`)
)
