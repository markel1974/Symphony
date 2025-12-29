package compiler

import (
	"go/ast"
	"go/token"
)

// Types represents a structure to manage import declarations and track imported modules in the build process.
type Types struct {
	declarations *Declarations
	container    []ast.Decl
	fileSet      *token.FileSet
	compile      func(node ast.Node) error
}

// NewTypes creates a new instance of Types with the provided Declarations.
func NewTypes(declarations *Declarations) *Types {
	return &Types{
		declarations: declarations,
		compile:      nil,
	}
}

// Setup initializes the Types instance with the provided compile function for processing AST nodes.
func (o *Types) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	o.fileSet = fileSet
	o.compile = compile
	return nil
}

// Declare adds the specified declaration to the container.
func (o *Types) Declare(decls ast.Decl) {
	o.container = append(o.container, decls)
}

// Prepare initializes the Types instance by setting up necessary state for processing import declarations.
func (o *Types) Prepare() error {
	return nil
}

// Compile processes all stored import declarations and validates them, returning an error if compilation fails.
func (o *Types) Compile() error {
	for _, decl := range o.container {
		if err := o.compile(decl); err != nil {
			return err
		}
	}
	return nil
}

// Finalize finalizes the Types structure by performing necessary cleanup or concluding operations. Returns an error if it fails.
func (c *Types) Finalize() error {
	return nil
}
