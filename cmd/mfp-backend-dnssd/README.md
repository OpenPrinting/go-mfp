# mfp-backend-dnssd

`mfp-backend-dnssd` is an alternative DNS-SD discovery backend for CUPS, designed to discover locally advertised printers, expose all available connection methods, and improve consistency across printer discovery backends.

## Motivation

CUPS uses several independent backends to discover network printers, notably `driverless`, `dnssd`, and `snmp`. They use different discovery mechanisms and selection policies, which can leave some printers undiscovered or produce inconsistent lists of available connection methods.

One limitation of the native `dnssd` backend is that it deliberately ignores services announced on localhost. This excludes simulated printers used for testing, as well as printer applications that expose virtual or application-managed printers through DNS-SD. Such services can be useful and perfectly legitimate discovery targets, but the native backend does not report them.

There is also an inconsistency in how CUPS discovery backends handle devices offering multiple connection methods:

* **`driverless`** discovers IPP printers that meet the requirements of the driverless printing ecosystem, including Bonjour, Mopria, and IPP Everywhere. It does not aim to discover every device that happens to offer IPP.
* **`dnssd`** discovers services advertised through DNS-SD, but selects a single connection method when a device offers several protocols.
* **`snmp`** discovers printers through SNMP and reports the connection methods it finds.

For example, a printer advertising both IPP and LPD through DNS-SD may be reported by `dnssd` only as an IPP device, with its LPD endpoint discarded. If the same printer is also discovered through SNMP, the LPD connection may nevertheless appear in the list of available devices.

Consequently, the connection methods presented to printer configuration software may depend on which backend discovers the device, rather than solely on the methods the device actually offers.

## Approach

`mfp-backend-dnssd` takes a different approach to DNS-SD discovery:

* **Discover locally advertised services.** Unlike the native `dnssd` backend, it does not ignore services announced on localhost. This allows local printer applications and simulated printers to be discovered.
* **Report all discovered connection methods.** If a device advertises multiple supported protocols, report each of them instead of selecting a single preferred endpoint. Choosing the appropriate connection method is left to the printer configuration software.
* **Improve consistency across backends.** By retaining connection methods that the native `dnssd` backend would discard, `mfp-backend-dnssd` can reduce inconsistencies between DNS-SD and SNMP discovery. When both backends discover the same printer and report equivalent connection methods, CUPS can deduplicate the resulting entries rather than presenting different connection options depending on which backend found the device.
* **Support multiple printing protocols.** Discover IPP, IPP over TLS, LPD, and AppSocket services.

## Supported DNS-SD service types

| Protocol     | DNS-SD service type    |
|--------------|------------------------|
| IPP          | `_ipp._tcp`            |
| IPP over TLS | `_ipps._tcp`           |
| LPD          | `_printer._tcp`        |
| AppSocket    | `_pdl-datastream._tcp` |

`mfp-backend-dnssd` is an alternative to the native CUPS `driverless`
and `dnssd` backends, providing a unified approach to DNS-SD printer
discovery, including locally advertised services and all available
connection methods.

<!-- vim:ts=8:sw=4:et:textwidth=72
-->
