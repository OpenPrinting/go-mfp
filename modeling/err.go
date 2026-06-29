// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Errors

package modeling

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/OpenPrinting/go-mfp/cpython"
)

// errImport represents the error that may occur during
// importing the Python object into the Go structure
type errImport struct {
	path []string // Path over attribute names
	err  error    // Underlying error
}

// errImportWrap wraps error into the errImport.
// name is the name of the attribute the error is related to.
func errImportWrap(name string, err error) error {
	if e, ok := err.(errImport); ok {
		return errImport{
			path: append([]string{name}, e.path...),
			err:  e.err,
		}
	}

	return errImport{
		path: []string{name},
		err:  err,
	}
}

// Error implements the error interface for errImport
func (e errImport) Error() string {
	buf := strings.Builder{}
	for _, p := range e.path {
		if buf.Len() > 0 && !strings.HasPrefix(p, "[") {
			buf.WriteByte('.')
		}
		buf.WriteString(p)
	}

	buf.Write([]byte(": "))
	buf.WriteString(e.err.Error())

	return buf.String()
}

// Unwrap "unwraps" the error.
func (e errImport) Unwrap() error {
	return e.err
}

// errPy2Go returns a conversion error for conversion
// from the Python object to the Go value
func errPy2Go(from *cpython.Object, to reflect.Value) error {
	return fmt.Errorf("can't convert %s to %s",
		from.TypeName(),
		to.Type().String(),
	)
}
