// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Coverage for Device.endpointsTree: tree construction, descriptor
// encoding, endpoint addressing, routing and lookup helpers.

package usbip

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

var (
	tBulkOut = usb.EndpointDescriptor{Type: usb.EndpointOut,
		BMAttributes: usb.XferBulk, WMaxPacketSize: 512}
	tBulkIn = usb.EndpointDescriptor{Type: usb.EndpointIn,
		BMAttributes: usb.XferBulk, WMaxPacketSize: 512}
	tIntrIn = usb.EndpointDescriptor{Type: usb.EndpointIn,
		BMAttributes: usb.XferInterrupt, WMaxPacketSize: 8}
)

// tAlt makes an InterfaceDescriptor
func tAlt(class, sub, proto uint8, eps ...usb.EndpointDescriptor) usb.InterfaceDescriptor {
	return usb.InterfaceDescriptor{
		BInterfaceClass:    class,
		BInterfaceSubClass: sub,
		BInterfaceProtocol: proto,
		Endpoints:          eps,
	}
}

// tIff makes an Interface out of alt settings
func tIff(alts ...usb.InterfaceDescriptor) usb.Interface {
	return usb.Interface{AltSettings: alts}
}

// tDev makes a DeviceDescriptor out of configurations,
// each given as a list of interfaces.
func tDev(confs ...[]usb.Interface) usb.DeviceDescriptor {
	desc := usb.DeviceDescriptor{
		BCDUSB: 0x0200, Speed: usb.SpeedHigh, BMaxPacketSize: 64,
	}
	for i, iffs := range confs {
		desc.Configurations = append(desc.Configurations,
			usb.ConfigurationDescriptor{
				BConfigurationValue: uint8(i + 1),
				Interfaces:          iffs,
			})
	}
	return desc
}

// endpointLayouts covers the layouts of modeling/examples/*.py
// plus a multi-configuration device.
var endpointLayouts = []struct {
	name string
	desc usb.DeviceDescriptor
}{
	{"simple", tDev([]usb.Interface{
		tIff(tAlt(7, 1, 2, tBulkIn, tBulkOut)),
	})},

	// HP-LaserJet-MFP-M426fdn.py: 18 endpoints, 11 IN, 7 OUT
	{"hp-m426fdn", tDev([]usb.Interface{
		tIff(tAlt(255, 2, 1, tBulkOut, tBulkIn, tIntrIn)),
		tIff(tAlt(7, 1, 2, tBulkOut, tBulkIn)),
		tIff(tAlt(255, 3, 1, tBulkOut, tBulkIn, tIntrIn)),
		tIff(tAlt(255, 4, 1, tBulkOut, tBulkIn, tIntrIn)),
		tIff(tAlt(255, 4, 1, tBulkOut, tBulkIn, tIntrIn)),
		tIff(tAlt(255, 9, 1, tBulkOut, tBulkIn)),
		tIff(tAlt(255, 9, 1, tBulkOut, tBulkIn)),
	})},

	// Kyocera-ECOSYS-M2040dn.py: alt settings, one alt without endpoints
	{"kyocera-m2040dn", tDev([]usb.Interface{
		tIff(tAlt(7, 1, 2, tBulkOut, tBulkIn),
			tAlt(7, 1, 4, tBulkOut, tBulkIn)),
		tIff(tAlt(255, 255, 255, tBulkOut, tBulkIn),
			tAlt(7, 1, 4, tBulkOut, tBulkIn)),
		tIff(tAlt(255, 255, 255),
			tAlt(7, 1, 4, tBulkOut, tBulkIn)),
	})},

	// Xerox-B235.py: interrupt-only interface in the middle
	{"xerox-b235", tDev([]usb.Interface{
		tIff(tAlt(255, 255, 255, tBulkOut, tBulkIn)),
		tIff(tAlt(7, 1, 2, tBulkOut, tBulkIn)),
		tIff(tAlt(3, 0, 0, tIntrIn)),
		tIff(tAlt(7, 1, 4, tBulkOut, tBulkIn)),
		tIff(tAlt(7, 1, 4, tBulkOut, tBulkIn)),
	})},

	{"two-configs", tDev(
		[]usb.Interface{
			tIff(tAlt(7, 1, 2, tBulkIn, tBulkOut)),
		},
		[]usb.Interface{
			tIff(tAlt(255, 0, 0, tBulkOut, tBulkIn, tIntrIn)),
			tIff(tAlt(7, 1, 4, tBulkOut, tBulkIn)),
		},
	)},
}

