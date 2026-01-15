package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/symphony/src/compilers/native/tables"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/sequencers/native"
)

// Expression represents a structure for handling and compiling abstract syntax trees using related tables and utilities.
type Expression struct {
	gk          objects.IGateKeeper
	fileSet     *token.FileSet
	constants   *tables.Constants
	scopes      *tables.Scopes
	imports     *Imports
	definitions *tables.DefinitionTable
	compile     func(node ast.Node) error
}

// NewExpression creates a new Expression instance using the provided gatekeeper, constants, scopes, imports, and definitions.
func NewExpression(gk objects.IGateKeeper, constants *tables.Constants, scopes *tables.Scopes, imports *Imports, definitions *tables.DefinitionTable) *Expression {
	return &Expression{
		gk:          gk,
		constants:   constants,
		imports:     imports,
		scopes:      scopes,
		definitions: definitions,
	}
}

// Setup initializes the Expression instance with a file set and a compile function to process AST nodes.
func (c *Expression) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Prepare initializes or validates the state of the Expression before compilation.
func (c *Expression) Prepare() error {
	return nil
}

// Compile performs the main compilation process for the expression, executing the provided compile function on AST nodes.
func (c *Expression) Compile() error {
	return nil
}

// Finalize completes final processing of the Expression and performs any necessary cleanup operations.
func (c *Expression) Finalize() error {
	return nil
}

// UnaryExpr processes a unary expression node, emits appropriate bytecode, and handles errors for unhandled operators or scopes.
func (c *Expression) UnaryExpr(node *ast.UnaryExpr) error {
	if node.Op == token.AND {
		switch operand := node.X.(type) {
		case *ast.Ident:
			// literal (es. '&h').
			symbol, ok := c.scopes.SymbolResolve(operand.Name)
			if !ok {
				return tables.NewCompilerError(c.fileSet, node, "undefined variable: %s", operand.Name)
			}
			opcodeId := native.OpNullId
			switch symbol.Scope() {
			case tables.LocalScope:
				opcodeId = native.OpLocalPtrGetId
			case tables.FreeScope:
				opcodeId = native.OpFreePtrGetId
			case tables.GlobalScope:
				opcodeId = native.OpGlobalPtrGetId
			default:
				return tables.NewCompilerError(c.fileSet, node, "cannot take the address of a unknown scope")
			}
			if _, err := c.scopes.SymbolEmit(node.Pos(), opcodeId, symbol.Index()); err != nil {
				return err
			}
		case *ast.CompositeLit:
			// literal (es. '&Home{...}').
			if err := c.compile(operand); err != nil {
				return err
			}
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_struct")
			if err != nil {
				return err
			}
			//tempSymbol.SetScope(tables.LocalScope)
			if err = c.scopes.SymbolEmitDefine(node.Pos(), tempSymbol); err != nil {
				return err
			}
			opcodeId := native.OpNullId
			switch tempSymbol.Scope() {
			case tables.LocalScope:
				opcodeId = native.OpLocalPtrGetId
			case tables.FreeScope:
				opcodeId = native.OpFreePtrGetId
			case tables.GlobalScope:
				opcodeId = native.OpGlobalPtrGetId
			default:
				return tables.NewCompilerError(c.fileSet, node, "cannot take the address of a unknown scope")
			}
			if _, err = c.scopes.SymbolEmit(node.Pos(), opcodeId, tempSymbol.Index()); err != nil {
				return err
			}
		default:
			return tables.NewCompilerError(c.fileSet, node, "cannot take the address of %T", node.X)
		}
		return nil
	}
	// logic for other unary operators (e.g. '!', '-', '^')
	if err := c.compile(node.X); err != nil {
		return err
	}
	z, ok := UnaryAdapterFor(node.Op)
	if !ok {
		return tables.NewCompilerError(c.fileSet, node, "unhandled unary op: %s", node.Op)
	}
	if _, err := c.scopes.SymbolEmit(node.Pos(), z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// BinaryExpr compiles a binary expression by evaluating both operands and applying the specified binary operator.
func (c *Expression) BinaryExpr(node *ast.BinaryExpr) error {
	if _, isCall := node.Y.(*ast.CallExpr); isCall {
		if err := c.compile(node.Y); err != nil {
			return err
		}
		tempSymbol, err := c.scopes.SymbolDefineUnique("__tmp_binary_rhs_")
		if err != nil {
			return err
		}
		if err = c.scopes.SymbolEmitSetAndPop(node.Pos(), tempSymbol); err != nil {
			return err
		}
		if err = c.compile(node.X); err != nil {
			return err
		}
		if err = c.scopes.SymbolEmitGet(node.Pos(), tempSymbol); err != nil {
			return err
		}
	} else {
		if err := c.compile(node.X); err != nil {
			return err
		}
		if err := c.compile(node.Y); err != nil {
			return err
		}
	}
	adapter, ok := BinaryAdapterFor(node.Op)
	if !ok {
		return tables.NewCompilerError(c.fileSet, node, "unhandled binary op: %s", node.Op)
	}
	if _, err := c.scopes.SymbolEmit(node.Pos(), adapter.op, adapter.arguments...); err != nil {
		return err
	}

	return nil
}

// SelectorExpr processes a selector expression (e.g., obj.field or pkg.Member) and emits necessary instructions.
func (c *Expression) SelectorExpr(node *ast.SelectorExpr) error {
	// Analyze the left-hand side of the dot to determine if it's a variable or a package.
	switch receiverIdent := node.X.(type) {
	case *ast.Ident:
		if c.imports.EmitPackage(node.Pos(), receiverIdent.Name, node.Sel.Name) {
			return nil
		}
		receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return tables.NewCompilerError(c.fileSet, node, "[SelectorExpr] undefined variable: %s", receiverIdent.Name)
		}
		if err := c.compile(node.X); err != nil {
			return err
		}
		// Interfaces
		if receiverSymbol.IsInterface() {
			return nil
		}
		// Symbols, imports etc...
		fieldName := node.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, keyConst); err != nil {
			return err
		}
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpIndexGetId); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr:
		rootExpr, selectors := c.flattenSelector(node)
		// 1. Root Analysis
		// We need to distinguish whether the root is a variable (Ident) or a computed expression
		if rootIdent, ok := rootExpr.(*ast.Ident); ok {
			// A. The root is an identifier (e.g. "user.name") - verify it exists in scope
			rootSymbol, ok := c.scopes.SymbolResolve(rootIdent.Name)
			if !ok {
				return tables.NewCompilerError(c.fileSet, node, "[SelectorExpr] undefined variable: %s", rootIdent.Name)
			}
			// Interface handling (as in the *ast.Ident case)
			if rootSymbol.IsInterface() {
				return nil
			}
		}
		// B. If it's not an Ident (e.g. "getFn().name"), we trust that c.compile(rootExpr)
		// will place the correct object on the stack.

		// 2. Root Compilation - this emits the LOAD of the variable or the Call/IndexExpr instructions
		if err := c.compile(rootExpr); err != nil {
			return err
		}

		// 3. Selector Chain Emission - perform the lookup of each field in sequence
		for _, fieldName := range selectors {
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
			// Push the field name
			if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, keyConst); err != nil {
				return err
			}
			// Field access (Stack: [Obj, "Key"] -> [Val])
			if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpIndexGetId); err != nil {
				return err
			}
		}
		return nil

	default:
		return tables.NewCompilerError(c.fileSet, node, "[SelectorExpr] unsupported receiver for selector expression: %T", node.X)
	}
}

