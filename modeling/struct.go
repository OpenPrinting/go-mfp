// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Conversion between Go and Python protocol structures

package modeling

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

var (
	// Reflection package paths to escl and wsscan modules
	pkgPathESCL = reflect.TypeOf(escl.ColorMode(0)).PkgPath()
	pkgPathWSD  = reflect.TypeOf(wsscan.ColorEntry(0)).PkgPath()
)

// structExport converts the protocol object, represented as Go
// structure or pointer to structure, into the Python dictionary.
//
// kwmap used to map Go struct field names into the
// resulting dictionary key
//
// s MUST be struct or pointer to struct.
func structExport(py *cpython.Python,
	kwmap map[string]string, s any) *cpython.Object {

	if legacyMode(py) {
		return legacyStructExport(py, kwmap, s)
	}

	return structExportInt(py, kwmap, s)
}

// structExportInt is the internal function behind the structExport.
func structExportInt(py *cpython.Python,
	kwmap map[string]string, s any) *cpython.Object {

	// Normalize input parameter and obtain the reflect.Value for it.
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct {
		v = v.Elem()
	}
	assert.Must((v.Kind() == reflect.Struct))

	// Roll over all struct fields
	flds := reflect.VisibleFields(v.Type())
	kwargs := make([]cpython.KWArg, 0, len(flds))
	for _, fld := range flds {
		// Skip non-exposed fields
		if !fld.IsExported() {
			continue
		}

		// Obtain and normalize field value
		f := v.FieldByName(fld.Name)
		switch f.Kind() {
		case reflect.Slice:
			// Skip nil slices
			if f.IsNil() {
				continue
			}
		case reflect.Pointer:
			// Skip nil pointers. Dereference others.
			if f.IsNil() {
				continue
			}
			f = f.Elem()
		}

		// Convert into the Python Object and add to the kw map
		item := structExportValue(py, kwmap, f)
		name := keywordNormalize(kwmap, fld.Name)
		kwargs = append(kwargs, cpython.KWArg{Name: name, Value: item})
	}

	// Compute Python-side type name
	name := v.Type().String()
	if i := strings.IndexByte(name, '.'); i >= 0 {
		prefix := name[:i]
		if prefix == "wsscan" {
			prefix = "wsd"
		}

		name = prefix + "." + keywordNormalize(kwmap, name[i+1:])
	} else {
		name = keywordNormalize(kwmap, name)
	}

	return py.Eval(name).CallKWArgs(kwargs)
}

// structExportSlice exports slice of values as the Python object.
func structExportSlice(py *cpython.Python,
	kwmap map[string]string, v reflect.Value) *cpython.Object {

	list := make([]*cpython.Object, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		list[i] = structExportValue(py, kwmap, elem)
	}

	return py.NewObject(list)
}

// structExportValue exports a value as the Python object.
func structExportValue(py *cpython.Python,
	kwmap map[string]string, v reflect.Value) *cpython.Object {

	// wsscan.TextWithLangList with the single element exported as
	// a single wsscan.TextWithLangElement
	if v2, ok := v.Interface().(wsscan.TextWithLangList); ok && len(v2) == 1 {
		v = reflect.ValueOf(v2[0])
	}

	// Handle known types
	data := v.Interface()
	switch v := data.(type) {
	case escl.Version:
		return py.NewObject(v.String())
	case uuid.UUID:
		return py.Get("UUID").Call(v.String())

	case wsscan.TextWithLangElement:
		if v.Lang == nil {
			return py.NewObject(v.Text)
		}

		kwargs := []cpython.KWArg{{Name: "lang", Value: *v.Lang}}
		return py.Eval("wsd.WithLang").CallKWArgs(kwargs, v.Text)

	case wsscan.WithOptionsGetter:
		kwargs := []cpython.KWArg{}
		if opt := v.GetMustHonor(); opt != nil {
			arg := cpython.KWArg{Name: "MustHonor", Value: *opt}
			kwargs = append(kwargs, arg)
		}
		if opt := v.GetOverride(); opt != nil {
			arg := cpython.KWArg{Name: "Override", Value: *opt}
			kwargs = append(kwargs, arg)
		}
		if opt := v.GetUsedDefault(); opt != nil {
			arg := cpython.KWArg{Name: "UsedDefault", Value: *opt}
			kwargs = append(kwargs, arg)
		}

		obj := structExportValue(py, kwmap,
			reflect.ValueOf(v.GetValue()))
		if len(kwargs) == 0 {
			return obj
		}

		return py.Eval("wsd.WithOptions").CallKWArgs(kwargs, obj)

	// fmt.Stringer becomes Python string
	case fmt.Stringer:
		// Use keyword, provided by module, if available
		module := ""
		switch reflect.TypeOf(v).PkgPath() {
		case pkgPathESCL:
			module = "escl"
		case pkgPathWSD:
			module = "wsd"
		}

		s := v.String()
		obj := py.Eval(module + "." + s)
		if py.Eval("helpers.iskeyword").Call(obj).IsTrue() {
			return obj
		}

		return py.NewObject(v.String())
	}

	// Switch by reflect.Kind
	switch v.Kind() {
	case reflect.Struct:
		return structExportInt(py, kwmap, data)

	case reflect.Slice:
		return structExportSlice(py, kwmap, v)
	}

	// Let Python handle default case
	return py.NewObject(data)
}

