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
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
)

// Model defines the whole characteristics of the MFP device being
// modeled, including the IPP printer attributes, eSCL and WSD
// scanner capabilities, scripting hooks, used to modify device
// behavior and the Python interpreter instance, used to execute
// these hooks.
type Model struct {
	py *cpython.Python

	// Scanner and printer capabilities, protocol-specific
	ippPrinterAttrs *ipp.PrinterAttributes
	esclScanCaps    *escl.ScannerCapabilities
	wsdScanCaps     *wsscan.GetScannerElementsResponse

	// Modules
	modHelpers *cpython.Object // helpers.py
	modQuery   *cpython.Object // query.py
	modIPP     *cpython.Object // ipp.py
	modEscl    *cpython.Object // escl.py
	modWSScan  *cpython.Object // wsd.py

	// Important Python class constructors
	clsHTTPMessage     *cpython.Object // query.HTTPMessage
	clsQuery           *cpython.Object // query.Query constructor
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
	model.modHelpers = py.Load(embedPyHelpers, "helpers", "helpers.py")
	if err := model.modHelpers.Err(); err != nil {
		return nil, err
	}

	model.modQuery = py.Load(embedPyQuery, "query", "query.py")
	if err := model.modQuery.Err(); err != nil {
		return nil, err
	}

	model.modIPP = py.Load(embedPyIPP, "ipp", "ipp.py")
	if err := model.modIPP.Err(); err != nil {
		return nil, err
	}

	model.modEscl = py.Load(embedPyEscl, "escl", "escl.py")
	if err := model.modEscl.Err(); err != nil {
		return nil, err
	}

	model.modWSScan = py.Load(embedPyWSD, "wsd", "wsd.py")
	if err := model.modWSScan.Err(); err != nil {
		return nil, err
	}

	// Load commonly used class constructors
	model.clsQuery = py.Eval("query.Query")
	if err := model.clsQuery.Err(); err != nil {
		return nil, err
	}

	model.clsHTTPMessage = py.Eval("query.HTTPMessage")
	if err := model.clsHTTPMessage.Err(); err != nil {
		return nil, err
	}

	model.clsDateTimeFromISO = py.Eval("datetime.fromisoformat")
	if err := model.clsDateTimeFromISO.Err(); err != nil {
		return nil, err
	}

	// Verify things
	assert.Must(model.clsHTTPMessage.IsCallable())
	assert.Must(model.clsQuery.IsCallable())
	assert.Must(model.clsDateTimeFromISO.IsCallable())

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
func (model *Model) Write(w io.Writer) (err error) {
	var ipp, escl, wsd string

	// Format parts
	if model.ippPrinterAttrs != nil {
		obj := ippExport(model.py, model.ippPrinterAttrs)
		ipp, err = formatPython(obj)
		if err != nil {
			return
		}
	}

	if model.esclScanCaps != nil {
		obj := structExport(model.py, keywordMapESCL, model.esclScanCaps)
		escl, err = formatPython(obj)
		if err != nil {
			return
		}
	}

	if model.wsdScanCaps != nil {
		obj := structExport(model.py, keywordMapWSD, model.wsdScanCaps)
		wsd, err = formatPython(obj)
		if err != nil {
			return
		}
	}

	// Expand callback
	expand := func(name string) string {
		switch name {
		case "IPP":
			return ipp
		case "ESCL":
			return escl
		case "WSD":
			return wsd
		}

		return ""
	}

	// Split template into lines. Trim terminating empty lines, if any.
	template := strings.Split(embedPyModel, "\n")
	if legacyMode(model.py) {
		template = strings.Split(embedPyLegacyModel, "\n")
	}

	for len(template) > 0 && template[len(template)-1] == "" {
		template = template[:len(template)-1]
	}

	// Expand template
	skip := true
	for _, t := range template {
		switch {
		case strings.HasPrefix(t, "#-ipp"):
			skip = model.ippPrinterAttrs == nil
		case strings.HasPrefix(t, "#-escl"):
			skip = model.esclScanCaps == nil
		case strings.HasPrefix(t, "#-wsd"):
			skip = model.wsdScanCaps == nil
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

	err = model.ippLoad()
	if err != nil {
		return err
	}

	err = model.esclLoad()
	if err != nil {
		return err
	}

	err = model.wsdLoad()
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
