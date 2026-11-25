// MFP - Multi-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Tests for LPD server

package lpd

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/abstract"
)

// capturedDoc is the document, captured by the captureBackend
type capturedDoc struct {
	params abstract.PrinterRequest
	body   string
}

// captureBackend implements abstract.Printer and collects
// received documents.
type captureBackend struct {
	lock sync.Mutex
	docs []capturedDoc
}

// PrintDocument implements abstract.Printer
func (cb *captureBackend) PrintDocument(params abstract.PrinterRequest,
	body io.Reader) error {

	data, err := io.ReadAll(body)

	cb.lock.Lock()
	cb.docs = append(cb.docs, capturedDoc{params, string(data)})
	cb.lock.Unlock()

	return err
}

// testClient is the LPD client side of the test session
type testClient struct {
	t    *testing.T
	conn net.Conn
}

// send sends raw bytes to the server
func (c *testClient) send(s string) {
	c.t.Helper()
	c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.WriteString(c.conn, s); err != nil {
		c.t.Fatalf("send %q: %s", s, err)
	}
}

// expect reads a single status byte and compares it with
// the expected value.
func (c *testClient) expect(status byte) {
	c.t.Helper()
	c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	var buf [1]byte
	if _, err := io.ReadFull(c.conn, buf[:]); err != nil {
		c.t.Fatalf("expect 0x%2.2x: %s", status, err)
	}
	if buf[0] != status {
		c.t.Fatalf("expect 0x%2.2x: got 0x%2.2x", status, buf[0])
	}
}

// file sends a control or data file subcommand, expecting
// all acknowledgements to be positive.
func (c *testClient) file(sub Subcommand, name, body string) {
	c.t.Helper()
	c.send(fmt.Sprintf("%c%d %s\n", byte(sub), len(body), name))
	c.expect(ack)
	c.send(body + "\x00")
	c.expect(ack)
}

// readAll reads everything until server closes connection
func (c *testClient) readAll() string {
	c.t.Helper()
	c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	data, _ := io.ReadAll(c.conn)
	return string(data)
}

// session runs the client script against the server over
// net.Pipe and waits until server finishes the connection.
func session(t *testing.T, srv *Server, script func(c *testClient)) {
	t.Helper()

	clnt, srvr := net.Pipe()
	done := make(chan struct{})

	go func() {
		srv.ServeConn(srvr)
		close(done)
	}()

	script(&testClient{t: t, conn: clnt})
	clnt.Close()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("server didn't finish the connection")
	}
}

// newTestServer creates a new Server with captureBackend
func newTestServer(opts ServerOptions) (*Server, *captureBackend) {
	cb := &captureBackend{}
	return NewServer(context.Background(), cb, opts), cb
}

