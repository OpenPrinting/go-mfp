package cupsraster

import (
	"bufio"
	"fmt"
	"io"
)

// decodeLines decodes the RLE-compressed page data shared by PWG Raster (PWG
// 5102.4 §4.2) and Apple URF into rows of bytesPerLine bytes.  Each line
// group starts with a line-repeat byte (the decoded line applies to repeat+1
// consecutive rows), followed by runs of pixel groups until bytesPerLine
// bytes are produced.  A pixel group is groupSize bytes (1 for bpp <= 8, 3
// for 24-bit RGB).  Control byte c:
//
//   - 0..127: the next group is repeated c+1 times;
//   - 129..255: 257-c literal groups follow;
//   - 128: reserved in PWG 5102.4; URF uses it as "fill the remainder of the
//     line with blank (white)".  Handled defensively for both.
//
// The fill byte for the 0x80 case depends on the colour space of the caller
// (0x00 is white in the black/K space, 0xff is white in sGray/sRGB), so it
// must be supplied explicitly.  setRow is called once per decoded row, in
// order, with a buffer that is only valid until the next call.
func decodeLines(r *bufio.Reader, height, bytesPerLine, groupSize int, fill byte, setRow func(y int, row []byte)) error {
	row := make([]byte, bytesPerLine)
	for y := 0; y < height; {
		repeat, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("row %d: reading line-repeat byte: %w", y, err)
		}
		if err := decodeLine(r, row, groupSize, fill); err != nil {
			return fmt.Errorf("row %d: %w", y, err)
		}
		for n := int(repeat) + 1; n > 0 && y < height; n-- {
			setRow(y, row)
			y++
		}
	}
	return nil
}

// decodeLine decodes a single RLE line into row.
func decodeLine(r *bufio.Reader, row []byte, groupSize int, fill byte) error {
	for pos := 0; pos < len(row); {
		c, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("reading control byte at offset %d: %w", pos, err)
		}
		switch {
		case c < 128:
			// repeated group
			if pos+groupSize > len(row) {
				return fmt.Errorf("group overflows line at offset %d", pos)
			}
			group := row[pos : pos+groupSize]
			if _, err := io.ReadFull(r, group); err != nil {
				return fmt.Errorf("reading repeated group at offset %d: %w", pos, err)
			}
			pos += groupSize
			for n := int(c); n > 0; n-- {
				if pos+groupSize > len(row) {
					return fmt.Errorf("repeated group overflows line at offset %d", pos)
				}
				copy(row[pos:pos+groupSize], group)
				pos += groupSize
			}
		case c > 128:
			// literal groups
			n := 257 - int(c)
			if pos+n*groupSize > len(row) {
				return fmt.Errorf("literal run of %d groups overflows line at offset %d", n, pos)
			}
			if _, err := io.ReadFull(r, row[pos:pos+n*groupSize]); err != nil {
				return fmt.Errorf("reading %d literal groups at offset %d: %w", n, pos, err)
			}
			pos += n * groupSize
		default: // c == 128
			// fill the remainder of the line with blank pixels.
			for ; pos < len(row); pos++ {
				row[pos] = fill
			}
		}
	}
	return nil
}
