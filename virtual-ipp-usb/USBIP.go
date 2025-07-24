package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

const (
	// USBIPDirOut indicates OUT direction in USB/IP
	USBIPDirOut = 0
	// USBIPDirIn indicates IN direction in USB/IP
	USBIPDirIn = 1
)

// USBIPHeader represents the common USB/IP protocol header used in communication.
type USBIPHeader struct {
	Version uint16
	Command uint16
	Status  uint32
}

// Pack serializes the USBIPHeader into a byte slice using big-endian encoding.
func (h *USBIPHeader) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h.Version)
	binary.Write(buf, binary.BigEndian, h.Command)
	binary.Write(buf, binary.BigEndian, h.Status)
	return buf.Bytes()
}

// Unpack deserializes the given byte slice into the USBIPHeader fields using big-endian encoding.
func (h *USBIPHeader) Unpack(data []byte) {
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.BigEndian, &h.Version)
	binary.Read(buf, binary.BigEndian, &h.Command)
	binary.Read(buf, binary.BigEndian, &h.Status)
}

// Size returns the fixed size in bytes of a packed USBIPHeader.
func (h *USBIPHeader) Size() int {
	return 8
}

// USBInterface represents a USB interface descriptor containing class and protocol information.
type USBInterface struct {
	BInterfaceClass    uint8
	BInterfaceSubClass uint8
	BInterfaceProtocol uint8
	Align              uint8
}

// Pack serializes the USBInterface into a byte slice using big-endian encoding.
func (u *USBInterface) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u.BInterfaceClass)
	binary.Write(buf, binary.BigEndian, u.BInterfaceSubClass)
	binary.Write(buf, binary.BigEndian, u.BInterfaceProtocol)
	binary.Write(buf, binary.BigEndian, u.Align)
	return buf.Bytes()
}

// OPREPDevList represents the device list reply structure in the USB/IP protocol.
type OPREPDevList struct {
	Base                USBIPHeader
	NExportedDevice     uint32
	UsbPath             [256]byte
	BusID               [32]byte
	Busnum              uint32
	Devnum              uint32
	Speed               uint32
	IDVendor            uint16
	IDProduct           uint16
	BcdDevice           uint16
	BDeviceClass        uint8
	BDeviceSubClass     uint8
	BDeviceProtocol     uint8
	BConfigurationValue uint8
	BNumConfigurations  uint8
	BNumInterfaces      uint8
	Interfaces          USBInterface
}

// Pack serializes the OPREPDevList into a byte slice using big-endian encoding.
func (o *OPREPDevList) Pack() []byte {
	buf := new(bytes.Buffer)
	buf.Write(o.Base.Pack())
	binary.Write(buf, binary.BigEndian, o.NExportedDevice)
	buf.Write(o.UsbPath[:])
	buf.Write(o.BusID[:])
	binary.Write(buf, binary.BigEndian, o.Busnum)
	binary.Write(buf, binary.BigEndian, o.Devnum)
	binary.Write(buf, binary.BigEndian, o.Speed)
	binary.Write(buf, binary.BigEndian, o.IDVendor)
	binary.Write(buf, binary.BigEndian, o.IDProduct)
	binary.Write(buf, binary.BigEndian, o.BcdDevice)
	binary.Write(buf, binary.BigEndian, o.BDeviceClass)
	binary.Write(buf, binary.BigEndian, o.BDeviceSubClass)
	binary.Write(buf, binary.BigEndian, o.BDeviceProtocol)
	binary.Write(buf, binary.BigEndian, o.BConfigurationValue)
	binary.Write(buf, binary.BigEndian, o.BNumConfigurations)
	binary.Write(buf, binary.BigEndian, o.BNumInterfaces)
	buf.Write(o.Interfaces.Pack())
	return buf.Bytes()
}

// OPREPImport represents the structure used to respond to an import request in the USB/IP protocol.
type OPREPImport struct {
	Base                USBIPHeader
	UsbPath             [256]byte
	BusID               [32]byte
	Busnum              uint32
	Devnum              uint32
	Speed               uint32
	IDVendor            uint16
	IDProduct           uint16
	BcdDevice           uint16
	BDeviceClass        uint8
	BDeviceSubClass     uint8
	BDeviceProtocol     uint8
	BConfigurationValue uint8
	BNumConfigurations  uint8
	BNumInterfaces      uint8
}

