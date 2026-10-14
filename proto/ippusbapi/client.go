// MFP - Multi-Function Printers and scanners toolkit
// ipp-usb API client
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// The Client implementation

package ippusbapi

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
)

// IPPUSBControlSocket defines path to the ipp-usb control socket.
const IPPUSBControlSocket = "/var/ipp-usb/ctrl"

// Client is the ipp-usb API client.
type Client struct {
	ctx        context.Context
	dialer     net.Dialer
	transport  http.Transport
	connection http.Client
}

// NewClient returns the new client.
func NewClient(ctx context.Context) *Client {
	client := &Client{
		ctx: log.WithPrefix(ctx, "ippusbapi"),
		dialer: net.Dialer{
			Timeout: 5 * time.Second,
		},
		connection: http.Client{},
	}

	dial := func(ctx context.Context, networt, addr string) (net.Conn, error) {
		return net.Dial("unix", IPPUSBControlSocket)
	}

	client.transport.DialContext = dial
	client.connection.Transport = &client.transport

	return client
}

// Status returns the formatted ipp-usb status test.
func (client *Client) Status() ([]byte, error) {
	return client.get("status")
}

// Version returns the ipp-usb API version string
func (client *Client) Version() (string, error) {
	version, err := client.get("/version")
	if err == ErrHTTP(http.StatusNotFound) {
		// Assume 0.9.22
		return "0.9.22", nil
	}
	return strings.TrimSpace(string(version)), err
}

// get performs a GET request to the ipp-usb
func (client *Client) get(path string) (data []byte, err error) {
	defer func() {
		if err != nil {
			err = errDecode(err)
			log.Trace(client.ctx, "GET %s -- %s", path, err)
		}
	}()

	dest := "http://localhost/" + strings.TrimLeft(path, "/")
	rq, err := http.NewRequestWithContext(client.ctx, "GET", dest, nil)
	if err != nil {
		return
	}

	rsp, err := client.connection.Do(rq)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode/100 != 2 {
		err = ErrHTTP(rsp.StatusCode)
		return
	}

	log.Trace(client.ctx, "GET %s -- %s", path, rsp.Status)
	data, err = io.ReadAll(rsp.Body)
	return
}
