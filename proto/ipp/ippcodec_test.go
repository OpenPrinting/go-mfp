// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// IPP codec test

package ipp

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/ipp/iana"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// TestIppCodecGenerate tests ippCodecGenerate() function
func TestIppCodecGenerate(t *testing.T) {
	type testData struct {
		data interface{} // Input structure
		err  string      // Expected error, "" if none
	}

	// GoodEmbedded used as embedded structure for ippCodecGenerate test
	type GoodEmbedded struct {
		X int
	}

	// BadEmbedded is bad embedded structure, for ippCodecGenerate test
	type BadEmbedded struct {
		X int `ipp:"x,boolean"`
	}

	tests := []testData{
		{
			data: struct {
				FldOk      int `ipp:"fld-ok,integer"`
				unexported string
				GoodEmbedded
			}{},
		},

		{
			data: struct {
				FldOk int `ipp:"fld-ok,integer"`
				BadEmbedded
			}{},
			err: `struct {...}.BadEmbedded.X: can't represent int as boolean`,
		},

		{
			data: struct {
				FldNoIPPTag int
			}{},
			err: `struct {...}: contains no IPP fields`,
		},

		{
			data: struct {
				FldBad int `ipp:""`
			}{},
			err: `struct {...}.FldBad: ipp:missed attribute name`,
		},

		{
			data: struct {
				FldBad float64 `ipp:"flg-bad,integer"`
			}{},
			err: `struct {...}.FldBad: float64 type not supported`,
		},

		{
			data: struct {
				Nested struct {
					FldBad float64 `ipp:"flg-bad,integer"`
				} `ipp:"flg-nested,collection"`
			}{},
			err: `struct {...}.Nested: struct {...}.Nested.FldBad: float64 type not supported`,
		},

		{
			data: struct {
				// ipp: tag contains unknown keyword
				FldBadTag int `ipp:"fld-bad-tag,unknown"`
			}{},
			err: `struct {...}.FldBadTag: ipp:"unknown": no data tags defined`,
		},

		{
			data: struct {
				// ipp: range constraint syntactically invalid
				FldBadTag int `ipp:"fld-bad-tag,integer(0:XXX)"`
			}{},
			err: `struct {...}.FldBadTag: ipp:"XXX": invalid limit`,
		},

		{
			data: struct {
				// ipp: range lower bound doesn't fit int32
				FldGoodTag int `ipp:"fld-good-tag,integer(4294967296:5)"`
			}{},
			err: `struct {...}.FldGoodTag: ipp:"4294967296": invalid limit`,
		},

		{
			data: struct {
				// ipp: range upper bound doesn't fit int32
				FldGoodTag int `ipp:"fld-good-tag,integer(5:4294967296)"`
			}{},
			err: `struct {...}.FldGoodTag: ipp:"4294967296": invalid limit`,
		},

		{
			data: struct {
				// ipp: range min > max
				FldGoodTag int `ipp:"fld-good-tag,integer(10:5)"`
			}{},
			err: `struct {...}.FldGoodTag: ipp:integer(10:5): min>max`,
		},

		{
			data: struct {
				// ipp: tag contains valid range constraint
				FldGoodTag int `ipp:"fld-good-tag,integer(0:100)"`
			}{},
		},

		{
			data: struct {
				FldConv int `ipp:"fld-conv,string"`
			}{},
			err: `struct {...}.FldConv: can't represent int as octetString`,
		},

		{
			data: struct {
				FldConv string `ipp:"fld-conv,enum"`
			}{},
			err: `struct {...}.FldConv: can't represent string as enum`,
		},

		{
			data: struct {
				FldConv bool `ipp:"fld-conv,keyword"`
			}{},
			err: `struct {...}.FldConv: can't represent bool as keyword`,
		},

		{
			data: struct {
				Fld1 int `ipp:"fld,integer"`
				Fld2 int `ipp:"fld,integer"`
			}{},
			err: `struct {...}.Fld2: attribute "fld" already used by Fld1`,
		},

		{
			data: struct {
				unexported int `ipp:"unexported"`
			}{},
			err: `struct {...}.unexported: ipp:tag used with unexported field`,
		},

		{
			data: struct{}{},
			err:  `struct {}: contains no IPP fields`,
		},
	}

	for _, test := range tests {
		_, err := ippCodecGenerate(reflect.TypeOf(test.data))

		var errExpected error
		if test.err != "" {
			errExpected = errors.New(test.err)
		}

		checkError(t, "TestIppCodecGenerate", err, errExpected)
	}
}

