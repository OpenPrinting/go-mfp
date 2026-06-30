// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Legacy Python->Go converters

package modeling

import (
	"fmt"
	"reflect"
	"time"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/util/uuid"
	"github.com/OpenPrinting/goipp"
)

// legacyStructExport converts the protocol object, represented as Go
// structure or pointer to structure, into the Python dictionary.
//
// kwmap used to map Go struct field names into the
// resulting dictionary key
//
// s MUST be struct or pointer to struct.
func legacyStructExport(py *cpython.Python,
	kwmap map[string]string, s any) *cpython.Object {

	// Create output cpython.Object (the empty dict).
	dict := py.NewObject(map[any]any(nil))

	// Normalize input parameter and obtain the reflect.Value for it.
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct {
		v = v.Elem()
	}
	assert.Must((v.Kind() == reflect.Struct))

	// Roll over all struct fields
	flgs := reflect.VisibleFields(v.Type())
	for _, fld := range flgs {
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

		// Convert into the Python Object and add to the dict,
		item := legacyStructExportValue(py, kwmap, f)
		err := dict.SetItem(keywordNormalize(kwmap, fld.Name), item)

		if err != nil {
			return py.NewError(err)
		}
	}

	return dict
}

// legacyStructExportSlice exports slice of values as the Python object.
func legacyStructExportSlice(py *cpython.Python,
	kwmap map[string]string, v reflect.Value) *cpython.Object {

	list := make([]*cpython.Object, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		list[i] = legacyStructExportValue(py, kwmap, elem)
	}

	return py.NewObject(list)
}

// legacyStructExportValue exports a value as the Python object.
func legacyStructExportValue(py *cpython.Python,
	kwmap map[string]string, v reflect.Value) *cpython.Object {

	// Unwrap wrapped values where possible.
	if wrapper, ok := v.Interface().(wsscan.Wrapper); ok {
		v = reflect.ValueOf(wrapper.Unwrap())
	}

	// Handle known types
	data := v.Interface()
	switch v := data.(type) {
	case escl.Version:
		return py.NewObject(v.String())
	case uuid.UUID:
		return py.Get("UUID").Call(v.String())

	// fmt.Stringer becomes Python string
	case fmt.Stringer:
		return py.NewObject(v.String())
	}

	// Switch by reflect.Kind
	switch v.Kind() {
	case reflect.Struct:
		return legacyStructExport(py, kwmap, data)

	case reflect.Slice:
		return legacyStructExportSlice(py, kwmap, v)
	}

	// Let Python handle default case
	return py.NewObject(data)
}

