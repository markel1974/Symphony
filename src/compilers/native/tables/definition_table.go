package tables

import (
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// DefinitionTable represents a central registry managing scopes, structure definitions, and interface relationships.
// It integrates with IGateKeeper to handle object creation, conversion, and adaptation efficiently.
type DefinitionTable struct {
	gk             objects.IGateKeeper
	scopes         *Scopes
	structTable    *StructTable
	interfaceTable *InterfaceTable
}

// NewDefinitionTable creates and initializes a new DefinitionTable instance with the provided gatekeeper, scopes, and tables.
func NewDefinitionTable(gk objects.IGateKeeper, scopes *Scopes, structTable *StructTable, interfaceTable *InterfaceTable) *DefinitionTable {
	return &DefinitionTable{
		gk:             gk,
		scopes:         scopes,
		structTable:    structTable,
		interfaceTable: interfaceTable,
	}
}

// SymbolDefine defines a new symbol with a specified name and type, associating it with a struct or interface if applicable.
func (f *DefinitionTable) SymbolDefine(name string, typeName string) (*Symbol, error) {
	symbol, err := f.scopes.SymbolDefine(name)
	if err != nil {
		return nil, err
	}
	isStruct := f.structTable.Has(typeName)
	isInterface := f.interfaceTable.Has(typeName)
	symbol.SetReturnTypes([]string{typeName})
	if isStruct {
		f.structTable.BindSymbol(symbol, typeName)
		symbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+symbol.Name()))
	} else if isInterface {
		symbol.SetInterface(typeName)
		symbol.SetObject(f.gk.NewString(objects.FrameStatic, "interface:"+symbol.Name()))
	}
	return symbol, nil
}

// StructAdd adds a new struct definition with the given name to the struct table.
func (f *DefinitionTable) StructAdd(name string) {
	f.structTable.AddStruct(name)
}

// StructAddField adds a new field to an existing struct, creating the struct if it does not exist.
func (f *DefinitionTable) StructAddField(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	f.structTable.AddField(name, fieldName, baseStruct, kind, node)
}

// StructFieldsFromLiteral extracts fields from a composite literal for a specified struct and returns them as StructField instances.
func (f *DefinitionTable) StructFieldsFromLiteral(structName string, eltS []ast.Expr) ([]*StructField, error) {
	return f.structTable.FieldsFromLiteral(structName, eltS)
}

// StructImplements checks if a struct with the given name implements the specified interface by querying the struct table.
func (f *DefinitionTable) StructImplements(structName string, interfaceName string) bool {
	return f.structTable.Implements(structName, interfaceName)
}

// InterfaceAdd adds a new interface definition to the interface table with the given name and AST node representation.
func (f *DefinitionTable) InterfaceAdd(name string, node *ast.InterfaceType) error {
	return f.interfaceTable.Add(name, node)
}

// InterfaceGet retrieves an InterfaceDescription by name from the DefinitionTable's interfaceTable and its existence as a boolean.
func (f *DefinitionTable) InterfaceGet(name string) (*InterfaceDescription, bool) {
	return f.interfaceTable.Get(name)
}
