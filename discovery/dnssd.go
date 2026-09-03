// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// DNS-SD definitions

package discovery

import (
	"slices"

	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// DNSSDDevice represents a DNS-SD advertising of the device.
// This information can be obtained from the discovered device,
// using the [Device.DNSSD] method.
type DNSSDDevice struct {
	Instance string         // Instance name
	UUID     uuid.UUID      // Device UUID
	Services []DNSSDService // Services behind the device
}

// DNSSDService represents a single DNS-SD service (e.g.,
// IPP print or eSCL scan part of device).
type DNSSDService struct {
	Types    []string // Service types (e.g. ["_ipp._tcp", "_ipps._tcp"])
	SubTypes []string // Service subtypes ("_universal._sub._ipp._tcp")
	TXT      []string // TXT record ("key=value"...)
}

// merge merges two services, assuming they are duplicated
// services.
//
// Technically, this duplication comes from the fact, that
// some DNS-SD services may represent more that a single
// unit, in the discovery.Device terms. For example, IPP printer
// and faxout service will be represented by the single service
// at the DNS-SD side, while discovery will return two separate units.
func (svc *DNSSDService) merge(svc2 *DNSSDService) {
	// Merge Types
	for _, t := range svc2.Types {
		if !slices.Contains(svc.Types, t) {
			svc.Types = append(svc.Types, t)
		}
	}

	slices.Sort(svc.Types)
}
