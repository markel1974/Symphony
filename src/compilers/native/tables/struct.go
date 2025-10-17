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
	name         string
	sType        StructType
	fields       []IStructField
	fieldsHelper map[string]IStructField

	fieldName string
	fieldNode ast.Node
}

// NewStruct creates and returns a pointer to a StructData instance, initializing its fields and type.
func NewStruct(name string, sType StructType) *Struct {
	return &Struct{
		sType:        sType,
		name:         name,
		fields:       []IStructField{},
		fieldsHelper: make(map[string]IStructField),
		fieldNode:    nil,
	}
}

// AddField appends a new field with the specified name, base type, kind, and associated AST node to the Struct.
func (sd *Struct) AddField(sf IStructField) {
	sd.fields = append(sd.fields, sf)
	sd.fieldsHelper[sf.FieldName()] = sf
}

// Fields returns a slice of StructField pointers representing the fields of the struct.
func (sd *Struct) Fields() []IStructField {
	out := make([]IStructField, len(sd.fields))
	for idx, v := range sd.fields {
		out[idx] = v.FieldClone()
	}
	return out
}

// FieldsName returns a slice of field names defined in the StructData instance.
func (sd *Struct) FieldsName() []string {
	names := make([]string, len(sd.fields))
	for x, field := range sd.fields {
		names[x] = field.FieldName()
	}
	return names
}

// Name returns the name of the Struct instance as a string.
func (sd *Struct) Name() string {
	return sd.name
}

// Type returns the StructType representing the kind or classification of the Struct instance.
func (sd *Struct) Type() StructType {
	return sd.sType
}

func (sd *Struct) FieldBase() string {
	return sd.name
}

func (sd *Struct) FieldName() string {
	return sd.fieldName
}

// SetFieldName assigns the provided AST node to the Struct instance.
func (sd *Struct) SetFieldName(name string) {
	sd.fieldName = name
}

// SetFieldNode assigns the provided AST node to the Struct instance.
func (sd *Struct) SetFieldNode(node ast.Node) {
	sd.fieldNode = node
}

// FieldNode returns the associated AST node of the Struct instance.
func (sd *Struct) FieldNode() ast.Node {
	return sd.fieldNode
}

// FieldClone creates and returns a deep copy of the Struct, including its fields, maintaining the original structure and properties.
func (sd *Struct) FieldClone() IStructField {
	out := NewStruct(sd.name, sd.sType)
	for _, field := range sd.fields {
		out.AddField(field.FieldClone())
	}
	return out
}

// IsBuiltin determines if the Struct instance represents a built-in type by comparing its kind to StructTypeBuiltin.
func (sd *Struct) IsBuiltin() bool {
	return sd.sType == StructTypeBuiltin
}

/*
func (sd *Struct) Walk(segments []string) (IWalker, bool) {
	if len(segments) == 0 {
		return nil, false
	}
	sf, ok := sd.fieldsHelper[segments[0]]
	if !ok {
		return nil, false
	}
	segments = segments[1:]
	if len(segments) == 0 {
		return sf.st, true
	}
	return sf.st.Walk(segments)
}
*/
