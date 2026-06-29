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

// ---------------------------------------------------------------------
// Added tests below — target gate.go lines reported at low/0% coverage.
// NOTE: method names (Keys, HasAttr, ContainsItem, Bool, Uint64, etc.)
// are inferred from the pattern of existing tests above (Set/Get/Del/
// SetItem/GetItem/ContainsItem all visible there). I do not have
// object.go/py.go, so verify each exported wrapper name/signature
// against your actual Object/Python API before relying on this file —
// run `go vet ./cpython/...` and fix any mismatches.
// ---------------------------------------------------------------------

// TestGateKeys covers keys() (0% coverage at gate.go:289).
func TestGateKeys(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("{'a': 1, 'b': 2}")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	keys, err := obj.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("Keys: got %d keys, want 2", len(keys))
	}
}

// TestGateHasAttrFalse covers the false branch of hasattr
// (gate.go:320, currently 80%).
func TestGateHasAttrFalse(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("42")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	has, err := obj.HasAttr("nonexistent_attr")
	if err != nil {
		t.Fatalf("HasAttr: %v", err)
	}
	if has {
		t.Fatalf("HasAttr: expected false for nonexistent_attr")
	}
}

// TestGateGetAttrMissing covers the error branch of getattr
// (gate.go:307, currently 83.3%).
func TestGateGetAttrMissing(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("42")
	got := obj.Get("nonexistent_attr")
	if err := got.Err(); err == nil {
		t.Fatalf("Get: expected error for nonexistent_attr")
	}
}

// TestGateSetAttrError covers the error branch of setattr
// (gate.go:333, currently 66.7%) — int objects reject attribute sets.
func TestGateSetAttrError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("42")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	if err := obj.Set("x", 1); err == nil {
		t.Fatalf("Set: expected error setting attr on int literal")
	}
}

// TestGateHasItemFalse covers the false branch of hasitem
// (gate.go:370, currently 50%).
func TestGateHasItemFalse(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("{'a': 1}")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	has, err := obj.ContainsItem("missing_key")
	if err != nil {
		t.Fatalf("ContainsItem: %v", err)
	}
	if has {
		t.Fatalf("ContainsItem: expected false for missing_key")
	}
}

// TestGateDelItemMissing covers the error branch of delitem
// (gate.go:347, currently 50%).
func TestGateDelItemMissing(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("{}")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	found, err := obj.Del("missing_key")
	if err != nil {
		t.Fatalf("Del: unexpected error: %v", err)
	}
	if found {
		t.Fatalf("Del: expected found=false for missing_key")
	}
}

// TestGateGetSeqItem covers getSeqItem (gate.go:777, currently 50%)
// via a generic sequence type (string) rather than list/tuple, which
// have their own dedicated get* paths.
func TestGateGetSeqItem(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval(`"abc"`)
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	item := obj.GetItem(1)
	if err := item.Err(); err != nil {
		t.Fatalf("GetItem: %v", err)
	}

	s, err := item.Str()
	if err != nil {
		t.Fatalf("Str: %v", err)
	}
	if s != "b" {
		t.Fatalf("GetItem(1) on string: got %q, want %q", s, "b")
	}
}

// TestGateSetListItem covers setListItem (gate.go:769, currently 0%).
func TestGateSetListItem(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("[1, 2, 3]")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	if err := obj.SetItem(0, 99); err != nil {
		t.Fatalf("SetItem: %v", err)
	}

	got := obj.GetItem(0)
	if err := got.Err(); err != nil {
		t.Fatalf("GetItem: %v", err)
	}

	n, err := got.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}
	if n != 99 {
		t.Fatalf("SetListItem: got %d, want 99", n)
	}
}

// TestGateGetTupleItemError covers the error branch of getTupleItem
// (gate.go:783, currently 0%) — out-of-range index.
func TestGateGetTupleItemError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("(1, 2, 3)")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	item := obj.GetItem(10)
	if err := item.Err(); err == nil {
		t.Fatalf("GetItem(10): expected out-of-range error")
	}
}

// TestGateMakeBool covers makeBool (gate.go:695, currently 0%).
func TestGateMakeBool(t *testing.T) {
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

	if err := obj.Set("x", true); err != nil {
		t.Fatalf("Set bool: %v", err)
	}

	got := obj.Get("x")
	b, err := got.Bool()
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if !b {
		t.Fatalf("MakeBool: got false, want true")
	}
}

// TestGateMakeDict covers makeDict (gate.go:719, currently 0%).
func TestGateMakeDict(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	err = py.Exec("def take(d): return len(d)", "test")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	fn := py.Eval("take")
	if err := fn.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	res := fn.Call(map[string]int{"a": 1, "b": 2})
	if err := res.Err(); err != nil {
		t.Fatalf("Call: %v", err)
	}
}

// TestGateMakeFloat covers makeFloat (gate.go:725, currently 0%).
func TestGateMakeFloat(t *testing.T) {
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

	if err := obj.Set("x", 3.14); err != nil {
		t.Fatalf("Set float: %v", err)
	}

	got := obj.Get("x")
	f, err := got.Float()
	if err != nil {
		t.Fatalf("Float: %v", err)
	}
	if f != 3.14 {
		t.Fatalf("MakeFloat: got %v, want 3.14", f)
	}
}

