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

func containsFatalPython(s string) bool {
	return strings.Contains(s, "sys.exit") ||
		strings.Contains(s, "os._exit") ||
		strings.Contains(s, "raise SystemExit")
}

func safeEval(py *Python, src string) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = py.Eval(src)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
	}
}

func safeExec(py *Python, src string) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = py.Exec(src, "")
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
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

// Interpreter fuzzing

func FuzzPythonEvalExec(f *testing.F) {
	f.Add("1 + 1")
	f.Add("print('hello')")
	f.Add("")
	f.Add("x = [1, 2, 3]\nx")
	f.Add("def f(x): return x * 2\nf(10)")
	f.Add("1/0")
	f.Add("this is not python")

	f.Fuzz(func(t *testing.T, src string) {
		if containsFatalPython(src) {
			t.Skip()
		}

		py := getFuzzPython(t)
		if py == nil {
			return
		}

		safeEval(py, src)
		safeExec(py, src)
	})
}

// Object conversion fuzzing

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
