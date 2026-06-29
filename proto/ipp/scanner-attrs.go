// MFP - Multi-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Scanner Attributes

package ipp

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// InputScanRegion represents a single region within the "input-scan-regions"
// member of [InputAttributes].
//
// All dimensions are in hundredths of a millimeter (1/100 mm).
// See PWG5100.15.
type InputScanRegion struct {
	XDimension optional.Val[int] `ipp:"x-dimension"`
	XOrigin    optional.Val[int] `ipp:"x-origin"`
	YDimension optional.Val[int] `ipp:"y-dimension"`
	YOrigin    optional.Val[int] `ipp:"y-origin"`
}

// OutputAttributes represents the "output-attributes" collection.
//
// It is used in scan job operation requests to specify per-job
// image-processing settings for the output document(s), and as the value
// type of "output-attributes-default" in printer description attributes.
type OutputAttributes struct {
	NoiseRemoval                   optional.Val[int] `ipp:"noise-removal"`
	OutputCompressionQualityFactor optional.Val[int] `ipp:"output-compression-quality-factor"`
}

// InputAttributes represents the "input-attributes" collection.
//
// It is used in scan job operation requests to specify per-job scanning
// parameters, and as the value type of "input-attributes-default" in
// printer description attributes.
//
// See PWG5100.15.
type InputAttributes struct {
	InputAutoExposure         optional.Val[bool]                        `ipp:"input-auto-exposure"`
	InputAutoScaling          optional.Val[bool]                        `ipp:"input-auto-scaling"`
	InputAutoSkewCorrection   optional.Val[bool]                        `ipp:"input-auto-skew-correction"`
	InputBrightness           optional.Val[int]                         `ipp:"input-brightness"`
	InputColorMode            optional.Val[KwInputColorMode]            `ipp:"input-color-mode"`
	InputContentType          optional.Val[KwInputContentType]          `ipp:"input-content-type"`
	InputContrast             optional.Val[int]                         `ipp:"input-contrast"`
	InputFilmScanMode         optional.Val[KwInputFilmScanMode]         `ipp:"input-film-scan-mode"`
	InputImagesToTransfer     optional.Val[int]                         `ipp:"input-images-to-transfer"`
	InputMedia                optional.Val[string]                      `ipp:"input-media"`
	InputOrientationRequested optional.Val[EnInputOrientationRequested] `ipp:"input-orientation-requested"`
	InputQuality              optional.Val[EnInputQuality]              `ipp:"input-quality"`
	InputResolution           optional.Val[goipp.Resolution]            `ipp:"input-resolution"`
	InputScalingHeight        optional.Val[int]                         `ipp:"input-scaling-height"`
	InputScalingWidth         optional.Val[int]                         `ipp:"input-scaling-width"`
	InputScanRegions          []InputScanRegion                         `ipp:"input-scan-regions"`
	InputSharpness            optional.Val[int]                         `ipp:"input-sharpness"`
	InputSides                optional.Val[KwSides]                     `ipp:"input-sides"`
	InputSource               optional.Val[KwInputSource]               `ipp:"input-source"`
}

// ScannerDescription contains scanner-specific printer description
// attributes returned by GetScanServiceElements.
//
// See PWG5100.17, 4.3.
type ScannerDescription struct {
	// PWG5100.15: destination URI schemes
	DestinationURISchemesSupported []string `ipp:"destination-uri-schemes-supported"`

	// PWG5100.15: default values for input-attributes members.
	// Present only if Add-Document-Images operation is supported.
	InputAttributesDefault optional.Val[InputAttributes] `ipp:"input-attributes-default"`

	// PWG5100.15: which input-attributes member attributes are supported.
	InputAttributesSupported []string `ipp:"input-attributes-supported"`

	// PWG5100.15: color mode
	InputColorModeSupported []KwInputColorMode `ipp:"input-color-mode-supported"`

	// PWG5100.15: media
	InputMediaSupported []string `ipp:"input-media-supported"`

	// PWG5100.15: orientation
	InputOrientationRequestedSupported []EnInputOrientationRequested `ipp:"input-orientation-requested-supported"`

	// PWG5100.15: quality (reuses print-quality values from RFC8011, 5.2.13)
	InputQualitySupported []EnInputQuality `ipp:"input-quality-supported"`

	// PWG5100.15: resolution
	InputResolutionSupported []goipp.Resolution `ipp:"input-resolution-supported"`

	// PWG5100.15: sides (reuses KwSides values from RFC8011, 5.2.8)
	InputSidesSupported []KwSides `ipp:"input-sides-supported"`

	// PWG5100.15: input source
	InputSourceSupported []KwInputSource `ipp:"input-source-supported"`

	// PWG5100.17: spooling behavior for scan job document data.
	JobDestinationSpoolingSupported optional.Val[KwJobSpooling] `ipp:"job-destination-spooling-supported"`

	// PWG5100.17: default values for output-attributes members.
	OutputAttributesDefault optional.Val[OutputAttributes] `ipp:"output-attributes-default"`

	// PWG5100.17: which output-attributes member attributes are supported.
	OutputAttributesSupported []string `ipp:"output-attributes-supported"`
}
