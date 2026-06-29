// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests for attribute access (Get, Set, Del, HasAttr) and container
// operations (GetItem, SetItem, Del, ContainsItem, Len, Slice, Keys).

package cpython

import (
	"errors"
	"fmt"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
)

// TestObjectAttributes tests Get, Set, DelAttr, HasAttr operations.
func TestObjectAttributes(t *testing.T) {
	script := `
class Dog:
    species = "Canis familiaris"

    def __init__(self, name, age):
        self.name = name
        self.age = age

dog = Dog("Archi", 4)
`
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	err = py.Exec(script, "")
	assert.NoError(err)

	obj := py.Eval("dog")
	assert.NoError(obj.Err())

	// Attribute "name" must exist for now
	found, err := obj.HasAttr("name")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !found {
		t.Errorf("Attribute %q expected but not found", "name")
	}

	// And "unknown" attribute must not exist
	found, err = obj.HasAttr("unknown")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if found {
		t.Errorf("Attribute %q not expected but found", "unknown")
	}

	// Check the returned attribute value
	attr := obj.Get("name")
	if err := attr.Err(); err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else {
		s, err := attr.Unicode()
		assert.NoError(err)
		if s != "Archi" {
			t.Errorf("obj.Get mismatch:\nexpected: %s\npresent:  %s\n", "Archi", s)
		}
	}

	// Attempt to delete "name" must succeed
	deleted, err := obj.DelAttr("name")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !deleted {
		t.Errorf("obj.DelAttr: attribute %s not deleted", "name")
	}

	// Now attribute must not exist
	found, err = obj.HasAttr("name")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if found {
		t.Errorf("Deleted attribute %q still present", "name")
	}

	// And attempt to delete "unknown" attribute must fail
	deleted, err = obj.DelAttr("unknown")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if deleted {
		t.Errorf("obj.DelAttr: attribute %s is deleted", "unknown")
	}
}

// TestObjectGetHasAttrError tests Get() when hasattr itself returns an error.
// A class with a raising __getattribute__ causes hasattr to propagate the
// error rather than returning (false, nil), covering the found=false, err!=nil
// branch in Get().
func TestObjectGetHasAttrError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	script := `
class Exploding:
    def __getattribute__(self, name):
        raise RuntimeError("no attrs for you")

class Empty:
    pass
`
	err = py.Exec(script, "")
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}

	// Access to non-existing attribute must return ErrNotFound{}
	obj := py.Eval("Empty()")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval exploding: %s", err)
	}

	result := obj.Get("anything")
	if !errors.Is(result.Err(), ErrNotFound{}) {
		t.Errorf("Get on Empty object:\n"+
			"expected: %s\n"+
			"present:  %s\n",
			result, ErrNotFound{})
	}

	// However, if object's __getattribute__ method throws an
	// exception, the error will be different...
	obj = py.Eval("Exploding()")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval exploding: %s", err)
	}

	result = obj.Get("anything")
	if !errors.Is(result.Err(), RuntimeError) {
		t.Errorf("Get on Exploding object:\n"+
			"expected: %s\n"+
			"present:  %s\n",
			result, RuntimeError)
	}
}

