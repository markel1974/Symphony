package compiler

import (
	"fmt"
	"go/ast"
)

// Others represents a structure to manage import declarations and track imported modules in the build process.
type Others struct {
	declarations *Declarations
	container    []ast.Decl
}

func NewOthers(declarations *Declarations) *Others {
	return &Others{
		declarations: declarations,
	}
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

// compile traverses the provided AST node and compiles it into bytecode, handling various node types in a switch block.
func (o *Others) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.GenDecl:
		err = o.declarations.GenDecl(node)
	default:
		err = fmt.Errorf("[compiler] unsupported expression type: %T", node)
	}
	return err
}
