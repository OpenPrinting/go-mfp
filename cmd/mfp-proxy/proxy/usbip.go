// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by Mohammed Imaduddin (mdimad005@gmail.com)
// See LICENSE for license terms and conditions
//
// IPP over USB emulation logic

package proxy

import (
	"context"
	"net"
	"net/http"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/proto/ieee1284"
	"github.com/OpenPrinting/go-mfp/proto/usb"
	"github.com/OpenPrinting/go-mfp/proto/usbip"
)

// newUsbipServer creates the new USBIP server representing
// an IPP over USB MFP device.
//
// The server accepts incoming USBIP connection on the provided
// address and forwards incoming IPP over USB requests (which
// are essentially the HTTP requests) to the provided http.Handler.
func newUsbipServer(ctx context.Context,
	addr net.Addr, handler http.Handler) *usbip.Server {

	// Obtain device descriptor
	desc := defaultUSBIPPDescriptor()

	// Create enough IEEE-1284 virtual printers
	match711 := usb.ClassID{Class: 7, SubClass: 1, Protocol: 1}
	match712 := usb.ClassID{Class: 7, SubClass: 1, Protocol: 2}

	n := desc.CntMatch(match711)
	n += desc.CntMatch(match712)

	ieeeprinters := make([]*ieee1284.Printer, n)
	for i := 0; i < n; i++ {
		ieeeprinters[i] = ieee1284.NewPrinter(ctx, nil)
	}

	// Create the USB printer device
	dev, err := usbip.NewPrinter(ctx, desc, handler, ieeeprinters)
	assert.NoError(err)

	// Create USBIP server
	usbipSrv := usbip.NewServer(ctx)
	err = usbipSrv.AddDevice(dev)
	assert.NoError(err)

	go usbipSrv.ListenAndServe(addr)

	return usbipSrv
}

// defaultUSBIPPDescriptor returns the [usb.DeviceDescriptor]
// for the typical IPP over USB device.
func defaultUSBIPPDescriptor() usb.DeviceDescriptor {
	ippusbEndpointIn := usb.EndpointDescriptor{
		Type:           usb.EndpointIn,
		BMAttributes:   usb.XferBulk,
		WMaxPacketSize: 512,
	}

	ippusbEndpointOut := usb.EndpointDescriptor{
		Type:           usb.EndpointOut,
		BMAttributes:   usb.XferBulk,
		WMaxPacketSize: 512,
	}

	const deviceID = "MFG:OpenPrinting;MDL:Virtual MFP;CMD:PJL,POSTSCRIPT;CLS:PRINTER;SERN:OP-0000223321;"

	desc := usb.DeviceDescriptor{
		BCDUSB:          0x0200,
		Speed:           usb.SpeedHigh,
		BDeviceClass:    0,
		BDeviceSubClass: 0,
		BDeviceProtocol: 0,
		BMaxPacketSize:  64,
		IDVendor:        0xdead,
		IDProduct:       0xbeaf,
		BCDDevice:       0x0100,
		IManufacturer:   "OpenPrinting",
		IProduct:        "Virtual MFP",
		ISerialNumber:   "NN-001122334455",
		Configurations: []usb.ConfigurationDescriptor{{
			BMAttributes: usb.ConfAttrSelfPowered,
			BMaxPower:    1,
			Interfaces: []usb.Interface{
				{
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    7,
							BInterfaceSubClass: 1,
							BInterfaceProtocol: 2,
							IEEE1284DeviceID:   deviceID,
							Endpoints: []usb.EndpointDescriptor{
								ippusbEndpointIn,
								ippusbEndpointOut,
							},
						},
					},
				},
				{
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    7,
							BInterfaceSubClass: 1,
							BInterfaceProtocol: 4,
							Endpoints: []usb.EndpointDescriptor{
								ippusbEndpointIn,
								ippusbEndpointOut,
							},
						},
					},
				},
				{
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    7,
							BInterfaceSubClass: 1,
							BInterfaceProtocol: 4,
							Endpoints: []usb.EndpointDescriptor{
								ippusbEndpointIn,
								ippusbEndpointOut,
							},
						},
					},
				},
				{
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    7,
							BInterfaceSubClass: 1,
							BInterfaceProtocol: 4,
							Endpoints: []usb.EndpointDescriptor{
								ippusbEndpointIn,
								ippusbEndpointOut,
							},
						},
					},
				},
			},
		}},
	}

	return desc
}
