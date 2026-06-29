// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Python exceptions

package cpython

// #include "cpython.h"
import "C"

// Except represents a Python exception by its name.
//
// Normally, Python values are bound to the [Python] interpreter,
// but using full exception objects would be inconvenient.
// Therefore, exceptions are represented by their Python name instead.
type Except string

// Standard exceptions
const (
	// Error exceptions
	ArithmeticError        Except = "ArithmeticError"
	AssertionError         Except = "AssertionError"
	AttributeError         Except = "AttributeError"
	BaseException          Except = "BaseExcept"
	BlockingIOError        Except = "BlockingIOError"
	BrokenPipeError        Except = "BrokenPipeError"
	BufferError            Except = "BufferError"
	ChildProcessError      Except = "ChildProcessError"
	ConnectionAbortedError Except = "ConnectionAbortedError"
	ConnectionError        Except = "ConnectionError"
	ConnectionRefusedError Except = "ConnectionRefusedError"
	ConnectionResetError   Except = "ConnectionResetError"
	EOFError               Except = "EOFError"
	Exception              Except = "Exception"
	FileExistsError        Except = "FileExistsError"
	FileNotFoundError      Except = "FileNotFoundError"
	FloatingPointError     Except = "FloatingPointError"
	GeneratorExit          Except = "GeneratorExit"
	ImportError            Except = "ImportError"
	IndentationError       Except = "IndentationError"
	IndexError             Except = "IndexError"
	InterruptedError       Except = "InterruptedError"
	IsADirectoryError      Except = "IsADirectoryError"
	KeyboardInterrupt      Except = "KeyboardInterrupt"
	KeyError               Except = "KeyError"
	LookupError            Except = "LookupError"
	MemoryError            Except = "MemoryError"
	ModuleNotFoundError    Except = "ModuleNotFoundError"
	NameError              Except = "NameError"
	NotADirectoryError     Except = "NotADirectoryError"
	NotImplementedError    Except = "NotImplementedError"
	OSError                Except = "OSError"
	OverflowError          Except = "OverflowError"
	PermissionError        Except = "PermissionError"
	ProcessLookupError     Except = "ProcessLookupError"
	RecursionError         Except = "RecursionError"
	ReferenceError         Except = "ReferenceError"
	RuntimeError           Except = "RuntimeError"
	StopAsyncIteration     Except = "StopAsyncIteration"
	StopIteration          Except = "StopIteration"
	SyntaxError            Except = "SyntaxError"
	SystemError            Except = "SystemError"
	SystemExit             Except = "SystemExit"
	TabError               Except = "TabError"
	TimeoutError           Except = "TimeoutError"
	TypeError              Except = "TypeError"
	UnboundLocalError      Except = "UnboundLocalError"
	UnicodeDecodeError     Except = "UnicodeDecodeError"
	UnicodeEncodeError     Except = "UnicodeEncodeError"
	UnicodeError           Except = "UnicodeError"
	UnicodeTranslateError  Except = "UnicodeTranslateError"
	ValueError             Except = "ValueError"
	ZeroDivisionError      Except = "ZeroDivisionError"

	// Warning exceptions
	BytesWarning              Except = "BytesWarning"
	DeprecationWarning        Except = "DeprecationWarning"
	FutureWarning             Except = "FutureWarning"
	ImportWarning             Except = "ImportWarning"
	PendingDeprecationWarning Except = "PendingDeprecationWarning"
	ResourceWarning           Except = "ResourceWarning"
	RuntimeWarning            Except = "RuntimeWarning"
	SyntaxWarning             Except = "SyntaxWarning"
	UnicodeWarning            Except = "UnicodeWarning"
	UserWarning               Except = "UserWarning"
	Warning                   Except = "Warning"
)

// Error implements error interface for the exception.
func (ex Except) Error() string {
	return string(ex)
}

// object returns Python object for the exception.
func (ex Except) object() pyObject {
	pyobj := exceptTable[ex]
	if pyobj == nil {
		// Fallback to the SystemError for unknown exception
		pyobj = C.PyExc_SystemError_p
	}
	return pyobj
}

// exceptTable maps Except names of the standard exceptions to the
// corresponding Python objects.
var exceptTable map[Except]pyObject

