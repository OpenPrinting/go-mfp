// MFP - Miulti-Function Printers and scanners toolkit
// DNS-SD service discovery
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// DNS-SD TXT records tests

package dnssd

import (
	"reflect"
	"testing"

	"github.com/OpenPrinting/go-mfp/abstract"
	"github.com/OpenPrinting/go-mfp/discovery"
	"github.com/OpenPrinting/go-mfp/util/optional"
)

// TestTxtParse tests txtParse function
func TestTxtParse(t *testing.T) {
	type testData struct {
		s       string // Input string
		key     string // Expected key
		value   string // Expected value
		novalue bool   // Expected novalue flag
	}

	tests := []testData{
		{"rp=ipp/print", "rp", "ipp/print", false},
		{"note=", "note", "", false},
		{"a=b=c", "a", "b=c", false},
		{"=value", "", "value", false},
		{"fax", "fax", "", true},
		{"", "", "", true},
	}

	for _, test := range tests {
		key, value, novalue := txtParse(test.s)
		if key != test.key || value != test.value || novalue != test.novalue {
			t.Errorf("%q: expected (%q, %q, %v), present (%q, %q, %v)",
				test.s, test.key, test.value, test.novalue,
				key, value, novalue)
		}
	}
}

// TestTxToLower tests txToLower function
func TestTxToLower(t *testing.T) {
	type testData struct {
		s, expected string
	}

	tests := []testData{
		{"", ""},
		{"abc", "abc"},
		{"ABC", "abc"},
		{"UsB_MfG", "usb_mfg"},
		{"isoC-A2", "isoc-a2"},
		{"Ä", "Ä"}, // Only ASCII letters are converted
	}

	for _, test := range tests {
		s := txToLower(test.s)
		if s != test.expected {
			t.Errorf("%q: expected %q, present %q",
				test.s, test.expected, s)
		}
	}
}

// TestTxtKeywords tests txtKeywords function
func TestTxtKeywords(t *testing.T) {
	type testData struct {
		s        string
		expected []string
	}

	tests := []testData{
		{"", []string{}},
		{"application/pdf", []string{"application/pdf"}},
		{"application/pdf,image/jpeg", []string{"application/pdf", "image/jpeg"}},
		{",application/pdf,,image/jpeg,", []string{"application/pdf", "image/jpeg"}},
	}

	for _, test := range tests {
		kw, err := txtKeywords(test.s)
		if err != nil {
			t.Errorf("%q: %s", test.s, err)
			continue
		}
		if !reflect.DeepEqual(kw, test.expected) {
			t.Errorf("%q: expected %q, present %q",
				test.s, test.expected, kw)
		}
	}
}

// TestTxtOption tests txtOption function
func TestTxtOption(t *testing.T) {
	type testData struct {
		s        string
		expected optional.Val[bool]
	}

	tests := []testData{
		{"T", optional.New(true)},
		{"t", optional.New(true)},
		{"F", optional.New(false)},
		{"f", optional.New(false)},
		{"U", nil},
		{"", nil},
	}

	for _, test := range tests {
		opt, err := txtOption(test.s)
		if err != nil {
			t.Errorf("%q: %s", test.s, err)
			continue
		}
		if !reflect.DeepEqual(opt, test.expected) {
			t.Errorf("%q: expected %v, present %v", test.s,
				optional.Get(test.expected), optional.Get(opt))
		}
	}
}

// TestTxtAuth tests txtAuth function
func TestTxtAuth(t *testing.T) {
	type testData struct {
		s        string
		expected discovery.AuthMode
	}

	tests := []testData{
		{"none", discovery.AuthNone},
		{"certificate", discovery.AuthCertificate},
		{"negotiate", discovery.AuthKerberos},
		{"oauth", discovery.AuthOAuth2},
		{"username,password", discovery.AuthPasswd},
		{"Username,Password", discovery.AuthPasswd},
		{"something-else", discovery.AuthOther},
	}

	for _, test := range tests {
		auth, err := txtAuth(test.s)
		if err != nil {
			t.Errorf("%q: %s", test.s, err)
			continue
		}
		if auth != test.expected {
			t.Errorf("%q: expected %v, present %v",
				test.s, test.expected, auth)
		}
	}
}