// structImport converts the Python object into the Go structure,
// that expected to be the protocol object.
//
// kwmap used to map Go struct field names into the
// resulting dictionary key
//
// p MUST be pointer to struct or pointer to pointer to struct.
func structImport(obj *cpython.Object, kwmap map[string]string, p any) error {
	if obj.IsDict() {
		return legacyStructImport(obj, kwmap, p)
	}

	err := structImportInt(obj, kwmap, p)
	if err != nil {
		name := reflect.TypeOf(p).Elem().String()
		if i := strings.IndexByte(name, '.'); i >= 0 {
			name = name[i+1:]
		}
		return errImportWrap(name, err)
	}

	return nil
}

// structImportInt is the internal function behind structImport.
func structImportInt(obj *cpython.Object, kwmap map[string]string, p any) error {
	// Validate argument
	t := reflect.TypeOf(p)

	msg := fmt.Sprintf("%s: invalid type", t)
	assert.MustMsg(t.Kind() == reflect.Pointer, msg)
	assert.MustMsg(p != nil, "nil pointer dereference")

	t = t.Elem()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	assert.MustMsg(t.Kind() == reflect.Struct, msg)

	// Create a new instance of the target structure
	v := reflect.New(t).Elem()

	// Import structure, field by field
	for _, fld := range reflect.VisibleFields(t) {
		// Lookup Python attribute
		kw := keywordNormalize(kwmap, fld.Name)
		item := obj.Get(kw)

		if err := item.Err(); err != nil {
			if item.NotFound() {
				continue
			}
			return errImportWrap(fld.Name, err)
		}

		// Decode the item, if found
		fldval := v.FieldByIndex(fld.Index)
		err := structImportValue(item, kwmap, fldval)
		if err != nil {
			return errImportWrap(fld.Name, err)
		}
	}

	// Save output
	out := reflect.ValueOf(p).Elem()
	if out.Type().Kind() == reflect.Pointer {
		out.Set(v.Addr())
	} else {
		out.Set(v)
	}

	return nil
}

// structImportSlice imports slice of values from the Python object.
func structImportSlice(obj *cpython.Object,
	kwmap map[string]string, v reflect.Value) error {

	// Obtain Python object items
	slice, err := obj.Slice()
	if err != nil {
		return err
	}

	// Allocate output memory
	v.Set(reflect.MakeSlice(v.Type(), len(slice), len(slice)))

	// Decode item by item
	for i, item := range slice {
		err = structImportValue(item, kwmap, v.Index(i))
		if err != nil {
			return errImportWrap(fmt.Sprintf("[%d]", i), err)
		}
	}

	return nil
}

// structImportValue imports a value from the Python object.
//
// It calls structImportValueInt, then post-processes the
// returned error, if any.
func structImportValue(obj *cpython.Object,
	kwmap map[string]string, v reflect.Value) error {

	err := structImportValueInt(obj, kwmap, v)
	if _, ok := err.(cpython.ErrTypeConversion); ok {
		err = errPy2Go(obj, v)
	}

	return err
}

