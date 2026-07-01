// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests for Python <-> Go type conversions (Object.Bool, Int, Float, etc.)
// and type inspection (Object.IsBool, IsSeq, etc.)

package cpython

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestObjectFromPython tests decoding of Python values into Go types.
func TestObjectFromPython(t *testing.T) {
	type testData struct {
		expr     string                     // Python expression
		val      any                        // Expected value
		unbox    func(*Object) (any, error) // Value unboxing function
		mustfail bool                       // Error expected
	}

	const verybig = "21267647892944572736998860269687930881"

	bigint := func(s string) *big.Int {
		v := big.NewInt(0)
		v.SetString(s, 0)
		return v
	}

	unboxNone := func(obj *Object) (any, error) { return obj.IsNone(), nil }
	unboxBigint := func(obj *Object) (any, error) { return obj.Bigint() }
	unboxBool := func(obj *Object) (any, error) { return obj.Bool() }
	unboxBytes := func(obj *Object) (any, error) { return obj.Bytes() }
	unboxComplex := func(obj *Object) (any, error) { return obj.Complex() }
	unboxFloat := func(obj *Object) (any, error) { return obj.Float() }
	unboxInt := func(obj *Object) (any, error) { return obj.Int() }
	unboxUint := func(obj *Object) (any, error) { return obj.Uint() }
	unboxUnicode := func(obj *Object) (any, error) { return obj.Unicode() }

	tests := []testData{
		{expr: `None`, val: true, unbox: unboxNone},
		{expr: `True`, val: true, unbox: unboxBool},
		{expr: `0`, val: false, unbox: unboxBool},
		{expr: `False`, val: false, unbox: unboxBool},
		{expr: `"hello"`, val: "hello", unbox: unboxUnicode},
		{expr: `"привет"`, val: "привет", unbox: unboxUnicode},
		{expr: `0`, val: "привет", unbox: unboxUnicode, mustfail: true},
		{expr: `""`, val: "", unbox: unboxUnicode},
		{expr: `0`, val: int64(0), unbox: unboxInt},
		{expr: `0`, val: uint64(0), unbox: unboxUint},
		{expr: `0x7fffffff`, val: int64(0x7fffffff), unbox: unboxInt},
		{expr: `0x7fffffff`, val: uint64(0x7fffffff), unbox: unboxUint},
		{expr: `-0x7fffffff`, val: int64(-0x7fffffff), unbox: unboxInt},
		{expr: `-0x7fffffff`, unbox: unboxUint, mustfail: true},
		{expr: `0xffffffff`, val: int64(0xffffffff), unbox: unboxInt},
		{expr: `0xffffffff`, val: uint64(0xffffffff), unbox: unboxUint},
		{expr: `0xffffffffffffffff`, val: uint64(0xffffffffffffffff), unbox: unboxUint},
		{expr: verybig, unbox: unboxInt, mustfail: true},
		{expr: verybig, unbox: unboxUint, mustfail: true},
		{expr: `0`, val: bigint("0"), unbox: unboxBigint},
		{expr: `1`, val: bigint("1"), unbox: unboxBigint},
		{expr: `-1`, val: bigint("-1"), unbox: unboxBigint},
		{expr: `0xffffffffffffffff`, val: bigint("0xffffffffffffffff"), unbox: unboxBigint},
		{expr: verybig, val: bigint(verybig), unbox: unboxBigint},
		{expr: `None`, unbox: unboxInt, mustfail: true},
		{expr: `None`, unbox: unboxUint, mustfail: true},
		{expr: `None`, unbox: unboxBigint, mustfail: true},
		{expr: `b'\x01\x02\x03'`, val: []byte{0x1, 0x2, 0x3}, unbox: unboxBytes},
		{expr: `bytearray(b'\x01\x02\x03')`, val: []byte{0x1, 0x2, 0x3}, unbox: unboxBytes},
		{expr: `None`, unbox: unboxBytes, mustfail: true},
		{expr: `0.5`, val: 0.5, unbox: unboxFloat},
		{expr: `None`, unbox: unboxFloat, mustfail: true},
		{expr: `complex(1,2)`, val: complex(1, 2), unbox: unboxComplex},
		{expr: `None`, unbox: unboxComplex, mustfail: true},

		// Corner cases for integers
		{expr: `-1`, val: int64(-1), unbox: unboxInt},
		{expr: `-9223372036854775808`, val: int64(-9223372036854775808), unbox: unboxInt},
		{expr: `9223372036854775807`, val: int64(9223372036854775807), unbox: unboxInt},
		{expr: `0xffffffffffffffff`, val: uint64(0xffffffffffffffff), unbox: unboxUint},
		{expr: `18446744073709551615`, val: uint64(18446744073709551615), unbox: unboxUint},

		// Numerical conversions -> complex
		{expr: `0`, val: complex(0, 0), unbox: unboxComplex},
		{expr: `1`, val: complex(1, 0), unbox: unboxComplex},
		{expr: `1.5`, val: complex(1.5, 0), unbox: unboxComplex},
		{expr: `"hello"`, unbox: unboxComplex, mustfail: true},

		// Numerical conversions -> float
		{expr: `0`, val: 0.0, unbox: unboxFloat},
		{expr: `1`, val: 1.0, unbox: unboxFloat},
		{expr: `1.5`, val: 1.5, unbox: unboxFloat},
		{expr: `0xffffffffffffffff`, val: 18446744073709551615., unbox: unboxFloat},
		{expr: `0xfffffffffffffffff`, unbox: unboxFloat, mustfail: true},
		{expr: `"hello"`, unbox: unboxFloat, mustfail: true},

		// Numerical conversions -> int64
		{expr: `0`, val: int64(0), unbox: unboxInt},
		{expr: `0.0`, val: int64(0), unbox: unboxInt},
		{expr: `2.0`, val: int64(2), unbox: unboxInt},
		{expr: `1.5`, val: int64(1), unbox: unboxInt},
		{expr: `-1.5`, val: int64(-1), unbox: unboxInt},
		{expr: fmt.Sprintf("%f", maxInt64Float), val: int64(maxInt64Float), unbox: unboxInt},
		{expr: fmt.Sprintf("%f", minInt64Float), val: int64(minInt64Float), unbox: unboxInt},
		{expr: fmt.Sprintf("%f", maxUint64Float), mustfail: true, unbox: unboxInt},

		// Numerical conversions -> uint64
		{expr: `0`, val: uint64(0), unbox: unboxUint},
		{expr: `0.0`, val: uint64(0), unbox: unboxUint},
		{expr: `2.0`, val: uint64(2), unbox: unboxUint},
		{expr: `1.5`, val: uint64(1), unbox: unboxUint},
		{expr: `-1.5`, mustfail: true, unbox: unboxUint},
		{expr: fmt.Sprintf("%f", maxInt64Float), val: uint64(maxInt64Float), unbox: unboxUint},
		{expr: fmt.Sprintf("%f", minInt64Float), mustfail: true, unbox: unboxUint},
		{expr: fmt.Sprintf("%f", maxUint64Float), val: uint64(maxUint64Float), unbox: unboxUint},
		{expr: fmt.Sprintf("%f", maxUint64Float*2), mustfail: true, unbox: unboxUint},
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj := py.Eval(test.expr)

		if err := obj.Err(); err != nil {
			t.Errorf("%s: unexpected error: %s", test.expr, err)
			continue
		}

		val, err := test.unbox(obj)
		if err != nil {
			if !test.mustfail {
				t.Errorf("%s: object value error: %s", test.expr, err)
			}
			continue
		}

		if test.mustfail {
			t.Errorf("%s: error expected but didn't occur", test.expr)
			continue
		}

		if !reflect.DeepEqual(val, test.val) {
			t.Errorf("%s: object value mismatch:\nexpected: %#v\npresent:  %#v\n",
				test.expr, test.val, val)
		}
	}
}

