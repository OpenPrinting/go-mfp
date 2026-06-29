// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Tests for Python exceptions
package cpython

import "testing"

// Compile-time assertions that Except implements the string type.
var _ = Except("")

// TestExceptConstants verifies that standard exception constants
// have the correct string values.
func TestExceptConstants(t *testing.T) {
	tests := []struct {
		except Except
		want   string
	}{
		{ArithmeticError, "ArithmeticError"},
		{AssertionError, "AssertionError"},
		{AttributeError, "AttributeError"},
		{BlockingIOError, "BlockingIOError"},
		{EOFError, "EOFError"},
		{Exception, "Exception"},
		{FileNotFoundError, "FileNotFoundError"},
		{ImportError, "ImportError"},
		{IndexError, "IndexError"},
		{KeyError, "KeyError"},
		{MemoryError, "MemoryError"},
		{NameError, "NameError"},
		{NotImplementedError, "NotImplementedError"},
		{OSError, "OSError"},
		{OverflowError, "OverflowError"},
		{RuntimeError, "RuntimeError"},
		{StopIteration, "StopIteration"},
		{SyntaxError, "SyntaxError"},
		{SystemError, "SystemError"},
		{TypeError, "TypeError"},
		{ValueError, "ValueError"},
		{ZeroDivisionError, "ZeroDivisionError"},
		{DeprecationWarning, "DeprecationWarning"},
		{RuntimeWarning, "RuntimeWarning"},
		{UserWarning, "UserWarning"},
		{Warning, "Warning"},
	}
	for _, tt := range tests {
		if string(tt.except) != tt.want {
			t.Fatalf("Except %q = %q, want %q",
				tt.except, string(tt.except), tt.want)
		}
	}
}

// TestExceptObjectKnown verifies that object() returns non-nil
// for all known standard exceptions.
func TestExceptObjectKnown(t *testing.T) {
	known := []Except{
		ArithmeticError, AssertionError, AttributeError,
		BlockingIOError, EOFError, Exception,
		FileNotFoundError, ImportError, IndexError,
		KeyError, MemoryError, NameError,
		NotImplementedError, OSError, OverflowError,
		RuntimeError, StopIteration, SyntaxError,
		SystemError, TypeError, ValueError,
		ZeroDivisionError, DeprecationWarning,
		RuntimeWarning, UserWarning, Warning,
	}
	for _, ex := range known {
		if ex.object() == nil {
			t.Fatalf("Except(%q).object() returned nil", ex)
		}
	}
}

// TestExceptObjectUnknown verifies that object() falls back to
// SystemError for unknown exception names.
func TestExceptObjectUnknown(t *testing.T) {
	unknown := Except("NoSuchException")
	obj := unknown.object()
	if obj == nil {
		t.Fatalf("Except(%q).object() returned nil, want SystemError fallback", unknown)
	}
}

// TestExceptError verifies that the Error() method returns the exception
// name as a string, satisfying the error interface.
func TestExceptError(t *testing.T) {
	tests := []struct {
		except Except
		want   string
	}{
		{ArithmeticError, "ArithmeticError"},
		{ValueError, "ValueError"},
		{Warning, "Warning"},
		{Except("CustomError"), "CustomError"},
	}

	for _, tt := range tests {
		got := tt.except.Error()
		if got != tt.want {
			t.Fatalf("Except(%q).Error() = %q, want %q", tt.except, got, tt.want)
		}
	}
}

// TestExceptImplementsError verifies Except satisfies the error interface
// at compile time and that the method is callable via the interface.
func TestExceptImplementsError(t *testing.T) {
	var err error = ValueError
	if err.Error() != "ValueError" {
		t.Fatalf("error.Error() = %q, want %q", err.Error(), "ValueError")
	}
}
