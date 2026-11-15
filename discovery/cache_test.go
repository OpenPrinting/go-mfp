// MFP - Multi-Function Printers and scanners toolkit
// Device discovery
//
// Copyright (C) 2026 and up by Anand Biju
// See LICENSE for license terms and conditions
//
// Discovery cache and output tests

package discovery

import (
	"context"
	"testing"
	"time"
)

// TestOutputCachedTTL checks that output.Cached returns the cached
// devices while the TTL is in the future, and nil once it has passed.
func TestOutputCachedTTL(t *testing.T) {
	out := output{ctx: context.Background()}

	out.Generate(time.Now().Add(time.Hour), nil)
	if out.Cached() == nil {
		t.Errorf("Cached() returned nil while TTL is still valid")
	}

	out.Generate(time.Now().Add(-time.Second), nil)
	if out.Cached() != nil {
		t.Errorf("Cached() returned devices after TTL expired")
	}

	out.Generate(time.Now().Add(time.Hour), nil)
	out.Invalidate()
	if out.Cached() != nil {
		t.Errorf("Cached() returned devices after Invalidate")
	}
}

// TestCacheExportUsesCachedOutput checks that cache.Export keeps
// its output cached between calls, until the cache is modified.
func TestCacheExportUsesCachedOutput(t *testing.T) {
	c := newCache(context.Background(), 0, 0)

	c.Export()
	if c.out.Cached() == nil {
		t.Fatalf("Export() result is not cached")
	}

	if !c.ReadyAt(ModeWaitIncomplete).IsZero() {
		t.Errorf("ReadyAt(ModeWaitIncomplete) waits, " +
			"though cached output is valid")
	}

	c.AddUnit(&EventAddUnit{ID: UnitID{DNSSDName: "test"}})
	if c.out.Cached() != nil {
		t.Errorf("cached output is not invalidated by AddUnit")
	}
}
