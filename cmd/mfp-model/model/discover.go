// MFP - Miulti-Function Printers and scanners toolkit
// The "model" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device discovery

package model

import (
	"context"

	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/discovery/dnssd"
	"github.com/OpenPrinting/go-mfp/discovery/wsdd"
)

// discoverDevices searches device by the DNS-SD name
func discoverByName(ctx context.Context, nm string) (*discovery.Device, error) {
	// Prepare discovery.Client
	clnt := discovery.NewClient(ctx)
	defer clnt.Close()

	backend, err := dnssd.NewBackend(ctx, "", 0)
	if err != nil {
		return nil, err
	}

	defer backend.Close()
	clnt.AddBackend(backend)

	backend, err = wsdd.NewBackend(ctx)
	if err != nil {
		return nil, err
	}

	defer backend.Close()
	clnt.AddBackend(backend)

	return clnt.GetByDNSSD(ctx, nm, discovery.ModeNormal)
}
