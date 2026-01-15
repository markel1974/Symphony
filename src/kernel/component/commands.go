package component

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// Command represents a structure that encapsulates an executable command, its metadata, and its expected inputs and outputs.
type Command struct {
	id          string
	description string
	command     interface{}
	argsType    []reflect.Type
	ret         []reflect.Kind
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
		description: desc,
		command:     command,
		exec:        reflect.ValueOf(command),
	}
	commandType := reflect.TypeOf(command)
	if commandType.Kind() != reflect.Func {
		panic(fmt.Errorf("%s: %s", id, "command isn't a function"))
	}
	//if commandType.NumOut() != 1 {
	//	panic(fmt.Errorf("%s: %s", id, "command must return a single value"))
	//}

	var retHelp []string
	for x := 0; x < commandType.NumOut(); x++ {
		cOut := commandType.Out(x).Kind()
		cmd.ret = append(cmd.ret, cOut)
		retHelp = append(retHelp, cOut.String())
	}
	var argsHelp []string
	for x := 0; x < commandType.NumIn(); x++ {
		cInType := commandType.In(x)
		cmd.argsType = append(cmd.argsType, cInType)
		argsHelp = append(argsHelp, cInType.String())
	}
	const sep = ", "
	cmd.signature = cmd.id + "(" + strings.Join(argsHelp, sep) + ")"
	if len(retHelp) == 0 {
		//nothing to do
	} else if len(retHelp) == 1 {
		cmd.signature += " " + strings.Join(retHelp, sep)
	} else {
		cmd.signature += " (" + strings.Join(retHelp, sep) + ")"
	}
	return cmd
}

// Id returns the unique identifier of the Command instance.
func (cmd *Command) Id() string {
	return cmd.id
}

// Description returns the description of the Command instance.
func (cmd *Command) Description() string {
	return cmd.description

}

// Exec invokes the encapsulated command with the provided arguments and returns the result or an error on failure.
func (cmd *Command) Exec(args []string) (interface{}, error) {
	if len(args) != len(cmd.argsType) {
		return nil, fmt.Errorf("wrong number of arguments")
	}
	var rArgs []reflect.Value
	for x := 0; x < len(cmd.argsType); x++ {
		val, ok := ConvertArgument(args[x], cmd.argsType[x])
		if !ok {
			return nil, fmt.Errorf("can't convert argument %d ('%s'): %v", x, args[x], cmd.argsType[x])
		}
		rArgs = append(rArgs, val)
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

// CreateShellCommand converts the Command instance into a shell-compatible command with execution and help functionality.
func (cmd *Command) CreateShellCommand() *process.Command {
	cmdExec := func(task interfaces.IUserProcess, args []string) error {
		v, err := cmd.Exec(args)
		if v != nil {
			task.Write(fmt.Sprint(v), true)
		}
		return err
	}
	childCmd := process.NewCommand(cmd.Id(), interfaces.CommandTypeFile, nil, false, cmdExec)
	childCmd.SetHelp(cmd.Description(), cmd.Description())
	return childCmd
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
func (c *Commands) Add(command *Command) error {
	if _, ok := c.commands[command.Id()]; ok {
		return fmt.Errorf("command '%s' already exists", command.Id())
	}
	c.commands[command.Id()] = command
	return nil
}

// Exists checks if a command with the specified Id exists in the collection. Returns true if it exists, otherwise false.
func (c *Commands) Exists(id string) bool {
	_, ok := c.commands[id]
	return ok
}

// Remove deletes a command from the commands map by its specified id.
func (c *Commands) Remove(id string) {
	delete(c.commands, id)
}

// Retrieve returns the Command associated with the specified id from the Commands collection. Returns nil if not found.
func (c *Commands) Retrieve(id string) *Command {
	return c.commands[id]
}

// Exec executes the command identified by the given id with the provided arguments and returns the result or an error.
func (c *Commands) Exec(id string, args []string) (interface{}, error) {
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

// List returns a sorted list of command signatures present in the Commands collection.
func (c *Commands) List() map[string]string {
	out := make(map[string]string)
	for _, v := range c.commands {
		out[v.signature] = v.description
	}
	return out
}

// CreateShellCommands generates and returns a list of shell-compatible commands based on the commands in the collection.
func (c *Commands) CreateShellCommands() []*process.Command {
	var out []*process.Command
	for _, cmd := range c.commands {
		childCmd := cmd.CreateShellCommand()
		out = append(out, childCmd)
	}
	return out
}
