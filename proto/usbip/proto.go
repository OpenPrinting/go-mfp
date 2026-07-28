// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// USBIP protocol

package usbip

import (
	"fmt"
	"syscall"

	"github.com/OpenPrinting/go-mfp/util/generic"
)

// protoVersion is the protocol version
const protoVersion = 0x0111

// protoMaxTransferSize limits the maximum transfer size
const protoMaxTransferSize = 65536

// protoHsCode represents the handshake operation code.
type protoHsCode uint32

// protoHsCode assigned values:
const (
	protoHsDevlistReq protoHsCode = 0x8005 // OP_REQ_DEVLIST request
	protoHsDevlistRep protoHsCode = 0x0005 // OP_REP_DEVLIST reply
	protoHsImportReq  protoHsCode = 0x8003 // OP_REQ_IMPORT request
	protoHsImportRep  protoHsCode = 0x0003 // OP_REP_IMPORT reply
)

// protoHsCode represents the handshake status code.
type protoHsStatus uint32

// protoHsStatus assigned values:
const (
	protoHsOK    protoHsStatus = 0
	protoHsError protoHsStatus = 1
)

// protoHandshakeRequest is the common interface for the handshake requests.
type protoHandshakeRequest interface {
	// String returns string representation of the request, for logging
	String() string

	// Recv reads message from the connection.
	Recv(*protoConn) error
}

// protoHandshakeResponse is the common interface for the handshake responses.
type protoHandshakeResponse interface {
	// String returns string representation of the request, for logging
	String() string

	// OpCode returns the message-specific protoHsCode
	OpCode() protoHsCode

	// OpStatus() returns message status code.
	OpStatus() protoHsStatus

	// Send writes message into the connection.
	Send(*protoConn) error
}

// protoDevlistRequest represents the OP_REQ_DEVLIST message.
type protoDevlistRequest struct {
}

// String returns string representation of the request, for logging
func (protoDevlistRequest) String() string {
	return "OP_REQ_DEVLIST"
}

// Recv reads message from the connection.
func (protoDevlistRequest) Recv(*protoConn) error {
	return nil
}

// protoDevlistResponse represents the OP_REP_DEVLIST message.
type protoDevlistResponse struct {
	Status  protoHsStatus // Response status
	Devlist []devInfo     // List of devices
}

// String returns string representation of the request, for logging
func (rsp *protoDevlistResponse) String() string {
	if rsp.Status != protoHsOK {
		return "OP_REP_DEVLIST: error"
	}

	return fmt.Sprintf("OP_REP_DEVLIST: %d devices", len(rsp.Devlist))
}

// OpCode returns the message-specific protoHsCode
func (protoDevlistResponse) OpCode() protoHsCode {
	return protoHsDevlistRep
}

// OpStatus() returns message status code.
func (rsp *protoDevlistResponse) OpStatus() protoHsStatus {
	return rsp.Status
}

// Send writes message into the connection.
func (rsp protoDevlistResponse) Send(pconn *protoConn) error {
	enc := newEncoder(1024)

	enc.PutBE32(uint32(len(rsp.Devlist)))
	for _, dev := range rsp.Devlist {
		path := makeDevPath(dev.Location)
		busid := dev.Location.BusID()
		desc := dev.Device.Descriptor
		conf := desc.Configurations[0]

		enc.PutBytes(path[:]...)
		enc.PutBytes(busid[:]...)

		enc.PutBE32(uint32(dev.Location.Bus))
		enc.PutBE32(uint32(dev.Location.Dev))
		enc.PutBE32(uint32(desc.Speed))
		enc.PutBE16(desc.IDVendor)
		enc.PutBE16(desc.IDProduct)
		enc.PutBE16(uint16(desc.BCDDevice))
		enc.PutU8(desc.BDeviceClass)
		enc.PutU8(desc.BDeviceSubClass)
		enc.PutU8(desc.BDeviceProtocol)
		enc.PutU8(1) // Current configuration
		enc.PutU8(byte(len(desc.Configurations)))
		enc.PutU8(byte(len(conf.Interfaces)))

		for _, iff := range conf.Interfaces {
			for _, alt := range iff.AltSettings {
				enc.PutU8(alt.BInterfaceClass)
				enc.PutU8(alt.BInterfaceSubClass)
				enc.PutU8(alt.BInterfaceProtocol)
				enc.PutU8(0) // Padding byte
			}
		}
	}

	return pconn.write(enc.Bytes())
}

