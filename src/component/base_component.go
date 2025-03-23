package component

import (
	"fmt"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/shell/cli"
	"io"
	"strconv"
	"strings"
)

// propertiesId is a constant string used as a key to reference component properties in state maps.
const propertiesId = "properties"

// childrenId is a constant key used to reference or identify child components in a hierarchical structure.
const childrenId = "children"

// detailsId is a constant used as a key to store or retrieve details of a component in a state map.
const detailsId = "details"

// BaseComponent provides a base implementation for composed components with properties, commands, and hierarchical structure.
type BaseComponent struct {
	cmd        *cli.Command
	id         string
	name       string
	instance   int
	kind       string
	factory    references.IComponentFactory
	node       references.INode
	properties *Properties
	commands   *Commands
	propagate  bool
}

// NewBaseComponent creates a new instance of BaseComponent with a unique ID, name, label, and initialized properties and commands.
func NewBaseComponent() *BaseComponent {
	bc := &BaseComponent{
		id:         "",
		factory:    nil,
		name:       "",
		instance:   0,
		kind:       "",
		properties: NewProperties(),
		commands:   NewCommands(),
		propagate:  true,
		cmd:        nil,
	}
	return bc
}

func (bc *BaseComponent) Register(f references.IComponentFactory, parent references.IComponent, name string, instance int, component references.IComponent, kind string) {
	bc.factory = f
	bc.name = name
	bc.instance = instance
	bc.kind = kind
	bc.id = bc.name + ":" + strconv.Itoa(instance) + ":" + bc.kind

	bc.cmd = cli.NewCommand()
	bc.cmd.Run = func(cmd *cli.Command, pid int, args []string) {
		_, _ = cmd.WriteLn([]byte{})
		_, _ = cmd.WriteLn([]byte(bc.name))
	}
	bc.cmd.Use = bc.name
	bc.cmd.Short = "Command " + bc.name
	bc.cmd.Long = "This is a command"

	if parent != nil && parent.GetNode() != nil {
		pNode := parent.GetNode()
		bc.node = newNode(pNode, component)
		pNode.AddComponent(component)
		_ = parent.GetCommand().AddCommand(bc.cmd)
	} else {
		bc.node = newNode(nil, component)
	}
}

func (bc *BaseComponent) DisablePropagate() {
	bc.propagate = false
}

func (bc *BaseComponent) Propagate() bool {
	return bc.propagate
}

func (bc *BaseComponent) Kind() string {
	return bc.kind
}

// GetFactory retrieves the INode instance associated with the BaseComponent.
func (bc *BaseComponent) GetFactory() references.IComponentFactory {
	return bc.factory
}

// GetNode retrieves the INode instance associated with the BaseComponent.
func (bc *BaseComponent) GetNode() references.INode {
	return bc.node
}

// Unregister removes the specified component from the parent's node, detaching their hierarchical relationship.
func (bc *BaseComponent) Unregister(component references.IComponent) {
	bc.node.RemoveComponent(component)
}

// GetId returns the unique identifier of the BaseComponent instance.
func (bc *BaseComponent) GetId() string {
	return bc.id
}

// GetChildren retrieves all child components of the current BaseComponent by traversing its associated node hierarchy.
func (bc *BaseComponent) GetChildren() []references.IComponent {
	var children []references.IComponent
	for _, child := range bc.node.GetChildren() {
		children = append(children, child.GetComponent())
	}
	return children
}

// GetChild retrieves a child component by its unique identifier. Returns nil if no child with the specified ID exists.
func (bc *BaseComponent) GetChild(id string) references.IComponent {
	for _, child := range bc.node.GetChildren() {
		if child.GetComponent().GetId() == id {
			return child.GetComponent()
		}
	}
	return nil
}

// GetComponentPath navigates the node structure using the specified path and returns the associated component, or nil if not found.
func (bc *BaseComponent) GetComponentPath(path string) references.IComponent {
	component := bc.node.Traverse(path)
	return component
}

// AddProperty registers a new property to the BaseComponent with the specified ID, description, read-only flag, getter, and setter.
func (bc *BaseComponent) AddProperty(id string, desc string, ro bool, get interface{}, set interface{}) {
	p := NewPropertyInfo(id, desc, ro, get, set)
	bc.properties.Add(p)
}

