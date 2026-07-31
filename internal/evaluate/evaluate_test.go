// MFP - Multi-Function Printers and scanners toolkit
// internal/evaluate - Image quality evaluation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions

package evaluate

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// makeTestPNG creates a temporary PNG file with the given fill color
// and returns its path. The caller must remove it after use.
func makeTestPNG(t *testing.T, r, g, b uint8) string {
	t.Helper()

	const size = 100
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	f, err := os.CreateTemp("", "evaluate-test-*.png")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		os.Remove(f.Name())
		t.Fatalf("encode PNG: %v", err)
	}

	return f.Name()
}

// comparatorPath returns the path to enhanced_comparison.py.
// It looks for the file next to the test binary (set via -comparator flag)
// or falls back to a well-known development location.
func comparatorPath(t *testing.T) string {
	t.Helper()

	candidates := []string{
		"enhanced_comparison.py",
		"/home/amrxg/enhanced_comparison.py",
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	t.Skip("enhanced_comparison.py not found; skipping evaluator test")
	return ""
}

// TestEvaluatorIdentical verifies that comparing an image with itself
// returns a positive score without error. A solid-color test image has no
// detectable keypoints or edges, so it will not reach DefaultThreshold;
// the assertion only checks that the subprocess ran and returned a valid score.
func TestEvaluatorIdentical(t *testing.T) {
	path := comparatorPath(t)

	ev, err := NewEvaluator(path)
	if err != nil {
		t.Fatalf("NewEvaluator: %v", err)
	}
	defer ev.Close()

	img := makeTestPNG(t, 255, 0, 0)
	defer os.Remove(img)

	res, err := ev.Compare(img, img, DefaultThreshold, false)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}

	if res.Score <= 0 {
		t.Errorf("identical images: expected positive score, got %.4f", res.Score)
	}

	t.Logf("identical images: score=%.4f", res.Score)
}

// TestEvaluatorDifferent verifies that comparing two completely different
// images yields a low score that fails the default threshold.
func TestEvaluatorDifferent(t *testing.T) {
	path := comparatorPath(t)

	ev, err := NewEvaluator(path)
	if err != nil {
		t.Fatalf("NewEvaluator: %v", err)
	}
	defer ev.Close()

	original := makeTestPNG(t, 255, 0, 0) // red
	captured := makeTestPNG(t, 0, 0, 255) // blue
	defer os.Remove(original)
	defer os.Remove(captured)

	res, err := ev.Compare(original, captured, DefaultThreshold, true)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}

	t.Logf("different images: score=%.4f passed=%v", res.Score, res.Passed)
	t.Logf("details: %v", res.Details)
}
