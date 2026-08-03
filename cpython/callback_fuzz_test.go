// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Python->Go callbacks fuzz test

package cpython

import (
	"errors"
	"testing"
	"unicode/utf8"
)

// FuzzCallbackName fuzzes the name used to register a Go callback as a
// Python callable. Registration goes through Python.Set (dict item
// assignment), not exec'd source, so name need not be a valid Python
// identifier and may contain arbitrary bytes, including NUL.
func FuzzCallbackName(f *testing.F) {
	f.Add("")
	f.Add("go_fn")
	f.Add("with space")
	f.Add("emb\x00edded")
	f.Add("日本語")

	f.Fuzz(func(t *testing.T, name string) {
		if !utf8.ValidString(name) {
			// makeString builds a real PyUnicode object, which
			// requires valid UTF-8; Save legitimately fails for
			// invalid input. Not a bug, so skip rather than fail.
			t.Skip("not valid UTF-8")
		}

		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(func() int { return 1 })
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		if err := obj.Save(name); err != nil {
			t.Fatalf("Save(%q): %s", name, err)
		}

		got := py.Get(name)
		if got.IsNone() || !got.IsCallable() {
			t.Fatalf("Get(%q): not found or not callable", name)
		}

		ret := got.Call()
		if err := ret.Err(); err != nil {
			t.Fatalf("Call after Save(%q): %s", name, err)
		}
		if v, err := ret.Int(); err != nil || v != 1 {
			t.Fatalf("Call after Save(%q): got %v, err %v", name, v, err)
		}
	})
}

// FuzzCallbackReturnInt fuzzes the value returned by a single-return-value
// Go callback, round-tripped through a real Python call.
func FuzzCallbackReturnInt(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1) << 62)

	f.Fuzz(func(t *testing.T, n int64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(func() int64 { return n })
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		if err := obj.Save("go_fn"); err != nil {
			t.Fatalf("Save: %s", err)
		}

		ret := py.Eval("go_fn()")
		if err := ret.Err(); err != nil {
			t.Fatalf("Eval: %s", err)
		}
		if v, err := ret.Int(); err != nil || v != n {
			t.Fatalf("expected %d, got %v (err %v)", n, v, err)
		}
	})
}

// FuzzCallbackErrPython fuzzes the (value, ErrPython) return path,
// exercising callbackSetError's ErrPython branch with an arbitrary,
// possibly-unrecognized exception name and message.
func FuzzCallbackErrPython(f *testing.F) {
	f.Add("ValueError", "boom")
	f.Add("", "")
	f.Add("NotARealException", "emb\x00edded\nmsg")
	f.Add("Warning", string([]byte{0xff, 0xfe}))

	f.Fuzz(func(t *testing.T, except, msg string) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		pyErr := ErrPython{except: Except(except), msg: msg}
		obj := py.NewObject(func() (int, error) { return 0, pyErr })
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		if err := obj.Save("go_fn"); err != nil {
			t.Fatalf("Save: %s", err)
		}

		ret := py.Eval("go_fn()")
		if ret.Err() == nil {
			t.Fatal("expected error to propagate, got nil")
		}
	})
}

// FuzzCallbackPlainError is like FuzzCallbackErrPython but wraps the
// fuzzed message in a plain Go error, exercising callbackSetError's
// non-ErrPython (generic SystemError-ish) message path.
func FuzzCallbackPlainError(f *testing.F) {
	f.Add("plain error")
	f.Add("")
	f.Add("emb\x00edded")

	f.Fuzz(func(t *testing.T, msg string) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(func() (int, error) { return 0, errors.New(msg) })
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		if err := obj.Save("go_fn"); err != nil {
			t.Fatalf("Save: %s", err)
		}

		ret := py.Eval("go_fn()")
		if ret.Err() == nil {
			t.Fatal("expected error to propagate, got nil")
		}
	})
}
