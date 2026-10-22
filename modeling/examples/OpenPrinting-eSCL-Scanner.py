# OpenPrinting eSCL Scanner
#
# This model describes a pure eSCL Scanner (without printer)
# with DNS-SD and IPP over USB support.

# eSCL scanner parameters:
escl.scanner = escl.ScannerCapabilities(
    Version = '2.9',
    MakeAndModel = 'OpenPrinting eSCL Scanner',
    SerialNumber = 'SN-314159265359',
    Uuid = UUID('4a6e56ce-b372-4fe0-bdd3-7d939b67bc03'),
    Platen = escl.Platen(
        PlatenInputCaps = escl.InputSourceCaps(
            MinWidth = 300,
            MaxWidth = 2550,
            MinHeight = 300,
            MaxHeight = 3508,
            MaxOpticalXResolution = 600,
            MaxOpticalYResolution = 600,
            MaxScanRegions = 1,
            MaxPhysicalWidth = 2550,
            MaxPhysicalHeight = 3508,
            SupportedIntents = [
                escl.Document,
                escl.Photo,
                escl.Preview,
                escl.TextAndGraphic,
            ],
            SettingProfiles = [
                escl.SettingProfile(
                    ColorModes = [
                        escl.RGB24,
                        escl.Grayscale8,
                        escl.BlackAndWhite1,
                    ],
                    ContentTypes = [
                        escl.Photo,
                        escl.Text,
                        escl.TextAndPhoto,
                    ],
                    DocumentFormats = [
                        'application/pdf',
                        'image/jpeg',
                        'image/tiff',
                    ],
                    DocumentFormatsExt = [
                        'application/pdf',
                        'image/jpeg',
                        'image/tiff',
                    ],
                    SupportedResolutions = [
                        escl.SupportedResolutions(
                            ColorMode = escl.RGB24,
                            DiscreteResolutions = [
                                escl.DiscreteResolution(
                                    XResolution = 75,
                                    YResolution = 75,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 150,
                                    YResolution = 150,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 200,
                                    YResolution = 200,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 300,
                                    YResolution = 300,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 400,
                                    YResolution = 400,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 600,
                                    YResolution = 600,
                                ),
                            ],
                        ),
                        escl.SupportedResolutions(
                            ColorMode = escl.Grayscale8,
                            DiscreteResolutions = [
                                escl.DiscreteResolution(
                                    XResolution = 75,
                                    YResolution = 75,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 150,
                                    YResolution = 150,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 200,
                                    YResolution = 200,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 300,
                                    YResolution = 300,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 400,
                                    YResolution = 400,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 600,
                                    YResolution = 600,
                                ),
                            ],
                        ),
                        escl.SupportedResolutions(
                            ColorMode = escl.BlackAndWhite1,
                            DiscreteResolutions = [
                                escl.DiscreteResolution(
                                    XResolution = 75,
                                    YResolution = 75,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 150,
                                    YResolution = 150,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 200,
                                    YResolution = 200,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 300,
                                    YResolution = 300,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 400,
                                    YResolution = 400,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 600,
                                    YResolution = 600,
                                ),
                            ],
                        ),
                    ],
                    ColorSpaces = [
                        'RGB',
                        escl.sRGB,
                    ],
                ),
            ],
            FeedDirections = [
                escl.LongEdgeFeed,
                escl.ShortEdgeFeed,
            ],
        ),
    ),
    Adf = escl.Adf(
        AdfSimplexInputCaps = escl.InputSourceCaps(
            MinWidth = 1230,
            MaxWidth = 2550,
            MinHeight = 1740,
            MaxHeight = 4200,
            MaxOpticalXResolution = 600,
            MaxOpticalYResolution = 600,
            MaxScanRegions = 1,
            SupportedIntents = [
                escl.Document,
                escl.Photo,
                escl.Preview,
                escl.TextAndGraphic,
            ],
            SettingProfiles = [
                escl.SettingProfile(
                    ColorModes = [
                        escl.RGB24,
                        escl.Grayscale8,
                        escl.BlackAndWhite1,
                    ],
                    DocumentFormats = [
                        'application/pdf',
                        'image/jpeg',
                        'image/tiff',
                    ],
                    DocumentFormatsExt = [
                        'application/pdf',
                        'image/jpeg',
                        'image/tiff',
                    ],
                    SupportedResolutions = [
                        escl.SupportedResolutions(
                            ColorMode = escl.RGB24,
                            DiscreteResolutions = [
                                escl.DiscreteResolution(
                                    XResolution = 75,
                                    YResolution = 75,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 150,
                                    YResolution = 150,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 200,
                                    YResolution = 200,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 300,
                                    YResolution = 300,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 400,
                                    YResolution = 400,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 600,
                                    YResolution = 600,
                                ),
                            ],
                        ),
                        escl.SupportedResolutions(
                            ColorMode = escl.Grayscale8,
                            DiscreteResolutions = [
                                escl.DiscreteResolution(
                                    XResolution = 75,
                                    YResolution = 75,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 150,
                                    YResolution = 150,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 200,
                                    YResolution = 200,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 300,
                                    YResolution = 300,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 400,
                                    YResolution = 400,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 600,
                                    YResolution = 600,
                                ),
                            ],
                        ),
                        escl.SupportedResolutions(
                            ColorMode = escl.BlackAndWhite1,
                            DiscreteResolutions = [
                                escl.DiscreteResolution(
                                    XResolution = 75,
                                    YResolution = 75,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 150,
                                    YResolution = 150,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 200,
                                    YResolution = 200,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 300,
                                    YResolution = 300,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 400,
                                    YResolution = 400,
                                ),
                                escl.DiscreteResolution(
                                    XResolution = 600,
                                    YResolution = 600,
                                ),
                            ],
                        ),
                    ],
                    ColorSpaces = [
                        'RGB',
                        escl.sRGB,
                    ],
                ),
            ],
            FeedDirections = [
                escl.LongEdgeFeed,
                escl.ShortEdgeFeed,
            ],
        ),
        AdfOptions = [
            escl.DetectPaperLoaded,
        ],
        Justification = escl.Justification(
            XImagePosition = escl.Left,
            YImagePosition = escl.Top,
        ),
    ),
    BrightnessSupport = escl.Range(
        Min = 1,
        Max = 9,
        Normal = 5,
        Step = 1,
    ),
    CompressionFactorSupport = escl.Range(
        Min = 5,
        Max = 95,
        Normal = 25,
        Step = 1,
    ),
    ContrastSupport = escl.Range(
        Min = 0,
        Max = 5,
        Normal = 0,
        Step = 1,
    ),
    BlankPageDetectionAndRemoval = False,
)


