// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Keyword arguments

package cpython

// KWArg represents a keyword argument (name = value) for the
// Python function call.
type KWArg struct {
	Name  string
	Value any
}
