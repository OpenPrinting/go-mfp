// MFP - Multi-Function Printers and scanners toolkit
// Abstract definition for printer and scanner interfaces
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Print sides (duplex) mode

package abstract

import "fmt"

// Sides controls how pages are imposed onto media sides during printing.
type Sides int

// Known sides values:
const (
	SidesUnset             Sides = iota // Not set
	SidesOneSided                       // Single-sided output
	SidesTwoSidedLongEdge               // Duplex, bound along the long edge
	SidesTwoSidedShortEdge              // Duplex, bound along the short edge
	sidesMax
)

// String returns the string representation of [Sides], for logging.
func (s Sides) String() string {
	switch s {
	case SidesUnset:
		return "Unset"
	case SidesOneSided:
		return "OneSided"
	case SidesTwoSidedLongEdge:
		return "TwoSidedLongEdge"
	case SidesTwoSidedShortEdge:
		return "TwoSidedShortEdge"
	}
	return fmt.Sprintf("Unknown (%d)", int(s))
}
