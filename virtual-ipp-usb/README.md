# IPP-over-USB Device emulator

Implementation of a **virtual USB printer device** using the **USB/IP protocol** that proxies print jobs to a remote IPP server (like CUPS). It simulates a real IPP-over-USB printer on the USB bus.

## Installation

```bash
git clone https://github.com/OpenPrinting/go-mfp.git
cd virtual-ipp-usb
go mod init virtual-ipp-usb
go mod tidy
```

---

## Configuration

When run for the first time, the emulator will auto-generate a default config file:
File: `ipp_usb_config.json`
This file can be modified to suit the environment.

---

## Running

Build and run:

```bash
go build -o ipp-printer ipp_printer.go USBIP.go
sudo ./ipp-printer
```

## Listing and Attaching (with usbip)

To use the device with a USB/IP client:

```bash
sudo usbip list -r 127.0.0.1
sudo usbip attach -r 127.0.0.1 -b 1-1
```

💡 *Note: Ensure usbip kernel modules are loaded:*

```bash
sudo modprobe usbip_core
sudo modprobe vhci_hcd
```
