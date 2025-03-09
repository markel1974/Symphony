package board

import (
	"fmt"
	"io"
	"strings"
)

// IComponent represents an interface for managing components within a hierarchical structure of nodes.
// GetId returns the unique identifier for the component.
// Reset initializes or resets the internal state of the component.
// SetNode associates the component with the provided node.
// GetNode retrieves the node currently associated with the component.
// GetProperty fetches the value of the specified property by its key.
// GetPropertyPath fetches the value of a property specified by a key in a node's path.
// SetProperty sets a property value identified by its key.
// SetPropertyPath sets a property value within the specified path.
// Dump returns a map representation of the component's current state.
// DumpPath retrieves the dumped state of the component at the specified path.
// DumpAll returns the state of the component and all related states in a hierarchical context.
// Restore restores the state of the component using the provided key-value map.
// RestoreAll restores the states of all components using the provided map.
// RestorePath restores the state of the component at a specific path using the provided map.
// RunCommand executes a command on the component with the given arguments.
// RunCommandPath executes a command on the component at a specific path with the given arguments.
// Path returns the hierarchical path of the component within its structure.
// Print outputs the component's details in human-readable form to the provided writer.
type IComponent interface {
	GetId() string

	Reset()

	SetNode(*Node)

	GetNode() *Node

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

// BaseComponent serves as a foundational structure for reusable components with unique identifiers and properties.
type BaseComponent struct {
	id         string
	node       *Node
	properties *Properties
}

// NewBaseComponent creates and initializes a new BaseComponent instance with a unique ID and associated properties.
// The ID is created by concatenating the given name and suffix.
// The provided RunFn is used for property operations.
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

// GetId returns the unique identifier of the BaseComponent.
func (b *BaseComponent) GetId() string {
	return b.id
}

// SetNode assigns the provided Node instance to the component, linking it within the node hierarchy.
func (b *BaseComponent) SetNode(n *Node) {
	b.node = n
}

// GetNode retrieves the Node associated with the BaseComponent.
func (b *BaseComponent) GetNode() *Node {
	return b.node
}

// AddProperty adds a new property to the BaseComponent with a unique id,
// description, read-only flag, and get/set functions.
func (b *BaseComponent) AddProperty(id string, desc string, ro bool, get interface{}, set interface{}) {
	p := CreatePropertyInfo(id, desc, ro, get, set)
	b.properties.Add(p)
}

// GetProperty retrieves the value of a property by its identifier. Returns the property value or an error if not found.
func (b *BaseComponent) GetProperty(prop string) (interface{}, error) {
	v, err := b.properties.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetPropertyPath retrieves the value of a specific property from a component located at a given path.
// Returns an error if the path or property is invalid or not found.
func (b *BaseComponent) GetPropertyPath(path string, prop string) (interface{}, error) {
	node := b.node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	v, err := node.GetComponent().GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// SetProperty updates the value of a specified property in the component. Returns an error if the operation fails.
func (b *BaseComponent) SetProperty(prop string, value interface{}) error {
	if err := b.properties.SetProperty(prop, value); err != nil {
		return err
	}
	return nil
}

// SetPropertyPath sets the specified property of a component identified by the given path to the provided value.
// It returns an error if the node at the given path is not found or if setting the property fails.
func (b *BaseComponent) SetPropertyPath(path string, prop string, val interface{}) error {
	node := b.node.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.GetComponent().SetProperty(prop, val)
	return err
}

// Dump retrieves the current state of the component's properties as a map. Returns an error if the operation fails.
func (b *BaseComponent) Dump() (map[string]interface{}, error) {
	state, err := b.properties.Dump()
	if err != nil {
		return state, err
	}
	return state, nil
}

// DumpPath retrieves and returns the state of the component located at the specified path as a map.
// Returns an error if the node is not found or if component dumping fails.
func (b *BaseComponent) DumpPath(path string) (map[string]interface{}, error) {
	node := b.node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	d, err := node.GetComponent().Dump()
	return d, err
}

// DumpAll recursively dumps the entire hierarchical state of all components into a map.
// Returns an error if dumping fails.
func (b *BaseComponent) DumpAll() (map[string]interface{}, error) {
	rootMap := make(map[string]interface{})
	if err := b.dump(b.node, rootMap); err != nil {
		return nil, err
	}
	return rootMap, nil
}

// Restore restores the component's state using the provided data map. Returns an error if the restoration fails.
func (b *BaseComponent) Restore(state map[string]interface{}) error {
	if err := b.properties.Restore(state); err != nil {
		return err
	}
	return nil
}

// RestorePath restores the state of a component identified by the given path using the provided data map.
func (b *BaseComponent) RestorePath(path string, d map[string]interface{}) error {
	node := b.node.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.GetComponent().Restore(d)
	if err != nil {
		return err
	}
	return nil
}

// RestoreAll restores the state of the entire component hierarchy recursively from the provided state map.
func (b *BaseComponent) RestoreAll(state map[string]interface{}) error {
	if err := b.restore(b.node, state); err != nil {
		return err
	}
	return nil
}

// RunCommand executes a specified command with the given arguments
// using the component's properties and returns the result.
func (b *BaseComponent) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	d, err := b.properties.Run(cmd, args)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// RunCommandPath executes a command on the component located at the specified path with the provided arguments.
func (b *BaseComponent) RunCommandPath(path string, cmd string, args string) (map[string]interface{}, error) {
	node := b.node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	values := strings.Split(args, " ")
	d, err := node.GetComponent().RunCommand(cmd, values)
	return d, err
}

// Path returns the full hierarchical path of the node by appending its parent's path to its own unique path.
func (b *BaseComponent) Path() string {
	id := b.GetId()
	parent := b.node.GetParent()
	if parent == nil || id == "" {
		return id
	}
	return parent.GetComponent().Path() + "." + id
}

// Print writes the `BaseComponent`'s ID and optional type to the provided writer, with indentation for readability.
// It iterates through child components, recursively calling their `Print` method for hierarchical output.
func (b *BaseComponent) Print(w io.Writer, indent string, showComponents bool) {
	_, _ = fmt.Fprintf(w, "%s%s", indent, b.GetId())
	if showComponents {
		_, _ = fmt.Fprintf(w, " (%T)", b)
	}
	_, _ = fmt.Fprintln(w)
	for _, child := range b.node.GetChildren() {
		child.GetComponent().Print(w, indent+"  ", showComponents)
	}
}

// dump recursively traverses the node hierarchy, extracts component data, and populates the rootMap with component states.
func (b *BaseComponent) dump(node *Node, rootMap map[string]interface{}) error {
	if node == nil {
		return nil
	}
	if component := node.GetComponent(); component != nil {
		data, err := component.Dump()
		if err != nil {
			return err
		}
		rootMap[node.component.GetId()] = data
	}
	for _, child := range node.GetChildren() {
		if err := b.dump(child, rootMap); err != nil {
			return err
		}
	}
	return nil
}

// restore traverses the node tree and restores the state of each component from the provided state map,
// returning errors if any.
func (b *BaseComponent) restore(node *Node, state map[string]interface{}) error {
	if node == nil {
		return nil
	}
	if component := node.GetComponent(); component != nil {
		componentState, ok := state[node.component.GetId()].(map[string]interface{})
		if ok {
			err := component.Restore(componentState)
			if err != nil {
				return fmt.Errorf("error restoring component %s: %w", node.component.GetId(), err)
			}
		}
	}
	for _, child := range node.GetChildren() {
		if err := b.restore(child, state); err != nil {
			return err
		}
	}
	return nil
}
