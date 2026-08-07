// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Concurrency and GIL-safety tests.
//
// These tests are meant to be run with -race. They don't assert much
// beyond "no crash / no data race / no deadlock", because that's the
// class of bug this file exists to catch: pyGateAcquire's LockOSThread
// + closelock interaction, and any object-map corruption under
// concurrent gate() calls from multiple goroutines.

package cpython

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// concurrencyWorkers returns a worker count that scales with
// available CPUs but stays reasonable in CI.
func concurrencyWorkers() int {
	n := runtime.NumCPU() * 4
	if n < 8 {
		n = 8
	}
	if n > 64 {
		n = 64
	}
	return n
}

// TestConcurrentIndependentInterpreters creates and drives many
// independent Python sub-interpreters in parallel goroutines. Each
// goroutine only ever touches its own interpreter, so this isolates
// bugs in per-interpreter state (as opposed to bugs in sharing one
// interpreter across goroutines, covered separately below).
func TestConcurrentIndependentInterpreters(t *testing.T) {
	workers := concurrencyWorkers()

	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			py, err := NewPython()
			if err != nil {
				errs <- fmt.Errorf("worker %d: NewPython: %w", id, err)
				return
			}
			defer py.Close()

			for j := 0; j < 200; j++ {
				err := py.Set("x", id*1000+j)
				if err != nil {
					errs <- fmt.Errorf("worker %d: Set: %w", id, err)
					return
				}

				obj := py.Eval("x * 2")
				if err := obj.Err(); err != nil {
					errs <- fmt.Errorf("worker %d: Eval: %w", id, err)
					return
				}

				got, err := obj.Int()
				if err != nil {
					errs <- fmt.Errorf("worker %d: Int: %w", id, err)
					return
				}

				want := int64((id*1000 + j) * 2)
				if got != want {
					errs <- fmt.Errorf("worker %d: got %d, want %d",
						id, got, want)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// TestConcurrentSharedInterpreter drives a single, shared *Python
// instance from many goroutines at once. This exercises gate()'s
// closelock/LockOSThread interaction and the objmap's internal
// locking under real concurrent pressure.
func TestConcurrentSharedInterpreter(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	workers := concurrencyWorkers()

	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				obj := py.NewObject(id*1000 + j)
				if err := obj.Err(); err != nil {
					errs <- fmt.Errorf("worker %d: NewObject: %w", id, err)
					return
				}

				got, err := obj.Int()
				if err != nil {
					errs <- fmt.Errorf("worker %d: Int: %w", id, err)
					return
				}

				if got != int64(id*1000+j) {
					errs <- fmt.Errorf("worker %d: round-trip mismatch: got %d",
						id, got)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// TestConcurrentCloseRace hammers NewPython/Close in a loop from one
// goroutine while other goroutines are actively using *separate*
// live interpreters, and periodically closes those too. This is
// aimed at the closelock synchronization specifically: Close() must
// never observe or cause a half-torn-down interpreter to be used by
// a concurrent gate() call.
func TestConcurrentCloseRace(t *testing.T) {
	const duration = 500 * time.Millisecond
	workers := concurrencyWorkers()

	var stop atomic.Bool
	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for !stop.Load() {
				py, err := NewPython()
				if err != nil {
					errs <- fmt.Errorf("worker %d: NewPython: %w", id, err)
					return
				}

				// Use the interpreter briefly, then close it
				// from the same goroutine that created it, racing
				// against every other goroutine's create/use/close
				// cycle on their own interpreters.
				obj := py.Eval("1 + 1")
				if err := obj.Err(); err != nil {
					errs <- fmt.Errorf("worker %d: Eval: %w", id, err)
					py.Close()
					return
				}

				py.Close()

				// Using it after Close must return ErrClosed, not
				// crash or hang.
				if err := py.Set("x", 1); err == nil {
					errs <- fmt.Errorf("worker %d: Set after Close: expected error", id)
					return
				}
			}
		}(i)
	}

	time.Sleep(duration)
	stop.Store(true)
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// TestConcurrentCallbacks verifies that Go callbacks registered as
// Python-callable functions behave correctly when invoked from
// multiple goroutines' interpreters concurrently, i.e. that the
// Go->Python->Go round trip doesn't corrupt shared Go-side state
// under concurrent GIL acquisition.
//
// NOTE: callback.call() does not yet translate Python-side call
// arguments through to the Go function (see the "TODO" in
// callback.go) -- only argument-less Go functions are supported
// today. This test intentionally only uses that supported subset;
// once argument passing lands, extend this test accordingly.
func TestConcurrentCallbacks(t *testing.T) {
	workers := concurrencyWorkers()

	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			py, err := NewPython()
			if err != nil {
				errs <- fmt.Errorf("worker %d: NewPython: %w", id, err)
				return
			}
			defer py.Close()

			var calls atomic.Int64
			cb := func() int {
				return int(calls.Add(1))
			}

			err = py.Set("next_call", cb)
			if err != nil {
				errs <- fmt.Errorf("worker %d: Set callback: %w", id, err)
				return
			}

			for j := 1; j <= 50; j++ {
				obj := py.Eval("next_call()")
				if err := obj.Err(); err != nil {
					errs <- fmt.Errorf("worker %d: Eval: %w", id, err)
					return
				}

				got, err := obj.Int()
				if err != nil {
					errs <- fmt.Errorf("worker %d: Int: %w", id, err)
					return
				}

				if got != int64(j) {
					errs <- fmt.Errorf("worker %d: got %d, want %d",
						id, got, j)
					return
				}
			}

			if calls.Load() != 50 {
				errs <- fmt.Errorf("worker %d: callback invoked %d times, want 50",
					id, calls.Load())
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// TestConcurrentGCPressure interleaves object creation across many
// goroutines/interpreters with explicit runtime.GC() calls, to
// surface any finalizer/objmap race between Go's collector reclaiming
// *Object wrappers and the interpreter being torn down or used
// concurrently.
func TestConcurrentGCPressure(t *testing.T) {
	workers := concurrencyWorkers()

	var wg, bgwg sync.WaitGroup
	errs := make(chan error, workers)
	var stop atomic.Bool

	// Background goroutine forcing GC cycles throughout the test.
	bgwg.Add(1)
	go func() {
		defer bgwg.Done()
		for !stop.Load() {
			runtime.GC()
			time.Sleep(time.Millisecond)
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			py, err := NewPython()
			if err != nil {
				errs <- fmt.Errorf("worker %d: NewPython: %w", id, err)
				return
			}
			defer py.Close()

			for j := 0; j < 500; j++ {
				// Deliberately drop the reference immediately so
				// the *Object becomes eligible for GC right away.
				py.NewObject(j)
			}

			// countObjID should never panic or deadlock even while
			// finalizers may be running concurrently in the background.
			_ = py.countObjID()
		}(i)
	}

	wg.Wait()
	stop.Store(true)
	bgwg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

