package compiler

import (
	"fmt"
	"go/ast"
)

// StructProperty represents metadata about a struct field, including its name, base type, full type, and AST node.
type StructProperty struct {
	name string
	base string
	kind string
	node ast.Node
}

// NewStructProperty creates a new instance of StructProperty with the provided name, base, kind, and AST node.
func NewStructProperty(name string, base string, kind string, node ast.Node) *StructProperty {
	return &StructProperty{
		name: name,
		base: base,
		kind: kind,
		node: node,
	}
}

// Structs is a collection that manages mappings of struct names to their associated properties.
type Structs struct {
	container map[string][]*StructProperty
}

// NewStructs initializes and returns a pointer to a Structs instance with an empty container map.
func NewStructs() *Structs {
	return &Structs{
		container: make(map[string][]*StructProperty),
	}
}

// Add adds a new struct definition with the given name and properties to the container map in the Structs instance.
func (s *Structs) Add(name string, def []*StructProperty) {
	s.container[name] = def
	fmt.Printf("Struct %s defined with fields:\n", name)
	for _, p := range def {
		fmt.Printf("\t%s %s %s\n", p.name, p.base, p.kind)
	}
	fmt.Printf("\n")
}

// Get retrieves a slice of StructProperty pointers associated with the given name from the container map.
func (s *Structs) Get(name string) []*StructProperty {
	return s.container[name]
}

/*
// Get recupera la definizione di uno struct tramite il suo nome.
func (s *Structs) Get(name string) (*StructDef, bool) {
	def, ok := s.container[name]
	return def, ok
}

*/