// TestObjectItems tests GetItem, SetItem, ContainsItem, Del on a dict.
func TestObjectItems(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.Eval(`{1:"1", 2:"2", 3:"3"}`)
	assert.NoError(obj.Err())

	// Items 1-3 must exist with correct values.
	for i := 1; i <= 3; i++ {
		found, err := obj.ContainsItem(i)
		if err != nil {
			t.Errorf("Object.ContainsItem(%v): %s", i, err)
			return
		}
		if !found {
			t.Errorf("Object.ContainsItem(%v): item not found", i)
		}

		item := obj.GetItem(i)
		if err := item.Err(); err != nil {
			t.Errorf("Object.GetItem(%v): %s", i, err)
			return
		}

		s, err := item.Str()
		assert.NoError(err)
		if expected := fmt.Sprintf("%d", i); s != expected {
			t.Errorf("Object.GetItem(%v):\nexpected: %s\npresent:  %s\n", i, expected, s)
		}
	}

	// Items 4-6 must not exist; GetItem must return ErrNotFound.
	for i := 4; i <= 6; i++ {
		found, err := obj.ContainsItem(i)
		if err != nil {
			t.Errorf("Object.ContainsItem(%v): %s", i, err)
			return
		}
		if found {
			t.Errorf("Object.ContainsItem(%v): item found unexpectedly", i)
		}

		item := obj.GetItem(i)
		if err := item.Err(); !errors.Is(err, ErrNotFound{}) {
			t.Errorf("Object.GetItem(%v):\n"+
				"expected: (%s)\n"+
				"present:  (%s)\n",
				i, ErrNotFound{}, err)
		}
	}

	// Add items 4-6 and verify all 6 are now present.
	for i := 4; i <= 6; i++ {
		if err := obj.SetItem(i, fmt.Sprintf("%d", i)); err != nil {
			t.Errorf("Object.SetItem(%v): %s", i, err)
			return
		}
	}
	for i := 1; i <= 6; i++ {
		item := obj.GetItem(i)
		if err := item.Err(); err != nil {
			t.Errorf("Object.GetItem(%v): %s", i, err)
			return
		}
		s, err := item.Str()
		assert.NoError(err)
		if expected := fmt.Sprintf("%d", i); s != expected {
			t.Errorf("Object.GetItem(%v):\nexpected: %s\npresent:  %s\n", i, expected, s)
		}
	}

	// Delete items 1-3; second delete must report not-found.
	for i := 1; i <= 3; i++ {
		found, err := obj.Del(i)
		if err != nil {
			t.Errorf("Object.Del(%v): %s", i, err)
			return
		}
		if !found {
			t.Errorf("Object.Del(%v): item not found on first delete", i)
		}

		found, err = obj.Del(i)
		if err != nil {
			t.Errorf("Object.Del(%v): %s", i, err)
			return
		}
		if found {
			t.Errorf("Object.Del(%v): item still found after delete", i)
		}
	}

	// Items 1-3 must now be absent.
	for i := 1; i <= 3; i++ {
		found, err := obj.ContainsItem(i)
		if err != nil {
			t.Errorf("Object.ContainsItem(%v): %s", i, err)
			return
		}
		if found {
			t.Errorf("Object.ContainsItem(%v): deleted item still present", i)
		}
	}
}

// TestObjectLen tests Object.Len on various container and scalar types.
func TestObjectLen(t *testing.T) {
	type testData struct {
		expr string
		l    int
		err  bool
	}

	tests := []testData{
		{expr: `[1,2,3]`, l: 3},
		{expr: `(1,2,3)`, l: 3},
		{expr: `{1:"1", 2:"2", 3:"3"}`, l: 3},
		{expr: `[]`, l: 0},
		{expr: `()`, l: 0},
		{expr: `{}`, l: 0},
		{expr: `5`, l: 0, err: true},
		{expr: `"hello"`, l: 5},
		{expr: `"привет"`, l: 6},
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj := py.Eval(test.expr)
		assert.NoError(obj.Err())

		l, err := obj.Len()
		switch {
		case err == nil && test.err:
			t.Errorf("Object.Len(%s): error not occurred", test.expr)
		case err != nil && !test.err:
			t.Errorf("Object.Len(%s): %s", test.expr, err)
		case l != test.l:
			t.Errorf("Object.Len(%s):\nexpected: %d\npresent:  %d\n",
				test.expr, test.l, l)
		}
	}
}

