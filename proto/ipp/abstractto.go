// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Conversions from IPP data structures to abstract types

package ipp

import (
	"slices"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/uuid"
	"github.com/OpenPrinting/goipp"
)

// ToAbstractScannerCapabilities converts [PrinterAttributes] of the
// IPP Scan Service into the [abstract.ScannerCapabilities].
//
// IPP doesn't define per-input scanner capabilities, so all inputs
// share the same geometry and settings profile.
func (pa *PrinterAttributes) ToAbstractScannerCapabilities() *abstract.ScannerCapabilities {
	abscaps := &abstract.ScannerCapabilities{
		MakeAndModel:    optional.Get(pa.PrinterMakeAndModel),
		DocumentFormats: slices.Clone(pa.DocumentFormatSupported),
	}

	if pa.PrinterUUID != nil {
		if u, err := uuid.Parse(*pa.PrinterUUID); err == nil {
			abscaps.UUID = u
		}
	}

	sd := &pa.ScannerDescription

	duplex := slices.ContainsFunc(sd.InputSidesSupported,
		func(kw KwSides) bool {
			return kw == KwSidesTwoSidedLongEdge ||
				kw == KwSidesTwoSidedShortEdge
		})

	for _, src := range sd.InputSourceSupported {
		switch src {
		case KwInputSourcePlaten:
			abscaps.Platen = sd.toAbstractInputCapabilities()

		case KwInputSourceADF:
			abscaps.ADFSimplex = sd.toAbstractInputCapabilities()
			if duplex {
				abscaps.ADFDuplex = sd.toAbstractInputCapabilities()
			}
		}
	}

	return abscaps
}

// toAbstractInputCapabilities converts [ScannerDescription] into the
// [abstract.InputCapabilities].
func (sd *ScannerDescription) toAbstractInputCapabilities() *abstract.InputCapabilities {
	inpcaps := &abstract.InputCapabilities{
		Profiles: []abstract.SettingsProfile{
			sd.toAbstractSettingsProfile(),
		},
	}

	if sd.InputScanRegionsSupported != nil {
		regions := *sd.InputScanRegionsSupported

		if regions.XDimension != nil {
			inpcaps.MinWidth = abstract.Dimension(regions.XDimension.Lower)
			inpcaps.MaxWidth = abstract.Dimension(regions.XDimension.Upper)
		}

		if regions.YDimension != nil {
			inpcaps.MinHeight = abstract.Dimension(regions.YDimension.Lower)
			inpcaps.MaxHeight = abstract.Dimension(regions.YDimension.Upper)
		}

		if regions.XOrigin != nil {
			inpcaps.MaxXOffset = abstract.Dimension(regions.XOrigin.Upper)
		}

		if regions.YOrigin != nil {
			inpcaps.MaxYOffset = abstract.Dimension(regions.YOrigin.Upper)
		}
	}

	return inpcaps
}

// toAbstractSettingsProfile converts [ScannerDescription] into the
// [abstract.SettingsProfile].
func (sd *ScannerDescription) toAbstractSettingsProfile() abstract.SettingsProfile {
	prof := abstract.SettingsProfile{}

	for _, cm := range sd.InputColorModeSupported {
		mode, depth := inputColorModeToAbstract(cm)

		switch mode {
		case abstract.ColorModeUnset:
			// "auto" and unknown modes
			continue

		case abstract.ColorModeBinary:
			// IPP doesn't define binary rendering modes,
			// assume halftone.
			prof.BinaryRenderings.Add(abstract.BinaryRenderingHalftone)

		default:
			// "color" and "monochrome" don't define
			// color depth, assume 8 bits per channel.
			if depth == abstract.ColorDepthUnset {
				depth = abstract.ColorDepth8
			}
			prof.Depths.Add(depth)
		}

		prof.ColorModes.Add(mode)
	}

	for _, res := range sd.InputResolutionSupported {
		x, y := res.Xres, res.Yres
		if res.Units == goipp.UnitsDpcm {
			x = x * 254 / 100
			y = y * 254 / 100
		}

		absres := abstract.Resolution{XResolution: x, YResolution: y}
		if !slices.Contains(prof.Resolutions, absres) {
			prof.Resolutions = append(prof.Resolutions, absres)
		}
	}

	return prof
}

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

// inputColorModeToAbstract maps a KwInputColorMode IPP keyword to
// abstract.ColorMode and abstract.ColorDepth.
func inputColorModeToAbstract(cm KwInputColorMode) (
	abstract.ColorMode, abstract.ColorDepth) {
	switch cm {
	case KwInputColorModeBiLevel:
		return abstract.ColorModeBinary, abstract.ColorDepthUnset
	case KwInputColorModeMonochrome:
		return abstract.ColorModeMono, abstract.ColorDepthUnset
	case KwInputColorModeMonochrome4, KwInputColorModeMonochrome8:
		return abstract.ColorModeMono, abstract.ColorDepth8
	case KwInputColorModeMonochrome16:
		return abstract.ColorModeMono, abstract.ColorDepth16
	case KwInputColorModeColor:
		return abstract.ColorModeColor, abstract.ColorDepthUnset
	case KwInputColorModeColor8, KwInputColorModeRGBA8, KwInputColorModeCMYK8:
		return abstract.ColorModeColor, abstract.ColorDepth8
	case KwInputColorModeRGB16, KwInputColorModeRGBA16, KwInputColorModeCMYK16:
		return abstract.ColorModeColor, abstract.ColorDepth16
	}
	// KwInputColorModeAuto and unknown values: let caps choose.
	return abstract.ColorModeUnset, abstract.ColorDepthUnset
}
