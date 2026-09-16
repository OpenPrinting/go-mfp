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

// sleepMonitor manages the D-Bus connection and monitors sleep/wake events.
type sleepMonitor struct {
	conn     *dbus.Conn        // Active D-Bus system connection.
	ch       chan sleepSignal  // Output channel for sleep/wake events and errors
	doneCh   chan struct{}     // Сhan for synchronization
	dbusChan chan *dbus.Signal // Chan for D-Bus signals
}

type sleepSignal struct {
	isSleep bool
	err     error
}

// newSleepMonitor creates a new SleepMonitor, connects to D-Bus,
// and starts the background monitoring goroutine
func newSleepMonitor(ctx context.Context) (*sleepMonitor, error) {
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

	sm := &sleepMonitor{
		conn:     conn,
		ch:       make(chan sleepSignal, 1),
		doneCh:   make(chan struct{}),
		dbusChan: dbusChan,
	}

	// Start the background goroutine
	go sm.run(ctx)

	return sm, nil
}

// run listens for D-Bus signals and sends them to the event channel
func (sm *sleepMonitor) run(ctx context.Context) {
	defer close(sm.doneCh)
	defer close(sm.ch)
	for {
		select {
		case <-ctx.Done():
			return
		case signal, ok := <-sm.dbusChan:
			if !ok {
				sm.ch <- sleepSignal{err: fmt.Errorf("D-Bus connection lost")}
				return
			}
			if signal.Name == "org.freedesktop.login1.Manager.PrepareForSleep" && len(signal.Body) > 0 {
				if isSleeping, ok := signal.Body[0].(bool); ok {
					select {
					case sm.ch <- sleepSignal{isSleep: isSleeping, err: nil}:
					default:
					}
				}
			}
		}
	}
}

// Chan returns the channel from which sleep events can be read
func (sm *sleepMonitor) Chan() <-chan sleepSignal {
	return sm.ch
}

// Close gracefully shuts down the sleep monitor and guarantees that the
// background goroutine is completely terminated before returning.
//
// Exact sequence of events:
//  1. sm.conn.Close() closes the underlying D-Bus socket connection
//  2. The godbus library automatically closes the sm.dbusChan channel as a result
//  3. The background run() goroutine, which is blocked reading from sm.dbusChan,
//     wakes up and receives ok == false, indicating the channel is closed
//  4. The run() goroutine exits its loop and executes its deferred functions,
//     most importantly: defer close(sm.doneCh)
//  5. The main thread, which is blocked waiting on <-sm.doneCh, is unblocked
//     and the Close() method safely returns, preventing any goroutine leaks
func (sm *sleepMonitor) Close() {

	sm.conn.Close()

	<-sm.doneCh
}
