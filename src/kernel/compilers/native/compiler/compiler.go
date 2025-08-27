package compiler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/markel1974/c64emu/src/kernel/compilers/native/common"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// maxScope defines the maximum allowable depth for compilation scopes to prevent excessive recursion or memory use.
const (
	maxScope = 1024
)

// Compiler represents a structure to manage the compilation process, including scopes and associated token file sets.
type Compiler struct {
	gk           objects.IGateKeeper
	fileSet      *token.FileSet
	loader       bytecode.ILoader
	scopes       *Scopes
	constants    *Constants
	references   *Constants
	imports      *Imports
	functions    *Functions
	types        *Types
	others       *Others
	declarations *Declarations
	structs      *StructTable
	rootNode     *ast.File
}

// New creates and returns a new instance of Compiler with initialized scopes using a standard library loader.
func New(gk objects.IGateKeeper, loader bytecode.ILoader, opcodes *bytecode.Opcodes) *Compiler {
	scopes := NewScopes(gk, opcodes)
	constants := NewConstants()
	references := NewConstants()
	structs := NewStructTable(gk, scopes)
	declarations := NewDeclarations(gk, references, constants, scopes, structs)
	imports := NewImports(gk, references, scopes)
	functions := NewFunctions(gk, constants, scopes, imports, declarations, structs)
	types := NewTypes(declarations)
	others := NewOthers(declarations)
	c := &Compiler{
		gk:           gk,
		loader:       loader,
		scopes:       scopes,
		constants:    constants,
		references:   references,
		imports:      imports,
		functions:    functions,
		declarations: declarations,
		types:        types,
		others:       others,
		structs:      structs,
		rootNode:     nil,
	}
	return c
}

// Id returns the unique identifier of the compiler as defined in the common package.
func (c *Compiler) Id() string {
	return common.Identifier
}

// Compile parses the provided source file and compiles it into bytecode. Returns compiled bytecode or an error.
func (c *Compiler) Compile(filename string, source any) error {
	c.fileSet = token.NewFileSet()
	if err := c.imports.Setup(c.loader, c.compile); err != nil {
		return err
	}
	if err := c.declarations.Setup(c.fileSet, c.compile); err != nil {
		return err
	}
	if err := c.functions.Setup(c.fileSet, c.compile); err != nil {
		return err
	}
	if err := c.others.Setup(c.fileSet, c.compile); err != nil {
		return err
	}
	if err := c.types.Setup(c.fileSet, c.compile); err != nil {
		return err
	}
	astFile, err := parser.ParseFile(c.fileSet, filename, source, 0)
	if err != nil {
		return err
	}
	if err = c.compile(astFile); err != nil {
		return err
	}
	return nil
}

// Constants retrieves a slice of IObject containing all constants stored in the current compiler scopes.
func (c *Compiler) Constants() []objects.IObject {
	return c.constants.Retrieve()
}

// References retrieves a list of IObject references from the current compiler scope.
func (c *Compiler) References() []objects.IObject {
	return c.references.Retrieve()
}

// Globals retrieves and returns all global objects from the root scope and any objects tracked by references.
func (c *Compiler) Globals() []objects.IObject {
	ret := make([]objects.IObject, len(c.scopes.initSymbolTable.symbols))
	for _, obj := range c.scopes.initSymbolTable.definitions {
		target := obj.GetObject()
		if target != nil {
			ret[obj.index] = target
		} else {
			ret[obj.index] = c.gk.NewString(objects.FrameStatic, obj.Name()+"_placeHolder")
			//ret[obj.index] = c.factory.UndefinedValue()
		}
	}
	return ret
}

// Print writes the content of the internal scopes to the provided writer, typically for debugging or inspection.
func (c *Compiler) Print(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "----- Constants -----")
	c.constants.Print(writer)
	_, _ = fmt.Fprintf(writer, "----- References -----")
	c.references.Print(writer)
	c.scopes.Print(writer)
}

