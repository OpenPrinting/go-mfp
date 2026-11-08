// MFP - Miulti-Function Printers and scanners toolkit
// The "mfp-proxy" command
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Protocol identifiers

package command

// proto identifies proxy protocol
type proto int

const (
	protoIPP proto = iota
	protoESCL
	protoWSD
)
