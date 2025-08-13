package compiler

// Symbol represents a variable or constant with an associated name, scope, and unique index in a program.
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// NewSymbol creates a new Symbol with the given name, index, and scope, and returns a pointer to the Symbol instance.
func NewSymbol(name string, index int, scope SymbolScope) *Symbol {
	symbol := &Symbol{Name: name, Index: index, Scope: scope}
	return symbol
}
