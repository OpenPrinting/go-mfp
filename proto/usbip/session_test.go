// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// USB/IP session lifecycle tests (detach / re-attach)

package usbip

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// ippUsbPrinter creates a printer with one IPP-over-USB (7/1/4)
// interface. dev.endpoints: [0]=nil, [1]=bulk IN, [2]=bulk OUT.
func ippUsbPrinter(t *testing.T, ctx context.Context) *Device {
	desc := usb.DeviceDescriptor{
		BCDUSB: 0x0200, Speed: usb.SpeedHigh, BMaxPacketSize: 64,
		Configurations: []usb.ConfigurationDescriptor{{
			BConfigurationValue: 1,
			Interfaces: []usb.Interface{{AltSettings: []usb.InterfaceDescriptor{{
				BInterfaceClass: 7, BInterfaceSubClass: 1, BInterfaceProtocol: 4,
				Endpoints: []usb.EndpointDescriptor{
					{Type: usb.EndpointIn, BMAttributes: usb.XferBulk, WMaxPacketSize: 512},
					{Type: usb.EndpointOut, BMAttributes: usb.XferBulk, WMaxPacketSize: 512},
				},
			}}}},
		}},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "0")
		w.Header().Set("X-Path", r.URL.Path)
	})

	dev, err := NewPrinter(ctx, desc, handler, nil)
	if err != nil {
		t.Fatal(err)
	}
	return dev
}

// attach connects to srv over net.Pipe and imports "1-1",
// exactly like `usbip attach -b 1-1`.
func attach(t *testing.T, srv *Server) net.Conn {
	cli, s := net.Pipe()
	go srv.serve(s)

	var req [40]byte
	binary.BigEndian.PutUint16(req[0:], 0x0111)
	binary.BigEndian.PutUint16(req[2:], 0x8003) // OP_REQ_IMPORT
	copy(req[8:], "1-1")
	if _, err := cli.Write(req[:]); err != nil {
		t.Fatal(err)
	}

	var rep [8]byte
	if _, err := io.ReadFull(cli, rep[:]); err != nil {
		t.Fatal(err)
	}
	if st := binary.BigEndian.Uint32(rep[4:]); st != 0 {
		t.Fatalf("OP_REP_IMPORT status %d", st)
	}
	var body [312]byte
	if _, err := io.ReadFull(cli, body[:]); err != nil {
		t.Fatal(err)
	}
	return cli
}

// readRetSubmit reads one USBIP_RET_SUBMIT, including IN data.
func readRetSubmit(t *testing.T, c net.Conn, in bool) (seq uint32, data []byte) {
	var b [48]byte
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.ReadFull(c, b[:]); err != nil {
		t.Fatalf("read RET_SUBMIT: %v", err)
	}
	if cmd := binary.BigEndian.Uint32(b[0:]); cmd != 3 {
		t.Fatalf("expected RET_SUBMIT, got cmd %d", cmd)
	}
	seq = binary.BigEndian.Uint32(b[4:])
	if in {
		data = make([]byte, binary.BigEndian.Uint32(b[24:]))
		if _, err := io.ReadFull(c, data); err != nil {
			t.Fatalf("read RET_SUBMIT data: %v", err)
		}
	}
	return
}

