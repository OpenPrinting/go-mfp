// MFP - Miulti-Function Printers and scanners toolkit
// USB host API
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Top-level API

package usbhost

import "sort"

var listDevices = libusbListDevices

// ListDevices returns list of all connected USB devices.
//
// If withIEEE1284Id is true, this function also loads
// IEEE-1284 device ID, where appropriate.
//
// Note setting the withIEEE1284id flag may have a side effect
// of changing the active USB device configuration.
func ListDevices(withIEEE1284id bool) ([]DeviceInfo, error) {
	// Obtain list of devices
	infos, err := listDevices(withIEEE1284id)
	if err != nil {
		return nil, err
	}

	// Sort by location
	sort.Slice(infos, func(i, j int) bool {
		loc1 := infos[i].Loc
		loc2 := infos[j].Loc

		switch {
		case loc1.Bus < loc2.Bus:
			return true
		case loc1.Bus > loc2.Bus:
			return false
		}
		return loc1.Dev < loc2.Dev
	})

	return infos, nil
}
