// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Conversion between Go and Python DNS-SD information

package modeling

import (
	"slices"
	"sort"
	"strings"

	"github.com/OpenPrinting/go-mfp/cpython"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/util/generic"
)

// DNSSDDevice represents DNS-SD advertising of the device
type DNSSDDevice struct {
	Instance string         // "E.g., "Kyocera ECOSYS M2040dn"
	UUID     string         // Device UUID
	Services []DNSSDService // Device's services
}

// DNSSDService represents DNS-SD service
type DNSSDService struct {
	Types []string // E.g. ["_ipp._tcp", "_ipps.tcp"]
	TXT   []string // E.g. ["txtvers=1", "ty=Kyocera ECOSYS M2040dn",...]
}

// merge merges two services, assuming they are duplicated
// services.
//
// Technically, this duplication comes from the fact, that
// some DNS-SD services may represent more that a single
// unit, in the discovery.Device terms. For example, IPP printer
// and faxout service will be represented by the single service
// at the DNS-SD side, while discovery will return two separate units.
func (svc *DNSSDService) merge(svc2 *DNSSDService) {
	// Merge Types
	for _, t := range svc2.Types {
		if !slices.Contains(svc.Types, t) {
			svc.Types = append(svc.Types, t)
		}
	}

	slices.Sort(svc.Types)
}

// dnssdExport converts DNS-SD information from discovery.Device
// into the Python object.
func dnssdExport(py *cpython.Python,
	discovered *discovery.Device) *cpython.Object {
	dev := dnssdFromDiscoveryDevice(discovered)
	return dnssdExportDevice(py, dev)
}

// dnssdExportDevice converts DNS-SD information from a DNSSDDevice
// into the Python object.
func dnssdExportDevice(py *cpython.Python, dev DNSSDDevice) *cpython.Object {
	obj := py.Eval("dnssd.Device()")

	err := obj.Set("instance", dev.Instance)
	if err != nil {
		return py.NewError(err)
	}

	services := []*cpython.Object{}
	for _, svc := range dev.Services {
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
func dnssdExportService(py *cpython.Python, svc DNSSDService) *cpython.Object {
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

// dnssdFromDiscoveryDevice builds DNSSDDevice from discovery.Device.
func dnssdFromDiscoveryDevice(discovered *discovery.Device) DNSSDDevice {
	// Services by DNS-SD type
	svcmap := make(map[string]*DNSSDService)

	// Gather print units
	for _, un := range discovered.PrintUnits {
		types := dnssdDiscoveredServiceTypes(un.Proto, un.Endpoints)
		if len(types) == 0 {
			// Not a DNS-SD unit
			continue
		}

		for _, t := range types {
			svc := &DNSSDService{Types: types, TXT: un.TXT}
			if found := svcmap[t]; found != nil {
				found.merge(svc)
				svc = found
			}
			for _, t := range svc.Types {
				svcmap[t] = svc
			}
		}
	}

	// Gather scan units
	for _, un := range discovered.ScanUnits {
		types := dnssdDiscoveredServiceTypes(un.Proto, un.Endpoints)
		if len(types) == 0 {
			// Not a DNS-SD unit
			continue
		}

		for _, t := range types {
			svc := &DNSSDService{Types: types, TXT: un.TXT}
			if found := svcmap[t]; found != nil {
				found.merge(svc)
				svc = found
			}
			for _, t := range svc.Types {
				svcmap[t] = svc
			}
		}
	}

	// Gather faxout units
	for _, un := range discovered.FaxoutUnits {
		types := dnssdDiscoveredServiceTypes(un.Proto, un.Endpoints)
		if len(types) == 0 {
			// Not a DNS-SD unit
			continue
		}

		for _, t := range types {
			svc := &DNSSDService{Types: types, TXT: un.TXT}
			if found := svcmap[t]; found != nil {
				found.merge(svc)
				svc = found
			}
			for _, t := range svc.Types {
				svcmap[t] = svc
			}
		}
	}

	// Build slice of services. Avoid duplicates.
	services := make([]DNSSDService, 0, len(svcmap))
	seen := generic.NewSet[*DNSSDService]()
	for _, svc := range svcmap {
		if seen.TestAndAdd(svc) {
			services = append(services, *svc)
		}
	}

	// Sort services by the first type name, to be deterministic
	sort.Slice(services, func(i, j int) bool {
		return services[i].Types[0] < services[j].Types[0]
	})

	dnssddev := DNSSDDevice{
		Instance: discovered.DNSSDName,
		UUID:     discovered.DNSSDUUID.String(),
		Services: services,
	}

	return dnssddev
}

// dnssdDiscoveredServiceTypes returns DNS-SD service types (e.g.,
// ["_ipp.tcp", "_ipps._tcp"], based on discovery.ServiceProto
// and endpoints.
//
// Endpoints are required here to distinguish IPP vs IPPS and similar.
func dnssdDiscoveredServiceTypes(
	proto discovery.ServiceProto, endpoints []string) []string {

	types := []string{}

	switch proto {
	case discovery.ServiceIPP:
		hasIPP := false
		hasIPPS := false

		for _, ep := range endpoints {
			if strings.HasPrefix(ep, "ipp:") {
				hasIPP = true
			}

			if strings.HasPrefix(ep, "ipps:") {
				hasIPPS = true
			}

			if hasIPP && hasIPPS {
				break
			}
		}

		if hasIPP {
			types = append(types, "_ipp._tcp")
		}

		if hasIPPS {
			types = append(types, "_ipps._tcp")
		}

	case discovery.ServiceESCL:
		hasESCL := false
		hasESCLS := false

		for _, ep := range endpoints {
			if strings.HasPrefix(ep, "http:") {
				hasESCL = true
			}

			if strings.HasPrefix(ep, "https:") {
				hasESCLS = true
			}

			if hasESCL && hasESCLS {
				break
			}
		}

		if hasESCL {
			types = append(types, "_uscan._tcp")
		}

		if hasESCLS {
			types = append(types, "_uscans._tcp")
		}

	case discovery.ServiceLPD:
		types = []string{"_printer._tcp"}

	case discovery.ServiceAppSocket:
		types = []string{"_pdl-datastream._tcp"}
	}

	return types
}
