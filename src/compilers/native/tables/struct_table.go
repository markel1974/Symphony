package tables

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// StructTable is a collection that manages mappings of struct names to their associated properties.
type StructTable struct {
	container       map[string]*Struct
	gk              objects.IGateKeeper
	scopes          *Scopes
	implementations map[string][]string
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable(gk objects.IGateKeeper, scopes *Scopes) *StructTable {
	st := &StructTable{
		container:       make(map[string]*Struct),
		implementations: make(map[string][]string),
		gk:              gk,
		scopes:          scopes,
	}
	builtins := []string{"error"}
	for _, builtin := range builtins {
		//z := NewStructField(internal, "", internal, nil)
		sd := NewStruct(builtin, StructTypeBuiltin)
		//sd.Add(z)
		st.container[sd.name] = sd
	}
	return st
}

// Keys returns a slice of struct names in the container map.
func (st *StructTable) Keys() []string {
	keys := make([]string, 0, len(st.container))
	for k := range st.container {
		keys = append(keys, k)
	}
	return keys
}

// SetImplementations sets the implementation mappings for structs to interfaces in the StructTable.
func (st *StructTable) SetImplementations(impls map[string][]string) {
	st.implementations = impls
}

// Implements verifica se uno struct implementa una data interfaccia.
func (st *StructTable) Implements(structName, interfaceName string) bool {
	if impls, ok := st.implementations[structName]; ok {
		for _, iName := range impls {
			if iName == interfaceName {
				return true
			}
		}
	}
	return false
}

// AddExternal adds a new package with the given name to the StructTable if it does not already exist.
func (st *StructTable) AddExternal(name string) {
	if _, ok := st.container[name]; !ok {
		sd := NewStruct(name, StructTypePackage)
		st.container[name] = sd
	}
}

// Add adds a new field description to a struct in the StructTable. If the struct does not exist, it creates it.
func (st *StructTable) Add(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	sd, ok := st.container[name]
	if !ok {
		sd = NewStruct(name, StructTypeDefined)
		st.container[name] = sd
	}
	sd.AddField(fieldName, baseStruct, kind, node)
}

// Has checks if a struct definition with the given name exists in the container map.
func (st *StructTable) Has(name string) bool {
	if _, ok := st.container[name]; ok {
		return true
	}
	return false
}

// TypeInference infers struct type information from the given AST expression and scope context.
// It returns a generated struct name, a list of associated base type names, and a boolean indicating success.
func (st *StructTable) TypeInference(expr ast.Expr) (string, bool) {
	var ret string
	switch rhs := expr.(type) {
	case *ast.Ident:
		//nothing to do
		return "", false
	case *ast.BinaryExpr:
		//nothing to do
		return "", false
	case *ast.BasicLit:
		//nothing to do
		return "", false
	case *ast.CompositeLit: // es. MyStruct{...}
		ret = st.ExtractBaseName(rhs.Type)
		//if baseName := st.ExtractBaseName(rhs.Type); len(baseName) > 0 {
		//	return baseName, []string{baseName}, true
		//}
	case *ast.CallExpr: // es. NewStruct()
		if ident, ok := rhs.Fun.(*ast.Ident); ok {
			if funcSymbol, ok := st.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.ReturnTypes()) > 0 {
				// We assume the first return type
				returnType := funcSymbol.ReturnTypes()[0]
				// Verify if the returned type is a struct
				if typeSymbol, ok := st.scopes.SymbolResolve(returnType); ok && typeSymbol.IsStruct() {
					ret = returnType
					//return returnType, []string{returnType}, true
				}
			}
		}
	case *ast.UnaryExpr: // es. &MyStruct{}
		if rhs.Op == token.AND {
			if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
				if ident, ok := compLit.Type.(*ast.Ident); ok {
					if returnSymbol, ok := st.scopes.SymbolResolve(ident.Name); ok && returnSymbol.IsStruct() {
						ret = returnSymbol.Name()
						//return returnSymbol.Name(), []string{returnSymbol.Name()}, true
					}
				}
			}
		}
	case *ast.TypeAssertExpr:
		targetTypeName := rhs.Type.(*ast.Ident).Name
		if returnSymbol, ok := st.scopes.SymbolResolve(targetTypeName); ok && returnSymbol.IsStruct() {
			ret = returnSymbol.Name()
			//return returnSymbol.Name(), []string{returnSymbol.Name()}, true
		}
	}
	if len(ret) == 0 {
		return "", false
	}
	return ret, true
	//return "", nil, false
}

