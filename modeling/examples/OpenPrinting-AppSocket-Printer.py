# OpenPrinting AppSocket Printer
#
# This model describes a printer, that supports only AppSocket protocol

# DNS-SD parameters
dnssd.device = dnssd.Device(
    instance = 'OpenPrinting AppSocket Printer',
    UUID = UUID('5d25553b-68b4-4383-a284-3cff1192972f'),
    services = [
        dnssd.Service(
            types = ['_pdl-datastream._tcp'],
            TXT = [
                'txtvers=1',
                'UUID=5d25553b-68b4-4383-a284-3cff1192972f',
                'note=',
                'pdl=application/postscript,application/pdf',
                'ty=OpenPrinting AppSocket Printer',
                'product=(OpenPrinting AppSocket Printer)',
                'qtotal=1',
                'priority=50',
            ],
            endpoints = [
                'socket://192.168.1.100',
                'socket://[fe80::ce64:1aff:fe71:e2d0]',
            ],
        ),
    ],
)
