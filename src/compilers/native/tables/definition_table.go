package tables

import (
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// DefinitionTable represents a container for managing the scopes, structures, and interfaces in a type system.
// It utilizes a gatekeeper interface to handle object-related operations such as creation, conversion, and adaptation.
// The struct facilitates the organization and resolution of types within a given context or environment.
type DefinitionTable struct {
	gk             objects.IGateKeeper
	scopes         *Scopes
	structTable    *StructTable
	interfaceTable *InterfaceTable
}

// NewDefinitionTable initializes and returns a pointer to a new DefinitionTable with the provided gatekeeper, scopes, and tables.
func NewDefinitionTable(gk objects.IGateKeeper, scopes *Scopes) *DefinitionTable {
	structTable := NewStructTable(gk, scopes)
	interfaceTable := NewInterfaceTable(gk, scopes)
	return &DefinitionTable{
		gk:             gk,
		scopes:         scopes,
		structTable:    structTable,
		interfaceTable: interfaceTable,
	}
}

// SymbolDefine defines a new symbol with the given name and type, assigning it to a struct or interface as appropriate.
// Returns the created symbol or an error if the symbol could not be defined.
func (f *DefinitionTable) SymbolDefine(name string, typeName string) (*Symbol, error) {
	symbol, err := f.scopes.SymbolDefine(name)
	if err != nil {
		return nil, err
	}
	isInterface := f.interfaceTable.Has(typeName)
	if isInterface {
		f.SymbolInterfaceAssign(symbol, typeName)
	} else {
		f.SymbolTypeAssign(symbol, typeName)
	}
	return symbol, nil
}

// SymbolInferAssign assigns a type to a symbol by inferring the type from the provided expression or using the given type name.
func (f *DefinitionTable) SymbolInferAssign(symbol *Symbol, inferredTypeName string, rhsIn ast.Expr) {
	if len(inferredTypeName) == 0 {
		if inferredTypeName, _ = f.structTable.TypeInference(rhsIn); len(inferredTypeName) == 0 {
			return
		}
	}
	f.SymbolTypeAssign(symbol, inferredTypeName)
}

// SymbolInterfaceAssign assigns an interface to the given symbol and configures its return types and related object.
func (f *DefinitionTable) SymbolInterfaceAssign(symbol *Symbol, interfaceName string) {
	symbol.SetReturnTypes([]string{interfaceName})
	symbol.SetInterface(interfaceName)
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, interfaceName+":"+symbol.Name()))
}

// SymbolTypeAssign assigns a type to the provided symbol, sets its return types, and binds it to a struct if applicable.
func (f *DefinitionTable) SymbolTypeAssign(symbol *Symbol, typeName string) {
	symbol.SetReturnTypes([]string{typeName})
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+symbol.Name()))
	//assign struct if present in struct table
	f.structTable.BindSymbol(symbol, typeName)
}

// StructBindSymbol binds a given symbol to a struct type by setting its object and registering it in the struct table.
func (f *DefinitionTable) StructBindSymbol(symbol *Symbol, typeName string) {
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+symbol.Name()))
	f.structTable.BindSymbol(symbol, typeName)
}

// StructHas checks if a struct definition with the provided name exists in the struct table and returns a boolean.
func (f *DefinitionTable) StructHas(name string) bool {
	return f.structTable.Has(name)
}

// StructIsBuiltin checks if a struct with the given name is considered a built-in type in the struct table.
func (f *DefinitionTable) StructIsBuiltin(name string) bool {
	return f.structTable.IsBuiltin(name)
}

// StructAdd defines a new struct by its name and adds it to the struct table if it doesn't already exist.
func (f *DefinitionTable) StructAdd(name string) {
	f.structTable.AddStruct(name)
}

// StructAddField adds a new field to a struct, specifying the struct name, field name, base struct, type, and AST node.
func (f *DefinitionTable) StructAddField(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	f.structTable.AddField(name, fieldName, baseStruct, kind, node)
}

// StructFieldsFromLiteral retrieves struct fields from a given struct name and a slice of AST expressions.
// It returns a slice of StructField pointers or an error if the operation fails.
func (f *DefinitionTable) StructFieldsFromLiteral(structName string, eltS []ast.Expr) ([]*StructField, error) {
	return f.structTable.FieldsFromLiteral(structName, eltS)
}

// StructImplements checks whether a struct implements a specific interface by delegating to the underlying struct table.
func (f *DefinitionTable) StructImplements(structName string, interfaceName string) bool {
	return f.structTable.Implements(structName, interfaceName)
}

// StructTypeInference infers the type of a struct from the provided AST expression and returns the type name and success flag.
func (f *DefinitionTable) StructTypeInference(expr ast.Expr) (string, bool) {
	return f.structTable.TypeInference(expr)
}

// StructKeys retrieves the list of struct names from the StructTable associated with the DefinitionTable.
func (f *DefinitionTable) StructKeys() []string {
	return f.structTable.Keys()
}

// StructSetImplementations sets the implementation mappings for structs to interfaces in the DefinitionTable.
func (f *DefinitionTable) StructSetImplementations(impls map[string][]string) {
	f.structTable.SetImplementations(impls)
}

// StructReturnTypeFromSymbol retrieves the return type of a symbol as a string and a success flag based on its name.
func (f *DefinitionTable) StructReturnTypeFromSymbol(name string) (string, bool) {
	return f.structTable.ReturnTypeFromSymbol(name)
}

// StructTypeNameFromSymbolField retrieves the type name of a field within a struct associated with a given symbol name.
// Returns the type as a string and a boolean indicating success.
func (f *DefinitionTable) StructTypeNameFromSymbolField(name string, fieldName string) (string, bool) {
	return f.structTable.TypeNameFromSymbolField(name, fieldName)
}

// InterfaceAdd registers a new interface in the interface table using its name and ast.InterfaceType node, returning an error if the addition fails.
func (f *DefinitionTable) InterfaceAdd(name string, node *ast.InterfaceType) error {
	return f.interfaceTable.Add(name, node)
}

// InterfaceGet retrieves an InterfaceDescription by name from the interface table. Returns the description and a boolean indicating existence.
func (f *DefinitionTable) InterfaceGet(name string) (*InterfaceDescription, bool) {
	return f.interfaceTable.Get(name)
}

// InterfaceContainer retrieves the map of all interface descriptions indexed by their names from the interface table.
func (f *DefinitionTable) InterfaceContainer() map[string]*InterfaceDescription {
	return f.interfaceTable.Container()
}
