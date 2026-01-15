package tables

import (
	"fmt"
	"go/token"
	"io"
	"strconv"

	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
	"github.com/markel1974/symphony/src/vm/sequencers/native"
)

// maxScope defines the maximum allowable depth of nested compilation scopes to prevent excessive recursion or overflow.
const (
	maxScope = 1024
)

// Scopes is a structure for managing compilation contexts, symbol resolution, and constant storage in program execution.
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

// NewScopes initializes and returns a new Scopes instance configured with the provided IGateKeeper, IOpcodes, and Constants.
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

// Enter creates a new compilation scope, updates the symbol table, and increments the scope index.
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

// Leave exits the current compilation scope, restoring the previous symbol table and returning scope data.
func (c *Scopes) Leave() ([]byte, map[int]int, error) {
	scopesLen := len(c.compilations)
	if scopesLen <= 0 {
		return nil, nil, fmt.Errorf("no scopes to leave")
	}
	scope, err := c.Current()
	if err != nil {
		return nil, nil, err
	}
	c.symbolTable = c.symbolTable.Outer()
	c.compilations = c.compilations[:scopesLen-1]
	c.scopeIndex--
	return scope.Instructions(), scope.Source(), nil
}

// Current returns the current CompilationScope based on the current scope index or an error if the index is invalid.
func (c *Scopes) Current() (*CompilationScope, error) {
	if c.scopeIndex < 0 || c.scopeIndex >= len(c.compilations) {
		return nil, fmt.Errorf("invalid scope index: %d", c.scopeIndex)
	}
	return c.compilations[c.scopeIndex], nil
}

// SetRootIndex resets the scope index of the Scopes object to its root (0).
func (c *Scopes) SetRootIndex() {
	c.scopeIndex = 0
}

// IsRootScope checks if the current scope is the root scope by verifying if the scope index is zero.
func (c *Scopes) IsRootScope() bool {
	return c.scopeIndex == 0
}

// InstructionsInit retrieves the instructions from the initial compilation scope and the count of symbols in the initial symbol table.
func (c *Scopes) InstructionsInit() ([]byte, int) {
	return c.initCompilationScope.Instructions(), c.initSymbolTable.Count()
}

// InstructionsAppend appends given instructions to the current scope and updates source mappings. Returns the position of the new instruction or an error if appending fails.
func (c *Scopes) InstructionsAppend(ins []byte, source map[int]int) (int, error) {
	scope, err := c.Current()
	if err != nil {
		return 0, err
	}
	posNewInstruction, err := scope.InstructionsAppend(ins)
	if err != nil {
		return 0, err
	}
	for opPos, tPos := range source {
		scope.SetSource(opPos, token.Pos(tPos))
	}
	return posNewInstruction, nil
}

// InstructionGet retrieves the OpcodeId at the specified position within the current compilation scope. Returns an error if unsuccessful.
func (c *Scopes) InstructionGet(opPos int) (opcodes.OpcodeId, error) {
	scope, err := c.Current()
	if err != nil {
		return 0, err
	}
	data, err := scope.InstructionsGet(opPos)
	if err != nil {
		return 0, err
	}
	return data, nil
}

// InstructionsChangeOperand modifies a specific instruction's operand at the given position in the current compilation scope.
func (c *Scopes) InstructionsChangeOperand(p token.Pos, opPos int, operand int) error {
	scope, err := c.Current()
	if err != nil {
		return err
	}
	data, err := scope.InstructionsGet(opPos)
	if err != nil {
		return err
	}
	newInstruction, err := c.op.Bytecode(data, operand)
	if err != nil {
		return err
	}
	if err = scope.InstructionsReplace(opPos, newInstruction); err != nil {
		return err
	}
	scope.SetSource(opPos, p)
	return nil
}

// SymbolDefine defines a new symbol in the current symbol table and returns it along with any error encountered.
func (c *Scopes) SymbolDefine(name string) (*Symbol, error) {
	return c.symbolTable.Define(name)
}

// SymbolDefineUnique creates a uniquely named symbol by appending a counter to the given name and defines it in the symbol table.
// It increments the unique counter to ensure uniqueness for subsequent definitions. Returns the newly defined symbol or an error.
func (c *Scopes) SymbolDefineUnique(name string) (*Symbol, error) {
	uniqueName := name + strconv.Itoa(c.uniqueCounter)
	c.uniqueCounter++
	return c.symbolTable.Define(uniqueName)
}

// SymbolDefineType defines a new type symbol with the specified name and returns the symbol or an error if it fails.
func (c *Scopes) SymbolDefineType(name string) (*Symbol, error) {
	return c.symbolTable.DefineType(name)
}

// SymbolDefineConst defines a constant symbol with the specified name and object, returning the created symbol or an error.
func (c *Scopes) SymbolDefineConst(name string, object objects.IObject) (*Symbol, error) {
	constIdx := c.constants.Add(name, object)
	return c.symbolTable.DefineConst(constIdx, name)
}

// SymbolRebuildScope updates the scope of the given symbol and returns the updated Symbol and a success flag.
func (c *Scopes) SymbolRebuildScope(symbol string, scope SymbolScope) (*Symbol, bool) {
	return c.symbolTable.RebuildScope(symbol, scope)
}

