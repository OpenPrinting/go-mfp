// MFP - Multi-Function Printers and scanners toolkit
// internal/evaluate - Image quality evaluation
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions

package evaluate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// DefaultThreshold is the minimum overall quality score to pass.
const DefaultThreshold = 0.95

// runner is the Python script executed as a subprocess. It loads
// ImageComparator from enhanced_comparison.py, runs the comparison,
// and prints JSON results to stdout.
//
// argv: runner.py <comparator-path> <original-path> <captured-path>
//
// Stdout is redirected to stderr during import and comparison so that
// any library warnings or progress messages do not corrupt the JSON output.
const runner = `
import sys, json, os, math

_stdout = sys.stdout
sys.stdout = sys.stderr

sys.path.insert(0, os.path.dirname(os.path.abspath(sys.argv[1])))
from enhanced_comparison import ImageComparator
comp = ImageComparator(sys.argv[2], sys.argv[3])
results = comp.run_all_comparisons()

sys.stdout = _stdout
out = {}
for k, v in results.items():
    try:
        f = float(v)
        if not (math.isinf(f) or math.isnan(f)):
            out[k] = f
    except (TypeError, ValueError):
        pass
print(json.dumps(out))
`

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

// Evaluator runs ImageComparator in a subprocess and provides
// a simple Go API for image quality comparison.
//
// Create with [NewEvaluator] and release with [Evaluator.Close].
type Evaluator struct {
	comparatorPath string // path to enhanced_comparison.py
	runnerPath     string // temp file containing the runner script
}

// NewEvaluator creates a new Evaluator for the given enhanced_comparison.py
// path. A small Python runner script is written to a temp file; it is removed
// when [Evaluator.Close] is called.
func NewEvaluator(comparatorPath string) (*Evaluator, error) {
	if _, err := os.Stat(comparatorPath); err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}

	f, err := os.CreateTemp("", "mfp-evaluate-*.py")
	if err != nil {
		return nil, fmt.Errorf("evaluate: create runner: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(runner); err != nil {
		os.Remove(f.Name())
		return nil, fmt.Errorf("evaluate: write runner: %w", err)
	}

	return &Evaluator{
		comparatorPath: comparatorPath,
		runnerPath:     f.Name(),
	}, nil
}

// Close removes the temporary runner script created by [NewEvaluator].
func (e *Evaluator) Close() {
	os.Remove(e.runnerPath)
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

	cmd := exec.Command("python3", e.runnerPath, e.comparatorPath, original, captured)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("evaluate: python3: %w: %s", err, ee.Stderr)
		}
		return nil, fmt.Errorf("evaluate: python3: %w", err)
	}

	var scores map[string]float64
	if err := json.Unmarshal(out, &scores); err != nil {
		return nil, fmt.Errorf("evaluate: parse results: %w", err)
	}

	score, ok := scores["overall_quality"]
	if !ok {
		return nil, fmt.Errorf("evaluate: overall_quality not found in results")
	}

	res := &Result{
		Score:  score,
		Passed: score >= threshold,
	}

	if verbose || !res.Passed {
		res.Details = scores
	}

	return res, nil
}
