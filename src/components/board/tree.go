package board

import (
	"fmt"
	"io"
	"strings"
)

// createPathFromKey splits a dot-separated key into a slice of strings, representing path segments.
// Returns an error if the key is empty or invalid.
func createPathFromKey(key string) ([]string, error) {
	v := strings.Split(key, ".")
	if len(v) == 0 {
		return nil, fmt.Errorf("invalid path")
	}
	return v, nil
}

// Node represents a hierarchical structure of components or intermediate nodes.
// It contains a reference to a component, a unique path identifier, a parent node, and a map of children nodes.
type Node struct {
	component IComponent       // Il componente (CPU, VIC-II, ecc.), o nil per i nodi intermedi.
	path      string           // L'ID univoco del componente (es: "c64.cia1").
	parent    *Node            // Puntatore al nodo genitore (nil per la radice).
	children  map[string]*Node // Mappa dei figli (chiave: ID del figlio, valore: puntatore al nodo figlio).
}

// Node rappresenta un nodo nell'albero dei componenti.

// NewNode creates a new Node instance with the given component and parent, initializing the path and children map.
func NewNode(component IComponent, parent *Node) *Node {
	path := ""
	if component != nil {
		path = component.GetId()
	}
	return &Node{
		component: component,
		path:      path,
		parent:    parent,
		children:  make(map[string]*Node),
	}
}

// GetPath returns the unique identifier of the node as a string.
func (n *Node) GetPath() string {
	return n.path
}

// GetComponent retrieves the component associated with the node. Returns nil if the node has no associated component.
func (n *Node) GetComponent() IComponent {
	return n.component
}

// AddChild adds a child node to the current node's children map. Returns an error if a child with the same path already exists.
func (n *Node) AddChild(child *Node) error {
	if _, ok := n.children[child.GetPath()]; ok {
		return fmt.Errorf("duplicate child ID: %s", child.GetPath())
	}
	n.children[child.GetPath()] = child
	return nil
}

// RemoveChild removes the specified child node from the children map of the current node.
func (n *Node) RemoveChild(child *Node) {
	delete(n.children, child.GetPath())
}

