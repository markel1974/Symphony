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
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"log"
	"sort"
	"strings"
)

// Command defines the structure for application commands, holding metadata, behavior, and configuration properties.
type Command struct {
	name                       string
	aliases                    []string
	shortHelp                  string
	longHelp                   string
	suggestionsMinimumDistance int
	daemon                     bool
	background                 bool
	template                   *Template
	run                        interfaces.RunFn
	timerEvent                 interfaces.TimerFn
	readEvent                  interfaces.ReadFn
	paintEvent                 interfaces.PaintFn
	commands                   []*Command
	parent                     *Command
	commandsMaxUseLen          int
	commandsMaxCommandPathLen  int
	commandsMaxNameLen         int
}

// NewCommand creates and returns a new instance of Command with a pre-defined template.
func NewCommand(name string, aliases []string, daemon bool, run interfaces.RunFn) *Command {
	if run == nil {
		run = func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
			return nil
		}
	}
	if i := strings.Index(name, " "); i >= 0 {
		name = name[:i]
	}
	return &Command{
		name:     name,
		aliases:  aliases,
		template: NewTemplate(),
		daemon:   daemon,
		run:      run,
	}
}

func (c *Command) PaintEvent() interfaces.PaintFn {
	return c.paintEvent
}

func (c *Command) ReadEvent() interfaces.ReadFn {
	return c.readEvent
}

func (c *Command) TimerEvent() interfaces.TimerFn {
	return c.timerEvent
}

func (c *Command) SetReadFn(fn interfaces.ReadFn) {
	c.readEvent = fn
}

func (c *Command) SetTimerFn(fn interfaces.TimerFn) {
	c.timerEvent = fn
}

func (c *Command) SetPaintFn(fn interfaces.PaintFn) {
	c.paintEvent = fn
}

func (c *Command) Daemon() bool {
	return c.daemon
}

func (c *Command) Background() bool {
	return c.background
}

func (c *Command) SetHelp(short string, long string) {
	c.shortHelp = short
	c.longHelp = long
}

// Name returns the name of the command.
func (c *Command) Name() string {
	return c.name
}

// Root navigates up the command hierarchy and returns the root command of the current command hierarchy.
func (c *Command) Root() interfaces.ICommand {
	if c.HasParent() {
		return c.Parent().Root()
	}
	return c
}

// Parent returns the immediate parent command of the current command, or nil if the command has no parent.
func (c *Command) Parent() interfaces.ICommand {
	return c.parent
}

// Childs returns a slice of pointers to the Command's child commands.
func (c *Command) Childs() []interfaces.ICommand {
	var commands []interfaces.ICommand
	for _, cmd := range c.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Help returns a string containing the help details for the command, generated using its associated help function.
func (c *Command) Help() string {
	help := c.template.Help()
	data, err := c.template.Exec(c, help)
	if err != nil {
		log.Printf("Help: %s", err.Error())
	}
	return data
}

// FindChildren searches for a child command by its name or alias and returns it. Returns nil if no matching command is found.
func (c *Command) FindChildren(name string) interfaces.ICommand {
	for _, cmd := range c.commands {
		if cmd.Name() == name || cmd.HasAlias(name) {
			return cmd
		}
	}
	return nil
}

// FindChildrenPrefix searches for a child command whose name starts with the given prefix and returns it if found.
func (c *Command) FindChildrenPrefix(prefix string) interfaces.ICommand {
	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd.Name(), prefix) {
			return cmd
		}
	}
	return nil
}

func (c *Command) Find(args []string) (interfaces.ICommand, []string, error) {
	var innerFind func(*Command, []string) (*Command, []string)
	innerFind = func(c *Command, innerArgs []string) (*Command, []string) {
		var argsWOFlags = innerArgs
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
	return commandFound, a, nil
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
func (c *Command) Traverse(args []string) (interfaces.ICommand, []string, error) {
	for i, arg := range args {
		cmd := c.findNext(arg)
		if cmd == nil {
			return c, args, nil
		}
		return cmd.Traverse(args[i+1:])
	}
	return c, args, nil
}

// SuggestionsFor returns a list of command suggestions based on the provided typedName, considering prefix and edit distance.
func (c *Command) SuggestionsFor(typedName string) []string {
	var suggestions []string
	distance := c.suggestionsMinimumDistance
	if distance <= 0 {
		distance = 2
	}
	for _, cmd := range c.commands {
		ld := levenshteinDistance(typedName, cmd.Name(), true)
		suggestByLevenshtein := ld <= distance
		suggestByPrefix := strings.HasPrefix(strings.ToLower(cmd.Name()), strings.ToLower(typedName))
		if suggestByLevenshtein || suggestByPrefix {
			suggestions = append(suggestions, cmd.Name())
		}
	}
	return suggestions
}

// FindSuggestions generates a string of suggestion messages for the given argument if there are relevant suggestions available.
// Suggestions are returned only if DisableSuggestions is false and relevant matches are found through SuggestionsFor.
func (c *Command) FindSuggestions(arg string) []string {
	var ret []string
	if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 {
		for _, s := range suggestions {
			ret = append(ret, s)
		}
	}
	return ret
}

// VisitParents calls the provided function for each parent command in the hierarchy, traversing upwards recursively.
func (c *Command) VisitParents(fn func(interfaces.ICommand)) {
	if c.HasParent() {
		fn(c.Parent())
		c.Parent().VisitParents(fn)
	}
}

// Execute runs the command with the provided context, arguments, and process ID, handling flags and validations.
func (c *Command) Execute(r interfaces.IContext, arg []string, pid int) error {
	if err := c.run(r, c, pid, arg); err != nil {
		return err
	}
	return nil
}

// Commands returns a sorted slice of subcommands by their names.
func (c *Command) Commands() []interfaces.ICommand {
	sort.Sort(commandSorterByName(c.commands))
	var commands []interfaces.ICommand
	for _, cmd := range c.commands {
		commands = append(commands, cmd)
	}
	return commands
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

// CommandPath returns the full path to the command by appending parent command names separated by a slash.
func (c *Command) CommandPath() string {
	if c.HasParent() {
		return c.Parent().CommandPath() + interfaces.PathSeparator + c.Name()
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

// HasSubCommands checks if the Command has any subcommands. Returns true if there are subcommands, otherwise false.
func (c *Command) HasSubCommands() bool {
	return len(c.commands) > 0
}

// IsAdditionalHelpTopicCommand determines if the command is a help topic by checking its runnability, visibility, and subcommands.
func (c *Command) IsAdditionalHelpTopicCommand() bool {
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

// HasParent checks if the command has a parent by verifying if the parent field is not nil. Returns true if parent exists.
func (c *Command) HasParent() bool {
	return c.parent != nil
}
