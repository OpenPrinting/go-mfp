// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device attribites, as returned by CUPS-Get-Devices

package ipp

import "github.com/OpenPrinting/go-mfp/util/optional"

// PPDAttributes represents PPD file attributes, as returned by
// the CUPS-Get-PPDs request
type PPDAttributes struct {
	ObjectRawAttrs
	CUPSPPDAttributesGroup

	PPDName            optional.Val[string] `ipp:"ppd-name"`
	PPDDeviceID        optional.Val[string] `ipp:"ppd-device-id"`
	PPDMakeAndModel    optional.Val[string] `ipp:"ppd-make-and-model"`
	PPDMake            optional.Val[string] `ipp:"ppd-make"`
	PPDModelNumber     optional.Val[int]    `ipp:"ppd-model-number"`
	PPDNaturalLanguage optional.Val[string] `ipp:"ppd-natural-language"`
	PPDProduct         optional.Val[string] `ipp:"ppd-product"`
	PPDPsVersion       optional.Val[string] `ipp:"ppd-psversion"`
	PPDType            optional.Val[string] `ipp:"ppd-type,keyword"`
}