// forEachLayout runs f as a subtest for every layout.
func forEachLayout(t *testing.T, f func(t *testing.T, desc usb.DeviceDescriptor)) {
	for _, l := range endpointLayouts {
		t.Run(l.name, func(t *testing.T) { f(t, l.desc) })
	}
}

// TestEndpointsTreeMirrorsDescriptor: endpointsTree has exactly one
// distinct Endpoint per EndpointDescriptor, at the same position.
// In particular, IN and OUT endpoints sharing an endpoint number
// never replace each other.
func TestEndpointsTreeMirrorsDescriptor(t *testing.T) {
	forEachLayout(t, func(t *testing.T, desc usb.DeviceDescriptor) {
		dev := MustNewDevice(desc)
		seen := make(map[*Endpoint]string)

		if len(dev.endpointsTree) != len(desc.Configurations) {
			t.Fatalf("configurations: %d, want %d",
				len(dev.endpointsTree), len(desc.Configurations))
		}

		for c, conf := range desc.Configurations {
			if len(dev.endpointsTree[c]) != len(conf.Interfaces) {
				t.Fatalf("conf %d: interfaces: %d, want %d", c,
					len(dev.endpointsTree[c]), len(conf.Interfaces))
			}
			for i, iff := range conf.Interfaces {
				if len(dev.endpointsTree[c][i]) != len(iff.AltSettings) {
					t.Fatalf("conf %d iff %d: alts: %d, want %d", c, i,
						len(dev.endpointsTree[c][i]), len(iff.AltSettings))
				}
				for a, alt := range iff.AltSettings {
					eps := dev.endpointsTree[c][i][a]
					if len(eps) != len(alt.Endpoints) {
						t.Fatalf("conf %d iff %d alt %d: endpoints: %d, want %d",
							c, i, a, len(eps), len(alt.Endpoints))
					}
					for e, ed := range alt.Endpoints {
						where := fmt.Sprintf("conf %d iff %d alt %d ep %d", c, i, a, e)
						ep := eps[e]
						switch {
						case ep == nil:
							t.Errorf("%s: nil", where)
						case seen[ep] != "":
							t.Errorf("%s: same Endpoint as %s", where, seen[ep])
						case ep.Type() != ed.Type ||
							ep.Attrs() != ed.BMAttributes ||
							ep.PktSize() != int(ed.WMaxPacketSize):
							t.Errorf("%s: doesn't match its descriptor", where)
						}
						seen[ep] = where
					}
				}
			}
		}
	})
}

// TestEndpointsLookupHelpers: EndpointsByFunc, EndpointsByClassID and
// EndpointsByClass return the endpointsTree entries of matching alts.
func TestEndpointsLookupHelpers(t *testing.T) {
	forEachLayout(t, func(t *testing.T, desc usb.DeviceDescriptor) {
		dev := MustNewDevice(desc)
		id := usb.ClassID{Class: 7, SubClass: 1, Protocol: 4}

		var want [][]*Endpoint
		var flat []*Endpoint
		for c, conf := range desc.Configurations {
			for i, iff := range conf.Interfaces {
				for a, alt := range iff.AltSettings {
					if alt.Match(id) {
						eps := dev.endpointsTree[c][i][a]
						want = append(want, eps)
						flat = append(flat, eps...)
					}
				}
			}
		}

		got := dev.EndpointsByClassID(id)
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("EndpointsByClassID: %v, want %v", got, want)
		}

		gotFlat := dev.EndpointsByClass(7, 1, 4)
		if fmt.Sprint(gotFlat) != fmt.Sprint(flat) {
			t.Errorf("EndpointsByClass: %v, want %v", gotFlat, flat)
		}

		for _, eps := range got {
			in, out := EndpointsInOut(eps)
			if in == nil || in.Type() != usb.EndpointIn ||
				out == nil || out.Type() != usb.EndpointOut {
				t.Errorf("EndpointsInOut(%v): %v, %v", eps, in, out)
			}
		}
	})
}

