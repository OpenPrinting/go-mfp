package cupsraster

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
)

// Apple Raster (URF) stream layout: the 8-byte magic "UNIRAST\x00" and a
// big-endian uint32 page count, then for each page a 32-byte header followed
// by RLE-compressed page data.  Unlike PWG raster there is no bytes-per-line
// field; it is derived from width and bits per pixel.

const (
	urfMagic          = "UNIRAST\x00"
	urfPageHeaderSize = 32
)

// URF color space values (subset the decoder understands).
const (
	urfCSSGray = 0 // sGray, luminance semantics: 0 = black
	urfCSSRGB  = 1 // sRGB
)

type urfHeader struct {
	BitsPerPixel  int
	ColorSpace    int
	Width, Height int
	DPI           int
}

func parseURFHeader(buf []byte) (urfHeader, error) {
	h := urfHeader{
		BitsPerPixel: int(buf[0]),
		ColorSpace:   int(buf[1]),
		Width:        int(binary.BigEndian.Uint32(buf[12:16])),
		Height:       int(binary.BigEndian.Uint32(buf[16:20])),
		DPI:          int(binary.BigEndian.Uint32(buf[20:24])),
	}
	if err := checkDimensions(h.Width, h.Height); err != nil {
		return h, err
	}
	switch h.ColorSpace {
	case urfCSSGray:
		if h.BitsPerPixel != 1 && h.BitsPerPixel != 8 {
			return h, fmt.Errorf("unsupported bits per pixel %d for sGray", h.BitsPerPixel)
		}
	case urfCSSRGB:
		if h.BitsPerPixel != 24 {
			return h, fmt.Errorf("unsupported bits per pixel %d for sRGB", h.BitsPerPixel)
		}
	default:
		return h, fmt.Errorf("unsupported color space %d", h.ColorSpace)
	}
	return h, nil
}

// DecodeURF decodes an Apple Raster (URF) stream into one image per page.
func DecodeURF(r io.Reader) ([]image.Image, error) {
	pages, err := decodeURFPages(bufio.NewReader(r))
	if err != nil {
		return nil, err
	}
	return images(pages), nil
}

func decodeURFPages(br *bufio.Reader) ([]Page, error) {
	head := make([]byte, len(urfMagic)+4)
	if _, err := io.ReadFull(br, head); err != nil {
		return nil, fmt.Errorf("reading URF header: %w", err)
	}
	if string(head[:len(urfMagic)]) != urfMagic {
		return nil, fmt.Errorf("not a URF stream: magic % x", head[:len(urfMagic)])
	}
	numPages := int(binary.BigEndian.Uint32(head[len(urfMagic):]))
	if numPages > 65535 {
		return nil, fmt.Errorf("invalid URF page count %d", numPages)
	}
	// Some CUPS versions stream URF without updating the page count header,
	// leaving it as 0. When that happens, read pages until EOF.
	pages := make([]Page, 0, max(numPages, 1))
	hdr := make([]byte, urfPageHeaderSize)
	for page := 1; numPages == 0 || page <= numPages; page++ {
		if _, err := io.ReadFull(br, hdr); err != nil {
			if numPages == 0 && (err == io.EOF || err == io.ErrUnexpectedEOF) {
				break // normal end-of-stream when count was 0
			}
			return nil, fmt.Errorf("page %d: reading page header: %w", page, err)
		}
		h, err := parseURFHeader(hdr)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}
		img, err := decodeURFPage(br, h)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}
		pages = append(pages, Page{Image: img, XDPI: h.DPI, YDPI: h.DPI})
	}
	if len(pages) == 0 {
		return nil, errors.New("no pages in URF stream")
	}
	return pages, nil
}

func decodeURFPage(br *bufio.Reader, h urfHeader) (image.Image, error) {
	bytesPerLine := (h.Width*h.BitsPerPixel + 7) / 8
	switch h.BitsPerPixel {
	case 1:
		return decodeGrayPage(br, h.Width, h.Height, h.BitsPerPixel, bytesPerLine, true)
	case 8:
		return decodeGrayPage(br, h.Width, h.Height, h.BitsPerPixel, bytesPerLine, false)
	case 24:
		return decodeRGBPage(br, h.Width, h.Height, bytesPerLine)
	}
	panic("unreachable: bpp validated in parseURFHeader")
}
