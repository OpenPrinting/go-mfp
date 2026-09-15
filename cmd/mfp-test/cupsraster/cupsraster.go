// Package cupsraster decodes the raster streams produced by client-side
// rasterisation in CUPS and macOS/iOS printing: PWG Raster (PWG 5102.4,
// image/pwg-raster) and Apple Raster (URF, image/urf).  Both formats carry
// one or more pre-rendered pages compressed with the same simple run-length
// scheme; the decoder converts each page to an image.Image.
//
// References:
//   - https://ftp.pwg.org/pub/pwg/candidates/cs-ippraster10-20120420-5102.4.pdf
//   - https://openprinting.github.io/driverless/01-standards-and-their-pdls/#apple-raster
//
// This package is derived from github.com/rusq/thermoprint/cupsraster,
// which is copyright (c) 2025 github.com/rusq, MIT License.
// See CREDITS.txt at the root of this repository.
package cupsraster

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
)

// Format identifies a supported raster stream format.
type Format int

const (
	FormatUnknown Format = iota
	FormatPWG            // PWG Raster (image/pwg-raster)
	FormatURF            // Apple Raster (image/urf)
)

func (f Format) String() string {
	switch f {
	case FormatPWG:
		return "PWG"
	case FormatURF:
		return "URF"
	}
	return "unknown"
}

// maxDim caps page dimensions to guard against corrupt headers.
const maxDim = 32768

// maxPixels caps total page size (width*height) to about 512 MiB of 8-bit
// gray, which is far beyond anything a label printer will see.
const maxPixels = 1 << 29

// Detect sniffs the magic bytes of data and reports the raster format, or
// FormatUnknown.  PWG raster is identified by both the sync word and the
// PwgRaster header magic, which distinguishes it from little-endian CUPS
// raster streams sharing the sync word.
func Detect(data []byte) Format {
	if len(data) >= len(pwgSyncWord)+len(pwgMagic) &&
		string(data[:len(pwgSyncWord)]) == pwgSyncWord &&
		string(data[len(pwgSyncWord):len(pwgSyncWord)+len(pwgMagic)]) == pwgMagic {
		return FormatPWG
	}
	if len(data) >= len(urfMagic) && string(data[:len(urfMagic)]) == urfMagic {
		return FormatURF
	}
	return FormatUnknown
}

// Page is a decoded raster page together with the resolution declared in
// its header.  A consumer printing at a different resolution must scale the
// image accordingly to preserve the physical size.
type Page struct {
	image.Image
	XDPI, YDPI int
}

// Decode sniffs the stream format and decodes all pages.
func Decode(r io.Reader) ([]image.Image, error) {
	pages, err := DecodePages(r)
	if err != nil {
		return nil, err
	}
	return images(pages), nil
}

// DecodePages sniffs the stream format and decodes all pages with their
// declared resolutions.
func DecodePages(r io.Reader) ([]Page, error) {
	br := bufio.NewReader(r)
	head, err := br.Peek(len(pwgSyncWord) + len(pwgMagic))
	if err != nil && len(head) == 0 {
		return nil, fmt.Errorf("reading stream header: %w", err)
	}
	switch Detect(head) {
	case FormatPWG:
		return decodePWGPages(br)
	case FormatURF:
		return decodeURFPages(br)
	}
	return nil, fmt.Errorf("unrecognised raster stream (header % x)", head)
}

func images(pages []Page) []image.Image {
	imgs := make([]image.Image, len(pages))
	for i, pg := range pages {
		imgs[i] = pg.Image
	}
	return imgs
}

func checkDimensions(width, height int) error {
	if width <= 0 || height <= 0 || width > maxDim || height > maxDim || width*height > maxPixels {
		return fmt.Errorf("invalid page dimensions %dx%d", width, height)
	}
	return nil
}

// decodeGrayPage decodes a 1- or 8-bit single-channel page into image.Gray.
// blackOne selects ink semantics (K color space: max value = black) as
// opposed to luminance semantics (sGray: zero = black).
func decodeGrayPage(br *bufio.Reader, width, height, bpp, bytesPerLine int, blackOne bool) (*image.Gray, error) {
	img := image.NewGray(image.Rect(0, 0, width, height))
	// The RLE blank filler must be white in the page's own semantics.
	fill := byte(0xff)
	if blackOne {
		fill = 0x00
	}
	err := decodeLines(br, height, bytesPerLine, 1, fill, func(y int, row []byte) {
		dst := img.Pix[y*img.Stride : y*img.Stride+width]
		switch bpp {
		case 1:
			for x := 0; x < width; x++ {
				bit := row[x/8]>>(7-x%8)&1 == 1
				if bit == blackOne {
					dst[x] = 0x00 // black
				} else {
					dst[x] = 0xff // white
				}
			}
		case 8:
			copy(dst, row)
			if blackOne {
				for x := 0; x < len(dst); x++ {
					dst[x] = 0xff - dst[x]
				}
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return img, nil
}

// decodeRGBPage decodes a 24-bit chunky RGB page into image.NRGBA.
func decodeRGBPage(br *bufio.Reader, width, height, bytesPerLine int) (*image.NRGBA, error) {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	err := decodeLines(br, height, bytesPerLine, 3, 0xff, func(y int, row []byte) {
		for x := 0; x < width; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: row[x*3], G: row[x*3+1], B: row[x*3+2], A: 0xff})
		}
	})
	if err != nil {
		return nil, err
	}
	return img, nil
}