// IncDecStmt processes increment and decrement statements, validating variables and emitting corresponding operations.
func (c *Expression) IncDecStmt(node *ast.IncDecStmt) error {
	ident, ok := node.X.(*ast.Ident)
	if !ok {
		return tables.NewCompilerError(c.fileSet, node, "unsupported IncDec statement for type %T", node.X)
	}
	symbol, ok := c.scopes.SymbolResolve(ident.Name)
	if !ok {
		return tables.NewCompilerError(c.fileSet, node, "undefined variable: %s", ident.Name)
	}
	if err := c.scopes.SymbolEmitGet(node.Pos(), symbol); err != nil {
		return err
	}
	// adds constant '1' to the stack
	constIndex := c.constants.Add("", c.gk.NewInt(objects.FrameStatic, 1))
	if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, constIndex); err != nil {
		return err
	}
	if node.Tok == token.INC {
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpArithmeticId, int(objects.OperatorAdd)); err != nil {
			return err
		}
	} else if node.Tok == token.DEC {
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpArithmeticId, int(objects.OperatorSub)); err != nil {
			return err
		}
	} else {
		return tables.NewCompilerError(c.fileSet, node, "unsupported IncDec token: %s", node.Tok)
	}
	if err := c.scopes.SymbolEmitSetAndPop(node.Pos(), symbol); err != nil {
		return err
	}
	return nil
}

// ParenExpr processes a parenthesized expression by compiling its inner expression.
func (c *Expression) ParenExpr(node *ast.ParenExpr) error {
	//panic("not implemented")
	return c.compile(node.X)
}

// SliceExpr compiles a slice expression, handling the base object and optional low/high index expressions.
func (c *Expression) SliceExpr(node *ast.SliceExpr) error {
	// 1. Compile the object to slice from (e.g. array)
	if err := c.compile(node.X); err != nil {
		return err
	}
	// 2. Compile the 'low' index (starting position)
	if node.Low != nil {
		if err := c.compile(node.Low); err != nil {
			return err
		}
	} else {
		// If 'low' index is omitted, push 'undefined' (OpNull)
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpNullId); err != nil {
			return err
		}
	}
	// 3. Compile the 'high' index (ending position)
	if node.High != nil {
		if err := c.compile(node.High); err != nil {
			return err
		}
	} else {
		// If 'high' index is omitted, push 'undefined' (OpNull)
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpNullId); err != nil {
			return err
		}
	}
	if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpIndexSliceId); err != nil {
		return err
	}
	return nil
}

// flattenSelector traverses a SelectorExpr recursively and extracts its base expression and all selector names.
func (c *Expression) flattenSelector(node *ast.SelectorExpr) (ast.Expr, []string) {
	var selectors []string
	curr := node
	for {
		selectors = append(selectors, curr.Sel.Name)
		if nested, ok := curr.X.(*ast.SelectorExpr); ok {
			curr = nested
		} else {
			break
		}
	}
	for i, j := 0, len(selectors)-1; i < j; i, j = i+1, j-1 {
		selectors[i], selectors[j] = selectors[j], selectors[i]
	}
	return curr.X, selectors
}
