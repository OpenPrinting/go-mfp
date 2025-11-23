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

// MockBackend is a mock implementation of the Backend interface used for testing.
// It simulates a backend driver that emits events to the discovery queue.
type MockBackend struct {
	name   string
	queue  *Eventqueue
	events []Event
}

// NewMockBackend creates a new instance of MockBackend with the specified name.
func NewMockBackend(name string) *MockBackend {
	return &MockBackend{
		name:   name,
		events: make([]Event, 0),
	}
}

// Name returns the name of the mock backend.
func (mb *MockBackend) Name() string {
	return mb.name
}

// Start initializes the backend and starts a goroutine to push queued events
// to the event queue.
func (mb *MockBackend) Start(q *Eventqueue) {
	mb.queue = q
	go func() {
		for _, e := range mb.events {
			mb.queue.Push(e)
		}
	}()
}

// Close cleans up backend resources. For the mock, this is a no-op.
func (mb *MockBackend) Close() {
	// No-op for mock
}

// AddEvent appends an event to the list of events that the backend will emit upon starting.
func (mb *MockBackend) AddEvent(e Event) {
	mb.events = append(mb.events, e)
}

// TestClient_NoDevices verifies that GetDevices returns an empty list when no devices are discovered.
func TestClient_NoDevices(t *testing.T) {
	// Reduce WarmUpTime for testing to speed up execution
	originalWarmUpTime := warmUpTime
	originalStabilizationTime := stabilizationTime
	warmUpTime = 100 * time.Millisecond
	stabilizationTime = 100 * time.Millisecond
	defer func() {
		warmUpTime = originalWarmUpTime
		stabilizationTime = originalStabilizationTime
	}()

	ctx := context.Background()
	client := NewClient(ctx)
	defer client.Close()

	backend := NewMockBackend("mock-backend")
	client.AddBackend(backend)

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
	originalWarmUpTime := warmUpTime
	originalStabilizationTime := stabilizationTime
	warmUpTime = 100 * time.Millisecond
	stabilizationTime = 100 * time.Millisecond
	defer func() {
		warmUpTime = originalWarmUpTime
		stabilizationTime = originalStabilizationTime
	}()

	ctx := context.Background()
	client := NewClient(ctx)
	defer client.Close()

	backend := NewMockBackend("mock-backend")
	
	uid := UnitID{
		DNSSDName: "Test Printer",
		UUID:      uuid.Must(uuid.Random()),
		SvcType:   ServicePrinter,
		SvcProto:  ServiceIPP,
	}

	backend.AddEvent(&EventAddUnit{ID: uid})
	backend.AddEvent(&EventPrinterParameters{
		ID:        uid,
		MakeModel: "Test Make Model",
		Printer: PrinterParameters{
			Queue: "test-queue",
		},
	})
	backend.AddEvent(&EventAddEndpoint{
		ID:       uid,
		Endpoint: "ipp://192.168.1.100/ipp/print",
	})

	client.AddBackend(backend)

	// Wait for discovery to complete (WarmUpTime + processing)
	time.Sleep(200 * time.Millisecond)

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
	originalWarmUpTime := warmUpTime
	originalStabilizationTime := stabilizationTime
	warmUpTime = 100 * time.Millisecond
	stabilizationTime = 100 * time.Millisecond
	defer func() {
		warmUpTime = originalWarmUpTime
		stabilizationTime = originalStabilizationTime
	}()

	ctx := context.Background()
	client := NewClient(ctx)
	defer client.Close()

	backend := NewMockBackend("mock-backend")
	
	uid := UnitID{
		DNSSDName: "Test Printer",
		UUID:      uuid.Must(uuid.Random()),
		SvcType:   ServicePrinter,
		SvcProto:  ServiceIPP,
	}

	// 1. Duplicate EventAddUnit
	backend.AddEvent(&EventAddUnit{ID: uid})
	backend.AddEvent(&EventAddUnit{ID: uid}) // Should be handled gracefully (logged error)

	// 2. EventPrinterParameters for unknown unit
	unknownUID := UnitID{DNSSDName: "Unknown", UUID: uuid.Must(uuid.Random())}
	backend.AddEvent(&EventPrinterParameters{
		ID:        unknownUID,
		MakeModel: "Unknown",
	})

	// 3. EventDelUnit for unknown unit
	backend.AddEvent(&EventDelUnit{ID: unknownUID})

	client.AddBackend(backend)
	time.Sleep(200 * time.Millisecond)

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
	originalWarmUpTime := warmUpTime
	warmUpTime = 5 * time.Second // Long enough to block
	defer func() { warmUpTime = originalWarmUpTime }()

	ctx, cancel := context.WithCancel(context.Background())
	client := NewClient(ctx)
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
	originalWarmUpTime := warmUpTime
	warmUpTime = 5 * time.Second // Long enough to block
	defer func() { warmUpTime = originalWarmUpTime }()

	ctx := context.Background()
	client := NewClient(ctx)
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
	originalWarmUpTime := warmUpTime
	originalStabilizationTime := stabilizationTime
	warmUpTime = 100 * time.Millisecond
	stabilizationTime = 100 * time.Millisecond
	defer func() {
		warmUpTime = originalWarmUpTime
		stabilizationTime = originalStabilizationTime
	}()

	ctx := context.Background()
	client := NewClient(ctx)
	defer client.Close()

	backend := NewMockBackend("mock-backend")
	
	uid := UnitID{
		DNSSDName: "Test Printer",
		UUID:      uuid.Must(uuid.Random()),
		SvcType:   ServicePrinter,
		SvcProto:  ServiceIPP,
	}

	backend.AddEvent(&EventAddUnit{ID: uid})
	// Missing MakeModel: explicit check for empty MakeModel scenario
	backend.AddEvent(&EventPrinterParameters{
		ID:        uid,
		MakeModel: "", // Empty
		Printer: PrinterParameters{
			Queue: "test-queue",
		},
	})
	backend.AddEvent(&EventAddEndpoint{
		ID:       uid,
		Endpoint: "ipp://192.168.1.100/ipp/print",
	})

	client.AddBackend(backend)
	time.Sleep(200 * time.Millisecond)

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

// TestClient_Unreachable verifies that no devices are returned if the backend
// is unresponsive or provides no events, resulting in an empty discovery.
func TestClient_Unreachable(t *testing.T) {
	originalWarmUpTime := warmUpTime
	originalStabilizationTime := stabilizationTime
	warmUpTime = 100 * time.Millisecond
	stabilizationTime = 100 * time.Millisecond
	defer func() {
		warmUpTime = originalWarmUpTime
		stabilizationTime = originalStabilizationTime
	}()

	ctx := context.Background()
	client := NewClient(ctx)
	defer client.Close()

	// Backend that sends nothing
	backend := NewMockBackend("mock-backend")
	client.AddBackend(backend)

	time.Sleep(200 * time.Millisecond)

	devices, err := client.GetDevices(ctx, ModeNormal)
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	if len(devices) != 0 {
		t.Errorf("Expected 0 devices, got %d", len(devices))
	}
}
