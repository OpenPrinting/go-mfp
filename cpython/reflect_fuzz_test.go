// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Reflect helper functions fuzz test

package cpython

import (
	"reflect"
	"testing"
	"unicode/utf8"
)

// FuzzReflectSortInt64 fuzzes reflectSort on a slice of int64 values,
// decoded 8 bytes at a time from raw fuzz input, and checks the result
// is sorted ascending.
func FuzzReflectSortInt64(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, raw []byte) {
		n := len(raw) / 8
		if n < 1 {
			t.Skip("too short")
		}

		values := make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			var v int64
			for j := 0; j < 8; j++ {
				v = v<<8 | int64(raw[i*8+j])
			}
			values[i] = reflect.ValueOf(v)
		}

		if !reflectSort(values) {
			t.Fatalf("reflectSort returned false for uniform int64 slice")
		}
		for i := 0; i < len(values)-1; i++ {
			if values[i].Int() > values[i+1].Int() {
				t.Fatalf("not sorted at %d: %d > %d",
					i, values[i].Int(), values[i+1].Int())
			}
		}
	})
}

// FuzzReflectSortStrings fuzzes reflectSort on a slice of strings split
// from one fuzzed input on '\n'.
func FuzzReflectSortStrings(f *testing.F) {
	f.Add("banana\napple\ncherry")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		parts := splitOnNewline(s)

		values := make([]reflect.Value, len(parts))
		for i, p := range parts {
			values[i] = reflect.ValueOf(p)
		}

		if !reflectSort(values) {
			t.Fatalf("reflectSort returned false for string slice")
		}
		for i := 0; i < len(values)-1; i++ {
			if values[i].String() > values[i+1].String() {
				t.Fatalf("not sorted at %d: %q > %q",
					i, values[i].String(), values[i+1].String())
			}
		}
	})
}

// splitOnNewline is a tiny helper avoiding an extra import for
// strings.Split.
func splitOnNewline(s string) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// FuzzReflectDictRoundTrip fuzzes the full Python.NewObject path for
// map[string]int64, which internally calls reflectSort to produce
// deterministic dict key ordering.
func FuzzReflectDictRoundTrip(f *testing.F) {
	f.Add("a", int64(1), "b", int64(2))

	f.Fuzz(func(t *testing.T, k1 string, v1 int64, k2 string, v2 int64) {
		if k1 == k2 || !utf8.ValidString(k1) || !utf8.ValidString(k2) {
			t.Skip("duplicate key or invalid UTF-8")
		}

		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		m := map[string]int64{k1: v1, k2: v2}
		obj := py.NewObject(m)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(map): %s", err)
		}

		n, err := obj.Len()
		if err != nil || n != 2 {
			t.Fatalf("Len: got %d, err %v, want 2", n, err)
		}
	})
}