// structImportValueInt is the internal function behind the structImportValue.
func structImportValueInt(obj *cpython.Object,
	kwmap map[string]string, v reflect.Value) error {

	// If we are decoding pointer to value, create a new
	// value instance and shift to it.
	if v.Kind() == reflect.Pointer {
		v2 := reflect.New(v.Type().Elem())
		v.Set(v2)
		v = v2.Elem()
	}

	// Handle known types
	switch v.Interface().(type) {

	// escl types
	case escl.ADFOption:
		return structDecodeEnum(obj, v, escl.DecodeADFOption)
	case escl.ADFState:
		return structDecodeEnum(obj, v, escl.DecodeADFState)
	case escl.BinaryRendering:
		return structDecodeEnum(obj, v, escl.DecodeBinaryRendering)
	case escl.CCDChannel:
		return structDecodeEnum(obj, v, escl.DecodeCCDChannel)
	case escl.ColorMode:
		return structDecodeEnum(obj, v, escl.DecodeColorMode)
	case escl.ColorSpace:
		return structDecodeEnum(obj, v, escl.DecodeColorSpace)
	case escl.ContentType:
		return structDecodeEnum(obj, v, escl.DecodeContentType)
	case escl.FeedDirection:
		return structDecodeEnum(obj, v, escl.DecodeFeedDirection)
	case escl.ImagePosition:
		return structDecodeEnum(obj, v, escl.DecodeImagePosition)
	case escl.InputSource:
		return structDecodeEnum(obj, v, escl.DecodeInputSource)
	case escl.Intent:
		return structDecodeEnum(obj, v, escl.DecodeIntent)
	case escl.JobState:
		return structDecodeEnum(obj, v, escl.DecodeJobState)
	case escl.Units:
		return structDecodeEnum(obj, v, escl.DecodeUnits)

	case escl.JobStateReason:
		rsn, err := esclDecodeJobStateReason(obj)
		if err == nil {
			v.Set(reflect.ValueOf(rsn))
		}
		return err

	case escl.Version:
		ver, err := esclDecodeVersion(obj)
		if err == nil {
			v.Set(reflect.ValueOf(ver))
		}
		return err

	// wsscan types
	case wsscan.ColorEntry:
		return structDecodeEnum(obj, v, wsscan.DecodeColorEntry)
	case wsscan.ContentTypeValue:
		return structDecodeEnum(obj, v, wsscan.DecodeContentTypeValue)
	case wsscan.FilmScanMode:
		return structDecodeEnum(obj, v, wsscan.DecodeFilmScanMode)
	case wsscan.InputSourceValue:
		return structDecodeEnum(obj, v, wsscan.DecodeInputSourceValue)
	case wsscan.JobElemName:
		return structDecodeEnum(obj, v, wsscan.DecodeJobElemName)
	case wsscan.JobStateReason:
		return structDecodeEnum(obj, v, wsscan.DecodeJobStateReason)
	case wsscan.JobState:
		return structDecodeEnum(obj, v, wsscan.DecodeJobState)
	case wsscan.RotationValue:
		return structDecodeEnum(obj, v, wsscan.DecodeRotationValue)
	case wsscan.ScannerElemName:
		return structDecodeEnum(obj, v, wsscan.DecodeScannerElemName)
	case wsscan.ScannerStateReason:
		return structDecodeEnum(obj, v, wsscan.DecodeScannerStateReason)
	case wsscan.ScannerState:
		return structDecodeEnum(obj, v, wsscan.DecodeScannerState)
	case wsscan.Severity:
		return structDecodeEnum(obj, v, wsscan.DecodeSeverity)

	case wsscan.TextWithLangElement:
		return structDecodeTextWithLangElement(obj, v)

	case wsscan.TextWithLangList:
		return structDecodeTextWithLangList(obj, v)

	// other types
	case uuid.UUID:
		s, err := obj.Str()
		if err != nil {
			return err
		}

		u, err := uuid.Parse(s)
		if err == nil {
			v.Set(reflect.ValueOf(u))
		}

		return err
	}

	// Handle interface types with pointer receiver
	switch p := v.Addr().Interface().(type) {
	case wsscan.WithOptions:
		return structDecodeValWithOptions(obj, p)
	}

	// Switch by reflect.Kind
	switch v.Kind() {
	case reflect.Struct:
		return structImportInt(obj, kwmap, v.Addr().Interface())

	case reflect.Slice:
		return structImportSlice(obj, kwmap, v)

	case reflect.Int:
		i, err := obj.Int()
		if err == nil {
			v.Set(reflect.ValueOf(int(i)).Convert(v.Type()))
		}
		return err

	case reflect.String:
		s, err := obj.Str()
		if err == nil {
			v.Set(reflect.ValueOf(s).Convert(v.Type()))
		}
		return err
	}

	return nil
}

