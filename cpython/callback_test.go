// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Python->Go callbacks test

package cpython

import (
	"errors"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestNewCallbackInvalidSignature tests that newCallback returns nil for
// functions with unsupported return signatures:
//   - two return values where the second is not error
//   - more than two return values
//
// NOTE: the following branches in callback.go are not tested because they
// are unreachable from Go-level tests:
//   - callback.object: makeCapsule failure (requires CPython OOM)
//   - callback.object: makeCfunction failure (requires CPython OOM)
//   - callbackCall: p==nil branch (requires passing a non-capsule PyObject
//     at the C level; calling callbackCall(nil,nil) directly causes SIGSEGV)
//   - callbackDestroy: p==nil branch (same reason as above)
//   - callbackDestroy: the function body itself is a CGo export attributed
//     to C by the coverage tool; Delete() is its Go-side body and is tested
//     directly by TestCallbackDeleteDirect.
//   - callback.call: len(ret)==2, err==nil fallthrough — triggers a panic
//     in callback.go due to nil interface assertion; not tested to avoid
//     masking the underlying bug without modifying callback.go.
func TestNewCallbackInvalidSignature(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// Two return values, second is NOT error.
	obj := py.NewObject(func() (int, int) { return 1, 2 })
	if obj.Err() == nil {
		t.Error("NewObject (int,int): expected error for invalid signature, got nil")
	}

	// Three return values.
	obj = py.NewObject(func() (int, int, int) { return 1, 2, 3 })
	if obj.Err() == nil {
		t.Error("NewObject (int,int,int): expected error for invalid signature, got nil")
	}
}

// TestNewCallbackValidSignatures tests all accepted return-value combinations.
func TestNewCallbackValidSignatures(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// 0 return values.
	obj := py.NewObject(func() {})
	if err := obj.Err(); err != nil {
		t.Errorf("0-return callback: unexpected error: %s", err)
	}

	// 1 return value.
	obj = py.NewObject(func() int { return 0 })
	if err := obj.Err(); err != nil {
		t.Errorf("1-return callback: unexpected error: %s", err)
	}

	// 2 return values (value + error).
	obj = py.NewObject(func() (int, error) { return 0, errors.New("x") })
	if err := obj.Err(); err != nil {
		t.Errorf("2-return (value+error) callback: unexpected error: %s", err)
	}
}

// TestCallbackVoidNoArgs covers the fast-path in callback.call:
// NumIn==0 && NumOut==0 — returns pyNone without entering the gate.
func TestCallbackVoidNoArgs(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	called := false
	obj := py.NewObject(func() { called = true })
	assert.NoError(obj.Err())

	ret := obj.Call()
	if err := ret.Err(); err != nil {
		t.Errorf("void callback: unexpected error: %s", err)
	}
	if !called {
		t.Error("void callback: function was not called")
	}
	if !ret.IsNone() {
		t.Errorf("void callback: expected None return, got %s", ret.String())
	}
}

// TestCallbackSingleReturn covers the len(ret)==1 branch in callback.call.
func TestCallbackSingleReturn(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(func() int { return 42 })
	assert.NoError(obj.Err())

	ret := obj.Call()
	if err := ret.Err(); err != nil {
		t.Errorf("single-return callback: unexpected error: %s", err)
	}

	v, err := ret.Int()
	if err != nil {
		t.Errorf("single-return callback: Int: %s", err)
	}
	if v != 42 {
		t.Errorf("single-return callback: expected 42, got %d", v)
	}
}

// TestCallbackValueAndError covers the len(ret)==2, err!=nil branch.
// The Go error propagates as a Python exception via callbackSetError.
func TestCallbackValueAndError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(func() (int, error) { return 0, errors.New("go error") })
	assert.NoError(obj.Err())

	ret := obj.Call()
	if ret.Err() == nil {
		t.Error("value+error callback: expected error to propagate, got nil")
	}
}

