// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// TCP ports information

package modeling

import (
	"net"
	"net/netip"
	"strings"

	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/transport/urilib"
)

// TCPPorts represents TCP ports allocation, used
// by the [Model].
type TCPPorts struct {
	HTTPPort          uint16       // HTTP port number
	HTTPListener      net.Listener // HTTP port Listener
	LPDPort           uint16       // LPD port number
	LPDListener       net.Listener // LPD port Listener
	AppSocketPort     uint16       // AppSocket port number
	AppSocketListener net.Listener // AppSocket port Listener
}

// OpenTCPPorts opens TCP ports, as defined by model.
// The min and max parameters specify available TCP port range.
func (model *Model) OpenTCPPorts(addr netip.Addr, min, max uint16) (
	*TCPPorts, error) {

	// Compute ports
	withHTTP := false
	withLPD := false
	withSOCKET := false

	hasHTTPServices := model.GetIPPPrinterAttrs() != nil ||
		model.GetIPPScannerAttrs() != nil ||
		model.GetESCLScanCaps() != nil ||
		model.GetWSDScanCaps() != nil

	if dnssddev := model.GetDNSSDDevice(); dnssddev != nil {
		for _, svc := range dnssddev.Services {
			for _, ep := range svc.Endpoints {
				scheme := strings.ToLower(urilib.Scheme(ep))
				switch {
				case urilib.IsHTTP(ep):
					withHTTP = hasHTTPServices
				case scheme == "lpd":
					withLPD = true
				case scheme == "socket":
					withSOCKET = true
				}
			}
		}
	} else {
		withHTTP = hasHTTPServices
	}

	portcount := 0
	for _, with := range []bool{withHTTP, withLPD, withSOCKET} {
		if with {
			portcount++
		}
	}

	// Allocate ports
	portnum, listeners, err := transport.ListenTCPRangeN(netip.Addr{},
		portcount, min, max)

	if err != nil {
		return nil, err
	}

	// Fill output structure
	ports := &TCPPorts{}
	if withHTTP {
		ports.HTTPPort = portnum
		ports.HTTPListener = listeners[0]
		listeners = listeners[1:]
		portnum++
	}

	if withLPD {
		ports.LPDPort = portnum
		ports.LPDListener = listeners[0]
		listeners = listeners[1:]
		portnum++
	}

	if withSOCKET {
		ports.AppSocketPort = portnum
		ports.AppSocketListener = listeners[0]
		listeners = listeners[1:]
		portnum++
	}

	return ports, nil
}
