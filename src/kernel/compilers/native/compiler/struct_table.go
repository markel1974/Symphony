package compiler

import (
	"go/ast"
)

// FieldDescription represents metadata about a struct field, including its name, base type, full type, and AST node.
type FieldDescription struct {
	name string
	base string
	kind string
	node ast.Node
}

// NewFieldDescription creates a new instance of StructProperty with the provided name, base, kind, and AST node.
func NewFieldDescription(name string, base string, kind string, node ast.Node) *FieldDescription {
	return &FieldDescription{
		name: name,
		base: base,
		kind: kind,
		node: node,
	}
}

// StructTable is a collection that manages mappings of struct names to their associated properties.
type StructTable struct {
	container map[string][]*FieldDescription
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable() *StructTable {
	return &StructTable{
		container: make(map[string][]*FieldDescription),
	}
}

func (s *StructTable) Add(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	// here we could add a check for duplicate fields.
	v := NewFieldDescription(fieldName, baseStruct, kind, node)
	fields, ok := s.container[name]
	if !ok {
		s.container[name] = []*FieldDescription{v}
		return
	}
	s.container[name] = append(fields, v)
}

// GetFields retrieves a slice of StructProperty pointers associated with the given name from the container map.
func (s *StructTable) GetFields(name string) ([]*FieldDescription, bool) {
	v, ok := s.container[name]
	return v, ok
}

// Has checks if a struct definition with the given name exists in the container map.
func (s *StructTable) Has(name string) bool {
	if _, ok := s.container[name]; ok {
		return true
	}
	return false
}
