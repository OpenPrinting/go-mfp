// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Cache of parsed URLs

package urlcache

import (
	"net"
	"net/netip"
	"net/url"
	"path"
	"strconv"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
)

// CacheSize is the size of cache of parsed URLs
const CacheSize = 16384

// cachedURL represents an URL cache entry.
// It contains the parsed URL.
type cachedURL struct {
	parsed    *url.URL // Parsed URL (nil for invalid URL)
	err       error    // URL parsing error, if any
	canonical URL      // Cached canonical form of the URL
}

// canonicalize precomputes cachedURL.canonical
func (cached *cachedURL) canonicalize() {
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

	cached.canonical = URL(parsed.String())
}

// Valid reports if cachedURL is valid
func (cached *cachedURL) Valid() bool {
	return cached.err == nil
}

// IsTCP reports if cachedURL uses TCP-based transport.
func (cached *cachedURL) IsTCP() bool {
	if cached.err == nil {
		switch strings.ToLower(cached.parsed.Scheme) {
		case "http", "https", "ipp", "ipps", "lpd", "socket":
			return true
		}
	}
	return false
}

// IsHTTP reports if cachedURL uses HTTP or HTTPS - based transport.
func (cached *cachedURL) IsHTTP() bool {
	if cached.err == nil {
		switch strings.ToLower(cached.parsed.Scheme) {
		case "http", "https", "ipp", "ipps":
			return true
		}
	}
	return false
}

// IsTLS reports if cachedURL uses TLS encryption.
func (cached *cachedURL) IsTLS() bool {
	if cached.err == nil {
		switch strings.ToLower(cached.parsed.Scheme) {
		case "https", "ipps":
			return true
		}
	}
	return false
}

// DefaultPort returns the default port, based on the URL scheme.
// If URL is not valid or scheme doesn't imply the port, it
// returns 0.
func (cached *cachedURL) DefaultPort() int {
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

// cache is a process-global cache of parsed URLs
var cache *lru.Cache[URL, *cachedURL]

// init creates a cache
func init() {
	var err error
	cache, err = lru.New[URL, *cachedURL](CacheSize)
	if err != nil {
		panic(err)
	}
}

// lookup returns cachedURL.
// In a case of cache miss, new cachedURL will be created on demand.
func lookup(u URL) *cachedURL {
	cached, _ := cache.Get(u)
	if cached != nil {
		return cached
	}

	parsed, err := parse(string(u))
	cached = &cachedURL{err: err, canonical: u}

	if err == nil {
		cached.parsed = parsed
		cached.canonicalize()
	}

	cache.Add(u, cached)
	return cached
}
