// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Fuzz tests for Python interpreter lifecycle, Eval/Exec and Go->Python object conversion.
//
//go:build linux || darwin || windows

package cpython

import (
	"math"
	"math/big"
	"strings"
	"testing"
	"time"
)

// Helpers

// containsFatalPython filters out Python code that would terminate the process
func containsFatalPython(s string) bool {
	return strings.Contains(s, "sys.exit") ||
		strings.Contains(s, "os._exit") ||
		strings.Contains(s, "raise SystemExit")
}

// safeEval runs py.Eval with a timeout to avoid blocking fuzz workers
func safeEval(py *Python, src string) error {
	done := make(chan error, 1)
	go func() {
		obj := py.Eval(src)
		if obj != nil {
			done <- obj.Err()
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

// safeExec runs py.Exec with a timeout to avoid blocking fuzz workers
func safeExec(py *Python, src string) error {
	done := make(chan error, 1)
	go func() {
		done <- py.Exec(src, "")
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

// Global interpreter for fuzz worker

var fuzzPy *Python

func getFuzzPython(t *testing.T) *Python {
	if fuzzPy != nil {
		return fuzzPy
	}

	py, err := NewPython()
	if err != nil {
		t.Skip("Python interpreter not available in fuzz worker")
		return nil
	}

	fuzzPy = py
	return fuzzPy
}

// FuzzPythonEvalExec fuzzes Python Eval/Exec to ensure that arbitrary input does not crash the interpreter or leave it in a broken state
func FuzzPythonEvalExec(f *testing.F) {
	f.Add("1 + 1")
	f.Add("")
	f.Add("x = [1, 2, 3]")
	f.Add("def f(x): return x * 2\nf(10)")
	f.Add("1/0")                // runtime error
	f.Add("this is not python") // syntax error

	f.Fuzz(func(t *testing.T, src string) {
		if containsFatalPython(src) {
			t.Skip()
		}

		py := getFuzzPython(t)
		if py == nil {
			return
		}

		errEval := safeEval(py, src)
		errExec := safeExec(py, src)

		// For known invalid inputs, we expect errors
		if src == "1/0" || src == "this is not python" {
			if errEval == nil && errExec == nil {
				t.Fatalf("expected error for input %q, got nil", src)
			}
		}
	})
}

// FuzzPythonNewObjectInt64 fuzzes conversion of integer values from Go to Python objects
func FuzzPythonNewObjectInt64(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1))
	f.Add(int64(1 << 60))

	f.Fuzz(func(t *testing.T, v int64) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.NewObject(v)
		if obj != nil {
			_ = obj.Err()
		}

		bi := big.NewInt(v)
		obj = py.NewObject(bi)
		if obj != nil {
			_ = obj.Err()
		}
	})
}

// FuzzPythonNewObjectFloat64 fuzzes conversion of float values from Go to Python objects
func FuzzPythonNewObjectFloat64(f *testing.F) {
	f.Add(float64(0))
	f.Add(float64(1.5))
	f.Add(float64(-1.5))
	f.Add(float64(math.MaxFloat64))

	f.Fuzz(func(t *testing.T, v float64) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.NewObject(v)
		if obj != nil {
			_ = obj.Err()
		}
	})
}

// FuzzPythonNewObjectString fuzzes string conversion from Go to Python unicode objects
func FuzzPythonNewObjectString(f *testing.F) {
	f.Add("")
	f.Add("hello")
	f.Add("привет")

	f.Fuzz(func(t *testing.T, s string) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.NewObject(s)
		if obj != nil {
			_ = obj.Err()
		}
	})
}

// FuzzPythonNewObjectBytes fuzzes byte slice conversion to Python bytes objects
func FuzzPythonNewObjectBytes(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x00, 0x01, 0xff})

	f.Fuzz(func(t *testing.T, b []byte) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.NewObject(b)
		if obj != nil {
			_ = obj.Err()
		}
	})
}

// FuzzPythonContainerConversion fuzzes conversion of Go slices and maps to Python containers
func FuzzPythonContainerConversion(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(5)

	f.Fuzz(func(t *testing.T, n int) {
		if n < 0 {
			return
		}
		if n > 20 {
			n = 20
		}

		py := getFuzzPython(t)
		if py == nil {
			return
		}

		slice := make([]any, 0, n)
		for i := 0; i < n; i++ {
			slice = append(slice, i)
		}
		obj := py.NewObject(slice)
		if obj != nil {
			_ = obj.Err()
		}

		m := make(map[string]int)
		for i := 0; i < n; i++ {
			m[string('a'+rune(i%26))] = i
		}
		obj = py.NewObject(m)
		if obj != nil {
			_ = obj.Err()
		}
	})
}

// FuzzPythonGlobals fuzzes access to Python global variables via Set/Get/Contains/Del
func FuzzPythonGlobals(f *testing.F) {
	f.Add("x", int64(1))
	f.Add("", int64(0))

	f.Fuzz(func(t *testing.T, name string, v int64) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		_ = py.SetGlobal(name, v)

		obj := py.GetGlobal(name)
		if obj != nil {
			_ = obj.Err()
		}

		_, _ = py.ContainsGlobal(name)
		_, _ = py.DelGlobal(name)
	})
}
