// MFP - Multi-Function Printers and scanners toolkit
// Abstract definition for printer and scanner interfaces
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Print media size

package abstract

// MediaSize represents the physical dimensions of print media.
// Width and Height are expressed as [Dimension] values (units of 1/100 mm).
// A zero value for both fields means the media size is unset.
type MediaSize struct {
	Width  Dimension // Media width
	Height Dimension // Media height
}
