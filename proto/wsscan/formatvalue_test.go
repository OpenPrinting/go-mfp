// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Test for format value

package wsscan

import "testing"

// TestFormatValue_KnownConstants verifies that known constants map to the
// expected string representations and round-trip through DecodeFormatValue.
func TestFormatValue_KnownConstants(t *testing.T) {
	cases := []struct {
		v       FormatValue
		strRepr string
	}{
		{DIB, "dib"},
		{EXIF, "exif"},
		{JBIG, "jbig"},
		{JFIF, "jfif"},
		{JPEG2K, "jpeg2k"},
		{PDFA, "pdf-a"},
		{PNG, "png"},
		{TIFFSingleUncompressed, "tiff-single-uncompressed"},
		{TIFFSingleG4, "tiff-single-g4"},
		{TIFFSingleG3MH, "tiff-single-g3mh"},
		{TIFFSingleJPEGTN2, "tiff-single-jpeg-tn2"},
		{TIFFMultiUncompressed, "tiff-multi-uncompressed"},
		{TIFFMultiG4, "tiff-multi-g4"},
		{TIFFMultiG3MH, "tiff-multi-g3mh"},
		{TIFFMultiJPEGTN2, "tiff-multi-jpeg-tn2"},
		{XPS, "xps"},
	}

	for _, c := range cases {
		if got := c.v.String(); got != c.strRepr {
			t.Errorf("FormatValue(%q).String(): expected %q, got %q",
				string(c.v), c.strRepr, got)
		}

		if decoded := DecodeFormatValue(c.strRepr); decoded != c.v {
			t.Errorf("DecodeFormatValue(%q): expected %v, got %v",
				c.strRepr, c.v, decoded)
		}
	}
}

// TestFormatValue_MimeType verifies that known format values map to the
// expected MIME types and unknown/vendor-defined values return "".
func TestFormatValue_MimeType(t *testing.T) {
	cases := []struct {
		v        FormatValue
		mimeType string
	}{
		{DIB, "image/bmp"},
		{EXIF, "image/jpeg"},
		{JFIF, "image/jpeg"},
		{JBIG, "image/jbig"},
		{JPEG2K, "image/jp2"},
		{PDFA, "application/pdf"},
		{PNG, "image/png"},
		{TIFFSingleUncompressed, "image/tiff"},
		{TIFFSingleG4, "image/tiff"},
		{TIFFSingleG3MH, "image/tiff"},
		{TIFFSingleJPEGTN2, "image/tiff"},
		{TIFFMultiUncompressed, "image/tiff"},
		{TIFFMultiG4, "image/tiff"},
		{TIFFMultiG3MH, "image/tiff"},
		{TIFFMultiJPEGTN2, "image/tiff"},
		{XPS, "application/vnd.ms-xpsdocument"},
		// Unknown/vendor-defined:
		{UnknownFormatValue, ""},
		{FormatValue("vendor/foo-format"), ""},
	}

	for _, c := range cases {
		if got := c.v.MimeType(); got != c.mimeType {
			t.Errorf("FormatValue(%q).MimeType(): expected %q, got %q",
				string(c.v), c.mimeType, got)
		}
	}
}

// TestFormatValue_VendorDefined verifies that arbitrary vendor-defined values
// are preserved, rather than collapsed into UnknownFormatValue.
func TestFormatValue_VendorDefined(t *testing.T) {
	const vendor = "vendor/foo-format"

	fv := DecodeFormatValue(vendor)
	if string(fv) != vendor {
		t.Fatalf("expected underlying value %q, got %q", vendor, string(fv))
	}

	if fv.String() != vendor {
		t.Fatalf("String(): expected %q, got %q", vendor, fv.String())
	}
}

// TestFormatValue_EmptyString verifies that an empty string is preserved and
// does not receive any special "Unknown" handling.
func TestFormatValue_EmptyString(t *testing.T) {
	fv := DecodeFormatValue("")

	if string(fv) != "" {
		t.Fatalf("expected underlying value %q, got %q", "", string(fv))
	}

	if fv.String() != "" {
		t.Fatalf("String(): expected %q, got %q", "", fv.String())
	}
}
