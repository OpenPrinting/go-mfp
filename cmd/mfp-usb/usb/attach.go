// MFP - Multi-Function Printers and scanners toolkit
// The "usb" command
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions

package usb

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/godbus/dbus/v5"
)

const (
	defaultIP    = "localhost:3240"
	defaultBUSID = "1-1"
	timeout      = 500 * time.Millisecond
	interval     = 500 * time.Millisecond
)

// cmdAttach defines the "attach" command that automatically attaches
// mfp-virtual in USB mode with sleep/wake event handling.
var cmdAttach = argv.Command{
	Name:    "attach",
	Help:    "Auto attach mfp-virtual in USB mode",
	Handler: cmdAttachHandler,
	Options: []argv.Option{
		{
			Name:      "-d",
			Aliases:   []string{"--debug"},
			Help:      "Enable debug output",
			Singleton: true,
		},
		{
			Name:      "-i",
			Aliases:   []string{"--ip"},
			HelpArg:   "address",
			Help:      fmt.Sprintf("Server address. Default: %s", defaultIP),
			Singleton: true,
		},
		{
			Name:      "-b",
			Aliases:   []string{"--busid"},
			HelpArg:   "busid",
			Help:      fmt.Sprintf("USB Bus ID. Default: %s", defaultBUSID),
			Singleton: true,
		},
		argv.HelpOption,
	},
}

// cmdAttachHandler parses command-line arguments and starts the monitoring process.
func cmdAttachHandler(ctx context.Context, inv *argv.Invocation) error {
	ip := defaultIP
	if inputIP, ok := inv.Get("-i"); ok {
		ip = inputIP
	}

	busid := defaultBUSID
	if inputBUSID, ok := inv.Get("-b"); ok {
		busid = inputBUSID
	}

	// Check if device is already attached
	isAvailable := false
	cmd := exec.Command("usbip", "port")
	if output, err := cmd.CombinedOutput(); err == nil {
		if strings.Contains(string(output), "Port 00:") {
			log.Info(ctx, "Device already attached")
			isAvailable = true
		}
	}

	sleepCh := make(chan bool, 1)
	errCh := make(chan error, 1)
	go func() {
		errCh <- monitorSleepEvents(ctx, sleepCh)
	}()

	return runMonitor(ctx, sleepCh, errCh, ip, busid, isAvailable)
}

// runMonitor executes the main monitoring cycle for server availability
// and sleep events, automatically attaching or detaching the device as needed.
func runMonitor(ctx context.Context, sleepCh <-chan bool, errCh <-chan error, ip, busid string, isAvailable bool) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	isSleeping := false

	for {
		select {
		case <-ctx.Done():
			if isAvailable {
				log.Info(ctx, "Shutting down...")
				_ = detach(ctx)
			}
			return ctx.Err()

		case errMsg := <-errCh:
			if isAvailable {
				log.Info(ctx, "monitorSleepEvents Error")
				_ = detach(ctx)
			}
			return errMsg

		// Sleep event from D-Bus
		case isSleep := <-sleepCh:
			if isSleep {
				log.Info(ctx, "System is going to sleep. Detaching...")
				if err := detach(ctx); err != nil {
					log.Info(ctx, "Warning: failed to detach: %v", err)
				}
				isAvailable = false
				isSleeping = true
			} else {
				log.Info(ctx, "System woke up.")
				isSleeping = false
			}

		case <-ticker.C:
			// Skip attach attempts if system is sleeping
			if isSleeping {
				continue
			}

			// Check if server is available
			serverUp := false
			if conn, err := net.DialTimeout("tcp", ip, timeout); err == nil {
				conn.Close()
				serverUp = true
			}

			if serverUp && !isAvailable {
				log.Info(ctx, "Server is up. Attempting to attach device...")
				if err := attach(busid); err != nil {
					log.Info(ctx, "Attach failed: %v", err)
				} else {
					isAvailable = true
					log.Info(ctx, "attached")
				}
			} else if !serverUp && isAvailable {
				log.Info(ctx, "Server is down. Detaching...")
				if err := detach(ctx); err != nil {
					log.Info(ctx, "Warning: failed to detach: %v", err)
				}
				isAvailable = false
			}
		}
	}
}

// monitorSleepEvents monitors system D-Bus for sleep/wake signals.
// It sends true to sleepCh when the system is going to sleep,
// and false when the system is waking up.
func monitorSleepEvents(ctx context.Context, sleepCh chan<- bool) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled before D-Bus connection: %w", ctx.Err())
	default:
	}

	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return fmt.Errorf("failed to connect to D-Bus: %w", err)
	}
	defer conn.Close()

	err = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.login1.Manager"),
		dbus.WithMatchMember("PrepareForSleep"),
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to D-Bus: %w", err)
	}

	dbusChan := make(chan *dbus.Signal, 10)
	conn.Signal(dbusChan)

	for {
		select {
		case <-ctx.Done():
			return nil
		case signal := <-dbusChan:
			if signal.Name == "org.freedesktop.login1.Manager.PrepareForSleep" && len(signal.Body) > 0 {
				isSleeping := signal.Body[0].(bool)

				select {
				case sleepCh <- isSleeping:
				default:
				}
			}
		}
	}
}

// attach loads the vhci_hcd kernel module and attaches the virtual USB device
// using the specified bus ID.
func attach(busID string) error {
	if out, err := exec.Command("sudo", "modprobe", "vhci_hcd").CombinedOutput(); err != nil {
		return fmt.Errorf("modprobe error: %w, output: %s", err, string(out))
	}

	cmd := exec.Command("sudo", "usbip", "attach", "-r", "localhost", "-b", busID)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usbip attach error: %w, output: %s", err, string(out))
	}

	return nil
}

// detach disconnects the virtual USB device from port 0.
func detach(ctx context.Context) error {
	cmd := exec.Command("sudo", "usbip", "detach", "-p", "0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usbip detach error: %w, output: %s", err, string(out))
	}

	log.Info(ctx, "detached")
	return nil
}
