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
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/util/uuid"
	"github.com/OpenPrinting/goipp"
)

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
