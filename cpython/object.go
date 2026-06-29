// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Python objects

package cpython

import (
	"errors"
	"math/big"
	"runtime"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// Object represents a Python value or Python error.
//
// All Python interpreter calls that return values encapsulate them
// into an Object. If a Python call fails, the error is also returned
// encapsulated into an Object (referred to as an "error Object").
//
// Any operation performed on an error Object will return the same error
// as the error Object contains. If an operation involves multiple
// Objects (e.g., [Object.Call]), the error is inherited from the first
// Object used, starting with the receiver.
//
// This enables chaining operations with error propagation through the chain.
//
// The error contained in an error Object is finally exposed when attempting
// to convert the Object into a Go value (using functions like [Object.Int])
// or can be inspected directly via the [Object.Err] method.
//
// Normally, Object lifetime is managed by Go's garbage collector,
// so explicit release is unnecessary. However, Python values may be
// released early for optimization by calling [Object.Invalidate].
type Object struct {
	py  *Python // interpreter that owns this Object
	err error   // error value for error Object
	oid objid   // object ID of the underlying *C.PyObject
}

// newObjectFromPython constructs new Object, decoded from PyObject.
func newObjectFromPython(py *Python, gate pyGate, pyobj pyObject) *Object {
	assert.Must(pyobj != nil)
	obj := &Object{
		py:  py,
		oid: py.newObjID(gate, pyobj),
	}

	runtime.SetFinalizer(obj, func(obj *Object) {
		obj.finalizer()
	})

	return obj
}

// newErrorObject constructs new Object, that represents the error
func newErrorObject(py *Python, err error) *Object {
	return &Object{
		py:  py,
		err: err,
		oid: 0,
	}
}

// finalizer is called when Object is garbage-collected.
// It released *C.PyObject, associated with the Object.
func (obj *Object) finalizer() {
	gate, err := obj.py.gate()
	if err == nil {
		obj.py.delObjID(gate, obj.oid)
		gate.release()
	}
}

// begin starts operations with Object by obtaining the pyGate
// and *C.PyObject reference to the Object's value at the
// cpython side.
func (obj *Object) begin() (pyGate, pyObject, error) {
	if obj.err != nil {
		return pyGate{}, nil, obj.err
	}

	gate, err := obj.py.gate()

	if err == nil {
		var pyobj pyObject
		pyobj, err = obj.py.lookupObjID(gate, obj.oid)
		if pyobj != nil {
			return gate, pyobj, nil
		}

		gate.release()
	}

	return pyGate{}, nil, err
}

// Invalidate marks Object as invalid and releases Python value
// that it might have.
//
// This is optimization call. You actually don't have to manage
// Object lifetime manually. Go GC does it automatically. However,
// this call may be useful to release Python value early, especially
// if the value is heavyweight at the Python side.
func (obj *Object) Invalidate() {
	gate, err := obj.py.gate()
	if err != nil {
		return
	}

	obj.py.delObjID(gate, obj.oid)
	gate.release()

	runtime.SetFinalizer(obj, nil) // Finalizer is not required anymore
}

// Py returns the [Python] interpreter Object belongs to.
func (obj *Object) Py() *Python {
	return obj.py
}

// Err returns error for the error Object, or nil if there
// is no error.
func (obj *Object) Err() error {
	return obj.err
}

// NotFound returns true, if [Object] is the error object
// and error type is [ErrNotFound].
func (obj *Object) NotFound() bool {
	return errors.Is(obj.err, ErrNotFound{})
}

// Len returns Object length, in items. It works with container
// objects (lists, tuples, dict, ...).
//
// In Python:
//
//	len(obj)
func (obj *Object) Len() (int, error) {
	return objDo(obj, pyGate.length)
}

// Save sets the global variable 'name' to the Object's value:
//
//	name = obj
//
// This is the same as obj.Py().Set(name, obj).
// This method is useful for operations chaining with error propagation:
//
//	// x = (2 * 3 + 5) % 7
//	err = py.NewObject(2).Mul(3).Add(5).Mod(7).Save("x")
func (obj *Object) Save(name string) error {
	return obj.py.Set(name, obj)
}

// SaveTo saves the global variable 'name' to the dest Object's attribute:
//
//	dest.name = obj
//
// See also [Object.Save].
func (obj *Object) SaveTo(dest *Object, name string) error {
	return dest.Set(name, obj)
}

// SaveItem saves the global variable 'name' to the dest's item with
// the specified key.
//
//	dest[key] = obj
//
// See also [Object.Save].
func (obj *Object) SaveItem(dest *Object, key any) error {
	return dest.SetItem(key, obj)
}

// String returns string representation of the Object, for debugging
// purposes.
//
// It tries [Object.Str], with fallback to obj.Err().Error()
func (obj *Object) String() string {
	s, err := obj.Str()
	if err != nil {
		s = err.Error()
	}
	return s
}

// Del deletes Object item with the specified key:
//
// In Python:
//
//	del(obj, key)
//
// The Object must be container (array, dict, etc).
// The key may be any value that [Python.NewObject] accepts.
func (obj *Object) Del(key any) (bool, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return false, err
	}
	defer gate.release()

	// Obtain *C.PyObject reference for key
	pykey, err := obj.py.newPyObject(gate, key)
	if err != nil {
		return false, err
	}
	defer gate.unref(pykey)

	// Check for item existence, then delete, if found.
	found, err := gate.hasitem(pyobj, pykey)
	if found {
		err = gate.delitem(pyobj, pykey)
	}

	return found, err
}

