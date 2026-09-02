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
	"strings"

	"github.com/OpenPrinting/go-mfp/util/uuid"
)

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
	Zone      string       // Namespace zone (e.g., network interface)
	Variant   string       // Variant of the same unit (e.g., http/https)
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
