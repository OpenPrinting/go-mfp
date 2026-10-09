# mfp-backend-dnssd

`mfp-backend-dnssd` is a drop-in replacement for the
`/usr/lib/cups/backend/dnssd` discovery backend.

It supports DNS-SD discovery of printers implementing one of the
following protocols:

* IPP (`_ipp._tcp` or `_ipps._tcp`)
* LPD (`_printer._tcp`)
* AppSocket (`_pdl-datastream._tcp`)

Compared to the native `dnssd` backend, it provides the following
improvements:

* Services announced on the localhost are not ignored, allowing
  locally running printer applications and simulated hardware to
  be discovered.
* All supported protocols are returned for each device, allowing
  printer configuration utilities to choose between them. The
  standard `dnssd` backend prioritizes the available protocols and
  returns only a single endpoint for each printer.

<!-- vim:ts=8:sw=4:et:textwidth=72
-->
