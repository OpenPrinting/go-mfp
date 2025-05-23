// MFP - Miulti-Function Printers and scanners toolkit
// WSD core protocol
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Message Body tests

package wsd

import (
	"reflect"
	"testing"
)

// TestBodyAction tests Body.Actions for types that implement
// the Body interface
func TestBodyAction(t *testing.T) {
	type testData struct {
		act  Action
		body Body
	}

	tests := []testData{
		{ActBye, Bye{}},
		{ActGet, Get{}},
		{ActGetResponse, Metadata{}},
		{ActHello, Hello{}},
		{ActProbe, Probe{}},
		{ActProbeMatches, ProbeMatches{}},
		{ActResolve, Resolve{}},
		{ActResolveMatches, ResolveMatches{}},
	}

	for _, test := range tests {
		act := test.body.Action()
		if act != test.act {
			t.Errorf("%s.Action: expected %s, present %s",
				reflect.TypeOf(test.body),
				test.act, act)
		}
	}
}
