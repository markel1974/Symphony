package compiler

import (
	"fmt"
	"go/ast"
)

// Types represents a structure to manage import declarations and track imported modules in the build process.
type Types struct {
	declarations *Declarations
	container    []ast.Decl
}

func NewTypes(declarations *Declarations) *Types {
	return &Types{
		declarations: declarations,
	}
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

// compile traverses the provided AST node and compiles it into bytecode, handling various node types in a switch block.
func (o *Types) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.GenDecl:
		err = o.declarations.GenDecl(node)
	default:
		err = fmt.Errorf("[compiler] unsupported expression type: %T", node)
	}
	return err
}
