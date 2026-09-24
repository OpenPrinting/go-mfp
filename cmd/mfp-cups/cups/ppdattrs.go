// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// PPD information pretty-printer

package cups

import (
	"fmt"
	"io"

	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// prnAttrsFormat pretty-prints [ipp.PrinterAttributes]
func ppdAttrsFormat(w io.Writer, prn *ipp.PPDAttributes) {
	fmt.Fprintf(w, "Name:           %s\n", optional.Get(prn.PPDName))
	fmt.Fprintf(w, "Language:       %q\n", optional.Get(prn.PPDNaturalLanguage))
	fmt.Fprintf(w, "Make:           %q\n", optional.Get(prn.PPDMake))
	fmt.Fprintf(w, "Make and Model: %q\n", optional.Get(prn.PPDMakeAndModel))
	fmt.Fprintf(w, "Model Number:   %d\n", optional.Get(prn.PPDModelNumber))
	fmt.Fprintf(w, "IEEE-1284 ID:   %q\n", optional.Get(prn.PPDDeviceID))
	fmt.Fprintf(w, "Product:        %q\n", optional.Get(prn.PPDProduct))
	fmt.Fprintf(w, "PS Version:     %q\n", optional.Get(prn.PPDPsVersion))
	fmt.Fprintf(w, "PPD Type:       %q\n", optional.Get(prn.PPDType))
}
