package tables

import "go/ast"

// StructField represents a field within a structure with associated metadata.
type StructField struct {
	base string

	fieldName string
	fieldNode ast.Node
}

// NewStructField creates and returns a new instance of StructField with the specified name, base, kind, and AST node.
func NewStructField(base string) *StructField {
	return &StructField{
		base: base,
	}
}

// FieldBase returns the base type of the struct field as a string.
func (fd *StructField) FieldBase() string {
	return fd.base
}

// FieldName returns the name of the field as a string.
func (fd *StructField) FieldName() string {
	return fd.fieldName
}

// SetFieldName assigns a new name to the StructField instance.
func (fd *StructField) SetFieldName(fieldName string) {
	fd.fieldName = fieldName
}

// FieldNode retrieves the associated AST node of the StructField instance.
func (fd *StructField) FieldNode() ast.Node {
	return fd.fieldNode
}

// SetFieldNode assigns the given AST node to the StructField's node property.
func (fd *StructField) SetFieldNode(node ast.Node) {
	fd.fieldNode = node
}

// FieldClone creates a new instance of StructField with the same name, base, and kind but without a node value.
func (fd *StructField) FieldClone() IStructField {
	out := NewStructField(fd.base)
	out.fieldName = fd.fieldName
	return out
}
