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
func (dev Device) MarshalLog() []byte {
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

// device is the internal representation of the Device
type device struct {
	backend Backend      // Backend that found the device (nil if not known)
	uuid    uuid.UUID    // Device's UUID
	units   []unit       // Device's units
	addrs   []netip.Addr // Device's IP addresses
}

// Export exports device as Device
func (dev device) Export() Device {
	out := Device{Addrs: dev.addrs}

	// Classify units
	var ippPrinters []*unit
	var lpdPrinters []*unit
	var appsockPrinters []*unit
	var wsdPrinters []*unit
	var usbPrinters []*unit
	var ippScanners []*unit
	var esclSanners []*unit
	var wsdScanners []*unit
	var ippFaxes []*unit

	for i := range dev.units {
		un := &dev.units[i]
		switch un.ID.SvcType {
		case ServicePrinter:
			switch un.ID.SvcProto {
			case ServiceIPP:
				ippPrinters = append(ippPrinters, un)
			case ServiceLPD:
				lpdPrinters = append(lpdPrinters, un)
			case ServiceAppSocket:
				appsockPrinters = append(appsockPrinters, un)
			case ServiceWSD:
				wsdPrinters = append(wsdPrinters, un)
			case ServiceUSB:
				usbPrinters = append(usbPrinters, un)
			}

		case ServiceScanner:
			switch un.ID.SvcProto {
			case ServiceIPP:
				ippScanners = append(ippScanners, un)
			case ServiceESCL:
				esclSanners = append(esclSanners, un)
			case ServiceWSD:
				wsdScanners = append(wsdScanners, un)
			}

		case ServiceFaxout:
			switch un.ID.SvcProto {
			case ServiceIPP:
				ippFaxes = append(ippFaxes, un)
			}
		}
	}

	// Merge by classes
	printUnits := generic.ConcatSlices(
		ippPrinters,
		lpdPrinters,
		appsockPrinters,
		wsdPrinters,
		usbPrinters,
	)

	scanUnits := generic.ConcatSlices(
		ippScanners,
		esclSanners,
		wsdScanners,
	)

	faxoutUnits := ippFaxes

	// Convert units to external representation and save to device.
	for _, un := range printUnits {
		exp := un.Export().(PrintUnit)
		out.PrintUnits = append(out.PrintUnits, exp)
	}

	for _, un := range scanUnits {
		exp := un.Export().(ScanUnit)
		out.ScanUnits = append(out.ScanUnits, exp)
	}

	for _, un := range faxoutUnits {
		exp := un.Export().(FaxoutUnit)
		out.FaxoutUnits = append(out.FaxoutUnits, exp)
	}

	// Extract metadata
	dnssdUnits := generic.ConcatSlices(
		ippPrinters,
		lpdPrinters,
		appsockPrinters,
		ippScanners,
		esclSanners,
		ippFaxes,
	)

	allUnits := generic.ConcatSlices(
		ippPrinters,
		lpdPrinters,
		appsockPrinters,
		wsdPrinters,
		usbPrinters,
		ippScanners,
		esclSanners,
		wsdScanners,
		ippFaxes,
	)

	for _, un := range dnssdUnits {
		if un.MakeModel != "" {
			out.MakeModel = un.MakeModel
			break
		}
	}

	for _, un := range allUnits {
		if un.PPDManufacturer != "" && un.PPDModel != "" {
			out.PPDManufacturer = un.PPDManufacturer
			out.PPDModel = un.PPDModel
			break
		}
	}

	for _, un := range dnssdUnits {
		if un.ID.DNSSDName != "" && un.ID.UUID != uuid.NilUUID {
			out.DNSSDName = un.ID.DNSSDName
			out.DNSSDUUID = un.ID.UUID
			break
		}
	}

	for _, un := range dnssdUnits {
		switch un.ID.SvcType {
		case ServicePrinter:
			if out.PrintAdminURL == "" {
				out.PrintAdminURL = un.AdminURL
			}
		case ServiceScanner:
			if out.ScanAdminURL == "" {
				out.ScanAdminURL = un.AdminURL
			}
		case ServiceFaxout:
			if out.FaxoutAdminURL == "" {
				out.FaxoutAdminURL = un.AdminURL
			}
		}
	}

	for _, un := range generic.ConcatSlices(printUnits, scanUnits, faxoutUnits) {
		if out.Location == "" && un.Location != "" {
			out.Location = un.Location
		}

		if out.IconURL == "" && un.IconURL != "" {
			out.IconURL = un.IconURL
		}

		if out.Location != "" && out.IconURL != "" {
			break
		}
	}

	for _, un := range allUnits {
		if out.MakeModel == "" && un.MakeModel != "" {
			out.MakeModel = un.MakeModel
		}

		if out.DNSSDUUID == uuid.NilUUID && un.ID.UUID != uuid.NilUUID {
			out.DNSSDUUID = un.ID.UUID
		}

		if un.ID.USBSerial != "" {
			out.USBSerial = un.ID.USBSerial
		}

		if un.ID.USBHWID != "" {
			out.USBHWID = un.ID.USBHWID
		}
	}

	return out
}
