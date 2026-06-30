// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Conversion between ipp.Object and cpython.Object

package modeling

import (
	"fmt"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/goipp"
)

// ippExport converts the [ipp.Object] into the [cpython.Object].
func ippExport(py *cpython.Python, s ipp.Object) *cpython.Object {
	if legacyMode(py) {
		return legacyIPPExport(py, s)
	}

	return ippExportAttrs(py, s.RawAttrs().All())
}

// ippExportAttrs exports IPP attributes into the [cpython.Object].
func ippExportAttrs(py *cpython.Python,
	attrs goipp.Attributes) *cpython.Object {

	// Roll over all IPP attributes
	kwargs := make([]cpython.KWArg, len(attrs))
	for i, attr := range attrs {
		vals := ippExportValues(py, attr)
		name := ippAttrNameToPython(attr.Name)
		kwargs[i] = cpython.KWArg{Name: name, Value: vals}
	}

	// Create collection object
	return py.Eval("ipp.COLLECTION").CallKWArgs(kwargs)
}

// ippExportValues exports IPP attribute values into the [cpython.Object].
func ippExportValues(py *cpython.Python,
	attr goipp.Attribute) *cpython.Object {

	objs := make([]*cpython.Object, 0, len(attr.Values))
	for _, v := range attr.Values {
		obj := ippExportValue(py, attr.Name, v.T, v.V)
		objs = append(objs, obj)
	}

	if len(objs) == 1 {
		return objs[0]
	}

	return py.NewObject(objs)
}