// exceptInit initializes exceptTable.
//
// It is called during libpython3.so initialization, when symbols
// already loaded from the library.
func exceptInit() {
	exceptTable = map[Except]pyObject{
		// Errors
		ArithmeticError:        C.PyExc_ArithmeticError_p,
		AssertionError:         C.PyExc_AssertionError_p,
		AttributeError:         C.PyExc_AttributeError_p,
		BaseException:          C.PyExc_BaseException_p,
		BlockingIOError:        C.PyExc_BlockingIOError_p,
		BrokenPipeError:        C.PyExc_BrokenPipeError_p,
		BufferError:            C.PyExc_BufferError_p,
		ChildProcessError:      C.PyExc_ChildProcessError_p,
		ConnectionAbortedError: C.PyExc_ConnectionAbortedError_p,
		ConnectionError:        C.PyExc_ConnectionError_p,
		ConnectionRefusedError: C.PyExc_ConnectionRefusedError_p,
		ConnectionResetError:   C.PyExc_ConnectionResetError_p,
		EOFError:               C.PyExc_EOFError_p,
		Exception:              C.PyExc_Exception_p,
		FileExistsError:        C.PyExc_FileExistsError_p,
		FileNotFoundError:      C.PyExc_FileNotFoundError_p,
		FloatingPointError:     C.PyExc_FloatingPointError_p,
		GeneratorExit:          C.PyExc_GeneratorExit_p,
		ImportError:            C.PyExc_ImportError_p,
		IndentationError:       C.PyExc_IndentationError_p,
		IndexError:             C.PyExc_IndexError_p,
		InterruptedError:       C.PyExc_InterruptedError_p,
		IsADirectoryError:      C.PyExc_IsADirectoryError_p,
		KeyboardInterrupt:      C.PyExc_KeyboardInterrupt_p,
		KeyError:               C.PyExc_KeyError_p,
		LookupError:            C.PyExc_LookupError_p,
		MemoryError:            C.PyExc_MemoryError_p,
		ModuleNotFoundError:    C.PyExc_ModuleNotFoundError_p,
		NameError:              C.PyExc_NameError_p,
		NotADirectoryError:     C.PyExc_NotADirectoryError_p,
		NotImplementedError:    C.PyExc_NotImplementedError_p,
		OSError:                C.PyExc_OSError_p,
		OverflowError:          C.PyExc_OverflowError_p,
		PermissionError:        C.PyExc_PermissionError_p,
		ProcessLookupError:     C.PyExc_ProcessLookupError_p,
		RecursionError:         C.PyExc_RecursionError_p,
		ReferenceError:         C.PyExc_ReferenceError_p,
		RuntimeError:           C.PyExc_RuntimeError_p,
		StopAsyncIteration:     C.PyExc_StopAsyncIteration_p,
		StopIteration:          C.PyExc_StopIteration_p,
		SyntaxError:            C.PyExc_SyntaxError_p,
		SystemError:            C.PyExc_SystemError_p,
		SystemExit:             C.PyExc_SystemExit_p,
		TabError:               C.PyExc_TabError_p,
		TimeoutError:           C.PyExc_TimeoutError_p,
		TypeError:              C.PyExc_TypeError_p,
		UnboundLocalError:      C.PyExc_UnboundLocalError_p,
		UnicodeDecodeError:     C.PyExc_UnicodeDecodeError_p,
		UnicodeEncodeError:     C.PyExc_UnicodeEncodeError_p,
		UnicodeError:           C.PyExc_UnicodeError_p,
		UnicodeTranslateError:  C.PyExc_UnicodeTranslateError_p,
		ValueError:             C.PyExc_ValueError_p,
		ZeroDivisionError:      C.PyExc_ZeroDivisionError_p,

		// Warnings
		BytesWarning:              C.PyExc_BytesWarning_p,
		DeprecationWarning:        C.PyExc_DeprecationWarning_p,
		FutureWarning:             C.PyExc_FutureWarning_p,
		ImportWarning:             C.PyExc_ImportWarning_p,
		PendingDeprecationWarning: C.PyExc_PendingDeprecationWarning_p,
		ResourceWarning:           C.PyExc_ResourceWarning_p,
		RuntimeWarning:            C.PyExc_RuntimeWarning_p,
		SyntaxWarning:             C.PyExc_SyntaxWarning_p,
		UnicodeWarning:            C.PyExc_UnicodeWarning_p,
		UserWarning:               C.PyExc_UserWarning_p,
		Warning:                   C.PyExc_Warning_p,
	}
}
