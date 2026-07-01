// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Call gate into the Python interpreter

package cpython

import (
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"unsafe"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// #include "cpython.h"
import "C"

// pyGate represents the locked (attached to the current thread
// and with the GIL acquired) state of the Python interpreter.
//
// It works as a call gate into the interpreter and implements
// all interpreter operations that require locking.
type pyGate struct{}

// pyGateAcquire temporary attaches the calling thread to the
// Python interpreter.
//
// The interpreter is represented by the reference to its
// main thread state.
//
// It returns the pyGate object, that must be released after
// use with the [pyGate.release] call.
func pyGateAcquire(interp pyInterpState) pyGate {
	runtime.LockOSThread()
	C.py_enter(interp)
	return pyGate{}
}

// release detaches the calling thread from the Python interpreter.
func (gate pyGate) release() {
	gate.lastError() // Reset pending error condition, if any
	C.py_leave()
	runtime.UnlockOSThread()
}

// lastError returns a last error, nil if none.
func (gate pyGate) lastError() error {
	return gate.lastErrorAt("", -1)
}

// lastErrorAt returns a last error, nil if none.
//
// If file != "" and line > 0, it overrides the error location information.
func (gate pyGate) lastErrorAt(file string, line int) error {
	var etype, evalue, trace pyObject

	// Fetch Python error information
	C.py_err_fetch(&etype, &evalue, &trace)
	if etype == nil && evalue == nil && trace == nil {
		return nil
	}

	defer C.py_obj_unref(etype)
	defer C.py_obj_unref(evalue)
	defer C.py_obj_unref(trace)

	// Decode the error
	except := SystemError // The default
	msg := "Unknown Python exception"
	var needLocation bool

	if etype != nil {
		var nm string
		if tmp, _ := gate.getattr(etype, "__name__"); tmp != nil {
			nm, _ := gate.str(tmp)
			C.py_obj_unref(tmp)
			if nm != "" {
				except = Except(nm)
			}
		}

		// Some exception types, like SyntaxError, already
		// come with the file:line information. Others require
		// additional effort...
		switch nm {
		case "SyntaxError":
		default:
			needLocation = true
		}
	}

	if evalue != nil {
		s, _ := gate.str(evalue)
		if s != "" {
			msg = s
		}
	}

	if needLocation && trace != nil {
		var err error
		if file == "" || line < 0 {
			file, line, err = gate.lastErrorLocation(trace)
		}
		if err == nil {
			msg += fmt.Sprintf(" (%s, line %d)", file, line)
		}
	}

	return ErrPython{except, msg}
}

// lastErrorLocation extracts file and line information out of the
// traceback object.
func (gate pyGate) lastErrorLocation(trace pyObject) (
	file string, line int, err error) {

	var lineno, frame, code, filename pyObject

	// follow trace = trace.tb_next, until latest frame is reached.
	for {
		var next pyObject
		next, err = gate.getattr(trace, "tb_next")
		if err != nil {
			return
		}
		defer C.py_obj_unref(next)

		// Note, we don't have here a convenient access to
		// the None object to compare, so just run the loop
		// while trace.tb_next is of the same type as trace.
		if C.py_obj_type(next) != C.py_obj_type(trace) {
			break
		}

		trace = next
	}

	// lineno = trace.tb_lineno
	lineno, err = gate.getattr(trace, "tb_lineno")
	if err != nil {
		return
	}
	defer C.py_obj_unref(lineno)

	// frame = trace.tb_frame
	frame, err = gate.getattr(trace, "tb_frame")
	if err != nil {
		return
	}
	defer C.py_obj_unref(frame)

	// code = frame.f_code
	code, err = gate.getattr(frame, "f_code")
	if err != nil {
		return
	}
	defer C.py_obj_unref(code)

	// filename = code.co_filename
	filename, err = gate.getattr(code, "co_filename")
	if err != nil {
		return
	}
	defer C.py_obj_unref(filename)

	// Now convert filename and lineno from Python to go
	file, err = gate.decodeUnicode(filename)
	if err != nil {
		return
	}

	n, err := gate.decodeUint64(lineno)
	if err != nil {
		return
	}

	line = int(n)
	return
}

// objOrLastError returns pyobj, if it is not nil, or gate.lastError()
func (gate pyGate) objOrLastError(pyobj pyObject) (pyObject, error) {
	if pyobj != nil {
		return pyobj, nil
	}

	return nil, gate.lastError()
}

// ref increments PyObject's reference count.
func (gate pyGate) ref(pyobj pyObject) {
	C.py_obj_ref(pyobj)
}

// unref decrements PyObject's reference count.
func (gate pyGate) unref(pyobj pyObject) {
	C.py_obj_unref(pyobj)
}

// str returns str(pyobj), decoded as Go string.
func (gate pyGate) str(pyobj pyObject) (s string, err error) {
	str := C.py_obj_str(pyobj)
	if str != nil {
		defer C.py_obj_unref(str)
		s, err = gate.decodeUnicode(str)
	}

	err = gate.lastError()
	return
}

// repr returns repr(pyobj), decoded as Go string.
func (gate pyGate) repr(pyobj pyObject) (s string, err error) {
	repr := C.py_obj_repr(pyobj)
	if repr != nil {
		defer C.py_obj_unref(repr)
		s, err = gate.decodeUnicode(repr)
	}

	err = gate.lastError()
	return
}

// typename returns name of the PyObject's type
func (gate pyGate) typename(pyobj pyObject) string {
	module := "unknown"
	name := "unknown"

	t := pyObject(unsafe.Pointer(C.py_obj_type(pyobj)))

	if tmp, err := gate.getattr(t, "__module__"); err == nil {
		s, err := gate.str(tmp)
		C.py_obj_unref(tmp)

		if err == nil {
			module = s
		}
	}

	if tmp, err := gate.getattr(t, "__name__"); err == nil {
		s, err := gate.str(tmp)
		C.py_obj_unref(tmp)

		if err == nil {
			name = s
		}
	}

	name, _ = strings.CutPrefix(module+"."+name, "builtins.")
	return name
}

// typemodulename returns name of the module where PyObject's
// type is defined.
func (gate pyGate) typemodulename(pyobj pyObject) string {
	name := "unknown"

	t := pyObject(unsafe.Pointer(C.py_obj_type(pyobj)))
	if tmp, err := gate.getattr(t, "__module__"); err == nil {
		s, err := gate.str(tmp)
		C.py_obj_unref(tmp)

		if err == nil {
			name = s
		}
	}

	return name
}

// length returns length of container object (list, tuple, dict, ...)
// in items.
func (gate pyGate) length(pyobj pyObject) (int, error) {
	l := int(C.py_obj_length(pyobj))
	if l < 0 {
		return 0, gate.lastError()
	}
	return l, nil
}

// keys returns keys of the object that support mapping (dict, ...)
// as a sequence object.
func (gate pyGate) keys(pyobj pyObject) (pyObject, error) {
	return gate.objOrLastError(C.py_obj_keys(pyobj))
}

// delattr deletes Object attribute with the specified name.
func (gate pyGate) delattr(pyobj pyObject, name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ok := bool(C.py_obj_delattr(pyobj, cname))
	if !ok {
		return gate.lastError()
	}

	return nil
}

// getattr returns Object attribute with the specified name.
func (gate pyGate) getattr(pyobj pyObject, name string) (pyObject, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var attr pyObject
	if !bool(C.py_obj_getattr(pyobj, cname, &attr)) {
		return nil, gate.lastError()
	}

	return attr, nil
}

// hasattr reports if Object has attribute with the specified name.
func (gate pyGate) hasattr(pyobj pyObject, name string) (bool, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var found C.bool
	if !bool(C.py_obj_hasattr(pyobj, cname, &found)) {
		return false, gate.lastError()
	}

	return bool(found), nil
}

// setattr sets Object attribute with the specified name.
func (gate pyGate) setattr(pyobj pyObject, name string, val pyObject) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if !bool(C.py_obj_setattr(pyobj, cname, val)) {
		return gate.lastError()
	}

	return nil
}

