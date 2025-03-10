package board

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Command represents a structure that encapsulates an executable command, its metadata, and its expected inputs and outputs.
type Command struct {
	id          string
	description string
	command     interface{}
	args        []reflect.Kind
	argsHelp    []string
	ret         []reflect.Kind
	retHelp     []string
	exec        reflect.Value
	signature   string
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

	for x := 0; x < commandType.NumOut(); x++ {
		cOut := commandType.Out(x).Kind()
		cmd.ret = append(cmd.ret, cOut)
		cmd.retHelp = append(cmd.retHelp, cOut.String())
	}

	for x := 0; x < commandType.NumIn(); x++ {
		cIn := commandType.In(x).Kind()
		cmd.args = append(cmd.args, cIn)
		cmd.argsHelp = append(cmd.argsHelp, cIn.String())
	}
	cmd.signature = cmd.id + "(" + strings.Join(cmd.argsHelp, ", ") + ")"
	if len(cmd.retHelp) == 0 {
	} else if len(cmd.retHelp) == 1 {
		cmd.signature += " " + strings.Join(cmd.retHelp, ", ")
	} else {
		cmd.signature += " (" + strings.Join(cmd.retHelp, ", ") + ")"
	}
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
	if len(results) != len(cmd.ret) {
		return nil, fmt.Errorf("wrong number of results")
	}
	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		return results[0].Interface(), nil
	} else {
		var result []interface{}
		for _, v := range results {
			result = append(result, v.Interface())
		}
		return result, nil
	}
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

// Add adds a new command to the Commands collection with a unique id, description, and function reference.
// Returns an error if a command with the same id already exists.
func (c *Commands) Add(id string, desc string, command interface{}) error {
	if _, ok := c.commands[id]; ok {
		return fmt.Errorf("command '%s' already exists", id)
	}
	c.commands[id] = NewCommand(id, desc, command)
	return nil
}

// Exists checks if a command with the specified ID exists in the collection. Returns true if it exists, otherwise false.
func (c *Commands) Exists(id string) bool {
	_, ok := c.commands[id]
	return ok
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

// Documentation returns a map where keys are command IDs and values are their corresponding signatures from the Commands collection.
func (c *Commands) Documentation() []string {
	var out []string
	for _, v := range c.commands {
		signature := v.signature + " : " + v.description
		out = append(out, signature)
	}
	sort.Strings(out)
	return out
}
