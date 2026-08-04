// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Binary operations on objects fuzz test

package cpython

import (
	"math"
	"math/big"
	"testing"
	"unicode/utf8"
)

// FuzzBinopArithmeticInt64 fuzzes Add/Sub/Mul against exact big.Int math.
func FuzzBinopArithmeticInt64(f *testing.F) {
	f.Add(int64(1), int64(2))
	f.Add(int64(-1), int64(1))
	f.Add(int64(math.MaxInt64), int64(math.MaxInt64))

	f.Fuzz(func(t *testing.T, a, b int64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(a)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}

		bigA, bigB := big.NewInt(a), big.NewInt(b)
		cases := []struct {
			name string
			op   func(*Object, any) *Object
			want *big.Int
		}{
			{"Add", (*Object).Add, new(big.Int).Add(bigA, bigB)},
			{"Sub", (*Object).Sub, new(big.Int).Sub(bigA, bigB)},
			{"Mul", (*Object).Mul, new(big.Int).Mul(bigA, bigB)},
		}

		for _, c := range cases {
			ret := c.op(obj, b)
			if err := ret.Err(); err != nil {
				t.Fatalf("%s(%d, %d): %s", c.name, a, b, err)
			}
			got, err := ret.Bigint()
			if err != nil || got.Cmp(c.want) != 0 {
				t.Fatalf("%s(%d, %d): got %v, err %v, want %v",
					c.name, a, b, got, err, c.want)
			}
		}
	})
}

// FuzzBinopDivMod fuzzes FloorDiv/TrueDiv/Mod; zero divisor must raise
// ZeroDivisionError, otherwise the call must complete without panicking.
func FuzzBinopDivMod(f *testing.F) {
	f.Add(int64(10), int64(3))
	f.Add(int64(1), int64(0))
	f.Add(int64(-7), int64(2))

	f.Fuzz(func(t *testing.T, a, b int64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(a)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}

		ops := map[string]func(*Object, any) *Object{
			"FloorDiv": (*Object).FloorDiv,
			"TrueDiv":  (*Object).TrueDiv,
			"Mod":      (*Object).Mod,
		}
		for name, op := range ops {
			ret := op(obj, b)
			if b == 0 {
				if ret.Err() == nil {
					t.Fatalf("%s(%d, 0): expected ZeroDivisionError, got nil",
						name, a)
				}
				continue
			}
			if err := ret.Err(); err != nil {
				t.Fatalf("%s(%d, %d): %s", name, a, b, err)
			}
		}
	})
}

// FuzzBinopCompare fuzzes Lt/Gt/Le/Ge/Eq/Ne on two int64 operands and
// checks results against Go's own comparisons.
func FuzzBinopCompare(f *testing.F) {
	f.Add(int64(1), int64(2))
	f.Add(int64(5), int64(5))

	f.Fuzz(func(t *testing.T, a, b int64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(a)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}

		check := func(name string, op func(*Object, any) *Object, want bool) {
			ret := op(obj, b)
			if err := ret.Err(); err != nil {
				t.Fatalf("%s(%d, %d): %s", name, a, b, err)
			}
			got, err := ret.Bool()
			if err != nil || got != want {
				t.Fatalf("%s(%d, %d): got %v, err %v, want %v",
					name, a, b, got, err, want)
			}
		}

		check("Lt", (*Object).Lt, a < b)
		check("Gt", (*Object).Gt, a > b)
		check("Le", (*Object).Le, a <= b)
		check("Ge", (*Object).Ge, a >= b)
		check("Eq", (*Object).Eq, a == b)
		check("Ne", (*Object).Ne, a != b)
	})
}

// FuzzBinopTypeMismatch fuzzes Add between an int Object and a fuzzed
// string. binop calls __add__ directly, not the full '+' protocol, so
// int.__add__(str) returns NotImplemented rather than raising (see
// binops_test.go's "5 * 'x'" case). Invalid UTF-8 errors earlier, in
// makeString. Either way, must never panic.
func FuzzBinopTypeMismatch(f *testing.F) {
	f.Add(int64(1), "x")
	f.Add(int64(0), "")

	f.Fuzz(func(t *testing.T, n int64, s string) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(n)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}

		ret := obj.Add(s)

		if !utf8.ValidString(s) {
			if ret.Err() == nil {
				t.Fatalf("Add(%d, invalid-UTF8): expected error, got none", n)
			}
			return
		}

		if err := ret.Err(); err != nil {
			t.Fatalf("Add(%d, %q): unexpected error: %s", n, s, err)
		}
		str, err := ret.Str()
		if err != nil || str != "NotImplemented" {
			t.Fatalf("Add(%d, %q): got %q, err %v, want %q",
				n, s, str, err, "NotImplemented")
		}
	})
}
