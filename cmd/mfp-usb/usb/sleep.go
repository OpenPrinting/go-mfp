// MFP - Multi-Function Printers and scanners toolkit
// The "usb" command
//
// Copyright (C) 2024 and up by Egor Ivashuta (ivashuta.ega@mail.ru)
// See LICENSE for license terms and conditions

package usb

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
)

// SleepMonitor manages the D-Bus connection and monitors sleep/wake events.
type SleepMonitor struct {
	conn     *dbus.Conn
	ch       chan bool
	doneCh   chan struct{}     // Сhan for synchronization
	dbusChan chan *dbus.Signal // Chan for D-Bus signals
}

// newSleepMonitor creates a new SleepMonitor, connects to D-Bus,
// and starts the background monitoring goroutine
func newSleepMonitor(ctx context.Context) (*SleepMonitor, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to D-Bus: %w", err)
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.login1.Manager"),
		dbus.WithMatchMember("PrepareForSleep"),
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to subscribe to D-Bus: %w", err)
	}

	dbusChan := make(chan *dbus.Signal, 10)
	conn.Signal(dbusChan)

	sm := &SleepMonitor{
		conn:     conn,
		ch:       make(chan bool, 1),
		doneCh:   make(chan struct{}),
		dbusChan: dbusChan,
	}

	// Start the background goroutine
	go sm.run(ctx)

	return sm, nil
}

// run listens for D-Bus signals and sends them to the event channel
func (sm *SleepMonitor) run(ctx context.Context) {
	defer close(sm.doneCh)
	defer close(sm.ch)
	for {
		select {
		case <-ctx.Done():
			return
		case signal, ok := <-sm.dbusChan:
			if !ok {
				// D-Bus connection lost
				return
			}
			if signal.Name == "org.freedesktop.login1.Manager.PrepareForSleep" && len(signal.Body) > 0 {
				if isSleeping, ok := signal.Body[0].(bool); ok {
					select {
					case sm.ch <- isSleeping:
					default:
					}
				}
			}
		}
	}
}

// Chan returns the channel from which sleep events can be read
func (sm *SleepMonitor) Chan() <-chan bool {
	return sm.ch
}

// Close stops the monitoring goroutine and waits for it to finish.
func (sm *SleepMonitor) Close() {

	sm.conn.Close()

	// Wait for the goroutine to finish
	<-sm.doneCh
}
