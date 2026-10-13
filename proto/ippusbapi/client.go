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
	"errors"
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

	dial := func(networt, addr string) (net.Conn, error) {
		return net.Dial("unix", IPPUSBControlSocket)
	}

	client.transport.Dial = dial
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
	return strings.TrimSpace(string(version)), err
}

// get performs a GET request to the ipp-usb
func (client *Client) get(path string) ([]byte, error) {
	dest := "http://localhost/" + strings.TrimLeft(path, "/")
	rq, err := http.NewRequestWithContext(client.ctx, "GET", dest, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := client.connection.Do(rq)
	if err != nil {
		log.Trace(client.ctx, "GET %s -- %s", path, err)
		return nil, err
	}
	defer rsp.Body.Close()

	log.Trace(client.ctx, "GET %s -- %s", path, rsp.Status)
	if rsp.StatusCode/100 != 2 {
		err := errors.New("HTTP: " + rsp.Status)
		return nil, err
	}

	data, err := io.ReadAll(rsp.Body)
	return data, err
}