// delitem deletes Object item with the specified key:
//
//	del(pyobj[key])
func (gate pyGate) delitem(pyobj, key pyObject) error {
	if !bool(C.py_obj_delitem(pyobj, key)) {
		return gate.lastError()
	}

	return nil
}

// getitem returns Object item with the specified key.
//
//	pyobj[key]
func (gate pyGate) getitem(pyobj, key pyObject) (pyObject, error) {
	var item pyObject
	if !bool(C.py_obj_getitem(pyobj, key, &item)) {
		return nil, gate.lastError()
	}

	return item, nil
}

// hasitem reports if Object has item with the specified key.
//
//	key in pyobj
func (gate pyGate) hasitem(pyobj, key pyObject) (bool, error) {
	var found C.bool
	if !bool(C.py_obj_hasitem(pyobj, key, &found)) {
		return false, gate.lastError()
	}

	return bool(found), nil
}

// setitem sets Object item with the specified key.
//
//	pyobj[key] = val
func (gate pyGate) setitem(pyobj, key, val pyObject) error {
	if !bool(C.py_obj_setitem(pyobj, key, val)) {
		return gate.lastError()
	}

	return nil
}

// call calls callable object as a function, with positional
// arguments, defined by args (must be PyTuple_Type) and keyword
// arguments, defined by kwargs (must be PyDict_Type or nil).
//
// It returns strong reference to result on success, nil on an error.
func (gate pyGate) call(callable, args, kwargs pyObject) (pyObject, error) {
	return gate.objOrLastError(C.py_obj_call(callable, args, kwargs))
}

