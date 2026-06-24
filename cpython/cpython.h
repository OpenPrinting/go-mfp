// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// CPython glue -- the C side

#ifndef cpython_h
#define cpython_h

#define Py_LIMITED_API  0x03090000
#define PY_SSIZE_T_CLEAN

#include <Python.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

// py_init initializes Python stuff.
// libpython3 must be a full path to the libpython3.XX.so library.
//
// It returns NULL on success or an error message in a case of errors.
// This function needs to be called only once.
//
// This function MUST be called by the main Python thread only.
const char *py_init (const char *libpython3);

// py_new_interp returns a new Python interpreter.
//
// This function MUST be called by the main Python thread only.
void py_new_interp (PyThreadState **tstate_p, PyInterpreterState **interp_p);

// py_interp_close closes the Python interpreter.
//
// This function MUST be called by the main Python thread only.
void py_interp_close (PyThreadState *tstate, PyInterpreterState *interp);

// py_enter temporary attaches the calling thread to the
// Python interpreter.
//
// It must be called before any operations with the interpreter
// are performed and must be paired with the py_leave.
void py_enter (PyInterpreterState *interp);

// py_leave detaches the calling thread from the Python interpreter.
void py_leave (void);

// py_interp_eval evaluates string as a Python statement or expression.
// It returns, via the 'res' pointer, the strong reference to the Python
// value of the executed statement.
//
// The name parameter is used for diagnostics messages and
// indicated the input file name.
//
// If expr is true, this function evaluates Python expression and
// saves its result into res. Otherwise, it evaluates a multi-line
// Python script and don't return any PyObject (sets *res to NULL).
//
// In a case of the execution exception, the file line that caused
// the exception is saved into lineno. If line cannot be determined,
// it will be set to -1.
bool py_interp_eval (const char *s, const char *file,
                     bool expr, PyObject **res, long *lineno);

// py_interp_load loads (imports) string as a Python module.
// It returns, via the 'res' pointer, the strong reference
// to the Python object of the loaded module.
//
// The name parameter becomes the module name, while the
// file parameter used for diagnostics messages and
// indicated the input file name.
bool py_interp_load (const char *s, const char *name, const char *file,
                     PyObject **res);

// py_obj_is_bool reports if PyObject is PyBool_Type
bool py_obj_is_bool (PyObject *x);

// py_obj_is_byte_array reports if PyObject is PyByteArray_Type or its subclass.
bool py_obj_is_byte_array (PyObject *x);

// py_obj_is_bytes reports if PyObject is PyBytes_Type or its subclass.
bool py_obj_is_bytes (PyObject *x);

// py_obj_is_complex reports if PyObject is PyComplex_Type or its subclass.
bool py_obj_is_complex (PyObject *x);

// py_obj_is_float reports if PyObject is PyFloat_Type or its subclass.
bool py_obj_is_float (PyObject *x);

// py_obj_is_long reports if PyObject is PyLong_Type or its subclass.
bool py_obj_is_long (PyObject *x);

// py_obj_is_map reports if PyObject is a map (dict, namedtyple, ...).
bool py_obj_is_map (PyObject *x);

// py_obj_is_seq reports if PyObject is a sequence (list, tuple, ...).
bool py_obj_is_seq (PyObject *x);

// py_obj_is_unicode reports if PyObject is PyUnicode_Type or its subclass.
bool py_obj_is_unicode (PyObject *x);

// py_obj_ref increments the PyObject's reference count.
void py_obj_ref (PyObject *x);

// py_obj_unref decrements the PyObject's reference count.
void py_obj_unref (PyObject *x);

// py_obj_type returns PyObject's type.
PyTypeObject *py_obj_type (PyObject *x);

// py_obj_str returns a string representation of the PyObject.
// This is the equivalent of the Python expression str(x).
//
// There is very subtle difference between py_obj_str and py_obj_repr.
// In general:
//   - Use py_obj_str if you want to print the string
//   - Use py_obj_repr if you want to process the string
PyObject *py_obj_str (PyObject *x);

