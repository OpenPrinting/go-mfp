// MFP - Multi-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// LPD server

package lpd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/ieee1284"
)

// Default values for ServerOptions:
const (
	DefaultMaxJobSize = 256 << 20        // 256 MiB per job
	DefaultTimeout    = 60 * time.Second // Idle timeout
)

// maxDataFiles limits the number of data files per job.
// RFC 1179 naming (dfA...dfZ, dfa...dfz) allows 52; cups-lpd
// allows 100.
const maxDataFiles = 100

// ServerOptions contains optional parameters for the Server.
// Zero value of any field means "use default".
type ServerOptions struct {
	// MaxJobSize limits the total size of all files (control
	// and data) received within a single job. Larger jobs
	// are rejected.
	MaxJobSize int64

	// Timeout is the connection idle timeout. If the client
	// sends nothing during this time, connection is closed.
	Timeout time.Duration
}

// Server implements the LPD server.
//
// Received data files are passed through the [ieee1284.Printer]
// stream parser, which detects the document format and extracts
// job parameters embedded as PJL commands, and then delivered to
// the [abstract.Printer] backend.
type Server struct {
	ctx       context.Context       // Logging context
	backend   abstract.Printer      // Print backend; may be nil
	opts      ServerOptions         // Server options
	lock      sync.Mutex            // Protects the following fields
	listeners map[net.Listener]bool // Active listeners
	conns     map[net.Conn]bool     // Active connections
	closed    bool                  // Server is closed
	wg        sync.WaitGroup        // Waits for connection handlers
}

// NewServer creates a new LPD server.
//
// If backend is nil, received documents are drained and discarded.
func NewServer(ctx context.Context, backend abstract.Printer,
	opts ServerOptions) *Server {

	if opts.MaxJobSize <= 0 {
		opts.MaxJobSize = DefaultMaxJobSize
	}

	if opts.Timeout <= 0 {
		opts.Timeout = DefaultTimeout
	}

	return &Server{
		ctx:       log.WithPrefix(ctx, "lpd"),
		backend:   backend,
		opts:      opts,
		listeners: make(map[net.Listener]bool),
		conns:     make(map[net.Conn]bool),
	}
}

// Serve accepts incoming connections on the listener l and
// serves them, each in its own goroutine.
//
// Serve always closes l before return. After [Server.Close]
// it returns nil; otherwise it returns the Accept error.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	srv.lock.Lock()
	if srv.closed {
		srv.lock.Unlock()
		return nil
	}
	srv.listeners[l] = true
	srv.lock.Unlock()

	defer func() {
		srv.lock.Lock()
		delete(srv.listeners, l)
		srv.lock.Unlock()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			srv.lock.Lock()
			closed := srv.closed
			srv.lock.Unlock()

			if closed {
				return nil
			}

			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}

			return err
		}

		if !srv.trackConn(conn) {
			conn.Close()
			return nil
		}

		go func() {
			defer srv.untrackConn(conn)
			srv.handleConn(conn)
		}()
	}
}

// ServeConn serves a single connection. It returns when the
// connection is done, and always closes conn.
func (srv *Server) ServeConn(conn net.Conn) {
	if !srv.trackConn(conn) {
		conn.Close()
		return
	}

	defer srv.untrackConn(conn)
	srv.handleConn(conn)
}

// Close closes the Server: closes all active listeners and
// connections and waits until all connection handlers finish.
func (srv *Server) Close() error {
	srv.lock.Lock()
	srv.closed = true

	for l := range srv.listeners {
		l.Close()
	}

	for conn := range srv.conns {
		conn.Close()
	}
	srv.lock.Unlock()

	srv.wg.Wait()
	return nil
}

// trackConn registers an active connection.
// It returns false if Server is already closed.
func (srv *Server) trackConn(conn net.Conn) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if srv.closed {
		return false
	}

	srv.conns[conn] = true
	srv.wg.Add(1)
	return true
}

// untrackConn unregisters an active connection.
func (srv *Server) untrackConn(conn net.Conn) {
	srv.lock.Lock()
	delete(srv.conns, conn)
	srv.lock.Unlock()

	srv.wg.Done()
}

// handleConn handles a single LPD connection.
func (srv *Server) handleConn(netconn net.Conn) {
	defer netconn.Close()

	ctx := log.WithPrefix(srv.ctx,
		fmt.Sprintf("lpd %s", netconn.RemoteAddr()))

	conn := &idleConn{Conn: netconn, timeout: srv.opts.Timeout}
	r := bufio.NewReader(conn)

	line, err := readLine(r, MaxLineLen)
	if err == io.EOF {
		log.Debug(ctx, "connection closed without command")
		return
	}

	var code byte
	var operands []string
	if err == nil {
		code, operands, err = parseLine(line)
	}

	if err != nil {
		log.Error(ctx, "%s", err)
		conn.reply(nack)
		return
	}

	cmd := Command(code)
	queue := ""
	if len(operands) > 0 {
		queue = operands[0]
	}

	log.Debug(ctx, "%s: queue=%q operands=%q", cmd, queue, operands)

	switch cmd {
	case CmdPrintWaiting:
		conn.reply(ack)

	case CmdReceiveJob:
		if queue == "" {
			log.Error(ctx, "%s: missed queue name", cmd)
			conn.reply(nack)
			return
		}

		srv.receiveJob(ctx, r, conn, queue)

	case CmdQueueStateShort, CmdQueueStateLong:
		// No status byte for these commands, just the text
		fmt.Fprintf(conn, "%s: no entries\n", queue)

	case CmdRemoveJobs:
		// Jobs are delivered immediately, so there is
		// nothing to remove.
		conn.reply(ack)

	default:
		log.Error(ctx, "%s: command not supported", cmd)
		conn.reply(nack)
	}
}