// callable reports if object is callable.
// This function always succeeds.
func (gate pyGate) callable(pyobj pyObject) bool {
	return bool(C.py_obj_callable(pyobj))
}

// isBool reports if PyObject is bool.
func (gate pyGate) isBool(pyobj pyObject) bool {
	return bool(C.py_obj_is_bool(pyobj))
}

// isByteArray reports if PyObject is bytearray.
func (gate pyGate) isByteArray(pyobj pyObject) bool {
	return bool(C.py_obj_is_byte_array(pyobj))
}

// isBytes reports if PyObject is bytes.
func (gate pyGate) isBytes(pyobj pyObject) bool {
	return bool(C.py_obj_is_bytes(pyobj))
}

// isComplex reports if PyObject is complex number.
func (gate pyGate) isComplex(pyobj pyObject) bool {
	return bool(C.py_obj_is_complex(pyobj))
}

// isDict reports if PyObject is dict or similar.
func (gate pyGate) isDict(pyobj pyObject) bool {
	return bool(C.py_obj_is_map(pyobj))
}

// isFloat reports if PyObject is floating point number.
func (gate pyGate) isFloat(pyobj pyObject) bool {
	return bool(C.py_obj_is_float(pyobj))
}

// isLong reports if PyObject is long integer.
func (gate pyGate) isLong(pyobj pyObject) bool {
	return bool(C.py_obj_is_long(pyobj))
}

// isSeq reports if PyObject is sequence.
func (gate pyGate) isSeq(pyobj pyObject) bool {
	return bool(C.py_obj_is_seq(pyobj))
}

// isUnicode reports if PyObject is unicode string.
func (gate pyGate) isUnicode(pyobj pyObject) bool {
	return bool(C.py_obj_is_unicode(pyobj))
}

// isUnicode reports if PyObject is type object.
func (gate pyGate) isType(pyobj pyObject) bool {
	return bool(C.py_obj_is_type(pyobj))
}

// decodeError returns [ErrTypeConversion] for Python->Go conversion.
func (gate pyGate) decodeError(pyobj pyObject, to string) error {
	return ErrTypeConversion{
		from: gate.typename(pyobj),
		to:   to,
	}
}

// decodeBigint decodes PyLong_Type object as *big.Int
func (gate pyGate) decodeBigint(pyobj pyObject) (*big.Int, error) {
	if !bool(C.py_obj_is_long(pyobj)) {
		return nil, gate.decodeError(pyobj, "*big.Int,")
	}

	s, err := gate.repr(pyobj)
	if err != nil {
		return nil, err
	}

	v := big.NewInt(0)
	_, ok := v.SetString(s, 10)
	assert.Must(ok) // FIXME

	return v, nil
}