// legacyStructImport converts the Python object into the Go structure,
// that expected to be the protocol object.
//
// kwmap used to map Go struct field names into the
// resulting dictionary key
//
// p MUST be pointer to struct or pointer to pointer to struct.
func legacyStructImport(obj *cpython.Object,
	kwmap map[string]string, p any) error {

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

	// Import the object
	if !obj.IsDict() {
		// If object is not dictionary, try to interpret it
		// as wsscan.Wrapper without options
		if wrapper, ok := v.Interface().(wsscan.Wrapper); ok {
			// Create the new value of the Wrapper's underlying
			// type
			t2 := reflect.TypeOf(wrapper.Unwrap())
			v2 := reflect.New(t2).Elem()

			// Import its value from Python
			err := legacyStructImportValue(obj, kwmap, v2)
			if err != nil {
				return err
			}

			// Wrap the value
			wrapped := wrapper.Wrap(v2.Interface())
			if wrapped == nil {
				return errPy2Go(obj, v)
			}

			// Replace v with the wrapped value
			v = reflect.ValueOf(wrapped)
		}
	} else {
		// Import structure, field by field
		for _, fld := range reflect.VisibleFields(t) {
			// Lookup python dictionary
			kw := keywordNormalize(kwmap, fld.Name)
			item := obj.GetItem(kw)

			if err := item.Err(); err != nil {
				if item.NotFound() {
					continue
				}
				return errImportWrap(fld.Name, err)
			}

			// Decode the item, if found
			fldval := v.FieldByIndex(fld.Index)
			err := legacyStructImportValue(item, kwmap, fldval)
			if err != nil {
				return errImportWrap(fld.Name, err)
			}
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

// legacyStructImportSlice imports slice of values from the Python object.
func legacyStructImportSlice(obj *cpython.Object,
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
		err = legacyStructImportValue(item, kwmap, v.Index(i))
		if err != nil {
			return errImportWrap(fmt.Sprintf("[%d]", i), err)
		}
	}

	return nil
}

// legacyStructImportValue imports a value from the Python object.
//
// It calls legacyStructImportValueInt, then post-processes the
// returned error, if any.
func legacyStructImportValue(obj *cpython.Object,
	kwmap map[string]string, v reflect.Value) error {

	err := legacyStructImportValueInt(obj, kwmap, v)
	if _, ok := err.(cpython.ErrTypeConversion); ok {
		err = errPy2Go(obj, v)
	}

	return err
}

// legacyStructImportValueInt is the internal function behind the legacyStructImportValue.
func legacyStructImportValueInt(obj *cpython.Object,
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

	// Switch by reflect.Kind
	switch v.Kind() {
	case reflect.Struct:
		return legacyStructImport(obj, kwmap, v.Addr().Interface())

	case reflect.Slice:
		return legacyStructImportSlice(obj, kwmap, v)

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

// legacyIPPExport converts the [ipp.Object] into the [cpython.Object].
func legacyIPPExport(py *cpython.Python, s ipp.Object) *cpython.Object {
	return legacyIPPExportAttrs(py, s.RawAttrs().All())
}

// legacyIPPExportAttrs exports IPP attributes into the [cpython.Object].
func legacyIPPExportAttrs(py *cpython.Python,
	attrs goipp.Attributes) *cpython.Object {

	// Create output cpython.Object (the empty dict).
	dict := py.NewObject(map[any]any(nil))

	// Roll over all IPP attributes
	for _, attr := range attrs {
		vals := legacyIPPExportValues(py, attr)
		err := dict.SetItem(attr.Name, vals)
		if err != nil {
			return py.NewError(err)
		}
	}

	return dict
}

// legacyIPPExportValues exports IPP attribute values into the [cpython.Object].
func legacyIPPExportValues(py *cpython.Python,
	attr goipp.Attribute) *cpython.Object {

	objs := make([]*cpython.Object, 0, len(attr.Values))
	for _, v := range attr.Values {
		obj := legacyIPPExportValue(py, attr.Name, v.T, v.V)
		objs = append(objs, obj)
	}

	if len(objs) == 1 {
		return objs[0]
	}

	return py.NewObject(objs)
}

// legacyIPPExportValue exports IPP value as [cpython.Object].
func legacyIPPExportValue(py *cpython.Python,
	attrname string, tag goipp.Tag, val goipp.Value) *cpython.Object {

	// Collections handled the special way
	if v, ok := val.(goipp.Collection); ok {
		return legacyIPPExportAttrs(py, goipp.Attributes(v))
	}

	// Some Enums are handled the special way
	if tag == goipp.TagEnum {
		switch attrname {
		case "operations-supported":
			op := goipp.Op(val.(goipp.Integer))
			obj := py.Eval(fmt.Sprintf("ipp.OP(0x%.2x)", int(op)))
			if obj.Err() == nil {
				return obj
			}

			// If we got an error here, just continue and
			// handle the value as the regular Enum
		}
	}

	// We represent IPP tag+value at the Python side by wrapping
	// value into the tag-specific Python type:
	//   ipp.ENUM(5)
	//   ipp.KEYWORD('auto')
	//
	// Here we obtain Python type name for the IPP tag
	pytypename := ippTagName[tag]
	if pytypename == "" {
		err := fmt.Errorf("invalid IPP tag %d", int(tag))
		return py.NewError(err)
	}

	// Obtain constructor
	pytype := py.Eval(pytypename)

	// Encode the value
	switch v := val.(type) {
	case goipp.Void:
		return pytype.Call()
	case goipp.Integer:
		return pytype.Call(v)
	case goipp.Boolean:
		return pytype.Call(bool(v))
	case goipp.String:
		return pytype.Call(v)
	case goipp.Time:
		return pytype.Call(v.Format(time.RFC3339))
	case goipp.Resolution:
		return pytype.Call(v.Xres, v.Yres, v.Units.String())
	case goipp.Range:
		return pytype.Call(v.Lower, v.Upper)
	case goipp.TextWithLang:
		return pytype.Call(v.Text, v.Lang)
	case goipp.Binary:
		return pytype.Call(string(v))
	}

	return py.None()
}

// legacyIPPImportAttrs imports IPP attributes from the [cpython.Object].
func legacyIPPImportAttrs(obj *cpython.Object) (
	attrs goipp.Attributes, err error) {

	// Retrieve dictionary keys
	var keyobjs []*cpython.Object
	keyobjs, err = obj.Keys()
	if err != nil {
		return
	}

	for i := range keyobjs {
		var key string
		var valobj *cpython.Object

		// Obtain key name and value
		key, err = keyobjs[i].Str()
		if err == nil {
			valobj = obj.GetItem(keyobjs[i])
			err = valobj.Err()
		}

		if err != nil {
			return
		}

		// Decode the value
		var vals goipp.Values
		vals, err = legacyIPPImportValues(valobj)
		if err != nil {
			return nil, errImportWrap(key, err)
		}

		// Append the attribute
		attrs.Add(goipp.Attribute{Name: key, Values: vals})
	}

	return
}

// legacyIPPImportValues imports IPP values from the [cpython.Object].
func legacyIPPImportValues(obj *cpython.Object) (
	goipp.Values, error) {

	// If obj is the list object, expand it
	var objs []*cpython.Object

	if obj.IsSeq() {
		sz, err := obj.Len()
		if err != nil {
			return nil, err
		}

		objs = make([]*cpython.Object, sz)
		for i := 0; i < sz; i++ {
			objs[i] = obj.GetItem(i)
		}
	} else {
		objs = []*cpython.Object{obj}
	}

	// Now decode each value
	vals := make(goipp.Values, len(objs))
	for i := 0; i < len(objs); i++ {
		tag, val, err := legacyIPPImportValue(objs[i])
		if err != nil {
			return nil, err
		}

		vals[i].T = tag
		vals[i].V = val
	}

	return vals, nil
}

// legacyIPPImportValue imports a single IPP value from the Python object
func legacyIPPImportValue(obj *cpython.Object) (
	tag goipp.Tag, val goipp.Value, err error) {

	if obj.TypeModuleName() == "ipp" {
		typename := obj.TypeName()

		tag = pyIPPTagByName[typename]
		if tag == goipp.TagZero {
			switch typename {
			case "ipp.OP":
				tag = goipp.TagEnum
			}
		}

		switch tag.Type() {
		case goipp.TypeVoid:
			val = goipp.Void{}
		case goipp.TypeInteger:
			var data int64
			data, err = obj.Int()
			val = goipp.Integer(data)
		case goipp.TypeBoolean:
			var data bool
			data, err = obj.Bool()
			val = goipp.Boolean(data)
		case goipp.TypeString, goipp.TypeBinary:
			var data string
			data, err = obj.Str()
			val = goipp.String(data)
		case goipp.TypeDateTime:
			var data string
			data, err = obj.Str()
			if err != nil {
				return
			}

			var t time.Time
			t, err = time.Parse(time.RFC3339, data)
			if err != nil {
				return
			}

			val = goipp.Time{Time: t}
		case goipp.TypeResolution:
			val, err = ippImportIPPResolution(obj)
		case goipp.TypeRange:
			val, err = ippImportIPPRange(obj)
		case goipp.TypeTextWithLang:
			val, err = ippImportIPPTextWithLang(obj, tag)
		default:
			err = fmt.Errorf("ipp.%s: unknown tag type", tag)
		}

		return
	}

	if obj.IsDict() {
		var attrs goipp.Attributes
		attrs, err = legacyIPPImportAttrs(obj)
		val = goipp.Collection(attrs)
		tag = goipp.TagBeginCollection
		return
	}

	err = fmt.Errorf("%s cannot be converted to IPP value", obj.TypeName())
	return
}

// legacyMode returns true of model generation works in legacy mode.
//
// This is temporary function that will be removed after completion
// of the legacy->new format transition.
func legacyMode(py *cpython.Python) bool {
	legacy := py.Get("__use_legacy_format")
	return legacy.Err() == nil
}