// TestServerReceiveJob tests successful job reception
func TestServerReceiveJob(t *testing.T) {
	type testData struct {
		name   string
		script func(c *testClient)
		docs   []capturedDoc
	}

	const pjlJob = "\x1b%-12345X@PJL\r\n" +
		"@PJL JOB NAME=\"pjl-name\"\r\n" +
		"@PJL SET DUPLEX=ON\r\n" +
		"@PJL SET COLORMODE=MONO\r\n" +
		"@PJL ENTER LANGUAGE=POSTSCRIPT\r\n" +
		"%!PS-Adobe-3.0\nshowpage\n" +
		"\x1b%-12345X@PJL EOJ\r\n\x1b%-12345X"

	tests := []testData{
		{
			// Control file first, like CUPS lpd backend
			name: "control-first",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdControlFile, "cfA001host",
					"Hhost\nPgo-mfp\nJhello.txt\n"+
						"ldfA001host\nUdfA001host\n"+
						"Nhello.txt\n")
				c.file(SubcmdDataFile, "dfA001host",
					"Hello, printer!\n")
			},
			docs: []capturedDoc{
				{
					abstract.PrinterRequest{
						Format:  "text/plain",
						JobName: "hello.txt",
					},
					"Hello, printer!\n",
				},
			},
		},

		{
			// Data file first, like BSD lpr; PJL-wrapped
			// PostScript; PJL job name wins
			name: "data-first-pjl",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdDataFile, "dfA002host", pjlJob)
				c.file(SubcmdControlFile, "cfA002host",
					"Hhost\nProot\nJlpd-name\n"+
						"ldfA002host\n")
			},
			docs: []capturedDoc{
				{
					abstract.PrinterRequest{
						Format:    "application/postscript",
						JobName:   "pjl-name",
						Sides:     abstract.SidesTwoSidedLongEdge,
						ColorMode: abstract.ColorModeMono,
					},
					"%!PS-Adobe-3.0\nshowpage\n",
				},
			},
		},

		{
			// BSD lpr -#2: repeated print line, two copies;
			// no J line, so document names are used
			name: "copies-and-names",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdDataFile, "dfA003host",
					"%PDF-1.4 a")
				c.file(SubcmdDataFile, "dfB003host",
					"%!PS b")
				c.file(SubcmdControlFile, "cfA003host",
					"Hhost\nProot\n"+
						"ldfA003host\nldfA003host\n"+
						"Na.pdf\n"+
						"ldfB003host\nNb.ps\n")
			},
			docs: []capturedDoc{
				{
					abstract.PrinterRequest{
						Format:  "application/pdf",
						JobName: "a.pdf",
					},
					"%PDF-1.4 a",
				},
				{
					abstract.PrinterRequest{
						Format:  "application/pdf",
						JobName: "a.pdf",
					},
					"%PDF-1.4 a",
				},
				{
					abstract.PrinterRequest{
						Format:  "application/postscript",
						JobName: "b.ps",
					},
					"%!PS b",
				},
			},
		},

		{
			// Unreferenced data file is ignored
			name: "unreferenced-data-file",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdDataFile, "dfA004host", "used\n")
				c.file(SubcmdDataFile, "dfB004host", "unused\n")
				c.file(SubcmdControlFile, "cfA004host",
					"Jjob\nfdfA004host\n")
			},
			docs: []capturedDoc{
				{
					abstract.PrinterRequest{
						Format:  "text/plain",
						JobName: "job",
					},
					"used\n",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv, cb := newTestServer(ServerOptions{})
			defer srv.Close()

			session(t, srv, test.script)

			if len(cb.docs) != len(test.docs) {
				t.Fatalf("expected %d documents, got %d: %+v",
					len(test.docs), len(cb.docs), cb.docs)
			}

			for i := range test.docs {
				if cb.docs[i] != test.docs[i] {
					t.Errorf("document %d mismatch:\n"+
						"expected: %+v\n"+
						"present:  %+v",
						i, test.docs[i], cb.docs[i])
				}
			}
		})
	}
}

// TestServerReceiveJobErrors tests job reception errors.
// In all cases nothing must be printed.
func TestServerReceiveJobErrors(t *testing.T) {
	type testData struct {
		name   string
		opts   ServerOptions
		script func(c *testClient)
	}

	const control = "Hhost\nProot\nldfA001host\n"

	tests := []testData{
		{
			name: "missed-queue",
			script: func(c *testClient) {
				c.send("\x02\n")
				c.expect(nack)
			},
		},

		{
			name: "abort",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdControlFile, "cfA001host", control)
				c.file(SubcmdDataFile, "dfA001host", "data")
				c.send("\x01\n")
				c.expect(nack)
			},
		},

		{
			name: "no-control-file",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdDataFile, "dfA001host", "data")
			},
		},

		{
			name: "missed-data-file",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdControlFile, "cfA001host", control)
			},
		},

		{
			name: "bad-control-file",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdControlFile, "cfA001host", "l\n")
			},
		},

		{
			name: "bad-file-header",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.send("\x02abc cfA001host\n")
				c.expect(nack)
			},
		},

		{
			name: "zero-length-file",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.send("\x030 dfA001host\n")
				c.expect(nack)
			},
		},

		{
			name: "unknown-subcommand",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.send("\x0912 foo\n")
				c.expect(nack)
			},
		},

		{
			name: "empty-subcommand",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.send("\n")
				c.expect(nack)
			},
		},

		{
			name: "job-too-large",
			opts: ServerOptions{MaxJobSize: 30},
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.file(SubcmdControlFile, "cfA001host", control)
				c.send("\x0310 dfA001host\n")
				c.expect(nack)
			},
		},

		{
			name: "missed-terminator",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.send("\x024 cfA001host\n")
				c.expect(ack)
				c.send("abcdX")
				c.expect(nack)
			},
		},

		{
			name: "truncated-file",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				c.send("\x02100 cfA001host\n")
				c.expect(ack)
				c.send("abcd")
			},
		},

		{
			name: "too-many-data-files",
			script: func(c *testClient) {
				c.send("\x02lp\n")
				c.expect(ack)
				for i := 0; i < maxDataFiles; i++ {
					c.file(SubcmdDataFile,
						fmt.Sprintf("df%d", i), "x")
				}
				c.send("\x031 dfX\n")
				c.expect(nack)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv, cb := newTestServer(test.opts)
			defer srv.Close()

			session(t, srv, test.script)

			if len(cb.docs) != 0 {
				t.Errorf("expected no documents, got %+v",
					cb.docs)
			}
		})
	}
}

