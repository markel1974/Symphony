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

func GetMangledName(identId string, fnName string) string {
	m := fmt.Sprintf("%s.%s", identId, fnName)
	return m
}
