// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Test matrix generation for mfp-test

package test

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/OpenPrinting/go-mfp/proto/ipp"
)

// printerCaps holds the queried printer capabilities used to generate
// the test matrix.
type printerCaps struct {
	Sides      []ipp.KwSides
	ColorModes []string
	Formats    []string
}

// testConfig represents one specific combination of print parameters
// to exercise in a test run.
type testConfig struct {
	Name      string
	Sides     ipp.KwSides
	ColorMode string
	Format    string
}

// queryPrinterCaps queries the virtual IPP printer for its supported
// attribute values and returns them as a printerCaps.
func queryPrinterCaps(ctx context.Context, printerURL string) (*printerCaps, error) {
	u, err := url.Parse(printerURL)
	if err != nil {
		return nil, fmt.Errorf("matrix: parse printer URL: %w", err)
	}

	client := ipp.NewClient(u, nil)
	attrs, err := client.GetPrinterAttributes(ctx,
		[]string{"job-template", "printer-description"}, "")
	if err != nil {
		return nil, fmt.Errorf("matrix: query printer attributes: %w", err)
	}

	caps := &printerCaps{
		Sides:      attrs.SidesSupported,
		ColorModes: attrs.PrintColorModeSupported,
		Formats:    attrs.DocumentFormatSupported,
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

// configName builds a deterministic, human-readable name for a test
// configuration from its three dimensions.
func configName(sides ipp.KwSides, color, format string) string {
	return fmt.Sprintf("%s/%s/%s", sides, color, format)
}

// batchMatrix returns every combination of sides × color mode × format.
// This is the exhaustive test matrix.
func batchMatrix(caps *printerCaps) []testConfig {
	var configs []testConfig
	for _, sides := range caps.Sides {
		for _, color := range caps.ColorModes {
			for _, format := range caps.Formats {
				configs = append(configs, testConfig{
					Name:      configName(sides, color, format),
					Sides:     sides,
					ColorMode: color,
					Format:    format,
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
	var configs []testConfig
	for _, sides := range caps.Sides {
		for _, color := range caps.ColorModes {
			configs = append(configs, testConfig{
				Name:      configName(sides, color, format),
				Sides:     sides,
				ColorMode: color,
				Format:    format,
			})
		}
	}
	return configs
}

// singleConfig parses a configuration name of the form
// "sides/color-mode/format" and returns the corresponding testConfig.
// This is used with --single to reproduce a specific known bug.
func singleConfig(spec string) (*testConfig, error) {
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
	}, nil
}
