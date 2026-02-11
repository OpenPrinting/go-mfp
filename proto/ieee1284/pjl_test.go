// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// PJL command handling tests

package ieee1284

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

// TestPJLInfoID tests that @PJL INFO ID produces the correct
// response with the printer model name.
func TestPJLInfoID(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)
	p.SetModel("Virtual MFP")

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL INFO ID\r\n")

	p.Write(job.Bytes())

	// Read the response
	buf := make([]byte, 4096)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	got := string(buf[:n])
	want := "@PJL INFO ID\r\n\"Virtual MFP\"\r\n\x0c"
	if got != want {
		t.Errorf("INFO ID response:\ngot:  %q\nwant: %q",
			got, want)
	}
}

// TestPJLInfoStatus tests that @PJL INFO STATUS produces the
// correct status response.
func TestPJLInfoStatus(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL INFO STATUS\r\n")

	p.Write(job.Bytes())

	buf := make([]byte, 4096)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	got := string(buf[:n])
	if !strings.HasPrefix(got, "@PJL INFO STATUS\r\n") {
		t.Errorf("response should start with @PJL INFO STATUS")
	}
	if !strings.Contains(got, "CODE=10001") {
		t.Errorf("response should contain CODE=10001")
	}
	if !strings.Contains(got, "ONLINE=TRUE") {
		t.Errorf("response should contain ONLINE=TRUE")
	}
	if !strings.HasSuffix(got, "\x0c") {
		t.Errorf("response should end with form feed")
	}
}

// TestPJLEcho tests that @PJL ECHO returns the text verbatim.
func TestPJLEcho(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL ECHO hello world\r\n")

	p.Write(job.Bytes())

	buf := make([]byte, 4096)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	got := string(buf[:n])
	want := "@PJL ECHO hello world\r\n"
	if got != want {
		t.Errorf("ECHO response:\ngot:  %q\nwant: %q",
			got, want)
	}
}

// TestPJLNoResponse tests that PJL commands that don't produce
// responses leave the response buffer empty.
func TestPJLNoResponse(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL SET COPIES=1\r\n")
	job.WriteString("@PJL JOB NAME=\"test\"\r\n")
	job.WriteString("@PJL EOJ NAME=\"test\"\r\n")

	p.Write(job.Bytes())

	// Close so Read won't block
	p.Close()

	buf := make([]byte, 4096)
	n, err := p.Read(buf)
	if n != 0 {
		t.Errorf("expected no response data, got %d bytes: %q",
			n, buf[:n])
	}
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestPJLFullJob tests a complete PJL job with INFO ID query
// before the document, verifying both the PJL response and
// document extraction work together.
func TestPJLFullJob(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))
	p.SetModel("Test Printer")

	psContent := "%!PS-Adobe-3.0\nshowpage\n%%EOF\n"

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL INFO ID\r\n")
	job.WriteString("@PJL ENTER LANGUAGE=POSTSCRIPT\r\n")
	job.WriteString(psContent)
	job.WriteByte(0x04)
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL EOJ NAME=\"test\"\r\n")
	job.WriteString(uel)

	writeInChunks(t, p, job.Bytes(), 512)

	// Verify document was extracted
	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}
	if results[0].format != DocFormatPostScript {
		t.Errorf("format = %v, want PostScript",
			results[0].format)
	}

	// Verify PJL response is readable
	buf := make([]byte, 4096)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	got := string(buf[:n])
	want := "@PJL INFO ID\r\n\"Test Printer\"\r\n\x0c"
	if got != want {
		t.Errorf("INFO ID response:\ngot:  %q\nwant: %q",
			got, want)
	}
}

// TestReadBlocksUntilData verifies that Read() blocks when no
// data is available and unblocks when a response is queued.
func TestReadBlocksUntilData(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)

	readDone := make(chan string, 1)
	go func() {
		buf := make([]byte, 4096)
		n, _ := p.Read(buf)
		readDone <- string(buf[:n])
	}()

	// Verify Read is blocking
	select {
	case <-readDone:
		t.Fatal("Read returned before data was available")
	case <-time.After(50 * time.Millisecond):
		// Good — Read is blocking
	}

	// Now send a PJL command that produces a response
	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL ECHO test\r\n")
	p.Write(job.Bytes())

	// Read should now return
	select {
	case got := <-readDone:
		want := "@PJL ECHO test\r\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("Read did not unblock after data was queued")
	}
}

// TestReadAfterClose verifies that Read() returns io.EOF
// after the printer is closed.
func TestReadAfterClose(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)

	p.Close()

	buf := make([]byte, 4096)
	n, err := p.Read(buf)
	if n != 0 {
		t.Errorf("expected 0 bytes, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestReadChunked verifies that large responses can be read
// in multiple small chunks.
func TestReadChunked(t *testing.T) {
	ctx := newTestContext()
	p := NewPrinter(ctx, nil)
	p.SetModel("Virtual MFP")

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL INFO ID\r\n")
	job.WriteString("@PJL INFO STATUS\r\n")

	p.Write(job.Bytes())

	// Read in small chunks
	var got []byte
	buf := make([]byte, 10)
	for i := 0; i < 20; i++ {
		n, err := p.Read(buf)
		if n > 0 {
			got = append(got, buf[:n]...)
		}
		if err != nil {
			break
		}

		// If we got two form feeds, we have both responses
		if bytes.Count(got, []byte("\x0c")) >= 2 {
			break
		}
	}

	result := string(got)
	if !strings.Contains(result, "@PJL INFO ID\r\n") {
		t.Error("missing INFO ID response")
	}
	if !strings.Contains(result, "@PJL INFO STATUS\r\n") {
		t.Error("missing INFO STATUS response")
	}
}
