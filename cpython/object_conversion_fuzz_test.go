// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Fuzz tests for Python <-> Go type conversions (Object.Bool, Int,
// Float, etc. and Python.NewObject).
//
// Each fuzz target here reuses a single, shared *Python interpreter
// across all iterations rather than creating one per call. Interpreter
// creation/teardown is the concern of python_fuzz_test.go
// (FuzzSubInterpreterLifecycle); this file is only exercising the
// marshaling logic, so sharing one interpreter lets the fuzzer run far
// more iterations per second.
//
// A periodic countObjID() check still runs across iterations, as a
// cheap backstop for objects that conversions might fail to release.
//
// IMPORTANT: countObjID() counts entries in the binding's own live-
// object table (see Python.objects), which are normally released by
// an Object's runtime.SetFinalizer when the Go GC collects it. Go
// finalizers run asynchronously and are *not* guaranteed to have
// completed just because runtime.GC() was called - under fuzzing load
// (many workers, tens of thousands of execs/sec) a backlog of
// not-yet-finalized Objects is expected and is NOT a leak; it just
// looks like unbounded growth if you sample the count too eagerly.
//
// To keep the leak check deterministic, every Object created in these
// fuzz targets is explicitly released with Object.Invalidate() as soon
// as it's no longer needed, instead of relying on the Go GC finalizer
// to get around to it. That removes the timing race entirely: after
// Invalidate(), the object is *synchronously* gone from py.objects, so
// countObjID() reflects only genuinely-leaked objects, never a
// finalizer backlog.
package cpython

import (
	"math"
	"runtime"
	"testing"
	"unicode/utf8"
)

// conversionFuzzGCCheckEvery controls how often (in fuzz executions,
// within a single worker process) we force a GC and check that the
// live-object count is not trending upward without bound.
const conversionFuzzGCCheckEvery = 500

// newConversionFuzzPython creates the single shared interpreter used by
// the conversion fuzz targets in this file, and registers its teardown
// with f.Cleanup.
func newConversionFuzzPython(f *testing.F) *Python {
	py, err := NewPython()
	if err != nil {
		f.Fatalf("NewPython: %s", err)
	}
	f.Cleanup(py.Close)
	return py
}

// conversionFuzzObjCounter tracks executions per fuzz target, for the
// periodic GC/leak sanity check. Separate targets get separate counters
// via closures, so each has its own baseline.
type conversionFuzzGCTracker struct {
	iterations int
	baseline   int
	haveBase   bool
}

// check runs a periodic countObjID() sanity check. It is intentionally
// lenient (a growing-without-bound trend, not an exact count) since
// legitimate short-lived helper objects may still be pending Go GC for
// anything not explicitly released via Invalidate().
func (g *conversionFuzzGCTracker) check(t *testing.T, py *Python) {
	g.iterations++
	if g.iterations%conversionFuzzGCCheckEvery != 0 {
		return
	}

	runtime.GC()
	runtime.GC() // second pass to let finalizers queued by the first run

	count := py.countObjID()
	if !g.haveBase {
		g.baseline = count
		g.haveBase = true
		return
	}

	// Generous multiplier: we're looking for unbounded growth, not
	// exact parity, since finalizers race with the Go GC.
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

		// Go strings are arbitrary byte sequences and need not be
		// valid UTF-8. Python strings must be. A non-UTF-8 Go string
		// is expected to fail conversion cleanly (a returned error,
		// not a panic or crash) - that is correct behavior, not a
		// leak or bug, so we only require an exact round trip for
		// valid UTF-8 input.
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

// FuzzObjectRoundTripList fuzzes []any (built from a fuzzed []int64) ->
// Python list -> []any, exercising newPyList's per-element ref/unref
// bookkeeping (python.go).
func FuzzObjectRoundTripList(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	f.Add([]byte{})
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})

	// Native fuzzing has no []int64 corpus type, so a []byte seed is
	// reinterpreted as a sequence of small ints - plenty to exercise
	// the list path without needing a custom corpus encoder.
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

// FuzzObjectRoundTripDict fuzzes map[string]any (keys derived from a
// fuzzed string, values from its byte length) -> Python dict -> back,
// exercising newPyDict's key-sort and per-entry ref/unref bookkeeping.
func FuzzObjectRoundTripDict(f *testing.F) {
	py := newConversionFuzzPython(f)
	tracker := &conversionFuzzGCTracker{}

	f.Add("")
	f.Add("abc")
	f.Add("hello world this is a longer key set for more entries")

	f.Fuzz(func(t *testing.T, s string) {
		// Turn each byte position into a distinct single-character key
		// (skipping duplicates), so map size and content vary with the
		// fuzzed input without needing a multi-argument corpus.
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

// FuzzObjectListConversionErrorPath deliberately poisons one element of
// an otherwise-valid list with an unconvertible Go value (a channel),
// at a fuzzed position. This repeatedly drives newPyList's mid-loop
// error return - the exact kind of early-exit cleanup path where a
// forgotten unref is easy to introduce and easy to miss in normal
// testing, since it only runs when conversion fails partway through.
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

		// Poison exactly one slot, wherever the fuzzed index lands
		// (mod n keeps it in range regardless of what the fuzzer
		// generates).
		idx := ((poisonAt % n) + n) % n
		list[idx] = make(chan int) // unsupported: reflect.Chan

		before := py.countObjID()

		obj := py.NewObject(list)
		if obj.Err() == nil {
			t.Fatalf("NewObject unexpectedly succeeded with a poisoned element at %d", idx)
		}

		// The conversion must fail cleanly, and must not leak any of
		// the successfully-converted elements that came before the
		// poisoned one.
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
// boundary-conversion logic in float.go via Int()/Uint().
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
		// whether they succeed depends on float.go's boundary
		// constants, which is exactly what we're probing here. We
		// only assert on crash-freedom, not on which side of the
		// boundary any given value falls.
		_, _ = obj.Int()
		_, _ = obj.Uint()

		obj.Invalidate()
		tracker.check(t, py)
	})
}

