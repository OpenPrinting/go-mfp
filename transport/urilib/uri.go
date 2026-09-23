// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URI type and methods

package urilib

import (
	"net"
	"net/netip"
	"strconv"
	"strings"
)

// URI represents an URI string.
type URI string

// New returns a new URI, based on a given string.
func New(s string) URI {
	return URI(s)
}

// Valid reports if URI is valid.
func (u URI) Valid() bool {
	return lookup(u).Valid()
}

// Err returns parse error for invalid URI or nil if URI is valid.
func (u URI) Err() error {
	return lookup(u).err
}

// Scheme returns URI scheme (e.g., "http" or "ipp").
// For invalid URIs it returns "".
func (u URI) Scheme() string {
	cached := lookup(u)
	if !cached.Valid() {
		return ""
	}
	return cached.parsed.Scheme
}

// Hostname returns Hostname of the URI.
//
// If the result is enclosed in square brackets, as literal IPv6 addresses are,
// the square brackets are removed from the result.
//
// For invalid URIs it returns "".
func (u URI) Hostname() string {
	cached := lookup(u)
	if !cached.Valid() {
		return ""
	}

	return cached.parsed.Hostname()
}

// Portname returns Portname of the URI.
func (u URI) Portname() string {
	cached := lookup(u)
	if !cached.Valid() {
		return ""
	}

	return cached.parsed.Port()
}

// PortNum returns port number, defined by the URI.
// It may return 0 in the following cases:
//   - URI scheme doesn't support port
//   - Port is present, but it is not numeric
//   - URI is invalid
func (u URI) PortNum() uint16 {
	p := u.Portname()
	if p == "" {
		return u.DefaultPort()
	}

	n, err := strconv.ParseUint(p, 10, 16)
	if err != nil {
		return 0
	}

	return uint16(n)
}

// DefaultPort returns the default port, based on the URI scheme.
// If URI is not valid or scheme doesn't imply the port, it
// returns 0.
//
// The following schemes are supported here:
//   - http      - 80
//   - https     - 443
//   - ipp, ipps - 631
//   - lpd       - 515
//   - socket    - 9100
func (u URI) DefaultPort() uint16 {
	return lookup(u).DefaultPort()
}

// Canonical returns the canonical form of the URI:
//   - Scheme, Host and Port converted to lower case
//   - Literal IP addressed are canonicalized
//   - Leading zeroes from port are trimmed
//   - Port is stripped, if matches the default
//   - Path is normalized, "." and ".." and repeated "/" are processed
//   - Empty path replaced with "/". Otherwise, trailing slash
//     is preserved
//   - unix URIs converted to the short form (unix:/path)
//
// Invalid URIs returned unmodified
func (u URI) Canonical() URI {
	return lookup(u).canonical
}

// WithHostname replaces Hostname part of the URI.
//
// Host MUST be valid hostname. IPv6 literals MUST NOT
// be enclosed into square braces.
//
// Invalid or non-TCP URIs returned unchanged.
func (u URI) WithHostname(newHostname string) URI {
	cached := lookup(u)
	if !cached.IsTCP() {
		return u
	}

	parsed := *cached.parsed
	_, port, err := net.SplitHostPort(parsed.Host)
	if err != nil {
		if strings.IndexByte(newHostname, ':') >= 0 {
			parsed.Host = "[" + newHostname + "]"
		} else {
			parsed.Host = newHostname
		}
	} else {
		parsed.Host = net.JoinHostPort(newHostname, port)
	}

	return New(parsed.String())
}

// WithPortNum replaces Port part of the URI.
//
// If newPortNum <= 0, port will be removed. Otherwise,
// it will be set, as specified.
//
// Invalid or non-TCP URIs returned unchanged.
func (u URI) WithPortNum(newPortNum uint16) URI {
	cached := lookup(u)
	if !cached.IsTCP() {
		return u
	}

	parsed := *cached.parsed
	host := u.Hostname()

	port := strconv.Itoa(int(newPortNum))
	parsed.Host = net.JoinHostPort(host, port)

	return New(parsed.String())
}

// WithoutPort removes Port part of the URI.
//
// Invalid or non-TCP URIs returned unchanged.
func (u URI) WithoutPort() URI {
	cached := lookup(u)
	if !cached.IsTCP() {
		return u
	}

	parsed := *cached.parsed
	host := u.Hostname()

	if strings.IndexByte(host, ':') >= 0 {
		parsed.Host = "[" + host + "]"
	} else {
		parsed.Host = host
	}

	return New(parsed.String())
}

// IPAddress returns URI's IP address.
// It doesn't do any name resolution, URI must contain literal IP address.
//
// If URI is invalid, or not network URI or its hostname is not literal,
// it returns a zero netip.AddrPort
func (u URI) IPAddress() netip.AddrPort {
	cached := lookup(u)
	if !cached.IsTCP() {
		return netip.AddrPort{}
	}

	if addrport, err := netip.ParseAddrPort(cached.parsed.Host); err == nil {
		return addrport
	}

	addr, err := netip.ParseAddr(cached.parsed.Hostname())
	port := u.PortNum()

	if err != nil || port == 0 {
		return netip.AddrPort{}
	}

	return netip.AddrPortFrom(addr, uint16(port))
}

// IsTCP reports if URI uses TCP-based transport.
func (u URI) IsTCP() bool {
	return lookup(u).IsTCP()
}

// IsIP4 returns true, if URI has literal IP address and this address is IPv4.
func (u URI) IsIP4() bool {
	return u.IPAddress().Addr().Is4()
}

// IsIP6 returns true, if URI has literal IP address and this address is IPv6.
func (u URI) IsIP6() bool {
	return u.IPAddress().Addr().Is6()
}

// IsHTTP reports if URI uses HTTP or HTTPS - based transport.
func (u URI) IsHTTP() bool {
	return lookup(u).IsHTTP()
}

// IsTLS reports if URI uses TLS encryption.
func (u URI) IsTLS() bool {
	return lookup(u).IsTLS()
}