// decodeBytes decodes PyBytes_Type or PyByteArray_Type object as []byte slice.
func (gate pyGate) decodeBytes(pyobj pyObject) ([]byte, error) {
	if bool(C.py_obj_is_bytes(pyobj)) {
		// PyBytes_Type are immutable, so we return a copy.
		var data unsafe.Pointer
		var size C.size_t

		ok := bool(C.py_bytes_get(pyobj, &data, &size))
		if !ok {
			return nil, gate.lastError()
		}

		bytes := make([]byte, size)
		src := unsafe.Slice((*byte)(data), size)
		copy(bytes, src)

		return bytes, nil
	}

	if bool(C.py_obj_is_byte_array(pyobj)) {
		// Python byte array is mutable object, so we return a slice,
		// backed by the Python memory.
		var data unsafe.Pointer
		var size C.size_t

		ok := bool(C.py_bytearray_get(pyobj, &data, &size))
		if !ok {
			return nil, gate.lastError()
		}

		bytes := unsafe.Slice((*byte)(data), size)
		return bytes, nil
	}

	return nil, gate.decodeError(pyobj, "[]byte")
}

// decodeComplex decodes pyObject as complex128.
// It handles conversion from the compatible types.
func (gate pyGate) decodeComplex(pyobj pyObject) (complex128, error) {
	if bool(C.py_obj_is_complex(pyobj)) {
		return gate.decodeExactComplex(pyobj)
	}

	v, err := gate.decodeFloat(pyobj)
	if err == nil {
		return complex(float64(v), 0), nil
	}

	return 0, err
}

// decodeFloat decodes pyObject as float64.
// It handles conversion from the compatible types.
func (gate pyGate) decodeFloat(pyobj pyObject) (float64, error) {
	if bool(C.py_obj_is_float(pyobj)) {
		return gate.decodeExactFloat(pyobj)
	}

	i64, err := gate.decodeExactInt64(pyobj)
	if err == nil {
		return float64(i64), nil
	}

	if _, ovf := err.(ErrOverflow); ovf {
		u64, err2 := gate.decodeExactUint64(pyobj)
		if err2 == nil {
			return float64(u64), nil
		}
	}

	return 0, err
}

// decodeInt64 decodes pyObject as int64.
// It handles conversion from the compatible types.
func (gate pyGate) decodeInt64(pyobj pyObject) (int64, error) {
	if bool(C.py_obj_is_float(pyobj)) {
		f64, err := gate.decodeExactFloat(pyobj)
		if err == nil {
			if minInt64Float <= f64 && f64 <= maxInt64Float {
				return int64(f64), nil
			}

			s, _ := gate.repr(pyobj)
			return 0, ErrOverflow{s}
		}
	}

	return gate.decodeExactInt64(pyobj)
}

// decodeUint64 decodes pyObject as uint64.
// It handles conversion from the compatible types.
func (gate pyGate) decodeUint64(pyobj pyObject) (uint64, error) {
	if bool(C.py_obj_is_float(pyobj)) {
		f64, err := gate.decodeExactFloat(pyobj)
		if err == nil {
			if 0 <= f64 && f64 <= maxUint64Float {
				return uint64(f64), nil
			}

			s, _ := gate.repr(pyobj)
			return 0, ErrOverflow{s}
		}
	}

	return gate.decodeExactUint64(pyobj)
}

// decodeExactComplex decodes pyObject as complex128.
// The object type must be PyComplex_Type.
func (gate pyGate) decodeExactComplex(pyobj pyObject) (complex128, error) {
	if !bool(C.py_obj_is_complex(pyobj)) {
		return 0, gate.decodeError(pyobj, "complex128")
	}

	var real, imag C.double

	ok := bool(C.py_complex_get(pyobj, &real, &imag))
	if !ok {
		return 0, gate.lastError()
	}

	return complex(float64(real), float64(imag)), nil
}

