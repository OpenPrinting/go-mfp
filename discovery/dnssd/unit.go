// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Discovered unit

package dnssd

import (
	"net/netip"
	"net/url"
	"strconv"

	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/util/generic"
)

// unit represents a discovered print or scan unit.
// It accepts RR updates and generates events.
type unit struct {
	queue   *discovery.Eventqueue   // Event queue
	id      discovery.UnitID        // Unit ID
	svcType string                  // Service type
	addrs   generic.Set[netip.Addr] // IP addresses of the unit
	port    uint16                  // IP port
	txtPrn  txtPrinter              // Parsed TXT for print unit
	txtScn  txtScanner              // Parsed TXT for scan unit
}

// newPrinterUnit creates a new printer unit
func newPrinterUnit(queue *discovery.Eventqueue,
	id discovery.UnitID, txt txtPrinter) *unit {

	un := &unit{
		queue:   queue,
		id:      id,
		svcType: txt.svcType,
		addrs:   generic.NewSet[netip.Addr](),
		txtPrn:  txt,
	}

	un.queue.Push(&discovery.EventAddUnit{ID: un.id})
	un.SetTxtPrinter(txt)

	return un
}

// newFaxoutUnit creates a new faxout unit
func newFaxoutUnit(queue *discovery.Eventqueue,
	id discovery.UnitID, txt txtPrinter) *unit {

	un := &unit{
		queue:   queue,
		id:      id,
		svcType: txt.svcType,
		addrs:   generic.NewSet[netip.Addr](),
		txtPrn:  txt,
	}

	un.queue.Push(&discovery.EventAddUnit{ID: un.id})
	un.SetTxtFaxout(txt)

	return un
}

// newScannerUnit creates a new scanner unit
func newScannerUnit(queue *discovery.Eventqueue,
	id discovery.UnitID, txt txtScanner) *unit {
	un := &unit{
		queue:   queue,
		id:      id,
		svcType: txt.svcType,
		addrs:   generic.NewSet[netip.Addr](),
		txtScn:  txt,
	}

	un.queue.Push(&discovery.EventAddUnit{ID: un.id})
	un.SetTxtScanner(txt)

	return un
}

// Delete deletes the unit
func (un *unit) Delete() {
	un.queue.Push(&discovery.EventDelUnit{ID: un.id})
}

// IsPrinter reports if unit is the print unit
func (un *unit) IsPrinter() bool {
	return un.txtPrn.params != nil
}

// setTxtPrinter changes TXT record for printer unit
func (un *unit) SetTxtPrinter(txt txtPrinter) {
	un.txtPrn = txt
	params := *txt.params
	params.Queue = txt.rp
	un.queue.Push(&discovery.EventPrinterParameters{
		ID:              un.id,
		TXT:             txt.txt,
		MakeModel:       txt.makeModel,
		Location:        txt.location,
		AdminURL:        txt.adminURL,
		IconURL:         txt.iconURL,
		PPDManufacturer: txt.usbMFG,
		PPDModel:        txt.usbMDL,
		Printer:         params,
	})
}

// setTxtFaxout changes TXT record for printer unit
func (un *unit) SetTxtFaxout(txt txtPrinter) {
	un.txtPrn = txt
	params := *txt.params
	params.Queue = txt.rp
	un.queue.Push(&discovery.EventFaxoutParameters{
		ID:              un.id,
		TXT:             txt.txt,
		MakeModel:       txt.makeModel,
		Location:        txt.location,
		AdminURL:        txt.adminURL,
		IconURL:         txt.iconURL,
		PPDManufacturer: txt.usbMFG,
		PPDModel:        txt.usbMDL,
		Faxout:          params,
	})
}

// setTxtScanner changes TXT record for scanner unit
func (un *unit) SetTxtScanner(txt txtScanner) {
	un.txtScn = txt
	un.queue.Push(&discovery.EventScannerParameters{
		ID:        un.id,
		TXT:       txt.txt,
		MakeModel: txt.makeModel,
		Location:  txt.location,
		AdminURL:  txt.adminURL,
		IconURL:   txt.iconURL,
		Scanner:   *txt.params,
	})
}

// addAddr sets unit's port
func (un *unit) SetPort(port uint16) {
	// Drop old endpoints
	if un.port != 0 {
		un.addrs.ForEach(func(addr netip.Addr) {
			evnt := &discovery.EventDelEndpoint{
				ID:       un.id,
				Endpoint: un.endpoint(addr),
			}
			un.queue.Push(evnt)
		})
	}

	un.port = port

	// Restore endpoints after port change
	if port != 0 {
		un.addrs.ForEach(func(addr netip.Addr) {
			evnt := &discovery.EventAddEndpoint{
				ID:       un.id,
				Endpoint: un.endpoint(addr),
			}
			un.queue.Push(evnt)
		})
	}
}

// addAddr adds unit's IP address
func (un *unit) AddAddr(addr netip.Addr) {
	if !un.addrs.Contains(addr) {
		un.addrs.Add(addr)
		if un.port != 0 {
			evnt := &discovery.EventAddEndpoint{
				ID:       un.id,
				Endpoint: un.endpoint(addr),
			}
			un.queue.Push(evnt)
		}
	}
}

// delAddr deletes unit's IP address
func (un *unit) DelAddr(addr netip.Addr) {
	if un.addrs.Contains(addr) {
		un.addrs.Del(addr)
		if un.port != 0 {
			evnt := &discovery.EventDelEndpoint{
				ID:       un.id,
				Endpoint: un.endpoint(addr),
			}
			un.queue.Push(evnt)
		}
	}
}

// endpoint creates an endpoint URL
func (un *unit) endpoint(addr netip.Addr) string {
	host := addr.WithZone("").String()
	if addr.Is6() {
		// We need to put address into square brackets and
		// append zone for the link-local addresses
		host = "[" + host
		if zone := addr.Zone(); zone != "" {
			host += "%" + zone
		}
		host += "]"
	}

	port := int(un.port)
	var url url.URL

	switch un.svcType {
	case svcTypeAppSocket:
		// socket://host[:port]
		// default port: 9100
		url.Scheme = "socket"
		url.Host = host
		if port != 9100 {
			url.Host += ":" + strconv.Itoa(port)
		}

	case svcTypeIPP, svcTypeIPPS:
		// ipp://host[:port]/queue or ipps://host[:port]/queue
		// default port: 631
		url.Scheme = "ipp"
		if un.svcType == svcTypeIPPS {
			url.Scheme = "ipps"
		}

		url.Host = host
		if port != 631 {
			url.Host += ":" + strconv.Itoa(port)
		}

		url.Path = un.id.Queue

	case svcTypeLPD:
		// lpd://host[:port]/queue
		// default port: 515
		url.Scheme = "lpd"
		url.Host = host
		if port != 515 {
			url.Host += ":" + strconv.Itoa(port)
		}
		url.Path = un.txtPrn.params.Queue

	case svcTypeESCL, svcTypeESCLS:
		// http://host[:port]/path or https://host[:port]/path
		// default port: 80
		url.Host = host
		if un.svcType == svcTypeESCL {
			url.Scheme = "http"
			if port != 80 {
				url.Host += ":" + strconv.Itoa(port)
			}
		} else {
			url.Scheme = "https"
			if port != 443 {
				url.Host += ":" + strconv.Itoa(port)
			}
		}
		url.Path = un.txtScn.uriPath
	}

	return url.String()
}
