// MFP                    - Miulti-Function Printers and scanners toolkit
// cmd/mfp-backend-dnssd  - /usr/lib/cups/backend/dnssd replacement
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// The main() function.

package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/OpenPrinting/go-avahi"
)

// main function for the mfp-backend-dnssd command
func main() {
	// If called with arguments, print usage and exit
	if len(os.Args) > 1 {
		usage()
	}

	// Perform the discovery
	services, err := avahi.SimpleServiceResolver(
		context.TODO(),
		avahi.IfIndexUnspec,
		avahi.ProtocolUnspec,
		[]string{
			"_ipp._tcp",
			"_ipps._tcp",
			"_printer._tcp",
			"_pdl-datastream._tcp",
		},
		"",
		avahi.LookupNoAddress,
		avahi.DefaultTimeoutInteractive,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	// Group services by instance name
	devices := servicesClassify(services)

	// Generate output
	destList := []string{}
	destSeen := make(map[string]struct{})

	for _, dev := range devices {
		for _, svc := range dev.Services {
			uri := svcURI(svc)
			if uri == "" {
				continue
			}

			mfg := dev.Mfg()
			mdl := dev.Mdl()

			dest := reportDestination(
				uri,
				svcMakeModel(svc),
				svcInfo(svc),
				svcDeviceID(svc, mfg, mdl),
				svcLocation(svc),
			)

			if _, seen := destSeen[dest]; !seen {
				destSeen[dest] = struct{}{}
				destList = append(destList, dest)
			}
		}
	}

	// Write to Stdout
	for _, dest := range destList {
		fmt.Println(dest)
	}
}

// servicesContainsType reports if there is a service of the
// specified type within a slice of services.
func servicesContainsType(services []*avahi.Service, t string) bool {
	for _, svc := range services {
		if svc.SvcType == t {
			return true
		}
	}

	return false
}

// servicesDropInactive removes services of the particular type from
// the list.
//
// The list is updated in place.
func servicesDropType(services []*avahi.Service, t string) []*avahi.Service {
	cnt := 0
	for _, svc := range services {
		if svc.SvcType != t {
			services[cnt] = svc
			cnt++
		}
	}

	return services[:cnt]
}

// servicesClassify group services by instance name
// and returns a slice of devices.
func servicesClassify(services []*avahi.Service) []device {
	byname := make(map[string][]*avahi.Service)

	// Classify services by instance names.
	// Skip inactive services.
	for _, svc := range services {
		name := svc.InstanceName
		if len(svc.Endpoints) > 0 {
			byname[name] = append(byname[name], svc)
		}
	}

	// Convert map into the slice of devices
	//
	// Inactive services are removed and devices without
	// active services are ignored.
	devices := make([]device, 0, len(byname))
	for name, services := range byname {
		dev := device{name, services}
		if servicesContainsType(dev.Services, "_ipps._tcp") {
			dev.Services = servicesDropType(dev.Services,
				"_ipp._tcp")
		}

		devices = append(devices, dev)
	}

	// Sort devices by name
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Name < devices[j].Name
	})

	// Sort services of each device by rank
	for i := range devices {
		dev := &devices[i]
		sort.Slice(dev.Services, func(i, j int) bool {
			svc1 := dev.Services[i]
			svc2 := dev.Services[j]
			return svcRank(svc1) < svcRank(svc2)
		})
	}

	return devices
}

// placing the most preferred services first.
//
// It returns an integer rank—the lower the number, the
// more preferred the service is.
func svcRank(svc *avahi.Service) int {
	switch svc.SvcType {
	case "_ipps._tcp":
		return 0
	case "_ipp._tcp":
		return 1
	case "_pdl-datastream._tcp":
		return 2
	case "_printer._tcp":
		return 3
	}

	return 1000
}

// svcURI returns an device-uri string, pointing to the service.
// If URI cannot be generated for some reason, it returns an empty string.
func svcURI(svc *avahi.Service) string {
	// Skip inactive services (i.e., services without endpoints)
	if len(svc.Endpoints) == 0 {
		return ""
	}

	scheme := ""
	switch svc.SvcType {
	case "_ipp._tcp":
		scheme = "ipp"
	case "_ipps._tcp":
		scheme = "ipps"
	case "_printer._tcp":
		scheme = "lpd"
	case "_pdl-datastream._tcp":
		scheme = "socket"
	}

	if scheme == "" {
		return ""
	}

	name, err := avahi.DomainServiceNameJoin(svc.InstanceName,
		svc.SvcType, svc.Domain)
	if err != nil {
		return ""
	}

	name = hostNameEncode(avahi.DomainServiceNameUnescape(name))

	return scheme + "://" + name + "/"
}