// TestServerCommands tests commands other than Receive-Job
func TestServerCommands(t *testing.T) {
	type testData struct {
		name     string
		send     string
		expected string
	}

	tests := []testData{
		{"print-waiting", "\x01lp\n", "\x00"},
		{"queue-state-short", "\x03lp\n", "lp: no entries\n"},
		{"queue-state-long", "\x04lp root 12\n", "lp: no entries\n"},
		{"remove-jobs", "\x05lp root 12\n", "\x00"},
		{"unknown", "\x09lp\n", "\x01"},
		{"empty", "\n", "\x01"},
		{"too-long", string(make([]byte, MaxLineLen+1)), "\x01"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv, _ := newTestServer(ServerOptions{})
			defer srv.Close()

			var reply string
			session(t, srv, func(c *testClient) {
				c.send(test.send)
				reply = c.readAll()
			})

			if reply != test.expected {
				t.Errorf("expected %q, present %q",
					test.expected, reply)
			}
		})
	}

	// Client connects and closes without sending anything
	srv, _ := newTestServer(ServerOptions{})
	defer srv.Close()
	session(t, srv, func(c *testClient) {})
}

// TestServerNilBackend tests that Server without backend
// drains documents without crashing.
func TestServerNilBackend(t *testing.T) {
	srv := NewServer(context.Background(), nil, ServerOptions{})
	defer srv.Close()

	session(t, srv, func(c *testClient) {
		c.send("\x02lp\n")
		c.expect(ack)
		c.file(SubcmdControlFile, "cfA001host", "ldfA001host\n")
		c.file(SubcmdDataFile, "dfA001host", "%!PS\nshowpage\n")
	})
}

// TestServerTimeout tests the idle timeout
func TestServerTimeout(t *testing.T) {
	srv, _ := newTestServer(ServerOptions{Timeout: 50 * time.Millisecond})
	defer srv.Close()

	start := time.Now()
	session(t, srv, func(c *testClient) {
		c.send("\x02lp\n")
		c.expect(ack)
		// Stay idle; server must close connection
		c.readAll()
	})

	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("timeout took too long: %s", elapsed)
	}
}

// TestServerServe tests Serve and Close over real TCP
func TestServerServe(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("can't listen: %s", err)
	}

	srv, cb := newTestServer(ServerOptions{})

	served := make(chan error, 1)
	go func() {
		served <- srv.Serve(l)
	}()

	// Print a job
	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %s", err)
	}

	c := &testClient{t: t, conn: conn}
	c.send("\x02lp\n")
	c.expect(ack)
	c.file(SubcmdControlFile, "cfA001host", "Jtcp\nldfA001host\n")
	c.file(SubcmdDataFile, "dfA001host", "tcp job\n")
	conn.Close()

	// Open an idle connection; Close must terminate it
	idle, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %s", err)
	}
	defer idle.Close()

	// Wait until the job is delivered
	deadline := time.Now().Add(5 * time.Second)
	for {
		cb.lock.Lock()
		n := len(cb.docs)
		cb.lock.Unlock()

		if n == 1 {
			break
		}

		if time.Now().After(deadline) {
			t.Fatalf("job not delivered")
		}
		time.Sleep(10 * time.Millisecond)
	}

	if cb.docs[0].params.JobName != "tcp" || cb.docs[0].body != "tcp job\n" {
		t.Errorf("unexpected document: %+v", cb.docs[0])
	}

	// Close the server
	closed := make(chan struct{})
	go func() {
		srv.Close()
		close(closed)
	}()

	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatalf("Close hangs")
	}

	if err := <-served; err != nil {
		t.Errorf("Serve returned %s", err)
	}

	// Serve and ServeConn after Close
	l2, err := net.Listen("tcp", "127.0.0.1:0")
	if err == nil {
		if err := srv.Serve(l2); err != nil {
			t.Errorf("Serve after Close returned %s", err)
		}
	}

	clnt, srvr := net.Pipe()
	srv.ServeConn(srvr)
	if _, err := clnt.Read(make([]byte, 1)); err == nil {
		t.Errorf("ServeConn after Close didn't close connection")
	}
}
