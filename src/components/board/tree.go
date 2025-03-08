package board

import (
	"fmt"
	"io"
	"strings"
)

// createPathFromKey splits a dot-separated string key into a slice of strings representing the path components.
// Returns an error if the key is invalid or empty.
func createPathFromKey(key string) ([]string, error) {
	v := strings.Split(key, ".")
	if len(v) == 0 {
		return nil, fmt.Errorf("invalid path")
	}
	return v, nil
}

// Node represents a hierarchical structure used to manage components and their relationships.
type Node struct {
	component IComponent       // Il componente (CPU, VIC-II, ecc.), o nil per i nodi intermedi.
	parent    *Node            // Puntatore al nodo genitore (nil per la radice).
	children  map[string]*Node // Mappa dei figli (chiave: ID del figlio, valore: puntatore al nodo figlio).
}

// Node rappresenta un nodo nell'albero dei componenti.

// NewNode creates a new Node with the given component and parent, initializing its path and children map.
func NewNode(component IComponent, parent *Node) *Node {
	if component == nil {
		panic("nil component")
	}
	return &Node{
		component: component,
		parent:    parent,
		children:  make(map[string]*Node),
	}
}

// GetComponent retrieves the component associated with the node, or nil if the node is an intermediate node.
func (n *Node) GetComponent() IComponent {
	return n.component
}

// AddComponent adds a new component as a child node to the current node. It returns an error if addition fails.
func (n *Node) AddComponent(component IComponent) *Node {
	child := NewNode(component, n)
	n.AddChild(child)
	return child
}

// AddChild adds a child node to the current node's children map. Returns an error if a child with the same ID already exists.
func (n *Node) AddChild(child *Node) {
	n.children[child.component.GetId()] = child
}

// RemoveChild removes the specified child node from the current node's children map based on its path.
func (n *Node) RemoveChild(child *Node) {
	delete(n.children, child.component.GetId())
}

// FindNode traverses the hierarchical structure of nodes to locate a specific node based on the given path string.
func (n *Node) FindNode(path string) *Node {
	parts, err := createPathFromKey(path)
	if err != nil {
		return nil
	}
	id := n.component.GetId()
	if parts[0] != id && id != "" {
		return nil // Il percorso non inizia con l'ID di questo nodo.
	}
	parts = parts[1:]
	if len(parts) == 0 {
		return nil
	}
	currentNode := n
	for _, part := range parts {
		nextNode, ok := currentNode.children[part]
		if !ok {
			return nil // Nodo non trovato.
		}
		currentNode = nextNode
	}
	return currentNode
}

// Path returns the full hierarchical path of the node by appending its parent's path to its own unique path.
func (n *Node) Path() string {
	id := n.component.GetId()
	if n.parent == nil || id == "" {
		return id
	}
	return n.parent.Path() + "." + id
}

// GetChild retrieves the child node with the specified ID from the current node's children. Returns nil if not found.
func (n *Node) GetChild(id string) *Node {
	return n.children[id]
}

// HasChild checks if the node has a child with the specified ID and returns true if it exists, otherwise false.
func (n *Node) HasChild(id string) bool {
	return n.children[id] != nil
}

// GetChildren returns a slice of pointers to all child nodes of the current node.
func (n *Node) GetChildren() []*Node {
	var children []*Node
	for _, child := range n.children {
		children = append(children, child)
	}
	return children
}

// Print writes the current node's path to the specified writer, optionally including component types, with indentation.
// If showComponents is true and the node has a component, the component's type is displayed.
// Children are printed recursively with increased indentation.
func (n *Node) Print(w io.Writer, indent string, showComponents bool) {
	if n == nil {
		return
	}
	_, _ = fmt.Fprintf(w, "%s%s", indent, n.component.GetId())
	if showComponents && n.GetComponent() != nil {
		_, _ = fmt.Fprintf(w, " (%T)", n.GetComponent())
	}
	_, _ = fmt.Fprintln(w)
	for _, child := range n.GetChildren() {
		child.Print(w, indent+"  ", showComponents)
	}
}

