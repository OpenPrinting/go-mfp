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
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// PrintUnit represents a print unit.
type PrintUnit struct {
	Proto     ServiceProto      // Printing protocol
	Params    PrinterParameters // Printer parameters
	Endpoints []string          // URLs of printer endpoints
}

// ScanUnit represents a scan unit.
type ScanUnit struct {
	Proto     ServiceProto      // Scanning protocol
	Params    ScannerParameters // Scanner parameters
	Endpoints []string          // URLs of printer endpoints
}

// FaxoutUnit represents a fax unit.
type FaxoutUnit struct {
	Proto     ServiceProto      // Faxing protocol
	Params    PrinterParameters // Printer parameters
	Endpoints []string          // URLs of printer endpoints
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
	Endpoints       []string     // Unit endpoints
	Addrs           []netip.Addr // Addresses that unit use
}

// Merge merges two units
func (un *unit) Merge(un2 unit) {
	un.Endpoints = endpointsMerge(un.Endpoints, un2.Endpoints)
	un.Addrs = addrsMerge(un.Addrs, un2.Addrs)
}

// Export exports unit ad PrintUnit, ScanUnit or FaxoutUnit
func (un unit) Export() any {
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
			}
		case ServiceFaxout:
			return FaxoutUnit{
				Proto:     un.ID.SvcProto,
				Params:    params,
				Endpoints: un.Endpoints,
			}
		}

	case ScannerParameters:
		return ScanUnit{
			Proto:     un.ID.SvcProto,
			Params:    params,
			Endpoints: un.Endpoints,
		}
	}

	return nil
}

// unitsToLogRecord writes a slice of units into the log.Record
func unitsToLogRecord(rec *log.Record, units []unit) {
	// Group units by DNSSDName and UUID
	type id struct {
		name string
		uuid uuid.UUID
	}

	grouped := make(map[id][]*unit)
	for i := range units {
		un := &units[i]
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

// SameDevice reports if two [UnitID]s belong to the same device.
func (id UnitID) SameDevice(id2 UnitID) bool {
	if id.UUID == id2.UUID {
		return true
	}

	if id.DNSSDName == id2.DNSSDName &&
		id.Backend == id2.Backend &&
		id.Zone == id2.Zone {
		return true
	}

	return false
}

// SameService reports if two [UnitID]s belong to the same service of
// the same device.
func (id UnitID) SameService(id2 UnitID) bool {
	return id.SvcType == id2.SvcType && id.SameDevice(id2)
}

// SameUnit reports if two [UnitID]s belong to the same unit of
// the same device.
func (id UnitID) SameUnit(id2 UnitID) bool {
	return id.Queue == id2.Queue && id.SameService(id2)
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
