// MFP  - Miulti-Function Printers and scanners toolkit
// argv - Argv parsing mini-library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Argv parser

package argv

import (
	"fmt"
	"math"
	"strings"
)

// parser implements command line parsing.
//
// Notes.
//
// optConflicts["--opt1"] -> "--opt2" means, that previously
// processed option "--opt2" has declared option "--opt1" as
// conflicting.
//
// optRequires["--opt1"] -> "--opt2" means, that previously
// processed option "-opt2" has declared option "--opt1" as
// required.
//
// optSeen["--opt1"] -> "--opt2" means that option "--opt1"
// already seen, using "--opt2" among all its possible aliases.
type parser struct {
	inv          *Invocation               // Invocation being parsed
	nextarg      int                       // Index of the next argument
	optConflicts map[string]string         // Conflicting options
	optRequired  map[string]string         // Required options
	optSeen      map[string]string         // Option names seen so far
	options      map[*Option]*parserOptVal // Actually parsed options
	parameters   []parserParamVal          // Parameters by number
}

// parserOptVal represents parsed option with value
type parserOptVal struct {
	opt    *Option  // Option description
	name   string   // Actual name being used
	values []string // Option values
}

// parserParamVal represents parsed positional parameter with value
type parserParamVal struct {
	param *Parameter
	value string
}

// newParser creates a new parser.
//
// It panics, if cmd.Verify() returns an error.
func newParser(cmd *Command, argv []string) *parser {
	err := cmd.Verify()
	if err != nil {
		panic(err)
	}

	return &parser{
		inv: &Invocation{
			cmd:    cmd,
			argv:   argv,
			byName: make(map[string][]string),
		},
		optConflicts: make(map[string]string),
		optRequired:  make(map[string]string),
		optSeen:      make(map[string]string),
		options:      make(map[*Option]*parserOptVal),
	}
}

// parse parses the argv and returns parsed Invocation
func (prs *parser) parse(parent *Invocation) (*Invocation, error) {
	// Call parser
	err := prs.doParse(parent)

	// Finally post-process results
	if err != nil && len(prs.inv.argv) == 0 {
		// If argv parsing fails and argv is empty, this
		// is definitely due to the missed parameters.
		//
		// Call HelpHandler directly instead of the
		// reporting error at this case.
		inv := prs.inv
		inv.immediate = HelpHandler
		return inv, nil
	}

	if err != nil {
		return nil, err
	}

	return prs.inv, nil
}

// doParse does the prs.parse work
func (prs *parser) doParse(parent *Invocation) error {
	// Parse arguments, one by one.
	var doneOptions bool
	var paramValues []string

	paramsMin, paramsMax := prs.paramsInfo()

	for !prs.done() {
		arg := prs.next()

		var err error

		switch {
		case !doneOptions && arg == "--":
			doneOptions = true

		case !doneOptions && prs.isShortOption(arg):
			err = prs.handleShortOption(arg)

		case !doneOptions && prs.isLongOption(arg):
			err = prs.handleLongOption(arg)

		case prs.inv.cmd.hasSubCommands():
			err = prs.handleSubCommand(arg)

		case len(paramValues) < paramsMax:
			if prs.inv.cmd.NoOptionsAfterParameters {
				doneOptions = true
			}

			paramValues = append(paramValues, arg)

		default:
			err = fmt.Errorf("unexpected parameter: %q", arg)
		}

		if err != nil {
			return err
		}
	}

	// Check that we have enough parameters. Note, in the
	// immediate mode this check is suppressed.
	if prs.inv.immediate == nil {
		if len(paramValues) < paramsMin {
			missed := &prs.inv.cmd.Parameters[len(paramValues)]
			err := fmt.Errorf("missed parameter: %q", missed.Name)
			return err
		}

		if prs.inv.cmd.hasSubCommands() && prs.inv.subcmd == nil {
			return fmt.Errorf("missed sub-command name")
		}
	}

	// Toss paramValues
	if prs.inv.cmd.hasParameters() {
		err := prs.handleParameters(paramValues)
		if err != nil {
			return err
		}
	}

	// Build prs.inv.byName map
	prs.buildByName()

	// Validate things
	if err := prs.validateThings(); err != nil {
		return err
	}

	// Finish Invocation
	inv := prs.inv

	inv.parent = parent
	inv.root = inv
	if parent != nil {
		inv.root = parent.root
	}

	inv.parameters = make([]string, len(prs.parameters))

	for i := range prs.parameters {
		inv.parameters[i] = prs.parameters[i].value
	}

	return nil
}

