// MFP - Miulti-Function Printers and scanners toolkit
// The "mfp-proxy" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// net.Listener wrapper

package command

import (
	"context"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
)

// listener wraps net.Listener for a reason of fine-tuning incoming
// TCP connections.
type listener struct {
	net.Listener                 // Underlying net.Listener
	ctx          context.Context // For logging and shutdown
	cancel       func()          // ctx cancel function
	closed       atomic.Bool     // listener.Close called
	closeWait    sync.WaitGroup  // Wait for listener.Close completion
}

// newListener creates a new listener
func newListener(ctx context.Context, port int) (net.Listener, error) {
	// Setup network and address
	network := "tcp"
	addr := ":" + strconv.Itoa(port)

	// Create net.Listener
	nl, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(ctx)

	// Wrap into Listener
	l := &listener{
		Listener: nl,
		ctx:      ctx,
		cancel:   cancel,
	}

	l.closeWait.Add(1)
	go l.kill()

	return l, nil
}

// kill closes underlying net.Listener when listener.ctx is canceled.
func (l *listener) kill() {
	<-l.ctx.Done()

	l.closed.Store(true)
	l.Listener.Close()

	l.closeWait.Done()
}

// Accept new connection.
func (l *listener) Accept() (net.Conn, error) {
	if l.closed.Load() {
		return nil, ErrShutdown
	}

	l.closeWait.Add(1)
	defer l.closeWait.Done()

	for {
		// Accept new connection
		conn, err := l.Listener.Accept()

		if l.closed.Load() {
			if conn != nil {
				conn.Close()
			}
			return nil, ErrShutdown
		}

		if err != nil {
			return nil, err
		}

		// Obtain underlying net.TCPConn
		tcpconn, ok := conn.(*net.TCPConn)
		if !ok {
			// Should never happen, actually
			conn.Close()
			continue
		}

		// Setup TCP parameters
		tcpconn.SetKeepAlive(true)
		tcpconn.SetKeepAlivePeriod(20 * time.Second)

		// Issue log message
		log.Debug(l.ctx, "PROXY %s: new connection from %s",
			tcpconn.LocalAddr(), tcpconn.RemoteAddr())

		return tcpconn, nil
	}
}

// Close closes the listener.
func (l *listener) Close() error {
	l.cancel()
	l.closeWait.Wait()
	return nil
}