// decodeExactFloat decodes pyObject as float64.
// The object type must be PyFloat_Type.
func (gate pyGate) decodeExactFloat(pyobj pyObject) (float64, error) {
	if !bool(C.py_obj_is_float(pyobj)) {
		return 0, gate.decodeError(pyobj, "float64")
	}

	var val C.double

	ok := bool(C.py_float_get(pyobj, &val))
	if !ok {
		return 0, gate.lastError()
	}

	return float64(val), nil
}

// decodeExactInt64 decodes pyObject as int64.
// The object type must be PyLong_Type.
func (gate pyGate) decodeExactInt64(pyobj pyObject) (int64, error) {
	if !bool(C.py_obj_is_long(pyobj)) {
		return 0, gate.decodeError(pyobj, "int64")
	}

	var val C.int64_t
	var ovf C.bool

	ok := bool(C.py_long_get_int64(pyobj, &val, &ovf))
	switch {
	case !ok:
		return 0, gate.lastError()
	case bool(ovf):
		s, _ := gate.repr(pyobj)
		return 0, ErrOverflow{s}
	}

	return int64(val), nil
}

// decodeExactUint64 decodes pyObject as uint64.
// The object type must be PyLong_Type.
func (gate pyGate) decodeExactUint64(pyobj pyObject) (uint64, error) {
	if !bool(C.py_obj_is_long(pyobj)) {
		return 0, gate.decodeError(pyobj, "uint64")
	}

	var val C.uint64_t
	var ovf C.bool

	ok := bool(C.py_long_get_uint64(pyobj, &val, &ovf))
	switch {
	case !ok:
		return 0, gate.lastError()
	case bool(ovf):
		s, _ := gate.repr(pyobj)
		return 0, ErrOverflow{s}
	}

	return uint64(val), nil
}

// decodeUnicode decodes PyUnicode_Type object as a string.
func (gate pyGate) decodeUnicode(pyobj pyObject) (string, error) {
	if !bool(C.py_obj_is_unicode(pyobj)) {
		return "", gate.decodeError(pyobj, "string")
	}

	sz := C.py_str_len(pyobj)
	assert.Must(sz >= 0) // It only happens if PyObject not Unicode

	s := ""
	if sz > 0 {
		buf := make([]rune, sz)
		p := (*C.Py_UCS4)(unsafe.Pointer(&buf[0]))
		C.py_str_get(pyobj, p, C.size_t(sz))
		s = string(buf)
	}

	return s, nil
}

// makeBigint makes a new PyLong_type object from *big.Int.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeBigint(v *big.Int) (pyObject, error) {
	cs := C.CString(v.String())
	defer C.free(unsafe.Pointer(cs))

	return gate.objOrLastError(C.py_long_from_string(cs))
}

// makeBool makes a new PyBool_Type object.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeBool(v bool) (pyObject, error) {
	return gate.objOrLastError(C.py_bool_make(C.bool(v)))
}

// makeBytes makes a new PyList_Bytes object.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeBytes(data []byte) (pyObject, error) {
	var p unsafe.Pointer
	if len(data) != 0 {
		p = unsafe.Pointer(&data[0])
	}

	return gate.objOrLastError(C.py_bytes_make(p, C.size_t(len(data))))
}

// makeComplex makes a new PyComlex_Type object.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeComplex(v complex128) (pyObject, error) {
	pyobj := C.py_complex_make(C.double(real(v)), C.double(imag(v)))
	return gate.objOrLastError(pyobj)
}

// makeDict makes a new PyDict_Type object.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeDict() (pyObject, error) {
	return gate.objOrLastError(C.py_dict_make())
}

// makeFloat makes a new PyFloatType object.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeFloat(v float64) (pyObject, error) {
	return gate.objOrLastError(C.py_float_make(C.double(v)))
}

// makeInt makes a new PyLong_type object from int64.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeInt(v int64) (pyObject, error) {
	return gate.objOrLastError(C.py_long_from_int64(C.int64_t(v)))
}

// makeList makes a new PyList_Type object of the given size.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeList(sz int) (pyObject, error) {
	return gate.objOrLastError(C.py_list_make(C.size_t(sz)))
}

