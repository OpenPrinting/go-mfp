// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// Configuration descriptor and endpoint addressing tests

package usbip

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// hpM426fdnDevice mirrors the USB layout of
// modeling/examples/HP-LaserJet-MFP-M426fdn.py: 7 interfaces,
// 18 endpoints (11 IN, 7 OUT), as reported by the real device.
func hpM426fdnDevice(t *testing.T) *Device {
	bulkOut := usb.EndpointDescriptor{Type: usb.EndpointOut,
		BMAttributes: usb.XferBulk, WMaxPacketSize: 512}
	bulkIn := usb.EndpointDescriptor{Type: usb.EndpointIn,
		BMAttributes: usb.XferBulk, WMaxPacketSize: 512}
	intrIn := usb.EndpointDescriptor{Type: usb.EndpointIn,
		BMAttributes: usb.XferInterrupt, WMaxPacketSize: 8}

	iff := func(class, sub, proto uint8, name string,
		eps ...usb.EndpointDescriptor) usb.Interface {
		return usb.Interface{AltSettings: []usb.InterfaceDescriptor{{
			BInterfaceClass:    class,
			BInterfaceSubClass: sub,
			BInterfaceProtocol: proto,
			IInterface:         name,
			Endpoints:          eps,
		}}}
	}

	desc := usb.DeviceDescriptor{
		BCDUSB: 0x0200, Speed: usb.SpeedHigh, BMaxPacketSize: 64,
		IDVendor: 0x03f0, IDProduct: 0x5a2a,
		Configurations: []usb.ConfigurationDescriptor{{
			BConfigurationValue: 1,
			Interfaces: []usb.Interface{
				iff(255, 2, 1, "HP SCAN", bulkOut, bulkIn, intrIn),
				iff(7, 1, 2, "Printer", bulkOut, bulkIn),
				iff(255, 3, 1, "HP FAX", bulkOut, bulkIn, intrIn),
				iff(255, 4, 1, "HP LEDM", bulkOut, bulkIn, intrIn),
				iff(255, 4, 1, "HP LEDM", bulkOut, bulkIn, intrIn),
				iff(255, 9, 1, "IPP Printer", bulkOut, bulkIn),
				iff(255, 9, 1, "IPP Printer", bulkOut, bulkIn),
			},
		}},
	}

	dev, err := NewDevice(desc)
	if err != nil {
		t.Fatal(err)
	}
	return dev
}

// epDesc is the endpoint descriptor, as seen by the USB host.
type epDesc struct {
	iff, alt, idx int   // Interface, alt setting, index within alt
	addr          uint8 // bEndpointAddress
}

// getConfigDescriptor sends GET_DESCRIPTOR(CONFIGURATION, 0) on ep0,
// exactly like the Linux USB core does when the device is attached,
// and returns all endpoint descriptors it contains.
func getConfigDescriptor(t *testing.T, c net.Conn) []epDesc {
	b := make([]byte, 48)
	binary.BigEndian.PutUint32(b[0:], 1) // USBIP_CMD_SUBMIT
	binary.BigEndian.PutUint32(b[4:], 1000)
	binary.BigEndian.PutUint32(b[8:], testDevID)
	binary.BigEndian.PutUint32(b[12:], 1)      // USBIP_DIR_IN
	binary.BigEndian.PutUint32(b[20:], 0x0200) // URB_DIR_IN
	binary.BigEndian.PutUint32(b[24:], 1024)
	setup := b[40:48]
	setup[0] = 0x80                                // IN, standard, device
	setup[1] = 6                                   // GET_DESCRIPTOR
	binary.LittleEndian.PutUint16(setup[2:], 2<<8) // CONFIGURATION, index 0
	binary.LittleEndian.PutUint16(setup[6:], 1024)

	go c.Write(b)

	var rsp [48]byte
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.ReadFull(c, rsp[:]); err != nil {
		t.Fatalf("read RET_SUBMIT: %v", err)
	}
	if st := int32(binary.BigEndian.Uint32(rsp[20:])); st != 0 {
		t.Fatalf("GET_DESCRIPTOR(CONFIGURATION) status %d", st)
	}
	data := make([]byte, binary.BigEndian.Uint32(rsp[24:]))
	if _, err := io.ReadFull(c, data); err != nil {
		t.Fatalf("read descriptor: %v", err)
	}

	var eps []epDesc
	iff, alt, idx := -1, -1, 0
	for off := 0; off+2 <= len(data); off += int(data[off]) {
		if data[off] == 0 {
			t.Fatalf("zero bLength at offset %d", off)
		}
		switch data[off+1] {
		case 4: // INTERFACE
			iff, alt, idx = int(data[off+2]), int(data[off+3]), 0
		case 5: // ENDPOINT
			eps = append(eps, epDesc{iff, alt, idx, data[off+2]})
			idx++
		}
	}
	return eps
}

