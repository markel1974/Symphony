package compiler

import (
	"fmt"
	"go/ast"
)

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

func GetReceiver(result *ast.FieldList) (string, error) {
	if result == nil {
		return "", nil
	}
	if len(result.List) == 0 {
		return "", nil
	}
	kind := result.List[0].Type
	switch v := kind.(type) {
	case *ast.Ident:
		return v.Name, nil
	case *ast.StarExpr:
		if ident, ok := v.X.(*ast.Ident); ok {
			return ident.Name, nil
		} else {
			// Questo caso gestirebbe tipi più complessi come '*[]Home'
			return "", fmt.Errorf("unsupported pointer return type: *%T", v.X)
		}
	default:
		return "", fmt.Errorf("unsupported return type %T", kind)
	}
}

func GetMangledName(identId string, fnName string) string {
	m := fmt.Sprintf("%s.%s", identId, fnName)
	return m
}
