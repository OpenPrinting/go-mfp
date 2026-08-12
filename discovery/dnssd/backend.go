// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Discovery Backend

package dnssd

import (
	"context"
	"fmt"
	"net/netip"
	"sync"

	"github.com/OpenPrinting/go-avahi"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/internal/zone"
	"github.com/OpenPrinting/go-mfp/log"
)

// backend is the [discovery.Backend] for DNS-SD discovery
type backend struct {
	ctx    context.Context       // For logging and backend.Close
	cancel context.CancelFunc    // Context's cancel function
	clnt   *avahiClient          // Avahi connection
	queue  *discovery.Eventqueue // Output queue
	done   sync.WaitGroup        // For backend.Close synchronization
}

// NewBackend creates a new [discovery.Backend] for DNS-SD discovery.
func NewBackend(ctx context.Context,
	domain string, flags LookupFlags) (discovery.Backend, error) {

	// Set log prefix
	ctx = log.WithPrefix(ctx, "dnssd")

	// Initialize loopback index on demand
	err := loopbackInit()
	if err != nil {
		log.Error(ctx, "%s", err)
		return nil, err
	}

	// Create Avahi client.
	clnt, err := newAvahiClient(domain, flags)
	if err != nil {
		log.Error(ctx, "%s", err)
		return nil, err
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(ctx)

	back := &backend{
		ctx:    ctx,
		cancel: cancel,
		clnt:   clnt,
	}

	return back, nil
}

// Name returns backend name.
func (back *backend) Name() string {
	return "dnssd"
}

// Start starts Backend operations.
func (back *backend) Start(queue *discovery.Eventqueue) {
	back.queue = queue

	back.done.Add(1)
	go back.proc()

	log.Debug(back.ctx, "backend started")
}

// Close closes the backend
func (back *backend) Close() {
	back.cancel()
	back.done.Wait()

	back.clnt.Close()
}

// proc runs the backend event loop on its separate goroutine.
func (back *backend) proc() {
	defer back.done.Done()

	var err error
	for err == nil {
		// Start service browsers.
		err = back.startServiceBrowsers()

		// Handle events until an error
		for err == nil {
			var evnt any
			evnt, err = back.clnt.poll(back.ctx)

			switch evnt := evnt.(type) {
			case *avahi.ClientEvent:
				err = back.onClientEvent(evnt)

			case *avahi.ServiceBrowserEvent:
				err = back.onServiceBrowserEvent(evnt)

			case *avahi.ServiceResolverEvent:
				err = back.onServiceResolverEvent(evnt)

			case *avahi.RecordBrowserEvent:
				switch evnt.RType {
				case avahi.DNSTypeTXT:
					err = back.onTxtBrowserEvent(evnt)
				case avahi.DNSTypeA, avahi.DNSTypeAAAA:
					err = back.onAddrBrowserEvent(evnt)
				}
			}
		}

		// Attempt error recovery.
		err = back.clnt.Restart(back.ctx)
		if err == nil {
			log.Debug(back.ctx, "avahi client: restarted")
		}
	}
}

// startServiceBrowsers starts service browsers for all service
// types mentioned in the svcTypes.
//
// The newly created browsers are added to the back.poller
func (back *backend) startServiceBrowsers() error {
	for _, svctype := range svcTypes {
		_, err := back.clnt.NewServiceBrowser(svctype)

		title := fmt.Sprintf("svc-browse: start %q", svctype)

		if err != nil {
			log.Error(back.ctx, "%s: %s", title, err)
			return err
		}

		log.Debug(back.ctx, "%s: OK", title)
	}

	return nil
}

// onClientEvent handles avahi.ClientEvent.
func (back *backend) onClientEvent(evnt *avahi.ClientEvent) error {
	log.Debug(back.ctx, "avahi client: %s", evnt.State)
	switch evnt.State {
	case avahi.ClientStateFailure:
		return evnt.Err
	}

	return nil
}

// onServiceBrowserEvent handles avahi.ServiceBrowserEvent
func (back *backend) onServiceBrowserEvent(
	evnt *avahi.ServiceBrowserEvent) error {

	switch evnt.Event {
	case avahi.BrowserNew:
		key := avahiServiceKeyFromServiceBrowserEvent(evnt)
		title := fmt.Sprintf("svc-browse: found %s", key)

		if !back.clnt.HasService(key) {
			log.Debug(back.ctx, "%s", title)
		} else {
			log.Debug(back.ctx, "%s (duplicate)", title)
			return nil
		}

		return back.addService(key)

	case avahi.BrowserRemove:
		key := avahiServiceKeyFromServiceBrowserEvent(evnt)
		title := fmt.Sprintf("svc-browse: removed %s", key)

		service := back.clnt.GetService(key)
		if service != nil {
			log.Debug(back.ctx, "%s", title)
			service.Delete()
		} else {
			log.Debug(back.ctx, "%s (not found)", title)
		}

	case avahi.BrowserFailure:
		key := avahiServiceKeyFromServiceBrowserEvent(evnt)
		title := fmt.Sprintf("svc-browse: failed  %s", key)

		log.Warning(back.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// onServiceResolverEvent handles avahi.ServiceResolverEvent
func (back *backend) onServiceResolverEvent(
	evnt *avahi.ServiceResolverEvent) error {
	switch evnt.Event {
	case avahi.ResolverFound:
		key := avahiServiceKeyFromResolverEvent(evnt)
		title := fmt.Sprintf("svc-resolve: found %s", key)

		service := back.clnt.GetService(key)
		if service == nil {
			// It may be out of order avahi.ResolverFound
			// event, received while service already removed,
			// so just log and return.
			log.Debug(back.ctx, "%s (unknown service)", title)
			return nil
		}

		log.Begin(back.ctx).
			Debug("%s:", title).
			Debug("  host: %s", evnt.Hostname).
			Debug("  port: %d", evnt.Port).
			Commit()

		service.SetPort(evnt.Port)
		back.setServiceHostname(service, evnt.Hostname)

	case avahi.ResolverFailure:
		key := avahiServiceKeyFromResolverEvent(evnt)
		title := fmt.Sprintf("svc-resolve: failed  %s", key)

		// Note, typically it's not fatal, just answer
		// doesn't want to come in time.
		log.Warning(back.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// onTxtBrowserEvent handles avahi.RecordBrowserEvent
// for per-service TXT record browser
func (back *backend) onTxtBrowserEvent(evnt *avahi.RecordBrowserEvent) error {
	switch evnt.Event {
	case avahi.BrowserNew:
		key := avahiServiceKeyFromRecordBrowserEvent(evnt)
		title := fmt.Sprintf("txt-browse: found %s", key)
		log.Debug(back.ctx, "%s", title)

		service := back.clnt.GetService(key)
		if service == nil {
			log.Debug(back.ctx, "%s: service not found", title)
			return nil
		}

		svcType := key.SvcType
		svcInstance := key.InstanceName
		txt := avahi.DNSDecodeTXT(evnt.RData)

		if key.IsPrinter() {
			txtPrinter, err := decodeTxtPrinter(svcType,
				svcInstance, txt)
			if err != nil {
				log.Debug(back.ctx, "%s: %s", title, err)
				return nil // Don't propagate the error
			}

			id := key.PrinterUnitID(txtPrinter)
			un := service.GetUnit(id.Queue)
			if un == nil {
				un = newPrinterUnit(back.queue, id, txtPrinter)
				service.AddUnit(id.Queue, un)
			} else {
				un.SetTxtPrinter(txtPrinter)
			}
		} else {
			txtScanner, err := decodeTxtScanner(svcType,
				svcInstance, txt)
			if err != nil {
				log.Debug(back.ctx, "%s: %s", title, err)
				return nil // Don't propagate the error
			}

			unName := "scan"
			un := service.GetUnit(unName)
			if un == nil {
				id := key.ScannerUnitID(txtScanner)
				un = newScannerUnit(back.queue, id, txtScanner)
				service.AddUnit(unName, un)
			} else {
				un.SetTxtScanner(txtScanner)
			}
		}

	case avahi.BrowserFailure:
		key := avahiServiceKeyFromRecordBrowserEvent(evnt)
		title := fmt.Sprintf("txt-browse: failed %s", key)

		log.Warning(back.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// onAddrBrowserEvent handles avahi.RecordBrowserEvent
// for per-hostname A and AAAA record browsers
func (back *backend) onAddrBrowserEvent(
	evnt *avahi.RecordBrowserEvent) error {

	switch evnt.Event {
	case avahi.BrowserNew, avahi.BrowserRemove:
		key := avahiHostnameKeyFromRecordBrowserEvent(evnt)
		var title string

		if evnt.Event == avahi.BrowserNew {
			title = fmt.Sprintf("addr-browse: found %s", key)
		} else {
			title = fmt.Sprintf("addr-browse: removed %s", key)
		}

		// Find hostname resolver
		hostname := back.clnt.GetHostname(key)
		if hostname == nil {
			log.Debug(back.ctx, "%s: unknown hostname", title)
			return nil
		}

		// Decode address
		var addr netip.Addr
		if key.Proto == avahi.ProtocolIP4 {
			addr = avahi.DNSDecodeA(evnt.RData)
		} else {
			addr = avahi.DNSDecodeAAAA(evnt.RData)
			addr = addr.WithZone(zone.Name(int(evnt.IfIdx)))
		}

		if addr == (netip.Addr{}) {
			log.Debug(back.ctx, "%s: invalid addr", title)
			return nil
		}

		// Add or delete the address
		hasAddr := hostname.HasAddr(addr)
		switch {
		case evnt.Event == avahi.BrowserNew && hasAddr:
			log.Debug(back.ctx, "%s: %s (duplicate)", title, addr)
			return nil
		case evnt.Event == avahi.BrowserRemove && !hasAddr:
			log.Debug(back.ctx, "%s: %s (unknown addr)", title,
				addr)
			return nil

		case evnt.Event == avahi.BrowserNew:
			hostname.AddAddr(addr, evnt.IfIdx)
		case evnt.Event == avahi.BrowserRemove:
			hostname.DelAddr(addr)
		}

		log.Debug(back.ctx, "%s: %s", title, addr)

	case avahi.BrowserFailure:
		title := fmt.Sprintf("addr-browse: failed %s", evnt.Name)

		log.Warning(back.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// addService creates a new avahiService and registers it
// in the back.clnt
func (back *backend) addService(key avahiServiceKey) error {
	// Create service resolver
	svcResolver, err := back.clnt.NewServiceResolver(key)

	title := fmt.Sprintf("svc-resolve: start %s", key)

	if err != nil {
		log.Error(back.ctx, "%s: %s", title, err)
		return err
	}

	log.Debug(back.ctx, "%s: OK", title)

	// Create TXT record browser
	txtBrowser, err := back.clnt.NewTxtBrowser(key)

	title = fmt.Sprintf("txt-browse: start %s", key)

	if err != nil {
		log.Error(back.ctx, "%s: %s", title, err)
		svcResolver.Close()
		return err
	}

	log.Debug(back.ctx, "%s: OK", title)

	// Add the service
	back.clnt.AddService(key, svcResolver, txtBrowser)
	return nil
}

// addHostname adds a new avahiHostname for the key
func (back *backend) addHostname(key avahiHostnameKey) (*avahiHostname, error) {
	// Create A/AAAA record browser
	addrBrowser, err := back.clnt.NewAddrBrowser(key)

	title := fmt.Sprintf("addr-browse: start %s", key)

	if err != nil {
		log.Error(back.ctx, "%s: %s", title, err)
		return nil, err
	}

	log.Debug(back.ctx, "%s: OK", title)

	// Add avahiHostname
	return back.clnt.AddHostname(key, addrBrowser), nil
}

// setServiceHostname sets or updates the service's hostname.
//
// On success, it initiates hostname resolving, if it is not
// active yet.
func (back *backend) setServiceHostname(service *avahiService,
	name string) error {

	// Do nothing if service already has a hostname and
	// name was not changed.
	if service.hostname != nil && service.hostname.key.Hostname == name {
		return nil
	}

	// Prepare avahiHostnameKey
	key := service.key.HostnameKey(name)

	// Find or create avahiHostname
	hostname := back.clnt.GetHostname(key)
	if hostname == nil {
		var err error
		hostname, err = back.addHostname(key)
		if err != nil {
			return err
		}
	}

	// Associate hostname with the service
	service.SetHostname(hostname)

	return nil
}
