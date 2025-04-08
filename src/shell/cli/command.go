/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"errors"
	"fmt"
	"github.com/markel1974/c64emu/src/shell/cli/mflag"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"io"
	"log"
	"sort"
	"strings"
)

// CompOneRequiredFlag is a marker used to annotate flags that are required for command execution.
var CompOneRequiredFlag = "completion_one_required_flag"

// ErrSubCommandRequired is the error returned when a subcommand is required but not specified in the command execution.
var ErrSubCommandRequired = errors.New("subcommand is required")

// DefaultEol defines the default end-of-line sequence used in the application, which is set to "\r\n".
var DefaultEol = "\r\n"

// FParseErrWhitelist is a type alias for mflag.ParseErrorsWhitelist, specifying parsing errors that can be ignored.
type FParseErrWhitelist mflag.ParseErrorsWhitelist

// minUsagePadding specifies the minimum padding length used for aligning command usage text in help output.
var _minUsagePadding = 25

// Command defines the structure for application commands, holding metadata, behavior, and configuration properties.
type Command struct {
	name                       string
	aliases                    []string
	SuggestFor                 []string
	ShortHelp                  string
	LongHelp                   string
	ValidArgs                  []string
	Args                       PositionalArgs
	ArgAliases                 []string
	Hidden                     bool
	version                    string
	SilenceErrors              bool
	SilenceUsage               bool
	DisableFlagParsing         bool
	DisableFlagsInUseLine      bool
	DisableSuggestions         bool
	SuggestionsMinimumDistance int
	FParseErrWhitelist         FParseErrWhitelist
	Activate                   bool
	Background                 bool
	Pid                        int
	template                   *Template
	Run                        func(r interfaces.IContext, cmd *Command, pid int, args []string) error
	TimerEvent                 func(r interfaces.IContext, cmd *Command, pid int, tid int, ctx interface{}, interval int)
	ReadEvent                  func(r interfaces.IContext, cmd *Command, pid int, ctx interface{}, code int, key rune)
	PaintEvent                 func(r interfaces.IContext, cmd *Command, pid int, ctx interface{}, surface interfaces.ISurface)
	commands                   []*Command // commands is the list of commands supported by this program.
	parent                     *Command   // parent is a parent command for this command.
	commandsMaxUseLen          int
	commandsMaxCommandPathLen  int
	commandsMaxNameLen         int
	commandsAreSorted          bool
	flags                      *mflag.FlagSet // flags is full set of flags.
	pFlags                     *mflag.FlagSet // pFlags contains persistent flags.
	lFlags                     *mflag.FlagSet // lFlags contains local flags.
	iFlags                     *mflag.FlagSet // iFlags contains inherited flags.
	parentsPFlags              *mflag.FlagSet // parentsPFlags is all persistent flags of cmd's parents.
	helpCommand                *Command       // helpCommand is command with usage 'help'. If it's not defined by user, default help command.
}

// NewCommand creates and returns a new instance of Command with a pre-defined template.
func NewCommand() *Command {
	return &Command{
		template: NewTemplate(),
	}
}

// Version returns the version of the command.
func (c *Command) Version() string {
	return c.version
}

// SetName sets the name of the command and its aliases. Name is truncated at the first space if present.
func (c *Command) SetName(name string, aliases []string) {
	c.name = name
	if i := strings.Index(name, " "); i >= 0 {
		c.name = c.name[:i]
	}
	c.aliases = aliases
}

// Name returns the name of the command.
func (c *Command) Name() string {
	return c.name
}

// Root navigates up the command hierarchy and returns the root command of the current command hierarchy.
func (c *Command) Root() *Command {
	if c.HasParent() {
		return c.Parent().Root()
	}
	return c
}

// Parent returns the immediate parent command of the current command, or nil if the command has no parent.
func (c *Command) Parent() *Command {
	return c.parent
}

// Childs returns a slice of pointers to the Command's child commands.
func (c *Command) Childs() []*Command {
	return c.commands
}