// TestIppCodecStandardTypes tests that ippCodecGenerate() successfully
// generates codecs for  the standard types.
func TestIppCodecStandardTypes(t *testing.T) {
	tests := []Object{
		&CUPSGetDefaultRequest{},
		&CUPSGetDefaultResponse{},
		&CUPSGetDevicesRequest{},
		&CUPSGetDevicesResponse{},
		&CUPSGetPPDRequest{},
		&CUPSGetPPDResponse{},
		&CUPSGetPPDsRequest{},
		&CUPSGetPPDsResponse{},
		&CUPSGetPrintersRequest{},
		&CUPSGetPrintersResponse{},

		&CreateJobRequest{},
		&CreateJobResponse{},
		&JobAttributes{},
		&PPDAttributes{},
		&PrinterAttributes{},
		&SendDocumentRequest{},
		&SendDocumentResponse{},
		&ValidateJobRequest{},
		&ValidateJobResponse{},

		&DeviceAttributes{},
		&DocumentStatus{},
		&GetNextDocumentDataRequest{},
		&GetNextDocumentDataResponse{},
		&GetPrinterAttributesRequest{},
		&GetPrinterAttributesResponse{},
		&JobStatus{},
	}

	for _, test := range tests {
		ty := reflect.TypeOf(test).Elem()
		_, err := ippCodecGenerate(ty)
		if err != nil {
			t.Errorf("ippCodecGenerate: %s: %s",
				reflect.TypeOf(test), err)
		}
	}
}

// TestDecodePanic tests panics in (*ippCodec) decode()
func TestDecodePanic(t *testing.T) {
	// Compile the codec
	ttype := reflect.TypeOf(ippTestStruct{})
	codec := ippCodecMustGenerate(ttype)

	var attrs goipp.Attributes
	attrs.Add(goipp.MakeAttribute("test", goipp.TagInteger, goipp.Integer(5)))

	p := &PrinterAttributes{}
	errExpected := errors.New(`Decoder for "*ipp.ippTestStruct" applied to "**ipp.PrinterAttributes"`)

	defer func() {
		p := recover()
		err, ok := p.(error)
		if !ok {
			panic(p)
		}

		checkError(t, "TestDecodePanic", err, errExpected)
	}()

	dec := NewDecoder(nil)
	codec.decodeAttrs(dec, &p, attrs)
}

// testFakeObject implements Object interface, but it is not structure.
// It is used for testing
type testFakeObject int

func (testFakeObject) RawAttrs() *ObjectRawAttrs               { return nil }
func (testFakeObject) Get(string) (goipp.Attribute, bool)      { return goipp.Attribute{}, false }
func (testFakeObject) Set(goipp.Attribute) error               { return nil }
func (testFakeObject) Errors() []error                         { return nil }
func (testFakeObject) registrations() map[string]*iana.DefAttr { return nil }

// TestIppEncodeDecodeAttrsPanic tests panic in
// ippEncodeAttrs and ippDecodeAttrs
func TestIppEncodeDecodePanic(t *testing.T) {
	doTest := func(action func(), errExpected error) {
		defer func() {
			p := recover()
			var err error
			if p != nil {
				var ok bool
				err, ok = p.(error)
				if !ok {
					panic(p)
				}
			}

			checkError(t, "TestDecodePanic", err, errExpected)
		}()

		action()
	}

	enc := ippEncoder{}
	dec := NewDecoder(nil)

	doTest(func() {}, nil)

	doTest(func() { enc.Encode(new(testFakeObject)) },
		errors.New("ipp.testFakeObject: is not struct"))

	doTest(func() { dec.Decode(new(testFakeObject), nil) },
		errors.New("ipp.testFakeObject: is not struct"))
}

// ----- Decode test -----

// ippTestCollection used to test collection members
// of ippTestStruct
type ippTestCollection struct {
	CollInt    goipp.IntegerOrRange `ipp:"coll-int,integer|rangeOfInteger"`
	CollString string               `ipp:"coll-string,text"`
	CollU16    uint16               `ipp:"coll-u16,integer"`
}

// TestEmbedded used to test embedding
type TestEmbedded struct {
	FldEmbedded int `ipp:"fld-embedded,integer"`
}

