// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// pyGate fuzz test
//
// pyGate has no exported surface; exercised via the public API, which
// routes conversions through gate.make*/decode* and gate.eval.

package cpython

import (
	"bytes"
	"math"
	"math/big"
	"testing"
	"unicode/utf8"
)

// FuzzGateInt64RoundTrip fuzzes int64 through Python.NewObject/Object.Int.
func FuzzGateInt64RoundTrip(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))

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
		got, err := obj.Int()
		if err != nil || got != n {
			t.Fatalf("Int round-trip: got %d, err %v, want %d", got, err, n)
		}
	})
}

// FuzzGateUint64RoundTrip fuzzes uint64 through Python.NewObject/Object.Uint.
func FuzzGateUint64RoundTrip(f *testing.F) {
	f.Add(uint64(0))
	f.Add(uint64(math.MaxUint64))

	f.Fuzz(func(t *testing.T, n uint64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(n)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Uint()
		if err != nil || got != n {
			t.Fatalf("Uint round-trip: got %d, err %v, want %d", got, err, n)
		}
	})
}

// FuzzGateFloatRoundTrip fuzzes float64 through Python.NewObject/Object.Float.
func FuzzGateFloatRoundTrip(f *testing.F) {
	f.Add(0.0)
	f.Add(-1.5)
	f.Add(math.Inf(1))
	f.Add(math.Inf(-1))
	f.Add(math.NaN())

	f.Fuzz(func(t *testing.T, n float64) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(n)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Float()
		if err != nil {
			t.Fatalf("Float: %s", err)
		}
		if math.IsNaN(n) {
			if !math.IsNaN(got) {
				t.Fatalf("Float round-trip: got %v, want NaN", got)
			}
			return
		}
		if got != n {
			t.Fatalf("Float round-trip: got %v, want %v", got, n)
		}
	})
}

// FuzzGateBytesRoundTrip fuzzes []byte through Python.NewObject/Object.Bytes.
func FuzzGateBytesRoundTrip(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0, 1, 2, 0xff})

	f.Fuzz(func(t *testing.T, data []byte) {
		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(data)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Bytes()
		if err != nil || !bytes.Equal(got, data) {
			t.Fatalf("Bytes round-trip: got %v, err %v, want %v", got, err, data)
		}
	})
}

// FuzzGateStringRoundTrip fuzzes valid-UTF-8 strings through
// Python.NewObject/Object.Unicode (invalid UTF-8 legitimately errors,
// see FuzzCallbackName).
func FuzzGateStringRoundTrip(f *testing.F) {
	f.Add("hello")
	f.Add("привет 👋")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip("not valid UTF-8")
		}

		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(s)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Unicode()
		if err != nil || got != s {
			t.Fatalf("Unicode round-trip: got %q, err %v, want %q", got, err, s)
		}
	})
}

// FuzzGateBigintRoundTrip fuzzes arbitrary-precision integers through
// Python.NewObject/Object.Bigint.
//
// makeBigint converts through a decimal string (PyLong_FromString),
// which since Python 3.11 (backported to some 3.9/3.10 patches)
// enforces a default 4300-digit limit on int<->str conversion
// (sys.set_int_max_str_digits()) as a DoS mitigation. 1024 magnitude
// bytes is at most ~2466 decimal digits, comfortably under that limit
// on every supported Python version.
func FuzzGateBigintRoundTrip(f *testing.F) {
	f.Add(false, []byte{1, 2, 3})
	f.Add(true, []byte{})

	f.Fuzz(func(t *testing.T, neg bool, mag []byte) {
		if len(mag) > 1024 {
			t.Skip("too large")
		}

		n := new(big.Int).SetBytes(mag)
		if neg {
			n.Neg(n)
		}

		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		obj := py.NewObject(n)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject: %s", err)
		}
		got, err := obj.Bigint()
		if err != nil || got.Cmp(n) != 0 {
			t.Fatalf("Bigint round-trip: got %v, err %v, want %v", got, err, n)
		}
	})
}
