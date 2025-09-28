package tables

import (
	"fmt"
	"go/ast"
	"go/token"
)

// CharDefinition represents the keyword for character data type.
// IntDefinition represents the keyword for integer data type.
// FloatDefinition represents the keyword for floating point data type.
// StringDefinition represents the keyword for string data type.
// FuncDefinition represents the keyword for function type.
// ArrayDefinition represents the keyword for array type syntax.
// MapDefinition represents the keyword for map type syntax.
// InterfaceDefinition represents the keyword for empty interface type syntax.
const (
	CharDefinition      = "char"
	IntDefinition       = "int"
	FloatDefinition     = "float"
	StringDefinition    = "string"
	FuncDefinition      = "func"
	ArrayDefinition     = "[]"
	MapDefinition       = "map"
	InterfaceDefinition = "interface{}"
)

// _charDefinition holds a reference to an identifier representing the char type constant definition.
var _charDefinition = &ast.Ident{NamePos: 0, Name: CharDefinition, Obj: nil}

// _intDefinition is a pre-defined identifier representing the "int" type for use in AST manipulation or type matching.
var _intDefinition = &ast.Ident{NamePos: 0, Name: IntDefinition, Obj: nil}

// _floatDefinition is an identifier representing the "float" type used in AST type handling.
var _floatDefinition = &ast.Ident{NamePos: 0, Name: FloatDefinition, Obj: nil}

// _stringDefinition is a predefined ast.Ident representing the "string" type definition in the abstract syntax tree.
var _stringDefinition = &ast.Ident{NamePos: 0, Name: StringDefinition, Obj: nil}

// _interfaceDefinition represents the identifier for the empty interface type ("interface{}") in the abstract syntax tree.
var _interfaceDefinition = &ast.Ident{NamePos: 0, Name: InterfaceDefinition, Obj: nil}

// GetSelectorData extracts the package and selector names from a selector expression, returning them with a success flag.
func GetSelectorData(expr ast.Expr) (string, string, bool) {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		name := GetIdent(t.X)
		if name == nil {
			return "", "", false
		}
		return name.Name, t.Sel.Name, true
	default:
		return "", "", false
	}
}

// GetIdentName returns the name of the identifier from the given AST expression, or an empty string if no identifier is found.
func GetIdentName(expr ast.Expr) string {
	ident := GetIdent(expr)
	if ident == nil {
		return ""
	}
	return ident.Name
}

// GetIdent extracts and returns the *ast.Ident representing a type or identifier within an ast.Expr.
// It traverses various expression types such as pointers, arrays, slices, maps, and composite literals.
// If no valid identifier is found, the function returns nil.
func GetIdent(expr ast.Expr) *ast.Ident {
	switch t := expr.(type) {
	case *ast.BasicLit:
		switch t.Kind {
		case token.INT:
			return _intDefinition
		case token.FLOAT:
			return _floatDefinition
		case token.CHAR:
			return _charDefinition
		case token.STRING:
			return _stringDefinition
		default:
			return nil
		}
	case *ast.Ident:
		// Base case: we found the identifier (e.g. "MyStruct")
		return t
	case *ast.StarExpr:
		// Pointer case (*MyType): continue searching on the pointed type
		return GetIdent(t.X)
	case *ast.SliceExpr:
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
	case *ast.BinaryExpr:
		return GetIdent(t.X) // or GetIdent(t.Y)
	case *ast.ParenExpr:
		// Parenthesized expression ((MyType)): continue search on the inner expression
		return GetIdent(t.X) // or GetIdent(t.Y)
	case *ast.InterfaceType:
		// Interface case: treat empty interface as its own type
		if len(t.Methods.List) == 0 {
			return _interfaceDefinition
		}
		return nil
	default:
		return nil
	}
}

// GetReceiver determines the string representation of a receiver type based on its AST expression node.
// It returns the type name as a string or an error for unsupported types.
func GetReceiver(expr ast.Expr) (string, error) {
	switch v := expr.(type) {
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

// GetReceivers extracts the names of types from a given *ast.FieldList and returns them as a slice of strings.
// Returns an error if type extraction encounters an issue.
func GetReceivers(result *ast.FieldList) ([]string, error) {
	if result == nil {
		return nil, nil
	}
	if len(result.List) == 0 {
		return nil, nil
	}
	var ret []string
	for _, res := range result.List {
		rec, err := GetReceiver(res.Type)
		if err != nil {
			return nil, err
		}
		ret = append(ret, rec)
	}
	return ret, nil
}

// GetFuncName extracts the function name or mangled name from an ast.Expr and returns it with a boolean indicating success.
func GetFuncName(expr ast.Expr) (string, bool) {
	switch fun := expr.(type) {
	case *ast.Ident:
		return fun.Name, true
	case *ast.SelectorExpr:
		receiverIdent, ok := fun.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		return GetMangledName(receiverIdent.Name, fun.Sel.Name), true
	default:
		return "", false
	}
}

// GetMangledName returns a string combining an identifier and a function name in the format "Identifier.FunctionName".
func GetMangledName(identId string, fnName string) string {
	m := fmt.Sprintf("%s.%s", identId, fnName)
	return m
}

// NewCompilerError creates a compile-time error message with the file and node position, formatted with the given string and arguments.
func NewCompilerError(fileSet *token.FileSet, node ast.Node, format string, args ...interface{}) error {
	// fileSet.Position() ci dà la posizione esatta del nodo nel file sorgente
	position := fileSet.Position(node.Pos())
	// Creiamo il messaggio di errore principale
	msg := fmt.Sprintf(format, args...)
	// Ritorniamo un errore formattato che include la posizione
	return fmt.Errorf("compile error at %s: %s", position.String(), msg)
}
