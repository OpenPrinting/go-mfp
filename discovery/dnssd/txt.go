// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// DNS-SD TXT records

package dnssd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// txtPrinter represents a decoded TXT record for printer
type txtPrinter struct {
	uuid      uuid.UUID                    // Device UUID
	svcType   string                       // Service type
	makeModel string                       // Manufacturer + Model
	location  string                       // E.g., "2nd Floor Computer Lab"
	adminURL  string                       // Device administration URL
	iconURL   string                       // Device icon URL
	usbMFG    string                       // I.e., "Hewlett Packard"
	usbMDL    string                       // Model name
	usbSerial string                       // usb_SER, "" if none
	usbHWID   string                       // usb_HWID, "" if none
	params    *discovery.PrinterParameters // Printer parameters
}

// txtScanner represents a decoded TXT record for scanner
type txtScanner struct {
	uuid      uuid.UUID                    // Device UUID
	svcType   string                       // Service type
	uriPath   string                       // Path part of URI
	makeModel string                       // Manufacturer + Model
	location  string                       // E.g., "2nd Floor Computer Lab"
	adminURL  string                       // Device administration URL
	iconURL   string                       // Device icon URL
	usbSerial string                       // usb_SER, "" if none
	usbHWID   string                       // usb_HWID, "" if none
	params    *discovery.ScannerParameters // Scanner parameters
}

// txtPrinter decodes record for printer
func decodeTxtPrinter(svcType, svcInstance string,
	txt []string) (txtPrinter, error) {

	// Set defaults
	p := txtPrinter{
		// Save service type
		svcType: svcType,

		// Default parameters
		params: &discovery.PrinterParameters{
			// 50 is the reasonable default
			Priority: 50,
		},
	}

	// Some defaults depend on service type
	switch svcType {
	case svcTypeIPP, svcTypeIPPS:
		p.params.Queue = "ipp/print"
	case svcTypeLPD:
		p.params.Queue = "auto"
	}

	// Parse key by key
	found := map[string]struct{}{}

	for _, t := range txt {
		// Split record into key and value.
		// Ignore keys without values -- CUPS does the same.
		key, value, novalue := txtParse(t)

		if novalue {
			continue
		}

		// Update found/missed sets and check for duplicates
		if _, dup := found[key]; dup {
			continue
		}

		// Decode the value
		var err error
		switch txToLower(key) {
		case "adminurl":
			p.adminURL = value
		case "air":
			p.params.Auth, err = txtAuth(value)
		case "bind":
			p.params.Bind, err = txtOption(value)
		case "color":
			p.params.Color, err = txtOption(value)
		case "copies":
			p.params.Copies, err = txtOption(value)
		case "duplex":
			p.params.Duplex, err = txtOption(value)
		case "fax":
		case "kind":
			p.params.Media, err = txtMediaKind(value)
		case "note":
			p.location = value
		case "papermax":
			p.params.Paper, err = txtPaperMax(value)
		case "pdl":
			p.params.PDL, err = txtKeywords(value)
		case "priority":
			p.params.Priority, err = strconv.Atoi(value)
			if err != nil {
				break
			}
			if p.params.Priority < 0 || p.params.Priority >= 100 {
				err = errors.New("out of range (0...99)")
			}
		case "product":
			p.params.PSProduct = value
		case "punch":
			// Punch can take values "0", "1", "2", "3",
			// "4" and "U", according to the number of holes
			// the puncher can make. "0" means "no punching"
			if value != "0" && txToLower(value) != "u" {
				p.params.Punch = optional.New(true)
			}
		case "rp":
			p.params.Queue = value
		case "sort":
			p.params.Sort, err = txtOption(value)
		case "staple":
			p.params.Staple, err = txtOption(value)
		case "txtvers":
			if value != "1" {
				err = fmt.Errorf("unknown version %q", value)
			}
		case "ty":
			p.makeModel = value
		case "usb_mdl":
			p.usbMDL = value
		case "usb_mfg":
			p.usbMFG = value
		case "uuid":
			p.uuid, err = uuid.Parse(value)

		// ipp-usb extensions
		case "usb_ser":
			p.usbSerial = value
		case "usb_hwid":
			p.usbHWID = value

		// These parameters are ignored so far, but it may change
		// in the future.
		//
		// TODO: review these parameters.
		case "print_wfds":
		case "qtotal":
		case "scan":
		case "tls":
		case "urf":
		}

		// Check for error
		if err != nil {
			err = fmt.Errorf("%s: %w", key, err)
			return txtPrinter{}, err
		}
	}

	return p, nil
}

