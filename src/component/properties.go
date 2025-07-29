package component

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/shell"
	"reflect"
	"sort"
	"strconv"
	"unicode"
)

// PropertyInfo represents metadata and behavior for a property, including its ID, type details, description, and functionality.
type PropertyInfo struct {
	id          string
	description string
	readOnly    bool
	getType     reflect.Type
	getValue    reflect.Value
	setType     reflect.Type
	setValue    reflect.Value
	set0Type    reflect.Type
}

// NewPropertyInfo initializes and returns a new PropertyInfo object with the given parameters.
// id specifies the unique identifier of the property.
// kind defines the data type of the property.
// desc provides a human-readable description of the property.
// ro indicates whether the property is read-only.
// get is a function to retrieve the property's value, with signature func() <ret>.
// set is a function to update the property's value, with signature func(v <arg>).
func NewPropertyInfo(id string, desc string, ro bool, get interface{}, set interface{}) *PropertyInfo {
	const getError = "wrong get signature must be func get() <ret>"
	const setError = "wrong set signature must be func set(v <arg>)error"
	setReference := reflect.TypeOf(func(int) error { return nil })
	getReference := reflect.TypeOf(func() int { return 0 })
	p := &PropertyInfo{id: id, description: desc, readOnly: ro}
	p.getType = reflect.TypeOf(get)
	p.getValue = reflect.ValueOf(get)
	p.setType = reflect.TypeOf(set)
	p.setValue = reflect.ValueOf(set)
	if get == nil || p.getType.Kind() != getReference.Kind() {
		panic(fmt.Errorf("%s: %s", id, getError))
	}
	if p.getType.NumIn() != getReference.NumIn() || p.getType.NumOut() != getReference.NumOut() {
		panic(fmt.Errorf("%s: %s", id, getError))
	}
	if set == nil || p.setType.Kind() != setReference.Kind() {
		panic(fmt.Errorf("%s: %s", id, setError))
	}
	if p.setType.NumIn() != setReference.NumIn() || p.setType.NumOut() != setReference.NumOut() {
		panic(fmt.Errorf("%s: %s", id, setError))
	}
	for x := 0; x < setReference.NumOut(); x++ {
		if setReference.Out(x) != p.setType.Out(x) {
			panic(fmt.Errorf("%s: %s", id, setError))
		}
	}
	p.set0Type = p.setType.In(0)
	return p
}

// Id returns the unique identifier of the PropertyInfo.
func (prop *PropertyInfo) Id() string {
	return prop.id
}

// Description returns the description of the property as a string.
func (prop *PropertyInfo) Description() string {
	return prop.description
}

// Set assigns a new value to the property if it is not read-only and the input type matches the expected kind.
func (prop *PropertyInfo) Set(arg interface{}) error {
	exec := func(argValue reflect.Value) error {
		if !argValue.Type().ConvertibleTo(prop.set0Type) {
			return fmt.Errorf("property '%s' expected type '%v', got '%T'", prop.id, prop.set0Type, arg)
		}
		convertedArg := argValue.Convert(prop.set0Type)
		results := prop.setValue.Call([]reflect.Value{convertedArg})
		if len(results) != 1 {
			return fmt.Errorf("property '%s' set failed", prop.id)
		}
		if res := results[0].Interface(); res == nil {
			return nil
		}
		err, ok := results[0].Interface().(error)
		if !ok {
			return fmt.Errorf("property '%s' set failed", prop.id)
		}
		return err
	}

	if prop.readOnly {
		return fmt.Errorf("property '%s' is read-only", prop.id)
	}
	argValue := reflect.ValueOf(arg)
	if argValue.Type().AssignableTo(prop.set0Type) {
		prop.setValue.Call([]reflect.Value{argValue})
		return nil
	}
	if err := exec(argValue); err == nil {
		return nil
	}
	//try explicit conversion
	switch argValue.Kind() {
	case reflect.String:
		if f, err := strconv.ParseFloat(argValue.Interface().(string), 64); err == nil {
			argValue = reflect.ValueOf(f)
			if err = exec(argValue); err != nil {
				return err
			}
			return nil
		}
	default:
		//nothing to do
	}

	return fmt.Errorf("property '%s' expected type '%v', got '%T'", prop.id, prop.set0Type, arg)
}

// Get retrieves the value of the property by calling its getter function and returns the result or an error if any occurs.
func (prop *PropertyInfo) Get() (interface{}, error) {
	results := prop.getValue.Call([]reflect.Value{})
	if len(results) == 0 {
		return nil, fmt.Errorf("property '%s' no results returned", prop.id)
	}
	result := results[0].Interface()
	return result, nil
}