// ippTestStruct is the structure, intended for testing
// of the IPP codec
type ippTestStruct struct {
	ObjectRawAttrs
	OperationGroup // FIXME

	TestEmbedded

	FldBooleanF     bool   `ipp:"fld-boolean-f,boolean"`
	FldBooleanSlice []bool `ipp:"fld-boolean-slice,1setOf boolean"`
	FldBooleanT     bool   `ipp:"fld-boolean-t,boolean"`

	FldCharset      string   `ipp:"fld-charset,charset"`
	FldCharsetSlice []string `ipp:"fld-charset-slice,1setOf charset"`

	FldColl         ippTestCollection   `ipp:"fld-coll,collection"`
	FldCollSlice    []ippTestCollection `ipp:"fld-coll-slice,1setOf collection"`
	FldCollNilSlice []ippTestCollection `ipp:"fld-coll-nil-slice,1setOf collection"`

	FldDateTime      time.Time   `ipp:"fld-datetime,datetime"`
	FldDateTimeSlice []time.Time `ipp:"fld-datetime-slice,1setOf datetime"`

	FldEnum      int   `ipp:"fld-enum,enum"`
	FldEnumSlice []int `ipp:"fld-enum-slice,1setOf enum"`

	FldInteger      int   `ipp:"fld-integer,integer"`
	FldIntegerSlice []int `ipp:"fld-integer-slice,1setOf integer"`

	FldKeyword      string   `ipp:"fld-keyword,keyword"`
	FldKeywordSlice []string `ipp:"fld-keyword-slice,1setOf keyword"`

	FldLanguage      string   `ipp:"fld-language,naturalLanguage"`
	FldLanguageSlice []string `ipp:"fld-language-slice,1setOf naturalLanguage"`

	FldMime      string   `ipp:"fld-mime,mimemediatype"`
	FldMimeSlice []string `ipp:"fld-mime-slice,1setOf mimemediatype"`

	FldName      string   `ipp:"fld-name,name"`
	FldNameSlice []string `ipp:"fld-name-slice,1setOf name"`

	FldNilSlice []int `ipp:"fld-nil-slice,1setOf integer"`

	FldRange      goipp.Range   `ipp:"fld-range,rangeOfInteger"`
	FldRangeSlice []goipp.Range `ipp:"fld-range-slice,1setOf rangeOfInteger"`

	FldResolution      goipp.Resolution   `ipp:"fld-resolution,resolution"`
	FldResolutionSlice []goipp.Resolution `ipp:"fld-resolution-slice,1setOf resolution"`

	FldTextWithLang      goipp.TextWithLang   `ipp:"fld-textwithlang,textwithlanguage"`
	FldTextWithLangSlice []goipp.TextWithLang `ipp:"fld-textwithlang-slice,1setOf textwithlanguage"`

	FldString      string   `ipp:"fld-string,string"`
	FldStringSlice []string `ipp:"fld-string-slice,1setOf string"`

	FldText      string   `ipp:"fld-text,text"`
	FldTextSlice []string `ipp:"fld-text-slice,1setOf text"`

	FldURI      string   `ipp:"fld-uri,uri"`
	FldURISlice []string `ipp:"fld-uri-slice,1setOf uri"`

	FldURIScheme      string   `ipp:"fld-urischeme,urischeme"`
	FldURISchemeSlice []string `ipp:"fld-urischeme-slice,1setOf urischeme"`

	FldUint16      uint16   `ipp:"fld-uint16,integer"`
	FldUint16Slice []uint16 `ipp:"fld-uint16-slice,1setOf integer"`

	FldVersion      goipp.Version   `ipp:"fld-version,keyword"`
	FldVersionSlice []goipp.Version `ipp:"fld-version-slice,1setOf keyword"`

	FltOptionalMissed  optional.Val[int] `ipp:"flt-optional-missed,integer"`
	FltOptionalPresent optional.Val[int] `ipp:"flt-optional-present,integer"`
}

// ippDecodeTest represents a single decode test data
type ippDecodeTest struct {
	attrs goipp.Attributes // Input attributes
	data  *ippTestStruct   // Expected decoded data
	err   error            // Expected error
}

