// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// PJL command handling

package ieee1284

import (
	"fmt"
	"strings"

	"github.com/OpenPrinting/go-mfp/log"
)

// pjlCommand represents a parsed PJL command.
type pjlCommand struct {
	name string // e.g. "INFO ID", "ECHO", "SET"
	args string // everything after the command name
}

// parsePJLCommand parses a PJL line like "@PJL INFO ID" into
// a command name and args.
//
// The line must already be trimmed of leading/trailing whitespace
// and have the "@PJL" prefix stripped.
func parsePJLCommand(after string) pjlCommand {
	after = strings.TrimSpace(after)
	if after == "" {
		return pjlCommand{}
	}

	// Two-word commands: INFO ID, INFO STATUS, ENTER LANGUAGE,
	// USTATUS DEVICE, USTATUS JOB
	fields := strings.Fields(after)
	if len(fields) >= 2 {
		twoWord := strings.ToUpper(fields[0] + " " + fields[1])
		switch {
		case twoWord == "INFO ID" ||
			twoWord == "INFO STATUS" ||
			twoWord == "ENTER LANGUAGE" ||
			twoWord == "USTATUS DEVICE" ||
			twoWord == "USTATUS JOB":
			rest := strings.TrimSpace(
				after[len(fields[0])+1+len(fields[1]):])
			return pjlCommand{
				name: twoWord,
				args: rest,
			}
		}
	}

	// Single-word command: ECHO, SET, JOB, EOJ, USTATUSOFF
	word := fields[0]
	rest := strings.TrimSpace(after[len(word):])
	return pjlCommand{
		name: strings.ToUpper(word),
		args: rest,
	}
}

// handlePJL dispatches a parsed PJL command, generating responses
// where required.
func (p *Printer) handlePJL(cmd pjlCommand) {
	switch cmd.name {
	case "INFO ID":
		log.Debug(p.ctx, "ieee1284: PJL INFO ID")
		p.queueResponse(pjlInfoID(p.model))

	case "INFO STATUS":
		log.Debug(p.ctx, "ieee1284: PJL INFO STATUS")
		p.queueResponse(pjlInfoStatus())

	case "ECHO":
		log.Debug(p.ctx, "ieee1284: PJL ECHO %s", cmd.args)
		p.queueResponse(pjlEcho(cmd.args))

	case "JOB":
		log.Debug(p.ctx, "ieee1284: PJL JOB %s", cmd.args)

	case "EOJ":
		log.Debug(p.ctx, "ieee1284: PJL EOJ %s", cmd.args)

	case "SET":
		log.Debug(p.ctx, "ieee1284: PJL SET %s", cmd.args)

	case "USTATUS DEVICE", "USTATUS JOB":
		log.Debug(p.ctx, "ieee1284: PJL %s %s",
			cmd.name, cmd.args)

	case "USTATUSOFF":
		log.Debug(p.ctx, "ieee1284: PJL USTATUSOFF")

	case "":
		// Bare @PJL — initialization marker, ignore

	default:
		log.Debug(p.ctx, "ieee1284: PJL %s %s (unhandled)",
			cmd.name, cmd.args)
	}
}

// pjlInfoID formats a PJL INFO ID response.
func pjlInfoID(model string) []byte {
	return []byte(fmt.Sprintf("@PJL INFO ID\r\n%q\r\n\x0c", model))
}

// pjlInfoStatus formats a PJL INFO STATUS response.
func pjlInfoStatus() []byte {
	return []byte("@PJL INFO STATUS\r\n" +
		"CODE=10001\r\n" +
		"DISPLAY=\"READY\"\r\n" +
		"ONLINE=TRUE\r\n" +
		"\x0c")
}

// pjlEcho formats a PJL ECHO response.
// Per PJL spec, ECHO responses do not have a form feed terminator.
func pjlEcho(text string) []byte {
	return []byte(fmt.Sprintf("@PJL ECHO %s\r\n", text))
}
