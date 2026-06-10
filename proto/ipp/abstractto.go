// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Conversions from IPP data structures to abstract types

package ipp

import "github.com/OpenPrinting/go-mfp/abstract"

// sidesToAbstract maps a KwSides IPP keyword to abstract.Sides.
func sidesToAbstract(kw KwSides) abstract.Sides {
	switch kw {
	case KwSidesOneSided:
		return abstract.SidesOneSided
	case KwSidesTwoSidedLongEdge:
		return abstract.SidesTwoSidedLongEdge
	case KwSidesTwoSidedShortEdge:
		return abstract.SidesTwoSidedShortEdge
	}
	return abstract.SidesUnset
}

// colorModeToAbstract maps an IPP print-color-mode string to abstract.ColorMode.
func colorModeToAbstract(s string) abstract.ColorMode {
	switch s {
	case "color":
		return abstract.ColorModeColor
	case "monochrome", "auto-monochrome", "process-monochrome",
		"highlight-monochrome":
		return abstract.ColorModeMono
	case "bi-level":
		return abstract.ColorModeBinary
	}
	return abstract.ColorModeUnset
}

// mediaSizeToAbstract maps a KwMedia IPP keyword to abstract.MediaSize.
// KwMedia.Size() returns dimensions in 1/100 mm, matching Dimension units.
func mediaSizeToAbstract(kw KwMedia) abstract.MediaSize {
	wid, hei := kw.Size()
	if wid <= 0 || hei <= 0 {
		return abstract.MediaSize{}
	}
	return abstract.MediaSize{
		Width:  abstract.Dimension(wid),
		Height: abstract.Dimension(hei),
	}
}
