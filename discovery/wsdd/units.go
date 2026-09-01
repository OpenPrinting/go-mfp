// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Table of discovered units

package wsdd

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/internal/zone"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/wsd"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// units manages a table of discovered units.
//
// Please note, WSD uses 3-level addressing architecture:
//  1. Each device has a stable device "address", persistent
//     across reboots. Formally this is "xsd:anyURI", but usually,
//     uri:uuid:...
//  2. Each "hosted service" (i.e., printer, scanner, ...) has
//     one or more "transport addresses" (HTTP URLs). The transport
//     addresses are suitable to query device metadata. Printer
//     and scanner XAddrs are often the same, but this is not
//     necessarily.
//  3. Querying device metadata via XAddrs return one or more
//     service endpoints, suitable for printing or scanning,
//     depending on a service type.
//
// Every unit represents a discovered unit (the "hosted service"
// in the WSD terms), localized to the particular network interface.
// If the same device is visible over multiple network interfaces,
// the discovery system when merge them together.
type units struct {
	src   *source                    // Parent backend
	table map[discovery.UnitID]*unit // Discovered units
	lock  sync.Mutex                 // units.table lock
}

// newUnits creates a new table of units
func newUnits(src *source) *units {
	// Create units structure
	ut := &units{
		src:   src,
		table: make(map[discovery.UnitID]*unit),
	}

	return ut
}

// Close closes the unit table and cancels all ongoing discovery activity,
// like fetching unit's metadata
func (ut *units) Close() {
	for _, un := range ut.table {
		un.close()
	}
}

// InputFromUSB handles WSD message, received from UDP
func (ut *units) InputFromUDP(msg wsd.Msg) {
	ut.lock.Lock()
	defer ut.lock.Unlock()

	switch msg.Body.(type) {
	case wsd.AnnouncesBody:
		ut.handleAnnounces(msg)
	case wsd.Bye:
		ut.handleBye(msg)
	}
}

// handleBye handles received [wsd.Bye] message.
//
// Called under units.lock.
func (ut *units) handleBye(msg wsd.Msg) {
}

// handleAnnounce is the common handler for WSD announce messages
// (i.e., Hello, ProbeMatch and ResolveMatch).
//
// Called under units.lock.
func (ut *units) handleAnnounces(msg wsd.Msg) {
	body := msg.Body.(wsd.AnnouncesBody)
	action := body.Action()
	anns := body.Announces()
	ifidx := msg.IfIdx
	zone := zone.Name(ifidx)

	logmsg := log.Begin(ut.src.ctx)
	defer logmsg.Commit()

	// Parse and dispatch XAddrs. Log the event.
	for _, ann := range anns {
		target := ann.EndpointReference.Address
		ver := ann.MetadataVersion

		logmsg.Debug("%s received:", action)
		logmsg.Debug("  IP From:        %s", msg.From)
		logmsg.Debug("  IP To:          %s", msg.To)
		logmsg.Debug("  IP Zone:        %s", zone)
		logmsg.Debug("  Address         %s", target)
		logmsg.Debug("  Types           %s", ann.Types)
		logmsg.Debug("  MetadataVersion %d", ver)

		printUnitID := ut.makeUnitID(msg.IfIdx,
			discovery.ServicePrinter, target)
		scanUnitID := ut.makeUnitID(msg.IfIdx,
			discovery.ServiceScanner, target)

		if len(ann.XAddrs) != 0 {
			logmsg.Debug("  Xaddrs:")

			// Parse and collect XAddrs
			var xaddrs []*url.URL
			for _, s := range ann.XAddrs {
				if u := urlParse(s); u != nil {
					u = urlWithZone(u, zone)
					logmsg.Debug("    %s", u)
					xaddrs = append(xaddrs, u)
				} else {
					logmsg.Warning("    %s (invalid)", s)
				}
			}

			// Dispatch XAddrs
			if len(xaddrs) != 0 {
				if ann.Types.Contains(wsd.PrinterServiceType) {
					un := ut.getUnit(printUnitID, true)
					un.handleXaddrs(ifidx, target,
						xaddrs, ver)
				}

				if ann.Types.Contains(wsd.ScannerServiceType) {
					un := ut.getUnit(scanUnitID, true)
					un.handleXaddrs(ifidx, target,
						xaddrs, ver)
				}
			}
		}
	}
}

// makeUnitID creates a discovery.UnitID for the discovered
// service
func (ut *units) makeUnitID(ifidx int, svctype discovery.ServiceType,
	addr wsd.AnyURI) discovery.UnitID {
	return discovery.UnitID{
		UUID:     addr.UUID(),
		Backend:  Backend,
		Zone:     zone.Name(ifidx),
		SvcType:  svctype,
		SvcProto: discovery.ServiceWSD,
	}
}

// getUnit returns an unit by id. If unit is not known yet,
// it can be created on demand.
//
// Called under units.lock.
func (ut *units) getUnit(id discovery.UnitID,
	create bool) *unit {

	un := ut.table[id]
	if un == nil && create {
		un = newUnit(id, ut)
		ut.table[id] = un

		ut.src.queue.Push(&discovery.EventAddUnit{ID: id})
	}
	return un
}

