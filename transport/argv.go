// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// argv validators for some common types

package transport

import (
	"fmt"
	"strconv"
	"strings"
)

// ArgvValidatePortRange is the [argv.Option.Validate] and
// [argv.Parameter.Validate] callback for range of TCP ports.
//
// It is intended for parameter checking for port range
// options, usable with [ListenTCPRange] and [ListenTCPRangeN]
// functions.
func ArgvValidatePortRange(s string) error {
	invalidPortRange := fmt.Errorf("%s: invalid port range", s)

	n := strings.IndexByte(s, '-')
	switch {
	case n < 0:
		v, err := strconv.ParseUint(s, 10, 16)
		if err != nil || v == 0 || v > 65535 {
			return invalidPortRange
		}
		return nil

	case n == 0 || n == len(s):
		return invalidPortRange
	}

	min := s[:n]
	max := s[n+1:]

	vmin, err := strconv.ParseUint(min, 10, 16)
	if err != nil || vmin == 0 || vmin > 65535 {
		return invalidPortRange
	}

	vmax, err := strconv.ParseUint(max, 10, 16)
	if err != nil || vmax == 0 || vmax > 65535 {
		return invalidPortRange
	}

	if vmin > vmax {
		return invalidPortRange
	}

	return nil
}
