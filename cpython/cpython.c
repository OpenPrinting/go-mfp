// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// CPython glue -- the C side

#include "cpython.h"

#include <assert.h>
#include <dlfcn.h>
#include <stdarg.h>
#include <stdatomic.h>

// py_libpython3 contains libpython3.so handle.
static void                                     *py_libpython3;

// py_error is set to non-NULL in a case of any initialization
// error. py_error_buf provides static buffer for that.
//
// This all is thread-safe, because used only at the initialization
// time where only a single thread has access to this stuff.
static const char                               *py_error;
static char                                     py_error_buf[1024];

// py_main_thread keeps reference to the main Python thread state
// between calls to functions that use that.
static PyThreadState                            *py_main_thread;

// py_enter_level counts py_enter nesting level when called
// from the same thread;
static __thread int                             py_enter_level;
static __thread PyThreadState                   *py_thread_current;

// The table of the libpython3 symbols.
//
// Most of the programs that use embedded Python interpreter
// link against the libpython3.NN.so.1.0 dynamic library.
//
// As result, they are tied to the particular Python version.
//
// We, instead, link against the libpython3.so and use only the
// Stable Python API (https://docs.python.org/3/c-api/stable.html).
//
// Unfortunately, libpython3.so doesn't expose any symbols and
// only contains a reference (ELF NEEDED) to the proper
// libpython3.NN.so.1.0 library.
//
// So we can't use libpython3.NN.so.1.0 symbols directly. Loading
// the libpython3.so implicitly puts these symbols into the main
// process namespace, but every symbols needs to be manually resolved
// with the explicit dlsym(py_libpython3, name) call.
//
// So these pointers keeps these resolved symbols...
static __typeof__(PyBool_FromLong)              *PyBool_FromLong_p;
static __typeof__(PyByteArray_AsString)         *PyByteArray_AsString_p;
static __typeof__(PyByteArray_Size)             *PyByteArray_Size_p;
static __typeof__(PyBytes_AsStringAndSize)      *PyBytes_AsStringAndSize_p;
static __typeof__(PyBytes_FromStringAndSize)    *PyBytes_FromStringAndSize_p;
static __typeof__(PyCallable_Check)             *PyCallable_Check_p;
static __typeof__(PyCapsule_GetPointer)         *PyCapsule_GetPointer_p;
static __typeof__(PyCapsule_New)                *PyCapsule_New_p;
static __typeof__(PyCFunction_NewEx)            *PyCFunction_NewEx_p;
static __typeof__(Py_CompileString)             *Py_CompileString_p;
static __typeof__(PyComplex_FromDoubles)        *PyComplex_FromDoubles_p;
static __typeof__(PyComplex_ImagAsDouble)       *PyComplex_ImagAsDouble_p;
static __typeof__(PyComplex_RealAsDouble)       *PyComplex_RealAsDouble_p;
static __typeof__(Py_DecRef)                    *Py_DecRef_p;
static __typeof__(PyDict_New)                   *PyDict_New_p;
static __typeof__(PyDict_SetItemString)         *PyDict_SetItemString_p;
static __typeof__(Py_EndInterpreter)            *Py_EndInterpreter_p;
static __typeof__(PyErr_Clear)                  *PyErr_Clear_p;
static __typeof__(PyErr_Fetch)                  *PyErr_Fetch_p;
static __typeof__(PyErr_NormalizeException)     *PyErr_NormalizeException_p;
static __typeof__(PyErr_Occurred)               *PyErr_Occurred_p;
static __typeof__(PyErr_Restore)                *PyErr_Restore_p;
static __typeof__(PyErr_SetString)              *PyErr_SetString_p;
static __typeof__(PyEval_AcquireThread)         *PyEval_AcquireThread_p;
static __typeof__(PyEval_EvalCode)              *PyEval_EvalCode_p;
static __typeof__(PyEval_ReleaseThread)         *PyEval_ReleaseThread_p;
static __typeof__(PyEval_RestoreThread)         *PyEval_RestoreThread_p;
static __typeof__(PyEval_SaveThread)            *PyEval_SaveThread_p;
static __typeof__(PyFloat_AsDouble)             *PyFloat_AsDouble_p;
static __typeof__(PyFloat_FromDouble)           *PyFloat_FromDouble_p;
static __typeof__(PyImport_AddModule)           *PyImport_AddModule_p;
static __typeof__(PyImport_ExecCodeModule)      *PyImport_ExecCodeModule_p;
static __typeof__(Py_IncRef)                    *Py_IncRef_p;
static __typeof__(Py_InitializeEx)              *Py_InitializeEx_p;
static __typeof__(PyList_GetItem)               *PyList_GetItem_p;
static __typeof__(PyList_New)                   *PyList_New_p;
static __typeof__(PyList_SetItem)               *PyList_SetItem_p;
static __typeof__(PyLong_AsLongLong)            *PyLong_AsLongLong_p;
static __typeof__(PyLong_AsUnsignedLongLong)    *PyLong_AsUnsignedLongLong_p;
static __typeof__(PyLong_FromLongLong)          *PyLong_FromLongLong_p;
static __typeof__(PyLong_FromString)            *PyLong_FromString_p;
static __typeof__(PyLong_FromUnsignedLongLong)  *PyLong_FromUnsignedLongLong_p;
static __typeof__(PyMapping_Check)              *PyMapping_Check_p;
static __typeof__(PyMapping_Keys)               *PyMapping_Keys_p;
static __typeof__(PyModule_GetDict)             *PyModule_GetDict_p;
static __typeof__(Py_NewInterpreter)            *Py_NewInterpreter_p;
static __typeof__(PyObject_Call)                *PyObject_Call_p;
static __typeof__(PyObject_DelItem)             *PyObject_DelItem_p;
static __typeof__(PyObject_GetAttrString)       *PyObject_GetAttrString_p;
static __typeof__(PyObject_GetItem)             *PyObject_GetItem_p;
static __typeof__(*PyObject_Length)             *PyObject_Length_p;
static __typeof__(PyObject_Repr)                *PyObject_Repr_p;
static __typeof__(PyObject_SetAttrString)       *PyObject_SetAttrString_p;
static __typeof__(PyObject_SetItem)             *PyObject_SetItem_p;
static __typeof__(PyObject_Str)                 *PyObject_Str_p;
static __typeof__(PyObject_Type)                *PyObject_Type_p;
static __typeof__(PySequence_Check)             *PySequence_Check_p;
static __typeof__(PySequence_GetItem)           *PySequence_GetItem_p;
static __typeof__(PyThreadState_Clear)          *PyThreadState_Clear_p;
static __typeof__(PyThreadState_Delete)         *PyThreadState_Delete_p;
static __typeof__(PyThreadState_New)            *PyThreadState_New_p;
static __typeof__(PyThreadState_Swap)           *PyThreadState_Swap_p;
static __typeof__(PyTuple_GetItem)              *PyTuple_GetItem_p;
static __typeof__(PyTuple_New)                  *PyTuple_New_p;
static __typeof__(PyTuple_SetItem)              *PyTuple_SetItem_p;
static __typeof__(*PyType_GetFlags)             *PyType_GetFlags_p;
static __typeof__(PyType_IsSubtype)             *PyType_IsSubtype_p;
static __typeof__(PyUnicode_AsUCS4)             *PyUnicode_AsUCS4_p;
static __typeof__(PyUnicode_FromStringAndSize)  *PyUnicode_FromStringAndSize_p;

