// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Discovery backend

package discovery

import "context"

// Backend scans/monitors its search realm and reports discovered
// devices by sending series of [Event] into the provided [Eventqueue].
//
// The following model of operation is assumed:
//   - The search realm, where Backend operates (local network, for example)
//     contains some connected devices.
//   - Each device contains one or more units, and each unit may be
//     either print unit or scan unit
//   - Device may expose multiple interfaces to the same physical
//     unit. For example, printer may support multiple protocols (say,
//     IPP and LPD). Different interfaces to the same physical unit
//     needs to be reported as different units.
//   - Even the same interface may be visible to Backend as one or
//     more distinct units. For example, the same printer may be visible
//     via WiFi and via the Ethernet connection. This is up to Backend,
//     either to merge these "virtual units" together, by reporting
//     them as a single unit with combined endpoints, or report them
//     separately. At the later case, Backend should use UnitID.SubRealm
//     to distinguish between these virtual units.
type Backend interface {
	// Name returns backend name.
	Name() string

	// DeviceID returns a subset of UnitID, that uniquely
	// identifies a physical device.
	DeviceID(UnitID) UnitID

	// Open returns a new [Source]
	Open(ctx context.Context, queue *Eventqueue) (Source, error)
}

// Source generates a stream of discovery [Event]s.
type Source interface {
	// Close closes the Source and releases resources it holds
	Close()
}
