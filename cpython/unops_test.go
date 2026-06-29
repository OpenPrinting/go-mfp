// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Unary operations on objects tests

package cpython

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestBinops tests binary operations on Objects
func TestUnops(t *testing.T) {
	type testData struct {
		name string                // Operation name
		in   any                   // Operand
		op   func(*Object) *Object // Operation
		out  any                   // Expected output
		err  string                // Expected error
	}

	tests := []testData{
		{
			name: "-",
			in:   1,
			op:   (*Object).Neg,
			out:  -1,
		},

		{
			name: "-",
			in:   "s",
			op:   (*Object).Neg,
			err:  `AttributeError: 'str' object has no attribute '__neg__'`,
		},

		{
			name: "+",
			in:   1,
			op:   (*Object).Pos,
			out:  1,
		},

		{
			name: "~",
			in:   1,
			op:   (*Object).Invert,
			out:  -2,
		},
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj := py.NewObject(test.in)
		assert.NoError(obj.Err())

		res := test.op(obj)
		exp := fmt.Sprintf("%v", test.out)
		if test.err != "" {
			exp = test.err
		}

		pres, err := res.Str()
		if err != nil {
			pres = err.Error()
		}

		if exp != pres {
			t.Errorf("%s %v:\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.name,
				test.in,
				exp, pres)
		}
	}
}

// TestUnopBeginError tests the unop() error path that fires when
// obj.begin() fails, i.e., when called on an Object that is already
// in the error state.
func TestUnopBeginError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// A channel has no defined conversion to a Python object,
	// so NewObject returns an error Object. obj.begin() returns
	// this same error without ever acquiring the gate.
	obj := py.NewObject(make(chan int))
	if obj.Err() == nil {
		t.Fatal("expected obj.Err() != nil for unconvertible value")
	}

	res := obj.Neg()
	if res.Err() == nil {
		t.Errorf("Neg on broken Object: expected error, got none")
	}
	if res.Err() != obj.Err() {
		t.Errorf("Neg on broken Object:\n"+
			"expected: %v\n"+
			"present:  %v\n",
			obj.Err(), res.Err())
	}
}

// TestUnopCallError tests the unop() error path that fires when
// gate.call() fails, i.e., the named dunder method exists and is
// callable but raises an exception when invoked.
func TestUnopCallError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	err = py.Exec(`
class Bad:
    def __neg__(self):
        raise ValueError("nope")
`, "")
	assert.NoError(err)

	cls := py.Get("Bad")
	assert.NoError(cls.Err())

	obj := cls.Call()
	assert.NoError(obj.Err())

	res := obj.Neg()
	if res.Err() == nil {
		t.Fatal("Neg on Bad(): expected error, got none")
	}

	const expect = "ValueError: nope"
	if !strings.HasPrefix(res.Err().Error(), expect) {
		t.Errorf("Neg on Bad():\n"+
			"expected prefix: %v\n"+
			"present:         %v\n",
			expect, res.Err())
	}
}
