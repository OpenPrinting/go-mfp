// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Tests for object identifier mapping

package cpython

import (
	"testing"
)

// TestObjmapNewEmpty verifies that newObjmap creates an empty map.
func TestObjmapNewEmpty(t *testing.T) {
	omap := newObjmap()
	if omap == nil {
		t.Fatalf("newObjmap() returned nil")
	}
	if omap.count() != 0 {
		t.Fatalf("newObjmap: count = %d, want 0", omap.count())
	}
}

// TestObjmapPutGet verifies that put and get work correctly.
func TestObjmapPutGet(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	// Verify count increases after adding objects
	before := py.countObjID()

	obj1 := py.Eval("1")
	if err := obj1.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	obj2 := py.Eval("2")
	if err := obj2.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	after := py.countObjID()
	if after != before+2 {
		t.Fatalf("count after put: got %d, want %d", after, before+2)
	}
}

// TestObjmapDel verifies that del removes objects from the map.
func TestObjmapDel(t *testing.T) {
	py, err := NewPython()
	if err != nil {
		t.Fatalf("NewPython: %v", err)
	}
	defer py.Close()

	before := py.countObjID()

	obj := py.Eval("42")
	if err := obj.Err(); err != nil {
		t.Fatalf("Eval: %v", err)
	}

	if py.countObjID() != before+1 {
		t.Fatalf("count after put: got %d, want %d",
			py.countObjID(), before+1)
	}

	obj.Invalidate()

	// After invalidation the object should be removed
	if py.countObjID() != before {
		t.Fatalf("count after del: got %d, want %d",
			py.countObjID(), before)
	}
}

// TestObjidInc verifies that objid.inc generates unique incrementing IDs.
func TestObjidInc(t *testing.T) {
	var oid objid

	prev := oid.inc()
	for i := 0; i < 100; i++ {
		next := oid.inc()
		if next <= prev {
			t.Fatalf("inc: next %d <= prev %d", next, prev)
		}
		prev = next
	}
}