// TestTxtPaperMax tests txtPaperMax function
func TestTxtPaperMax(t *testing.T) {
	type testData struct {
		s        string
		expected discovery.PaperSize
	}

	// Values are spelled as in the Bonjour Printing Specification
	tests := []testData{
		{"<legal-A4", discovery.PaperA4Minus},
		{"legal-A4", discovery.PaperA4},
		{"tabloid-A3", discovery.PaperA3},
		{"isoC-A2", discovery.PaperA2},
		{">isoC-A2", discovery.PaperA2Plus},
		{"ISOC-A2", discovery.PaperA2},
		{"huge", discovery.PaperUnknown},
	}

	for _, test := range tests {
		paper, err := txtPaperMax(test.s)
		if err != nil {
			t.Errorf("%q: %s", test.s, err)
			continue
		}
		if paper != test.expected {
			t.Errorf("%q: expected %v, present %v",
				test.s, test.expected, paper)
		}
	}
}

// TestTxtMediaKind tests txtMediaKind function
func TestTxtMediaKind(t *testing.T) {
	type testData struct {
		s        string
		expected discovery.MediaKind
	}

	tests := []testData{
		{"", 0},
		{"document", discovery.MediaDocument},
		{"Document,Envelope", discovery.MediaDocument | discovery.MediaEnvelope},
		{"disc,label,photo", discovery.MediaDisk | discovery.MediaLabel | discovery.MediaPhoto},
		{"roll,unknown", discovery.MediaRoll | discovery.MediaOther},
	}

	for _, test := range tests {
		media, err := txtMediaKind(test.s)
		if err != nil {
			t.Errorf("%q: %s", test.s, err)
			continue
		}
		if media != test.expected {
			t.Errorf("%q: expected %v, present %v",
				test.s, test.expected, media)
		}
	}
}

// TestTxtColorsInputs tests txtColors and txtInputs functions
func TestTxtColorsInputs(t *testing.T) {
	colors, _ := txtColors("color,Grayscale,binary,sepia")
	for _, c := range []abstract.ColorMode{abstract.ColorModeColor,
		abstract.ColorModeMono, abstract.ColorModeBinary} {
		if !colors.Contains(c) {
			t.Errorf("txtColors: %v missed in %v", c, colors)
		}
	}
	if colors.Count() != 3 {
		t.Errorf("txtColors: expected 3 modes, present %v", colors)
	}

	inputs, _ := txtInputs("platen,ADF,camera")
	if !inputs.Contains(abstract.InputPlaten) ||
		!inputs.Contains(abstract.InputADF) || inputs.Count() != 2 {
		t.Errorf("txtInputs: expected platen and ADF, present %v",
			inputs)
	}
}

// TestDecodeTxtPrinter tests decodeTxtPrinter function
func TestDecodeTxtPrinter(t *testing.T) {
	txt := []string{
		"txtvers=1",
		"rp=ipp/print",
		"ty=Example Printer 1000",
		"note=2nd Floor",
		"UUID=4509a320-00a0-008f-00b6-002507510eca",
		"Color=T",
		"Duplex=F",
		"priority=10",
		"pdl=application/pdf,image/urf",
		"PaperMax=isoC-A2",
		"usb_MFG=Example",
		"usb_MDL=Printer 1000",
		"fax",
	}

	p, err := decodeTxtPrinter(svcTypeIPP, "Example", txt)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if !p.IsPrinter() || p.IsFaxout() {
		t.Errorf("expected printer only, present printer=%v faxout=%v",
			p.IsPrinter(), p.IsFaxout())
	}

	check := func(name string, present, expected any) {
		if !reflect.DeepEqual(present, expected) {
			t.Errorf("%s: expected %v, present %v",
				name, expected, present)
		}
	}

	check("rp", p.rp, "ipp/print")
	check("makeModel", p.makeModel, "Example Printer 1000")
	check("location", p.location, "2nd Floor")
	check("uuid", p.uuid.String(), "4509a320-00a0-008f-00b6-002507510eca")
	check("usbMFG", p.usbMFG, "Example")
	check("usbMDL", p.usbMDL, "Printer 1000")
	check("Color", p.params.Color, optional.Val[bool](optional.New(true)))
	check("Duplex", p.params.Duplex, optional.Val[bool](optional.New(false)))
	check("Priority", p.params.Priority, 10)
	check("PDL", p.params.PDL, []string{"application/pdf", "image/urf"})
	check("Paper", p.params.Paper, discovery.PaperA2)
	check("Queue", p.params.Queue, "")
}