// ippExportValue exports IPP value as [cpython.Object].
func ippExportValue(py *cpython.Python,
	attrname string, tag goipp.Tag, val goipp.Value) *cpython.Object {

	// Collections handled the special way
	if v, ok := val.(goipp.Collection); ok {
		return ippExportAttrs(py, goipp.Attributes(v))
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

// ippImportPrinterAppributes imports IPP printer attributes from the
// Python representation
func ippImportPrinterAppributes(obj *cpython.Object) (
	*ipp.PrinterAttributes, error) {

	attrs, err := ippImportIPPAttrs(obj)
	if err != nil {
		return nil, err
	}

	opt := &ipp.DecoderOptions{
		KeepTrying: true,
	}

	return ipp.DecodePrinterAttributes(attrs, opt)
}

// ippImportIPPAttrs imports IPP attributes from the [cpython.Object].
func ippImportIPPAttrs(obj *cpython.Object) (
	attrs goipp.Attributes, err error) {

	if obj.IsDict() {
		return legacyIPPImportAttrs(obj)
	}

	return ippImportIPPAttrsInt(obj)
}

// ippImportIPPAttrsInt is the internal function behind ippImportIPPAttrs.
func ippImportIPPAttrsInt(obj *cpython.Object) (
	attrs goipp.Attributes, err error) {

	// Retrieve attribute names as a slice of obj.__dict__ key objects
	var keyobjs []*cpython.Object
	keyobjs, err = obj.Get("__dict__").Keys()
	if err != nil {
		return
	}

	for i := range keyobjs {
		var key string
		var valobj *cpython.Object

		// Obtain key name and value
		key, err = keyobjs[i].Str()
		if err == nil {
			valobj = obj.Get(key)
			err = valobj.Err()
		}

		if err != nil {
			return
		}

		// Decode the value
		var vals goipp.Values
		vals, err = ippImportIPPValues(valobj)
		if err != nil {
			return nil, errImportWrap(key, err)
		}

		// Append the attribute
		name := ippPythonToAttrName(key)
		attrs.Add(goipp.Attribute{Name: name, Values: vals})
	}

	return
}

// ippImportIPPValues imports IPP values from the [cpython.Object].
func ippImportIPPValues(obj *cpython.Object) (
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
		tag, val, err := ippImportIPPValue(objs[i])
		if err != nil {
			return nil, err
		}

		vals[i].T = tag
		vals[i].V = val
	}

	return vals, nil
}

// ippImportIPPValue imports IPP value from the Python object
func ippImportIPPValue(obj *cpython.Object) (
	tag goipp.Tag, val goipp.Value, err error) {

	if obj.TypeName() == "ipp.COLLECTION" {
		var attrs goipp.Attributes
		attrs, err = ippImportIPPAttrsInt(obj)
		val = goipp.Collection(attrs)
		tag = goipp.TagBeginCollection
		return
	}

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

	err = fmt.Errorf("%s cannot be converted to IPP value", obj.TypeName())
	return
}

// ippImportIPPResolution imports IPP resolution from the Python object
func ippImportIPPResolution(obj *cpython.Object) (
	res goipp.Resolution, err error) {

	var x, y int64
	var units goipp.Units

	// Load Xres
	x, err = obj.Get("X").Int()
	if err != nil {
		err = errImportWrap("X", err)
		return
	}

	// Load Yres
	y, err = obj.Get("Y").Int()
	if err != nil {
		err = errImportWrap("Y", err)
		return
	}

	// Load Units
	unitsName, err := obj.Get("Units").Str()
	if err == nil {
		switch unitsName {
		case "dpi":
			units = goipp.UnitsDpi
		case "dpcm":
			units = goipp.UnitsDpcm
		default:
			err = fmt.Errorf("%s: invalid resolution units", unitsName)
		}
	}

	if err != nil {
		return
	}

	res = goipp.Resolution{
		Xres:  int(x),
		Yres:  int(y),
		Units: units,
	}

	return
}

// ippImportIPPRange imports IPP range from the Python object
func ippImportIPPRange(obj *cpython.Object) (
	rng goipp.Range, err error) {

	var lower, upper int64

	// Load Lower
	lower, err = obj.Get("Lower").Int()
	if err != nil {
		err = errImportWrap("Lower", err)
		return
	}

	// Load Upper
	upper, err = obj.Get("Upper").Int()
	if err != nil {
		err = errImportWrap("Upper", err)
		return
	}

	rng = goipp.Range{
		Lower: int(lower),
		Upper: int(upper),
	}

	return
}

// ippImportIPPTextWithLang imports IPP text with language from the Python object
func ippImportIPPTextWithLang(obj *cpython.Object, tag goipp.Tag) (
	txt goipp.TextWithLang, err error) {

	var lang, text string

	// Load lang
	lang, err = obj.Get("Lang").Str()
	if err != nil {
		err = errImportWrap("Lang", err)
		return
	}

	// Load Text or Name
	nm := "Text"
	if tag == goipp.TagNameLang {
		nm = "Name"
	}

	text, err = obj.Get(nm).Str()
	if err != nil {
		err = errImportWrap(nm, err)
		return
	}

	txt = goipp.TextWithLang{
		Lang: lang,
		Text: text,
	}

	return
}

// ippAttrNameToPython translates IPP attribute name to Python identifier
func ippAttrNameToPython(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

// ippAttrNameToPython translates IPP attribute name from Python identifier
func ippPythonToAttrName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}

// ippTagName maps goipp.Tag to its Python name
var ippTagName = map[goipp.Tag]string{
	// Delimiters
	goipp.TagZero: "ipp.ZERO",
	goipp.TagEnd:  "ipp.END",

	// Groups of attributes
	goipp.TagOperationGroup:         "ipp.OPERATION",
	goipp.TagJobGroup:               "ipp.JOB",
	goipp.TagPrinterGroup:           "ipp.PRINTER",
	goipp.TagUnsupportedGroup:       "ipp.UNSUPPORTED_GROUP",
	goipp.TagSubscriptionGroup:      "ipp.SUBSCRIPTION",
	goipp.TagEventNotificationGroup: "ipp.EVENT_NOTIFICATION",
	goipp.TagResourceGroup:          "ipp.RESOURCE",
	goipp.TagDocumentGroup:          "ipp.DOCUMENT",
	goipp.TagSystemGroup:            "ipp.SYSTEM",

	// Special values
	goipp.TagUnsupportedValue: "ipp.UNSUPPORTED_VALUE",
	goipp.TagDefault:          "ipp.DEFAULT",
	goipp.TagUnknown:          "ipp.UNKNOWN",
	goipp.TagNoValue:          "ipp.NOVALUE",
	goipp.TagNotSettable:      "ipp.NOTSETTABLE",
	goipp.TagDeleteAttr:       "ipp.DELETEATTR",
	goipp.TagAdminDefine:      "ipp.ADMINDEFINE",

	// Values
	goipp.TagInteger:    "ipp.INTEGER",
	goipp.TagBoolean:    "ipp.BOOLEAN",
	goipp.TagEnum:       "ipp.ENUM",
	goipp.TagString:     "ipp.STRING",
	goipp.TagDateTime:   "ipp.DATE",
	goipp.TagResolution: "ipp.RESOLUTION",
	goipp.TagRange:      "ipp.RANGE",
	goipp.TagTextLang:   "ipp.TEXTLANG",
	goipp.TagNameLang:   "ipp.NAMELANG",
	goipp.TagText:       "ipp.TEXT",
	goipp.TagName:       "ipp.NAME",
	goipp.TagKeyword:    "ipp.KEYWORD",
	goipp.TagURI:        "ipp.URI",
	goipp.TagURIScheme:  "ipp.URISCHEME",
	goipp.TagCharset:    "ipp.CHARSET",
	goipp.TagLanguage:   "ipp.LANGUAGE",
	goipp.TagMimeType:   "ipp.MIMETYPE",
	goipp.TagExtension:  "ipp.EXTENSION",

	// Collections
	goipp.TagBeginCollection: "ipp.BEGIN_COLLECTION",
	goipp.TagEndCollection:   "ipp.END_COLLECTION",
	goipp.TagMemberName:      "ipp.MEMBERNAME",
}

// pyIPPTagByName maps goipp.Tag's Python name to its value
var pyIPPTagByName = map[string]goipp.Tag{}

// init populates the pyIPPTagByName name
func init() {
	for tag, name := range ippTagName {
		pyIPPTagByName[name] = tag
	}
}
