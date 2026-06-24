# GO-MFP

[![godoc.org](https://godoc.org/github.com/OpenPrinting/go-mfp?status.svg)](https://godoc.org/github.com/OpenPrinting/go-mfp)
![GitHub](https://img.shields.io/github/license/OpenPrinting/go-mfp)
[![Go Report Card](https://goreportcard.com/badge/github.com/OpenPrinting/go-mfp)](https://goreportcard.com/report/github.com/OpenPrinting/go-mfp)

**Warning:** this project is work in progress and under the intensive
development.

This project contains the foundation of the behavior-accurate simulator for
multi-function printers (MFPs) that supports at least IPP 2.x for
printing, eSCL and WSD for scanning, and DNS-SD and WS-Discovery for
device discovery.

This simulator will consist of a core simulation engine that provides
reference implementations of the aforementioned protocols, along with a
customization engine that allows for the expression of implementation
details specific to individual devices without the need to reimplement
common functionalities repeatedly.

For the purpose of implementation of the MFP core simulation engine it
contains many libraries (Go packages) that are aimed to be a part of the
simulator, but can be useful as a general-purpose libraries as well.

It contains (partially done for now) the fairly complete eSCL and WSD
WS-Scan implementations, DNS-SD and WS-Discovery published and querier,
smart device discovery engine, able to bind various instances of the
same device together, representing them as a sub-units of a single
physical device (IPP/IPPS, AppSocket AKA JetDirect and LPD for printing,
eSCL and WSD for scanning), based on the similar work previously done for
scanners in the [sane-airscan](https://github.com/alexpevzner/sane-airscan)
project and much more.

The ultimate goal of this project, among writing the MFP simulator, is
to implement a comprehensive Go toolkit for printing and scanning on
Linux and UNIX-like systems.

## Building

Building is quite straightforward: just type `make` at the root
directory.

It requires the following development packages to be installed:

  * Golang compiler
  * gcc compiler
  * avahi-libs for the DNS-SD
  * libjpeg, libpng for image conversion
  * libusb
  * libppd
  * Python3 (cpython) development package

It also requires some external Go libraries, but if your build computer
is connected to the network, they will be downloaded and cached
automatically at the first build attrmpt:

  * github.com/OpenPrinting/go-avahi - Go binding for Avahi
  * github.com/OpenPrinting/goipp - Low-level IPP implementation
  * github.com/google/go-cmp - used only for testing
  * github.com/peterh/liner - line editor for the interactive shell

## Source tree organization

### import github.com/OpenPrinting/go-mfp/abstract"

This package contains definitions for the abstract
(protocol-independent) printer and scanner.

### import github.com/OpenPrinting/go-mfp/argv"

This is a general-purpose library for parsing command options.

### import github.com/OpenPrinting/go-mfp/cmd"

This is a root directory for all commands (executable files) build as a
part of this project.

#### import github.com/OpenPrinting/go-mfp/cmd/mfp"

TODO

#### import github.com/OpenPrinting/go-mfp/cmd/mfp-cups"

This is the command-line CUPS client.

#### import github.com/OpenPrinting/go-mfp/cmd/mfp-discover"

This is the command-line device discovery tool.

#### import github.com/OpenPrinting/go-mfp/cmd/mfp-proxy"

This is the MFP masquerading tool. It allows to re-expose existent MFP
device under the different name or different IP address and to sniff its
traffic (which is convenient tool for troubleshooting).

#### import github.com/OpenPrinting/go-mfp/cmd/mfp-shell"

This is interactive shell for all mfp-XXX commands.

#### import github.com/OpenPrinting/go-mfp/cmd/mfp-virtual"

This is the vurtual MFP simulator.

#### import github.com/OpenPrinting/go-mfp/cups"

This is the CUPS API library. It uses the IPP protocol to communicate
with the CUPS server.

#### import github.com/OpenPrinting/go-mfp/discovery"

This is the smart MFP discovery library.

#### import github.com/OpenPrinting/go-mfp/internal"

This is the collection of internally used support libraries. These
libraries are not exposed via the GO-MFP API.

#### import github.com/OpenPrinting/go-mfp/log"

This is the logging library, used across the project.

### import github.com/OpenPrinting/go-mfp/proto"

This is the collection of packages, that implement MFP protocols.

#### import github.com/OpenPrinting/go-mfp/proto/escl"

This is the eSCL protocol implementation.

#### import github.com/OpenPrinting/go-mfp/proto/ipp"

This is the high-level IPP protocol implementation.

#### import github.com/OpenPrinting/go-mfp/proto/wsd"

This is the WS-Discovery protocol implementation, querier and responder.

### import github.com/OpenPrinting/go-mfp/transport"

This is the low-level transport library, used for communication with
devices.

### import github.com/OpenPrinting/go-mfp/util"

This is the collection of internally used support libraries. These
libraries are exposed via the GO-MFP API. 

<!-- vim:ts=8:sw=4:et:textwidth=72
-->