// Usage returns the usage information of the command as a string.
func (c *Command) Usage(w io.Writer) {
	c.mergePersistentFlags(w)
	err := c.template.Exec(w, c, c.UsageTemplate())
	if err != nil {
		log.Printf("Usage: %s", err.Error())
	}
}

// Help returns a string containing the help details for the command, generated using its associated help function.
func (c *Command) Help(w io.Writer) {
	c.mergePersistentFlags(w)
	err := c.template.Exec(w, c, c.templateHelp())
	if err != nil {
		log.Printf("Help: %s", err.Error())
	}
}

// _flagErrorFunc retrieves the flag error function for the command, cascading to the parent command if not defined.
func (c *Command) _flagErrorFunc() (f func(*Command, error) error) {
	if c.HasParent() {
		return c.parent._flagErrorFunc()
	}
	return func(c *Command, err error) error {
		return err
	}
}

// UsagePadding calculates and returns the appropriate padding value for a command's usage output alignment.
func (c *Command) UsagePadding() int {
	if c.parent == nil || _minUsagePadding > c.parent.commandsMaxUseLen {
		return _minUsagePadding
	}
	return c.parent.commandsMaxUseLen
}

// minCommandPathPadding defines the minimum padding length for command path alignment in output formatting.
var minCommandPathPadding = 11

// CommandPathPadding calculates and returns the padding size for the command path based on its parent attributes.
func (c *Command) CommandPathPadding() int {
	if c.parent == nil || minCommandPathPadding > c.parent.commandsMaxCommandPathLen {
		return minCommandPathPadding
	}
	return c.parent.commandsMaxCommandPathLen
}

// minNamePadding defines the minimum padding for command names to ensure proper alignment in output formatting.
var minNamePadding = 11

// NamePadding returns the padding value for command names based on their length or a minimum padding value.
func (c *Command) NamePadding() int {
	if c.parent == nil || minNamePadding > c.parent.commandsMaxNameLen {
		return minNamePadding
	}
	return c.parent.commandsMaxNameLen
}

// hasNoOptDefVal checks whether a flag with the given name has a non-empty NoOptDefVal in the provided FlagSet.
func hasNoOptDefVal(name string, fs *mflag.FlagSet) bool {
	xFlag := fs.Lookup(name)
	if xFlag == nil {
		return false
	}
	return xFlag.NoOptDefVal != ""
}

// shortHasNoOptDefVal checks if a shorthand flag exists within the provided FlagSet and has a non-empty NoOptDefVal.
func shortHasNoOptDefVal(writer io.Writer, name string, fs *mflag.FlagSet) bool {
	if len(name) == 0 {
		return false
	}
	xFlag := fs.ShorthandLookup(writer, name[:1])
	if xFlag == nil {
		return false
	}
	return xFlag.NoOptDefVal != ""
}

// FindChildren searches for a child command by its name or alias and returns it. Returns nil if no matching command is found.
func (c *Command) FindChildren(name string) *Command {
	for _, cmd := range c.commands {
		if cmd.Name() == name || cmd.HasAlias(name) {
			return cmd
		}
	}
	return nil
}

// FindChildrenPrefix searches for a child command whose name starts with the given prefix and returns it if found.
func (c *Command) FindChildrenPrefix(prefix string) *Command {
	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd.Name(), prefix) {
			return cmd
		}
	}
	return nil
}

func (c *Command) Find(writer io.Writer, args []string) (*Command, []string, error) {
	var innerFind func(*Command, []string) (*Command, []string)
	innerFind = func(c *Command, innerArgs []string) (*Command, []string) {
		var argsWOFlags = c.stripFlags(writer, innerArgs)
		if len(argsWOFlags) == 0 {
			return c, innerArgs
		}
		nextSubCmd := argsWOFlags[0]
		var cmd = c.findNext(nextSubCmd)
		if cmd != nil {
			return innerFind(cmd, argsMinusFirstX(innerArgs, nextSubCmd))
		}
		return c, innerArgs
	}
	commandFound, a := innerFind(c, args)
	if commandFound.Args == nil {
		return commandFound, a, legacyArgs(commandFound, commandFound.stripFlags(writer, a))
	}
	return commandFound, a, nil
}

