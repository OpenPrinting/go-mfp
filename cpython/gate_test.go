// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Tests for pyGate operations

package cpython

import (
	"testing"
)

// TestGateSetAttr verifies that Set correctly sets an object attribute.
func TestGateSetAttr(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	err = py.Exec("class Obj: pass\nobj = Obj()", "test")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	obj := py.Eval("obj")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	if err := obj.Set("x", 42); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got := obj.Get("x")
	if err := got.Err(); err != nil {
		t.Fatalf("Get: %v", err)
	}

	n, err := got.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}

	if n != 42 {
		t.Fatalf("Set/Get attr: got %d, want 42", n)
	}
}

// TestGateTypeModuleName verifies that TypeModuleName returns
// a non-empty string for standard Python objects.
func TestGateTypeModuleName(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("42")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	name := obj.TypeModuleName()
	if name == "" {
		t.Fatalf("TypeModuleName() returned empty string")
	}
}

// TestGateGetListItem verifies that list item retrieval works correctly.
func TestGateGetListItem(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("[10, 20, 30]")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	item := obj.GetItem(1)
	if err := item.Err(); err != nil {
		t.Fatalf("GetItem: %v", err)
	}

	n, err := item.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}

	if n != 20 {
		t.Fatalf("GetItem(1): got %d, want 20", n)
	}
}

// TestGateGetTupleItem verifies that tuple item retrieval works correctly.
func TestGateGetTupleItem(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("(10, 20, 30)")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	item := obj.GetItem(2)
	if err := item.Err(); err != nil {
		t.Fatalf("GetItem: %v", err)
	}

	n, err := item.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}

	if n != 30 {
		t.Fatalf("GetItem(2): got %d, want 30", n)
	}
}

// TestGateDelItem verifies that dict item deletion works correctly.
func TestGateDelItem(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("{'a': 1, 'b': 2}")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	found, err := obj.Del("a")
	if err != nil {
		t.Fatalf("Del: %v", err)
	}
	if !found {
		t.Fatalf("Del: key 'a' not found")
	}

	has, err := obj.ContainsItem("a")
	if err != nil {
		t.Fatalf("ContainsItem: %v", err)
	}
	if has {
		t.Fatalf("Del: key 'a' still present after deletion")
	}
}

// TestGateSetItem verifies that dict item setting works correctly.
func TestGateSetItem(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("{}")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	if err := obj.SetItem("mykey", 99); err != nil {
		t.Fatalf("SetItem: %v", err)
	}

	got := obj.GetItem("mykey")
	if err := got.Err(); err != nil {
		t.Fatalf("GetItem: %v", err)
	}

	n, err := got.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}

	if n != 99 {
		t.Fatalf("SetItem: got %d, want 99", n)
	}
}
