// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by Mohammed Imaduddin (mdimad005@gmail.com)
// See LICENSE for license terms and conditions
//
// IPP over USB emulation logic

package usbip

import (
	"context"
	"errors"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/proto/usb"
	"github.com/OpenPrinting/go-mfp/util/generic"
)

// Server implements the USBIP server.
type Server struct {
	ctx        context.Context             // Server context
	devices    [usb.MaxDevices + 1]devslot // Device slots
	nextslot   int                         // Next free slot
	byid       map[devBusID]*devslot       // Slots by devBusID
	bylocation map[devLocation]*devslot    // Slots by devLocation
	lock       sync.Mutex                  // Access lock
}

// devslot represents a device slot.
type devslot struct {
	dev      *Device
	busid    devBusID
	location devLocation
	pconn    *protoConn
}

// NewServer creates a new [Server]
func NewServer(ctx context.Context) *Server {
	srv := &Server{
		ctx:        ctx,
		nextslot:   1,
		byid:       make(map[devBusID]*devslot, usb.MaxDevices+1),
		bylocation: make(map[devLocation]*devslot, usb.MaxDevices+1),
	}

	for i := range srv.devices {
		slot := &srv.devices[i]
		slot.location.Bus = 1 // Hardcoded for now
		slot.location.Dev = i
		slot.busid = slot.location.BusID()

		srv.byid[slot.busid] = slot
		srv.bylocation[slot.location] = slot
	}

	return srv
}

// AddDevice adds a USBDevice to the container.
func (srv *Server) AddDevice(dev *Device) error {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	// Look for the empty slot.
	for i := range srv.devices {
		next := (srv.nextslot + i) % len(srv.devices)
		if next != 0 && srv.devices[next].dev == nil {
			srv.devices[next].dev = dev
			srv.nextslot = (next + 1) % len(srv.devices)
			return nil
		}
	}

	return errors.New("Can't add new device. All slots are busy")
}

// doImport handles a device import request from the USB/IP client.
// It returns true if device was actually attached and the protocol
// response message in any case.
func (srv *Server) doImport(pconn *protoConn, rq *protoImportRequest) *protoImportResponse {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	// Find the device
	slot := srv.byid[rq.BusID]
	if slot == nil || slot.dev == nil || slot.pconn != nil {
		return &protoImportResponse{Status: protoHsError}
	}

	// Attach device to connection
	slot.pconn = pconn

	// Build the response
	rsp := &protoImportResponse{
		Status: protoHsOK,
		DevInfo: &devInfo{
			Location: slot.location,
			Device:   slot.dev,
		},
	}

	return rsp
}

// doDevlist responds with a list of available USB devices.
// It returns the protocol response message.
func (srv *Server) doDevlist() *protoDevlistResponse {
	// Gather connected but not yet exported slots
	srv.lock.Lock()

	slots := make([]devslot, 0, len(srv.devices))
	for _, slot := range srv.devices {
		if slot.dev != nil && slot.pconn == nil {
			slots = append(slots, slot)
		}
	}

	srv.lock.Unlock()

	// Generate a response
	rsp := &protoDevlistResponse{
		Status:  protoHsOK,
		Devlist: make([]devInfo, len(slots)),
	}

	for i := range slots {
		rsp.Devlist[i].Location = slots[i].location
		rsp.Devlist[i].Device = slots[i].dev
	}

	return rsp
}

// doDisconnect handles the client disconnect event.
func (srv *Server) doDisconnect(pconn *protoConn) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	for i := range srv.devices {
		if srv.devices[i].pconn == pconn {
			srv.devices[i].pconn = nil
		}
	}
}

// ListenAndServe listens on the TCP network address and
// then calls [Server.Serve] to handle incoming connections.
func (srv *Server) ListenAndServe(addr net.Addr) error {
	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	return srv.Serve(l)
}

// Serve accepts and serves incoming connections.
func (srv *Server) Serve(l net.Listener) error {
	const tempDelayMin = 5 * time.Millisecond
	const tempDelayMax = 100 * time.Millisecond

	tempDelay := tempDelayMin

	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				time.Sleep(tempDelay)
				tempDelay *= 2
				tempDelay = generic.Min(tempDelay, tempDelayMax)
				continue
			}
			return err
		}

		tempDelay = tempDelayMin
		go srv.serve(conn)
	}
}

// Serve manages a client connection.
func (srv *Server) serve(conn net.Conn) {
	log.Debug(srv.ctx, "Connection: from %s", conn.RemoteAddr())

	// Create protocol connection
	pconn := newProtoConn(conn)

	defer func() {
		srv.doDisconnect(pconn)
		log.Debug(srv.ctx, "Connection: closed")
		pconn.Close()
	}()

	devinfo := srv.handshake(pconn)
	if devinfo == nil {
		return
	}

	srv.io(pconn, devinfo)
}

// handshake performs the USBIP handshake exchange.
func (srv *Server) handshake(pconn *protoConn) *devInfo {
	rq, err := pconn.RecvHandshake()
	if err != nil {
		log.Debug(srv.ctx, "Handshake: %s", err)
		return nil
	}

	log.Debug(srv.ctx, "Handshake: < %s", rq)

	var rsp protoHandshakeResponse

	switch rq := rq.(type) {
	case *protoDevlistRequest:
		rsp = srv.doDevlist()
	case *protoImportRequest:
		rsp = srv.doImport(pconn, rq)
	}

	if rsp != nil {
		log.Debug(srv.ctx, "Handshake: > %s", rsp)
		err = pconn.SendHandshake(rsp)
		if err != nil {
			log.Debug(srv.ctx, "Handshake: %s", err)
			return nil
		}
	}

	if rsp, ok := rsp.(*protoImportResponse); ok && rsp.Status == protoHsOK {
		return rsp.DevInfo
	}

	return nil
}

