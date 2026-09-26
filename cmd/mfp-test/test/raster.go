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

// convertToPNG converts captured document bytes to a PNG image.
// The format argument is the MIME type of the document (e.g. "image/pwg-raster").
// dpi is used for PostScript rendering; if zero, 300 DPI is used.
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
		"image/jpeg",
		"image/tiff",
		"image/webp",
		"image/gif",
		"image/png":
		return convertVipsToPNG(data)
	case "application/postscript",
		"application/vnd.cups-postscript":
		return convertPSToPNG(data, dpi)
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

// convertPSToPNG calls Ghostscript directly to convert the first page of a
// PostScript document to PNG. Using gs avoids the ImageMagick dependency and
// the Ubuntu policy.xml reconfiguration it requires.
// dpi is the render resolution; if zero, 300 DPI is used as a safe default.
func convertPSToPNG(data []byte, dpi int) ([]byte, error) {
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
		return convertVipsToPNG(data)
	case n >= 2 && string(data[:2]) == "%!":
		return convertPSToPNG(data, dpi)
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