// handleShortOption handles a short option
func (prs *parser) handleShortOption(arg string) error {
	// Split into name and value and try to find Option
	name, val, novalue := prs.splitOptVal(arg)
	opt := prs.findOption(name)
	if opt == nil {
		err := fmt.Errorf("unknown option: %q", name)
		return err
	}

	// Two simple cases:
	//   - option argument doesn't contain a value (i.e., -c, not -cXXX)
	//   - option requires a value, so argument cannot be treated as
	//     a multi-options argument
	//
	// These cases are handled the same way: we attempt to fetch
	// the next argument as option value, if value is required, and
	// let prs.appendOptVal() to do the rest.
	if novalue || opt.withValue() {
		if novalue && opt.withValue() {
			val, novalue = prs.nextValue()
		}

		return prs.appendOptVal(opt, name, val, novalue)
	}

	// Short options without value can be combined:
	//
	//   -cru equals to -c -r -u
	//
	// If we are here, we have a fist option without the value
	// and non-empty value.
	//
	// So try to consider value as a sequence of short options
	for _, c := range name[1:] + val {
		name2 := "-" + string(c)

		opt2 := prs.findOption(name2)
		if opt2 == nil {
			err := fmt.Errorf(
				"unknown option: %q",
				name2)
			return err
		}

		err := prs.appendOptVal(opt2, name2, "", true)
		if err != nil {
			return err
		}
	}

	return nil
}

// handleLongOption handles a long option
func (prs *parser) handleLongOption(arg string) error {
	name, val, novalue := prs.splitOptVal(arg)

	opt := prs.findOption(name)
	if opt == nil {
		err := fmt.Errorf("unknown option: %q", name)
		return err
	}

	if novalue && opt.withValue() {
		val, novalue = prs.nextValue()
	}

	err := prs.appendOptVal(opt, name, val, novalue)
	if err != nil {
		return err
	}

	return nil
}

// handleParameters handles positional parameters
func (prs *parser) handleParameters(paramValues []string) error {
	// Build slice of parameters' descriptors
	paramDescs := make([]*Parameter, len(paramValues))
	rept := -1

	for i := 0; i < len(paramValues); i++ {
		paramDescs[i] = &prs.inv.cmd.Parameters[i]
		if paramDescs[i].repeated() {
			rept = i
			break
		}
	}

	if rept >= 0 {
		i := len(paramDescs) - 1
		j := len(prs.inv.cmd.Parameters) - 1

		for !prs.inv.cmd.Parameters[j].repeated() {
			paramDescs[i] = &prs.inv.cmd.Parameters[j]
			i--
			j--
		}

		for i := rept + 1; i < len(paramDescs); i++ {
			if paramDescs[i] == nil {
				paramDescs[i] = paramDescs[rept]
			}
		}
	}

	// Validate parameters one by one
	for i := range paramValues {
		val := paramValues[i]
		desc := paramDescs[i]

		if desc.Validate != nil {
			err := desc.Validate(val)
			if err != nil {
				return fmt.Errorf("%q: %w %q", desc.Name, err, val)
			}
		}
	}

	// Save parameters
	prs.parameters = make([]parserParamVal, len(paramValues))
	for i := range paramValues {
		prs.parameters[i].param = paramDescs[i]
		prs.parameters[i].value = paramValues[i]
	}

	return nil
}

// handleSubCommand handles a sub-command
func (prs *parser) handleSubCommand(arg string) error {
	subcmd, err := prs.inv.cmd.FindSubCommand(arg)
	if err != nil {
		return err
	}

	prs.inv.subcmd = subcmd
	prs.inv.subargv = prs.inv.argv[prs.nextarg:]

	return nil
}

