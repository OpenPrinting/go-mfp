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
	"errors"
	"fmt"
	"net"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/internal/env"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/transport"
)

// simulate runs scanner simulator.
//
// If argv is not empty, it specifies the external command that will
// be run under the simulator.
func simulate(ctx context.Context, model *modeling.Model,
	portnum int, usbip bool, argv []string) error {

	// Create the PathMux
	mux := transport.NewPathMux()
	runner := env.Runner{}

	// Add eSCL handler
	if esclcaps := model.GetESCLScanCaps(); esclcaps != nil {
		s := &abstract.VirtualScanner{
			ScanCaps: esclcaps.ToAbstract(),
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

		handler := model.NewESCLServer(s)
		mux.Add("/eSCL", handler)

		runner.ESCLName = "Virtual MFP Scanner"
		runner.ESCLPort = portnum
		runner.ESCLPath = "/eSCL"
	}

	// Add WS-Scan handler
	if wsdcaps := model.GetWSDScanCaps(); wsdcaps != nil {
		s := &abstract.VirtualScanner{
			ScanCaps: wsdcaps.ToAbstract(),
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

		handler := model.NewWSDServer(s)
		mux.Add("/WSScan", handler)

		runner.WSDName = "Virtual MFP Scanner"
		runner.WSDPort = portnum
		runner.WSDPath = "/WSScan"
	}

	// Add IPP handler
	if handler := model.NewIPPServer(); handler != nil {
		mux.Add("/ipp/print", handler)
		runner.CUPSPort = portnum
	}

	// Check that we have added at least something
	if mux.Empty() {
		return errors.New("model is emoty")
	}

	// Create server for incoming connections.
	if !usbip {
		addr := fmt.Sprintf("localhost:%d", portnum)

		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		srvr := transport.NewServer(ctx, nil, mux)
		log.Info(ctx, "starting virtual MFP at http://%s", addr)
		go srvr.Serve(ln)

		defer srvr.Close()
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
		return runner.Run(ctx, argv[0], argv[1:]...)
	}

	// Wait for termination signal
	<-ctx.Done()
	log.Info(ctx, "Exiting...")

	return nil
}
