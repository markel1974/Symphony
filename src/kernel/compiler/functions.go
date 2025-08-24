package compiler

import (
	"go/ast"
)

// FunctionDescription represents the metadata of a function including its name, associated struct, parameters, and receiver info.
type FunctionDescription struct {
	Name       string
	Struct     string
	Types      []string
	Params     []string
	Recv       []string
	FuncDecl   *ast.FuncDecl
	StructType bool
}

// NewFunctionDescription creates a new instance of FunctionDescription with the provided function declaration.
func NewFunctionDescription(funcDecl *ast.FuncDecl) *FunctionDescription {
	return &FunctionDescription{
		FuncDecl: funcDecl,
	}
}

// Functions is a collection that manages a list of function descriptions.
type Functions struct {
	container []*FunctionDescription
}

// NewFunctions initializes and returns a new Functions instance.
func NewFunctions() *Functions {
	return &Functions{}
}

// Declare adds a function description derived from the provided function declaration to the Functions container.
func (f *Functions) Declare(funcDecl *ast.FuncDecl) {
	fd := NewFunctionDescription(funcDecl)
	f.container = append(f.container, fd)
}
