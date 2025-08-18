package compiler

import (
	"go/ast"
)

// FieldDef represents a definition of a field with a name, type, and associated AST node.
type FieldDef struct {
	name string
	kind string
	node ast.Node
}

// NewFieldDef creates and returns a new instance of FieldDef with the specified name, type kind, and associated AST node.
func NewFieldDef(name string, kind string, node ast.Node) *FieldDef {
	return &FieldDef{
		name: name,
		kind: kind,
		node: node,
	}
}

// Name returns the name of the field defined in FieldDef.
func (f *FieldDef) Name() string {
	return f.name
}

// Type returns the kind of the field as a string.
func (f *FieldDef) Type() string {
	return f.kind
}

// Node returns the associated AST node for the FieldDef instance.
func (f *FieldDef) Node() ast.Node {
	return f.node
}

// SetNode assigns the provided AST node to the FieldDef. It updates the internal `node` field with the given value.
func (f *FieldDef) SetNode(node ast.Node) {
	f.node = node
}

// Symbol represents a named entity with a specific scope, type, index, and optional metadata like fields or kind.
type Symbol struct {
	Name   string
	Scope  SymbolScope
	Index  int
	Fields []*FieldDef
	object string
	types  []string
}

// NewSymbol constructs and returns a new Symbol with the given name, index, and scope, initializing its fields as empty.
func NewSymbol(name string, index int, scope SymbolScope, obj string) *Symbol {
	symbol := &Symbol{
		Name:   name,
		Index:  index,
		Scope:  scope,
		Fields: []*FieldDef{},
		object: obj,
		types:  []string{},
	}
	return symbol
}

// SetTypes assigns a type to the Symbol, emitting a warning if the Symbol already has a type.
func (s *Symbol) SetTypes(t []string) {
	s.types = t
}

// Types return the type of the symbol as a string.
func (s *Symbol) Types() []string {
	return s.types
}

func (s *Symbol) Object() string {
	return s.object
}
