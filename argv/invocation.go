// MFP  - Miulti-Function Printers and scanners toolkit
// argv - Argv parsing mini-library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Invocation -- contains parsed Command's arguments.

package argv

import "context"

// Invocation represents a particular [Command] invocation.
//
// It contains a whole [Command] execution context, like parsed
// options and arguments.
//
// If [Command] contains sub-commands, and sub-command is invoked,
// the following hierarchy of Invocations is constructed:
//
//   - root Invocation, that represents top-level Command invocation
//     with its own parameters
//   - child Invocation, that represents invocation of the sub-command
//   - and so-on, if sub-command has its own nested subcommands
//
// The [Invocation.Parent] method allows to access parent of current
// Invocation. For the top-level Invocation it returns nil.
type Invocation struct {
	// parent is the upper-level Invocation for sub-command
	// Invocation, nil for the root Invocation.
	//
	// root is the top-level Invocation.
	parent, root *Invocation

	// cmd contains back reference to invoked Command.
	cmd *Command

	// argv contains original argv being used when Command
	// was invoked.
	argv []string

	// byName contains options and parameters values,
	// indexed by name.
	//
	// Because options and parameters names are incompatible,
	// we can use a single structure for both.
	byName map[string][]string

	// parameters contains parameters values, indexed by numbers.
	parameters []string

	// subcmd is the Command's SubCommand and subargv is its arguments.
	subcmd  *Command
	subargv []string

	// immediate is the first Option's Immediate callback, if any
	immediate func(context.Context, *Invocation) error
}

// Parent returns Invocation's parent, which is the upper-level
// Invocation in a case of sun-command execution, nil otherwise.
func (inv *Invocation) Parent() *Invocation {
	return inv.parent
}

// Root returns the roo Invocation, which is the top-level
// Invocation in a case of sun-command execution. For the
// root Invocation this function returns pointer to self.
func (inv *Invocation) Root() *Invocation {
	return inv.root
}

// IsImmediate returns true, if Invocation contains some active [Option]
// with non-nil Immediate callback (see Option description for more
// information).
//
// See [Option] description for details.
func (inv *Invocation) IsImmediate() bool {
	return inv.immediate != nil
}

// Cmd returns a reference to [Command], invoked by this Invocation
func (inv *Invocation) Cmd() *Command {
	return inv.cmd
}

// Argv returns Command's argv, as passed to [Command.Parse] or
// [Command.Run].
//
// For sub-commands it contains only the tail part of argv,
// related to the sub-command being invoked.
func (inv *Invocation) Argv() []string {
	return inv.argv
}

// Get returns the first value of option or parameter by its name.
//
// [Option] name is the same as used in option definition (as Name
// or Aliases). I.e., "-h", "--help" and so on.
//
// [Parameter] name is the Parameter.Name without "decorations"
// (i.e., with removed square braces and elipses:
//
//	"name", "[name]", "name...", "[name...]  - access as "name"
//
// The value of flag options (options that don't expect explicit
// value) considered to be an empty string.
func (inv *Invocation) Get(name string) (val string, found bool) {
	vals, found := inv.byName[name]
	if found && len(vals) > 0 {
		val = vals[0]
	}

	return
}

// GetString is a simplified form of [Invocation.Get].
//
// It returns the first value of option or parameter by its name
// or "" if parameter is not found.
func (inv *Invocation) GetString(name string) string {
	s, _ := inv.Get(name)
	return s
}

// Flag returns true if option is present in the command line.
func (inv *Invocation) Flag(name string) bool {
	_, found := inv.Get(name)
	return found
}

// Values returns a slice of all values of option or parameter by
// its name.
//
// [Option] name is the same as used in option definition (as Name
// or Aliases). I.e., "-h", "--help" and so on.
//
// [Parameter] name is the Parameter.Name without "decorations"
// (i.e., with removed square braces and elipses:
//
//	"name", "[name]", "name...", "[name...]  - access as "name"
//
// The value of flag options (options that don't expect explicit
// value) considered to be an empty string.
//
// For repeated flag options, the returned slice will contain one
// empty string per each occurrence.
func (inv *Invocation) Values(name string) []string {
	return inv.byName[name]
}

// ParamCount returns count of positional parameters.
func (inv *Invocation) ParamCount() int {
	return len(inv.parameters)
}

// ParamGet returns value of the n-th positional parameter.
// If n is our of range, it returns empty string ("").
//
// Parameters are numbered, starting from zero, by their actual position
// within command's arguments, not by their position withing Parameters
// slice in [Command] description. This difference becomes important when
// we come to repeated parameters. Repeated parameter will take only one slot
// in the Parameters slice, but may be repeated (and take many positions)
// in the [Command]'s arguments.
func (inv *Invocation) ParamGet(n int) (param string) {
	if 0 <= n && n < len(inv.parameters) {
		param = inv.parameters[n]
	}

	return
}

// SubCommand returns Command's SubCommand and its arguments. If Command
// doesn't have SubCommands, this function returns (nil, nil).
func (inv *Invocation) SubCommand() (*Command, []string) {
	return inv.subcmd, inv.subargv
}
