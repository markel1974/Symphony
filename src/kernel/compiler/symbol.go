package compiler

import "go/ast"

type FieldDef struct {
	Name string
	Type string
	Node ast.Node
}

// Symbol represents a variable or constant with an associated name, scope, and unique index in a program.
type Symbol struct {
	Name    string
	Scope   SymbolScope
	Index   int
	Fields  []FieldDef
	Methods map[string]int
	Type    string
}

// NewSymbol creates a new Symbol with the given name, index, and scope, and returns a pointer to the Symbol instance.
func NewSymbol(name string, index int, scope SymbolScope) *Symbol {
	symbol := &Symbol{
		Name:    name,
		Index:   index,
		Scope:   scope,
		Fields:  []FieldDef{},
		Methods: make(map[string]int),
	}
	return symbol
}
