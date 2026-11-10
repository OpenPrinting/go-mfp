// MFP - Multi-Function Printers and scanners toolkit
// The "mfp-usb" command
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions

package command

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/usbip"
)

const (
	// defaultIP is the default server address and port.
	defaultIP = "localhost:3240"
	// defaultBUSID is the default USB bus ID for the virtual device.
	defaultBUSID = "1-1"
	// frequency is the duration between consecutive server availability checks.
	frequency = 2 * time.Second
	// lockDir is the absolute path to the directory where process lock files are stored.
	lockDir = "/var/run/mfp-usb"
)

// AttachLock represents an exclusive advisory lock for a specific server address (host:port).
type AttachLock struct {
	file *os.File
}

// unlockFile releases the advisory lock and closes the underlying file descriptor.
func (l *AttachLock) unlockFile() {
	if l.file != nil {
		_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		_ = l.file.Close()
	}
}

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
	host, port, _ := net.SplitHostPort(ip)

	busid := defaultBUSID
	if inputBUSID, ok := inv.Get("-b"); ok {
		busid = inputBUSID
	}

	lock, err := isAttached(ip)
	if err != nil {
		return err
	}
	defer lock.unlockFile()

	sm, err := newSleepMonitor(ctx)
	if err != nil {
		return err
	}
	defer sm.Close()

	return runMonitor(ctx, sm.Chan(), ip, host, port, busid)
}

// runMonitor executes the main monitoring cycle for server availability
// and sleep events, automatically attaching or detaching the device as needed.
func runMonitor(ctx context.Context, sleepCh <-chan sleepSignal, ip, host, port, busid string) error {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	deathCh := make(chan struct{}, 1)
	deviceAvailable := false
	isSleeping := false

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "Shutting down...")
			_ = detach(ctx)
			return nil

		// Sleep event from D-Bus
		case signal := <-sleepCh:
			if signal.err != nil {
				log.Error(ctx, "Sleep monitor connection lost: %v", signal.err)
				_ = detach(ctx)
				return fmt.Errorf("sleep monitor failed %w", signal.err)
			}
			if signal.isSleep {
				log.Debug(ctx, "System is going to sleep. Detaching...")
				if err := detach(ctx); err != nil {
					log.Debug(ctx, "Warning: failed to detach: %v", err)
				}
				deviceAvailable = false
				isSleeping = true
			} else {
				log.Debug(ctx, "System woke up.")
				isSleeping = false
			}

		case <-deathCh:
			log.Info(ctx, "Server connection lost. Detaching...")
			if err := detach(ctx); err != nil {
				log.Debug(ctx, "Warning: failed to detach: %v", err)
			}
			deviceAvailable = false

		case <-ticker.C:
			if isSleeping || deviceAvailable {
				continue
			}

			conn, err := net.DialTimeout("tcp", ip, 500*time.Millisecond)
			if err != nil {
				continue
			}
			conn.Close()

			available, err := isDeviceAvailable(ctx, ip, busid)
			if err != nil {
				log.Debug(ctx, "Device availability check failed: %v", err)
				continue
			}

			if available {
				log.Info(ctx, "Server is up and device is available. Attaching...")
				if err := attach(host, busid); err != nil {
					log.Debug(ctx, "Attach failed: %v", err)
				} else {
					deviceAvailable = true
					log.Info(ctx, "Attached successfully")
					go waitServerDeath(ctx, deathCh, ip)
				}
			}
		}
	}
}

// waitServerDeath monitors server liveness via a blocking TCP "canary" connection.
// It avoids periodic polling and exec overhead, signaling the channel immediately
// when the server terminates (connection closed) or the context is cancelled.
func waitServerDeath(ctx context.Context, ch chan struct{}, ip string) {
	conn, err := net.DialTimeout("tcp", ip, time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 1)
		// blocking while server is alive
		_, err := conn.Read(buf)
		done <- err
	}()

	select {
	case <-done:
		//server died
		select {
		case ch <- struct{}{}:
		default:
		}
	case <-ctx.Done():
		return
	}
}

// isDeviceAvailable checks if the specified busid is exported on the remote USB/IP server.
// It resolves the host to an IP to satisfy the client's netip.AddrPort requirement.
func isDeviceAvailable(ctx context.Context, ip, busid string) (bool, error) {
	// Resolve host to a valid TCP address.
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		return false, fmt.Errorf("resolve failed: %w", err)
	}

	// Parse the resolved address to satisfy netip.AddrPort.
	addr, err := netip.ParseAddrPort(tcpAddr.String())
	if err != nil {
		return false, fmt.Errorf("parse failed: %w", err)
	}

	// Dial the server and request the device list.
	clnt, err := usbip.Dial(ctx, addr)
	if err != nil {
		return false, fmt.Errorf("dial failed: %w", err)
	}
	defer clnt.Close()

	entries, err := clnt.List()
	if err != nil {
		return false, fmt.Errorf("List() failed: %w", err)
	}

	// Return true if the target busid is found in the list.
	for _, entry := range entries {
		if entry.Busid == busid {
			return true, nil
		}
	}

	return false, nil
}

// isAttached acquires an exclusive, non-blocking file lock for the given IP to prevent concurrent management.
// The caller must invoke lock.unlockFile() (via defer) to release the lock upon exit.
func isAttached(ip string) (*AttachLock, error) {
	// Ensure the lock directory exists.
	if err := os.MkdirAll(lockDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create lock directory %s: %w", lockDir, err)
	}

	// Sanitize the IP string to create a safe filename.
	safeName := strings.ReplaceAll(ip, ":", "_") + ".lock"
	lockPath := filepath.Join(lockDir, safeName)

	// Open or create the lock file.
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file %s: %w", lockPath, err)
	}

	// Attempt to acquire an exclusive, non-blocking lock.
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == syscall.EWOULDBLOCK {
		file.Close()
		return nil, fmt.Errorf("another mfp-usb instance is already managing %s", ip)
	}
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to lock file %s: %w", lockPath, err)
	}

	return &AttachLock{file: file}, nil
}

// attach loads the vhci_hcd kernel module and attaches the virtual USB device
// using the specified bus ID.
func attach(ip, busID string) error {
	if out, err := exec.Command("modprobe", "vhci_hcd").CombinedOutput(); err != nil {
		return fmt.Errorf("modprobe error: %w, output: %s", err, string(out))
	}

	cmd := exec.Command("usbip", "attach", "-r", ip, "-b", busID)
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
