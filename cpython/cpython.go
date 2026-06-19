// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// CPython glue -- the Go side

package cpython

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"unsafe"
)

// #cgo pkg-config: python3
// #cgo LDFLAGS: -l dl
//
// #include "cpython.h"
import "C"

type (
	// pyThreadState is the Go name for the *C.PyThreadState
	pyThreadState = *C.PyThreadState

	// pyInterpState is the Go name for the *C.PyInterpreterState
	pyInterpState = *C.PyInterpreterState

	// pyObject is the Go name for the *C.PyObject
	pyObject = *C.PyObject

	// pyObject is the Go name for the *C.PyTypeObject
	pyTypeObject = *C.PyTypeObject
)

var (
	// pyInitError holds Python initialization error, if any.
	pyInitError error
)

// pyInterpNewRequestChan is the channel where requests to create
// new Python sub-interpretes are sent to.
//
// These requests are handled by the dedicated thread. The request
// itself is the channel of pyInterpHandles, where response is sent.
var pyInterpNewRequestChan = make(chan chan pyInterpHandles)

var pyInterpDeleteRequestChan = make(chan pyInterpHandles)

// pyInterpHandles contains pointer to PyInterpreterState
// and its main PyThreadState
type pyInterpHandles struct {
	tstate pyThreadState
	interp pyInterpState
}

// pyNewInterp creates a new Python sub-interpreter and returns
// pointer to its main thread state.
func pyNewInterp() (pyThreadState, pyInterpState, error) {
	if pyInitError != nil {
		return nil, nil, pyInitError
	}

	rsp := make(chan pyInterpHandles)
	pyInterpNewRequestChan <- rsp
	response := <-rsp

	return response.tstate, response.interp, nil
}

// pyInterpDelete releases the Python sub-interpreter
func pyInterpDelete(tstate pyThreadState, interp pyInterpState) {
	pyInterpDeleteRequestChan <- pyInterpHandles{tstate, interp}
}

// pyLocateLibPython locates the full path to the libpython3.XX.so library
func pyLocateLibPython() (string, error) {
	script := ""
	script += "import sysconfig;"
	script += "import os;"
	script += "dir=sysconfig.get_config_var('LIBDIR');"
	script += "lib=sysconfig.get_config_var('INSTSONAME');"
	script += "print(os.path.join(dir,lib),end='');"

	python := os.Getenv("MFP_PYTHON")
	if python == "" {
		python = "python3"
	}

	cmd := exec.Command(python, "-c", script)
	out, err := cmd.Output()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			err = errors.New(string(bytes.TrimSpace(ee.Stderr)))
		}
	}

	return string(out), err
}

// pyInterpThread runs Python dedicated thread.
//
// We need this thread, because CPython pollutes the thread local
// storage with its own stuff while doing operations like Python
// initialization and new interpreter creation, and we don't want
// this litter to be distributed across all Go threads.
//
// Once initialization is completed, pyInterpThread signals it
// over the supplied sync.WaitGroup parameter.
func pyInterpThread(initilized *sync.WaitGroup) {
	// We need the dedicated thread...
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Locate and initialize Python library
	var lib string
	lib, pyInitError = pyLocateLibPython()
	if pyInitError == nil {
		clib := C.CString(lib)
		msg := C.py_init(clib)
		C.free(unsafe.Pointer(clib))

		pyInitErrorCheck(msg)
	}

	// Initialize exceptions
	if pyInitError == nil {
		exceptInit()
	}

	// Notify caller that initialization is completed
	initilized.Done()

	// If no error, serve incoming requests
	if pyInitError == nil {
		for {
			select {
			case rq := <-pyInterpNewRequestChan:
				var tstate pyThreadState
				var interp pyInterpState
				C.py_new_interp(&tstate, &interp)
				rq <- pyInterpHandles{tstate, interp}

			case rq := <-pyInterpDeleteRequestChan:
				C.py_interp_close(rq.tstate, rq.interp)
			}
		}
	}
}

// pyInitErrorCheck sets the pyInitError, if error message is not nil.
func pyInitErrorCheck(msg *C.char) {
	if msg != nil {
		pyInitError = errors.New(C.GoString(msg))
	}
}

// pyInitErrorCheckTest is the test interface to the pyInitErrorCheck
// function. We cannot call pyInitErrorCheck directly from tests, because
// CGo is not available in the _test.go files, so we cannot construct
// *C.char string out of there.
func pyInitErrorCheckTest(s string) {
	msg := C.CString(s)
	defer C.free(unsafe.Pointer(msg))
	pyInitErrorCheck(msg)
}

// init starts a dedicated Python thread.
func init() {
	var initilized sync.WaitGroup

	initilized.Add(1)
	go pyInterpThread(&initilized)
	initilized.Wait()
}