// PyInterpreterState_Get is not officially exposed on Python 3.8,
// so it needs the special care.
static PyInterpreterState *(*PyInterpreterState_Get_p)(void);

// Directly exposed functions
__typeof__(PyUnicode_GetLength)                  *PyUnicode_GetLength_p;

// Python build-in (primitive) types:
PyTypeObject *PyBool_Type_p;
PyTypeObject *PyByteArray_Type_p;
PyTypeObject *PyBytes_Type_p;
PyTypeObject *PyCFunction_Type_p;
PyTypeObject *PyComplex_Type_p;
PyTypeObject *PyDictKeys_Type_p;
PyTypeObject *PyDict_Type_p;
PyTypeObject *PyFloat_Type_p;
PyTypeObject *PyFrozenSet_Type_p;
PyTypeObject *PyList_Type_p;
PyTypeObject *PyLong_Type_p;
PyTypeObject *PyMemoryView_Type_p;
PyTypeObject *PyModule_Type_p;
PyTypeObject *PySet_Type_p;
PyTypeObject *PySlice_Type_p;
PyTypeObject *PyTuple_Type_p;
PyTypeObject *PyType_Type_p;
PyTypeObject *PyUnicode_Type_p;

// Standard exceptions
PyObject *PyExc_ArithmeticError_p;
PyObject *PyExc_AssertionError_p;
PyObject *PyExc_AttributeError_p;
PyObject *PyExc_BaseException_p;
PyObject *PyExc_BlockingIOError_p;
PyObject *PyExc_BrokenPipeError_p;
PyObject *PyExc_BufferError_p;
PyObject *PyExc_ChildProcessError_p;
PyObject *PyExc_ConnectionAbortedError_p;
PyObject *PyExc_ConnectionError_p;
PyObject *PyExc_ConnectionRefusedError_p;
PyObject *PyExc_ConnectionResetError_p;
PyObject *PyExc_EOFError_p;
PyObject *PyExc_Exception_p;
PyObject *PyExc_FileExistsError_p;
PyObject *PyExc_FileNotFoundError_p;
PyObject *PyExc_FloatingPointError_p;
PyObject *PyExc_GeneratorExit_p;
PyObject *PyExc_ImportError_p;
PyObject *PyExc_IndentationError_p;
PyObject *PyExc_IndexError_p;
PyObject *PyExc_InterruptedError_p;
PyObject *PyExc_IsADirectoryError_p;
PyObject *PyExc_KeyboardInterrupt_p;
PyObject *PyExc_KeyError_p;
PyObject *PyExc_LookupError_p;
PyObject *PyExc_MemoryError_p;
PyObject *PyExc_ModuleNotFoundError_p;
PyObject *PyExc_NameError_p;
PyObject *PyExc_NotADirectoryError_p;
PyObject *PyExc_NotImplementedError_p;
PyObject *PyExc_OSError_p;
PyObject *PyExc_OverflowError_p;
PyObject *PyExc_PermissionError_p;
PyObject *PyExc_ProcessLookupError_p;
PyObject *PyExc_RecursionError_p;
PyObject *PyExc_ReferenceError_p;
PyObject *PyExc_RuntimeError_p;
PyObject *PyExc_StopAsyncIteration_p;
PyObject *PyExc_StopIteration_p;
PyObject *PyExc_SyntaxError_p;
PyObject *PyExc_SystemError_p;
PyObject *PyExc_SystemExit_p;
PyObject *PyExc_TabError_p;
PyObject *PyExc_TimeoutError_p;
PyObject *PyExc_TypeError_p;
PyObject *PyExc_UnboundLocalError_p;
PyObject *PyExc_UnicodeDecodeError_p;
PyObject *PyExc_UnicodeEncodeError_p;
PyObject *PyExc_UnicodeError_p;
PyObject *PyExc_UnicodeTranslateError_p;
PyObject *PyExc_ValueError_p;
PyObject *PyExc_ZeroDivisionError_p;

// Standard warnings
PyObject *PyExc_BytesWarning_p;
PyObject *PyExc_DeprecationWarning_p;
PyObject *PyExc_FutureWarning_p;
PyObject *PyExc_ImportWarning_p;
PyObject *PyExc_PendingDeprecationWarning_p;
PyObject *PyExc_ResourceWarning_p;
PyObject *PyExc_RuntimeWarning_p;
PyObject *PyExc_SyntaxWarning_p;
PyObject *PyExc_UnicodeWarning_p;
PyObject *PyExc_UserWarning_p;
PyObject *PyExc_Warning_p;

// Forward declarations
static PyObject *py_obj_getattr_int(PyObject *x, const char *name);

// py_set_error formats and sets py_error.
static void py_set_error (const char *fmt, ...) {
    va_list ap;

    if (py_error == NULL) {
        va_start(ap, fmt);
        vsnprintf(py_error_buf, sizeof(py_error_buf), fmt, ap);
        va_end(ap);
        py_error = py_error_buf;
    }
}

