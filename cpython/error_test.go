// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests for error types

package cpython

import (
	"strings"
	"testing"
)

// TestErrPython verifies that ErrPython implements the error
// interface and returns a non-empty message.
func TestErrPython(t *testing.T) {
	e := ErrPython{msg: "something went wrong"}
	if e.Error() != "something went wrong" {
		t.Fatalf("ErrPython.Error() = %q, want %q",
			e.Error(), "something went wrong")
	}
}

// TestErrPythonEmpty verifies that ErrPython works with empty message.
func TestErrPythonEmpty(t *testing.T) {
	e := ErrPython{msg: ""}
	if e.Error() != "" {
		t.Fatalf("ErrPython.Error() = %q, want empty string", e.Error())
	}
}

// TestErrTypeConversion verifies that ErrTypeConversion formats
// its message correctly.
func TestErrTypeConversion(t *testing.T) {
	e := ErrTypeConversion{from: "int", to: "string"}
	got := e.Error()
	if !strings.Contains(got, "int") || !strings.Contains(got, "string") {
		t.Fatalf("ErrTypeConversion.Error() = %q, "+
			"want it to contain 'int' and 'string'", got)
	}
}

// TestErrOverflow verifies that ErrOverflow formats its message correctly.
func TestErrOverflow(t *testing.T) {
	e := ErrOverflow{val: "99999999999999999999"}
	got := e.Error()
	if !strings.Contains(got, "99999999999999999999") {
		t.Fatalf("ErrOverflow.Error() = %q, "+
			"want it to contain the value", got)
	}
}

// TestErrClosed verifies that ErrClosed returns a non-empty message.
func TestErrClosed(t *testing.T) {
	e := ErrClosed{}
	if e.Error() == "" {
		t.Fatalf("ErrClosed.Error() returned empty string")
	}
}

// TestErrInvalidObject verifies that ErrInvalidObject returns
// a non-empty message.
func TestErrInvalidObject(t *testing.T) {
	e := ErrInvalidObject{}
	if e.Error() == "" {
		t.Fatalf("ErrInvalidObject.Error() returned empty string")
	}
}

// TestErrNotFound verifies that ErrNotFound returns a non-empty message.
func TestErrNotFound(t *testing.T) {
	e := ErrNotFound{}
	if e.Error() == "" {
		t.Fatalf("ErrNotFound.Error() returned empty string")
	}
}

// TestErrorsImplementInterface is a compile-time assertion that all
// error types implement the error interface.
func TestErrorsImplementInterface(t *testing.T) {
	var _ error = ErrPython{}
	var _ error = ErrTypeConversion{}
	var _ error = ErrOverflow{}
	var _ error = ErrClosed{}
	var _ error = ErrInvalidObject{}
	var _ error = ErrNotFound{}
}
