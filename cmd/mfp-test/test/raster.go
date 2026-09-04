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

	"github.com/rusq/thermoprint/cupsraster"
)

// convertToPNG converts captured document bytes to a PNG image.
// The format argument is the MIME type of the document (e.g. "image/pwg-raster").
// For multi-page documents, the first page is returned.
func convertToPNG(data []byte, format string) ([]byte, error) {
	switch format {
	case "image/pwg-raster", "image/urf":
		return convertRasterToPNG(data)
	default:
		return nil, fmt.Errorf("raster: unsupported format %q", format)
	}
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
