package tables

import (
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

type DefinitionTable struct {
	gk             objects.IGateKeeper
	scopes         *Scopes
	structTable    *StructTable
	interfaceTable *InterfaceTable
}

func NewDefinitionTable(gk objects.IGateKeeper, scopes *Scopes, structTable *StructTable, interfaceTable *InterfaceTable) *DefinitionTable {
	return &DefinitionTable{
		gk:             gk,
		scopes:         scopes,
		structTable:    structTable,
		interfaceTable: interfaceTable,
	}
}

// SymbolDefine defines a new Symbol in the current scope with the specified name and type.
// It associates the symbol with a struct or interface if applicable.
// Returns the defined Symbol or an error if the operation fails.
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

// StructAdd registers a new struct by its name in the struct table.
func (f *DefinitionTable) StructAdd(name string) {
	f.structTable.AddStruct(name)
}

// StructAddField adds a new field to a struct using the specified field name, base struct, kind, and AST node.
func (f *DefinitionTable) StructAddField(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	f.structTable.AddField(name, fieldName, baseStruct, kind, node)
}

// StructFieldsFromLiteral extracts fields of a struct from a composite literal and returns them along with any errors encountered.
func (f *DefinitionTable) StructFieldsFromLiteral(structName string, eltS []ast.Expr) ([]*StructField, error) {
	return f.structTable.FieldsFromLiteral(structName, eltS)
}

// StructImplements determines if a struct implements a given interface based on the definition table.
func (f *DefinitionTable) StructImplements(structName string, interfaceName string) bool {
	return f.structTable.Implements(structName, interfaceName)
}

// InterfaceAdd registers a new interface in the DefinitionTable using the provided name and AST node; returns an error on failure.
func (f *DefinitionTable) InterfaceAdd(name string, node *ast.InterfaceType) error {
	return f.interfaceTable.Add(name, node)
}

// InterfaceGet retrieves an InterfaceDescription and a boolean indicating its existence by the provided name.
// Returns the InterfaceDescription if found, otherwise returns false along with a nil description.
func (f *DefinitionTable) InterfaceGet(name string) (*InterfaceDescription, bool) {
	return f.interfaceTable.Get(name)
}