// svcMakeModel returns the service's device-make-and-model string.
func svcMakeModel(svc *avahi.Service) string {
	// Try "ty" record
	if ty := svcTXT(svc, "ty"); ty != "" {
		return ty
	}

	// Try "product" record
	if product := svcTXT(svc, "product"); product != "" {
		if len(product) > 0 && product[0] == '(' {
			product = product[1:]
		}

		if len(product) > 0 && product[len(product)-1] == ')' {
			product = product[:len(product)-1]
		}

		if product != "" {
			return product
		}
	}

	// Try usb_MFG + usb_MDL
	mfg := svcTXT(svc, "usb_MFG")
	mdl := svcTXT(svc, "usb_MDL")
	switch {
	case mfg != "" && mdl != "":
		if strings.HasPrefix(mdl, mfg) {
			return mdl
		}
		return mfg + " " + mdl

	case mfg == "" && mdl != "":
		return mdl
	case mfg != "" && mdl == "":
		return mfg
	}

	// Give up
	return ""
}

// svcInfo returns the service's device-info string.
func svcInfo(svc *avahi.Service) string {
	return svc.InstanceName
}

// svcDeviceID returns the service's IEEE 1284 device ID.
//
// The mfg and mdl parameters provide fallback values for services
// that do not advertise their own manufacturer and model.
func svcDeviceID(svc *avahi.Service, mfg, mdl string) string {
	attrs := []string{}

	if s := svcTXT(svc, "usb_MFG"); s != "" {
		mfg = s
	}

	if s := svcTXT(svc, "usb_MDL"); s != "" {
		mdl = s
	}

	if mfg := sanitizeDeviceIDValue(mfg); mfg != "" {
		attrs = append(attrs, "MFG:"+mfg+";")
	}

	if mdl := sanitizeDeviceIDValue(mdl); mdl != "" {
		attrs = append(attrs, "MDL:"+mdl+";")
	}

	if cmd := svcCMD(svc); cmd != "" {
		attrs = append(attrs, "CMD:"+cmd+";")
	}

	return strings.Join(attrs, "")
}

// svcCMD returns the CMD value of the IEEE-1284 device ID
func svcCMD(svc *avahi.Service) string {
	// Try usb_CMD
	cmd := svcTXT(svc, "usb_CMD")
	if cmd != "" {
		return cmd
	}

	// Try to generate cmd from PDL
	pdl := svcTXT(svc, "pdl")
	if pdl == "" {
		return ""
	}

	formats := []string{}
	formatsSeen := make(map[string]struct{})

	for _, mime := range strings.Split(pdl, ",") {
		var add []string
		switch strings.ToLower(strings.TrimSpace(mime)) {
		case "application/pdf":
			add = []string{"PDF"}
		case "application/postscript":
			add = []string{"PS"}

		case "application/vnd.epson.esc_p":
			add = []string{"ESCP"}
		case "application/vnd.epson.esc_p-r":
			add = []string{"ESCPR"}

		case "application/vnd.hp-pcl":
			add = []string{"PCL"}
		case "application/vnd.hp-pclxl":
			add = []string{"PCLXL"}
		case "application/vnd.hp-xqx":
			add = []string{"XQX"}

		case "application/vnd.canon-bj":
			add = []string{"BJ"}
		case "application/vnd.canon-capt":
			add = []string{"CAPT"}
		case "application/vnd.canon-cpdl":
			add = []string{"CPDL"}
		case "application/vnd.canon-lips":
			add = []string{"LIPS"}

		case "application/vnd.ms-xpsdocument":
			add = []string{"XPS"}
		case "application/oxps":
			add = []string{"XPS"}
		case "application/pclm":
			add = []string{"PCLM"}

		case "image/jpeg":
			add = []string{"JPEG"}
		case "image/pwg-raster":
			add = []string{"PWGRaster", "PWG"}
		case "image/urf":
			add = []string{"AppleRaster", "URF"}
		case "image/tiff":
			add = []string{"TIFF"}
		case "image/png":
			add = []string{"PNG"}

		case "text/plain":
			add = []string{"TEXT"}
		}

		for _, fmt := range add {
			if _, seen := formatsSeen[fmt]; !seen {
				formatsSeen[fmt] = struct{}{}
				formats = append(formats, fmt)
			}
		}
	}

	if len(formats) > 0 {
		return strings.Join(formats, ",")
	}

	// Give up
	return ""
}