// GetProperty retrieves the value of the specified property by its name. It returns the value and an error, if any occur.
func (bc *BaseComponent) GetProperty(prop string) (interface{}, error) {
	v, err := bc.properties.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetPropertyPath retrieves the value of a specified property from a component located at a given path in the node hierarchy.
// Returns the property value and an error if the path or property is invalid.
func (bc *BaseComponent) GetPropertyPath(path string, prop string) (interface{}, error) {
	component := bc.node.Traverse(path)
	if component == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	v, err := component.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// SetProperty updates the value of a specified property if it exists. Returns an error if the update fails.
func (bc *BaseComponent) SetProperty(prop string, value interface{}) error {
	if err := bc.properties.SetProperty(prop, value); err != nil {
		return err
	}
	return nil
}

// SetPropertyPath sets a specific property value of a component identified by its path. Returns an error if the component is not found or setting the property fails.
func (bc *BaseComponent) SetPropertyPath(path string, prop string, val interface{}) error {
	component := bc.node.Traverse(path)
	if component == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := component.SetProperty(prop, val)
	return err
}

// Dump retrieves and returns the complete state of the component's properties as a map. Returns an error if dumping fails.
func (bc *BaseComponent) Dump() (map[string]interface{}, error) {
	state, err := bc.properties.Dump()
	if err != nil {
		return state, err
	}
	return state, nil
}

// DumpPath retrieves and returns the state of a component at a specified path in the node hierarchy. Returns an error if not found.
func (bc *BaseComponent) DumpPath(path string) (map[string]interface{}, error) {
	component := bc.node.Traverse(path)
	if component == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	d, err := component.Dump()
	return d, err
}

// DumpAll traverses the entire component tree, dumping the state of each component into a nested map structure. It returns the map and any error that occurs.
func (bc *BaseComponent) DumpAll() (map[string]interface{}, error) {
	rootMap := make(map[string]interface{})
	if err := bc.dump(bc.node.GetComponent(), rootMap); err != nil {
		return nil, err
	}
	return rootMap, nil
}

// Restore reinstates the internal state of the BaseComponent using the provided state map. Returns an error if unsuccessful.
func (bc *BaseComponent) Restore(state map[string]interface{}) error {
	if err := bc.properties.Restore(state); err != nil {
		return err
	}
	return nil
}

// RestorePath restores the state of a component found at the given path using the provided data map.
// Returns an error if the component is not found or if the restoration process fails.
func (bc *BaseComponent) RestorePath(path string, d map[string]interface{}) error {
	component := bc.node.Traverse(path)
	if component == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := component.Restore(d)
	if err != nil {
		return err
	}
	return nil
}

// RestoreAll restores the entire state of all components in the hierarchy from the given state map. Returns an error if failed.
func (bc *BaseComponent) RestoreAll(state map[string]interface{}) error {
	if err := bc.restore(bc.node.GetComponent(), state); err != nil {
		return err
	}
	return nil
}

func (bc *BaseComponent) GetCommand() *cli.Command {
	return bc.cmd
}

// CommandAdd registers a new command with an identifier, description, and implementation; returns an error if the id exists.
func (bc *BaseComponent) CommandAdd(id string, desc string, cmd interface{}) error {
	return bc.commands.Add(id, desc, cmd)
}

// CommandExec executes the specified command with the given arguments and returns the result or an error if execution fails.
func (bc *BaseComponent) CommandExec(cmd string, args ...interface{}) (interface{}, error) {
	d, err := bc.commands.Exec(cmd, args)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CommandExecPath executes a command on a component located at the specified path, passing optional arguments to it.
func (bc *BaseComponent) CommandExecPath(path string, cmd string, args ...interface{}) (interface{}, error) {
	component := bc.node.Traverse(path)
	if component == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	d, err := component.CommandExec(cmd, args...)
	return d, err
}

// CommandDocumentation recursively collects command documentation for the current component and its children.
func (bc *BaseComponent) CommandDocumentation(data map[string]interface{}) {
	data[bc.GetId()] = bc.commands.Documentation()
	for _, child := range bc.node.GetChildren() {
		child.GetComponent().CommandDocumentation(data)
	}
}

// Print generates a textual representation of the component and its children, writing it to the provided writer.
// w is the io.Writer where output is written, indent specifies the indentation level, and showComponents toggles type display.
func (bc *BaseComponent) Print(w io.Writer, indent string, showComponents bool) {
	_, _ = fmt.Fprintf(w, "%s%s", indent, bc.GetId())
	if showComponents {
		_, _ = fmt.Fprintf(w, " (%T)", bc.node.GetComponent())
	}
	_, _ = fmt.Fprintln(w)
	for _, child := range bc.node.GetChildren() {
		child.GetComponent().Print(w, indent+"  ", showComponents)
	}
}

// dump recursively extracts and organizes the state of a component and its children into a hierarchical map structure.
// It handles component properties and children relationships, and returns an error if a component state fails to dump.
func (bc *BaseComponent) dump(component references.IComponent, state map[string]interface{}) error {
	if component == nil {
		return nil
	}
	id := component.GetId()
	if len(id) == 0 {
		return nil
	}
	componentState, err := component.Dump()
	if err != nil {
		return err
	}
	componentSegment := make(map[string]interface{})
	state[id] = componentSegment
	componentSegment[detailsId] = map[string]interface{}{
		"id": component.GetId(),
	}
	if len(componentState) > 0 {
		componentSegment[propertiesId] = componentState
	}
	if children := component.GetChildren(); len(children) > 0 {
		childrenSegment := make(map[string]interface{})
		componentSegment[childrenId] = childrenSegment
		for _, child := range children {
			if err = bc.dump(child, childrenSegment); err != nil {
				return err
			}
		}
	}
	return nil
}

// restore restores the state of the provided component from the given state map recursively, handling children and properties.
// Returns an error if the restoration process encounters any issues.
func (bc *BaseComponent) restore(component references.IComponent, state map[string]interface{}) error {
	if component == nil {
		return nil
	}
	id := component.GetId()
	if len(id) == 0 {
		return nil
	}
	componentI, ok := state[id]
	if !ok {
		return nil
	}
	componentSegment, ok := componentI.(map[string]interface{})
	if !ok || componentSegment == nil {
		return fmt.Errorf("error restoring component %s: %s", id, "missing component node")
	}
	detailsI, ok := componentSegment[detailsId]
	if !ok {
		return fmt.Errorf("error restoring component %s: %s", id, "missing detail node")
	}
	detailsSegment, ok := detailsI.(map[string]interface{})
	if !ok || detailsSegment == nil {
		return fmt.Errorf("error restoring component %s: %s", id, "unknown details node")
	}
	//TODO Details Handler
	if propertiesI, ok := componentSegment[propertiesId]; ok {
		if propertiesSegment, ok := propertiesI.(map[string]interface{}); ok && len(propertiesSegment) > 0 {
			if err := component.Restore(propertiesSegment); err != nil {
				return fmt.Errorf("error restoring component %s: %w", id, err)
			}
		}
	}
	if childrenI, ok := componentSegment[childrenId]; ok {
		childrenSegment, ok := childrenI.(map[string]interface{})
		if !ok {
			return fmt.Errorf("error restoring component %s: %s", id, "unknown children node")
		}
		for _, child := range component.GetChildren() {
			if err := bc.restore(child, childrenSegment); err != nil {
				return err
			}
		}
	}
	return nil
}

// Restore reconstructs an IComponent with its hierarchical state using a factory, parentComponent, and a state map.
// It returns the restored IComponent or an error if the process fails.
func Restore(factory references.IComponentFactory, parentComponent references.IComponent, component references.IComponent, state map[string]interface{}) (references.IComponent, error) {
	root, err := _restore(factory, parentComponent, component, state)
	if err != nil {
		return nil, err
	}
	return root, nil
}

// _restore reconstructs an IComponent and its hierarchy from a serialized state using the provided component factory.
// It initializes a component if nil, restores its properties and child components, and handles state validation.
// Returns the restored component or any error encountered during the restoration process.
func _restore(factory references.IComponentFactory, parentComponent references.IComponent, component references.IComponent, s map[string]interface{}) (references.IComponent, error) {
	if component == nil {
		keys, err := GetSegmentKeys(s)
		if err != nil {
			return nil, err
		}
		if len(keys) != 1 {
			return nil, fmt.Errorf("error restoring component: %s", "invalid key")
		}
		id := keys[0]
		var instance int
		if pos := strings.LastIndex(id, ":"); pos > 0 {
			v := id[pos+1:]
			if instance, err = strconv.Atoi(v); err != nil {
				return nil, fmt.Errorf("error restoring component: %s", err.Error())
			}
			id = id[:pos]
		}
		if component, err = factory.Create(parentComponent, id, instance); err != nil {
			return nil, err
		}
	}
	id := component.GetId()
	if len(id) == 0 {
		return nil, nil
	}
	stateI := s[id]
	if stateI == nil {
		return nil, fmt.Errorf("error restoring component %s: %s", id, "missing component node")
	}
	//TODO Details Handler
	//detailsSegment, err := GetSegment(detailsId, stateI)
	//if err != nil {
	//	return nil, err
	//}
	if propertiesSegment, _ := GetSegment(propertiesId, stateI); len(propertiesSegment) > 0 {
		if err := component.Restore(propertiesSegment); err != nil {
			return nil, fmt.Errorf("error restoring component %s: %w", id, err)
		}
	}

	if !component.Propagate() {
		return component, nil
	}

	if childrenSegment, _ := GetSegment(childrenId, stateI); len(childrenSegment) > 0 {
		for k, childI := range childrenSegment {
			if childI == nil {
				continue
			}
			child, ok := childI.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("error restoring component %s: %s", id, "unknown children node")
			}
			fullChild := map[string]interface{}{k: child}
			c := component.GetChild(k)
			if _, err := _restore(factory, component, c, fullChild); err != nil {
				return nil, err
			}
		}
	}
	return component, nil
}