// protoImportRequest represents the OP_REQ_IMPORT message.
type protoImportRequest struct {
	BusID devBusID // Target BusID
}

// String returns string representation of the request, for logging
func (rq *protoImportRequest) String() string {
	return fmt.Sprintf("OP_REQ_IMPORT: %s", rq.BusID)
}

// Recv reads message from the connection.
func (rq *protoImportRequest) Recv(pconn *protoConn) error {
	return pconn.read(rq.BusID[:])
}

// protoImportResponse represents the OP_REP_IMPORT message.
type protoImportResponse struct {
	Status  protoHsStatus // Response status
	DevInfo *devInfo      // Attached device
}

// String returns string representation of the request, for logging
func (rsp *protoImportResponse) String() string {
	if rsp.Status != protoHsOK {
		return "OP_REP_IMPORT: error"
	}

	return fmt.Sprintf("OP_REP_IMPORT: %s OK", rsp.DevInfo.Location)
}

// OpCode returns the message-specific protoHsCode
func (protoImportResponse) OpCode() protoHsCode {
	return protoHsImportRep
}

// OpStatus() returns message status code.
func (rsp protoImportResponse) OpStatus() protoHsStatus {
	return rsp.Status
}

// Send writes message into the connection.
func (rsp protoImportResponse) Send(pconn *protoConn) error {
	enc := newEncoder(0x140)

	path := makeDevPath(rsp.DevInfo.Location)
	busid := rsp.DevInfo.Location.BusID()
	desc := rsp.DevInfo.Device.Descriptor

	enc.PutBytes(path[:]...)
	enc.PutBytes(busid[:]...)

	enc.PutBE32(uint32(rsp.DevInfo.Location.Bus))
	enc.PutBE32(uint32(rsp.DevInfo.Location.Dev))
	enc.PutBE32(uint32(desc.Speed))
	enc.PutBE16(desc.IDVendor)
	enc.PutBE16(desc.IDProduct)
	enc.PutBE16(uint16(desc.BCDDevice))
	enc.PutU8(desc.BDeviceClass)
	enc.PutU8(desc.BDeviceSubClass)
	enc.PutU8(desc.BDeviceProtocol)
	enc.PutU8(1) // Current configuration
	enc.PutU8(byte(len(desc.Configurations)))
	enc.PutU8(byte(len(desc.Configurations[0].Interfaces)))

	return pconn.write(enc.Bytes())
}

// protoIORequest is the common interface for the I/O requests.
type protoIORequest interface {
	// String returns string representation of the request, for logging
	String() string

	// Header returns the request protocol header.
	Header() *protoIOHeader

	// Recv reads the request body from the connection.
	Recv(*protoConn) error
}

// protoIOResponse is the common interface for the I/O requests.
type protoIOResponse interface {
	// Header returns the response protocol header.
	Header() *protoIOHeader

	// Send writes the response body from the connection.
	Send(*protoConn) error
}

// protoIOCommand is the command code for the I/O messages.
type protoIOCommand uint32

// protoIOCommand assigned values:
const (
	protoIOSubmitCmd protoIOCommand = 0x00000001 // USBIP_CMD_SUBMIT
	protoIOSubmitRet protoIOCommand = 0x00000003 // USBIP_RET_SUBMIT
	protoIOUnlinkCmd protoIOCommand = 0x00000002 // USBIP_CMD_UNLINK
	protoIOUnlinkRet protoIOCommand = 0x00000004 // USBIP_RET_UNLINK
)

// protoTransferFlags is the USBIP transfer flags
type protoTransferFlags uint32

// protoTransferFlags assigned bits:
const (
	// Short reads not allowed
	protoTransferShortNotOk protoTransferFlags = 0x0001

	// ISO only: do transfer ASAP
	protoTransferIsoASAP protoTransferFlags = 0x0002

	// Finish bulk OUT with short packet
	protoTransferZeroPacket protoTransferFlags = 0x0040

	// Hint: non-error interrupt is not needed
	protoTransferNoInterrupt protoTransferFlags = 0x0080

	// Transfer direction: input if set
	protoTransferDirIn protoTransferFlags = 0x0200
)

// protoIOHeader is the common header for all I/O requests and
// responses.
type protoIOHeader struct {
	Command  protoIOCommand // Command code
	Seqnum   uint32         // Sequence number
	Location devLocation    // Device location
	Input    bool           // Direction, true for device->host input
	Endpoint uint8          // Endpoint number (client only)
}

