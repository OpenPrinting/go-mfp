// MFP - Miulti-Function Printers and scanners toolkit
// Abstract definition for printer and scanner interfaces
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Settings profile

package abstract

import "github.com/OpenPrinting/go-mfp/util/generic"

// SettingsProfile defines a valid combination of scan parameters setting.
//
// Each [InputCapabilities] may include multiple instances of
// SettingsProfile.
type SettingsProfile struct {
	ColorModes       generic.Bitset[ColorMode]       // Allowed color modes
	Depths           generic.Bitset[ColorDepth]      // Allowed depths
	BinaryRenderings generic.Bitset[BinaryRendering] // For 1-bit B&W
	CCDChannels      generic.Bitset[CCDChannel]      // Allowed CCD channel
	Resolutions      []Resolution                    // Allowed resolutions
	ResolutionRange  ResolutionRange                 // Zero if unset
}

// AllowsColorMode reports if the [SettingsProfile] compatible
// with the given combination of the [ColorMode], [Depth],
// and [BinaryRendering] combination.
func (prof SettingsProfile) AllowsColorMode(
	cm ColorMode, d ColorDepth, r BinaryRendering) bool {
	switch {
	case cm == ColorModeUnset:
		return true
	case !prof.ColorModes.Contains(cm):
		return false
	case cm == ColorModeBinary:
		return r == BinaryRenderingUnset ||
			prof.BinaryRenderings.Contains(r)
	}

	return d == ColorDepthUnset || prof.Depths.Contains(d)
}

// AllowsCCDChannel reports if the [SettingsProfile] compatible
// with the given [CCDChannel].
func (prof SettingsProfile) AllowsCCDChannel(ccd CCDChannel) bool {
	if ccd == CCDChannelUnset {
		return true
	}

	return prof.CCDChannels.Contains(ccd)
}

// AllowsResolution reports if the [SettingsProfile] compatible
// with the given [Resolution].
func (prof SettingsProfile) AllowsResolution(res Resolution) bool {
	if res.IsZero() {
		return true
	}

	for _, supported := range prof.Resolutions {
		if res == supported {
			return true
		}
	}

	if !prof.ResolutionRange.IsZero() {
		rr := prof.ResolutionRange
		if res.XResolution >= rr.XMin && res.XResolution <= rr.XMax &&
			res.YResolution >= rr.YMin && res.YResolution <= rr.YMax {
			if (rr.XStep == 0 || (res.XResolution-rr.XMin)%rr.XStep == 0) &&
				(rr.YStep == 0 || (res.YResolution-rr.YMin)%rr.YStep == 0) {
				return true
			}
		}
	}

	return false
}

func defaultResolution(profiles []SettingsProfile) Resolution {
	// Choose profiles that allows the desired resolution
	desiredRes := Resolution{300, 300}

	// Lookup profile that directly supports the desired resolution
	for _, prof := range profiles {
		if prof.AllowsResolution(desiredRes) {
			return desiredRes
		}
	}

	// If nothing found, choose the best resolution from available ones.
	// First off all, gather all available resolutions
	resolutions := make(map[Resolution]struct{})
	square := make(map[Resolution]struct{})

	for _, prof := range profiles {
		for _, res := range prof.Resolutions {
			resolutions[res] = struct{}{}
			if res.XResolution == res.YResolution {
				square[res] = struct{}{}
			}
		}
	}

	// Prefer square resolutions if available with fallback
	// to all possible resolutions
	if len(square) > 0 {
		resolutions = square
	}

	// Try to use only resolutions that exceeds the desired, if any
	use := make([]Resolution, 0, len(resolutions))
	for res := range resolutions {
		if res.XResolution >= desiredRes.XResolution &&
			res.YResolution >= desiredRes.YResolution {
			use = append(use, res)
		}
	}

	// If nothing found, use all available resolutions
	if len(use) == 0 {
		for res := range resolutions {
			use = append(use, res)
		}
	}

	// Locate nearest resolution from available ones.
	//
	// We range resolutions by the euclidean distance,
	// and prefer more square resolution, if euclidean
	// distances are equal.
	if len(use) == 0 {
		// Should not happen
		return Resolution{}
	}

	bestRes := use[0]
	bestDest := bestRes.euclideanDistance(desiredRes)

	for _, res := range use {
		dest := res.euclideanDistance(desiredRes)
		if dest < bestDest ||
			dest == bestDest && bestRes.lessSquare(res) {
			bestRes, bestDest = res, dest
		}
	}

	return bestRes
}

// summarizeColorModes returns summary ColorMode supported
// by given profiles.
func summarizeColorModes(profiles []SettingsProfile) generic.Bitset[ColorMode] {
	var sum generic.Bitset[ColorMode]
	for i := range profiles {
		sum = sum.Union(profiles[i].ColorModes)
	}
	return sum
}

// summarizeColorDepths returns summary ColorDepth supported
// by given profiles.
func summarizeColorDepths(profiles []SettingsProfile) generic.Bitset[ColorDepth] {
	var sum generic.Bitset[ColorDepth]
	for i := range profiles {
		sum = sum.Union(profiles[i].Depths)
	}
	return sum
}

// summarizeBinaryRenderings returns summary BinaryRendering supported
// by given profiles.
func summarizeBinaryRenderings(
	profiles []SettingsProfile) generic.Bitset[BinaryRendering] {

	var sum generic.Bitset[BinaryRendering]
	for i := range profiles {
		sum = sum.Union(profiles[i].BinaryRenderings)
	}
	return sum
}

// profilesByColorMode filters profiles by ColorMode.
// It modifies profiles in place and returns updated version.
func profilesByColorMode(profiles []SettingsProfile,
	mode ColorMode) []SettingsProfile {

	o := 0
	for i := range profiles {
		if profiles[i].ColorModes.Contains(mode) {
			profiles[o] = profiles[i]
			o++
		}
	}

	return profiles[:o]
}

// profilesByColorDepth filters profiles by ColorDepth.
// It modifies profiles in place and returns updated version.
func profilesByColorDepth(profiles []SettingsProfile,
	depth ColorDepth) []SettingsProfile {

	o := 0
	for i := range profiles {
		if profiles[i].Depths.Contains(depth) {
			profiles[o] = profiles[i]
			o++
		}
	}

	return profiles[:o]
}

// profilesByBinaryRendering filters profiles by BinaryRendering.
// It modifies profiles in place and returns updated version.
func profilesByBinaryRendering(profiles []SettingsProfile,
	rend BinaryRendering) []SettingsProfile {

	o := 0
	for i := range profiles {
		if profiles[i].BinaryRenderings.Contains(rend) {
			profiles[o] = profiles[i]
			o++
		}
	}

	return profiles[:o]
}

// profilesByCCDChannel filters profiles by CCDChannel.
// It modifies profiles in place and returns updated version.
func profilesByCCDChannel(profiles []SettingsProfile,
	ccd CCDChannel) []SettingsProfile {

	o := 0
	for i := range profiles {
		if profiles[i].CCDChannels.Contains(ccd) {
			profiles[o] = profiles[i]
			o++
		}
	}

	return profiles[:o]
}
