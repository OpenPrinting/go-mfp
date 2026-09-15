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
	"os"
	"os/exec"

	"github.com/OpenPrinting/go-mfp/cmd/mfp-test/cupsraster"
	"github.com/h2non/bimg"
)

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
func convertPSToPNG(data []byte) ([]byte, error) {
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

	// Run Ghostscript: render only the first page at 150 dpi.
	cmd := exec.Command("gs",
		"-dBATCH", "-dNOPAUSE", "-dQUIET",
		"-sDEVICE=png16m",
		"-r150",
		"-dFirstPage=1", "-dLastPage=1",
		"-sOutputFile="+outPath,
		inFile.Name(),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("raster: gs: %w: %s", err, out)
	}

	return os.ReadFile(outPath)
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
