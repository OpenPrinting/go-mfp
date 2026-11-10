# OpenPrinting LPD Printer
#
# This model describes a printer, that supports only LPD protocol

# DNS-SD parameters
dnssd.device = dnssd.Device(
    instance = 'OpenPrinting LPD Printer',
    UUID = UUID('cbca9fde-d208-4156-9a2f-82857dd0be5f'),
    services = [
        dnssd.Service(
            types = ['_printer._tcp'],
            TXT = [
                'txtvers=1',
                'UUID=cbca9fde-d208-4156-9a2f-82857dd0be5f',
                'note='
                'pdl=application/postscript,application/pdf',
                'ty=OpenPrinting LPD Printer',
                'product=(OpenPrinting LPD Printer)',
                'qtotal=1',
                'priority=50',
                'usb_MFG=OpenPrinting',
                'usb_MDL=LPD Printer',
            ],
            endpoints = [
                'lpd://192.168.1.100',
                'lpd://[fe80::ce64:1aff:fe71:e2d0]',
            ],
        ),
    ],
)
