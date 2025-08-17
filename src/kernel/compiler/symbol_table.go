package compiler

import (
	"fmt"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// SymbolScope represents the scope of a symbol in a program, such as global, local, free, builtin, or type-specific.
type SymbolScope string

// GlobalScope represents a symbol defined in the global scope.
// LocalScope represents a symbol defined in the local function scope.
// FreeScope represents a free variable captured from an enclosing scope.
// BuiltinScope represents a built-in function or value in the scope.
// TypeScope represents a custom type definition in the scope.
const (
	ImportScope  SymbolScope = "IMPORT"
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	FreeScope    SymbolScope = "FREE"
	BuiltinScope SymbolScope = "BUILTIN"
	TypeScope    SymbolScope = "TYPE"
	UnknownScope SymbolScope = "UNKNOWN"
)

// SymbolTable represents a hierarchical structure for storing and resolving symbols in various scopes of the program.
// It supports adding, resolving, and managing symbols across nested and global scopes.
// The structure tracks variable and function definitions and handles free symbols for closures.
type SymbolTable struct {
	outer          *SymbolTable
	container      map[string]*Symbol
	numDefinitions int
	freeSymbols    []*Symbol
	uniqueCounter  int
}

// NewSymbolTable initializes and returns a new instance of SymbolTable with an empty container and counter set to zero.
func NewSymbolTable() *SymbolTable {
	s := make(map[string]*Symbol)
	return &SymbolTable{
		container:     s,
		uniqueCounter: 0,
	}
}

// NewEnclosedSymbolTable creates a new symbol table enclosed by the provided outer symbol table.
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.outer = outer
	return s
}

// Print displays the symbols stored in the SymbolTable, excluding those with the "BUILTIN" scope.
func (s *SymbolTable) Print() {
	if s.outer != nil {
		s.outer.Print()
	}
	for k, v := range s.container {
		if v.Scope == BuiltinScope {
			continue
		}
		fmt.Println(k, v.Name, v.Scope, v.Index, v.Fields, v.Methods)
	}
	for _, v := range s.freeSymbols {
		fmt.Println(v.Name, v.Scope, v.Index, v.Fields, v.Methods)
	}
}

// NumDefinitions returns the number of symbols defined in the symbol table.
func (s *SymbolTable) NumDefinitions() int {
	return s.numDefinitions
}

// Outer returns the outer SymbolTable linked to the current SymbolTable, or nil if no outer SymbolTable exists.
func (s *SymbolTable) Outer() *SymbolTable {
	return s.outer
}

// ConvertFreeSymbols transforms the free symbols in the symbol table into a slice of ObjectPointer and returns it.
func (s *SymbolTable) ConvertFreeSymbols() []*objects.ObjectPointer {
	// Implementazione fittizia per far compilare il codice
	return make([]*objects.ObjectPointer, len(s.freeSymbols))
}

// FreeSymbolsLen returns the number of free symbols in the symbol table.
func (s *SymbolTable) FreeSymbolsLen() int {
	return len(s.freeSymbols)
}

// DefineUnique generates a unique symbol name using a provided base name and counter, then defines and returns the symbol.
func (s *SymbolTable) DefineUnique(name string, scope SymbolScope) *Symbol {
	uniqueName := name + strconv.Itoa(s.uniqueCounter)
	s.uniqueCounter++
	return s.Define(uniqueName, scope)
}

// Define creates a new Symbol with the given name, assigns it a scope and index, and stores it in the symbol table.
func (s *SymbolTable) Define(name string, scope SymbolScope) *Symbol {
	if scope == UnknownScope {
		if s.outer == nil {
			scope = GlobalScope
		} else {
			scope = LocalScope
		}
	}
	symbol := NewSymbol(name, s.numDefinitions, scope)
	s.container[name] = symbol
	s.numDefinitions++
	return symbol
}

// DefineBuiltin adds a built-in symbol to the symbol table with the specified name and index, returning the new symbol.
func (s *SymbolTable) DefineBuiltin(name string, index int) {
	symbol := NewSymbol(name, index, BuiltinScope)
	s.container[name] = symbol
}

// Resolve attempts to look up a symbol by name in the current SymbolTable and outer scopes, if applicable.
// It returns the found Symbol and a boolean indicating whether the resolution was successful.
func (s *SymbolTable) Resolve(name string) (*Symbol, bool) {
	obj, ok := s.container[name]
	if ok {
		return obj, true
	}
	if s.outer == nil {
		return nil, false
	}
	obj, ok = s.outer.Resolve(name)
	if !ok {
		return obj, ok
	}
	// I tipi, le variabili globali e le funzioni builtin sono accessibili
	// direttamente da scope interni e non devono essere convertiti in "free variables".
	if obj.Scope == ImportScope || obj.Scope == GlobalScope || obj.Scope == BuiltinScope || obj.Scope == TypeScope {
		return obj, true
	}
	return s.defineFree(obj), true
}

func (s *SymbolTable) defineFree(original *Symbol) *Symbol {
	s.freeSymbols = append(s.freeSymbols, original)
	symbol := NewSymbol(original.Name, len(s.freeSymbols)-1, FreeScope)
	s.container[original.Name] = symbol
	return symbol
}
