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
	"path"
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
	return u.Err() == nil
}

// Err returns parse error for invalid URL or nil if URL is valid.
func (u URL) Err() error {
	return lookup(u).err
}

// Scheme returns URL scheme (e.g., "http" or "ipp").
// For invalid URLs it returns "".
func (u URL) Scheme() string {
	cached := lookup(u)
	if cached.err != nil {
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
	if cached.err != nil {
		return ""
	}

	return cached.parsed.Hostname()
}

// Portname returns Portname of the URL.
func (u URL) Portname() string {
	cached := lookup(u)
	if cached.err != nil {
		return ""
	}

	return cached.parsed.Port()
}

// PortNum returns port number, defined by the URL.
// It may return 0 in the following cases:
//   - URL scheme doesn't support port
//   - Port is present, but it is not numeric
//   - URL is invalid
func (u URL) PortNum() int {
	p := u.Portname()
	if p == "" {
		return u.DefaultPort()
	}

	n, err := strconv.ParseUint(p, 10, 16)
	if err != nil {
		return 0
	}

	return int(n)
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
func (u URL) DefaultPort() int {
	switch strings.ToLower(u.Scheme()) {
	case "http":
		return 80
	case "https":
		return 443
	case "ipp", "ipps":
		return 631
	case "lpd":
		return 515
	case "socket":
		return 9100
	}
	return 0
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
	cached := lookup(u)
	if cached.err != nil {
		return u
	}

	parsed := *cached.parsed

	// Lowercase Scheme and Host
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)

	// Canonicalize literal IP addresses
	if addr, err := netip.ParseAddr(parsed.Hostname()); err == nil {
		addr = addr.Unmap()
		host := addr.String()
		if addr.Is6() {
			host = "[" + host + "]"
		}

		port := u.Portname()
		parsed.Host = host
		if port != "" {
			parsed.Host += ":" + port
		}
	}

	// Strip leading 0s from the port
	if port := u.Portname(); port != "" && port != "0" {
		if port2 := strings.TrimLeft(port, "0"); port2 != port {
			parsed.Host, _ = strings.CutSuffix(parsed.Host, port)
			parsed.Host += port2
		}
	}

	// Remove unneeded port
	if host, port, err := net.SplitHostPort(parsed.Host); err == nil {
		portnum, err := strconv.ParseUint(port, 10, 16)
		if err == nil && portnum == uint64(u.DefaultPort()) {
			parsed.Host = host
		}
	}

	// Drop Host for unix scheme
	if parsed.Scheme == "unix" {
		parsed.Host = ""
		parsed.OmitHost = true
	}

	// Normalize path
	preserveSlash := strings.HasSuffix(parsed.Path, "/")

	switch parsed.Path {
	case "", ".":
		parsed.Path = "/"

	default:
		parsed.Path = path.Clean(parsed.Path)
		if preserveSlash && !strings.HasSuffix(parsed.Path, "/") {
			parsed.Path += "/"
		}
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
	if cached.err != nil {
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
	switch strings.ToLower(u.Scheme()) {
	case "http", "https", "ipp", "ipps", "lpd", "socket":
		return true
	}
	return false
}

// IsHTTP reports if URL uses HTTP or HTTPS - based transport.
func (u URL) IsHTTP() bool {
	switch strings.ToLower(u.Scheme()) {
	case "http", "https", "ipp", "ipps":
		return true
	}
	return false
}

// IsTLS reports if URL uses TLS encryption.
func (u URL) IsTLS() bool {
	switch strings.ToLower(u.Scheme()) {
	case "https", "ipps":
		return true
	}
	return false
}
