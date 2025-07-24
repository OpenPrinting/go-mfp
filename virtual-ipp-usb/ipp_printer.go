package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config holds settings used to initialize the virtual IPP-over-USB device.
type Config struct {
	IPPServerURL string `json:"ipp_server_url"`
	DeviceName   string `json:"device_name"`
	VendorID     string `json:"vendor_id"`
	ProductID    string `json:"product_id"`
	Manufacturer string `json:"manufacturer"`
	Product      string `json:"product"`
	Serial       string `json:"serial"`
	ListenIP     string `json:"listen_ip"`
	ListenPort   int    `json:"listen_port"`
	Debug        bool   `json:"debug"`
}

// IPPOverUSBDevice simulates a USB printer that forwards data to an IPP server.
type IPPOverUSBDevice struct {
	BaseUSBDevice
	config           Config
	serverURL        string
	deviceName       string
	vendorID         uint16
	productID        uint16
	deviceDescriptor DeviceDescriptor
	configurations   []DeviceConfiguration
	tcpConnection    net.Conn
	tcpConnected     bool
	connectionLock   sync.Mutex
	pendingResponse  []byte
}

// NewIPPOverUSBDevice creates a new IPP-over-USB device from the given config file.
func NewIPPOverUSBDevice(configFile string) (*IPPOverUSBDevice, error) {

	device := &IPPOverUSBDevice{}

	config, err := device.loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	device.config = config
	device.serverURL = config.IPPServerURL
	device.deviceName = config.DeviceName

	// Handle hex string vendor/product IDs from config
	vendorID, err := parseID(config.VendorID, 0x03F0)
	if err != nil {
		return nil, fmt.Errorf("invalid vendor_id: %v", err)
	}
	device.vendorID = vendorID

	productID, err := parseID(config.ProductID, 0x1234)
	if err != nil {
		return nil, fmt.Errorf("invalid product_id: %v", err)
	}
	device.productID = productID

	device.deviceDescriptor = device.createDeviceDescriptor()
	device.configurations = device.createConfigurations()

	device.GenerateRawConfiguration(device)

	fmt.Printf("IPP over USB Proxy Device\n")
	fmt.Printf("Configuration: %s\n", configFile)
	fmt.Printf("IPP Server URL: %s\n", device.serverURL)
	fmt.Printf("Device: %s\n", device.deviceName)
	fmt.Printf("Vendor ID: 0x%04X, Product ID: 0x%04X\n", device.vendorID, device.productID)

	return device, nil
}

func (d *IPPOverUSBDevice) loadConfig(configFile string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		// Creating a default config if file doesn't exist
		defaultConfig := Config{
			IPPServerURL: "http://localhost:631/ipp/print",
			DeviceName:   "Virtual IPP Printer",
			VendorID:     "0x03F0",
			ProductID:    "0x1234",
			Manufacturer: "Virtual",
			Product:      "IPP-USB Proxy",
			Serial:       "VIP001",
			ListenIP:     "0.0.0.0",
			ListenPort:   3240,
			Debug:        true,
		}

		jsonData, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return defaultConfig, err
		}

		err = ioutil.WriteFile(configFile, jsonData, 0644)
		if err != nil {
			return defaultConfig, err
		}

		fmt.Printf("Created default config file: %s\n", configFile)
		return defaultConfig, nil
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func parseID(idStr string, defaultValue uint16) (uint16, error) {
	if idStr == "" {
		return defaultValue, nil
	}

	if strings.HasPrefix(idStr, "0x") || strings.HasPrefix(idStr, "0X") {
		val, err := strconv.ParseUint(idStr[2:], 16, 16)
		if err != nil {
			return 0, err
		}
		return uint16(val), nil
	}

	val, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(val), nil
}

func (d *IPPOverUSBDevice) createDeviceDescriptor() DeviceDescriptor {
	return DeviceDescriptor{
		BLength:            18,
		BDescriptorType:    1,
		BcdUSB:             0x0200,
		BDeviceClass:       0x00,
		BDeviceSubClass:    0x00,
		BDeviceProtocol:    0x00,
		BMaxPacketSize0:    0x40,
		IDVendor:           d.vendorID,
		IDProduct:          d.productID,
		BcdDevice:          0x0100,
		IManufacturer:      1,
		IProduct:           2,
		ISerialNumber:      3,
		BNumConfigurations: 1,
	}
}