// FindSuggestions generates a string of suggestion messages for the given argument if there are relevant suggestions available.
// Suggestions are returned only if DisableSuggestions is false and relevant matches are found through SuggestionsFor.
func (c *Command) FindSuggestions(arg string) string {
	if c.DisableSuggestions {
		return ""
	}
	suggestionsString := ""
	if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 {
		suggestionsString += DefaultEol + DefaultEol + "Did you mean this?" + DefaultEol
		for _, s := range suggestions {
			suggestionsString += fmt.Sprintf("\t%v"+DefaultEol, s)
		}
	}
	return suggestionsString
}

// findNext searches for the next command matching the specified name or alias in the list of subcommands.
// It returns the matching command or nil if no suitable match is found.
// If only one command partially matches the name during prefix matching, it will return that command.
func (c *Command) findNext(next string) *Command {
	matches := make([]*Command, 0)
	for _, cmd := range c.commands {
		if cmd.Name() == next || cmd.HasAlias(next) {
			return cmd
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}

	return nil
}

// Traverse navigates through commands and flags based on the provided arguments and returns the matched command, remaining arguments, and error if any.
func (c *Command) Traverse(writer io.Writer, args []string) (*Command, []string, error) {
	var flags []string
	inFlag := false
	for i, arg := range args {
		switch {
		// A long mFlag with a space separated value
		case strings.HasPrefix(arg, "--") && !strings.Contains(arg, "="):
			// TODO: this isn't quite right, we should really check ahead for 'true' or 'false'
			inFlag = !hasNoOptDefVal(arg[2:], c.Flags())
			flags = append(flags, arg)
			continue
		// A short mFlag with a space separated value
		case strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") && len(arg) == 2 && !shortHasNoOptDefVal(writer, arg[1:], c.Flags()):
			inFlag = true
			flags = append(flags, arg)
			continue
		// The value for a mFlag
		case inFlag:
			inFlag = false
			flags = append(flags, arg)
			continue
		// A mFlag without a value, or with an `=` separated value
		case isFlagArg(arg):
			flags = append(flags, arg)
			continue
		}
		cmd := c.findNext(arg)
		if cmd == nil {
			return c, args, nil
		}
		if err := c.ParseFlags(writer, flags); err != nil {
			return nil, args, err
		}
		return cmd.Traverse(writer, args[i+1:])
	}
	return c, args, nil
}

// SuggestionsFor returns a list of command suggestions based on the provided typedName, considering prefix and edit distance.
func (c *Command) SuggestionsFor(typedName string) []string {
	var suggestions []string
	distance := c.SuggestionsMinimumDistance
	if distance <= 0 {
		distance = 2
	}
	for _, cmd := range c.commands {
		if cmd.IsAvailableCommand() {
			ld := levenshteinDistance(typedName, cmd.Name(), true)
			suggestByLevenshtein := ld <= distance
			suggestByPrefix := strings.HasPrefix(strings.ToLower(cmd.Name()), strings.ToLower(typedName))
			if suggestByLevenshtein || suggestByPrefix {
				suggestions = append(suggestions, cmd.Name())
			}
			for _, explicitSuggestion := range cmd.SuggestFor {
				if strings.EqualFold(typedName, explicitSuggestion) {
					suggestions = append(suggestions, cmd.Name())
				}
			}
		}
	}
	return suggestions
}

// VisitParents calls the provided function for each parent command in the hierarchy, traversing upwards recursively.
func (c *Command) VisitParents(fn func(*Command)) {
	if c.HasParent() {
		fn(c.Parent())
		c.Parent().VisitParents(fn)
	}
}

// ArgsLenAtDash returns the number of arguments before the first dash in the command-line arguments.
func (c *Command) ArgsLenAtDash() int {
	return c.Flags().ArgsLenAtDash()
}

// Execute runs the command with the provided context, arguments, and process ID, handling flags and validations.
func (c *Command) Execute(r interfaces.IContext, a []string, pid int) error {
	c.InitDefaultHelpFlag(r.GetWriter())
	c.InitDefaultVersionFlag(r.GetWriter())

	if err := c.ParseFlags(r.GetWriter(), a); err != nil {
		return c._flagErrorFunc()(c, err)
	}

	helpVal, errHelp := c.Flags().GetBool("help")
	if errHelp != nil {
		log.Println("\"help\" flag declared as non-bool. Please correct your code")
		return errHelp
	}

	if helpVal {
		_ = c.Flags().Set(r.GetWriter(), "help", "false")
		return mflag.ErrHelp
	}

	// for back-compat, only add version mFlag behavior if version is defined
	if c.Version() != "" {
		versionVal, err := c.Flags().GetBool("version")
		if err != nil {
			log.Println("\"version\" flag declared as non-bool. Please correct your code")
			return err
		}
		if versionVal {
			err = c.template.Exec(r.GetWriter(), c, c.templateVersion())
			if err != nil {
				log.Println(err)
			}
			return err
		}
	}

	if !c.Runnable() {
		return ErrSubCommandRequired
	}
	c.preRun()
	argWoFlags := c.Flags().Args()
	if c.DisableFlagParsing {
		argWoFlags = a
	}
	if err := c.ValidateArgs(argWoFlags); err != nil {
		return err
	}
	if err := c.validateRequiredFlags(); err != nil {
		return err
	}
	if err := c.Run(r, c, pid, argWoFlags); err != nil {
		return err
	}
	return nil
}

// preRun executes all initializer functions provided in the _initializers slice before the command is run.
func (c *Command) preRun() {
	for _, x := range _initializers {
		x()
	}
}

// ValidateArgs validates the provided arguments against the defined Args function, returning an error on failure.
func (c *Command) ValidateArgs(args []string) error {
	if c.Args == nil {
		return nil
	}
	return c.Args(c, args)
}

// validateRequiredFlags verifies that all flags marked as required are set and returns an error if any are missing.
func (c *Command) validateRequiredFlags() error {
	flags := c.Flags()
	var missingFlagNames []string
	flags.VisitAll(func(xFlag *mflag.Flag) {
		requiredAnnotation, found := xFlag.Annotations[CompOneRequiredFlag]
		if !found {
			return
		}
		if (requiredAnnotation[0] == "true") && !xFlag.Changed {
			missingFlagNames = append(missingFlagNames, xFlag.Name)
		}
	})
	if len(missingFlagNames) > 0 {
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, `", "`))
	}
	return nil
}

