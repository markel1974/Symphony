package tables

import "github.com/markel1974/c64emu/src/vm/objects"

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
