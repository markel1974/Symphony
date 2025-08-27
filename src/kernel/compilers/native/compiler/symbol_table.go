package compiler

import (
	"fmt"
	"io"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// SymbolScope represents the scope of a symbol in a program, such as global, local, free, builtin, or type-specific.
type SymbolScope string

// GlobalScope represents a symbol defined in the global scope.
// LocalScope represents a symbol defined in the local function scope.
// FreeScope represents a free variable captured from an enclosing scope.
// TypeScope represents a custom type definition in the scope.
const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	FreeScope    SymbolScope = "FREE"
	UnknownScope SymbolScope = "UNKNOWN"
)

// SymbolTable represents a hierarchical structure for storing and resolving symbols in various scopes of the program.
// It supports adding, resolving, and managing symbols across nested and global scopes.
// The structure tracks variable and function definitions and handles free symbols for closures.
type SymbolTable struct {
	outer         *SymbolTable
	symbols       map[string]*Symbol
	definitions   []*Symbol
	freeSymbols   []*Symbol
	uniqueCounter int
	funcName      string
}

// NewSymbolTable initializes and returns a new instance of SymbolTable with an empty container and counter set to zero.
func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{
		symbols:       make(map[string]*Symbol),
		uniqueCounter: 0,
	}
	return s
}

// NewEnclosedSymbolTable creates a new symbol table enclosed by the provided outer symbol table.
func NewEnclosedSymbolTable(outer *SymbolTable, funcName string) *SymbolTable {
	s := NewSymbolTable()
	s.outer = outer
	s.funcName = funcName
	return s
}

func (s *SymbolTable) Count() int {
	return len(s.definitions)
}

// Print displays the symbols stored in the SymbolTable, excluding those with the "BUILTIN" scope.
func (s *SymbolTable) Print(writer io.Writer) {
	if s.outer != nil {
		s.outer.Print(writer)
	}
	for k, v := range s.symbols {
		_, _ = fmt.Fprintf(writer, "%s %s %v %d %v", k, v.Name(), v.Scope(), v.Index(), v.Fields)
	}
	for idx, v := range s.freeSymbols {
		_, _ = fmt.Fprintf(writer, "%d %s %v %d %v", idx, v.Name(), v.Scope(), v.Index(), v.Fields)
	}
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
func (s *SymbolTable) DefineUnique(name string) (*Symbol, error) {
	uniqueName := name + strconv.Itoa(s.uniqueCounter)
	s.uniqueCounter++
	return s.Define(uniqueName)
}

// Define creates a new Symbol with the given name, assigns it a scope and index, and stores it in the symbol table.
func (s *SymbolTable) Define(name string) (*Symbol, error) {
	if _, ok := s.symbols[name]; ok {
		return nil, fmt.Errorf("symbol '%s' already defined", name)
	}
	scope := s.computeScope(UnknownScope)
	symbol := NewSymbol(name, len(s.definitions), scope)
	s.definitions = append(s.definitions, symbol)
	s.symbols[name] = symbol
	return symbol, nil
}

// Resolve attempts to look up a symbol by name in the current SymbolTable and outer scopes, if applicable.
// It returns the found Symbol and a boolean indicating whether the resolution was successful.
func (s *SymbolTable) Resolve(name string) (*Symbol, bool) {
	if obj, ok := s.symbols[name]; ok {
		return obj, true
	}
	if s.outer == nil {
		return nil, false
	}
	obj, ok := s.outer.Resolve(name)
	if !ok {
		return obj, ok
	}
	// Types, global variables, and builtin functions are directly accessible
	// from inner scopes and should not be converted to "free variables".
	if obj.Scope() == GlobalScope {
		return obj, true
	}
	s.freeSymbols = append(s.freeSymbols, obj)
	symbol := NewSymbol(obj.Name(), len(s.freeSymbols)-1, FreeScope)
	symbol.SetIsStruct(obj.StructName(), obj.IsStruct())
	s.symbols[obj.Name()] = symbol
	return symbol, true
}

// RebuildScope updates the scope of a Symbol if it exists and returns the updated Symbol alongside a boolean for success.
func (s *SymbolTable) RebuildScope(name string, scope SymbolScope) (*Symbol, bool) {
	symbol, ok := s.Resolve(name)
	if !ok {
		return nil, false
	}
	scope = s.computeScope(scope)
	symbol.SetScope(scope)
	return symbol, true
}

func (s *SymbolTable) computeScope(scope SymbolScope) SymbolScope {
	if scope == UnknownScope {
		if s.outer == nil {
			scope = GlobalScope
		} else {
			scope = LocalScope
		}
	}
	return scope
}
