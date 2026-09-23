// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URI parser

package urilib

import (
	"fmt"
	"net/url"
	"strings"
)

// parse parses the URI string.
//
// In comparison to the [url.Parse] from the standard library,
// it has the following differences:
//
//   - URI must be absolute (i.e., must contain scheme)
//   - only "unix" URIs allowed to miss host.
//   - for the "unix" scheme host must be "" or "localhost".
//
// The "unix" URI schema is similar to the "file" schema, as defined
// in the [RFC 8089] (surprisingly, there are still no official registration
// for the "unix" schema), with the following notes:
//
//   - "authority" part of URI may be set or omitted. If set, it
//     must be either empty or "localhost" (case-insensitive). So
//     valid forms are: "unix:/path" (no authority), "unix:///path" (empty
//     authority) or "unix://localhost/path" (localhost authority).
//   - in any case, the "unix" URI is normalized into "no authority"
//     short form (i.e., "unix:/path")
func parse(s string) (*url.URL, error) {
	// Test some corner cases
	if s == "" {
		return nil, ErrURIInvalid
	}

	// Parse the URI string
	parsed, err := url.Parse(s)
	if err != nil {
		fmt.Println(s, err)
		return nil, ErrURIInvalid
	}

	// Do schema-specific checks
	switch parsed.Scheme {
	case "":
		return nil, ErrURISchemeMissed

	case "unix":
		// For unix URIs, Scheme must be omitted or "localhost"
		switch strings.ToLower(parsed.Host) {
		case "", "localhost":
		default:
			return nil, ErrURIUNIXHost
		}

		parsed.Host = ""
		parsed.OmitHost = true

	default:
		if parsed.Host == "" {
			return nil, ErrURIHostMissed
		}
	}

	return parsed, nil
}
