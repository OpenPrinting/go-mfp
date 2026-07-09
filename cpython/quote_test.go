// MFP - Miulti-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Python-style strings quoting test

package cpython

import (
	"testing"
)

// TestQuoteUTF8 validates that QuoteSingle and QuoteDouble correctly wrap strings,
// preserve printable UTF-8 characters (like Cyrillic text and emojis) as readable text,
// and escape standard control characters and quotes according to Python rules.
func TestQuoteUTF8(t *testing.T) {
	type testData struct {
		input      string
		wantSingle string
		wantDouble string
	}

	tests := []testData{
		// Standard ASCII string with no special characters
		{"hello", "'hello'", `"hello"`},
		// Multi-byte Cyrillic characters must remain unescaped and fully readable
		{"привет", "'привет'", `"привет"`},
		// Multi-byte Emoji characters must remain unescaped and fully readable
		{"hello 👋", "'hello 👋'", `"hello 👋"`},
		// Control sequences like newlines must be escaped explicitly
		{"line1\nline2", `'line1\nline2'`, `"line1\nline2"`},
		// Double quotes inside single-quoted string should not be escaped, and vice-versa
		{`He said "Привет"`, `'He said "Привет"'`, `"He said \"Привет\""`},
		// Single quotes inside double-quoted string should not be escaped
		{"It's fine", "'It\\'s fine'", `"It's fine"`},
		// Invalid UTF-8 bytes must fall back to exact raw hex escapes without breaking the string
		{"bad\x80byte", `'bad\x80byte'`, `"bad\x80byte"`},
		// All control sequences
		{"---\n\r\t\a\b\f\\---", `'---\n\r\t\x07\x08\x0c\\---'`, `"---\n\r\t\x07\x08\x0c\\---"`},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			// Using %q formatting instead of %s to avoid raw character literal printing anomalies
			if got := QuoteSingle(test.input); got != test.wantSingle {
				t.Errorf("QuoteSingle():\n"+
					"expected: %q\n"+
					"present:  %q",
					got, test.wantSingle)
			}
			if got := QuoteDouble(test.input); got != test.wantDouble {
				t.Errorf("QuoteSingle():\n"+
					"expected: %q\n"+
					"present:  %q",
					got, test.wantDouble)
			}
		})
	}
}

// TestUnquoteUTF8 validates that Unquote accurately parses Python string literals
// back into regular Go strings. It ensures that standard escaped control codes,
// readable multi-byte UTF-8 text, legacy 16-bit unicode escapes (\uXXXX),
// and 32-bit unicode escapes (\UXXXXXXXX) are all parsed correctly.
// It also tests error handling for corrupted inputs, mismatched quotes, and bad escape patterns.
func TestUnquoteUTF8(t *testing.T) {
	type testData struct {
		name    string
		input   string
		want    string
		wantErr bool
	}

	tests := []testData{
		// Literal single-quoted readable Cyrillic text
		{"Readable Cyrillic Single", `'привет'`, "привет", false},
		// Literal double-quoted readable Cyrillic text
		{"Readable Cyrillic Double", `"привет"`, "привет", false},
		// Inline printable Emoji character handling
		{"Emoji Single", `'👋'`, "👋", false},
		// Python's 16-bit hex Unicode sequence handling
		{"Legacy Hex Unicode Escape", `'\u043f\u0440'`, "пр", false},
		// Python's 32-bit hex Unicode sequence handling (commonly used for emojis)
		{"32bit Emoji Escape", `'\U0001F44B'`, "👋", false},
		// Error case: Unmatched outer quote types
		{"Invalid Quotes Mismatch - 1", `"hello'`, "", true},
		{"Invalid Quotes Mismatch - 1", `'hello"`, "", true},
		// Error case: Input sequence too short to be a valid Python literal
		{"Invalid Empty String", ``, "", true},
		{"Invalid Short String", `'`, "", true},
		// Error case: Hex digits truncated or containing non-hex characters
		{"Invalid Unicode Escape", `'\u123'`, "", true},
		// Non-standard/Unknown python escapes should preserve the backslash literally
		{"Unknown escape stays literal", `'\z'`, `\z`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Unquote(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("%s: error = %v, wantErr %v",
					test.name, err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("%s: expected %q, present %q",
					test.name, got, test.want)
			}
		})
	}
}
