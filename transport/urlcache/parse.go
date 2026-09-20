// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URL parser

package urlcache

import (
	"fmt"
	"net/url"
	"strings"
)

// parse parses the URL string.
//
// In comparison to the [url.Parse] from the standard library,
// it has the following differences:
//
//   - URL must be absolute (i.e., must contain scheme)
//   - only "unix" URLs allowed to miss host.
//   - for the "unix" scheme host must be "" or "localhost".
//
// The "unix" URL schema is similar to the "file" schema, as defined
// in the [RFC 8089] (surprisingly, there are still no official registration
// for the "unix" schema), with the following notes:
//
//   - "authority" part of URL may be set or omitted. If set, it
//     must be either empty or "localhost" (case-insensitive). So
//     valid forms are: "unix:/path" (no authority), "unix:///path" (empty
//     authority) or "unix://localhost/path" (localhost authority).
//   - in any case, the "unix" URL is normalized into "no authority"
//     short form (i.e., "unix:/path")
func parse(s string) (*url.URL, error) {
	// Test some corner cases
	if s == "" {
		return nil, ErrURLInvalid
	}

	// Parse the URL string
	parsed, err := url.Parse(s)
	if err != nil {
		fmt.Println(s, err)
		return nil, ErrURLInvalid
	}

	// Do schema-specific checks
	switch parsed.Scheme {
	case "":
		return nil, ErrURLSchemeMissed

	case "unix":
		// For unix URLs, Scheme must be omitted or "localhost"
		switch strings.ToLower(parsed.Host) {
		case "", "localhost":
		default:
			return nil, ErrURLUNIXHost
		}

		parsed.Host = ""
		parsed.OmitHost = true

	default:
		if parsed.Host == "" {
			return nil, ErrURLHostMissed
		}
	}

	return parsed, nil
}
