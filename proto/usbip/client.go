// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// USBIP client

package usbip

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/netip"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// DefaultPort is the default port for TCP connections
const DefaultPort = 3240

// Client represents the USBIP client.
type Client struct {
	conn net.Conn         // Underlying connection
	rw   bufio.ReadWriter // Buffered reader/writer
}

// ListEntry represents a per-device record, returned by [Client.List].
type ListEntry struct {
	// Path is the path of the device on the host exporting
	// the USB device, e.g. "/sys/devices/pci0000:00/0000:00:1d.1/usb3/3-2"
	Path string

	// Busid is the Bus ID of the exported device, e.g. "3-2".
	Busid string

	// Bus and Dev are the bus and device number on the exporting host.
	Bus, Dev int

	// Descriptor is the USB device descriptor.
	Descriptor *usb.DeviceDescriptor
}

// Dial connects to the USBIP server.
func Dial(ctx context.Context, addr netip.AddrPort) (*Client, error) {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr.String())
	if err != nil {
		return nil, err
	}

	rw := bufio.ReadWriter{
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
	clnt := &Client{conn: conn, rw: rw}

	return clnt, nil
}

// Close closes the [Client].
func (clnt *Client) Close() error {
	return clnt.conn.Close()
}

// List returns list of USB devices, available on server.
func (clnt *Client) List() ([]ListEntry, error) {
	// Send the OP_REQ_DEVLIST request
	err := binary.Write(&clnt.rw, binary.BigEndian, uint16(protoVersion))
	if err == nil {
		err = binary.Write(&clnt.rw, binary.BigEndian, uint16(protoHsDevlistReq))
	}
	if err == nil {
		err = binary.Write(&clnt.rw, binary.BigEndian, uint32(0))
	}
	if err == nil {
		err = clnt.rw.Flush()
	}

	if err != nil {
		return nil, err
	}

	// Starting from now, interpret io.EOF as io.ErrUnexpectedEOF
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	// Read the response header
	var ver uint16
	var code uint16
	var status uint32
	var ndevices uint32

	err = binary.Read(&clnt.rw, binary.BigEndian, &ver)
	if err == nil {
		err = binary.Read(&clnt.rw, binary.BigEndian, &code)
	}
	if err == nil {
		err = binary.Read(&clnt.rw, binary.BigEndian, &status)
	}
	if err == nil {
		err = binary.Read(&clnt.rw, binary.BigEndian, &ndevices)
	}

	if err != nil {
		return nil, err
	}

	// Validate response header
	if ver != protoVersion {
		err := fmt.Errorf(
			"protocol version mismatch: expected 0x%4.4x, present 0x%4.4x",
			protoVersion, ver)
		return nil, err
	}

	if protoHsCode(code) != protoHsDevlistRep {
		err := fmt.Errorf(
			"operation code mismatch: expected 0x%4.4x, present 0x%4.4x",
			protoHsDevlistRep, code)
		return nil, err
	}

	if protoHsStatus(status) != protoHsOK {
		err := fmt.Errorf(
			"status mismatch: expected 0x%8.8x, present 0x%8.8x",
			protoHsOK, status)
		return nil, err
	}

	if ndevices == 0 {
		return nil, nil
	}

	// Read device descriptors, one by one
	devlist := make([]ListEntry, ndevices)
	for i := uint32(0); i < ndevices; i++ {
		var path devPath
		var busid devBusID
		var busnum uint32
		var devnum uint32

		_, err = io.ReadFull(&clnt.rw, path[:])
		if err == nil {
			_, err = io.ReadFull(&clnt.rw, busid[:])
		}
		if err == nil {
			err = binary.Read(&clnt.rw, binary.BigEndian, &busnum)
		}
		if err == nil {
			err = binary.Read(&clnt.rw, binary.BigEndian, &devnum)
		}
	}

	return devlist, nil
}