// TestServerConfigDescriptorEndpointAddresses: every bEndpointAddress
// must be a valid USB 2.0 address: number in 1...15, bits 6..4
// reserved (zero), and unique per direction.
func TestServerConfigDescriptorEndpointAddresses(t *testing.T) {
	cli := startServer(t, hpM426fdnDevice(t))
	defer cli.Close()

	seen := make(map[uint8]epDesc)
	for _, ep := range getConfigDescriptor(t, cli) {
		num := ep.addr & 0x0f
		if ep.addr&0x70 != 0 || num == 0 {
			t.Errorf("interface %d ep %d: invalid bEndpointAddress 0x%2.2x "+
				"(number %d does not fit 4 bits)",
				ep.iff, ep.idx, ep.addr, ep.addr&0x7f)
			continue
		}
		if prev, dup := seen[ep.addr]; dup {
			t.Errorf("interface %d ep %d: bEndpointAddress 0x%2.2x "+
				"already used by interface %d", ep.iff, ep.idx,
				ep.addr, prev.iff)
		}
		seen[ep.addr] = ep
	}
}

// TestServerBulkRoutingMatchesDescriptor: a URB that the host submits
// to the endpoint advertised in the descriptor must reach that
// interface's Endpoint. The USB/IP ep field comes from
// usb_pipeendpoint(), which is 4 bits wide.
func TestServerBulkRoutingMatchesDescriptor(t *testing.T) {
	dev := hpM426fdnDevice(t)
	cli := startServer(t, dev)
	defer cli.Close()

	eps := getConfigDescriptor(t, cli)

	// Drain RET_SUBMITs (EPIPE for misrouted URBs) in background.
	go io.Copy(io.Discard, cli)

	for seq, ep := range eps {
		want := dev.endpointsTree[0][ep.iff][ep.alt][ep.idx]
		in := ep.addr&0x80 != 0

		var data []byte
		if !in {
			data = []byte("DATA")
		}
		cli.Write(cmdSubmit(uint32(seq+1), in, uint32(ep.addr&0x0f), 4, data))
		time.Sleep(10 * time.Millisecond)

		if where := findQueued(dev, uint32(seq+1)); where != want {
			t.Errorf("interface %d (%s) ep 0x%2.2x: URB routed to %s",
				ep.iff, dev.Descriptor.Configurations[0].
					Interfaces[ep.iff].AltSettings[ep.alt].IInterface,
				ep.addr, describeEndpoint(dev, where))
		}
	}
}

// findQueued returns the Endpoint that holds the URB with given seqnum.
func findQueued(dev *Device, seq uint32) *Endpoint {
	for _, iff := range dev.endpointsTree[0] {
		for _, alt := range iff {
			for _, ep := range alt {
				ep.lock.Lock()
				q := append(append([]*protoIOSubmitRequest{},
					ep.inqueue...), ep.outqueue...)
				ep.lock.Unlock()
				for _, rq := range q {
					if rq.Seqnum == seq {
						return ep
					}
				}
			}
		}
	}
	return nil
}

// describeEndpoint names the interface that owns the Endpoint.
func describeEndpoint(dev *Device, ep *Endpoint) string {
	if ep == nil {
		return "nowhere (rejected)"
	}
	conf := dev.Descriptor.Configurations[0]
	for iffno, iff := range dev.endpointsTree[0] {
		for altno, alt := range iff {
			for i, e := range alt {
				if e == ep {
					return fmt.Sprintf("interface %d (%s) ep %d",
						iffno, conf.Interfaces[iffno].
							AltSettings[altno].IInterface, i)
				}
			}
		}
	}
	return "unknown endpoint"
}

// TestNewDeviceTooManyEndpoints: USB allows at most 15 IN and 15 OUT
// endpoints per configuration; NewDevice must reject more instead of
// generating invalid endpoint addresses.
func TestNewDeviceTooManyEndpoints(t *testing.T) {
	alt := usb.InterfaceDescriptor{BInterfaceClass: 255}
	for i := 0; i < 16; i++ {
		alt.Endpoints = append(alt.Endpoints, usb.EndpointDescriptor{
			Type: usb.EndpointIn, BMAttributes: usb.XferBulk,
			WMaxPacketSize: 512})
	}

	desc := usb.DeviceDescriptor{
		BCDUSB: 0x0200, BMaxPacketSize: 64,
		Configurations: []usb.ConfigurationDescriptor{{
			BConfigurationValue: 1,
			Interfaces: []usb.Interface{
				{AltSettings: []usb.InterfaceDescriptor{alt}}},
		}},
	}

	if _, err := NewDevice(desc); err == nil {
		t.Fatal("NewDevice accepted 16 IN endpoints")
	}

	alt.Endpoints = alt.Endpoints[:15]
	desc.Configurations[0].Interfaces[0].AltSettings[0] = alt
	if _, err := NewDevice(desc); err != nil {
		t.Fatalf("NewDevice rejected 15 IN endpoints: %v", err)
	}
}
