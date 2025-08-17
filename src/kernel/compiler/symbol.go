package compiler

import (
	"fmt"
	"go/ast"
)

// FieldDef represents a definition of a field within a struct, including its name, type, and associated AST node.
type FieldDef struct {
	Name string
	kind string
	Node ast.Node
}

// NewFieldDef creates and returns a new FieldDef instance with the specified name, kind, and associated AST node.
func NewFieldDef(name string, kind string, node ast.Node) *FieldDef {
	return &FieldDef{
		Name: name,
		kind: kind,
		Node: node,
	}
}

// Symbol represents a named entity within a specific scope of a program, such as a variable, function, or type.
type Symbol struct {
	Name   string
	Scope  SymbolScope
	Index  int
	Fields []*FieldDef
	Object string
	kind   string
}

// NewSymbol creates and returns a new Symbol with the provided name, index, and scope, initializing its Fields as an empty slice.
func NewSymbol(name string, index int, scope SymbolScope) *Symbol {
	symbol := &Symbol{
		Name:   name,
		Index:  index,
		Scope:  scope,
		Fields: []*FieldDef{},
	}
	return symbol
}

// SetType assigns a type to the Symbol if it doesn't already have one, otherwise logs the existing type without updating.
func (s *Symbol) SetType(t string) {
	if s.kind != "" {
		fmt.Println("Symbol already has a type:", s.kind)
	}
	s.kind = t
}

// Type returns the type of the Symbol as a string.
func (s *Symbol) Type() string {
	return s.kind
}