// unit represents a discovered unit (a device)
type unit struct {
	parent        *units                     // Parent units table
	ctx           context.Context            // Cancelable context
	cancel        context.CancelFunc         // Its cancel function
	id            discovery.UnitID           // Unit ID
	xaddrsSeen    *generic.LockedSet[string] // Known XAddrs
	endpointsSeen *generic.LockedSet[string] // Known endpoints
	paramsSent    atomic.Bool                // EventXXXParameters reported
	closing       atomic.Bool                // unit.close in progress
	closewait     sync.WaitGroup             // for unit.close
}

// newUnit creates a new unit
func newUnit(id discovery.UnitID, parent *units) *unit {
	ctx, cancel := context.WithCancel(parent.src.ctx)

	un := &unit{
		parent:        parent,
		ctx:           ctx,
		cancel:        cancel,
		id:            id,
		xaddrsSeen:    generic.NewLockedSet[string](),
		endpointsSeen: generic.NewLockedSet[string](),
	}

	return un
}

// close closes the unit.
func (un *unit) close() {
	un.closing.Store(true)

	un.parent.lock.Lock()
	delete(un.parent.table, un.id)
	un.parent.lock.Unlock()

	un.cancel()
	un.closewait.Wait()
}

// handleXaddrs handles newly discovered XAddrs
//
// Called under units.lock.
func (un *unit) handleXaddrs(ifidx int, target wsd.AnyURI,
	xaddrs []*url.URL, ver uint64) {

	src := un.parent.src

	for _, xaddr := range xaddrs {
		if !un.xaddrsSeen.TestAndAdd(xaddr.String()) {
			continue
		}

		un.closewait.Add(1)
		go func(xaddr2 *url.URL) {
			meta := src.mex.Get(un.ctx, ifidx, target, xaddr2, ver)
			un.handleMetadata(meta)
			un.closewait.Done()
		}(xaddr)
	}
}

// handleMetadata handles the WSD metadata for the unit
func (un *unit) handleMetadata(metadata []mexData) {
	zone := un.id.Zone

	for _, meta := range metadata {
		mfg := meta.ThisModel.Manufacturer.NeutralLang().String
		mdl := meta.ThisModel.ModelName.NeutralLang().String
		adm := meta.ThisModel.PresentationURL

		endpoints := un.extractMetadataEndpoints(meta)

		logmsg := log.Begin(un.ctx)
		logmsg.Debug("Got %s metadata (from %s)",
			un.id.SvcType, meta.from)
		logmsg.Debug("  Manufacturer: %q", mfg)
		logmsg.Debug("  Model:        %q", mdl)
		if adm != nil {
			logmsg.Debug("  Admin URL:    %q", *adm)
		}

		if len(endpoints) > 0 {
			logmsg.Debug("  Endpoints:")

			for _, endpoint := range endpoints {
				u := urlParse(endpoint)
				if u == nil {
					logmsg.Warning("    %s (bad)", endpoint)
					continue
				}

				u = urlWithZone(u, zone)
				s := u.String()

				if !un.xaddrsSeen.TestAndAdd(s) {
					logmsg.Debug("    %s (dup)", s)
					continue
				}

				logmsg.Debug("    %s", s)
				un.sendEndpoint(u)
			}
		}

		logmsg.Commit()

		un.sendParameters(mfg, mdl, adm)
	}
}

// extractMetadataEndpoints extract and returns service endpoint URLs
func (un *unit) extractMetadataEndpoints(meta mexData) []string {
	t := un.wsdType()

	var endpoints []string

	for _, hosted := range meta.Relationship.Hosted {
		if !hosted.Types.Contains(t) {
			continue
		}

		for _, endpoint := range hosted.EndpointReference {
			endpoints = append(endpoints, string(endpoint.Address))
		}
	}

	return endpoints
}

// wsdType returns WSD service type for the unit
func (un *unit) wsdType() wsd.Type {
	switch un.id.SvcType {
	case discovery.ServicePrinter:
		return wsd.PrinterServiceType
	case discovery.ServiceScanner:
		return wsd.ScannerServiceType
	}
	return wsd.UnknownType
}

// sendParameters sends EventPrinterParameters or EventScannerParameters
// to the discovery system.
func (un *unit) sendParameters(mfg, mdl string, adm optional.Val[string]) {
	if un.paramsSent.Swap(true) {
		return
	}

	var evnt discovery.Event

	mkmodel := mdl
	if !strings.HasPrefix(mkmodel, mfg) {
		mkmodel = mfg + " " + mdl
	}

	switch un.id.SvcType {
	case discovery.ServicePrinter:
		var adminURL string
		if adm != nil {
			adminURL = *adm
		}

		evnt = &discovery.EventPrinterParameters{
			ID:              un.id,
			MakeModel:       mkmodel,
			AdminURL:        adminURL,
			PPDManufacturer: mfg,
			PPDModel:        mdl,
			Printer: discovery.PrinterParameters{
				PSProduct: "(" + mdl + ")",
			},
		}

	case discovery.ServiceScanner:
		evnt = &discovery.EventScannerParameters{
			ID:        un.id,
			MakeModel: mkmodel,
		}

	default:
		panic("internal error")
	}

	un.parent.src.queue.Push(evnt)
}

// sendEndpoint sends EventAddEndpoint to the discovery system.
func (un *unit) sendEndpoint(u *url.URL) {
	s := u.String()
	if !un.endpointsSeen.TestAndAdd(s) {
		return
	}

	evnt := &discovery.EventAddEndpoint{
		ID:       un.id,
		Endpoint: s,
	}

	un.parent.src.queue.Push(evnt)
}
