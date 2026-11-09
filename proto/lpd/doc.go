// MFP - Multi-Function Printers and scanners toolkit
// LPD protocol implementation
//
// Copyright (C) 2026 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// Package documentation

// Package lpd implements the server side of the Line Printer Daemon
// protocol, as defined by RFC 1179.
//
// Supported subset of the protocol:
//
//   - \x01 queue LF: print any waiting jobs (acknowledged, no-op)
//   - \x02 queue LF: receive a print job, with subcommands:
//     \x01 (abort job), \x02 (control file), \x03 (data file)
//   - \x03 and \x04: send short/long queue state
//   - \x05: remove jobs (acknowledged, no-op)
package lpd
