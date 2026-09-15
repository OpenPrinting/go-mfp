// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions

package test

import _ "embed"

// defaultComparatorScript is the embedded enhanced_comparison.py bundled
// into the mfp-test binary so users do not need to provide a separate file.
//
//go:embed imgeval/enhanced_comparison.py
var defaultComparatorScript []byte
