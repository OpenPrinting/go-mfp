// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Discovery cache

package discovery

import (
	"errors"
	"fmt"
	"time"
)

// cache represents the discovery cache
//
// Note, cache API is not reentrant and requires external locking.
type cache struct {
	readyAt time.Time            // Time when cache is warmed up and ready
	entries map[UnitID]*cacheEnt // Cache entries
	out     output               // Cached output
}

// cacheEnt is the cache entry for print/scan/faxout units.
type cacheEnt struct {
	unit
	hasParams        bool      // Parameters are received
	stagingEndpoints []string  // Newly discovered endpoints, on quarantine
	stagingDoneAt    time.Time // End of staging time. Zero if no staging.
}

// newCache creates the new discovery cache
func newCache() *cache {
	return &cache{
		readyAt: time.Now().Add(warmUpTime),
		entries: make(map[UnitID]*cacheEnt),
	}
}

// ReadyAt returns time when cache is ready to be exported, according to
// the cache state and export Mode
func (c *cache) ReadyAt(m Mode) time.Time {
	// Handle simple cases
	switch m {
	case ModeNormal:
		return c.readyAt
	case ModeSnapshot:
		return time.Time{}
	}

	// If cache is ready, no need to wait
	if c.out.Cached() != nil {
		return time.Time{}
	}

	// for ModeWaitIncomplete we need to do some more work
	ready := c.readyAt

	for _, ent := range c.entries {
		ent.stagingCheck()
		if ent.stagingInProgress() {
			ready = timeLatest(ready, ent.stagingDoneAt)
		}
	}

	return ready
}

// Export exports the cached data.
func (c *cache) Export() []Device {
	// If cached output available, just return it now
	if devices := c.out.Cached(); devices != nil {
		return devices
	}

	// Re-generate the output
	units := make([]unit, 0, len(c.entries))
	ttl := time.Now().Add(365 * 24 * time.Hour) // Far in a future

	for _, ent := range c.entries {
		if ent.ready() {
			units = append(units, ent.unit)
			if ent.stagingInProgress() {
				ttl = timeEarliest(ttl, ent.stagingDoneAt)
			}
		}
	}

	return c.out.Generate(ttl, units)
}

// Snapshot exports the cached data in the ModeSnapshot mode.
func (c *cache) Snapshot() []Device {
	units := make([]unit, 0, len(c.entries))
	ttl := time.Now().Add(365 * 24 * time.Hour) // Far in a future

	for _, ent := range c.entries {
		unit, ok := ent.snapshot()
		if ok {
			units = append(units, unit)
		}
	}

	var out output
	return out.Generate(ttl, units)
}

// AddUnit adds new unit. Called when EventAddUnit is received.
func (c *cache) AddUnit(evnt *EventAddUnit) error {
	if c.entries[evnt.ID] != nil {
		return errors.New("unit already added")
	}

	c.entries[evnt.ID] = &cacheEnt{unit: unit{ID: evnt.ID}}
	c.out.Invalidate()

	return nil
}

// DelUnit deletes the unit. Called when EventDelUnit is received.
func (c *cache) DelUnit(evnt *EventDelUnit) error {
	if c.entries[evnt.ID] == nil {
		return errors.New("unknown UnitID")
	}

	delete(c.entries, evnt.ID)
	c.out.Invalidate()

	return nil
}

// SetPrinterParameters saves printer parameters.
// Called when EventPrinterParameters is received
func (c *cache) SetPrinterParameters(evnt *EventPrinterParameters) error {
	ent, err := c.setParametersBegin(evnt.ID, ServicePrinter)
	if err != nil {
		return err
	}

	params := evnt.Printer
	params.fixup()

	ent.MakeModel = evnt.MakeModel
	ent.Location = evnt.Location
	ent.AdminURL = evnt.AdminURL
	ent.IconURL = evnt.IconURL
	ent.PPDManufacturer = evnt.PPDManufacturer
	ent.PPDModel = evnt.PPDModel
	ent.Params = params

	c.setParametersCommit(ent)

	return nil
}

// SetScannerParameters saves scanner parameters.
// Called when EventScannerParameters is received
func (c *cache) SetScannerParameters(evnt *EventScannerParameters) error {
	ent, err := c.setParametersBegin(evnt.ID, ServiceScanner)
	if err != nil {
		return err
	}

	params := evnt.Scanner

	ent.MakeModel = evnt.MakeModel
	ent.Location = evnt.Location
	ent.AdminURL = evnt.AdminURL
	ent.IconURL = evnt.IconURL
	ent.Params = params

	c.setParametersCommit(ent)

	return nil
}

