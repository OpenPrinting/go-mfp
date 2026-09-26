// MFP  - Miulti-Function Printers and scanners toolkit
// argv - Argv parsing mini-library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Value completers.

package argv

import (
	"io/fs"
	"os"
	"os/user"
	"strings"
)

// Completer is a callback called for auto-completion
//
// Any [Option] or [Parameter] may have its own Completer.
//
// It receives the Option's value prefix, already typed
// by user, and must return a slice of completion candidates
// that match the prefix.
//
// For example, if possible Option or Parameter values are "Richard",
// "Roger" and  "Robert", then, depending of supplied prefix, the following
// output is expected:
//
//	"R"   -> ["Richard", "Roger", "Robert"]
//	"Ro"  -> ["Roger", "Robert"]
//	"Rog" -> ["Roger"]
//	"Rol" -> []
type Completer func(string) []Completion

// CompleteStrings returns a [Completer], that performs auto-completion,
// choosing from a set of supplied strings.
func CompleteStrings(s []string) Completer {
	// Create a copy of input, to protect from callers
	// that may change the slice after the call.
	set := make([]string, len(s))
	copy(set, s)

	// Create completer
	return func(in string) (compl []Completion) {
		for _, member := range set {
			if len(in) < len(member) &&
				strings.HasPrefix(member, in) {
				compl = append(compl, Completion{member, false})
			}
		}
		return compl
	}
}

// CompleteOSPath is the [Completer] that completes the operating
// system file paths.
func CompleteOSPath(s string) []Completion {
	if len(s) > 0 && s[0] == '~' {
		// Need to complete user name itself?
		i := strings.IndexByte(s, '/')
		if i < 0 {
			users := localUsers(s[1:])
			compl := make([]Completion, len(users))
			for i := range users {
				compl[i].String = "~" + users[i]
				compl[i].NoSpace = true
			}

			if len(compl) == 1 {
				compl[0].String += "/"
			}

			return compl
		}

		// Resolve tilde prefix into the home directory
		var usr *user.User
		var err error
		if username := s[1:i]; username != "" {
			usr, err = user.Lookup(username)
		} else {
			usr, err = user.Current()
		}
		if err != nil {
			return nil
		}

		// Replace tilde prefix with home directory
		// and call completeOSPath
		home, _ := strings.CutSuffix(usr.HomeDir, "/")
		compl := completeOSPath(home + s[i:])

		// Restore tilde prefix on every returned path
		prefix := s[:i]
		for i := range compl {
			str, ok := strings.CutPrefix(
				compl[i].String, home)

			if ok {
				compl[i].String = prefix + str
			}
		}

		return compl
	}

	return completeOSPath(s)
}

var completeOSPath = CompleteFs(os.DirFS("/"), os.Getwd)

// CompleteFs returns a [Completer], that performs file name auto-completion
// on a top of a virtual (or real) filesystem, represented as fs.FS,
//
// getwd callback returns a current directory within that file system.
// It's signature is compatible with os.Getwd(), so this function can
// be used directly.
//
// If getwd is nil, current directory assumed to be "/"
func CompleteFs(fsys fs.FS, getwd func() (string, error)) Completer {
	fscompl := newFscompleter(fsys, getwd)
	return func(arg string) []Completion {
		return fscompl.complete(arg)
	}
}
