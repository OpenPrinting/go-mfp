// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// TCP listeners

package transport

import (
	"net"
	"net/netip"
	"syscall"
)

// ListenTCPRange attempts to allocate TCP port in range,
// and returns [net.Listener] on success.
func ListenTCPRange(addr netip.Addr, min, max uint16) (net.Listener, error) {
	addr = addr.Unmap()
	tcpaddr := net.TCPAddr{}

	if min == 0 {
		min = 1 // TCP port 0 doesn't exist
	}

	if min >= max {
		return nil, syscall.EADDRINUSE
	}

	if addr.IsValid() {
		tcpaddr.IP = net.IP(addr.AsSlice())
	}

	if addr.Is6() {
		tcpaddr.Zone = addr.Zone()
	}

	var err error
	for p := min; ; p++ {
		tcpaddr.Port = int(p)

		ln, err2 := net.ListenTCP("tcp", &tcpaddr)
		switch {
		case ln != nil:
			return ln, nil
		case err == nil:
			// Save the first occurred error
			err = err2
		}

		if p == max {
			break
		}
	}

	return nil, err
}

// ListenTCPRangeN allocates N consecutive TCP ports within the range from min
// to max and, on success, returns a slice of [net.Listener]s.
func ListenTCPRangeN(addr netip.Addr, N int, min, max uint16) (
	[]net.Listener, error) {

	// Handle trivial cases
	if min == 0 {
		min = 1 // TCP port 0 doesn't exist
	}

	switch {
	case N <= 0:
		return []net.Listener{}, nil
	case min > max:
		return nil, syscall.EADDRINUSE
	case N > int(max-min)+1:
		return nil, syscall.EADDRINUSE
	}

	// Try to allocate N consecutive ports
	ports := make([]net.Listener, 0, N)
	var err error

	for beg := min; ; beg++ {
		end := beg + uint16(N-1)
		for p := beg; p <= end; p++ {
			ln, err2 := ListenTCPRange(addr, p, p)
			if ln != nil {
				ports = append(ports, ln)
				if len(ports) == N {
					return ports, nil
				}
			} else {
				// Save the first occurred error
				if err == nil {
					err = err2
				}

				// Close already allocated ports
				for _, ln := range ports {
					ln.Close()
				}

				ports = ports[:0]
				break
			}
		}

		if beg == max-uint16(N-1) {
			break
		}
	}

	return nil, err
}
