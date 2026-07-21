// MFP - Multi-Function Printers and scanners toolkit
// Virtual USB/IP device emulator for testing and fuzzing
//
// Copyright (C) 2025 and up by GO-MFP authors.
// See LICENSE for license terms and conditions
//
// The virtual device

package usbip

import (
	"errors"
	"fmt"
	"math"
	"syscall"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// Device represents the virtual USB device.
//
// Device lifetime is the following:
//   - use [NewDevice] to create the new Device.
//   - fill required callbacks
//   - use [Server.AddDevice] to add device to the server
type Device struct {
	// Device description
	Descriptor usb.DeviceDescriptor // Device descriptor

	// User-defined callbacks
	OnConfigurationChange func() // Called when configuration changes

	// Device state
	configuration  int               // Current configuration
	altSettings    []int             // Per-interface current alt setting
	maxAltSettings []int             // Max alt setting, by interface
	endpoints      []*Endpoint       // Endpoints of current configuration
	endpointsTree  [][][][]*Endpoint // All Endpoints, defined by descriptor
	strings        []string          // Content of string descriptors
	stringsmap     map[string]int    // Indices of strings descriptors
}

// NewDevice creates a new device, based on the provided device descriptor.
func NewDevice(desc usb.DeviceDescriptor) (*Device, error) {
	// Validate device descriptor
	confs := len(desc.Configurations)

	switch {
	case confs == 0:
		return nil, errors.New("Device has no configurations")

	case confs > usb.MaxConfigurations:
		err := fmt.Errorf("Device has too many (%d) configurations",
			confs)
		return nil, err
	}

	for confno, conf := range desc.Configurations {
		iffs := len(conf.Interfaces)

		switch {
		case iffs == 0:
			err := fmt.Errorf(
				"Configuration %d has no interfaces",
				confno)
			return nil, err

		case iffs > usb.MaxInterfaces:
			err := fmt.Errorf(
				"Configuration %d has too many [%d] interfaces",
				confno, iffs)
			return nil, err
		}

		cntEndpoints := 1 // Reserved for configuration endpoint
		for _, iff := range conf.Interfaces {
			cntEndpoints += iff.CntEndpoints()
			if cntEndpoints > usb.MaxEndpoints {
				err := fmt.Errorf(
					"Configuration %d has too many (%d) endpoints",
					confno, cntEndpoints)
				return nil, err
			}
		}
	}

	// Initialize the Device structure
	dev := &Device{
		Descriptor: desc,
		stringsmap: make(map[string]int),
	}

	// Populate endpointsTree
	dev.endpointsTree = make([][][][]*Endpoint, len(desc.Configurations))
	for confno, conf := range desc.Configurations {
		iffs := len(conf.Interfaces)
		dev.endpointsTree[confno] = make([][][]*Endpoint, iffs)

		for iffno, iff := range conf.Interfaces {
			alts := len(iff.AltSettings)
			dev.endpointsTree[confno][iffno] =
				make([][]*Endpoint, alts)

			for altno, alt := range iff.AltSettings {
				eps := len(alt.Endpoints)
				dev.endpointsTree[confno][iffno][altno] =
					make([]*Endpoint, eps)

				for epno, ep := range alt.Endpoints {
					dev.endpointsTree[confno][iffno][altno][epno] =
						NewEndpoint(ep)
				}
			}
		}
	}

	// Populate strings
	dev.addstring("") // Reserve index 0
	dev.addstring(desc.IManufacturer)
	dev.addstring(desc.IProduct)
	dev.addstring(desc.ISerialNumber)

	for _, conf := range desc.Configurations {
		dev.addstring(conf.IConfiguration)

		for _, iff := range conf.Interfaces {
			for _, alt := range iff.AltSettings {
				dev.addstring(alt.IInterface)
			}
		}
	}

	// Initialize current configuration
	dev.setConfiguration(1)

	return dev, nil
}

// MustNewDevice works like [NewDevice], but it panics in a case
// of an error. Note, the only reason here to fail is the malformed
// [USBDeviceDescriptor].
func MustNewDevice(desc usb.DeviceDescriptor) *Device {
	dev, err := NewDevice(desc)
	assert.NoError(err)
	return dev
}

// EndpointsByFunc returns all device endpoints for which the matching
// function returns true.
func (dev *Device) EndpointsByFunc(
	match func(usb.InterfaceDescriptor) bool) []*Endpoint {

	found := []*Endpoint{}

	for confno, conf := range dev.Descriptor.Configurations {
		for iffno, iff := range conf.Interfaces {
			for altno, alt := range iff.AltSettings {
				if match(alt) {
					endpoints := dev.endpointsTree[confno][iffno][altno]
					found = append(found, endpoints...)
				}
			}
		}
	}

	return found
}

// EndpointsByClassID returns all device endpoints that match the
// specified [usb.ClassID]
func (dev *Device) EndpointsByClassID(id usb.ClassID) []*Endpoint {
	match := func(alt usb.InterfaceDescriptor) bool {
		return alt.Match(id)
	}

	return dev.EndpointsByFunc(match)
}

// EndpointsByClass returns all device endpoints that belongs to the
// device interfaces with the specified class/subclass/protocol combination.
//
// It makes sense if all interfaces of the same class/subclass/protocol
// are functionally equal.
func (dev *Device) EndpointsByClass(class, subclass, proto uint8) []*Endpoint {
	found := []*Endpoint{}

	for confno, conf := range dev.Descriptor.Configurations {
		for iffno, iff := range conf.Interfaces {
			for altno, alt := range iff.AltSettings {
				if class == alt.BInterfaceClass &&
					subclass == alt.BInterfaceSubClass &&
					proto == alt.BInterfaceProtocol {

					endpoints := dev.endpointsTree[confno][iffno][altno]
					found = append(found, endpoints...)
				}
			}
		}
	}

	return found
}

// GetStatus returns the USB device status.
func (dev *Device) getStatus() ([]byte, syscall.Errno) {
	data := make([]byte, 2)

	conf := dev.Descriptor.Configurations[dev.configuration-1]

	if conf.BMAttributes&usb.ConfAttrSelfPowered != 0 {
		data[1] |= 0x01
	}

	if conf.BMAttributes&usb.ConfAttrRemoteWakeup != 0 {
		data[1] |= 0x02
	}

	return data, 0
}

// GetDescriptor returns USB descriptor, in the USB wire representation.
// If type or index is invalid, it returns nil and syscall.EPIPE.
func (dev *Device) getDescriptor(t usb.DescriptorType, i int) ([]byte, syscall.Errno) {
	switch t {
	case usb.DescriptorDevice:
		if i > 0 {
			return nil, syscall.EPIPE
		}

		desc := dev.Descriptor
		enc := newEncoder(18)

		enc.PutU8(18)                          // bLength
		enc.PutU8(uint8(usb.DescriptorDevice)) // bDescriptorType
		enc.PutLE16(uint16(desc.BCDUSB))       // bcdUSB
		enc.PutU8(desc.BDeviceClass)           // bDeviceClass
		enc.PutU8(desc.BDeviceSubClass)        // bDeviceSubClass
		enc.PutU8(desc.BDeviceProtocol)        // bDeviceProtocol
		enc.PutU8(desc.BMaxPacketSize)         // bMaxPacketSize
		enc.PutLE16(desc.IDVendor)             // idVendor
		enc.PutLE16(desc.IDProduct)            // idProduct
		enc.PutLE16(uint16(desc.BCDDevice))    // bcdDevice

		i := dev.getstring(desc.IManufacturer)
		enc.PutU8(uint8(i)) // iManufacturer

		i = dev.getstring(desc.IProduct)
		enc.PutU8(uint8(i)) // iProduct

		i = dev.getstring(desc.ISerialNumber)
		enc.PutU8(uint8(i)) // iSerialNumber

		enc.PutU8(uint8(len(desc.Configurations))) // bNumConfigurations

		return enc.Bytes(), 0

	case usb.DescriptorConfiguration:
		if i > len(dev.Descriptor.Configurations) {
			return nil, syscall.EPIPE
		}

		// Encode Configuration Descriptor
		conf := dev.Descriptor.Configurations[i]
		enc := newEncoder(256)

		enc.PutU8(9)                                  // bLength
		enc.PutU8(uint8(usb.DescriptorConfiguration)) // bDescriptorType
		enc.PutLE16(0)                                // wTotalLength, reserved
		enc.PutU8(uint8(len(conf.Interfaces)))        // bNumInterfaces
		enc.PutU8(uint8(i + 1))                       // bConfigurationValue

		i := dev.getstring(conf.IConfiguration)
		enc.PutU8(uint8(i)) // iConfiguration

		attrs := conf.BMAttributes
		attrs |= usb.ConfAttrReserved
		enc.PutU8(uint8(attrs)) // bmAttributes

		enc.PutU8(conf.BMaxPower) // bMaxPower

		// Encode interfaces and endpoints
		epnum := 1
		for iffno, iff := range conf.Interfaces {
			for altno, alt := range iff.AltSettings {
				cntEndpoints := iff.CntEndpoints()

				enc.PutU8(9)                              // bLength
				enc.PutU8(uint8(usb.DescriptorInterface)) // bDescriptorType
				enc.PutU8(uint8(iffno))                   // bInterfaceNumber
				enc.PutU8(uint8(altno))                   // bAlternateSetting
				enc.PutU8(uint8(cntEndpoints))            // bNumEndpoints
				enc.PutU8(alt.BInterfaceClass)            // bInterfaceClass
				enc.PutU8(alt.BInterfaceSubClass)         // bInterfaceSubClass
				enc.PutU8(alt.BInterfaceProtocol)         // bInterfaceProtocol

				i := dev.getstring(alt.IInterface)
				enc.PutU8(uint8(i)) // iInterface

				endpoints := dev.endpointsTree[i][iffno][altno]
				for _, ep := range endpoints {
					ty := ep.Type()
					if ty == usb.EndpointIn || ty == usb.EndpointInOut {
						enc.PutU8(7) // bLength
						enc.PutU8(   // bDescriptorType
							uint8(usb.DescriptorEndpoint))

						addr := epnum | 0x80
						epnum++

						enc.PutU8(uint8(addr))            // bEndpointAddress
						enc.PutU8(uint8(ep.Attrs()))      // bmAttributes
						enc.PutLE16(uint16(ep.PktSize())) // bmAttributes
						enc.PutU8(0)                      // bInterval
					}

					if ty == usb.EndpointOut || ty == usb.EndpointInOut {
						enc.PutU8(7) // bLength
						enc.PutU8(   // bDescriptorType
							uint8(usb.DescriptorEndpoint))

						addr := epnum
						epnum++

						enc.PutU8(uint8(addr))            // bEndpointAddress
						enc.PutU8(uint8(ep.Attrs()))      // bmAttributes
						enc.PutLE16(uint16(ep.PktSize())) // bmAttributes
						enc.PutU8(0)                      // bInterval
					}
				}
			}
		}

		// Patch wTotalLength
		data := enc.Bytes()
		length := len(data)
		data[2] = uint8(length)
		data[3] = uint8(length >> 8)

		return data, 0

	case usb.DescriptorString:
		// String Descriptor Zero is special: it contains list of
		// supported languages.
		if i == 0 {
			enc := newEncoder(4)
			enc.PutU8(8)                           // bLength
			enc.PutU8(uint8(usb.DescriptorString)) // bDescriptorType
			enc.PutLE16(0x0409)                    // wLANGID[0]
			return enc.Bytes(), 0
		}

		if i >= len(dev.strings) {
			return nil, syscall.EPIPE
		}

		s := []byte(dev.strings[i])
		enc := newEncoder(4)
		enc.PutU8(uint8(len(s)*2 + 2))         // bLength
		enc.PutU8(uint8(usb.DescriptorString)) // bDescriptorType
		for _, c := range s {
			enc.PutLE16(uint16(c)) // bString
		}
		return enc.Bytes(), 0

	}

	return nil, syscall.EPIPE
}

// GetConfiguration returns the current configuration.
func (dev *Device) GetConfiguration() ([]byte, syscall.Errno) {
	return []byte{uint8(dev.configuration)}, 0
}

// SetConfiguration selects the current configuration.
func (dev *Device) setConfiguration(n int) ([]byte, syscall.Errno) {
	if n > len(dev.Descriptor.Configurations) {
		return nil, syscall.EPIPE
	}

	if n != 0 {
		// Save Configuration
		dev.configuration = n

		confno := n - 1
		conf := dev.Descriptor.Configurations[confno]

		// Update dev.altSettings
		dev.altSettings = make([]int, len(conf.Interfaces))
		dev.maxAltSettings = make([]int, len(conf.Interfaces))

		for i, iff := range conf.Interfaces {
			dev.maxAltSettings[i] = len(iff.AltSettings)
		}

		// Update dev.endpoints
		// dev.endpoints[0] is reserved.
		dev.endpoints = make([]*Endpoint, 0, 16)
		dev.endpoints = append(dev.endpoints, nil)

		for iffno, iff := range conf.Interfaces {
			for altno, alt := range iff.AltSettings {
				for epno, ep := range alt.Endpoints {
					epp := dev.endpointsTree[confno][iffno][altno][epno]

					switch ep.Type {
					case usb.EndpointIn, usb.EndpointOut:
						dev.endpoints = append(dev.endpoints,
							epp)
					case usb.EndpointInOut:
						dev.endpoints = append(dev.endpoints,
							epp, epp)
					}
				}
			}
		}

		// Call user callback
		if dev.OnConfigurationChange != nil {
			dev.OnConfigurationChange()
		}
	}

	return []byte{}, 0
}

// submit routes protoIOSubmitRequest to the endpoint
func (dev *Device) submit(rq *protoIOSubmitRequest) syscall.Errno {
	// Lookup the endpoint
	var ep *Endpoint
	if int(rq.Endpoint) <= len(dev.endpoints) {
		ep = dev.endpoints[rq.Endpoint]
	}

	if ep == nil {
		return syscall.EPIPE
	}

	// Check for direction mismatch
	if (rq.Input && ep.Type() == usb.EndpointOut) ||
		(!rq.Input && ep.Type() == usb.EndpointIn) {
		return syscall.EPIPE
	}

	// Route the request
	return ep.submit(rq)
}

// unlink routes protoIOUnlinkRequest to the endpoint
func (dev *Device) unlink(rq *protoIOUnlinkRequest) syscall.Errno {
	// Lookup the endpoint
	var ep *Endpoint
	if int(rq.Endpoint) <= len(dev.endpoints) {
		ep = dev.endpoints[rq.Endpoint]
	}

	if ep == nil {
		return syscall.EPIPE
	}

	// Route the request
	return ep.unlink(rq)
}

// shutdown cancels app pending protoIOSubmitRequest at all endpoints.
func (dev *Device) shutdown() {
	for _, ep := range dev.endpoints[1:] {
		ep.shutdown()
	}
}

// confIndex returns *usb.ConfigurationDescriptor by index
// (not by bConfigurationValue!), or nil if index is out of range.
func (dev *Device) confIndex(confno int) *usb.ConfigurationDescriptor {
	desc := &dev.Descriptor
	if confno >= 0 && confno < len(desc.Configurations) {
		return &desc.Configurations[confno]
	}
	return nil
}

// ifIndex returns *usb.Interface by index (not by bConfigurationValue
// and bInterfaceNumber), or nil, if some index is out of range.
func (dev *Device) ifIndex(confno, iffno int) *usb.Interface {
	conf := dev.confIndex(confno)
	if conf != nil && iffno >= 0 && iffno < len(conf.Interfaces) {
		return &conf.Interfaces[iffno]
	}
	return nil
}

// altIndex returns *usb.InterfaceDescriptor by index (not by
// bConfigurationValue, bInterfaceNumber and bAlternateSetting),
// or nil, if some index is out of range.
func (dev *Device) altIndex(confno, iffno, altno int) *usb.InterfaceDescriptor {
	iff := dev.ifIndex(confno, iffno)
	if iff != nil && altno >= 0 && altno < len(iff.AltSettings) {
		return &iff.AltSettings[altno]
	}
	return nil
}

// GetInterfaceStatus returns the USB interface status.
func (dev *Device) getInterfaceStatus(ifn int) ([]byte, syscall.Errno) {
	if ifn >= len(dev.altSettings) {
		return nil, syscall.EPIPE
	}

	return make([]byte, 2), 0
}

// GetInterface returns alternate setting selection on the given interface.
func (dev *Device) getInterface(ifn int) ([]byte, syscall.Errno) {
	if ifn >= len(dev.altSettings) {
		return nil, syscall.EPIPE
	}

	return []byte{uint8(dev.altSettings[ifn])}, 0
}

// SetInterface selects alternate setting selection on the given interface.
func (dev *Device) setInterface(ifn, alt int) ([]byte, syscall.Errno) {
	if ifn >= len(dev.altSettings) {
		return nil, syscall.EPIPE
	}

	if alt > dev.maxAltSettings[ifn] {
		return nil, syscall.EPIPE
	}

	dev.altSettings[ifn] = alt
	return []byte{}, 0
}

// addstring adds string to strings descriptor.
func (dev *Device) addstring(s string) {
	if len(dev.strings) < math.MaxUint8 {
		if len(s) > usb.MaxStringLength {
			s = s[:usb.MaxStringLength]
		}

		if _, found := dev.stringsmap[s]; !found {
			dev.stringsmap[s] = len(dev.strings)
			dev.strings = append(dev.strings, s)
		}
	}
}

// getstring returns the string descriptor index for the string.
// If index is not available, 0 will be returned (which is OK for USB).
func (dev *Device) getstring(s string) int {
	if len(s) > usb.MaxStringLength {
		s = s[:usb.MaxStringLength]
	}

	return dev.stringsmap[s]
}
