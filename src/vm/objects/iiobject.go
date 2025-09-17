package objects

import (
	"time"
)

// IAllocator defines methods for memory management and object lifecycle, including reference counting and frame tracking.
type IAllocator interface {
	GateKeeper() IGateKeeper
	AddRef() int
	ReleaseRef() int
	RefCount() int
	setFrame(int)
	Frame() int
	SetStatic()
}

// IObject defines a polymorphic interface for managing various operations on objects, including conversions and evaluations.
// TypeName provides the name of the object's type as a string.
// AsBool converts and returns the object as a boolean.
// AsInt64 converts and returns the object as a 64-bit integer.
// AsFloat64 converts and returns the object as a 64-bit floating-point Code.
// AsString converts and returns the object as a string.
// AssignValue assigns the Code of another IObject to the implementing object instance.
// Nil checks if the object is nil or uninitialized.
// LogicalOp performs a logical operation with another IObject and returns the result.
// ArithmeticOp performs an arithmetic operation with another IObject and returns the result.
// Falsy evaluates the "truthiness" of the object and returns false if it is equivalently falsy.
// Equals determines whether another IObject is equal to the current instance.
// Copy creates a deep copy of the object, with specified frame and depth.
// IndexGet retrieves an element from the object using another IObject as an index.
// IndexSet assigns a Code to an index on the object using another IObject as the key.
// Iterate provides an iterator over the object if it supports iteration.
// Iterable checks whether the object can be iterated.
// Call invokes the object as a callable with the provided arguments.
// Length retrieves the len of the object, if applicable (e.g., arrays, strings).
type IObject interface {
	IAllocator
	setAllocator(allocator IAllocator)
	TypeName() string
	AsBool() bool
	AsInt64() int64
	AsFloat64() float64
	AsString() string
	AssignValue(object IObject) error
	Nil() bool
	LogicalOp(frame int, op LogicalOperator, rightHandSide IObject) (IObject, error)
	ArithmeticOp(frame int, op ArithmeticOperator, rightHandSide IObject) (IObject, error)
	Falsy() bool
	Equals(other IObject) bool
	Copy(frame int, depth int) IObject
	IndexGet(frame int, index IObject) (IObject, error)
	IndexSet(index IObject, value IObject) error
	Iterate(frame int) IIterator
	Iterable() bool
	Call(frame int, args ...IObject) (uint, IObject, error)
	Length() int
	Count() int
}

// IGateAllocator provides methods for creating and managing various IObject instances within a specific execution frame.
type IGateAllocator interface {
	Reset()
	FalseValue() IObject
	TrueValue() IObject
	Boolean(v bool) IObject
	UndefinedValue() IObject
	ReleaseObjects(int, []IObject)
	ReleaseAll()
	AllocatedObjects() uint64
	AssignAllocator(object IObject)
	SetPointer(ptr *ObjectPointer, value IObject)
	NewFuncInternal(frame int, id CallId) IObject
	NewFunc(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) IObject
	NewFuncImport(frame int, name string, args int, fn FuncCallable) IObject
	NewFuncJit(frame int, name string, data []byte) IObject
	NewArray(frame int, values []IObject) IObject
	NewBool(frame int, value bool) IObject
	NewBytes(frame int, value []byte) IObject
	NewChar(frame int, value rune) IObject
	NewError(frame int, e string) IObject
	NewFloat(frame int, v float64) IObject
	NewInt(frame int, v int64) IObject
	NewObjectPointer(frame int, value *IObject) IObject
	NewInterface(frame int, value IObject, iTable map[string]IObject) IObject
	NewMap(frame int, v map[string]IObject) IObject
	NewString(frame int, value string) IObject
	NewStruct(frame int, name string, value map[string]IObject) IObject
	NewAny(frame int, value interface{}) IObject
	NewTime(frame int, value time.Time) IObject
	NewMapIterator(frame int, v map[string]IObject, index int) IIterator
	NewStructIterator(frame int, v map[string]IObject, index int) IIterator
	NewStringIterator(frame int, v []rune, index int) IIterator
	NewArrayIterator(frame int, values []IObject, index int) IIterator
	NewBytesIterator(frame int, v []byte, index int) IIterator
}

// IGateConverter provides methods to convert IObject types to and from various native Go types and data structures.
type IGateConverter interface {
	ToBool(o IObject) (v bool, ok bool)
	ToInt64(o IObject) (int64, bool)
	ToFloat64(o IObject) (float64, bool)
	ToRune(o IObject) (v rune, ok bool)
	ToString(o IObject) (string, bool)
	ToTime(o IObject) (time.Time, bool)
	ToBytes(o IObject) ([]byte, bool)
	ToInterface(in IObject) (res interface{})
	ToMap(o IObject) (res map[string]interface{})

	ToBoolArg(index int, in []IObject) (bool, error)
	ToInt64Arg(index int, in []IObject) (int64, error)
	ToFloat64Arg(index int, in []IObject) (float64, error)
	ToTimeArg(index int, in []IObject) (time.Time, error)
	ToStringArg(index int, in []IObject) (string, error)
	ToBytesArg(index int, in []IObject) ([]byte, error)

	FromInterface(frame int, in interface{}) IObject
	FromMap(frame int, v map[string]interface{}) map[string]IObject
	FromBool(v bool) IObject
	FromStringArray(frame int, in []string) (IObject, error)

	AssignInt(val int64, dstObj IObject) error
	AssignBool(val bool, dstObj IObject) error
}

// IGateAdapter defines an interface for adapting various function signatures into FuncCallable instances.
type IGateAdapter interface {
	LogicalOpInt64(op LogicalOperator, lhs int64, rhs int64) (bool, error)
	ArithmeticOpInt64(op ArithmeticOperator, lhs int64, rhs int64) (int64, error)
	CreateSlice(frameId int, lowObj IObject, highObj IObject, target IObject) (IObject, error)
	IndexAssign(frame int, dst IObject, src IObject, selectors []IObject) error
}

// IGateKeeper combines IGateAllocator, IGateConverter, and IGateAdapter to manage object creation, conversion, and adaptation.
type IGateKeeper interface {
	IGateAllocator
	IGateConverter
	IGateAdapter
}
