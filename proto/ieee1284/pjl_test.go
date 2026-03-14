// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman(officialmdarman@gmail.com)
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

// TestPJLJobParams tests that @PJL JOB and @PJL SET commands
// are captured and passed to the DocumentHandler.
func TestPJLJobParams(t *testing.T) {
	ctx := newTestContext()
	var results []docResult
	p := NewPrinter(ctx, testHandler(&results))

	psContent := "%!PS-Adobe-3.0\nshowpage\n%%EOF\n"

	var job bytes.Buffer
	job.WriteString(uel + "@PJL\r\n")
	job.WriteString("@PJL JOB NAME = \"My Print Job\"\r\n")
	job.WriteString("@PJL SET SMOOTHING=ON\r\n")
	job.WriteString("@PJL SET ECONOMODE = OFF\r\n")
	job.WriteString("@PJL SET USERNAME = \"alice\"\r\n")
	job.WriteString("@PJL ENTER LANGUAGE=POSTSCRIPT\r\n")
	job.WriteString(psContent)
	job.WriteByte(0x04)
	job.WriteString(uel)

	writeInChunks(t, p, job.Bytes(), 512)

	if len(results) != 1 {
		t.Fatalf("expected 1 document, got %d", len(results))
	}

	params := results[0].params
	if params.JobName != "My Print Job" {
		t.Errorf("JobName = %q, want %q",
			params.JobName, "My Print Job")
	}
	if v := params.Variables["SMOOTHING"]; v != "ON" {
		t.Errorf("SMOOTHING = %q, want %q", v, "ON")
	}
	if v := params.Variables["ECONOMODE"]; v != "OFF" {
		t.Errorf("ECONOMODE = %q, want %q", v, "OFF")
	}
	if v := params.Variables["USERNAME"]; v != "alice" {
		t.Errorf("USERNAME = %q, want %q", v, "alice")
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

// TestParsePJLCommand tests parsing of PJL command lines,
// including commands with multiple spaces between words.
func TestParsePJLCommand(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect pjlCommand
	}{
		{
			"INFO ID single space",
			" INFO ID",
			pjlCommand{"INFO ID", ""},
		},
		{
			"INFO ID multiple spaces",
			" INFO   ID",
			pjlCommand{"INFO ID", ""},
		},
		{
			"INFO STATUS multiple spaces",
			" INFO   STATUS",
			pjlCommand{"INFO STATUS", ""},
		},
		{
			"ENTER LANGUAGE with args",
			" ENTER LANGUAGE=POSTSCRIPT",
			pjlCommand{"ENTER LANGUAGE", "=POSTSCRIPT"},
		},
		{
			"ENTER LANGUAGE multiple spaces",
			" ENTER   LANGUAGE=POSTSCRIPT",
			pjlCommand{"ENTER LANGUAGE", "=POSTSCRIPT"},
		},
		{
			"ECHO with text",
			" ECHO hello world",
			pjlCommand{"ECHO", "hello world"},
		},
		{
			"SET with value",
			" SET COPIES=2",
			pjlCommand{"SET", "COPIES=2"},
		},
		{
			"bare empty",
			"",
			pjlCommand{},
		},
		{
			"USTATUS DEVICE multiple spaces",
			" USTATUS   DEVICE  ON",
			pjlCommand{"USTATUS DEVICE", "ON"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePJLCommand(tt.input)
			if got.name != tt.expect.name ||
				got.args != tt.expect.args {
				t.Errorf("parsePJLCommand(%q):\n"+
					"  got  {name:%q, args:%q}\n"+
					"  want {name:%q, args:%q}",
					tt.input,
					got.name, got.args,
					tt.expect.name, tt.expect.args)
			}
		})
	}
}

// TestParsePJLKeyValue tests parsing of PJL SET key=value pairs.
func TestParsePJLKeyValue(t *testing.T) {
	tests := []struct {
		input     string
		wantKey   string
		wantValue string
	}{
		{"COPIES=1", "COPIES", "1"},
		{"SMOOTHING = ON", "SMOOTHING", "ON"},
		{"USERNAME = \"alice\"", "USERNAME", "alice"},
		{"ECONOMODE=OFF", "ECONOMODE", "OFF"},
		{"", "", ""},
		{"NOVALUE", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			k, v := parsePJLKeyValue(tt.input)
			if k != tt.wantKey || v != tt.wantValue {
				t.Errorf("parsePJLKeyValue(%q) = (%q, %q), want (%q, %q)",
					tt.input, k, v, tt.wantKey, tt.wantValue)
			}
		})
	}
}

// TestParsePJLQuotedValue tests extraction of named parameters.
func TestParsePJLQuotedValue(t *testing.T) {
	tests := []struct {
		args string
		key  string
		want string
	}{
		{`NAME = "test job"`, "NAME", "test job"},
		{`NAME = "job1" DISPLAY = "hello"`, "NAME", "job1"},
		{`NAME = "job1" DISPLAY = "hello"`, "DISPLAY", "hello"},
		{`NAME=unquoted`, "NAME", "unquoted"},
		{`FOO = "bar"`, "MISSING", ""},
		{``, "NAME", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key+"_"+tt.args, func(t *testing.T) {
			got := parsePJLQuotedValue(tt.args, tt.key)
			if got != tt.want {
				t.Errorf("parsePJLQuotedValue(%q, %q) = %q, want %q",
					tt.args, tt.key, got, tt.want)
			}
		})
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
