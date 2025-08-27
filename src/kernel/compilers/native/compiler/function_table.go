package compiler

import (
	"fmt"
	"go/ast"
)

// FunctionDescription represents the structure to describe a Go function, its metadata, and related declarations.
type FunctionDescription struct {
	Name            string
	ReturnValues    []string
	InputParams     []string
	StructName      string
	StructReceivers []string
	FuncDecl        *ast.FuncDecl
}

// NewFunctionDescription creates a new FunctionDescription using the provided *ast.FuncDecl.
func NewFunctionDescription(funcDecl *ast.FuncDecl) *FunctionDescription {
	return &FunctionDescription{
		FuncDecl: funcDecl,
	}
}

// FunctionTable is a type designed to manage a collection of function descriptions.
type FunctionTable struct {
	container []*FunctionDescription
}

// NewFunctionTable initializes and returns a new instance of FunctionTable.
func NewFunctionTable() *FunctionTable {
	return &FunctionTable{}
}

// Add appends a new function description, created from the provided function declaration, to the container.
func (f *FunctionTable) Add(node *ast.FuncDecl) {
	fd := NewFunctionDescription(node)
	f.container = append(f.container, fd)
}

// Len returns the number of elements currently stored in the container of the FunctionTable.
func (f *FunctionTable) Len() int {
	return len(f.container)
}

// Get retrieves a FunctionDescription at the specified index. Returns an error if the index is out of bounds.
func (f *FunctionTable) Get(index int) (*FunctionDescription, error) {
	if index < 0 || index >= f.Len() {
		return nil, fmt.Errorf("function handler: index %d out of bounds", index)
	}
	return f.container[index], nil
}
