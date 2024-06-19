// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// CUPS requests and responses

package ipp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/OpenPrinting/goipp"
	"github.com/alexpevzner/mfp/transport"
)

var (
	_ Request  = &CUPSGetDefaultRequest{}
	_ Response = &CUPSGetDefaultResponse{}
)

// TestCupsRequests tests CUPS requests
func TestCupsRequests(t *testing.T) {
	const (
		ippVersion   = goipp.DefaultVersion
		ippRequestID = 12345
	)

	hdr := RequestHeader{
		Version:                   ippVersion,
		RequestID:                 ippRequestID,
		AttributesCharset:         DefaultCharset,
		AttributesNaturalLanguage: DefaultNaturalLanguage,
	}

	type testData struct {
		op  goipp.Op       // Operation code
		rq  Request        // Pointer to Request structure
		msg *goipp.Message // Its IPP representation
		err string         // expected decode error
	}

	tests := []testData{
		// ----- CUPSGetDefaultRequest tests -----
		{
			op: 0x4001,

			rq: &CUPSGetDefaultRequest{
				RequestHeader:       hdr,
				RequestedAttributes: []string{"all"},
			},

			msg: goipp.NewMessageWithGroups(
				ippVersion,
				goipp.Code(goipp.OpCupsGetDefault),
				ippRequestID,
				goipp.Groups{
					{
						Tag: goipp.TagOperationGroup,
						Attrs: []goipp.Attribute{
							goipp.MakeAttribute(
								"attributes-charset",
								goipp.TagCharset,
								goipp.String(DefaultCharset)),
							goipp.MakeAttribute(
								"attributes-natural-language",
								goipp.TagLanguage,
								goipp.String(DefaultNaturalLanguage)),
							goipp.MakeAttribute(
								"requested-attributes",
								goipp.TagKeyword,
								goipp.String("all")),
						},
					},
				},
			),
		},

		{
			rq: &CUPSGetDefaultRequest{},

			msg: goipp.NewMessageWithGroups(
				ippVersion,
				goipp.Code(goipp.OpCupsGetDefault),
				ippRequestID,

				goipp.Groups{
					{
						Tag: goipp.TagOperationGroup,
						Attrs: []goipp.Attribute{
							goipp.MakeAttribute(
								"attributes-charset",
								goipp.TagInteger,
								goipp.Integer(111)),
							goipp.MakeAttribute(
								"attributes-natural-language",
								goipp.TagLanguage,
								goipp.String(DefaultNaturalLanguage)),
							goipp.MakeAttribute(
								"requested-attributes",
								goipp.TagKeyword,
								goipp.String("all")),
						},
					},
				},
			),

			err: `IPP decode ipp.CUPSGetDefaultRequest: "attributes-charset": can't convert integer to String`,
		},
	}

	for _, test := range tests {
		rqType := reflect.TypeOf(test.rq).Elem()

		// Test of Encode and Request interface
		if test.err == "" {
			if test.rq.GetOp() != test.op {
				t.Errorf("%s: GetOp(): expected 0x%x, present 0x%x",
					rqType, int(test.op), int(test.rq.GetOp()))
			}

			if test.rq.Header().Version != ippVersion {
				t.Errorf("%s: Header().Version: expected %s, present %s",
					rqType, ippVersion, test.rq.Header().Version)
			}

			if test.rq.Header().RequestID != ippRequestID {
				t.Errorf("%s: Header().RequestID: expected %d, present %d",
					rqType, ippRequestID, test.rq.Header().RequestID)
			}

			msg := test.rq.Encode()
			if !msg.Similar(*test.msg) {
				buf := &bytes.Buffer{}

				t.Errorf("Encode error. Message expected:")
				test.msg.Print(buf, true)
				t.Errorf("Message expected:\n%s", buf)

				buf.Reset()
				msg.Print(buf, true)
				t.Errorf("Message received:\n%s", buf)
			}
		}

		// Decode test
		rq := reflect.New(rqType).Interface().(Request)

		err := rq.Decode(test.msg)
		if err == nil {
			err = errors.New("")
		}

		if err.Error() != test.err {
			t.Errorf("Error mismatch:\n  <<< %s\n  >>> %s\n", test.err, err)
		} else if test.err == "" {
			diff := testDiffStruct(test.rq, rq)
			if diff != "" {
				t.Errorf("Decoded data doesn't match:\n%s", diff)
			}
		}
	}
}

