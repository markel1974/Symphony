package board

import (
	"fmt"
	"io"
	"strings"
)

// IHardware defines an interface for hardware components with a Reset method.
type IHardware interface {
	Reset()
}

// INavigate provides an interface for hierarchical navigation and manipulation of node structures and properties.
type INavigate interface {
	GetId() string

	GetNode() *Node

	SetNode(node *Node)

	GetChildren() []IComponent

	GetComponentPath(string) IComponent

	GetProperty(string) (interface{}, error)

	GetPropertyPath(string, string) (interface{}, error)

	SetProperty(string, interface{}) error

	SetPropertyPath(string, string, interface{}) error

	Dump() (map[string]interface{}, error)

	DumpPath(string) (map[string]interface{}, error)

	DumpAll() (map[string]interface{}, error)

	Restore(map[string]interface{}) error

	RestoreAll(map[string]interface{}) error

	RestorePath(string, map[string]interface{}) error

	RunCommand(string, []string) (map[string]interface{}, error)

	RunCommandPath(string, string, string) (map[string]interface{}, error)

	Path() string

	Print(io.Writer, string, bool)
}

// IComponent represents a composite interface that combines IHardware and INavigate capabilities.
type IComponent interface {
	IHardware
	INavigate
}

// BaseComponent represents a foundational structure that implements IComponent, encapsulating an id, properties, and hierarchy.
type BaseComponent struct {
	id         string
	node       *Node
	properties *Properties
}

// Register associates the current BaseComponent instance with a specified parent Node, creating a new node if needed.
func Register(parent IComponent, component IComponent) {
	var parentNode *Node = nil
	if parent != nil {
		parentNode = parent.GetNode()
	}
	if parentNode != nil {
		component.SetNode(newNode(parentNode, component))
		parentNode.AddComponent(component)
	} else {
		component.SetNode(newNode(nil, component))
	}
}

// NewBaseComponent creates and returns a new instance of BaseComponent with the specified name, suffix, and RunFn.
// The id is constructed by concatenating the name and suffix with an underscore if the suffix is non-empty.
// Initializes the BaseComponent with a properties map using the provided RunFn.
func NewBaseComponent(name string, suffix string, runFn RunFn) *BaseComponent {
	id := name
	if len(suffix) > 0 {
		id += "_" + suffix
	}
	bc := &BaseComponent{
		id:         id,
		properties: NewProperties(runFn),
	}
	return bc
}

func (bc *BaseComponent) GetNode() *Node {
	return bc.node
}

func (bc *BaseComponent) SetNode(node *Node) {
	bc.node = node
}

// GetId returns the unique identifier of the BaseComponent instance.
func (bc *BaseComponent) GetId() string {
	return bc.id
}

func (bc *BaseComponent) GetChildren() []IComponent {
	var children []IComponent
	for _, child := range bc.node.GetChildren() {
		children = append(children, child.GetComponent())
	}
	return children
}

func (bc *BaseComponent) GetComponentPath(path string) IComponent {
	node := bc.node.FindNode(path)
	if node == nil {
		return nil
	}
	return node.GetComponent()
}

// AddProperty adds a new property to the BaseComponent with the specified ID, description, read-only flag, getter, and setter.
func (bc *BaseComponent) AddProperty(id string, desc string, ro bool, get interface{}, set interface{}) {
	p := CreatePropertyInfo(id, desc, ro, get, set)
	bc.properties.Add(p)
}

