// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// eSCL part of Model

package modeling

import (
	"fmt"
	"io"
	"net/http"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// eSCL hook names
const (
	esclOnScanJobsRequestName      = "escl_onScanJobsRequest"
	esclOnNextDocumentResponseName = "escl_onNextDocumentResponse"
)

// esclImageFilter defines image filtering options for the
// images, received from the eSCL scanner.
type esclImageFilter struct {
	OutputFormat optional.Val[string]         // Output format
	XResolution  optional.Val[int]            // X resolution, DPI
	YResolution  optional.Val[int]            // Y resolution, DPI
	ColorMode    optional.Val[escl.ColorMode] // Desired color mode
}

// FilterOptions exports esclImageFilter settings as abstract.FilterOptions.
func (flt *esclImageFilter) FilterOptions() (opt abstract.FilterOptions) {
	if flt.OutputFormat != nil {
		opt.OutputFormat = *flt.OutputFormat
	}

	if flt.XResolution != nil && flt.YResolution != nil {
		opt.Res.XResolution = *flt.XResolution
		opt.Res.YResolution = *flt.YResolution
	}

	if flt.ColorMode != nil {
		switch *flt.ColorMode {
		case escl.BlackAndWhite1:
			opt.Mode = abstract.ColorModeBinary
		case escl.Grayscale8:
			opt.Mode = abstract.ColorModeMono
			opt.Depth = abstract.ColorDepth8
		case escl.Grayscale16:
			opt.Mode = abstract.ColorModeMono
			opt.Depth = abstract.ColorDepth16
		case escl.RGB24:
			opt.Mode = abstract.ColorModeColor
			opt.Depth = abstract.ColorDepth8
		case escl.RGB48:
			opt.Mode = abstract.ColorModeColor
			opt.Depth = abstract.ColorDepth16
		}
	}

	return
}

// SetESCLScanCaps sets the [escl.ScannerCapabilities].
func (model *Model) SetESCLScanCaps(caps *escl.ScannerCapabilities) {
	model.esclScanCaps = caps
}

// GetESCLScanCaps returns the [escl.ScannerCapabilities].
func (model *Model) GetESCLScanCaps() *escl.ScannerCapabilities {
	return model.esclScanCaps
}

// esclLoad decodes eSCL part of model. The model file assumed to
// be preloaded into the Model's Python interpreter (model.py).
func (model *Model) esclLoad() error {
	// Load and decode ScannerCapabilities
	name := "escl.scanner"
	obj := model.py.Eval(name)
	if obj.Err() != nil {
		name = "escl.caps"
		obj = model.py.Eval(name)
	}

	if err := obj.Err(); err != nil {
		err = fmt.Errorf("%s: %w", name, err)
		return err
	}

	if !obj.IsNone() {
		caps, err := esclImportScannerCapabilities(obj)
		if err != nil {
			err = fmt.Errorf("%s: %w", name, err)
			return err
		}

		model.esclScanCaps = caps
	}

	// Load eSCL hooks -- escl_onScanJobsRequest
	obj = model.py.GetGlobal(esclOnScanJobsRequestName)
	switch {
	case obj.NotFound():
	case obj.Err() != nil:
		err := obj.Err()
		return fmt.Errorf("%s: %w", esclOnScanJobsRequestName, err)
	case !obj.IsCallable():
		return fmt.Errorf("%s is not function",
			esclOnScanJobsRequestName)
	default:
		model.esclOnScanJobsRequestScriptlet = obj
	}

	// Load eSCL hooks -- escl_onNextDocumentResponse
	obj = model.py.GetGlobal(esclOnNextDocumentResponseName)
	switch {
	case obj.NotFound():
	case obj.Err() != nil:
		err := obj.Err()
		return fmt.Errorf("%s: %w", esclOnNextDocumentResponseName, err)
	case !obj.IsCallable():
		return fmt.Errorf("%s is not function",
			esclOnNextDocumentResponseName)
	default:
		model.esclOnNextDocumentResponseScriptlet = obj
	}

	return nil
}

// NewESCLServer creates a virtual eSCL server based on eSCL scanner
// capabilities, defined by the model (see also [Model.SetESCLScanCaps]).
//
// The actual scanning facilities provided by the supplied [abstract.Scanner].
//
// It will return nil, if model doesn't have the eSCL scanner capabilities.
func (model *Model) NewESCLServer(
	scanner abstract.Scanner) *escl.AbstractServer {

	// Obtain scanner capabilities
	caps := model.GetESCLScanCaps()
	if caps == nil {
		return nil
	}

	// Setup hooks
	hooks := escl.ServerHooks{
		OnScannerCapabilitiesResponse: model.esclOnScannerCapabilitiesResponse,
		OnScanJobsRequest:             model.esclOnScanJobsRequest,
		OnScanJobsResponse:            model.esclOnScanJobsResponse,
		OnNextDocumentResponse:        model.esclOnNextDocumentResponse,
	}

	// Setup options
	options := escl.AbstractServerOptions{
		Version:  caps.Version,
		Scanner:  scanner,
		BasePath: "/eSCL",
		Hooks:    hooks,
	}

	// Create the eSCL server
	return escl.NewAbstractServer(options)
}

// esclOnScannerCapabilitiesResponse implements the
// [escl.ServerHooks.OnScannerCapabilitiesResponse] hook
// for the modeled eSCL scanner.
func (model *Model) esclOnScannerCapabilitiesResponse(
	query *transport.ServerQuery,
	caps *escl.ScannerCapabilities) *escl.ScannerCapabilities {

	caps2 := model.GetESCLScanCaps()
	if caps2 == nil {
		query.Reject(http.StatusServiceUnavailable, nil)
	}

	return caps2
}

// esclOnScanJobsRequest implements the [escl.ServerHooks.OnScanJobsRequest]
// hook for he modeled eSCL scanner.
func (model *Model) esclOnScanJobsRequest(
	query *transport.ServerQuery,
	ss *escl.ScanSettings) *escl.ScanSettings {

	// Nothing to do if scriptlet not in use
	if model.esclOnScanJobsRequestScriptlet == nil {
		return nil
	}

	// Setup logging
	ctx := query.RequestContext()
	log.Debug(ctx, "MODEL: calling %s", esclOnScanJobsRequestName)

	var err error

	defer func() {
		println(err)
		if err != nil {
			log.Begin(ctx).
				Error("MODEL: on %s:", esclOnScanJobsRequestName).
				Error("MODEL:   %s", err).
				Commit()
		}
	}()

	// Export request to Python
	q := model.queryToPython(query)
	if err := q.Err(); err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	rq := structExport(model.py, keywordMapESCL, ss)
	if err := rq.Err(); err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	// Call the hook
	err = model.esclOnScanJobsRequestScriptlet.Call(q, rq).Err()
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	// Convert possibly modified request back to Go
	err = model.queryFromPython(query, q)
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	ss2, err := esclImportScanSettings(rq)
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	return ss2
}

// esclOnScanJobsResponse implements the [escl.ServerHooks.OnScanJobsResponse]
// hook for he modeled eSCL scanner.
func (model *Model) esclOnScanJobsResponse(
	query *transport.ServerQuery,
	ss *escl.ScanSettings, joburi string) string {

	// Save ScanSettings of the last successful ScanJobs request.
	// We will need it later for image filtering.
	//
	// Any explicit synchronization is not required here,
	// because only a single scan request can be active
	// in time.
	model.esclScanSettings = *ss

	return ""
}

// esclOnScanJobsRequest implements the [escl.ServerHooks.model.OnNextDocumentResponse]
// hook for he modeled eSCL scanner.
func (model *Model) esclOnNextDocumentResponse(
	query *transport.ServerQuery,
	body io.ReadCloser) io.ReadCloser {

	// Nothing to do if scriptlet not in use
	if model.esclOnNextDocumentResponseScriptlet == nil {
		return nil
	}

	// Setup logging
	ctx := query.RequestContext()
	log.Debug(ctx, "MODEL: calling %s", esclOnNextDocumentResponseName)

	var err error

	defer func() {
		println(err)
		if err != nil {
			log.Begin(ctx).
				Error("MODEL: on %s:", esclOnNextDocumentResponseName).
				Error("MODEL:   %s", err).
				Commit()
		}
	}()

	// Export request to Python
	q := model.queryToPython(query)
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	flt := model.py.NewObject(map[any]any(nil))
	if err := flt.Err(); err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	// Call the hook
	err = model.esclOnNextDocumentResponseScriptlet.Call(q, flt).Err()
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	// Convert possibly modified request back to Go
	err = model.queryFromPython(query, q)
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	var filter esclImageFilter
	err = structImport(flt, keywordMapESCL, &filter)
	if err != nil {
		query.Reject(http.StatusServiceUnavailable, err)
		return nil
	}

	opt := filter.FilterOptions()
	var res abstract.Resolution

	if model.esclScanSettings.XResolution != nil &&
		model.esclScanSettings.YResolution != nil {
		res.XResolution = *model.esclScanSettings.XResolution
		res.YResolution = *model.esclScanSettings.YResolution
	}

	return abstract.NewStreamFilter(body, res, "", opt)
}
