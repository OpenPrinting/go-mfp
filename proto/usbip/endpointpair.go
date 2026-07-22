// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Pair of related In/Out endpoints

package usbip

import "context"

// EndpointPair represents a pair of related In/Out [Endpoint]s.
type EndpointPair struct {
	In  *Endpoint // Device->Host
	Out *Endpoint // Host->Device
}

// Read returns data that was sent to [EndpointPair.Out] from the USB side.
func (epp EndpointPair) Read(buf []byte) (int, error) {
	return epp.Out.Read(buf)
}

// ReadContext is [EndpointPair.Read] with [context.Context].
func (epp EndpointPair) ReadContext(ctx context.Context, buf []byte) (int, error) {
	return epp.Out.ReadContext(ctx, buf)
}

// Write writes data that will arrive at the [EndpointPair.In]
// at the USB side.
func (epp EndpointPair) Write(buf []byte) (int, error) {
	return epp.In.Write(buf)
}

// WriteContext is [EndpointPair.Write] with [context.Context].
func (epp EndpointPair) WriteContext(ctx context.Context, buf []byte) (int, error) {
	return epp.In.Write(buf)
}
