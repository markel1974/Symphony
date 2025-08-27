package compiler

import (
	"go/ast"
	"strconv"
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
	container        map[string][]*FieldDescription
	anonymousCounter int
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable() *StructTable {
	return &StructTable{
		container:        make(map[string][]*FieldDescription),
		anonymousCounter: 0,
	}
}

// CreateStructName creates a unique name for a struct based on the provided name.
func (s *StructTable) CreateStructName(name string) string {
	if len(name) == 0 {
		r := "<anonymous_" + strconv.Itoa(s.anonymousCounter) + ">"
		s.anonymousCounter++
		return r
	}
	return name
}

// Add adds a new field description to a struct in the StructTable. If the struct does not exist, it creates it.
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
	fields, ok := s.container[name]
	if !ok {
		return nil, false
	}
	out := make([]*FieldDescription, len(fields))
	for idx, v := range fields {
		out[idx] = NewFieldDescription(v.name, v.base, v.kind, nil)
	}
	return out, true
}

// Has checks if a struct definition with the given name exists in the container map.
func (s *StructTable) Has(name string) bool {
	if _, ok := s.container[name]; ok {
		return true
	}
	return false
}
