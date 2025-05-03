// MFP - Miulti-Function Printers and scanners toolkit
// The "cups" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Common options

package cups

import (
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/cups"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/transport"
)

// optAttrs describes the --attrs option.
// It specifies a list of requested attributes.
var optAttrs = argv.Option{
	Name:     "--attrs",
	Help:     "Additional attributes",
	HelpArg:  "attr,...",
	Validate: argv.ValidateAny,
	Complete: optAttrsComplete,
}

// optAttrsGet returns --attrs option (list of requested attributes).
func optAttrsGet(inv *argv.Invocation) (attrs []string) {
	params := inv.Values("--attrs")
	for _, param := range params {
		for _, name := range strings.Split(param, ",") {
			if name != "" {
				attrs = append(attrs, name)
			}
		}
	}

	return
}

// optAttrsGet is the completion callback for the --attrs option.
func optAttrsComplete(arg string) (compl []argv.Completion) {
	infos := ((*ipp.PrinterAttributes)(nil)).KnownAttrs()

	attrName := arg
	prefix := ""

	if i := strings.LastIndex(attrName, ","); i >= 0 {
		attrName = arg[i+1:]
		prefix = arg[:i+1]
	}

	for _, info := range infos {
		if strings.HasPrefix(info.Name, attrName) {
			c := argv.Completion{
				String: prefix + info.Name + ",",
				Flags:  argv.CompletionNoSpace,
			}
			compl = append(compl, c)
		}
	}

	return
}

// optID describes the --id option.
// It specifies the printer-id.
var optID = argv.Option{
	Name:     "--id",
	Help:     "Printer ID (1...65535)",
	HelpArg:  "id",
	Validate: argv.ValidateIntRange(0, 1, 65535),
}

// optIDGet returns --id option value.
func optIDGet(inv *argv.Invocation) int {
	id := 0
	if opt, ok := inv.Get("--id"); ok {
		id, _ = strconv.Atoi(opt)
	}
	return id
}

// optLimit describes the --limit option.
// It specifies the maximum number of returned printers
var optLimit = argv.Option{
	Name:     "--limit",
	Help:     "Maximum number of printers",
	HelpArg:  "N",
	Validate: argv.ValidateIntRange(0, 1, math.MaxInt32),
}

// optLimitGet returns --limit option value.
func optLimitGet(inv *argv.Invocation) int {
	lim := 0
	if opt, ok := inv.Get("--limit"); ok {
		lim, _ = strconv.Atoi(opt)
	}
	return lim
}

// optLocation describes the --location option.
// It specified the desired printer location (e.g. "2nd Floor Computer Lab")
var optLocation = argv.Option{
	Name: "--location",
	Help: "" +
		`Printer location ` +
		`(e.g., "2nd Floor Computer Lab")`,
	HelpArg:  "where",
	Validate: argv.ValidateAny,
}

// optLocationGet returns --location option value.
func optLocationGet(inv *argv.Invocation) string {
	opt, _ := inv.Get("--location")
	return opt
}

// optSchemesExclude describes the --exclude-schemes=scheme,... option
// It specifies URL schemes to be excluded
var optSchemesExclude = argv.Option{
	Name:     "--exclude-schemes",
	Help:     "URL schemes to exclude",
	HelpArg:  "scheme,...",
	Validate: argv.ValidateAny,
	Complete: optSchemesComplete,
}

// optSchemesExcludeGet returns --exclude-schemes option value
func optSchemesExcludeGet(inv *argv.Invocation) []string {
	return optSchemesGet(inv, "--exclude-schemes")
}

// optSchemesInclude describes the --include-schemes=scheme,... option
// It specifies URL schemes to be included.
var optSchemesInclude = argv.Option{
	Name:     "--include-schemes",
	Help:     "URL schemes to include",
	HelpArg:  "scheme,...",
	Validate: argv.ValidateAny,
	Complete: optSchemesComplete,
}

// optSchemesIncludeGet returns --include-schemes option value
func optSchemesIncludeGet(inv *argv.Invocation) []string {
	return optSchemesGet(inv, "--include-schemes")
}

// optSchemesGet is the common function for optSchemesExcludeGet
// and optSchemesIncludeGet
func optSchemesGet(inv *argv.Invocation, optname string) (schemes []string) {
	for _, val := range inv.Values(optname) {
		for _, s := range strings.Split(val, ",") {
			if s != "" {
				schemes = append(schemes, s)
			}
		}
	}

	return
}

// optSchemesComplete is a common completer for --exclude-schemes and
// --include-schemes options
func optSchemesComplete(arg string) (compl []argv.Completion) {
	var prefix, scheme string

	if i := strings.LastIndex(arg, ","); i >= 0 {
		scheme = arg[i+1:]
		prefix = arg[:i+1]
	}

	candidates := []string{
		"http",
		"https",
		"ipp",
		"ipps",
		"lpd",
		"smb",
		"socket",
	}

	for _, candidate := range candidates {
		if strings.HasPrefix(candidate, scheme) {
			c := argv.Completion{
				String: prefix + candidate + ",",
				Flags:  argv.CompletionNoSpace,
			}
			compl = append(compl, c)
		}
	}

	return
}

// optTimeout describes the --timeout=seconds option.
// It specifies operation timeout
var optTimeout = argv.Option{
	Name:     "--timeout",
	Help:     "operation timeout",
	HelpArg:  "seconds",
	Validate: argv.ValidateIntRange(0, 1, math.MaxInt32),
}

// optTimeoutGet returns --timeout option value.
func optTimeoutGet(inv *argv.Invocation) time.Duration {
	if opt, ok := inv.Get(optTimeout.Name); ok {
		v, _ := strconv.Atoi(opt)
		return time.Duration(v) * time.Second
	}
	return cups.DefaultGetDevicesTimeout
}

// optUser describes the --user option.
// It allows to filter printers by the user name
// Only printers accessible to that user will be returned.
var optUser = argv.Option{
	Name:     "--user",
	Help:     "Show only printers accessible to that user",
	HelpArg:  "name",
	Validate: argv.ValidateAny,
}

// optUserGet returns --user option value.
func optUserGet(inv *argv.Invocation) string {
	opt, _ := inv.Get("--user")
	return opt
}

// optCUPSURL returns CUPS URL (-u/--cups option).
// If option is not set, it uses default destination.
func optCUPSURL(inv *argv.Invocation) *url.URL {
	dest := transport.DefaultCupsUNIX

	if addr, ok := inv.Parent().Get("-u"); ok {
		dest = transport.MustParseAddr(addr, "ipp://localhost/")
	}

	return dest
}
