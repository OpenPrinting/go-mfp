// MFP - Multi-Function Printers and scanners toolkit
// The "mfp-test" command
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// CUPS queue management

package test

import (
	"context"
	"fmt"
	"os/exec"
)

// CreateCUPSQueue creates a CUPS printer queue named name,
// pointing to the virtual IPP printer at ippURL.
//
// It uses lpadmin to register the queue with the IPP Everywhere
// driver, which works with any standards-compliant IPP printer.
func CreateCUPSQueue(ctx context.Context, name, ippURL string) error {
	cmd := exec.CommandContext(ctx,
		"lpadmin",
		"-p", name,
		"-E",
		"-v", ippURL,
		"-m", "everywhere",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lpadmin -p %s: %w: %s", name, err, out)
	}

	return nil
}

// RemoveCUPSQueue removes the CUPS printer queue named name.
// It is safe to call even if the queue does not exist.
func RemoveCUPSQueue(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx,
		"lpadmin",
		"-x", name,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lpadmin -x %s: %w: %s", name, err, out)
	}

	return nil
}
