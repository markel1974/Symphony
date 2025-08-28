package compiler

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Symbol represents an identifier with associated metadata such as name, scope, index, fields, and type information.
type Symbol struct {
	name         string
	scope        SymbolScope
	index        int
	typeName     string
	structName   string
	structFields []string
	returnTypes  []string
	object       objects.IObject
	isInterface  bool
}

// NewSymbol creates a new Symbol instance with the provided name, index, scope, and associated object.
func NewSymbol(name string, index int, scope SymbolScope) *Symbol {
	symbol := &Symbol{
		name:         name,
		index:        index,
		scope:        scope,
		typeName:     "",
		object:       nil,
		isInterface:  false,
		structFields: []string{},
		returnTypes:  []string{},
	}
	return symbol
}

func (s *Symbol) Clone() *Symbol {
	return &Symbol{
		name:         s.name,
		scope:        s.scope,
		index:        s.index,
		typeName:     s.typeName,
		structName:   s.structName,
		structFields: s.structFields,
		returnTypes:  s.returnTypes,
		object:       s.object,
		isInterface:  s.isInterface,
	}
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

// TypeName returns the type name associated with the Symbol.
func (s *Symbol) TypeName() string {
	return s.typeName
}

// SetTypeName updates the typeName field of the Symbol with the provided name.
func (s *Symbol) SetTypeName(name string) {
	s.typeName = name
}

// SetStruct assigns a struct name and a slice of field names to the Symbol, updating its associated metadata.
func (s *Symbol) SetStruct(structName string, fields []string) {
	//fmt.Printf("SetStruct %s => %s\n", s.Name(), structName)
	s.structName = structName
	s.structFields = fields
	s.isInterface = false
}

// StructName returns the name of the container associated with the Symbol.
func (s *Symbol) StructName() string {
	return s.structName
}

// StructFields returns the list of field names associated with the Symbol.
func (s *Symbol) StructFields() []string {
	return s.structFields
}

// IsStruct returns a boolean indicating whether the Symbol represents an struct.
func (s *Symbol) IsStruct() bool {
	return len(s.structName) > 0
}

func (s *Symbol) SetInterface(name string) {
	s.isInterface = true
	s.typeName = name
	s.structName = "" // Non può essere uno struct
}

func (s *Symbol) IsInterface() bool {
	return s.isInterface
}

// SetReturnTypes assigns a slice of type names to the Symbol's types field. Use this to define the types associated with a Symbol.
func (s *Symbol) SetReturnTypes(t []string) {
	s.returnTypes = t
}

// ReturnTypes returns the list of type names associated with the Symbol.
func (s *Symbol) ReturnTypes() []string {
	return s.returnTypes
}

// Scope returns the scope of the Symbol, indicating its visibility and context within the program.
func (s *Symbol) Scope() SymbolScope {
	return s.scope
}

// SetScope assigns a new scope to the Symbol.
func (s *Symbol) SetScope(scope SymbolScope) {
	s.scope = scope
}
