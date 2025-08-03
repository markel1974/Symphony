package process

import (
	"errors"
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"sort"
	"strings"
)

// Command represents a structured command with metadata, associated functions, and subcommands in a CLI application.
type Command struct {
	name                       string
	kind                       interfaces.CommandType
	aliases                    []string
	help                       string
	longHelp                   string
	suggestionsMinimumDistance int
	daemon                     bool
	background                 bool
	run                        interfaces.RunFn
	timerEvent                 interfaces.TimerFn
	readEvent                  interfaces.ReadFn
	readBroadcastEvent         interfaces.ReadFn
	paintEvent                 interfaces.PaintFn
	commands                   []*Command
	parent                     *Command
	maxUsageLen                int
	maxPathLen                 int
	maxNameLen                 int
}

// NewCommand creates a new Command instance with the specified name, type, aliases, daemon status, and execution function.
func NewCommand(name string, kind interfaces.CommandType, aliases []string, daemon bool, run interfaces.RunFn) *Command {
	if run == nil {
		run = func(task interfaces.IProcess, args []string) error {
			return nil
		}
	}
	if i := strings.Index(name, " "); i >= 0 {
		name = name[:i]
	}
	return &Command{
		name:    name,
		kind:    kind,
		aliases: aliases,
		daemon:  daemon,
		run:     run,
	}
}

// Type returns the structural nature of the command represented as interfaces.CommandType.
func (c *Command) Type() interfaces.CommandType {
	return c.kind
}

// DirectoryListing returns a list of all child command names, property names, and executable names associated with the command.
func (c *Command) DirectoryListing() []string {
	//TODO return an array of object <name, type....>
	var out []string
	for _, cmd := range c.Childs() {
		out = append(out, cmd.Name())
	}
	return out
}

// PaintEvent returns the PaintFn associated with the command, which handles rendering or drawing operations.
func (c *Command) PaintEvent() interfaces.PaintFn {
	return c.paintEvent
}

// ReadEvent retrieves the function assigned to handle read events for the command.
func (c *Command) ReadEvent() interfaces.ReadFn {
	return c.readEvent
}

// ReadBroadcastEvent returns a function for handling broadcast event messages within the command context.
func (c *Command) ReadBroadcastEvent() interfaces.ReadFn {
	return c.readBroadcastEvent
}

// TimerEvent returns the TimerFn associated with the Command, used to define the timer-based behavior for the command.
func (c *Command) TimerEvent() interfaces.TimerFn {
	return c.timerEvent
}

// SetReadFn sets the function to handle read-related events for the command.
func (c *Command) SetReadFn(fn interfaces.ReadFn) {
	c.readEvent = fn
}

func (c *Command) SetReadBroadcastFn(fn interfaces.ReadFn) {
	c.readBroadcastEvent = fn
}

// SetTimerFn sets the TimerFn callback function for the command's timer event.
func (c *Command) SetTimerFn(fn interfaces.TimerFn) {
	c.timerEvent = fn
}

// SetPaintFn configures a custom function to handle paint events for the command.
func (c *Command) SetPaintFn(fn interfaces.PaintFn) {
	c.paintEvent = fn
}

// Daemon returns whether the command is configured to run in daemon mode.
func (c *Command) Daemon() bool {
	return c.daemon
}

// Background returns true if the command is marked to run in the background.
func (c *Command) Background() bool {
	return c.background
}

// SetHelp sets the short and long help text for the command. Short help is used for summaries, long help for details.
func (c *Command) SetHelp(short string, long string) {
	c.help = short
	c.longHelp = long
}

// Name returns the name of the Command.
func (c *Command) Name() string {
	return c.name
}

// Root returns the root command by recursively traversing the parent commands. If no parent exists, it returns itself.
func (c *Command) Root() interfaces.ICommand {
	if c.HasParent() {
		return c.Parent().Root()
	}
	return c
}

// Parent returns the parent command of the current command, or nil if the command does not have a parent.
func (c *Command) Parent() interfaces.ICommand {
	return c.parent
}

// Childs returns a slice of ICommand representing all direct child commands of the current command.
func (c *Command) Childs() []interfaces.ICommand {
	var commands []interfaces.ICommand
	for _, cmd := range c.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Help returns detailed help information for the command. If long help is unavailable, short help is returned.
func (c *Command) Help() string {
	if len(c.longHelp) > 0 {
		return c.longHelp
	}
	if len(c.help) > 0 {
		return c.help
	}
	return ""
}

// FindChildren searches for a sub-command by its name or alias and returns it if found; otherwise, returns nil.
func (c *Command) FindChildren(name string) interfaces.ICommand {
	for _, cmd := range c.commands {
		if cmd.Name() == name || cmd.HasAlias(name) {
			return cmd
		}
	}
	return nil
}

// FindChildrenPrefix searches for a child command whose name starts with the given prefix and returns it if found.
// Returns nil if no matching child command is found.
func (c *Command) FindChildrenPrefix(prefix string) interfaces.ICommand {
	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd.Name(), prefix) {
			return cmd
		}
	}
	return nil
}

// Find locates a subcommand and its remaining arguments by traversing the command's hierarchy recursively.
// It returns the matched command, remaining arguments, and an error if something goes wrong.
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