// structDecodeEnum decodes enum-alike value from the Python str object,
// using the supplied parse function.
//
// The parse function assumed to return the zero value of the
// target type if string cannot be decoded.
func structDecodeEnum[T comparable](obj *cpython.Object,
	v reflect.Value, parse func(string) T) error {

	var zero T

	s, err := obj.Str()
	if err != nil {
		return err
	}

	val := parse(s)
	if val == zero {
		err := fmt.Errorf("%s: invalid %s", s, reflect.TypeOf(zero))
		return err
	}

	v.Set(reflect.ValueOf(val))
	return nil
}

// structDecodeTextWithLangElement decodes wsscan.TextWithLangElement value
// from the Python object.
//
// Python Object can be of the 'str' or 'wsd.WithLang' type.
func structDecodeTextWithLangElement(obj *cpython.Object, v reflect.Value) error {
	text, err := obj.Unicode()
	if err != nil {
		return err
	}

	langobj := obj.Get("lang")
	var lang optional.Val[string]

	switch {
	case langobj.NotFound():
	case langobj.Err() != nil:
		return langobj.Err()
	case langobj.IsNone():
	default:
		s, err := langobj.Unicode()
		if err != nil {
			return err
		}

		lang = optional.New(s)
	}

	v.Set(reflect.ValueOf(
		wsscan.TextWithLangElement{Text: text, Lang: lang},
	))

	return nil
}

// structDecodeTextWithLangList decodes wsscan.TextWithLangList value from the
// Python object.
func structDecodeTextWithLangList(obj *cpython.Object, v reflect.Value) error {
	switch obj.TypeName() {
	case "str", "wsd.WithLang":
		// wsscan.TextWithLangList containing a single
		// element can be represented by a single value
		// of the str or wsd.WithLang type (not by list
		// of values).
		//
		// Handle it here.
		out := wsscan.TextWithLangElement{}
		err := structDecodeTextWithLangElement(obj,
			reflect.ValueOf(&out).Elem())
		if err != nil {
			return err
		}

		v.Set(reflect.ValueOf(wsscan.TextWithLangList{out}))

		return nil
	}

	return structImportSlice(obj, keywordMapWSD, v)
}

// structDecodeEnum decodes wsscan.ValWithOptions value from the
// Python object.
func structDecodeValWithOptions(obj *cpython.Object, v wsscan.WithOptions) error {
	// Obtain reflect.Type and reflect.Value for underlying value
	t2 := reflect.TypeOf(v.GetValue())
	v2 := reflect.New(t2).Elem()

	// Import its value from Python
	err := structImportValue(obj, keywordMapWSD, v2)
	if err != nil {
		return err
	}

	// Save the value
	if !v.SetValue(v2.Interface()) {
		return errPy2Go(obj, reflect.ValueOf(v))
	}

	// Import options
	type option struct {
		name string
		set  func(optional.Val[wsscan.BooleanElement])
	}

	options := []option{
		{name: "MustHonor", set: v.SetMustHonor},
		{name: "Override", set: v.SetOverride},
		{name: "UsedDefault", set: v.SetUsedDefault},
	}

	for _, opt := range options {
		optobj := obj.Get(opt.name)

		switch {
		case optobj.NotFound():
		case optobj.Err() != nil:
			err = optobj.Err()
		case optobj.IsNone():
		default:
			var s string
			s, err = optobj.Unicode()
			if err == nil {
				opt.set(optional.New(wsscan.BooleanElement(s)))
			}
		}

		if err != nil {
			return errImportWrap(opt.name, err)
		}
	}

	return nil
}
