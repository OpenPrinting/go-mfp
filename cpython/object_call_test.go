// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests for callable detection and invocation: IsCallable, Call, CallKW.

package cpython

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestObjectCallable tests IsCallable on callable and non-callable objects.
func TestObjectCallable(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.Eval("5")
	assert.NoError(obj.Err())
	if obj.IsCallable() {
		t.Errorf("Object.IsCallable: false positive on integer")
	}

	obj = py.Eval("min")
	assert.NoError(obj.Err())
	if !obj.IsCallable() {
		t.Errorf("Object.IsCallable: false negative on builtin function")
	}
}

// TestObjectCall tests Call with positional args and CallKW with keyword args.
func TestObjectCall(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.Eval("min")
	assert.NoError(obj.Err())

	// Positional args: min(1, 2) == 1
	res := obj.Call(1, 2)
	if err := res.Err(); err != nil {
		t.Errorf("Object.Call (positional): %s", err)
		return
	}
	val, err := res.Int()
	if err != nil {
		t.Errorf("Object.Call (positional): Int: %s", err)
		return
	}
	if val != 1 {
		t.Errorf("Object.Call (positional):\nexpected: 1\npresent:  %d\n", val)
	}

	// Keyword args: min([], default=5) == 5
	res = obj.CallKW(map[string]any{"default": 5}, []int{})
	if err := res.Err(); err != nil {
		t.Errorf("Object.Call (keyword): %s", err)
		return
	}
	val, err = res.Int()
	if err != nil {
		t.Errorf("Object.Call (keyword): Int: %s", err)
		return
	}
	if val != 5 {
		t.Errorf("Object.Call (keyword):\nexpected: 5\npresent:  %d\n", val)
	}
}

// TestObjectCallKWCoverage tests CallKW with an empty (non-nil) kw map,
// which exercises the len(kw)==0 branch where pykwargs stays nil.
func TestObjectCallKWCoverage(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`def f(*args, **kwargs): return len(args)`, ""))
	fn := py.Get("f")
	assert.NoError(fn.Err())

	res := fn.CallKW(map[string]any{}, 1, 2, 3)
	if err := res.Err(); err != nil {
		t.Errorf("CallKW empty kw: unexpected error: %s", err)
	}
	v, err := res.Int()
	if err != nil || v != 3 {
		t.Errorf("CallKW empty kw: expected 3, got %v (err: %v)", v, err)
	}
}

// TestObjectCallKWArgConversionError tests CallKW when a positional arg
// cannot be converted to a Python object.
func TestObjectCallKWArgConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`def f(x): return x`, ""))
	fn := py.Get("f")
	assert.NoError(fn.Err())

	result := fn.Call(make(chan int))
	if result.Err() == nil {
		t.Error("CallKW: expected error for unconvertible positional arg")
	}
}

// TestObjectCallKWKwargConversionError tests CallKW when a keyword arg value
// cannot be converted to a Python object.
func TestObjectCallKWKwargConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`def f(x=None): return x`, ""))
	fn := py.Get("f")
	assert.NoError(fn.Err())

	result := fn.CallKW(map[string]any{"x": make(chan int)})
	if result.Err() == nil {
		t.Error("CallKW: expected error for unconvertible keyword arg")
	}
}
