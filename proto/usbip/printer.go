// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Virtual printer device

package usbip

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/ieee1284"
	"github.com/OpenPrinting/go-mfp/proto/usb"
	"github.com/OpenPrinting/go-mfp/transport"
)

// NewPrinter creates the new USBIP virtual printer.
//
// It forwards the IPP over USB requests to the provided
// [http.Handler] and legacy 7/1/1 and 7/1/2 print requests
// to the provided [ieee1284.Printer]s.
func NewPrinter(ctx context.Context,
	desc usb.DeviceDescriptor,
	handler http.Handler,
	printers []*ieee1284.Printer) (*Device, error) {

	// Validate things
	match711 := usb.ClassID{Class: 7, SubClass: 1, Protocol: 1}
	match712 := usb.ClassID{Class: 7, SubClass: 1, Protocol: 2}
	match714 := usb.ClassID{Class: 7, SubClass: 1, Protocol: 4}

	cnt711 := desc.CntMatch(match711)
	cnt712 := desc.CntMatch(match712)
	cnt714 := desc.CntMatch(match714)

	if cnt711+cnt712 > len(printers) {
		err := fmt.Errorf("not enough ieee1284.Printer (wand %d have %d)",
			cnt711+cnt712, len(printers))
		return nil, err
	}

	if cnt714 > 9 && handler == nil {
		err := errors.New("device has 7/1/4 interface but handler not provided")
		return nil, err
	}

	// Create USB device
	dev, err := NewDevice(desc)
	if err != nil {
		return nil, err
	}

	// Initialize I/O on IEEE-1284 interfaces
	if cnt711+cnt712 > 0 {
		match := func(alt usb.InterfaceDescriptor) bool {
			return alt.Match(match711) || alt.Match(match712)
		}
		endpoints := dev.EndpointsByFunc(match)

		for i, endpoint := range endpoints {
			if i < len(printers) {
				go runLegacyPrinter(ctx, endpoint,
					printers[i])
			}
		}
	}

	// IPP over USB is the HTTP-based protocol.
	// Create HTTP server on a top of the USB device endpoints.
	if cnt714 > 0 {
		srv := transport.NewServer(ctx, nil, handler)

		endpoints := dev.EndpointsByClassID(match714)
		listener := NewEndpointListener(endpoints)

		go srv.Serve(listener)
	}

	return dev, nil
}

// runLegacyPrinter pipes data between a USB endpoint and an
// IEEE 1284 printer in both directions:
//   - Host to printer: endpoint.Read() → printer.Write()
//   - Printer to host: printer.Read() → endpoint.Write()
func runLegacyPrinter(ctx context.Context, endpoint *Endpoint,
	printer *ieee1284.Printer) {

	// Host → Printer
	go func() {
		buf := make([]byte, 512)
		for {
			n, err := endpoint.Read(buf)
			if err != nil {
				log.Debug(ctx, "usbip: legacy printer: "+
					"endpoint read: %s", err)
				printer.Close()
				return
			}
			printer.Write(buf[:n])
		}
	}()

	// Printer → Host
	go func() {
		buf := make([]byte, 512)
		for {
			n, err := printer.Read(buf)
			if err != nil {
				log.Debug(ctx, "usbip: legacy printer: "+
					"printer read: %s", err)
				return
			}
			endpoint.Write(buf[:n])
		}
	}()
}
