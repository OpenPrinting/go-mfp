// MFP - Miulti-Function Printers and scanners toolkit
// Utility functions and data BLOBs for testing
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Test data examples for Kyocera printers

package testutils

import (
	// Import "embed" for its side effects
	_ "embed"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// Kyocera contains data samples taken from the Kyocera printers
var Kyocera struct {
	// ECOSYS series
	ECOSYS struct {
		// M2040dn model
		M2040dn struct {
			// IPP protocol samples
			IPP struct {
				PrinterAttributes []byte
			}
			// ESCL protocol samples
			ESCL struct {
				ScannerCapabilities []byte
				ScannerStatus       []byte
			}
			// WS-Scan protocol samples
			WSD struct {
				GetScannerElementsResponse []byte
			}
			// USB samples
			USB struct {
				DeviceDescriptor usb.DeviceDescriptor
			}
		}
	}
}

func init() {
	Kyocera.ECOSYS.M2040dn.IPP.PrinterAttributes =
		kyoceraECOSYSM2040dnPrinterAttributes
	Kyocera.ECOSYS.M2040dn.ESCL.ScannerCapabilities =
		kyoceraECOSYSM2040dnScannerCapabilities
	Kyocera.ECOSYS.M2040dn.ESCL.ScannerStatus =
		kyoceraECOSYSM2040dnScannerStatus
	Kyocera.ECOSYS.M2040dn.WSD.GetScannerElementsResponse =
		kyoceraECOSYSM2040dnWSDGetScannerElementsResponse
	Kyocera.ECOSYS.M2040dn.USB.DeviceDescriptor = kyoceraECOSYSM2040dnUSBDeviceDescriptor

}

//go:embed "data/Kyocera-ECOSYS-M2040dn-Printer-Attributes.ipp"
var kyoceraECOSYSM2040dnPrinterAttributes []byte

//go:embed "data/Kyocera-ECOSYS-M2040dn-ScannerCapabilities.xml"
var kyoceraECOSYSM2040dnScannerCapabilities []byte

//go:embed "data/Kyocera-ECOSYS-M2040dn-ScannerStatus.xml"
var kyoceraECOSYSM2040dnScannerStatus []byte

//go:embed "data/Kyocera-ECOSYS-M2040dn-WSD-GetScannerElementsResponse.xml"
var kyoceraECOSYSM2040dnWSDGetScannerElementsResponse []byte

// kyoceraECOSYSM2040dnUSBDeviceDescriptor is the captured USB
// Device Descriptor for the Kyocera ECOSYS M2040dn MFP.
var kyoceraECOSYSM2040dnUSBDeviceDescriptor = usb.DeviceDescriptor{
	BCDUSB:          0x200,
	Speed:           3,
	BDeviceClass:    0x0,
	BDeviceSubClass: 0x0,
	BDeviceProtocol: 0x0,
	BMaxPacketSize:  0x40,
	IDVendor:        0x482,
	IDProduct:       0x69d,
	BCDDevice:       0x0,
	IManufacturer:   "Kyocera",
	IProduct:        "Kyocera ECOSYS M2040dn",
	ISerialNumber:   "VCF9192281",
	Configurations: []usb.ConfigurationDescriptor{
		{
			BConfigurationValue: 0x1,
			IConfiguration:      "",
			BMAttributes:        0xc0,
			BMaxPower:           0x1,
			Interfaces: []usb.Interface{
				{
					BInterfaceNumber: 0x0,
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    0x7,
							BInterfaceSubClass: 0x1,
							BInterfaceProtocol: 0x2,
							BAlternateSetting:  0x0,
							IInterface:         "",
							IEEE1284DeviceID:   "ID:ECOSYS M2040dn;MFG:Kyocera;CMD:PCLXL,PostScript Emulation,PCL5E,PJL;MDL:ECOSYS M2040dn;CLS:PRINTER;DES:Kyocera ECOSYS M2040dn;CID:KY_XPS_MonoA4FID;SER:VCF9192281;",
							Endpoints: []usb.EndpointDescriptor{
								{
									Type:           usb.EndpointOut,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
								{
									Type:           usb.EndpointIn,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
							},
						},
						{
							BInterfaceClass:    0x7,
							BInterfaceSubClass: 0x1,
							BInterfaceProtocol: 0x4,
							BAlternateSetting:  0x1,
							IInterface:         "",
							IEEE1284DeviceID:   "",
							Endpoints: []usb.EndpointDescriptor{
								{
									Type:           usb.EndpointOut,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
								{
									Type:           usb.EndpointIn,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
							},
						},
					},
				},
				{
					BInterfaceNumber: 0x1,
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    0xff,
							BInterfaceSubClass: 0xff,
							BInterfaceProtocol: 0xff,
							BAlternateSetting:  0x0,
							IInterface:         "",
							IEEE1284DeviceID:   "",
							Endpoints: []usb.EndpointDescriptor{
								{
									Type:           usb.EndpointOut,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
								{
									Type:           usb.EndpointIn,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
							},
						},
						{
							BInterfaceClass:    0x7,
							BInterfaceSubClass: 0x1,
							BInterfaceProtocol: 0x4,
							BAlternateSetting:  0x1,
							IInterface:         "",
							IEEE1284DeviceID:   "",
							Endpoints: []usb.EndpointDescriptor{
								{
									Type:           usb.EndpointOut,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
								{
									Type:           usb.EndpointIn,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
							},
						},
					},
				},
				{
					BInterfaceNumber: 0x2,
					AltSettings: []usb.InterfaceDescriptor{
						{
							BInterfaceClass:    0xff,
							BInterfaceSubClass: 0xff,
							BInterfaceProtocol: 0xff,
							BAlternateSetting:  0x0,
							IInterface:         "",
							IEEE1284DeviceID:   "",
							Endpoints:          []usb.EndpointDescriptor{},
						},
						{
							BInterfaceClass:    0x7,
							BInterfaceSubClass: 0x1,
							BInterfaceProtocol: 0x4,
							BAlternateSetting:  0x1,
							IInterface:         "",
							IEEE1284DeviceID:   "",
							Endpoints: []usb.EndpointDescriptor{
								{
									Type:           usb.EndpointOut,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
								{
									Type:           usb.EndpointIn,
									BMAttributes:   0x2,
									WMaxPacketSize: 0x200,
								},
							},
						},
					},
				},
			},
		},
	},
}