// buildByName populates prs.inv.byName map
func (prs *parser) buildByName() {
	// Save options values
	for _, optval := range prs.options {
		opt := optval.opt

		for _, name := range opt.names() {
			prs.inv.byName[name] = optval.values
		}
	}

	// Save parameters values
	//
	// Note, repeated parameters may have multiple values associated
	// with the same parameter
	for _, paramval := range prs.parameters {
		name := paramval.param.name()
		values := prs.inv.byName[name]
		values = append(values, paramval.value)

		prs.inv.byName[name] = values
	}
}

// validateThings validates things that can only be verified
// when parsing is done, like missed options requirements etc
func (prs *parser) validateThings() error {
	// These tests are suppressed if immediate option is in use
	if prs.inv.immediate != nil {
		return nil
	}

	// Check for missed required options
	for i := range prs.inv.cmd.Options {
		opt := &prs.inv.cmd.Options[i]
		if opt.Required {
			_, found := prs.inv.byName[opt.Name]
			if !found {
				return fmt.Errorf("missed option %q",
					opt.Name)
			}
		}
	}

	for required, byWhom := range prs.optRequired {
		if _, found := prs.inv.byName[required]; !found {
			return fmt.Errorf("missed option %q, required by %q",
				required, byWhom)
		}
	}
	return nil
}

// complete handles command auto-completion
func (prs *parser) complete() (compl []Completion) {
	doneOptions := false
	paramCount := 0

	// Roll over all arguments. Here our goals are:
	//
	//   - handle "--" switch, so we know when options end
	//   - count parameters, so we know what positional parameter
	//     to complete
	//   - skip short options arguments, so we don't mess them
	//     with parameters.
	//
	// If we find unknown option in the middle we can't know if
	// the subsequent argument is the option value or separate
	// parameter or option - so completion fails at this point.
	for !prs.done() {
		arg := prs.next()

		switch {
		case !doneOptions && arg == "--":
			// Handle the special case: last argument is "--"
			// and command have some long options. We don't know at
			// this point what is the meaning of the "--": beginning
			// of the long option or end of the option list.
			//
			// So if we have some long options, offer them here,
			// otherwise consider "--" as the end of the options
			// list.
			if prs.done() {
				compl := prs.completeOptionName(arg)
				if len(compl) > 0 {
					return compl
				}
			}

			doneOptions = true

		case !doneOptions && strings.HasPrefix(arg, "-"):
			name, val, novalue := prs.splitOptVal(arg)
			opt := prs.findOption(name)

			needvalue := false
			if opt != nil {
				needvalue = novalue && opt.withValue()
			}

			switch {
			case !prs.done() && opt == nil:
				// Abort if we have unknown option
				// in the middle
				return nil

			case !prs.done() && needvalue:
				// Complete or skip the value
				val = prs.next()
				if prs.done() {
					return prs.completeOptionValue(opt, val)
				}

			case prs.done():
				if novalue && opt == nil {
					// We have not reached the option
					// value and option name is not
					// known.
					//
					// Looks like we have a truncated
					// option name. Try to complete.
					return prs.completeOptionName(arg)
				}

				if opt != nil {
					// If option is not unknown, we may
					// try to complete the value.
					compl = prs.completeOptionValue(opt,
						val)

					prefix := name
					if prs.isLongOption(name) {
						prefix += "="
					}

					prs.completePrepend(compl, prefix)

					return
				}

				return nil
			}

		case prs.inv.cmd.hasSubCommands():
			subcmd, _ := prs.inv.cmd.FindSubCommand(arg)

			// If we have a sub-command here and there are more
			// argv[] arguments, simply let sub-command to
			// complete self
			if subcmd != nil && !prs.done() {
				argv := prs.inv.argv[prs.nextarg:]
				return subcmd.Complete(argv)
			}

			// If we are at the end of argv, complete
			// sub-command name
			if prs.done() {
				return prs.completeSubCommandName(arg)
			}

			// We are still at the middle of argv and
			// encountered unknown sub-command. Just abort
			// at this case.
			return nil

		default:
			// This is positional parameter. Count passed
			// parameters and complete the latest.
			if !prs.done() {
				paramCount++
			} else {
				return prs.completeParameter(arg, paramCount)
			}
		}
	}

	switch {
	case prs.inv.cmd.hasParameters():
		compl = prs.completeParameter("", paramCount)

	case prs.inv.cmd.hasSubCommands():
		compl = prs.completeSubCommandName("")

	case prs.inv.cmd.hasOptions():
		compl = prs.completeOptionName("")
	}

	return
}

