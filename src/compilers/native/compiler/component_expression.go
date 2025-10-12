package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

type Expression struct {
	gk        objects.IGateKeeper
	fileSet   *token.FileSet
	constants *tables.Constants
	scopes    *tables.Scopes
	imports   *Imports
	compile   func(node ast.Node) error
}

func NewExpression(gk objects.IGateKeeper, constants *tables.Constants, scopes *tables.Scopes, imports *Imports) *Expression {
	return &Expression{
		gk:        gk,
		constants: constants,
		imports:   imports,
		scopes:    scopes,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
func (c *Expression) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Prepare initializes the ControlFlow structure, ensuring it is ready for subsequent compilation tasks and operations.
func (c *Expression) Prepare() error {
	return nil
}

// Compile compiles the AST nodes using the configured compile function and returns an error if the process fails.
func (c *Expression) Compile() error {
	return nil
}

// UnaryExpr compiles a unary expression by evaluating the operand and applying the specified unary operator.
// It handles special cases for the address-of operator '&', ensuring correct pointer behavior based on operand type.
// Emits appropriate bytecode instructions for each unary operation or returns an error on unsupported cases.
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
			if _, err := c.scopes.Emit(node.Pos(), opcodeId, symbol.Index()); err != nil {
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
			if err = c.scopes.EmitSymbolDefine(node.Pos(), tempSymbol); err != nil {
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
			if _, err = c.scopes.Emit(node.Pos(), opcodeId, tempSymbol.Index()); err != nil {
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
	if _, err := c.scopes.Emit(node.Pos(), z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// BinaryExpr processes a binary expression node, compiling both operands and emitting the corresponding binary operation.
func (c *Expression) BinaryExpr(node *ast.BinaryExpr) error {
	if _, isCall := node.Y.(*ast.CallExpr); isCall {
		if err := c.compile(node.Y); err != nil {
			return err
		}
		tempSymbol, err := c.scopes.SymbolDefineUnique("__tmp_binary_rhs_")
		if err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSetAndPop(node.Pos(), tempSymbol); err != nil {
			return err
		}
		if err = c.compile(node.X); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolGet(node.Pos(), tempSymbol); err != nil {
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
	if _, err := c.scopes.Emit(node.Pos(), adapter.op, adapter.arguments...); err != nil {
		return err
	}

	return nil
}

// SelectorExpr processes a selector expression, resolving fields, methods, or package attributes.
// It distinguishes between struct field accesses and package-level selectors.
// Emits appropriate bytecode instructions for each case or returns an error if unsupported.
func (c *Expression) SelectorExpr(node *ast.SelectorExpr) error {
	// analyze the left-hand side of the dot to determine if it's a variable or a package.
	receiverIdent, ok := node.X.(*ast.Ident)
	if !ok {
		// currently not handling complex cases like a[0].field
		return tables.NewCompilerError(c.fileSet, node, "[SelectorExpr] unsupported receiver for selector expression: %T", node.X)
	}
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
	//interfaces
	if receiverSymbol.IsInterface() {
		return nil
	}
	//Symbols, imports etc...
	fieldName := node.Sel.Name
	keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
	if _, err := c.scopes.Emit(node.Pos(), native.OpConstantId, keyConst); err != nil {
		return err
	}
	if _, err := c.scopes.Emit(node.Pos(), native.OpIndexGetId); err != nil {
		return err
	}
	return nil
}

// IncDecStmt handles increment and decrement statements for identifiers, updating the corresponding variables and cleaning the stack.
func (c *Expression) IncDecStmt(node *ast.IncDecStmt) error {
	ident, ok := node.X.(*ast.Ident)
	if !ok {
		return tables.NewCompilerError(c.fileSet, node, "unsupported IncDec statement for type %T", node.X)
	}
	symbol, ok := c.scopes.SymbolResolve(ident.Name)
	if !ok {
		return tables.NewCompilerError(c.fileSet, node, "undefined variable: %s", ident.Name)
	}
	if err := c.scopes.EmitSymbolGet(node.Pos(), symbol); err != nil {
		return err
	}
	// adds constant '1' to the stack
	constIndex := c.constants.Add("", c.gk.NewInt(objects.FrameStatic, 1))
	if _, err := c.scopes.Emit(node.Pos(), native.OpConstantId, constIndex); err != nil {
		return err
	}
	if node.Tok == token.INC {
		if _, err := c.scopes.Emit(node.Pos(), native.OpArithmeticId, int(objects.OperatorAdd)); err != nil {
			return err
		}
	} else if node.Tok == token.DEC {
		if _, err := c.scopes.Emit(node.Pos(), native.OpArithmeticId, int(objects.OperatorSub)); err != nil {
			return err
		}
	} else {
		return tables.NewCompilerError(c.fileSet, node, "unsupported IncDec token: %s", node.Tok)
	}
	if err := c.scopes.EmitSymbolSetAndPop(node.Pos(), symbol); err != nil {
		return err
	}
	return nil
}

// ParenExpr processes a parenthesized expression and delegates the compilation to the contained sub-expression.
func (c *Expression) ParenExpr(node *ast.ParenExpr) error {
	//panic("not implemented")
	return c.compile(node.X)
}

// SliceExpr compiles a slice expression, processing the target, low, and high indices, and emits the slice bytecode.
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
		if _, err := c.scopes.Emit(node.Pos(), native.OpNullId); err != nil {
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
		if _, err := c.scopes.Emit(node.Pos(), native.OpNullId); err != nil {
			return err
		}
	}
	if _, err := c.scopes.Emit(node.Pos(), native.OpIndexSliceId); err != nil {
		return err
	}
	return nil
}
