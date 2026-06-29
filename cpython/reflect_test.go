// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Tests for reflect helper functions

package cpython

import (
	"reflect"
	"testing"
)

// TestReflectSortInts verifies that reflectSort sorts int slices correctly.
func TestReflectSortInts(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf(int(3)),
		reflect.ValueOf(int(1)),
		reflect.ValueOf(int(2)),
	}

	if !reflectSort(slice) {
		t.Fatalf("reflectSort returned false for int slice")
	}

	for i := 0; i < len(slice)-1; i++ {
		if slice[i].Int() > slice[i+1].Int() {
			t.Fatalf("slice not sorted at index %d: %d > %d",
				i, slice[i].Int(), slice[i+1].Int())
		}
	}
}

// TestReflectSortUints verifies that reflectSort sorts uint slices correctly.
func TestReflectSortUints(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf(uint(30)),
		reflect.ValueOf(uint(10)),
		reflect.ValueOf(uint(20)),
	}

	if !reflectSort(slice) {
		t.Fatalf("reflectSort returned false for uint slice")
	}

	for i := 0; i < len(slice)-1; i++ {
		if slice[i].Uint() > slice[i+1].Uint() {
			t.Fatalf("slice not sorted at index %d: %d > %d",
				i, slice[i].Uint(), slice[i+1].Uint())
		}
	}
}

// TestReflectSortFloats verifies that reflectSort sorts float slices correctly.
func TestReflectSortFloats(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf(float64(3.3)),
		reflect.ValueOf(float64(1.1)),
		reflect.ValueOf(float64(2.2)),
	}

	if !reflectSort(slice) {
		t.Fatalf("reflectSort returned false for float slice")
	}

	for i := 0; i < len(slice)-1; i++ {
		if slice[i].Float() > slice[i+1].Float() {
			t.Fatalf("slice not sorted at index %d: %f > %f",
				i, slice[i].Float(), slice[i+1].Float())
		}
	}
}

// TestReflectSortStrings verifies that reflectSort sorts string slices correctly.
func TestReflectSortStrings(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf("cherry"),
		reflect.ValueOf("apple"),
		reflect.ValueOf("banana"),
	}

	if !reflectSort(slice) {
		t.Fatalf("reflectSort returned false for string slice")
	}

	for i := 0; i < len(slice)-1; i++ {
		if slice[i].String() > slice[i+1].String() {
			t.Fatalf("slice not sorted at index %d: %q > %q",
				i, slice[i].String(), slice[i+1].String())
		}
	}
}

// TestReflectSortBools verifies that reflectSort sorts bool slices correctly.
func TestReflectSortBools(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf(true),
		reflect.ValueOf(false),
		reflect.ValueOf(true),
	}

	if !reflectSort(slice) {
		t.Fatalf("reflectSort returned false for bool slice")
	}

	// false should come before true
	if slice[0].Bool() {
		t.Fatalf("bool sort: expected false first, got true")
	}
}

// TestReflectSortEmpty verifies that reflectSort returns true for empty slice.
func TestReflectSortEmpty(t *testing.T) {
	if !reflectSort([]reflect.Value{}) {
		t.Fatalf("reflectSort returned false for empty slice")
	}
}

// TestReflectSortMixedTypes verifies that reflectSort returns false
// for slices with mixed types.
func TestReflectSortMixedTypes(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf("hello"),
	}

	if reflectSort(slice) {
		t.Fatalf("reflectSort returned true for mixed-type slice")
	}
}

// TestReflectSortUnsupportedType verifies that reflectSort returns false
// for unsupported types.
func TestReflectSortUnsupportedType(t *testing.T) {
	slice := []reflect.Value{
		reflect.ValueOf([]int{1, 2}),
		reflect.ValueOf([]int{3, 4}),
	}

	if reflectSort(slice) {
		t.Fatalf("reflectSort returned true for unsupported type")
	}
}
