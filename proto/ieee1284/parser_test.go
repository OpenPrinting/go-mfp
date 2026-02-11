// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Print job stream parser tests

package ieee1284

import (
	"bytes"
	"testing"
)

// TestParserPJLPostScript tests PJL-wrapped PostScript parsing.
func TestParserPJLPostScript(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	psContent := "%!PS-Adobe-3.0\n" +
		"%%Title: Test Page\n" +
		"/Helvetica findfont 24 scalefont setfont\n" +
		"72 720 moveto (Hello, World!) show\n" +
		"showpage\n" +
		"%%EOF\n"

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL SET COPIES=1\r\n")
	job.WriteString("@PJL ENTER LANGUAGE=POSTSCRIPT\r\n")
	job.WriteString(psContent)
	job.WriteByte(0x04)
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL EOJ NAME=\"test\"\r\n")
	job.WriteString(uel)

	writeInChunks(t, p, job.Bytes(), 512)

	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}
	if results[0].format != DocFormatPostScript {
		t.Errorf("format = %v, want PostScript", results[0].format)
	}
	if !bytes.Equal(results[0].data, []byte(psContent)) {
		t.Errorf("document data mismatch:\ngot:  %q\nwant: %q",
			results[0].data, psContent)
	}
}

// TestParserPJLPDF tests PJL-wrapped PDF parsing.
func TestParserPJLPDF(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	pdfContent := "%PDF-1.4\n" +
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n" +
		"2 0 obj\n<< /Type /Pages /Kids [] /Count 0 >>\nendobj\n" +
		"xref\n0 3\n" +
		"trailer\n<< /Size 3 /Root 1 0 R >>\n" +
		"startxref\n0\n" +
		"%%EOF\n"

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL SET COPIES=1\r\n")
	job.WriteString("@PJL ENTER LANGUAGE=PDF\r\n")
	job.WriteString(pdfContent)
	job.WriteByte(0x04)
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL EOJ NAME=\"pdf-test\"\r\n")
	job.WriteString(uel)

	writeInChunks(t, p, job.Bytes(), 512)

	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}
	if results[0].format != DocFormatPDF {
		t.Errorf("format = %v, want PDF", results[0].format)
	}
	if !bytes.Equal(results[0].data, []byte(pdfContent)) {
		t.Errorf("document data mismatch:\ngot:  %q\nwant: %q",
			results[0].data, pdfContent)
	}
}

// TestParserRawPostScript tests a raw PostScript document
// without PJL wrapping.
func TestParserRawPostScript(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	psContent := "%!PS-Adobe-3.0\n" +
		"%%Title: Raw Test\n" +
		"/Courier findfont 12 scalefont setfont\n" +
		"72 720 moveto (Raw PostScript test) show\n" +
		"showpage\n" +
		"%%EOF\n"

	var job bytes.Buffer
	job.WriteString(psContent)
	job.WriteByte(0x04) // Ctrl-D terminates the document

	writeInChunks(t, p, job.Bytes(), 512)

	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}
	if results[0].format != DocFormatPostScript {
		t.Errorf("format = %v, want PostScript", results[0].format)
	}
	if !bytes.Equal(results[0].data, []byte(psContent)) {
		t.Errorf("document data mismatch:\ngot:  %q\nwant: %q",
			results[0].data, psContent)
	}
}

// TestParserMultipleJobs tests two back-to-back PJL jobs in a
// single stream.
func TestParserMultipleJobs(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	psContent := "%!PS-Adobe-3.0\n" +
		"72 720 moveto (Job 1) show showpage\n" +
		"%%EOF\n"

	pdfContent := "%PDF-1.4\n" +
		"... pdf content ...\n" +
		"%%EOF\n"

	var stream bytes.Buffer
	// Job 1: PostScript
	stream.WriteString(uel + "@PJL\r\n")
	stream.WriteString("@PJL JOB NAME=\"job1\"\r\n")
	stream.WriteString("@PJL ENTER LANGUAGE=POSTSCRIPT\r\n")
	stream.WriteString(psContent)
	stream.WriteByte(0x04)
	stream.WriteString(uel + "@PJL\r\n")
	stream.WriteString("@PJL EOJ NAME=\"job1\"\r\n")

	// Job 2: PDF
	stream.WriteString(uel + "@PJL\r\n")
	stream.WriteString("@PJL JOB NAME=\"job2\"\r\n")
	stream.WriteString("@PJL ENTER LANGUAGE=PDF\r\n")
	stream.WriteString(pdfContent)
	stream.WriteByte(0x04)
	stream.WriteString(uel + "@PJL\r\n")
	stream.WriteString("@PJL EOJ NAME=\"job2\"\r\n")
	stream.WriteString(uel)

	writeInChunks(t, p, stream.Bytes(), 512)

	if len(results) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(results))
	}

	if results[0].format != DocFormatPostScript {
		t.Errorf("job 1: format = %v, want PostScript",
			results[0].format)
	}
	if !bytes.Equal(results[0].data, []byte(psContent)) {
		t.Errorf("job 1: data mismatch:\ngot:  %q\nwant: %q",
			results[0].data, psContent)
	}

	if results[1].format != DocFormatPDF {
		t.Errorf("job 2: format = %v, want PDF",
			results[1].format)
	}
	if !bytes.Equal(results[1].data, []byte(pdfContent)) {
		t.Errorf("job 2: data mismatch:\ngot:  %q\nwant: %q",
			results[1].data, pdfContent)
	}
}

