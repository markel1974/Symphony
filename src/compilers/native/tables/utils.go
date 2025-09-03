package tables

import (
	"fmt"
	"go/ast"
	"go/token"
)

// GetIdent retrieves the *ast.Ident type from an *ast.Field's Type if it exists, including handling pointer types.
func GetIdent(expr *ast.Field) *ast.Ident {
	switch t := expr.Type.(type) {
	case *ast.Ident:
		return t
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident
		}
	}
	return nil
}

// GetIdentName extracts and returns the name of the identifier from a given ast.Expr.
// It supports *ast.Ident, *ast.StarExpr, and *ast.SelectorExpr types.
func GetIdentName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.SelectorExpr:
		if receiverIdent, ok := t.X.(*ast.Ident); ok {
			return receiverIdent.Name
		}
	}
	return ""
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
		switch v := res.Type.(type) {
		case *ast.Ident:
			ret = append(ret, v.Name)
		case *ast.StarExpr:
			name := GetIdentName(v.X)
			ret = append(ret, name)
			//if ident, ok := v.X.(*ast.Ident); ok {
			//	ret = append(ret, ident.Name)
			//} else {
			// This case would handle more complex types like '*[]Home'
			//	return nil, fmt.Errorf("unsupported pointer return type: *%T", v.X)
			//}
		case *ast.FuncType:
			ret = append(ret, "func")
		case *ast.InterfaceType:
			// Handles the case of `interface{}` as return type
			if len(v.Methods.List) == 0 {
				ret = append(ret, "interface{}")
			} else {
				return nil, fmt.Errorf("unsupported non-empty interface return type")
			}
		case *ast.ArrayType:
			name := GetIdentName(v.Elt)
			ret = append(ret, "[]"+name)
		case *ast.MapType:
			key := GetIdentName(v.Key)
			value := GetIdentName(v.Value)
			ret = append(ret, "map["+key+"]"+value)
		default:
			return nil, fmt.Errorf("unsupported return type %T", v)
		}
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
