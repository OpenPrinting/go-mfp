// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// Report generation for mfp-test

package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Report is the top-level JSON report written by --output.
type Report struct {
	GeneratedAt time.Time     `json:"generated_at"`
	Total       int           `json:"total"`
	Passed      int           `json:"passed"`
	Failed      int           `json:"failed"`
	Results     []reportEntry `json:"results"`
}

// reportEntry is one row in the JSON report.
type reportEntry struct {
	Name    string             `json:"name"`
	Passed  bool               `json:"passed"`
	Score   float64            `json:"score"`
	Details map[string]float64 `json:"details,omitempty"`
}

// buildReport converts a slice of testResults into a Report.
func buildReport(results []*testResult) *Report {
	r := &Report{
		GeneratedAt: time.Now().UTC(),
		Total:       len(results),
	}
	for _, res := range results {
		if res.Passed {
			r.Passed++
		} else {
			r.Failed++
		}
		r.Results = append(r.Results, reportEntry{
			Name:    res.Config.Name,
			Passed:  res.Passed,
			Score:   res.Score,
			Details: res.Details,
		})
	}
	return r
}

// writeReport serialises r as indented JSON to path.
func writeReport(path string, r *Report) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("report: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("report: write %q: %w", path, err)
	}
	return nil
}

// printSummaryTable prints a compact pass/fail table to stdout.
func printSummaryTable(results []*testResult) {
	const colWidth = 55

	fmt.Println()
	fmt.Printf("%-*s  %8s  %s\n", colWidth, "TEST", "SCORE", "RESULT")
	fmt.Println(strings.Repeat("-", colWidth+20))

	for _, res := range results {
		status := "PASS"
		if !res.Passed {
			status = "FAIL"
		}
		fmt.Printf("%-*s  %8.4f  %s\n", colWidth, res.Config.Name, res.Score, status)
	}

	fmt.Println(strings.Repeat("-", colWidth+20))
	passed := 0
	for _, res := range results {
		if res.Passed {
			passed++
		}
	}
	fmt.Printf("%-*s  %d/%d passed\n", colWidth, "TOTAL", passed, len(results))
	fmt.Println()
}

// printVerboseDetails prints per-metric scores for a failed or verbose result.
func printVerboseDetails(res *testResult) {
	if len(res.Details) == 0 {
		return
	}
	fmt.Printf("  details for %s:\n", res.Config.Name)
	for k, v := range res.Details {
		fmt.Printf("    %-40s %.4f\n", k, v)
	}
}