// py_obj_repr returns a string representation of the PyObject.
// This is the equivalent of the Python expression repr(x).
//
// There is very subtle difference between py_obj_str and py_obj_repr.
// In general:
//   - Use py_obj_str if you want to print the string
//   - Use py_obj_repr if you want to process the string
PyObject *py_obj_repr (PyObject *x);

// py_obj_length returns PyObject, in items. It works with any
// container objects (lists, tuples, dictionaries, ...).
// Returns -1 on a error.
ssize_t py_obj_length (PyObject *x);

// py_obj_keys returns PyObject mapping keys. It works for objects
// that supports mapping (see py_obj_is_map), i.e., dict etc.
//
// On success it returns PyList_Type or PyTuple_Type object that
// contains the keys. On error it returns NULL.
PyObject *py_obj_keys (PyObject *x);

// py_obj_hasattr reports if PyObject has the attribute with the
// specified name.
//
// It returns true on success, false on error and puts answer into
// its third parameter.
bool py_obj_hasattr(PyObject *x, const char *name, bool *answer);

// py_obj_delattr deletes the attribute with the specified name.
// It returns true on success, false on error.
bool py_obj_delattr(PyObject *x, const char *name);

// py_obj_getattr retrieves the attribute with the specified name.
// The returned answer, on success, contains a string reference to PyObject.
// It returns true on success, false on error.
bool py_obj_getattr(PyObject *x, const char *name, PyObject **answer);

// py_obj_setattr sets the attribute with the specified name.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_obj_setattr(PyObject *x, const char *name, PyObject *value);

// py_obj_hasitem reports if PyObject contains the item with the
// specified key.
//
// It returns true on success, false on error and puts answer into
// its third parameter.
bool py_obj_hasitem(PyObject *x, PyObject *key, bool *answer);

// py_obj_delitem deletes the item with the specified key.
// It returns true on success, false on error.
bool py_obj_delitem(PyObject *x, PyObject *key);

// py_obj_getitem retrieves the item with the specified key.
// The returned answer, on success, contains a strong reference to PyObject.
// It returns true on success, false on error.
bool py_obj_getitem(PyObject *x, PyObject *key, PyObject **answer);

// py_obj_setitem sets the item with the specified key.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_obj_setitem(PyObject *x, PyObject *key, PyObject *value);

// py_obj_call calls callable object (i.e., function, method, ...)
// with the specified arguments.
//
// The args parameter must be of PyTuple_Type object and it specified
// the function parameters. It must not be NULL.
//
// The kwargs must be PyDict_Type and it specifies keyword arguments.
// It can be NULL, if keyword arguments are not used.
//
// It returns strong object reference on success, NULL on an error.
PyObject *py_obj_call(PyObject *x, PyObject *args, PyObject *kwargs);

// py_obj_callable reports if object is callable.
// This function always succeeds.
bool py_obj_callable(PyObject *x);

// py_err_fetch fetches and clears last error.
// If there is no pending error, all pointers will be set to NULL.
void py_err_fetch (PyObject **etype, PyObject **evalue, PyObject **trace);

// py_bool_make makes a new PyBool_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_bool_make(bool val);

// py_bytes_get obtains content of the Python bytes object.
// It returns true on success, false on error.
bool py_bytes_get (PyObject *x, void **data, size_t *size);

// py_bytes_make makes a new PyBytes_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_bytes_make(const void *data, size_t size);

// py_bytearray_get obtains content of the Python bytearray object.
// It returns true on success, false on error.
bool py_bytearray_get (PyObject *x, void **data, size_t *size);

// py_complex_get obtains content of the Python complex object.
// It returns true on success, false on error.
bool py_complex_get (PyObject *x, double *real, double *imag);

// py_complex_make makes a new PyComlex_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_complex_make(double real, double imag);

// py_dict_make makes a new PyDict_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_dict_make(void);

// py_float_get obtains content of the Python float object.
// It returns true on success, false on error.
bool py_float_get (PyObject *x, double *val);

// py_float_make makes a new PyFloat_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_float_make(double val);

// py_list_make makes a new PyList_Type object of the specified size.
// The newly created list MUST be fully populated with the py_list_set
// calls before it can be safely passed to Python interpreter.
// It returns strong object reference on success, NULL on an error.
PyObject *py_list_make(size_t len);