// py_load loads Python symbol by name.
static void *py_load (const char *name) {
    void *p = NULL;

    if (py_error == NULL) {
        p = dlsym(py_libpython3, name);
        if (p == NULL) {
            py_set_error("%s", dlerror());
        }
    }

    return p;
}

// py_load2 loads Python symbol by name, trying two variants
// of the name
static void *py_load2 (const char *name1, const char *name2) {
    void *p = NULL;

    if (py_error == NULL) {
        p = dlsym(RTLD_DEFAULT, name1);
        if (p == NULL) {
            p = dlsym(RTLD_DEFAULT, name2);
        }

        if (p == NULL) {
            py_set_error("%s", dlerror());
        }
    }

    return p;
}

// py_load loads and dereferences pointer from the libpython3.so.
static void *py_load_ptr (const char *name) {
    void **pp = py_load(name);

    if (pp != NULL) {
        return *pp;
    }

    return NULL;
}

// py_load_library loads the libpython3.so library
static void py_load_library (const char *libpython3) {
    py_libpython3 = dlopen(libpython3, RTLD_NOW | RTLD_GLOBAL);
    if (py_libpython3 == NULL) {
        py_set_error("%s", dlerror());
    }
}

// py_load_all loads all Python symbols.
static void py_load_all (const char *libpython3) {
    // Load the libpython3.so library
    py_load_library(libpython3);

    // Load all required symbols
    PyBool_FromLong_p = py_load("PyBool_FromLong");
    PyByteArray_AsString_p = py_load("PyByteArray_AsString");
    PyByteArray_Size_p = py_load("PyByteArray_Size");
    PyBytes_AsStringAndSize_p = py_load("PyBytes_AsStringAndSize");
    PyBytes_FromStringAndSize_p = py_load("PyBytes_FromStringAndSize");
    PyCallable_Check_p = py_load("PyCallable_Check");
    PyCapsule_GetPointer_p = py_load("PyCapsule_GetPointer");
    PyCapsule_New_p = py_load("PyCapsule_New");
    PyCFunction_NewEx_p = py_load("PyCFunction_NewEx");
    Py_CompileString_p = py_load("Py_CompileString");
    PyComplex_FromDoubles_p = py_load("PyComplex_FromDoubles");
    PyComplex_ImagAsDouble_p = py_load("PyComplex_ImagAsDouble");
    PyComplex_RealAsDouble_p = py_load("PyComplex_RealAsDouble");
    Py_DecRef_p = py_load("Py_DecRef");
    PyDict_New_p = py_load("PyDict_New");
    PyDict_SetItemString_p = py_load("PyDict_SetItemString");
    Py_EndInterpreter_p = py_load("Py_EndInterpreter");
    PyErr_Clear_p = py_load("PyErr_Clear");
    PyErr_Fetch_p = py_load("PyErr_Fetch");
    PyErr_NormalizeException_p = py_load("PyErr_NormalizeException");
    PyErr_Occurred_p = py_load("PyErr_Occurred");
    PyErr_Restore_p = py_load("PyErr_Restore");
    PyErr_SetString_p = py_load("PyErr_SetString");
    PyEval_AcquireThread_p = py_load("PyEval_AcquireThread");
    PyEval_EvalCode_p = py_load("PyEval_EvalCode");
    PyEval_ReleaseThread_p = py_load("PyEval_ReleaseThread");
    PyEval_RestoreThread_p = py_load("PyEval_RestoreThread");
    PyEval_SaveThread_p = py_load("PyEval_SaveThread");
    PyFloat_AsDouble_p = py_load("PyFloat_AsDouble");
    PyFloat_FromDouble_p = py_load("PyFloat_FromDouble");
    PyImport_AddModule_p = py_load("PyImport_AddModule");
    PyImport_ExecCodeModule_p = py_load("PyImport_ExecCodeModule");
    Py_IncRef_p = py_load("Py_IncRef");
    Py_InitializeEx_p = py_load("Py_InitializeEx");
    PyList_GetItem_p = py_load("PyList_GetItem");
    PyList_New_p = py_load("PyList_New");
    PyList_SetItem_p = py_load("PyList_SetItem");
    PyLong_AsLongLong_p = py_load("PyLong_AsLongLong");
    PyLong_AsUnsignedLongLong_p = py_load("PyLong_AsUnsignedLongLong");
    PyLong_FromLongLong_p = py_load("PyLong_FromLongLong");
    PyLong_FromString_p = py_load("PyLong_FromString");
    PyLong_FromUnsignedLongLong_p = py_load("PyLong_FromUnsignedLongLong");
    PyMapping_Check_p = py_load("PyMapping_Check");
    PyMapping_Keys_p = py_load("PyMapping_Keys");
    PyModule_GetDict_p = py_load("PyModule_GetDict");
    Py_NewInterpreter_p = py_load("Py_NewInterpreter");
    PyObject_Call_p = py_load("PyObject_Call");
    PyObject_DelItem_p = py_load("PyObject_DelItem");
    PyObject_GetAttrString_p = py_load("PyObject_GetAttrString");
    PyObject_GetItem_p = py_load("PyObject_GetItem");
    PyObject_Length_p = py_load("PyObject_Length");
    PyObject_Repr_p = py_load("PyObject_Repr");
    PyObject_SetAttrString_p = py_load("PyObject_SetAttrString");
    PyObject_SetItem_p = py_load("PyObject_SetItem");
    PyObject_Str_p = py_load("PyObject_Str");
    PyObject_Type_p = py_load("PyObject_Type");
    PySequence_Check_p = py_load("PySequence_Check");
    PySequence_GetItem_p = py_load("PySequence_GetItem");
    PyThreadState_Clear_p = py_load("PyThreadState_Clear");
    PyThreadState_Delete_p = py_load("PyThreadState_Delete");
    PyThreadState_New_p = py_load("PyThreadState_New");
    PyThreadState_Swap_p = py_load("PyThreadState_Swap");
    PyTuple_GetItem_p = py_load("PyTuple_GetItem");
    PyTuple_New_p = py_load("PyTuple_New");
    PyTuple_SetItem_p = py_load("PyTuple_SetItem");
    PyType_GetFlags_p = py_load("PyType_GetFlags");
    PyType_IsSubtype_p = py_load("PyType_IsSubtype");
    PyUnicode_AsUCS4_p = py_load("PyUnicode_AsUCS4");
    PyUnicode_FromStringAndSize_p = py_load("PyUnicode_FromStringAndSize");

    // Python 3.8 "unoficcially" exposes this symbol as
    // _PyInterpreterState_Get.
    //
    // Python 3.9 and up exposes it as PyInterpreterState_Get, and
    // now it's a part of the stable API.
    //
    // So this trickery is safe.
    PyInterpreterState_Get_p = py_load2(
        "_PyInterpreterState_Get", "PyInterpreterState_Get");

    PyUnicode_GetLength_p = py_load("PyUnicode_GetLength");

    PyBool_Type_p = py_load("PyBool_Type");
    PyByteArray_Type_p = py_load("PyByteArray_Type");
    PyBytes_Type_p = py_load("PyBytes_Type");
    PyCFunction_Type_p = py_load("PyCFunction_Type");
    PyComplex_Type_p = py_load("PyComplex_Type");
    PyDictKeys_Type_p = py_load("PyDictKeys_Type");
    PyDict_Type_p = py_load("PyDict_Type");
    PyFloat_Type_p = py_load("PyFloat_Type");
    PyFrozenSet_Type_p = py_load("PyFrozenSet_Type");
    PyList_Type_p = py_load("PyList_Type");
    PyLong_Type_p = py_load("PyLong_Type");
    PyMemoryView_Type_p = py_load("PyMemoryView_Type");
    PyModule_Type_p = py_load("PyModule_Type");
    PySet_Type_p = py_load("PySet_Type");
    PySlice_Type_p = py_load("PySlice_Type");
    PyTuple_Type_p = py_load("PyTuple_Type");
    PyType_Type_p = py_load("PyType_Type");
    PyUnicode_Type_p = py_load("PyUnicode_Type");

    PyExc_ArithmeticError_p = py_load_ptr("PyExc_ArithmeticError");
    PyExc_AssertionError_p = py_load_ptr("PyExc_AssertionError");
    PyExc_AttributeError_p = py_load_ptr("PyExc_AttributeError");
    PyExc_BaseException_p = py_load_ptr("PyExc_BaseException");
    PyExc_BlockingIOError_p = py_load_ptr("PyExc_BlockingIOError");
    PyExc_BrokenPipeError_p = py_load_ptr("PyExc_BrokenPipeError");
    PyExc_BufferError_p = py_load_ptr("PyExc_BufferError");
    PyExc_ChildProcessError_p = py_load_ptr("PyExc_ChildProcessError");
    PyExc_ConnectionAbortedError_p = py_load_ptr("PyExc_ConnectionAbortedError");
    PyExc_ConnectionError_p = py_load_ptr("PyExc_ConnectionError");
    PyExc_ConnectionRefusedError_p = py_load_ptr("PyExc_ConnectionRefusedError");
    PyExc_ConnectionResetError_p = py_load_ptr("PyExc_ConnectionResetError");
    PyExc_EOFError_p = py_load_ptr("PyExc_EOFError");
    PyExc_Exception_p = py_load_ptr("PyExc_Exception");
    PyExc_FileExistsError_p = py_load_ptr("PyExc_FileExistsError");
    PyExc_FileNotFoundError_p = py_load_ptr("PyExc_FileNotFoundError");
    PyExc_FloatingPointError_p = py_load_ptr("PyExc_FloatingPointError");
    PyExc_GeneratorExit_p = py_load_ptr("PyExc_GeneratorExit");
    PyExc_ImportError_p = py_load_ptr("PyExc_ImportError");
    PyExc_IndentationError_p = py_load_ptr("PyExc_IndentationError");
    PyExc_IndexError_p = py_load_ptr("PyExc_IndexError");
    PyExc_InterruptedError_p = py_load_ptr("PyExc_InterruptedError");
    PyExc_IsADirectoryError_p = py_load_ptr("PyExc_IsADirectoryError");
    PyExc_KeyboardInterrupt_p = py_load_ptr("PyExc_KeyboardInterrupt");
    PyExc_KeyError_p = py_load_ptr("PyExc_KeyError");
    PyExc_LookupError_p = py_load_ptr("PyExc_LookupError");
    PyExc_MemoryError_p = py_load_ptr("PyExc_MemoryError");
    PyExc_ModuleNotFoundError_p = py_load_ptr("PyExc_ModuleNotFoundError");
    PyExc_NameError_p = py_load_ptr("PyExc_NameError");
    PyExc_NotADirectoryError_p = py_load_ptr("PyExc_NotADirectoryError");
    PyExc_NotImplementedError_p = py_load_ptr("PyExc_NotImplementedError");
    PyExc_OSError_p = py_load_ptr("PyExc_OSError");
    PyExc_OverflowError_p = py_load_ptr("PyExc_OverflowError");
    PyExc_PermissionError_p = py_load_ptr("PyExc_PermissionError");
    PyExc_ProcessLookupError_p = py_load_ptr("PyExc_ProcessLookupError");
    PyExc_RecursionError_p = py_load_ptr("PyExc_RecursionError");
    PyExc_ReferenceError_p = py_load_ptr("PyExc_ReferenceError");
    PyExc_RuntimeError_p = py_load_ptr("PyExc_RuntimeError");
    PyExc_StopAsyncIteration_p = py_load_ptr("PyExc_RuntimeError");
    PyExc_StopIteration_p = py_load_ptr("PyExc_StopIteration");
    PyExc_SyntaxError_p = py_load_ptr("PyExc_SyntaxError");
    PyExc_SystemError_p = py_load_ptr("PyExc_SystemError");
    PyExc_SystemExit_p = py_load_ptr("PyExc_SystemExit");
    PyExc_TabError_p = py_load_ptr("PyExc_TabError");
    PyExc_TimeoutError_p = py_load_ptr("PyExc_TimeoutError");
    PyExc_TypeError_p = py_load_ptr("PyExc_TypeError");
    PyExc_UnboundLocalError_p = py_load_ptr("PyExc_UnboundLocalError");
    PyExc_UnicodeDecodeError_p = py_load_ptr("PyExc_UnicodeDecodeError");
    PyExc_UnicodeEncodeError_p = py_load_ptr("PyExc_UnicodeEncodeError");
    PyExc_UnicodeError_p = py_load_ptr("PyExc_UnicodeError");
    PyExc_UnicodeTranslateError_p = py_load_ptr("PyExc_UnicodeTranslateError");
    PyExc_ValueError_p = py_load_ptr("PyExc_ValueError");
    PyExc_ZeroDivisionError_p = py_load_ptr("PyExc_ZeroDivisionError");

    PyExc_BytesWarning_p = py_load_ptr("PyExc_BytesWarning");
    PyExc_DeprecationWarning_p = py_load_ptr("PyExc_DeprecationWarning");
    PyExc_FutureWarning_p = py_load_ptr("PyExc_FutureWarning");
    PyExc_ImportWarning_p = py_load_ptr("PyExc_ImportWarning");
    PyExc_PendingDeprecationWarning_p = py_load_ptr("PyExc_PendingDeprecationWarning");
    PyExc_ResourceWarning_p = py_load_ptr("PyExc_ResourceWarning");
    PyExc_RuntimeWarning_p = py_load_ptr("PyExc_RuntimeWarning");
    PyExc_SyntaxWarning_p = py_load_ptr("PyExc_SyntaxWarning");
    PyExc_UnicodeWarning_p = py_load_ptr("PyExc_UnicodeWarning");
    PyExc_UserWarning_p = py_load_ptr("PyExc_UserWarning");
    PyExc_Warning_p = py_load_ptr("PyExc_Warning");
}