// TestCallbackSetErrorPython covers the ErrPython branch in callbackSetError.
func TestCallbackSetErrorPython(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	pyErr := ErrPython{
		except: Except("ValueError"),
		msg:    "python-style error",
	}
	obj := py.NewObject(func() (int, error) { return 0, pyErr })
	assert.NoError(obj.Err())

	ret := obj.Call()
	if ret.Err() == nil {
		t.Error("ErrPython callback: expected error to propagate, got nil")
	}
}

// TestCallbackRoundTrip registers a void Go callback as a Python global and
// calls it via py.Eval, verifying the complete call chain end-to-end.
func TestCallbackRoundTrip(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	called := false
	obj := py.NewObject(func() { called = true })
	assert.NoError(obj.Err())

	if err = obj.Save("go_callback"); err != nil {
		t.Fatalf("Save: %s", err)
	}

	ret := py.Eval("go_callback()")
	if err := ret.Err(); err != nil {
		t.Errorf("round-trip call: %s", err)
	}
	if !called {
		t.Error("round-trip call: Go function was not called")
	}
}

// TestCallbackRoundTripReturn verifies a callback returning a value can be
// called from Python and its return value retrieved.
func TestCallbackRoundTripReturn(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(func() int { return 7 })
	assert.NoError(obj.Err())

	if err = obj.Save("go_fn"); err != nil {
		t.Fatalf("Save: %s", err)
	}

	ret := py.Eval("go_fn()")
	if err := ret.Err(); err != nil {
		t.Errorf("round-trip return: %s", err)
	}

	v, err := ret.Int()
	if err != nil || v != 7 {
		t.Errorf("round-trip return: expected 7, got %v (err: %v)", v, err)
	}
}

// TestCallbackDelete covers callback.Delete and callbackDestroy.
// The callback object is a PyCFunction wrapping an internal PyCapsule.
// When the PyCFunction's refcount drops to zero, CPython frees it, which
// decrefs the capsule, firing callbackDestroy → cb.Delete().
// We delete the Python global and force a CPython GC cycle to ensure
// the destructor fires before the test returns.
func TestCallbackDelete(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(func() {})
	assert.NoError(obj.Err())

	if err = obj.Save("cb_to_delete"); err != nil {
		t.Fatalf("Save: %s", err)
	}

	// Drop the Go-side reference first.
	obj.Invalidate()

	// Delete the Python global and run CPython's cyclic GC to ensure
	// the PyCFunction and its internal capsule are freed synchronously,
	// firing callbackDestroy → cb.Delete().
	assert.NoError(py.Exec(`
import gc
del cb_to_delete
gc.collect()
`, ""))
}

// TestCallbackDeleteDirect covers callback.Delete directly, exercising
// the C.free calls for ml_name and the PyMethodDef allocation.
// (callbackDestroy is a CGo export attributed to C by the coverage tool;
// Delete() is its Go-side body and is what actually needs coverage.)
func TestCallbackDeleteDirect(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	cb := newCallback(py, "direct_delete", func() {})
	if cb == nil {
		t.Fatal("newCallback returned nil unexpectedly")
	}
	// Call Delete directly — must not panic or crash.
	cb.Delete()
}

// TestCallbackCallGateError covers the cb.py.gate() failure path in
// callback.call. This path is reached when the interpreter is closed
// between the time the callback object was created and when it is called.
// We simulate this by closing the interpreter and then directly invoking
// callback.call on a callback whose interpreter is gone.
func TestCallbackCallGateError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)

	// Create a callback with NumIn>0 or NumOut>0 so it does NOT take the
	// void fast-path and instead reaches the gate() call.
	// A func() int has NumOut==1, NumIn==0 — it skips the void path and
	// calls gate() before anything else.
	cb := newCallback(py, "test", func() int { return 1 })
	if cb == nil {
		t.Fatal("newCallback returned nil")
	}

	// Close the interpreter so gate() will return an error.
	py.Close()

	// Directly call cb.call — gate() must fail and return an error.
	_, err = cb.call(nil)
	if err == nil {
		t.Error("call after interpreter close: expected error, got nil")
	}
}

