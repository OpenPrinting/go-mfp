// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// discovery.Source implementation for DNS-SD

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

// source is the [discovery.Source] for DNS-SD discovery
type source struct {
	ctx    context.Context       // For logging and source.Close
	cancel context.CancelFunc    // Context's cancel function
	clnt   *avahiClient          // Avahi connection
	queue  *discovery.Eventqueue // Output queue
	done   sync.WaitGroup        // For source.Close synchronization
}

// newSource creates a new source for DNS-SD discovery.
func newSource(ctx context.Context,
	queue *discovery.Eventqueue) (*source, error) {

	// Set log prefix
	ctx = log.WithPrefix(ctx, "dnssd")

	// Initialize loopback index on demand
	err := loopbackInit()
	if err != nil {
		log.Error(ctx, "%s", err)
		return nil, err
	}

	// Create Avahi client.
	clnt, err := newAvahiClient("", 0)
	if err != nil {
		log.Error(ctx, "%s", err)
		return nil, err
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(ctx)

	// Create source structure
	src := &source{
		ctx:    ctx,
		cancel: cancel,
		clnt:   clnt,
		queue:  queue,
	}

	// Start source operations
	src.done.Add(1)
	go src.proc()

	log.Debug(src.ctx, "started")

	return src, nil
}

// Close closes the source
func (src *source) Close() {
	src.cancel()
	src.done.Wait()

	src.clnt.Close()
}

// proc runs the source event loop on its separate goroutine.
func (src *source) proc() {
	defer src.done.Done()

	var err error
	for err == nil {
		// Start service browsers.
		err = src.startServiceBrowsers()

		// Handle events until an error
		for err == nil {
			var evnt any
			evnt, err = src.clnt.poll(src.ctx)

			switch evnt := evnt.(type) {
			case *avahi.ClientEvent:
				err = src.onClientEvent(evnt)

			case *avahi.ServiceBrowserEvent:
				err = src.onServiceBrowserEvent(evnt)

			case *avahi.ServiceResolverEvent:
				err = src.onServiceResolverEvent(evnt)

			case *avahi.RecordBrowserEvent:
				switch evnt.RType {
				case avahi.DNSTypeTXT:
					err = src.onTxtBrowserEvent(evnt)
				case avahi.DNSTypeA, avahi.DNSTypeAAAA:
					err = src.onAddrBrowserEvent(evnt)
				}
			}
		}

		// Attempt error recovery.
		err = src.clnt.Restart(src.ctx)
		if err == nil {
			log.Debug(src.ctx, "avahi client: restarted")
		}
	}
}

// startServiceBrowsers starts service browsers for all service
// types mentioned in the svcTypes.
//
// The newly created browsers are added to the src.poller
func (src *source) startServiceBrowsers() error {
	for _, svctype := range svcTypes {
		_, err := src.clnt.NewServiceBrowser(svctype)

		title := fmt.Sprintf("svc-browse: start %q", svctype)

		if err != nil {
			log.Error(src.ctx, "%s: %s", title, err)
			return err
		}

		log.Debug(src.ctx, "%s: OK", title)
	}

	return nil
}

// onClientEvent handles avahi.ClientEvent.
func (src *source) onClientEvent(evnt *avahi.ClientEvent) error {
	log.Debug(src.ctx, "avahi client: %s", evnt.State)
	switch evnt.State {
	case avahi.ClientStateFailure:
		return evnt.Err
	}

	return nil
}

// onServiceBrowserEvent handles avahi.ServiceBrowserEvent
func (src *source) onServiceBrowserEvent(
	evnt *avahi.ServiceBrowserEvent) error {

	switch evnt.Event {
	case avahi.BrowserNew:
		key := avahiServiceKeyFromServiceBrowserEvent(evnt)
		title := fmt.Sprintf("svc-browse: found %s", key)

		if !src.clnt.HasService(key) {
			log.Debug(src.ctx, "%s", title)
		} else {
			log.Debug(src.ctx, "%s (duplicate)", title)
			return nil
		}

		return src.addService(key)

	case avahi.BrowserRemove:
		key := avahiServiceKeyFromServiceBrowserEvent(evnt)
		title := fmt.Sprintf("svc-browse: removed %s", key)

		service := src.clnt.GetService(key)
		if service != nil {
			log.Debug(src.ctx, "%s", title)
			service.Delete()
		} else {
			log.Debug(src.ctx, "%s (not found)", title)
		}

	case avahi.BrowserFailure:
		key := avahiServiceKeyFromServiceBrowserEvent(evnt)
		title := fmt.Sprintf("svc-browse: failed  %s", key)

		log.Warning(src.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// onServiceResolverEvent handles avahi.ServiceResolverEvent
func (src *source) onServiceResolverEvent(
	evnt *avahi.ServiceResolverEvent) error {
	switch evnt.Event {
	case avahi.ResolverFound:
		key := avahiServiceKeyFromResolverEvent(evnt)
		title := fmt.Sprintf("svc-resolve: found %s", key)

		service := src.clnt.GetService(key)
		if service == nil {
			// It may be out of order avahi.ResolverFound
			// event, received while service already removed,
			// so just log and return.
			log.Debug(src.ctx, "%s (unknown service)", title)
			return nil
		}

		log.Begin(src.ctx).
			Debug("%s:", title).
			Debug("  host: %s", evnt.Hostname).
			Debug("  port: %d", evnt.Port).
			Commit()

		service.SetPort(evnt.Port)
		src.setServiceHostname(service, evnt.Hostname)

	case avahi.ResolverFailure:
		key := avahiServiceKeyFromResolverEvent(evnt)
		title := fmt.Sprintf("svc-resolve: failed  %s", key)

		// Note, typically it's not fatal, just answer
		// doesn't want to come in time.
		log.Warning(src.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// onTxtBrowserEvent handles avahi.RecordBrowserEvent
// for per-service TXT record browser
func (src *source) onTxtBrowserEvent(evnt *avahi.RecordBrowserEvent) error {
	switch evnt.Event {
	case avahi.BrowserNew:
		key := avahiServiceKeyFromRecordBrowserEvent(evnt)
		title := fmt.Sprintf("txt-browse: found %s", key)
		log.Debug(src.ctx, "%s", title)

		service := src.clnt.GetService(key)
		if service == nil {
			log.Debug(src.ctx, "%s: service not found", title)
			return nil
		}

		svcType := key.SvcType
		svcInstance := key.InstanceName
		txt := avahi.DNSDecodeTXT(evnt.RData)

		if key.IsPrinter() {
			txtPrinter, err := decodeTxtPrinter(svcType,
				svcInstance, txt)
			if err != nil {
				log.Debug(src.ctx, "%s: %s", title, err)
				return nil // Don't propagate the error
			}

			if txtPrinter.IsPrinter() {
				id := key.PrinterUnitID(Backend, txtPrinter)
				un := service.GetUnit(id.Queue)
				if un == nil {
					un = newPrinterUnit(src.queue, id,
						txtPrinter)
					service.AddUnit(id.Queue, un)
				} else {
					un.SetTxtPrinter(txtPrinter)
				}
			}

			if txtPrinter.IsFaxout() {
				id := key.FaxoutUnitID(Backend, txtPrinter)
				un := service.GetUnit(id.Queue)
				if un == nil {
					un = newFaxoutUnit(src.queue, id,
						txtPrinter)
					service.AddUnit(id.Queue, un)
				} else {
					un.SetTxtFaxout(txtPrinter)
				}
			}
		} else {
			txtScanner, err := decodeTxtScanner(svcType,
				svcInstance, txt)
			if err != nil {
				log.Debug(src.ctx, "%s: %s", title, err)
				return nil // Don't propagate the error
			}

			unName := "scan"
			un := service.GetUnit(unName)
			if un == nil {
				id := key.ScannerUnitID(Backend, txtScanner)
				un = newScannerUnit(src.queue, id, txtScanner)
				service.AddUnit(unName, un)
			} else {
				un.SetTxtScanner(txtScanner)
			}
		}

	case avahi.BrowserFailure:
		key := avahiServiceKeyFromRecordBrowserEvent(evnt)
		title := fmt.Sprintf("txt-browse: failed %s", key)

		log.Warning(src.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// onAddrBrowserEvent handles avahi.RecordBrowserEvent
// for per-hostname A and AAAA record browsers
func (src *source) onAddrBrowserEvent(
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
		hostname := src.clnt.GetHostname(key)
		if hostname == nil {
			log.Debug(src.ctx, "%s: unknown hostname", title)
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
			log.Debug(src.ctx, "%s: invalid addr", title)
			return nil
		}

		// Add or delete the address
		hasAddr := hostname.HasAddr(addr)
		switch {
		case evnt.Event == avahi.BrowserNew && hasAddr:
			log.Debug(src.ctx, "%s: %s (duplicate)", title, addr)
			return nil
		case evnt.Event == avahi.BrowserRemove && !hasAddr:
			log.Debug(src.ctx, "%s: %s (unknown addr)", title,
				addr)
			return nil

		case evnt.Event == avahi.BrowserNew:
			hostname.AddAddr(addr, evnt.IfIdx)
		case evnt.Event == avahi.BrowserRemove:
			hostname.DelAddr(addr)
		}

		log.Debug(src.ctx, "%s: %s", title, addr)

	case avahi.BrowserFailure:
		title := fmt.Sprintf("addr-browse: failed %s", evnt.Name)

		log.Warning(src.ctx, "%s: %s", title, evnt.Err)
		return nil
	}

	return nil
}

// addService creates a new avahiService and registers it
// in the src.clnt
func (src *source) addService(key avahiServiceKey) error {
	// Create service resolver
	svcResolver, err := src.clnt.NewServiceResolver(key)

	title := fmt.Sprintf("svc-resolve: start %s", key)

	if err != nil {
		log.Error(src.ctx, "%s: %s", title, err)
		return err
	}

	log.Debug(src.ctx, "%s: OK", title)

	// Create TXT record browser
	txtBrowser, err := src.clnt.NewTxtBrowser(key)

	title = fmt.Sprintf("txt-browse: start %s", key)

	if err != nil {
		log.Error(src.ctx, "%s: %s", title, err)
		svcResolver.Close()
		return err
	}

	log.Debug(src.ctx, "%s: OK", title)

	// Add the service
	src.clnt.AddService(key, svcResolver, txtBrowser)
	return nil
}

// addHostname adds a new avahiHostname for the key
func (src *source) addHostname(key avahiHostnameKey) (*avahiHostname, error) {
	// Create A/AAAA record browser
	addrBrowser, err := src.clnt.NewAddrBrowser(key)

	title := fmt.Sprintf("addr-browse: start %s", key)

	if err != nil {
		log.Error(src.ctx, "%s: %s", title, err)
		return nil, err
	}

	log.Debug(src.ctx, "%s: OK", title)

	// Add avahiHostname
	return src.clnt.AddHostname(key, addrBrowser), nil
}

// setServiceHostname sets or updates the service's hostname.
//
// On success, it initiates hostname resolving, if it is not
// active yet.
func (src *source) setServiceHostname(service *avahiService,
	name string) error {

	// Do nothing if service already has a hostname and
	// name was not changed.
	if service.hostname != nil && service.hostname.key.Hostname == name {
		return nil
	}

	// Prepare avahiHostnameKey
	key := service.key.HostnameKey(name)

	// Find or create avahiHostname
	hostname := src.clnt.GetHostname(key)
	if hostname == nil {
		var err error
		hostname, err = src.addHostname(key)
		if err != nil {
			return err
		}
	}

	// Associate hostname with the service
	service.SetHostname(hostname)

	return nil
}
