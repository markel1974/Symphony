package compiler

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// SymbolScope defines the scope of a symbol within a program, such as global, local, free, or built-in.
type SymbolScope string

// GlobalScope represents symbols defined in the global scope.
// LocalScope represents symbols defined in the local scope.
// FreeScope represents symbols captured from outer scopes.
// BuiltinScope represents symbols related to built-in identifiers.
const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	FreeScope    SymbolScope = "FREE"
	BuiltinScope SymbolScope = "BUILTIN"
)

// SymbolTable manages variable definitions and resolutions across nested scopes during compilation.
type SymbolTable struct {
	outer          *SymbolTable
	store          map[string]*Symbol
	numDefinitions int
	freeSymbols    []*Symbol
	symbolsCounter int
}

// NewSymbolTable creates and returns a pointer to a new, empty SymbolTable with initialized storage.
func NewSymbolTable() *SymbolTable {
	s := make(map[string]*Symbol)
	return &SymbolTable{
		store:          s,
		symbolsCounter: 0,
	}
}

// NewEnclosedSymbolTable creates a new SymbolTable and sets the provided table as its outer scope.
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.outer = outer
	return s
}

func (s *SymbolTable) NumDefinitions() int {
	return s.numDefinitions
}

func (s *SymbolTable) Outer() *SymbolTable {
	return s.outer
}

// ConvertFreeSymbols transforms a slice of Symbol into a slice of ObjectPointer, used for managing free variables.
func (s *SymbolTable) ConvertFreeSymbols() []*objects.ObjectPointer {
	// Implementazione fittizia per far compilare il codice
	return make([]*objects.ObjectPointer, len(s.freeSymbols))
}

// FreeSymbolsLen returns the number of free symbols in the symbol table.
func (s *SymbolTable) FreeSymbolsLen() int {
	return len(s.freeSymbols)
}

// DefineUnique adds a new symbol with a unique name to the symbol table and assigns it a unique index and scope.
// It returns a pointer to the newly created Symbol.
func (s *SymbolTable) DefineUnique(name string) *Symbol {
	uniqueName := name + strconv.Itoa(s.symbolsCounter)
	s.symbolsCounter++
	return s.Define(uniqueName)
}

// Define adds a new symbol with the provided name to the symbol table and assigns it a unique index and scope.
// It returns a pointer to the newly created Symbol.
func (s *SymbolTable) Define(name string) *Symbol {
	scope := GlobalScope
	if s.outer != nil {
		scope = LocalScope
	}
	symbol := NewSymbol(name, s.numDefinitions, scope)
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// DefineBuiltin defines a built-in symbol in the current symbol table with a given name and index.
func (s *SymbolTable) DefineBuiltin(name string, index int) *Symbol {
	symbol := NewSymbol(name, index, BuiltinScope)
	s.store[name] = symbol
	return symbol
}

// Resolve searches for a symbol by name within the symbol table and its outer scopes, if defined.
// It returns the symbol and a boolean indicating whether it was found.
// If the symbol is found in an outer local scope, it's defined as a free variable in the current scope.
// Symbols in the global or builtin scope are returned directly without modification.
func (s *SymbolTable) Resolve(name string) (*Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.outer != nil {
		obj, ok = s.outer.Resolve(name)
		if !ok {
			return obj, ok
		}
		// Se il simbolo è stato trovato in uno scope esterno (e non è globale/builtin),
		// diventa una "free variable" per lo scope corrente.
		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, true
		}
		return s.defineFree(obj), true
	}
	return obj, ok
}

// defineFree promotes a symbol from an outer scope to be a free variable in the current scope and adds it to freeSymbols.
func (s *SymbolTable) defineFree(original *Symbol) *Symbol {
	s.freeSymbols = append(s.freeSymbols, original)
	symbol := NewSymbol(original.Name, len(s.freeSymbols)-1, FreeScope)
	s.store[original.Name] = symbol
	return symbol
}
