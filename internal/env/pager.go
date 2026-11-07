// MFP - Miulti-Function Printers and scanners toolkit
// Execution environment
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Pager -- scroll text in a terminal

package env

import (
	"fmt"
	"os"
	"os/exec"
)

// Parameters for external Pager program:
var (
	DefaultPagerCommand     = "less"
	DefaultPagerEnvironment = []string{"LESS=FRX", "LV=-c"}
)

// Pager scrolls potentially large text in a terminal.
//
// Implements [io.WriteCloser] interface.
type Pager struct {
	pipe   *os.File  // Pipe to pager process
	stdout *os.File  // Saved os.Stdout
	stderr *os.File  // Saved os.Stderr
	cmd    *exec.Cmd // Pager command
}

// NewPager creates a new [Pager].
//
// Use [Pager.Redirect] to redirect process's stdout and stderr
// to the Pager.
//
// Caller MUST call [Pager.Close] before the process exit.
func NewPager() (pager *Pager) {
	pager = &Pager{}

	// Create pager's pipe stream
	rpipe, wpipe, err := os.Pipe()
	if err != nil {
		return
	}

	defer func() {
		if rpipe != nil {
			rpipe.Close()
		}
		if wpipe != nil {
			wpipe.Close()
		}
	}()

	// Prepare pager command
	command := os.Getenv("PAGER")
	environment := []string{}
	if command == "" {
		command = DefaultPagerCommand
		environment = DefaultPagerEnvironment
	}

	cmd := exec.Command(command)
	cmd.Stdin = rpipe
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), environment...)

	// Run the command
	err = cmd.Start()
	if err != nil {
		return
	}

	// Update pager.
	pager.pipe = wpipe
	pager.cmd = cmd

	// Cleanup pipe
	rpipe.Close()
	rpipe, wpipe = nil, nil

	return pager
}

// Redirect redirect the calling process's stdin and stdout
// into the pager.
func (p *Pager) Redirect() {
	if p.pipe != nil && os.Stdout != p.pipe {
		p.stdout, p.stderr = p.pipe, p.pipe
		os.Stdout, os.Stderr = p.pipe, p.pipe
	}
}

// Write collects text for display. All subsequent writes
// are concatenated and displayed as a whole.
//
// Implements [io.Writer] interface.
func (p *Pager) Write(text []byte) (n int, err error) {
	if p.pipe != nil {
		return p.pipe.Write(text)
	}
	return os.Stdout.Write(text)
}

// Printf writes formatted line into the [Pager], adding newline
// character after end of the line.
//
// It returns the number of bytes written and any write error encountered.
func (p *Pager) Printf(format string, args ...any) (n int, err error) {
	return fmt.Fprintf(p, format+"\n", args...)
}

// Close closes the [Pager].
func (p *Pager) Close() error {
	var err error
	if p.pipe != nil {
		err = p.pipe.Close()
		os.Stdout, os.Stderr = p.stdout, p.stderr
		p.cmd.Wait()
		p.pipe = nil
	}
	return err
}