// wireRet is the USBIP_RET_SUBMIT or USBIP_RET_UNLINK
type wireRet struct {
	cmd, seq uint32
	status   syscall.Errno
	data     []byte
}

// wireHost is the minimal vhci-like USB/IP host.
type wireHost struct {
	t    *testing.T
	c    net.Conn
	seq  uint32
	lock sync.Mutex
	in   map[uint32]bool // IN requests by seqnum, to know if data follows
	rets chan wireRet
}

// newWireHost attaches to a new Server, serving the device.
func newWireHost(t *testing.T, dev *Device) *wireHost {
	h := &wireHost{
		t:    t,
		c:    startServer(t, dev),
		in:   make(map[uint32]bool),
		rets: make(chan wireRet, 256),
	}
	t.Cleanup(func() { h.c.Close() })
	go h.reader()
	return h
}

// reader receives RET_SUBMIT/RET_UNLINK messages
func (h *wireHost) reader() {
	for {
		var b [48]byte
		if _, err := io.ReadFull(h.c, b[:]); err != nil {
			close(h.rets)
			return
		}
		r := wireRet{
			cmd:    binary.BigEndian.Uint32(b[0:]),
			seq:    binary.BigEndian.Uint32(b[4:]),
			status: syscall.Errno(-int32(binary.BigEndian.Uint32(b[20:]))),
		}
		h.lock.Lock()
		in := h.in[r.seq]
		h.lock.Unlock()
		if r.cmd == 3 && in {
			r.data = make([]byte, binary.BigEndian.Uint32(b[24:]))
			if _, err := io.ReadFull(h.c, r.data); err != nil {
				close(h.rets)
				return
			}
		}
		h.rets <- r
	}
}

// send sends USBIP_CMD_SUBMIT and returns its seqnum.
func (h *wireHost) submit(in bool, ep uint8, setup []byte,
	length uint32, data []byte) uint32 {

	h.lock.Lock()
	h.seq++
	seq := h.seq
	h.in[seq] = in
	h.lock.Unlock()

	b := cmdSubmit(seq, in, uint32(ep), length, data)
	copy(b[40:48], setup)
	if _, err := h.c.Write(b); err != nil {
		h.t.Fatalf("write CMD_SUBMIT: %v", err)
	}
	return seq
}

// unlink sends USBIP_CMD_UNLINK and returns its seqnum.
func (h *wireHost) unlink(target uint32) uint32 {
	h.lock.Lock()
	h.seq++
	seq := h.seq
	h.lock.Unlock()

	if _, err := h.c.Write(cmdUnlink(seq, target)); err != nil {
		h.t.Fatalf("write CMD_UNLINK: %v", err)
	}
	return seq
}

// wait waits for RET with the specified seqnum.
func (h *wireHost) wait(seq uint32) wireRet {
	timeout := time.After(5 * time.Second)
	for {
		select {
		case r, ok := <-h.rets:
			if !ok {
				h.t.Fatalf("connection closed while waiting for #%d", seq)
			}
			if r.seq == seq {
				return r
			}
		case <-timeout:
			h.t.Fatalf("timeout waiting for #%d", seq)
		}
	}
}

// control performs a control transfer on ep0.
func (h *wireHost) control(reqType, req uint8, value, index,
	length uint16) wireRet {

	setup := make([]byte, 8)
	setup[0], setup[1] = reqType, req
	binary.LittleEndian.PutUint16(setup[2:], value)
	binary.LittleEndian.PutUint16(setup[4:], index)
	binary.LittleEndian.PutUint16(setup[6:], length)

	in := reqType&0x80 != 0
	return h.wait(h.submit(in, 0, setup, uint32(length), nil))
}