// CreateShellSetCommand generates two shell commands for setting and getting the property's value.
func (prop *PropertyInfo) CreateShellSetCommand() *shell.Command {
	id := []rune(prop.Id())
	if len(id) > 0 {
		id[0] = unicode.ToUpper(id[0])
	}
	setProp := func(task interfaces.ITask, args []string) error {
		if len(args) == 0 {
			task.WriteLn("no argument provided")
			return nil
		}
		if err := prop.Set(args[0]); err != nil {
			task.WriteLn(err.Error())
		}
		return nil
	}
	childSet := shell.NewCommand("set"+string(id), interfaces.CommandTypeFile, nil, false, setProp)
	childSet.SetHelp(prop.Description(), prop.Description())
	return childSet
}

// CreateShellGetCommand generates two shell commands for setting and getting the property's value.
func (prop *PropertyInfo) CreateShellGetCommand() *shell.Command {
	id := []rune(prop.Id())
	if len(id) > 0 {
		id[0] = unicode.ToUpper(id[0])
	}
	getProp := func(task interfaces.ITask, args []string) error {
		v, err := prop.Get()
		task.WriteLn(fmt.Sprint(v))
		return err
	}
	childGet := shell.NewCommand("get"+string(id), interfaces.CommandTypeFile, nil, false, getProp)
	childGet.SetHelp(prop.Description(), prop.Description())
	return childGet
}

// CreateShellCommand generates two shell commands for setting and getting the property's value.
func (prop *PropertyInfo) CreateShellCommand() []*shell.Command {
	var out []*shell.Command
	out = append(out, prop.CreateShellGetCommand())
	out = append(out, prop.CreateShellSetCommand())
	return out
}

// Properties represent a collection of property definitions mapped by identifiers and a function to execute commands.
type Properties struct {
	properties map[string]*PropertyInfo
}

// NewProperties creates a new instance of Properties and initializes an empty properties map.
func NewProperties() *Properties {
	return &Properties{
		properties: make(map[string]*PropertyInfo),
	}
}

// Add inserts a new PropertyInfo instance into the Properties map using its id as the key.
// It panics if the provided PropertyInfo argument is nil.
func (p *Properties) Add(prop *PropertyInfo) {
	if prop == nil {
		panic(fmt.Errorf("property info is nil"))
	}
	p.properties[prop.Id()] = prop
}

// Retrieve retrieves the PropertyInfo associated with the given id from the property map. Returns nil if not found.
func (p *Properties) Retrieve(id string) *PropertyInfo {
	return p.properties[id]
}

// List returns a slice of string identifiers for all properties in the Properties instance.
func (p *Properties) List() map[string]string {
	out := make(map[string]string)
	for _, v := range p.properties {
		out[v.id] = v.description
	}
	return out
}

// GetProperty retrieves the value of a property by its id from the Properties instance. Returns the value or an error if not found.
func (p *Properties) GetProperty(id string) (interface{}, error) {
	prop, ok := p.properties[id]
	if !ok || prop == nil {
		return nil, fmt.Errorf("property '%s' not found", id)
	}
	ret, err := prop.Get()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// SetProperty updates the value of a property with the given id. Returns an error if the property is not found or invalid.
func (p *Properties) SetProperty(id string, arg interface{}) error {
	prop, ok := p.properties[id]
	if !ok || prop == nil {
		return fmt.Errorf("property '%s' not found", id)
	}
	if err := prop.Set(arg); err != nil {
		return err
	}
	return nil
}

// Dump returns a map of all properties with their current values. Returns an error if any property's retrieval fails.
func (p *Properties) Dump() (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for _, prop := range p.properties {
		res, err := prop.Get()
		if err != nil {
			return nil, err
		}
		out[prop.Id()] = res
	}
	return out, nil
}

// Restore updates the properties in the receiver using the provided map. Returns an error if a property cannot be restored.
func (p *Properties) Restore(d map[string]interface{}) error {
	for k, v := range d {
		prop, ok := p.properties[k]
		if !ok || prop == nil {
			return fmt.Errorf("property '%s' not found", k)
		}
		if err := prop.Set(v); err != nil {
			return err
		}
	}
	return nil
}

// CreateShellDump generates a slice of shell commands to get and set property values in the Properties instance.
func (p *Properties) CreateShellDump(name string) *shell.Command {
	dumpFn := func(task interfaces.ITask, args []string) error {
		for _, prop := range p.sort() {
			if v, err := prop.Get(); err == nil {
				task.WriteLn(prop.Id() + ": " + fmt.Sprint(v))
			}
		}
		return nil
	}
	c := shell.NewCommand(name, interfaces.CommandTypeFile, nil, false, dumpFn)
	c.SetHelp("dump available properties", "dump available properties")
	return c
}

// sort orders the properties by their identifier and returns a slice of PropertyInfo sorted in ascending order.
func (p *Properties) sort() []*PropertyInfo {
	var props []*PropertyInfo
	for _, prop := range p.properties {
		props = append(props, prop)
	}
	sort.Slice(props, func(i, j int) bool {
		return props[i].Id() < props[j].Id()
	})
	return props
}
