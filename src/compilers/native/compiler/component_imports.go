package compiler

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// Imports is a structure that manages a set of imported items, built-in functions, and related compilation resources.
type Imports struct {
	gk        objects.IGateKeeper
	loader    bytecode.ILoader
	imports   *tables.Constants
	scopes    *tables.Scopes
	modules   map[string]bool
	helper    map[string]int
	container []ast.Decl
	fileSet   *token.FileSet
	compile   func(node ast.Node) error
}

// NewImports creates and initializes a new Imports instance with provided GateKeeper, Constants, and Scopes references.
func NewImports(gk objects.IGateKeeper, loader bytecode.ILoader, imports *tables.Constants, scopes *tables.Scopes) *Imports {
	i := &Imports{
		gk:      gk,
		loader:  loader,
		imports: imports,
		scopes:  scopes,
		modules: make(map[string]bool),
		helper:  make(map[string]int),
		compile: nil,
	}
	return i
}

// Setup initializes the Imports instance by configuring compilation and loading built-in functions from the provided loader.
func (i *Imports) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	i.fileSet = fileSet
	i.compile = compile
	return nil
}

// Prepare initializes the Imports instance, ensuring it is ready for further use in the compilation process.
func (i *Imports) Prepare() error {
	return nil
}

// Compile processes and compiles all stored declarations. Returns an error if any declaration fails to compile.
func (i *Imports) Compile() error {
	for _, decl := range i.container {
		if err := i.compile(decl); err != nil {
			return err
		}
	}
	return nil
}

// Declare appends the provided declaration to the container slice of the Imports structure.
func (i *Imports) Declare(decls ast.Decl) {
	i.container = append(i.container, decls)
}

// Emit attempts to attach a function reference or emit a builtin reference, returning true if successful.
func (i *Imports) Emit(pos token.Pos, name string, selName string) bool {
	if len(name) == 0 {
		return false
	}
	var target string
	if len(selName) > 0 {
		_, ok := i.modules[name]
		if !ok {
			return false
		}
		target = tables.GetMangledName(name, selName)
	} else {
		target = name
	}
	var index int
	if v, ok := i.helper[target]; ok {
		index = v
	} else {
		index = i.imports.Add(target, i.gk.NewString(objects.FrameStatic, target))
		i.helper[target] = index
	}
	if _, err := i.scopes.Emit(pos, native.OpImportId, index); err != nil {
		return false
	}
	return true
}

// doGenDecl processes a generic declaration node, compiling each specification it contains. Returns an error if compilation fails.
func (i *Imports) doGenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := i.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// ImportSpec processes an ast.ImportSpec node, extracts the module name, and adds it to the imports map as a key.
func (i *Imports) ImportSpec(node *ast.ImportSpec) error {
	moduleName := node.Path.Value
	moduleName = strings.Trim(moduleName, "\"'")
	i.modules[moduleName] = true
	return nil
}