// wireEndpoint is the endpoint, as parsed from the configuration
// descriptor received over the wire.
type wireEndpoint struct {
	iff, alt, idx int
	addr, attrs   uint8
	pktsize       uint16
}

// getConfig fetches and validates the configuration descriptor
// with index confno, and returns its endpoints.
func (h *wireHost) getConfig(conf usb.ConfigurationDescriptor,
	confno int) []wireEndpoint {

	t := h.t
	r := h.control(0x80, 6, 2<<8|uint16(confno), 0, 4096)
	if r.status != 0 {
		t.Fatalf("GET_DESCRIPTOR(CONFIGURATION, %d): %v", confno, r.status)
	}

	d := r.data
	if len(d) < 9 || d[1] != 2 {
		t.Fatalf("not a configuration descriptor: % x", d)
	}
	if total := int(binary.LittleEndian.Uint16(d[2:])); total != len(d) {
		t.Errorf("wTotalLength %d, received %d bytes", total, len(d))
	}
	if int(d[4]) != len(conf.Interfaces) {
		t.Errorf("bNumInterfaces %d, want %d", d[4], len(conf.Interfaces))
	}
	if int(d[5]) != confno+1 {
		t.Errorf("bConfigurationValue %d, want %d", d[5], confno+1)
	}

	var eps []wireEndpoint
	iff, alt, idx, numEps := -1, -1, 0, -1
	checkCount := func() {
		if numEps >= 0 && idx != numEps {
			t.Errorf("interface %d alt %d: bNumEndpoints %d, "+
				"but %d endpoint descriptors follow",
				iff, alt, numEps, idx)
		}
	}

	for off := int(d[0]); off < len(d); off += int(d[off]) {
		if d[off] < 2 || off+int(d[off]) > len(d) {
			t.Fatalf("bad bLength %d at offset %d", d[off], off)
		}
		switch d[off+1] {
		case 4: // INTERFACE
			checkCount()
			iff, alt, idx, numEps = int(d[off+2]), int(d[off+3]), 0, int(d[off+4])
			want := conf.Interfaces[iff].AltSettings[alt]
			if d[off+5] != want.BInterfaceClass ||
				d[off+6] != want.BInterfaceSubClass ||
				d[off+7] != want.BInterfaceProtocol {
				t.Errorf("interface %d alt %d: class %d/%d/%d", iff, alt,
					d[off+5], d[off+6], d[off+7])
			}
		case 5: // ENDPOINT
			eps = append(eps, wireEndpoint{iff, alt, idx,
				d[off+2], d[off+3],
				binary.LittleEndian.Uint16(d[off+4:])})
			idx++
		}
	}
	checkCount()

	return eps
}

// TestEndpointsDescriptorOnWire: configuration descriptors, as seen by
// the USB host, are well-formed, list every endpoint of every alt
// setting exactly once, and use valid, unique endpoint addresses.
func TestEndpointsDescriptorOnWire(t *testing.T) {
	forEachLayout(t, func(t *testing.T, desc usb.DeviceDescriptor) {
		dev := MustNewDevice(desc)
		h := newWireHost(t, dev)

		r := h.control(0x80, 6, 1<<8, 0, 18)
		if r.status != 0 || len(r.data) != 18 {
			t.Fatalf("GET_DESCRIPTOR(DEVICE): %v, %d bytes", r.status, len(r.data))
		}
		if int(r.data[17]) != len(desc.Configurations) {
			t.Errorf("bNumConfigurations %d, want %d",
				r.data[17], len(desc.Configurations))
		}

		for c, conf := range desc.Configurations {
			eps := h.getConfig(conf, c)
			seen := make(map[uint8]wireEndpoint)
			cnt := 0

			for _, ep := range eps {
				want := conf.Interfaces[ep.iff].AltSettings[ep.alt].Endpoints[ep.idx]
				where := fmt.Sprintf("conf %d iff %d alt %d ep %d",
					c, ep.iff, ep.alt, ep.idx)

				num := ep.addr & 0x0f
				in := ep.addr&0x80 != 0
				switch {
				case ep.addr&0x70 != 0 || num == 0:
					t.Errorf("%s: invalid bEndpointAddress 0x%2.2x", where, ep.addr)
				case in != (want.Type == usb.EndpointIn):
					t.Errorf("%s: bEndpointAddress 0x%2.2x: wrong direction",
						where, ep.addr)
				case usb.EndpointAttributes(ep.attrs) != want.BMAttributes ||
					ep.pktsize != want.WMaxPacketSize:
					t.Errorf("%s: bmAttributes/wMaxPacketSize mismatch", where)
				}

				if prev, dup := seen[ep.addr]; dup {
					t.Errorf("%s: bEndpointAddress 0x%2.2x already used "+
						"by iff %d alt %d", where, ep.addr, prev.iff, prev.alt)
				}
				seen[ep.addr] = ep
				cnt++
			}

			want := 0
			for _, iff := range conf.Interfaces {
				for _, alt := range iff.AltSettings {
					want += len(alt.Endpoints)
				}
			}
			if cnt != want {
				t.Errorf("conf %d: %d endpoint descriptors, want %d", c, cnt, want)
			}
		}

		// Non-existent configuration must stall, not crash
		n := len(desc.Configurations)
		r = h.control(0x80, 6, 2<<8|uint16(n), 0, 4096)
		if r.status != syscall.EPIPE {
			t.Errorf("GET_DESCRIPTOR(CONFIGURATION, %d): %v, want EPIPE", n, r.status)
		}
	})
}