// GetItem returns Object item with the specified key:
//
//	ret = obj[key]
//
// The Object must be container (array, dict, etc).
// The key may be any value that [Python.NewObject] accepts.
func (obj *Object) GetItem(key any) *Object {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return newErrorObject(obj.py, err)
	}
	defer gate.release()

	// Obtain *C.PyObject references for key
	pykey, err := obj.py.newPyObject(gate, key)
	if err != nil {
		return newErrorObject(obj.py, err)
	}
	defer gate.unref(pykey)

	// Check if item exists, then retrieve
	found, err := gate.hasitem(pyobj, pykey)
	var pyitem pyObject
	if found {
		pyitem, err = gate.getitem(pyobj, pykey)
	} else if err == nil {
		name, _ := gate.str(pykey)
		err = ErrNotFound{name}
	}

	if err != nil {
		return newErrorObject(obj.py, err)
	}

	// Create the item object
	return newObjectFromPython(obj.py, gate, pyitem)
}

// ContainsItem reports if Object has the item with the specified key:
//
//	ret = key in obj
//
// The Object must be container (array, dict, etc).
// The key may be any value that [Python.NewObject] accepts.
func (obj *Object) ContainsItem(key any) (bool, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return false, err
	}
	defer gate.release()

	// Obtain *C.PyIbject references for key
	pykey, err := obj.py.newPyObject(gate, key)
	if err != nil {
		return false, err
	}
	defer gate.unref(pykey)

	// Check for item existence
	return gate.hasitem(pyobj, pykey)
}

// SetItem sets Object item with the specified name:
//
//	obj[name] = val
//
// The Object must be container (array, dict, etc).
// The key and val may be any value that [Python.NewObject] accepts.
func (obj *Object) SetItem(key, val any) error {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return err
	}
	defer gate.release()

	// Obtain *C.PyObject references for key and val
	pykey, err := obj.py.newPyObject(gate, key)
	if err != nil {
		return err
	}
	defer gate.unref(pykey)

	pyval, err := obj.py.newPyObject(gate, val)
	if err != nil {
		return err
	}
	defer gate.unref(pyval)

	// Set the item
	return gate.setitem(pyobj, pykey, pyval)
}

// DelAttr deletes Object attribute with the specified name:
//
//	delattr(obj, name)
//
// It returns:
//   - (true, nil) if attribute was found and deleted
//   - (false, nil) if attribute was not found
//   - (false, error) in a case of error
func (obj *Object) DelAttr(name string) (bool, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return false, err
	}
	defer gate.release()

	// Delete the attribute. If attribute didn't exist,
	// don't threat it as an error.
	err = gate.delattr(pyobj, name)
	if errors.Is(err, AttributeError) {
		return false, nil
	}

	return err == nil, err
}

// Get returns Object attribute with the specified name:
//
//	ret = obj.name
func (obj *Object) Get(name string) *Object {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return newErrorObject(obj.py, err)
	}
	defer gate.release()

	// Fetch the attribute
	pyattr, err := gate.getattr(pyobj, name)
	if err != nil {
		if errors.Is(err, AttributeError) {
			err = ErrNotFound{name}
		}

		return newErrorObject(obj.py, err)
	}

	// Create the attribute object
	return newObjectFromPython(obj.py, gate, pyattr)
}

// HasAttr reports if Object has the attribute with the specified name.
//
//	hasattr(obj, name)
func (obj *Object) HasAttr(name string) (bool, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return false, err
	}
	defer gate.release()

	return gate.hasattr(pyobj, name)
}

// Set sets Object attribute with the specified name:
//
//	obj.name = val
//
// The val may be any value that [Python.NewObject] accepts.
func (obj *Object) Set(name string, val any) error {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return err
	}
	defer gate.release()

	// Obtain *C.PyIbject references val
	pyval, err := obj.py.newPyObject(gate, val)
	if err != nil {
		return err
	}
	defer gate.unref(pyval)

	// Set the attribute
	return gate.setattr(pyobj, name, pyval)
}

