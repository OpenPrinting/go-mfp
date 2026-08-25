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
	"time"

	"github.com/OpenPrinting/go-mfp/log"
)

// DefaultThreshold is the minimum similarity score required to pass a test.
const DefaultThreshold = 0.95

// TestResult holds the outcome of a single test run.
type TestResult struct {
	Config  TestConfig
	Score   float64
	Passed  bool
	Details map[string]float64
}

// RunTest runs a single print test using the given configuration:
// generates a test PNG, sends it to the CUPS queue with the specified
// job attributes, waits for capture, and returns the test result.
//
// Image evaluation is not yet implemented; the function currently
// reports success if the document was captured within the timeout.
func RunTest(ctx context.Context, cfg TestConfig, queueName string,
	capture *DocumentCapture, threshold float64, verbose bool) (*TestResult, error) {

	// Reset capture so we get only this job's document.
	capture.Reset()

	// Generate a fresh test image for this run.
	imgPath, err := generateTestPNG()
	if err != nil {
		return nil, fmt.Errorf("generate test image: %w", err)
	}
	defer os.Remove(imgPath)

	// Build lp options for the test configuration.
	lpArgs := []string{"-d", queueName}
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
	case <-capture.OnDocument():
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout: no document received after 30s")
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	docs := capture.Docs()
	if len(docs) == 0 {
		return nil, fmt.Errorf("capture returned no documents")
	}

	d := docs[len(docs)-1]
	if verbose {
		log.Info(ctx, "captured %d bytes format=%q job=%q",
			len(d.Data), d.Params.Format, d.Params.JobName)
	}

	// Image evaluation will be wired here in Phase 5 once raster
	// conversion (captured bytes → PNG) is implemented. For now,
	// a successful capture counts as a pass with a placeholder score.
	score := 1.0
	return &TestResult{
		Config: cfg,
		Score:  score,
		Passed: score >= threshold,
	}, nil
}
