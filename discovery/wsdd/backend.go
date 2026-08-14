// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// WSDD backend

package wsdd

import (
	"context"
	"net/netip"

	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/wsd"
)

// backend is the [discovery.Backend] for WSD device discovery.
type backend struct {
	ctx   context.Context       // For logging and backend.Close
	queue *discovery.Eventqueue // Event queue
	links *links                // Per-local address links
	units *units                // Discovered units
	mex   *mexGetter            // Metadata getter
	res   *urlResolver          // URL resolver
}

// NewBackend creates a new [discovery.Backend] for WSD device discovery.
func NewBackend(ctx context.Context) (discovery.Backend, error) {
	// Set log prefix
	ctx = log.WithPrefix(ctx, "wsdd")

	// Create backend structure
	back := &backend{
		ctx: ctx,
	}

	// Create links
	var err error
	back.links, err = newLinks(back)
	if err != nil {
		return nil, err
	}

	// Create other stuff
	back.units = newUnits(back)
	back.mex = newMexGetter(back)
	back.res = newURLResolver(back)

	return back, nil
}

// Name returns backend name.
func (back *backend) Name() string {
	return "wsdd"
}

// DeviceID returns a subset of UnitID, that uniquely
// identifies a physical device.
func (back *backend) DeviceID(id discovery.UnitID) discovery.UnitID {
	return discovery.UnitID{
		UUID:    id.UUID,
		Zone:    id.Zone,
		Backend: id.Backend,
	}
}

// Start starts Backend operations.
func (back *backend) Start(queue *discovery.Eventqueue) {
	back.queue = queue
	back.links.Start()

	log.Debug(back.ctx, "backend started")
}

// Close closes the backend
func (back *backend) Close() {
	back.links.Close()
	back.units.Close()
	back.res.Close()
}

// input handles received UDP messages.
func (back *backend) input(data []byte, from, to netip.AddrPort, ifidx int) {
	// Silently drop looped packets
	if back.links.IsLocalPort(from) {
		return
	}

	// Decode the message
	back.debug("%d bytes received from %s%%%d", len(data), from, ifidx)

	msg, err := wsd.DecodeMsg(data)
	if err != nil {
		back.warning("%s", err)
		return
	}

	// Fill Msg.From, Msg.To and Msg.IfIdx
	msg.From = from
	msg.To = to
	msg.IfIdx = ifidx

	// Dispatch the message
	back.debug("%s message received", msg.Header.Action)

	switch msg.Header.Action {
	case wsd.ActHello, wsd.ActBye, wsd.ActProbeMatches,
		wsd.ActResolveMatches:
		back.units.InputFromUDP(msg)
	}
}

// Debug writes a LevelDebug message on behalf of the backend.
func (back *backend) debug(format string, args ...any) {
	log.Debug(back.ctx, format, args...)
}

// Warning writes a LevelWarning message on behalf of the backend.
func (back *backend) warning(format string, args ...any) {
	log.Warning(back.ctx, format, args...)
}

// Error writes a LevelError message on behalf of the backend.
func (back *backend) error(format string, args ...any) {
	log.Error(back.ctx, format, args...)
}