// ippDecodeTestData is a collection of decode tests
var ippDecodeTestData = []ippDecodeTest{
	// ----- Test for errors -----
	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-boolean-f",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-boolean-f": can't use integer as boolean`),
	},

	{
		attrs: goipp.Attributes{
			goipp.Attribute{
				Name: "fld-boolean-slice",
				Values: goipp.Values{
					{
						T: goipp.TagBoolean,
						V: goipp.Boolean(true),
					},
					{
						T: goipp.TagBoolean,
						V: goipp.Boolean(false),
					},
					{
						T: goipp.TagString,
						V: goipp.String("hello"),
					},
				},
			},
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-boolean-slice[2]": can't use octetString as 1setOf boolean`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-datetime",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-datetime": can't use integer as dateTime`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-coll",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-coll": can't use integer as collection`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-coll",
				goipp.TagBeginCollection,
				goipp.Collection{
					goipp.MakeAttribute("coll-int",
						goipp.TagBoolean, goipp.Boolean(true)),
					goipp.MakeAttribute("coll-string",
						goipp.TagText, goipp.String("hello")),
					goipp.MakeAttribute("coll-u16",
						goipp.TagInteger, goipp.Integer(5)),
				},
			),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-coll/coll-int": can't use boolean as integer | rangeOfInteger`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-coll-slice",
				goipp.TagBeginCollection,
				goipp.Collection{
					goipp.MakeAttribute("coll-int",
						goipp.TagBoolean, goipp.Boolean(true)),
					goipp.MakeAttribute("coll-string",
						goipp.TagText, goipp.String("hello")),
					goipp.MakeAttribute("coll-u16",
						goipp.TagInteger, goipp.Integer(5)),
				},
			),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-coll-slice[0]/coll-int": can't use boolean as integer | rangeOfInteger`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-coll-slice",
				goipp.TagBeginCollection,
				goipp.Collection{
					goipp.MakeAttribute("coll-int",
						goipp.TagInteger, goipp.Integer(0)),
					goipp.MakeAttribute("coll-string",
						goipp.TagText, goipp.String("hello")),
					goipp.MakeAttribute("coll-u16",
						goipp.TagInteger, goipp.Integer(65536)),
				},
			),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-coll-slice[0]/coll-u16": Value 65536 out of range`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-coll-slice",
				goipp.TagBeginCollection,
				goipp.Collection{
					goipp.MakeAttribute("coll-int",
						goipp.TagRange,
						goipp.Range{Lower: 5, Upper: 7}),
					goipp.MakeAttribute("coll-string",
						goipp.TagText, goipp.String("hello")),
					goipp.MakeAttribute("coll-u16",
						goipp.TagInteger, goipp.Integer(65536)),
				},
			),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-coll-slice[0]/coll-u16": Value 65536 out of range`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-enum",
				goipp.TagText, goipp.String("12345")),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-enum": can't use textWithoutLanguage as enum`),
	},

	{
		attrs: goipp.Attributes{
			goipp.Attribute{
				Name: "fld-integer",
			},
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-integer": at least 1 value required`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-range",
				goipp.TagText, goipp.String("12345")),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-range": can't use textWithoutLanguage as rangeOfInteger`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-resolution",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-resolution": can't use integer as resolution`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-textwithlang",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-textwithlang": can't use integer as text`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-string",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-string": can't use integer as octetString`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-uint16",
				goipp.TagInteger, goipp.Integer(65536)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-uint16": Value 65536 out of range`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-uint16",
				goipp.TagInteger, goipp.Integer(-1)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-uint16": Value -1 out of range`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-uint16",
				goipp.TagText, goipp.String("12345")),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-uint16": can't use textWithoutLanguage as integer`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-version",
				goipp.TagKeyword, goipp.String("12345")),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-version": "12345": invalid version string`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-version",
				goipp.TagKeyword, goipp.String("aaa.bbb")),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-version": "aaa.bbb": invalid version string`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-version",
				goipp.TagKeyword, goipp.String("123.bbb")),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-version": "123.bbb": invalid version string`),
	},

	{
		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-version",
				goipp.TagInteger, goipp.Integer(12345)),
		},

		err: errors.New(`IPP decode ipp.ippTestStruct: "fld-version": can't use integer as keyword`),
	},

	// ----- Big test of successful decoding -----
	{

		attrs: goipp.Attributes{
			goipp.MakeAttribute("fld-embedded",
				goipp.TagInteger, goipp.Integer(123)),

			goipp.MakeAttribute("fld-boolean-f",
				goipp.TagBoolean, goipp.Boolean(false)),
			goipp.MakeAttribute("fld-boolean-t",
				goipp.TagBoolean, goipp.Boolean(true)),
			goipp.Attribute{
				Name: "fld-boolean-slice",
				Values: goipp.Values{
					{
						T: goipp.TagBoolean,
						V: goipp.Boolean(true),
					},
					{
						T: goipp.TagBoolean,
						V: goipp.Boolean(false),
					},
				},
			},

			goipp.MakeAttribute("fld-charset",
				goipp.TagCharset, goipp.String("utf-8")),
			goipp.Attribute{
				Name: "fld-charset-slice",
				Values: goipp.Values{
					{
						T: goipp.TagCharset,
						V: goipp.String("ibm866"),
					},
					{
						T: goipp.TagCharset,
						V: goipp.String("iso-8859-5"),
					},
					{
						T: goipp.TagCharset,
						V: goipp.String("windows-1251"),
					},
				},
			},

			goipp.MakeAttribute("fld-coll",
				goipp.TagBeginCollection,
				goipp.Collection{
					goipp.MakeAttribute("coll-int",
						goipp.TagInteger, goipp.Integer(5)),
					goipp.MakeAttribute("coll-string",
						goipp.TagText, goipp.String("hello")),
					goipp.MakeAttribute("coll-u16",
						goipp.TagInteger, goipp.Integer(15)),
				},
			),

			goipp.Attribute{
				Name: "fld-coll-slice",
				Values: goipp.Values{
					{
						T: goipp.TagBeginCollection,
						V: goipp.Collection{
							goipp.MakeAttribute("coll-int",
								goipp.TagInteger, goipp.Integer(1)),
							goipp.MakeAttribute("coll-string",
								goipp.TagText, goipp.String("one")),
							goipp.MakeAttribute("coll-u16",
								goipp.TagInteger, goipp.Integer(10)),
						},
					},
					{
						T: goipp.TagBeginCollection,
						V: goipp.Collection{
							goipp.MakeAttribute("coll-int",
								goipp.TagInteger, goipp.Integer(2)),
							goipp.MakeAttribute("coll-string",
								goipp.TagText, goipp.String("two")),
							goipp.MakeAttribute("coll-u16",
								goipp.TagInteger, goipp.Integer(20)),
						},
					},
					{
						T: goipp.TagBeginCollection,
						V: goipp.Collection{
							goipp.MakeAttribute("coll-int",
								goipp.TagRange,
								goipp.Range{Lower: 5, Upper: 7}),
							goipp.MakeAttribute("coll-string",
								goipp.TagText, goipp.String("many")),
							goipp.MakeAttribute("coll-u16",
								goipp.TagInteger, goipp.Integer(30)),
						},
					},
				},
			},

			goipp.MakeAttribute("fld-coll-slice",
				goipp.TagBeginCollection,
				goipp.Collection{
					goipp.MakeAttribute("coll-int",
						goipp.TagInteger, goipp.Integer(5)),
					goipp.MakeAttribute("coll-string",
						goipp.TagText, goipp.String("hello")),
				},
			),

			goipp.MakeAttribute("fld-datetime",
				goipp.TagDateTime, goipp.Time{Time: testTime1}),
			goipp.Attribute{
				Name: "fld-datetime-slice",
				Values: goipp.Values{
					{
						T: goipp.TagDateTime,
						V: goipp.Time{Time: testTime2},
					},
					{
						T: goipp.TagDateTime,
						V: goipp.Time{Time: testTime3},
					},
					{
						T: goipp.TagDateTime,
						V: goipp.Time{Time: testTime4},
					},
				},
			},

			goipp.MakeAttribute("fld-enum",
				goipp.TagEnum, goipp.Integer(4321)),
			goipp.Attribute{
				Name: "fld-enum-slice",
				Values: goipp.Values{
					{
						T: goipp.TagEnum,
						V: goipp.Integer(3),
					},
					{
						T: goipp.TagEnum,
						V: goipp.Integer(2),
					},
					{
						T: goipp.TagEnum,
						V: goipp.Integer(1),
					},
				},
			},

			goipp.MakeAttribute("fld-integer",
				goipp.TagInteger, goipp.Integer(1234)),
			goipp.Attribute{
				Name: "fld-integer-slice",
				Values: goipp.Values{
					{
						T: goipp.TagInteger,
						V: goipp.Integer(1),
					},
					{
						T: goipp.TagInteger,
						V: goipp.Integer(2),
					},
					{
						T: goipp.TagInteger,
						V: goipp.Integer(3),
					},
				},
			},

			goipp.MakeAttribute("fld-keyword",
				goipp.TagKeyword, goipp.String("document")),
			goipp.Attribute{
				Name: "fld-keyword-slice",
				Values: goipp.Values{
					{
						T: goipp.TagKeyword,
						V: goipp.String("one-sided"),
					},
					{
						T: goipp.TagKeyword,
						V: goipp.String("two-sided-short-edge"),
					},
					{
						T: goipp.TagKeyword,
						V: goipp.String("two-sided-long-edge"),
					},
				},
			},

			goipp.MakeAttribute("fld-language",
				goipp.TagLanguage, goipp.String("en-US")),
			goipp.Attribute{
				Name: "fld-language-slice",
				Values: goipp.Values{
					{
						T: goipp.TagLanguage,
						V: goipp.String("be-BY"),
					},
					{
						T: goipp.TagLanguage,
						V: goipp.String("ru-RU"),
					},
					{
						T: goipp.TagLanguage,
						V: goipp.String("uk-UA"),
					},
				},
			},

			goipp.MakeAttribute("fld-mime",
				goipp.TagMimeType, goipp.String("application/pdf")),
			goipp.Attribute{
				Name: "fld-mime-slice",
				Values: goipp.Values{
					{
						T: goipp.TagMimeType,
						V: goipp.String("image/tiff"),
					},
					{
						T: goipp.TagMimeType,
						V: goipp.String("image/jpeg"),
					},
					{
						T: goipp.TagMimeType,
						V: goipp.String("image/urf"),
					},
				},
			},

			goipp.MakeAttribute("fld-name",
				goipp.TagName, goipp.String("Printer in a classroom")),
			goipp.Attribute{
				Name: "fld-name-slice",
				Values: goipp.Values{
					{
						T: goipp.TagName,
						V: goipp.String("Job0001"),
					},
					{
						T: goipp.TagName,
						V: goipp.String("Job0002"),
					},
					{
						T: goipp.TagName,
						V: goipp.String("Job0003"),
					},
				},
			},

			goipp.MakeAttribute("fld-range",
				goipp.TagRange, goipp.Range{Lower: 1, Upper: 99}),
			goipp.Attribute{
				Name: "fld-range-slice",
				Values: goipp.Values{
					{
						T: goipp.TagRange,
						V: goipp.Range{Lower: 10000, Upper: 14800},
					},
					{
						T: goipp.TagRange,
						V: goipp.Range{Lower: 21600, Upper: 35600},
					},
				},
			},

			goipp.MakeAttribute("fld-resolution",
				goipp.TagResolution,
				goipp.Resolution{Xres: 100, Yres: 150, Units: goipp.UnitsDpi}),
			goipp.Attribute{
				Name: "fld-resolution-slice",
				Values: goipp.Values{
					{
						T: goipp.TagResolution,
						V: goipp.Resolution{
							Xres:  200,
							Yres:  300,
							Units: goipp.UnitsDpi,
						},
					},
					{T: goipp.TagResolution,
						V: goipp.Resolution{
							Xres:  400,
							Yres:  500,
							Units: goipp.UnitsDpcm,
						}},
				},
			},

			goipp.MakeAttribute("fld-textwithlang",
				goipp.TagTextLang,
				goipp.TextWithLang{Lang: "en-US", Text: "Hello"}),
			goipp.Attribute{
				Name: "fld-textwithlang-slice",
				Values: goipp.Values{
					{
						T: goipp.TagTextLang,
						V: goipp.TextWithLang{
							Lang: "be-BY",
							Text: "Прывітанне",
						},
					},
					{
						T: goipp.TagTextLang,
						V: goipp.TextWithLang{
							Lang: "ru-RU",
							Text: "Привет"},
					},
					{
						T: goipp.TagTextLang,
						V: goipp.TextWithLang{
							Lang: "uk-UA",
							Text: "Привіт"},
					},
				},
			},

			goipp.MakeAttribute("fld-string",
				goipp.TagString, goipp.String("hello, world")),
			goipp.Attribute{
				Name: "fld-string-slice",
				Values: goipp.Values{
					{
						T: goipp.TagString,
						V: goipp.String("A"),
					},
					{
						T: goipp.TagString,
						V: goipp.String("B"),
					},
					{
						T: goipp.TagString,
						V: goipp.String("C"),
					},
				},
			},

			goipp.MakeAttribute("fld-text",
				goipp.TagText, goipp.String("ping pong")),
			goipp.Attribute{
				Name: "fld-text-slice",
				Values: goipp.Values{
					{
						T: goipp.TagText,
						V: goipp.String("X"),
					},
					{
						T: goipp.TagText,
						V: goipp.String("Y"),
					},
					{
						T: goipp.TagText,
						V: goipp.String("Z"),
					},
				},
			},

			goipp.MakeAttribute("fld-uri",
				goipp.TagURI, goipp.String("http://example.com")),
			goipp.Attribute{
				Name: "fld-uri-slice",
				Values: goipp.Values{
					{
						T: goipp.TagURI,
						V: goipp.String("http://example.com/print"),
					},
					{
						T: goipp.TagURI,
						V: goipp.String("http://example.com/scan"),
					},
				},
			},

			goipp.MakeAttribute("fld-urischeme",
				goipp.TagURIScheme, goipp.String("http")),
			goipp.Attribute{
				Name: "fld-urischeme-slice",
				Values: goipp.Values{
					{
						T: goipp.TagURIScheme,
						V: goipp.String("tel"),
					},
					{
						T: goipp.TagURIScheme,
						V: goipp.String("mailto"),
					},
				},
			},

			goipp.MakeAttribute("fld-uint16",
				goipp.TagInteger, goipp.Integer(4567)),
			goipp.Attribute{
				Name: "fld-uint16-slice",
				Values: goipp.Values{
					{
						T: goipp.TagInteger,
						V: goipp.Integer(11),
					},
					{
						T: goipp.TagInteger,
						V: goipp.Integer(22),
					},
					{
						T: goipp.TagInteger,
						V: goipp.Integer(33),
					},
				},
			},

			goipp.MakeAttribute("fld-version",
				goipp.TagKeyword, goipp.String("2.0")),
			goipp.Attribute{
				Name: "fld-version-slice",
				Values: goipp.Values{
					{
						T: goipp.TagKeyword,
						V: goipp.String("2.0"),
					},
					{
						T: goipp.TagKeyword,
						V: goipp.String("1.1"),
					},
					{
						T: goipp.TagKeyword,
						V: goipp.String("1.0"),
					},
				},
			},

			goipp.Attribute{
				// Note: "fld-version-slice" purposely duplicated
				Name: "fld-version-slice",
				Values: goipp.Values{
					{
						T: goipp.TagKeyword,
						V: goipp.String("0.0"),
					},
					{
						T: goipp.TagKeyword,
						V: goipp.String("0.1"),
					},
				},
			},

			goipp.Attribute{
				Name: "flt-optional-present",
				Values: goipp.Values{
					{
						T: goipp.TagInteger,
						V: goipp.Integer(5),
					},
				},
			},
		},
		data: &ippTestStruct{
			TestEmbedded: TestEmbedded{123},

			FldBooleanF:     false,
			FldBooleanT:     true,
			FldBooleanSlice: []bool{true, false},

			FldColl: ippTestCollection{
				CollInt:    goipp.Integer(5),
				CollU16:    15,
				CollString: "hello",
			},
			FldCollSlice: []ippTestCollection{
				{
					CollInt:    goipp.Integer(1),
					CollU16:    10,
					CollString: "one",
				},
				{
					CollInt:    goipp.Integer(2),
					CollU16:    20,
					CollString: "two",
				},
				{
					CollInt:    goipp.Range{Lower: 5, Upper: 7},
					CollU16:    30,
					CollString: "many",
				},
			},

			FldCharset:      "utf-8",
			FldCharsetSlice: []string{"ibm866", "iso-8859-5", "windows-1251"},

			FldDateTime: testTime1,
			FldDateTimeSlice: []time.Time{
				testTime2, testTime3, testTime4,
			},

			FldEnum:      4321,
			FldEnumSlice: []int{3, 2, 1},

			FldInteger:      1234,
			FldIntegerSlice: []int{1, 2, 3},

			FldKeyword: "document",
			FldKeywordSlice: []string{
				"one-sided",
				"two-sided-short-edge",
				"two-sided-long-edge"},

			FldLanguage:      "en-US",
			FldLanguageSlice: []string{"be-BY", "ru-RU", "uk-UA"},

			FldMime: "application/pdf",
			FldMimeSlice: []string{
				"image/tiff",
				"image/jpeg",
				"image/urf"},

			FldName: "Printer in a classroom",
			FldNameSlice: []string{
				"Job0001",
				"Job0002",
				"Job0003"},

			FldRange: goipp.Range{Lower: 1, Upper: 99},
			FldRangeSlice: []goipp.Range{
				{Lower: 10000, Upper: 14800},
				{Lower: 21600, Upper: 35600}},

			FldResolution: goipp.Resolution{Xres: 100, Yres: 150, Units: goipp.UnitsDpi},
			FldResolutionSlice: []goipp.Resolution{
				goipp.Resolution{Xres: 200, Yres: 300, Units: goipp.UnitsDpi},
				goipp.Resolution{Xres: 400, Yres: 500, Units: goipp.UnitsDpcm}},

			FldTextWithLang: goipp.TextWithLang{Lang: "en-US", Text: "Hello"},
			FldTextWithLangSlice: []goipp.TextWithLang{
				goipp.TextWithLang{Lang: "be-BY", Text: "Прывітанне"},
				goipp.TextWithLang{Lang: "ru-RU", Text: "Привет"},
				goipp.TextWithLang{Lang: "uk-UA", Text: "Привіт"}},

			FldString:      "hello, world",
			FldStringSlice: []string{"A", "B", "C"},

			FldText:      "ping pong",
			FldTextSlice: []string{"X", "Y", "Z"},

			FldURI: "http://example.com",
			FldURISlice: []string{
				"http://example.com/print",
				"http://example.com/scan"},

			FldURIScheme:      "http",
			FldURISchemeSlice: []string{"tel", "mailto"},

			FldUint16:      4567,
			FldUint16Slice: []uint16{11, 22, 33},

			FldVersion: goipp.MakeVersion(2, 0),
			FldVersionSlice: []goipp.Version{
				goipp.MakeVersion(2, 0),
				goipp.MakeVersion(1, 1),
				goipp.MakeVersion(1, 0),
			},

			FltOptionalPresent: optional.New(5),
		},
	},
}