func (d *IPPOverUSBDevice) createConfigurations() []DeviceConfiguration {
	// First interface - IPP
	ippInterface := InterfaceDescriptor{
		BLength:            9,
		BDescriptorType:    4,
		BInterfaceNumber:   0,
		BAlternateSetting:  0,
		BNumEndpoints:      2,
		BInterfaceClass:    0x07, // Printer class
		BInterfaceSubClass: 0x01,
		BInterfaceProtocol: 0x04, // IPP-over-USB protocol
		IInterface:         0,
	}

	bulkOutEndpoint := EndpointDescriptor{
		BLength:          7,
		BDescriptorType:  5,
		BEndpointAddress: 0x01,
		BmAttributes:     0x02, // Bulk
		WMaxPacketSize:   0x0200,
		BInterval:        0x00,
	}

	bulkInEndpoint := EndpointDescriptor{
		BLength:          7,
		BDescriptorType:  5,
		BEndpointAddress: 0x82,
		BmAttributes:     0x02, // Bulk
		WMaxPacketSize:   0x0200,
		BInterval:        0x00,
	}

	ippInterface.Endpoints = []EndpointDescriptor{bulkOutEndpoint, bulkInEndpoint}

	// Second interface
	ippInterface2 := InterfaceDescriptor{
		BLength:            9,
		BDescriptorType:    4,
		BInterfaceNumber:   1,
		BAlternateSetting:  0,
		BNumEndpoints:      2,
		BInterfaceClass:    0x07, // Printer class
		BInterfaceSubClass: 0x01,
		BInterfaceProtocol: 0x04, // IPP-over-USB protocol
		IInterface:         0,
	}

	bulkOutEndpoint2 := EndpointDescriptor{
		BLength:          7,
		BDescriptorType:  5,
		BEndpointAddress: 0x03,
		BmAttributes:     0x02, // Bulk
		WMaxPacketSize:   0x0200,
		BInterval:        0x00,
	}

	bulkInEndpoint2 := EndpointDescriptor{
		BLength:          7,
		BDescriptorType:  5,
		BEndpointAddress: 0x84,
		BmAttributes:     0x02, // Bulk
		WMaxPacketSize:   0x0200,
		BInterval:        0x00,
	}

	ippInterface2.Endpoints = []EndpointDescriptor{bulkOutEndpoint2, bulkInEndpoint2}

	config := DeviceConfiguration{
		BLength:             9,
		BDescriptorType:     2,
		WTotalLength:        0x0039,
		BNumInterfaces:      2,
		BConfigurationValue: 1,
		IConfiguration:      0,
		BmAttributes:        0xC0, // Self-powered
		BMaxPower:           0x32,
	}

	config.Interfaces = [][]InterfaceDescriptor{
		{ippInterface},
		{ippInterface2},
	}

	return []DeviceConfiguration{config}
}

// GetDeviceDescriptor returns the USB device descriptor for the virtual printer.
func (d *IPPOverUSBDevice) GetDeviceDescriptor() DeviceDescriptor {

	return d.deviceDescriptor
}

// GetConfigurations returns all USB configurations supported by the device.
func (d *IPPOverUSBDevice) GetConfigurations() []DeviceConfiguration {

	return d.configurations
}

func (d *IPPOverUSBDevice) connectToServer() bool {
	parsedURL, err := url.Parse(d.serverURL)
	if err != nil {
		fmt.Printf("Failed to parse server URL: %v\n", err)
		return false
	}

	host := parsedURL.Hostname()
	if host == "" {
		host = "localhost"
	}

	port := parsedURL.Port()
	if port == "" {
		port = "631"
	}

	fmt.Printf("Connecting to IPP server at %s:%s\n", host, port)

	d.connectionLock.Lock()
	defer d.connectionLock.Unlock()

	if d.tcpConnection != nil {
		d.tcpConnection.Close()
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 10*time.Second)
	if err != nil {
		fmt.Printf("Failed to connect to IPP server: %v\n", err)
		d.tcpConnected = false
		d.tcpConnection = nil
		return false
	}

	d.tcpConnection = conn
	d.tcpConnected = true

	fmt.Println("Successfully connected to IPP server")
	return true
}

func (d *IPPOverUSBDevice) disconnectFromServer() {
	d.connectionLock.Lock()
	defer d.connectionLock.Unlock()

	if d.tcpConnection != nil {
		d.tcpConnection.Close()
		d.tcpConnection = nil
	}
	d.tcpConnected = false
}

// HandleData processes bulk IN/OUT USB data requests sent to the device.
func (d *IPPOverUSBDevice) HandleData(usbReq USBRequest) {

	switch usbReq.Ep {
	case 0x01, 0x03:
		d.handleBulkOut(usbReq)
	case 0x82, 0x84:
		d.handleBulkIn(usbReq)
	default:
		fmt.Printf("Unknown endpoint: %02x\n", usbReq.Ep)
		d.SendUSBRet(usbReq, []byte{}, 0, 1)
	}
}

