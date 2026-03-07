// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// format value

package wsscan

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// FormatValue defines the file format and compression
// type supported by the scan device.
//
// It is backed by string so that it can represent both standard values defined
// in the WS-Scan specification and vendor-specific extensions that may appear
// in the wild.
type FormatValue string

// Known format values (constants for well-known types).
// Additional values may exist in real devices; those are represented by
// arbitrary [FormatValue] values constructed from their string form.
const (
	// UnknownFormatValue is used when the format string is empty or otherwise
	// invalid. For any non-empty string we keep the original value instead of
	// collapsing it to UnknownFormatValue, so that vendor-defined values are
	// preserved.
	UnknownFormatValue FormatValue = ""

	DIB                    FormatValue = "dib"
	EXIF                   FormatValue = "exif"
	JBIG                   FormatValue = "jbig"
	JFIF                   FormatValue = "jfif"
	JPEG2K                 FormatValue = "jpeg2k"
	PDFA                   FormatValue = "pdf-a"
	PNG                    FormatValue = "png"
	TIFFSingleUncompressed FormatValue = "tiff-single-uncompressed"
	TIFFSingleG4           FormatValue = "tiff-single-g4"
	TIFFSingleG3MH         FormatValue = "tiff-single-g3mh"
	TIFFSingleJPEGTN2      FormatValue = "tiff-single-jpeg-tn2"
	TIFFMultiUncompressed  FormatValue = "tiff-multi-uncompressed"
	TIFFMultiG4            FormatValue = "tiff-multi-g4"
	TIFFMultiG3MH          FormatValue = "tiff-multi-g3mh"
	TIFFMultiJPEGTN2       FormatValue = "tiff-multi-jpeg-tn2"
	XPS                    FormatValue = "xps"
)

// decodeFormatValue decodes [FormatValue] from the XML tree.
func decodeFormatValue(root xmldoc.Element) (fv FormatValue, err error) {
	if root.Text == "" {
		err = fmt.Errorf("invalid FormatValue: empty")
		err = xmldoc.XMLErrWrap(root, err)
		return UnknownFormatValue, err
	}
	return FormatValue(root.Text), nil
}

// toXML generates XML tree for the [FormatValue].
func (fv FormatValue) toXML(name string) xmldoc.Element {
	return xmldoc.Element{
		Name: name,
		Text: fv.String(),
	}
}

// String returns a string representation of the [FormatValue].
func (fv FormatValue) String() string {
	return string(fv)
}

// MimeType returns MIME type string for the value, or
// empty string for the unknown or vendor-defined value.
func (fv FormatValue) MimeType() string {
	switch fv {
	case DIB:
		return "image/bmp"
	case EXIF, JFIF:
		return "image/jpeg"
	case JBIG:
		return "image/jbig"
	case JPEG2K:
		return "image/jp2"
	case PDFA:
		return "application/pdf"
	case PNG:
		return "image/png"
	case TIFFSingleUncompressed,
		TIFFSingleG4,
		TIFFSingleG3MH,
		TIFFSingleJPEGTN2,
		TIFFMultiUncompressed,
		TIFFMultiG4,
		TIFFMultiG3MH,
		TIFFMultiJPEGTN2:
		return "image/tiff"
	case XPS:
		return "application/vnd.ms-xpsdocument"
	default:
		return ""
	}
}

// DecodeFormatValue decodes [FormatValue] out of its XML string representation.
func DecodeFormatValue(s string) FormatValue {
	return FormatValue(s)
}
