package compiler

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
		case *ast.FuncType:
			// Se il tipo di ritorno è una funzione, aggiungiamo un placeholder
			// generico. Potrebbe essere raffinato per includere i tipi dei parametri
			// se il tuo symbol system lo richiede.
			ret = append(ret, "func")
		case *ast.InterfaceType:
			// Gestisce il caso di `interface{}` come tipo di ritorno.
			if len(v.Methods.List) == 0 {
				ret = append(ret, "interface{}")
			} else {
				return nil, fmt.Errorf("unsupported non-empty interface return type")
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

func ExtractBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Caso base: abbiamo trovato l'identificatore del tipo (es. "MyStruct").
		return t.Name
	case *ast.StarExpr:
		// Caso puntatore (*MyType): continuiamo la ricerca sul tipo puntato.
		return ExtractBaseName(t.X)
	case *ast.ArrayType:
		// Caso array/slice ([]MyType): continuiamo la ricerca sul tipo dell'elemento.
		return ExtractBaseName(t.Elt)
	case *ast.MapType:
		// Caso mappa (map[KeyType]ValueType): ci interessa il tipo del valore.
		return ExtractBaseName(t.Value)
	case *ast.SelectorExpr:
		// Caso tipo qualificato (es. package.Type): restituiamo il nome del tipo.
		// Una logica più avanzata potrebbe restituire "package.Type".
		return t.Sel.Name
	case *ast.InterfaceType:
		// Caso interfaccia: se è un'interfaccia vuota, la trattiamo come un tipo a sé.
		if len(t.Methods.List) == 0 {
			return "interface{}"
		}
		// Le interfacce non vuote non sono attualmente supportate per l'estrazione.
		return ""
	default:
		// Altri tipi complessi non sono gestiti.
		return ""
	}
}

func NewCompilerError(fileSet *token.FileSet, node ast.Node, format string, args ...interface{}) error {
	// fileSet.Position() ci dà la posizione esatta del nodo nel file sorgente
	position := fileSet.Position(node.Pos())
	// Creiamo il messaggio di errore principale
	msg := fmt.Sprintf(format, args...)
	// Ritorniamo un errore formattato che include la posizione
	return fmt.Errorf("compile error at %s: %s", position.String(), msg)
}