func (c *Compiler) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.File:
		err = c.File(node)
	case *ast.ImportSpec:
		err = c.imports.ImportSpec(node)
	case *ast.DeclStmt:
		err = c.declarations.DeclStmt(node)
	case *ast.ValueSpec:
		err = c.declarations.ValueSpec(node)
	case *ast.GenDecl:
		err = c.declarations.GenDecl(node)
	case *ast.TypeSpec:
		err = c.declarations.TypeSpec(node)
	case *ast.Ident:
		err = c.declarations.Ident(node)
	case *ast.AssignStmt:
		err = c.declarations.AssignStmt(node)
	case *ast.BasicLit:
		err = c.declarations.BasicLit(node)
	case *ast.CompositeLit:
		err = c.declarations.CompositeLit(node)
	case *ast.KeyValueExpr:
		err = c.declarations.KeyValueExpr(node)
	case *ast.StarExpr:
		err = c.declarations.StarExpr(node)
	case *ast.IndexExpr:
		err = c.declarations.IndexExpr(node)
	case *ast.BlockStmt:
		err = c.functions.BlockStmt(node)
	case *ast.ExprStmt:
		err = c.functions.ExprStmt(node)
	case *ast.IfStmt:
		err = c.functions.IfStmt(node)
	case *ast.RangeStmt:
		err = c.functions.RangeStmt(node)
	case *ast.ForStmt:
		err = c.functions.ForStmt(node)
	case *ast.IncDecStmt:
		err = c.functions.IncDecStmt(node)
	case *ast.BinaryExpr:
		err = c.functions.BinaryExpr(node)
	case *ast.UnaryExpr:
		err = c.functions.UnaryExpr(node)
	case *ast.FuncLit:
		err = c.functions.FuncLit(node)
	case *ast.FuncDecl:
		err = c.functions.FuncDecl(node)
	case *ast.CallExpr:
		err = c.functions.CallExpr(node)
	case *ast.ReturnStmt:
		err = c.functions.ReturnStmt(node)
	case *ast.SelectorExpr:
		err = c.functions.SelectorExpr(node)
	default:
		err = NewCompilerError(c.fileSet, node, "unsupported expression type: %T", node)
	}
	return err
}

// defaultPipeline returns a default compilation pipeline for the compiler.
func (c *Compiler) defaultPipeline() []func() error {
	pipeline := []func() error{
		c.collectDeclarations,
		c.imports.Prepare,
		c.imports.Compile,
		c.types.Prepare,
		c.types.Compile,
		c.functions.Prepare,
		c.others.Prepare,
		c.others.Compile,
		c.functions.Compile,
		c.createInit,
	}
	return pipeline
}

// File processes the given AST file node, categorizes declarations, and compiles them in a defined order while handling errors.
func (c *Compiler) File(node *ast.File) error {
	c.rootNode = node
	pipeline := c.defaultPipeline()
	for _, fn := range pipeline {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// collectDeclarations separates declarations by category and stores them in the appropriate category.
func (c *Compiler) collectDeclarations() error {
	if c.rootNode == nil {
		return fmt.Errorf("nil file node")
	}
	// step 1: Separate declarations by category
	for _, decl := range c.rootNode.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			c.functions.Declare(d)
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				c.imports.Declare(d)
			} else if d.Tok == token.TYPE {
				c.types.Declare(d)
			} else {
				c.others.Declare(d)
			}
		default:
			c.others.Declare(d)
		}
	}
	return nil
}

// createInit creates a default __init__ function for the root scope.
func (c *Compiler) createInit() error {
	c.scopes.scopeIndex = 0
	if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
		return err
	}
	initFuncCode := c.scopes.compilations[0].Instructions()
	numLocals := c.scopes.SymbolCount()
	initSymbols, err := c.scopes.SymbolDefine("__init__")
	if err != nil {
		return err
	}
	compiledInitFn := c.gk.NewFuncCompiled(objects.FrameStatic, initSymbols.Name(), initFuncCode, numLocals, 0, false, nil, nil)
	initSymbols.SetObject(compiledInitFn)
	return nil
}