// TestServerReattachGetsFreshHTTPState: the host goes away in the
// middle of an IPP-over-USB HTTP request (usbip detach, TCP drop,
// host reboot, ipp-usb killed), then a host attaches again. The new
// session must not inherit HTTP state of the old one, like a real
// printer after USB disconnect/reset.
func TestServerReattachGetsFreshHTTPState(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := NewServer(ctx)
	if err := srv.AddDevice(ippUsbPrinter(t, ctx)); err != nil {
		t.Fatal(err)
	}

	// Session 1: request headers are only partially sent
	// when the host disappears.
	c1 := attach(t, srv)
	partial := "POST /ipp/print HTTP/1.1\r\nHost: localhost\r\n"
	go c1.Write(cmdSubmit(1, false, 2, uint32(len(partial)), []byte(partial)))
	readRetSubmit(t, c1, false)
	c1.Close()
	time.Sleep(100 * time.Millisecond)

	// Session 2: a fresh, complete request.
	c2 := attach(t, srv)
	defer c2.Close()
	req := "GET /second HTTP/1.1\r\nHost: localhost\r\n\r\n"
	go func() {
		c2.Write(cmdSubmit(1, false, 2, uint32(len(req)), []byte(req)))
		c2.Write(cmdSubmit(2, true, 1, 16384, nil))
	}()

	var rsp string
	for i := 0; i < 2; i++ {
		var b [48]byte
		c2.SetReadDeadline(time.Now().Add(5 * time.Second))
		if _, err := io.ReadFull(c2, b[:]); err != nil {
			t.Fatalf("read RET_SUBMIT: %v", err)
		}
		if binary.BigEndian.Uint32(b[4:]) == 2 {
			data := make([]byte, binary.BigEndian.Uint32(b[24:]))
			io.ReadFull(c2, data)
			rsp = string(data)
		}
	}

	status, _, _ := strings.Cut(rsp, "\r\n")
	if !strings.Contains(rsp, "X-Path: /second") {
		t.Fatalf("session 2 got a response that is not for its request:\n"+
			"  status line: %q\n  full response: %q", status, rsp)
	}
}

// TestServerReattachNoStaleResponse: the host detaches after sending a
// complete request but before posting the bulk-IN URB for the answer.
// The next session must not receive that orphaned response.
func TestServerReattachNoStaleResponse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := NewServer(ctx)
	if err := srv.AddDevice(ippUsbPrinter(t, ctx)); err != nil {
		t.Fatal(err)
	}

	// Session 1: full request, host goes away before reading reply.
	c1 := attach(t, srv)
	req1 := "GET /first HTTP/1.1\r\nHost: localhost\r\n\r\n"
	go c1.Write(cmdSubmit(1, false, 2, uint32(len(req1)), []byte(req1)))
	readRetSubmit(t, c1, false)
	time.Sleep(50 * time.Millisecond) // handler runs, reply is pending
	c1.Close()
	time.Sleep(100 * time.Millisecond)

	// Session 2: its own request and a bulk-IN URB for the reply.
	c2 := attach(t, srv)
	defer c2.Close()
	req2 := "GET /second HTTP/1.1\r\nHost: localhost\r\n\r\n"
	go func() {
		c2.Write(cmdSubmit(1, false, 2, uint32(len(req2)), []byte(req2)))
		c2.Write(cmdSubmit(2, true, 1, 16384, nil))
	}()

	var rsp string
	for i := 0; i < 2; i++ {
		var b [48]byte
		c2.SetReadDeadline(time.Now().Add(5 * time.Second))
		if _, err := io.ReadFull(c2, b[:]); err != nil {
			t.Fatalf("read RET_SUBMIT: %v", err)
		}
		if binary.BigEndian.Uint32(b[4:]) == 2 {
			data := make([]byte, binary.BigEndian.Uint32(b[24:]))
			io.ReadFull(c2, data)
			rsp = string(data)
		}
	}

	if !strings.Contains(rsp, "X-Path: /second") {
		t.Fatalf("session 2 got a response that is not for its request:\n"+
			"  full response: %q", rsp)
	}
}

// TestEndpointConnCloseUnblocksWrite: EndpointConn.Close() must unblock
// a Write waiting for a bulk-IN URB, same as it does for Read.
func TestEndpointConnCloseUnblocksWrite(t *testing.T) {
	_, in, out := testDevice(t)
	conn := NewEndpointConn(EndpointPair{In: in, Out: out})

	errc := make(chan error, 1)
	go func() {
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		errc <- err
	}()
	time.Sleep(50 * time.Millisecond)
	conn.Close()

	select {
	case err := <-errc:
		if err != net.ErrClosed {
			t.Fatalf("Write: %v, want net.ErrClosed", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Write not unblocked by Close")
	}
}
