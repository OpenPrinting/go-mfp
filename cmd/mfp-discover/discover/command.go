// MFP - Miulti-Function Printers and scanners toolkit
// The "discover" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Command description.

package discover

import (
	"context"
	"strings"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/discovery/dnssd"
	"github.com/OpenPrinting/go-mfp/discovery/wsdd"
	"github.com/OpenPrinting/go-mfp/internal/env"
	"github.com/OpenPrinting/go-mfp/log"
)

// Command is the 'cups' command description
var Command = argv.Command{
	Name: "discover",
	Help: "search for printers and scanners",
	Options: []argv.Option{
		argv.Option{
			Name:    "-d",
			Aliases: []string{"--debug"},
			Help:    "Enable debug output",
		},
		argv.Option{
			Name:    "-v",
			Aliases: []string{"--verbose"},
			Help:    "Enable verbose debug output",
		},
		argv.Option{
			Name:    "-p",
			Aliases: []string{"--printers"},
			Help:    "Search for printers",
		},
		argv.Option{
			Name:    "-s",
			Aliases: []string{"--scanners"},
			Help:    "Search for scanners",
		},
		argv.HelpOption,
	},
	Handler: cmdDiscoverHandler,
}

// cmdCupsHandler is the handler for the 'discover' command.
func cmdDiscoverHandler(ctx context.Context, inv *argv.Invocation) error {
	// Setup logging
	_, dbg := inv.Get("-d")
	_, vrb := inv.Get("-v")

	dbg = true // FIXME

	level := log.LevelInfo
	if dbg {
		level = log.LevelDebug
	}
	if vrb {
		level = log.LevelTrace
	}

	logger := log.NewLogger(level, log.Console)
	ctx = log.NewContext(ctx, logger)

	// Prepare discovery.Client
	clnt, err := discovery.NewClient(ctx, dnssd.Backend, wsdd.Backend)
	if err != nil {
		return err
	}
	defer clnt.Close()

	// Perform device discovery
	devices, err := clnt.GetDevices(ctx, discovery.ModeNormal)
	if err != nil {
		return err
	}

	// Format output
	pager := env.NewPager()
	defer pager.Display()

	if len(devices) == 0 {
		pager.Printf("No devices found.")
	}

	for _, dev := range devices {
		pager.Printf("================================")

		pager.Printf("  MakeModel:        %q", dev.MakeModel)
		pager.Printf("  Location:         %q", dev.Location)
		pager.Printf("  DNS-SD name:      %q", dev.DNSSDName)
		pager.Printf("  DNS-SD UUID:      %q", dev.DNSSDUUID)
		pager.Printf("  WSD UUID:         %q", dev.WSDUUID)
		pager.Printf("  Print Admin URL:  %s", dev.PrintAdminURL)
		pager.Printf("  Scan Admin URL:   %s", dev.ScanAdminURL)
		pager.Printf("  Faxout Admin URL: %s", dev.FaxoutAdminURL)
		pager.Printf("  Icon URL:         %s", dev.IconURL)
		pager.Printf("  PPD Manufacturer: %q", dev.PPDManufacturer)
		pager.Printf("  PPD Model:        %q", dev.PPDModel)
		pager.Printf("  USB Serial:       %q", dev.USBSerial)
		pager.Printf("  USB HWID:         %q", dev.USBHWID)

		s := []string{}
		for _, addr := range dev.Addrs {
			s = append(s, addr.String())
		}
		pager.Printf("  IP addresses:     %s", strings.Join(s, ", "))
		pager.Printf("")

		if len(dev.PrintUnits) != 0 {
			pager.Printf("  Print units:")
			for i, un := range dev.PrintUnits {
				if i != 0 {
					pager.Printf("")
				}

				p := un.Params

				pager.Printf("    Type:       %s printer", un.Proto)
				pager.Printf("    Auth:       %s", p.Auth)

				if p.Paper != discovery.PaperUnknown {
					pager.Printf("    Paper Size: %s",
						p.Paper)
				}

				if p.Media != 0 {
					pager.Printf("    Media Type: %s",
						p.Media)
				}

				pager.Printf("    Flags:      %s", p.Flags())

				pager.Printf("    PSProduct:  %q", p.PSProduct)
				pager.Printf("    PDL:        %s",
					strings.Join(p.PDL, ","))
				pager.Printf("    Queue:      %s", p.Queue)
				pager.Printf("    Priority:   %d", p.Priority)

				pager.Printf("    Endpoints:")
				for _, ep := range un.Endpoints {
					pager.Printf("      %s", ep)
				}

				pager.Printf("    TXT:        %q", un.TXT)
			}
			pager.Printf("")
		}

		if len(dev.ScanUnits) != 0 {
			pager.Printf("  Scan units:")
			for i, un := range dev.ScanUnits {
				if i != 0 {
					pager.Printf("")
				}

				p := un.Params
				pager.Printf("    Type:       %s scanner",
					un.Proto)
				if p.Duplex != nil {
					pager.Printf("    Duplex:     %v",
						*p.Duplex)
				}
				if !p.Inputs.IsEmpty() {
					pager.Printf("    Inputs:     %s",
						p.Inputs)
				}
				if !p.Colors.IsEmpty() {
					var modes []string
					if p.Colors.Contains(
						abstract.ColorModeColor) {
						modes = append(modes, "color")
					}
					if p.Colors.Contains(
						abstract.ColorModeMono) {
						modes = append(modes, "mono")
					}
					if p.Colors.Contains(
						abstract.ColorModeBinary) {
						modes = append(modes, "bin")
					}
					pager.Printf("    ColorModes: %s",
						strings.Join(modes, ","))
				}
				if len(p.PDL) != 0 {
					pager.Printf("    PDL:        %s",
						strings.Join(p.PDL, ","))
				}

				pager.Printf("    Endpoints:")
				for _, ep := range un.Endpoints {
					pager.Printf("      %s", ep)
				}

				pager.Printf("    TXT:        %q", un.TXT)
			}
			pager.Printf("")
		}

		if len(dev.FaxoutUnits) != 0 {
			pager.Printf("  Faxout units:")
			for i, un := range dev.FaxoutUnits {
				if i != 0 {
					pager.Printf("")
				}

				p := un.Params

				pager.Printf("    Type:       %s fax",
					un.Proto)
				pager.Printf("    Auth:       %s", p.Auth)
				pager.Printf("    Paper Size: %s", p.Paper)
				pager.Printf("    Media Type: %s", p.Media)

				pager.Printf("    Flags:      %s", p.Flags())

				pager.Printf("    PSProduct:  %q", p.PSProduct)
				pager.Printf("    PDL:        %s",
					strings.Join(p.PDL, ","))
				pager.Printf("    Queue:      %s", p.Queue)
				pager.Printf("    Priority:   %d", p.Priority)

				pager.Printf("    Endpoints:")
				for _, ep := range un.Endpoints {
					pager.Printf("      %s", ep)
				}

				pager.Printf("    TXT:        %q", un.TXT)
			}
			pager.Printf("")
		}
	}

	return nil
}
