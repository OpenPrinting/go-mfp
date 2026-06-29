// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// UDP connection

package wsdd

import (
	"fmt"
	"net"
	"net/netip"
	"sync/atomic"

	"github.com/OpenPrinting/go-mfp/internal/netstate"
)

// uconn wraps net.UDPConn, binds it to the particular network interface
// and prepares it to be used for sending unicast and multicast packets.
type uconn struct {
	*net.UDPConn               // Underlying connection
	local        netstate.Addr // Local address
	closed       atomic.Bool   // Connection is closed
}

// newUconn creates a new unicast connection
func newUconn(local netstate.Addr, port uint16) (*uconn, error) {
	// Address must be unicast
	if local.Addr().IsMulticast() {
		err := fmt.Errorf("%s not unicast", local.Addr())
		return nil, err
	}

	// Prepare net.UDPAddr structure
	addr := &net.UDPAddr{
		IP:   net.IP(local.Addr().AsSlice()),
		Port: int(port),
		Zone: local.Addr().Zone(),
	}

	// Open UDP connection.
	network := "udp4"
	if local.Addr().Is6() {
		network = "udp6"
	}

	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		return nil, err
	}

	// Fill the uconn structure
	uc := &uconn{
		UDPConn: conn,
		local:   local,
	}

	// Do system-specific setup
	if uc.Is4() {
		err = uc.sysSetSockOptIP4()
	} else {
		err = uc.sysSetSockOptIP6()
	}

	if err != nil {
		uc.Close()
		return nil, err
	}

	return uc, nil
}

// Close closes the connection
func (uc *uconn) Close() {
	uc.closed.Store(true)
	uc.UDPConn.Close()
}

// LocalAddrPort returns connection's local address and port
func (uc *uconn) LocalAddrPort() netip.AddrPort {
	addr := uc.UDPConn.LocalAddr().(*net.UDPAddr)
	return addr.AddrPort()
}

// IsClosed reports if connection is closed
func (uc *uconn) IsClosed() bool {
	return uc.closed.Load()
}

// Is4 reports if connection uses IPv4 address family
func (uc *uconn) Is4() bool {
	return uc.local.Addr().Is4()
}

// Is6 reports if connection uses IPv6 address family
func (uc *uconn) Is6() bool {
	return uc.local.Addr().Is6()
}

// RecvFrom receives a packet from the UDP connection
func (uc *uconn) RecvFrom(b []byte) (n int, from netip.AddrPort, err error) {
	n, from, err = uc.UDPConn.ReadFromUDPAddrPort(b)
	if err != nil {
		return
	}

	from = netip.AddrPortFrom(from.Addr().Unmap(), from.Port())

	return
}

// control invokes f on the underlying connection's
// file descriptor.
func (uc *uconn) control(f func(fd int) error) error {
	rawconn, err := uc.SyscallConn()
	if err != nil {
		return err
	}

	var err2 error
	err = rawconn.Control(func(fd uintptr) {
		err2 = f(int(fd))
	})

	if err != nil {
		return err
	}

	return err2
}
