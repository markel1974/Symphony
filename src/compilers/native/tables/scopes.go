package tables

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// maxScope defines the maximum allowable depth for compilation scopes to prevent excessive recursion or memory use.
const (
	maxScope = 1024
)

// Scopes manages a collection of compilation scopes and the associated symbol table for nested compilation contexts.
type Scopes struct {
	gk                   objects.IGateKeeper
	op                   opcodes.IOpcodes
	constants            *Constants
	symbolTable          *SymbolTable
	initSymbolTable      *SymbolTable
	scopeIndex           int
	compilations         []*CompilationScope
	initCompilationScope *CompilationScope
	uniqueCounter        int
}

// NewScopes initializes and returns a Scopes structure with a new symbol table, main compilation scope, and scope index set to 0.
func NewScopes(gk objects.IGateKeeper, op opcodes.IOpcodes, constants *Constants) *Scopes {
	c := &Scopes{
		gk:                   gk,
		op:                   op,
		constants:            constants,
		initSymbolTable:      NewSymbolTable(UnknownScope),
		symbolTable:          nil,
		scopeIndex:           0,
		compilations:         []*CompilationScope{},
		initCompilationScope: NewCompilationScope(),
		uniqueCounter:        0,
	}
	c.compilations = append(c.compilations, c.initCompilationScope)
	c.symbolTable = c.initSymbolTable
	return c
}

func (c *Scopes) CreateGlobals() []objects.IObject {
	ret := make([]objects.IObject, len(c.initSymbolTable.definitions))
	for _, obj := range c.initSymbolTable.definitions {
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

// SetRootIndex resets the scope index to 0, designating it as the root index for the current Scopes instance.
func (c *Scopes) SetRootIndex() {
	c.scopeIndex = 0
}

// IsRootScope checks if the current scope is the root scope by verifying if the scope index is zero.
func (c *Scopes) IsRootScope() bool {
	return c.scopeIndex == 0
}

// SymbolDefine defines a new symbol within the symbol table and returns the created symbol or an error if it fails.
func (c *Scopes) SymbolDefine(name string) (*Symbol, error) {
	return c.symbolTable.Define(name)
}

// SymbolDefineUnique ensures the given symbol is uniquely defined and returns it or an error if the operation fails.
func (c *Scopes) SymbolDefineUnique(name string) (*Symbol, error) {
	uniqueName := name + strconv.Itoa(c.uniqueCounter)
	c.uniqueCounter++
	return c.symbolTable.Define(uniqueName)
}

// SymbolDefineType defines a new type symbol with the given name in the symbol table and returns the created symbol or an error.
func (c *Scopes) SymbolDefineType(name string) (*Symbol, error) {
	return c.symbolTable.DefineType(name)
}

// SymbolDefineConst defines a constant in the symbol table with the given index and symbol name. Returns the defined symbol or an error.
func (c *Scopes) SymbolDefineConst(name string, object objects.IObject) (*Symbol, error) {
	constIdx := c.constants.Add(name, object)
	return c.symbolTable.DefineConst(constIdx, name)
}

// SymbolRebuildScope rebuilds and updates the specified symbol with a new scope in the symbol table, returning the updated symbol.
func (c *Scopes) SymbolRebuildScope(symbol string, scope SymbolScope) (*Symbol, bool) {
	return c.symbolTable.RebuildScope(symbol, scope)
}

// SymbolResolve attempts to find a symbol in the current scope and returns it along with a boolean indicating success.
func (c *Scopes) SymbolResolve(symbol string) (*Symbol, bool) {
	if obj, ok := c.initSymbolTable.Resolve(symbol); ok {
		return obj, true
	}
	return c.symbolTable.Resolve(symbol)
}

// SymbolCount returns the number of symbol definitions in the symbol table.
func (c *Scopes) SymbolCount() int {
	return c.symbolTable.Count()
}

// SymbolFree retrieves a slice of integers representing free symbols from the internal symbol table within Scopes.
func (c *Scopes) SymbolFree() []int {
	return c.symbolTable.FreeSymbols()
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
func (c *Scopes) InstructionSetLast(op opcodes.OpcodeId, pos int) error {
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
func (c *Scopes) InstructionGet(pos int) (opcodes.OpcodeId, error) {
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
func (c *Scopes) Enter(defaultScope SymbolScope, funcName string) error {
	if c.scopeIndex > maxScope {
		return fmt.Errorf("maximum scope depth exceeded: %d", maxScope)
	}
	scope := NewCompilationScope()
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable, defaultScope, funcName)
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
	newInstruction, err := c.op.Compile(op, operand)
	if err != nil {
		return err
	}
	if err = c.InstructionReplace(opPos, newInstruction); err != nil {
		return err
	}
	return nil
}

// Emit generates and adds a new instruction to the current scope and updates the last emitted instruction info.
func (c *Scopes) Emit(op opcodes.OpcodeId, operands ...int) (int, error) {
	ins, err := c.op.Compile(op, operands...)
	if err != nil {
		return 0, err
	}
	pos, err := c.InstructionsAdd(ins)
	if err != nil {
		return 0, err
	}
	if err = c.InstructionSetLast(op, pos); err != nil {
		return 0, err
	}
	return pos, nil
}

// EmitAndPop emits the given opcode with operands, followed by a pop operation, and returns an error if either fails.
func (c *Scopes) EmitAndPop(op opcodes.OpcodeId, operands ...int) error {
	if _, err := c.Emit(op, operands...); err != nil {
		return err
	}
	if _, err := c.Emit(native.OpPopId); err != nil {
		return err
	}
	return nil
}

// EmitSymbolDefineAndPop defines a symbol within the current scope, emits its bytecode, and pops it off the stack.
func (c *Scopes) EmitSymbolDefineAndPop(s *Symbol) error {
	if err := c.EmitSymbolDefine(s); err != nil {
		return err
	}
	if _, err := c.Emit(native.OpPopId); err != nil {
		return err
	}
	return nil
}

// EmitSymbolDefine emits the opcode for *defining* a variable.
func (c *Scopes) EmitSymbolDefine(s *Symbol) error {
	if s.Constant() {
		return fmt.Errorf("cannot define constant symbol: %s", s.Name())
	}
	var op opcodes.OpcodeId
	switch s.Scope() {
	case GlobalScope:
		op = native.OpGlobalDefineId
	case LocalScope:
		op = native.OpLocalDefineId // Use new opcode for local variables
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
	if s.Constant() {
		return fmt.Errorf("cannot set constant symbol: %s", s.Name())
	}
	var op opcodes.OpcodeId
	switch s.Scope() {
	case GlobalScope:
		op = native.OpGlobalSetId
	case LocalScope:
		op = native.OpLocalSetId
	case FreeScope:
		op = native.OpFreeSetId
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
	if s.Constant() {
		_, err := c.Emit(native.OpConstantId, s.Index())
		return err
	}
	var op opcodes.OpcodeId
	switch s.Scope() {
	case GlobalScope:
		op = native.OpGlobalGetId
	case LocalScope:
		op = native.OpLocalGetId
	case FreeScope:
		op = native.OpFreeGetId
	default:
		return fmt.Errorf("unsupported symbol scope: %v", s.Scope())
	}
	if _, err := c.Emit(op, s.Index()); err != nil {
		return err
	}
	return nil
}

// EmitSymbolSetAndPop sets a symbol and then emits a pop operation, returning an error if either fails.
func (c *Scopes) EmitSymbolSetAndPop(s *Symbol) error {
	if err := c.EmitSymbolSet(s); err != nil {
		return err
	}
	if _, err := c.Emit(native.OpPopId); err != nil {
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
