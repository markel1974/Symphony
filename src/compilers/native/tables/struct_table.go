package tables

import (
	"fmt"
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

type IWalker interface {
	Name() string
	Walk(segments []string) (IWalker, bool)
}

// StructTable is a collection that manages mappings of struct names to their associated properties.
type StructTable struct {
	container       map[string]*Struct
	gk              objects.IGateKeeper
	scopes          *Scopes
	implementations map[string][]string
}

// NewStructTable initializes and returns a pointer to a StructTable instance with an empty container map.
func NewStructTable(gk objects.IGateKeeper, scopes *Scopes) *StructTable {
	st := &StructTable{
		container:       make(map[string]*Struct),
		implementations: make(map[string][]string),
		gk:              gk,
		scopes:          scopes,
	}
	builtins := []string{"error"}
	for _, builtin := range builtins {
		//z := NewStructField(internal, "", internal, nil)
		st.AddStructBuiltin(builtin)
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

// SetImplementations sets the implementation mappings for structs to interfaces in the StructTable.
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

// AddStructBuiltin retrieves or creates a new Struct with the given name in the StructTable and returns it.
func (st *StructTable) AddStructBuiltin(name string) *Struct {
	return st.tryAddStruct(name, StructTypeBuiltin)
}

// AddExternal adds a new package with the given name to the StructTable if it does not already exist.
func (st *StructTable) AddExternal(name string) *Struct {
	return st.tryAddStruct(name, StructTypePackage)
}

// AddStruct retrieves or creates a new Struct with the given name in the StructTable and returns it.
func (st *StructTable) AddStruct(name string) *Struct {
	return st.tryAddStruct(name, StructTypeDefined)
}

// tryAddStruct retrieves an existing Struct or creates a new one with the given name and type classification.
func (st *StructTable) tryAddStruct(name string, kind StructType) *Struct {
	sd, ok := st.container[name]
	if !ok {
		sd = NewStruct(name, kind)
		st.container[name] = sd
	}
	return sd
}

// Has checks if a struct definition with the given name exists in the container map.
func (st *StructTable) Has(name string) bool {
	if _, ok := st.container[name]; ok {
		return true
	}
	return false
}

// Walk traverses through the nested structure of the specified name using the given path, returning true if valid.
func (st *StructTable) Walk(name string, path []string) bool {
	root, ok := st.container[name]
	if !ok {
		return false
	}
	if len(path) == 0 {
		return true
	}
	v, ok := root.fieldsHelper[path[0]]
	if !ok {
		return false
	}
	walker := v.st
	if walker == nil {
		return false
	}
	if _, ok = walker.Walk(path[1:]); ok {
		return ok
	}
	return false
}

// FieldsFromLiteral extracts and assigns struct fields from a given composite literal node, handling both keyed and positional formats.
func (st *StructTable) FieldsFromLiteral(structName string, eltS []ast.Expr) ([]*StructField, error) {
	sd, ok := st.container[structName]
	if !ok {
		return nil, fmt.Errorf("unknown composite literal type: %st", structName)
	}
	structFields := sd.Fields()
	if len(eltS) > len(structFields) {
		return nil, fmt.Errorf("too many values in positional struct literal for type '%st'", structName)
	}
	isKeyed := false
	if len(eltS) > 0 {
		if _, ok := eltS[0].(*ast.KeyValueExpr); ok {
			isKeyed = true
		}
	}
	if isKeyed {
		// key literal (es. Home{Name: "Alfa", Address: "Shanghai"})
		providedFields := make(map[string]ast.Expr)
		for _, elt := range eltS {
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
		for i, elt := range eltS {
			structFields[i].node = elt
		}
	}
	return structFields, nil
}

// TypeNameFromSymbolField retrieves the base type of a field within a struct using its name and returns it with a success flag.
func (st *StructTable) TypeNameFromSymbolField(name string, fieldName string) (string, bool) {
	receiverSymbol, ok := st.scopes.SymbolResolve(name)
	if !ok {
		return "", false
	}
	sd, ok := st.container[receiverSymbol.StructName()]
	if !ok {
		return "", false
	}
	for _, receiverField := range sd.Fields() {
		if receiverField.name == fieldName {
			return receiverField.base, true
		}
	}
	return "", false
}

// BindSymbol assigns a struct name and types to a Symbol, validates the struct, and creates a description object.
func (st *StructTable) BindSymbol(symbol *Symbol, typeName string) {
	if sd, ok := st.container[typeName]; ok {
		fields := sd.FieldsName()
		symbol.SetStruct(typeName, fields)
		return
	}
}

// IsBuiltin returns true if the given name is a struct internal to the compiler.
func (st *StructTable) IsBuiltin(name string) bool {
	fd, ok := st.container[name]
	if !ok {
		return false
	}
	return fd.IsBuiltin()
}

// Walker retrieves an IWalker associated with the given name from the StructTable container and returns it with a success flag.
func (st *StructTable) Walker(name string) (IWalker, bool) {
	w, ok := st.container[name]
	return w, ok
}
