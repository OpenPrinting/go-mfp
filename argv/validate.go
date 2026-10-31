// MFP  - Miulti-Function Printers and scanners toolkit
// argv - Argv parsing mini-library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Value validators

package argv

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ValidateAny is the Option.Validate and Parameter.Validate callback.
//
// It accepts any string.
func ValidateAny(string) error {
	return nil
}

// ValidateInt8 is the Option.Validate and Parameter.Validate callback.
//
// It accepts signed 8-bit integers.
func ValidateInt8(in string) error {
	return validateInt8(in)
}

// ValidateInt16 is the Option.Validate and Parameter.Validate callback.
//
// It accepts signed 16-bit integers.
func ValidateInt16(in string) error {
	return validateInt16(in)
}

// ValidateInt32 is the Option.Validate and Parameter.Validate callback.
//
// It accepts signed 32-bit integers.
func ValidateInt32(in string) error {
	return validateInt32(in)
}

// ValidateInt64 is the Option.Validate and Parameter.Validate callback.
//
// It accepts signed 64-bit integers.
func ValidateInt64(in string) error {
	return validateInt64(in)
}

// ValidateUint8 is the Option.Validate and Parameter.Validate callback.
//
// It accepts unsigned 8-bit integers.
func ValidateUint8(in string) error {
	return validateUint8(in)
}

// ValidateUint16 is the Option.Validate and Parameter.Validate callback.
//
// It accepts unsigned 16-bit integers.
func ValidateUint16(in string) error {
	return validateUint16(in)
}

// ValidateUint32 is the Option.Validate and Parameter.Validate callback.
//
// It accepts unsigned 32-bit integers.
func ValidateUint32(in string) error {
	return validateUint32(in)
}

// ValidateUint64 is the Option.Validate and Parameter.Validate callback.
//
// It accepts unsigned 64-bit integers.
func ValidateUint64(in string) error {
	return validateUint64(in)
}

// ValidateBool is the Option.Validate and Parameter.Validate callback.
//
// It assepts boolean values.
//
// Valid boolean values are (case-insensetive):
//   - t, true  - for true
//   - f, false - for false
//   - any unsigned integer. 0 is false, everything else is true.
func ValidateBool(in string) error {
	switch strings.ToLower(in) {
	case "t", "true":
		return nil
	case "f", "false":
		return nil
	}

	return ValidateUint64(in)
}

// DecodeBool decodes boolean value of the option.
// This is assumed that value has passed the [ValidateBool] validation.
func DecodeBool(in string) bool {
	switch strings.ToLower(in) {
	case "t", "true":
		return true
	case "f", "false":
		return false
	}

	v, _ := strconv.ParseUint(in, 0, 64)
	return v != 0
}

// ValidateStrings creates the Option.Validate and Parameter.Validate callback.
//
// It returns validator that accepts any of supplied strings.
func ValidateStrings(s []string) func(string) error {
	// Create a copy of input, to protect from callers
	// that may change the slice after the call.
	set := make([]string, len(s))
	copy(set, s)

	// Create validator
	return func(in string) error {
		for _, member := range set {
			if in == member {
				return nil
			}
		}

		return errors.New("invalid argument")
	}
}

// ValidateIntRange creates the Option.Validate and Parameter.Validate
// callback.
//
// It returns validator, that accepts signed integer in range
// (min <= x && x <= max). If base is 0, the actual base is implied
// automatically based on a string prefix following the sign
// (2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise).
func ValidateIntRange(base int, min, max int64) func(string) error {
	return func(in string) error {
		v, err := strconv.ParseInt(in, base, 64)
		if err != nil {
			return errors.New("invalid integer")
		}

		if v < min || v > max {
			return fmt.Errorf(
				"value out of range (%d...%d)", min, max)
		}

		return nil
	}
}

// ValidateUintRange creates the Option.Validate and Parameter.Validate
// callback.
//
// It returns validator, that accepts unsigned integer in range
// (min <= x && x <= max). If base is 0, the actual base is implied
// automatically based on a string prefix following the sign
// (2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise).
//
// It is like [ValidateIntRange] but for unsigned numbers. The sign
// prefix is not accepted.
func ValidateUintRange(base int, min, max uint64) func(string) error {
	return func(in string) error {
		v, err := strconv.ParseUint(in, base, 64)
		if err != nil {
			return errors.New("invalid integer")
		}

		if v < min || v > max {
			return fmt.Errorf(
				"value out of range (%d...%d)", min, max)
		}

		return nil
	}
}

// ValidateIntBits returns validator, that accepts signed integer
// that fits specified number of bits (1 to 64).  If base is 0, the
// actual base is implied automatically based on a string prefix
// (2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise).
//
// It panics, if number of bits is out of range
func ValidateIntBits(base, bits int) func(string) error {
	if bits < 1 || bits > 64 {
		err := fmt.Errorf(
			"ValidateIntBits: bits (%d) out of range (1...64)",
			bits)
		panic(err)
	}

	var min, max int64
	switch bits {
	case 1:
		min, max = 0, 1
	default:
		max = int64(1<<uint64(bits-1) - 1)
		min = -max - 1
	}

	return func(in string) error {
		v, err := strconv.ParseInt(in, base, 64)
		if err != nil {
			return errors.New("invalid integer")
		}

		if v < min || v > max {
			return fmt.Errorf(
				"value doesn't fit %d bits", bits)
		}

		return nil
	}
}

// ValidateUintBits returns validator, that accepts unsigned integer
// that fits specified number of bits (1 to 64).  If base is 0, the
// actual base is implied automatically based on a string prefix
// (2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise).
//
// It panics, if number of bits is out of range
func ValidateUintBits(base, bits int) func(string) error {
	if bits < 1 || bits > 64 {
		err := fmt.Errorf(
			"ValidateUintBits: bits (%d) out of range (1...64)",
			bits)
		panic(err)
	}

	var max uint64
	switch bits {
	case 64:
		max = math.MaxUint64

	default:
		max = (1 << uint64(bits)) - 1
	}

	return func(in string) error {
		v, err := strconv.ParseUint(in, base, 64)
		if err != nil {
			return errors.New("invalid integer")
		}

		if v > max {
			return fmt.Errorf(
				"value doesn't fit %d bits", bits)
		}

		return nil
	}
}

var (
	// Some commonly used precomputed validators.
	validateInt8   = ValidateIntBits(0, 8)
	validateInt16  = ValidateIntBits(0, 16)
	validateInt32  = ValidateIntBits(0, 32)
	validateInt64  = ValidateIntBits(0, 64)
	validateUint8  = ValidateUintBits(0, 8)
	validateUint16 = ValidateUintBits(0, 16)
	validateUint32 = ValidateUintBits(0, 32)
	validateUint64 = ValidateUintBits(0, 64)
)
