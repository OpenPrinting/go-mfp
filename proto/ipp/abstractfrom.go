// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Conversions from abstract.Scanner to IPP data structures

package ipp

import (
	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// fromAbstractScannerDescription translates [abstract.ScannerCapabilities]
// into an IPP [ScannerDescription].
func fromAbstractScannerDescription(
	caps *abstract.ScannerCapabilities) ScannerDescription {

	sd := ScannerDescription{
		InputSourceSupported:     abstractInputSources(caps),
		InputColorModeSupported:  abstractInputColorModes(caps),
		InputResolutionSupported: abstractInputResolutions(caps),
		InputSidesSupported:      abstractInputSides(caps),
		InputAttributesSupported: abstractInputAttributesSupported(caps),
	}

	if req := caps.DefaultRequest(); req != nil {
		sd.InputAttributesDefault = optional.New(fromAbstractInputAttributes(req))
	}

	return sd
}

// abstractInputSources returns InputSourceSupported values derived
// from the abstract scanner capabilities.
func abstractInputSources(
	caps *abstract.ScannerCapabilities) []KwInputSource {

	var out []KwInputSource

	if caps.Platen != nil {
		out = append(out, KwInputSourcePlaten)
	}

	if caps.ADFSimplex != nil || caps.ADFDuplex != nil {
		out = append(out, KwInputSourceADF)
	}

	return out
}

// abstractInputColorModes collects unique IPP color modes from all
// inputs of the abstract scanner capabilities.
func abstractInputColorModes(
	caps *abstract.ScannerCapabilities) []KwInputColorMode {

	seen := generic.NewSet[KwInputColorMode]()
	var out []KwInputColorMode

	add := func(m KwInputColorMode) {
		if seen.TestAndAdd(m) {
			out = append(out, m)
		}
	}

	for _, inp := range abstractAllInputs(caps) {
		for _, prof := range inp.Profiles {
			if prof.ColorModes.Contains(abstract.ColorModeBinary) {
				add(KwInputColorModeBiLevel)
			}
			if prof.ColorModes.Contains(abstract.ColorModeMono) {
				if prof.Depths.IsEmpty() {
					add(KwInputColorModeMonochrome)
				} else {
					if prof.Depths.Contains(abstract.ColorDepth8) {
						add(KwInputColorModeMonochrome8)
					}
					if prof.Depths.Contains(abstract.ColorDepth16) {
						add(KwInputColorModeMonochrome16)
					}
				}
			}
			if prof.ColorModes.Contains(abstract.ColorModeColor) {
				if prof.Depths.IsEmpty() {
					add(KwInputColorModeColor)
				} else {
					if prof.Depths.Contains(abstract.ColorDepth8) {
						add(KwInputColorModeColor8)
					}
					if prof.Depths.Contains(abstract.ColorDepth16) {
						add(KwInputColorModeRGB16)
					}
				}
			}
		}
	}

	return out
}

// abstractInputResolutions collects unique IPP resolutions from all
// inputs of the abstract scanner capabilities.
func abstractInputResolutions(
	caps *abstract.ScannerCapabilities) []goipp.Resolution {

	seen := generic.NewSet[goipp.Resolution]()
	var out []goipp.Resolution

	for _, inp := range abstractAllInputs(caps) {
		for _, res := range inp.Resolutions() {
			r := goipp.Resolution{
				Xres:  res.XResolution,
				Yres:  res.YResolution,
				Units: goipp.UnitsDpi,
			}
			if seen.TestAndAdd(r) {
				out = append(out, r)
			}
		}
	}

	return out
}

// abstractInputSides returns InputSidesSupported values derived
// from the abstract scanner capabilities.
func abstractInputSides(caps *abstract.ScannerCapabilities) []KwSides {
	if abstractAllInputs(caps) == nil {
		return nil
	}

	out := []KwSides{KwSidesOneSided}

	if caps.ADFDuplex != nil {
		out = append(out, KwSidesTwoSidedLongEdge)
	}

	return out
}

// abstractInputAttributesSupported returns the list of input-attributes
// member names supported based on the scanner capabilities.
func abstractInputAttributesSupported(
	caps *abstract.ScannerCapabilities) []string {

	attrs := []string{
		"input-color-mode",
		"input-resolution",
		"input-scan-regions",
		"input-source",
	}

	if caps.ADFDuplex != nil {
		attrs = append(attrs, "input-sides")
	}

	if !caps.BrightnessRange.IsZero() {
		attrs = append(attrs, "input-brightness")
	}

	if !caps.ContrastRange.IsZero() {
		attrs = append(attrs, "input-contrast")
	}

	if !caps.SharpenRange.IsZero() {
		attrs = append(attrs, "input-sharpness")
	}

	return attrs
}

// fromAbstractInputAttributes converts an [abstract.ScannerRequest] into
// an IPP [InputAttributes].
func fromAbstractInputAttributes(req *abstract.ScannerRequest) InputAttributes {
	attrs := InputAttributes{}

	// Input source
	switch req.Input {
	case abstract.InputPlaten:
		attrs.InputSource = optional.New(KwInputSourcePlaten)
	case abstract.InputADF:
		attrs.InputSource = optional.New(KwInputSourceADF)
	}

	// Input sides
	switch {
	case req.Input == abstract.InputADF && req.ADFMode == abstract.ADFModeDuplex:
		attrs.InputSides = optional.New(KwSidesTwoSidedLongEdge)
	case req.Input != abstract.InputUnset:
		attrs.InputSides = optional.New(KwSidesOneSided)
	}

	// Color mode
	switch req.ColorMode {
	case abstract.ColorModeBinary:
		attrs.InputColorMode = optional.New(KwInputColorModeBiLevel)
	case abstract.ColorModeMono:
		switch req.ColorDepth {
		case abstract.ColorDepth8:
			attrs.InputColorMode = optional.New(KwInputColorModeMonochrome8)
		case abstract.ColorDepth16:
			attrs.InputColorMode = optional.New(KwInputColorModeMonochrome16)
		default:
			attrs.InputColorMode = optional.New(KwInputColorModeMonochrome)
		}
	case abstract.ColorModeColor:
		switch req.ColorDepth {
		case abstract.ColorDepth8:
			attrs.InputColorMode = optional.New(KwInputColorModeColor8)
		case abstract.ColorDepth16:
			attrs.InputColorMode = optional.New(KwInputColorModeRGB16)
		default:
			attrs.InputColorMode = optional.New(KwInputColorModeColor)
		}
	}

	// Resolution
	if !req.Resolution.IsZero() {
		attrs.InputResolution = optional.New(goipp.Resolution{
			Xres:  req.Resolution.XResolution,
			Yres:  req.Resolution.YResolution,
			Units: goipp.UnitsDpi,
		})
	}

	// Scan region
	if !req.Region.IsZero() {
		attrs.InputScanRegions = []InputScanRegion{
			{
				XOrigin:    optional.New(int(req.Region.XOffset)),
				YOrigin:    optional.New(int(req.Region.YOffset)),
				XDimension: optional.New(int(req.Region.Width)),
				YDimension: optional.New(int(req.Region.Height)),
			},
		}
	}

	// Brightness, contrast, sharpness
	if req.Brightness != nil {
		attrs.InputBrightness = optional.New(optional.Get(req.Brightness))
	}

	if req.Contrast != nil {
		attrs.InputContrast = optional.New(optional.Get(req.Contrast))
	}

	if req.Sharpen != nil {
		attrs.InputSharpness = optional.New(optional.Get(req.Sharpen))
	}

	return attrs
}

// abstractAllInputs returns all non-nil InputCapabilities from caps.
func abstractAllInputs(
	caps *abstract.ScannerCapabilities) []*abstract.InputCapabilities {

	var inputs []*abstract.InputCapabilities

	if caps.Platen != nil {
		inputs = append(inputs, caps.Platen)
	}

	if caps.ADFSimplex != nil {
		inputs = append(inputs, caps.ADFSimplex)
	}

	if caps.ADFDuplex != nil {
		inputs = append(inputs, caps.ADFDuplex)
	}

	return inputs
}
