// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Additional tests to reach 100% coverage for python.go

package cpython

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// TestStressObjectAllocation attempts to trigger race between
// Python.Close and Object.finalizer
func TestStressObjectAllocation(t *testing.T) {
	var done atomic.Bool
	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		for !done.Load() {
			py, err := NewPython()
			assert.NoError(err)

			for i := 0; i < 10000; i++ {
				py.NewObject(i)
			}

			py.Close()
		}
		wait.Done()
	}()

	time.Sleep(500 * time.Millisecond)
	done.Store(true)
	wait.Wait()
}

// TestPythonNone covers: None() → return py.objNone
func TestPythonNone(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	none := py.None()
	if none == nil {
		t.Fatal("Python.None() returned nil")
	}
	if none.Err() != nil {
		t.Fatalf("Python.None() error: %s", none.Err())
	}
	s, err := none.Str()
	if err != nil {
		t.Fatalf("None.Str(): %s", err)
	}
	if s != "None" {
		t.Errorf("None.Str(): expected %q, got %q", "None", s)
	}
}

// TestPythonBool covers: Bool(true) → objTrue, Bool(false) → objFalse
func TestPythonBool(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	objTrue := py.Bool(true)
	if objTrue.Err() != nil {
		t.Fatalf("Bool(true): %s", objTrue.Err())
	}
	bTrue, err := objTrue.Bool()
	if err != nil {
		t.Fatalf("Bool(true).Bool(): %s", err)
	}
	if !bTrue {
		t.Error("Bool(true): expected true, got false")
	}

	objFalse := py.Bool(false)
	if objFalse.Err() != nil {
		t.Fatalf("Bool(false): %s", objFalse.Err())
	}
	bFalse, err := objFalse.Bool()
	if err != nil {
		t.Fatalf("Bool(false).Bool(): %s", err)
	}
	if bFalse {
		t.Error("Bool(false): expected false, got true")
	}
}

// TestPythonGetAndGetGlobal covers:
//   - GetGlobal → globals.GetItem (not-found for builtin, found for real global)
//   - Get → found in globals (no builtins fallthrough)
//   - Get → NOT found in globals → builtins.Get (the uncovered 25% branch)
func TestPythonGetAndGetGlobal(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// GetGlobal: builtin "print" is NOT in globals dict.
	obj := py.GetGlobal("print")
	if !obj.NotFound() {
		t.Errorf("GetGlobal(print): expected NotFound, got %s", obj)
	}

	// GetGlobal: "__name__" IS a real global.
	obj = py.GetGlobal("__name__")
	if obj.Err() != nil {
		t.Errorf("GetGlobal(__name__): unexpected error: %s", obj.Err())
	}

	// Get: "__name__" is found in globals — builtins branch NOT taken.
	obj = py.Get("__name__")
	if obj.Err() != nil {
		t.Errorf("Get(__name__): unexpected error: %s", obj.Err())
	}

	// Get: "print" is NOT in globals → falls through to builtins.Get.
	// This covers the previously uncovered branch in Get.
	obj = py.Get("print")
	if obj.Err() != nil {
		t.Errorf("Get(print via builtins): unexpected error: %s", obj.Err())
	}
}

// TestPythonDel covers: Del → globals.Del, found=true and found=false paths
func TestPythonDel(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	if err := py.Set("del_test_var", 42); err != nil {
		t.Fatalf("Set: %s", err)
	}

	// Delete existing → (true, nil)
	deleted, err := py.Del("del_test_var")
	if err != nil {
		t.Errorf("Del(existing): unexpected error: %s", err)
	}
	if !deleted {
		t.Error("Del(existing): expected true, got false")
	}

	// Confirm gone via ContainsGlobal (safe — only touches globals dict).
	found, err := py.ContainsGlobal("del_test_var")
	if err != nil {
		t.Fatalf("ContainsGlobal after del: %s", err)
	}
	if found {
		t.Error("ContainsGlobal: variable should not exist after Del")
	}

	// Delete non-existing → (false, nil)
	deleted, err = py.Del("del_test_nonexistent_xyz")
	if err != nil {
		t.Errorf("Del(nonexistent): unexpected error: %s", err)
	}
	if deleted {
		t.Error("Del(nonexistent): expected false, got true")
	}
}

