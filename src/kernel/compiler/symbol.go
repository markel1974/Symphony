package compiler

import (
	"go/ast"
)

// FieldDef represents the definition of a field, including its name, type, and associated AST node.
type FieldDef struct {
	name string
	kind string
	node ast.Node
}

// NewFieldDef creates a new FieldDef with the specified name, type, and associated AST node.
func NewFieldDef(name string, kind string, node ast.Node) *FieldDef {
	return &FieldDef{
		name: name,
		kind: kind,
		node: node,
	}
}

// Name returns the name of the field defined by the FieldDef.
func (f *FieldDef) Name() string {
	return f.name
}

// Type returns the kind of the field represented by the FieldDef.
func (f *FieldDef) Type() string {
	return f.kind
}

// Node returns the underlying ast.Node associated with the FieldDef instance.
func (f *FieldDef) Node() ast.Node {
	return f.node
}

// SetNode assigns the provided AST node to the FieldDef instance, representing the associated syntax tree element.
func (f *FieldDef) SetNode(node ast.Node) {
	f.node = node
}

// Symbol represents an identifier with associated metadata such as name, scope, index, fields, and type information.
type Symbol struct {
	name   string
	scope  SymbolScope
	index  int
	object string
	types  []string
	Fields []*FieldDef
}

// NewSymbol creates a new Symbol instance with the provided name, index, scope, and associated object.
func NewSymbol(name string, index int, scope SymbolScope, obj string) *Symbol {
	symbol := &Symbol{
		name:   name,
		index:  index,
		scope:  scope,
		object: obj,
		types:  []string{},
		Fields: []*FieldDef{},
	}
	return symbol
}

// Name returns the name of the Symbol.
func (s *Symbol) Name() string {
	return s.name
}

// Index returns the index of the Symbol within the program.
func (s *Symbol) Index() int {
	return s.index
}

// SetTypes assigns a slice of type names to the Symbol's types field. Use this to define the types associated with a Symbol.
func (s *Symbol) SetTypes(t []string) {
	s.types = t
}

// Types returns the list of type names associated with the Symbol.
func (s *Symbol) Types() []string {
	return s.types
}

// Object returns the internal object string associated with the Symbol.
func (s *Symbol) Object() string {
	return s.object
}

// Scope returns the scope of the Symbol, indicating its visibility and context within the program.
func (s *Symbol) Scope() SymbolScope {
	return s.scope
}

// SetScope assigns a new scope to the Symbol.
//func (s *Symbol) SetScope(scope SymbolScope) {
//	s.scope = scope
//}
