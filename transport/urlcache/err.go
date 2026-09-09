// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URL errors

package urlcache

import "errors"

// URL errors:
var (
	ErrURLInvalid       = errors.New(`URL: syntax error`)
	ErrURLSchemeMissed  = errors.New(`URL: missed scheme`)
	ErrURLSchemeInvalid = errors.New(`URL: invalid scheme`)
	ErrURLHostMissed    = errors.New(`URL: missed host`)
	ErrURLUNIXHost      = errors.New(`URL: host must be "localhost" or empty`)
)
