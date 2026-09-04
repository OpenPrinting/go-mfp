// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Raster conversion for mfp-test

package test

import (
	"bytes"
	"fmt"
	"image/png"
	"sync"

	"github.com/h2non/bimg"
	"github.com/rusq/thermoprint/cupsraster"
	"gopkg.in/gographics/imagick.v2/imagick"
)

var imagickOnce sync.Once

func initImagick() {
	imagickOnce.Do(imagick.Initialize)
}

// convertToPNG converts captured document bytes to a PNG image.
// The format argument is the MIME type of the document (e.g. "image/pwg-raster").
// For multi-page documents, the first page is returned.
func convertToPNG(data []byte, format string) ([]byte, error) {
	switch format {
	case "image/pwg-raster", "image/urf":
		return convertRasterToPNG(data)
	case "application/pdf",
		"application/vnd.cups-pdf",
		"image/jpeg",
		"image/tiff",
		"image/webp",
		"image/gif",
		"image/png":
		return convertVipsToPNG(data)
	case "application/postscript",
		"application/vnd.cups-postscript":
		return convertPSToPNG(data)
	default:
		return nil, fmt.Errorf("raster: unsupported format %q", format)
	}
}

// convertVipsToPNG uses bimg (libvips) to convert PDF, JPEG, TIFF and other
// common formats to PNG. For multi-page documents, the first page is used.
func convertVipsToPNG(data []byte) ([]byte, error) {
	out, err := bimg.NewImage(data).Convert(bimg.PNG)
	if err != nil {
		return nil, fmt.Errorf("raster: bimg convert: %w", err)
	}
	return out, nil
}

// convertPSToPNG uses ImageMagick (via Ghostscript delegate) to convert
// PostScript to PNG. Only the first page is returned.
func convertPSToPNG(data []byte) ([]byte, error) {
	initImagick()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(data); err != nil {
		return nil, fmt.Errorf("raster: imagick read PS: %w", err)
	}
	if !mw.SetIteratorIndex(0) {
		return nil, fmt.Errorf("raster: imagick: no pages in document")
	}
	if err := mw.SetImageFormat("PNG"); err != nil {
		return nil, fmt.Errorf("raster: imagick set format: %w", err)
	}
	out, err := mw.GetImageBlob()
	if err != nil {
		return nil, fmt.Errorf("raster: imagick get blob: %w", err)
	}
	return out, nil
}

// convertRasterToPNG decodes a PWG Raster or Apple URF stream and encodes
// the first page as PNG.
func convertRasterToPNG(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	pages, err := cupsraster.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("raster: decode: %w", err)
	}
	if len(pages) == 0 {
		return nil, fmt.Errorf("raster: no pages in raster stream")
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, pages[0]); err != nil {
		return nil, fmt.Errorf("raster: encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}
