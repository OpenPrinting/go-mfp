// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Binary operations on objects tests

package cpython

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestBinops tests binary operations on Objects
func TestBinops(t *testing.T) {
	type testData struct {
		name     string                     // Operation name
		in1, in2 any                        // Operands
		op       func(*Object, any) *Object // Operation
		out      any                        // Expected output
		err      string                     // Expected error
	}

	tests := []testData{
		{
			name: "+",
			in1:  1,
			in2:  2,
			op:   (*Object).Add,
			out:  3,
		},

		{
			name: "-",
			in1:  1,
			in2:  2,
			op:   (*Object).Sub,
			out:  -1,
		},

		{
			name: "*",
			in1:  2,
			in2:  3,
			op:   (*Object).Mul,
			out:  6,
		},

		{
			name: "*",
			in1:  "x",
			in2:  5,
			op:   (*Object).Mul,
			out:  "xxxxx",
		},

		{
			name: "*",
			in1:  5,
			in2:  "x",
			op:   (*Object).Mul,
			err:  "NotImplemented",
		},

		{
			name: "/",
			in1:  2,
			in2:  4,
			op:   (*Object).TrueDiv,
			out:  0.5,
		},

		{
			name: "//",
			in1:  2,
			in2:  4,
			op:   (*Object).FloorDiv,
			out:  0,
		},

		{
			name: "%",
			in1:  10,
			in2:  3,
			op:   (*Object).Mod,
			out:  1,
		},

		{
			name: "**",
			in1:  2,
			in2:  10,
			op:   (*Object).Pow,
			out:  1024,
		},

		{
			name: "<",
			in1:  2,
			in2:  4,
			op:   (*Object).Lt,
			out:  "True",
		},

		{
			name: "<",
			in1:  2,
			in2:  2,
			op:   (*Object).Lt,
			out:  "False",
		},

		{
			name: "<",
			in1:  4,
			in2:  2,
			op:   (*Object).Lt,
			out:  "False",
		},

		{
			name: ">",
			in1:  2,
			in2:  4,
			op:   (*Object).Gt,
			out:  "False",
		},

		{
			name: ">",
			in1:  2,
			in2:  2,
			op:   (*Object).Gt,
			out:  "False",
		},

		{
			name: ">",
			in1:  4,
			in2:  2,
			op:   (*Object).Gt,
			out:  "True",
		},

		{
			name: "<=",
			in1:  2,
			in2:  4,
			op:   (*Object).Le,
			out:  "True",
		},

		{
			name: "<=",
			in1:  2,
			in2:  2,
			op:   (*Object).Le,
			out:  "True",
		},

		{
			name: "<=",
			in1:  4,
			in2:  2,
			op:   (*Object).Le,
			out:  "False",
		},

		{
			name: ">=",
			in1:  2,
			in2:  4,
			op:   (*Object).Ge,
			out:  "False",
		},

		{
			name: ">=",
			in1:  2,
			in2:  2,
			op:   (*Object).Ge,
			out:  "True",
		},

		{
			name: ">=",
			in1:  4,
			in2:  2,
			op:   (*Object).Ge,
			out:  "True",
		},

		{
			name: "==",
			in1:  2,
			in2:  2,
			op:   (*Object).Eq,
			out:  "True",
		},

		{
			name: "==",
			in1:  2,
			in2:  4,
			op:   (*Object).Eq,
			out:  "False",
		},

		{
			name: "!=",
			in1:  2,
			in2:  2,
			op:   (*Object).Ne,
			out:  "False",
		},

		{
			name: "!=",
			in1:  2,
			in2:  4,
			op:   (*Object).Ne,
			out:  "True",
		},
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj1 := py.NewObject(test.in1)
		assert.NoError(obj1.Err())

		obj2 := py.NewObject(test.in2)
		assert.NoError(obj2.Err())

		res := test.op(obj1, obj2)

		exp := fmt.Sprintf("%v", test.out)
		if test.err != "" {
			exp = test.err
		}

		s, err := res.Str()
		if err != nil || s != exp {
			pres := s
			if err != nil {
				pres = err.Error()
			}

			t.Errorf("%v %s %v:\n"+
				"expected: %v\n"+
				"present: %v\n",
				test.in1,
				test.name,
				test.in2,
				exp, pres)
		}
	}
}

// TestBinopBadAttr covers the getattr() failure branch inside binop:
// requesting a dunder method that doesn't exist on the operand's type
// causes gate.getattr to fail, and binop must wrap that via
// newErrorObject instead of panicking or returning a nil result.
func TestBinopBadAttr(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj1 := py.NewObject(1)
	assert.NoError(obj1.Err())

	res := obj1.binop("__nonexistent__", 2)

	s, err := res.Str()
	if err == nil {
		t.Fatalf("expected an error, got result: %q", s)
	}

	if !strings.Contains(err.Error(), "AttributeError") {
		t.Errorf("expected error to contain %q, got %q",
			"AttributeError", err.Error())
	}
}

// TestBinopCallError covers the gate.call() failure branch inside
// binop: dividing by zero causes the underlying Python call to raise
// ZeroDivisionError, which gate.call reports as a genuine Go error
// (unlike the NotImplemented sentinel object returned by unsupported
// operand types), and binop must wrap that via newErrorObject.
func TestBinopCallError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj1 := py.NewObject(1)
	assert.NoError(obj1.Err())

	res := obj1.TrueDiv(0)

	if res.Err() == nil {
		s, _ := res.Str()
		t.Fatalf("expected an error Object, got: %q", s)
	}

	if !strings.Contains(res.Err().Error(), "ZeroDivisionError") {
		t.Errorf("expected error to contain %q, got %q",
			"ZeroDivisionError", res.Err().Error())
	}
}

// TestBinopBeginError covers the obj.begin() failure branch inside
// binop: if the receiver is already an error Object (obj.err != nil),
// begin() returns that error immediately, and binop must propagate it
// via newErrorObject.
func TestBinopBeginError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// obj1 is deliberately an error Object: Get() on a fresh int
	// Object for a nonexistent attribute returns ErrNotFound.
	bad := py.NewObject(1).Get("__nonexistent_attr__")
	assert.Must(bad.Err() != nil)

	res := bad.binop("__add__", 2)

	if res.Err() == nil {
		t.Fatalf("expected an error Object, got: %v", res)
	}

	if res.Err().Error() != bad.Err().Error() {
		t.Errorf("expected error to be propagated unchanged:\n"+
			"expected: %v\npresent: %v", bad.Err(), res.Err())
	}
}

// TestBinopNewPyObjectError covers the obj.py.newPyObject() failure
// branch inside binop: passing a Go value of a type that cannot be
// converted to a Python object causes newPyObject to fail, and binop
// must wrap that via newErrorObject.
func TestBinopNewPyObjectError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj1 := py.NewObject(1)
	assert.NoError(obj1.Err())

	// A channel has no Python equivalent and newPyObject must
	// reject it.
	res := obj1.Add(make(chan int))

	if res.Err() == nil {
		t.Fatalf("expected an error Object, got: %v", res)
	}
}
