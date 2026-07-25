// MFP - Multi-Function Printers and scanners toolkit
// internal/evaluate - Image quality evaluation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions

package evaluate

import (
	"fmt"
	"os"

	"github.com/OpenPrinting/go-mfp/cpython"
)

// DefaultThreshold is the minimum overall quality score to pass.
const DefaultThreshold = 0.95

// Result holds the outcome of a single image comparison.
type Result struct {
	// Score is the overall quality score returned by ImageComparator
	// (0.0 = worst, 1.0 = perfect match).
	Score float64

	// Passed reports whether Score >= the requested threshold.
	Passed bool

	// Details contains per-metric scores. It is only populated when
	// verbose is requested or when the comparison fails (Score < threshold),
	// to help diagnose which metric caused the failure.
	Details map[string]float64
}

// Evaluator wraps the ImageComparator Python class and provides
// a simple Go API for image quality comparison.
//
// Create with [NewEvaluator] and release with [Evaluator.Close].
type Evaluator struct {
	py  *cpython.Python
	cls *cpython.Object // ImageComparator class
}

// NewEvaluator creates a new Evaluator by loading the ImageComparator
// class from the given Python file path.
//
// comparatorPath must point to enhanced_comparison.py from the
// OpenPrinting-Image-Evaluation project.
func NewEvaluator(comparatorPath string) (*Evaluator, error) {
	// Read the Python source file
	data, err := os.ReadFile(comparatorPath)
	if err != nil {
		return nil, fmt.Errorf("evaluate: read %q: %w", comparatorPath, err)
	}

	// Create a new Python interpreter
	py, err := cpython.NewPython()
	if err != nil {
		return nil, fmt.Errorf("evaluate: init Python: %w", err)
	}

	// Clean up interpreter on failure
	defer func() {
		if err != nil {
			py.Close()
		}
	}()

	// Load enhanced_comparison.py as a Python module
	mod := py.Load(string(data), "enhanced_comparison", comparatorPath)
	if err = mod.Err(); err != nil {
		return nil, fmt.Errorf("evaluate: load %q: %w", comparatorPath, err)
	}

	// Get the ImageComparator class from the module
	cls := mod.Get("ImageComparator")
	if err = cls.Err(); err != nil {
		return nil, fmt.Errorf("evaluate: ImageComparator class not found: %w", err)
	}

	return &Evaluator{py: py, cls: cls}, nil
}

// Close releases the Python interpreter and all resources held
// by the Evaluator.
func (e *Evaluator) Close() {
	e.py.Close()
}

// Compare compares the captured image against the original and
// returns a [Result].
//
//   - original: path to the reference PNG image
//   - captured: path to the captured/printed PNG image
//   - threshold: minimum score to pass (use [DefaultThreshold] if unsure)
//   - verbose: if true, always populate Result.Details;
//     if false, Details is only populated when the test fails
func (e *Evaluator) Compare(original, captured string,
	threshold float64, verbose bool) (*Result, error) {

	// Instantiate: comp = ImageComparator(original, captured)
	instance := e.cls.Call(original, captured)
	if err := instance.Err(); err != nil {
		return nil, fmt.Errorf("evaluate: ImageComparator(%q, %q): %w",
			original, captured, err)
	}

	// Call: results = comp.run_all_comparisons()
	method := instance.Get("run_all_comparisons")
	if err := method.Err(); err != nil {
		return nil, fmt.Errorf("evaluate: run_all_comparisons not found: %w", err)
	}

	results := method.Call()
	if err := results.Err(); err != nil {
		return nil, fmt.Errorf("evaluate: run_all_comparisons(): %w", err)
	}

	// Extract overall_quality score
	scoreObj := results.GetItem("overall_quality")
	if err := scoreObj.Err(); err != nil {
		return nil, fmt.Errorf("evaluate: overall_quality not found: %w", err)
	}

	score, err := scoreObj.Float()
	if err != nil {
		return nil, fmt.Errorf("evaluate: overall_quality to float: %w", err)
	}

	passed := score >= threshold

	res := &Result{
		Score:  score,
		Passed: passed,
	}

	// Populate per-metric details when verbose or on failure
	if verbose || !passed {
		details, err := extractDetails(results)
		if err != nil {
			return nil, fmt.Errorf("evaluate: extract metrics: %w", err)
		}
		res.Details = details
	}

	return res, nil
}

// extractDetails reads all key/value pairs from the Python results
// dict and returns them as a Go map[string]float64.
func extractDetails(results *cpython.Object) (map[string]float64, error) {
	keys, err := results.Keys()
	if err != nil {
		return nil, err
	}

	details := make(map[string]float64, len(keys))
	for _, k := range keys {
		name, err := k.Unicode()
		if err != nil {
			continue
		}

		val := results.GetItem(name)
		if val.Err() != nil {
			continue
		}

		f, err := val.Float()
		if err != nil {
			continue
		}

		details[name] = f
	}

	return details, nil
}