// Pack serializes the OPREPImport into a byte slice using big-endian encoding.
func (o *OPREPImport) Pack() []byte {
	buf := new(bytes.Buffer)
	buf.Write(o.Base.Pack())
	buf.Write(o.UsbPath[:])
	buf.Write(o.BusID[:])
	binary.Write(buf, binary.BigEndian, o.Busnum)
	binary.Write(buf, binary.BigEndian, o.Devnum)
	binary.Write(buf, binary.BigEndian, o.Speed)
	binary.Write(buf, binary.BigEndian, o.IDVendor)
	binary.Write(buf, binary.BigEndian, o.IDProduct)
	binary.Write(buf, binary.BigEndian, o.BcdDevice)
	binary.Write(buf, binary.BigEndian, o.BDeviceClass)
	binary.Write(buf, binary.BigEndian, o.BDeviceSubClass)
	binary.Write(buf, binary.BigEndian, o.BDeviceProtocol)
	binary.Write(buf, binary.BigEndian, o.BConfigurationValue)
	binary.Write(buf, binary.BigEndian, o.BNumConfigurations)
	binary.Write(buf, binary.BigEndian, o.BNumInterfaces)
	return buf.Bytes()
}

// USBIPRETSubmit represents the USB/IP return submit packet structure.
type USBIPRETSubmit struct {
	Command         uint32
	Seqnum          uint32
	Devid           uint32
	Direction       uint32
	Ep              uint32
	Status          uint32
	ActualLength    uint32
	StartFrame      uint32
	NumberOfPackets uint32
	ErrorCount      uint32
	Padding         uint64
	Data            []byte
}

// Pack serializes the USBIPRETSubmit structure into a byte slice using big-endian encoding.
func (u *USBIPRETSubmit) Pack() []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u.Command)
	binary.Write(buf, binary.BigEndian, u.Seqnum)
	binary.Write(buf, binary.BigEndian, u.Devid)
	binary.Write(buf, binary.BigEndian, u.Direction)
	binary.Write(buf, binary.BigEndian, u.Ep)
	binary.Write(buf, binary.BigEndian, u.Status)
	binary.Write(buf, binary.BigEndian, u.ActualLength)
	binary.Write(buf, binary.BigEndian, u.StartFrame)
	binary.Write(buf, binary.BigEndian, u.NumberOfPackets)
	binary.Write(buf, binary.BigEndian, u.ErrorCount)
	binary.Write(buf, binary.BigEndian, u.Padding)
	buf.Write(u.Data)
	return buf.Bytes()
}

// Size returns the fixed size (in bytes) of the USBIPRETSubmit header.
func (u *USBIPRETSubmit) Size() int {

	return 48
}

// USBIPCMDSubmit represents the USB/IP command submit packet structure.
type USBIPCMDSubmit struct {
	Command              uint32
	Seqnum               uint32
	Devid                uint32
	Direction            uint32
	Ep                   uint32
	TransferFlags        uint32
	TransferBufferLength uint32
	StartFrame           uint32
	NumberOfPackets      uint32
	Interval             uint32
	Setup                [8]byte
}

// Unpack deserializes the byte slice into a USBIPCMDSubmit structure using big-endian encoding.
func (u *USBIPCMDSubmit) Unpack(data []byte) {

	buf := bytes.NewReader(data)
	binary.Read(buf, binary.BigEndian, &u.Command)
	binary.Read(buf, binary.BigEndian, &u.Seqnum)
	binary.Read(buf, binary.BigEndian, &u.Devid)
	binary.Read(buf, binary.BigEndian, &u.Direction)
	binary.Read(buf, binary.BigEndian, &u.Ep)
	binary.Read(buf, binary.BigEndian, &u.TransferFlags)
	binary.Read(buf, binary.BigEndian, &u.TransferBufferLength)
	binary.Read(buf, binary.BigEndian, &u.StartFrame)
	binary.Read(buf, binary.BigEndian, &u.NumberOfPackets)
	binary.Read(buf, binary.BigEndian, &u.Interval)
	copy(u.Setup[:], data[40:48])
}

