// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URL tests

package urlcache

import (
	"testing"
)

// TestURLValid tests URL.Valid function
func TestURLValid(t *testing.T) {
	type testData struct {
		u     URL
		valid bool
	}

	tests := []testData{
		{"http://example.com", true},
		{"invalid URL", false},
		{"localhost", false},
		{"", false},
		{"unix://localhost", true},
		{"unix://localhost/", true},
		{"unix://example.com/", false},
		{"unix:/path", true},
		{"unix:///path", true},
		{"unix:", true},
	}

	for _, test := range tests {
		valid := test.u.Valid()
		if valid != test.valid {
			t.Errorf("%q: URL.Valid failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.valid, valid)
		}
	}
}

// TestURLCanonical tests URL.Canonical function
func TestURLCanonical(t *testing.T) {
	type testData struct {
		u         URL
		canonical URL
	}

	tests := []testData{
		// Invalid URLs are not modified
		{"", ""},
		{"invalid URL", "invalid URL"},

		// Scheme and Host are lowercased
		{"HtTp://example.com", "http://example.com/"},
		{"http://EXAMPLE.com", "http://example.com/"},

		// Literal IP addresses are canonicalized
		{"http://[::1]/", "http://[::1]/"},
		{"http://[0::1]/", "http://[::1]/"},
		{"http://127.0.0.1:8080", "http://127.0.0.1:8080/"},

		// Leading zeroes from port are trimmed
		{"http://127.0.0.1:08080", "http://127.0.0.1:8080/"},
		{"http://example.com:08080", "http://example.com:8080/"},
		{"http://example.com:0", "http://example.com:0/"},

		// Port is stripped, if matches the default
		{"http://example.com:80", "http://example.com/"},
		{"http://example.com:080", "http://example.com/"},
		{"http://example.com:8080", "http://example.com:8080/"},
		{"https://example.com:443", "https://example.com/"},
		{"https://example.com:4430", "https://example.com:4430/"},
		{"ipp://example.com:631", "ipp://example.com/"},
		{"ipp://example.com:80", "ipp://example.com:80/"},
		{"lpd://example.com:515", "lpd://example.com/"},
		{"lpd://example.com:5150", "lpd://example.com:5150/"},
		{"socket://example.com:9100", "socket://example.com/"},
		{"socket://example.com:910", "socket://example.com:910/"},

		// . and .. and repeated / in Path are processed
		{"http://example.com/foo///bar", "http://example.com/foo/bar"},
		{"http://example.com/foo/./bar", "http://example.com/foo/bar"},
		{"http://example.com/foo/../bar", "http://example.com/bar"},
		{"http://example.com/foo/../../../bar", "http://example.com/bar"},

		// Empty path replaced with "/"
		{"http://example.com", "http://example.com/"},
		{"http://example.com/", "http://example.com/"},

		// Otherwise, trailing "/" is preserved
		{"http://example.com/path", "http://example.com/path"},
		{"http://example.com/path/", "http://example.com/path/"},

		// unix URLs converted to the short form (unix:/path)
		{"unix://localhost", "unix:/"},
		{"unix://localhost/", "unix:/"},
		{"unix://localhost/path", "unix:/path"},
	}

	for _, test := range tests {
		c := test.u.Canonical()
		if c != test.canonical {
			t.Errorf("%q: URL.Canonical failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.canonical, c)
		}
	}
}

// TestURLHostPortName tests URL.Hostname and URL.Portname
func TestURLHostPortName(t *testing.T) {
	type testData struct {
		u          URL
		host, port string
	}

	tests := []testData{
		{"http://example.com", "example.com", ""},
		{"http://example.com:80", "example.com", "80"},
		{"http://127.0.0.1", "127.0.0.1", ""},
		{"http://127.0.0.1:80", "127.0.0.1", "80"},
		{"http://[::1]:80", "::1", "80"},
		{"http://[fe80::217:c8ff:fe7b:6a91%252]:9095/eSCL/", "fe80::217:c8ff:fe7b:6a91%2", "9095"},
	}

	for _, test := range tests {
		host := test.u.Hostname()
		port := test.u.Portname()

		if host != test.host {
			t.Errorf("%q: URL.Hostname failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.host, host)
		}

		if port != test.port {
			t.Errorf("%q: URL.Portname failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.port, port)
		}
	}
}
