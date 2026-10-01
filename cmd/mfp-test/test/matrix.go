// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Test matrix generation for mfp-test

package test

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/proto/ipp"
	"github.com/OpenPrinting/go-mfp/util/optional"
	goipp "github.com/OpenPrinting/goipp"
)

// printerCaps holds the queried printer capabilities used to generate
// the test matrix.
type printerCaps struct {
	Sides       []ipp.KwSides
	ColorModes  []string
	Formats     []string
	Resolutions []goipp.Resolution // printer-resolution-supported
	ModelName   string             // printer-make-and-model
}

// testConfig represents one specific combination of print parameters
// to exercise in a test run.
type testConfig struct {
	Name      string
	Sides     ipp.KwSides
	ColorMode string
	Format    string
	DPI       int  // negotiated print resolution (X DPI); 0 means use default
	Mono      bool // true when the effective color mode is monochrome
}

// capsFromModel builds a printerCaps directly from the model's IPP printer
// attributes, avoiding an HTTP round-trip to the virtual printer.
func capsFromModel(model *modeling.Model) (*printerCaps, error) {
	attrs := model.GetIPPPrinterAttrs()
	if attrs == nil {
		return nil, fmt.Errorf("matrix: model has no IPP printer attributes")
	}

	caps := &printerCaps{
		Sides:       attrs.SidesSupported,
		ColorModes:  attrs.PrintColorModeSupported,
		Formats:     attrs.DocumentFormatSupported,
		Resolutions: attrs.PrinterResolutionSupported,
		ModelName:   optional.Get(attrs.PrinterMakeAndModel),
	}

	if len(caps.Sides) == 0 {
		caps.Sides = []ipp.KwSides{ipp.KwSidesOneSided}
	}
	if len(caps.ColorModes) == 0 {
		caps.ColorModes = []string{"color"}
	}
	if len(caps.Formats) == 0 {
		caps.Formats = []string{"application/octet-stream"}
	}

	return caps, nil
}

// safeModelName converts a printer-make-and-model string into a safe
// CUPS queue name segment by lowercasing and replacing non-alphanumeric
// characters with hyphens.
var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func safeModelName(model string) string {
	s := strings.ToLower(model)
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "printer"
	}
	return s
}

// configName builds a deterministic, human-readable name for a test
// configuration from its three dimensions.
func configName(sides ipp.KwSides, color, format string) string {
	return fmt.Sprintf("%s/%s/%s", sides, color, format)
}

// isMonoMode reports whether the effective print color mode is monochrome.
// It returns true when colorMode is "monochrome" or "auto-monochrome", or
// when colorMode is "auto" and the printer only supports monochrome modes
// (no "color" or "highlight-color" in its capability list).
func isMonoMode(caps *printerCaps, colorMode string) bool {
	if colorMode == "monochrome" || colorMode == "auto-monochrome" {
		return true
	}
	if colorMode == "auto" {
		for _, m := range caps.ColorModes {
			if m == "color" || m == "highlight-color" {
				return false
			}
		}
		return true
	}
	return false
}

// defaultDPI returns the first X resolution from the printer's supported
// resolutions, or 300 if none are advertised.
func defaultDPI(caps *printerCaps) int {
	if len(caps.Resolutions) > 0 {
		return caps.Resolutions[0].Xres
	}
	return 300
}

// batchMatrix returns every combination of sides × color mode × format.
// This is the exhaustive test matrix.
func batchMatrix(caps *printerCaps) []testConfig {
	dpi := defaultDPI(caps)
	var configs []testConfig
	for _, sides := range caps.Sides {
		for _, color := range caps.ColorModes {
			for _, format := range caps.Formats {
				configs = append(configs, testConfig{
					Name:      configName(sides, color, format),
					Sides:     sides,
					ColorMode: color,
					Format:    format,
					DPI:       dpi,
					Mono:      isMonoMode(caps, color),
				})
			}
		}
	}
	return configs
}

// quickMatrix returns a reduced matrix: all sides × all color modes,
// but only the first document format. Duplex/simplex and color/mono
// are tested independently; format variation is omitted to keep the
// run short.
func quickMatrix(caps *printerCaps) []testConfig {
	format := caps.Formats[0]
	dpi := defaultDPI(caps)
	var configs []testConfig
	for _, sides := range caps.Sides {
		for _, color := range caps.ColorModes {
			configs = append(configs, testConfig{
				Name:      configName(sides, color, format),
				Sides:     sides,
				ColorMode: color,
				Format:    format,
				DPI:       dpi,
				Mono:      isMonoMode(caps, color),
			})
		}
	}
	return configs
}

// singleConfig parses a configuration name of the form
// "sides/color-mode/format" and returns the corresponding testConfig.
// This is used with --single to reproduce a specific known bug.
func singleConfig(spec string, caps *printerCaps) (*testConfig, error) {
	parts := strings.SplitN(spec, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("matrix: --single requires sides/color-mode/format, got %q", spec)
	}
	sides := ipp.KwSides(parts[0])
	switch sides {
	case ipp.KwSidesOneSided,
		ipp.KwSidesTwoSidedLongEdge,
		ipp.KwSidesTwoSidedShortEdge:
	default:
		return nil, fmt.Errorf("matrix: unknown sides value %q", sides)
	}
	return &testConfig{
		Name:      spec,
		Sides:     sides,
		ColorMode: parts[1],
		Format:    parts[2],
		DPI:       defaultDPI(caps),
		Mono:      isMonoMode(caps, parts[1]),
	}, nil
}
