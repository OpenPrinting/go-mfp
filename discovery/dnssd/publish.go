// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// DNS-SD publisher

package dnssd

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/OpenPrinting/go-avahi"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/log"
)

// Publisher publishes (advertises) a [discovery.DNSSDDevice]
// to DND-SD.
type Publisher struct {
	ctx      context.Context    // For logging and Close
	cancel   context.CancelFunc // Used by Publisher.Close
	done     sync.WaitGroup     // Used to wait for goroutine termination
	services []*avahi.Service   // DNS-SD services being published
	instance string             // Instance name, for logging
}

// NewPublisher creates a new Publisher for the [discovery.DNSSDDevice].
func NewPublisher(ctx context.Context, dev *discovery.DNSSDDevice) *Publisher {
	// Setup DNS-SD services
	services := make([]*avahi.Service, 0, len(dev.Services))
	for _, dnssdsvc := range dev.Services {
		for _, svctype := range dnssdsvc.Types {
			svc := &avahi.Service{
				IfIdx:        avahi.IfIndexUnspec,
				Flags:        0,
				SvcType:      svctype,
				SvcSubTypes:  nil,
				InstanceName: dev.Instance,
				Domain:       "",
				Hostnames:    nil,
				Endpoints:    nil,
				Txt:          dnssdsvc.TXT,
			}
			services = append(services, svc)
		}
	}

	// Configure Context
	ctx = log.WithPrefix(ctx, "dnssd")
	ctx, cancel := context.WithCancel(ctx)

	// Create a Publisher
	pub := &Publisher{
		ctx:      ctx,
		cancel:   cancel,
		services: services,
		instance: dev.Instance,
	}

	// Issue log message
	rec := log.Begin(ctx)
	rec.Debug("%q: publishing requested", pub.instance)
	rec.Verbose("%q: services are:", pub.instance)
	for _, svc := range services {
		rec.Verbose("  --------------------------------")
		rec.Verbose("  svctype:   %q", svc.SvcType)
		rec.Verbose("  subtypes:  %q", svc.SvcSubTypes)
		rec.Verbose("  TXT:       %q", svc.Txt)

		var endpoints []string
		for _, endpoint := range svc.Endpoints {
			endpoints = append(endpoints, endpoint.String())
		}

		s := "none"
		if endpoints != nil {
			s = strings.Join(endpoints, ", ")
		}

		rec.Verbose("  endpoints: %s\n", s)

	}
	rec.Commit()

	// Start Publisher goroutine
	pub.done.Add(1)
	go pub.proc()

	return pub
}

// Close closes the Publisher and withdraws DNS-SD advertising
// being published.
func (pub *Publisher) Close() {
	pub.cancel()
	pub.done.Wait()
	log.Debug(pub.ctx, "%q: publishing closed", pub.instance)
}

// proc runs on its own goroutine and manages Publisher operations.
func (pub *Publisher) proc() {
	defer pub.done.Done()

	for {
		err := avahi.SimpleServicePublisher(pub.ctx,
			avahi.ProtocolUnspec, 0, nil, pub.services)

		log.Verbose(pub.ctx, "publishing interrupted: %s: ", err)

		select {
		case <-pub.ctx.Done():
			return
		case <-time.After(avahiClientRestartInterval):
		}
	}
}
