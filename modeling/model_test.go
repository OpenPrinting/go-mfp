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
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/proto/usb"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
	"github.com/OpenPrinting/goipp"
)

// TestKyoceraIPPPrinterAttributes is the real-world test, that
// verifies that the real Kyocera ECOSYS M2040dn IPP Printer Attributes
// is properly handled.
func TestKyoceraIPPPrinterAttributes(t *testing.T) {
	// Decode Kyocera PrinterAttributes
	var msg goipp.Message
	err := msg.DecodeBytes(testutils.Kyocera.ECOSYS.M2040dn.
		IPP.PrinterAttributes)
	assert.NoError(err)

	pa, err := ipp.DecodePrinterAttributes(msg.Printer, nil)
	assert.NoError(err)

	// Create a new, empty Model
	model, err := NewModel()
	assert.NoError(err)

	defer model.Close()

	// Roll over ippExport/ippImportPrinterAppributes
	obj := ippExport(model.py, pa)
	if err := obj.Err(); err != nil {
		t.Fatalf("ippExport: %s", err)
		return
	}

	pa2, err := ippImportPrinterAppributes(obj)
	if err != nil {
		t.Fatalf("ippImportPrinterAppributes: %s", err)
		return
	}

	attrs := pa.RawAttrs().All()
	attrs2 := pa2.RawAttrs().All()
	if !attrs.Equal(attrs2) {
		diff := testutils.IPPDiffAttributes("expected", attrs, "present", attrs2)
		t.Errorf("ippExport/ippImportPrinterAppributes:\n%s", diff)
	}

	// Roll over Model.Write/Model.Read
	buf := &bytes.Buffer{}

	model.SetIPPPrinterAttrs(pa)
	err = model.Write(buf)
	if err != nil {
		t.Fatalf("Model.Write: %s", err)
	}

	model2, err := NewModel()
	assert.NoError(err)

	defer model2.Close()

	err = model2.Read("test", buf)
	if err != nil {
		t.Fatalf("Model.Read: %s", err)
	}

	pa2 = model2.GetIPPPrinterAttrs()
	if pa2 == nil {
		t.Fatalf("Model.Read: missed IPP printer attributes")
	}

	attrs = pa.RawAttrs().All()
	attrs2 = pa2.RawAttrs().All()
	if !attrs.Equal(attrs2) {
		diff := testutils.IPPDiffAttributes("expected", attrs, "present", attrs2)
		t.Errorf("Model.Write/Model.Read:\n%s", diff)
	}
}

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
		t.Fatalf("structExport: %s", err)
		return
	}

	var scancaps2 *escl.ScannerCapabilities
	err = structImport(obj, keywordMapESCL, &scancaps2)
	if err != nil {
		t.Fatalf("structImport: %s", err)
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
		t.Fatalf("Model.Write: %s", err)
	}

	model2, err := NewModel()
	assert.NoError(err)

	defer model2.Close()

	err = model2.Read("test", buf)
	if err != nil {
		t.Fatalf("Model.Read: %s", err)
	}

	scancaps2 = model2.GetESCLScanCaps()
	if scancaps2 == nil {
		t.Fatalf("Model.Read: missed eSCL scanner capabilities")
	}

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
		t.Fatalf("structExport: %s", err)
		return
	}

	var scancaps2 *wsscan.GetScannerElementsResponse
	err = structImport(obj, keywordMapWSD, &scancaps2)
	if err != nil {
		t.Fatalf("structImport: %s", err)
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
		t.Fatalf("Model.Write: %s", err)
	}

	model2, err := NewModel()
	assert.NoError(err)

	defer model2.Close()

	err = model2.Read("test", buf)
	if err != nil {
		t.Fatalf("Model.Read: %s", err)
	}

	scancaps2 = model2.GetWSDScanCaps()
	if scancaps2 == nil {
		t.Fatalf("Model.Read: missed WSD scanner capabilities")
	}

	diff = testutils.Diff(scancaps, scancaps2)
	if diff != "" {
		t.Errorf("Model.Write/Model.Read:\n%s", diff)
	}
}

