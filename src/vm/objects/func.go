package objects

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/markel1974/c64emu/src/vm/opcodes"
)

const (
	FuncDef = "func"
)

// init registers the Func type with the gob package to allow for serialization and deserialization.
func init() {
	gob.Register(&Func{})
}

// Func represents a compiled function object with associated instructions, metadata, and runtime state.
// It includes memory allocation management, local variable tracking, parameter handling, and a source map.
type Func struct {
	IAllocator
	name          string
	instructions  *opcodes.Instructions
	numLocals     int
	numParameters int
	varArgs       bool
	source        map[int]int
	free          []*ObjectPointer
	freeIndices   []int
	async         bool
	serializer    []any
}

// newFunc creates and returns a new instance of Func using provided parameters and default initializations.
// factory handles object allocation and management while frame specifies the execution context frame.
// name defines the function name, and instructions specifies the bytecode to be executed.
// numLocals is the number of local variables, with numParameters defining the expected argument count.
// varArgs indicates if the function accepts a variable number of arguments.
// source maps instruction indices to source code, defaulting to an empty map if nil is provided.
// free holds externally closed variables as object pointers for use within the function's scope.
func newFunc(allocator IAllocator, name string, data []byte, numLocals int, numParameters int, varArgs bool, source map[int]int) IObject {
	if source == nil {
		source = make(map[int]int)
	}
	f := &Func{
		IAllocator:    allocator,
		name:          name,
		instructions:  opcodes.NewInstructions(data),
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		source:        source,
		free:          nil,
		async:         false,
	}
	f.serializer = []any{&f.name, &f.instructions, &f.numLocals, &f.numParameters, &f.varArgs, &f.source, &f.freeIndices, &f.async}
	return f
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Func) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Func) AsInterface() interface{} {
	return nil
}

// AsBool returns a boolean representation of the Func object, always returning false.
func (o *Func) AsBool() bool {
	return false
}

// AsInt64 converts the Func instance to an int64 representation, always returning 0.
func (o *Func) AsInt64() int64 {
	return 0
}

// AsFloat64 converts the Func instance to a float64 representation. Always returns 0.
func (o *Func) AsFloat64() float64 {
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Func) AsBytes() []byte {
	return nil
}

// AsString returns a string representation of the Func instance, specifically the constant FuncCompiledLabel.
func (o *Func) AsString() string {
	return ""
}

// AssignValue attempts to assign a Code to the instance but always returns ErrNotAssignable as it is not supported.
func (o *Func) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil returns false indicating the object is not nil.
func (o *Func) Nil() bool {
	return false
}

// LogicalOp performs a logical operation between the calling Func instance and the provided IObject operand.
// Returns the result of the logical operation or an error if the operator is invalid or unsupported.
func (o *Func) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation with the specified operator and IObject operand, returning the result or an error.
func (o *Func) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// UnaryOp performs a unary operation on the Func object and always returns ErrInvalidOperator.
func (o *Func) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false, indicating the object is considered falsy in boolean contexts.
func (o *Func) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve an element by index but always returns ErrIndexNotIndexable as the object is not indexable.
func (o *Func) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet assigns a specified Code to an index on the object but always returns ErrIndexUnsupported, indicating the operation is not supported.
func (o *Func) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns nil, as iteration is not supported by Func.
func (o *Func) Iterate(_ int) IIterator {
	return nil
}

// Iterable checks whether the Func instance is iterable and always returns false.
func (o *Func) Iterable() bool {
	return false
}

// Call invokes the compiled function with the given frame and arguments, returning the result count, object, and error if any.
func (o *Func) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Func object, typically representing the size of its associated instructions or data.
func (o *Func) Length() int {
	return 0
}

// Name returns the name of the compiled function as a string.
func (o *Func) Name() string {
	return o.name
}

// Code retrieves the compiled bytecode instructions associated with the Func object.
func (o *Func) Code() []byte {
	return o.instructions.Code()
}

// Instructions return a pointer to the Instructions associated with the Func instance.
func (o *Func) Instructions() *opcodes.Instructions {
	return o.instructions
}

// NumLocals returns the number of local variables required by the compiled function instance.
func (o *Func) NumLocals() int {
	return o.numLocals
}

// NumParameters returns the number of parameters required by the compiled function.
func (o *Func) NumParameters() int {
	return o.numParameters
}

// TypeName returns the string identifier "func_compiled" representing the type of the Func object.
func (o *Func) TypeName() string {
	return FuncDef
}

// Copy creates a deep copy of the Func object, replicating its properties and associated instructions.
func (o *Func) Copy(frame int, _ int) IObject {
	obj := o.GateKeeper().NewFunc(frame, o.name, o.instructions.Code(), o.numLocals, o.numParameters, o.varArgs, o.source)
	ret, ok := obj.(*Func)
	if !ok {
		return o.GateKeeper().UndefinedValue()
	}
	//ret.instructions = o.instructions.Copy()
	if len(o.freeIndices) > 0 {
		ret.freeIndices = append([]int{}, o.freeIndices...)
	}
	if len(o.free) > 0 {
		ret.free = append([]*ObjectPointer{}, o.free...)
	}
	return ret
}

// Equals determine whether the current Func instance is equal to the provided IObject. Returns false.
func (o *Func) Equals(_ IObject) bool {
	return false
}

// SourcePos returns the source code position corresponding to a given instruction pointer (ip) by consulting the source.
func (o *Func) SourcePos(ip int) int {
	for ip >= 0 {
		if p, ok := o.source[ip]; ok {
			return p
		}
		ip--
	}
	return 0
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Func) Count() int {
	return 1
}

// VarArgs returns true if the function is variadic, allowing it to accept a variable number of arguments.
func (o *Func) VarArgs() bool {
	return o.varArgs
}

// Async returns the async state of the Func, indicating whether it is set to operate asynchronously.
func (o *Func) Async() bool {
	return o.async
}

// SetAsync sets the async flag for the Func instance.
// Pass true to enable async operation or false to disable it.
func (o *Func) SetAsync(async bool) {
	o.async = async
}

// Free returns the slice of ObjectPointer instances that represent the free variables of the compiled function.
func (o *Func) Free() []*ObjectPointer {
	return o.free
}

// FreeIndices returns a slice of integers representing the indexes of free variables within the compiled function.
func (o *Func) FreeIndices() []int {
	return o.freeIndices
}

// FreeInitialize initializes free variables with provided indices and assigns default undefined values to their object pointers.
func (o *Func) FreeInitialize(freeIdx []int) {
	o.freeIndices = make([]int, len(freeIdx))
	copy(o.freeIndices, freeIdx)
}

// FreeSet sets the free variable pointers for a Func instance, ensuring the count matches the expected free variable indices.
func (o *Func) FreeSet(frameId int, required []IObject) error {
	if len(o.freeIndices) != len(required) {
		return fmt.Errorf("invalid free variable count: %d != %d", len(required), len(o.freeIndices))
	}
	o.free = make([]*ObjectPointer, len(required))
	for idx, obj := range required {
		freeObjPtr, err := o.GateKeeper().CreateObjectPointer(frameId, obj)
		if err != nil {
			return err
		}
		o.free[idx] = freeObjPtr
	}
	return nil
}

// GobEncode encodes the receiver into a byte slice using the Gob encoding scheme and returns the result or an error.
func (o *Func) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	for _, v := range o.serializer {
		if err := encoder.Encode(v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the receiver using the gob package.
func (o *Func) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	for _, v := range o.serializer {
		if err := decoder.Decode(v); err != nil {
			return err
		}
	}
	return nil
}
