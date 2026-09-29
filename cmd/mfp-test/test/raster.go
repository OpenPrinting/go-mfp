// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Raster conversion for mfp-test

package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"image/png"
	"io"
	"os"
	"os/exec"

	"github.com/OpenPrinting/go-mfp/cmd/mfp-test/cupsraster"
	"github.com/h2non/bimg"
)

// cropCapturedPNG detects the non-white content bounding box of a full-page
// raster image and crops to that region. CUPS centres the printed image on
// the paper with white margins; this function removes those margins so the
// evaluator compares only the printed content against the reference.
// A full raster scan finds the tightest bounding box of non-white pixels,
// making detection robust even when the content has a light or white centre.
// If no non-white pixels are found (blank page), the data is returned unchanged.
func cropCapturedPNG(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("raster: crop: decode: %w", err)
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	isWhite := func(x, y int) bool {
		r, g, bv, _ := img.At(x, y).RGBA()
		// RGBA returns [0, 65535]; threshold ≈ 245/255
		const t uint32 = 62000
		return r > t && g > t && bv > t
	}

	// Full raster scan: find bounding box of all non-white pixels.
	xMin, xMax, yMin, yMax := w, 0, h, 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !isWhite(x, y) {
				if x < xMin {
					xMin = x
				}
				if x > xMax {
					xMax = x
				}
				if y < yMin {
					yMin = y
				}
				if y > yMax {
					yMax = y
				}
			}
		}
	}

	if xMax <= xMin || yMax <= yMin {
		return data, nil // blank page — nothing to crop
	}

	cropped, err := bimg.NewImage(data).Extract(yMin, xMin, xMax-xMin+1, yMax-yMin+1)
	if err != nil {
		return nil, fmt.Errorf("raster: crop: extract: %w", err)
	}
	return cropped, nil
}

// convertToPNG converts captured document bytes to a PNG image.
// The format argument is the MIME type of the document (e.g. "image/pwg-raster").
// dpi is used for PostScript and PDF rendering via Ghostscript; if zero, 300 DPI is used.
// For multi-page documents, the first page is returned.
// CUPS may gzip-compress documents before delivery; gzip is transparently
// decompressed before format-specific conversion.
func convertToPNG(data []byte, format string, dpi int) ([]byte, error) {
	// Decompress gzip-wrapped data (magic bytes 1f 8b). Loop to handle
	// the case where CUPS applies gzip at both the HTTP and IPP levels.
	for len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("raster: gzip: %w", err)
		}
		data, err = io.ReadAll(r)
		r.Close()
		if err != nil {
			return nil, fmt.Errorf("raster: gzip decompress: %w", err)
		}
	}
	switch format {
	case "image/pwg-raster", "image/urf":
		return convertRasterToPNG(data)
	case "application/pdf",
		"application/vnd.cups-pdf",
		"application/postscript",
		"application/vnd.cups-postscript":
		// Use Ghostscript so the negotiated DPI is honoured. libvips renders
		// PDF at 72 DPI by default, producing a thumbnail-sized image that
		// compares poorly against the 300 DPI reference.
		return convertGSToPNG(data, dpi)
	case "image/jpeg",
		"image/tiff",
		"image/webp",
		"image/gif",
		"image/png":
		return convertVipsToPNG(data)
	case "application/octet-stream":
		// CUPS may deliver any format under the generic octet-stream type.
		// Detect the actual format from the magic bytes.
		return detectAndConvert(data, dpi)
	default:
		// image/vnd.cups-raster and image/jpeg+gzip are not yet supported
		// for image evaluation; captured bytes are still saved with --keep.
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

// convertGSToPNG calls Ghostscript to render the first page of a PostScript
// or PDF document to PNG at the negotiated printer DPI. Using gs ensures the
// captured image is the same physical size as the reference (both at dpi pixels
// per inch), which is required for accurate histogram and SSIM comparison.
// dpi is the render resolution; if zero, 300 DPI is used as a safe default.
func convertGSToPNG(data []byte, dpi int) ([]byte, error) {
	if dpi <= 0 {
		dpi = 300
	}
	// Write PostScript data to a temp input file.
	inFile, err := os.CreateTemp("", "mfp-ps-*.ps")
	if err != nil {
		return nil, fmt.Errorf("raster: gs: create input temp: %w", err)
	}
	defer os.Remove(inFile.Name())
	if _, err := inFile.Write(data); err != nil {
		inFile.Close()
		return nil, fmt.Errorf("raster: gs: write PS: %w", err)
	}
	if err := inFile.Close(); err != nil {
		return nil, fmt.Errorf("raster: gs: close input: %w", err)
	}

	// Create a temp file path for the PNG output.
	outFile, err := os.CreateTemp("", "mfp-png-*.png")
	if err != nil {
		return nil, fmt.Errorf("raster: gs: create output temp: %w", err)
	}
	outPath := outFile.Name()
	outFile.Close()
	defer os.Remove(outPath)

	// Run Ghostscript: render only the first page at the negotiated DPI.
	cmd := exec.Command("gs",
		"-dBATCH", "-dNOPAUSE", "-dQUIET",
		"-sDEVICE=png16m",
		fmt.Sprintf("-r%d", dpi),
		"-dFirstPage=1", "-dLastPage=1",
		"-sOutputFile="+outPath,
		inFile.Name(),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("raster: gs: %w: %s", err, out)
	}

	return os.ReadFile(outPath)
}

// detectAndConvert detects the actual format of data from its magic bytes
// and dispatches to the appropriate converter. Used for application/octet-stream
// where CUPS may send any format.
func detectAndConvert(data []byte, dpi int) ([]byte, error) {
	n := len(data)
	switch {
	case n >= 4 && string(data[:4]) == "RaS2":
		return convertRasterToPNG(data)
	case n >= 8 && string(data[:8]) == "UNIRAST\x00":
		return convertRasterToPNG(data)
	case n >= 4 && string(data[:4]) == "%PDF":
		return convertGSToPNG(data, dpi)
	case n >= 2 && string(data[:2]) == "%!":
		return convertGSToPNG(data, dpi)
	case n >= 2 && data[0] == 0xff && data[1] == 0xd8:
		return convertVipsToPNG(data) // JPEG
	case n >= 4 && string(data[:4]) == "\x89PNG":
		return convertVipsToPNG(data)
	}
	return nil, fmt.Errorf("raster: cannot detect format of application/octet-stream data")
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

