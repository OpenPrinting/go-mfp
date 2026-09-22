// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions

package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// TestBuildReport verifies that buildReport correctly counts passed/failed
// results and populates all fields.
func TestBuildReport(t *testing.T) {
	results := []*testResult{
		{Config: testConfig{Name: "one-sided/color/application/pdf"}, Score: 0.98, Passed: true},
		{Config: testConfig{Name: "one-sided/monochrome/image/pwg-raster"}, Score: 0.91, Passed: false},
		{Config: testConfig{Name: "two-sided-long-edge/color/application/pdf"}, Score: 0.97, Passed: true},
	}

	r := buildReport(results)

	if r.Total != 3 {
		t.Errorf("Total: want 3, got %d", r.Total)
	}
	if r.Passed != 2 {
		t.Errorf("Passed: want 2, got %d", r.Passed)
	}
	if r.Failed != 1 {
		t.Errorf("Failed: want 1, got %d", r.Failed)
	}
	if len(r.Results) != 3 {
		t.Errorf("Results length: want 3, got %d", len(r.Results))
	}
	if r.Results[1].Name != "one-sided/monochrome/image/pwg-raster" {
		t.Errorf("Results[1].Name: got %q", r.Results[1].Name)
	}
}

// TestWriteReport verifies that writeReport produces valid JSON that
// round-trips back to the same Report.
func TestWriteReport(t *testing.T) {
	results := []*testResult{
		{Config: testConfig{Name: "one-sided/color/application/pdf"}, Score: 0.98, Passed: true},
	}
	report := buildReport(results)

	f, err := os.CreateTemp("", "mfp-test-report-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	if err := writeReport(path, report); err != nil {
		t.Fatalf("writeReport: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}

	var got Report
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Total != 1 || got.Passed != 1 || got.Failed != 0 {
		t.Errorf("round-trip mismatch: %+v", got)
	}
}

// TestMimeToExt verifies the MIME → file extension mapping for known types.
func TestMimeToExt(t *testing.T) {
	cases := []struct {
		mime string
		want string
	}{
		{"application/pdf", ".pdf"},
		{"application/postscript", ".ps"},
		{"image/pwg-raster", ".pwg"},
		{"image/urf", ".urf"},
		{"image/jpeg", ".jpg"},
		{"image/png", ".png"},
		{"unknown/type", ".bin"},
	}

	for _, c := range cases {
		if got := mimeToExt(c.mime); got != c.want {
			t.Errorf("mimeToExt(%q) = %q, want %q", c.mime, got, c.want)
		}
	}
}

// TestConfigName verifies that configName produces the expected
// slash-separated string.
func TestConfigName(t *testing.T) {
	got := configName(ipp.KwSidesOneSided, "color", "application/pdf")
	want := "one-sided/color/application/pdf"
	if got != want {
		t.Errorf("configName = %q, want %q", got, want)
	}
}

// TestGenerateTestPNG verifies that generateTestPNG creates a non-empty
// PNG file that can be removed afterwards.
func TestGenerateTestPNG(t *testing.T) {
	path, err := generateTestPNG()
	if err != nil {
		t.Fatalf("generateTestPNG: %v", err)
	}
	defer os.Remove(path)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() == 0 {
		t.Error("generated PNG is empty")
	}
}
