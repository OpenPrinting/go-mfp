// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Conversion between Go and Python DNS-SD information

package modeling

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// dnssdLoad decodes DNS-SD part of model. The model file assumed to be
// preloaded into the Model's Python interpreter (model.py).
func (model *Model) dnssdLoad() error {
	// Obtain Python object for "dnssd.device"
	obj := model.py.Eval("dnssd.device")

	if err := obj.Err(); err != nil {
		err = fmt.Errorf("dnssd.device: %w", err)
		return err
	}

	if obj.IsNone() {
		return nil
	}

	// Decode dnssd.device
	dnssddev, err := dnssdImport(obj)
	if err != nil {
		err = fmt.Errorf("dnssd.device: %w", err)
		return err
	}

	model.dnssd = dnssddev
	return nil
}

// dnssdExport converts DNS-SD information from discovery.Device
// into the Python object.
func dnssdExport(py *cpython.Python,
	dnssddev *discovery.DNSSDDevice) *cpython.Object {

	obj := py.Eval("dnssd.Device()")

	err := obj.Set("instance", dnssddev.Instance)
	if err != nil {
		return py.NewError(err)
	}

	err = obj.Set("UUID", py.Get("UUID").Call(dnssddev.UUID.String()))
	if err != nil {
		return py.NewError(err)
	}

	services := []*cpython.Object{}
	for _, svc := range dnssddev.Services {
		svcobj := dnssdExportService(py, svc)
		switch {
		case svcobj == nil:
		case svcobj.Err() != nil:
			return svcobj
		default:
			services = append(services, svcobj)
		}
	}

	err = obj.Set("services", services)
	if err != nil {
		return py.NewError(err)
	}

	return obj
}

// dnssdExportService converts DNSSDService into the Python object
func dnssdExportService(py *cpython.Python, svc discovery.DNSSDService) *cpython.Object {
	// Create service object
	svcobj := py.Eval("dnssd.Service()")

	err := svcobj.Set("types", svc.Types)
	if err == nil {
		err = svcobj.Set("TXT", svc.TXT)
	}

	if err != nil {
		return py.NewError(err)
	}

	return svcobj
}

// dnssdImport converts DNS-SD information from Python object
// into the discovery.Device.
func dnssdImport(obj *cpython.Object) (*discovery.DNSSDDevice, error) {
	dnssddev := &discovery.DNSSDDevice{}
	var err error

	dnssddev.Instance, err = obj.Get("instance").Unicode()
	if err != nil {
		err = errImportWrap("Instance", err)
		return nil, err
	}

	s, err := obj.Get("UUID").Str()
	if err != nil {
		err = errImportWrap("UUID", err)
		return nil, err
	}

	dnssddev.UUID, err = uuid.Parse(s)
	if err != nil {
		err = fmt.Errorf("%q: invalid UUID", s)
		err = errImportWrap("UUID", err)
		return nil, err
	}

	svcobjects, err := obj.Get("services").Slice()
	if err != nil {
		err = errImportWrap("Services", err)
		return nil, err
	}

	dnssddev.Services = make([]discovery.DNSSDService, 0, len(svcobjects))
	for i, svcobj := range svcobjects {
		svc, err := dnssdImportService(svcobj)
		if err != nil {
			err = errImportWrap(fmt.Sprintf("%d", i), err)
			err = errImportWrap("Services", err)
			return nil, err
		}

		dnssddev.Services = append(dnssddev.Services, svc)
	}

	return dnssddev, nil
}

// dnssdImport converts Python object into discovery.DNSSDService.
func dnssdImportService(obj *cpython.Object) (discovery.DNSSDService, error) {
	svc := discovery.DNSSDService{}

	// Decode Types
	sliceobjects, err := obj.Get("types").Slice()
	if err != nil {
		err = errImportWrap("Types", err)
		return svc, err
	}

	svc.Types = make([]string, 0, len(sliceobjects))
	for i, item := range sliceobjects {
		s, err := item.Unicode()
		if err != nil {
			err = errImportWrap(fmt.Sprintf("%d", i), err)
			err = errImportWrap("Types", err)
			return svc, err
		}

		svc.Types = append(svc.Types, s)
	}

	// Decode TXT
	sliceobjects, err = obj.Get("TXT").Slice()
	if err != nil {
		err = errImportWrap("TXT", err)
		return svc, err
	}

	svc.TXT = make([]string, 0, len(sliceobjects))
	for i, item := range sliceobjects {
		s, err := item.Unicode()
		if err != nil {
			err = errImportWrap(fmt.Sprintf("%d", i), err)
			err = errImportWrap("TXT", err)
			return svc, err
		}

		svc.TXT = append(svc.TXT, s)
	}

	return svc, nil
}
