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
	"errors"
	"fmt"
	"net/url"

	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/discovery/dnssd"
	"github.com/OpenPrinting/go-mfp/discovery/wsdd"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/escl"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/proto/usbhost"
	"github.com/OpenPrinting/go-mfp/proto/wsscan"
	"github.com/OpenPrinting/go-mfp/transport"
)

// DownloadIPPPrinterAttrs downloads IPP Printer Attributes from
// the provided endpoints (assuming they all are aliases of the same
// device).
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadIPPPrinterAttrs(ctx context.Context,
	endpoints []string) error {

	attrs, err := model.fetchIPPPrinterAttrs(ctx, endpoints)
	if err == nil {
		model.SetIPPPrinterAttrs(attrs)
	}

	return err
}

// DownloadESCLScannerCapabilities downloads eSCL Scanner Capabilities from
// the provided endpoints (assuming they all are aliases of the same
// device).
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadESCLScannerCapabilities(ctx context.Context,
	endpoints []string) error {

	caps, err := model.fetchESCLScannerCapabilities(ctx, endpoints)
	if err == nil {
		model.SetESCLScanCaps(caps)
	}

	return err
}

// DownloadWSDScannerCapabilities downloads WS-Scan scanner configuration.
// from the provided endpoints (assuming they all are aliases of the same
// device).
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadWSDScannerCapabilities(ctx context.Context,
	endpoints []string) error {

	caps, err := model.fetchWSDScannerCapabilities(ctx, endpoints)
	if err == nil {
		model.SetWSDScanCaps(caps)
	}
	return err
}

// DownloadByDNSSDName locates device by its DNS-SD name and
// downloads its printer/scanner capabilities.
func (model *Model) DownloadByDNSSDName(ctx context.Context,
	name string) error {
	// Prepare discovery.Client
	clnt := discovery.NewClient(ctx)
	defer clnt.Close()

	backend, err := dnssd.NewBackend(ctx, "", 0)
	if err != nil {
		return err
	}

	defer backend.Close()
	clnt.AddBackend(backend)

	backend, err = wsdd.NewBackend(ctx)
	if err != nil {
		return err
	}

	defer backend.Close()
	clnt.AddBackend(backend)

	// Search for device
	dev, err := clnt.GetByDNSSD(ctx, name, discovery.ModeNormal)
	if err != nil {
		return err
	}

	// Gather endpoints
	var endpointsIPP, endpointsESCL, endpointsWSD []string

	for _, unit := range dev.PrintUnits {
		if unit.Proto == discovery.ServiceIPP {
			endpointsIPP = append(endpointsIPP, unit.Endpoints...)
		}
	}

	for _, unit := range dev.ScanUnits {
		switch unit.Proto {
		case discovery.ServiceESCL:
			endpointsESCL = append(endpointsESCL, unit.Endpoints...)
		case discovery.ServiceWSD:
			endpointsWSD = append(endpointsWSD, unit.Endpoints...)
		}
	}

	if endpointsIPP == nil && endpointsESCL == nil && endpointsWSD == nil {
		err := errors.New("no eSCL/IPP/WSD endpoints discovered")
		return err
	}

	// Fetch configuration
	attrsIPP, err := model.fetchIPPPrinterAttrs(ctx, endpointsIPP)
	if err != nil {
		return err
	}

	capsESCL, err := model.fetchESCLScannerCapabilities(ctx, endpointsESCL)
	if err != nil {
		return err
	}

	capsWSD, err := model.fetchWSDScannerCapabilities(ctx, endpointsWSD)
	if err != nil {
		return err
	}

	// Update the Model
	model.SetIPPPrinterAttrs(attrsIPP)
	model.SetESCLScanCaps(capsESCL)
	model.SetWSDScanCaps(capsWSD)
	model.discovered = dev

	return nil
}

// DownloadUSBDeviceDescriptor downloads USB device configuration.
//
// Upon successful completion, Model is updated.
func (model *Model) DownloadUSBDeviceDescriptor(ctx context.Context,
	serial string) error {

	devices, err := usbhost.ListDevices(true)
	if err != nil {
		return err
	}

	for _, dev := range devices {
		if dev.Desc.ISerialNumber == serial {
			model.SetUSBDeviceDescriptor(&dev.Desc)
			return nil
		}
	}

	return fmt.Errorf("%s: device not found", serial)
}

// fetchIPPPrinterAttrs does a hard work of Model.DownloadIPPPrinterAttrs,
// but doesn't update the Model on success.
func (model *Model) fetchIPPPrinterAttrs(ctx context.Context,
	endpoints []string) (*ipp.PrinterAttributes, error) {

	var err error
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

		return attrs, nil
	}

	return nil, err
}

// fetchESCLScannerCapabilities does a hard work of
// Model.DownloadESCLScannerCapabilities, but doesn't update the
// Model on success.
func (model *Model) fetchESCLScannerCapabilities(ctx context.Context,
	endpoints []string) (*escl.ScannerCapabilities, error) {

	var err error
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

		return caps, nil
	}

	return nil, err
}

// fetchWSDScannerCapabilities does a hard work of
// Model.DownloadWSDScannerCapabilities, but doesn't update the
// Model on success.
func (model *Model) fetchWSDScannerCapabilities(ctx context.Context,
	endpoints []string) (*wsscan.GetScannerElementsResponse, error) {

	var err error
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

		return caps, nil
	}

	return nil, err
}
