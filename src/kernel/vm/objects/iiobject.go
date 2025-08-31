package objects

import (
	"time"
)

// IObject defines an interface for handling various object types in a flexible and extensible manner.
// GateKeeper returns the IGateKeeper instance associated with the object for managing object operations.
// TypeName returns the name of the object's type as a string.
// AsBool converts the object to a boolean value, if possible.
// AsInt64 converts the object to a 64-bit integer, if possible.
// AsFloat64 converts the object to a 64-bit floating-point number, if possible.
// AsString converts the object to a string value, if possible.
// Assign assigns the current object with attributes of another IObject, returning an error if the operation fails.
// Frame retrieves the execution frame ID associated with the object.
// LogicalOp performs a logical operation between the object and another IObject, returning a new result or an error.
// ArithmeticOp performs an arithmetic operation between the object and another IObject, returning a new result or an error.
// Falsy checks if the object evaluates to a false-like value.
// Equals determines if the current object is equal to another IObject.
// Copy creates a copy of the current object within the specified execution frame and recursion depth.
// IndexGet retrieves a value from the object at the specified index and returns it, or an error if unavailable.
// IndexSet updates the object at the specified index with a new value, potentially returning an error.
// Iterate returns an iterator for traversing the object if it's iterable.
// CanIterate checks if the object supports iteration.
// Call invokes the object as if it's a callable function, passing arguments and returning the result or an error.
// CanCall checks if the object can be invoked as a callable entity.
// Length returns the length of the object, typically for iterable or sized entities.
type IObject interface {
	GateKeeper() IGateKeeper
	TypeName() string
	AsBool() bool
	AsInt64() int64
	AsFloat64() float64
	AsString() string
	AssignValue(object IObject) error
	Frame() int
	LogicalOp(frame int, op LogicalOperator, rightHandSide IObject) (IObject, error)
	ArithmeticOp(frame int, op ArithmeticOperator, rightHandSide IObject) (IObject, error)
	Falsy() bool
	Equals(other IObject) bool
	Copy(frame int, depth int) IObject
	IndexGet(frame int, index IObject) (value IObject, err error)
	IndexSet(index, value IObject) error
	Iterate(frame int) IIterator
	CanIterate() bool
	Call(frame int, args ...IObject) (ret IObject, err error)
	CanCall() bool
	Length() int
}

// IGateAllocator provides methods for creating and managing various IObject instances within a specific execution frame.
type IGateAllocator interface {
	Reset()
	FalseValue() IObject
	TrueValue() IObject
	UndefinedValue() IObject
	ReleaseObject(IObject)
	ReleaseObjects([]IObject)
	NewFuncCompiled(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) IObject
	NewFuncPackage(kind string, name string, fn FuncCallable) IObject
	NewFuncPackageFrame(frame int, kind string, name string, fn FuncCallable) IObject
	NewBuiltin(frame int, name string, index int) IObject
	NewFuncJit(kind string, name string, data []byte) IObject
	NewFuncJitFrame(frame int, kind string, name string, data []byte) IObject
	NewArray(frame int, values []IObject) IObject
	NewBool(frame int, value bool) IObject
	NewBytes(frame int, value []byte) IObject
	NewChar(frame int, value rune) IObject
	NewError(frame int, e string) IObject
	NewFloat(frame int, v float64) IObject
	NewInt(frame int, v int64) IObject
	NewObjectPointer(frame int, value *IObject) IObject
	NewMap(frame int, v map[string]IObject) IObject
	NewString(frame int, value string) IObject
	NewStruct(frame int, value map[string]IObject) IObject
	NewTime(frame int, value time.Time) IObject
	NewMapIterator(frame int, v map[string]IObject, index int) IIterator
	NewStructIterator(frame int, v map[string]IObject, index int) IIterator
	NewStringIterator(frame int, v []rune, index int) IIterator
	NewArrayIterator(frame int, values []IObject, index int) IIterator
	NewBytesIterator(frame int, v []byte, index int) IIterator
	NewInterface(frame int, value IObject, iTable map[string]IObject) IObject
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

	ToBoolArg(index int, o IObject) (bool, error)
	ToInt64Arg(index int, o IObject) (int64, error)
	ToFloat64Arg(index int, o IObject) (float64, error)
	ToTimeArg(index int, o IObject) (time.Time, error)
	ToStringArg(index int, o IObject) (string, error)
	ToStringArrayArg(index int, arr []IObject) ([]string, error)
	ToBytesArg(index int, o IObject) ([]byte, error)

	FromInterface(frame int, in interface{}) IObject
	FromMap(frame int, v map[string]interface{}) map[string]IObject
	FromBool(v bool) IObject
	FromStringArray(frame int, in []string) (IObject, error)
}

// IGateAdapter defines an interface for adapting various function signatures into FuncCallable instances.
type IGateAdapter interface {
	FuncIsOs(fn func(string) string) FuncCallable
	FuncIsOsS(fn func(string) []string) FuncCallable
	FuncIsOse(fn func(string) (string, error)) FuncCallable
	FuncIsOie(fn func(string) (int, error)) FuncCallable
	FuncIsObSe(fn func(string) ([]byte, error)) FuncCallable
	FuncIbSOs(fn func([]byte) string) FuncCallable
	FuncIssOsS(fn func(string, string) []string) FuncCallable
	FuncIssiOsS(fn func(string, string, int) []string) FuncCallable
	FuncIssOi(fn func(string, string) int) FuncCallable
	FuncIssOs(fn func(string, string) string) FuncCallable
	FuncIssOb(fn func(string, string) bool) FuncCallable
	FuncIiOs(fn func(int) string) FuncCallable
	FuncIiOiS(fn func(int) []int) FuncCallable

	FuncIiOf64(fn func(int) float64) FuncCallable
	FuncIi64On(fn func(int64)) FuncCallable
	FuncIf64Of64(fn func(float64) float64) FuncCallable
	FuncIf64Oi(fn func(float64) int) FuncCallable
	FuncIf64Ob(fn func(float64) bool) FuncCallable
	FuncIf64f64Of64(fn func(float64, float64) float64) FuncCallable
	FuncIif64Of64(fn func(int, float64) float64) FuncCallable
	FuncIf64iOf64(fn func(float64, int) float64) FuncCallable
	FuncIf64iOb(fn func(float64, int) bool) FuncCallable

	LogicalOpInt64(op LogicalOperator, lhs int64, rhs int64) (bool, error)
	ArithmeticOpInt64(op ArithmeticOperator, lhs int64, rhs int64) (int64, error)
	BoundsCheck(lowStack IObject, highStack IObject, numElements int64) (int64, int64, error)
	IndexAssign(frame int, dst IObject, src IObject, selectors []IObject) error
}

// IGateKeeper combines IGateAllocator, IGateConverter, and IGateAdapter to manage object creation, conversion, and adaptation.
type IGateKeeper interface {
	IGateAllocator
	IGateConverter
	IGateAdapter
}