// Size returns the fixed size (in bytes) of a USBIPCMDSubmit structure.
func (u *USBIPCMDSubmit) Size() int {

	return 48
}

// StandardDeviceRequest represents a standard USB control request structure.
type StandardDeviceRequest struct {
	BmRequestType uint8
	BRequest      uint8
	WValue        uint16
	WIndex        uint16
	WLength       uint16
}

// Unpack parses a StandardDeviceRequest from a byte slice.
func (s *StandardDeviceRequest) Unpack(data []byte) {
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &s.BmRequestType)
	binary.Read(buf, binary.LittleEndian, &s.BRequest)
	binary.Read(buf, binary.LittleEndian, &s.WValue)
	binary.Read(buf, binary.LittleEndian, &s.WIndex)
	binary.Read(buf, binary.LittleEndian, &s.WLength)
}

// DeviceDescriptor describes a USB device, including vendor/product IDs and configurations.
type DeviceDescriptor struct {
	BLength            uint8
	BDescriptorType    uint8
	BcdUSB             uint16
	BDeviceClass       uint8
	BDeviceSubClass    uint8
	BDeviceProtocol    uint8
	BMaxPacketSize0    uint8
	IDVendor           uint16
	IDProduct          uint16
	BcdDevice          uint16
	IManufacturer      uint8
	IProduct           uint8
	ISerialNumber      uint8
	BNumConfigurations uint8
}

// Pack serializes the DeviceDescriptor into a byte slice using little-endian format.
func (d DeviceDescriptor) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, d.BLength)
	binary.Write(buf, binary.LittleEndian, d.BDescriptorType)
	binary.Write(buf, binary.LittleEndian, d.BcdUSB)
	binary.Write(buf, binary.LittleEndian, d.BDeviceClass)
	binary.Write(buf, binary.LittleEndian, d.BDeviceSubClass)
	binary.Write(buf, binary.LittleEndian, d.BDeviceProtocol)
	binary.Write(buf, binary.LittleEndian, d.BMaxPacketSize0)
	binary.Write(buf, binary.LittleEndian, d.IDVendor)
	binary.Write(buf, binary.LittleEndian, d.IDProduct)
	binary.Write(buf, binary.LittleEndian, d.BcdDevice)
	binary.Write(buf, binary.LittleEndian, d.IManufacturer)
	binary.Write(buf, binary.LittleEndian, d.IProduct)
	binary.Write(buf, binary.LittleEndian, d.ISerialNumber)
	binary.Write(buf, binary.LittleEndian, d.BNumConfigurations)
	return buf.Bytes()
}

// DeviceConfiguration defines a configuration for a USB device including interfaces.
type DeviceConfiguration struct {
	BLength             uint8
	BDescriptorType     uint8
	WTotalLength        uint16
	BNumInterfaces      uint8
	BConfigurationValue uint8
	IConfiguration      uint8
	BmAttributes        uint8
	BMaxPower           uint8
	Interfaces          [][]InterfaceDescriptor
}

// Pack serializes the DeviceConfiguration into a byte slice.
func (d *DeviceConfiguration) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, d.BLength)
	binary.Write(buf, binary.LittleEndian, d.BDescriptorType)
	binary.Write(buf, binary.LittleEndian, d.WTotalLength)
	binary.Write(buf, binary.LittleEndian, d.BNumInterfaces)
	binary.Write(buf, binary.LittleEndian, d.BConfigurationValue)
	binary.Write(buf, binary.LittleEndian, d.IConfiguration)
	binary.Write(buf, binary.LittleEndian, d.BmAttributes)
	binary.Write(buf, binary.LittleEndian, d.BMaxPower)
	return buf.Bytes()
}

// BOSDescriptor defines a Binary Object Store descriptor for a USB device.
type BOSDescriptor struct {
	BLength         uint8
	BDescriptorType uint8
	WTotalLength    uint16
	BNumDeviceCaps  uint8
}

// Pack serializes the BOSDescriptor into a byte slice.
func (b *BOSDescriptor) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, b.BLength)
	binary.Write(buf, binary.LittleEndian, b.BDescriptorType)
	binary.Write(buf, binary.LittleEndian, b.WTotalLength)
	binary.Write(buf, binary.LittleEndian, b.BNumDeviceCaps)
	return buf.Bytes()
}