# ----- PUT YOUR ESCL HOOKS HERE -----

# Called on request:  POST /{root}/ScanJobs
#
# def escl_onScanJobsRequest (q: query.Query, rq: escl.ScanSettings):

# Called on response: GET /{JobUri}/NextDocument
#
# def escl_onNextDocumentResponse (q: query.Query, flt: escl.ImageFilter):

# DNS-SD parameters
dnssd.device = dnssd.Device(
    instance = 'OpenPrinting eSCL Scanner',
    UUID = UUID('4a6e56ce-b372-4fe0-bdd3-7d939b67bc03'),
    services = [
        dnssd.Service(
            types = [
                '_uscan._tcp',
                '_uscans._tcp',
            ],
            TXT = [
                'txtvers=1',
                'UUID=4a6e56ce-b372-4fe0-bdd3-7d939b67bc03',
                'cs=color,grayscale',
                'duplex=F',
                'is=platen,adf',
                'mopria-certified-scan=1.4',
                'note=',
                'pdl=application/pdf,image/jpeg,image/tiff',
                'rs=eSCL',
                'ty=OpenPrinting eSCL Scanner',
                'usb_MDL=eSCL Scanner',
                'usb_MFG=OpenPrinting',
                'vers=2.9',
            ],
            endpoints = [
                'http://192.168.1.100/eSCL',
                'https://192.168.1.100/eSCL',
                'http://[fe80::ce64:1aff:fe71:e2d0]/eSCL/',
                'https://[fe80::ce64:1aff:fe71:e2d0]/eSCL/',
            ],
        ),
    ],
)

# USB device descriptor
usb.device = usb.DeviceDescriptor(
    bcdUSB = '2.0',
    Speed = 3,
    bDeviceClass = 0,
    bDeviceSubClass = 0,
    bDeviceProtocol = 0,
    bMaxPacketSize = 64,
    idVendor = 0x1d6b,
    idProduct = 0xeeff,
    bcdDevice = '0.0',
    iManufacturer = 'OpenPrinting',
    iProduct = 'eSCL Scanner',
    iSerialNumber = 'SN-314159265359',
    Configurations = [
        usb.ConfigurationDescriptor(
            bConfigurationValue = 1,
            iConfiguration = '',
            bmAttributes = 192,
            bMaxPower = 1,
            Interfaces = [
                usb.Interface(
                    bInterfaceNumber = 0,
                    AltSettings = [
                        usb.InterfaceDescriptor(
                            bInterfaceClass = 7,
                            bInterfaceSubClass = 1,
                            bInterfaceProtocol = 4,
                            bAlternateSetting = 0,
                            iInterface = 'eSCL scanner',
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
                usb.Interface(
                    bInterfaceNumber = 1,
                    AltSettings = [
                        usb.InterfaceDescriptor(
                            bInterfaceClass = 7,
                            bInterfaceSubClass = 1,
                            bInterfaceProtocol = 4,
                            bAlternateSetting = 0,
                            iInterface = 'eSCL scanner',
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
                usb.Interface(
                    bInterfaceNumber = 2,
                    AltSettings = [
                        usb.InterfaceDescriptor(
                            bInterfaceClass = 7,
                            bInterfaceSubClass = 1,
                            bInterfaceProtocol = 4,
                            bAlternateSetting = 0,
                            iInterface = 'eSCL scanner',
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