// TestParserChunkedUEL tests that UEL split across two Write()
// calls is still detected.
func TestParserChunkedUEL(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	psContent := "%!PS-Adobe-3.0\nshowpage\n%%EOF\n"

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL ENTER LANGUAGE=POSTSCRIPT\r\n")
	job.WriteString(psContent)

	// Split the UEL across two Write() calls.
	// UEL = \x1b%-12345X (9 bytes). Split at byte 4.
	endUEL := []byte(uel)
	part1 := endUEL[:4] // \x1b%-1
	part2 := endUEL[4:] // 2345X

	jobBytes := job.Bytes()

	// Write the job content + first part of UEL
	chunk1 := make([]byte, len(jobBytes)+len(part1))
	copy(chunk1, jobBytes)
	copy(chunk1[len(jobBytes):], part1)

	n, err := p.Write(chunk1)
	if err != nil {
		t.Fatalf("Write chunk1 failed: %v", err)
	}
	if n != len(chunk1) {
		t.Fatalf("Write chunk1: got %d, want %d", n, len(chunk1))
	}

	// No document should be emitted yet (partial UEL)
	if len(results) != 0 {
		t.Fatalf("expected 0 documents after chunk1, got %d",
			len(results))
	}

	// Write the rest of UEL + PJL EOJ + final UEL
	var chunk2 bytes.Buffer
	chunk2.Write(part2)
	chunk2.WriteString("@PJL\r\n")
	chunk2.WriteString("@PJL EOJ NAME=\"test\"\r\n")
	chunk2.WriteString(uel)

	n, err = p.Write(chunk2.Bytes())
	if err != nil {
		t.Fatalf("Write chunk2 failed: %v", err)
	}
	if n != chunk2.Len() {
		t.Fatalf("Write chunk2: got %d, want %d", n, chunk2.Len())
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}
	if results[0].format != DocFormatPostScript {
		t.Errorf("format = %v, want PostScript", results[0].format)
	}
	if !bytes.Equal(results[0].data, []byte(psContent)) {
		t.Errorf("document data mismatch:\ngot:  %q\nwant: %q",
			results[0].data, psContent)
	}
}

// TestParserChunkedPJLLine tests that a PJL line split across
// two Write() calls is still parsed.
func TestParserChunkedPJLLine(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	psContent := "%!PS-Adobe-3.0\nshowpage\n%%EOF\n"

	// First chunk: UEL + @PJL + partial ENTER line
	enterLine := "@PJL ENTER LANGUAGE=POSTSCRIPT\r\n"
	splitAt := 15 // Split in the middle of "LANGUAGE"

	var chunk1 bytes.Buffer
	chunk1.WriteString(uel + "@PJL\r\n")
	chunk1.WriteString(enterLine[:splitAt])

	n, err := p.Write(chunk1.Bytes())
	if err != nil {
		t.Fatalf("Write chunk1 failed: %v", err)
	}
	if n != chunk1.Len() {
		t.Fatalf("Write chunk1: got %d, want %d", n, chunk1.Len())
	}

	// No document yet
	if len(results) != 0 {
		t.Fatalf("expected 0 documents after chunk1, got %d",
			len(results))
	}

	// Second chunk: rest of ENTER line + PS content + end
	var chunk2 bytes.Buffer
	chunk2.WriteString(enterLine[splitAt:])
	chunk2.WriteString(psContent)
	chunk2.WriteByte(0x04)
	chunk2.WriteString(uel)

	n, err = p.Write(chunk2.Bytes())
	if err != nil {
		t.Fatalf("Write chunk2 failed: %v", err)
	}
	if n != chunk2.Len() {
		t.Fatalf("Write chunk2: got %d, want %d", n, chunk2.Len())
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}
	if results[0].format != DocFormatPostScript {
		t.Errorf("format = %v, want PostScript", results[0].format)
	}
	if !bytes.Equal(results[0].data, []byte(psContent)) {
		t.Errorf("document data mismatch:\ngot:  %q\nwant: %q",
			results[0].data, psContent)
	}
}
