// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Error types fuzz test
//
// Pure Go logic, no interpreter needed. Except.object()'s exceptTable
// lookup is covered separately by FuzzCallbackErrPython.

package cpython

import (
	"fmt"
	"testing"
)

// FuzzErrPythonError fuzzes ErrPython.Error's "<except>: <msg>" format.
func FuzzErrPythonError(f *testing.F) {
	f.Add("ValueError", "boom")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, except, msg string) {
		e := ErrPython{except: Except(except), msg: msg}
		want := except + ": " + msg
		if got := e.Error(); got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}
	})
}

// FuzzErrPythonIs fuzzes ErrPython.Is against Except and ErrPython targets.
func FuzzErrPythonIs(f *testing.F) {
	f.Add("ValueError", "a", "ValueError", "b")
	f.Add("TypeError", "x", "ValueError", "x")

	f.Fuzz(func(t *testing.T, except1, msg1, except2, msg2 string) {
		e := ErrPython{except: Except(except1), msg: msg1}

		if got, want := e.Is(Except(except2)), except1 == except2; got != want {
			t.Fatalf("Is(Except(%q)) = %v, want %v", except2, got, want)
		}

		target := ErrPython{except: Except(except2), msg: msg2}
		want := except1 == except2 && msg1 == msg2
		if got := e.Is(target); got != want {
			t.Fatalf("Is(ErrPython{%q,%q}) = %v, want %v",
				except2, msg2, got, want)
		}
	})
}

// FuzzErrTypeConversionError fuzzes ErrTypeConversion.Error's message
// format.
func FuzzErrTypeConversionError(f *testing.F) {
	f.Add("int", "str")

	f.Fuzz(func(t *testing.T, from, to string) {
		e := ErrTypeConversion{from: from, to: to}
		want := fmt.Sprintf("can't convert %s to %s", from, to)
		if got := e.Error(); got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}
	})
}

// FuzzErrNotFound fuzzes ErrNotFound.Error and ErrNotFound.Is.
func FuzzErrNotFound(f *testing.F) {
	f.Add("key1", "key2")
	f.Add("", "key")

	f.Fuzz(func(t *testing.T, name1, name2 string) {
		e := ErrNotFound{name: name1}

		want := "item not found"
		if name1 != "" {
			want = fmt.Sprintf("%s: %s", name1, want)
		}
		if got := e.Error(); got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}

		wantIs := name1 == name2 || name2 == ""
		if got := e.Is(ErrNotFound{name: name2}); got != wantIs {
			t.Fatalf("Is(ErrNotFound{%q}) = %v, want %v",
				name2, got, wantIs)
		}
	})
}