// TestObjectIsMapSeq tests type-inspection methods: IsBool, IsBytes, IsSeq, etc.
func TestObjectIsMapSeq(t *testing.T) {
	type testData struct {
		expr        string
		isbool      bool
		isbytearray bool
		isbytes     bool
		iscallable  bool
		iscomplex   bool
		isdict      bool
		isfloat     bool
		islong      bool
		isnone      bool
		isseq       bool
		isunicode   bool
		istype      bool
		istrue      bool
		isfalse     bool
		iserror     bool
	}

	tests := []testData{
		{expr: "True", isbool: true, istrue: true},
		{expr: "False", isbool: true, isfalse: true},
		{expr: `bytearray(b'\x01\x02\x03')`, isbytearray: true},
		{expr: `bytes(b'\x01\x02\x03')`, isbytes: true},
		{expr: `min`, iscallable: true},
		{expr: `0.5 + 0.25j`, iscomplex: true, istrue: true},
		{expr: `{}`, isdict: true},
		{expr: `0.5`, isfloat: true, istrue: true},
		{expr: `5`, islong: true, istrue: true},
		{expr: `None`, isnone: true, isfalse: true},
		{expr: `[]`, isseq: true},
		{expr: `()`, isseq: true},
		{expr: `"hello"`, isunicode: true},
		{expr: `int`, iscallable: true, istype: true},
		{expr: `str`, iscallable: true, istype: true},
		{expr: `undefined`, iserror: true}, // Force an error
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	for _, test := range tests {
		obj := py.Eval(test.expr)

		if v := obj.IsBool(); v != test.isbool {
			t.Errorf("%#v: Object.IsBool: expected %v, present %v",
				test.expr, test.isbool, v)
		}
		if v := obj.IsByteArray(); v != test.isbytearray {
			t.Errorf("%#v: Object.IsByteArray: expected %v, present %v",
				test.expr, test.isbytearray, v)
		}
		if v := obj.IsBytes(); v != test.isbytes {
			t.Errorf("%#v: Object.IsBytes: expected %v, present %v",
				test.expr, test.isbytes, v)
		}
		if v := obj.IsCallable(); v != test.iscallable {
			t.Errorf("%#v: Object.IsCallable: expected %v, present %v",
				test.expr, test.iscallable, v)
		}
		if v := obj.IsComplex(); v != test.iscomplex {
			t.Errorf("%#v: Object.IsComplex: expected %v, present %v",
				test.expr, test.iscomplex, v)
		}
		if v := obj.IsDict(); v != test.isdict {
			t.Errorf("%#v: Object.IsDict: expected %v, present %v",
				test.expr, test.isdict, v)
		}
		if v := obj.IsFloat(); v != test.isfloat {
			t.Errorf("%#v: Object.IsFloat: expected %v, present %v",
				test.expr, test.isfloat, v)
		}
		if v := obj.IsLong(); v != test.islong {
			t.Errorf("%#v: Object.IsLong: expected %v, present %v",
				test.expr, test.islong, v)
		}
		if v := obj.IsNone(); v != test.isnone {
			t.Errorf("%#v: Object.IsNone: expected %v, present %v",
				test.expr, test.isnone, v)
		}
		if v := obj.IsSeq(); v != test.isseq {
			t.Errorf("%#v: Object.IsSeq: expected %v, present %v",
				test.expr, test.isseq, v)
		}
		if v := obj.IsUnicode(); v != test.isunicode {
			t.Errorf("%#v: Object.IsUnicode: expected %v, present %v",
				test.expr, test.isunicode, v)
		}
		if v := obj.IsType(); v != test.istype {
			t.Errorf("%#v: Object.IsType: expected %v, present %v",
				test.expr, test.istype, v)
		}
		if v := obj.IsTrue(); v != test.istrue {
			t.Errorf("%#v: Object.IsTrue: expected %v, present %v",
				test.expr, test.istrue, v)
		}
		if v := obj.IsFalse(); v != test.isfalse {
			t.Errorf("%#v: Object.IsFalse: expected %v, present %v",
				test.expr, test.isfalse, v)
		}
		if v := obj.IsError(); v != test.iserror {
			t.Errorf("%#v: Object.IsError: expected %v, present %v",
				test.expr, test.iserror, v)
		}
	}
}

// TestNewObject tests Python.NewObject — encoding Go values to Python objects.
func TestNewObject(t *testing.T) {
	type testData struct {
		in  any    // Input value
		out string // Expected output
	}

	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	bigint := func(s string) *big.Int {
		v := big.NewInt(0)
		v.SetString(s, 0)
		return v
	}

	eval := func(s string) *Object {
		obj := py.Eval(s)
		assert.NoError(obj.Err())
		return obj
	}

	tests := []testData{
		{in: nil, out: `None`},
		{in: true, out: `True`},
		{in: false, out: `False`},
		{in: 0, out: `0`},
		{in: 1, out: `1`},
		{in: -1, out: `-1`},
		{in: int8(0), out: `0`},
		{in: int8(math.MaxInt8), out: fmt.Sprintf("%d", math.MaxInt8)},
		{in: int8(math.MinInt8), out: fmt.Sprintf("%d", math.MinInt8)},
		{in: int16(0), out: `0`},
		{in: int16(math.MaxInt16), out: fmt.Sprintf("%d", math.MaxInt16)},
		{in: int16(math.MinInt16), out: fmt.Sprintf("%d", math.MinInt16)},
		{in: int32(0), out: `0`},
		{in: int32(math.MaxInt32), out: fmt.Sprintf("%d", math.MaxInt32)},
		{in: int32(math.MinInt32), out: fmt.Sprintf("%d", math.MinInt32)},
		{in: int64(0), out: `0`},
		{in: int64(math.MaxInt64), out: fmt.Sprintf("%d", math.MaxInt64)},
		{in: int64(math.MinInt64), out: fmt.Sprintf("%d", math.MinInt64)},
		{in: uint(0), out: `0`},
		{in: uint(1), out: `1`},
		{in: uint8(0), out: `0`},
		{in: uint8(math.MaxUint8), out: fmt.Sprintf("%d", math.MaxUint8)},
		{in: uint16(0), out: `0`},
		{in: uint16(math.MaxUint16), out: fmt.Sprintf("%d", math.MaxUint16)},
		{in: uint32(0), out: `0`},
		{in: uint32(math.MaxUint32), out: fmt.Sprintf("%d", math.MaxUint32)},
		{in: uint64(0), out: `0`},
		{in: uint64(math.MaxUint64), out: fmt.Sprintf("%d", uint64(math.MaxUint64))},
		{in: bigint("0"), out: `0`},
		{in: bigint("1"), out: `1`},
		{in: bigint("-1"), out: `-1`},
		{in: bigint("340282366920938463426481119284349108225"),
			out: `340282366920938463426481119284349108225`},
		{in: bigint("-340282366920938463426481119284349108225"),
			out: `-340282366920938463426481119284349108225`},
		{in: 0.5, out: `0.5`},
		{in: -0.5, out: `-0.5`},
		{in: 0.25 + 0.25i, out: `(0.25+0.25j)`},
		{in: -0.25 - 0.25i, out: `(-0.25-0.25j)`},
		{in: "", out: ``},
		{in: "Hello, world!", out: `Hello, world!`},
		{in: "Здравствуй, мир!", out: `Здравствуй, мир!`},
		{in: eval("12345"), out: `12345`},
		{in: []int{1, 2, 3}, out: `[1, 2, 3]`},
		{in: [3]int{1, 2, 3}, out: `[1, 2, 3]`},
		{in: []byte("ABC"), out: `b'ABC'`},
		{in: [3]byte{'A', 'B', 'C'}, out: `b'ABC'`},
		{in: []byte{}, out: `b''`},
		{in: [0]byte{}, out: `b''`},
		{in: map[int]int{}, out: `{}`},
		{in: map[bool]string{true: "T", false: "F"}, out: `{False: 'F', True: 'T'}`},
		{in: map[int]string{1: "1", 2: "2", 3: "3"}, out: `{1: '1', 2: '2', 3: '3'}`},
		{in: map[uint]string{1: "1", 2: "2", 3: "3"}, out: `{1: '1', 2: '2', 3: '3'}`},
		{in: map[float64]string{0.25: "1/4", 0.5: "1/2"}, out: `{0.25: '1/4', 0.5: '1/2'}`},
		{in: map[string]int{"one": 1, "two": 2, "three": 3},
			out: `{'one': 1, 'three': 3, 'two': 2}`},
		{in: map[string]int{"раз": 1, "два": 2, "три": 3},
			out: `{'два': 2, 'раз': 1, 'три': 3}`},
		{in: map[int]int(nil), out: `{}`},
		{in: map[any]any(nil), out: `{}`},
	}

	for _, test := range tests {
		obj := py.NewObject(test.in)
		if err := obj.Err(); err != nil {
			t.Errorf("%v: Python.NewObject: %s", test.in, err)
			continue
		}

		s, err := obj.Str()
		if err != nil {
			t.Errorf("%v: Object.Str: %s", test.in, err)
			continue
		}

		if s != test.out {
			t.Errorf("%v: Python.NewObject:\nexpected: %s\npresent:  %s\n",
				test.in, test.out, s)
			fmt.Printf("%v\n", obj)
		}
	}
}
