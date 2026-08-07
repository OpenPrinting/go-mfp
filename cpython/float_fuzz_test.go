// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Floating point limits fuzz test

package cpython

import (
	"math"
	"testing"
)

// FuzzFloatToInt64Boundary checks Object.Int() succeeds with an exact
// int64 exactly when the value is within [minInt64Float, maxInt64Float].
func FuzzFloatToInt64Boundary(f *testing.F) {
	f.Add(0.0)
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(math.NaN())
	f.Add(math.Inf(1))
	f.Add(math.Inf(-1))
	f.Add(maxInt64Float)
	f.Add(minInt64Float)
	f.Fuzz(func(t *testing.T, val float64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()
		obj := py.NewObject(val)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Int()
		inRange := minInt64Float <= val && val <= maxInt64Float
		if inRange {
			if err != nil {
				t.Fatalf("Int(%v): unexpected error: %s", val, err)
			}
			if got != int64(val) {
				t.Fatalf("Int(%v): got %d, want %d", val, got, int64(val))
			}
			return
		}
		if err == nil {
			t.Fatalf("Int(%v): expected overflow error, got %d", val, got)
		}
	})
}

// FuzzFloatToUint64Boundary is FuzzFloatToInt64Boundary for Object.Uint().
func FuzzFloatToUint64Boundary(f *testing.F) {
	f.Add(0.0)
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(math.NaN())
	f.Add(math.Inf(1))
	f.Add(math.Inf(-1))
	f.Add(maxUint64Float)
	f.Fuzz(func(t *testing.T, val float64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()
		obj := py.NewObject(val)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Uint()
		inRange := 0 <= val && val <= maxUint64Float
		if inRange {
			if err != nil {
				t.Fatalf("Uint(%v): unexpected error: %s", val, err)
			}
			if got != uint64(val) {
				t.Fatalf("Uint(%v): got %d, want %d", val, got, uint64(val))
			}
			return
		}
		if err == nil {
			t.Fatalf("Uint(%v): expected overflow error, got %d", val, got)
		}
	})
}
