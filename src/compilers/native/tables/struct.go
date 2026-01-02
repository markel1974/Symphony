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
	pointer      bool
	container    string
	kind         string
	fieldName    string
	fieldNode    ast.Node
	totalSize    int
	finalized    bool
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

// SetOptions updates the pointer status, container type, and kind classification of the Struct instance.
func (sd *Struct) SetOptions(pointer bool, container string, kind string) {
	sd.pointer = pointer
	sd.container = container
	sd.kind = kind
}

// Options returns the pointer status, container type, and kind of the Struct instance as a tuple.
func (sd *Struct) Options() (bool, string, string) {
	return sd.pointer, sd.container, sd.kind
}

// FieldClone creates and returns a deep copy of the Struct, including its fields, maintaining the original structure and properties.
func (sd *Struct) FieldClone() IStructField {
	out := NewStruct(sd.name, sd.sType)
	out.fieldName = sd.fieldName
	out.fieldNode = sd.fieldNode
	out.pointer = sd.pointer
	out.container = sd.container
	out.kind = sd.kind
	out.totalSize = sd.totalSize
	out.finalized = sd.finalized
	for _, field := range sd.fields {
		out.AddField(field.FieldClone())
	}
	return out
}

// IsBuiltin indicates whether the Struct is of the built-in type StructTypeBuiltin.
func (sd *Struct) IsBuiltin() bool {
	return sd.sType == StructTypeBuiltin
}

// Offset returns the offset value for the Struct instance, typically used for memory or layout calculations.
func (sd *Struct) Offset() int {
	return 0
}

// SetOffset sets the offset value for the struct, but does not perform any operation for top-level structs.
func (sd *Struct) SetOffset(offset int) {
	// No-op for top-level
}

// Definition returns the Struct instance itself as an interface{} type.
func (sd *Struct) Definition() IStructField {
	return sd
}

// BindDefinition associates an external definition with the Struct instance. This method is currently a no-op.
func (sd *Struct) BindDefinition(_ IStructField) {
	// No-op
}

// IsPlaceholder determines if the Struct instance acts as a placeholder and always returns false.
func (sd *Struct) IsPlaceholder() bool {
	return false // Una struct definita non è mai un placeholder
}

// IsFinalized returns true if the Struct instance has completed initialization and its metadata is finalized.
func (sd *Struct) IsFinalized() bool {
	return sd.finalized
}

// SetFinalized sets the finalized status of the Struct to the specified boolean value.
func (sd *Struct) SetFinalized(finalized bool) {
	sd.finalized = finalized
}

// TotalSize returns the total computed size of the Struct.
func (sd *Struct) TotalSize() int {
	return sd.totalSize
}

// SetTotalSize updates the total size of the struct instance with the provided size value.
func (sd *Struct) SetTotalSize(size int) {
	sd.totalSize = size
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
