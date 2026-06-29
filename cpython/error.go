// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package cpython

import (
	"fmt"
)

// ErrPython represents a Python exception.
type ErrPython struct {
	except Except
	msg    string
}

// Error returns error message. It implements the [error] interface.
func (e ErrPython) Error() string {
	return string(e.except) + ": " + e.msg
}

// Is reports if ErrPython matches the target error.
//
// ErrPython matches on the following cases:
//   - target is ErrPython with the same exception type and message
//   - target is [Except] and exception type is the same
func (e ErrPython) Is(target error) bool {
	switch target := target.(type) {
	case Except:
		return e.except == target
	case ErrPython:
		return e == target
	}

	return false
}

// ErrTypeConversion represents Go<->Python type conversion error.
type ErrTypeConversion struct {
	from, to string // from/to types that can't be converted
}

// Error returns error message. It implements the [error] interface.
func (e ErrTypeConversion) Error() string {
	return fmt.Sprintf("can't convert %s to %s", e.from, e.to)
}

// ErrOverflow represents the integer overflow error.
type ErrOverflow struct {
	val string
}

// Error returns error message. It implements the [error] interface.
func (e ErrOverflow) Error() string {
	return fmt.Sprintf("integer overflow: %s", e.val)
}

// ErrClosed represent the error that occurs when [Python]
// interpreter or [Object] that it owns accessed after call
// to [Python.Close].
type ErrClosed struct{}

// Error returns error message. It implements the [error] interface.
func (e ErrClosed) Error() string {
	return "use Python interpreter after Python.Close"
}

// ErrInvalidObject represents the error that occurs when [Object]
// accessed after call to [Object.Invalidate]
type ErrInvalidObject struct{}

// Error returns error message. It implements the [error] interface.
func (e ErrInvalidObject) Error() string {
	return "use Object after Object.Invalidate"
}

// ErrNotFound represents the error that occurs when retrieving
// value from container (dict, array, ...) fails because key was
// not found in the container
type ErrNotFound struct{ name string }

// Error returns error message. It implements the [error] interface.
func (e ErrNotFound) Error() string {
	const s = "item not found"
	if e.name == "" {
		return s
	}
	return fmt.Sprintf("%s: %s", e.name, s)
}

// Is reports if ErrNotFound matches the target error.
//
// ErrNotFound matches on the following cases:
//   - target is ErrNotFound with the same item name
//   - target is ErrNotFound with the empty item name
func (e ErrNotFound) Is(target error) bool {
	switch target := target.(type) {
	case ErrNotFound:
		return e.name == target.name || target.name == ""
	}

	return false
}