// completeOptionName returns slice of completion candidates for
// Option name
func (prs *parser) completeOptionName(arg string) (compl []Completion) {
	for i := range prs.inv.cmd.Options {
		opt := &prs.inv.cmd.Options[i]

		for _, name := range opt.names() {
			if strings.HasPrefix(name, arg) {
				c := Completion{name, false}
				if opt.withValue() && prs.isLongOption(name) {
					c.String += "="
					c.NoSpace = true
				}
				compl = append(compl, c)
			}
		}
	}

	return prs.completePostProcess(arg, compl)
}

// completeOption handles auto-completion for options.
func (prs *parser) completeOptionValue(opt *Option, arg string) (
	compl []Completion) {

	compl = opt.complete(arg)
	compl = prs.completePostProcess(arg, compl)

	return
}

// completeParameter handles auto-completion for positional
// Parameters. 'n' is the count of preceding Parameters.
func (prs *parser) completeParameter(arg string, n int) (compl []Completion) {

	var paramFound *Parameter

	for i := range prs.inv.cmd.Parameters {
		param := &prs.inv.cmd.Parameters[i]
		if i == n || param.repeated() {
			paramFound = param
			break
		}
	}

	if paramFound != nil {
		compl := paramFound.complete(arg)
		return prs.completePostProcess(arg, compl)
	}

	return
}

// completeShortOption handles auto-completion for sub-commands
func (prs *parser) completeSubCommandName(arg string) (compl []Completion) {

	for i := range prs.inv.cmd.SubCommands {
		subcmd := &prs.inv.cmd.SubCommands[i]
		names := subcmd.names()

		for _, name := range names {
			if strings.HasPrefix(name, arg) {
				compl = append(compl, Completion{name, false})
			}
		}
	}

	return prs.completePostProcess(arg, compl)
}

// completePostProcess post-processes completion candidates for
// the given argument.
//
// Returns (done/compl/flags) tuple.
func (prs *parser) completePostProcess(arg string,
	compl []Completion) []Completion {
	// For sanity: drop candidates, that doesn't contain arg
	// as prefix. It actually should never happen, but..
	var complCount int
	for i := range compl {
		if strings.HasPrefix(compl[i].String, arg) {
			compl[complCount] = compl[i]
			complCount++
		}
	}

	compl = compl[:complCount]

	// If we have multiple candidates with the common prefix
	// that is longer that arg, just return that common prefix
	// as a single candidate, so when user presses Tab, the
	// completion will go to that common point
	if len(compl) > 1 {
		complStrings := make([]string, len(compl))
		for i := range compl {
			complStrings[i] = compl[i].String
		}

		prefix := strCommonPrefixSlice(complStrings)

		// If prefix is longer that arg, it means that all possible
		// candidates has common prefix, so just return it as a single
		// suggestion, so completion will go to that point.
		if len(prefix) > len(arg) {
			compl = []Completion{{prefix, true}}
		}
	}

	return compl
}

// completePrepend prepends a prefix to each completion candidate
func (prs *parser) completePrepend(compl []Completion, prefix string) {
	for i := range compl {
		compl[i].String = prefix + compl[i].String
	}
}

// isShortOption tells if argument is a short option
func (prs *parser) isShortOption(arg string) bool {
	return len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'
}

// isShortOption tells if argument is a long option
func (prs *parser) isLongOption(arg string) bool {
	return len(arg) >= 3 && arg[0] == '-' && arg[1] == '-'
}