// TestDecodeTxtPrinterDefaults tests decodeTxtPrinter defaults
func TestDecodeTxtPrinterDefaults(t *testing.T) {
	type testData struct {
		svcType  string
		txt      []string
		priority int
		queue    string
		printer  bool
		faxout   bool
	}

	tests := []testData{
		{svcTypeIPP, nil, 50, "ipp/print", false, false},
		{svcTypeIPPS, []string{"rfo=ipp/faxout"}, 50, "", false, true},
		{svcTypeIPP, []string{"rp=ipp/print", "rfo=ipp/faxout"}, 50, "", true, true},
		{svcTypeLPD, nil, 50, "auto", true, false},
	}

	for _, test := range tests {
		p, err := decodeTxtPrinter(test.svcType, "Example", test.txt)
		if err != nil {
			t.Errorf("%s %q: %s", test.svcType, test.txt, err)
			continue
		}

		if p.params.Priority != test.priority ||
			p.params.Queue != test.queue ||
			p.IsPrinter() != test.printer ||
			p.IsFaxout() != test.faxout {
			t.Errorf("%s %q: expected (%d, %q, %v, %v), "+
				"present (%d, %q, %v, %v)",
				test.svcType, test.txt,
				test.priority, test.queue, test.printer, test.faxout,
				p.params.Priority, p.params.Queue,
				p.IsPrinter(), p.IsFaxout())
		}
	}
}

// TestDecodeTxtPrinterErrors tests decodeTxtPrinter error handling
func TestDecodeTxtPrinterErrors(t *testing.T) {
	tests := [][]string{
		{"priority=100"},
		{"priority=-1"},
		{"priority=high"},
		{"txtvers=2"},
		{"uuid=not-a-uuid"},
	}

	for _, txt := range tests {
		_, err := decodeTxtPrinter(svcTypeIPP, "Example", txt)
		if err == nil {
			t.Errorf("%q: error expected", txt)
		}
	}
}

// TestDecodeTxtDuplicates tests that only the first occurrence of
// a key is used. Keys are case-insensitive (RFC 6763, 6.4).
func TestDecodeTxtDuplicates(t *testing.T) {
	p, err := decodeTxtPrinter(svcTypeIPP, "Example",
		[]string{"ty=First", "ty=Second", "TY=Third", "priority=10",
			"Priority=bad"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	if p.makeModel != "First" {
		t.Errorf("printer ty: expected %q, present %q",
			"First", p.makeModel)
	}

	s, err := decodeTxtScanner(svcTypeESCL, "Example",
		[]string{"rs=eSCL", "RS=other"})
	if err != nil {
		t.Fatalf("%s", err)
	}

	if s.uriPath != "eSCL" {
		t.Errorf("scanner rs: expected %q, present %q",
			"eSCL", s.uriPath)
	}
}

// TestDecodeTxtScanner tests decodeTxtScanner function
func TestDecodeTxtScanner(t *testing.T) {
	type testData struct {
		rs   []string // Optional rs=... TXT entry
		path string   // Expected uriPath
	}

	tests := []testData{
		{nil, "eSCL"},
		{[]string{"rs=eSCL"}, "eSCL"},
		{[]string{"rs=/eSCL"}, "eSCL"},
		{[]string{"rs=/eSCL/"}, "eSCL"},
		{[]string{"rs=//scan/eSCL//"}, "scan/eSCL"},
	}

	for _, test := range tests {
		txt := append([]string{
			"ty=Example Scanner",
			"cs=color,grayscale",
			"is=platen,adf",
			"duplex=T",
			"representation=http://example.com/icon.png",
		}, test.rs...)

		s, err := decodeTxtScanner(svcTypeESCL, "Example", txt)
		if err != nil {
			t.Errorf("%q: %s", txt, err)
			continue
		}

		if s.uriPath != test.path {
			t.Errorf("%q: uriPath: expected %q, present %q",
				test.rs, test.path, s.uriPath)
		}

		if s.makeModel != "Example Scanner" ||
			s.iconURL != "http://example.com/icon.png" ||
			s.params.Colors.Count() != 2 ||
			s.params.Inputs.Count() != 2 ||
			!reflect.DeepEqual(s.params.Duplex, optional.Val[bool](optional.New(true))) {
			t.Errorf("%q: unexpected decoding: %+v %+v",
				txt, s, s.params)
		}
	}
}
