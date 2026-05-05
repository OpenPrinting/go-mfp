// MFP - Miulti-Function Printers and scanners toolkit
// Protocol tracer
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Trace writer

package trace

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/transport"
)

// Writer writes a protocol trace
type Writer struct {
	ctx      context.Context // Logging context
	name     string          // file name
	fp       *os.File        // Underlying file
	tar      *tar.Writer     // TAR writer
	lock     sync.Mutex      // Access lock
	err      error           // First error
	donewait sync.WaitGroup  // Wait for async activities
}

// NewWriter creates a new trace writer
func NewWriter(ctx context.Context, name string) (*Writer, error) {
	nameLog := name + ".log"
	nameTar := name + ".tar"

	// Create name.log
	os.Remove(nameLog)
	backend := log.NewFileBackend(nameLog, 0, 0)
	log.CtxLogger(ctx).Attach(log.LevelTrace, backend)

	// Create name.tar
	const flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	fp, err := os.OpenFile(nameTar, flags, 0644)
	if err != nil {
		return nil, err
	}

	writer := &Writer{
		ctx:  ctx,
		name: name,
		fp:   fp,
		tar:  tar.NewWriter(fp),
	}

	return writer, nil
}

// Close closes the Writer
func (writer *Writer) Close() {
	writer.donewait.Wait()

	writer.lock.Lock()
	defer writer.lock.Unlock()

	err := writer.tar.Close()
	if err != nil {
		writer.setError(err)
	}
	err = writer.fp.Close()
	if err != nil {
		writer.setError(err)
	}
}

// OnRequest needs to be called by protocol being traced
// to write the request message.
func (writer *Writer) OnRequest(query *transport.ServerQuery,
	msg Message, body io.Reader) {

	name := fmt.Sprintf("%8.8d/00-%s", query.ID(), msg.Name())

	writer.Send(name+".http", query.DumpRequest())
	writer.Send(name+"."+msg.Ext(), msg.MarshalTrace())

	writer.donewait.Add(1)
	go func() {
		data, _ := io.ReadAll(body)

		if len(data) != 0 {
			name := fmt.Sprintf("%8.8d/01-odata.%s",
				query.ID(), magic(data))

			writer.Send(name, data)
		}

		writer.donewait.Done()
	}()
}

// OnResponse needs to be called by protocol being traced
// to write the response message.
func (writer *Writer) OnResponse(query *transport.ServerQuery,
	msg Message, body io.Reader) {

	name := fmt.Sprintf("%8.8d/02-%s.%s", query.ID(), msg.Name(), msg.Ext())
	writer.Send(name, msg.MarshalTrace())

	writer.donewait.Add(1)
	go func() {
		data, _ := io.ReadAll(body)

		if len(data) != 0 {
			name := fmt.Sprintf("%8.8d/03-rdata.%s",
				query.ID(), magic(data))

			writer.Send(name, data)
		}

		writer.donewait.Done()
	}()
}

// Send writes a new record (a file) into the writer archive.
func (writer *Writer) Send(name string, data []byte) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	log.Debug(writer.ctx, "%s: %d bytes saved", name, len(data))

	hdr := tar.Header{
		Typeflag: tar.TypeReg,
		Name:     name,
		Size:     int64(len(data)),
		Mode:     0644,
		ModTime:  time.Now(),
		Devmajor: 1,
		Devminor: 1,
	}

	if writer.err == nil {
		err := writer.tar.WriteHeader(&hdr)
		if err != nil {
			writer.setError(err)
		}
	}

	if writer.err == nil {
		_, err := writer.tar.Write(data)
		if err != nil {
			writer.setError(err)
		}
	}

	if writer.err == nil {
		err := writer.tar.Flush()
		if err != nil {
			writer.setError(err)
		}
	}
}

// setError sets writer.err, when error occurs for the first time.
// When it happens, the event is logged.
//
// This function must be called under writer.lock
func (writer *Writer) setError(err error) {
	if writer.err == nil {
		writer.err = err
		log.Error(writer.ctx, "%s: %s", writer.name, err)
	}
}
