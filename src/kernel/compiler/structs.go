package compiler

import (
	"go/ast"
	"log"
)

// StructField represents metadata about a struct field, including its name, base type, full type, and AST node.
type StructField struct {
	name string
	base string
	kind string
	node ast.Node
}

// NewStructProperty creates a new instance of StructProperty with the provided name, base, kind, and AST node.
func NewStructProperty(name string, base string, kind string, node ast.Node) *StructField {
	return &StructField{
		name: name,
		base: base,
		kind: kind,
		node: node,
	}
}

// Structs is a collection that manages mappings of struct names to their associated properties.
type Structs struct {
	container map[string][]*StructField
}

// NewStructs initializes and returns a pointer to a Structs instance with an empty container map.
func NewStructs() *Structs {
	return &Structs{
		container: make(map[string][]*StructField),
	}
}

// Add adds a new struct definition with the given name and properties to the container map in the Structs instance.
func (s *Structs) Add(name string, def []*StructField) {
	s.container[name] = def
	log.Printf("Struct %s defined with fields:\n", name)
	for _, p := range def {
		log.Printf("\t%s %s %s\n", p.name, p.base, p.kind)
	}
	log.Printf("\n")
}

// Get retrieves a slice of StructProperty pointers associated with the given name from the container map.
func (s *Structs) Get(name string) []*StructField {
	return s.container[name]
}

/*
// Get recupera la definizione di uno struct tramite il suo nome.
func (s *Structs) Get(name string) (*StructDef, bool) {
	def, ok := s.container[name]
	return def, ok
}

*/
