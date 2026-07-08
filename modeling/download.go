// MFP - Miulti-Function Printers and scanners toolkit
// The "model" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Download printer and scanner attributes

package modeling

import (
	"context"
	"net/url"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/transport"
)

// DownloadIPPPrinterAttrs downloads IPP Printer Attributes from
// the provided endpoints (assuming they all are aliases of the same
// device).
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadIPPPrinterAttrs(ctx context.Context,
	endpoints []string) (err error) {

	for _, ep := range endpoints {
		log.Debug(ctx, "ipp: trying %q", ep)

		var u *url.URL
		u, err2 := transport.ParseAddr(ep, "ipp://localhost")
		if err2 != nil {
			if err == nil {
				err = err2
			}

			log.Debug(ctx, "ipp: %q: %s", ep, err2)
			continue
		}

		clnt := ipp.NewClient(u, nil)
		clnt.SetDecoderOptions(
			&ipp.DecoderOptions{KeepTrying: true},
		)

		attrs, err2 := clnt.GetPrinterAttributes(ctx,
			[]string{
				ipp.GetPrinterAttributesAll,
				ipp.GetPrinterAttributesMediaColDatabase,
			},
			"",
		)

		if err2 != nil {
			if err == nil {
				err = err2
			}

			log.Debug(ctx, "ipp: %q: %s", ep, err2)
			continue
		}

		model.SetIPPPrinterAttrs(attrs)
		return nil
	}

	return err
}

// DownloadESCLScannerCapabilities downloads eSCL Scanner Capabilities from
// the provided endpoints (assuming they all are aliases of the same
// device).
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadESCLScannerCapabilities(ctx context.Context,
	endpoints []string) (err error) {

	for _, ep := range endpoints {
		log.Debug(ctx, "escl: trying %q", ep)

		var u *url.URL
		u, err2 := transport.ParseAddr(ep, "")
		if err2 != nil {
			if err == nil {
				err = err2
			}

			log.Debug(ctx, "escl: %q: %s", ep, err2)
			continue
		}

		clnt := escl.NewClient(u, nil)
		caps, _, err2 := clnt.GetScannerCapabilities(ctx)

		if err2 != nil {
			if err == nil {
				err = err2
			}

			log.Debug(ctx, "escl: %q: %s", ep, err2)
			continue
		}

		model.SetESCLScanCaps(caps)
		return nil
	}

	return err
}

// DownloadWSDScannerCapabilities downloads WS-Scan scanner configuration.
// from the provided endpoints (assuming they all are aliases of the same
// device).
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadWSDScannerCapabilities(ctx context.Context,
	endpoints []string) (err error) {

	for _, ep := range endpoints {
		log.Debug(ctx, "wsscan: trying %q", ep)

		var u *url.URL
		u, err2 := transport.ParseAddr(ep, "")
		if err2 != nil {
			if err == nil {
				err = err2
			}

			log.Debug(ctx, "wsscan: %q: %s", ep, err2)
			continue
		}

		clnt := wsscan.NewClient(u, nil)
		caps, err2 := clnt.GetScannerElements(
			ctx,
			wsscan.ScannerElemDescription,
			wsscan.ScannerElemConfiguration,
			wsscan.ScannerElemDefaultScanTicket,
		)

		if err2 != nil {
			if err == nil {
				err = err2
			}

			log.Debug(ctx, "wsscan: %q: %s", ep, err2)
			continue
		}

		model.SetWSDScanCaps(caps)
		return nil
	}

	return err
}
