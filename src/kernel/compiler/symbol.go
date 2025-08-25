package compiler

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Symbol represents an identifier with associated metadata such as name, scope, index, fields, and type information.
type Symbol struct {
	name       string
	scope      SymbolScope
	index      int
	structName string
	funcName   string
	types      []string
	isStruct   bool
	Fields     []*StructField
	object     objects.IObject
}

// NewSymbol creates a new Symbol instance with the provided name, index, scope, and associated object.
func NewSymbol(name string, index int, scope SymbolScope, structName string, funcName string, isStruct bool) *Symbol {
	symbol := &Symbol{}
	symbol.Reset(name, index, scope, structName, funcName, isStruct)
	return symbol
}

// Reset resets the Symbol to its initial state with the provided name, index, scope, and associated object.
func (s *Symbol) Reset(name string, index int, scope SymbolScope, structName string, funcName string, isStruct bool) {
	s.name = name
	s.index = index
	s.scope = scope
	s.structName = structName
	s.funcName = funcName
	s.types = []string{}
	s.object = nil
	s.isStruct = isStruct
	s.Fields = []*StructField{}
}

// Name returns the name of the Symbol.
func (s *Symbol) Name() string {
	return s.name
}

// Index returns the index of the Symbol within the program.
func (s *Symbol) Index() int {
	return s.index
}

// SetObject assigns the provided IObject implementation to the Symbol, allowing it to associate with specific metadata.
func (s *Symbol) SetObject(obj objects.IObject) {
	s.object = obj
}

// GetObject retrieves the associated IObject instance from the Symbol.
func (s *Symbol) GetObject() objects.IObject {
	return s.object
}

// IsStruct returns a boolean indicating whether the Symbol represents an struct.
func (s *Symbol) IsStruct() bool {
	//if len(s.structName) > 0 {
	//	return true
	//}
	if s.isStruct {
		return true
	}
	return false
	//return s.isStruct
}

// SetTypes assigns a slice of type names to the Symbol's types field. Use this to define the types associated with a Symbol.
func (s *Symbol) SetTypes(t []string) {
	s.types = t
}

func (s *Symbol) StructPropertyAssign(fields []*StructField) {
	s.Fields = make([]*StructField, len(fields))
	for i, f := range fields {
		s.Fields[i] = NewStructProperty(f.name, f.base, f.kind, nil)
	}
}

// Types returns the list of type names associated with the Symbol.
func (s *Symbol) Types() []string {
	return s.types
}

// StructName returns the name of the container associated with the Symbol.
func (s *Symbol) StructName() string {
	return s.structName
}

// Scope returns the scope of the Symbol, indicating its visibility and context within the program.
func (s *Symbol) Scope() SymbolScope {
	return s.scope
}

// SetScope assigns a new scope to the Symbol.
func (s *Symbol) SetScope(scope SymbolScope) {
	s.scope = scope
}
