// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device model test

package modeling

import (
	"bytes"
	"testing"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// TestKyoceraESCLScannerCapabilities is the real-world test, that
// verifies that the real Kyocera ECOSYS M2040dn eSCL ScannerCapabilities
// is properly handled.
func TestKyoceraESCLScannerCapabilities(t *testing.T) {
	// Decode Kyocera ScannerCapabilities
	rd := bytes.NewReader(testutils.Kyocera.
		ECOSYS.M2040dn.ESCL.ScannerCapabilities)
	xml, err := xmldoc.Decode(escl.NsMap, rd)
	assert.NoError(err)

	scancaps, err := escl.DecodeScannerCapabilities(xml)
	assert.NoError(err)

	// Create a new, empty Model
	model, err := NewModel()
	assert.NoError(err)

	defer model.Close()

	// Roll over structExport/structImport
	obj := structExport(model.py, keywordMapESCL, scancaps)
	if err := obj.Err(); err != nil {
		t.Errorf("structExport: %s", err)
		return
	}

	var scancaps2 *escl.ScannerCapabilities
	err = structImport(obj, keywordMapESCL, &scancaps2)
	if err != nil {
		t.Errorf("structImport: %s", err)
		return
	}

	diff := testutils.Diff(scancaps, scancaps2)
	if diff != "" {
		t.Errorf("structExport/structImport:\n%s", diff)
	}

	// Roll over Model.Write/Model.Read
	buf := &bytes.Buffer{}

	model.SetESCLScanCaps(scancaps)
	err = model.Write(buf)
	if err != nil {
		t.Errorf("Model.Write: %s", err)
	}

	model2, err := NewModel()
	assert.NoError(err)

	defer model2.Close()

	err = model2.Read("test", buf)
	if err != nil {
		t.Errorf("Model.Read: %s", err)
	}

	scancaps2 = model2.GetESCLScanCaps()
	diff = testutils.Diff(scancaps, scancaps2)
	if diff != "" {
		t.Errorf("Model.Write/Model.Read:\n%s", diff)
	}
}

// TestKyoceraWSDScannerCapabilities is the real-world test, that
// verifies that the real Kyocera ECOSYS M2040dn WSD ScannerCapabilities
// is properly handled.
func TestKyoceraWSDScannerCapabilities(t *testing.T) {
	// Decode Kyocera ScannerCapabilities
	rd := bytes.NewReader(testutils.Kyocera.
		ECOSYS.M2040dn.WSD.GetScannerElementsResponse)
	xml, err := xmldoc.Decode(wsscan.NsMap, rd)
	assert.NoError(err)

	msg, err := wsscan.DecodeMessage(xml)
	assert.NoError(err)

	scancaps := msg.Body.(*wsscan.GetScannerElementsResponse)

	// Create a new, empty Model
	model, err := NewModel()
	assert.NoError(err)

	defer model.Close()

	// Roll over structExport/structImport
	obj := structExport(model.py, keywordMapWSD, scancaps)
	if err := obj.Err(); err != nil {
		t.Errorf("structExport: %s", err)
		return
	}

	var scancaps2 *wsscan.GetScannerElementsResponse
	err = structImport(obj, keywordMapWSD, &scancaps2)
	if err != nil {
		t.Errorf("structImport: %s", err)
		return
	}

	diff := testutils.Diff(scancaps, scancaps2)
	if diff != "" {
		t.Errorf("structExport/structImport:\n%s", diff)
	}

	// Roll over Model.Write/Model.Read
	buf := &bytes.Buffer{}

	model.SetWSDScanCaps(scancaps)
	err = model.Write(buf)
	if err != nil {
		t.Errorf("Model.Write: %s", err)
	}

	model2, err := NewModel()
	assert.NoError(err)

	defer model2.Close()

	err = model2.Read("test", buf)
	if err != nil {
		t.Errorf("Model.Read: %s", err)
	}

	scancaps2 = model2.GetWSDScanCaps()
	diff = testutils.Diff(scancaps, scancaps2)
	if diff != "" {
		t.Errorf("Model.Write/Model.Read:\n%s", diff)
	}
}

func TestWSDTextWithLang(t *testing.T) {
	model, err := NewModel()
	assert.NoError(err)
	py := model.py
	defer model.Close()

	pyScannerDescription := py.Eval("wsd.ScannerDescription")
	pyWithLang := py.Eval("wsd.WithLang")

	type testData struct {
		name string          // Test name
		in   any             // Input data (Go struct)
		obj  *cpython.Object // Expected Python output
	}

	tests := []testData{
		testData{
			name: "Single TextWithLangList, no language",
			in: wsscan.ScannerDescription{
				ScannerInfo: wsscan.TextWithLangList{
					wsscan.TextWithLangElement{
						Text: "Sample scanner",
					},
				},
			},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": "Sample scanner",
				},
			),
		},

		testData{
			name: "Single TextWithLangList with language",
			in: wsscan.ScannerDescription{
				ScannerInfo: wsscan.TextWithLangList{
					wsscan.TextWithLangElement{
						Text: "Sample scanner",
						Lang: optional.New("en"),
					},
				},
			},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": pyWithLang.CallKW(
						map[string]any{
							"lang": "en",
						},
						"Sample scanner",
					),
				},
			),
		},

		testData{
			name: "Multiple TextWithLangList with language",
			in: wsscan.ScannerDescription{
				ScannerInfo: wsscan.TextWithLangList{
					wsscan.TextWithLangElement{
						Text: "Sample scanner",
						Lang: optional.New("en"),
					},
					wsscan.TextWithLangElement{
						Text: "Простой сканер",
						Lang: optional.New("ru"),
					},
				},
			},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": []any{
						pyWithLang.CallKW(
							map[string]any{
								"lang": "en",
							},
							"Sample scanner",
						),
						pyWithLang.CallKW(
							map[string]any{
								"lang": "ru",
							},
							"Простой сканер",
						),
					},
				},
			),
		},

		testData{
			name: "Multiple TextWithLangList with mixed language",
			in: wsscan.ScannerDescription{
				ScannerInfo: wsscan.TextWithLangList{
					wsscan.TextWithLangElement{
						Text: "Sample scanner",
					},
					wsscan.TextWithLangElement{
						Text: "Простой сканер",
						Lang: optional.New("ru"),
					},
				},
			},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": []any{
						"Sample scanner",
						pyWithLang.CallKW(
							map[string]any{
								"lang": "ru",
							},
							"Простой сканер",
						),
					},
				},
			),
		},
	}

	for _, test := range tests {
		// Encode Go->Python and check result against expectation
		obj := structExport(model.py, keywordMapWSD, test.in)

		expected := test.obj.String()
		present := obj.String()

		if expected != present {
			t.Errorf("%s: encode error\n%s",
				test.name, testutils.Diff(expected, present))
		}
	}

}
