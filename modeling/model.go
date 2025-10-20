// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device model

package modeling

import (
	"io"
	"os"
	"strings"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// Model defines the whole characteristics of the MFP device being
// modeled, including the IPP printer attributes, eSCL and WSD
// scanner capabilities, scripting hooks, used to modify device
// behavior and the Python interpreter instance, used to execute
// these hooks.
type Model struct {
	py              *cpython.Python
	ippPrinterAttrs *ipp.PrinterAttributes
	esclScanCaps    *escl.ScannerCapabilities

	// Modules
	modQuery *cpython.Object // query.py
	modEscl  *cpython.Object // escl.py
	modIPP   *cpython.Object // IPP.py

	// Important Python class constructors
	clsDict            *cpython.Object // dict
	clsHTTPMessage     *cpython.Object // query.HTTPMessage
	clsQuery           *cpython.Object // query.Query constructor
	clsUUID            *cpython.Object // uuid.UUID
	clsDateTimeFromISO *cpython.Object // datetime.datetime.fromisoformat

	// Python hooks for eSCL
	esclOnScanJobsRequestScriptlet      *cpython.Object
	esclOnNextDocumentResponseScriptlet *cpython.Object

	// eSCL state
	esclScanSettings escl.ScanSettings
}

// NewModel creates a new Model with empty printer/scanner parameters.
// Use [Model.Close] to release resources owned by the Model.
func NewModel() (*Model, error) {
	// Create Python interpreter
	py, err := cpython.NewPython()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			py.Close()
		}
	}()

	// Create Model structure
	model := &Model{py: py}

	// Load startup script
	err = py.Exec(embedPyInit, "init.py")
	if err != nil {
		return nil, err
	}

	// Load modules
	model.modQuery, err = py.Load(embedPyQuery, "query", "query.py")
	if err != nil {
		return nil, err
	}

	model.modEscl, err = py.Load(embedPyEscl, "escl", "escl.py")
	if err != nil {
		return nil, err
	}

	model.modIPP, err = py.Load(embedPyIPP, "ipp", "ipp.py")
	if err != nil {
		return nil, err
	}

	// Load commonly used class constructors
	model.clsDict, err = py.Eval("dict")
	if err != nil {
		return nil, err
	}

	model.clsQuery, err = py.Eval("query.Query")
	if err != nil {
		return nil, err
	}

	model.clsHTTPMessage, err = py.Eval("query.HTTPMessage")
	if err != nil {
		return nil, err
	}

	model.clsUUID, err = py.Eval("UUID")
	if err != nil {
		return nil, err
	}

	model.clsDateTimeFromISO, err = py.Eval("datetime.fromisoformat")
	if err != nil {
		return nil, err
	}

	// Verify things
	assert.Must(model.clsDict.IsCallable())
	assert.Must(model.clsHTTPMessage.IsCallable())
	assert.Must(model.clsQuery.IsCallable())
	assert.Must(model.clsUUID.IsCallable())

	return model, nil
}

// Close closes the Model and releases all resources associated
// with it.
func (model *Model) Close() {
	model.py.Close()
	model.py = nil
}

// Reset resets the Modal into its initial state.
func (model *Model) Reset() error {
	model2, err := NewModel()
	if err != nil {
		return err
	}

	model.py.Close()
	*model = *model2
	return nil
}

// Write writes model into the [io.Writer]
func (model *Model) Write(w io.Writer) error {
	var escl, ipp string

	// Format parts
	if model.esclScanCaps != nil {
		obj, err := model.pyExportStruct(model.esclScanCaps)
		if err != nil {
			return err
		}

		escl, err = formatPython(obj)
		if err != nil {
			return err
		}
	}

	if model.ippPrinterAttrs != nil {
		obj, err := model.pyExportIPP(model.ippPrinterAttrs)
		if err != nil {
			return err
		}

		ipp, err = formatPython(obj)
		if err != nil {
			return err
		}
	}

	// Expand callback
	expand := func(name string) string {
		switch name {
		case "ESCL":
			return escl
		case "IPP":
			return ipp
		}

		return ""
	}

	// Split template into lines. Trim terminating empty lines, if any.
	template := strings.Split(embedPyModel, "\n")
	for len(template) > 0 && template[len(template)-1] == "" {
		template = template[:len(template)-1]
	}

	// Expand template
	skip := true
	for _, t := range template {
		switch {
		case strings.HasPrefix(t, "#-escl"):
			skip = model.esclScanCaps == nil
		case strings.HasPrefix(t, "#-ipp"):
			skip = model.ippPrinterAttrs == nil
		case strings.HasPrefix(t, "#-"):
			skip = false
		default:
			if !skip {
				s := os.Expand(t, expand)
				_, err := w.Write(([]byte)(s + "\n"))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Read reads model from the [io.Reader]
// The filename parameter required for the diagnostics messages.
func (model *Model) Read(filename string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	err = model.py.Exec(string(data), filename)
	if err != nil {
		return err
	}

	err = model.esclLoad()
	if err != nil {
		return err
	}

	return nil
}

// Save writes model into the disk file.
func (model *Model) Save(file string) error {
	// Open the file
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	fp, err := os.OpenFile(file, flags, 0644)
	if err != nil {
		return err
	}

	// Write model data
	err = model.Write(fp)
	err2 := fp.Close()

	if err == nil {
		err = err2
	}

	return err
}

// Load reads model from the disk file.
func (model *Model) Load(file string) error {
	fp, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fp.Close()

	return model.Read(file, fp)
}
