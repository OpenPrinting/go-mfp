// MFP - Miulti-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Discovery client

package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
)

// Client implements a client side of devices discovery.
type Client struct {
	ctx      context.Context
	cancel   context.CancelFunc
	queue    *Eventqueue
	backends map[Backend]struct{}
	cache    *cache
	lock     sync.Mutex
	done     sync.WaitGroup
}

// NewClient creates a new discovery [Client].
//
// The provided [context.Context] is used for two purposes:
//   - For logging
//   - Client will terminate its operations, if context is canceled.
func NewClient(ctx context.Context) *Client {
	return NewClientTm(ctx, WarmUpTime, StabilizationTime)
}

// NewClientTm creates a new discovery [Client] with the
// warm-up time and stabilization time explicitly set.
//
// This interface is primary intended for testing but exported due
// to its general usability in some cases.
//
// Think carefully when choosing the time intervals, or use the
// simplified [NewClient] if not sure.
func NewClientTm(ctx context.Context,
	warmUpTime, stabilizationTime time.Duration) *Client {

	// Set log prefix
	ctx = log.WithPrefix(ctx, "discovery")

	// Create cancelable context
	ctx, cancel := context.WithCancel(ctx)

	// Create client structure
	clnt := &Client{
		ctx:      ctx,
		cancel:   cancel,
		queue:    NewEventqueue(),
		cache:    newCache(ctx, warmUpTime, stabilizationTime),
		backends: make(map[Backend]struct{}),
	}

	// Start work thread
	clnt.done.Add(1)
	go clnt.proc()

	return clnt
}

// Close closes all attached backends and then closes the Client
// and releases all resources it holds.
func (clnt *Client) Close() {
	// Close attached backends
	clnt.lock.Lock()
	backents := make([]Backend, 0, len(clnt.backends))
	for bk := range clnt.backends {
		backents = append(backents, bk)
	}
	clear(clnt.backends)
	clnt.lock.Unlock()

	for _, bk := range backents {
		bk.Close()
	}

	// Close the client itself
	clnt.cancel()
	clnt.done.Wait()
}

// AddBackend adds a discovery [Backend] to the [Client].
func (clnt *Client) AddBackend(bk Backend) {
	clnt.lock.Lock()
	defer clnt.lock.Unlock()

	if _, found := clnt.backends[bk]; found {
		err := fmt.Errorf("backend %s already added", bk.Name())
		panic(err)
	}

	log.Debug(clnt.ctx, "%s: backend added", bk.Name())
	clnt.backends[bk] = struct{}{}
	bk.Start(clnt.queue)
}

// GetDevices returns a list of discovered devices.
//
// Depending on [Mode] parameter and present discovery state,
// it may wait for some time or return immediately.
//
// If GetDevices decides to wait, expiration of either Context,
// given to this function as argument, or Context, using as [NewClient]
// argument during the Client creation will cause this function to
// return immediately with the appropriate error. And this is the
// only case when error is returned.
func (clnt *Client) GetDevices(ctx context.Context, m Mode) ([]*Device, error) {
	// Lock the client
	clnt.lock.Lock()
	defer clnt.lock.Unlock()

	// If snapshot is requested, take it immediately
	if m == ModeSnapshot {
		return clnt.cache.Snapshot(), nil
	}

	// Wait until ready
	ready := clnt.cache.ReadyAt(m)
	now := time.Now()
	for ready.After(now) {
		// As OS sleep is imprecise, pause for a slightly more
		// time to avoid spurious wakeups
		delay := ready.Sub(now) + time.Millisecond
		timer := time.NewTimer(delay)
		var err error

		clnt.lock.Unlock()
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-clnt.ctx.Done():
			err = clnt.ctx.Err()
		case now = <-timer.C:
		}
		clnt.lock.Lock()

		timer.Stop()
		if err != nil {
			return nil, err
		}
	}

	// And now read the cache
	return clnt.cache.Export(), nil
}

// GetByDNSSD returns a single device, selected by its DNS-SD name.
//
// See [GetDevices] for the [Mode] parameter explanation.
func (clnt *Client) GetByDNSSD(ctx context.Context,
	name string, m Mode) (*Device, error) {

	devices, err := clnt.GetDevices(ctx, m)
	if err != nil {
		return nil, err
	}

	for _, dev := range devices {
		if dev.DNSSDName == name {
			return dev, nil
		}
	}

	err = fmt.Errorf("%q: device not found", name)
	return nil, err
}

// Refresh causes [Client] to forcibly refresh its vision of
// discovered devices.
//
// The Refresh call returns immediately, but the subsequent call
// to the [Client.GetDevices] may wait until refresh completion,
// depending on mode.
func (clnt *Client) Refresh() {
}

// proc runs the discovery event loop on its separate goroutine.
func (clnt *Client) proc() {
	defer clnt.done.Done()

	var err error
	for err == nil {
		err = clnt.nextEvent()
	}
}

// nextEvent pulls and handles the next event
func (clnt *Client) nextEvent() error {
	evnt, err := clnt.queue.pull(clnt.ctx)
	if err != nil {
		return err
	}

	clnt.lock.Lock()
	defer clnt.lock.Unlock()

	rec := log.Begin(clnt.ctx)
	defer rec.Commit()

	rec.Trace("%s:", evnt.Name())
	rec.Object(log.LevelTrace, 2, evnt.GetID())

	switch evnt := evnt.(type) {
	case *EventAddUnit:
		err = clnt.cache.AddUnit(evnt)
	case *EventDelUnit:
		err = clnt.cache.DelUnit(evnt)
	case *EventPrinterParameters:
		err = clnt.cache.SetPrinterParameters(evnt)
	case *EventScannerParameters:
		err = clnt.cache.SetScannerParameters(evnt)
	case *EventFaxoutParameters:
		err = clnt.cache.SetFaxoutParameters(evnt)
	case *EventAddEndpoint:
		rec.Trace("  Endpoint:  %s", evnt.Endpoint)
		err = clnt.cache.AddEndpoint(evnt)
	case *EventDelEndpoint:
		rec.Trace("  Endpoint:  %s", evnt.Endpoint)
		err = clnt.cache.DelEndpoint(evnt)
	}

	if err != nil {
		// Log backend error and don't propagate it up the stack
		rec.Error("%s", err)
		err = nil
	}

	return err
}

// flush waits until all queued events are processed.
//
// This is the testing interface and it is not recommended
// for the general use
func (clnt *Client) flush() {
	clnt.lock.Lock()
	defer clnt.lock.Unlock()

	for {
		switch {
		case clnt.queue.Count() == 0:
			return
		case clnt.ctx.Err() != nil:
			return
		}

		clnt.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
		clnt.lock.Lock()
	}
}
