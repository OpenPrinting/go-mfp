// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Cache of parsed URLs

package urlcache

import (
	"net/url"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
)

// CacheSize is the size of cache of parsed URLs
const CacheSize = 16384

// cachedURL represents an URL cache entry.
// It contains the parsed URL.
type cachedURL struct {
	parsed *url.URL
	err    error
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
	cached = &cachedURL{err: err}

	if err == nil {
		cached.parsed = parsed
	}

	cache.Add(u, cached)
	return cached
}
