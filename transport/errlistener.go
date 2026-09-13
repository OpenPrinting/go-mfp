// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// net.Listener that always return an error

package transport

import "net"

// ErrListener implements [net.Listener] interface.
// Its Accept method always returns the error.
type ErrListener struct {
	Err error
}

var _ = net.Listener(ErrListener{})

// Accept waits for and returns the next connection to the listener.
// ErrListener always returns an error instead.
func (l ErrListener) Accept() (net.Conn, error) {
	err := l.Err
	if err == nil {
		err = net.ErrClosed
	}
	return nil, err
}

// Close closes the listener.
func (l ErrListener) Close() error {
	return nil
}

// Addr returns the listener's network address.
func (l ErrListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4zero, Port: 0}
}
