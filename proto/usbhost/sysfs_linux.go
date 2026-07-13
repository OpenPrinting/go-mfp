// MFP - Miulti-Function Printers and scanners toolkit
// USB host API
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// USB host API, linux sysfs version

package usbhost

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	"github.com/OpenPrinting/go-mfp/internal/assert"
	"github.com/OpenPrinting/go-mfp/proto/usb"
)

const (
	// sysfsUSB is the path to USB hierarchy under the sysfs
	sysfsUSB = "/sys/bus/usb/devices/"
)

// Use sysfsListDevices on Linux if we run under the regular
// user, not as root.
//
// It is almost as accurate as libusb, but doesn't require root
// privileges.
//
// The only missed information at this case is the USB strings
// for non-default configuration. As we are concentrated on MFP,
// this limitation is very unlikely to affect us.
func init() {
	if os.Geteuid() != 0 {
		listDevices = sysfsListDevices
	}
}

// sysfsListDevices returns list of all connected USB devices.
//
// This is the version of ListDevices, that works on a top of
// the Linux sysfs.
func sysfsListDevices(withIEEE1284id bool) ([]DeviceInfo, error) {
	// Obtain list of devices
	devices, err := sysfsGetDeviceList()
	if err != nil {
		return nil, err
	}

	// Decode device list
	infos := make([]DeviceInfo, 0, len(devices))
	for _, name := range devices {
		info, err := name.LoadDeviceInfo(withIEEE1284id)

		if err != nil && sysfsIsFatal(err) {
			return nil, err
		}

		if err == nil {
			infos = append(infos, info)
		}
	}

	return infos, nil
}

// sysfsGetDeviceList returns list of USB devices under
// the /sys/bus/usb/devices/ directory
func sysfsGetDeviceList() ([]devname, error) {
	// List everything under /sys/bus/usb/devices/
	entries, err := os.ReadDir(sysfsUSB)
	if err != nil {
		return nil, err
	}

	// Filter out everything that is not device
	devices := []devname{}
	for _, ent := range entries {
		name := devname(ent.Name())
		if name.Valid() {
			devices = append(devices, name)
		}
	}

	return devices, nil
}

// sysfsIsFatal determines whether the USB error should be considered fatal
// (i.e., whether it should interrupt a major operation such as enumerating
// devices).
//
// Since devices may be unplugged during descriptor decoding or I/O errors may
// occur, we silently ignore certain error conditions related to these
// scenarios.
func sysfsIsFatal(err error) bool {
	if err == nil {
		return false
	}

	// Check for standard "File Does Not Exist" (translates to ENOENT)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	// Check for underlying low-level Linux syscall errors
	// (ENODEV or ENOENT)
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		var errno syscall.Errno
		if errors.As(pathErr.Err, &errno) {
			switch errno {
			case syscall.ENODEV, syscall.ENOENT:
				return false
			}

			return true
		}
	}

	return false
}

func sysfsOptional[T any](val T, err error) (T, error) {
	if sysfsIsFatal(err) {
		var zero T
		return zero, err
	}

	return val, nil
}

// devname contains an USB device name, relative to /sys/bus/usb/devices/
type devname string

