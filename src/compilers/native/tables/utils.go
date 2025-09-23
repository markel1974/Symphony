package tables

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	FuncDefinition      = "func"
	ArrayDefinition     = "[]"
	MapDefinition       = "map"
	InterfaceDefinition = "interface{}"
)

// GetIdent traverses an AST expression and returns the identifier (*ast.Ident) of the type, if found.
// It handles identifiers, pointers, arrays/slices, and maps, returning the type's identifier or nil if not applicable.
func GetIdent(expr ast.Expr) *ast.Ident {
	switch t := expr.(type) {
	case *ast.Ident:
		// Base case: we found the identifier (e.g. "MyStruct")
		return t
	case *ast.StarExpr:
		// Pointer case (*MyType): continue searching on the pointed type
		return GetIdent(t.X)
	case *ast.ArrayType:
		// Array/slice case ([]MyType): continue searching on the element type
		return GetIdent(t.Elt)
	case *ast.MapType:
		// Map case (map[KeyType]ValueType): we are interested in the value type
		return GetIdent(t.Value)
	case *ast.CompositeLit:
		// Composite literal case (MyStruct{...}): continue searching on the type identifier
		return GetIdent(t.Type)
	case *ast.SelectorExpr:
		// Qualified type case (package.Type): return the type identifier
		return t.Sel
	case *ast.CallExpr:
		return GetIdent(t.Fun)
	case *ast.UnaryExpr:
		return GetIdent(t.X)
	case *ast.ParenExpr:
		// Parenthesized expression ((MyType)): continue search on the inner expression
		return GetIdent(t.X)
	}
	return nil
}

// GetReceiver extracts the type name from the given AST Field's Type and returns it as a string or an error for unsupported types.
func GetReceiver(result *ast.Field) (string, error) {
	switch v := result.Type.(type) {
	case *ast.Ident:
		return v.Name, nil
	case *ast.StarExpr:
		z := GetIdent(v.X)
		if z == nil {
			return "", fmt.Errorf("unsupported pointer return type: *%T", v.X)
		}
		return z.Name, nil
	case *ast.FuncType:
		return FuncDefinition, nil
	case *ast.InterfaceType:
		if len(v.Methods.List) == 0 {
			return InterfaceDefinition, nil
		} else {
			return "", fmt.Errorf("unsupported non-empty interface return type")
		}
	case *ast.ArrayType:
		z := GetIdent(v.Elt)
		if z == nil {
			return "", fmt.Errorf("unsupported array return type: %T", v.Elt)
		}
		return ArrayDefinition + z.Name, nil
	case *ast.MapType:
		key := GetIdent(v.Key)
		value := GetIdent(v.Value)
		if key == nil || value == nil {
			return "", fmt.Errorf("unsupported map return type: %T", v)
		}
		return MapDefinition + "[" + key.Name + "]" + value.Name, nil
	default:
		return "", fmt.Errorf("unsupported return type %T", v)
	}
}

// GetBaseName extracts the base type name from an AST expression, handling pointers, arrays, maps, and selectors.
func GetBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Base case: we found the type identifier (e.g. "MyStruct")
		return t.Name
	case *ast.StarExpr:
		// Pointer case (*MyType): continue search on pointed type
		return GetBaseName(t.X)
	case *ast.ArrayType:
		// Array/slice case ([]MyType): continue search on element type
		return GetBaseName(t.Elt)
	case *ast.MapType:
		// Map case (map[KeyType]ValueType): we care about the value type
		return GetBaseName(t.Value)
	case *ast.SelectorExpr:
		// Qualified type case (e.g. package.Type): return the type name
		// More advanced logic could return "package.Type"
		return t.Sel.Name
	case *ast.InterfaceType:
		// Interface case: treat empty interface as its own type
		if len(t.Methods.List) == 0 {
			return InterfaceDefinition
		}
		// Non-empty interfaces are not currently supported for extraction
		return ""
	default:
		// Other complex types are not handled
		return ""
	}
}

// GetReceivers extracts and returns the list of type names from the given AST FieldList result.
// It handles both non-pointer and pointer type fields and returns an error for unsupported types.
func GetReceivers(result *ast.FieldList) ([]string, error) {
	if result == nil {
		return nil, nil
	}
	if len(result.List) == 0 {
		return nil, nil
	}
	var ret []string
	for _, res := range result.List {
		rec, err := GetReceiver(res)
		if err != nil {
			return nil, err
		}
		ret = append(ret, rec)
	}
	return ret, nil
}

// GetMangledName combines an identifier and function name to generate a mangled name in the format "identifier.function".
func GetMangledName(identId string, fnName string) string {
	m := fmt.Sprintf("%s.%s", identId, fnName)
	return m
}

func NewCompilerError(fileSet *token.FileSet, node ast.Node, format string, args ...interface{}) error {
	// fileSet.Position() ci dà la posizione esatta del nodo nel file sorgente
	position := fileSet.Position(node.Pos())
	// Creiamo il messaggio di errore principale
	msg := fmt.Sprintf(format, args...)
	// Ritorniamo un errore formattato che include la posizione
	return fmt.Errorf("compile error at %s: %s", position.String(), msg)
}