// io performs USBIO I/O exchange.
func (srv *Server) io(pconn *protoConn, devinfo *devInfo) {
	dev := devinfo.Device
	defer dev.shutdown()

	// Map of the pending protoIOSubmitRequest by Seqnum.
	//
	// We need it, because protoIOUnlinkRequest lacks some
	// important information (endpoint number, for example).
	submitsBySeqnum := make(map[uint32]*protoIOSubmitRequest)

	// Completion callback for protoIOSubmitRequest
	completion := func(rq *protoIOSubmitRequest, err syscall.Errno) {
		delete(submitsBySeqnum, rq.Seqnum)

		rsp := rq.Response(err)
		log.Debug(srv.ctx, "IO:   > %s", rsp)
		pconn.SendIO(rsp)
	}

	for {
		rq, err := pconn.RecvIO()
		if err != nil {
			log.Debug(srv.ctx, "IO:   < %s", err)
			return
		}

		log.Debug(srv.ctx, "IO:   < %s", rq)

		hdr := rq.Header()
		if hdr.Location != devinfo.Location {
			log.Debug(srv.ctx,
				"IO:   < Device: expected %s, present %s",
				devinfo.Location, hdr.Location)

			// FIXME - reject request and continue
			return
		}

		switch rq := rq.(type) {
		case *protoIOSubmitRequest:
			// Control requests are processed instantly.
			if rq.Endpoint == 0 {
				data, err := srv.control(dev, rq)

				rq.Buffer = data
				rq.actualLength = len(data)

				rsp := rq.Response(err)
				pconn.SendIO(rsp)
				break
			}

			submitsBySeqnum[rq.Seqnum] = rq

			rq.completion = completion
			err := dev.submit(rq)
			if err != 0 {
				rsp := rq.Response(err)
				pconn.SendIO(rsp)
			}

		case *protoIOUnlinkRequest:
			err := syscall.Errno(0)

			submit := submitsBySeqnum[rq.Seqnum]
			if submit != nil {
				// This information is missed in
				// the USBIP_CMD_SUBMIT.
				rq.Endpoint = submit.Endpoint

				err = dev.unlink(rq)
			}

			rsp := rq.Response(err)
			pconn.SendIO(rsp)
		}
	}
}

// control handles USB control requests.
func (srv *Server) control(dev *Device,
	rq *protoIOSubmitRequest) ([]byte, syscall.Errno) {

	var setup usb.SetupPacket
	setup.Decode(rq.Setup)

	log.Debug(srv.ctx, "CTRL: < %s", setup)

	var data []byte
	var err syscall.Errno

	switch setup.RequestType & usb.RequestTypeTypeMask {
	case usb.RequestTypeStandard:
		data, err = srv.controlStandard(dev, setup)

	case usb.RequestTypeClass:
		data, err = srv.controlClassSpecific(dev, setup)

	default:
		err = syscall.EPIPE
	}

	if data != nil {
		log.Debug(srv.ctx, "CTRL: > %d bytes returned", len(data))
	} else {
		log.Debug(srv.ctx, "CTRL: > %s", err)
	}

	return data, err
}

// classSpecific handles USB standard control requests.
func (srv *Server) controlStandard(dev *Device,
	setup usb.SetupPacket) ([]byte, syscall.Errno) {

	switch setup.RequestType & usb.RecipientMask {
	case usb.RecipientDevice:
		switch setup.Request {
		case usb.RequestGetStatus:
			return dev.getStatus()

		case usb.RequestGetDescriptor:
			t := usb.DescriptorType(setup.WValue >> 8)
			i := int(setup.WValue & 255)
			return dev.getDescriptor(t, i)

		case usb.RequestGetConfiguration:
			return dev.GetConfiguration()

		case usb.RequestSetConfiguration:
			n := setup.WValue
			return dev.setConfiguration(int(n))
		}

	case usb.RecipientInterface:
		ifn := int(setup.WIndex)
		alt := int(setup.WValue)

		switch setup.Request {
		case usb.RequestGetStatus:
			return dev.getInterfaceStatus(ifn)

		case usb.RequestGetInterface:
			return dev.getInterface(ifn)

		case usb.RequestSetInterface:
			return dev.setInterface(ifn, alt)
		}
	}

	return nil, syscall.EPIPE
}

// classSpecific handles USB class-specific control requests.
func (srv *Server) controlClassSpecific(dev *Device,
	setup usb.SetupPacket) ([]byte, syscall.Errno) {

	// Get interface descriptor
	confno := int(setup.WValue)
	iffno := int(setup.WIndex >> 8)
	altno := int(setup.WIndex & 255)

	alt := dev.altIndex(confno, iffno, altno)

	if alt == nil {
		log.Debug(srv.ctx, "CTRL: %d-%d-%d: interface not found",
			confno, iffno, altno)

		return nil, syscall.EPIPE
	}

	// Check interface class and protocol
	if alt.BInterfaceClass != 7 || alt.BInterfaceSubClass != 1 {
		return nil, syscall.EPIPE
	}

	if alt.BInterfaceProtocol != 1 && alt.BInterfaceProtocol != 2 {
		return nil, syscall.EPIPE
	}

	// Handle class-specific request for the printer class
	if setup.RequestType&usb.SetupIn != 0 {
		switch setup.Request {
		case 0:
			// GET_DEVICE_ID
			devid := alt.IEEE1284DeviceID
			if devid == "" {
				return nil, syscall.EPIPE
			}

			length := len(devid) + 2
			data := make([]byte, length)
			data[0] = uint8(length >> 8)
			data[1] = uint8(length)
			copy(data[2:], devid)

			return data, 0
		}
	}

	return nil, syscall.EPIPE
}
