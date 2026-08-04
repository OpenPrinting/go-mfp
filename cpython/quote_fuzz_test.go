// MFP - Multi-Function Printers and scanners toolkit
// CPython binding.
//
// Copyright (C) 2026 and up by Abhishrestha Tiwari
// See LICENSE for license terms and conditions
//
// Python-style strings quoting fuzz test

package cpython

import "testing"

// FuzzQuoteSingleRoundTrip fuzzes QuoteSingle+Unquote: any Go string,
// including invalid UTF-8, must survive quoting and unquoting unchanged.
func FuzzQuoteSingleRoundTrip(f *testing.F) {
	f.Add("hello")
	f.Add("It's fine")
	f.Add("bad\x80byte")
	f.Add("emb\x00edded\nmsg")
	f.Add("\U0001F44B привет")

	f.Fuzz(func(t *testing.T, s string) {
		q := QuoteSingle(s)
		got, err := Unquote(q)
		if err != nil {
			t.Fatalf("Unquote(%q): %s", q, err)
		}
		if got != s {
			t.Fatalf("round-trip mismatch: %q != %q", got, s)
		}
	})
}

// FuzzQuoteDoubleRoundTrip is FuzzQuoteSingleRoundTrip for QuoteDouble.
func FuzzQuoteDoubleRoundTrip(f *testing.F) {
	f.Add("hello")
	f.Add(`He said "hi"`)
	f.Add("bad\x80byte")
	f.Add("emb\x00edded\nmsg")

	f.Fuzz(func(t *testing.T, s string) {
		q := QuoteDouble(s)
		got, err := Unquote(q)
		if err != nil {
			t.Fatalf("Unquote(%q): %s", q, err)
		}
		if got != s {
			t.Fatalf("round-trip mismatch: %q != %q", got, s)
		}
	})
}

// FuzzUnquoteNoPanic feeds arbitrary, likely-malformed literals directly
// into Unquote. It must only ever return an error, never panic, on any
// input (index-out-of-range in hex-escape parsing is the main risk).
func FuzzUnquoteNoPanic(f *testing.F) {
	f.Add(`'\x'`)
	f.Add(`'\u12'`)
	f.Add(`'\U1234567'`)
	f.Add(`"mismatched'`)
	f.Add(``)
	f.Add(`'`)
	f.Add(`\`)
	f.Add(`'\`)

	f.Fuzz(func(t *testing.T, s string) {
		_, _ = Unquote(s) // must not panic
	})
}
