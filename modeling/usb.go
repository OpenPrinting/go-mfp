// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// USB part of Model

package modeling

import (
	"fmt"

	"github.com/OpenPrinting/go-mfp/proto/usb"
)

// SetUSBDeviceDescriptor sets the [usb.DeviceDescriptor].
func (model *Model) SetUSBDeviceDescriptor(desc *usb.DeviceDescriptor) {
	model.usbDevice = desc
}

// GetUSBDeviceDescriptor returns the [usb.DeviceDescriptor].
func (model *Model) GetUSBDeviceDescriptor() *usb.DeviceDescriptor {
	return model.usbDevice
}

// usbLoad decodes USB part of model. The model file assumed to be
// preloaded into the Model's Python interpreter (model.py).
func (model *Model) usbLoad() error {
	// Obtain Python object for "usb.device"
	obj := model.py.Eval("usb.device")

	if err := obj.Err(); err != nil {
		err = fmt.Errorf("usb.device: %w", err)
		return err
	}

	if obj.IsNone() {
		return nil
	}

	// Decode the usb.DeviceDescriptor
	var desc usb.DeviceDescriptor
	err := structImport(obj, keywordMapUSB, &desc)
	if err != nil {
		return err
	}

	// Update the Model
	model.usbDevice = &desc
	return nil
}
