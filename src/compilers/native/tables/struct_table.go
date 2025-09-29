package tables

import (
	"fmt"
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

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
		sd := NewStruct(builtin, StructTypeBuiltin)
		//sd.Add(z)
		st.container[sd.name] = sd
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

// AddExternal adds a new package with the given name to the StructTable if it does not already exist.
func (st *StructTable) AddExternal(name string) {
	if _, ok := st.container[name]; !ok {
		sd := NewStruct(name, StructTypePackage)
		st.container[name] = sd
	}
}

// AddStruct retrieves or creates a new Struct with the given name in the StructTable and returns it.
func (st *StructTable) AddStruct(name string) *Struct {
	sd, ok := st.container[name]
	if !ok {
		sd = NewStruct(name, StructTypeDefined)
		st.container[name] = sd
	}
	return sd
}

// AddField adds a new field description to a struct in the StructTable. If the struct does not exist, it creates it.
func (st *StructTable) AddField(name string, fieldName string, baseStruct string, kind string, node ast.Node) {
	sd := st.AddStruct(name)
	sd.AddField(fieldName, baseStruct, kind, node)
}

// Has checks if a struct definition with the given name exists in the container map.
func (st *StructTable) Has(name string) bool {
	if _, ok := st.container[name]; ok {
		return true
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