// Valid reports if the device name is the really valid name of
// device.
//
// The /sys/bus/usb/devices/ directory contains many various entries:
//
//	$ ls /sys/bus/usb/devices/
//	1-0:1.0   1-11:1.0  1-7:1.2  1-9.4:1.0  3-1      3-1:1.3  usb3
//	1-10      1-7       1-9      1-9.4:1.1  3-1:1.0  4-0:1.0  usb4
//	1-10:1.0  1-7:1.0   1-9:1.0  2-0:1.0    3-1:1.1  usb1
//	1-11      1-7:1.1   1-9.4    3-0:1.0    3-1:1.2  usb2
//
// This function filters out the names that doesn't represent devices.
func (name devname) Valid() bool {
	// These are USB host controllers.
	if strings.HasPrefix(string(name), "usb") {
		// The "usb" prefix must be accompanied by a
		// non-empty sequence of digits.
		suffix := string(name[3:])
		if len(suffix) == 0 {
			return false
		}

		for _, r := range suffix {
			if !unicode.IsDigit(r) {
				return false
			}
		}

		return true
	}

	// Locate the single dash separating the bus number and
	// the port chain
	dashIdx := strings.IndexByte(string(name), '-')
	if dashIdx <= 0 || dashIdx == len(name)-1 {
		// No dash, dash is at the start, or dash is at the very end
		return false
	}
	if strings.IndexByte(string(name[dashIdx+1:]), '-') != -1 {
		// More than one dash found
		return false
	}

	// Validate the left part: Bus number (must be digits only)
	busPart := name[:dashIdx]
	for _, r := range busPart {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	// Validate the right part: Port chain (must be digits and dots only)
	portsPart := name[dashIdx+1:]

	// A dot cannot immediately follow the dash, nor can it be at the
	// very end
	if portsPart[0] == '.' || portsPart[len(portsPart)-1] == '.' {
		return false
	}

	lastWasDot := false
	for _, r := range portsPart {
		if r == '.' {
			if lastWasDot {
				// Consecutive dots are invalid (e.g., "1-1..2")
				return false
			}
			lastWasDot = true
		} else if unicode.IsDigit(r) {
			lastWasDot = false
		} else {
			// Any other character (letters, colons, etc.) is
			// invalid
			return false
		}
	}

	return true
}

// LoadDeviceInfo loads DeviceInfo of the USB device.
func (name devname) LoadDeviceInfo(withIEEE1284id bool) (
	info DeviceInfo, err error) {

	// Decode location and usb.DeviceDescriptor
	info.Loc, err = name.LoadLocation()
	if err == nil {
		info.Desc.BCDUSB, err = name.LoadBSDUSB()
	}
	if err == nil {
		info.Desc.Speed, err = name.LoadSpeed()
	}
	if err == nil {
		info.Desc.BDeviceClass, err = name.LoadHex8("/bDeviceClass")
	}
	if err == nil {
		info.Desc.BDeviceSubClass, err = name.LoadHex8("/bDeviceSubClass")
	}
	if err == nil {
		info.Desc.BDeviceProtocol, err = name.LoadHex8("/bDeviceProtocol")
	}
	if err == nil {
		var sz int
		sz, err = name.LoadInt("/bMaxPacketSize0")
		info.Desc.BMaxPacketSize = uint8(sz)
	}
	if err == nil {
		info.Desc.IDVendor, err = name.LoadHex16("/idVendor")
	}
	if err == nil {
		info.Desc.IDProduct, err = name.LoadHex16("/idProduct")
	}
	if err == nil {
		var bcd uint16
		bcd, err = name.LoadHex16("/bcdDevice")
		info.Desc.BCDDevice = usb.Version(bcd)
	}
	if err == nil {
		info.Desc.IManufacturer, err = sysfsOptional(
			name.LoadString("/manufacturer"))
	}
	if err == nil {
		info.Desc.IProduct, err = sysfsOptional(
			name.LoadString("/product"))
	}
	if err == nil {
		info.Desc.ISerialNumber, err = sysfsOptional(
			name.LoadString("/serial"))
	}
	if err == nil {
		info.Desc.Configurations, err = sysfsOptional(
			name.LoadConfigurations(withIEEE1284id))
	}

	assert.NoError(err)

	return info, err
}

// LoadConfigurations reads a sysfs binary 'descriptors' file and builds the descriptor hierarchy.
func (name devname) LoadConfigurations(withIEEE1284id bool) (
	[]usb.ConfigurationDescriptor, error) {

	// Get bConfigurationValue of current configuration
	bConfigurationValue, err := name.LoadInt("/bConfigurationValue")
	if err != nil {
		return nil, err
	}

	// Load /descriptors file
	const path = "/descriptors"
	data, err := name.LoadBinary(path)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	// 1. Skip Device Descriptor (always the first 18 bytes)
	if reader.Len() < 18 {
		return nil, fmt.Errorf("%s: file too short", path)
	}
	reader.Seek(18, io.SeekStart)

	var configs []usb.ConfigurationDescriptor

	// Temporary map to group interfaces by bInterfaceNumber within the
	// current configuration
	var currentInterfaces map[uint8]*usb.Interface
	var currentConfigIdx = -1
	var currentConfigIsActive = false
	var currentInterface *usb.Interface

	// 2. Sequentially parse the remaining descriptor stream
	for reader.Len() > 2 {
		// Peek at the length and type of the next descriptor
		length, _ := reader.ReadByte()
		descType, _ := reader.ReadByte()

		// Move the reader pointer back to read the whole descriptor via binary.Read
		reader.Seek(-2, io.SeekCurrent)

		if reader.Len() < int(length) {
			return nil, fmt.Errorf("%s: configuration descriptor truncated", path)
		}

		// Allocate a dedicated buffer for the current descriptor block
		descBuf := make([]byte, length)
		reader.Read(descBuf)
		bufReader := bytes.NewReader(descBuf)

		switch usb.DescriptorType(descType) {
		case usb.DescriptorConfiguration:
			// Temporary structure matching the USB specification layout
			var raw struct {
				BLength             uint8
				BDescriptorType     uint8
				WTotalLength        uint16
				BNumInterfaces      uint8
				BConfigurationValue uint8
				IConfiguration      uint8
				BMAttributes        uint8
				MaxPower            uint8
			}
			if err := binary.Read(bufReader, binary.LittleEndian, &raw); err != nil {
				return nil, fmt.Errorf("%s: configuration descriptor truncated", path)
			}

			config := usb.ConfigurationDescriptor{
				BConfigurationValue: raw.BConfigurationValue,
				IConfiguration:      "", // Kept empty as requested
				BMAttributes:        usb.ConfAttributes(raw.BMAttributes),
				MaxPower:            raw.MaxPower,
				Interfaces:          []usb.Interface{},
			}
			configs = append(configs, config)

			currentConfigIdx = len(configs) - 1
			currentInterfaces = make(map[uint8]*usb.Interface)
			currentInterface = nil

			// Configuration name available only for currently active configuration
			currentConfigIsActive = int(raw.BConfigurationValue) == bConfigurationValue
			if currentConfigIsActive {
				config.IConfiguration, err = sysfsOptional(
					name.LoadString("./configuration"))
				if err != nil {
					return nil, err
				}
			}

		case usb.DescriptorInterface:
			if currentConfigIdx == -1 {
				err := fmt.Errorf("%s: interface descriptor comes before configuration descriptor", path)
				return nil, err
			}

			var raw struct {
				BLength            uint8
				BDescriptorType    uint8
				BInterfaceNumber   uint8
				BAlternateSetting  uint8
				BNumEndpoints      uint8
				BInterfaceClass    uint8
				BInterfaceSubClass uint8
				BInterfaceProtocol uint8
				IInterface         uint8
			}
			if err := binary.Read(bufReader, binary.LittleEndian, &raw); err != nil {
				err := fmt.Errorf("%s: interface descriptor truncated", path)
				return nil, err
			}

			ifaceDesc := usb.InterfaceDescriptor{
				BInterfaceClass:    raw.BInterfaceClass,
				BInterfaceSubClass: raw.BInterfaceSubClass,
				BInterfaceProtocol: raw.BInterfaceProtocol,
				BAlternateSetting:  raw.BAlternateSetting,
				Endpoints:          []usb.EndpointDescriptor{},
			}

			// Strings are only available for the currently active configuration
			if currentConfigIsActive {
				subpath := fmt.Sprintf("/%s:%d.%d",
					name, bConfigurationValue, raw.BInterfaceNumber)

				// Interface name is only related to the zero alt
				// setting
				if raw.BAlternateSetting == 0 {
					ifaceDesc.IInterface, err = sysfsOptional(
						name.LoadString(subpath + "/interface"))
					if err != nil {
						return nil, err
					}
				}

				// IEEE-1284 device ID is common for all alts, but
				// returned only for 7/1/1 and 7/1/2 devices
				if withIEEE1284id &&
					raw.BInterfaceClass == 7 &&
					raw.BInterfaceSubClass == 1 &&
					(raw.BInterfaceProtocol == 1 ||
						raw.BInterfaceProtocol == 2) {
					ifaceDesc.IEEE1284DeviceID, err = sysfsOptional(
						name.LoadString(subpath + "/ieee1284_id"))
					if err != nil {
						return nil, err
					}
				}
			}

			// If this interface number hasn't been encountered yet, initialize it
			currentInterface = currentInterfaces[raw.BInterfaceNumber]
			if currentInterface == nil {
				currentInterface = &usb.Interface{
					BInterfaceNumber: raw.BInterfaceNumber,
					AltSettings:      []usb.InterfaceDescriptor{},
				}
				currentInterfaces[raw.BInterfaceNumber] = currentInterface
			}

			// And append the AltSetting
			currentInterface.AltSettings = append(currentInterface.AltSettings, ifaceDesc)

		case usb.DescriptorEndpoint:
			if currentInterface == nil {
				err := fmt.Errorf("%s: endpoint descriptor comes before interface descriptor", path)
				return nil, err
			}

			var raw struct {
				BLength          uint8
				BDescriptorType  uint8
				BEndpointAddress uint8
				BMAttributes     uint8
				WMaxPacketSize   uint16
				BInterval        uint8
			}
			if err := binary.Read(bufReader, binary.LittleEndian, &raw); err != nil {
				err := fmt.Errorf("%s: endpoint descriptor truncated", path)
				return nil, err
			}

			epDesc := usb.EndpointDescriptor{
				BMAttributes:   usb.EndpointAttributes(raw.BMAttributes),
				WMaxPacketSize: raw.WMaxPacketSize & 0x07FF, // Masking bits 0-10 for actual packet size in bytes
			}

			if (raw.BEndpointAddress & 0x80) == 0 {
				epDesc.Type = usb.EndpointOut
			} else {
				epDesc.Type = usb.EndpointIn
			}

			// Since descriptors are linear streams, find the latest added AltSetting
			// across current active interfaces to append this endpoint to
			lastIdx := len(currentInterface.AltSettings) - 1
			currentInterface.AltSettings[lastIdx].Endpoints = append(currentInterface.AltSettings[lastIdx].Endpoints, epDesc)

		default:
			// Automatically skip class-specific or vendor-specific descriptors (HID, Hub, Audio, etc.)
			continue
		}

		// Save the accumulated and grouped interfaces back to the configuration block
		if currentConfigIdx != -1 {
			var sortedIfaces []usb.Interface
			for _, iface := range currentInterfaces {
				sortedIfaces = append(sortedIfaces, *iface)
			}
			sort.Slice(sortedIfaces, func(i, j int) bool {
				return sortedIfaces[i].BInterfaceNumber < sortedIfaces[j].BInterfaceNumber
			})
			configs[currentConfigIdx].Interfaces = sortedIfaces
		}
	}

	return configs, nil
}

// LoadBSDUSB loads USB version (bsdUSB in USB terms)
func (name devname) LoadBSDUSB() (ver usb.Version, err error) {
	// Load the /version file
	versionStr, err := name.LoadString("/version")
	if err != nil {
		return 0, err
	}

	// Pre-format an error
	err = fmt.Errorf("%s: invalid /version", versionStr)

	// Split version string (major.minor) into major and minor parts
	dotIdx := strings.IndexByte(versionStr, '.')
	if dotIdx <= 0 || dotIdx == len(versionStr)-1 {
		return 0, err
	}

	majorStr := versionStr[:dotIdx]
	minorStr := versionStr[dotIdx+1:]

	// Parse Major and minor versions
	var major, minor uint16

	if tmp, err2 := strconv.ParseUint(majorStr, 16, 16); err2 == nil {
		major = uint16(tmp)
	} else {
		return 0, err
	}

	if tmp, err2 := strconv.ParseUint(minorStr, 16, 16); err2 == nil {
		minor = uint16(tmp)
	} else {
		return 0, err
	}

	// Combine into BCD layout: Major byte (0xMM) and Minor byte (0xNN)
	bcd := usb.Version((major << 8) | (minor & 0xFF))
	return bcd, nil
}

// LoadBSDUSB loads USB Speed.
func (name devname) LoadSpeed() (speed usb.Speed, err error) {
	s, err := name.LoadString("/speed")
	if err != nil {
		return 0, err
	}

	speed = usb.SpeedUnknown
	switch s {
	case "1.5":
		speed = usb.SpeedLow
	case "12":
		speed = usb.SpeedFull
	case "480":
		speed = usb.SpeedHigh
	case "5000":
		speed = usb.SpeedSuper
	case "10000":
		speed = usb.SpeedSuperPlus
	case "20000":
		speed = usb.SpeedSuperX2
	}

	return
}

// sysfsLoadLocation loads USB device location.
func (name devname) LoadLocation() (Location, error) {
	busnum, err := name.LoadInt("/busnum")
	if err != nil {
		return Location{}, err
	}

	devnum, err := name.LoadInt("/devnum")
	if err != nil {
		return Location{}, err
	}

	return Location{Bus: busnum, Dev: devnum}, nil
}

// LoadBibary loads binary data from USB device
func (name devname) LoadBinary(path string) ([]byte, error) {
	return os.ReadFile(sysfsUSB + string(name) + path)
}

// LoadString loads string from USB device.
// It automatically calls strings.TrimSpace on success.
func (name devname) LoadString(path string) (string, error) {
	data, err := name.LoadBinary(path)
	return strings.TrimSpace(string(data)), err
}

// LoadInt loads integer parameter.
func (name devname) LoadInt(path string) (int, error) {
	s, err := name.LoadString(path)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		err = fmt.Errorf("%s: bad integer", path)
	}

	return i, err
}

// LoadHex8 loads hex-encoded uint8 parameter.
func (name devname) LoadHex8(path string) (uint8, error) {
	s, err := name.LoadString(path)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		err = fmt.Errorf("%s: bad uint8", path)
	}

	return uint8(i), err
}

// LoadInt loads hex-encoded uint16 parameter.
func (name devname) LoadHex16(path string) (uint16, error) {
	s, err := name.LoadString(path)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		err = fmt.Errorf("%s: bad hex-16", path)
	}

	return uint16(i), err
}
