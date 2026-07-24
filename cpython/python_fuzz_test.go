// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Fuzz test for Python sub-interpreter lifecycle (create/exec/close),
// with RSS-based leak detection.
//
// Runs the same script through many lifecycle cycles to catch slow
// leaks a single cycle wouldn't reveal, checking two signals:
//  1. PythonInstancesCount() returns to baseline after every Close().
//  2. Process RSS, sampled periodically, to catch native-memory leaks.
package cpython

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const (
	// Create/exec/close cycles per fuzz input.
	subinterpFuzzIterations = 120

	// How often (in iterations) RSS is sampled.
	subinterpRSSSampleEvery = 10

	// Iterations excluded from the leak-trend check to let allocator/OS
	// settle first.
	subinterpRSSWarmupIterations = 30

	// Max tolerated RSS growth (KB) before flagging a likely leak.
	subinterpRSSGrowthLimitKB = 20 * 1024 // 20 MiB
)

// FuzzSubInterpreterLifecycle fuzzes Python interpreter creation, script
// execution, and teardown, looking for leaked interpreter handles and
// slow native-memory growth.
func FuzzSubInterpreterLifecycle(f *testing.F) {
	// Seed corpus: eval vs exec, success vs error, syntax errors,
	// imports, exceptions, and object-allocating scripts.
	seeds := []struct {
		script string
		exec   bool
	}{
		{"", false},
		{"1 + 1", false},
		{"[x for x in range(1000)]", false},
		{"import sys\nsys.version", true},
		{"raise ValueError('boom')", true},
		{"def f():\n    raise KeyError('x')\nf()", true},
		{"class C:\n    def __init__(self):\n        self.data = [0] * 1000\nc = C()", true},
		{"(", false}, // syntax error
		{"import os\nos.getpid()", true},
		{"while False:\n    pass", true},
	}
	for _, s := range seeds {
		f.Add(s.script, s.exec)
	}

	f.Fuzz(func(t *testing.T, script string, exec bool) {
		var baselineRSS, lastRSS uint64
		var haveBaseline bool

		for i := 0; i < subinterpFuzzIterations; i++ {
			before := PythonInstancesCount()

			py, err := NewPython()
			if err != nil {
				// Not the bug we're hunting; skip and keep looping.
				continue
			}

			if exec {
				_ = py.Exec(script, "fuzz.py")
			} else {
				_ = py.Eval(script)
			}

			py.Close()

			after := PythonInstancesCount()
			if after != before {
				t.Fatalf(
					"interpreter instance count leaked: before=%d after=%d (script=%q exec=%v)",
					before, after, script, exec)
			}

			if i%subinterpRSSSampleEvery == 0 && i >= subinterpRSSWarmupIterations {
				runtime.GC()
				rss, ok := readRSSKB()
				if !ok {
					continue // /proc unavailable on this platform
				}
				if !haveBaseline {
					baselineRSS = rss
					haveBaseline = true
				}
				lastRSS = rss
			}
		}

		if haveBaseline && lastRSS > baselineRSS {
			growth := lastRSS - baselineRSS
			if growth > subinterpRSSGrowthLimitKB {
				t.Fatalf(
					"possible memory leak: RSS grew by %d KB over %d iterations (script=%q exec=%v)",
					growth, subinterpFuzzIterations-subinterpRSSWarmupIterations,
					script, exec)
			}
		}
	})
}

// readRSSKB reads current RSS (KB) from /proc/self/status.
// Returns ok=false if unavailable (e.g. non-Linux).
func readRSSKB() (kb uint64, ok bool) {
	file, err := os.Open("/proc/self/status")
	if err != nil {
		return 0, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "VmRSS:") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			return 0, false
		}

		v, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0, false
		}

		return v, true
	}

	return 0, false
}
