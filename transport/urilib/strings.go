// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URI adapter for regular strings

package urilib

import (
	"net/netip"
)

// Valid reports if URI is valid.
func Valid(s string) bool {
	return URI(s).Valid()
}

// Err returns parse error for invalid URI or nil if URI is valid.
func Err(s string) error {
	return URI(s).Err()
}

// Scheme returns URI scheme (e.g., "http" or "ipp").
// For invalid URIs it returns "".
func Scheme(s string) string {
	return URI(s).Scheme()
}

// Hostname returns Hostname of the URI.
//
// If the result is enclosed in square brackets, as literal IPv6 addresses are,
// the square brackets are removed from the result.
//
// For invalid URIs it returns "".
func Hostname(s string) string {
	return URI(s).Hostname()
}

// Portname returns Portname of the URI.
func Portname(s string) string {
	return URI(s).Portname()
}

// PortNum returns port number, defined by the URI.
// It may return 0 in the following cases:
//   - URI scheme doesn't support port
//   - Port is present, but it is not numeric
//   - URI is invalid
func PortNum(s string) uint16 {
	return URI(s).PortNum()
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
func DefaultPort(s string) uint16 {
	return URI(s).DefaultPort()
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
func Canonical(s string) URI {
	return URI(s).Canonical()
}

// WithHostname replaces Hostname part of the URI.
//
// Host MUST be valid hostname. IPv6 literals MUST NOT
// be enclosed into square braces.
//
// Invalid or non-TCP URIs returned unchanged.
func WithHostname(s string, newHostname string) URI {
	return URI(s).WithHostname(newHostname)
}

// WithPortNum replaces Port part of the URI.
//
// If newPortNum <= 0, port will be removed. Otherwise,
// it will be set, as specified.
//
// Invalid or non-TCP URIs returned unchanged.
func WithPortNum(s string, newPortNum uint16) URI {
	return URI(s).WithPortNum(newPortNum)
}

// WithoutPort removes Port part of the URI.
//
// Invalid or non-TCP URIs returned unchanged.
func WithoutPort(s string) URI {
	return URI(s).WithoutPort()
}

// IPAddress returns URI's IP address.
// It doesn't do any name resolution, URI must contain literal IP address.
//
// If URI is invalid, or not network URI or its hostname is not literal,
// it returns a zero netip.AddrPort
func IPAddress(s string) netip.AddrPort {
	return URI(s).IPAddress()
}

// IsTCP reports if URI uses TCP-based transport.
func IsTCP(s string) bool {
	return URI(s).IsTCP()
}

// IsIP4 returns true, if URI has literal IP address and this address is IPv4.
func IsIP4(s string) bool {
	return URI(s).IsIP4()
}

// IsIP6 returns true, if URI has literal IP address and this address is IPv6.
func IsIP6(s string) bool {
	return URI(s).IsIP6()
}

// IsHTTP reports if URI uses HTTP or HTTPS - based transport.
func IsHTTP(s string) bool {
	return URI(s).IsHTTP()
}

// IsTLS reports if URI uses TLS encryption.
func IsTLS(s string) bool {
	return URI(s).IsTLS()
}
