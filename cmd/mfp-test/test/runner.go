// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Test runner for mfp-test

package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/OpenPrinting/go-mfp/log"
)

// defaultThreshold is the minimum similarity score required to pass a test.
const defaultThreshold = 0.95

// defaultTimeout is the default time to wait for a captured document.
const defaultTimeout = 30 * time.Second

// testResult holds the outcome of a single test run.
type testResult struct {
	Config  testConfig
	Score   float64
	Passed  bool
	Details map[string]float64
}

// runTest runs a single print test using the given configuration:
// generates a test PNG, sends it to the CUPS queue with the specified
// job attributes, waits for capture, and returns the test result.
//
// Image evaluation is not yet implemented; the function currently
// reports success if the document was captured within the timeout.
func runTest(ctx context.Context, cfg testConfig, queueName string,
	capture *documentCapture, threshold float64, timeout time.Duration, keep, verbose bool) (*testResult, error) {

	// Reset capture so we get only this job's document.
	capture.reset()

	// Generate a fresh test image for this run.
	imgPath, err := generateTestPNG()
	if err != nil {
		return nil, fmt.Errorf("generate test image: %w", err)
	}
	defer os.Remove(imgPath)

	// Build lp options for the test configuration.
	lpArgs := []string{"-d", queueName, "-t", "mfp-test/" + cfg.Name}
	if cfg.Sides != "" {
		lpArgs = append(lpArgs, "-o", "sides="+string(cfg.Sides))
	}
	if cfg.ColorMode != "" {
		lpArgs = append(lpArgs, "-o", "print-color-mode="+cfg.ColorMode)
	}
	lpArgs = append(lpArgs, imgPath)

	log.Info(ctx, "lp %v", lpArgs)
	lpCmd := exec.CommandContext(ctx, "lp", lpArgs...)
	if out, err := lpCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("lp: %w: %s", err, out)
	}

	// Wait for the document to arrive.
	select {
	case <-capture.onDocument():
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout: no document received after %s", timeout)
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	captured := capture.docs()
	if len(captured) == 0 {
		return nil, fmt.Errorf("capture returned no documents")
	}

	d := captured[len(captured)-1]
	if verbose {
		log.Info(ctx, "captured %d bytes format=%q job=%q",
			len(d.Data), d.Params.Format, d.Params.JobName)
	}

	if keep {
		safeName := strings.ReplaceAll(cfg.Name, "/", "-")
		outPath := safeName + mimeToExt(d.Params.Format)
		if err := os.WriteFile(outPath, d.Data, 0644); err != nil {
			log.Info(ctx, "keep: failed to save %s: %v", outPath, err)
		} else {
			log.Info(ctx, "keep: saved captured document to %s", outPath)
		}
	}

	// Image evaluation will be wired here in Phase 5 once raster
	// conversion (captured bytes → PNG) is implemented. For now,
	// a successful capture counts as a pass with a placeholder score.
	score := 1.0
	return &testResult{
		Config: cfg,
		Score:  score,
		Passed: score >= threshold,
	}, nil
}

// mimeToExt returns a file extension for a MIME type.
// Gzip-compressed variants get a double extension (e.g. .jpg.gz).
func mimeToExt(mime string) string {
	switch mime {
	case "application/pdf":
		return ".pdf"
	case "application/postscript":
		return ".ps"
	case "image/jpeg":
		return ".jpg"
	case "image/jpeg+gzip":
		return ".jpg.gz"
	case "image/png":
		return ".png"
	case "image/pwg-raster":
		return ".pwg"
	case "image/urf":
		return ".urf"
	case "image/vnd.cups-raster":
		return ".ras"
	case "application/vnd.cups-pdf":
		return ".cups.pdf"
	case "application/vnd.cups-postscript":
		return ".cups.ps"
	default:
		return ".bin"
	}
}
