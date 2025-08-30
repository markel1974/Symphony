package tables

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

func (fd *FieldDescription) Name() string {
	return fd.name
}

func (fd *FieldDescription) Base() string {
	return fd.base
}

func (fd *FieldDescription) Kind() string {
	return fd.kind
}

func (fd *FieldDescription) Node() ast.Node {
	return fd.node
}

// StructTable is a collection that manages mappings of struct names to their associated properties.
type StructTable struct {
	container       map[string][]*FieldDescription
	gk              objects.IGateKeeper
	scopes          *Scopes
	implementations map[string][]string
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable(gk objects.IGateKeeper, scopes *Scopes) *StructTable {
	st := &StructTable{
		container:       make(map[string][]*FieldDescription),
		implementations: make(map[string][]string),
		gk:              gk,
		scopes:          scopes,
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

// Container returns the internal map associating struct names with slices of their field descriptions.
func (st *StructTable) Container() map[string][]*FieldDescription {
	return st.container
}

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
func (st *StructTable) Inference(expr ast.Expr, inferredTypeName string) (string, []string, bool) {
	if len(inferredTypeName) > 0 {
		return inferredTypeName, []string{inferredTypeName}, true
	}
	switch rhs := expr.(type) {
	case *ast.BinaryExpr:
		//nothing to do
		return "", nil, false
	case *ast.BasicLit:
		//nothing to do
		return "", nil, false
	case *ast.CompositeLit: // es. MyStruct{...}
		if baseName := st.ExtractBaseName(rhs.Type); len(baseName) > 0 {
			return baseName, []string{baseName}, true
		}
	case *ast.CallExpr: // es. NewStruct()
		if ident, ok := rhs.Fun.(*ast.Ident); ok {
			if funcSymbol, ok := st.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.ReturnTypes()) > 0 {
				// We assume the first return type
				returnType := funcSymbol.ReturnTypes()[0]
				// Verify if the returned type is a struct
				if typeSymbol, ok := st.scopes.SymbolResolve(returnType); ok && typeSymbol.IsStruct() {
					return returnType, []string{returnType}, true
				}
			}
		}
	case *ast.UnaryExpr: // es. &MyStruct{}
		if rhs.Op == token.AND {
			if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
				if ident, ok := compLit.Type.(*ast.Ident); ok {
					if returnSymbol, ok := st.scopes.SymbolResolve(ident.Name); ok && returnSymbol.IsStruct() {
						return returnSymbol.Name(), []string{returnSymbol.Name()}, true
					}
				}
			}
		}
	}
	return "", nil, false
}

// FieldsFromLiteral extracts and assigns struct fields from a given composite literal node, handling both keyed and positional formats.
func (st *StructTable) FieldsFromLiteral(node *ast.CompositeLit) ([]*FieldDescription, error) {
	// struct literal (es. MyStruct{...})
	t, ok := node.Type.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("unsupported composite literal type: %T", node)
	}
	structFields, ok := st.getFields(t.Name)
	if !ok {
		return nil, fmt.Errorf("unknown composite literal type: %st", t.Name)
	}
	if len(node.Elts) > len(structFields) {
		return nil, fmt.Errorf("too many values in positional struct literal for type '%st'", t.Name)
	}
	symbol, ok := st.scopes.SymbolResolve(t.Name)
	if !ok {
		var err error
		if symbol, err = st.scopes.SymbolDefine(t.Name); err != nil {
			return nil, err
		}
	}
	if err := st.AssignSymbol(symbol, t.Name, []string{t.Name}); err != nil {
		return nil, err
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
		for i, elt := range node.Elts {
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
	receiverStructFields, ok := st.getFields(receiverSymbol.StructName())
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
func (st *StructTable) AssignSymbol(symbol *Symbol, structName string, returnTypes []string) error {
	if len(structName) == 0 {
		return nil
	}
	structFields, ok := st.container[structName]
	if !ok {
		return nil
	}
	fields := make([]string, len(structFields))
	for x, field := range structFields {
		fields[x] = field.name
	}
	description := structName + "=>" + symbol.Name() + ":" + strings.Join(returnTypes, " ")
	symbol.SetStruct(structName, fields)
	symbol.SetReturnTypes(returnTypes)
	symbol.SetObject(st.gk.NewString(objects.FrameStatic, description))
	return nil
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
