// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Fuzz tests for Python Object operations and lifecycle.
//
//go:build linux || darwin || windows

package cpython

import (
	"math"
	"testing"
)

// FuzzObjectLifecycle stresses object creation, early invalidation and repeated Object allocation to exercise reference management.
func FuzzObjectLifecycle(f *testing.F) {
	f.Add(1)
	f.Add(2)
	f.Add(10)

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

		objs := make([]*Object, 0, n)
		for i := 0; i < n; i++ {
			obj := py.Eval("1")
			_ = obj.Err()
			objs = append(objs, obj)
		}

		// Invalidate some objects early
		for i := 0; i < len(objs); i += 2 {
			objs[i].Invalidate()
		}

	})
}

// FuzzObjectAttributes fuzzes GetAttr/SetAttr/DelAttr/HasAttr with random names
func FuzzObjectAttributes(f *testing.F) {
	f.Add("x")
	f.Add("")
	f.Add("__class__")
	f.Add("does_not_exist")

	f.Fuzz(func(t *testing.T, name string) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.Eval(`{"a": 1}`)
		if obj.Err() != nil {
			return
		}

		// Attribute operations must never panic
		_, _ = obj.HasAttr(name)
		attr := obj.GetAttr(name)
		_ = attr.Err()
		_, _ = obj.DelAttr(name)
		_ = obj.SetAttr(name, 1)
	})
}

// FuzzObjectItems fuzzes Get/Set/Del/Contains on container objects
func FuzzObjectItems(f *testing.F) {
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

		obj := py.Eval(`{}`)
		if obj.Err() != nil {
			return
		}

		// Insert values
		for i := 0; i < n; i++ {
			_ = obj.Set(i, i)
		}

		// Randomly query and delete
		for i := 0; i < n; i++ {
			_, _ = obj.Contains(i)
			item := obj.Get(i)
			_ = item.Err()
			_, _ = obj.Del(i)
			_, _ = obj.Del(i)
		}
	})
}

// FuzzObjectCall fuzzes calling objects as functions
func FuzzObjectCall(f *testing.F) {
	f.Add(1)
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, n int) {
		if n < 0 {
			return
		}
		if n > 10 {
			n = 10
		}

		py := getFuzzPython(t)
		if py == nil {
			return
		}

		fn := py.Eval("min")
		if fn.Err() != nil {
			return
		}

		for i := 0; i < n; i++ {
			res := fn.Call(i, i+1)
			_ = res.Err()

			res = fn.CallKW(map[string]any{"default": i}, []int{})
			_ = res.Err()
		}
	})
}

// FuzzObjectDecoders fuzzes type conversion accessors
func FuzzObjectDecoders(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1))
	f.Add(int64(42))

	f.Fuzz(func(t *testing.T, v int64) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.NewObject(v)
		if obj.Err() != nil {
			return
		}

		_, _ = obj.Int()
		_, _ = obj.Uint()
		_, _ = obj.Float()
		_, _ = obj.Complex()
		_, _ = obj.Unicode()
		_, _ = obj.Bytes()
		_, _ = obj.Bool()
	})
}

// FuzzObjectTypeChecks fuzzes Is* predicates
func FuzzObjectTypeChecks(f *testing.F) {
	f.Add("None")
	f.Add("[]")
	f.Add("{}")
	f.Add("1")
	f.Add("1.5")
	f.Add("b'abc'")
	f.Add("bytearray(b'abc')")
	f.Add("min")

	f.Fuzz(func(t *testing.T, expr string) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.Eval(expr)
		if obj.Err() != nil {
			return
		}

		_ = obj.IsBool()
		_ = obj.IsByteArray()
		_ = obj.IsBytes()
		_ = obj.IsCallable()
		_ = obj.IsComplex()
		_ = obj.IsDict()
		_ = obj.IsFloat()
		_ = obj.IsLong()
		_ = obj.IsNone()
		_ = obj.IsSeq()
		_ = obj.IsUnicode()
	})
}

// FuzzObjectReprStr fuzzes Str() and Repr() for stability
func FuzzObjectReprStr(f *testing.F) {
	f.Add("1")
	f.Add("'hello'")
	f.Add("[]")
	f.Add("{}")

	f.Fuzz(func(t *testing.T, expr string) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.Eval(expr)
		if obj.Err() != nil {
			return
		}

		_, _ = obj.Str()
		_, _ = obj.Repr()
	})
}

// FuzzObjectNumericEdgeCases fuzzes numeric edge cases that often break bindings
func FuzzObjectNumericEdgeCases(f *testing.F) {
	f.Add(float64(0))
	f.Add(float64(math.MaxFloat64))
	f.Add(float64(math.SmallestNonzeroFloat64))

	f.Fuzz(func(t *testing.T, v float64) {
		py := getFuzzPython(t)
		if py == nil {
			return
		}

		obj := py.NewObject(v)
		if obj.Err() != nil {
			return
		}

		_, _ = obj.Float()
		_, _ = obj.Int()
		_, _ = obj.Uint()
	})
}
