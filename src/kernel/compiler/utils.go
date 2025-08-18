package compiler

import (
	"fmt"
	"go/ast"
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
			if ident, ok := v.X.(*ast.Ident); ok {
				ret = append(ret, ident.Name)
			} else {
				// Questo caso gestirebbe tipi più complessi come '*[]Home'
				return nil, fmt.Errorf("unsupported pointer return type: *%T", v.X)
			}
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
