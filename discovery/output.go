// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Output generation

package discovery

import (
	"context"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// output generates and manages the final discovery output from
// the internal representation of the discovered information,
// gathered in the cache
type output struct {
	ctx     context.Context // Logging context
	devices []Device        // Cached output data
	ttl     time.Time       // Cache valid until this time
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
	// Setup logging
	rec := log.Begin(out.ctx)
	defer rec.Commit()

	rec.Trace("preparing discovery output")

	rec.Trace("units discovered:")
	unitsToLogRecord(rec, units)

	// Extract IP addresses
	out.genExtractIPAddresses(units)

	// Merge variants
	units = out.genMergeUnitCrossVariandsAndZones(units)

	rec.Trace("units variants merged:")
	unitsToLogRecord(rec, units)

	// Classify units by DeviceID
	groups := out.groupUnitsByDeviceID(units)
	groups = out.mergeUnitGrousByUUID(groups)

	// Generate final output, save and return
	outdevs := make([]Device, len(groups))
	for i := range groups {
		outdevs[i] = groups[i].Export()
	}

	out.devices = outdevs
	out.ttl = ttl

	// Log result
	rec.Trace("devices reported:")
	for _, dev := range outdevs {
		rec.Object(log.LevelTrace, 2, dev)
	}

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

// genMergeUnitCrossVariandsAndZones merges together units with
// distinct UnitID.Variant or UnitID.Zone, but otherwise equal.
//
// The same device may be visible in different variants (say, IP4
// vs IP6, or HTTP vs HTTPS). This is why variants are merged.
//
// Also, the same device may be visible from different network
// interfaces (say, Ethernet vs WiFi), which will make them
// identical but with different zones. This is why cross-zone
// merge is also performed.
func (out *output) genMergeUnitCrossVariandsAndZones(units []unit) []unit {
	scratchpad := make(map[UnitID]unit)
	for _, un := range units {
		un.ID.Variant = ""
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

// groupUnitsByDeviceID groups units by their DeviceID,
// as reported by Backend.DeviceID()
func (out *output) groupUnitsByDeviceID(units []unit) []unitgroup {
	// Classify units by DeviceID
	scratchpad := make(map[UnitID][]unit)
	for _, un := range units {
		key := un.ID.Backend.DeviceID(un.ID)
		scratchpad[key] = append(scratchpad[key], un)
	}

	// Build slice of unit groups
	groups := make([]unitgroup, 0, len(scratchpad))
	for _, grp := range scratchpad {
		groups = append(groups, grp)
	}

	return groups
}

// mergeUnitGrousByUUID merges unit groups by UUID.
//
// It allows to merge units of the same device, discovered
// by different backends (say, dnssd vs wsdd).
func (out *output) mergeUnitGrousByUUID(groups []unitgroup) []unitgroup {
	// Classify groups by UUID
	//
	// We keep with indices of groups within the original groups slice,
	// because indices are comparable.
	byUUID := make(map[uuid.UUID][]int)
	for idx, grp := range groups {
		for _, uu := range grp.UUIDs() {
			byUUID[uu] = append(byUUID[uu], idx)
		}
	}

	// Now match groups for each UUID
	consumed := generic.NewSet[int]()
	for _, indices := range byUUID {
		for i, idx1 := range indices {
			grp1 := groups[idx1]
			for _, idx2 := range indices[i+1:] {
				if consumed.Contains(idx2) {
					continue
				}

				grp2 := groups[idx2]
				if grp1.CanMergeByUUID(grp2) {
					grp1 = append(grp1, grp2...)
					groups[idx1] = grp1
					consumed.Add(idx2)
				}
			}
		}
	}

	// And rebuild list of groups
	result := []unitgroup{}
	for _, indices := range byUUID {
		for _, idx := range indices {
			if consumed.TestAndAdd(idx) {
				grp := groups[idx]
				result = append(result, grp)
			}
		}
	}

	return result
}
