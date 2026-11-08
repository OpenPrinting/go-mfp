// MFP - Miulti-Function Printers and scanners toolkit
// The "mfp-cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device information pretty-printer

package command

import (
	"fmt"
	"io"

	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// devAttrsFormat pretty-prints [ipp.DeviceAttributes]
func devAttrsFormat(w io.Writer, dev *ipp.DeviceAttributes) {
	fmt.Fprintf(w, "Device information:\n")
	if dev.DeviceClass != nil {
		fmt.Fprintf(w, "  Class:          %s\n", *dev.DeviceClass)
	}
	if dev.DeviceInfo != nil {
		fmt.Fprintf(w, "  Info:           %s\n", *dev.DeviceInfo)
	}
	if dev.DeviceMakeAndModel != nil {
		fmt.Fprintf(w, "  Make and Model: %s\n", *dev.DeviceMakeAndModel)
	}
	if dev.DeviceURI != nil {
		fmt.Fprintf(w, "  Device URI:     %s\n", *dev.DeviceURI)
	}
	if dev.DeviceID != nil {
		fmt.Fprintf(w, "  IEEE-1284 ID:   %s\n", *dev.DeviceID)
	}
	if dev.DeviceLocation != nil {
		fmt.Fprintf(w, "  Location:       %s\n", *dev.DeviceLocation)
	}
}
