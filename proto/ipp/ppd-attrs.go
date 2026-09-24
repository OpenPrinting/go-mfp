// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Device attribites, as returned by CUPS-Get-Devices

package ipp

// PPDAttributes represents PPD file attributes, as returned by
// the CUPS-Get-PPDs request
type PPDAttributes struct {
	ObjectRawAttrs
	CUPSPPDAttributesGroup
}