// DeviceQualifierDescriptor describes a high-speed capable device in other speeds.
type DeviceQualifierDescriptor struct {
	BLength            uint8
	BDescriptorType    uint8
	BcdUSB             uint16
	BDeviceClass       uint8
	BDeviceSubClass    uint8
	BDeviceProtocol    uint8
	BMaxPacketSize0    uint8
	BNumConfigurations uint8
	BReserved          uint8
}

// Pack serializes the DeviceQualifierDescriptor into a byte slice.
func (d *DeviceQualifierDescriptor) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, d.BLength)
	binary.Write(buf, binary.LittleEndian, d.BDescriptorType)
	binary.Write(buf, binary.LittleEndian, d.BcdUSB)
	binary.Write(buf, binary.LittleEndian, d.BDeviceClass)
	binary.Write(buf, binary.LittleEndian, d.BDeviceSubClass)
	binary.Write(buf, binary.LittleEndian, d.BDeviceProtocol)
	binary.Write(buf, binary.LittleEndian, d.BMaxPacketSize0)
	binary.Write(buf, binary.LittleEndian, d.BNumConfigurations)
	binary.Write(buf, binary.LittleEndian, d.BReserved)
	return buf.Bytes()
}

// InterfaceDescriptor describes a single interface within a configuration.
type InterfaceDescriptor struct {
	BLength            uint8
	BDescriptorType    uint8
	BInterfaceNumber   uint8
	BAlternateSetting  uint8
	BNumEndpoints      uint8
	BInterfaceClass    uint8
	BInterfaceSubClass uint8
	BInterfaceProtocol uint8
	IInterface         uint8
	ClassDescriptor    interface{ Pack() []byte }
	Endpoints          []EndpointDescriptor
}

// Pack serializes the InterfaceDescriptor into a byte slice.
func (i *InterfaceDescriptor) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, i.BLength)
	binary.Write(buf, binary.LittleEndian, i.BDescriptorType)
	binary.Write(buf, binary.LittleEndian, i.BInterfaceNumber)
	binary.Write(buf, binary.LittleEndian, i.BAlternateSetting)
	binary.Write(buf, binary.LittleEndian, i.BNumEndpoints)
	binary.Write(buf, binary.LittleEndian, i.BInterfaceClass)
	binary.Write(buf, binary.LittleEndian, i.BInterfaceSubClass)
	binary.Write(buf, binary.LittleEndian, i.BInterfaceProtocol)
	binary.Write(buf, binary.LittleEndian, i.IInterface)
	return buf.Bytes()
}

// EndpointDescriptor describes a USB endpoint within an interface.
type EndpointDescriptor struct {
	BLength          uint8
	BDescriptorType  uint8
	BEndpointAddress uint8
	BmAttributes     uint8
	WMaxPacketSize   uint16
	BInterval        uint8
	ClassDescriptor  interface{ Pack() []byte }
}

// Pack serializes the EndpointDescriptor into a byte slice.
func (e *EndpointDescriptor) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, e.BLength)
	binary.Write(buf, binary.LittleEndian, e.BDescriptorType)
	binary.Write(buf, binary.LittleEndian, e.BEndpointAddress)
	binary.Write(buf, binary.LittleEndian, e.BmAttributes)
	binary.Write(buf, binary.LittleEndian, e.WMaxPacketSize)
	binary.Write(buf, binary.LittleEndian, e.BInterval)
	return buf.Bytes()
}

// USBRequest represents a USB transfer request received from the host.
type USBRequest struct {
	Seqnum               uint32
	Devid                uint32
	Direction            uint32
	Ep                   uint32
	Flags                uint32
	TransferBufferLength uint32
	NumberOfPackets      uint32
	Interval             uint32
	Setup                [8]byte
	TransferBuffer       []byte
}

// USBDevice defines an interface for emulated USB devices.
type USBDevice interface {
	GetConfigurations() []DeviceConfiguration
	GetDeviceDescriptor() DeviceDescriptor
	HandleData(usbReq USBRequest)
	HandleDeviceSpecificControl(controlReq StandardDeviceRequest, usbReq USBRequest)
	SetConnection(conn net.Conn)
}

// BaseUSBDevice provides a default implementation for core USB device behaviors.
type BaseUSBDevice struct {
	Connection        net.Conn
	AllConfigurations []byte
}

