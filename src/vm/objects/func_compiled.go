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
func (i *Instructions) Get8(base uint) (uint8, error) {
	if base >= uint(len(i.data)) {
		return 0, fmt.Errorf("invalid instruction index: %d", base)
	}
	return i.data[base], nil
}

// Get16 retrieves a 16-bit unsigned integer composed of two bytes from the instructions at indices low and low-1.
// Returns an error if the provided indices are out of bounds.
func (i *Instructions) Get16(base uint) (uint16, error) {
	high := base - 1
	if base >= uint(len(i.data)) {
		return 0, fmt.Errorf("invalid instruction low index: %d", base)
	}
	if high >= uint(len(i.data)) {
		return 0, fmt.Errorf("invalid instruction high index: %d", high)
	}
	return uint16(i.data[base]) | uint16(i.data[high])<<8, nil
}

// Get32 retrieves a 32-bit unsigned integer composed of four bytes from the instructions at base and the preceding three indices.
// The base index must be greater than or equal to 3 and within the bounds of the instructions' data.
// Returns an error if the base index is out of bounds.
func (i *Instructions) Get32(base uint) (uint32, error) {
	if base < 3 || base >= uint(len(i.data)) {
		return 0, fmt.Errorf("invalid instruction index for 32-bit read: base %d is out of bounds", base)
	}
	byte1 := i.data[base-3] // MSB (Most Significant Byte)
	byte2 := i.data[base-2]
	byte3 := i.data[base-1]
	byte4 := i.data[base] // LSB (Least Significant Byte)
	return uint32(byte4) | uint32(byte3)<<8 | uint32(byte2)<<16 | uint32(byte1)<<24, nil
}

// Get64 retrieves a 64-bit unsigned integer from the instructions at base and the preceding seven indices.
// The base index must be greater than or equal to 7 and within the bounds of the instructions' data.
// Returns an error if the base index is out of bounds.
func (i *Instructions) Get64(base uint) (uint64, error) {
	if base < 7 || base >= uint(len(i.data)) {
		return 0, fmt.Errorf("invalid instruction index for 64-bit read: base %d is out of bounds", base)
	}
	byte1 := i.data[base-7] // MSB (Most Significant Byte)
	byte2 := i.data[base-6]
	byte3 := i.data[base-5]
	byte4 := i.data[base-4]
	byte5 := i.data[base-3]
	byte6 := i.data[base-2]
	byte7 := i.data[base-1]
	byte8 := i.data[base] // LSB (Least Significant Byte)

	// Combina i byte in un valore uint64
	return uint64(byte8) | uint64(byte7)<<8 | uint64(byte6)<<16 | uint64(byte5)<<24 |
		uint64(byte4)<<32 | uint64(byte3)<<40 | uint64(byte2)<<48 | uint64(byte1)<<56, nil
}

// FuncCompiled represents a compiled function with bytecode instructions, metadata, and associated free variables.
type FuncCompiled struct {
	Allocator
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
		Allocator:     Allocator{gk: factory, frame: frame},
		name:          name,
		instructions:  NewInstructions(instructions),
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		sourceMap:     sourceMap,
		free:          free,
	}
}

// AsBool returns a boolean representation of the FuncCompiled object, always returning false.
func (o *FuncCompiled) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (o *FuncCompiled) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *FuncCompiled) AsFloat64() float64 {
	return 0
}

// AsString returns a human-readable string representation of the FuncCompiled object.
func (o *FuncCompiled) AsString() string {
	return FuncCompiledLabel
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *FuncCompiled) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *FuncCompiled) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the object using the specified operator and right-hand-side operand.
// Returns an error if the operator is invalid.
func (o *FuncCompiled) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
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

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *FuncCompiled) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *FuncCompiled) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
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