// receiveJob handles the Receive-Job command and its
// subcommands.
func (srv *Server) receiveJob(ctx context.Context,
	r *bufio.Reader, conn *idleConn, queue string) {

	if !conn.reply(ack) {
		return
	}

	var control []byte          // Control file(s), concatenated
	data := map[string][]byte{} // Data files, by name
	remaining := srv.opts.MaxJobSize

	for {
		line, err := readLine(r, MaxLineLen)
		if err == io.EOF {
			break // Client closed connection: job complete
		}

		var code byte
		var operands []string
		if err == nil {
			code, operands, err = parseLine(line)
		}

		if err != nil {
			log.Error(ctx, "receive-job: %s", err)
			conn.reply(nack)
			return
		}

		sub := Subcommand(code)
		switch sub {
		case SubcmdAbortJob:
			log.Info(ctx, "receive-job: aborted by client")
			conn.reply(nack) // Same as cups-lpd
			return

		case SubcmdControlFile, SubcmdDataFile:
		default:
			log.Error(ctx, "receive-job: %s: subcommand not supported",
				sub)
			conn.reply(nack)
			return
		}

		count, name, err := parseFileHeader(operands, remaining)
		if err == nil && sub == SubcmdDataFile &&
			len(data) >= maxDataFiles {
			err = errors.New("lpd: too many data files")
		}

		if err != nil {
			log.Error(ctx, "receive-job: %s: %s", sub, err)
			conn.reply(nack)
			return
		}

		if !conn.reply(ack) {
			return
		}

		body, err := readFileBody(r, count)
		if err != nil {
			log.Error(ctx, "receive-job: %s %q: %s", sub, name, err)
			conn.reply(nack)
			return
		}

		log.Debug(ctx, "receive-job: %s %q: %d bytes", sub, name, count)
		remaining -= count

		if sub == SubcmdControlFile {
			// Some clients (OS/2) send multiple control
			// files per connection; cups-lpd appends them.
			control = append(control, body...)
		} else {
			if _, dup := data[name]; dup {
				log.Warning(ctx,
					"receive-job: data file %q received twice",
					name)
			}
			data[name] = body
		}

		if !conn.reply(ack) {
			return
		}
	}

	srv.printJob(ctx, queue, control, data)
}

// printJob processes the completely received job.
func (srv *Server) printJob(ctx context.Context, queue string,
	control []byte, data map[string][]byte) {

	if control == nil {
		log.Error(ctx, "job discarded: no control file")
		return
	}

	cf, err := ParseControlFile(control)
	if err != nil {
		log.Error(ctx, "job discarded: control file: %s", err)
		return
	}

	log.Info(ctx, "job received: queue=%q user=%q host=%q title=%q",
		queue, cf.User, cf.Host, cf.Title())

	// Check that all referenced data files have been received,
	// the same way as cups-lpd does.
	for _, name := range cf.DataFiles() {
		if _, found := data[name]; !found {
			log.Error(ctx, "job discarded: data file %q missed",
				name)
			return
		}
	}

	if len(cf.Documents) == 0 {
		log.Warning(ctx, "job has no documents to print")
	}

	for i, doc := range cf.Documents {
		title := doc.Name
		if title == "" {
			title = cf.Title()
		}

		log.Debug(ctx, "document %d/%d: format=%q file=%q name=%q",
			i+1, len(cf.Documents), doc.Format, doc.DataFile,
			doc.Name)

		jb := &jobBackend{ctx: ctx, backend: srv.backend, title: title}
		p := ieee1284.NewPrinter(ctx, jb)
		p.Write(data[doc.DataFile])
		p.Flush()
		p.Close()

		if jb.count == 0 {
			log.Warning(ctx, "document %d/%d: no printable data",
				i+1, len(cf.Documents))
		}
	}
}

// jobBackend wraps the Server's backend for a single LPD
// document. It supplies the job name from the control file,
// if PJL didn't provide one, and drains the document if the
// Server has no backend.
type jobBackend struct {
	ctx     context.Context  // Logging context
	backend abstract.Printer // Real backend; may be nil
	title   string           // Fallback job name
	count   int              // Count of delivered documents
}

// PrintDocument implements the [abstract.Printer] interface.
func (jb *jobBackend) PrintDocument(params abstract.PrinterRequest,
	body io.Reader) error {

	jb.count++

	if params.JobName == "" {
		params.JobName = jb.title
	}

	log.Debug(jb.ctx, "print: format=%q name=%q copies=%d",
		params.Format, params.JobName, params.Copies)

	if jb.backend == nil {
		n, err := io.Copy(io.Discard, body)
		log.Debug(jb.ctx, "print: %d bytes discarded (no backend)", n)
		return err
	}

	return jb.backend.PrintDocument(params, body)
}

// idleConn wraps net.Conn and refreshes the deadline before
// each Read and Write, so the timeout applies to idle periods
// rather than to the whole connection.
type idleConn struct {
	net.Conn
	timeout time.Duration
}

// Read implements the io.Reader interface.
func (conn *idleConn) Read(buf []byte) (int, error) {
	conn.SetDeadline(time.Now().Add(conn.timeout))
	return conn.Conn.Read(buf)
}

// Write implements the io.Writer interface.
func (conn *idleConn) Write(buf []byte) (int, error) {
	conn.SetDeadline(time.Now().Add(conn.timeout))
	return conn.Conn.Write(buf)
}

// reply sends the single-byte acknowledgement.
// It returns false on error.
func (conn *idleConn) reply(status byte) bool {
	_, err := conn.Write([]byte{status})
	return err == nil
}