// svcLocation returns the service's device-location string.
func svcLocation(svc *avahi.Service) string {
	return svcTXT(svc, "note")
}

// svcTXT returns the service's TXT record's value by key name.
func svcTXT(svc *avahi.Service, key string) string {
	key = strings.ToLower(key) + "="
	for _, txt := range svc.Txt {
		if strings.HasPrefix(strings.ToLower(txt), key) {
			return txt[len(key):]
		}
	}

	return ""
}

// device represents services, grouped by DNS-SD Instance name.
type device struct {
	Name     string           // Instance name
	Services []*avahi.Service // Available services
}

// Mfg returns device's manufacturer, trying to look to all
// available services.
func (dev device) Mfg() string {
	for _, svc := range dev.Services {
		mfg := svcTXT(svc, "usb_MFG")
		if mfg != "" {
			return mfg
		}
	}

	return ""
}

// Mdl returns device's model, trying to look to all
// available services.
func (dev device) Mdl() string {
	for _, svc := range dev.Services {
		mdl := svcTXT(svc, "usb_MDL")
		if mdl != "" {
			return mdl
		}
	}

	return ""
}

// reportDestination formats a destination report string.
func reportDestination(uri, model, info, devid, location string) string {
	buf := strings.Builder{}

	// Write device-class and device-uri
	buf.WriteString("network ")
	buf.WriteString(uri)

	// Write device-make-and-model
	buf.WriteByte(' ')
	if model != "" {
		buf.WriteString(quote(model))
	} else {
		buf.WriteString(quote("unknown"))
	}

	// Write device-info
	buf.WriteByte(' ')
	buf.WriteString(quote(info))

	// Write device-id, if available or not last
	if devid != "" || location != "" {
		buf.WriteByte(' ')
		buf.WriteString(quote(devid))
	}

	// Write device-location, if available
	if location != "" {
		buf.WriteByte(' ')
		buf.WriteString(quote(location))
	}

	return buf.String()
}

// quote returns the quoted string.
// special characters are \-escaped.
func quote(s string) string {
	buf := strings.Builder{}
	buf.WriteByte('"')

	for _, c := range []byte(s) {
		switch {
		case c == '\\' || c == '"':
			buf.WriteByte('\\')
			buf.WriteByte(c)
		case (c < ' ' && c != '\t') || c == 0x7f:
			buf.WriteByte(' ')
		default:
			buf.WriteByte(c)
		}
	}

	buf.WriteByte('"')
	return buf.String()
}

// hostNameEncode encodes a hostname string by replacing reserved characters
// with their percent-encoded equivalents.
func hostNameEncode(s string) string {
	buf := strings.Builder{}

	for _, c := range []byte(s) {
		escape := false
		switch c {
		case '%', ' ', '/', '?', '#', '[', ']', '@':
			escape = true
		default:
			escape = c >= 128
		}

		if escape {
			buf.WriteString(fmt.Sprintf("%%%2.2X", c))
		} else {
			buf.WriteByte(c)
		}
	}

	return buf.String()
}

// sanitizeDeviceIDValue sanitizes a string value for use in
// an IEEE 1284 Device ID.
func sanitizeDeviceIDValue(s string) string {
	buf := strings.Builder{}

	for _, c := range []byte(s) {
		switch {
		case c == ':' || c == ',' || c == ';':
			// Delimiter characters are not allowed in values.

		case c < ' ' || c == 0x7f:
			// Control characters are not allowed.

		default:
			buf.WriteByte(c)
		}
	}

	return buf.String()
}

// usage prints usage information and exits.
func usage() {
	const text = `usage: mfp-backend-dnssd

mfp-backend-dnssd is a drop-in replacement for the
/usr/lib/cups/backend/dnssd discovery backend.

It supports DNS-SD discovery of printers implementing one of the
following protocols:

    * IPP (_ipp._tcp or _ipps._tcp)
    * LPD (_printer._tcp)
    * AppSocket (_pdl-datastream._tcp)

Compared to the native dnssd backend, it provides the following
improvements:

    * Services announced on the localhost are not ignored, allowing
      locally running printer applications and simulated hardware to
      be discovered.
    * All supported protocols are returned for each device, allowing
      printer configuration utilities to choose between them. The
      standard dnssd backend prioritizes the available protocols and
      returns only a single endpoint for each printer.
`

	os.Stdout.Write([]byte(text))
	os.Exit(0)
}
