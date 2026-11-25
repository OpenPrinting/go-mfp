// MFP - Multi-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Tests for control file parser

package lpd

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestParseControlFile tests ParseControlFile
func TestParseControlFile(t *testing.T) {
	type testData struct {
		name string
		in   string
		cf   *ControlFile
		err  error
	}

	tests := []testData{
		{
			name: "empty",
			in:   "",
			cf:   &ControlFile{},
		},

		{
			name: "cups",
			in: "Hmyhost\n" +
				"Pgo-mfp\n" +
				"Jreport.pdf\n" +
				"ldfA123myhost\n" +
				"UdfA123myhost\n" +
				"Nreport.pdf\n",
			cf: &ControlFile{
				Host:    "myhost",
				User:    "go-mfp",
				JobName: "report.pdf",
				Documents: []Document{
					{'l', "dfA123myhost", "report.pdf"},
				},
			},
		},

		{
			// BSD lpr -#2 with two files: N after print lines
			name: "bsd-copies",
			in: "Hhost\n" +
				"Proot\n" +
				"Jjob\n" +
				"L root\n" +
				"fdfA001host\n" +
				"fdfA001host\n" +
				"UdfA001host\n" +
				"Na.txt\n" +
				"fdfB001host\n" +
				"fdfB001host\n" +
				"UdfB001host\n" +
				"Nb.txt\n",
			cf: &ControlFile{
				Host:    "host",
				User:    "root",
				JobName: "job",
				Banner:  true,
				Documents: []Document{
					{'f', "dfA001host", "a.txt"},
					{'f', "dfA001host", "a.txt"},
					{'f', "dfB001host", "b.txt"},
					{'f', "dfB001host", "b.txt"},
				},
			},
		},

		{
			// N before print lines
			name: "name-before",
			in: "Na.ps\n" +
				"odfA001host\n" +
				"odfA001host\n" +
				"Nb.txt\n" +
				"pdfB001host\n",
			cf: &ControlFile{
				Documents: []Document{
					{'o', "dfA001host", "a.ps"},
					{'o', "dfA001host", "a.ps"},
					{'p', "dfB001host", "b.txt"},
				},
			},
		},

		{
			// CRLF line endings, empty lines, unknown commands
			name: "crlf-and-junk",
			in: "Hhost\r\n" +
				"\r\n" +
				"\n" +
				"Mroot\r\n" +
				"W132\n" +
				"I8\n" +
				"Cclass\n" +
				"Ttitle\n" +
				"1font\n" +
				"Z?\n" +
				"vdfA001host\r\n" +
				"Proot",
			cf: &ControlFile{
				Host: "host",
				User: "root",
				Documents: []Document{
					{'v', "dfA001host", ""},
				},
			},
		},

		{
			name: "all-formats",
			in: "cdf1\nddf2\nfdf3\ngdf4\nldf5\nndf6\n" +
				"odf7\npdf8\nrdf9\ntdf10\nvdf11\n",
			cf: &ControlFile{
				Documents: []Document{
					{'c', "df1", ""},
					{'d', "df2", ""},
					{'f', "df3", ""},
					{'g', "df4", ""},
					{'l', "df5", ""},
					{'n', "df6", ""},
					{'o', "df7", ""},
					{'p', "df8", ""},
					{'r', "df9", ""},
					{'t', "df10", ""},
					{'v', "df11", ""},
				},
			},
		},

		{
			name: "no-data-file-name",
			in:   "Hhost\nl\n",
			err:  ErrNoDataFileName,
		},
	}

	for _, test := range tests {
		cf, err := ParseControlFile([]byte(test.in))

		if !errors.Is(err, test.err) {
			t.Errorf("%s: error mismatch:\n"+
				"expected: %v\n"+
				"present:  %v", test.name, test.err, err)
			continue
		}

		if diff := cmp.Diff(test.cf, cf); diff != "" {
			t.Errorf("%s: result mismatch (-expected +present):\n%s",
				test.name, diff)
		}
	}
}

// TestControlFileTitle tests ControlFile.Title
func TestControlFileTitle(t *testing.T) {
	type testData struct {
		cf    ControlFile
		title string
	}

	tests := []testData{
		{ControlFile{}, ""},
		{ControlFile{JobName: "job"}, "job"},
		{
			ControlFile{
				JobName:   "job",
				Documents: []Document{{'l', "df1", "doc"}},
			},
			"job",
		},
		{
			ControlFile{
				Documents: []Document{
					{'l', "df1", ""},
					{'l', "df2", "doc"},
				},
			},
			"doc",
		},
	}

	for _, test := range tests {
		title := test.cf.Title()
		if title != test.title {
			t.Errorf("%+v: expected %q, present %q",
				test.cf, test.title, title)
		}
	}
}

// TestControlFileDataFiles tests ControlFile.DataFiles
func TestControlFileDataFiles(t *testing.T) {
	cf := ControlFile{
		Documents: []Document{
			{'l', "dfA", ""},
			{'l', "dfA", ""},
			{'f', "dfB", ""},
			{'l', "dfA", ""},
			{'o', "dfC", ""},
		},
	}

	expected := []string{"dfA", "dfB", "dfC"}
	present := cf.DataFiles()

	if !reflect.DeepEqual(expected, present) {
		t.Errorf("expected %q, present %q", expected, present)
	}

	if names := (&ControlFile{}).DataFiles(); names != nil {
		t.Errorf("empty: expected nil, present %q", names)
	}
}

// TestDocumentIsRaw tests Document.IsRaw
func TestDocumentIsRaw(t *testing.T) {
	if !(Document{Format: 'l'}).IsRaw() {
		t.Errorf("'l' must be raw")
	}

	for _, c := range []byte("cdfgnoprtv") {
		if (Document{Format: c}).IsRaw() {
			t.Errorf("%q must not be raw", c)
		}
	}
}
