package compiler

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/compiler/stdlib"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Scopes manages a collection of compilation scopes and the associated symbol table for nested compilation contexts.
type Scopes struct {
	constants   *Constants
	references  *Constants
	symbolTable *SymbolTable
	scopeIndex  int
	scopes      []*CompilationScope
}

// NewScopes initializes and returns a Scopes structure with a new symbol table, main compilation scope, and scope index set to 0.
func NewScopes() *Scopes {
	symbolTable := NewSymbolTable()
	for i, fn := range stdlib.GetAllBuiltinFunctions() {
		symbolTable.DefineBuiltin(fn.Name(), i)
	}
	return &Scopes{
		constants:   NewConstants(),
		references:  NewConstants(),
		symbolTable: symbolTable,
		scopeIndex:  0,
		scopes:      []*CompilationScope{NewCompilationScope()},
	}
}

// ReferencesAdd adds a constant object with the given id to the scope and returns its internal index.
func (c *Scopes) ReferencesAdd(id string, obj objects.IObject) int {
	return c.references.Add(id, obj)
}

// ReferencesGet retrieves the constant identified by the provided id and returns its value along with a boolean for existence.
func (c *Scopes) ReferencesGet(id string) (int, bool) {
	return c.references.Get(id)
}

// ReferencesRetrieve retrieves a slice of IObject constants from the scope's internal constants collection.
func (c *Scopes) ReferencesRetrieve() []objects.IObject {
	return c.references.Retrieve()
}

// ConstantsAdd adds a constant object with the given id to the scope and returns its internal index.
func (c *Scopes) ConstantsAdd(obj objects.IObject) int {
	return c.constants.Add("", obj)
}

// ConstantsRetrieve retrieves a slice of IObject constants from the scope's internal constants collection.
func (c *Scopes) ConstantsRetrieve() []objects.IObject {
	return c.constants.Retrieve()
}

// SymbolDefine defines a new symbol in the current scope and adds it to the symbol table. Returns the defined Symbol.
func (c *Scopes) SymbolDefine(symbol string) *Symbol {
	return c.symbolTable.Define(symbol)
}

// SymbolResolve attempts to find a symbol in the current scope and returns it along with a boolean indicating success.
func (c *Scopes) SymbolResolve(symbol string) (*Symbol, bool) {
	return c.symbolTable.Resolve(symbol)
}

// SymbolCount returns the number of symbol definitions in the symbol table.
func (c *Scopes) SymbolCount() int {
	return c.symbolTable.NumDefinitions()
}

// SymbolFreeConvert converts and retrieves free symbols from the symbol table as a slice of ObjectPointer.
func (c *Scopes) SymbolFreeConvert() []*objects.ObjectPointer {
	return c.symbolTable.ConvertFreeSymbols()
}

// SymbolFreeCount returns the count of free symbols in the current scope's symbol table.
func (c *Scopes) SymbolFreeCount() int {
	return c.symbolTable.FreeSymbolsLen()
}

// Current returns the current CompilationScope based on the internal scope index. Returns an error if the index is invalid.
func (c *Scopes) Current() (*CompilationScope, error) {
	if c.scopeIndex < 0 || c.scopeIndex >= len(c.scopes) {
		return nil, fmt.Errorf("invalid scope index: %d", c.scopeIndex)
	}
	return c.scopes[c.scopeIndex], nil
}

// AddInstructions appends the given byte slice to the current scope's instructions and returns the starting position or an error.
func (c *Scopes) AddInstructions(ins []byte) (int, error) {
	scope, err := c.Current()
	if err != nil {
		return 0, err
	}
	posNewInstruction := scope.InstructionsLen()
	if err = scope.InstructionsAppend(ins); err != nil {
		return 0, err
	}
	return posNewInstruction, nil
}

// SetLastInstruction updates the last emitted instruction for the current scope and tracks the previous one.
// It returns an error if the current scope cannot be retrieved.
func (c *Scopes) SetLastInstruction(op bytecode.Opcode, pos int) error {
	scope, err := c.Current()
	if err != nil {
		return err
	}
	previous := scope.LastInstruction()
	last := NewEmittedInstruction(op, pos)
	scope.SetPreviousInstruction(previous)
	scope.SetLastInstruction(last)
	return nil
}

// InstructionReplace replaces a section of instructions at the given position with new instructions in the current scope.
func (c *Scopes) InstructionReplace(pos int, newInstruction []byte) error {
	scope, err := c.Current()
	if err != nil {
		return err
	}
	if err = scope.InstructionsReplace(pos, newInstruction); err != nil {
		return err
	}
	return nil
}

// InstructionSet updates the instruction at the specified position in the current scope with the given byte value. Returns an error if the operation fails.
func (c *Scopes) InstructionSet(pos int, instruction byte) error {
	scope, err := c.Current()
	if err != nil {
		return err
	}
	if err = scope.InstructionsSet(pos, instruction); err != nil {
		return err
	}
	return nil
}

// InstructionGet retrieves an instruction from the current scope at the specified position. Returns the instruction or an error.
func (c *Scopes) InstructionGet(pos int) (byte, error) {
	scope, err := c.Current()
	if err != nil {
		return 0, err
	}
	data, err := scope.InstructionsGet(pos)
	if err != nil {
		return 0, err
	}
	return data, nil
}

// Enter creates a new compilation scope, updates the symbol table to be enclosed, and increments the scope index.
func (c *Scopes) Enter() error {
	if c.scopeIndex > maxScope {
		return fmt.Errorf("maximum scope depth exceeded: %d", maxScope)
	}
	scope := NewCompilationScope()
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
	return nil
}

