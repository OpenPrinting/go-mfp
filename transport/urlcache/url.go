// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package urlcache

import (
	"net"
	"net/netip"
	"strconv"
	"strings"
)

// URL represents an URL string.
type URL string

// New returns a new URL, based on a given string.
func New(s string) URL {
	return URL(s)
}

// Valid reports if URL is valid.
func (u URL) Valid() bool {
	return lookup(u).Valid()
}

// Err returns parse error for invalid URL or nil if URL is valid.
func (u URL) Err() error {
	return lookup(u).err
}

// Scheme returns URL scheme (e.g., "http" or "ipp").
// For invalid URLs it returns "".
func (u URL) Scheme() string {
	cached := lookup(u)
	if !cached.Valid() {
		return ""
	}
	return cached.parsed.Scheme
}

// Hostname returns Hostname of the URL.
//
// If the result is enclosed in square brackets, as literal IPv6 addresses are,
// the square brackets are removed from the result.
//
// For invalid URLs it returns "".
func (u URL) Hostname() string {
	cached := lookup(u)
	if !cached.Valid() {
		return ""
	}

	return cached.parsed.Hostname()
}

// Portname returns Portname of the URL.
func (u URL) Portname() string {
	cached := lookup(u)
	if !cached.Valid() {
		return ""
	}

	return cached.parsed.Port()
}

// PortNum returns port number, defined by the URL.
// It may return 0 in the following cases:
//   - URL scheme doesn't support port
//   - Port is present, but it is not numeric
//   - URL is invalid
func (u URL) PortNum() uint16 {
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

// DefaultPort returns the default port, based on the URL scheme.
// If URL is not valid or scheme doesn't imply the port, it
// returns 0.
//
// The following schemes are supported here:
//   - http      - 80
//   - https     - 443
//   - ipp, ipps - 631
//   - lpd       - 515
//   - socket    - 9100
func (u URL) DefaultPort() uint16 {
	return lookup(u).DefaultPort()
}

// Canonical returns the canonical form of the URL:
//   - Scheme, Host and Port converted to lower case
//   - Literal IP addressed are canonicalized
//   - Leading zeroes from port are trimmed
//   - Port is stripped, if matches the default
//   - Path is normalized, "." and ".." and repeated "/" are processed
//   - Empty path replaced with "/". Otherwise, trailing slash
//     is preserved
//   - unix URLs converted to the short form (unix:/path)
//
// Invalid URLs returned unmodified
func (u URL) Canonical() URL {
	return lookup(u).canonical
}

// WithHostname replaces Hostname part of the URL.
//
// Host MUST be valid hostname. IPv6 literals MUST NOT
// be enclosed into square braces.
//
// Invalid or non-TCP URLs returned unchanged.
func (u URL) WithHostname(newHostname string) URL {
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

// WithPortNum replaces Port part of the URL.
//
// If newPortNum <= 0, port will be removed. Otherwise,
// it will be set, as specified.
//
// Invalid or non-TCP URLs returned unchanged.
func (u URL) WithPortNum(newPortNum uint16) URL {
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

// WithoutPort removes Port part of the URL.
//
// Invalid or non-TCP URLs returned unchanged.
func (u URL) WithoutPort() URL {
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

// IPAddress returns URL's IP address.
// It doesn't do any name resolution, URL must contain literal IP address.
//
// If URL is invalid, or not network URL or its hostname is not literal,
// it returns a zero netip.AddrPort
func (u URL) IPAddress() netip.AddrPort {
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

// IsTCP reports if URL uses TCP-based transport.
func (u URL) IsTCP() bool {
	return lookup(u).IsTCP()
}

// IsIP4 returns true, if URL has literal IP address and this address is IPv4.
func (u URL) IsIP4() bool {
	return u.IPAddress().Addr().Is4()
}

// IsIP6 returns true, if URL has literal IP address and this address is IPv6.
func (u URL) IsIP6() bool {
	return u.IPAddress().Addr().Is6()
}

// IsHTTP reports if URL uses HTTP or HTTPS - based transport.
func (u URL) IsHTTP() bool {
	return lookup(u).IsHTTP()
}

// IsTLS reports if URL uses TLS encryption.
func (u URL) IsTLS() bool {
	return lookup(u).IsTLS()
}