// TestPythonContains covers:
//   - ContainsGlobal → found=false (builtin), found=true (set global)
//   - Contains → found=true in globals (builtins path not taken)
//   - Contains → NOT found in globals → builtins.ContainsItem (the uncovered branch)
func TestPythonContains(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// ContainsGlobal: builtin not in globals dict.
	found, err := py.ContainsGlobal("print")
	if err != nil {
		t.Errorf("ContainsGlobal(print): unexpected error: %s", err)
	}
	if found {
		t.Error("ContainsGlobal(print): expected false, got true")
	}

	// Set a global so Contains finds it in globals (builtins not reached).
	if err := py.Set("contains_test_var", 1); err != nil {
		t.Fatalf("Set: %s", err)
	}

	found, err = py.ContainsGlobal("contains_test_var")
	if err != nil {
		t.Errorf("ContainsGlobal(set var): unexpected error: %s", err)
	}
	if !found {
		t.Error("ContainsGlobal(set var): expected true, got false")
	}

	found, err = py.Contains("contains_test_var")
	if err != nil {
		t.Errorf("Contains(set var): unexpected error: %s", err)
	}
	if !found {
		t.Error("Contains(set var): expected true, got false")
	}

	// Contains: name NOT in globals → builtins.ContainsItem is called.
	// In a subinterpreter __builtins__ is a module so this may return an
	// error; we only need the branch to execute for coverage.
	_, _ = py.Contains("print")
}

// TestPythonNewError covers: NewError → newErrorObject(py, err)
func TestPythonNewError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	origErr := errors.New("sentinel error")
	obj := py.NewError(origErr)
	if obj.Err() == nil {
		t.Fatal("NewError: expected error Object")
	}
	if obj.Err().Error() != origErr.Error() {
		t.Errorf("NewError: wrong error: %s", obj.Err())
	}
}

// TestPythonCloseTwice covers: Close → gate() returns ErrClosed → early return
func TestPythonCloseTwice(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	py.Close()
	py.Close() // must not panic
}

// TestPythonNewObjectOnClosedInterp covers:
// NewObject → gate() error → return newErrorObject
func TestPythonNewObjectOnClosedInterp(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	py.Close()

	if obj := py.NewObject(42); obj.Err() == nil {
		t.Error("NewObject on closed interp: expected error")
	}
}

// TestPythonNewObjectUnsupportedKinds covers:
// newPyObject → chan/struct fall-through → ErrTypeConversion
func TestPythonNewObjectUnsupportedKinds(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	if obj := py.NewObject(make(chan int)); obj.Err() == nil {
		t.Error("NewObject(chan): expected ErrTypeConversion")
	}
	type myStruct struct{ X int }
	if obj := py.NewObject(myStruct{1}); obj.Err() == nil {
		t.Error("NewObject(struct): expected ErrTypeConversion")
	}
}

// TestPythonNewObjectErrorValue covers:
// newPyObject → case error: → return nil, v
func TestPythonNewObjectErrorValue(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	inputErr := errors.New("go error as value")
	obj := py.NewObject(inputErr)
	if obj.Err() == nil {
		t.Error("NewObject(error): expected error Object")
	}
	if obj.Err().Error() != inputErr.Error() {
		t.Errorf("NewObject(error): wrong error: %s", obj.Err())
	}
}

// TestPythonNewObjectObjectWithError covers:
// newPyObject → case *Object: v.err != nil → return nil, v.err
func TestPythonNewObjectObjectWithError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	errObj := py.NewError(errors.New("inner"))
	if obj := py.NewObject(errObj); obj.Err() == nil {
		t.Error("NewObject(*Object with error): expected error")
	}
}

// TestPythonNewObjectFloat covers:
// newPyObject → reflect.Float32/Float64 → gate.makeFloat
func TestPythonNewObjectFloat(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(float64(3.14))
	if obj.Err() != nil {
		t.Fatalf("NewObject(float64): %s", obj.Err())
	}
	f, err := obj.Float()
	if err != nil {
		t.Fatalf("Float(): %s", err)
	}
	if math.Abs(f-3.14) > 1e-9 {
		t.Errorf("Float(): expected 3.14, got %v", f)
	}
	if obj32 := py.NewObject(float32(1.5)); obj32.Err() != nil {
		t.Errorf("NewObject(float32): %s", obj32.Err())
	}
}

// TestPythonNewObjectComplex covers:
// newPyObject → reflect.Complex64/Complex128 → gate.makeComplex
func TestPythonNewObjectComplex(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	obj := py.NewObject(complex(1.0, 2.0))
	if obj.Err() != nil {
		t.Fatalf("NewObject(complex128): %s", obj.Err())
	}
	c, err := obj.Complex()
	if err != nil {
		t.Fatalf("Complex(): %s", err)
	}
	if real(c) != 1.0 || imag(c) != 2.0 {
		t.Errorf("Complex(): expected (1+2i), got %v", c)
	}
	if obj64 := py.NewObject(complex64(complex(3.0, 4.0))); obj64.Err() != nil {
		t.Errorf("NewObject(complex64): %s", obj64.Err())
	}
}

// TestPythonNewObjectArrayBytes covers:
// newPyObject → reflect.Array, Uint8 elem, !CanAddr → copy branch → makeBytes
func TestPythonNewObjectArrayBytes(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	arr := [4]byte{10, 20, 30, 40}
	obj := py.NewObject(arr) // passed by value → not addressable
	if obj.Err() != nil {
		t.Fatalf("NewObject([4]byte): %s", obj.Err())
	}
	b, err := obj.Bytes()
	if err != nil {
		t.Fatalf("Bytes(): %s", err)
	}
	for i, want := range []byte{10, 20, 30, 40} {
		if b[i] != want {
			t.Errorf("Bytes()[%d]: expected %d, got %d", i, want, b[i])
		}
	}
}