// py_init initializes Python stuff.
// libpython3 must be a full path to the libpython3.XX.so library.
//
// It returns NULL on success or an error message in a case of errors.
// This function needs to be called only once.
//
// This function MUST be called by the main Python thread only.
const char *py_init (const char *libpython3) {
    py_load_all(libpython3);

    if (py_error != NULL) {
        return py_error;
    }

    Py_InitializeEx_p(0);
    py_main_thread = PyEval_SaveThread_p();

    return py_error;
}

// py_new_interp returns a new Python interpreter.
//
// This function MUST be called by the main Python thread only.
void py_new_interp (PyThreadState **tstate_p, PyInterpreterState **interp_p) {
    // This attaches the current OS thread to py_main_thread and
    // locks the GIL
    PyEval_RestoreThread_p(py_main_thread);

    // This creates an interpreter and the new thread state for it
    // and switches the current OS thread to the newly created
    // thread state
    PyThreadState *tstate = Py_NewInterpreter_p();
    PyInterpreterState *interp = PyInterpreterState_Get_p();

    // Here we switch back to the py_main_thread, which detaches
    // the newly created thread state from the OS thread
    PyThreadState_Swap_p(py_main_thread);

    // And now we finally release the GIL
    PyEval_SaveThread_p();

    *tstate_p = tstate;
    *interp_p = interp;
}