// Leave removes the current scope and reverts to the previous one, returning the instructions of the removed scope.
func (c *Scopes) Leave() ([]byte, error) {
	scopesL := len(c.scopes)
	if scopesL <= 0 {
		return nil, errors.New("no scopes to leave")
	}
	scope, err := c.Current()
	if err != nil {
		return nil, err
	}
	c.scopes = c.scopes[:scopesL-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer()
	return scope.Instructions(), nil
}

// ChangeOperand modifies the operand of an instruction at the specified position within the current scope.
func (c *Scopes) ChangeOperand(opPos int, operand int) error {
	op, err := c.InstructionGet(opPos)
	if err != nil {
		return err
	}
	newInstruction := bytecode.MakeInstruction(op, operand)
	if err = c.InstructionReplace(opPos, newInstruction); err != nil {
		return err
	}
	return nil
}

// Emit generates and adds a new instruction to the current scope and updates the last emitted instruction info.
func (c *Scopes) Emit(op bytecode.Opcode, operands ...int) (int, error) {
	ins := bytecode.MakeInstruction(op, operands...)
	pos, err := c.AddInstructions(ins)
	if err != nil {
		return 0, err
	}
	if err = c.SetLastInstruction(op, pos); err != nil {
		return 0, err
	}
	return pos, nil
}

// EmitLiteral compiles a basic literal node into bytecode and adds the constant to the constant pool.
// It handles integer, float, and string literal types, emitting the corresponding constant operation.
// Returns an error if the literal type is not supported or an issue occurs during bytecode emission.
func (c *Scopes) EmitLiteral(node *ast.BasicLit) error {
	switch node.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(node.Value, 0, 64)
		if _, err := c.Emit(bytecode.OpConstant, c.ConstantsAdd(objects.NewInt(val))); err != nil {
			return err
		}
	case token.FLOAT:
		val, _ := strconv.ParseFloat(node.Value, 64)
		if _, err := c.Emit(bytecode.OpConstant, c.ConstantsAdd(objects.NewFloat(val))); err != nil {
			return err
		}
	case token.STRING:
		val, _ := strconv.Unquote(node.Value)
		s, err := objects.NewString(val)
		if err != nil {
			return err
		}
		if _, err = c.Emit(bytecode.OpConstant, c.ConstantsAdd(s)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled literal: %s", node.Kind)
	}
	return nil
}

// EmitSymbolSet generates bytecode instructions to set the value of a symbol in its appropriate scope (global, local, or free).
func (c *Scopes) EmitSymbolSet(s *Symbol) error {
	//fmt.Println("Emitting Symbol:", s)
	switch s.Scope {
	case GlobalScope:
		if _, err := c.Emit(bytecode.OpSetGlobal, s.Index); err != nil {
			return err
		}
	case LocalScope:
		if _, err := c.Emit(bytecode.OpSetLocal, s.Index); err != nil {
			return err
		}
	case FreeScope:
		if _, err := c.Emit(bytecode.OpSetFree, s.Index); err != nil {
			return err
		}
	}
	return nil
}

// EmitSymbolDefine emits the opcode for *defining* a variable.
func (c *Scopes) EmitSymbolDefine(s *Symbol) error {
	switch s.Scope {
	case GlobalScope:
		// For global scope, define is same as assign
		if _, err := c.Emit(bytecode.OpSetGlobal, s.Index); err != nil {
			return err
		}
	case LocalScope:
		// Use new opcode for local variables
		if _, err := c.Emit(bytecode.OpDefineLocal, s.Index); err != nil {
			return err
		}
	}
	return nil
}

// EmitSymbolGet generates bytecode instructions to retrieve a symbol's value based on its scope and index.
func (c *Scopes) EmitSymbolGet(s *Symbol) error {
	switch s.Scope {
	case GlobalScope:
		if _, err := c.Emit(bytecode.OpGetGlobal, s.Index); err != nil {
			return err
		}
	case LocalScope:
		if _, err := c.Emit(bytecode.OpGetLocal, s.Index); err != nil {
			return err
		}
	case BuiltinScope:
		if _, err := c.Emit(bytecode.OpGetBuiltin, s.Index); err != nil {
			return err
		}
	case FreeScope:
		if _, err := c.Emit(bytecode.OpGetFree, s.Index); err != nil {
			return err
		}
	}
	return nil
}

// EmitBinaryOp compiles a binary operation by emitting the corresponding bytecode based on the provided token operator.
func (c *Scopes) EmitBinaryOp(op token.Token) error {
	switch op {
	case token.ADD:
		if _, err := c.Emit(bytecode.OpBinaryOp, int(objects.OperatorAdd)); err != nil {
			return err
		}
	case token.SUB:
		if _, err := c.Emit(bytecode.OpBinaryOp, int(objects.OperatorSub)); err != nil {
			return err
		}
	case token.MUL:
		if _, err := c.Emit(bytecode.OpBinaryOp, int(objects.OperatorMul)); err != nil {
			return err
		}
	case token.QUO:
		if _, err := c.Emit(bytecode.OpBinaryOp, int(objects.OperatorQuo)); err != nil {
			return err
		}
	case token.EQL:
		if _, err := c.Emit(bytecode.OpEqual); err != nil {
			return err
		}
	case token.NEQ:
		if _, err := c.Emit(bytecode.OpNotEqual); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled binary op: %s", op)
	}
	return nil
}

// EmitUnaryOp compiles unary operations by emitting the corresponding bytecode for the given operator.
func (c *Scopes) EmitUnaryOp(op token.Token) error {
	switch op {
	case token.SUB:
		if _, err := c.Emit(bytecode.OpMinus); err != nil {
			return err
		}
	case token.NOT: // !
		if _, err := c.Emit(bytecode.OpLNot); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled unary op: %s", op)
	}
	return nil
}
