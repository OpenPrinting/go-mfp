// MFP - Miulti-Function Printers and scanners toolkit
// Abstract definition for printer and scanner interfaces
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Scan region

package abstract

import "fmt"

// Region defines a scanning region.
// See [Dimension] description for definition of the coordinate system.
type Region struct {
	XOffset Dimension // Horizontal offset, 0-based
	YOffset Dimension // Vertical offset, 0-based
	Width   Dimension // Region width
	Height  Dimension // Region height
}

// String returns string representation of [Region], for logging.
func (reg Region) String() string {
	return fmt.Sprintf("%dx%d+%d+%d",
		reg.Width, reg.Height, reg.XOffset, reg.YOffset)
}

// IsZero reports if Region has a zero value.
func (reg Region) IsZero() bool {
	return reg == Region{}
}

// Valid reports overall validity of the region.
func (reg Region) Valid() bool {
	if reg.IsZero() {
		return true
	}

	return reg.XOffset >= 0 && reg.YOffset >= 0 &&
		reg.Width > 0 && reg.Height > 0
}

// FitsCapabilities reports if scan Region fits scanner
// input capabilities.
func (reg Region) FitsCapabilities(caps *InputCapabilities) bool {
	return caps.MinWidth <= reg.Width && reg.Width <= caps.MaxWidth &&
		caps.MinHeight <= reg.Height && reg.Height <= caps.MaxHeight &&
		(caps.MaxXOffset == 0 || reg.XOffset <= caps.MaxXOffset) &&
		(caps.MaxYOffset == 0 || reg.YOffset <= caps.MaxYOffset) &&
		reg.XOffset+reg.Width <= caps.MaxWidth &&
		reg.YOffset+reg.Height <= caps.MaxHeight
}