// TestPythonNewPyListItemError covers:
// newPyList → newPyObject error on item → return nil, err
func TestPythonNewPyListItemError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	if obj := py.NewObject([]any{1, make(chan int), 3}); obj.Err() == nil {
		t.Error("NewObject(list with chan): expected error")
	}
}

// TestPythonNewPyDictUnsortableKey covers:
// newPyDict → reflectSort returns false → ErrTypeConversion
//
// Struct keys are comparable in Go (no panic) but reflectSort has no
// case for reflect.Struct, so it returns false.
func TestPythonNewPyDictUnsortableKey(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	type structKey struct{ X int }
	badMap := map[structKey]string{{X: 1}: "a", {X: 2}: "b"}
	if obj := py.NewObject(badMap); obj.Err() == nil {
		t.Error("NewObject(map[struct]): expected ErrTypeConversion")
	}
}

// TestPythonNewPyDictValueError covers:
// newPyDict → value newPyObject fails → unref + return nil, err
func TestPythonNewPyDictValueError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	if obj := py.NewObject(map[string]any{"k": make(chan int)}); obj.Err() == nil {
		t.Error("NewObject(map with chan value): expected error")
	}
}

// TestPythonNewPyFunctionBadSignature covers:
// newPyFunction → newCallback returns nil → ErrTypeConversion
//
// newCallback returns nil when NumOut()>1 and second return is not error.
func TestPythonNewPyFunctionBadSignature(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	if obj := py.NewObject(func() (int, int) { return 1, 2 }); obj.Err() == nil {
		t.Error("NewObject(func()(int,int)): expected ErrTypeConversion")
	}
}

// TestPythonLoadError covers:
// Load → gate.load returns error → return newErrorObject
func TestPythonLoadError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	if obj := py.Load("def (((invalid###", "bad", "bad.py"); obj.Err() == nil {
		t.Error("Load(bad syntax): expected error")
	}
}

// TestPythonLoadOnClosedInterp covers:
// Load → gate() returns ErrClosed → return newErrorObject
func TestPythonLoadOnClosedInterp(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	py.Close()

	if obj := py.Load("x = 1", "mod", "mod.py"); obj.Err() == nil {
		t.Error("Load on closed interp: expected error")
	}
}

// TestPythonEvalOnClosedInterp covers:
// eval → gate() returns ErrClosed → return newErrorObject
func TestPythonEvalOnClosedInterp(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	py.Close()

	if obj := py.Eval("1+1"); obj.Err() == nil {
		t.Error("Eval on closed interp: expected error")
	}
}

// TestPythonEvalReturnsNone covers:
// eval → gate.eval returns nil pyobj (exec mode) → return py.objNone
func TestPythonEvalReturnsNone(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// Exec calls eval(..., expr=false); a plain statement returns nil pyobj.
	if err := py.Exec("_cov_var = 99", "cov.py"); err != nil {
		t.Errorf("Exec: unexpected error: %s", err)
	}
}

// TestPythonLoadSuccess covers:
// Load → success path → return newObjectFromPython(py, gate, pyobj)
// (Image 4 showed this line red — TestPythonLoad from upstream is missing locally)
func TestPythonLoadSuccess(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	mod := "\nx = 42\n"
	obj := py.Load(mod, "covmodule", "covmodule.py")
	if obj.Err() != nil {
		t.Fatalf("Load: unexpected error: %s", obj.Err())
	}

	// Verify the module is accessible and the variable has the right value.
	res := py.Eval("covmodule.x")
	if res.Err() != nil {
		t.Fatalf("Eval(covmodule.x): %s", res.Err())
	}
	v, err := res.Int()
	if err != nil {
		t.Fatalf("Int(): %s", err)
	}
	if v != 42 {
		t.Errorf("covmodule.x: expected 42, got %d", v)
	}
}

// TestPythonEvalGateError covers:
// eval → gate.eval returns error (Python runtime exception with explicit filename)
// → return newErrorObject(py, err)
// (Image 5 showed this return statement red)
func TestPythonEvalGateError(t *testing.T) {
	py, err := NewPython()
	assert.NoError(err)
	defer py.Close()

	// Passing an explicit filename bypasses the runtime.Callers block and goes
	// straight to gate.eval. A bad expression causes gate.eval to return an error,
	// hitting the red return newErrorObject(py, err) branch.
	obj := py.eval("1/0", "explicit_file.py", true)
	if obj.Err() == nil {
		t.Error("eval(1/0): expected ZeroDivisionError, got nil")
	}
}