// TestObjectSlice tests Object.Slice on lists, tuples, and a non-sequence.
func TestObjectSlice(t *testing.T) {
	type testData struct {
		expr     string
		expected []string
		mustfail bool
	}

	tests := []testData{
		{expr: "()", expected: []string{}},
		{expr: "(1,2,3)", expected: []string{"1", "2", "3"}},
		{expr: "[]", expected: []string{}},
		{expr: "[1,2,3]", expected: []string{"1", "2", "3"}},
		{expr: "5", mustfail: true},
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj := py.Eval(test.expr)
		assert.NoError(obj.Err())

		slice, err := obj.Slice()
		if err != nil {
			if !test.mustfail {
				t.Errorf("%s: Object.Slice: %s", test.expr, err)
			}
			continue
		}
		if test.mustfail {
			t.Errorf("%s: Object.Slice: expected error didn't occur", test.expr)
			continue
		}

		result := make([]string, len(slice))
		for i := range slice {
			result[i], err = slice[i].Str()
			assert.NoError(err)
		}
		if diff := testutils.Diff(test.expected, result); diff != "" {
			t.Errorf("%s: Object.Slice:\n%s", test.expr, diff)
		}
	}
}

// TestObjectKeys tests Object.Keys on dicts and non-mapping types.
func TestObjectKeys(t *testing.T) {
	type testData struct {
		expr     string
		expected []string
		mustfail bool
	}

	tests := []testData{
		{expr: `{1:"one",2:"two",3:"three"}`, expected: []string{"1", "2", "3"}},
		{expr: `{}`, expected: []string{}},
		{expr: "()", mustfail: true},
		{expr: "5", mustfail: true},
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj := py.Eval(test.expr)
		assert.NoError(obj.Err())

		slice, err := obj.Keys()
		if err != nil {
			if !test.mustfail {
				t.Errorf("%s: Object.Keys: %s", test.expr, err)
			}
			continue
		}
		if test.mustfail {
			t.Errorf("%s: Object.Keys: expected error didn't occur", test.expr)
			continue
		}

		result := make([]string, len(slice))
		for i := range slice {
			result[i], err = slice[i].Str()
			assert.NoError(err)
		}
		if diff := testutils.Diff(test.expected, result); diff != "" {
			t.Errorf("%s: Object.Keys:\n%s", test.expr, diff)
		}
	}
}

// TestObjSliceNonSequence tests objSlice rejects a non-sequence object.
func TestObjSliceNonSequence(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(42)
	assert.NoError(obj.Err())

	_, err = obj.Slice()
	if err == nil {
		t.Error("Slice: expected error for non-sequence object")
	}
}

// TestObjSliceLengthError covers the gate.length failure path in objSlice.
// A list subclass whose __len__ raises passes isSeq() but fails at
// gate.length(), hitting the "return nil, err" branch.
func TestObjSliceLengthError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	script := `
class BadLen(list):
    def __len__(self):
        raise RuntimeError("len exploded")

badlen = BadLen([1, 2, 3])
`
	err = py.Exec(script, "")
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}

	obj := py.Eval("badlen")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval badlen: %s", err)
	}

	_, err = obj.Slice()
	if err == nil {
		t.Error("Slice on BadLen: expected error from __len__, got nil")
	}
}

// TestObjSliceGetSeqItemError tests the mid-loop cleanup path in objSlice.
// A list subclass that raises on index >= 1 causes objSlice to fetch item[0]
// successfully, then fail, triggering the unref-cleanup loop.
func TestObjSliceGetSeqItemError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	script := `
class BrokenList(list):
    def __getitem__(self, idx):
        if idx >= 1:
            raise IndexError("boom at index 1")
        return super().__getitem__(idx)

broken = BrokenList([10, 20, 30])
`
	err = py.Exec(script, "")
	if err != nil {
		t.Fatalf("setup failed: %s", err)
	}

	obj := py.Eval("broken")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval broken: %s", err)
	}

	_, err = obj.Slice()
	if err == nil {
		t.Error("Slice on BrokenList: expected error, got nil")
	}
}

// TestObjSliceGetItemErrorPath verifies Keys() works correctly on a populated
// dict, exercising the objSlice path via gate.keys().
func TestObjSliceGetItemErrorPath(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	assert.NoError(py.Exec(`d = {"a": 1, "b": 2, "c": 3}`, ""))
	obj := py.Get("d")
	assert.NoError(obj.Err())

	keys, err := obj.Keys()
	if err != nil {
		t.Fatalf("Keys: unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("Keys: expected 3 keys, got %d", len(keys))
	}
}