// py_interp_close closes the Python interpreter.
//
// This function MUST be called by the main Python thread only.
void py_interp_close (PyThreadState *tstate, PyInterpreterState *interp) {
    // This attaches the current OS thread to py_main_thread and
    // locks the GIL
    PyEval_RestoreThread_p(py_main_thread);

    // Now switch to the thread state of the interpreter being closed
    // and call Py_EndInterpreter to close it
    PyThreadState_Swap_p(tstate);
    Py_EndInterpreter_p(tstate);

    // Switch back to py_main_thread
    PyThreadState_Swap_p(py_main_thread);

    // And now we finally release the GIL
    PyEval_SaveThread_p();
}

// py_enter temporary attaches the calling thread to the
// Python interpreter.
//
// It must be called before any operations with the interpreter
// are performed and must be paired with the py_leave.
void py_enter(PyInterpreterState *interp) {
    py_enter_level ++;
    if (py_enter_level == 1) {
        PyThreadState *tstate = PyThreadState_New_p(interp);

        assert(py_thread_current == NULL);
        py_thread_current = tstate;

        PyEval_AcquireThread_p(tstate);
    }
}

// py_leave detaches the calling thread from the Python interpreter.
void py_leave (void) {
    assert(py_thread_current != NULL);

    if (py_enter_level == 1) {
        PyThreadState *tstate = py_thread_current;
        py_thread_current = NULL;

        PyThreadState_Clear_p(tstate);
        PyEval_ReleaseThread_p(tstate);
        PyThreadState_Delete_p(tstate);
    }
    py_enter_level --;
}

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
                     bool expr, PyObject **res, long *lineno) {

    // Initialize lineno into the safe default
    *lineno = -1;

    // Obtain the __main__ module reference and its namespace
    PyObject *main_module = PyImport_AddModule_p("__main__");
    if (main_module == NULL) {
        return false;
    }

    PyObject *dict = PyModule_GetDict_p(main_module);

    // Compile the statement
    int      mode = expr ? Py_eval_input : Py_file_input;
    PyObject *code = Py_CompileString_p(s, file, mode);
    if (code == NULL) {
        return false;
    }

    // Execute the statement, release code object
    PyObject *ret = PyEval_EvalCode_p(code, dict, dict);

    // In a case of error, locate the last frame that belongs
    // to the our code.
    if (PyErr_Occurred_p()) {
        PyObject     *exc, *val, *tb;

        PyErr_Fetch_p(&exc, &val, &tb);

        if (tb != NULL) {
            PyObject     *frame = tb;
            PyTypeObject *frame_t = py_obj_type(frame);

            do {
                PyObject *tb_frame = py_obj_getattr_int(frame, "tb_frame");
                PyObject *f_code = py_obj_getattr_int(tb_frame, "f_code");

                if (f_code == code) {
                    PyObject *loc = py_obj_getattr_int(frame, "tb_lineno");
                    uint64_t tmp = (uint64_t) PyLong_AsUnsignedLongLong_p(loc);

                    if (PyErr_Occurred_p()) {
                        PyErr_Clear_p();
                    } else {
                        *lineno = (long) tmp;
                    }
                }

                frame = py_obj_getattr_int(frame, "tb_next");
                if (frame == NULL) {
                    break;
                }
            } while (frame != NULL && py_obj_type(frame) == frame_t);

            PyErr_Restore_p(exc, val, tb);
        }
    }

    // Release the code object
    Py_DecRef_p(code);

    // Now interpret the result
    if (ret == NULL) {
        return false;
    }

    if (!expr) {
        Py_DecRef_p(ret);
        ret = NULL;
    }

    *res = ret;
    return true;
}

