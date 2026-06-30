// MFP - Miulti-Function Printers and scanners toolkit
// Utility functions and data BLOBs for testing
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

package testutils

import (
	diff "github.com/thepudds/patience-diff"
)

// DiffLines returns diff of the two texts old and new in the “unified diff” format.
// If texts are identical, it returns "" (no output)
func DiffLines(oldName string, old string, newName string, new string) string {
	return string(DiffLineBytes(oldName, []byte(old), newName, []byte(new)))
}

// DiffLineBytes is like [DiffLines], but operates with byte slices instead of
// strings.
func DiffLineBytes(oldName string, old []byte, newName string, new []byte) []byte {
	return diff.Diff(oldName, old, newName, new)
}
