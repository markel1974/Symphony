package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

type Expression struct {
	gk        objects.IGateKeeper
	fileSet   *token.FileSet
	constants *Constants
	scopes    *Scopes
	imports   *Imports
	compile   func(node ast.Node) error
}

func NewExpression(gk objects.IGateKeeper, constants *Constants, scopes *Scopes, imports *Imports) *Expression {
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
				return NewCompilerError(c.fileSet, node, "undefined variable: %s", operand.Name)
			}
			switch symbol.Scope() {
			case LocalScope:
				if _, err := c.scopes.Emit(bytecode.OpLocalPtrGet, symbol.Index()); err != nil {
					return err
				}
			case FreeScope:
				if _, err := c.scopes.Emit(bytecode.OpFreePtrGet, symbol.Index()); err != nil {
					return err
				}
			default:
				return NewCompilerError(c.fileSet, node, "cannot take the address of a global variable")
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
			tempSymbol.SetScope(LocalScope)
			if err = c.scopes.EmitSymbolDefine(tempSymbol); err != nil {
				return err
			}
			if _, err = c.scopes.Emit(bytecode.OpLocalPtrGet, tempSymbol.Index()); err != nil {
				return err
			}
		default:
			return NewCompilerError(c.fileSet, node, "cannot take the address of %T", node.X)
		}
		return nil
	}
	// existing logic for other unary operators (e.g. '!', '-', '^')
	if err := c.compile(node.X); err != nil {
		return err
	}
	z, ok := UnaryAdapterFor(node.Op)
	if !ok {
		return NewCompilerError(c.fileSet, node, "unhandled unary op: %s", node.Op)
	}
	if _, err := c.scopes.Emit(z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// BinaryExpr processes a binary expression node, compiling both operands and emitting the corresponding binary operation.
func (c *Expression) BinaryExpr(node *ast.BinaryExpr) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if err := c.compile(node.Y); err != nil {
		return err
	}
	z, ok := BinaryAdapterFor(node.Op)
	if !ok {
		return NewCompilerError(c.fileSet, node, "unhandled binary op: %s", node.Op)
	}
	if _, err := c.scopes.Emit(z.op, z.arguments...); err != nil {
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
		return NewCompilerError(c.fileSet, node, "[SelectorExpr] unsupported receiver for selector expression: %T", node.X)
	}
	if c.imports.Emit(receiverIdent.Name, node.Sel.Name) {
		return nil
	}
	receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
	if !ok {
		return NewCompilerError(c.fileSet, node, "[SelectorExpr] undefined variable: %s", receiverIdent.Name)
	}
	if receiverSymbol.IsStruct() { // struct
		if err := c.compile(node.X); err != nil {
			return err
		}
		fieldName := node.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpIndex); err != nil {
			return err
		}
		return nil
	}
	return NewCompilerError(c.fileSet, node, "[SelectorExpr] unsupported selector expression for symbol %s", receiverSymbol.Name())
}

// IncDecStmt handles increment and decrement statements for identifiers, updating the corresponding variables and cleaning the stack.
func (c *Expression) IncDecStmt(node *ast.IncDecStmt) error {
	ident, ok := node.X.(*ast.Ident)
	if !ok {
		return NewCompilerError(c.fileSet, node, "unsupported IncDec statement for type %T", node.X)
	}
	symbol, ok := c.scopes.SymbolResolve(ident.Name)
	if !ok {
		return NewCompilerError(c.fileSet, node, "undefined variable: %s", ident.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	// adds constant '1' to the stack
	constIndex := c.constants.Add("", c.gk.NewInt(objects.FrameStatic, 1))
	if _, err := c.scopes.Emit(bytecode.OpConstant, constIndex); err != nil {
		return err
	}
	if node.Tok == token.INC {
		if _, err := c.scopes.Emit(bytecode.OpBinary, int(objects.OperatorAdd)); err != nil {
			return err
		}
	} else if node.Tok == token.DEC {
		if _, err := c.scopes.Emit(bytecode.OpBinary, int(objects.OperatorSub)); err != nil {
			return err
		}
	} else {
		return NewCompilerError(c.fileSet, node, "unsupported IncDec token: %s", node.Tok)
	}
	if err := c.scopes.EmitSymbolSet(symbol); err != nil {
		return err
	}
	// the increment/decrement operation leaves the result on the stack. Since this is a statement, we need to clean this value.
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}
