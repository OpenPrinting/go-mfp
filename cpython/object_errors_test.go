// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests for error propagation through Object methods, conversion errors,
// and miscellaneous methods: Py, Save, SaveTo, SaveItem, String, Repr,
// TypeName, TypeModuleName.

package cpython

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestObjectMiscMethods tests Py, Save, SaveTo, SaveItem, String, Repr,
// TypeName, TypeModuleName, and error-object behaviour.
func TestObjectMiscMethods(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.Eval("42")
	assert.NoError(obj.Err())

	if obj.Py() != py {
		t.Errorf("Object.Py: returned wrong interpreter")
	}

	r, err := obj.Repr()
	if err != nil {
		t.Errorf("Object.Repr: %s", err)
	} else if r != "42" {
		t.Errorf("Object.Repr:\nexpected: 42\npresent:  %s", r)
	}

	if tn := obj.TypeName(); tn != "builtins.int" {
		t.Errorf("Object.TypeName:\nexpected: int\npresent:  %s", tn)
	}

	if tmn := obj.TypeModuleName(); tmn != "builtins" {
		t.Errorf("Object.TypeModuleName:\nexpected: builtins\npresent:  %s", tmn)
	}

	if s := obj.String(); s != "42" {
		t.Errorf("Object.String:\nexpected: 42\npresent:  %s", s)
	}

	// String() on an error object returns the error message.
	errObj := newErrorObject(py, ErrNotFound{})
	if s := errObj.String(); s != (ErrNotFound{}).Error() {
		t.Errorf("Object.String (error):\n"+
			"expected: %s\n"+
			"present:  %s",
			ErrNotFound{}.Error(), s)
	}

	// Save: store as a global variable.
	num := py.Eval("100")
	assert.NoError(num.Err())
	if err = num.Save("saved_num"); err != nil {
		t.Errorf("Object.Save: %s", err)
	}
	got := py.Eval("saved_num")
	assert.NoError(got.Err())
	if v, err := got.Int(); err != nil || v != 100 {
		t.Errorf("Object.Save: expected 100, got %v (err: %v)", v, err)
	}

	// SaveTo: store as an attribute of another object.
	err = py.Exec("class Box:\n    pass\nbox = Box()\n", "")
	if err != nil {
		t.Fatalf("Object.SaveTo setup: %s", err)
	}
	box := py.Eval("box")
	if err := box.Err(); err != nil {
		t.Fatalf("Object.SaveTo setup: Eval box: %s", err)
	}
	saveToVal := py.NewObject(999)
	assert.NoError(saveToVal.Err())
	if err = saveToVal.SaveTo(box, "contents"); err != nil {
		t.Errorf("Object.SaveTo: %s", err)
	}
	attr := box.Get("contents")
	if err := attr.Err(); err != nil {
		t.Errorf("Object.SaveTo: attribute not set: %s", err)
	} else if n, err := attr.Int(); err != nil || n != 999 {
		t.Errorf("Object.SaveTo: expected 999, got %v (err: %v)", n, err)
	}

	// SaveItem: store as an item in a dict.
	d := py.Eval("{}")
	assert.NoError(d.Err())
	saveItemVal := py.NewObject(777)
	assert.NoError(saveItemVal.Err())
	if err = saveItemVal.SaveItem(d, "key"); err != nil {
		t.Errorf("Object.SaveItem: %s", err)
	}
	gotItem := d.GetItem("key")
	if err := gotItem.Err(); err != nil {
		t.Errorf("Object.SaveItem: item not set: %s", err)
	} else if n, err := gotItem.Int(); err != nil || n != 777 {
		t.Errorf("Object.SaveItem: expected 777, got %v (err: %v)", n, err)
	}

	// Error-object: type/callable/seq checks return zero values.
	e := newErrorObject(py, ErrNotFound{})
	if e.IsNone() {
		t.Error("Object.IsNone on error obj: expected false")
	}
	if e.TypeName() != "" {
		t.Error("Object.TypeName on error obj: expected empty string")
	}
	if e.TypeModuleName() != "" {
		t.Error("Object.TypeModuleName on error obj: expected empty string")
	}
	if e.IsCallable() {
		t.Error("Object.IsCallable on error obj: expected false")
	}
	if e.IsSeq() {
		t.Error("Object.IsSeq on error obj: expected false")
	}
}

