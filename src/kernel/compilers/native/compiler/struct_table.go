package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
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
	container        map[string][]*FieldDescription
	anonymousCounter int
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable() *StructTable {
	return &StructTable{
		container:        make(map[string][]*FieldDescription),
		anonymousCounter: 0,
	}
}

// CreateStructName creates a unique name for a struct based on the provided name.
func (s *StructTable) CreateStructName(name string) string {
	if len(name) == 0 {
		r := "<anonymous_" + strconv.Itoa(s.anonymousCounter) + ">"
		s.anonymousCounter++
		return r
	}
	return name
}

// Add adds a new field description to a struct in the StructTable. If the struct does not exist, it creates it.
func (s *StructTable) Add(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	// here we could add a check for duplicate fields.
	v := NewFieldDescription(fieldName, baseStruct, kind, node)
	fields, ok := s.container[name]
	if !ok {
		s.container[name] = []*FieldDescription{v}
		return
	}
	s.container[name] = append(fields, v)
}

// getFields retrieves a slice of StructProperty pointers associated with the given name from the container map.
func (s *StructTable) getFields(name string) ([]*FieldDescription, bool) {
	fields, ok := s.container[name]
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
func (s *StructTable) Has(name string) bool {
	if _, ok := s.container[name]; ok {
		return true
	}
	return false
}

// Inference infers struct type information from the given AST expression and scope context.
// It returns a generated struct name, a list of associated base type names, and a boolean indicating success.
func (s *StructTable) Inference(expr ast.Expr, scopes *Scopes) (string, []string, bool) {
	switch rhs := expr.(type) {
	case *ast.BinaryExpr:
		//nothing to do
		return "", nil, false
	case *ast.BasicLit:
		//nothing to do
		return "", nil, false
	case *ast.CompositeLit: // es. MyStruct{...}
		if baseName := ExtractBaseName(rhs.Type); len(baseName) > 0 {
			return s.CreateStructName(baseName), []string{baseName}, true
		}
	case *ast.CallExpr: // es. NewStruct()
		if ident, ok := rhs.Fun.(*ast.Ident); ok {
			if funcSymbol, ok := scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.Types()) > 0 {
				// Assumiamo il primo tipo di ritorno
				typeName := funcSymbol.Types()[0]
				// Verifichiamo se il tipo restituito è uno struct
				if typeSymbol, ok := scopes.SymbolResolve(typeName); ok && typeSymbol.IsStruct() {
					return s.CreateStructName(typeName), []string{typeName}, true
				}
			}
		}
	case *ast.UnaryExpr: // es. &MyStruct{}
		if rhs.Op == token.AND {
			if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
				if ident, ok := compLit.Type.(*ast.Ident); ok {
					if typeSymbol, ok := scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
						return s.CreateStructName(typeSymbol.Name()), []string{typeSymbol.Name()}, true
					}
				}
			}
		}
	}
	return "", nil, false
}

// CreateSymbolFromLiteral creates a symbol and field descriptions from a given composite literal and scope context.
func (s *StructTable) CreateSymbolFromLiteral(node *ast.CompositeLit, scopes *Scopes) (*Symbol, []*FieldDescription, error) {
	// struct literal (es. MyStruct{...})
	t, ok := node.Type.(*ast.Ident)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported composite literal type: %T", node)
	}
	structFields, ok := s.getFields(t.Name)
	if !ok {
		return nil, nil, fmt.Errorf("unknown composite literal type: %s", t.Name)
	}
	if len(node.Elts) > len(structFields) {
		return nil, nil, fmt.Errorf("too many values in positional struct literal for type '%s'", t.Name)
	}
	symbol, ok := scopes.SymbolResolve(t.Name)
	if !ok {
		var err error
		if symbol, err = scopes.SymbolDefine(t.Name); err != nil {
			return nil, nil, err
		}
	}
	structName := s.CreateStructName(t.Name)
	symbol.SetStruct(structName)
	symbol.SetTypes([]string{t.Name})
	//TODO
	//symbol.SetObject()
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
func (s *StructTable) GetTypeNameFromFields(structName string, fieldName string) (string, bool) {
	receiverStructFields, ok := s.getFields(structName)
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