// py_list_set retrieves value of the list item at the given position.
// It returns strong object reference on success, NULL on an error.
PyObject *py_list_get(PyObject *list, int index);

// py_list_set sets value of the list item at the given position.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_list_set(PyObject *list, int index, PyObject *val);

// py_long_get obtains PyObject's value as int64_t.
// If value doesn't fit C long, overflow flag is set.
//
// It returns true on success, false on error.
bool py_long_get_int64 (PyObject *x, int64_t *val, bool *overflow);

// py_long_get obtains PyObject's value as uint64_t.
// If value doesn't fit C long, overflow flag is set.
//
// It returns true on success, false on error.
bool py_long_get_uint64 (PyObject *x, uint64_t *val, bool *overflow);

// py_long_from_int64 makes a new PyLong_Type object from int64_t value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_long_from_int64(int64_t val);
//
// py_long_from_uint64 makes a new PyLong_Type object from uint64_t value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_long_from_uint64(uint64_t val);

// py_long_from_string makes a new PyLong_Type object from string value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_long_from_string(const char *val);

// py_seq_set retrieves value of the sequence item at the given position.
// It returns strong object reference on success, NULL on an error.
PyObject *py_seq_get(PyObject *tuple, int index);

// py_str_len returns length of Unicode string, in code points.
// If PyObject is not Unicode, it returns -1.
static inline ssize_t py_str_len (PyObject *str) {
    extern Py_ssize_t (*PyUnicode_GetLength_p)(PyObject *);
    return (ssize_t) PyUnicode_GetLength_p(str);
}

// py_str_get copies Unicode string data as a sequence of the Py_UCS4
// characters.
//
// On success it returns buf, otherwise returns NULL.
// The function may fail if PyObject is not Unicode of if buffer
// is too short.
//
// The trailing '\0' is not copied.
//
// Use py_str_len to obtain the correct string length.
Py_UCS4 *py_str_get (PyObject *str, Py_UCS4 *buf, size_t len);

// py_str_make makes a new PyLong_Type object from string value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_str_make(const char *val, size_t len);

// py_tuple_make makes a new PyTuple_Type object of the specified size.
// The newly created tuple MUST be fully populated with the py_tuple_set
// calls before it can be safely passed to Python interpreter.
// It returns strong object reference on success, NULL on an error.
PyObject *py_tuple_make(size_t len);

// py_tuple_set retrieves value of the tuple item at the given position.
// It returns strong object reference on success, NULL on an error.
PyObject *py_tuple_get(PyObject *tuple, int index);

// py_tuple_set sets value of the tuple item at the given position.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_tuple_set(PyObject *tuple, int index, PyObject *val);

// py_cfunction_make makes a new PyCFunction object.
// The caller must ensure that ml outlives the returned object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_cfunction_make(PyMethodDef *ml, PyObject *self);

// py_cfunction_make makes a new PyCFunction object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_capsule_make(void *pointer, const char *name,
                          PyCapsule_Destructor destructor);

// py_capsule_get_ptr returns pointer, encapsulated into
// the capsule.
//
// The name must match exactly the name used for py_capsule_make.
void *py_capsule_get_ptr(PyObject *capsule, const char *name);

// py_set_err sets the error condition within the Python interpreter,
// which raises exception at the Python side.
//
// The callback function, registered by the py_cfunction_make, MUST
// set the error, if it returns NULL.
void py_set_err(PyObject *extype, const char *msg);

// Python build-in (primitive) types:
extern PyTypeObject *PyBool_Type_p;
extern PyTypeObject *PyByteArray_Type_p;
extern PyTypeObject *PyBytes_Type_p;
extern PyTypeObject *PyCFunction_Type_p;
extern PyTypeObject *PyComplex_Type_p;
extern PyTypeObject *PyDict_Type_p;
extern PyTypeObject *PyDictKeys_Type_p;
extern PyTypeObject *PyFloat_Type_p;
extern PyTypeObject *PyFrozenSet_Type_p;
extern PyTypeObject *PyList_Type_p;
extern PyTypeObject *PyLong_Type_p;
extern PyTypeObject *PyMemoryView_Type_p;
extern PyTypeObject *PyModule_Type_p;
extern PyTypeObject *PySet_Type_p;
extern PyTypeObject *PySlice_Type_p;
extern PyTypeObject *PyTuple_Type_p;
extern PyTypeObject *PyType_Type_p;
extern PyTypeObject *PyUnicode_Type_p;

