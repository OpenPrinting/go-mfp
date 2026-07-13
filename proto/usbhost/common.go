// MFP - Miulti-Function Printers and scanners toolkit
// USB host API
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Common definitions

package usbhost

import "github.com/OpenPrinting/go-mfp/proto/usb"

// Location represents the device location as Bus and Dev numbers
type Location struct {
	Bus int
	Dev int
}

// DeviceInfo contains USB device location and descriptor
type DeviceInfo struct {
	Loc  Location             // Device location
	Desc usb.DeviceDescriptor // Device descriptor
}

// IsPrinter reports whether the device is a printer. More precisely, it
// returns true if the device contains interfaces of the Printer class. Note
// that this means an IPP-USB scanner without printing capabilities would
// still be considered a printer.
func (info *DeviceInfo) IsPrinter() bool {
	return info.Desc.Contains(7, 1, -1)
}