// Recv reads protoIOHeader from the connection.
func (hdr *protoIOHeader) Recv(pconn *protoConn) error {
	var buf [20]byte
	err := pconn.read(buf[:])
	if err != nil {
		return err
	}

	dec := newDecoder(buf[:])
	hdr.Command = protoIOCommand(dec.GetBE32())
	hdr.Seqnum = dec.GetBE32()
	bus := dec.GetBE16()
	dev := dec.GetBE16()
	hdr.Location = devLocation{Bus: int(bus), Dev: int(dev)}
	inp := dec.GetBE32()
	hdr.Input = inp != 0
	hdr.Endpoint = uint8(dec.GetBE32())

	return nil
}

// Send writes protoIOHeader into the connection.
func (hdr *protoIOHeader) Send(pconn *protoConn) error {
	enc := newEncoder(20)
	enc.PutBE32(uint32(hdr.Command))
	enc.PutBE32(hdr.Seqnum)
	enc.PutBE16(uint16(hdr.Location.Bus))
	enc.PutBE16(uint16(hdr.Location.Dev))
	inp := uint32(0)
	if hdr.Input {
		inp = 1
	}
	enc.PutBE32(inp)
	enc.PutBE32(uint32(hdr.Endpoint))

	return pconn.write(enc.Bytes())
}

// protoIOSubmitRequest represents the USBIP_CMD_SUBMIT message.
type protoIOSubmitRequest struct {
	// These are protocol items
	protoIOHeader                    // Request header
	Flags         protoTransferFlags // Transfer flags
	Length        uint32             // Transfer buffer length
	ISOStartFrame uint32             // ISO start frame
	ISONumPackets uint32             // ISO number of packets
	Interval      uint32             // ISO/Interrupt max time
	Setup         [8]byte            // USB setup bytes
	Buffer        []byte             // Data buffer

	// Processing state
	actualLength int   // Length of data transferred so far
	completion   func( // Completion callback
		*protoIOSubmitRequest, syscall.Errno)
}

// String returns string representation of the request, for logging
func (rq *protoIOSubmitRequest) String() string {
	dir := "in"
	if (rq.Flags & protoTransferDirIn) == 0 {
		dir = "out"
	}

	return fmt.Sprintf("USBIP_CMD_SUBMIT #%8.8x: %s %d bytes, ep: %d",
		rq.Seqnum, dir, rq.Length, rq.Endpoint)
}

// Recv reads protoIOSubmitRequest from the connection.
// The protoIOHeader assumed to be already received at
// this point.
func (rq *protoIOSubmitRequest) Recv(pconn *protoConn) error {
	var buf [28]byte
	err := pconn.read(buf[:])
	if err != nil {
		return err
	}

	dec := newDecoder(buf[:])
	rq.Flags = protoTransferFlags(dec.GetBE32())
	rq.Length = dec.GetBE32()
	rq.ISOStartFrame = dec.GetBE32()
	rq.ISONumPackets = dec.GetBE32()
	rq.Interval = dec.GetBE32()
	dec.GetData(rq.Setup[:])

	if (rq.Flags & protoTransferDirIn) == 0 {
		sz := generic.Min(protoMaxTransferSize, rq.Length)

		rq.Buffer = make([]byte, sz)
		err = pconn.read(rq.Buffer)
		if err == nil && rq.Length > sz {
			err = pconn.drain(rq.Length - sz)
		}
		if err != nil {
			return err
		}

		rq.Length = sz
	}

	return nil
}

// Header returns the Request protocol header.
func (rq *protoIOSubmitRequest) Header() *protoIOHeader {
	return &rq.protoIOHeader
}

// Response creates protoIOSubmitResponse in reply to the
// protoIOSubmitRequest
func (rq *protoIOSubmitRequest) Response(
	err syscall.Errno) *protoIOSubmitResponse {

	rsp := &protoIOSubmitResponse{
		protoIOHeader: rq.protoIOHeader,
		ActualLength:  uint32(rq.actualLength),
		Status:        err,
		Buffer:        rq.Buffer[:rq.actualLength],
	}
	rsp.protoIOHeader.Command = protoIOSubmitRet
	rsp.protoIOHeader.Input = false
	rsp.protoIOHeader.Endpoint = 0

	if len(rsp.Buffer) > int(rq.Length) {
		// Truncate the response data if it exceeds request.
		rsp.ActualLength = rq.Length
		rsp.Buffer = rsp.Buffer[:rq.Length]
	}

	if !rq.Input && rq.Endpoint != 0 {
		rsp.Buffer = nil
	}

	return rsp
}

