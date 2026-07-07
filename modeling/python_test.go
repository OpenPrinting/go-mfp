// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Tests of Python helper scrips

package modeling

import (
	"testing"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// Test_ipp_py tests ipp.py
func Test_ipp_py(t *testing.T) {
	model, err := NewModel()
	assert.NoError(err)
	defer model.Close()
	py := model.py

	type testData struct {
		expr string
	}

	tests := []testData{
		// No-value tags
		{`ipp.UNSUPPORTED_VALUE()`},
		{`ipp.DEFAULT()`},
		{`ipp.UNKNOWN()`},
		{`ipp.NOVALUE()`},
		{`ipp.NOTSETTABLE()`},
		{`ipp.DELETEATTR()`},
		{`ipp.ADMINDEFINE()`},

		// Simple types
		{`ipp.INTEGER(12345)`},
		{`ipp.BOOLEAN(True)`},
		{`ipp.BOOLEAN(False)`},
		{`ipp.ENUM(12345)`},
		{`ipp.STRING('hello')`},
		{`ipp.TEXT('hello')`},
		{`ipp.NAME('hello')`},
		{`ipp.KEYWORD('hello')`},
		{`ipp.URI('http://example.com')`},
		{`ipp.URISCHEME('ipp')`},
		{`ipp.CHARSET('utf-8')`},
		{`ipp.LANGUAGE('en-US')`},
		{`ipp.MIMETYPE('image/jpeg')`},
		{`ipp.DATE('2026-07-06T16:35:48+03:00')`},

		// Composite types
		{`ipp.RANGE(1, 999)`},
		{`ipp.TEXTLANG('привет', 'ru-RU')`},
		{`ipp.NAMELANG('привет', 'ru-RU')`},

		// Collections
		{`ipp.COLLECTION()`},
		{
			"" +
				"ipp.COLLECTION(\n" +
				"    x = ipp.INTEGER(1),\n" +
				"    y = ipp.INTEGER(2),\n" +
				")",
		},
	}

	for _, test := range tests {
		obj := py.Eval(test.expr)
		if obj.Err() != nil {
			t.Errorf("%s: %s", test.expr, obj)
			continue
		}

		repr, err := obj.Repr()
		if err != nil {
			t.Errorf("%s: repr() error: %s", test.expr, err)
			continue
		}

		if repr != test.expr {
			t.Errorf("%s: repr() mismatch: %s", test.expr, repr)
		}
	}
}