func (test ippDecodeTest) exec(t *testing.T) {
	// Compile the codec
	ttype := reflect.TypeOf(ippTestStruct{})
	codec := ippCodecMustGenerate(ttype)

	// Decode IPP attributes
	dec := NewDecoder(nil)

	out := &ippTestStruct{}
	err := dec.Decode(out, test.attrs)

	checkError(t, "TestIppDecode", err, test.err)
	if err != nil || test.err != nil {
		return
	}

	// Compare result against expected
	diff := testDiffStruct(test.data, out)
	if diff != "" {
		t.Errorf("decode: input/output mismatch:\n%s", diff)
		return
	}

	// Now encode it back
	enc := ippEncoder{}
	attrs := codec.encodeAttrs(&enc, out)

	diff = testDiffAttrs(test.attrs, attrs)
	if diff != "" {
		t.Errorf("encode: input/output mismatch:\n%s", diff)
	}
}

func TestIppDecode(t *testing.T) {
	for _, test := range ippDecodeTestData {
		test.exec(t)
	}
}

// ----- IPP encode/decode test -----

// ippEncodeDecodeTest represents a single IPP encode/decode test
type ippEncodeDecodeTest struct {
	name  string       // Test name, for logging
	t     reflect.Type // Input type
	data  interface{}  // Input data
	panic error        // Expected panic
	err   error        // Expected error
}