func (d *IPPOverUSBDevice) handleBulkOut(usbReq USBRequest) {
	if len(usbReq.TransferBuffer) == 0 {
		d.SendUSBRet(usbReq, []byte{}, 0, 0)
		return
	}

	fmt.Printf("Received %d bytes from host\n", len(usbReq.TransferBuffer))

	if !d.tcpConnected {
		if !d.connectToServer() {
			d.SendUSBRet(usbReq, []byte{}, 0, 1)
			return
		}
	}

	d.connectionLock.Lock()
	defer d.connectionLock.Unlock()

	if d.tcpConnection == nil || !d.tcpConnected {
		d.SendUSBRet(usbReq, []byte{}, 0, 1)
		return
	}

	_, err := d.tcpConnection.Write(usbReq.TransferBuffer)
	if err != nil {
		fmt.Printf("Error forwarding to IPP server: %v\n", err)
		d.tcpConnection.Close()
		d.tcpConnection = nil
		d.tcpConnected = false
		d.SendUSBRet(usbReq, []byte{}, 0, 1)
		return
	}

	fmt.Printf("Forwarded %d bytes to IPP server\n", len(usbReq.TransferBuffer))

	// Trying to read immediate response
	d.tcpConnection.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	buffer := make([]byte, 8192)
	n, err := d.tcpConnection.Read(buffer)
	if err == nil && n > 0 {
		fmt.Printf("Received immediate response: %d bytes\n", n)
		d.pendingResponse = append(d.pendingResponse, buffer[:n]...)
	}
	d.tcpConnection.SetReadDeadline(time.Now().Add(10 * time.Second))

	d.SendUSBRet(usbReq, []byte{}, len(usbReq.TransferBuffer), 0)
}

func (d *IPPOverUSBDevice) handleBulkIn(usbReq USBRequest) {
	// Checking if we have pending response data
	if len(d.pendingResponse) > 0 {
		dataToSend := d.pendingResponse
		if len(dataToSend) > int(usbReq.TransferBufferLength) {
			dataToSend = dataToSend[:usbReq.TransferBufferLength]
			d.pendingResponse = d.pendingResponse[len(dataToSend):]
		} else {
			d.pendingResponse = nil
		}

		fmt.Printf("Sending %d bytes to host from buffer\n", len(dataToSend))
		d.SendUSBRet(usbReq, dataToSend, len(dataToSend), 0)
		return
	}

	if !d.tcpConnected {
		d.SendUSBRet(usbReq, []byte{}, 0, 0)
		return
	}

	d.connectionLock.Lock()
	defer d.connectionLock.Unlock()

	if d.tcpConnection == nil || !d.tcpConnected {
		d.SendUSBRet(usbReq, []byte{}, 0, 0)
		return
	}

	d.tcpConnection.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	buffer := make([]byte, usbReq.TransferBufferLength)
	n, err := d.tcpConnection.Read(buffer)
	if err == nil && n > 0 {
		fmt.Printf("Received %d bytes from IPP server\n", n)
		d.SendUSBRet(usbReq, buffer[:n], n, 0)
	} else {
		d.SendUSBRet(usbReq, []byte{}, 0, 0)
	}
	d.tcpConnection.SetReadDeadline(time.Now().Add(10 * time.Second))
}

// HandleDeviceSpecificControl processes USB class or vendor-specific control requests.
func (d *IPPOverUSBDevice) HandleDeviceSpecificControl(controlReq StandardDeviceRequest, usbReq USBRequest) {

	if controlReq.BmRequestType == 0xA1 {
		if controlReq.BRequest == 0x01 { // GET_DEVICE_ID
			deviceID := fmt.Sprintf("MFG:%s;CMD:PostScript,PDF;MDL:%s;CLS:PRINTER;",
				d.config.Manufacturer, d.config.Product)

			deviceIDBytes := []byte(deviceID)
			lengthBytes := []byte{byte(len(deviceIDBytes) >> 8), byte(len(deviceIDBytes) & 0xFF)}
			response := append(lengthBytes, deviceIDBytes...)

			d.SendUSBRet(usbReq, response, len(response), 0)
			return
		} else if controlReq.BRequest == 0x02 { // GET_PORT_STATUS
			status := byte(0x18)
			d.SendUSBRet(usbReq, []byte{status}, 1, 0)
			return
		}
	} else if controlReq.BmRequestType == 0x21 {
		if controlReq.BRequest == 0x02 { // SOFT_RESET
			fmt.Println("Printer soft reset requested")
			d.SendUSBRet(usbReq, []byte{}, 0, 0)
			return
		}
	}

	fmt.Printf("Unhandled control request: %02x %02x\n", controlReq.BmRequestType, controlReq.BRequest)
	d.SendUSBRet(usbReq, []byte{}, 0, 1)
}

func main() {
	configFile := "ipp_usb_config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	ippDevice, err := NewIPPOverUSBDevice(configFile)
	if err != nil {
		fmt.Printf("Error creating device: %v\n", err)
		return
	}

	defer ippDevice.disconnectFromServer()

	container := &USBContainer{}
	container.AddUSBDevice(ippDevice)

	// Get listen settings from config
	listenIP := ippDevice.config.ListenIP
	listenPort := ippDevice.config.ListenPort

	if listenIP == "" {
		listenIP = "0.0.0.0"
	}
	if listenPort == 0 {
		listenPort = 3240
	}

	fmt.Printf("Listening on %s:%d\n", listenIP, listenPort)
	fmt.Println("Press Ctrl+C to stop")

	container.Run(listenIP, listenPort)
}
