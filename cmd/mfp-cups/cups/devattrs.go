// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device information pretty-printer

package cups

import (
	"fmt"
	"io"

	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// devAttrsFormat pretty-prints [ipp.DeviceAttributes]
func devAttrsFormat(w io.Writer, dev *ipp.DeviceAttributes) {
	fmt.Fprintf(w, "Device information:\n")
	fmt.Fprintf(w, "  Class:          %s\n", dev.DeviceClass)
	fmt.Fprintf(w, "  Info:           %s\n", dev.DeviceInfo)
	fmt.Fprintf(w, "  Make and Model: %s\n", dev.DeviceMakeAndModel)
	fmt.Fprintf(w, "  Device URI:     %s\n", dev.DeviceURI)
	fmt.Fprintf(w, "  IEEE-1284 ID:   %s\n", dev.DeviceID)
	fmt.Fprintf(w, "  Location:       %s\n", dev.DeviceLocation)
}
