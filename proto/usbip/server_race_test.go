// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Wire-level Server tests (vhci-style client over net.Pipe)

package usbip

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// Endpoint layout of testDevice: bulk IN 0x81, bulk OUT 0x01
func testDevice(t *testing.T) (*Device, *Endpoint, *Endpoint) {
	desc := usb.DeviceDescriptor{
		BCDUSB: 0x0200, Speed: usb.SpeedHigh, BMaxPacketSize: 64,
		Configurations: []usb.ConfigurationDescriptor{{
			BConfigurationValue: 1,
			Interfaces: []usb.Interface{{AltSettings: []usb.InterfaceDescriptor{{
				BInterfaceClass: 7, BInterfaceSubClass: 1, BInterfaceProtocol: 2,
				Endpoints: []usb.EndpointDescriptor{
					{Type: usb.EndpointIn, BMAttributes: usb.XferBulk, WMaxPacketSize: 512},
					{Type: usb.EndpointOut, BMAttributes: usb.XferBulk, WMaxPacketSize: 512},
				},
			}}}},
		}},
	}
	dev, err := NewDevice(desc)
	if err != nil {
		t.Fatal(err)
	}
	in, out := EndpointsInOut(dev.endpointsTree[0][0][0])
	return dev, in, out
}

const testDevID = 1<<16 | 1 // bus 1, dev 1 -> busid "1-1"

// startServer runs the real Server.serve() on one end of net.Pipe and
// performs OP_REQ_IMPORT of "1-1" on the other, like `usbip attach`.
func startServer(t *testing.T, dev *Device) net.Conn {
	srv := NewServer(context.Background())
	if err := srv.AddDevice(dev); err != nil {
		t.Fatal(err)
	}
	cli, s := net.Pipe()
	go srv.serve(s)

	var req [40]byte
	binary.BigEndian.PutUint16(req[0:], 0x0111)
	binary.BigEndian.PutUint16(req[2:], 0x8003) // OP_REQ_IMPORT
	copy(req[8:], "1-1")
	if _, err := cli.Write(req[:]); err != nil {
		t.Fatal(err)
	}
	var rep [320]byte // 8 hdr + 312 usbip_usb_device
	if _, err := io.ReadFull(cli, rep[:]); err != nil {
		t.Fatal(err)
	}
	if st := binary.BigEndian.Uint32(rep[4:]); st != 0 {
		t.Fatalf("OP_REP_IMPORT status %d", st)
	}
	return cli
}

func cmdSubmit(seq uint32, in bool, ep uint32, length uint32, data []byte) []byte {
	b := make([]byte, 48, 48+len(data))
	binary.BigEndian.PutUint32(b[0:], 1) // USBIP_CMD_SUBMIT
	binary.BigEndian.PutUint32(b[4:], seq)
	binary.BigEndian.PutUint32(b[8:], testDevID)
	flags := uint32(0)
	if in {
		binary.BigEndian.PutUint32(b[12:], 1)
		flags = 0x0200
	}
	binary.BigEndian.PutUint32(b[16:], ep)
	binary.BigEndian.PutUint32(b[20:], flags)
	binary.BigEndian.PutUint32(b[24:], length)
	return append(b, data...)
}

// cmdUnlink is encoded like vhci_tx.c: header memset to 0, then
// command/seqnum/devid and the target seqnum; direction=0, ep=0.
func cmdUnlink(seq, target uint32) []byte {
	b := make([]byte, 48)
	binary.BigEndian.PutUint32(b[0:], 2) // USBIP_CMD_UNLINK
	binary.BigEndian.PutUint32(b[4:], seq)
	binary.BigEndian.PutUint32(b[8:], testDevID)
	binary.BigEndian.PutUint32(b[20:], target)
	return b
}

type ret struct {
	cmd, seq uint32
	status   int32
	data     []byte
}

func readRet(t *testing.T, c net.Conn) ret {
	var b [48]byte
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.ReadFull(c, b[:]); err != nil {
		t.Fatalf("read RET: %v", err)
	}
	r := ret{
		cmd:    binary.BigEndian.Uint32(b[0:]),
		seq:    binary.BigEndian.Uint32(b[4:]),
		status: int32(binary.BigEndian.Uint32(b[20:])),
	}
	return r
}

