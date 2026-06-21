// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Tests for error types
package cpython

import (
	"errors"
	"strings"
	"testing"
)

// Compile-time assertions that all error types implement the error interface.
var (
	_ error = ErrPython{}
	_ error = ErrTypeConversion{}
	_ error = ErrOverflow{}
	_ error = ErrClosed{}
	_ error = ErrInvalidObject{}
	_ error = ErrNotFound{}
)

// TestErrPython verifies that ErrPython returns a correctly formatted message.
func TestErrPython(t *testing.T) {
	e := ErrPython{except: "RuntimeError", msg: "something went wrong"}
	got := e.Error()
	if !strings.Contains(got, "RuntimeError") || !strings.Contains(got, "something went wrong") {
		t.Fatalf("ErrPython.Error() = %q, want it to contain exception and message", got)
	}
}

// TestErrPythonIs verifies all branches of ErrPython.Is directly.
func TestErrPythonIs(t *testing.T) {
	e := ErrPython{except: "RuntimeError", msg: "something went wrong"}

	// case Except: → true  (hit via errors.Is which calls e.Is with Except as error)
	if !errors.Is(e, Except("RuntimeError")) {
		t.Error("expected ErrPython to match Except with same type")
	}

	// case Except: → false
	if errors.Is(e, Except("ValueError")) {
		t.Error("expected ErrPython not to match Except with different type")
	}

	// case ErrPython: → true  (call Is directly so type switch sees ErrPython, not error)
	same := ErrPython{except: "RuntimeError", msg: "something went wrong"}
	if !e.Is(same) {
		t.Error("expected ErrPython.Is to return true for identical ErrPython")
	}

	// case ErrPython: → false
	diff := ErrPython{except: "RuntimeError", msg: "different message"}
	if e.Is(diff) {
		t.Error("expected ErrPython.Is to return false for different ErrPython")
	}

	// default return false — unrelated type
	if e.Is(ErrClosed{}) {
		t.Error("expected ErrPython.Is to return false for unrelated error type")
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
		t.Fatalf("ErrOverflow.Error() = %q, want it to contain the value", got)
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