// SetFaxoutParameters saves faxout parameters.
// Called when EventFaxoutParameters is received.
//
// Note, Faxout parameters are represented by the PrinterParameters type.
func (c *cache) SetFaxoutParameters(evnt *EventFaxoutParameters) error {
	ent, err := c.setParametersBegin(evnt.ID, ServiceFaxout)
	if err != nil {
		return err
	}

	params := evnt.Faxout
	params.fixup()

	ent.MakeModel = evnt.MakeModel
	ent.Location = evnt.Location
	ent.AdminURL = evnt.AdminURL
	ent.IconURL = evnt.IconURL
	ent.PPDManufacturer = evnt.PPDManufacturer
	ent.PPDModel = evnt.PPDModel
	ent.Params = params

	c.setParametersCommit(ent)

	return nil
}

// setParametersBegin begins operation of setting unit parameters.
// It returns cache entry to be modified or an error.
func (c *cache) setParametersBegin(id UnitID,
	svcMustBe ServiceType) (*cacheEnt, error) {

	ent := c.entries[id]
	if ent == nil {
		return nil, errors.New("unknown UnitID")
	}

	if ent.ID.SvcType != svcMustBe {
		return nil, fmt.Errorf("unit is %s, must be %s",
			ent.ID.SvcType, svcMustBe)
	}

	return ent, nil
}

// setParametersCommit finishes operation of setting unit parameters
func (c *cache) setParametersCommit(ent *cacheEnt) {
	ent.hasParams = true
	c.out.Invalidate()
}

// AddEndpoint adds unit endpoint.
func (c *cache) AddEndpoint(evnt *EventAddEndpoint) error {
	ent := c.entries[evnt.ID]
	if ent == nil {
		return errors.New("unknown UnitID")
	}

	endpoint := evnt.Endpoint
	if endpointsContain(ent.Endpoints, endpoint) ||
		endpointsContain(ent.stagingEndpoints, endpoint) {
		return errors.New("endpoint already added")
	}

	ent.stagingBegin()
	ent.stagingEndpoints, _ = endpointsAdd(ent.stagingEndpoints, endpoint)

	c.out.Invalidate()

	return nil
}

// AddEndpoint deletes unit endpoint.
func (c *cache) DelEndpoint(evnt *EventDelEndpoint) error {
	ent := c.entries[evnt.ID]
	if ent == nil {
		return errors.New("unknown UnitID")
	}

	endpoint := evnt.Endpoint

	switch {
	case endpointsContain(ent.Endpoints, endpoint):
		ent.Endpoints, _ = endpointsDel(ent.Endpoints, endpoint)

	case endpointsContain(ent.stagingEndpoints, endpoint):
		ent.stagingEndpoints, _ = endpointsDel(ent.stagingEndpoints,
			endpoint)
	default:
		return errors.New("unknown endpoint")
	}

	c.out.Invalidate()

	return nil
}

// ready checks if cache entry is ready to be exported.
func (ent *cacheEnt) ready() bool {
	if ent.hasParams {
		ent.stagingCheck() // Merge staged endpoints, if any
		return len(ent.Endpoints) > 0
	}
	return false
}

// snapshot makes a snapshot from the cache entry.
func (ent *cacheEnt) snapshot() (un unit, ok bool) {
	if ent.hasParams {
		un = ent.unit
		un.Endpoints = endpointsMerge(un.Endpoints,
			ent.stagingEndpoints)
		if len(un.Endpoints) > 0 {
			return un, true
		}
	}

	return unit{}, false
}

// stagingBegin starts staging interval for discovered endpoints.
//
// The purpose of this interval is to merge together endpoints
// discovery events, which may come not in the same time (for example,
// if device has IP4 and IP6 addresses, these addresses may come
// spread out over time).
//
// So when new endpoint is discovered, we don't publish it immediately,
// but instead save into the staging area and start the timer. The subsequent
// endpoints are added to the staging area, but already running timer is
// not restarted. When the timer expires, all endpoints collected at
// the staging area is published.
func (ent *cacheEnt) stagingBegin() {
	if !ent.stagingInProgress() {
		ent.stagingDoneAt = time.Now().Add(stabilizationTime)
	}
}

// stagingEnd publishes all endpoints collected in the staging area.
func (ent *cacheEnt) stagingEnd() {
	ent.Endpoints = endpointsMerge(ent.Endpoints, ent.stagingEndpoints)
	ent.stagingEndpoints = ent.stagingEndpoints[:0]
	ent.stagingDoneAt = time.Time{}
}

// stagingInProgress tells if staging interval is in progress.
func (ent *cacheEnt) stagingInProgress() bool {
	return !ent.stagingDoneAt.IsZero()
}

// stagingCheck checks if staging interval is expired. If so,
// it merges publishes the staged endpoints and clears staging area.
//
// It returns 'true' if the previously started staging interval still
// in progress.
func (ent *cacheEnt) stagingCheck() {
	if ent.stagingInProgress() && !ent.stagingDoneAt.After(time.Now()) {
		ent.stagingEnd()
	}
}
