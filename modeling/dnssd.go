// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Conversion between Go and Python DNS-SD information

package modeling

import (
	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/discovery"
)

// dnssdExport converts DNS-SD information from discovery.Device
// into the Python object.
func dnssdExport(py *cpython.Python,
	dnssddev *discovery.DNSSDDevice) *cpython.Object {

	obj := py.Eval("dnssd.Device()")

	err := obj.Set("instance", dnssddev.Instance)
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