// SymbolResolve attempts to resolve the given symbol in the current or initial symbol table and returns it if found.
func (c *Scopes) SymbolResolve(symbol string) (*Symbol, bool) {
	if obj, ok := c.initSymbolTable.Resolve(symbol); ok {
		return obj, true
	}
	return c.symbolTable.Resolve(symbol)
}

// SymbolResolveOrDefine resolves a symbol by name or defines it if not found, returning the symbol or an error if defining fails.
func (c *Scopes) SymbolResolveOrDefine(symbol string) (*Symbol, error) {
	if s, ok := c.SymbolResolve(symbol); ok {
		return s, nil
	}
	s, err := c.SymbolDefine(symbol)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// SymbolCount returns the total number of symbols currently defined in the symbol table of the Scopes instance.
func (c *Scopes) SymbolCount() int {
	return c.symbolTable.Count()
}

// SymbolFree returns a slice of indices representing the free symbols defined in the current symbol table.
func (c *Scopes) SymbolFree() []int {
	return c.symbolTable.FreeSymbols()
}

// SymbolGlobalsCreate creates global symbols by initializing a slice with objects from the initial symbol table definitions.
func (c *Scopes) SymbolGlobalsCreate() []objects.IObject {
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

// SymbolEmit generates bytecode instructions for the specified opcode and operands, appends them to the current scope,
// and associates the instruction position with the given source position. Returns the position of the new instruction
// or an error if the operation fails.
func (c *Scopes) SymbolEmit(p token.Pos, op opcodes.OpcodeId, operands ...int) (int, error) {
	scope, err := c.Current()
	if err != nil {
		return 0, err
	}
	ins, err := c.op.Bytecode(op, operands...)
	if err != nil {
		return 0, err
	}
	opPos, err := scope.InstructionsAppend(ins)
	if err != nil {
		return 0, err
	}
	previous := scope.LastInstruction()
	last := NewEmittedInstruction(op, opPos)
	scope.SetPreviousInstruction(previous)
	scope.SetLastInstruction(last)

	scope.SetSource(opPos, p)

	return opPos, nil
}

// SymbolEmitAndPop emits a given opcode and operands and immediately emits a pop operation, returning any encountered error.
func (c *Scopes) SymbolEmitAndPop(p token.Pos, op opcodes.OpcodeId, operands ...int) error {
	if _, err := c.SymbolEmit(p, op, operands...); err != nil {
		return err
	}
	if _, err := c.SymbolEmit(p, native.OpPopId); err != nil {
		return err
	}
	return nil
}

// SymbolEmitDefineAndPop defines a symbol in the current scope, emits the associated bytecode, and pops it from the stack.
func (c *Scopes) SymbolEmitDefineAndPop(p token.Pos, s *Symbol) error {
	if err := c.SymbolEmitDefine(p, s); err != nil {
		return err
	}
	if _, err := c.SymbolEmit(p, native.OpPopId); err != nil {
		return err
	}
	return nil
}

// SymbolEmitDefine emits bytecode to define a symbol within the current scope. Returns an error for unsupported scopes.
func (c *Scopes) SymbolEmitDefine(p token.Pos, s *Symbol) error {
	if s.Constant() {
		return fmt.Errorf("cannot define constant symbol: %s", s.Name())
	}
	var op opcodes.OpcodeId
	switch s.Scope() {
	case GlobalScope:
		op = native.OpGlobalDefineId
	case LocalScope:
		op = native.OpLocalDefineId
	default:
		return fmt.Errorf("unsupported symbol scope: %v", s.Scope())
	}
	if _, err := c.SymbolEmit(p, op, s.Index()); err != nil {
		return err
	}
	return nil
}

// SymbolEmitSet emits the opcode to set the value of a given symbol, returning an error if the symbol is a constant or unsupported.
func (c *Scopes) SymbolEmitSet(p token.Pos, s *Symbol) error {
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
	if _, err := c.SymbolEmit(p, op, s.Index()); err != nil {
		return err
	}
	return nil
}

// SymbolEmitGet emits bytecode for retrieving the value of the given symbol based on its scope and type.
func (c *Scopes) SymbolEmitGet(p token.Pos, s *Symbol) error {
	if s.Constant() {
		_, err := c.SymbolEmit(p, native.OpConstantId, s.Index())
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
	if _, err := c.SymbolEmit(p, op, s.Index()); err != nil {
		return err
	}
	return nil
}

// SymbolEmitSetAndPop emits the instructions to set a symbol and then pops it from the stack. Returns an error if any step fails.
func (c *Scopes) SymbolEmitSetAndPop(p token.Pos, s *Symbol) error {
	if err := c.SymbolEmitSet(p, s); err != nil {
		return err
	}
	if _, err := c.SymbolEmit(p, native.OpPopId); err != nil {
		return err
	}
	return nil
}

// Print writes the current state of the symbols and their table to the provided writer, formatted for readability.
func (c *Scopes) Print(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "----- Symbols -----")
	c.symbolTable.Print(writer)
	_, _ = fmt.Fprintf(writer, "--------------------")
}