// txtPrinter decodes record for printer
func decodeTxtScanner(svcType, svcInstance string,
	txt []string) (txtScanner, error) {
	s := txtScanner{
		// Save service type
		svcType: svcType,

		// The default path to the eSCL endpoint
		uriPath: "eSCL",

		// Default parameters
		params: &discovery.ScannerParameters{},
	}

	found := map[string]struct{}{}

	for _, t := range txt {
		// Split record into key and value.
		// Ignore keys without values -- CUPS does the same.
		key, value, novalue := txtParse(t)

		if novalue {
			continue
		}

		// Update found/missed sets and check for duplicates
		if _, dup := found[key]; dup {
			continue
		}

		// Decode the value
		var err error
		switch txToLower(key) {
		case "adminurl":
			s.adminURL = value
		case "cs":
			s.params.Colors, err = txtColors(value)
		case "duplex":
			s.params.Duplex, err = txtOption(value)
		case "is":
			s.params.Sources, err = txtSources(value)
		case "note":
			s.location = value
		case "pdl":
			s.params.PDL, err = txtKeywords(value)
		case "representation":
			s.iconURL = value
		case "rs":
			// Strip leading and trailing '/'.
			// sane-airscan does the same.
			rs := value
			for strings.HasPrefix(rs, "/") {
				rs = rs[1:]
			}
			for strings.HasSuffix(rs, "/") {
				rs = rs[:len(rs)-1]
			}

			s.uriPath = value
		case "ty":
			s.makeModel = value
		case "uuid":
			s.uuid, err = uuid.Parse(value)

		// ipp-usb extensions
		case "usb_ser":
			s.usbSerial = value
		case "usb_hwid":
			s.usbHWID = value
		}

		// Check for error
		if err != nil {
			err = fmt.Errorf("%s: %w", key, err)
			return txtScanner{}, err
		}
	}

	return s, nil
}

// txtAuth decodes an authentication mode
func txtAuth(value string) (discovery.AuthMode, error) {
	switch txToLower(value) {
	case "none":
		return discovery.AuthNone, nil
	case "certificate":
		return discovery.AuthCertificate, nil
	case "negotiate":
		return discovery.AuthKerberos, nil
	case "oauth":
		return discovery.AuthOAuth2, nil
	case "username,password":
		return discovery.AuthPasswd, nil
	}

	return discovery.AuthOther, nil
}

// txtOption decodes an Option value value
func txtOption(value string) (optional.Val[bool], error) {
	switch txToLower(value) {
	case "f":
		return optional.New(false), nil
	case "t":
		return optional.New(true), nil
	}
	return nil, nil
}

// txtColors decodes discovery.ColorMode bits
func txtColors(value string) (generic.Bitset[abstract.ColorMode], error) {
	keywords, _ := txtKeywords(value)

	var colors generic.Bitset[abstract.ColorMode]
	for _, kw := range keywords {
		switch txToLower(kw) {
		case "color":
			colors.Add(abstract.ColorModeColor)
		case "grayscale":
			colors.Add(abstract.ColorModeMono)
		case "binary":
			colors.Add(abstract.ColorModeBinary)
		}
	}

	return colors, nil
}

// txtKeywords decodes comma-separated list of keywords
func txtKeywords(value string) ([]string, error) {
	// Split value into comma-separated list
	pdl := strings.Split(value, ",")

	// Drop empty items (just in case)
	o := 0
	for i := range pdl {
		if pdl[i] != "" {
			pdl[o] = pdl[i]
			o++
		}
	}

	pdl = pdl[:o]
	return pdl, nil
}

// txtMediaKind decodes discovery.MediaKind bits
func txtMediaKind(value string) (discovery.MediaKind, error) {
	keywords, _ := txtKeywords(value)

	var media discovery.MediaKind
	for _, kw := range keywords {
		switch txToLower(kw) {
		case "disc":
			media |= discovery.MediaDisk
		case "document":
			media |= discovery.MediaDocument
		case "envelope":
			media |= discovery.MediaEnvelope
		case "label":
			media |= discovery.MediaLabel
		case "large-format":
			media |= discovery.MediaLargeFormat
		case "photo":
			media |= discovery.MediaPhoto
		case "postcard":
			media |= discovery.MediaPostcard
		case "receipt":
			media |= discovery.MediaReceipt
		case "roll":
			media |= discovery.MediaRoll
		default:
			media |= discovery.MediaOther
		}
	}

	return media, nil
}

// txtPPD decodes the max paper size
func txtPaperMax(value string) (discovery.PaperSize, error) {
	switch txToLower(value) {
	case "<legal-a4":
		return discovery.PaperA4Minus, nil
	case "legal-a4":
		return discovery.PaperA4, nil
	case "tabloid-a3":
		return discovery.PaperA3, nil
	case "isoC-a2":
		return discovery.PaperA2, nil
	case ">isoC-a2":
		return discovery.PaperA2Plus, nil
	}

	return discovery.PaperUnknown, nil
}

// txtSources decodes discovery.ScanSource bits
func txtSources(value string) (discovery.ScanSource, error) {
	keywords, _ := txtKeywords(value)

	var sources discovery.ScanSource
	for _, kw := range keywords {
		switch txToLower(kw) {
		case "adf":
			sources |= discovery.ScanADF
		case "platen":
			sources |= discovery.ScanPlaten
		default:
			sources |= discovery.ScanOther
		}
	}

	return sources, nil
}

// txtParse parses "key=value" string into key and value parts.
func txtParse(s string) (key, value string, novalue bool) {
	i := strings.IndexByte(s, '=')
	if i >= 0 {
		return s[0:i], s[i+1:], false
	}

	return s, "", true
}

// txToLower converts ASCII string to lowercase.
func txToLower(s string) string {
	buf := []byte(s)

	for i := range buf {
		if c := buf[i]; 'A' <= c && c <= 'Z' {
			buf[i] = c - 'A' + 'a'
		}
	}

	return string(buf)
}
