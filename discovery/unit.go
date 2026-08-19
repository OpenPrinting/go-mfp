// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device units

package discovery

import (
	"fmt"
	"net/netip"
	"sort"
	"strings"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// PrintUnit represents a print unit.
type PrintUnit struct {
	Proto     ServiceProto      // Printing protocol
	Params    PrinterParameters // Printer parameters
	Endpoints []string          // URLs of printer endpoints
	TXT       []string          // DND-SD TXT record, if available
}

// ScanUnit represents a scan unit.
type ScanUnit struct {
	Proto     ServiceProto      // Scanning protocol
	Params    ScannerParameters // Scanner parameters
	Endpoints []string          // URLs of printer endpoints
	TXT       []string          // DND-SD TXT record, if available
}

// FaxoutUnit represents a fax unit.
type FaxoutUnit struct {
	Proto     ServiceProto      // Faxing protocol
	Params    PrinterParameters // Printer parameters
	Endpoints []string          // URLs of printer endpoints
	TXT       []string          // DND-SD TXT record, if available
}

// unit is the internal representation of the PrintUnit, ScanUnit
// or FaxoutUnit
type unit struct {
	ID              UnitID       // Unit identity
	MakeModel       string       // Manufacturer + Model
	Location        string       // E.g., "2nd Floor Computer Lab"
	AdminURL        string       // Device administration URL
	IconURL         string       // Device icon URL
	PPDManufacturer string       // I.e., "Hewlett Packard" or "Canon"
	PPDModel        string       // Model name
	Params          any          // PrinterParameters or ScannerParameters
	TXT             []string     // DND-SD TXT record, if available
	Endpoints       []string     // Unit endpoints
	Addrs           []netip.Addr // Addresses that unit use
}

// Merge merges two units
func (un *unit) Merge(un2 *unit) {
	un.Endpoints = endpointsMerge(un.Endpoints, un2.Endpoints)
	un.Addrs = addrsMerge(un.Addrs, un2.Addrs)

	txtSeen := generic.NewSet[string]()
	for _, t := range un.TXT {
		txtSeen.Add(t)
	}

	for _, t := range un2.TXT {
		if txtSeen.TestAndAdd(t) {
			un.TXT = append(un.TXT, t)
		}
	}
}

// Export exports unit ad PrintUnit, ScanUnit or FaxoutUnit
func (un *unit) Export() any {
	switch params := un.Params.(type) {
	case PrinterParameters:
		// PrinterParameters can be used either with PrintUnit
		// or FaxoutUnit
		switch un.ID.SvcType {
		case ServicePrinter:
			return PrintUnit{
				Proto:     un.ID.SvcProto,
				Params:    params,
				Endpoints: un.Endpoints,
				TXT:       un.TXT,
			}
		case ServiceFaxout:
			return FaxoutUnit{
				Proto:     un.ID.SvcProto,
				Params:    params,
				Endpoints: un.Endpoints,
				TXT:       un.TXT,
			}
		}

	case ScannerParameters:
		return ScanUnit{
			Proto:     un.ID.SvcProto,
			Params:    params,
			Endpoints: un.Endpoints,
			TXT:       un.TXT,
		}
	}

	return nil
}

// unitgroup is the collection of units, most likely to belong
// to the same device
type unitgroup []*unit

// Add adds unit to the group
func (grp *unitgroup) Add(un *unit) {
	*grp = append(*grp, un)
}

// Addrs returns IP addresses, found in the group.
func (grp unitgroup) Addrs() []netip.Addr {
	var addrs []netip.Addr
	for _, un := range grp {
		addrs = addrsMerge(addrs, un.Addrs)
	}
	return addrs
}

// UUIDs returns collection of UUIDs, owned by members
// of the group. uuid.NilUUID not included.
func (grp unitgroup) UUIDs() []uuid.UUID {
	uuids := make([]uuid.UUID, 0, len(grp))
	for _, un := range grp {
		uu := un.ID.UUID
		if uu != uuid.NilUUID {
			uuids, _ = uuidsAdd(uuids, uu)
		}
	}
	return uuids
}

// CanMergeByUUID decides if two unitgroups can be merged
// by UUID.
func (grp unitgroup) CanMergeByUUID(grp2 unitgroup) bool {
	// Don't merge groups, if they don't have common addresses.
	// Otherwise, IPP over USB printer may steal WSD units from
	// the same printer, reachable over network.
	if !addrsOverlap(grp.Addrs(), grp2.Addrs()) {
		return false
	}

	ok := false
	for _, un1 := range grp {
		// Ignore units without UUID
		if un1.ID.UUID == uuid.NilUUID {
			continue
		}

		for _, un2 := range grp2 {
			// Ignore units without UUID
			if un2.ID.UUID == uuid.NilUUID {
				continue
			}

			// If UUIDs match:
			//   - Same backend: merge not possible (backend
			//     distinguishes devices)
			//   - Different backends: merge possible, unless
			//     blocked by the previous condition
			if un1.ID.UUID == un2.ID.UUID {
				if un1.ID.Backend == un2.ID.Backend {
					return false
				}

				ok = true
			}
		}
	}

	return ok
}

// Export exports unitgroup as Device
func (grp unitgroup) Export() *Device {
	// Gather IP addresses
	out := &Device{}
	for _, un := range grp {
		out.Addrs = addrsMerge(out.Addrs, un.Addrs)
	}

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

	for _, un := range grp {
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

	for _, un := range append(wsdScanners, wsdPrinters...) {
		if un.ID.UUID != uuid.NilUUID {
			out.WSDUUID = un.ID.UUID
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

// UnitID contains combination of parameters that identifies a device.
//
// Please note, depending on a discovery protocol being used, not
// all the fields of the following structure will have any sense.
//
// Note also, that device UUID is not necessary the same between
// protocols. Some Canon devices known to use different UUID for
// DNS-SD and WS-Discovery.
//
// The intended fields usage is the following:
//
//	DeviceName - realm-unique device name, in the DNS-SD sense.
//	             E.g., "Kyocera ECOSYS M2040dn",
//	UUID       - device UUID
//	Queue      - Job queue name for units with logical sub-units,
//	             like LPD server with multiple queues
//	Backend    - backend that found the unit. Different backends
//	             are treated as independent namespaces.
//	Zone       - allows backend to further divide its namespace
//	             (for example, to split it between network interfaces)
//	Variant    - used to distinguish between logically equivalent
//	             variants of discovered units, that backend sees as
//	             independent instances (for example IP4/IP6, HTTP/HTTPS)
//	SvcType    - service type, printer/scanner/faxout
//	SvcProto   - service protocol, i.e., IPP, LPD, eSCL etc
//	Serial     - device serial number, if appropriate (i.e., for USB)
type UnitID struct {
	DNSSDName string       // DNS-SD name, "" if not available
	UUID      uuid.UUID    // uuid.NilUUID if not available
	Queue     string       // Logical unit within a device
	Backend   Backend      // Backend that discovered the unit
	Zone      string       // Namespace zone within the Realm
	Variant   string       // Finding variant of the same unit
	SvcType   ServiceType  // Service type
	SvcProto  ServiceProto // Service protocol
	USBSerial string       // "" if not available
	USBHWID   string       // "" if not available
}

// MarshalLog dumps [UnitID] as text, for [log.Object].
// It implements [log.Marshaler].
func (id UnitID) MarshalLog() []byte {
	var line string
	lines := make([]string, 0, 6)

	if id.DNSSDName != "" {
		line = fmt.Sprintf("DNSSDName: %q", id.DNSSDName)
		lines = append(lines, line)
	}
	if id.UUID != uuid.NilUUID {
		line = fmt.Sprintf("UUID:      %s", id.UUID)
		lines = append(lines, line)
	}
	if id.Queue != "" {
		line = fmt.Sprintf("Queue:     %q", id.Queue)
		lines = append(lines, line)
	}

	if id.Backend != nil {
		line = fmt.Sprintf("Backend:   %s", id.Backend.Name())
		lines = append(lines, line)
	}

	if id.Zone != "" {
		line = fmt.Sprintf("Zone:      %s", id.Zone)
		lines = append(lines, line)
	}

	if id.Variant != "" {
		line = fmt.Sprintf("Variant:   %s", id.Variant)
		lines = append(lines, line)
	}

	line = fmt.Sprintf("Service:   %s %s", id.SvcProto, id.SvcType)
	lines = append(lines, line)

	if id.USBSerial != "" {
		line = fmt.Sprintf("Serial:    %s", id.USBSerial)
		lines = append(lines, line)
	}

	if id.USBHWID != "" {
		line = fmt.Sprintf("HWID:      %s", id.USBHWID)
		lines = append(lines, line)
	}

	return []byte(strings.Join(lines, "\n"))
}

// unitsToLogRecord writes a slice of units into the log.Record
func unitsToLogRecord(rec *log.Record, units []*unit) {
	// Group units by DNSSDName and UUID
	type id struct {
		name string
		uuid uuid.UUID
	}

	grouped := make(map[id][]*unit)
	for _, un := range units {
		id := id{un.ID.DNSSDName, un.ID.UUID}
		grouped[id] = append(grouped[id], un)
	}

	// Sort groups
	sorted := make([][]*unit, 0, len(grouped))
	for _, grp := range grouped {
		sorted = append(sorted, grp)
	}

	sort.Slice(sorted, func(i, j int) bool {
		g1 := sorted[i]
		g2 := sorted[j]

		switch {
		case g1[0].ID.DNSSDName < g2[0].ID.DNSSDName:
			return true
		case g1[0].ID.DNSSDName > g2[0].ID.DNSSDName:
			return false
		}

		return g1[0].ID.UUID.Less(g2[0].ID.UUID)
	})

	// Write to log record
	for _, grp := range sorted {
		rec.Trace("  %q (%s)",
			grp[0].ID.DNSSDName,
			grp[0].ID.UUID)

		for _, un := range grp {
			rec.Trace("    %s %s: %s",
				un.ID.SvcProto, un.ID.SvcType, un.Endpoints)
		}
	}
}
