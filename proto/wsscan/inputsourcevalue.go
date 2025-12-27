// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// InputSource value

package wsscan

// InputSourceValue defines the source of the original document on the scanning device.
type InputSourceValue int

// Known input source values
const (
	UnknownInputSourceValue InputSourceValue = iota // Unknown input source
	InputSourceADF                                  // Document delivered by ADF, front side only
	InputSourceADFDuplex                            // Document delivered by ADF, both sides
	InputSourceFilm                                 // Document scanned using film scanning option
	InputSourcePlaten                               // Document scanned from scanner platen
)

// String returns a string representation of the [InputSourceValue].
func (isv InputSourceValue) String() string {
	switch isv {
	case InputSourceADF:
		return "ADF"
	case InputSourceADFDuplex:
		return "ADFDuplex"
	case InputSourceFilm:
		return "Film"
	case InputSourcePlaten:
		return "Platen"
	}
	return "Unknown"
}

// DecodeInputSourceValue decodes [InputSourceValue] out of its XML string representation.
func DecodeInputSourceValue(s string) InputSourceValue {
	switch s {
	case "ADF":
		return InputSourceADF
	case "ADFDuplex":
		return InputSourceADFDuplex
	case "Film":
		return InputSourceFilm
	case "Platen":
		return InputSourcePlaten
	}
	return UnknownInputSourceValue
}