// SetConnection sets the TCP connection used to communicate with the host.
func (b *BaseUSBDevice) SetConnection(conn net.Conn) {
	b.Connection = conn
}

// GenerateRawConfiguration flattens and serializes the configuration descriptors.
func (b *BaseUSBDevice) GenerateRawConfiguration(device USBDevice) {
	var allConfigurations []byte
	for _, configuration := range device.GetConfigurations() {
		allConfigurations = append(allConfigurations, configuration.Pack()...)
		for _, interfaceGroup := range configuration.Interfaces {
			for _, interfaceAlt := range interfaceGroup {
				allConfigurations = append(allConfigurations, interfaceAlt.Pack()...)
				if interfaceAlt.ClassDescriptor != nil {
					allConfigurations = append(allConfigurations, interfaceAlt.ClassDescriptor.Pack()...)
				}
				for _, endpoint := range interfaceAlt.Endpoints {
					allConfigurations = append(allConfigurations, endpoint.Pack()...)
					if endpoint.ClassDescriptor != nil {
						allConfigurations = append(allConfigurations, endpoint.ClassDescriptor.Pack()...)
					}
				}
			}
		}
	}
	b.AllConfigurations = allConfigurations
}

// SendUSBRet sends a USBIP_RET_Submit response back to the host.
func (b *BaseUSBDevice) SendUSBRet(usbReq USBRequest, usbRes []byte, usbLen int, status uint32) {
	fmt.Printf("Sending %s\n", BytesToString(usbRes))
	ret := &USBIPRETSubmit{
		Command:         0x3,
		Seqnum:          usbReq.Seqnum,
		Devid:           0,
		Direction:       0,
		Ep:              0,
		Status:          status,
		ActualLength:    uint32(usbLen),
		StartFrame:      0,
		NumberOfPackets: 0xffffffff,
		ErrorCount:      0,
		Padding:         0,
		Data:            usbRes,
	}
	b.Connection.Write(ret.Pack())
}

// HandleGetDescriptor processes a GET_DESCRIPTOR request from the host.
func (b *BaseUSBDevice) HandleGetDescriptor(device USBDevice, controlReq StandardDeviceRequest, usbReq USBRequest) bool {
	descriptorType := uint8(controlReq.WValue >> 8)
	descriptorIndex := uint8(controlReq.WValue & 0xff)
	fmt.Printf("handle_get_descriptor %d %d\n", descriptorType, descriptorIndex)

	if descriptorType == 0x01 { // Device Descriptor
		ret := device.GetDeviceDescriptor().Pack()
		b.SendUSBRet(usbReq, ret, len(ret), 0)
		return true
	} else if descriptorType == 0x02 { // Configuration Descriptor
		ret := b.AllConfigurations
		if int(controlReq.WLength) < len(ret) {
			ret = ret[:controlReq.WLength]
		}
		b.SendUSBRet(usbReq, ret, len(ret), 0)
		return true
	}
	return false
}

// HandleSetConfiguration processes a SET_CONFIGURATION request from the host.
func (b *BaseUSBDevice) HandleSetConfiguration(device USBDevice, controlReq StandardDeviceRequest, usbReq USBRequest) bool {
	fmt.Printf("handle_set_configuration %d\n", controlReq.WValue)
	b.SendUSBRet(usbReq, []byte{}, 0, 0)
	return true
}

