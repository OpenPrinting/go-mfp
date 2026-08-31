// MFP - Multi-Function Printers and scanners toolkit
// The "usb" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
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
	"github.com/godbus/dbus/v5"
)

const (
	DefaultIP    = "localhost:3240"
	DefaultBUSID = "1-1"
	Timeout      = 500 * time.Millisecond
	Interval     = 2 * time.Second
)

var cmdAuto = argv.Command{
	Name:    "auto",
	Help:    "Auto attach for mfp-virtual in USB mode",
	Handler: cmdAutoHandler,
	Options: []argv.Option{
		{
			Name:      "-ip",
			Aliases:   []string{"--ip"},
			HelpArg:   "IP adress with a port",
			Help:      fmt.Sprintf("IP adress. Default: %s", DefaultIP),
			Singleton: true,
		},
		{
			Name:      "-b",
			Aliases:   []string{"--busid"},
			HelpArg:   "Busid for usbip",
			Help:      fmt.Sprintf("Default: %s", DefaultBUSID),
			Singleton: true,
		},
		argv.HelpOption,
	},
}

func cmdAutoHandler(ctx context.Context, inv *argv.Invocation) error {

	IP := DefaultIP
	if inputIP, ok := inv.Get("-ip"); ok {
		IP = inputIP
	}

	BUSID := DefaultBUSID
	if inputBUSID, ok := inv.Get("-b"); ok {
		BUSID = inputBUSID
	}

	// Check for device is already attached
	isAvailable := false
	cmd := exec.Command("usbip", "port")
	if output, err := cmd.CombinedOutput(); err == nil {
		if strings.Contains(string(output), "Port 00:") {
			fmt.Println("Device already attached")
			isAvailable = true
		}
	}

	sleepCh := make(chan bool, 1)
	go monitorSleepEvents(ctx, sleepCh)

	ticker := time.NewTicker(Interval)
	defer ticker.Stop()

	isSleeping := false

	for {
		select {
		case <-ctx.Done():
			if isAvailable {
				fmt.Println("Shutting down...")
				_ = detach()
			}
			return ctx.Err()

		// Sleep event from dbus
		case isSleep := <-sleepCh:
			if isSleep {
				fmt.Println("System is going to sleep. Detaching...")
				if err := detach(); err != nil {
					fmt.Printf("Warning: failed to detach: %v\n", err)
				}
				isAvailable = false
				isSleeping = true
			} else {
				fmt.Println("System woke up.")
				isSleeping = false
			}

		case <-ticker.C:

			// skip attach attempts if system is sleeping
			if isSleeping {
				continue
			}

			// check if server is available
			serverUp := false
			if conn, err := net.DialTimeout("tcp", IP, Timeout); err == nil {
				conn.Close()
				serverUp = true
			}

			if serverUp && !isAvailable {
				fmt.Println("Server is up. Attempting to attach device...")
				if err := attach(BUSID); err != nil {
					fmt.Printf("Attach failed: %v\n", err)
				} else {
					isAvailable = true
					fmt.Println("attached")
				}
			} else if !serverUp && isAvailable {
				fmt.Println("Server is down. Detaching...")
				if err := detach(); err != nil {
					fmt.Printf("Warning: failed to detach: %v\n", err)
				}
				isAvailable = false
			}
		}
	}
}

// listens to system D-Bus for sleep signals
// Send 'true' when system is going to sleep, 'false' when waking up
func monitorSleepEvents(ctx context.Context, sleepCh chan<- bool) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		fmt.Printf("Warning: failed to connect to D-Bus: %v. Sleep events will be ignored.\n", err)
		return
	}
	defer conn.Close()

	err = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.login1.Manager"),
		dbus.WithMatchMember("PrepareForSleep"),
	)
	if err != nil {
		fmt.Printf("Warning: failed to subscribe to D-Bus: %v\n", err)
		return
	}

	dbusChan := make(chan *dbus.Signal, 10)
	conn.Signal(dbusChan)

	for {
		select {
		case <-ctx.Done():
			return
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

func attach(BUSID string) error {
	if out, err := exec.Command("sudo", "modprobe", "vhci_hcd").CombinedOutput(); err != nil {
		return fmt.Errorf("modprobe error: %w, output: %s", err, string(out))
	}

	cmd := exec.Command("sudo", "usbip", "attach", "-r", "localhost", "-b", BUSID)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usbip attach error: %w, output: %s", err, string(out))
	}

	return nil
}

func detach() error {
	cmd := exec.Command("sudo", "usbip", "detach", "-p", "0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usbip detach error: %w, output: %s", err, string(out))
	}

	fmt.Println("detached")
	return nil
}