// GetProperty retrieves the value of the specified property from the component associated with the node.
// It returns the property value or an error if the property cannot be retrieved.
func (n *Node) GetProperty(prop string) (interface{}, error) {
	v, err := n.component.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// SetProperty sets the specified property to the given value for the node's component. Returns an error if the operation fails.
func (n *Node) SetProperty(prop string, value interface{}) error {
	if err := n.component.SetProperty(prop, value); err != nil {
		return err
	}
	return nil
}

// Dump retrieves the current state of the associated component and returns it as a map, along with any potential error.
func (n *Node) Dump() (map[string]interface{}, error) {
	state, err := n.component.Dump()
	if err != nil {
		return state, err
	}
	return state, nil
}

// Restore restores the state of the Node's component using the given state map and returns an error if restoration fails.
func (n *Node) Restore(state map[string]interface{}) error {
	if err := n.component.Restore(state); err != nil {
		return err
	}
	return nil
}

// RunCommand executes a given command with arguments on the component associated with the node and returns the result.
func (n *Node) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	values, err := n.component.RunCommand(cmd, args)
	if err != nil {
		return nil, err
	}
	return values, nil
}

// Tree represents a hierarchical structure starting from a root Node.
type Tree struct {
	*Node
}

// NewTree creates a new Tree with the specified IComponent as the root node.
func NewTree(root *Node) *Tree {
	return &Tree{Node: root}
}

// RunCommandPath executes a command on a component identified by the given path, passing the provided arguments as a string.
// It returns the result as a map and an error if the component is not found or the command execution fails.
func (t *Tree) RunCommandPath(path string, cmd string, args string) (map[string]interface{}, error) {
	node := t.Node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	values := strings.Split(args, " ")
	d, err := node.RunCommand(cmd, values)
	return d, err
}

// SetPropertyPath sets the property `prop` of the component at the specified `path` to the value `val`.
// Returns an error if the component is not found or the property cannot be set.
func (t *Tree) SetPropertyPath(path string, prop string, val interface{}) error {
	node := t.Node.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.SetProperty(prop, val)
	return err
}

// GetPropertyPath retrieves the property value identified by 'prop' from the node at the given 'path' in the tree. It returns an error if the path or property is not found.
func (t *Tree) GetPropertyPath(path string, prop string) (interface{}, error) {
	node := t.Node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	v, err := node.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// DumpComponentPath retrieves the state of a specific component by its path within the tree structure and returns it as a map.
func (t *Tree) DumpComponentPath(path string) (map[string]interface{}, error) {
	node := t.Node.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	d, err := node.Dump()
	return d, err
}

// RestoreComponentPath restores the state of a component at the given path using the provided data map.
func (t *Tree) RestoreComponentPath(path string, d map[string]interface{}) error {
	node := t.Node.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.Restore(d)
	return err
}

// DumpAll generates a deep hierarchical map representation of the tree, starting from its root node. Returns the map and error.
func (t *Tree) DumpAll() (map[string]interface{}, error) {
	rootMap := make(map[string]interface{})
	if err := t.dumpNode(t.Node, rootMap); err != nil {
		return nil, err
	}
	return rootMap, nil
}

// RestoreAll updates the tree structure and its components to a previous state using the provided state map.
func (t *Tree) RestoreAll(state map[string]interface{}) error {
	if err := t.restoreNode(t.Node, state); err != nil {
		return err
	}
	return nil
}

// dumpNode traverses the tree starting from the given node and populates the rootMap with serialized data of each component.
// Returns an error if any component fails to serialize during the traversal.
func (t *Tree) dumpNode(node *Node, rootMap map[string]interface{}) error {
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
		if err := t.dumpNode(child, rootMap); err != nil {
			return err
		}
	}
	return nil
}

// restoreNode recursively restores the state of a node and its children using the provided state map.
// It returns an error if any component restoration fails.
func (t *Tree) restoreNode(node *Node, state map[string]interface{}) error {
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
		if err := t.restoreNode(child, state); err != nil {
			return err
		}
	}
	return nil
}

/*
// Add inserts the given IComponent into the tree at the correct path location, creating a hierarchical structure.
// Returns an error if the path is invalid or intermediate nodes are missing.
func (t *Tree) Add(components IComponent) error {
	path := components.GetId()
	parts, err := createPathFromKey(path)
	if err != nil {
		return err
	}
	node := t.root
	for _, part := range parts[:len(parts)-1] {
		if !node.HasChild(part) {
			return fmt.Errorf("missing intermediate node")
		}
	}
	if err = node.AddComponent(components); err != nil {
		return err
	}
	return nil
}
*/
