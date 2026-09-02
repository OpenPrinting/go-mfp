// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// discovery.Source implementation for WS-Discovery

package wsdd

import (
	"context"
	"net/netip"

	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/wsd"
)

// source is the [discovery.Backend] for WSD device discovery.
type source struct {
	ctx   context.Context       // For logging and source.Close
	queue *discovery.Eventqueue // Event queue
	links *links                // Per-local address links
	units *units                // Discovered units
	mex   *mexGetter            // Metadata getter
	res   *urlResolver          // URL resolver
}

// NewBackend creates a new [discovery.Source] for WS-Discovery.
func newSource(ctx context.Context,
	queue *discovery.Eventqueue) (*source, error) {

	// Set log prefix
	ctx = log.WithPrefix(ctx, "wsdd")

	// Create source structure
	src := &source{
		ctx:   ctx,
		queue: queue,
	}

	// Create links
	var err error
	src.links, err = newLinks(src)
	if err != nil {
		return nil, err
	}

	// Create other stuff
	src.units = newUnits(src)
	src.mex = newMexGetter(src)
	src.res = newURLResolver(src)

	// Start source operations
	src.links.Start()
	src.debug("started")

	return src, nil
}

// Close closes the source
func (src *source) Close() {
	src.links.Close()
	src.units.Close()
	src.res.Close()
}

// input handles received UDP messages.
func (src *source) input(data []byte, from, to netip.AddrPort, ifidx int) {
	// Silently drop looped packets
	if src.links.IsLocalPort(from) {
		return
	}

	// Decode the message
	src.trace("%d bytes received from %s%%%d", len(data), from, ifidx)

	msg, err := wsd.DecodeMsg(data)
	if err != nil {
		src.warning("%s", err)
		return
	}

	// Fill Msg.From, Msg.To and Msg.IfIdx
	msg.From = from
	msg.To = to
	msg.IfIdx = ifidx

	// Dispatch the message
	src.trace("%s message received", msg.Header.Action)

	switch msg.Header.Action {
	case wsd.ActHello, wsd.ActBye, wsd.ActProbeMatches,
		wsd.ActResolveMatches:
		src.units.InputFromUDP(msg)
	}
}

// trace writes a LevelTrace message on behalf of the source.
func (src *source) trace(format string, args ...any) {
	log.Trace(src.ctx, format, args...)
}

// verbose writes a LevelVerbose message on behalf of the source.
func (src *source) verbose(format string, args ...any) {
	log.Verbose(src.ctx, format, args...)
}

// debug writes a LevelDebug message on behalf of the source.
func (src *source) debug(format string, args ...any) {
	log.Debug(src.ctx, format, args...)
}

// Warning writes a LevelWarning message on behalf of the source.
func (src *source) warning(format string, args ...any) {
	log.Warning(src.ctx, format, args...)
}

// Error writes a LevelError message on behalf of the source.
func (src *source) error(format string, args ...any) {
	log.Error(src.ctx, format, args...)
}
