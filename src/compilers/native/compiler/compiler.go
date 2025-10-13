package compiler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// Compiler represents a structure to manage the compilation process, including scopes and associated token file sets.
type Compiler struct {
	gk                objects.IGateKeeper
	components        []IComponent
	fileSet           *token.FileSet
	loader            bytecode.ILoader
	scopes            *tables.Scopes
	constants         *tables.Constants
	importConstants   *tables.Constants
	imports           *Imports
	functions         *Functions
	types             *Types
	others            *Others
	expressions       *Expression
	declarations      *Declarations
	controlFlow       *ControlFlow
	loops             *Loops
	definitionTables  *tables.DefinitionTable
	functionTable     *tables.FunctionTable
	typeCompatibility *TypeCompatibility
	rootNode          *ast.File
}

// New creates and returns a new instance of Compiler with initialized scopes using a standard library loader.
func New(gk objects.IGateKeeper, loader bytecode.ILoader, opcodes opcodes.IOpcodes) *Compiler {
	var components []IComponent
	constants := tables.NewConstants()
	scopes := tables.NewScopes(gk, opcodes, constants)
	definitionTable := tables.NewDefinitionTable(gk, scopes)
	functionTable := tables.NewFunctionTable(gk, scopes, definitionTable)
	importConstants := tables.NewConstants()
	imports := NewImports(gk, loader, importConstants, constants, scopes)
	components = append(components, imports)
	declarations := NewDeclarations(gk, importConstants, constants, scopes, imports, definitionTable, functionTable)
	components = append(components, declarations)
	expressions := NewExpression(gk, constants, scopes, imports)
	components = append(components, expressions)
	functions := NewFunctions(gk, constants, scopes, imports, declarations, definitionTable, functionTable)
	components = append(components, functions)
	controlFlow := NewControlFlow(gk, constants, scopes, definitionTable)
	components = append(components, controlFlow)
	loops := NewLoops(gk, scopes, definitionTable, functionTable)
	components = append(components, loops)
	types := NewTypes(declarations)
	components = append(components, types)
	others := NewOthers(declarations)
	components = append(components, others)
	typeCompatibility := NewTypeCompatibility(definitionTable, functionTable)
	components = append(components, typeCompatibility)

	c := &Compiler{
		gk:                gk,
		components:        components,
		loader:            loader,
		scopes:            scopes,
		constants:         constants,
		importConstants:   importConstants,
		definitionTables:  definitionTable,
		functionTable:     functionTable,
		imports:           imports,
		functions:         functions,
		declarations:      declarations,
		controlFlow:       controlFlow,
		expressions:       expressions,
		loops:             loops,
		types:             types,
		others:            others,
		typeCompatibility: typeCompatibility,
		rootNode:          nil,
	}
	return c
}

// FileSet returns the token file set associated with the compiler.
func (c *Compiler) FileSet() *bytecode.FileSet {
	return bytecode.NewFileSet(c.fileSet)
}

// Compile parses the provided source file and compiles it into bytecode. Returns compiled bytecode or an error.
func (c *Compiler) Compile(filename string, source any) error {
	c.fileSet = token.NewFileSet()
	for _, component := range c.components {
		if err := component.Setup(c.fileSet, c.compile); err != nil {
			return err
		}
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

// Imports retrieves a list of IObject references from the current compiler scope.
func (c *Compiler) Imports() []objects.IObject {
	return c.importConstants.Retrieve()
}

// Globals retrieves and returns all global objects from the root scope and any objects tracked by references.
func (c *Compiler) Globals() []objects.IObject {
	ret := c.scopes.SymbolGlobalsCreate()
	return ret
}

// Print writes the content of the internal scopes to the provided writer, typically for debugging or inspection.
func (c *Compiler) Print(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "----- Constants -----")
	c.constants.Print(writer)
	_, _ = fmt.Fprintf(writer, "----- Imports -----")
	c.importConstants.Print(writer)
	c.scopes.Print(writer)
}

func (c *Compiler) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.File:
		err = c.File(node)
	case *ast.ParenExpr:
		return c.expressions.ParenExpr(node)
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
	case *ast.MapType:
		err = c.declarations.MapType(node)
	case *ast.ArrayType:
		err = c.declarations.ArrayType(node)
	case *ast.BlockStmt:
		err = c.functions.BlockStmt(node)
	case *ast.ExprStmt:
		err = c.functions.ExprStmt(node)
	case *ast.RangeStmt:
		err = c.loops.RangeStmt(node)
	case *ast.ForStmt:
		err = c.loops.ForStmt(node)
	case *ast.IncDecStmt:
		err = c.expressions.IncDecStmt(node)
	case *ast.BinaryExpr:
		err = c.expressions.BinaryExpr(node)
	case *ast.UnaryExpr:
		err = c.expressions.UnaryExpr(node)
	case *ast.SelectorExpr:
		err = c.expressions.SelectorExpr(node)
	case *ast.SliceExpr:
		return c.expressions.SliceExpr(node)
	case *ast.FuncLit:
		err = c.functions.FuncLit(node)
	case *ast.FuncDecl:
		err = c.functions.FuncDecl(node)
	case *ast.CallExpr:
		err = c.functions.CallExpr(node)
	case *ast.ReturnStmt:
		err = c.functions.ReturnStmt(node)
	case *ast.DeferStmt:
		err = c.functions.DeferStmt(node)
	case *ast.GoStmt:
		err = c.functions.GoStmt(node)
	case *ast.IfStmt:
		err = c.controlFlow.IfStmt(node)
	case *ast.BranchStmt:
		err = c.controlFlow.BranchStmt(node)
	case *ast.SwitchStmt:
		err = c.controlFlow.SwitchStmt(node)
	case *ast.TypeSwitchStmt:
		err = c.controlFlow.TypeSwitchStmt(node)
	default:
		err = tables.NewCompilerError(c.fileSet, node, "unsupported expression type: %T", node)
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
		c.typeCompatibility.Prepare,
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
	c.scopes.SetRootIndex()
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	if scope.InstructionsLen() == 0 {
		return nil
	}
	if _, err = c.scopes.SymbolEmit(token.Pos(0), native.OpReturnId, 0); err != nil {
		return err
	}
	initFuncCode := scope.Instructions()
	initSource := scope.Source()
	numLocals := c.scopes.SymbolCount()
	compiledInitFn := c.gk.NewFunc(objects.FrameStatic, bytecode.PreInitFunction, initFuncCode, numLocals, 0, false, initSource)
	_, err = c.scopes.SymbolDefineConst(bytecode.PreInitFunction, compiledInitFn)
	if err != nil {
		return err
	}
	//	initSymbols.SetObject(compiledInitFn)
	return nil
}
