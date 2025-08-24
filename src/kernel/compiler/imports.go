package compiler

import (
	"fmt"
	"go/ast"
	"strings"
)

// Imports represents a structure to manage import declarations and track imported modules in the build process.
type Imports struct {
	scopes    *Scopes
	imports   map[string]bool
	container []ast.Decl
}

// NewImports creates and returns a new instance of Imports with an initialized map to manage import declarations.
func NewImports(scopes *Scopes) *Imports {
	return &Imports{
		scopes:  scopes,
		imports: make(map[string]bool),
	}
}

// Contains checks if the specified import name exists in the Imports map and returns true if it is present, otherwise false.
func (i *Imports) Contains(name string) bool {
	return i.imports[name]
}

// Declare adds the specified declaration to the container.
func (i *Imports) Declare(decls ast.Decl) {
	i.container = append(i.container, decls)
}

// Compile processes all stored import declarations and validates them, returning an error if compilation fails.
func (i *Imports) Compile() error {
	for _, decl := range i.container {
		if err := i.compile(decl); err != nil {
			return err
		}
	}
	return nil
}

// compile processes the given AST node and applies import-specific handling or returns an error for unsupported types.
func (i *Imports) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.ImportSpec:
		err = i.doImportSpec(node)
	case *ast.GenDecl:
		err = i.doGenDecl(node)
	default:
		err = fmt.Errorf("unsupported expression type: %T", node)
	}
	return err
}

// doGenDecl processes a general declaration node by compiling each specification within the node. It returns an error if any occur.
func (i *Imports) doGenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := i.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// doImportSpec processes an ast.ImportSpec node and adds the module name to the imports map after sanitization.
func (i *Imports) doImportSpec(node *ast.ImportSpec) error {
	moduleName := node.Path.Value
	moduleName = strings.Trim(moduleName, "\"'")
	i.imports[moduleName] = true
	return nil
}
