package objects

import (
	"time"
)

// IObject represents a generic interface for any object in the system.
// TypeName returns the name of the type of the object.
// String returns a string representation of the object.
// Frame retrieves the current execution frame associated with the object.
// BinaryOp performs a binary operation using the specified operator and operands, returning the result or an error.
// Boolean evaluates and returns the boolean value of the object.
// Equals returns true if the object is equal to another given object.
// Copy creates and returns a deep copy of the object within the given frame.
// IndexGet retrieves the value at a given index from the object, returning the value or an error.
// IndexSet updates the value at a given index within the object, returning an error if the operation fails.
// Iterate returns an iterator for the object, enabling traversal over its elements.
// CanIterate checks whether the object supports iteration.
// Call invokes the object as a callable function with the provided arguments, returning the result or an error.
// CanCall checks whether the object supports being called as a function.
// Length returns the length of the object, if applicable.
type IObject interface {
	GateKeeper() IGateKeeper
	TypeName() string
	String() string
	Frame() int
	BinaryOp(frame int, op Operator, rightHandSide IObject) (IObject, error)
	Boolean() bool
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
	NewFuncCompiled(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) IObject
	NewFuncPackage(kind string, name string, fn FuncCallable) IObject
	NewFuncPackageFrame(frame int, kind string, name string, fn FuncCallable) IObject
	NewBuiltin(frame int, name string, index int) IObject
	NewArray(frame int, values []IObject) IObject
	NewArrayImmutable(frame int, values []IObject) IObject
	NewBool(frame int, value bool) IObject
	NewBytes(frame int, value []byte) IObject
	NewChar(frame int, value rune) IObject
	NewError(frame int, e string) IObject
	NewFloat(frame int, v float64) IObject
	NewInt(frame int, v int64) IObject
	NewObjectPointer(frame int, value *IObject) IObject
	NewMap(frame int, v map[string]IObject) IObject
	NewMapImmutable(frame int, v map[string]IObject) IObject
	NewString(frame int, value string) IObject
	NewStruct(frame int, value map[string]IObject) IObject
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

	BinaryOpInt64(op Operator, lhs int64, rhs int64) (int64, error)
	BoundsCheck(lowStack IObject, highStack IObject, numElements int64) (int64, int64, error)
	IndexAssign(frame int, dst IObject, src IObject, selectors []IObject) error
}

// IGateKeeper combines IGateAllocator, IGateConverter, and IGateAdapter to manage object creation, conversion, and adaptation.
type IGateKeeper interface {
	IGateAllocator
	IGateConverter
	IGateAdapter
}
