package board

import (
	"fmt"
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

// AssignNode assigns a component to a node, creating a new child node if a parentNode is provided, or a root node otherwise.
// Updates the component to reference the assigned node in the hierarchical structure.
func AssignNode(parentNode *Node, component IComponent) {
	var node *Node
	if parentNode != nil {
		node = parentNode.AddComponent(component)
	} else {
		node = newNode(component, nil)
	}
	component.SetNode(node)
}

// NewNode creates a new Node with the given component and parent, initializing its path and children map.
func newNode(component IComponent, parent *Node) *Node {
	if component == nil {
		panic("nil component")
	}
	return &Node{
		component: component,
		parent:    parent,
		children:  make(map[string]*Node),
	}
}

func (n *Node) GetParent() *Node {
	return n.parent
}

// GetComponent retrieves the component associated with the node, or nil if the node is an intermediate node.
func (n *Node) GetComponent() IComponent {
	return n.component
}

// RemoveChild removes the specified child node from the current node's children map based on its path.
func (n *Node) RemoveChild(child *Node) {
	delete(n.children, child.component.GetId())
}

// AddChild adds a child node to the current node's children map. Returns an error if a child with the same ID already exists.
func (n *Node) AddChild(child *Node) {
	n.children[child.component.GetId()] = child
}

// AddComponent adds a new component as a child node to the current node. It returns an error if addition fails.
func (n *Node) AddComponent(component IComponent) *Node {
	child := newNode(component, n)
	n.AddChild(child)
	return child
}

// FindNode traverses the hierarchical structure of nodes to locate a specific node based on the given path string.
func (n *Node) FindNode(path string) *Node {
	parts, err := createPathFromKey(path)
	if err != nil {
		return nil
	}
	id := n.component.GetId()
	if parts[0] != id && id != "" {
		return nil
	}
	parts = parts[1:]
	if len(parts) == 0 {
		return nil
	}
	currentNode := n
	for _, part := range parts {
		nextNode, ok := currentNode.children[part]
		if !ok {
			return nil
		}
		currentNode = nextNode
	}
	return currentNode
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