// TestObjectErrorPaths tests that every Object method propagates an error
// object's error unchanged, and covers begin() failure paths.
func TestObjectErrorPaths(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	sentinel := ErrNotFound{}
	e := newErrorObject(py, sentinel)

	// Invalidate on a closed interpreter hits the gate() error path.
	py2, err := NewPython()
	assert.NoError(err)
	obj2 := py2.Eval("1")
	assert.NoError(obj2.Err())
	py2.Close()
	obj2.Invalidate()

	check := func(name string, got error) {
		t.Helper()
		if got != sentinel {
			t.Errorf("%s: expected sentinel error, got %v", name, got)
		}
	}

	_, err = e.Del("key")
	check("Del", err)

	if got := e.GetItem("key"); got.Err() != sentinel {
		t.Errorf("GetItem: expected sentinel error, got %v", got.Err())
	}

	_, err = e.ContainsItem("key")
	check("ContainsItem", err)

	check("SetItem", e.SetItem("key", 1))

	_, err = e.DelAttr("name")
	check("DelAttr", err)

	if got := e.Get("name"); got.Err() != sentinel {
		t.Errorf("Get: expected sentinel error, got %v", got.Err())
	}

	_, err = e.HasAttr("name")
	check("HasAttr", err)

	check("Set", e.Set("name", 1))

	if got := e.CallKW(nil, 1, 2); got.Err() != sentinel {
		t.Errorf("CallKW: expected sentinel error, got %v", got.Err())
	}

	_, err = e.Bool()
	check("Bool", err)

	_, err = e.Keys()
	check("Keys", err)

	_, err = e.Slice()
	check("Slice", err)

	// fastBool on a plain integer (not True/False) must return an error.
	obj := py.Eval("5")
	assert.NoError(obj.Err())
	if _, err = obj.fastBool(); err == nil {
		t.Error("fastBool on int: expected error, got nil")
	}
}

// TestObjectBeginClosed tests begin() after the interpreter has been closed.
func TestObjectBeginClosed(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)

	obj := py.NewObject(42)
	assert.NoError(obj.Err())
	py.Close()

	_, err = obj.Int()
	if err == nil {
		t.Error("begin: expected error after interpreter closed, got nil")
	}
}

// TestObjectBeginInvalidatedOID tests begin() after Invalidate removes the
// oid from the interpreter's map; lookupObjID returns nil and must error.
func TestObjectBeginInvalidatedOID(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.Eval("999")
	assert.NoError(obj.Err())
	obj.Invalidate()

	_, err = obj.Int()
	if err == nil {
		t.Error("begin after oid removal: expected error, got nil")
	}

	if s := obj.String(); s == "" {
		t.Error("String after oid removal: expected non-empty error string")
	}
}

// TestObjectGetErrorPropagation tests that Get on an error object propagates
// the error and NotFound() reports correctly.
func TestObjectGetErrorPropagation(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	errObj := newErrorObject(py, ErrNotFound{})
	result := errObj.Get("anything")
	if result.Err() == nil {
		t.Error("Get on error object: expected error propagation")
	}
	if !result.NotFound() {
		t.Errorf("Get on error object: expected ErrNotFound, got %v", result.Err())
	}
}

// TestObjectDelKeyConversionError tests Del when the key cannot be converted.
func TestObjectDelKeyConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`d = {}`, ""))
	obj := py.Get("d")
	assert.NoError(obj.Err())

	_, err = obj.Del(make(chan int))
	if err == nil {
		t.Error("Del: expected error for unconvertible key")
	}
}

// TestObjectGetItemKeyConversionError tests GetItem when the key cannot be converted.
func TestObjectGetItemKeyConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`d = {"a": 1}`, ""))
	obj := py.Get("d")
	assert.NoError(obj.Err())

	if result := obj.GetItem(make(chan int)); result.Err() == nil {
		t.Error("GetItem: expected error for unconvertible key")
	}
}

// TestObjectContainsItemKeyConversionError tests ContainsItem when the key
// cannot be converted.
func TestObjectContainsItemKeyConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`d = {"a": 1}`, ""))
	obj := py.Get("d")
	assert.NoError(obj.Err())

	_, err = obj.ContainsItem(make(chan int))
	if err == nil {
		t.Error("ContainsItem: expected error for unconvertible key")
	}
}

// TestObjectSetItemKeyConversionError tests SetItem when the key cannot be converted.
func TestObjectSetItemKeyConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`d = {}`, ""))
	obj := py.Get("d")
	assert.NoError(obj.Err())

	if err = obj.SetItem(make(chan int), 1); err == nil {
		t.Error("SetItem: expected error for unconvertible key")
	}
}

// TestObjectSetItemValConversionError tests SetItem when the value cannot be converted.
func TestObjectSetItemValConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`d = {}`, ""))
	obj := py.Get("d")
	assert.NoError(obj.Err())

	if err = obj.SetItem("key", make(chan int)); err == nil {
		t.Error("SetItem: expected error for unconvertible value")
	}
}

// TestObjectSetValConversionError tests Set when the value cannot be converted.
func TestObjectSetValConversionError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`
class Obj:
    pass
o = Obj()
`, ""))
	obj := py.Get("o")
	assert.NoError(obj.Err())

	if err = obj.Set("attr", make(chan int)); err == nil {
		t.Error("Set: expected error for unconvertible value")
	}
}
