// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URL resolver

package wsdd

import (
	"context"
	"errors"
	"net/url"
	"sync"

	"github.com/OpenPrinting/go-mfp/discovery/dnssd"
	"github.com/OpenPrinting/go-mfp/log"
)

// urlResolver resolves symbolic URLs into the IP literal URLs.
type urlResolver struct {
	res       *dnssd.Resolver // DNSSD Resolver, shared
	usecnt    int             // Use count
	lock      sync.Mutex      // Access lock
	closing   bool            // Closing in progress
	closewait sync.WaitGroup  // For urlResolver.Close
}

// newURLResolver creates new urlResolver.
func newURLResolver(src *source) *urlResolver {
	return &urlResolver{}
}

// Close closes the urlResolver.
func (urlres *urlResolver) Close() {
	urlres.lock.Lock()
	urlres.closing = true
	urlres.lock.Unlock()

	urlres.closewait.Wait()
}

// Resolve resolves symbolic URLs into the IP literal URLs.
// As the hostname may have multiple IP addresses, resolving
// may return multiple URLs.
func (urlres *urlResolver) Resolve(ctx context.Context,
	ifidx int, u *url.URL) []*url.URL {

	hostname := u.Hostname()

	// Create a resolver
	res, err := urlres.getDNSSDResolver()
	if err != nil {
		log.Warning(ctx, "resolve %q: %s", hostname, err)
		return nil
	}

	defer urlres.putDNSSDResolver()

	// Resolve hostname
	addrs, err := res.LookupHost(ctx, ifidx, hostname)
	if err != nil {
		log.Warning(ctx, "resolve %q: %s", hostname, err)
		return nil
	}

	if len(addrs) == 0 {
		err = errors.New("no addresses")
		log.Warning(ctx, "resolve %q: %s", hostname, err)
		return nil
	}

	// Create literal URLs
	urls := make([]*url.URL, len(addrs))
	for i, addr := range addrs {
		urls[i] = urlWithHostname(u, addr.String())
	}

	return urls
}

// getDNSSDResolver returns the dnssd.Resolver.
func (urlres *urlResolver) getDNSSDResolver() (*dnssd.Resolver, error) {
	urlres.lock.Lock()
	defer urlres.lock.Unlock()

	if urlres.closing {
		err := errors.New("close in progress")
		return nil, err
	}

	if urlres.res == nil {
		var err error
		urlres.res, err = dnssd.NewResolver()
		if err != nil {
			return nil, err
		}
	}

	urlres.usecnt++
	urlres.closewait.Add(1)

	return urlres.res, nil
}

// putDNSSDResolver releases the dnssd.Resolver.
func (urlres *urlResolver) putDNSSDResolver() {
	urlres.lock.Lock()
	defer urlres.lock.Unlock()

	urlres.usecnt--
	if urlres.usecnt == 0 {
		urlres.res.Close()
		urlres.res = nil
		urlres.closewait.Done()
	}
}
