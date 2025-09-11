package tables

import "go/ast"

type StructType int

const (
	StructTypeBuiltin StructType = iota
	StructTypeDefined
	StructTypePackage
)

// Struct represents a structured data type containing a name, type classification, and a collection of field descriptions.
type Struct struct {
	name   string
	kind   StructType
	fields []*StructField
}

// NewStruct creates and returns a pointer to a StructData instance, initializing its fields and type.
func NewStruct(name string, kind StructType) *Struct {
	return &Struct{
		kind:   kind,
		name:   name,
		fields: []*StructField{},
	}
}

// AddField appends a new field with the specified name, base type, kind, and associated AST node to the Struct.
func (sd *Struct) AddField(name string, base string, kind string, node ast.Node) {
	sf := NewStructField(name, base, kind, node)
	sd.fields = append(sd.fields, sf)
}

// Fields returns a slice of StructField pointers representing the fields of the struct.
func (sd *Struct) Fields() []*StructField {
	out := make([]*StructField, len(sd.fields))
	for idx, v := range sd.fields {
		out[idx] = NewStructField(v.name, v.base, v.kind, nil)
	}
	return out
}

// FieldsName returns a slice of field names defined in the StructData instance.
func (sd *Struct) FieldsName() []string {
	names := make([]string, len(sd.fields))
	for x, field := range sd.fields {
		names[x] = field.name
	}
	return names
}

// Name returns the name of the Struct instance as a string.
func (sd *Struct) Name() string {
	return sd.name
}

// Kind returns the StructType representing the kind or classification of the Struct instance.
func (sd *Struct) Kind() StructType {
	return sd.kind
}

// IsBuiltin determines if the Struct instance represents a built-in type by comparing its kind to StructTypeBuiltin.
func (sd *Struct) IsBuiltin() bool {
	return sd.kind == StructTypeBuiltin
}