// py_interp_load loads (imports) string as a Python module.
// It returns, via the 'res' pointer, the strong reference
// to the Python object of the loaded module.
//
// The name parameter becomes the module name, while the
// file parameter used for diagnostics messages and
// indicated the input file name.
bool py_interp_load (const char *s, const char *name, const char *file,
                    PyObject **res) {

    // Obtain the __main__ module reference and its namespace
    PyObject *main_module = PyImport_AddModule_p("__main__");
    if (main_module == NULL) {
        return false;
    }

    PyObject *dict = PyModule_GetDict_p(main_module);

    // Compile the statement
    PyObject *code = Py_CompileString_p(s, file, Py_file_input);
    if (code == NULL) {
        return false;
    }

    // Add the module
    PyObject *module = PyImport_ExecCodeModule_p(name, code);
    if (module == NULL) {
        return false;
    }

    int rc = PyDict_SetItemString_p(dict, name, module);
    if (rc == 0) {
        *res = module;
    } else {
        Py_DecRef_p(module);
    }

    return rc == 0;
}

// py_obj_is_bool reports if PyObject is PyBool_Type
bool py_obj_is_bool (PyObject *x) {
    return Py_TYPE(x) == PyBool_Type_p;
}

// py_obj_is_byte_array reports if PyObject is PyByteArray_Type or its subclass.
bool py_obj_is_byte_array (PyObject *x) {
    return PyType_IsSubtype_p(Py_TYPE(x), PyByteArray_Type_p) != 0;
}

// py_obj_is_bytes reports if PyObject is PyBytes_Type or its subclass.
bool py_obj_is_bytes (PyObject *x) {
    return PyType_IsSubtype_p(Py_TYPE(x), PyBytes_Type_p) != 0;
}

// py_obj_is_complex reports if PyObject is PyComplex_Type or its subclass.
bool py_obj_is_complex (PyObject *x) {
    return PyType_IsSubtype_p(Py_TYPE(x), PyComplex_Type_p) != 0;
}

// py_obj_is_float reports if PyObject is PyFloat_Type or its subclass.
bool py_obj_is_float (PyObject *x) {
    return PyType_IsSubtype_p(Py_TYPE(x), PyFloat_Type_p) != 0;
}

// py_obj_is_long reports if PyObject is PyLong_Type or its subclass.
bool py_obj_is_long (PyObject *x) {
    if (py_obj_is_bool(x)) {
        // Python PyBool_Type is subclass of PyLong_Type which we
        // want to hide here.
        return false;
    }

    unsigned long flags = PyType_GetFlags_p(Py_TYPE(x));
    return (flags & Py_TPFLAGS_LONG_SUBCLASS) != 0;
}

// py_obj_is_map reports if PyObject is a map (dict, namedtyple, ...).
bool py_obj_is_map (PyObject *x) {
    unsigned long flags = PyType_GetFlags_p(Py_TYPE(x));
    return (flags & Py_TPFLAGS_DICT_SUBCLASS) != 0;
}

// py_obj_is_seq reports if PyObject is a sequence (list, tuple, ...).
bool py_obj_is_seq (PyObject *x) {
    unsigned long flags = PyType_GetFlags_p(Py_TYPE(x));
    unsigned long mask = Py_TPFLAGS_LIST_SUBCLASS | Py_TPFLAGS_TUPLE_SUBCLASS;
    return (flags & mask) != 0;
}

// py_obj_is_unicode reports if PyObject is PyUnicode_Type or its subclass.
bool py_obj_is_unicode (PyObject *x) {
    unsigned long flags = PyType_GetFlags_p(Py_TYPE(x));
    return (flags & Py_TPFLAGS_UNICODE_SUBCLASS) != 0;
}

// py_obj_ref increments the PyObject's reference count.
void py_obj_ref (PyObject *x) {
        Py_IncRef_p(x);
}

// py_obj_unref decrements the PyObject's reference count.
void py_obj_unref (PyObject *x) {
    Py_DecRef_p(x);
}

// py_obj_type returns PyObject's type.
PyTypeObject *py_obj_type (PyObject *x) {
    // Python 3.8 doesn't expose Py_TYPE() as a part of its stable API,
    // so we have to use the PyObject_Type() function.
    //
    // Note, PyObject_Type() returns the strong reference to the type
    // object while we need the borrowed reference, hence Py_DecRef_p()
    PyObject *t = PyObject_Type_p(x);
    Py_DecRef_p(t);
    return (PyTypeObject*) t;
}

// py_obj_str returns a string representation of the PyObject.
// This is the equivalent of the Python expression str(x).
//
// There is very subtle difference between py_obj_str and py_obj_repr.
// In general:
//   - Use py_obj_str if you want to print the string
//   - Use py_obj_repr if you want to process the string
PyObject *py_obj_str (PyObject *x) {
    return PyObject_Str_p(x);
}

// py_obj_repr returns a string representation of the PyObject.
// This is the equivalent of the Python expression repr(x).
//
// There is very subtle difference between py_obj_str and py_obj_repr.
// In general:
//   - Use py_obj_str if you want to print the string
//   - Use py_obj_repr if you want to process the string
PyObject *py_obj_repr (PyObject *x) {
    return PyObject_Repr_p(x);
}

// py_obj_length returns PyObject, in items. It works with any
// container objects (lists, tuples, dictionaries, ...).
// Returns -1 on a error.
ssize_t py_obj_length (PyObject *x) {
    return PyObject_Length_p(x);
}

// py_obj_keys returns PyObject mapping keys. It works for objects
// that supports mapping (see py_obj_is_map), i.e., dict etc.
//
// On success it returns PyList_Type or PyTuple_Type object that
// contains the keys. On error it returns NULL.
PyObject *py_obj_keys (PyObject *x) {
    return PyMapping_Keys_p(x);
}