// HandleUSBControl handles a standard USB control transfer on endpoint 0.
func (b *BaseUSBDevice) HandleUSBControl(device USBDevice, usbReq USBRequest) {
	controlReq := StandardDeviceRequest{}
	controlReq.Unpack(usbReq.Setup[:])
	handled := false
	fmt.Printf("  UC Request Type %d\n", controlReq.BmRequestType)
	fmt.Printf("  UC Request %d\n", controlReq.BRequest)
	fmt.Printf("  UC Value  %d\n", controlReq.WValue)
	fmt.Printf("  UC Index  %d\n", controlReq.WIndex)
	fmt.Printf("  UC Length %d\n", controlReq.WLength)

	if controlReq.BmRequestType == 0x80 { // Data flows IN, from Device to Host
		if controlReq.BRequest == 0x00 { // GET_STATUS
			configurations := device.GetConfigurations()
			attributes := configurations[0].BmAttributes
			isSelfPowered := (attributes & (1 << 6)) != 0
			isRemoteWakeup := (attributes & (1 << 5)) != 0
			ret := uint16(0x0000)
			if isRemoteWakeup {
				ret |= (1 << 1)
			}
			if isSelfPowered {
				ret |= 1
			}
			retBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(retBytes, ret)
			b.SendUSBRet(usbReq, retBytes, 2, 0)
			handled = true
		} else if controlReq.BRequest == 0x06 { // GET_DESCRIPTOR
			handled = b.HandleGetDescriptor(device, controlReq, usbReq)
		}
	} else if controlReq.BmRequestType == 0x00 { // Data flows OUT, from Host to Device
		if controlReq.BRequest == 0x09 { // Set Configuration
			handled = b.HandleSetConfiguration(device, controlReq, usbReq)
		}
	}

	if !handled {
		device.HandleDeviceSpecificControl(controlReq, usbReq)
	}
}

// HandleUSBRequest routes a USB request to the appropriate handler.
func (b *BaseUSBDevice) HandleUSBRequest(device USBDevice, usbReq USBRequest) {
	if usbReq.Ep == 0 { // Endpoint 0 is always the control endpoint
		b.HandleUSBControl(device, usbReq)
	} else {
		device.HandleData(usbReq)
	}
}

// BytesToString returns a hex-escaped string for a byte slice
func BytesToString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	var result strings.Builder
	for _, b := range data {
		result.WriteString(fmt.Sprintf("\\x%02x", b))
	}
	return result.String()
}

// USBContainer manages a collection of virtual USB devices.
type USBContainer struct {
	USBDevices []USBDevice
}

// AddUSBDevice adds a USBDevice to the container.
func (c *USBContainer) AddUSBDevice(device USBDevice) {
	c.USBDevices = append(c.USBDevices, device)
}

// HandleAttach handles a device import request from the USB/IP client.
func (c *USBContainer) HandleAttach() *OPREPImport {
	usbDev := c.USBDevices[0]
	deviceDescriptor := usbDev.GetDeviceDescriptor()
	configurations := usbDev.GetConfigurations()

	var usbPath [256]byte
	copy(usbPath[:], "/sys/devices/pci0000:00/0000:00:01.2/usb1/1-1")

	var busID [32]byte
	copy(busID[:], "1-1")

	return &OPREPImport{
		Base:                USBIPHeader{Version: 0x0111, Command: 3, Status: 0},
		UsbPath:             usbPath,
		BusID:               busID,
		Busnum:              1,
		Devnum:              2,
		Speed:               2,
		IDVendor:            deviceDescriptor.IDVendor,
		IDProduct:           deviceDescriptor.IDProduct,
		BcdDevice:           deviceDescriptor.BcdDevice,
		BDeviceClass:        deviceDescriptor.BDeviceClass,
		BDeviceSubClass:     deviceDescriptor.BDeviceSubClass,
		BDeviceProtocol:     deviceDescriptor.BDeviceProtocol,
		BConfigurationValue: configurations[0].BConfigurationValue,
		BNumConfigurations:  deviceDescriptor.BNumConfigurations,
		BNumInterfaces:      configurations[0].BNumInterfaces,
	}
}

// HandleDeviceList responds with a list of available USB devices.
func (c *USBContainer) HandleDeviceList() *OPREPDevList {
	usbDev := c.USBDevices[0]
	deviceDescriptor := usbDev.GetDeviceDescriptor()
	configurations := usbDev.GetConfigurations()

	var usbPath [256]byte
	copy(usbPath[:], "/sys/devices/pci0000:00/0000:00:01.2/usb1/1-1")

	var busID [32]byte
	copy(busID[:], "1-1")

	return &OPREPDevList{
		Base:                USBIPHeader{Version: 0x0111, Command: 5, Status: 0},
		NExportedDevice:     1,
		UsbPath:             usbPath,
		BusID:               busID,
		Busnum:              1,
		Devnum:              2,
		Speed:               2,
		IDVendor:            deviceDescriptor.IDVendor,
		IDProduct:           deviceDescriptor.IDProduct,
		BcdDevice:           deviceDescriptor.BcdDevice,
		BDeviceClass:        deviceDescriptor.BDeviceClass,
		BDeviceSubClass:     deviceDescriptor.BDeviceSubClass,
		BDeviceProtocol:     deviceDescriptor.BDeviceProtocol,
		BConfigurationValue: configurations[0].BConfigurationValue,
		BNumConfigurations:  deviceDescriptor.BNumConfigurations,
		BNumInterfaces:      configurations[0].BNumInterfaces,
		Interfaces: USBInterface{
			BInterfaceClass:    configurations[0].Interfaces[0][0].BInterfaceClass,
			BInterfaceSubClass: configurations[0].Interfaces[0][0].BInterfaceSubClass,
			BInterfaceProtocol: configurations[0].Interfaces[0][0].BInterfaceProtocol,
			Align:              0,
		},
	}
}

