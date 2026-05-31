// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests for Bool/fastBool paths and object lifecycle:
// GC collection, Invalidate, finalizer after interpreter close.

package cpython

import (
	"runtime"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestObjectBoolFallback tests Bool() via the __bool__ method fallback path.
func TestObjectBoolFallback(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	script := `
class Truthy:
    def __bool__(self):
        return True

class Falsy:
    def __bool__(self):
        return False

truthy = Truthy()
falsy  = Falsy()
`
	err = py.Exec(script, "")
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}

	obj := py.Eval("truthy")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval truthy: %s", err)
	}
	b, err := obj.Bool()
	if err != nil {
		t.Errorf("Object.Bool (Truthy): %s", err)
	} else if !b {
		t.Errorf("Object.Bool (Truthy): expected true, got false")
	}

	obj = py.Eval("falsy")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval falsy: %s", err)
	}
	b, err = obj.Bool()
	if err != nil {
		t.Errorf("Object.Bool (Falsy): %s", err)
	} else if b {
		t.Errorf("Object.Bool (Falsy): expected false, got true")
	}
}

// TestObjectBoolNoMethod tests Bool() when __bool__ exists but is not callable
// (set to None). fastBool fails, toBool.Call() also fails, so Bool() must
// return the original fastBool error.
func TestObjectBoolNoMethod(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	script := `
class NoBool:
    __bool__ = None

nobool = NoBool()
`
	err = py.Exec(script, "")
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}

	obj := py.Eval("nobool")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval nobool: %s", err)
	}

	_, err = obj.Bool()
	if err == nil {
		t.Error("Object.Bool (NoBool): expected error, got nil")
	}
}

// TestObjectGC tests that Objects are properly tracked and collected by GC.
func TestObjectGC(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)

	runtime.GC()
	base := py.countObjID()

	err = py.Eval("5").Err()
	assert.NoError(err)

	if len(py.objects.mapped) != base+1 {
		t.Errorf("TestObjectGC: object not properly mapped")
	}

	runtime.GC()
	runtime.GC()

	if py.countObjID() != base {
		t.Errorf("TestObjectGC: GC did not collect object")
	}
}

// TestObjectInvalidateValid tests Invalidate on a live object.
// Verifies the gate-success branch of Invalidate, and that any subsequent
// operation on the invalidated object returns an error.
func TestObjectInvalidateValid(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.Eval("123")
	assert.NoError(obj.Err())

	v, err := obj.Int()
	if err != nil || v != 123 {
		t.Fatalf("pre-Invalidate: expected 123, got %v (err: %v)", v, err)
	}

	obj.Invalidate()

	_, err = obj.Int()
	if err == nil {
		t.Error("Int after Invalidate: expected error, got nil")
	}
}

// TestObjectFinalizerAfterClose tests that the GC finalizer handles the case
// where the Python interpreter is already closed when it fires.
// This covers the obj.py.closed() == true branch inside finalizer().
func TestObjectFinalizerAfterClose(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)

	{
		obj := py.Eval("42")
		assert.NoError(obj.Err())
		py.Close()
		_ = obj
	}

	// Two GC passes: first enqueues finalizers, second runs them.
	runtime.GC()
	runtime.GC()
	// Reaching here without panic means the branch is covered.
}

