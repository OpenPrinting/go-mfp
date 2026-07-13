// MFP - Miulti-Function Printers and scanners toolkit
// The "usb" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The "list" command.

package usb

import (
	"context"
	"strings"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/internal/env"
	"github.com/OpenPrinting/go-mfp/proto/usb"
	"github.com/OpenPrinting/go-mfp/proto/usbhost"
)

// cmdList defines the "list" sub-command
var cmdList = argv.Command{
	Name:    "list",
	Help:    "List connected devices",
	Handler: cmdListHandler,
	Options: []argv.Option{
		argv.Option{
			Name: "--all",
			Help: "List all USB devices, not only MFPs",
		},
		argv.HelpOption,
	},
}

// cmdListHandler is the "list" command handler.
func cmdListHandler(ctx context.Context, inv *argv.Invocation) error {
	list, err := usbhost.ListDevices(true)
	if err != nil {
		return err
	}

	all := inv.Flag("--all")
	pager := env.NewPager()
	for _, info := range list {
		desc := info.Desc

		if !all &&
			desc.BDeviceClass != 7 &&
			!desc.Contains(7, -1, -1) {
			continue
		}

		pager.Printf("Bus %3.3d Device %3.3d ID %4.4x:%4.4x",
			info.Loc.Bus, info.Loc.Dev,
			desc.IDVendor, desc.IDProduct)
		pager.Printf("  Manufacturer:    %s", desc.IManufacturer)
		pager.Printf("  Product:         %s", desc.IProduct)
		pager.Printf("  SerialNumber:    %s", desc.ISerialNumber)
		pager.Printf("  USB Version:     %s", desc.BCDUSB)
		pager.Printf("  Device Version:  %s", desc.BCDUSB)
		pager.Printf("  Class:           %d/%d/%d",
			desc.BDeviceClass,
			desc.BDeviceSubClass,
			desc.BDeviceProtocol)
		pager.Printf("  Speed:           %d", desc.Speed)
		pager.Printf("  Max Packet size: %d", desc.BMaxPacketSize)

		for confno, conf := range desc.Configurations {
			pager.Printf("")
			pager.Printf("  Configuration %d:", confno)
			pager.Printf("    description: %s", conf.IConfiguration)

			attrs := []string{}
			if conf.BMAttributes&usb.ConfAttrSelfPowered != 0 {
				attrs = append(attrs, "self powered")
			}

			if conf.BMAttributes&usb.ConfAttrRemoteWakeup != 0 {
				attrs = append(attrs, "remote wakeup")
			}

			pager.Printf("    Attributes:  0x%2.2x (%s)",
				conf.BMAttributes,
				strings.Join(attrs, ","))

			pager.Printf("    Max Power:   %d mA", int(conf.MaxPower)*2)

			for iffno, iff := range conf.Interfaces {
				pager.Printf("")
				pager.Printf("    Interface %d:", iffno)

				for altno, alt := range iff.AltSettings {
					pager.Printf("      Alt Setting %d:",
						altno)
					pager.Printf("        Class:       %d/%d/%d",
						alt.BInterfaceClass,
						alt.BInterfaceSubClass,
						alt.BInterfaceProtocol)
					pager.Printf("        Description: %s",
						alt.IInterface)
					pager.Printf("        Device ID:   %q",
						alt.IEEE1284DeviceID)
					pager.Printf("        Endpoints:")

					for _, ep := range alt.Endpoints {
						dir := "OUT:"
						if ep.Type == usb.EndpointIn {
							dir = "IN: "
						}

						xferbits := ep.BMAttributes
						xferbits &= usb.XferMask

						xfer := ""
						switch xferbits {
						case usb.XferControl:
							xfer = "control"
						case usb.XferIsochronous:
							xfer = "iso"
						case usb.XferBulk:
							xfer = "bulk"
						case usb.XferInterrupt:
							xfer = "interrupt"
						}

						pager.Printf(
							"          %s %s, pktsize=%d",
							dir, xfer,
							ep.WMaxPacketSize)
					}
				}
			}
		}

		pager.Printf("")
	}

	return pager.Display()
}
