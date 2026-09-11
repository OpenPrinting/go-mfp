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
