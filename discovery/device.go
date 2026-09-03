// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Common device information

package discovery

import (
	"bytes"
	"fmt"
	"net/netip"
	"sort"
	"strings"

	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// Device consist of the multiple functional units. There are
// three types of units:
//   - [PrintUnit], for printing
//   - [ScanUnit], for scanning
//   - [FaxoutUnit], for sending faxes.
//
// Multiple units of each type may exist, and depending on the device,
// they may have different parameters.
//
// Each unit has its unique [UnitID], the combination of parameters,
// that uniquely identifies the unit.
type Device struct {
	// Device metadata
	MakeModel      string    // Device make and model
	Location       string    // E.g., "2nd Floor Computer Lab"
	DNSSDName      string    // DNS-SD name, "" if none
	DNSSDUUID      uuid.UUID // DNS-SD UUID, uuid.NilUUID if n/a
	WSDUUID        uuid.UUID // WSD UUID
	PrintAdminURL  string    // Admin URL for printer
	ScanAdminURL   string    // Admin URL for scanner
	FaxoutAdminURL string    // Admin URL for faxout
	IconURL        string    // Device icon URL

	// PPDManufacturer and PPDModel are matched against Manufacturer
	// and Model parameters in the PPD file when searching for the
	// appropriate driver for the legacy printer.
	//
	// Please notice, it is not necessary true that MakeModel
	// is the exact concatenation of these two strings.
	PPDManufacturer string // Manufacturer name
	PPDModel        string // Model name

	// USBSerial may be available for the ipp-usb devices too.
	USBSerial string // USB serial number, "" if n/a
	USBHWID   string // USB hardware ID, "" if n/a

	// Connectivity
	Addrs []netip.Addr // Device's IP addresses

	// Device units
	PrintUnits  []PrintUnit  // Print units
	ScanUnits   []ScanUnit   // Scan units
	FaxoutUnits []FaxoutUnit // Faxout units
}

// MarshalLog formats [Device] for logging.
// It implements the [log.Marshaler] interface.
func (dev *Device) MarshalLog() []byte {
	buf := bytes.Buffer{}

	fmt.Fprintf(&buf, "%q (%s)\n", dev.DNSSDName, dev.DNSSDUUID)
	for _, un := range dev.PrintUnits {
		fmt.Fprintf(&buf, "  %s printer:\n", un.Proto)
		for _, ep := range un.Endpoints {
			fmt.Fprintf(&buf, "    %s\n", ep)
		}
	}
	for _, un := range dev.ScanUnits {
		fmt.Fprintf(&buf, "  %s scanner:\n", un.Proto)
		for _, ep := range un.Endpoints {
			fmt.Fprintf(&buf, "    %s\n", ep)
		}
	}
	for _, un := range dev.FaxoutUnits {
		fmt.Fprintf(&buf, "  %s faxout:\n", un.Proto)
		for _, ep := range un.Endpoints {
			fmt.Fprintf(&buf, "    %s\n", ep)
		}
	}

	return buf.Bytes()
}

// less compares devices for sorting
func (dev *Device) less(dev2 *Device) bool {
	// Put devices with DNSSDName first.
	// Compare by DNSSDName
	switch {
	case dev.DNSSDName != "" && dev2.DNSSDName == "":
		return true
	case dev.DNSSDName == "" && dev2.DNSSDName != "":
		return false
	case dev.DNSSDName != "" && dev2.DNSSDName != "":
		return dev.DNSSDName < dev2.DNSSDName
	}

	// Now, if MakeModel is available, compare
	// by MakeModel + DNSSDUUID
	switch {
	case dev.MakeModel != "" && dev2.MakeModel == "":
		return true
	case dev.MakeModel == "" && dev2.MakeModel != "":
		return false
	case dev.MakeModel != "" && dev2.MakeModel != "":
		return dev.DNSSDName+dev.DNSSDUUID.String() <
			dev2.DNSSDName+dev2.DNSSDUUID.String()
	}

	// Compare by DNSSDUUID, if available
	switch {
	case !dev.DNSSDUUID.IsZero() && dev2.DNSSDUUID.IsZero():
		return true
	case dev.DNSSDUUID.IsZero() && !dev2.DNSSDUUID.IsZero():
		return false
	case !dev.DNSSDUUID.IsZero() && !dev2.DNSSDUUID.IsZero():
		return dev.DNSSDUUID.Less(dev2.DNSSDUUID)
	}

	// Compare by USBSerial, if available
	switch {
	case dev.USBSerial != "" && dev2.USBSerial == "":
		return true
	case dev.USBSerial == "" && dev2.USBSerial != "":
		return false
	case dev.USBSerial != "" && dev2.USBSerial != "":
		return dev.USBSerial < dev2.USBSerial
	}

	// Give up
	return false
}

// DNSSD returns DNS-SD portion of the discovery information
// for the discovered [Device].
func (dev *Device) DNSSD() *DNSSDDevice {
	// Services by DNS-SD type
	svcmap := make(map[string]*DNSSDService)

	// Gather print units
	for _, un := range dev.PrintUnits {
		types := dev.dnssdServiceTypes(un.Proto, un.Endpoints)
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
	for _, un := range dev.ScanUnits {
		types := dev.dnssdServiceTypes(un.Proto, un.Endpoints)
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
	for _, un := range dev.FaxoutUnits {
		types := dev.dnssdServiceTypes(un.Proto, un.Endpoints)
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

	dnssddev := &DNSSDDevice{
		Instance: dev.DNSSDName,
		UUID:     dev.DNSSDUUID,
		Services: services,
	}

	return dnssddev
}

// dnssdServiceTypes returns DNS-SD service types (e.g.,
// ["_ipp.tcp", "_ipps._tcp"], based on discovery.ServiceProto
// and endpoints.
//
// Endpoints are required here to distinguish IPP vs IPPS and similar.
func (dev *Device) dnssdServiceTypes(
	proto ServiceProto, endpoints []string) []string {

	types := []string{}

	switch proto {
	case ServiceIPP:
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

	case ServiceESCL:
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

	case ServiceLPD:
		types = []string{"_printer._tcp"}

	case ServiceAppSocket:
		types = []string{"_pdl-datastream._tcp"}
	}

	return types
}