// TestEndpointsRoutingOnWire: in every configuration, a bulk/interrupt
// URB submitted to the address advertised in the descriptor reaches
// the corresponding endpointsTree entry, and CMD_UNLINK removes it
// from there. Endpoints of opposite directions that share a number
// are routed independently.
func TestEndpointsRoutingOnWire(t *testing.T) {
	forEachLayout(t, func(t *testing.T, desc usb.DeviceDescriptor) {
		dev := MustNewDevice(desc)
		h := newWireHost(t, dev)

		for c, conf := range desc.Configurations {
			eps := h.getConfig(conf, c)

			r := h.control(0x00, 9, uint16(c+1), 0, 0) // SET_CONFIGURATION
			if r.status != 0 {
				t.Fatalf("SET_CONFIGURATION(%d): %v", c+1, r.status)
			}

			// Queue one URB on every endpoint, all at once
			seqs := make([]uint32, len(eps))
			for i, ep := range eps {
				in := ep.addr&0x80 != 0
				var data []byte
				if !in {
					data = []byte("DATA")
				}
				seqs[i] = h.submit(in, ep.addr&0x0f, nil, 4, data)
			}
			time.Sleep(20 * time.Millisecond)

			for i, ep := range eps {
				want := dev.endpointsTree[c][ep.iff][ep.alt][ep.idx]
				if got := findQueued(dev, seqs[i]); got != want {
					t.Errorf("conf %d iff %d alt %d ep 0x%2.2x: URB routed to %s",
						c, ep.iff, ep.alt, ep.addr, describeEndpoint(dev, got))
				}
			}

			// Unlink them; each must be found where it was queued
			for i, ep := range eps {
				r := h.wait(h.unlink(seqs[i]))
				if r.status != syscall.ECONNRESET {
					t.Errorf("conf %d iff %d alt %d ep 0x%2.2x: "+
						"RET_UNLINK %v, want ECONNRESET",
						c, ep.iff, ep.alt, ep.addr, r.status)
				}
				if got := findQueued(dev, seqs[i]); got != nil {
					t.Errorf("conf %d ep 0x%2.2x: URB still queued after unlink",
						c, ep.addr)
				}
			}
		}

		// Endpoint numbers that don't exist are rejected
		r := h.wait(h.submit(false, 15, nil, 4, []byte("DATA")))
		if r.status != syscall.EPIPE {
			t.Errorf("CMD_SUBMIT to OUT 15: %v, want EPIPE", r.status)
		}
		r = h.wait(h.submit(true, 200, nil, 4, nil))
		if r.status != syscall.EPIPE {
			t.Errorf("CMD_SUBMIT to IN 200: %v, want EPIPE", r.status)
		}
	})
}