// Call calls Object as function.
//
// Arguments are automatically converted from Go to Python.
// See [Python.NewObject] for details.
//
// Use [Object.CallKW] for call with keyword arguments.
func (obj *Object) Call(args ...any) *Object {
	return obj.CallKW(nil, args...)
}

// CallKW calls Object as function with keyword arguments defined
// by the kw parameter and positional arguments defined by the
// args parameter (variadic).
//
// Arguments are automatically converted from Go to Python.
// See [Python.NewObject] for details.
//
// If keyword arguments are not used, kw may be nil.
//
// It returns the function's return value.
func (obj *Object) CallKW(kw map[string]any, args ...any) *Object {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return newErrorObject(obj.py, err)
	}
	defer gate.release()

	// Convert positional arguments
	pyargs, err := gate.makeTuple(len(args))
	if err != nil {
		return newErrorObject(obj.py, err)
	}

	defer gate.unref(pyargs)

	for i, arg := range args {
		pyarg, err := obj.py.newPyObject(gate, arg)
		if err != nil {
			return newErrorObject(obj.py, err)
		}

		err = gate.setTupleItem(pyargs, pyarg, i)
		if err != nil {
			gate.unref(pyarg)
			return newErrorObject(obj.py, err)
		}
	}

	// Convert keyword arguments
	var pykwargs pyObject
	if len(kw) > 0 {
		var err error
		pykwargs, err = obj.py.newPyObject(gate, kw)
		if err != nil {
			return newErrorObject(obj.py, err)
		}

		defer gate.unref(pykwargs)
	}

	// Perform a call
	pyret, err := gate.call(pyobj, pyargs, pykwargs)
	if err != nil {
		return newErrorObject(obj.py, err)
	}

	// Create response Object
	return newObjectFromPython(obj.py, gate, pyret)
}

// Str returns string representation of the Object.
// This is the equivalent of the Python expression str(x).
func (obj *Object) Str() (string, error) {
	return objDo(obj, pyGate.str)
}

// Repr returns string representation of the Object.
// This is the equivalent of the Python expression repr(x).
func (obj *Object) Repr() (string, error) {
	return objDo(obj, pyGate.repr)
}

// Bigint returns Object value as *[big.Int] or an error.
func (obj *Object) Bigint() (*big.Int, error) {
	return objDo(obj, pyGate.decodeBigint)
}

// Bool returns Object value as bool or an error.
func (obj *Object) Bool() (bool, error) {
	// Try the fast conversion
	val, err := obj.fastBool()
	if err == nil {
		return val, nil
	}

	// Try to call the __bool__ method
	toBool := obj.Get("__bool__")
	if toBool.Err() == nil {
		val, err2 := toBool.Call().fastBool()
		if err2 == nil {
			return val, nil
		}
	}

	return false, err
}

// fastBool returns Object value as bool or an error.
// This is the fast path for Obkect.Bool.
func (obj *Object) fastBool() (bool, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return false, err
	}
	defer gate.release()

	switch pyobj {
	case obj.py.pyTrue:
		return true, nil
	case obj.py.pyFalse:
		return false, nil
	}

	return false, gate.decodeError(pyobj, "bool")
}

// Bytes returns Object value as []byte slice or an error.
func (obj *Object) Bytes() ([]byte, error) {
	return objDo(obj, pyGate.decodeBytes)
}

// Complex returns Object value as complex128 number or an error.
func (obj *Object) Complex() (complex128, error) {
	return objDo(obj, pyGate.decodeComplex)
}

// Float returns Object value as float64 number or an error.
func (obj *Object) Float() (float64, error) {
	return objDo(obj, pyGate.decodeFloat)
}

// Int returns Object value as int64 number or an error.
func (obj *Object) Int() (int64, error) {
	return objDo(obj, pyGate.decodeInt64)
}

// Keys returns Object mapping keys as the []*Object slice.
// The Object must support mapping.
func (obj *Object) Keys() ([]*Object, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return nil, err
	}
	defer gate.release()

	// Obtain keys as list (or tuple)
	pykeys, err := gate.keys(pyobj)
	if err != nil {
		return nil, err
	}

	defer gate.unref(pykeys)

	// Convert into []*Object slice
	return objSlice(obj.py, gate, pykeys)
}

// Slice returns Object value as []*Object slice or an error.
// It works with sequence objects (lists, tuples, ...).
func (obj *Object) Slice() ([]*Object, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return nil, err
	}
	defer gate.release()

	return objSlice(obj.py, gate, pyobj)
}

