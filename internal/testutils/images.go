// MFP - Miulti-Function Printers and scanners toolkit
// Utility functions and data BLOBs for testing
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package testutils

import (
	// Import "embed" for its side effects
	_ "embed"
)

// Images contains samples of test images of various formats and sizes
var Images struct {
	// Small images in various formats
	BMP100x75             []byte
	JPEG100x75            []byte
	PDF100x75             []byte
	PNG100x75rgb8         []byte
	PNG100x75rgb8i        []byte
	PNG100x75rgb8paletted []byte
	PNG100x75rgb16        []byte
	PNG100x75gray1        []byte
	PNG100x75gray8        []byte
	PNG100x75gray16       []byte
	TIFF100x75            []byte

	// This page is suitable as both A4 and Letter image sample
	// at 600 DPI:
	//   Letter: 5100x6600
	//   A4:     4960x7016
	PNG5100x7016 []byte
}

//go:embed "data/UEIT-100x75.bmp"
var imageBMP100x75 []byte

//go:embed "data/UEIT-100x75.jpeg"
var imageJPEG100x75 []byte

//go:embed "data/UEIT-100x75.pdf"
var imagePDF100x75 []byte

//go:embed "data/UEIT-100x75rgb8.png"
var imagePNG100x75rgb8 []byte

//go:embed "data/UEIT-100x75rgb8i.png"
var imagePNG100x75rgb8i []byte

//go:embed "data/UEIT-100x75rgb8paletted.png"
var imagePNG100x75rgb8paletted []byte

//go:embed "data/UEIT-100x75rgb16.png"
var imagePNG100x75rgb16 []byte

//go:embed "data/UEIT-100x75gray1.png"
var imagePNG100x75gray1 []byte

//go:embed "data/UEIT-100x75gray8.png"
var imagePNG100x75gray8 []byte

//go:embed "data/UEIT-100x75gray16.png"
var imagePNG100x75gray16 []byte

//go:embed "data/UEIT-100x75.tiff"
var imageTIFF100x75 []byte

//go:embed "data/testpage-5100x7016.png"
var imagePNG5100x7016 []byte

func init() {
	Images.BMP100x75 = imageBMP100x75
	Images.JPEG100x75 = imageJPEG100x75
	Images.PDF100x75 = imagePDF100x75
	Images.PNG100x75rgb8 = imagePNG100x75rgb8
	Images.PNG100x75rgb8i = imagePNG100x75rgb8i
	Images.PNG100x75rgb8paletted = imagePNG100x75rgb8paletted
	Images.PNG100x75rgb16 = imagePNG100x75rgb16
	Images.PNG100x75gray1 = imagePNG100x75gray1
	Images.PNG100x75gray8 = imagePNG100x75gray8
	Images.PNG100x75gray16 = imagePNG100x75gray16
	Images.TIFF100x75 = imageTIFF100x75
	Images.PNG5100x7016 = imagePNG5100x7016
}
