package tables

import "go/ast"

// StructField represents a single field of a struct, capturing its name, base type, kind, and corresponding AST node.
type StructField struct {
	name string
	base string
	kind string
	node ast.Node
}

// NewStructField creates and returns a new instance of StructField with the specified name, base, kind, and node.
func NewStructField(name string, base string, kind string, node ast.Node) *StructField {
	return &StructField{
		name: name,
		base: base,
		kind: kind,
		node: node,
	}
}

// Name returns the name of the struct field.
func (fd *StructField) Name() string {
	return fd.name
}

// Base retrieves the base representation of the StructField.
func (fd *StructField) Base() string {
	return fd.base
}

// Kind returns the string representing the kind of the StructField.
func (fd *StructField) Kind() string {
	return fd.kind
}

// Node returns the associated ast.Node of the StructField instance.
func (fd *StructField) Node() ast.Node {
	return fd.node
}
