// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// HTTP Transport implementation

package transport

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/OpenPrinting/go-mfp/util/missed"
)

var defaultDiaaler net.Dialer

// Transport wraps [http.Transport] and adds the following functionality:
//
//   - "ipp", "ipps" schemes support.
//   - "unix" schema support for connecting via AF_UNIX sockets.
type Transport struct {
	*http.Transport
	templateDialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// NewTransport creates a new Transport. Provided [http.Transport]
// is only used as a configuration template.
func NewTransport(template *http.Transport) *Transport {
	if template == nil {
		template = http.DefaultTransport.(*http.Transport).Clone()
		template.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	tr := &Transport{
		Transport:           template.Clone(),
		templateDialContext: template.DialContext,
	}

	tr.DialContext = tr.dialContext

	return tr
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (tr *Transport) RoundTrip(rq *http.Request) (*http.Response, error) {
	oldURL := rq.URL
	if oldURL == nil {
		return tr.Transport.RoundTrip(rq)
	}

	newURL := &url.URL{}
	*newURL = *oldURL

	// Request may use target URL that looks as follows:
	//   - http://example.com/path - normal HTTP over TCP
	//   - ipp://example.com/path  - IPP over TCP
	//   - unix:///path            - HTTP over UNIX socket
	//
	// Please notice, at that moment the schema part of URL
	// is used to distinguish between protocols (HTTP vs IPP)
	// or transport layer (TCP vs UNIX socket)
	//
	// However, server expects schemes like "http" or "https",
	// while dialContext needs the information, represented
	// by the scheme.
	//
	// So here we hack the Request URL as follows:
	//   - scheme always set to "http" or "https"
	//   - underlying socket-level protocol ("tcp" or "unix")
	//     embedded into the Host (i.e., "tcp+example.com")
	//   - for "unix", path also embedded into the Host, because
	//     dialContext doesn't see the path part of the URL
	//
	// Then dialContext() can decode this information from the
	// supplied address and use appropriately.
	host := newURL.Hostname()
	port := newURL.Port()
	defaultPort := ""
	proto := "tcp"

	switch newURL.Scheme {
	case "ipp":
		newURL.Scheme = "http"
		defaultPort = "631"

	case "ipps":
		newURL.Scheme = "https"
		defaultPort = "631"

	case "http":
		defaultPort = "80"
		proto = "tcp"

	case "https":
		defaultPort = "443"
		proto = "tcp"

	case "unix":
		newURL.Scheme = "http"
		host = escapePath(newURL.Path)
		port = "80"
		proto = "unix"

	default:
		return tr.Transport.RoundTrip(rq)
	}

	if port == "" {
		port = defaultPort
	}

	newURL.Host = net.JoinHostPort(proto+"+"+host, port)

	// Replace Request URL with the hacked URL. Restore after use
	defer func() { rq.URL = oldURL }()
	rq.URL = newURL

	return tr.Transport.RoundTrip(rq)
}

// dialContext implements DialContext callback for underlying
// http.Transport.
func (tr *Transport) dialContext(ctx context.Context,
	_, addr string) (net.Conn, error) {

	host, port, _ := net.SplitHostPort(addr)
	network, host, _ := strings.Cut(host, "+")

	addr = net.JoinHostPort(host, port)

	if network == "unix" {
		addr, _ = missed.StringsCutSuffix(addr, ":"+port)
		addr = unescapePath(addr)
	}

	dial := tr.templateDialContext
	if dial == nil {
		dial = defaultDiaaler.DialContext
	}

	return dial(ctx, network, addr)
}

// escapePath encodes path so it becomes syntactically correct
// when passed as address to dialContext.
//
// Indented for unix:/ scheme implementation.
func escapePath(in string) string {
	const hex = "0123456789abcdef"

	out := make([]byte, 0, len(in)*2)
	for i := 0; i < len(in); i++ {
		c := in[i]

		if 'a' <= c && c <= 'z' ||
			'A' <= c && c <= 'Z' ||
			'0' <= c && c <= '9' {
			out = append(out, c)
		} else {
			out = append(out, '-', hex[c>>4], hex[c&0xf])
		}
	}

	return string(out)
}

// unescapePath decodes path, encoded by escapePath().
func unescapePath(in string) string {
	type state int
	const (
		stateNorm state = iota
		stateEsc1
		stateEsc2
	)

	out := make([]byte, 0, len(in))
	decodeState := stateNorm
	escval := 0

	for i := 0; i < len(in); i++ {
		c := in[i]

		switch decodeState {
		case stateNorm:
			if c != '-' {
				out = append(out, c)
			} else {
				decodeState = stateEsc1
			}

		case stateEsc1, stateEsc2:
			v := hexval(c)
			switch {
			case v < 0:
				out = append(out, c)
				decodeState = stateNorm

			case decodeState == stateEsc1:
				escval = v
				decodeState = stateEsc2

			case decodeState == stateEsc2:
				escval = (escval << 4) + v
				out = append(out, byte(escval))
				decodeState = stateNorm
			}
		}
	}

	return string(out)
}

// hexval returns hexadecimal value of c, or -1 if c is not
// hexadecimal digit.
func hexval(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c-'a') + 10
	case 'A' <= c && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}
