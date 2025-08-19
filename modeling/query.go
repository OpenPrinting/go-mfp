// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// transport.ServerQuery to Python conversion

package modeling

import (
	"os"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/transport"
)

// queryToPython converts [transport.ServerQuery] into the [cpython.Object].
func (model *Model) queryToPython(query *transport.ServerQuery) (
	*cpython.Object, error) {

	// Create the query.Query Object
	obj, err := model.clsQuery.Call()
	if err != nil {
		return nil, err
	}

	// Convert request and response HTTP headers
	request, err := model.httpHeaderToPython(query.RequestHeader())
	if err != nil {
		return nil, err
	}

	response, err := model.httpHeaderToPython(query.ResponseHeader())
	if err != nil {
		return nil, err
	}

	// Add them to the query Object
	err = obj.SetAttr("request", request)
	if err != nil {
		return nil, err
	}

	err = obj.SetAttr("response", response)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// queryFromPython updates [transport.ServerQuery] from the
// [cpython.Object].
func (model *Model) queryFromPython(query *transport.ServerQuery,
	obj *cpython.Object) error {

	// Extract request and response
	request, err := obj.GetAttr("request")
	if err != nil {
		return err
	}

	response, err := obj.GetAttr("response")
	if err != nil {
		return err
	}

	// Convert both to the http.Header
	requestHdr, err := model.httpHeaderFromPython(request)
	if err != nil {
		return err
	}

	responseHdr, err := model.httpHeaderFromPython(response)
	if err != nil {
		return err
	}

	// Update query headers
	transport.HTTPPurgeHeaders(query.RequestHeader())
	transport.HTTPCopyHeaders(query.RequestHeader(), requestHdr)

	requestHdr.WriteSubset(os.Stdout, nil)

	transport.HTTPPurgeHeaders(query.ResponseHeader())
	transport.HTTPCopyHeaders(query.ResponseHeader(), responseHdr)

	return nil
}
