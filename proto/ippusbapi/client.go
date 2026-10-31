// MFP - Multi-Function Printers and scanners toolkit
// ipp-usb API client
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// The Client implementation

package ippusbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/util/optional"
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

// DeviceName represents set of parameters to identify
// a particular device.
type DeviceName struct {
	MFG  optional.Val[string] `json:"usb_MFG"`  // Manufacturer name
	MDL  optional.Val[string] `json:"usb_MDL"`  // Model name
	SN   optional.Val[string] `json:"usb_SN"`   // Serial number
	HWID optional.Val[string] `json:"usb_HWID"` // Vendor:Product, hex
	URI  optional.Val[string] `json:"usb_URI"`  // CUPS-compatible URI
}

// DeviceConf represents a device configuration.
type DeviceConf struct {
	DeviceName                    // Device name
	Disabled   optional.Val[bool] `json:"disabled"` // True if device is disabled
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
	return client.get("/status")
}

// Version returns the ipp-usb API version string
func (client *Client) Version() (string, error) {
	version, err := client.get("/version")
	if err == ErrNotSupported {
		// Assume 0.9.22
		return "0.9.22", nil
	}
	return strings.TrimSpace(string(version)), err
}

// GetConf returns configuration of devices, searched by name.
func (client *Client) GetConf(name DeviceName) ([]DeviceConf, error) {
	out, err := client.post("/get/conf", name)
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, nil
	}

	items, ok := out.([]json.RawMessage)
	if !ok {
		err = errors.New("invalid response received")
		return nil, err
	}

	configs := make([]DeviceConf, len(items))
	for i, item := range items {
		decoder := json.NewDecoder(bytes.NewReader(item))
		err = decoder.Decode(&configs[i])
		if err != nil {
			return nil, err
		}
	}
	return configs, nil
}

// get performs a GET request to the ipp-usb
func (client *Client) get(path string) (data []byte, err error) {
	// Ensure uniform error handling
	defer func() {
		if err != nil {
			err = errDecode(err)
			log.Verbose(client.ctx, "GET %s -- %s", path, err)
		}
	}()

	// Build GET request
	dest := "http://localhost/" + strings.TrimLeft(path, "/")
	rq, err := http.NewRequestWithContext(client.ctx, "GET", dest, nil)
	if err != nil {
		return
	}

	// Dump request
	dump, _ := httputil.DumpRequestOut(rq, true)
	log.Begin(client.ctx).
		Trace("Request sent:").
		Indent(log.LevelTrace, 2, dump).
		Commit()

	// Perform server query
	rsp, err := client.connection.Do(rq)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	// Dump response
	dump, _ = httputil.DumpResponse(rsp, true)
	log.Begin(client.ctx).
		Trace("Response received:").
		Indent(log.LevelTrace, 2, dump).
		Commit()

	// Decode response
	if rsp.StatusCode/100 != 2 {
		err = ErrHTTP(rsp.StatusCode)
		return
	}

	log.Verbose(client.ctx, "GET %s -- %s", path, rsp.Status)
	data, err = io.ReadAll(rsp.Body)
	return
}

// get performs a POST request to the ipp-usb
//
// The 'in' parameter is any JSON-encodable value.
// The 'out' value is always:
//   - nil for no data
//   - json.RawMessage for JSON object
//   - []json.RawMessage for JSON array
func (client *Client) post(path string, in any) (out any, err error) {
	// Ensure uniform error handling
	defer func() {
		if err != nil {
			err = errDecode(err)
			log.Verbose(client.ctx, "GET %s -- %s", path, err)
		}
	}()

	// Encode JSON request
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(in)
	assert.NoError(err)

	// Build GET request
	dest := "http://localhost/" + strings.TrimLeft(path, "/")
	rq, err := http.NewRequestWithContext(client.ctx, "POST", dest, nil)
	if err != nil {
		return
	}

	rq.Header.Set("Content-Type", "application/json")
	rq.Body = io.NopCloser(&buf)
	rq.ContentLength = int64(buf.Len())

	// Dump request
	log.Begin(client.ctx).
		Trace("Request sent:").
		Object(log.LevelTrace, 2, (*httpOutRequestLogger)(rq)).
		Object(log.LevelTrace, 2, jsonLogger(buf.Bytes())).
		Commit()

	// Perform server query
	rsp, err := client.connection.Do(rq)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	buf.Reset()
	_, err = io.Copy(&buf, rsp.Body)
	if err != nil {
		return
	}

	// Dump response
	log.Begin(client.ctx).
		Trace("Response received:").
		Object(log.LevelTrace, 2, (*httpResponseLogger)(rsp)).
		Object(log.LevelTrace, 2, jsonLogger(buf.Bytes())).
		Commit()

	// Decode response
	if rsp.StatusCode/100 != 2 {
		err = ErrHTTP(rsp.StatusCode)
		return
	}

	type envelope struct {
		Err  *string         `json:"err"`
		Data json.RawMessage `json:"data"`
	}

	var env envelope

	err = json.NewDecoder(&buf).Decode(&env)
	if err != nil {
		return
	}

	if env.Err != nil {
		err = errors.New(*env.Err)
		return
	}

	// Normalize output.
	//
	// It must be:
	//   - nil for no data
	//   - json.RawMessage for JSON object
	//   - []json.RawMessage for JSON array
	env.Data = bytes.TrimSpace(env.Data)
	switch {
	case len(env.Data) == 0:
		out = nil
		out = env.Data

	case bytes.HasPrefix(env.Data, []byte{'['}):
		var array []json.RawMessage

		err = json.NewDecoder(bytes.NewReader(env.Data)).Decode(&array)
		assert.NoError(err)

		out = array

	default:
		out = env.Data
	}

	// Log the successful transaction
	log.Verbose(client.ctx, "POST %s -- %s", path, rsp.Status)

	return
}