// TestServerPipelinedBulkOut: vhci pipelines many bulk-OUT URBs
// (usblp/ipp-usb sending a print job) while the device-side consumer
// drains the OUT endpoint on its own goroutine.
func TestServerPipelinedBulkOut(t *testing.T) {
	dev, _, out := testDevice(t)
	cli := startServer(t, dev)
	defer cli.Close()

	const n = 20000
	go func() { // device-side consumer, like runLegacyPrinter
		buf := make([]byte, 512)
		for {
			if _, err := out.Read(buf); err != nil {
				return
			}
		}
	}()

	go func() {
		payload := make([]byte, 64)
		for i := uint32(1); i <= n; i++ {
			if _, err := cli.Write(cmdSubmit(i, false, 1, 64, payload)); err != nil {
				return
			}
		}
	}()

	for i := 0; i < n; i++ {
		r := readRet(t, cli)
		if r.cmd != 3 || r.status != 0 {
			t.Fatalf("bad RET: %+v", r)
		}
	}
}

// TestServerUnlinkPendingBulkIn: ipp-usb/usblp cancels a pending
// bulk-IN URB (timeout / close); vhci sends CMD_UNLINK.
func TestServerUnlinkPendingBulkIn(t *testing.T) {
	dev, in, _ := testDevice(t)
	cli := startServer(t, dev)
	defer cli.Close()

	go cli.Write(cmdSubmit(1, true, 1, 512, nil))
	time.Sleep(50 * time.Millisecond)
	go cli.Write(cmdUnlink(2, 1))

	r := readRet(t, cli)
	if r.cmd != 4 || r.seq != 2 {
		t.Fatalf("expected RET_UNLINK #2, got %+v", r)
	}
	if syscall.Errno(-r.status) != syscall.ECONNRESET {
		t.Errorf("RET_UNLINK status=%d, want -ECONNRESET(%d): URB was not unlinked",
			r.status, -int32(syscall.ECONNRESET))
	}

	// The device later produces data. A correctly unlinked URB must
	// not receive it, and no RET_SUBMIT for #1 may follow RET_UNLINK.
	done := make(chan struct{})
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		in.WriteContext(ctx, []byte("IPP-RESPONSE"))
		close(done)
	}()
	<-done
	cli.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	var b [48]byte
	if _, err := io.ReadFull(cli, b[:]); err == nil {
		t.Errorf("RET_SUBMIT for unlinked seqnum #%d sent after RET_UNLINK "+
			"(cmd=%d); %q consumed by a dead URB",
			binary.BigEndian.Uint32(b[4:]), binary.BigEndian.Uint32(b[0:]),
			"IPP-RESPONSE")
	}
}

// TestServerSubmitBadEndpoint: CMD_SUBMIT to a non-existent endpoint.
func TestServerSubmitBadEndpoint(t *testing.T) {
	dev, _, _ := testDevice(t)
	cli := startServer(t, dev)
	defer cli.Close()

	go cli.Write(cmdSubmit(1, false, 3, 4, []byte("ABCD")))
	r := readRet(t, cli)
	if r.cmd != 3 || syscall.Errno(-r.status) != syscall.EPIPE {
		t.Fatalf("want RET_SUBMIT -EPIPE, got %+v", r)
	}
}

// TestEndpointConnCloseUnblocksRead: EndpointConn.Close() is documented
// to "unblock all pending Reads and Writes" (used by the IPP-over-USB
// http.Server via EndpointListener).
func TestEndpointConnCloseUnblocksRead(t *testing.T) {
	_, in, out := testDevice(t)
	conn := NewEndpointConn(EndpointPair{In: in, Out: out})

	errc := make(chan error, 1)
	go func() {
		_, err := conn.Read(make([]byte, 16))
		errc <- err
	}()
	time.Sleep(50 * time.Millisecond)
	conn.Close()

	select {
	case err := <-errc:
		if err != net.ErrClosed {
			t.Fatalf("Read: %v, want net.ErrClosed", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Read not unblocked by Close")
	}
}
