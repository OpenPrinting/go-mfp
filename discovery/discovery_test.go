// MFP  - Multi-Function Printers and scanners toolkit
// discovery - Discovery module test suite
//
// Copyright (C) 2025 and up by SinghCod3r
// See LICENSE for license terms and conditions
//
// Test suite for discovery functionality

package discovery

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// mockBackend is a mock implementation of the Backend interface used
// for testing.
//
// It simulates a backend driver that emits events to the discovery queue.
type mockBackend struct{}

// mockSource implements Source interface for mockBackend.
type mockSource struct {
	queue *Eventqueue
}

// mockContextKey used as a key for context.WithValue to
// pass events to the mockSource
type mockContextKey struct{}

// Name returns backend name.
func (mockBackend) Name() string {
	return "mock-backend"
}

// DeviceID returns a subset of UnitID, that uniquely
// identifies a physical device.
func (mockBackend) DeviceID(id UnitID) UnitID {
	return UnitID{
		DNSSDName: id.DNSSDName,
		Backend:   id.Backend,
	}
}

// Open creates a new [discovery.Source] for DNS-SD discovery.
func (mockBackend) Open(ctx context.Context,
	queue *Eventqueue) (Source, error) {

	src := &mockSource{queue: queue}
	if events := ctx.Value(mockContextKey{}); events != nil {
		for _, e := range events.([]Event) {
			queue.Push(e)
		}
	}

	return src, nil
}

// Close cleans up backend resources. For the mock, this is a no-op.
func (mb *mockSource) Close() {
	// No-op for mock
}

// AddEvent appends an event to the list of events that the backend will emit upon starting.
func (mb *mockBackend) AddEvent(e Event) {
	//mb.events = append(mb.events, e)
}

// TestClient_NoDevices verifies that GetDevices returns an empty list
// when no devices are discovered.
func TestClient_NoDevices(t *testing.T) {
	ctx := context.Background()
	client, _ := NewClientTm(
		ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		&mockBackend{})

	defer client.Close()

	devices, err := client.GetDevices(ctx, ModeNormal)
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	if len(devices) != 0 {
		t.Errorf("Expected 0 devices, got %d", len(devices))
	}
}

// TestClient_Discovery verifies the successful discovery of a printer device
// when a backend emits valid AddUnit, PrinterParameters, and AddEndpoint events.
func TestClient_Discovery(t *testing.T) {
	// Prepare events
	uid := UnitID{
		DNSSDName: "Test Printer",
		UUID:      uuid.Random(),
		Backend:   mockBackend{},
		SvcType:   ServicePrinter,
		SvcProto:  ServiceIPP,
	}

	events := []Event{
		&EventAddUnit{ID: uid},
		&EventPrinterParameters{
			ID:        uid,
			MakeModel: "Test Make Model",
			Printer: PrinterParameters{
				Queue: "test-queue",
			},
		},
		&EventAddEndpoint{
			ID:        uid,
			Endpoints: []string{"ipp://192.168.1.100/ipp/print"},
		},
	}

	// Open the client
	ctx := context.Background()
	ctx = context.WithValue(ctx, mockContextKey{}, events)
	client, _ := NewClientTm(ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		mockBackend{})
	defer client.Close()

	// Wait for discovery to complete (WarmUpTime + processing)
	client.flush()

	devices, err := client.GetDevices(ctx, ModeNormal)
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	if len(devices) != 1 {
		t.Errorf("Expected 1 device, got %d", len(devices))
	} else {
		dev := devices[0]
		if dev.MakeModel != "Test Make Model" {
			t.Errorf("Expected MakeModel 'Test Make Model', got '%s'", dev.MakeModel)
		}
	}
}

// TestClient_InvalidEvents verifies the client's robustness against duplicate or unknown events.
// It checks that such events do not cause panics or incorrect device listings.
func TestClient_InvalidEvents(t *testing.T) {
	// Prepare events
	uid := UnitID{
		DNSSDName: "Test Printer",
		UUID:      uuid.Random(),
		Backend:   mockBackend{},
		SvcType:   ServicePrinter,
		SvcProto:  ServiceIPP,
	}

	unknownUID := UnitID{DNSSDName: "Unknown", UUID: uuid.Random()}

	events := []Event{
		// 1. Duplicate EventAddUnit
		&EventAddUnit{ID: uid},
		&EventAddUnit{ID: uid}, // Should be handled gracefully (logged error)

		// 2. EventPrinterParameters for unknown unit
		&EventPrinterParameters{
			ID:        unknownUID,
			MakeModel: "Unknown",
		},

		// 3. EventDelUnit for unknown unit
		&EventDelUnit{ID: unknownUID},
	}

	// Open the client
	ctx := context.Background()
	ctx = context.WithValue(ctx, mockContextKey{}, events)
	client, _ := NewClientTm(ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		mockBackend{})
	defer client.Close()

	client.flush()
	devices, err := client.GetDevices(ctx, ModeNormal)
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	if len(devices) != 0 {
		t.Errorf("Expected 0 devices, got %d", len(devices))
	}
}