// py_obj_hasattr reports if PyObject has the attribute with the
// specified name.
//
// It returns true on success, false on error and puts answer into
// its third parameter.
bool py_obj_hasattr(PyObject *x, const char *name, bool *answer) {
    PyObject *attr = NULL;

    // Try to retrieve an attribute
    if (py_obj_getattr(x, name, &attr)) {
        *answer = true;
        Py_DecRef_p(attr);
        return true;
    }

    // Don't treat PyExc_AttributeError as an error.
    // Just indicate that attribute didn't exist
    if (PyErr_Occurred_p() == PyExc_AttributeError_p) {
        *answer = false;
        PyErr_Clear_p();
        return true;
    }

    // Otherwise, propagate the error
    *answer = false;
    return false;
}

// py_obj_delattr deletes the attribute with the specified name.
// It returns true on success, false on error.
bool py_obj_delattr(PyObject *x, const char *name) {
    return py_obj_setattr(x, name, NULL);
}

// py_obj_getattr retrieves the attribute with the specified name.
// The returned answer, on success, contains a strong reference to PyObject.
// It returns true on success, false on error.
bool py_obj_getattr(PyObject *x, const char *name, PyObject **answer) {
    PyObject *attr = PyObject_GetAttrString_p(x, name);
    *answer = attr;
    return attr != NULL;
}

// py_obj_getattr_int is the py_obj_getattr version for the
// internal use.
//
// It has the following differences from the normal py_obj_getattr:
//   - the returned reference is borrowed, not strong
//   - input object may be NULL
//   - error is cleared if any occurs
static PyObject *py_obj_getattr_int(PyObject *x, const char *name) {
    PyObject *attr = NULL;

    if (x != NULL) {
        attr = PyObject_GetAttrString_p(x, name);

        if (attr != NULL) {
            Py_DecRef_p(attr);
        }else{
            PyErr_Clear_p();
        }
    }

    return attr;
}

// py_obj_getattr sets the attribute with the specified name.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_obj_setattr(PyObject *x, const char *name, PyObject *value) {
    return PyObject_SetAttrString_p(x, name, value) == 0;
}

// py_obj_hasitem reports if PyObject contains the item with the
// specified key.
//
// It returns true on success, false on error and puts answer into
// its third parameter.
bool py_obj_hasitem(PyObject *x, PyObject *key, bool *answer) {
    PyObject *item;
    bool     ok = py_obj_getitem(x, key, &item);

    *answer = false;
    if (item != NULL) {
        *answer = true;
        Py_DecRef_p(item);
    }

    return ok;
}

// py_obj_delitem deletes the item with the specified key.
// It returns true on success, false on error.
bool py_obj_delitem(PyObject *x, PyObject *key) {
    return PyObject_DelItem_p(x, key) == 0;
}

// py_obj_getitem retrieves the item with the specified key.
// The returned answer, on success, contains a strong reference to PyObject.
// It returns true on success, false on error.
bool py_obj_getitem(PyObject *x, PyObject *key, PyObject **answer) {
    *answer = PyObject_GetItem_p(x, key);

    if (*answer != NULL) {
        return true;
    }

    PyObject *err = PyErr_Occurred_p();
    if (err == NULL || err == PyExc_KeyError_p) {
        PyErr_Clear_p();
        return true;
    }

    return false;
}

// py_obj_setitem sets the item with the specified key.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_obj_setitem(PyObject *x, PyObject *key, PyObject *value) {
    return PyObject_SetItem_p(x, key, value) == 0;
}

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
PyObject *py_obj_call(PyObject *x, PyObject *args, PyObject *kwargs) {
    return PyObject_Call_p(x, args, kwargs);
}

// py_obj_callable reports if object is callable.
// This function always succeeds.
bool py_obj_callable(PyObject *x) {
    return PyCallable_Check_p(x) != 0;
}

// py_err_fetch fetches and clears last error.a
// If there is no pending error, all pointers will be set to NULL.
void py_err_fetch (PyObject **etype, PyObject **evalue, PyObject **trace) {
    PyObject *exc, *val, *tb;

    PyErr_Fetch_p(&exc, &val, &tb);
    if ((exc != NULL) || (val != NULL) || (tb != NULL)) {
        PyErr_NormalizeException_p(&exc, &val, &tb);
    }

    *etype = exc;
    *evalue = val;
    *trace = tb;
}

// py_bool_make makes a new PyBool_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_bool_make(bool val) {
    return PyBool_FromLong_p((long) val);
}

// py_bytes_get obtains content of the Python bytes object.
// It returns true on success, false on error.
bool py_bytes_get (PyObject *x, void **data, size_t *size) {
    Py_ssize_t sz = 0;
    int        rc = PyBytes_AsStringAndSize_p(x, (char**) data, &sz);
    *size = (size_t) sz;
    return rc == 0;
}

// py_bytes_make makes a new PyBytes_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_bytes_make(const void *data, size_t size) {
    return PyBytes_FromStringAndSize_p(data, size);
}

// py_bytearray_get obtains content of the Python bytearray object.
// It returns true on success, false on error.
bool py_bytearray_get (PyObject *x, void **data, size_t *size) {
    Py_ssize_t sz = PyByteArray_Size_p(x);
    char       *bytes = PyByteArray_AsString_p(x);

    if (sz >= 0 && bytes != NULL) {
        *data = (void*) bytes;
        *size = (size_t) sz;
        return true;
    }

    return false;
}

// py_complex_get obtains content of the Python complex object.
// It returns true on success, false on error.
bool py_complex_get (PyObject *x, double *real, double *imag) {
    double r = PyComplex_RealAsDouble_p(x);
    double i = PyComplex_ImagAsDouble_p(x);

    if ((r == -1.0 || i == -1.0) && PyErr_Occurred_p() != NULL) {
        PyErr_Clear_p();
        return false;
    }

    *real = r;
    *imag = i;

    return true;
}

// py_complex_make makes a new PyComlex_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_complex_make(double real, double imag) {
    return PyComplex_FromDoubles_p(real, imag);
}

// py_dict_make makes a new PyDict_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_dict_make(void) {
    return PyDict_New_p();
}