// splitOptVal splits option argument into name and value in a case
// when they are placed into the single argument:
//
//	-cVAL     - short option case
//	-long=val - long option case
//
// It returns option name and value and the indication that value
// is actually missed (which is not necessarily the same as "" value).
func (prs *parser) splitOptVal(arg string) (name, val string, novalue bool) {
	switch {
	case prs.isShortOption(arg):
		name = arg[:2]
		val = arg[2:]
		novalue = val == ""

	case prs.isLongOption(arg):
		// For --name=value, pick out the name
		idx := strings.IndexByte(arg, '=')
		if idx >= 0 {
			name = arg[:idx]
			val = arg[idx+1:]
			novalue = false
		} else {
			name = arg
			novalue = true
		}

	default:
		novalue = true
	}

	return
}

// findOption finds Command's Option by name.
func (prs *parser) findOption(name string) *Option {
	for i := range prs.inv.cmd.Options {
		opt := &prs.inv.cmd.Options[i]
		for _, n := range opt.names() {
			if name == n {
				return opt
			}
		}
	}

	return nil
}

// paramsInfo returns information on a command parameters:
//
//	paramsMin - minimal count of parameters
//	paramsMax - maximal count of parameters
//
// If Command can accept unlimited amount of parameters
// (i.e., it has repeated parameters), paramsMax will be
// reported as math.MaxInt
func (prs *parser) paramsInfo() (paramsMin, paramsMax int) {
	for i := range prs.inv.cmd.Parameters {
		param := &prs.inv.cmd.Parameters[i]

		if param.required() {
			paramsMin++
		}

		if param.repeated() {
			paramsMax = math.MaxInt
		}
	}

	if paramsMax != math.MaxInt {
		paramsMax = len(prs.inv.cmd.Parameters)
	}

	return
}

// appendOptVal validates option value and appends
// it to the prs.options
func (prs *parser) appendOptVal(opt *Option, name, value string,
	novalue bool) error {

	// Validate things
	if novalue && opt.withValue() {
		err := fmt.Errorf("option requires operand: %q", name)
		return err
	}

	if !novalue && !opt.withValue() {
		err := fmt.Errorf("unrecognized option argument: %q",
			name+"="+value)
		return err
	}

	if !novalue {
		err := opt.Validate(value)
		if err != nil {
			return fmt.Errorf("%w: %s %q", err, name, value)
		}
	}

	if conflict, found := prs.optConflicts[name]; found {
		return fmt.Errorf("option %q conflicts with %q",
			name, conflict)
	}

	for _, conflict := range opt.Conflicts {
		if seen := prs.optSeen[conflict]; seen != "" {
			return fmt.Errorf("option %q conflicts with %q",
				name, seen)
		}
	}

	// Save the option
	optval := prs.options[opt]
	if optval == nil {
		optval = &parserOptVal{
			opt:  opt,
			name: name,
		}

		prs.options[opt] = optval
	} else if opt.Singleton {
		return fmt.Errorf("option %q cannot be repeated", opt.Name)
	}

	optval.values = append(optval.values, value)

	if opt.Immediate != nil && prs.inv.immediate == nil {
		prs.inv.immediate = opt.Immediate
	}

	// Update optConflicts, optRequired and optSeen
	for _, conflict := range opt.Conflicts {
		if _, found := prs.optConflicts[conflict]; !found {
			prs.optConflicts[conflict] = name
		}
	}

	for _, required := range opt.Requires {
		if _, found := prs.optRequired[required]; !found {
			prs.optRequired[required] = name
		}
	}

	if prs.optSeen[name] == "" {
		for _, alias := range opt.names() {
			prs.optSeen[alias] = name
		}
	}

	return nil
}

// done returns true if all arguments are consumed
func (prs *parser) done() bool {
	return prs.nextarg == len(prs.inv.argv) || prs.inv.subcmd != nil
}

// next returns the next argument.
func (prs *parser) next() (arg string) {
	if prs.nextarg < len(prs.inv.argv) {
		arg = prs.inv.argv[prs.nextarg]
		prs.nextarg++
	}

	return
}

// nextValue returns the next argument, of one exist.
func (prs *parser) nextValue() (val string, novalue bool) {
	if !prs.done() {
		return prs.next(), false
	}

	return "", true
}