// TestClient_ContextCancel verifies that the client handles context cancellation appropriately.
func TestClient_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client, _ := NewClientTm(ctx, 5*time.Second,
		StabilizationTime, mockBackend{})
	defer client.Close()

	// Cancel context immediately
	cancel()

	_, err := client.GetDevices(ctx, ModeNormal)
	if err == nil {
		t.Error("Expected error due to context cancellation, got nil")
	}
}

// TestClient_Timeout verifies that the client returns a deadline exceeded error
// when the context times out before discovery completes.
func TestClient_Timeout(t *testing.T) {
	ctx := context.Background()
	client, _ := NewClientTm(ctx, 5*time.Second,
		StabilizationTime, mockBackend{})
	defer client.Close()

	// Create a context with a short timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	_, err := client.GetDevices(timeoutCtx, ModeNormal)
	if err == nil {
		t.Error("Expected error due to timeout, got nil")
	} else if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded, got %v", err)
	}
}

// TestClient_MissingFields verifies behavior when events are missing optional fields (like MakeModel).
func TestClient_MissingFields(t *testing.T) {
	// Prepare events
	uid := UnitID{
		DNSSDName: "Test Printer",
		UUID:      uuid.Random(),
		Backend:   mockBackend{},
		SvcType:   ServicePrinter,
		SvcProto:  ServiceIPP,
	}

	events := []Event{
		&EventAddUnit{ID: uid},
		// Missing MakeModel: explicit check for empty MakeModel scenario
		&EventPrinterParameters{
			ID:        uid,
			MakeModel: "", // Empty
			Printer: PrinterParameters{
				Queue: "test-queue",
			},
		},
		&EventAddEndpoint{
			ID:        uid,
			Endpoints: []string{"ipp://192.168.1.100/ipp/print"},
		},
	}

	// Open the client
	ctx := context.Background()
	ctx = context.WithValue(ctx, mockContextKey{}, events)
	client, _ := NewClientTm(ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		mockBackend{})
	defer client.Close()

	client.flush()

	devices, err := client.GetDevices(ctx, ModeNormal)
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	if len(devices) != 1 {
		t.Errorf("Expected 1 device, got %d", len(devices))
	} else {
		if devices[0].MakeModel != "" {
			t.Errorf("Expected empty MakeModel, got '%s'", devices[0].MakeModel)
		}
	}
}

// okBackend119 is a Backend that opens successfully, returning a
// mockSource bound to the discovery queue.
type okBackend119 struct{}

func (okBackend119) Name() string              { return "ok119" }
func (okBackend119) DeviceID(id UnitID) UnitID { return id }
func (okBackend119) Open(ctx context.Context, queue *Eventqueue) (Source, error) {
	return &mockSource{queue: queue}, nil
}

// failBackend119 is a Backend whose Open always fails.
type failBackend119 struct{}

func (failBackend119) Name() string              { return "fail119" }
func (failBackend119) DeviceID(id UnitID) UnitID { return id }
func (failBackend119) Open(ctx context.Context, queue *Eventqueue) (Source, error) {
	return nil, errors.New("boom119")
}

// countSource119 counts Close calls, so tests can observe cleanup.
// (mockSource.Close is a no-op and cannot be observed.)
type countSource119 struct{ closed *atomic.Int32 }

func (s *countSource119) Close() { s.closed.Add(1) }

// countBackend119 opens successfully, returning an observable source.
type countBackend119 struct{ closed *atomic.Int32 }

func (countBackend119) Name() string                { return "count119" }
func (b countBackend119) DeviceID(id UnitID) UnitID { return id }
func (b countBackend119) Open(ctx context.Context, queue *Eventqueue) (Source, error) {
	return &countSource119{closed: b.closed}, nil
}

// TestNewClientTmAttachesSources verifies that successfully opened
// sources are retained for later Close (see issue #119).
func TestNewClientTmAttachesSources(t *testing.T) {
	ctx := context.Background()
	client, err := NewClientTm(ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		okBackend119{})
	if err != nil {
		t.Fatalf("NewClientTm failed: %v", err)
	}
	defer client.Close()

	if len(client.sources) != 1 {
		t.Errorf("len(sources) = %d, want 1", len(client.sources))
	}
}

// TestNewClientTmPropagatesOpenError verifies that a failing first
// backend surfaces its error instead of a healthy client (see #119).
func TestNewClientTmPropagatesOpenError(t *testing.T) {
	ctx := context.Background()
	client, err := NewClientTm(ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		failBackend119{})
	if err == nil {
		if client != nil {
			client.Close()
		}
		t.Fatalf("expected error, got nil")
	}
	if client != nil {
		t.Errorf("expected nil client on failure")
	}
}

// TestNewClientTmClosesPriorSources verifies that a later backend
// failure closes already-opened sources (see issue #119).
func TestNewClientTmClosesPriorSources(t *testing.T) {
	var closed atomic.Int32

	ctx := context.Background()
	client, err := NewClientTm(ctx,
		100*time.Millisecond,
		100*time.Millisecond,
		countBackend119{closed: &closed}, failBackend119{})
	if err == nil {
		if client != nil {
			client.Close()
		}
		t.Fatalf("expected error, got nil")
	}
	if client != nil {
		t.Errorf("expected nil client on failure")
	}
	if got := closed.Load(); got != 1 {
		t.Errorf("prior source closed %d times, want 1", got)
	}
}
