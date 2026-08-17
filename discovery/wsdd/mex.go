// MFP - Miulti-Function Printers and scanners toolkit
// WSD device discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Metadate Exchange (MEX)

package wsdd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/wsd"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// mexData wraps wsd.Metadata and adds few additional fields
type mexData struct {
	wsd.Metadata          // The metadata itself
	from         *url.URL // URL it comes from
}

// mexGetter retrieves WSD metadata by XAddr URL.
type mexGetter struct {
	back  *backend                    // Parent backend
	http  http.Client                 // HTTP client
	cache map[mexCacheID]*mexCacheEnt // Cached metadata
	lock  sync.Mutex                  // Access lock
}

// newMexgetter creates a new mexGetter
func newMexGetter(back *backend) *mexGetter {
	mg := &mexGetter{
		back: back,
		http: http.Client{
			Timeout: 5 * time.Second,
		},
		cache: make(map[mexCacheID]*mexCacheEnt),
	}

	return mg
}

// Get retrives the mexData from the device.
//
// Parameters are:
//   - ifidx is the index of the local interface where request is initiated
//   - target is the WSD target service address
//   - ver is the MetadataVersion. It comes together with the
//     XAddr URLs as a part of the wsd.Announce structure.
//
// The MetadataVersion affects metadata caching. The never metadata
// overrides the cached version.
func (mg *mexGetter) Get(ctx context.Context,
	ifidx int, target wsd.AnyURI,
	xaddr *url.URL, ver uint64) []mexData {

	// Create mexCacheID
	id := mexCacheID{ifidx, target, xaddr.String()}

	literal := urlIsLiteral(xaddr)
	if literal {
		// Interface index doesn't affect processing of the
		// literal URLs
		id.ifidx = 0
	}

	// Obtain cache entry or create a new one
	ent, justwait := mg.cacheLookup(id, ver)

	// Fetching already done?
	if ent.isDone() {
		return ent.metadata
	}

	// Just wait?
	if justwait {
		ent.wait()
		return ent.metadata
	}

	// We were the first here with this XAddr, so it is our
	// responsibility to fetch the metadata.
	var xaddrs []*url.URL
	if literal {
		xaddrs = []*url.URL{xaddr}
	} else {
		xaddrs = mg.back.res.Resolve(ctx, ifidx, xaddr)
	}

	var metadata []mexData
	if len(xaddrs) > 0 {
		metadata = mg.fetch(ctx, xaddrs, target)
	}

	// Update the cache entry
	mg.cacheUpdate(id, ent, metadata)

	// Return whatever we have
	return ent.metadata
}

// cacheLookup lookups the metadata cache.
//
// It returns new or existing cache entry and 'true' as a second
// returned value, if existent cache entry was found for this if.
func (mg *mexGetter) cacheLookup(id mexCacheID,
	ver uint64) (*mexCacheEnt, bool) {

	mg.lock.Lock()
	defer mg.lock.Unlock()

	ent := mg.cache[id]
	if ent != nil {
		return ent, true
	}

	ent = &mexCacheEnt{
		ver:      ver,
		waitchan: make(chan struct{}),
	}
	mg.cache[id] = ent

	return ent, false
}

// cacheUpdate updates metadata cache entry.
func (mg *mexGetter) cacheUpdate(id mexCacheID,
	ent *mexCacheEnt, metadata []mexData) {

	mg.lock.Lock()
	defer mg.lock.Unlock()

	if len(metadata) > 0 {
		ent.metadata = metadata
	} else {
		delete(mg.cache, id)
	}

	ent.done()
}

// fetch fetches the metadata
func (mg *mexGetter) fetch(ctx context.Context,
	xaddrs []*url.URL, target wsd.AnyURI) []mexData {

	// Fetch metadata
	var wait sync.WaitGroup
	var lock sync.Mutex
	var metadata []mexData

	wait.Add(len(xaddrs))

	for _, xaddr := range xaddrs {
		go func(xaddr2 *url.URL) {
			meta, err := mg.fetchHTTP(ctx, target, xaddr2)
			if err == nil {
				lock.Lock()
				metadata = append(metadata, meta)
				lock.Unlock()
			}

			wait.Done()
		}(xaddr)
	}

	wait.Wait()

	return metadata
}

// fetchHTTP performs HTTP query for the WSD metadata
func (mg *mexGetter) fetchHTTP(ctx context.Context,
	target wsd.AnyURI, xaddr *url.URL) (meta mexData, err error) {

	// Create a request
	msgid := wsd.AnyURI(uuid.Random().URN())
	msg := wsd.Msg{
		Header: wsd.Header{
			Action:    wsd.ActGet,
			MessageID: msgid,
			To:        optional.New(target),
			ReplyTo: optional.New(
				wsd.EndpointReference{Address: wsd.ToAnonymous},
			),
		},
		Body: wsd.Get{},
	}

	data := msg.Encode()

	rq := &http.Request{
		Method: "POST",
		URL:    xaddr,
		Header: http.Header{
			"Content-Type": []string{
				"application/soap+xml; charset=utf-8",
			},
		},
		Body:          io.NopCloser(bytes.NewReader(data)),
		ContentLength: int64(len(data)),
		Close:         true,
	}

	rq = rq.WithContext(ctx)

	// Perform HTTP query
	mg.back.debug("POST %s\n%s", xaddr, msg.Format())

	rsp, err := mg.http.Do(rq)
	if err != nil {
		mg.back.warning("POST %s: %s", xaddr, err)
		return
	}

	defer rsp.Body.Close()

	if rsp.StatusCode/100 != 2 {
		err = fmt.Errorf("Unexpected HTTP status: %s", rsp.Status)
		mg.back.warning("POST %s: %s", xaddr, err)
		return
	}

	// Fetch HTTP response body
	data, err = io.ReadAll(io.LimitReader(rsp.Body,
		int64(wsddMetadataGetMaxResponse+1)))

	if err != nil {
		mg.back.warning("POST %s: %s", xaddr, err)
		return
	}

	if len(data) > wsddMetadataGetMaxResponse {
		err = fmt.Errorf("HTTP response too large")
		mg.back.warning("POST %s: %s", xaddr, err)
		return
	}

	mg.back.debug("POST %s: %s", xaddr, rsp.Status)

	// Decode response
	msg, err = wsd.DecodeMsg(data)
	if err != nil {
		mg.back.warning("POST %s: %s", xaddr, err)
		return
	}

	metadata, ok := msg.Body.(wsd.Metadata)
	if !ok {
		err = fmt.Errorf("Unexpected WSD response: %s",
			msg.Header.Action)

		mg.back.warning("POST %s: %s", xaddr, err)
		return
	}

	meta.Metadata = metadata
	meta.from = xaddr

	return
}

// mexCacheID is the MEX cache entry ID.
type mexCacheID struct {
	ifidx  int        // Local interface address
	target wsd.AnyURI // Target service WSD "address"
	xaddr  string     // Requested URL string
}

// mexCacheEnt is the MEX cache entry.
type mexCacheEnt struct {
	ver      uint64        // Requested metadata version
	metadata []mexData     // Cached metadata
	waitchan chan struct{} // Closed when fetch is complete
}

// done marks cache entry as done downloading and awakes all
// goroutines blocked at mexCacheEnt.wait
func (ent *mexCacheEnt) done() {
	close(ent.waitchan)
}

// isDone reports if entry processing is done.
func (ent *mexCacheEnt) isDone() bool {
	select {
	case <-ent.waitchan:
		return true
	default:
		return false
	}
}

// wait waits until metadata retrieval completion.
func (ent *mexCacheEnt) wait() {
	<-ent.waitchan
}
