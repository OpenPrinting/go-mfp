// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Fuzz tests for Python <-> Go type conversions (Object.Bool, Int,
// Float, etc. and Python.NewObject).
//
// Each target reuses a single shared *Python interpreter across all
// iterations instead of creating one per call, so more iterations/sec
// are possible. Interpreter create/teardown itself is covered by a
// separate fuzz target.
//
// A periodic countObjID() check is a cheap backstop for objects that
// conversions might fail to release. Every Object is explicitly
// released with Invalidate() rather than left to the GC finalizer, to
// keep that check deterministic.
package cpython

import (
	"math"
	"runtime"
	"testing"
	"unicode/utf8"
)

// How often (in executions) to force a GC and check the live-object
// count isn't trending upward without bound.
const conversionFuzzGCCheckEvery = 500

// newConversionFuzzPython creates the shared interpreter used by the
// conversion fuzz targets, and registers teardown via f.Cleanup.
func newConversionFuzzPython(f *testing.F) *Python {
	py, err := NewPython()
	if err != nil {
		f.Fatalf("NewPython: %s", err)
	}
	f.Cleanup(py.Close)
	return py
}

// conversionFuzzGCTracker tracks executions per fuzz target for the
// periodic GC/leak check. Each target gets its own tracker/baseline.
type conversionFuzzGCTracker struct {
	iterations int
	baseline   int
	haveBase   bool
}

// check runs a periodic countObjID() sanity check. Intentionally
// lenient (growth trend, not exact count), since short-lived helper
// objects may still be pending GC if not explicitly Invalidate()'d.
func (g *conversionFuzzGCTracker) check(t *testing.T, py *Python) {
	g.iterations++
	if g.iterations%conversionFuzzGCCheckEvery != 0 {
		return
	}

	runtime.GC()
	runtime.GC() // second pass for finalizers queued by the first

	count := py.countObjID()
	if !g.haveBase {
		g.baseline = count
		g.haveBase = true
		return
	}

	// Generous multiplier: looking for unbounded growth, not exact
	// parity, since finalizers race with the Go GC.
	if count > g.baseline*4+100 {
		t.Fatalf(
			"possible object leak in conversions: object count grew from %d to %d after %d iterations",
			g.baseline, count, g.iterations)
	}
}

