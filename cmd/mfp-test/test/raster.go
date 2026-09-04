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

	"github.com/h2non/bimg"
	"github.com/rusq/thermoprint/cupsraster"
)

// convertToPNG converts captured document bytes to a PNG image.
// The format argument is the MIME type of the document (e.g. "image/pwg-raster").
// For multi-page documents, the first page is returned.
func convertToPNG(data []byte, format string) ([]byte, error) {
	switch format {
	case "image/pwg-raster", "image/urf":
		return convertRasterToPNG(data)
	case "application/pdf",
		"image/jpeg",
		"image/tiff",
		"image/webp",
		"image/gif",
		"image/png":
		return convertVipsToPNG(data)
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