// ippEncodeDecodeTestData is the test data for the IPP encode/decode test
var ippEncodeDecodeTestData = []ippEncodeDecodeTest{
	{
		name:  "panic expected: ippCodecGenerate() with invalid type",
		t:     reflect.TypeOf(int(0)),
		data:  1,
		panic: errors.New(`int: is not struct`),
	},
	{
		name:  "panic expected: ippCodec applied to wrong type",
		t:     reflect.TypeOf(PrinterAttributes{}),
		data:  1,
		panic: errors.New(`Encoder for "*ipp.PrinterAttributes" applied to "int"`),
	},
	{
		name: "success expected",
		t:    reflect.TypeOf(PrinterAttributes{}),
		data: &testdataPrinterAttributes,
	},
}

func (test ippEncodeDecodeTest) exec(t *testing.T) {
	// This function catches the possible panic
	defer func() {
		// Panic not expected - let it go its way
		if test.panic == nil {
			return
		}

		p := recover()
		if p == nil && test.panic != nil {
			t.Errorf("in test %q:", test.name)
			t.Errorf("panic expected but didn't happen: %s",
				test.panic)
			return
		}

		if p != nil {
			err, ok := p.(error)
			if !ok {
				panic(p)
			}

			if err.Error() != test.panic.Error() {
				t.Errorf("in test %q:", test.name)
				t.Errorf("panic mismatch:\n"+
					"expected: %s\n"+
					"present:  %s",
					test.panic, err)
			}
		}
	}()

	// Generate codec
	codec := ippCodecMustGenerate(test.t)

	// Test encoding
	enc := ippEncoder{}
	attrs := codec.encodeAttrs(&enc, test.data)

	// Test decoding
	dec := NewDecoder(nil)
	out := reflect.New(test.t).Interface()
	err := codec.decodeAttrs(dec, out, attrs)

	checkError(t, test.name, err, test.err)
	if err != nil {
		return
	}

	if !reflect.DeepEqual(test.data, out) {
		t.Errorf("in test %q:", test.name)
		t.Errorf("input/output mismatch")
		t.Errorf("expected: %#v\n", test.data)
		t.Errorf("present: %#v\n", out)
	}
}