// Standard exceptions
extern PyObject *PyExc_ArithmeticError_p;
extern PyObject *PyExc_AssertionError_p;
extern PyObject *PyExc_AttributeError_p;
extern PyObject *PyExc_BaseException_p;
extern PyObject *PyExc_BlockingIOError_p;
extern PyObject *PyExc_BrokenPipeError_p;
extern PyObject *PyExc_BufferError_p;
extern PyObject *PyExc_ChildProcessError_p;
extern PyObject *PyExc_ConnectionAbortedError_p;
extern PyObject *PyExc_ConnectionError_p;
extern PyObject *PyExc_ConnectionRefusedError_p;
extern PyObject *PyExc_ConnectionResetError_p;
extern PyObject *PyExc_EOFError_p;
extern PyObject *PyExc_Exception_p;
extern PyObject *PyExc_FileExistsError_p;
extern PyObject *PyExc_FileNotFoundError_p;
extern PyObject *PyExc_FloatingPointError_p;
extern PyObject *PyExc_GeneratorExit_p;
extern PyObject *PyExc_ImportError_p;
extern PyObject *PyExc_IndentationError_p;
extern PyObject *PyExc_IndexError_p;
extern PyObject *PyExc_InterruptedError_p;
extern PyObject *PyExc_IsADirectoryError_p;
extern PyObject *PyExc_KeyboardInterrupt_p;
extern PyObject *PyExc_KeyError_p;
extern PyObject *PyExc_LookupError_p;
extern PyObject *PyExc_MemoryError_p;
extern PyObject *PyExc_ModuleNotFoundError_p;
extern PyObject *PyExc_NameError_p;
extern PyObject *PyExc_NotADirectoryError_p;
extern PyObject *PyExc_NotImplementedError_p;
extern PyObject *PyExc_OSError_p;
extern PyObject *PyExc_OverflowError_p;
extern PyObject *PyExc_PermissionError_p;
extern PyObject *PyExc_ProcessLookupError_p;
extern PyObject *PyExc_RecursionError_p;
extern PyObject *PyExc_ReferenceError_p;
extern PyObject *PyExc_RuntimeError_p;
extern PyObject *PyExc_StopAsyncIteration_p;
extern PyObject *PyExc_StopIteration_p;
extern PyObject *PyExc_SyntaxError_p;
extern PyObject *PyExc_SystemError_p;
extern PyObject *PyExc_SystemExit_p;
extern PyObject *PyExc_TabError_p;
extern PyObject *PyExc_TimeoutError_p;
extern PyObject *PyExc_TypeError_p;
extern PyObject *PyExc_UnboundLocalError_p;
extern PyObject *PyExc_UnicodeDecodeError_p;
extern PyObject *PyExc_UnicodeEncodeError_p;
extern PyObject *PyExc_UnicodeError_p;
extern PyObject *PyExc_UnicodeTranslateError_p;
extern PyObject *PyExc_ValueError_p;
extern PyObject *PyExc_ZeroDivisionError_p;

// Standard warnings
extern PyObject *PyExc_BytesWarning_p;
extern PyObject *PyExc_DeprecationWarning_p;
extern PyObject *PyExc_FutureWarning_p;
extern PyObject *PyExc_ImportWarning_p;
extern PyObject *PyExc_PendingDeprecationWarning_p;
extern PyObject *PyExc_ResourceWarning_p;
extern PyObject *PyExc_RuntimeWarning_p;
extern PyObject *PyExc_SyntaxWarning_p;
extern PyObject *PyExc_UnicodeWarning_p;
extern PyObject *PyExc_UserWarning_p;
extern PyObject *PyExc_Warning_p;

#endif

// vim:ts=8:sw=4:et