// TestKyoceraUSBDeviceDescriptor is the real-world test, that
// verifies that the real Kyocera ECOSYS M2040dn USB Device Descriptor
// is properly handled.
func TestKyoceraUSBDeviceDescriptor(t *testing.T) {
	// FIXME -- disabled for now
	return

	// Obtain USB device descriptor
	desc := &testutils.Kyocera.ECOSYS.M2040dn.USB.DeviceDescriptor

	// Create a new, empty Model
	model, err := NewModel()
	assert.NoError(err)

	defer model.Close()

	// Roll over structExport/structImport
	obj := structExport(model.py, keywordMapUSB, desc)
	if err := obj.Err(); err != nil {
		t.Fatalf("structExport: %s", err)
		return
	}

	var desc2 *usb.DeviceDescriptor
	err = structImport(obj, keywordMapUSB, &desc2)
	if err != nil {
		t.Fatalf("structImport: %s", err)
		return
	}

	diff := testutils.Diff(desc, desc2)
	if diff != "" {
		t.Errorf("structExport/structImport:\n%s", diff)
	}

	// Roll over Model.Write/Model.Read
	buf := &bytes.Buffer{}

	model.SetUSBDeviceDescriptor(desc)
	err = model.Write(buf)
	if err != nil {
		t.Fatalf("Model.Write: %s", err)
	}

	model2, err := NewModel()
	assert.NoError(err)

	defer model2.Close()

	err = model2.Read("test", buf)
	if err != nil {
		t.Fatalf("Model.Read: %s", err)
	}

	desc2 = model2.GetUSBDeviceDescriptor()
	if desc2 == nil {
		t.Fatalf("Model.Read: missed USB Device Descriptor")
	}

	diff = testutils.Diff(desc, desc2)
	if diff != "" {
		t.Errorf("Model.Write/Model.Read:\n%s", diff)
	}
}

// TestWSDTextWithLang tests Go<->Python export/import conversions
// for structures that contain wsscan.WSDTextWithLangList and
// wsscan.WSDTextWithLangElement fields
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
			t.Errorf("%s: export error\n%s",
				test.name, testutils.Diff(expected, present))
			continue
		}

		// Decode Python->Go
		out := reflect.New(reflect.TypeOf(test.in)).Interface()
		err := structImport(obj, keywordMapWSD, out)
		if err != nil {
			t.Errorf("%s: import error\n%s", test.name, err)
			continue
		}

		diff := testutils.Diff(test.in,
			reflect.ValueOf(out).Elem().Interface())

		if diff != "" {
			t.Errorf("%s: import error\n%s", test.name, diff)
		}
	}

}

// TestWSDTextWithLangDecodeErrors tests Python->Go import errors,
// specific for  wsscan.WSDTextWithLangList and wsscan.WSDTextWithLangElement
func TestWSDTextWithLangDecodeErrors(t *testing.T) {
	model, err := NewModel()
	assert.NoError(err)
	py := model.py
	defer model.Close()

	pyScannerDescription := py.Eval("wsd.ScannerDescription")

	type testData struct {
		name string          // Test name
		out  any             // Output data (empty Go struct)
		obj  *cpython.Object // Expected Python output
		err  string          // Expected error
	}

	tests := []testData{
		testData{
			name: "Single invalid element",
			out:  &wsscan.ScannerDescription{},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": 25,
				},
			),
			err: `ScannerDescription.ScannerInfo: can't convert int to wsscan.TextWithLangList`,
		},

		testData{
			name: "Multiple invalid elements",
			out:  &wsscan.ScannerDescription{},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": []any{1, 2},
				},
			),
			err: `ScannerDescription.ScannerInfo[0]: can't convert int to wsscan.TextWithLangElement`,
		},

		testData{
			name: "Mixed valid/invalid elements",
			out:  &wsscan.ScannerDescription{},
			obj: pyScannerDescription.CallKW(
				map[string]any{
					"ScannerInfo": []any{"OK", 2.5},
				},
			),
			err: `ScannerDescription.ScannerInfo[1]: can't convert float to wsscan.TextWithLangElement`,
		},
	}

	for _, test := range tests {
		err = structImport(test.obj, keywordMapWSD, test.out)

		expected := test.err
		present := ""
		if err != nil {
			present = err.Error()
		}

		if present != expected {
			t.Errorf("%s: error mismatch:\n"+
				"expected: %q\n"+
				"present:  %q\n", test.name, expected, present)
		}
	}
}

