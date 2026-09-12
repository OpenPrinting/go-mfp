// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// URL tests

package urlcache

import (
	"net/netip"
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

// TestURLWithHostname tests URL.WithHostname
func TestURLWithHostname(t *testing.T) {
	type testData struct {
		u    URL
		host string
		out  URL
	}

	tests := []testData{
		// Replace name with name. Path is preserved.
		{"http://example.com", "localhost", "http://localhost"},
		{"http://example.com/", "localhost", "http://localhost/"},

		// Replace name with name. Port is preserved.
		{"http://example.com:80", "localhost", "http://localhost:80"},

		// Name vs IP4 literal
		{"http://example.com", "127.0.0.1", "http://127.0.0.1"},
		{"http://example.com:8080", "127.0.0.1", "http://127.0.0.1:8080"},
		{"http://127.0.0.1", "example.com", "http://example.com"},
		{"http://127.0.0.1:8080", "example.com", "http://example.com:8080"},

		// Name vs IP6 literal
		{"http://example.com", "::1", "http://[::1]"},
		{"http://example.com:80", "::1", "http://[::1]:80"},
		{"http://[::1]", "example.com", "http://example.com"},
		{"http://[::1]:80", "example.com", "http://example.com:80"},

		// IP4 vs IP6 literal
		{"http://127.0.0.1", "::1", "http://[::1]"},
		{"http://127.0.0.1:80", "::1", "http://[::1]:80"},
		{"http://[::1]", "127.0.0.1", "http://127.0.0.1"},
		{"http://[::1]:80", "127.0.0.1", "http://127.0.0.1:80"},

		// Invalid or non-TCP URL
		{"invalid URL", "localhost", "invalid URL"},
		{"unix:/path", "localhost", "unix:/path"},
	}

	for _, test := range tests {
		out := test.u.WithHostname(test.host)
		if out != test.out {
			t.Errorf("%q: URL.WithHostname(%q) failed:\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.u, test.host, test.out, out)
		}
	}
}

// TestURLWithPortNum tests URL.WithPortNum
func TestURLWithPortNum(t *testing.T) {
	type testData struct {
		u    URL
		port int
		out  URL
	}

	tests := []testData{
		// Symbolic URLs
		{"http://example.com", 1234, "http://example.com:1234"},
		{"http://example.com", 80, "http://example.com:80"},
		{"http://example.com:8888", 1234, "http://example.com:1234"},
		{"http://example.com:8888", 80, "http://example.com:80"},
		{"http://example.com:8888", 0, "http://example.com"},

		// IP4 literal
		{"http://192.168.0.1", 1234, "http://192.168.0.1:1234"},
		{"http://192.168.0.1:80", 1234, "http://192.168.0.1:1234"},
		{"http://192.168.0.1:80", 0, "http://192.168.0.1"},

		// IP6 literal
		{"http://[::1]", 1234, "http://[::1]:1234"},
		{"http://[::1]:80", 1234, "http://[::1]:1234"},
		{"http://[::1]:80", 0, "http://[::1]"},

		// Invalid or non-TCP URL
		{"invalid URL", 80, "invalid URL"},
		{"unix:/path", 80, "unix:/path"},
	}

	for _, test := range tests {
		out := test.u.WithPortNum(test.port)
		if out != test.out {
			t.Errorf("%q: URL.WithPortNum(%v) failed:\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.u, test.port, test.out, out)
		}
	}
}

// TestURLWithoutPort tests URL.WithoutPort
func TestURLWithoutPort(t *testing.T) {
	type testData struct {
		u   URL
		out URL
	}

	tests := []testData{
		{"http://example.com:8888", "http://example.com"},
		{"http://192.168.0.1:80", "http://192.168.0.1"},
		{"http://[::1]:80", "http://[::1]"},
		{"invalid URL", "invalid URL"},
		{"unix:/path", "unix:/path"},
	}

	for _, test := range tests {
		out := test.u.WithoutPort()
		if out != test.out {
			t.Errorf("%q: URL.WithoutPort() failed:\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.u, test.out, out)
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

// TestURLIPAddress tests URL.IPAddress
func TestURLIPAddress(t *testing.T) {
	type testData struct {
		u    URL
		addr netip.AddrPort
	}

	tests := []testData{
		{"http://127.0.0.1", netip.MustParseAddrPort("127.0.0.1:80")},
		{"http://127.0.0.1:80", netip.MustParseAddrPort("127.0.0.1:80")},
		{"http://[::1]", netip.MustParseAddrPort("[::1]:80")},
		{"http://[::1]:80", netip.MustParseAddrPort("[::1]:80")},
		{"http://example.com", netip.AddrPort{}},
	}

	for _, test := range tests {
		addr := test.u.IPAddress()
		if addr != test.addr {
			t.Errorf("%q: URL.IPAddress failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.addr, addr)
		}
	}
}

// TestURLIs tests URL.IsTCP, URL.IsHTTP, URL.IsTLS, ... methods
func TestURLIs(t *testing.T) {
	type testData struct {
		u      URL
		valid  bool
		isTCP  bool
		isHTTP bool
		isTLS  bool
	}

	tests := []testData{
		{
			u:      "http://example.com",
			valid:  true,
			isTCP:  true,
			isHTTP: true,
			isTLS:  false,
		},

		{
			u:      "https://example.com",
			valid:  true,
			isTCP:  true,
			isHTTP: true,
			isTLS:  true,
		},

		{
			u:      "ipp://example.com",
			valid:  true,
			isTCP:  true,
			isHTTP: true,
			isTLS:  false,
		},

		{
			u:      "ipps://example.com",
			valid:  true,
			isTCP:  true,
			isHTTP: true,
			isTLS:  true,
		},

		{
			u:      "lpd://example.com",
			valid:  true,
			isTCP:  true,
			isHTTP: false,
			isTLS:  false,
		},

		{
			u:      "socket://example.com",
			valid:  true,
			isTCP:  true,
			isHTTP: false,
			isTLS:  false,
		},

		{
			u:      "unix:/path",
			valid:  true,
			isTCP:  false,
			isHTTP: false,
			isTLS:  false,
		},

		{
			u:      "invalid URL",
			valid:  false,
			isTCP:  false,
			isHTTP: false,
			isTLS:  false,
		},
	}

	for _, test := range tests {
		valid := test.u.Valid()
		isTCP := test.u.IsTCP()
		isHTTP := test.u.IsHTTP()
		isTLS := test.u.IsTLS()

		if valid != test.valid {
			t.Errorf("%q: URL.Valid failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.valid, valid)
		}

		if isTCP != test.isTCP {
			t.Errorf("%q: URL.IsTCP failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.isTCP, isTCP)
		}

		if isHTTP != test.isHTTP {
			t.Errorf("%q: URL.IsHTTP failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.isHTTP, isHTTP)
		}

		if isTLS != test.isTLS {
			t.Errorf("%q: URL.IsTLS failed:\n"+
				"expected: %v\n"+
				"present:  %v\n", test.u, test.isTLS, isTLS)
		}
	}
}