// FindNode traverses the tree from the current node to locate a node identified by the given path. Returns nil if not found.
func (n *Node) FindNode(path string) *Node {
	parts, err := createPathFromKey(path)
	if err != nil {
		return nil
	}
	if parts[0] != n.GetPath() && n.GetPath() != "" {
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

// Path returns the hierarchical path of the node by concatenating its parent's path with its own, separated by a dot.
func (n *Node) Path() string {
	if n.parent == nil || n.parent.GetPath() == "" {
		return n.GetPath()
	}
	return n.parent.Path() + "." + n.GetPath()
}

// GetChild retrieves the child node associated with the given ID from the current node's children. Returns nil if not found.
func (n *Node) GetChild(id string) *Node {
	return n.children[id]
}

// GetChildren retrieves all child nodes of the current node and returns them as a slice of pointers to Node.
func (n *Node) GetChildren() []*Node {
	var children []*Node
	for _, child := range n.children {
		children = append(children, child)
	}
	return children
}

// Print prints the current Node along with its children to the provided writer, using the given indentation and options.
func (n *Node) Print(w io.Writer, indent string, showComponents bool) {
	if n == nil {
		return
	}
	_, _ = fmt.Fprintf(w, "%s%s", indent, n.GetPath())
	if showComponents && n.GetComponent() != nil {
		_, _ = fmt.Fprintf(w, " (%T)", n.GetComponent())
	}
	_, _ = fmt.Fprintln(w) // A capo
	for _, child := range n.GetChildren() {
		child.Print(w, indent+"  ", showComponents)
	}
}

// Tree represents a hierarchical structure with a root node of type *Node.
type Tree struct {
	root *Node
}

// NewTree creates a new Tree structure initialized with the given root component.
func NewTree(root IComponent) *Tree {
	return &Tree{root: NewNode(root, nil)}
}

// Add integrates a component into the tree by creating its corresponding node and associating it with the correct parent node.
// Returns an error if the operation fails, such as when intermediate nodes are missing or a child cannot be added.
func (t *Tree) Add(components IComponent) error {
	path := components.GetId()
	parts, err := createPathFromKey(path)
	if err != nil {
		return err
	}
	node := t.root
	for _, part := range parts[:len(parts)-1] {
		if child := node.GetChild(part); child == nil {
			return fmt.Errorf("missing intermediate node")
		}
	}
	newNode := NewNode(components, node)
	if err = node.AddChild(newNode); err != nil {
		return err
	}
	return nil
}

// Print writes the tree structure to the specified writer, optionally displaying the components if showComponents is true.
func (t *Tree) Print(w io.Writer, showComponents bool) {
	t.root.Print(w, "", showComponents)
}

// RunCommand executes a command on the component at the specified path with the given arguments and returns the result.
func (t *Tree) RunCommand(path string, cmd string, args string) (map[string]interface{}, error) {
	node := t.root.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	values := strings.Split(args, " ")
	d, err := node.GetComponent().GetProperties().Run(cmd, values)
	return d, err
}

// SetProperty sets a property identified by `prop` to the value `val` for the component located at the given `path`.
// Returns an error if the component is not found or if setting the property fails.
func (t *Tree) SetProperty(path string, prop string, val interface{}) error {
	node := t.root.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.GetComponent().GetProperties().SetProperty(prop, val)
	return err
}

// GetProperty retrieves the value of a specific property from a component located at a given path in the Tree.
// Returns the property value and an error if the component or property is not found.
func (t *Tree) GetProperty(path string, prop string) (interface{}, error) {
	node := t.root.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	v, err := node.GetComponent().GetProperties().GetProperty(prop)
	return v, err
}

// DumpComponent retrieves and returns the properties of the component located at the specified path as a map.
func (t *Tree) DumpComponent(path string) (map[string]interface{}, error) {
	node := t.root.FindNode(path)
	if node == nil {
		return nil, fmt.Errorf("component %s not found", path)
	}
	d, err := node.GetComponent().GetProperties().Dump()
	return d, err
}

// RestoreComponent attempts to restore the properties of a component at the given path using the provided data map.
func (t *Tree) RestoreComponent(path string, d map[string]interface{}) error {
	node := t.root.FindNode(path)
	if node == nil {
		return fmt.Errorf("component %s not found", path)
	}
	err := node.GetComponent().GetProperties().Restore(d)
	return err
}

// Dump constructs and returns a map representing the entire tree structure and its state. Returns an error if any occurs.
func (t *Tree) Dump() (map[string]interface{}, error) {
	rootMap := make(map[string]interface{})
	if err := t.dumpNode(t.root, rootMap); err != nil {
		return nil, err
	}
	return rootMap, nil
}

// Restore restores the state of the tree and its nodes from the provided state map. Returns an error if restoration fails.
func (t *Tree) Restore(state map[string]interface{}) error {
	if err := t.restoreNode(t.root, state); err != nil {
		return err
	}
	return nil
}

// dumpNode traverses a node and its children, extracting component data and populating the provided map with the results.
func (t *Tree) dumpNode(node *Node, rootMap map[string]interface{}) error {
	if node == nil {
		return nil
	}
	if component := node.GetComponent(); component != nil {
		properties := component.GetProperties()
		data, err := properties.Dump()
		if err != nil {
			return err
		}
		rootMap[node.GetPath()] = data
	}
	for _, child := range node.GetChildren() {
		if err := t.dumpNode(child, rootMap); err != nil {
			return err
		}
	}
	return nil
}

// restoreNode restores the state of a node and its children using the given state map.
// Returns an error if restoring the state fails for any component.
func (t *Tree) restoreNode(node *Node, state map[string]interface{}) error {
	if node == nil {
		return nil
	}
	if component := node.GetComponent(); component != nil {
		componentState, ok := state[node.GetPath()].(map[string]interface{})
		if ok {
			properties := component.GetProperties()
			err := properties.Restore(componentState)
			if err != nil {
				return fmt.Errorf("error restoring component %s: %w", node.GetPath(), err)
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
