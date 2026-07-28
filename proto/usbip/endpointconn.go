// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// USBIP protocol connection

package usbip

import (
	"context"
	"net"
	"sync/atomic"
	"time"
)

// EndpointConn implements the [net.Conn] interface on a top of [Endpoint].
//
// This is suitable for usage with the [http.Server] to implement the
// IPP over USB protocol or for the similar tasks.
//
// Data, that was written into the endpoint from the USB side appears
// as data received from the connection and visa versa.
type EndpointConn struct {
	epp         EndpointPair       // Underlying endpoint
	closectx    context.Context    // Canceled by Close()
	closecancel context.CancelFunc // Cancels closectx
	listener    atomic.Pointer[    // Parent listener, nil if none
	EndpointListener]
}

// NewEndpointConn creates the [EndpointConn] on a top of the existent
// [EndpointPair].
func NewEndpointConn(epp EndpointPair) *EndpointConn {
	conn := &EndpointConn{epp: epp}

	ctx := context.Background()
	conn.closectx, conn.closecancel = context.WithCancel(ctx)

	return conn
}

// Close closes the connection and unblocks all pending Reads and Writes.
func (conn *EndpointConn) Close() error {
	// Unblock readers and writers
	conn.closecancel()

	// Return connection to the listener's pool
	listener := conn.listener.Swap(nil)
	if listener != nil {
		listener.endpoints <- conn.epp
	}

	return nil
}

// Read received data from the underlying [Endpoint].
func (conn *EndpointConn) Read(buf []byte) (int, error) {
	n, err := conn.epp.ReadContext(conn.closectx, buf)
	if err == context.Canceled {
		err = net.ErrClosed
	}
	return n, err
}

// Write sends data into the underlying [Endpoint].
func (conn *EndpointConn) Write(buf []byte) (int, error) {
	n, err := conn.epp.WriteContext(conn.closectx, buf)
	if err == context.Canceled {
		err = net.ErrClosed
	}
	return n, err
}

// LocalAddr returns the local network address, if known.
func (conn *EndpointConn) LocalAddr() net.Addr {
	return localhostAddr{}
}

// RemoteAddr returns the remote network address, if known.
func (conn *EndpointConn) RemoteAddr() net.Addr {
	return localhostAddr{}
}

// SetDeadline sets the read and write deadlines associated
// with the connection.
func (conn *EndpointConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
func (conn *EndpointConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
func (conn *EndpointConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// EndpointListener implements the [net.Listener] interface on
// a top of the group of [Endpoint]s.
//
// This is suitable for usage with the [http.Server] to implement the
// IPP over USB protocol or for the similar tasks.
//
// Its Accept methods waits until and returns some idle [Endpoint],
// wrapped into the [EndpointConn] type. The [EndpointConn.Close]
// returns connection into the pool.
type EndpointListener struct {
	endpoints chan EndpointPair // Queue of available endpoint pairs
	closechan chan struct{}     // Closed by EndpointListener.Close
}

// NewEndpointListener creates the new [EndpointListener] on a top
// of the group of existing [EndpointPair]s.
func NewEndpointListener(endpoints []EndpointPair) *EndpointListener {
	listener := &EndpointListener{
		endpoints: make(chan EndpointPair, len(endpoints)),
		closechan: make(chan struct{}),
	}

	for _, epp := range endpoints {
		listener.endpoints <- epp
	}

	return listener
}

// Accept waits for and returns the next connection to the listener.
func (listener *EndpointListener) Accept() (net.Conn, error) {
	// Don't even try to accept, if listener already closed
	select {
	case <-listener.closechan:
		return nil, net.ErrClosed
	default:
	}

	// Wait for connection or close
	select {
	case ep := <-listener.endpoints:
		conn := NewEndpointConn(ep)
		conn.listener.Store(listener)
		return conn, nil
	case <-listener.closechan:
		return nil, net.ErrClosed
	}
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (listener *EndpointListener) Close() error {
	close(listener.closechan)
	return nil
}

// Addr returns the listener's network address.
func (listener *EndpointListener) Addr() net.Addr {
	return localhostAddr{}
}

// localhostAddr implements the [net.Addr] interface for addresses
// used by EndpointConn.
//
// EndpointConn implements the [net.Conn] interface on a top of [Entpoint],
// which implies that it must be able to report its "local" and "remote"
// addresses. These addresses are actually meaningless for the USN
// endpoints, but interface requirements forces us to implement this.
//
// So we use just a dummy type that implements the required functionality.
type localhostAddr struct{}

func (localhostAddr) Network() string {
	return "tcp"
}

func (localhostAddr) String() string {
	return "localhost"
}