// makeString makes a new PyUnicoder_type object from string.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeString(v string) (pyObject, error) {
	cs := C.CString(v)
	defer C.free(unsafe.Pointer(cs))
	return gate.objOrLastError(C.py_str_make(cs, C.size_t(len(v))))
}

// makeTuple makes a new PyTuple_Type object of the given size.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeTuple(sz int) (pyObject, error) {
	return gate.objOrLastError(C.py_tuple_make(C.size_t(sz)))
}

// makeUint makes a new PyLong_type object from uint64.
// It returns strong object reference on success, nil on an error.
func (gate pyGate) makeUint(v uint64) (pyObject, error) {
	return gate.objOrLastError(C.py_long_from_uint64(C.uint64_t(v)))
}

// getListItem retrieves the item of the PyList_Type the specified position.
func (gate pyGate) getListItem(list pyObject, idx int) (pyObject, error) {
	item := C.py_list_get(list, C.int(idx))
	return gate.objOrLastError(item)
}

// setListItem sets the item of the PyList_Type the specified position.
// Internally, it creates a new strong reference to the item object.
func (gate pyGate) setListItem(list, item pyObject, idx int) error {
	if !bool(C.py_list_set(list, C.int(idx), item)) {
		return gate.lastError()
	}
	return nil
}

// getSeqItem retrieves the item of the sequence the specified position.
func (gate pyGate) getSeqItem(tuple pyObject, idx int) (pyObject, error) {
	item := C.py_seq_get(tuple, C.int(idx))
	return gate.objOrLastError(item)
}

// getTupleItem retrieves the item of the PyTuple_Type the specified position.
func (gate pyGate) getTupleItem(tuple pyObject, idx int) (pyObject, error) {
	item := C.py_tuple_get(tuple, C.int(idx))
	return gate.objOrLastError(item)
}

// setTupleItem sets the item of the PyTuple_Type the specified position.
// Internally, it creates a new strong reference to the item object.
func (gate pyGate) setTupleItem(tuple, item pyObject, idx int) error {
	if !bool(C.py_tuple_set(tuple, C.int(idx), item)) {
		return gate.lastError()
	}
	return nil
}

// eval evaluates string as a Python statement.
//
// The name parameter indicates the Python source file name and
// used only for diagnostic messages.
//
// If expr is true, input interpreted as a Python expression and
// on success, its result returned as a *Object. Otherwise, input
// interpreted as a multi-line Python script, and returned *Object
// will be nil.
func (gate pyGate) eval(s, name string, expr bool) (pyObject, error) {
	// Convert expression and filename to the C strings
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	// Execute the expression
	var pyobj pyObject
	var lineno C.long
	ok := bool(C.py_interp_eval(cs, cname, C.bool(expr), &pyobj, &lineno))
	if !ok {
		return nil, gate.lastErrorAt(name, int(lineno))
	}

	return pyobj, nil
}

// load loads (imports) string as a Python module.
//
// Module name is specified by the 'name' parameter and
// the 'file' parameter is used as a file name for the
// diagnostics messages.
func (gate pyGate) load(s, name, file string) (pyObject, error) {
	// Convert strings to C
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	// Import the module
	var pyobj pyObject
	ok := bool(C.py_interp_load(cs, cname, cfile, &pyobj))
	if !ok {
		return nil, gate.lastError()
	}

	return pyobj, nil
}

// py_capsule_make makes a new PyCapsule_Type object.
// This object encapsulates an opaque C pointer with optional destructor.
func (gate pyGate) makeCapsule(p unsafe.Pointer, name *C.char,
	destructor C.PyCapsule_Destructor) (pyObject, error) {

	capsule := C.py_capsule_make(p, name, destructor)
	if capsule == nil {
		return nil, gate.lastError()
	}

	return capsule, nil
}

// makeCfunction makes a new PyCFunction_Type object.
// This object represents a Python->C callback.
func (gate pyGate) makeCfunction(def *C.PyMethodDef,
	self pyObject) (pyObject, error) {

	cfunction := C.py_cfunction_make(def, self)
	if cfunction == nil {
		err := gate.lastError()
		return nil, err
	}

	return cfunction, nil
}
