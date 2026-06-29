// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Tests for floating point limits

package cpython

import (
	"math"
	"testing"
)

// TestMaxInt64Float verifies that maxInt64Float is a valid boundary
// for float64 to int64 conversion.
func TestMaxInt64Float(t *testing.T) {
	if maxInt64Float <= 0 {
		t.Fatalf("maxInt64Float = %v, want positive value", maxInt64Float)
	}

	if maxInt64Float > math.MaxInt64 {
		t.Fatalf("maxInt64Float = %v exceeds math.MaxInt64", maxInt64Float)
	}

	// Verify conversion is lossless at the boundary
	if int64(maxInt64Float) != int64(float64(int64(maxInt64Float))) {
		t.Fatalf("maxInt64Float boundary conversion is not lossless")
	}
}

// TestMinInt64Float verifies that minInt64Float is a valid boundary
// for float64 to int64 conversion.
func TestMinInt64Float(t *testing.T) {
	if minInt64Float >= 0 {
		t.Fatalf("minInt64Float = %v, want negative value", minInt64Float)
	}

	if minInt64Float < math.MinInt64 {
		t.Fatalf("minInt64Float = %v is below math.MinInt64", minInt64Float)
	}

	// Verify conversion is lossless at the boundary
	if int64(minInt64Float) != int64(float64(int64(minInt64Float))) {
		t.Fatalf("minInt64Float boundary conversion is not lossless")
	}
}

// TestMaxUint64Float verifies that maxUint64Float is a valid boundary
// for float64 to uint64 conversion.
func TestMaxUint64Float(t *testing.T) {
	if maxUint64Float <= 0 {
		t.Fatalf("maxUint64Float = %v, want positive value", maxUint64Float)
	}

	if maxUint64Float > math.MaxUint64 {
		t.Fatalf("maxUint64Float = %v exceeds math.MaxUint64", maxUint64Float)
	}

	// Verify conversion is lossless at the boundary
	if uint64(maxUint64Float) != uint64(float64(uint64(maxUint64Float))) {
		t.Fatalf("maxUint64Float boundary conversion is not lossless")
	}
}

// TestFloatBoundariesOrder verifies the relative ordering of boundaries.
func TestFloatBoundariesOrder(t *testing.T) {
	if minInt64Float >= maxInt64Float {
		t.Fatalf("minInt64Float (%v) >= maxInt64Float (%v)",
			minInt64Float, maxInt64Float)
	}

	if maxInt64Float > maxUint64Float {
		t.Fatalf("maxInt64Float (%v) > maxUint64Float (%v)",
			maxInt64Float, maxUint64Float)
	}
}
