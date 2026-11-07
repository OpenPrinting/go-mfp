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
	"fmt"
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
			Help:    "Verbose logging (-vv for very verbose)",
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
	level := log.LevelInfo
	switch {
	case len(inv.Values("-v")) > 1:
		level = log.LevelTrace
	case len(inv.Values("-v")) > 0:
		level = log.LevelVerbose
	case inv.Flag("-d"):
		level = log.LevelDebug
	}

	logger := log.NewLogger(level, log.Console)
	ctx = log.NewContext(ctx, logger)

	// Setup pager
	pager := env.NewPager()
	pager.Redirect()
	defer pager.Close()

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
	if len(devices) == 0 {
		fmt.Printf("No devices found.\n")
	}

	for _, dev := range devices {
		fmt.Printf("================================\n")

		fmt.Printf("  MakeModel:        %q\n", dev.MakeModel)
		fmt.Printf("  Location:         %q\n", dev.Location)
		fmt.Printf("  DNS-SD name:      %q\n", dev.DNSSDName)
		fmt.Printf("  DNS-SD UUID:      %q\n", dev.DNSSDUUID)
		fmt.Printf("  WSD UUID:         %q\n", dev.WSDUUID)
		fmt.Printf("  Print Admin URL:  %s\n", dev.PrintAdminURL)
		fmt.Printf("  Scan Admin URL:   %s\n", dev.ScanAdminURL)
		fmt.Printf("  Faxout Admin URL: %s\n", dev.FaxoutAdminURL)
		fmt.Printf("  Icon URL:         %s\n", dev.IconURL)
		fmt.Printf("  PPD Manufacturer: %q\n", dev.PPDManufacturer)
		fmt.Printf("  PPD Model:        %q\n", dev.PPDModel)
		fmt.Printf("  USB Serial:       %q\n", dev.USBSerial)
		fmt.Printf("  USB HWID:         %q\n", dev.USBHWID)

		s := []string{}
		for _, addr := range dev.Addrs {
			s = append(s, addr.String())
		}
		fmt.Printf("  IP addresses:     %s\n", strings.Join(s, ", "))
		fmt.Printf("\n")

		if len(dev.PrintUnits) != 0 {
			fmt.Printf("  Print units:\n")
			for i, un := range dev.PrintUnits {
				if i != 0 {
					fmt.Printf("\n")
				}

				p := un.Params

				fmt.Printf("    Type:       %s printer\n", un.Proto)
				fmt.Printf("    Auth:       %s\n", p.Auth)

				if p.Paper != discovery.PaperUnknown {
					fmt.Printf("    Paper Size: %s\n",
						p.Paper)
				}

				if p.Media != 0 {
					fmt.Printf("    Media Type: %s\n",
						p.Media)
				}

				fmt.Printf("    Flags:      %s\n", p.Flags())

				fmt.Printf("    PSProduct:  %q\n", p.PSProduct)
				fmt.Printf("    PDL:        %s\n",
					strings.Join(p.PDL, ","))
				fmt.Printf("    Queue:      %s\n", p.Queue)
				fmt.Printf("    Priority:   %d\n", p.Priority)

				fmt.Printf("    Endpoints:\n")
				for _, ep := range un.Endpoints {
					fmt.Printf("      %s\n", ep)
				}

				fmt.Printf("    TXT:        %q\n", un.TXT)
			}
			fmt.Printf("\n")
		}

		if len(dev.ScanUnits) != 0 {
			fmt.Printf("  Scan units:\n")
			for i, un := range dev.ScanUnits {
				if i != 0 {
					fmt.Printf("\n")
				}

				p := un.Params
				fmt.Printf("    Type:       %s scanner\n",
					un.Proto)
				if p.Duplex != nil {
					fmt.Printf("    Duplex:     %v\n",
						*p.Duplex)
				}
				if !p.Inputs.IsEmpty() {
					fmt.Printf("    Inputs:     %s\n",
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
					fmt.Printf("    ColorModes: %s\n",
						strings.Join(modes, ","))
				}
				if len(p.PDL) != 0 {
					fmt.Printf("    PDL:        %s\n",
						strings.Join(p.PDL, ","))
				}

				fmt.Printf("    Endpoints:\n")
				for _, ep := range un.Endpoints {
					fmt.Printf("      %s\n", ep)
				}

				fmt.Printf("    TXT:        %q\n", un.TXT)
			}
			fmt.Printf("\n")
		}

		if len(dev.FaxoutUnits) != 0 {
			fmt.Printf("  Faxout units:\n")
			for i, un := range dev.FaxoutUnits {
				if i != 0 {
					fmt.Printf("\n")
				}

				p := un.Params

				fmt.Printf("    Type:       %s fax\n",
					un.Proto)
				fmt.Printf("    Auth:       %s\n", p.Auth)
				fmt.Printf("    Paper Size: %s\n", p.Paper)
				fmt.Printf("    Media Type: %s\n", p.Media)

				fmt.Printf("    Flags:      %s\n", p.Flags())

				fmt.Printf("    PSProduct:  %q\n", p.PSProduct)
				fmt.Printf("    PDL:        %s\n",
					strings.Join(p.PDL, ","))
				fmt.Printf("    Queue:      %s\n", p.Queue)
				fmt.Printf("    Priority:   %d\n", p.Priority)

				fmt.Printf("    Endpoints:\n")
				for _, ep := range un.Endpoints {
					fmt.Printf("      %s\n", ep)
				}

				fmt.Printf("    TXT:        %q\n", un.TXT)
			}
			fmt.Printf("\n")
		}
	}

	return nil
}
