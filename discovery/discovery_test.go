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