// Uint returns Object value as uint64 number or an error.
func (obj *Object) Uint() (uint64, error) {
	return objDo(obj, pyGate.decodeUint64)
}

// Unicode returns Object value as UNICODE string or an error.
func (obj *Object) Unicode() (string, error) {
	return objDo(obj, pyGate.decodeUnicode)
}

// TypeName returns name of the Object's Python type
func (obj *Object) TypeName() string {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return ""
	}
	defer gate.release()

	return gate.typename(pyobj)
}

// TypeModuleName returns name of the module where Object's Python type
// is defined.
func (obj *Object) TypeModuleName() string {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return ""
	}
	defer gate.release()

	return gate.typemodulename(pyobj)
}

// IsCallable reports if Object is callable.
func (obj *Object) IsCallable() bool {
	return objDoNoError(obj, pyGate.callable)
}

// IsBool reports if Object is boolean.
func (obj *Object) IsBool() bool {
	return objDoNoError(obj, pyGate.isBool)
}

// IsByteArray reports if Object is bytearray.
// Python bytearray objects are writable sequences of bytes.
// See also [Object.IsBytes].
func (obj *Object) IsByteArray() bool {
	return objDoNoError(obj, pyGate.isByteArray)
}

// IsBytes reports if Object is bytearray.
// Python bytearray objects are immutable sequences of bytes.
// See also [Object.IsByteArray].
func (obj *Object) IsBytes() bool {
	return objDoNoError(obj, pyGate.isBytes)
}

// IsComplex reports if Object is the complex number.
func (obj *Object) IsComplex() bool {
	return objDoNoError(obj, pyGate.isComplex)
}

// IsFloat reports if Object is the floating point number.
func (obj *Object) IsFloat() bool {
	return objDoNoError(obj, pyGate.isFloat)
}

// IsDict reports if Object is dict or similar.
func (obj *Object) IsDict() bool {
	return objDoNoError(obj, pyGate.isDict)
}

// IsLong reports if Object is the integer number.
func (obj *Object) IsLong() bool {
	return objDoNoError(obj, pyGate.isLong)
}

// IsNone reports if Object is Python None.
func (obj *Object) IsNone() bool {
	gate, pyobj, err := obj.begin()
	if err != nil {
		return false
	}
	defer gate.release()

	return pyobj == obj.py.pyNone
}

// IsSeq reports if Object is sequence (i.e., list, tuple, ...).
func (obj *Object) IsSeq() bool {
	return objDoNoError(obj, pyGate.isSeq)
}

// IsUnicode reports if Object is the Unicode string
func (obj *Object) IsUnicode() bool {
	return objDoNoError(obj, pyGate.isUnicode)
}

// objDo is the convenience wrapper for the pyGate methods
// with the following signature:
//
//	pyGate.method(pyObject) (T, error)
//
// It acquires the pyGate for the Python interpreter that owns
// the Object, calls the method, releases the gate and return
// result.
func objDo[T any](obj *Object, f func(pyGate, pyObject) (T, error)) (T, error) {
	gate, pyobj, err := obj.begin()
	if err != nil {
		var t T
		return t, err
	}
	defer gate.release()

	return f(gate, pyobj)
}

// objDoNoError is the convenience wrapper for the pyGate methods
// with the following signature:
//
//	pyGate.method(pyObject) T
//
// It acquires the pyGate for the Python interpreter that owns
// the Object, calls the method, releases the gate and return
// result.
func objDoNoError[T any](obj *Object, f func(pyGate, pyObject) T) T {
	gate, pyobj, err := obj.begin()
	if err != nil {
		var t T
		return t
	}
	defer gate.release()

	return f(gate, pyobj)
}

// objSlice is the helper function that extracts contained objects from the
// sequence object and converts extracted objects into the []*Object slice.
func objSlice(py *Python, gate pyGate, pyobj pyObject) ([]*Object, error) {
	// Check that object is sequence
	if !gate.isSeq(pyobj) {
		return nil, gate.decodeError(pyobj, "[]*Object")
	}

	// Obtain length
	length, err := gate.length(pyobj)
	if err != nil {
		return nil, err
	}

	// Obtain items
	pyobjects := make([]pyObject, length)
	for i := range pyobjects {
		pyobjects[i], err = gate.getSeqItem(pyobj, i)
		if err != nil {
			for j := 0; j < i; j++ {
				gate.unref(pyobjects[j])
			}
			return nil, err
		}
	}

	// Convert into []*Object
	objects := make([]*Object, length)
	for i := range pyobjects {
		objects[i] = newObjectFromPython(py, gate, pyobjects[i])
	}

	return objects, nil
}
