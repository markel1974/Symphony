package compiler

import (
	"go/ast"
	"go/token"
)

// Others represents a structure to manage import declarations and track imported modules in the build process.
type Others struct {
	declarations *Declarations
	container    []ast.Decl
	fileSet      *token.FileSet
	compile      func(node ast.Node) error
}

func NewOthers(declarations *Declarations) *Others {
	return &Others{
		declarations: declarations,
		compile:      nil,
	}
}

// Setup initializes the `Others` instance with a compile function used for processing AST nodes.
func (o *Others) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	o.fileSet = fileSet
	o.compile = compile
	return nil
}

// Declare adds the specified declaration to the container.
func (o *Others) Declare(decls ast.Decl) {
	o.container = append(o.container, decls)
}

// Prepare initializes the Others instance by setting up necessary state for processing import declarations.
func (o *Others) Prepare() error {
	return nil
}

// Compile processes all stored import declarations and validates them, returning an error if compilation fails.
func (o *Others) Compile() error {
	for _, decl := range o.container {
		if err := o.compile(decl); err != nil {
			return err
		}
	}
	return nil
}
