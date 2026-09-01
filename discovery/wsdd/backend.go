// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// discovery.Backend implementation for WS-Discovery

package wsdd

import (
	"context"

	"github.com/OpenPrinting/go-mfp/discovery"
)

// Backend is the DNS-SD backend for [discovery.NewClient].
var Backend = backend{}

// backend implements the discovery.Backend implementation for DNS-SD
type backend struct{}

// Name returns backend name.
func (backend) Name() string {
	return "wsdd"
}

// DeviceID returns a subset of UnitID, that uniquely
// identifies a physical device.
func (backend) DeviceID(id discovery.UnitID) discovery.UnitID {
	return discovery.UnitID{
		UUID:    id.UUID,
		Zone:    id.Zone,
		Backend: id.Backend,
	}
}

// Open creates a new [discovery.Source] for WS-Discovery.
func (backend) Open(ctx context.Context,
	queue *discovery.Eventqueue) (discovery.Source, error) {

	return newSource(ctx, queue)
}