// InitDefaultHelpFlag initializes the default "help" flag for the command if it does not already exist.
func (c *Command) InitDefaultHelpFlag(writer io.Writer) {
	c.mergePersistentFlags(writer)
	if c.Flags().Lookup("help") == nil {
		usage := "help for "
		if c.Name() == "" {
			usage += "this command"
		} else {
			usage += c.Name()
		}
		c.Flags().BoolP(writer, "help", "h", false, usage)
	}
}

// InitDefaultVersionFlag initializes the default "version" flag for the command if a version is defined.
func (c *Command) InitDefaultVersionFlag(writer io.Writer) {
	if c.Version() == "" {
		return
	}
	c.mergePersistentFlags(writer)
	if c.Flags().Lookup("version") == nil {
		usage := "version for "
		if c.Name() == "" {
			usage += "this command"
		} else {
			usage += c.Name()
		}
		c.Flags().Bool(writer, "version", false, usage)
	}
}

/*
// InitDefaultHelpCmd initializes a default "help" command for the root or a command with subcommands.
// The default "help" command provides usage details and help topics for the commands in the application.
func (c *Command) InitDefaultHelpCmd(writer io.Writer) {
	if !c.HasSubCommands() {
		return
	}
	if c.helpCommand == nil {
		c.helpCommand = &Command{
			name:      "help [command]",
			ShortHelp: "Help about any command",
			LongHelp:  `Help provides help for any command in the application. Simply type ` + c.Name() + ` help [path to command] for full details.`,
			Run: func(r interfaces.IContext, c *Command, pid int, args []string) error {
				cmd, _, e := c.Root().Find(writer, args)
				if cmd == nil || e != nil {
					r.WriteLn("Unknown help topic %#q" + DefaultEol + strings.Join(args, " "))
					r.WriteLn(c.Root().Usage())
				} else {
					cmd.InitDefaultHelpFlag(writer)
					r.WriteLn(cmd.Help(args))
				}
				return nil
			},
		}
	}
	c.RemoveCommand(c.helpCommand)
	_ = c.AddCommand(c.helpCommand)
}
*/

