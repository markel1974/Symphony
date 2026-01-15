package tables

import (
	"github.com/markel1974/symphony/src/vm/objects"
)

// Symbol represents a named entity with a defined scope, associated with specific attributes and metadata.
type Symbol struct {
	constant      bool
	name          string
	scope         SymbolScope
	index         int
	interfaceName string
	structName    string
	structFields  []string
	inputType     []string
	inputName     []string
	returnTypes   []string
	object        objects.IObject
}

// NewSymbol creates a new Symbol instance with the specified constant index, name, index, and scope.
func NewSymbol(constant bool, name string, index int, scope SymbolScope) *Symbol {
	symbol := &Symbol{
		constant:      constant,
		name:          name,
		index:         index,
		scope:         scope,
		interfaceName: "",
		object:        nil,
		structFields:  []string{},
		returnTypes:   []string{},
	}
	return symbol
}

// Clone creates and returns a deep copy of the current Symbol, ensuring all fields are duplicated.
func (s *Symbol) Clone() *Symbol {
	return &Symbol{
		constant:      s.constant,
		name:          s.name,
		scope:         s.scope,
		index:         s.index,
		interfaceName: s.interfaceName,
		structName:    s.structName,
		structFields:  s.structFields,
		returnTypes:   s.returnTypes,
		object:        s.object,
	}
}

// Name returns the name of the Symbol.
func (s *Symbol) Name() string {
	return s.name
}

// Constant checks if the Symbol represents a constant value and returns true if it does.
func (s *Symbol) Constant() bool {
	return s.constant
}

// Index returns the unique index of the symbol within its scope. It is used to reference the symbol in bytecode.
func (s *Symbol) Index() int {
	return s.index
}

// SetObject assigns the provided object to the Symbol instance.
func (s *Symbol) SetObject(obj objects.IObject) {
	s.object = obj
}

// GetObject retrieves the IObject instance associated with the Symbol.
func (s *Symbol) GetObject() objects.IObject {
	return s.object
}

// InterfaceName returns the name of the interface associated with the symbol.
func (s *Symbol) InterfaceName() string {
	return s.interfaceName
}

// SetStruct assigns a struct name and its fields to the Symbol, clearing any associated interface name.
func (s *Symbol) SetStruct(structName string, fields []string) {
	//fmt.Printf("SetStruct %s => %s\n", s.Name(), structName)
	s.structName = structName
	s.structFields = fields
	s.interfaceName = ""
}

// StructName returns the name of the struct associated with the Symbol.
func (s *Symbol) StructName() string {
	return s.structName
}

// StructFields returns the list of field names for a struct associated with the Symbol instance.
func (s *Symbol) StructFields() []string {
	return s.structFields
}

// IsStruct determines whether the current Symbol instance represents a struct by checking if structName is non-empty.
func (s *Symbol) IsStruct() bool {
	return len(s.structName) > 0
}

// SetInterface sets the symbol's interface name and ensures it is not simultaneously a struct by clearing the struct name.
func (s *Symbol) SetInterface(name string) {
	s.interfaceName = name
	s.structName = "" // cant' be a struct and an interface at the same time
}

// IsInterface checks if a Symbol represents an interface by verifying if its interfaceName field is non-empty.
func (s *Symbol) IsInterface() bool {
	return len(s.interfaceName) > 0
}

// SetReturnTypes assigns the specified slice of return types to the Symbol.
func (s *Symbol) SetReturnTypes(t []string) {
	s.returnTypes = t
}

// SetInputTypes sets the input parameter names and their corresponding types for the symbol.
func (s *Symbol) SetInputTypes(name []string, kind []string) {
	s.inputName = name
	s.inputType = kind
}

// ReturnTypes retrieves the list of return types associated with the Symbol instance.
func (s *Symbol) ReturnTypes() []string {
	return s.returnTypes
}

// ReturnTypeFirst retrieves the first return type of the Symbol and a boolean indicating its presence or absence.
func (s *Symbol) ReturnTypeFirst() (string, bool) {
	if len(s.ReturnTypes()) == 0 {
		return "", false
	}
	return s.ReturnTypes()[0], true
}

// Scope returns the scope of the current symbol as a SymbolScope.
func (s *Symbol) Scope() SymbolScope {
	return s.scope
}

// SetScope updates the scope of the Symbol to the provided SymbolScope.
func (s *Symbol) SetScope(scope SymbolScope) {
	s.scope = scope
}
