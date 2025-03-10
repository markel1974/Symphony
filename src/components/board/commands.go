package board

import (
	"fmt"
	"reflect"
	"strings"
)

// Command represents a structure that encapsulates an executable command, its metadata, and its expected inputs and outputs.
type Command struct {
	id           string
	description  string
	command      interface{}
	args         []reflect.Kind
	argsHelp     []string
	retValue     reflect.Kind
	retValueHelp string
	exec         reflect.Value
	signature    string
}

// NewCommand creates a new Command instance with the provided id, description, and function reference.
// It validates that the command is a function, accepts arguments, and returns a single value.
// Panics occur for invalid inputs such as nil commands, non-function references, or functions with invalid signatures.
func NewCommand(id string, desc string, command interface{}) *Command {
	if command == nil {
		panic(fmt.Errorf("%s: %s", id, "command is nil"))
	}
	cmd := &Command{
		id:          id,
		argsHelp:    []string{},
		description: desc,
		command:     command,
		exec:        reflect.ValueOf(command),
	}
	commandType := reflect.TypeOf(command)
	if commandType.Kind() != reflect.Func {
		panic(fmt.Errorf("%s: %s", id, "command isn't a function"))
	}
	if commandType.NumOut() != 1 {
		panic(fmt.Errorf("%s: %s", id, "command must return a single value"))
	}
	cOut := commandType.Out(0).Kind()
	cmd.retValue = cOut
	cmd.retValueHelp = cOut.String()
	for x := 0; x < commandType.NumIn(); x++ {
		cIn := commandType.In(x).Kind()
		cmd.args = append(cmd.args, cIn)
		cmd.argsHelp = append(cmd.argsHelp, cIn.String())
	}
	cmd.signature = cmd.id + "(" + strings.Join(cmd.argsHelp, ", ") + ") " + cmd.retValueHelp
	return cmd
}

// Exec invokes the encapsulated command with the provided arguments and returns the result or an error on failure.
func (cmd *Command) Exec(args []interface{}) (interface{}, error) {
	if len(args) != len(cmd.args) {
		return nil, fmt.Errorf("wrong number of arguments")
	}
	var rArgs []reflect.Value
	for x := 0; x < len(cmd.args); x++ {
		t := reflect.TypeOf(args[x])
		if t.Kind() != cmd.args[x] {
			return nil, fmt.Errorf("wrong argument type")
		}
		rArgs = append(rArgs, reflect.ValueOf(args[x]))
	}
	results := cmd.exec.Call(rArgs)
	if len(results) != 1 {
		return nil, fmt.Errorf("wrong number of results")
	}
	result := results[0].Interface()
	return result, nil
}

// Commands is a collection that manages a set of Command instances identified by unique string IDs.
type Commands struct {
	commands map[string]*Command
}

// NewCommands initializes and returns a new Commands instance with an empty map for storing command definitions.
func NewCommands() *Commands {
	return &Commands{
		commands: make(map[string]*Command),
	}
}

// Add registers a new command with a unique ID, description, and executable function. Panics if the ID already exists.
func (c *Commands) Add(id string, desc string, command interface{}) {
	if _, ok := c.commands[id]; ok {
		panic(fmt.Errorf("command '%s' already exists", id))
	}
	c.commands[id] = NewCommand(id, desc, command)
}

// Remove deletes a command from the commands map by its specified id.
func (c *Commands) Remove(id string) {
	delete(c.commands, id)
}

// Exec executes the command identified by the given id with the provided arguments and returns the result or an error.
func (c *Commands) Exec(id string, args []interface{}) (interface{}, error) {
	cmd, ok := c.commands[id]
	if !ok {
		return nil, fmt.Errorf("command '%s' not found", id)
	}
	v, err := cmd.Exec(args)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// PrintHelp returns a list of strings describing available commands, including their signature and description.
func (c *Commands) PrintHelp() []string {
	var out []string
	for _, v := range c.commands {
		out = append(out, v.signature+" : "+v.description)
	}
	return out
}
