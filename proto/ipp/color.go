// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Color

package ipp

import "image/color"

// Color represents a traditional 32-bit alpha-premultiplied color
// with 8 bits per channel.
//
// It has the following layout:
//
//	0xRRGGBBAA
//	   | | | |
//	   | | | `- Alpha
//	   | | `--- Blue
//	   | `----- Green
//	   `------- Red
//
// For convenience, it implements a color.Color interface.
type Color uint32

// RGBA splits Color into separate channels, with 16 bits per channel
//
// It implements color.Color interface.
func (c Color) RGBA() (r, g, b, a uint32) {
	rgba := color.RGBA{
		R: uint8(c >> 24),
		G: uint8(c >> 16),
		B: uint8(c >> 8),
		A: uint8(c),
	}

	return rgba.RGBA()
}
