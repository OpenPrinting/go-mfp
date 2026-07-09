// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Python-style strings quoting

package cpython

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// QuoteSingle wraps a string in single quotes (') and escapes it like a Python string literal.
// Valid UTF-8 characters remain readable.
func QuoteSingle(s string) string {
	var sb strings.Builder
	sb.WriteByte('\'')

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			// Handle invalid UTF-8 byte as a raw hex escape safely using the exact byte index
			sb.WriteString(fmt.Sprintf(`\x%02x`, s[0]))
			s = s[1:]
			continue
		}
		s = s[size:]

		switch r {
		case '\'':
			sb.WriteString(`\'`)
		case '\\':
			sb.WriteString(`\\`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			// Keep printable UTF-8 characters readable (e.g., Cyrillic, Emojis)
			if unicode.IsPrint(r) {
				sb.WriteRune(r)
			} else {
				// Escape non-printable control characters
				if r <= 0xFF {
					sb.WriteString(fmt.Sprintf(`\x%02x`, r))
				} else if r <= 0xFFFF {
					sb.WriteString(fmt.Sprintf(`\u%04x`, r))
				} else {
					sb.WriteString(fmt.Sprintf(`\U%08x`, r))
				}
			}
		}
	}

	sb.WriteByte('\'')
	return sb.String()
}

// QuoteDouble wraps a string in double quotes (") and escapes it like a Python string literal.
// Valid UTF-8 characters remain readable.
func QuoteDouble(s string) string {
	var sb strings.Builder
	sb.WriteByte('"')

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			// Handle invalid UTF-8 byte as a raw hex escape safely using the exact byte index
			sb.WriteString(fmt.Sprintf(`\x%02x`, s[0]))
			s = s[1:]
			continue
		}
		s = s[size:]

		switch r {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			if unicode.IsPrint(r) {
				sb.WriteRune(r)
			} else {
				if r <= 0xFF {
					sb.WriteString(fmt.Sprintf(`\x%02x`, r))
				} else if r <= 0xFFFF {
					sb.WriteString(fmt.Sprintf(`\u%04x`, r))
				} else {
					sb.WriteString(fmt.Sprintf(`\U%08x`, r))
				}
			}
		}
	}

	sb.WriteByte('"')
	return sb.String()
}

// Unquote converts a Python string literal back to a Go UTF-8 string.
// Supports \x, \u, and \U Python escape sequences.
func Unquote(s string) (string, error) {
	if len(s) < 2 {
		return "", fmt.Errorf("string literal too short")
	}

	// Capture the first byte as the quote identifier
	quote := s[0]
	if quote != '\'' && quote != '"' {
		return "", fmt.Errorf("invalid Python string literal quote character")
	}
	if s[len(s)-1] != quote {
		return "", fmt.Errorf("mismatched closing quote")
	}

	// Strip the outer quotes
	content := s[1 : len(s)-1]
	var sb strings.Builder

	for i := 0; i < len(content); i++ {
		if content[i] == '\\' {
			if i+1 >= len(content) {
				return "", fmt.Errorf("trailing backslash in literal")
			}
			next := content[i+1]
			i++ // consume the backslash

			switch next {
			case quote:
				sb.WriteByte(quote)
			// Handle cases where the alternate quote character was escaped anyway
			case '\'', '"':
				sb.WriteByte(next)
			case '\\':
				sb.WriteByte('\\')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'a':
				sb.WriteByte('\a')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'v':
				sb.WriteByte('\v')
			case 'x': // 8-bit hex escape: \xhh
				if i+2 >= len(content) {
					return "", fmt.Errorf("invalid hex escape sequence")
				}
				hexStr := content[i+1 : i+3]
				i += 2
				val, err := strconv.ParseUint(hexStr, 16, 8)
				if err != nil {
					return "", fmt.Errorf("invalid hex value: %s", hexStr)
				}
				sb.WriteByte(byte(val))
			case 'u': // 16-bit unicode escape: \uxxxx
				if i+4 >= len(content) {
					return "", fmt.Errorf("invalid u16 escape sequence")
				}
				hexStr := content[i+1 : i+5]
				i += 4
				val, err := strconv.ParseUint(hexStr, 16, 16)
				if err != nil {
					return "", fmt.Errorf("invalid u16 value: %s", hexStr)
				}
				sb.WriteRune(rune(val))
			case 'U': // 32-bit unicode escape: \Uxxxxxxxx
				if i+8 >= len(content) {
					return "", fmt.Errorf("invalid u32 escape sequence")
				}
				hexStr := content[i+1 : i+9]
				i += 8
				val, err := strconv.ParseUint(hexStr, 16, 32)
				if err != nil || val > unicode.MaxRune {
					return "", fmt.Errorf("invalid u32 value: %s", hexStr)
				}
				sb.WriteRune(rune(val))
			default:
				sb.WriteByte('\\')
				sb.WriteByte(next)
			}
		} else {
			sb.WriteByte(content[i])
		}
	}

	return sb.String(), nil
}
