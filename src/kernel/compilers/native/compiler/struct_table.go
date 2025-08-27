package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// FieldDescription represents metadata about a struct field, including its name, base type, full type, and AST node.
type FieldDescription struct {
	name string
	base string
	kind string
	node ast.Node
}

// NewFieldDescription creates a new instance of StructProperty with the provided name, base, kind, and AST node.
func NewFieldDescription(name string, base string, kind string, node ast.Node) *FieldDescription {
	return &FieldDescription{
		name: name,
		base: base,
		kind: kind,
		node: node,
	}
}

// StructTable is a collection that manages mappings of struct names to their associated properties.
type StructTable struct {
	container map[string][]*FieldDescription
	gk        objects.IGateKeeper
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable(gk objects.IGateKeeper) *StructTable {
	st := &StructTable{
		container: make(map[string][]*FieldDescription),
		gk:        gk,
	}
	return st
}

// Add adds a new field description to a struct in the StructTable. If the struct does not exist, it creates it.
func (st *StructTable) Add(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	// here we could add a check for duplicate fields.
	v := NewFieldDescription(fieldName, baseStruct, kind, node)
	fields, ok := st.container[name]
	if !ok {
		st.container[name] = []*FieldDescription{v}
		return
	}
	st.container[name] = append(fields, v)
}

// getFields retrieves a slice of StructProperty pointers associated with the given name from the container map.
func (st *StructTable) getFields(name string) ([]*FieldDescription, bool) {
	fields, ok := st.container[name]
	if !ok {
		return nil, false
	}
	out := make([]*FieldDescription, len(fields))
	for idx, v := range fields {
		out[idx] = NewFieldDescription(v.name, v.base, v.kind, nil)
	}
	return out, true
}

// Has checks if a struct definition with the given name exists in the container map.
func (st *StructTable) Has(name string) bool {
	if _, ok := st.container[name]; ok {
		return true
	}
	return false
}

// Inference infers struct type information from the given AST expression and scope context.
// It returns a generated struct name, a list of associated base type names, and a boolean indicating success.
func (st *StructTable) Inference(expr ast.Expr, scopes *Scopes) (string, []string, bool) {
	switch rhs := expr.(type) {
	case *ast.BinaryExpr:
		//nothing to do
		return "", nil, false
	case *ast.BasicLit:
		//nothing to do
		return "", nil, false
	case *ast.CompositeLit: // es. MyStruct{...}
		if baseName := ExtractBaseName(rhs.Type); len(baseName) > 0 {
			return baseName, []string{baseName}, true
		}
	case *ast.CallExpr: // es. NewStruct()
		if ident, ok := rhs.Fun.(*ast.Ident); ok {
			if funcSymbol, ok := scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.Types()) > 0 {
				// Assumiamo il primo tipo di ritorno
				typeName := funcSymbol.Types()[0]
				// Verifichiamo se il tipo restituito è uno struct
				if typeSymbol, ok := scopes.SymbolResolve(typeName); ok && typeSymbol.IsStruct() {
					return typeName, []string{typeName}, true
				}
			}
		}
	case *ast.UnaryExpr: // es. &MyStruct{}
		if rhs.Op == token.AND {
			if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
				if ident, ok := compLit.Type.(*ast.Ident); ok {
					if typeSymbol, ok := scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
						return typeSymbol.Name(), []string{typeSymbol.Name()}, true
					}
				}
			}
		}
	}
	return "", nil, false
}

// SymbolFromLiteral creates a symbol and field descriptions from a given composite literal and scope context.
func (st *StructTable) SymbolFromLiteral(node *ast.CompositeLit, scopes *Scopes) (*Symbol, []*FieldDescription, error) {
	// struct literal (es. MyStruct{...})
	t, ok := node.Type.(*ast.Ident)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported composite literal type: %T", node)
	}
	structFields, ok := st.getFields(t.Name)
	if !ok {
		return nil, nil, fmt.Errorf("unknown composite literal type: %st", t.Name)
	}
	if len(node.Elts) > len(structFields) {
		return nil, nil, fmt.Errorf("too many values in positional struct literal for type '%st'", t.Name)
	}
	symbol, ok := scopes.SymbolResolve(t.Name)
	if !ok {
		var err error
		if symbol, err = scopes.SymbolDefine(t.Name); err != nil {
			return nil, nil, err
		}
	}
	if err := st.AssignSymbol(symbol, t.Name, []string{t.Name}); err != nil {
		return nil, nil, err
	}
	isKeyed := false
	if len(node.Elts) > 0 {
		if _, ok := node.Elts[0].(*ast.KeyValueExpr); ok {
			isKeyed = true
		}
	}
	if isKeyed {
		// key literal (es. Home{Name: "Alfa", Address: "Shanghai"})
		providedFields := make(map[string]ast.Expr)
		for _, elt := range node.Elts {
			kvExpr, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				return nil, nil, fmt.Errorf("cannot mix keyed and unkeyed values in struct literal")
			}
			keyIdent, ok := kvExpr.Key.(*ast.Ident)
			if !ok {
				return nil, nil, fmt.Errorf("invalid field name in struct literal")
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
		for i, elt := range node.Elts {
			structFields[i].node = elt
		}
	}
	return symbol, structFields, nil
}

// GetTypeNameFromFields retrieves the base type of a field within a struct using its name and returns it with a success flag.
func (st *StructTable) GetTypeNameFromFields(structName string, fieldName string) (string, bool) {
	receiverStructFields, ok := st.getFields(structName)
	if !ok {
		return "", false
	}
	for _, receiverField := range receiverStructFields {
		if receiverField.name == fieldName {
			return receiverField.base, true
		}
	}
	return "", false
}

// AssignSymbol assigns a struct name and types to a Symbol, validates the struct, and creates a description object.
func (st *StructTable) AssignSymbol(symbol *Symbol, structName string, types []string) error {
	if len(structName) == 0 {
		return fmt.Errorf("empty struct type")
	}
	//if structName != "interface{}" {
	//	if !st.Has(structName) {
	//		return fmt.Errorf("unknown struct type: %s", structName)
	//	}
	//}
	description := structName + "=>" + symbol.Name() + ":" + strings.Join(types, " ")
	symbol.SetStruct(structName)
	symbol.SetTypes(types)
	symbol.SetObject(st.gk.NewString(objects.FrameStatic, description))
	return nil
}