// FuzzObjectRoundTripString fuzzes Go string -> Python str -> Go string.
func FuzzObjectRoundTripString(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	for _, s := range []string{"", "hello", "привет", "\x00", "\ufffd", "a\nb\tc"} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		obj := py.NewObject(s)

		// Go strings need not be valid UTF-8; Python strings must be.
		// Invalid UTF-8 should fail conversion cleanly (an error, not
		// a panic), so we only require an exact round trip otherwise.
		if !utf8.ValidString(s) {
			if obj.Err() == nil {
				t.Fatalf(
					"NewObject(%q) unexpectedly succeeded on invalid UTF-8",
					s)
			}
			obj.Invalidate()
			tracker.check(t, py)
			return
		}

		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%q) failed: %s", s, err)
		}

		got, err := obj.Unicode()
		if err != nil {
			t.Fatalf("Unicode() failed for %q: %s", s, err)
		}
		if got != s {
			t.Fatalf("round trip mismatch: sent %q, got %q", s, got)
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripBytes fuzzes Go []byte -> Python bytes -> Go []byte.
func FuzzObjectRoundTripBytes(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	for _, b := range [][]byte{{}, {0x00}, {0xff, 0xfe}, []byte("hello")} {
		f.Add(b)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		obj := py.NewObject(b)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%v) failed: %s", b, err)
		}

		got, err := obj.Bytes()
		if err != nil {
			t.Fatalf("Bytes() failed for %v: %s", b, err)
		}
		if string(got) != string(b) {
			t.Fatalf("round trip mismatch: sent %v, got %v", b, got)
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripInt fuzzes Go int64 -> Python int -> Go int64, and
// checks Uint()'s error behavior is consistent with the sign.
func FuzzObjectRoundTripInt(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	for _, v := range []int64{0, 1, -1, math.MaxInt64, math.MinInt64} {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v int64) {
		obj := py.NewObject(v)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%d) failed: %s", v, err)
		}

		got, err := obj.Int()
		if err != nil {
			t.Fatalf("Int() failed for %d: %s", v, err)
		}
		if got != v {
			t.Fatalf("round trip mismatch: sent %d, got %d", v, got)
		}

		// A negative value must never be accepted by Uint().
		if v < 0 {
			if _, err := obj.Uint(); err == nil {
				t.Fatalf("Uint() unexpectedly succeeded for negative value %d", v)
			}
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripUint fuzzes Go uint64 -> Python int -> Go uint64.
func FuzzObjectRoundTripUint(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	for _, v := range []uint64{0, 1, math.MaxUint64, math.MaxInt64, math.MaxInt64 + 1} {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v uint64) {
		obj := py.NewObject(v)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%d) failed: %s", v, err)
		}

		got, err := obj.Uint()
		if err != nil {
			t.Fatalf("Uint() failed for %d: %s", v, err)
		}
		if got != v {
			t.Fatalf("round trip mismatch: sent %d, got %d", v, got)
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripBool fuzzes Go bool -> Python bool -> Go bool.
func FuzzObjectRoundTripBool(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	f.Add(true)
	f.Add(false)

	f.Fuzz(func(t *testing.T, v bool) {
		obj := py.NewObject(v)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%v) failed: %s", v, err)
		}

		got, err := obj.Bool()
		if err != nil {
			t.Fatalf("Bool() failed for %v: %s", v, err)
		}
		if got != v {
			t.Fatalf("round trip mismatch: sent %v, got %v", v, got)
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripList fuzzes []any (from a fuzzed []int64) -> Python
// list -> []any, exercising newPyList's per-element ref/unref
// bookkeeping.
func FuzzObjectRoundTripList(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	f.Add([]byte{})
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})

	// No []int64 corpus type exists, so a []byte seed is reinterpreted
	// as small ints - exercises the list path without a custom encoder.
	f.Fuzz(func(t *testing.T, raw []byte) {
		want := make([]any, len(raw))
		for i, b := range raw {
			want[i] = int64(b)
		}

		obj := py.NewObject(want)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%v) failed: %s", want, err)
		}

		items, err := obj.Slice()
		if err != nil {
			t.Fatalf("Slice() failed for %v: %s", want, err)
		}
		if len(items) != len(want) {
			t.Fatalf("length mismatch: sent %d items, got %d", len(want), len(items))
		}

		for i, item := range items {
			got, err := item.Int()
			if err != nil {
				t.Fatalf("element %d: Int() failed: %s", i, err)
			}
			if got != want[i] {
				t.Fatalf("element %d mismatch: sent %v, got %v", i, want[i], got)
			}
			item.Invalidate()
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripDict fuzzes map[string]any (keys from a fuzzed
// string, values from byte position) -> Python dict -> back,
// exercising newPyDict's key-sort and per-entry ref/unref bookkeeping.
func FuzzObjectRoundTripDict(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	f.Add("")
	f.Add("abc")
	f.Add("hello world this is a longer key set for more entries")

	f.Fuzz(func(t *testing.T, s string) {
		// Turn each byte position into a distinct key (skipping
		// duplicates), so map size/content vary without needing a
		// multi-argument corpus.
		want := map[string]any{}
		for i, r := range s {
			key := string(r) + string(rune('a'+(i%26)))
			want[key] = int64(i)
		}

		obj := py.NewObject(want)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%v) failed: %s", want, err)
		}

		for k, v := range want {
			item := obj.GetItem(k)
			if err := item.Err(); err != nil {
				t.Fatalf("GetItem(%q) failed: %s", k, err)
			}

			got, err := item.Int()
			if err != nil {
				t.Fatalf("GetItem(%q).Int() failed: %s", k, err)
			}
			if got != v {
				t.Fatalf("value mismatch for key %q: sent %v, got %v", k, v, got)
			}
			item.Invalidate()
		}

		obj.Invalidate()
		tracker.check(t, py)
	})
}

// FuzzObjectListConversionErrorPath poisons one element of a valid list
// with an unconvertible Go value (a channel) at a fuzzed position,
// driving newPyList's mid-loop error return - an early-exit cleanup
// path where a forgotten unref is easy to miss in normal testing.
func FuzzObjectListConversionErrorPath(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	f.Add(uint8(0), 0)
	f.Add(uint8(5), 3)
	f.Add(uint8(20), 19)

	f.Fuzz(func(t *testing.T, size uint8, poisonAt int) {
		n := int(size)
		if n == 0 {
			return
		}

		list := make([]any, n)
		for i := range list {
			list[i] = int64(i)
		}

		// Poison exactly one slot (mod n keeps it in range).
		idx := ((poisonAt % n) + n) % n
		list[idx] = make(chan int) // unsupported: reflect.Chan

		before := py.countObjID()

		obj := py.NewObject(list)
		if obj.Err() == nil {
			t.Fatalf("NewObject unexpectedly succeeded with a poisoned element at %d", idx)
		}

		// Must fail cleanly without leaking elements converted before
		// the poisoned one.
		runtime.GC()
		after := py.countObjID()
		if after > before {
			t.Fatalf(
				"list conversion error path leaked objects: before=%d after=%d (size=%d poisonAt=%d)",
				before, after, n, idx)
		}

		tracker.check(t, py)
	})
}

// FuzzObjectRoundTripFloat fuzzes Go float64 -> Python float -> Go
// float64, including NaN/Inf, and cross-checks the int64/uint64
// boundary-conversion logic via Int()/Uint().
func FuzzObjectRoundTripFloat(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	for _, v := range []float64{
		0, 1, -1, 0.5, -0.5,
		math.MaxFloat64, -math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		math.NaN(), math.Inf(1), math.Inf(-1),
		float64(math.MaxInt64), float64(math.MinInt64),
	} {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v float64) {
		obj := py.NewObject(v)
		if err := obj.Err(); err != nil {
			t.Fatalf("NewObject(%v) failed: %s", v, err)
		}

		got, err := obj.Float()
		if err != nil {
			t.Fatalf("Float() failed for %v: %s", v, err)
		}

		switch {
		case math.IsNaN(v):
			if !math.IsNaN(got) {
				t.Fatalf("round trip mismatch: sent NaN, got %v", got)
			}
		default:
			if got != v {
				t.Fatalf("round trip mismatch: sent %v, got %v", v, got)
			}
		}

		// Int()/Uint() must not panic or hang regardless of magnitude;
		// whether they succeed depends on internal boundary constants,
		// which is what we're probing - we only assert on crash-freedom.
		_, _ = obj.Int()
		_, _ = obj.Uint()

		obj.Invalidate()
		tracker.check(t, py)
	})
}