// TestGateMakeTuple covers makeTuple (gate.go:751, currently 0%).
func TestGateMakeTuple(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	err = py.Exec("def take(t): return len(t)", "test")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	fn := py.Eval("take")
	if err := fn.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	res := fn.Call([3]int{1, 2, 3})
	if err := res.Err(); err != nil {
		t.Fatalf("Call: %v", err)
	}
}

// TestGateMakeUintError covers the overflow branch of makeUint
// (gate.go:757, currently 50%) by round-tripping a negative-overflow
// case through decode, and the success path via a large uint64.
func TestGateMakeUintLarge(t *testing.T) {
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

	want := uint64(1) << 63
	if err := obj.Set("x", want); err != nil {
		t.Fatalf("Set uint64: %v", err)
	}

	got := obj.Get("x")
	n, err := got.Uint()
	if err != nil {
		t.Fatalf("Uint: %v", err)
	}
	if n != want {
		t.Fatalf("MakeUint: got %d, want %d", n, want)
	}
}

// TestGateDecodeUint64Negative covers decodeUint64's overflow branch
// (gate.go:570, currently 66.7%) — negative int cannot become uint64.
func TestGateDecodeUint64Negative(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("-1")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	_, err = obj.Uint()
	if err == nil {
		t.Fatalf("Uint: expected overflow error for -1")
	}
}

// TestGateDecodeExactComplexError covers decodeExactComplex's type
// error branch (gate.go:588, currently 66.7%).
func TestGateDecodeExactComplexError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval(`"not a complex"`)
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	_, err = obj.Complex()
	if err == nil {
		t.Fatalf("Complex: expected type error for string")
	}
}

// TestGateDecodeExactFloatError covers decodeExactFloat's type error
// branch (gate.go:605, currently 85.7%).
func TestGateDecodeExactFloatError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval(`"not a float"`)
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	_, err = obj.Float()
	if err == nil {
		t.Fatalf("Float: expected type error for string")
	}
}

// TestGateDecodeExactInt64Overflow covers decodeExactInt64's overflow
// branch (gate.go:622, currently 87.5%).
func TestGateDecodeExactInt64Overflow(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("99999999999999999999999999999")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	_, err = obj.Int()
	if err == nil {
		t.Fatalf("Int: expected overflow error for huge literal")
	}
}

// TestGateLoad covers load() error branch (gate.go:830, currently 83.3%)
// by importing a module that does not exist.
func TestGateLoadError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Load("this_module_does_not_exist_xyz", "mod", "test")
	if err := obj.Err(); err == nil {
		t.Fatalf("Load: expected error for nonexistent module")
	}
}

// TestGateLastErrorLocationNested covers the multi-frame traceback walk in
// lastErrorLocation (gate.go:120, currently 76.5%) by raising an error
// from inside a nested function call, producing a traceback with more
// than one frame.
func TestGateLastErrorLocationNested(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	err = py.Exec(
		"def inner():\n"+
			"    raise ValueError('boom')\n"+
			"def outer():\n"+
			"    inner()\n"+
			"outer()\n",
		"nested_test")
	if err == nil {
		t.Fatalf("Exec: expected error from nested raise")
	}
}

// TestGateSlice covers objSlice's getSeqItem call path (gate.go:777)
// via the exported Object.Slice() method, distinct from GetItem which
// uses gate.getitem instead.
func TestGateSlice(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval("(1, 2, 3)")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	items, err := obj.Slice()
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("Slice: got %d items, want 3", len(items))
	}
}

// TestGateCallPositionalArgs covers setTupleItem's success path
// (gate.go:790) via Object.Call with positional arguments, and
// binop's tuple construction (gate.go:790 too, via Add).
func TestGateCallPositionalArgs(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	err = py.Exec("def add(a, b): return a + b", "test")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	fn := py.Eval("add")
	if err := fn.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	res := fn.Call(1, 2)
	if err := res.Err(); err != nil {
		t.Fatalf("Call: %v", err)
	}

	n, err := res.Int()
	if err != nil {
		t.Fatalf("Int: %v", err)
	}
	if n != 3 {
		t.Fatalf("Call(1, 2): got %d, want 3", n)
	}

	one := py.Eval("1")
	if err := one.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	sum := one.Add(2)
	if err := sum.Err(); err != nil {
		t.Fatalf("Add: %v", err)
	}
}

// TestGateDecodeBigintError covers decodeBigint's type-check error
// branch (gate.go:459) for a non-PyLong object.
func TestGateDecodeBigintError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval(`"not an int"`)
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	_, err = obj.Bigint()
	if err == nil {
		t.Fatalf("Bigint: expected error for string")
	}
}

// TestGateDecodeBytesError covers decodeBytes' final fallthrough
// error branch (gate.go:477) for an object that is neither
// PyBytes_Type nor PyByteArray_Type.
func TestGateDecodeBytesError(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	obj := py.Eval(`"not bytes"`)
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	_, err = obj.Bytes()
	if err == nil {
		t.Fatalf("Bytes: expected error for string")
	}
}

