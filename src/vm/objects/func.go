package objects

import (
	"encoding/gob"

	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// FuncCompiledDef represents a constant value for a compiled function definition.
// FuncCompiledLabel represents a formatted label using FuncCompiledDef.
const (
	FuncCompiledDef   = "func_compiled"
	FuncCompiledLabel = "<" + FuncCompiledDef + ">"
)

// init registers the Func type with the gob package to allow for serialization and deserialization.
func init() {
	gob.Register(&Func{})
}

// Func represents a compiled function object with associated instructions, metadata, and runtime state.
// It includes memory allocation management, local variable tracking, parameter handling, and a source map.
type Func struct {
	Allocator
	name          string
	instructions  *opcodes.Instructions
	numLocals     int
	numParameters int
	varArgs       bool
	source        map[int]int
	free          []*ObjectPointer
}

// newFunc creates and returns a new instance of Func using provided parameters and default initializations.
// factory handles object allocation and management while frame specifies the execution context frame.
// name defines the function name, and instructions specifies the bytecode to be executed.
// numLocals is the number of local variables, with numParameters defining the expected argument count.
// varArgs indicates if the function accepts a variable number of arguments.
// source maps instruction indices to source code, defaulting to an empty map if nil is provided.
// free holds externally closed variables as object pointers for use within the function's scope.
func newFunc(factory IGateKeeper, frame int, name string, data []byte, numLocals int, numParameters int, varArgs bool, source map[int]int, free []*ObjectPointer) IObject {
	if source == nil {
		source = make(map[int]int)
	}
	return &Func{
		Allocator:     Allocator{gk: factory, frame: frame},
		name:          name,
		instructions:  opcodes.NewInstructions(data),
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		source:        source,
		free:          free,
	}
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

// AsString returns a string representation of the Func instance, specifically the constant FuncCompiledLabel.
func (o *Func) AsString() string {
	return FuncCompiledLabel
}

// AssignValue attempts to assign a value to the instance but always returns ErrNotAssignable as it is not supported.
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
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation with the specified operator and IObject operand, returning the result or an error.
func (o *Func) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false, indicating the object is considered falsy in boolean contexts.
func (o *Func) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve an element by index but always returns ErrIndexNotIndexable as the object is not indexable.
func (o *Func) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.gk.UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet assigns a specified value to an index on the object but always returns ErrIndexUnsupported, indicating the operation is not supported.
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

// Length returns the length of the Func object, typically representing the size of its associated instructions or data.
func (o *Func) Length() int {
	return 0
}

// Name returns the name of the compiled function as a string.
func (o *Func) Name() string {
	return o.name
}

// Data retrieves the compiled bytecode instructions associated with the Func object.
func (o *Func) Data() []byte {
	return o.instructions.Data()
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

// VarArgs returns true if the function is variadic, allowing it to accept a variable number of arguments.
func (o *Func) VarArgs() bool {
	return o.varArgs
}

// Free returns the slice of ObjectPointer instances that represent the free variables of the compiled function.
func (o *Func) Free() []*ObjectPointer {
	return o.free
}

// TypeName returns the string identifier "func_compiled" representing the type of the Func object.
func (o *Func) TypeName() string {
	return FuncCompiledDef
}

// Copy creates a deep copy of the Func object, replicating its properties and associated instructions.
func (o *Func) Copy(frame int, _ int) IObject {
	obj := o.GateKeeper().NewFunc(frame, o.name, nil, o.numLocals, o.numParameters, o.varArgs, nil, nil)
	ret, ok := obj.(*Func)
	if !ok {
		return o.GateKeeper().UndefinedValue()
	}
	ret.instructions = o.instructions.Copy()
	ret.free = append([]*ObjectPointer{}, o.free...)
	return ret
}

// Equals determines whether the current Func instance is equal to the provided IObject. Returns false.
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