// GetProperty retrieves the value of a specified property by its identifier. Returns the property value or an error.
func (bc *BaseComponent) GetProperty(prop string) (interface{}, error) {
	v, err := bc.properties.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetPropertyPath retrieves the value of a property from a node identified by a path within the component hierarchy.
func (bc *BaseComponent) GetPropertyPath(path string, prop string) (interface{}, error) {
	node := bc.node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	v, err := node.GetComponent().GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// SetProperty sets the value of the specified property in the component. Returns an error if the operation fails.
func (bc *BaseComponent) SetProperty(prop string, value interface{}) error {
	if err := bc.properties.SetProperty(prop, value); err != nil {
		return err
	}
	return nil
}

// SetPropertyPath sets a property value for a component identified by a path. Returns an error if the node or property is not found.
func (bc *BaseComponent) SetPropertyPath(path string, prop string, val interface{}) error {
	node := bc.node.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.GetComponent().SetProperty(prop, val)
	return err
}

// Dump returns a map representing the current state of the component's properties or an error if retrieval fails.
func (bc *BaseComponent) Dump() (map[string]interface{}, error) {
	state, err := bc.properties.Dump()
	if err != nil {
		return state, err
	}
	return state, nil
}

// DumpPath retrieves and dumps the state of a component located at the specified path, or returns an error if not found.
func (bc *BaseComponent) DumpPath(path string) (map[string]interface{}, error) {
	node := bc.node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	d, err := node.GetComponent().Dump()
	return d, err
}

// DumpAll collects and returns the state of the current component and all its children in a nested map structure.
func (bc *BaseComponent) DumpAll() (map[string]interface{}, error) {
	rootMap := make(map[string]interface{})
	if err := bc.dump(bc.node.GetComponent(), rootMap); err != nil {
		return nil, err
	}
	return rootMap, nil
}

// Restore restores the state of the BaseComponent using the provided map. It returns an error if the restoration fails.
func (bc *BaseComponent) Restore(state map[string]interface{}) error {
	if err := bc.properties.Restore(state); err != nil {
		return err
	}
	return nil
}

// RestorePath restores the state of a component identified by the given path using the provided data map.
func (bc *BaseComponent) RestorePath(path string, d map[string]interface{}) error {
	node := bc.node.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.GetComponent().Restore(d)
	if err != nil {
		return err
	}
	return nil
}

// RestoreAll recursively restores the state of all components in the hierarchy from the provided state map.
func (bc *BaseComponent) RestoreAll(state map[string]interface{}) error {
	if err := bc.restore(bc.node.GetComponent(), state); err != nil {
		return err
	}
	return nil
}

// RunCommand executes a specified command with provided arguments using the component's properties and returns the result or an error.
func (bc *BaseComponent) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	d, err := bc.properties.Run(cmd, args)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// RunCommandPath executes a command on a specified path node with provided arguments and returns the result or an error.
func (bc *BaseComponent) RunCommandPath(path string, cmd string, args string) (map[string]interface{}, error) {
	node := bc.node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	values := strings.Split(args, " ")
	d, err := node.GetComponent().RunCommand(cmd, values)
	return d, err
}

// Path constructs and returns the hierarchical path of the component based on its parent's path and its own id.
func (bc *BaseComponent) Path() string {
	id := bc.GetId()
	parent := bc.node.GetParent()
	if parent == nil || id == "" {
		return id
	}
	return parent.GetComponent().Path() + "." + id
}

// Print writes the component's ID to the provided writer with optional formatting for indent and child components.
func (bc *BaseComponent) Print(w io.Writer, indent string, showComponents bool) {
	_, _ = fmt.Fprintf(w, "%s%s", indent, bc.GetId())
	if showComponents {
		_, _ = fmt.Fprintf(w, " (%T)", bc)
	}
	_, _ = fmt.Fprintln(w)
	for _, child := range bc.node.GetChildren() {
		child.GetComponent().Print(w, indent+"  ", showComponents)
	}
}

// dump recursively processes nodes and their components, creating a map with the component's ID as the key and its data as the value.
// Returns an error if the dump operation for any component fails.
func (bc *BaseComponent) dump(component IComponent, stateMap map[string]interface{}) error {
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
	stateMap[id] = componentState
	if children := component.GetChildren(); len(children) > 0 {
		if componentState == nil {
			componentState = make(map[string]interface{})
		}
		for _, child := range children {
			if err = bc.dump(child, componentState); err != nil {
				return err
			}
		}
	}
	return nil
}

// restore traverses the node tree recursively, restoring the state of each component from the provided state map.
// Returns an error if the restoration of a component's state fails or if any issues occur during traversal.
func (bc *BaseComponent) restore(component IComponent, state map[string]interface{}) error {
	if component == nil {
		return nil
	}
	id := component.GetId()
	if len(id) == 0 {
		return nil
	}
	componentState, ok := state[id].(map[string]interface{})
	if ok {
		if err := component.Restore(componentState); err != nil {
			return fmt.Errorf("error restoring component %s: %w", id, err)
		}
	}
	for _, child := range component.GetChildren() {
		if err := bc.restore(child, state); err != nil {
			return err
		}
	}
	return nil
}
