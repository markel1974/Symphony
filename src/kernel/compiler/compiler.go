package compiler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/markel1974/c64emu/src/kernel/compiler/sdk"
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
	loader       *sdk.Loader
	scopes       *Scopes
	constants    *Constants
	references   *Constants
	imports      *Imports
	functions    *Functions
	types        *Types
	others       *Others
	declarations *Declarations
	rootNode     *ast.File
}

// New creates and returns a new instance of Compiler with initialized scopes using a standard library loader.
func New(gk objects.IGateKeeper) *Compiler {
	loader := sdk.NewLoader(gk)
	op := bytecode.NewOpcodes(gk)
	scopes := NewScopes(gk, op)
	constants := NewConstants()
	references := NewConstants()
	declarations := NewDeclarations(gk, references, constants, scopes)
	imports := NewImports(gk, references, scopes)
	functions := NewFunctions(gk, constants, scopes, imports, declarations)
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
		rootNode:     nil,
	}
	return c
}

// Compile parses the provided source file and compiles it into bytecode. Returns compiled bytecode or an error.
func (c *Compiler) Compile(filename string, source any) error {
	for idx := 0; idx < c.loader.BuiltinLen(); idx++ {
		bi := c.loader.Builtin(idx)
		if bi == nil {
			return fmt.Errorf("builtin %d not found", idx)
		}
		c.constants.Add(bi.Name(), bi)
		//c.scopes.SymbolDefine(bi.Name(), GlobalScope, false)
	}

	fileSet := token.NewFileSet()
	c.declarations.Initialize(fileSet)
	astFile, err := parser.ParseFile(fileSet, filename, source, 0)
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

// Global retrieves and returns all global objects from the root scope and any objects tracked by references.
func (c *Compiler) Global() []objects.IObject {
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

// compile traverses the provided AST node and compiles it into bytecode, handling various node types in a switch block.
func (c *Compiler) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.File:
		err = c.doFile(node)
	default:
		err = fmt.Errorf("[compiler] unsupported expression type: %T", node)
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

// doFile processes the given AST file node, categorizes declarations, and compiles them in a defined order while handling errors.
func (c *Compiler) doFile(node *ast.File) error {
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
	initSymbols, err := c.scopes.SymbolDefine("__init__", UnknownScope, false)
	if err != nil {
		return err
	}
	compiledInitFn := c.gk.NewFuncCompiled(objects.FrameStatic, initSymbols.Name(), initFuncCode, numLocals, 0, false, nil, nil)
	initSymbols.SetObject(compiledInitFn)
	return nil
}