// ResetCommands clears all state in the Command, including parent, sub-commands, help command, and persistent flags.
func (c *Command) ResetCommands() {
	c.parent = nil
	c.commands = nil
	c.helpCommand = nil
	c.parentsPFlags = nil
}

// Commands returns a sorted slice of subcommands by their names.
func (c *Command) Commands() []*Command {
	sort.Sort(commandSorterByName(c.commands))
	return c.commands
}

// AddCommand adds one or more subcommands to the current command.
// It ensures a command cannot be added as a child of itself.
func (c *Command) AddCommand(cx ...*Command) error {
	for i, x := range cx {
		if cx[i] == c {
			return errors.New("command can't be a child of itself")
		}
		cx[i].parent = c
		usageLen := len(x.name)
		if usageLen > c.commandsMaxUseLen {
			c.commandsMaxUseLen = usageLen
		}
		commandPathLen := len(x.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen {
			c.commandsMaxCommandPathLen = commandPathLen
		}
		nameLen := len(x.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
		}
		c.commands = append(c.commands, x)
		c.commandsAreSorted = false
	}

	return nil
}

// RemoveCommand removes one or more specified subcommands from the current command's list of subcommands.
func (c *Command) RemoveCommand(cx ...*Command) {
	var commands []*Command
main:
	for _, command := range c.commands {
		for _, cmd := range cx {
			if command == cmd {
				command.parent = nil
				continue main
			}
		}
		commands = append(commands, command)
	}
	c.commands = commands
	// recompute all lengths
	c.commandsMaxUseLen = 0
	c.commandsMaxCommandPathLen = 0
	c.commandsMaxNameLen = 0
	for _, command := range c.commands {
		usageLen := len(command.name)
		if usageLen > c.commandsMaxUseLen {
			c.commandsMaxUseLen = usageLen
		}
		commandPathLen := len(command.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen {
			c.commandsMaxCommandPathLen = commandPathLen
		}
		nameLen := len(command.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
		}
	}
}

const pathSeparator = "/"

// CommandPath returns the full path to the command by appending parent command names separated by a slash.
func (c *Command) CommandPath() string {
	if c.HasParent() {
		return c.Parent().CommandPath() + pathSeparator + c.Name()
	}
	return c.Name()
}

func (c *Command) Path() []string {
	if c.HasParent() {
		var out []string
		if path := c.Parent().Path(); len(path) > 0 {
			out = append(out, path...)
		}
		if len(c.Name()) > 0 {
			out = append(out, c.Name())
		}
		return out
	}
	if len(c.Name()) == 0 {
		return nil
	}
	return []string{c.Name()}
}

// DebugFlags writes debugging information about command flags and subcommands to the provided writer.
func (c *Command) DebugFlags(writer io.Writer) {
	_, _ = writer.Write([]byte("DebugFlags called on " + c.Name()))
	var debugFlags func(*Command)
	debugFlags = func(x *Command) {
		if x.HasFlags() || x.HasPersistentFlags() {
			_, _ = writer.Write([]byte(x.Name()))
		}
		if x.HasFlags() {
			x.flags.VisitAll(func(f *mflag.Flag) {
				if x.HasPersistentFlags() && x.persistentFlag(writer, f.Name) != nil {
					_, _ = writer.Write([]byte("  -" + f.Shorthand + "," + "--" + f.Name + "[" + f.DefValue + "]" + "" + f.Value.String() + "  [LP]"))
				} else {
					_, _ = writer.Write([]byte("  -" + f.Shorthand + "," + "--" + f.Name + "[" + f.DefValue + "]" + "" + f.Value.String() + "  [L]"))
				}
			})
		}
		if x.HasPersistentFlags() {
			x.pFlags.VisitAll(func(f *mflag.Flag) {
				if x.HasFlags() {
					if x.flags.Lookup(f.Name) == nil {
						_, _ = writer.Write([]byte("  -" + f.Shorthand + "," + "--" + f.Name + "[" + f.DefValue + "]" + "" + f.Value.String() + "  [P]"))
					}
				} else {
					_, _ = writer.Write([]byte("  -" + f.Shorthand + "," + "--" + f.Name + "[" + f.DefValue + "]" + "" + f.Value.String() + "  [P]"))
				}
			})
		}
		if x.HasSubCommands() {
			for _, y := range x.commands {
				debugFlags(y)
			}
		}
	}
	debugFlags(c)
}

// HasAlias checks if the given string matches any of the aliases of the command and returns true if a match is found.
func (c *Command) HasAlias(s string) bool {
	for _, a := range c.aliases {
		if a == s {
			return true
		}
	}
	return false
}

// hasNameOrAliasPrefix checks if the given prefix matches the start of the command's name or any of its aliases.
func (c *Command) hasNameOrAliasPrefix(prefix string) bool {
	if strings.HasPrefix(c.Name(), prefix) {
		return true
	}
	for _, alias := range c.aliases {
		if strings.HasPrefix(alias, prefix) {
			return true
		}
	}
	return false
}

func (c *Command) NameAndAliases() string {
	return strings.Join(append([]string{c.Name()}, c.aliases...), ", ")
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

// HasSubCommands checks if the Command has any subcommands. Returns true if there are subcommands, otherwise false.
func (c *Command) HasSubCommands() bool {
	return len(c.commands) > 0
}

// IsAvailableCommand determines if the command is visible and runnable or has available subcommands, excluding hidden commands.
func (c *Command) IsAvailableCommand() bool {
	if c.Hidden {
		return false
	}
	if c.HasParent() && c.Parent().helpCommand == c {
		return false
	}
	if c.Runnable() || c.HasAvailableSubCommands() {
		return true
	}
	return false
}

// IsAdditionalHelpTopicCommand determines if the command is a help topic by checking its runnability, visibility, and subcommands.
func (c *Command) IsAdditionalHelpTopicCommand() bool {
	// if a command is runnable, deprecated, or hidden, it is not a 'help' command
	if c.Runnable() || c.Hidden {
		return false
	}

	// if any non-help sub commands are found, the command is not a 'help' command
	for _, sub := range c.commands {
		if !sub.IsAdditionalHelpTopicCommand() {
			return false
		}
	}
	return true
}

// HasHelpSubCommands checks if the command has any subcommands that are additional help topics.
func (c *Command) HasHelpSubCommands() bool {
	for _, sub := range c.commands {
		if sub.IsAdditionalHelpTopicCommand() {
			return true
		}
	}
	return false
}

// HasAvailableSubCommands checks if the command has any non-deprecated, non-hidden, or non-help subcommands available.
func (c *Command) HasAvailableSubCommands() bool {
	// return true on the first found available (non deprecated/help/hidden) sub command
	for _, sub := range c.commands {
		if sub.IsAvailableCommand() {
			return true
		}
	}
	// the command either has no sub commands, or no available (non deprecated/help/hidden) sub commands
	return false
}

// HasParent checks if the command has a parent by verifying if the parent field is not nil. Returns true if parent exists.
func (c *Command) HasParent() bool {
	return c.parent != nil
}

// Flags returns the FlagSet associated with the Command. It initializes the FlagSet if it is not already created.
func (c *Command) Flags() *mflag.FlagSet {
	if c.flags == nil {
		c.flags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	}
	return c.flags
}

// LocalNonPersistentFlags returns a flag set containing local flags that are not marked as persistent for the command.
func (c *Command) LocalNonPersistentFlags(writer io.Writer) *mflag.FlagSet {
	persistentFlags := c.PersistentFlags()
	out := mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	c.LocalFlags(writer).VisitAll(func(f *mflag.Flag) {
		if persistentFlags.Lookup(f.Name) == nil {
			out.AddFlag(writer, f)
		}
	})
	return out
}

func (c *Command) LocalFlags(writer io.Writer) *mflag.FlagSet {
	c.mergePersistentFlags(writer)
	if c.lFlags == nil {
		c.lFlags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	}
	c.lFlags.SortFlags = c.Flags().SortFlags
	addToLocal := func(f *mflag.Flag) {
		if c.lFlags.Lookup(f.Name) == nil && c.parentsPFlags.Lookup(f.Name) == nil {
			c.lFlags.AddFlag(writer, f)
		}
	}
	c.Flags().VisitAll(addToLocal)
	c.PersistentFlags().VisitAll(addToLocal)
	return c.lFlags
}

// InheritedFlags merges persistent and parent flags into a single FlagSet and returns it.
func (c *Command) InheritedFlags(writer io.Writer) *mflag.FlagSet {
	c.mergePersistentFlags(writer)
	if c.iFlags == nil {
		c.iFlags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	}
	local := c.LocalFlags(writer)
	c.parentsPFlags.VisitAll(func(f *mflag.Flag) {
		if c.iFlags.Lookup(f.Name) == nil && local.Lookup(f.Name) == nil {
			c.iFlags.AddFlag(writer, f)
		}
	})
	return c.iFlags
}

// NonInheritedFlags returns the FlagSet containing flags specific to the command, excluding those inherited from parent commands.
func (c *Command) NonInheritedFlags(writer io.Writer) *mflag.FlagSet {
	return c.LocalFlags(writer)
}

// PersistentFlags returns a FlagSet containing the persistent flags of the command, initializing it if necessary.
func (c *Command) PersistentFlags() *mflag.FlagSet {
	if c.pFlags == nil {
		c.pFlags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	}
	return c.pFlags
}

// ResetFlags reinitializes the flags and persistent flags for the Command by creating new empty flag sets.
func (c *Command) ResetFlags() {
	c.flags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	c.pFlags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
	c.lFlags = nil
	c.iFlags = nil
	c.parentsPFlags = nil
}

func (c *Command) HasFlags() bool {
	return c.Flags().HasFlags()
}

// HasPersistentFlags checks if the command has persistent flags defined and returns true if any are present.
func (c *Command) HasPersistentFlags() bool {
	return c.PersistentFlags().HasFlags()
}

func (c *Command) HasLocalFlags(writer io.Writer) bool {
	return c.LocalFlags(writer).HasFlags()
}

// HasInheritedFlags checks if any inherited flags are present for the command by analyzing the provided writer.
func (c *Command) HasInheritedFlags(writer io.Writer) bool {
	return c.InheritedFlags(writer).HasFlags()
}

// HasAvailableFlags checks if the command has any non-hidden flags available for use.
func (c *Command) HasAvailableFlags() bool {
	return c.Flags().HasAvailableFlags()
}

// HasAvailablePersistentFlags checks if the command has any persistent flags that are available for use.
func (c *Command) HasAvailablePersistentFlags() bool {
	return c.PersistentFlags().HasAvailableFlags()
}

// HasAvailableLocalFlags returns true if there are any local flags available for the command.
func (c *Command) HasAvailableLocalFlags(writer io.Writer) bool {
	return c.LocalFlags(writer).HasAvailableFlags()
}

// HasAvailableInheritedFlags checks if the command has inherited flags that are available for use, writing output to the provided writer.
func (c *Command) HasAvailableInheritedFlags(writer io.Writer) bool {
	return c.InheritedFlags(writer).HasAvailableFlags()
}

// Flag retrieves a flag by name from the command's flag set or creates a persistent flag if it doesn't exist.
func (c *Command) Flag(writer io.Writer, name string) *mflag.Flag {
	xFlag := c.Flags().Lookup(name)
	if xFlag == nil {
		xFlag = c.persistentFlag(writer, name)
	}
	return xFlag
}

// persistentFlag retrieves a persistent flag from the current command or parent commands by its name.
// It first checks the current command's persistent flags and updates parent flags if necessary.
// Returns a pointer to the mflag.Flag if found, or nil if not.
func (c *Command) persistentFlag(writer io.Writer, name string) *mflag.Flag {
	var pFlag *mflag.Flag = nil
	if c.HasPersistentFlags() {
		pFlag = c.PersistentFlags().Lookup(name)
	}
	if pFlag == nil {
		c.updateParentsFlags(writer)
		pFlag = c.parentsPFlags.Lookup(name)
	}
	return pFlag
}

// ParseFlags parses command-line flags using the provided writer for output and the given arguments slice.
func (c *Command) ParseFlags(writer io.Writer, args []string) error {
	if c.DisableFlagParsing {
		return nil
	}
	c.mergePersistentFlags(writer)
	c.Flags().ParseErrorsWhitelist = mflag.ParseErrorsWhitelist(c.FParseErrWhitelist)
	if err := c.Flags().Parse(writer, args); err != nil {
		return err
	}
	return nil
}

// mergePersistentFlags merges persistent flags from the command and its parents into the command's flag set.
func (c *Command) mergePersistentFlags(writer io.Writer) {
	c.updateParentsFlags(writer)
	c.Flags().AddFlagSet(writer, c.PersistentFlags())
	c.Flags().AddFlagSet(writer, c.parentsPFlags)
}

// updateParentsFlags updates the parent's persistent flag set for the command, merging it with the global normalization function.
func (c *Command) updateParentsFlags(writer io.Writer) {
	if c.parentsPFlags == nil {
		c.parentsPFlags = mflag.NewFlagSet(c.Name(), mflag.ContinueOnError)
		c.parentsPFlags.SortFlags = false
	}
	c.Root().PersistentFlags().AddFlagSet(writer, mflag.CommandLine)
	c.VisitParents(func(parent *Command) {
		c.parentsPFlags.AddFlagSet(writer, parent.PersistentFlags())
	})
}

// templateHelp retrieves the help template of the command, delegating to the parent command if one exists.
func (c *Command) templateHelp() string {
	if c.HasParent() {
		return c.parent.templateHelp()
	}
	return c.template.Help()
}

// templateVersion retrieves the template version of the command, delegating to the parent if one exists.
func (c *Command) templateVersion() string {
	if c.HasParent() {
		return c.parent.templateVersion()
	}
	return c.template.Version()
}

// UsageTemplate returns the usage template string for the command. If the command has a parent, it delegates to the parent.
func (c *Command) UsageTemplate() string {
	if c.HasParent() {
		return c.parent.UsageTemplate()
	}
	return c.template.Usage()
}

// stripFlags removes flag-like arguments from the provided args slice and returns the remaining commands.
// It processes flags and their optional argument values, handling long flags, short flags, and "--" to terminate flags.
func (c *Command) stripFlags(writer io.Writer, args []string) []string {
	if len(args) == 0 {
		return args
	}
	c.mergePersistentFlags(writer)
	var commands []string
	flags := c.Flags()
Loop:
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		switch {
		case s == "--":
			// "--" terminates the flags
			break Loop
		case strings.HasPrefix(s, "--") && !strings.Contains(s, "=") && !hasNoOptDefVal(s[2:], flags):
			// If '--flag arg' then
			// delete arg from args.
			fallthrough // (do the same as below)
		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2 && !shortHasNoOptDefVal(writer, s[1:], flags):
			// If '-f arg' then
			// delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 {
				break Loop
			} else {
				args = args[1:]
				continue
			}
		case s != "" && !strings.HasPrefix(s, "-"):
			commands = append(commands, s)
		}
	}
	return commands
}