// IPP encode/decode test
func TestIppEncodeDecode(t *testing.T) {
	for _, test := range ippEncodeDecodeTestData {
		test.exec(t)
	}
}

// ----- IPP decode test -----

// ----- Common stuff -----

// Check error against expected
func checkError(t *testing.T, name string, err, expected error) {
	switch {
	case err == nil && expected != nil:
		t.Errorf("in test %q:", name)
		t.Errorf("error expected but didn't happen: %s", expected)
	case err != nil && expected == nil:
		t.Errorf("in test %q:", name)
		t.Errorf("error not expected: %s", err)
	case err != nil && expected != nil && err.Error() != expected.Error():
		t.Errorf("in test %q:", name)
		t.Errorf("error mismatch:\n"+
			"expected: %s\n"+
			"present:  %s", expected, err)
	}
}

// attrsDedup removes duplicated attributes
func attrsDedup(attrs goipp.Attributes) goipp.Attributes {
	newattrs := make(goipp.Attributes, 0, len(attrs))

	for _, attr := range attrs {
		if len(newattrs) == 0 || newattrs[len(newattrs)-1].Name != attr.Name {
			newattrs = append(newattrs, attr)
		}
	}

	return newattrs
}

var (
	testTime1 = time.Date(1970, time.January, 9, 23, 0, 0, 0, time.UTC)
	testTime2 = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	testTime3 = time.Date(2019, time.April, 12, 15, 0, 0, 0, time.UTC)
	testTime4 = time.Date(2025, time.May, 17, 45, 0, 0, 0, time.UTC)
)

var testdataPrinterAttributes = PrinterAttributes{
	PrinterDescription: PrinterDescription{
		CharsetConfigured:    optional.New(DefaultCharset),
		CharsetSupported:     DefaultCharsetSupported,
		CompressionSupported: []KwCompression{KwCompressionNone},
		IppFeaturesSupported: []string{
			"airprint-1.7",
			"airprint-1.6",
			"airprint-1.5",
			"airprint-1.4",
		},
		IppVersionsSupported: DefaultIppVersionsSupported,
		MediaSizeSupported: []MediaSizeRange{
			{
				XDimension: goipp.Integer(21590),
				YDimension: goipp.Integer(27940),
			},
			{
				XDimension: goipp.Integer(21000),
				YDimension: goipp.Integer(29700),
			},
			{
				XDimension: goipp.Range{Lower: 10000, Upper: 14800},
				YDimension: goipp.Range{Lower: 21600, Upper: 35600},
			},
		},
		OperationsSupported: []goipp.Op{
			goipp.OpGetPrinterAttributes,
		},
	},
}