// py_float_get obtains content of the Python float object.
// It returns true on success, false on error.
bool py_float_get (PyObject *x, double *val) {
    double v = PyFloat_AsDouble_p(x);

    if (v == -1.0 && PyErr_Occurred_p() != NULL) {
        PyErr_Clear_p();
        return false;
    }

    *val = v;

    return true;
}

// py_float_make makes a new PyFloat_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_float_make(double val) {
    return PyFloat_FromDouble_p(val);
}

// py_list_make makes a new PyList_Type object of the specified size.
// The newly created list MUST be fully populated with the py_list_set
// calls before it can be safely passed to Python interpreter.
// It returns strong object reference on success, NULL on an error.
PyObject *py_list_make(size_t len) {
    return PyList_New_p(len);
}

// py_list_set retrieves value of the list item at the given position.
// It returns strong object reference on success, NULL on an error.
PyObject *py_list_get(PyObject *list, int index) {
    PyObject *item = PyList_GetItem_p(list, index);
    if (item != NULL) {
        Py_IncRef_p(item);
    }
    return item;
}

// py_list_set sets value of the list item at the given position.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_list_set(PyObject *list, int index, PyObject *val) {
    Py_IncRef_p(val);
    return PyList_SetItem_p(list, index, val) == 0;
}

// py_long_get obtains PyObject's value as int64_t.
// If value doesn't fit C long, overflow flag is set.
//
// It returns true on success, false on error.
bool py_long_get_int64 (PyObject *x, int64_t *val, bool *overflow) {
    long long tmp;
    bool      ok = true, ovf = false;

    tmp = (int64_t) PyLong_AsLongLong_p(x);
    if (tmp == -1) {
        PyObject *err = PyErr_Occurred_p();
        if (err == PyExc_OverflowError_p) {
            ovf = true;
            PyErr_Clear_p();
        } else if (err != NULL) {
            ok = false;
        }
    }

    if (!(INT64_MIN <= tmp && tmp <= INT64_MAX)) {
        ovf = true;
    }

    *val = (int64_t) tmp;
    *overflow = ovf;

    return ok;
}

// py_long_get obtains PyObject's value as uint64_t.
// If value doesn't fit C long, overflow flag is set.
//
// It returns true on success, false on error.
bool py_long_get_uint64 (PyObject *x, uint64_t *val, bool *overflow) {
    unsigned long long  tmp;
    bool                ok = true, ovf = false;

    tmp = (uint64_t) PyLong_AsUnsignedLongLong_p(x);
    if (tmp == (unsigned long long) -1) {
        PyObject *err = PyErr_Occurred_p();
        if (err == PyExc_OverflowError_p) {
            ovf = true;
            PyErr_Clear_p();
        } else if (err != NULL) {
            ok = false;
        }
    }

    if (!(tmp <= UINT64_MAX)) {
        ovf = true;
    }

    *val = (int64_t) tmp;
    *overflow = ovf;

    return ok;
}

// py_long_from_int64 makes a new PyLong_Type object from int64_t value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_long_from_int64(int64_t val) {
    return PyLong_FromLongLong_p((long long) val);
}

// py_long_from_uint64 makes a new PyLong_Type object from uint64_t value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_long_from_uint64(uint64_t val) {
    return PyLong_FromUnsignedLongLong_p((unsigned long long) val);
}

// py_long_from_string makes a new PyLong_Type object from string value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_long_from_string(const char *val) {
    return PyLong_FromString_p(val, NULL, 0);
}

// py_seq_set retrieves value of the sequence item at the given position.
// It returns strong object reference on success, NULL on an error.
PyObject *py_seq_get(PyObject *tuple, int index) {
    PyObject *item = PySequence_GetItem_p(tuple, index);
    if (item != NULL) {
        Py_IncRef_p(item);
    }
    return item;
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
Py_UCS4 *py_str_get (PyObject *str, Py_UCS4 *buf, size_t len) {
    return PyUnicode_AsUCS4_p(str, buf, len, 0);
}

// py_str_make makes a new PyLong_Type object from string value.
// It returns strong object reference on success, NULL on an error.
PyObject *py_str_make(const char *val, size_t len) {
    return PyUnicode_FromStringAndSize_p(val, len);
}

// py_tuple_make makes a new PyTuple_Type object of the specified size.
// The newly created tuple MUST be fully populated with the py_tuple_set
// calls before it can be safely passed to Python interpreter.
// It returns strong object reference on success, NULL on an error.
PyObject *py_tuple_make(size_t len) {
    return PyTuple_New_p(len);
}

// py_tuple_get retrieves value of the tuple item at the given position.
// It returns strong object reference on success, NULL on an error.
PyObject *py_tuple_get(PyObject *tuple, int index) {
    PyObject *item = PyTuple_GetItem_p(tuple, index);
    if (item != NULL) {
        Py_IncRef_p(item);
    }
    return item;
}

// py_tuple_set sets value of the tuple item at the given position.
// Internally, it creates a new strong reference to the object.
// It returns true on success, false on error.
bool py_tuple_set(PyObject *tuple, int index, PyObject *val) {
    Py_IncRef_p(val);
    return PyTuple_SetItem_p(tuple, index, val) == 0;
}

// py_cfunction_make makes a new PyCFunction_Type object.
// The caller must ensure that ml outlives the returned object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_cfunction_make(PyMethodDef *ml, PyObject *self) {
    return PyCFunction_NewEx_p(ml, self, NULL);
}

// py_capsule_make makes a new PyCapsule_Type object.
// It returns strong object reference on success, NULL on an error.
PyObject *py_capsule_make(void *pointer, const char *name,
                          PyCapsule_Destructor destructor) {
    return PyCapsule_New_p(pointer, name, destructor);
}

// py_capsule_get_ptr returns pointer, encapsulated into
// the capsule.
//
// The name must match exactly the name used for py_capsule_make.
void *py_capsule_get_ptr(PyObject *capsule, const char *name) {
    return PyCapsule_GetPointer_p(capsule, name);
}

// py_set_err sets the error condition within the Python interpreter,
// which raises exception at the Python side.
//
// The callback function, registered by the py_cfunction_make, MUST
// set the error, if it returns NULL.
void py_set_err(PyObject *extype, const char *msg) {
    PyErr_SetString_p(extype, msg);
}

// vim:ts=8:sw=4:et