// TestCupsRequests tests CUPS responses
func TestCupsRequesponses(t *testing.T) {
	const (
		ippVersion   = goipp.DefaultVersion
		ippRequestID = 12345
	)

	hdr := ResponseHeader{
		Version:                   ippVersion,
		RequestID:                 ippRequestID,
		AttributesCharset:         DefaultCharset,
		AttributesNaturalLanguage: DefaultNaturalLanguage,
		StatusMessage:             "success",
	}

	type testData struct {
		rsp Response       // Pointer to Response structure
		msg *goipp.Message // Its IPP representation
		err string         // Expected decode error
	}

	tests := []testData{
		// ----- CUPSGetDefaultResponse tests -----
		{
			rsp: &CUPSGetDefaultResponse{
				ResponseHeader: hdr,
			},

			msg: goipp.NewMessageWithGroups(
				ippVersion,
				goipp.Code(goipp.StatusOk),
				ippRequestID,

				goipp.Groups{
					{
						Tag: goipp.TagOperationGroup,
						Attrs: []goipp.Attribute{
							goipp.MakeAttribute(
								"attributes-charset",
								goipp.TagCharset,
								goipp.String(DefaultCharset)),
							goipp.MakeAttribute(
								"attributes-natural-language",
								goipp.TagLanguage,
								goipp.String(DefaultNaturalLanguage)),
							goipp.MakeAttribute(
								"status-message",
								goipp.TagText,
								goipp.String("success")),
						},
					},
				},
			),
		},

		{
			rsp: &CUPSGetDefaultResponse{},

			msg: goipp.NewMessageWithGroups(
				ippVersion,
				goipp.Code(goipp.StatusOk),
				ippRequestID,

				goipp.Groups{
					{
						Tag: goipp.TagOperationGroup,
						Attrs: []goipp.Attribute{
							goipp.MakeAttribute(
								"attributes-charset",
								goipp.TagCharset,
								goipp.String(DefaultCharset)),
							goipp.MakeAttribute(
								"attributes-natural-language",
								goipp.TagBoolean,
								goipp.Boolean(true)),
							goipp.MakeAttribute(
								"status-message",
								goipp.TagText,
								goipp.String("success")),
						},
					},
				},
			),

			err: `IPP decode ipp.CUPSGetDefaultResponse: "attributes-natural-language": can't convert boolean to String`,
		},
	}

	for _, test := range tests {
		rspType := reflect.TypeOf(test.rsp).Elem()

		// Test of Encode and Response interface
		if test.err == "" {
			if test.rsp.Header().Version != ippVersion {
				t.Errorf("%s: Header().Version: expected %s, present %s",
					rspType, ippVersion, test.rsp.Header().Version)
			}

			if test.rsp.Header().RequestID != ippRequestID {
				t.Errorf("%s: Header().RequestID: expected %d, present %d",
					rspType, ippRequestID, test.rsp.Header().RequestID)
			}

			msg := test.rsp.Encode()
			if !msg.Similar(*test.msg) {
				buf := &bytes.Buffer{}

				t.Errorf("Encode error. Message expected:")
				test.msg.Print(buf, false)
				t.Errorf("Message expected:\n%s", buf)

				buf.Reset()
				msg.Print(buf, false)
				t.Errorf("Message received:\n%s", buf)
			}
		}

		// Decode test
		rsp := reflect.New(rspType).Interface().(Response)

		err := rsp.Decode(test.msg)
		if err == nil {
			err = errors.New("")
		}

		if err.Error() != test.err {
			t.Errorf("Error mismatch:\n  <<< %s\n  >>> %s\n", test.err, err)
		} else if test.err == "" {
			diff := testDiffStruct(test.rsp, rsp)
			if diff != "" {
				t.Errorf("Decoded data doesn't match:\n%s", diff)
			}
		}
	}
}

func TestCUPS(t *testing.T) {
	//	return

	c := NewClient(transport.DefaultCupsUNIX, nil)
	rq := &CUPSGetDefaultRequest{
		RequestHeader: DefaultRequestHeader,
	}
	rsp := &CUPSGetDefaultResponse{}
	err := c.Do(context.Background(), rq, rsp)
	println(err)
	if err == nil {
		fmt.Printf("%#v", rsp)
	}
}
