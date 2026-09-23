// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Cache of parsed URIs

package urilib

import (
	"net"
	"net/netip"
	"net/url"
	"path"
	"strconv"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
)

// CacheSize is the size of cache of parsed URIs
const CacheSize = 16384

// cachedURI represents an URI cache entry.
// It contains the parsed URI.
type cachedURI struct {
	parsed    *url.URL // Parsed URI (nil for invalid URI)
	err       error    // URI parsing error, if any
	canonical URI      // Cached canonical form of the URI
}

// canonicalize precomputes cachedURI.canonical
func (cached *cachedURI) canonicalize() {
	// Lowercase Scheme and Host
	parsed := *cached.parsed
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)

	// Canonicalize literal IP addresses
	if addr, err := netip.ParseAddr(parsed.Hostname()); err == nil {
		addr = addr.Unmap()
		host := addr.String()
		if addr.Is6() {
			host = "[" + host + "]"
		}

		port := parsed.Port()
		parsed.Host = host
		if port != "" {
			parsed.Host += ":" + port
		}
	}

	// Strip leading 0s from the port
	if port := parsed.Port(); port != "" && port != "0" {
		if port2 := strings.TrimLeft(port, "0"); port2 != port {
			parsed.Host, _ = strings.CutSuffix(parsed.Host, port)
			parsed.Host += port2
		}
	}

	// Remove unneeded port
	if host, port, err := net.SplitHostPort(parsed.Host); err == nil {
		portnum, err := strconv.ParseUint(port, 10, 16)
		if err == nil && portnum == uint64(cached.DefaultPort()) {
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

	cached.canonical = URI(parsed.String())
}

// Valid reports if cachedURI is valid
func (cached *cachedURI) Valid() bool {
	return cached.err == nil
}

// IsTCP reports if cachedURI uses TCP-based transport.
func (cached *cachedURI) IsTCP() bool {
	if cached.err == nil {
		switch strings.ToLower(cached.parsed.Scheme) {
		case "http", "https", "ipp", "ipps", "lpd", "socket":
			return true
		}
	}
	return false
}

// IsHTTP reports if cachedURI uses HTTP or HTTPS - based transport.
func (cached *cachedURI) IsHTTP() bool {
	if cached.err == nil {
		switch strings.ToLower(cached.parsed.Scheme) {
		case "http", "https", "ipp", "ipps":
			return true
		}
	}
	return false
}

// IsTLS reports if cachedURI uses TLS encryption.
func (cached *cachedURI) IsTLS() bool {
	if cached.err == nil {
		switch strings.ToLower(cached.parsed.Scheme) {
		case "https", "ipps":
			return true
		}
	}
	return false
}

// DefaultPort returns the default port, based on the URI scheme.
// If URI is not valid or scheme doesn't imply the port, it
// returns 0.
func (cached *cachedURI) DefaultPort() uint16 {
	if cached.IsTCP() {
		switch strings.ToLower(cached.parsed.Scheme) {
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
	}

	return 0
}

// cache is a process-global cache of parsed URIs
var cache *lru.Cache[URI, *cachedURI]

// init creates a cache
func init() {
	var err error
	cache, err = lru.New[URI, *cachedURI](CacheSize)
	if err != nil {
		panic(err)
	}
}

// lookup returns cachedURI.
// In a case of cache miss, new cachedURI will be created on demand.
func lookup(u URI) *cachedURI {
	cached, _ := cache.Get(u)
	if cached != nil {
		return cached
	}

	parsed, err := parse(string(u))
	cached = &cachedURI{err: err, canonical: u}

	if err == nil {
		cached.parsed = parsed
		cached.canonicalize()
	}

	cache.Add(u, cached)
	return cached
}
