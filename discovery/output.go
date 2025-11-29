// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Output generation

package discovery

import (
	"sort"
	"time"

	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// output generates and manages the final discovery output from
// the internal representation of the discovered information,
// gathered in the cache
type output struct {
	devices []Device  // Cached output data
	ttl     time.Time // Cache valid until this time
}

// Cached returns the cached output data (created by latest output.Generate)
// It may return nil, if this information is not available.
func (out *output) Cached() []Device {
	if out.devices != nil && !out.ttl.After(time.Now()) {
		return out.devices
	}
	return nil
}

// Invalidate drops the cached output
func (out *output) Invalidate() {
	out.devices = nil
}

// Generate generates the discovery output from the discovery
// information, gathered in the cache.
func (out *output) Generate(ttl time.Time, units []unit) []Device {
	// Extract IP addresses
	out.genExtractIPAddresses(units)

	// Merge variants
	units = out.genMergeUnitVariants(units)

	// Classify units by DeviceName+UUID+Realm
	devices := out.genMergeDevicesByNameUUID(units)

	// Merge devices by UUID
	devices = out.genMergeDevicesByUUID(devices)

	// Generate final output, save and returns
	outdevs := make([]Device, len(devices))
	for i := range devices {
		outdevs[i] = devices[i].Export()
	}

	out.devices = outdevs
	out.ttl = ttl

	return outdevs
}

// genExtractIPAddresses extracts IP addresses from endpoints.
// It modifies slice of units in place.
func (out *output) genExtractIPAddresses(units []unit) {
	for i := range units {
		un := &units[i]
		un.Addrs = addrsFromEndpoints(un.Endpoints)
	}
}

// genMergeUnitVariants merges together units with distinct UnitID.Variant,
// but otherwise similar.
//
// We need it because if some variant (say, ip6) will disappear
// from the network, remaining variants still make the unit
// visible.
//
// Some units may come several times but in distinct variants
// (say ip4/ip6). This function merges them together, effectively
// remove duplicates.
func (out *output) genMergeUnitVariants(units []unit) []unit {
	scratchpad := make(map[UnitID]unit)
	for _, un := range units {
		un.ID.Variant = ""
		key := un.ID

		if prev, found := scratchpad[key]; found {
			// Keep the first found unit; merge endpoints
			prev.Merge(un)
			scratchpad[key] = prev
		} else {
			// Add new unit
			scratchpad[key] = un
		}
	}

	units = make([]unit, 0, len(scratchpad))
	for _, un := range scratchpad {
		units = append(units, un)
	}

	return units
}

// genMergeUnitSameFunction merges together units that differ
// in UnitID.Zone, but otherwise similar.
//
// It may happen, for example, when units from different
// network protocols appeared to belong to the same device.
func (out *output) genMergeUnitCrossZones(units []unit) []unit {
	scratchpad := make(map[UnitID]unit)

	for _, un := range units {
		un.ID.Zone = ""
		key := un.ID

		if prev, found := scratchpad[key]; found {
			// Keep the first found unit; merge endpoints
			prev.Merge(un)
			scratchpad[key] = prev
		} else {
			// Add new unit
			scratchpad[key] = un
		}
	}

	units = make([]unit, 0, len(scratchpad))
	for _, un := range scratchpad {
		units = append(units, un)
	}

	return units
}

// genMergeDevicesByNameUUID merges units by DeviceName+UUID+Realm
// and returns result as a slice of device-s.
func (out *output) genMergeDevicesByNameUUID(units []unit) []device {
	// Classify units by DeviceName+UUID+Realm
	scratchpad := make(map[UnitID][]unit)
	for _, un := range units {
		key := UnitID{
			DNSSDName: un.ID.DNSSDName,
			UUID:      un.ID.UUID,
			Realm:     un.ID.Realm,
		}
		scratchpad[key] = append(scratchpad[key], un)
	}

	// TODO -- merge by addr

	// Build slice of devices
	devices := make([]device, 0, len(scratchpad))
	for key, devunits := range scratchpad {
		dev := device{realm: key.Realm, uuid: key.UUID, units: devunits}
		for _, un := range devunits {
			dev.addrs = addrsMerge(dev.addrs, un.Addrs)
		}

		devices = append(devices, dev)
	}

	// Post-process devices
	for i := range devices {
		dev := &devices[i]
		dev.units = out.genMergeUnitCrossZones(dev.units)

	}

	// Sort devices by realm
	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].realm < devices[j].realm
	})

	return devices
}

// genMergeDevicesByUUID merges devices with the same UUID
func (out *output) genMergeDevicesByUUID(devices []device) []device {
	scratchpad := make(map[uuid.UUID]device)
	noUUID := make([]device, 0)

	for _, dev := range devices {
		if dev.uuid == uuid.NilUUID {
			noUUID = append(noUUID, dev)
			continue
		}

		if prev, found := scratchpad[dev.uuid]; found {
			prev.realm = RealmInvalid
			prev.units = append(prev.units, dev.units...)
			prev.addrs = addrsMerge(prev.addrs, dev.addrs)
			scratchpad[dev.uuid] = prev
		} else {
			scratchpad[dev.uuid] = dev
		}
	}

	devices = make([]device, 0, len(scratchpad)+len(noUUID))
	for _, dev := range scratchpad {
		devices = append(devices, dev)
	}
	devices = append(devices, noUUID...)

	return devices
}