// findNext searches for the next command by name or alias and returns it, returning the parent command if next is "..".
func (c *Command) findNext(next string) *Command {
	if next == ".." {
		return c.parent
	}
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

// Traverse navigates through a command hierarchy based on the provided path and returns the matching ICommand or nil if not found.
func (c *Command) Traverse(path []string) interfaces.ICommand {
	for i, arg := range path {
		cmd := c.findNext(arg)
		if cmd == nil {
			return nil
		}
		if next := path[i+1:]; len(next) > 0 {
			return cmd.Traverse(next)
		}
		return cmd
	}
	return nil
}

// Execute runs the command using the provided task and arguments, returning an error if execution fails.
func (c *Command) Execute(task interfaces.IProcess, arg []string) error {
	if err := c.run(task, arg); err != nil {
		return err
	}
	return nil
}

// Commands returns a sorted list of ICommand instances associated with the Command.
func (c *Command) Commands() []interfaces.ICommand {
	sort.Sort(commandSorterByName(c.commands))
	var commands []interfaces.ICommand
	for _, cmd := range c.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// AddCommand adds one or more subcommands to the current command if it is of type CommandTypeDirectory.
// It returns an error if attempting to add subcommands to a non-directory command or if a command is its own child.
func (c *Command) AddCommand(cx interfaces.ICommand) error {
	if cx == nil {
		return fmt.Errorf("nil command")
	}
	if c.kind != interfaces.CommandTypeDirectory {
		return fmt.Errorf("can't add subcommands to a non-directory command: %s", c.Name())
	}
	cy, ok := cx.(*Command)
	if !ok {
		return fmt.Errorf("invalid command type command: %s", cx.Name())
	}
	if cy == c {
		return errors.New("command can't be a child of itself")
	}
	cy.parent = c
	usageLen := len(cy.Name())
	if usageLen > c.maxUsageLen {
		c.maxUsageLen = usageLen
	}
	pathLen := len(cy.CommandPath())
	if pathLen > c.maxPathLen {
		c.maxPathLen = pathLen
	}
	nameLen := len(cy.Name())
	if nameLen > c.maxNameLen {
		c.maxNameLen = nameLen
	}
	c.commands = append(c.commands, cy)
	return nil
}

// RemoveCommand removes one or more subcommands from the parent command and updates parent-related references.
func (c *Command) RemoveCommand(cx interfaces.ICommand) error {
	if cx == nil {
		return fmt.Errorf("nil command")
	}
	cy, ok := cx.(*Command)
	if !ok {
		return fmt.Errorf("invalid command type command: %s", cx.Name())
	}
	var commands []*Command
	for _, command := range c.commands {
		if cy == command {
			command.parent = nil
			break
		}
		commands = append(commands, command)
	}
	c.commands = commands

	//compute length
	c.maxUsageLen = 0
	c.maxPathLen = 0
	c.maxNameLen = 0
	for _, command := range c.commands {
		usageLen := len(command.Name())
		if usageLen > c.maxUsageLen {
			c.maxUsageLen = usageLen
		}
		pathLen := len(command.CommandPath())
		if pathLen > c.maxPathLen {
			c.maxPathLen = pathLen
		}
		nameLen := len(command.Name())
		if nameLen > c.maxNameLen {
			c.maxNameLen = nameLen
		}
	}
	return nil
}

// CommandPath returns the full path to the command including parent commands, separated by the defined path separator.
func (c *Command) CommandPath() string {
	if !c.HasParent() {
		return interfaces.PathSeparator
	}
	parentPath := c.Parent().CommandPath()
	if !c.Parent().HasParent() {
		return parentPath + c.Name()
	}
	return parentPath + interfaces.PathSeparator + c.Name()
}

// Path returns the hierarchical path of command names, starting from the root, representing the command structure.
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

// HasAlias checks if the given string matches any alias of the command. Returns true if a match is found.
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

// NameAndAliases returns a comma-separated string of the command name followed by its aliases.
func (c *Command) NameAndAliases() string {
	return strings.Join(append([]string{c.Name()}, c.aliases...), ", ")
}

// IsAdditionalHelpTopicCommand checks if the command and its subcommands are categorized as 'help' commands.
func (c *Command) IsAdditionalHelpTopicCommand() bool {
	// if any non-help sub commands are found, the command is not a 'help' command
	for _, sub := range c.commands {
		if !sub.IsAdditionalHelpTopicCommand() {
			return false
		}
	}
	return true
}

// HasHelpSubCommands checks if the command has any subcommands that are categorized as additional help topics and returns true if so.
func (c *Command) HasHelpSubCommands() bool {
	for _, sub := range c.commands {
		if sub.IsAdditionalHelpTopicCommand() {
			return true
		}
	}
	return false
}

// HasParent checks if the command has a parent command and returns true if a parent exists, otherwise false.
func (c *Command) HasParent() bool {
	return c.parent != nil
}

// SuggestionsFor generates a sorted list of suggestions based on the provided prefix.
// It uses exact prefix matching and Levenshtein distance for approximate matches.
// Suggestions include command names, properties, and other registered items.
func (c *Command) SuggestionsFor(prefix string) []string {
	const levenshteinThreshold = 2
	const levenshteinMin = 2

	if strings.Contains(prefix, " ") {
		return nil
	}
	suggestionsMap := make(map[string]bool)
	itemToComplete := prefix
	itemToCompleteLower := strings.ToLower(itemToComplete)
	var items []string
	for _, cmd := range c.commands {
		items = append(items, cmd.Name())
	}
	for _, item := range items {
		if strings.HasPrefix(strings.ToLower(item), itemToCompleteLower) {
			suggestionsMap[item] = true
		}
	}
	if len(suggestionsMap) == 0 && len(itemToComplete) >= levenshteinMin {
		for _, item := range items {
			if len(item) >= levenshteinMin {
				if d := levenshteinDistance(itemToComplete, item, true); d <= levenshteinThreshold {
					suggestionsMap[item] = true
				}
			}
		}
	}
	ret := make([]string, 0, len(suggestionsMap))
	for suggestion := range suggestionsMap {
		ret = append(ret, suggestion)
	}
	sort.Strings(ret)
	return ret
}
