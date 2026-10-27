# This is the generated MFP model file.
# You probably need to edit it appropriately before use.

# USB device descriptor
usb.device = usb.DeviceDescriptor(
    bcdUSB = '2.0',
    Speed = 3,
    bDeviceClass = 0,
    bDeviceSubClass = 0,
    bDeviceProtocol = 0,
    bMaxPacketSize = 64,
    idVendor = 1008,
    idProduct = 16663,
    bcdDevice = '1.0',
    iManufacturer = 'Hewlett-Packard',
    iProduct = 'HP LaserJet 1018',
    iSerialNumber = 'KP2G65T',
    Configurations = [
        usb.ConfigurationDescriptor(
            bConfigurationValue = 1,
            iConfiguration = '',
            bmAttributes = 192,
            bMaxPower = 49,
            Interfaces = [
                usb.Interface(
                    bInterfaceNumber = 0,
                    AltSettings = [
                        usb.InterfaceDescriptor(
                            bInterfaceClass = 7,
                            bInterfaceSubClass = 1,
                            bInterfaceProtocol = 2,
                            bAlternateSetting = 0,
                            iInterface = '',
                            IEEE1284DeviceID = 'MFG:Hewlett-Packard;MDL:HP LaserJet 1018;CMD:ACL;CLS:PRINTER;DES:HP LaserJet 1018;',
                            Endpoints = [
                                usb.EndpointDescriptor(
                                    Type = usb.OUT,
                                    bmAttributes = 2,
                                    wMaxPacketSize = 512,
                                ),
                                usb.EndpointDescriptor(
                                    Type = usb.IN,
                                    bmAttributes = 2,
                                    wMaxPacketSize = 512,
                                ),
                            ],
                        ),
                    ],
                ),
            ],
        ),
    ],
)

