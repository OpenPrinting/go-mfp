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

// MockBackend is a mock implementation of the Backend interface
type MockBackend struct {
	name   string
	queue  *Eventqueue
	events []Event
}

func NewMockBackend(name string) *MockBackend {
	return &MockBackend{
		name:   name,
		events: make([]Event, 0),
	}
}

func (mb *MockBackend) Name() string {
	return mb.name
}

func (mb *MockBackend) Start(q *Eventqueue) {
	mb.queue = q
	go func() {
		for _, e := range mb.events {
			mb.queue.Push(e)
		}
	}()
}

func (mb *MockBackend) Close() {
	// No-op for mock
}

func (mb *MockBackend) AddEvent(e Event) {
	mb.events = append(mb.events, e)
}

func TestClient_NoDevices(t *testing.T) {
	// Reduce WarmUpTime for testing
	originalWarmUpTime := WarmUpTime
	originalStabilizationTime := StabilizationTime
	WarmUpTime = 100 * time.Millisecond
	StabilizationTime = 100 * time.Millisecond
	defer func() {
		WarmUpTime = originalWarmUpTime
		StabilizationTime = originalStabilizationTime
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

func TestClient_Discovery(t *testing.T) {
	originalWarmUpTime := WarmUpTime
	originalStabilizationTime := StabilizationTime
	WarmUpTime = 100 * time.Millisecond
	StabilizationTime = 100 * time.Millisecond
	defer func() {
		WarmUpTime = originalWarmUpTime
		StabilizationTime = originalStabilizationTime
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

func TestClient_InvalidEvents(t *testing.T) {
	originalWarmUpTime := WarmUpTime
	originalStabilizationTime := StabilizationTime
	WarmUpTime = 100 * time.Millisecond
	StabilizationTime = 100 * time.Millisecond
	defer func() {
		WarmUpTime = originalWarmUpTime
		StabilizationTime = originalStabilizationTime
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

func TestClient_ContextCancel(t *testing.T) {
	originalWarmUpTime := WarmUpTime
	WarmUpTime = 5 * time.Second // Long enough to block
	defer func() { WarmUpTime = originalWarmUpTime }()

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

func TestClient_Timeout(t *testing.T) {
	originalWarmUpTime := WarmUpTime
	WarmUpTime = 5 * time.Second // Long enough to block
	defer func() { WarmUpTime = originalWarmUpTime }()

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

func TestClient_MissingFields(t *testing.T) {
	originalWarmUpTime := WarmUpTime
	originalStabilizationTime := StabilizationTime
	WarmUpTime = 100 * time.Millisecond
	StabilizationTime = 100 * time.Millisecond
	defer func() {
		WarmUpTime = originalWarmUpTime
		StabilizationTime = originalStabilizationTime
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
	// Missing MakeModel
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

func TestClient_Unreachable(t *testing.T) {
	originalWarmUpTime := WarmUpTime
	originalStabilizationTime := StabilizationTime
	WarmUpTime = 100 * time.Millisecond
	StabilizationTime = 100 * time.Millisecond
	defer func() {
		WarmUpTime = originalWarmUpTime
		StabilizationTime = originalStabilizationTime
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