// FieldsFromLiteral extracts and assigns struct fields from a given composite literal node, handling both keyed and positional formats.
func (st *StructTable) FieldsFromLiteral(structName string, eltS []ast.Expr) ([]*StructField, error) {
	sd, ok := st.container[structName]
	if !ok {
		return nil, fmt.Errorf("unknown composite literal type: %st", structName)
	}
	structFields := sd.Fields()
	if len(eltS) > len(structFields) {
		return nil, fmt.Errorf("too many values in positional struct literal for type '%st'", structName)
	}
	symbol, ok := st.scopes.SymbolResolve(structName)
	if !ok {
		var err error
		if symbol, err = st.scopes.SymbolDefine(structName); err != nil {
			return nil, err
		}
	}
	symbol.SetReturnTypes([]string{structName})
	symbol.SetObject(st.gk.NewString(objects.FrameStatic, structName+":"+symbol.Name()))
	st.BindSymbol(symbol, structName)

	isKeyed := false
	if len(eltS) > 0 {
		if _, ok := eltS[0].(*ast.KeyValueExpr); ok {
			isKeyed = true
		}
	}
	if isKeyed {
		// key literal (es. Home{Name: "Alfa", Address: "Shanghai"})
		providedFields := make(map[string]ast.Expr)
		for _, elt := range eltS {
			kvExpr, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				return nil, fmt.Errorf("cannot mix keyed and unkeyed values in struct literal")
			}
			keyIdent, ok := kvExpr.Key.(*ast.Ident)
			if !ok {
				return nil, fmt.Errorf("invalid field name in struct literal")
			}
			providedFields[keyIdent.Name] = kvExpr.Value
		}
		for idx := range structFields {
			if valueExpr, ok := providedFields[structFields[idx].name]; ok {
				structFields[idx].node = valueExpr
			}
		}
	} else {
		// positional literal (es. Home{"Alfa", 20, "Shanghai"}) ---
		for i, elt := range eltS {
			structFields[i].node = elt
		}
	}
	return structFields, nil
}

// ReturnTypeFromSymbol resolves the return type of a symbol by its name and returns it along with a success flag.
func (st *StructTable) ReturnTypeFromSymbol(name string) (string, bool) {
	symbol, ok := st.scopes.SymbolResolve(name)
	if ok && len(symbol.ReturnTypes()) > 0 {
		return symbol.ReturnTypes()[0], true
	}
	return "", false
}

// TypeNameFromSymbolField retrieves the base type of a field within a struct using its name and returns it with a success flag.
func (st *StructTable) TypeNameFromSymbolField(name string, fieldName string) (string, bool) {
	receiverSymbol, ok := st.scopes.SymbolResolve(name)
	if !ok {
		return "", false
	}
	sd, ok := st.container[receiverSymbol.StructName()]
	if !ok {
		return "", false
	}
	for _, receiverField := range sd.Fields() {
		if receiverField.name == fieldName {
			return receiverField.base, true
		}
	}
	return "", false
}

// BindSymbol assigns a struct name and types to a Symbol, validates the struct, and creates a description object.
func (st *StructTable) BindSymbol(symbol *Symbol, typeName string) {
	if sd, ok := st.container[typeName]; ok {
		fields := sd.FieldsName()
		symbol.SetStruct(typeName, fields)
		return
	}
}

// IsBuiltin returns true if the given name is a struct internal to the compiler.
func (st *StructTable) IsBuiltin(name string) bool {
	fd, ok := st.container[name]
	if !ok {
		return false
	}
	return fd.IsBuiltin()
}

// ExtractBaseName extracts the base type name from an AST expression, handling pointers, arrays, maps, and selectors.
func (st *StructTable) ExtractBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Base case: we found the type identifier (e.g. "MyStruct")
		return t.Name
	case *ast.StarExpr:
		// Pointer case (*MyType): continue search on pointed type
		return st.ExtractBaseName(t.X)
	case *ast.ArrayType:
		// Array/slice case ([]MyType): continue search on element type
		return st.ExtractBaseName(t.Elt)
	case *ast.MapType:
		// Map case (map[KeyType]ValueType): we care about the value type
		return st.ExtractBaseName(t.Value)
	case *ast.SelectorExpr:
		// Qualified type case (e.g. package.Type): return the type name
		// More advanced logic could return "package.Type"
		return t.Sel.Name
	case *ast.InterfaceType:
		// Interface case: treat empty interface as its own type
		if len(t.Methods.List) == 0 {
			return "interface{}"
		}
		// Non-empty interfaces are not currently supported for extraction
		return ""
	default:
		// Other complex types are not handled
		return ""
	}
}
