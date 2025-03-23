package component

import (
	"fmt"
	"github.com/markel1974/c64emu/src/references"
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
	component references.IComponent // Il componente (CPU, VIC-II, ecc.), o nil per i nodi intermedi.
	parent    *Node                 // Puntatore al nodo genitore (nil per la radice).
	children  map[string]*Node      // Mappa dei figli (chiave: ID del figlio, valore: puntatore al nodo figlio).
}

// NewNode creates a new Node with the given component and parent, initializing its path and children map.
func newNode(parent references.INode, component references.IComponent) *Node {
	var parentNode *Node = nil
	if parent != nil {
		var ok bool
		parentNode, ok = parent.(*Node)
		if !ok {
			panic(fmt.Sprintf("invalid parent node: %T", parent))
		}
	}

	if component == nil {
		panic("nil component")
	}
	n := &Node{
		component: component,
		parent:    parentNode,
		children:  make(map[string]*Node),
	}
	return n
}

func (n *Node) GetParent() references.INode {
	return n.parent
}

// GetComponent retrieves the component associated with the node, or nil if the node is an intermediate node.
func (n *Node) GetComponent() references.IComponent {
	return n.component
}

// RemoveComponent removes the specified child node from the current node's children map based on its path.
func (n *Node) RemoveComponent(component references.IComponent) bool {
	if component == nil {
		return false
	}
	id := component.GetId()
	if _, ok := n.children[component.GetId()]; ok {
		delete(n.children, id)
		return true
	}
	return false
}

// AddComponent adds a new component as a child node to the current node. It returns an error if addition fails.
func (n *Node) AddComponent(component references.IComponent) references.INode {
	child := newNode(n, component)
	n.children[child.component.GetId()] = child
	return child
}

// Traverse traverses the hierarchical structure of nodes to locate a specific node based on the given path string.
func (n *Node) Traverse(path string) references.IComponent {
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
	if currentNode == nil {
		return nil
	}
	return currentNode.GetComponent()
}

// GetChild retrieves the child node with the specified ID from the current node's children. Returns nil if not found.
func (n *Node) GetChild(id string) references.INode {
	return n.children[id]
}

// HasChild checks if the node has a child with the specified ID and returns true if it exists, otherwise false.
func (n *Node) HasChild(id string) bool {
	return n.children[id] != nil
}

// GetChildren returns a slice of pointers to all child nodes of the current node.
func (n *Node) GetChildren() []references.INode {
	var children []references.INode
	for _, child := range n.children {
		children = append(children, child)
	}
	return children
}
