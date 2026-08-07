// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Unary operations on objects fuzz test

package cpython

import (
	"math/big"
	"testing"
)

// FuzzUnopsRoundTrip fuzzes Neg/Pos/Invert on an int64 Object and checks
// results against exact big.Int arithmetic (Python ints don't overflow,
// unlike the Go int64 they were built from).
func FuzzUnopsRoundTrip(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1) << 62)

	f.Fuzz(func(t *testing.T, n int64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(n)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}

		bigN := big.NewInt(n)

		neg := obj.Neg()
		if err := neg.Err(); err != nil {
			t.Fatalf("Neg(%d): %s", n, err)
		}
		if got, err := neg.Bigint(); err != nil ||
			got.Cmp(new(big.Int).Neg(bigN)) != 0 {
			t.Fatalf("Neg(%d): got %v, err %v", n, got, err)
		}

		pos := obj.Pos()
		if err := pos.Err(); err != nil {
			t.Fatalf("Pos(%d): %s", n, err)
		}
		if got, err := pos.Bigint(); err != nil || got.Cmp(bigN) != 0 {
			t.Fatalf("Pos(%d): got %v, err %v", n, got, err)
		}

		inv := obj.Invert()
		if err := inv.Err(); err != nil {
			t.Fatalf("Invert(%d): %s", n, err)
		}
		// Python's ~x == -x-1, same identity as Go's ^n for int64.
		want := new(big.Int).Sub(new(big.Int).Neg(bigN), big.NewInt(1))
		if got, err := inv.Bigint(); err != nil || got.Cmp(want) != 0 {
			t.Fatalf("Invert(%d): got %v, err %v, want %v", n, got, err, want)
		}
	})
}