// TestWSDValWithOptions Go<->Python export/import conversions
// for structures that contain wsscan.WSDValWithOptins fields
func TestWSDValWithOptions(t *testing.T) {
	model, err := NewModel()
	assert.NoError(err)
	py := model.py
	defer model.Close()

	pyDocumentParameters := py.Eval("wsd.DocumentParameters")
	pyWithOptions := py.Eval("wsd.WithOptions")

	type testData struct {
		name string          // Test name
		in   any             // Input data (Go struct)
		obj  *cpython.Object // Expected Python output
	}

	tests := []testData{
		testData{
			name: "wsscan.ValWithOptions[int], no options",
			in: wsscan.DocumentParameters{
				CompressionQualityFactor: optional.New(
					wsscan.ValWithOptions[int]{Val: 5}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": 5,
				},
			),
		},

		testData{
			name: "wsscan.ValWithOptions[int], MustHonor = true",
			in: wsscan.DocumentParameters{
				CompressionQualityFactor: optional.New(
					wsscan.ValWithOptions[int]{
						Val:       5,
						MustHonor: optional.New(wsscan.BooleanElement("true")),
					}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": pyWithOptions.CallKW(
						map[string]any{
							"MustHonor": "true",
						},
						5,
					),
				},
			),
		},

		testData{
			name: "wsscan.ValWithOptions[int], Override = true",
			in: wsscan.DocumentParameters{
				CompressionQualityFactor: optional.New(
					wsscan.ValWithOptions[int]{
						Val:      5,
						Override: optional.New(wsscan.BooleanElement("true")),
					}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": pyWithOptions.CallKW(
						map[string]any{
							"Override": "true",
						},
						5,
					),
				},
			),
		},

		testData{
			name: "wsscan.ValWithOptions[int], UsedDefault = true",
			in: wsscan.DocumentParameters{
				CompressionQualityFactor: optional.New(
					wsscan.ValWithOptions[int]{
						Val:         5,
						UsedDefault: optional.New(wsscan.BooleanElement("true")),
					}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": pyWithOptions.CallKW(
						map[string]any{
							"UsedDefault": "true",
						},
						5,
					),
				},
			),
		},

		testData{
			name: "wsscan.ValWithOptions[int], MustHonor = 1, Override = 1, UsedDefault = 1",
			in: wsscan.DocumentParameters{
				CompressionQualityFactor: optional.New(
					wsscan.ValWithOptions[int]{
						Val:         5,
						MustHonor:   optional.New(wsscan.BooleanElement("1")),
						Override:    optional.New(wsscan.BooleanElement("1")),
						UsedDefault: optional.New(wsscan.BooleanElement("1")),
					}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": pyWithOptions.CallKW(
						map[string]any{
							"MustHonor":   "1",
							"Override":    "1",
							"UsedDefault": "1",
						},
						5,
					),
				},
			),
		},

		testData{
			name: "wsscan.ValWithOptions[ContentTypeValue], no options",
			in: wsscan.DocumentParameters{
				ContentType: optional.New(
					wsscan.ValWithOptions[wsscan.ContentTypeValue]{
						Val: wsscan.Text,
					}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"ContentType": py.Eval("wsd.Text"),
				},
			),
		},

		testData{
			name: "wsscan.ValWithOptions[ContentTypeValue], MustHonor = True",
			in: wsscan.DocumentParameters{
				ContentType: optional.New(
					wsscan.ValWithOptions[wsscan.ContentTypeValue]{
						Val:       wsscan.Text,
						MustHonor: optional.New(wsscan.BooleanElement("True")),
					}),
			},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"ContentType": pyWithOptions.CallKW(
						map[string]any{
							"MustHonor": "True",
						},
						py.Eval("wsd.Text"),
					),
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
			t.Errorf("%s: export error\n%s",
				test.name,
				testutils.DiffLines(
					"expected", expected,
					"present", present))
			continue
		}

		// Decode Python->Go
		out := reflect.New(reflect.TypeOf(test.in)).Interface()
		err := structImport(obj, keywordMapWSD, out)
		if err != nil {
			t.Errorf("%s: import error\n%s", test.name, err)
			continue
		}

		diff := testutils.Diff(test.in,
			reflect.ValueOf(out).Elem().Interface())

		if diff != "" {
			t.Errorf("%s: import error\n%s", test.name, diff)
		}
	}
}

// TestWSDValWithOptinsDecodeErrors tests Python->Go import errors,
// specific for wsscan.WSDValWithOptins
func TestWSDValWithOptinsDecodeErrors(t *testing.T) {
	model, err := NewModel()
	assert.NoError(err)
	py := model.py
	defer model.Close()

	pyDocumentParameters := py.Eval("wsd.DocumentParameters")
	pyWithOptions := py.Eval("wsd.WithOptions")

	type testData struct {
		name string          // Test name
		out  any             // Output data (empty Go struct)
		obj  *cpython.Object // Expected Python output
		err  string          // Expected error
	}

	tests := []testData{
		testData{
			name: "wsscan.ValWithOptions[int]: mismatching value type",
			out:  &wsscan.DocumentParameters{},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": "xxx",
				},
			),
			err: `DocumentParameters.CompressionQualityFactor: can't convert str to int`,
		},

		testData{
			name: "wsscan.ValWithOptions[int]: mismatching MustHonor type",
			out:  &wsscan.DocumentParameters{},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"CompressionQualityFactor": pyWithOptions.CallKW(
						map[string]any{
							"MustHonor": 12345,
						},
						5,
					),
				},
			),
			err: `DocumentParameters.CompressionQualityFactor.MustHonor: can't convert int to string`,
		},

		testData{
			name: "wsscan.ValWithOptions[ContentTypeValue]: mismatching value type",
			out:  &wsscan.DocumentParameters{},
			obj: pyDocumentParameters.CallKW(
				map[string]any{
					"ContentType": "bad value",
				},
			),
			err: `DocumentParameters.ContentType: bad value: invalid wsscan.ContentTypeValue`,
		},
	}

	for _, test := range tests {
		err = structImport(test.obj, keywordMapWSD, test.out)

		expected := test.err
		present := ""
		if err != nil {
			present = err.Error()
		}

		if present != expected {
			t.Errorf("%s: error mismatch:\n"+
				"expected: %q\n"+
				"present:  %q\n", test.name, expected, present)
		}
	}
}
