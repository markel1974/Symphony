package compiler

import (
	"errors"
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Scopes manages a collection of compilation scopes and the associated symbol table for nested compilation contexts.
type Scopes struct {
	factory              objects.IGateKeeper
	op                   *bytecode.Opcodes
	symbolTable          *SymbolTable
	initSymbolTable      *SymbolTable
	scopeIndex           int
	compilations         []*CompilationScope
	initCompilationScope *CompilationScope
}

// NewScopes initializes and returns a Scopes structure with a new symbol table, main compilation scope, and scope index set to 0.
func NewScopes(factory objects.IGateKeeper, op *bytecode.Opcodes) *Scopes {
	c := &Scopes{
		factory:              factory,
		op:                   op,
		initSymbolTable:      NewSymbolTable(),
		symbolTable:          nil,
		scopeIndex:           0,
		compilations:         []*CompilationScope{},
		initCompilationScope: NewCompilationScope(),
	}
	c.compilations = append(c.compilations, c.initCompilationScope)
	c.symbolTable = c.initSymbolTable
	return c
}

// SymbolDefineUnique defines a unique symbol in the current symbol table with the specified scope and object type.
func (c *Scopes) SymbolDefineUnique(symbol string, scope SymbolScope, isObj bool) (*Symbol, error) {
	return c.symbolTable.DefineUnique(symbol, scope, isObj)
}

// SymbolDefine defines a new symbol in the current symbol table with the specified name, scope, and struct flag.
// It returns the created symbol or an error if the operation fails.
func (c *Scopes) SymbolDefine(symbol string, scope SymbolScope, isStruct bool) (*Symbol, error) {
	//if symbol == "item" {
	//	fmt.Println("symbol item found!!!!")
	//}
	return c.symbolTable.Define(symbol, scope, isStruct)
}

// SymbolReset resets the symbol in the current scope with the specified name, scope, and struct flag.
func (c *Scopes) SymbolReset(symbol string, scope SymbolScope, isStruct bool) (*Symbol, error) {
	return c.symbolTable.Reset(symbol, scope, isStruct)
}

// SymbolResolve attempts to find a symbol in the current scope and returns it along with a boolean indicating success.
func (c *Scopes) SymbolResolve(symbol string) (*Symbol, bool) {
	if obj, ok := c.initSymbolTable.Resolve(symbol); ok {
		return obj, true
	}
	return c.symbolTable.Resolve(symbol)
}

// SymbolScope returns the current scope's symbol table's scope.
func (c *Scopes) SymbolScope() SymbolScope {
	return c.symbolTable.Scope()
}

// SymbolCount returns the number of symbol definitions in the symbol table.
func (c *Scopes) SymbolCount() int {
	return c.symbolTable.Count()
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
	if c.scopeIndex < 0 || c.scopeIndex >= len(c.compilations) {
		return nil, fmt.Errorf("invalid scope index: %d", c.scopeIndex)
	}
	return c.compilations[c.scopeIndex], nil
}

func (c *Scopes) InstructionsInit() ([]byte, int) {
	return c.initCompilationScope.Instructions(), c.initSymbolTable.Count()
}

// InstructionsAdd appends the given byte slice to the current scope's instructions and returns the starting position or an error.
func (c *Scopes) InstructionsAdd(ins []byte) (int, error) {
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

// InstructionSetLast updates the last emitted instruction for the current scope and tracks the previous one.
// It returns an error if the current scope cannot be retrieved.
func (c *Scopes) InstructionSetLast(op bytecode.Opcode, pos int) error {
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
func (c *Scopes) Enter(structName string, funcName string) error {
	if c.scopeIndex > maxScope {
		return fmt.Errorf("maximum scope depth exceeded: %d", maxScope)
	}
	scope := NewCompilationScope()
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable, structName, funcName)
	c.compilations = append(c.compilations, scope)
	c.scopeIndex++
	return nil
}

// Leave removes the current scope and reverts to the previous one, returning the instructions of the removed scope.
func (c *Scopes) Leave() ([]byte, error) {
	scopesLen := len(c.compilations)
	if scopesLen <= 0 {
		return nil, errors.New("no scopes to leave")
	}
	scope, err := c.Current()
	if err != nil {
		return nil, err
	}
	c.symbolTable = c.symbolTable.Outer()
	c.compilations = c.compilations[:scopesLen-1]
	c.scopeIndex--
	return scope.Instructions(), nil
}

// ChangeOperand modifies the operand of an instruction at the specified position within the current scope.
func (c *Scopes) ChangeOperand(opPos int, operand int) error {
	op, err := c.InstructionGet(opPos)
	if err != nil {
		return err
	}
	newInstruction := c.op.CompileInstruction(op, operand)
	if err = c.InstructionReplace(opPos, newInstruction); err != nil {
		return err
	}
	return nil
}

// Emit generates and adds a new instruction to the current scope and updates the last emitted instruction info.
func (c *Scopes) Emit(op bytecode.Opcode, operands ...int) (int, error) {
	ins := c.op.CompileInstruction(op, operands...)
	pos, err := c.InstructionsAdd(ins)
	if err != nil {
		return 0, err
	}
	if err = c.InstructionSetLast(op, pos); err != nil {
		return 0, err
	}
	return pos, nil
}

// EmitSymbolDefine emits the opcode for *defining* a variable.
func (c *Scopes) EmitSymbolDefine(s *Symbol) error {
	var op bytecode.Opcode
	switch s.Scope() {
	case GlobalScope:
		op = bytecode.OpSetGlobal
	case LocalScope:
		op = bytecode.OpDefineLocal // Use new opcode for local variables
	default:
		return fmt.Errorf("unsupported symbol scope: %v", s.Scope())
	}
	if _, err := c.Emit(op, s.Index()); err != nil {
		return err
	}
	return nil
}

// EmitSymbolSet generates bytecode instructions to set the value of a symbol in its appropriate scope (global, local, or free).
func (c *Scopes) EmitSymbolSet(s *Symbol) error {
	var op bytecode.Opcode
	switch s.Scope() {
	case GlobalScope:
		op = bytecode.OpSetGlobal
	case LocalScope:
		op = bytecode.OpSetLocal
	case FreeScope:
		op = bytecode.OpSetFree
	default:
		return fmt.Errorf("unsupported symbol scope: %v", s.Scope())
	}
	if _, err := c.Emit(op, s.Index()); err != nil {
		return err
	}
	return nil
}

// EmitSymbolGet generates bytecode instructions to retrieve a symbol's value based on its scope and index.
func (c *Scopes) EmitSymbolGet(s *Symbol) error {
	var op bytecode.Opcode
	switch s.Scope() {
	case GlobalScope:
		op = bytecode.OpGetGlobal
	case LocalScope:
		op = bytecode.OpGetLocal
	case FreeScope:
		op = bytecode.OpGetFree
	default:
		return fmt.Errorf("unsupported symbol scope: %v", s.Scope())
	}
	if _, err := c.Emit(op, s.Index()); err != nil {
		return err
	}
	return nil
}

// Print prints the contents of the Scopes structure to the console.
func (c *Scopes) Print(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "----- Symbols -----")
	c.symbolTable.Print(writer)
	_, _ = fmt.Fprintf(writer, "--------------------")
}
