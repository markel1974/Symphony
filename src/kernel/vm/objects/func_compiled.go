package objects

import (
	"encoding/gob"
	"fmt"
)

const (
	FuncCompiledDef   = "func_compiled"
	FuncCompiledLabel = "<" + FuncCompiledDef + ">"
)

func init() {
	gob.Register(&FuncCompiled{})
}

// Instructions represent a collection of bytecode instructions stored as a byte slice.
type Instructions struct {
	data []byte
}

// NewInstructions creates a new instance of Instructions with the provided byte slice data.
func NewInstructions(data []byte) *Instructions {
	return &Instructions{data: data}
}

// Copy creates and returns a new Instructions instance with a duplicated copy of the original data slice.
func (i *Instructions) Copy() *Instructions {
	out := NewInstructions(nil)
	out.data = append([]byte{}, i.data...)
	return out
}

// Data returns the internal byte slice representing the instructions.
func (i *Instructions) Data() []byte {
	return i.data
}

// Length returns the number of elements in the Instructions' data slice.
func (i *Instructions) Length() int {
	return len(i.data)
}

// Get8 retrieves a single byte from the instructions at the specified index.
// Returns an error if the index is out of bounds.
func (i *Instructions) Get8(base int) (uint8, error) {
	if base < 0 || base >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction index: %d", base)
	}
	return i.data[base], nil
}

// Get16 retrieves a 16-bit unsigned integer composed of two bytes from the instructions at indices low and low-1.
// Returns an error if the provided indices are out of bounds.
func (i *Instructions) Get16(base int) (uint16, error) {
	high := base - 1
	if base < 0 || base >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction low index: %d", base)
	}
	if high < 0 || high >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction high index: %d", high)
	}
	return uint16(i.data[base]) | uint16(i.data[high])<<8, nil
}

// FuncCompiled represents a compiled function with bytecode instructions, metadata, and associated free variables.
type FuncCompiled struct {
	gk            IGateKeeper
	frame         int
	name          string
	instructions  *Instructions
	numLocals     int
	numParameters int
	varArgs       bool
	sourceMap     map[int]int
	free          []*ObjectPointer
}

// NewFunctionCompiled creates a new instance of FuncCompiled with the given instructions, locals, parameters, varArgs, sourceMap, and free vars.
func newFuncCompiled(factory IGateKeeper, frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) IObject {
	if sourceMap == nil {
		sourceMap = make(map[int]int)
	}
	return &FuncCompiled{
		gk:            factory,
		frame:         frame,
		name:          name,
		instructions:  NewInstructions(instructions),
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		sourceMap:     sourceMap,
		free:          free,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *FuncCompiled) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *FuncCompiled) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation on the object using the specified operator and right-hand-side operand.
// Returns an error if the operator is invalid.
func (o *FuncCompiled) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the provided operator and right-hand side operand.
// Returns an error if the operation is unsupported.
func (o *FuncCompiled) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *FuncCompiled) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *FuncCompiled) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *FuncCompiled) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncCompiled) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *FuncCompiled) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *FuncCompiled) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// Length returns the length of the Int object.
func (o *FuncCompiled) Length() int {
	return 0
}

// Name returns the name of the compiled function.
func (o *FuncCompiled) Name() string {
	return o.name
}

// Data returns the bytecode instructions of the compiled function.
func (o *FuncCompiled) Data() []byte {
	return o.instructions.Data()
}

// Instructions return the bytecode instructions associated with the compiled function.
func (o *FuncCompiled) Instructions() *Instructions {
	return o.instructions
}

// NumLocals returns the number of local variables required by the compiled function.
func (o *FuncCompiled) NumLocals() int {
	return o.numLocals
}

// NumParameters returns the total number of parameters required by the compiled function.
func (o *FuncCompiled) NumParameters() int {
	return o.numParameters
}

// VarArgs returns true if the function accepts a variable number of arguments, otherwise false.
func (o *FuncCompiled) VarArgs() bool {
	return o.varArgs
}

// Free returns the slice of *ObjectPointer representing the free variables captured by the compiled function.
func (o *FuncCompiled) Free() []*ObjectPointer {
	return o.free
}

// TypeName returns the type name of the object as a string, specifically "compiled-function".
func (o *FuncCompiled) TypeName() string {
	return FuncCompiledDef
}

// String returns a human-readable string representation of the FuncCompiled object.
func (o *FuncCompiled) String() string {
	return FuncCompiledLabel
}

// Copy creates and returns a new instance of FuncCompiled, duplicating its state, except for its variable pointers.
func (o *FuncCompiled) Copy(frame int, _ int) IObject {
	obj := o.GateKeeper().NewFuncCompiled(frame, o.name, nil, o.numLocals, o.numParameters, o.varArgs, nil, nil)
	ret, ok := obj.(*FuncCompiled)
	if !ok {
		return o.GateKeeper().UndefinedValue()
	}
	ret.instructions = o.instructions.Copy()
	ret.free = append([]*ObjectPointer{}, o.free...)
	return ret
}

// Equals checks if the current FuncCompiled is equal to the provided IObject. Always returns false.
func (o *FuncCompiled) Equals(_ IObject) bool {
	return false
}

// SourcePos retrieves the source position for the instruction pointer (ip) in the compiled function's source map.
// If the ip is not found in the map, it decrements ip until a valid position is found or returns NoPos if none exists.
func (o *FuncCompiled) SourcePos(ip int) int {
	for ip >= 0 {
		if p, ok := o.sourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return 0
}

// CanCall determines if the compiled function object can be invoked as a callable, always returning true.
func (o *FuncCompiled) CanCall() bool {
	return true
}
