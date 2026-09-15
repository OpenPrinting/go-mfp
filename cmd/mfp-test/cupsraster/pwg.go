package cupsraster

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
)

// PWG Raster stream layout (PWG 5102.4): a 4-byte synchronisation word
// "RaS2", then for each page a 1796-byte big-endian header followed by
// RLE-compressed page data.

const (
	pwgSyncWord   = "RaS2"
	pwgMagic      = "PwgRaster\x00"
	pwgHeaderSize = 1796
)

// Header field byte offsets (PWG 5102.4 §4.3, empirically verified against
// macOS cupsfilter output).
const (
	pwgOffHWResolutionX = 276
	pwgOffHWResolutionY = 280
	pwgOffWidth         = 372
	pwgOffHeight        = 376
	pwgOffBitsPerColor  = 384
	pwgOffBitsPerPixel  = 388
	pwgOffBytesPerLine  = 392
	pwgOffColorOrder    = 396
	pwgOffColorSpace    = 400
)

// PWG cupsColorSpace values (subset the decoder understands).
const (
	pwgCSBlack    = 3  // K, ink semantics: 1 = black
	pwgCSSGray    = 18 // sGray, luminance semantics: 0 = black
	pwgCSSRGB     = 19 // sRGB
	pwgCSAdobeRGB = 20 // AdobeRGB, treated as sRGB
)

type pwgHeader struct {
	Width, Height              int
	BitsPerColor, BitsPerPixel int
	BytesPerLine               int
	ColorOrder                 int
	ColorSpace                 int
	XRes, YRes                 int
}

func parsePWGHeader(buf []byte) (pwgHeader, error) {
	u32 := func(off int) int {
		return int(binary.BigEndian.Uint32(buf[off : off+4]))
	}
	h := pwgHeader{
		Width:        u32(pwgOffWidth),
		Height:       u32(pwgOffHeight),
		BitsPerColor: u32(pwgOffBitsPerColor),
		BitsPerPixel: u32(pwgOffBitsPerPixel),
		BytesPerLine: u32(pwgOffBytesPerLine),
		ColorOrder:   u32(pwgOffColorOrder),
		ColorSpace:   u32(pwgOffColorSpace),
		XRes:         u32(pwgOffHWResolutionX),
		YRes:         u32(pwgOffHWResolutionY),
	}
	if err := checkDimensions(h.Width, h.Height); err != nil {
		return h, err
	}
	if h.ColorOrder != 0 {
		return h, fmt.Errorf("unsupported color order %d (only chunky is valid)", h.ColorOrder)
	}
	switch h.ColorSpace {
	case pwgCSBlack, pwgCSSGray:
		if h.BitsPerPixel != 1 && h.BitsPerPixel != 8 {
			return h, fmt.Errorf("unsupported bits per pixel %d for color space %d", h.BitsPerPixel, h.ColorSpace)
		}
	case pwgCSSRGB, pwgCSAdobeRGB:
		if h.BitsPerPixel != 24 {
			return h, fmt.Errorf("unsupported bits per pixel %d for color space %d", h.BitsPerPixel, h.ColorSpace)
		}
	default:
		return h, fmt.Errorf("unsupported color space %d", h.ColorSpace)
	}
	if want := (h.Width*h.BitsPerPixel + 7) / 8; h.BytesPerLine != want {
		return h, fmt.Errorf("bytes per line %d inconsistent with width %d at %d bpp (want %d)", h.BytesPerLine, h.Width, h.BitsPerPixel, want)
	}
	return h, nil
}

// DecodePWG decodes a PWG Raster stream into one image per page.
func DecodePWG(r io.Reader) ([]image.Image, error) {
	pages, err := decodePWGPages(bufio.NewReader(r))
	if err != nil {
		return nil, err
	}
	return images(pages), nil
}

func decodePWGPages(br *bufio.Reader) ([]Page, error) {
	sync := make([]byte, len(pwgSyncWord))
	if _, err := io.ReadFull(br, sync); err != nil {
		return nil, fmt.Errorf("reading sync word: %w", err)
	}
	if string(sync) != pwgSyncWord {
		return nil, fmt.Errorf("not a PWG raster stream: sync word %q", sync)
	}
	var pages []Page
	hdr := make([]byte, pwgHeaderSize)
	for page := 1; ; page++ {
		if _, err := io.ReadFull(br, hdr); err != nil {
			if err == io.EOF && page > 1 {
				break // clean end of stream
			}
			return nil, fmt.Errorf("page %d: reading header: %w", page, err)
		}
		if string(hdr[:len(pwgMagic)]) != pwgMagic {
			return nil, fmt.Errorf("page %d: header does not start with %q", page, pwgMagic[:len(pwgMagic)-1])
		}
		h, err := parsePWGHeader(hdr)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}
		img, err := decodePWGPage(br, h)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}
		pages = append(pages, Page{Image: img, XDPI: h.XRes, YDPI: h.YRes})
	}
	if len(pages) == 0 {
		return nil, errors.New("no pages in PWG raster stream")
	}
	return pages, nil
}

func decodePWGPage(br *bufio.Reader, h pwgHeader) (image.Image, error) {
	blackOne := h.ColorSpace == pwgCSBlack
	switch h.BitsPerPixel {
	case 1, 8:
		return decodeGrayPage(br, h.Width, h.Height, h.BitsPerPixel, h.BytesPerLine, blackOne)
	case 24:
		return decodeRGBPage(br, h.Width, h.Height, h.BytesPerLine)
	}
	panic("unreachable: bpp validated in parsePWGHeader")
}
