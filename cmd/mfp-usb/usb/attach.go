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
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/godbus/dbus/v5"
)

const (
	// defaultIP is the default server address and port.
	defaultIP = "localhost:3240"
	// defaultBUSID is the default USB bus ID for the virtual device.
	defaultBUSID = "1-1"
	// timeout is the maximum duration allowed for a single TCP connection attempt.
	timeout = 500 * time.Millisecond
	// frequency is the duration between consecutive server availability checks.
	frequency = 500 * time.Millisecond
)

// cmdAttach defines the "attach" command that automatically attaches
// mfp-virtual in USB mode with sleep/wake event handling.
var cmdAttach = argv.Command{
	Name:    "attach",
	Help:    "Auto attach mfp-virtual in USB mode",
	Handler: cmdAttachHandler,
	Options: []argv.Option{
		{
			Name:      "-i",
			Aliases:   []string{"--ip"},
			HelpArg:   "address",
			Help:      fmt.Sprintf("Server address. Default: %s", defaultIP),
			Singleton: true,
			Validate:  validateAddress,
		},
		{
			Name:      "-b",
			Aliases:   []string{"--busid"},
			HelpArg:   "busid",
			Help:      fmt.Sprintf("USB Bus ID. Default: %s", defaultBUSID),
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		argv.HelpOption,
	},
}

// cmdAttachHandler parses command-line arguments and starts the monitoring process.
func cmdAttachHandler(ctx context.Context, inv *argv.Invocation) error {
	// Check for root privileges before doing anything
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command requires root privileges, please run with sudo")
	}

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
			log.Debug(ctx, "Device already attached")
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
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	isSleeping := false

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "Shutting down...")
			if isAvailable {
				_ = detach(ctx)
			}
			return nil

		case errMsg := <-errCh:
			if isAvailable {
				log.Debug(ctx, "monitorSleepEvents Error")
				_ = detach(ctx)
			}
			return errMsg

		// Sleep event from D-Bus
		case isSleep := <-sleepCh:
			if isSleep {
				log.Debug(ctx, "System is going to sleep. Detaching...")
				if err := detach(ctx); err != nil {
					log.Debug(ctx, "Warning: failed to detach: %v", err)
				}
				isAvailable = false
				isSleeping = true
			} else {
				log.Debug(ctx, "System woke up.")
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
				log.Debug(ctx, "Server is up. Attempting to attach device...")
				if err := attach(busid); err != nil {
					log.Debug(ctx, "Attach failed: %v", err)
				} else {
					isAvailable = true
					log.Debug(ctx, "attached")
				}
			} else if !serverUp && isAvailable {
				log.Debug(ctx, "Server is down. Detaching...")
				if err := detach(ctx); err != nil {
					log.Debug(ctx, "Warning: failed to detach: %v", err)
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
	if out, err := exec.Command("modprobe", "vhci_hcd").CombinedOutput(); err != nil {
		return fmt.Errorf("modprobe error: %w, output: %s", err, string(out))
	}

	cmd := exec.Command("usbip", "attach", "-r", "localhost", "-b", busID)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usbip attach error: %w, output: %s", err, string(out))
	}

	return nil
}

// detach disconnects the virtual USB device from port 0.
func detach(ctx context.Context) error {
	cmd := exec.Command("usbip", "detach", "-p", "0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usbip detach error: %w, output: %s", err, string(out))
	}

	log.Debug(ctx, "detached")
	return nil
}

// validateAddress checks whether the given string represents a valid network address
func validateAddress(value string) error {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	_, err = net.ResolveTCPAddr("tcp", value)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	return nil
}
