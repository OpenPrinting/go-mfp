// MFP - Miulti-Function Printers and scanners toolkit
// The "virtual" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Virtual MFP simulator

package virtual

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/discovery/dnssd"
	"github.com/OpenPrinting/go-mfp/internal/env"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/OpenPrinting/go-mfp/transport/urilib"
)

// simulate runs scanner simulator.
//
// If argv is not empty, it specifies the external command that will
// be run under the simulator.
func simulate(ctx context.Context, model *modeling.Model,
	portmin, portmax uint16, usbip bool, argv []string) error {

	// Create the PathMux
	mux := transport.NewPathMux()

	// Virtual scanner template, common for all protocols
	templateScanner := abstract.VirtualScanner{
		Resolution: abstract.Resolution{
			XResolution: 600,
			YResolution: 600,
		},
		PlatenImage: testutils.Images.PNG5100x7016,
		ADFImages: [][]byte{
			testutils.Images.PNG5100x7016,
			testutils.Images.PNG5100x7016,
			testutils.Images.PNG5100x7016,
		},
	}

	// Add eSCL scanner
	esclcaps := model.GetESCLScanCaps()
	if esclcaps != nil {
		s := templateScanner
		s.ScanCaps = esclcaps.ToAbstract()

		handler := model.NewESCLServer(&s)
		mux.Add("/eSCL", handler)
	}

	// Add WS-Scan scanner
	wsdcaps := model.GetWSDScanCaps()
	if wsdcaps != nil {
		s := templateScanner
		s.ScanCaps = wsdcaps.ToAbstract()

		handler := model.NewWSDServer(&s)
		mux.Add("/WSScan", handler)
	}

	// Add IPP scanner
	ippScanCaps := model.GetIPPScannerAttrs()
	if ippScanCaps != nil {
		s := templateScanner
		s.ScanCaps = ippScanCaps.ToAbstractScannerCapabilities()

		handler := model.NewIPPScanner(&s)
		mux.Add("/ipp/scan", handler)
	}

	// Add IPP printer
	if handler := model.NewIPPPrinter(); handler != nil {
		mux.Add("/ipp/print", handler)
	}

	// Check that we have added at least something
	if mux.Empty() {
		return errors.New("model is empty")
	}

	// Create server for incoming connections.
	httpport := 0
	if !usbip {
		// Allocate TCP ports
		ports, err := model.OpenTCPPorts(netip.Addr{}, portmin, portmax)
		if err != nil {
			return err
		}

		// Create TLS certificate
		cert := model.GetTLSCertificate()
		template := http.Server{
			TLSConfig: &tls.Config{
				// Allow TLS 1.2 and 1.3
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,

				// Let clients use their proffered
				// cipher suites
				PreferServerCipherSuites: false,

				// Don't require TLS authentication
				ClientAuth: tls.NoClientCert,

				// Allow common curves
				CurvePreferences: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				},

				// Specify the TLS certificate
				Certificates: []tls.Certificate{cert},
			},
		}

		// Start servers
		log.Info(ctx, "Starting protocol servers:")
		if ports.HTTPListener != nil {
			httpport = int(ports.HTTPPort)

			srvr := transport.NewServer(ctx, &template, mux)
			addr := fmt.Sprintf("localhost:%d", ports.HTTPPort)
			log.Info(ctx, "  http://%s", addr)
			go srvr.ServeAutoTLS(ports.HTTPListener)

			defer srvr.Close()
		}

		if ports.LPDListener != nil {
			addr := fmt.Sprintf("localhost:%d", ports.LPDPort)
			log.Info(ctx, "  lpd://%s", addr)

			defer ports.LPDListener.Close()
		}

		if ports.AppSocketListener != nil {
			addr := fmt.Sprintf("localhost:%d", ports.AppSocketPort)
			log.Info(ctx, "  socket://%s", addr)

			defer ports.AppSocketListener.Close()
		}

		// Start DNS-SD advertising
		if dnssddev := model.GetDNSSDDevice(); dnssddev != nil {
			dnssddev = dnssdDeviceRewrite(dnssddev, ports)
			pub := dnssd.NewPublisher(ctx, dnssddev)
			defer pub.Close()
		}
	} else {
		desc := model.GetUSBDeviceDescriptor()
		if desc == nil {
			return errors.New("model lack USB support")
		}

		addr := &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 3240,
		}

		log.Info(ctx, "starting USBIP server at %s", addr)
		log.Info(ctx, "to connect the USB printer, run the following commands:")
		log.Info(ctx, "  sudo modprobe vhci-hcd")
		log.Info(ctx, "  sudo usbip attach -r localhost -b 1-1")

		_, err := newUsbipServer(ctx, desc, addr, mux)
		if err != nil {
			return err
		}
	}

	// Run external command if specified
	if len(argv) != 0 {
		runner := env.Runner{}
		if httpport > 0 {
			if esclcaps != nil {
				runner.ESCLName = "Virtual MFP Scanner"
				runner.ESCLPort = httpport
				runner.ESCLPath = "/eSCL"
			}

			if wsdcaps != nil {
				runner.WSDName = "Virtual MFP Scanner"
				runner.WSDPort = httpport
				runner.WSDPath = "/WSScan"
			}

			if model.GetIPPPrinterAttrs() != nil {
				runner.CUPSPort = httpport
			}
		}

		return runner.Run(ctx, argv[0], argv[1:]...)
	}

	// Wait for termination signal
	<-ctx.Done()
	log.Info(ctx, "Exiting...")

	return nil
}

// dnssdDeviceRewrite rewrites DNSSDDevice to point to the simulator's
// host and port.
func dnssdDeviceRewrite(dnssddev *discovery.DNSSDDevice,
	ports *modeling.TCPPorts) *discovery.DNSSDDevice {

	dnssddev = dnssddev.Clone()
	for i := range dnssddev.Services {
		svc := &dnssddev.Services[i]

		out := 0
		for _, ep := range svc.Endpoints {
			u := urilib.New(ep)
			scheme := strings.ToLower(urilib.Scheme(ep))
			port := uint16(0)

			switch {
			case u.IsHTTP():
				port = ports.HTTPPort
			case scheme == "lpd":
				port = ports.LPDPort
			case scheme == "socket":
				port = ports.AppSocketPort
			}

			if port != 0 {
				u = u.WithPortNum(port)
				if u.IsIP4() {
					u = u.WithHostname("127.0.0.1")
				} else {
					u = u.WithHostname("::1")
				}

				svc.Endpoints[out] = string(u)
				out++
			}
		}
		svc.Endpoints = svc.Endpoints[:out]
	}

	return dnssddev
}