// Run starts the USB/IP server and handles incoming connections.
func (c *USBContainer) Run(ip string, port int) {
	if ip == "" {
		ip = "0.0.0.0"
	}
	if port == 0 {
		port = 3240
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	attached := false
	req := USBIPHeader{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("Connection address: %s\n", conn.RemoteAddr().String())

		for {
			if !attached {
				data := make([]byte, 8)
				n, err := conn.Read(data)
				if err != nil || n == 0 {
					break
				}
				req.Unpack(data)
				fmt.Println("Header Packet")
				fmt.Printf("command: %x\n", req.Command)

				if req.Command == 0x8005 { // OP_REQ_DEVLIST
					fmt.Println("list of devices")
					conn.Write(c.HandleDeviceList().Pack())
				} else if req.Command == 0x8003 { // OP_REQ_IMPORT
					fmt.Println("attach device")
					busIDData := make([]byte, 32)
					conn.Read(busIDData) // receive bus id
					conn.Write(c.HandleAttach().Pack())
					attached = true
				}
			} else {
				fmt.Println("----------------")
				fmt.Println("handles requests")
				cmd := USBIPCMDSubmit{}
				cmdHeaderData := make([]byte, cmd.Size())
				n, err := conn.Read(cmdHeaderData)
				if err != nil || n == 0 {
					break
				}
				cmd.Unpack(cmdHeaderData)

				var transferBuffer []byte
				if cmd.Direction == USBIPDirOut && cmd.TransferBufferLength > 0 {
					transferBuffer = make([]byte, cmd.TransferBufferLength)
					conn.Read(transferBuffer)
				}

				fmt.Printf("usbip cmd %x\n", cmd.Command)
				fmt.Printf("usbip seqnum %x\n", cmd.Seqnum)
				fmt.Printf("usbip devid %x\n", cmd.Devid)
				fmt.Printf("usbip direction %x\n", cmd.Direction)
				fmt.Printf("usbip ep %x\n", cmd.Ep)
				fmt.Printf("usbip flags %x\n", cmd.TransferFlags)
				fmt.Printf("usbip transfer buffer length %x\n", cmd.TransferBufferLength)
				fmt.Printf("usbip start %x\n", cmd.StartFrame)
				fmt.Printf("usbip number of packets %x\n", cmd.NumberOfPackets)
				fmt.Printf("usbip interval %x\n", cmd.Interval)
				fmt.Printf("usbip setup %s\n", BytesToString(cmd.Setup[:]))
				fmt.Printf("usbip transfer buffer %s\n", BytesToString(transferBuffer))

				usbReq := USBRequest{
					Seqnum:               cmd.Seqnum,
					Devid:                cmd.Devid,
					Direction:            cmd.Direction,
					Ep:                   cmd.Ep,
					Flags:                cmd.TransferFlags,
					TransferBufferLength: cmd.TransferBufferLength,
					NumberOfPackets:      cmd.NumberOfPackets,
					Interval:             cmd.Interval,
					Setup:                cmd.Setup,
					TransferBuffer:       transferBuffer,
				}

				c.USBDevices[0].SetConnection(conn)
				if baseDevice, ok := c.USBDevices[0].(interface{ HandleUSBRequest(USBDevice, USBRequest) }); ok {
					baseDevice.HandleUSBRequest(c.USBDevices[0], usbReq)
				}
			}
		}
		fmt.Println("Close connection")
		conn.Close()
	}
}
