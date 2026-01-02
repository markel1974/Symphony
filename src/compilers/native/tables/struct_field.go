package tables

import "go/ast"

// StructField represents a field within a struct, containing metadata and other attributes about the field.
type StructField struct {
	base      string
	fieldName string
	fieldNode ast.Node

	pointer   bool
	container string
	kind      string

	offset     int         // Byte offset from the beginning of the struct
	definition interface{} // Reference to the actual definition (Struct or Interface)
}

// NewStructField creates and returns a pointer to a new StructField with the specified base string.
func NewStructField(base string) *StructField {
	return &StructField{
		base: base,
	}
}

// FieldBase returns the base string value of the StructField.
func (fd *StructField) FieldBase() string {
	return fd.base
}

// FieldName returns the name of the field associated with the StructField instance.
func (fd *StructField) FieldName() string {
	return fd.fieldName
}

// SetFieldName sets the name of the field for the StructField instance.
func (fd *StructField) SetFieldName(fieldName string) {
	fd.fieldName = fieldName
}

// FieldNode returns the AST node associated with the StructField.
func (fd *StructField) FieldNode() ast.Node {
	return fd.fieldNode
}

// SetFieldNode assigns the specified AST node to the fieldNode property of the StructField.
func (fd *StructField) SetFieldNode(node ast.Node) {
	fd.fieldNode = node
}

// SetOptions updates the pointer, container, and kind attributes of the StructField.
func (fd *StructField) SetOptions(pointer bool, container string, kind string) {
	fd.pointer = pointer
	fd.container = container
	fd.kind = kind
}

// Options retrieves the pointer state, container type, and kind type of the StructField.
func (fd *StructField) Options() (bool, string, string) {
	return fd.pointer, fd.container, fd.kind
}

// FieldClone creates a deep copy of the current StructField, preserving all its attributes and returning it as IStructField.
func (fd *StructField) FieldClone() IStructField {
	out := NewStructField(fd.base)
	out.fieldName = fd.fieldName
	out.pointer = fd.pointer
	out.container = fd.container
	out.kind = fd.kind
	out.offset = fd.offset
	out.definition = fd.definition
	return out
}

// SetOffset sets the byte offset of the field from the beginning of the struct.
func (fd *StructField) SetOffset(offset int) {
	fd.offset = offset
}

// Offset returns the byte offset of the field from the beginning of the struct.
func (fd *StructField) Offset() int {
	return fd.offset
}

// BindDefinition assigns a definition, such as a struct or interface, to the field for reference or further processing.
func (fd *StructField) BindDefinition(def interface{}) {
	fd.definition = def
}

// Definition returns the definition associated with the struct field. It may reference a struct or interface definition.
func (fd *StructField) Definition() interface{} {
	return fd.definition
}

// IsPlaceholder checks if the StructField has no associated definition, indicating it is a placeholder.
func (fd *StructField) IsPlaceholder() bool {
	return fd.definition == nil
}

// IsFinalized checks if the field's definition is finalized. Returns true if finalized, otherwise false.
func (fd *StructField) IsFinalized() bool {
	if fd.definition != nil {
		// Optional: check if the definition is finalized
		if def, ok := fd.definition.(IStructField); ok {
			return def.IsFinalized()
		}
	}
	return false
}

// SetFinalized sets the finalized state for the field. This is typically managed at the struct level, not field level.
func (fd *StructField) SetFinalized(finalized bool) {
	// No-op for a simple field, or you might want to propagate it.
	// Generally, it is the Struct that maintains the "Finalized" state.
}
