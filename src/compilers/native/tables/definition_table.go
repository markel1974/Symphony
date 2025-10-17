package tables

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// DefinitionTable represents a container for managing the scopes, structures, and interfaces in a type system.
// It uses a gatekeeper interface to handle object-related operations such as creation, conversion, and adaptation.
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

// Assign assigns a type or interface to the given symbol based on its presence in the interface or struct table.
func (f *DefinitionTable) Assign(symbol *Symbol, typeName string) error {
	isInterface := f.interfaceTable.Has(typeName)
	if isInterface {
		f.InterfaceAssign(symbol, typeName)
	} else {
		f.TypeAssign(symbol, typeName)
	}
	return nil
}

// InferAssign assigns a type to a symbol by inferring the type from the provided expression or using the given type name.
func (f *DefinitionTable) InferAssign(symbol *Symbol, inferredTypeName string, rhsIn ast.Expr) {
	if len(inferredTypeName) == 0 {
		if inferredTypeName, _ = f.TypeInference(rhsIn); len(inferredTypeName) == 0 {
			return
		}
	}
	f.TypeAssign(symbol, inferredTypeName)
}

// InterfaceAssign assigns an interface to the given symbol and configures its return types and related object.
func (f *DefinitionTable) InterfaceAssign(symbol *Symbol, interfaceName string) {
	symbol.SetReturnTypes([]string{interfaceName})
	symbol.SetInterface(interfaceName)
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, interfaceName+":"+symbol.Name()))
}

// TypeAssign assigns a type to the provided symbol, sets its return types, and binds it to a struct if applicable.
func (f *DefinitionTable) TypeAssign(symbol *Symbol, typeName string) {
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
	sd := f.structTable.AddStruct(name)
	var walker IWalker = nil
	if w, ok := f.structTable.Walker(baseStruct); ok {
		walker = w
	} else if w, ok = f.interfaceTable.Walker(baseStruct); ok {
		walker = w
	}
	sf := NewStructField(fieldName, baseStruct, kind, walker, node)
	sd.AddField(sf)
}

// StructKeys retrieves the list of struct names from the StructTable associated with the DefinitionTable.
func (f *DefinitionTable) StructKeys() []string {
	return f.structTable.Keys()
}

// StructSetImplementations sets the implementation mappings for structs to interfaces in the DefinitionTable.
func (f *DefinitionTable) StructSetImplementations(impls map[string][]string) {
	f.structTable.SetImplementations(impls)
}

// StructFieldsFromLiteral extracts struct fields from a list of AST expressions for a given struct name, returning them or an error.
func (f *DefinitionTable) StructFieldsFromLiteral(structName string, eltS []ast.Expr) ([]*StructField, error) {
	return f.structTable.FieldsFromLiteral(structName, eltS)
}

// StructTypeNameFromSymbolField retrieves the type name of a field within a struct associated with a given symbol name.
// Returns the type as a string and a boolean indicating success.
func (f *DefinitionTable) StructTypeNameFromSymbolField(name string, fieldName string) (string, bool) {
	return f.structTable.TypeNameFromSymbolField(name, fieldName)
}

// StructImplements checks whether a struct implements a specific interface by delegating to the underlying struct table.
func (f *DefinitionTable) StructImplements(structName string, interfaceName string) bool {
	return f.structTable.Implements(structName, interfaceName)
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

// TypeInference infers struct type information from the given AST expression and scope context.
// It returns a generated struct name, a list of associated base type names, and a boolean indicating success.
func (f *DefinitionTable) TypeInference(expr ast.Expr) (string, bool) {
	var ret string
	switch rhs := expr.(type) {
	case *ast.Ident:
		if v, ok := f.scopes.SymbolResolve(rhs.Name); ok {
			ret = v.StructName()
		}
	case *ast.CompositeLit: // es. MyStruct{...}
		if ident := GetIdent(rhs.Type); ident != nil {
			ret = ident.Name
		}
	case *ast.CallExpr: // es. NewStruct()
		if ident := GetIdent(rhs.Fun); ident != nil {
			if funcSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.ReturnTypes()) > 0 {
				// We assume the first return type
				returnType := funcSymbol.ReturnTypes()[0]
				// Verify if the returned type is a struct
				if typeSymbol, ok := f.scopes.SymbolResolve(returnType); ok && typeSymbol.IsStruct() {
					ret = returnType
				}
			}
		}
	case *ast.UnaryExpr: // es. &MyStruct{}
		if rhs.Op == token.AND {
			if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
				if ident := GetIdent(compLit.Type); ident != nil {
					if returnSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok && returnSymbol.IsStruct() {
						ret = returnSymbol.Name()
					}
				}
			}
		}
	case *ast.TypeAssertExpr:
		if ident := GetIdent(rhs.Type); ident != nil {
			if returnSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok && returnSymbol.IsStruct() {
				ret = returnSymbol.Name()
			}
		}
	}
	if len(ret) == 0 {
		return "", false
	}
	return ret, true
	//return "", nil, false
}
