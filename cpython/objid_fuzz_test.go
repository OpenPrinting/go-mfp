// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Object identifiers fuzz test

package cpython

import "testing"

// FuzzObjmapLifecycle fuzzes a sequence of create/invalidate operations
// against the real objid map (via Python.NewObject / Object.Invalidate),
// checking that the live object count always matches what this test
// independently tracks. This exercises objmap.put/get/del/count under
// an adversarial mix of interleaved operations.
func FuzzObjmapLifecycle(f *testing.F) {
	f.Add([]byte{0, 1, 2, 1, 0})
	f.Add([]byte{2, 2, 2})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, ops []byte) {
		if len(ops) > 4096 {
			t.Skip("too long")
		}

		py, err := NewPython()
		if err != nil {
			t.Skip(err)
		}
		defer py.Close()

		base := py.countObjID()
		var live []*Object

		for i, op := range ops {
			switch op % 3 {
			case 0, 1:
				// Create a new object.
				obj := py.NewObject(int(op))
				if err := obj.Err(); err != nil {
					t.Fatalf("op %d: NewObject: %s", i, err)
				}
				live = append(live, obj)

			case 2:
				// Invalidate one, if any are live.
				if len(live) > 0 {
					idx := int(op) % len(live)
					live[idx].Invalidate()
					live = append(live[:idx], live[idx+1:]...)
				}
			}

			if got, want := py.countObjID(), base+len(live); got != want {
				t.Fatalf("op %d: countObjID = %d, want %d", i, got, want)
			}
		}
	})
}