// protoIOSubmitResponse represents the USBIP_RET_SUBMIT message.
type protoIOSubmitResponse struct {
	protoIOHeader
	Status        syscall.Errno // 0 - OK, otherwise errno
	ActualLength  uint32        // Actual transfer length
	ISOStartFrame uint32        // ISO start frame
	ISONumPackets uint32        // ISO number of packets
	ErrCount      uint32        // Count of errors
	Buffer        []byte        // Data buffer
}

// String returns string representation of the response, for logging
func (rsp *protoIOSubmitResponse) String() string {
	dir := "out"
	if rsp.Buffer != nil {
		dir = "in"
	}

	return fmt.Sprintf("USBIP_RET_SUBMIT #%8.8x: %s %d bytes, ep: %d, %s",
		rsp.Seqnum,
		dir, rsp.ActualLength, rsp.Endpoint,
		syscall.Errno(rsp.Status))
}

// Header returns the response protocol header.
func (rsp *protoIOSubmitResponse) Header() *protoIOHeader {
	return &rsp.protoIOHeader
}

// Send sends protoIOSubmitResponse into the connection.
func (rsp *protoIOSubmitResponse) Send(pconn *protoConn) error {
	enc := newEncoder(28)
	enc.PutBE32(uint32(-rsp.Status))

	enc.PutBE32(rsp.ActualLength)
	enc.PutBE32(rsp.ISOStartFrame)
	enc.PutBE32(rsp.ISONumPackets)
	enc.PutBE32(rsp.ErrCount)
	var padding [8]byte
	enc.PutBytes(padding[:]...)

	err := pconn.write(enc.Bytes())
	if err == nil && len(rsp.Buffer) > 0 {
		err = pconn.write(rsp.Buffer)
	}

	return err
}

// protoIOUnlinkRequest represents the USBIP_CMD_UNLINK message.
type protoIOUnlinkRequest struct {
	protoIOHeader
	UnlinkSeqnum uint32 // Target sequence number
}

// Header returns the Request protocol header.
func (rq *protoIOUnlinkRequest) Header() *protoIOHeader {
	return &rq.protoIOHeader
}

// String returns string representation of the request, for logging
func (rq *protoIOUnlinkRequest) String() string {
	return fmt.Sprintf("USBIP_CMD_UNLINK #%8.8x: target=%8.8x",
		rq.Seqnum, rq.UnlinkSeqnum)
}

// Recv reads protoIOSubmitRequest from the connection.
// The protoIOHeader assumed to be already received at
// this point.
func (rq *protoIOUnlinkRequest) Recv(pconn *protoConn) error {
	var buf [28]byte
	err := pconn.read(buf[:])
	if err != nil {
		return err
	}

	dec := newDecoder(buf[:])
	rq.UnlinkSeqnum = dec.GetBE32()

	return nil
}

// Response creates protoIOUnlinkResponse in reply to the
// protoIOUnlinkRequest
func (rq *protoIOUnlinkRequest) Response(
	err syscall.Errno) *protoIOUnlinkResponse {

	rsp := &protoIOUnlinkResponse{
		protoIOHeader: rq.protoIOHeader,
		Status:        err,
	}

	rsp.protoIOHeader.Command = protoIOUnlinkRet
	rsp.protoIOHeader.Input = false
	rsp.protoIOHeader.Endpoint = 0

	return rsp
}

// protoIOUnlinkResponse represents the USBIP_RET_UNLINK message.
type protoIOUnlinkResponse struct {
	protoIOHeader
	Status syscall.Errno // OK - ECONREFUSED, 0 - not found or errno
}

// Header returns the response protocol header.
func (rsp *protoIOUnlinkResponse) Header() *protoIOHeader {
	return &rsp.protoIOHeader
}

// Send sends protoIOUnlinkResponse into the connection.
func (rsp *protoIOUnlinkResponse) Send(pconn *protoConn) error {
	enc := newEncoder(28)
	enc.PutBE32(uint32(-rsp.Status))
	var padding [24]byte
	enc.PutBytes(padding[:]...)

	err := pconn.write(enc.Bytes())
	return err
}
