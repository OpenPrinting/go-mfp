// MFP - Miulti-Function Printers and scanners toolkit
// IEEE 1284 definitions
//
// Copyright (C) 2024 and up by Mohammad Arman(officialmdarman@gmail.com)
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

	upper := strings.ToUpper(after)

	// Two-word commands: INFO ID, INFO STATUS, ENTER LANGUAGE,
	// USTATUS DEVICE, USTATUS JOB.
	//
	// We match against the uppercased string to handle
	// case-insensitive input and multiple spaces between words.
	twoWordCmds := []string{
		"INFO ID",
		"INFO STATUS",
		"ENTER LANGUAGE",
		"USTATUS DEVICE",
		"USTATUS JOB",
	}

	for _, cmd := range twoWordCmds {
		words := strings.Fields(cmd)
		// Check that the first word matches
		if !strings.HasPrefix(upper, words[0]) {
			continue
		}
		// Skip whitespace after the first word
		rest := upper[len(words[0]):]
		trimmed := strings.TrimLeft(rest, " \t")
		if rest == trimmed {
			// No whitespace after first word — not a match
			continue
		}
		// Check that the second word matches
		if !strings.HasPrefix(trimmed, words[1]) {
			continue
		}
		// Compute where the args start in the original string
		consumed := len(after) - len(rest) +
			(len(rest) - len(trimmed)) +
			len(words[1])
		args := strings.TrimSpace(after[consumed:])
		return pjlCommand{
			name: cmd,
			args: args,
		}
	}

	// Single-word command: ECHO, SET, JOB, EOJ, USTATUSOFF
	fields := strings.Fields(after)
	word := fields[0]
	rest := strings.TrimSpace(after[len(word):])
	return pjlCommand{
		name: strings.ToUpper(word),
		args: rest,
	}
}

// handlePJL dispatches a parsed PJL command, generating responses
// where required.
//
// Returns true if the command triggers a state transition
// (i.e. ENTER LANGUAGE), false otherwise.
func (p *Printer) handlePJL(cmd pjlCommand) bool {
	switch cmd.name {
	case "ENTER LANGUAGE":
		// Strip leading "=" from args (e.g. "=POSTSCRIPT")
		lang := strings.TrimPrefix(cmd.args, "=")
		lang = strings.TrimSpace(lang)
		log.Debug(p.ctx, "PJL ENTER LANGUAGE=%s", lang)
		p.format = detectFormatByLanguage(lang)
		p.state = stateDocument
		p.docBuf = nil
		return true

	case "INFO ID":
		log.Debug(p.ctx, "PJL INFO ID")
		p.queueResponse(pjlInfoID(p.model))

	case "INFO STATUS":
		log.Debug(p.ctx, "PJL INFO STATUS")
		p.queueResponse(pjlInfoStatus())

	case "ECHO":
		log.Debug(p.ctx, "PJL ECHO %s", cmd.args)
		p.queueResponse(pjlEcho(cmd.args))

	case "JOB":
		log.Debug(p.ctx, "PJL JOB %s", cmd.args)
		p.params = JobParams{
			Variables: make(map[string]string),
		}
		p.params.JobName = parsePJLQuotedValue(cmd.args, "NAME")

	case "EOJ":
		log.Debug(p.ctx, "PJL EOJ %s", cmd.args)

	case "SET":
		log.Debug(p.ctx, "PJL SET %s", cmd.args)
		key, val := parsePJLKeyValue(cmd.args)
		if key != "" {
			if p.params.Variables == nil {
				p.params.Variables = make(map[string]string)
			}
			p.params.Variables[key] = val
		}

	case "USTATUS DEVICE", "USTATUS JOB":
		log.Debug(p.ctx, "PJL %s %s",
			cmd.name, cmd.args)

	case "USTATUSOFF":
		log.Debug(p.ctx, "PJL USTATUSOFF")

	case "":
		// Bare @PJL — initialization marker, ignore

	default:
		log.Debug(p.ctx, "PJL %s %s (unhandled)",
			cmd.name, cmd.args)
	}

	return false
}

// pjlInfoID formats a PJL INFO ID response.
func pjlInfoID(model string) []byte {
	return []byte(fmt.Sprintf("@PJL INFO ID\r\n%q\r\n\x0c", model))
}

// pjlInfoStatus formats a PJL INFO STATUS response.
// CODE=10001 means "Ready".
func pjlInfoStatus() []byte {
	return []byte("@PJL INFO STATUS\r\n" +
		"CODE=10001\r\n" +
		"DISPLAY=\"READY\"\r\n" +
		"ONLINE=TRUE\r\n" +
		"\x0c")
}

// pjlEcho formats a PJL ECHO response.
// ECHO responses do not have a form feed terminator.
func pjlEcho(text string) []byte {
	return []byte(fmt.Sprintf("@PJL ECHO %s\r\n", text))
}

// parsePJLKeyValue parses "KEY=VALUE" or "KEY = VALUE" from
// a PJL SET argument string. Returns uppercase key and trimmed value.
func parsePJLKeyValue(args string) (string, string) {
	eq := strings.IndexByte(args, '=')
	if eq < 0 {
		return "", ""
	}
	key := strings.TrimSpace(args[:eq])
	val := strings.TrimSpace(args[eq+1:])
	// Strip surrounding quotes if present
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
	}
	return strings.ToUpper(key), val
}

// parsePJLQuotedValue extracts the value for a named parameter
// from a PJL argument string. For example, given
// `NAME = "test" DISPLAY = "hello"` and key "NAME",
// it returns "test".
func parsePJLQuotedValue(args, key string) string {
	upper := strings.ToUpper(args)
	key = strings.ToUpper(key)
	idx := strings.Index(upper, key)
	if idx < 0 {
		return ""
	}
	rest := args[idx+len(key):]
	rest = strings.TrimSpace(rest)
	if len(rest) == 0 || rest[0] != '=' {
		return ""
	}
	rest = strings.TrimSpace(rest[1:])
	if len(rest) == 0 {
		return ""
	}
	if rest[0] == '"' {
		end := strings.IndexByte(rest[1:], '"')
		if end < 0 {
			return rest[1:]
		}
		return rest[1 : end+1]
	}
	// Unquoted value: take until whitespace
	fields := strings.Fields(rest)
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}
