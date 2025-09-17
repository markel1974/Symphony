package objects

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)

// CallId represents a unique identifier for function calls or various operation types in a system.
type CallId int

// CallIdLen represents the operation identifier for determining the length.
// CallIdCopy represents the operation identifier for copying data.
// CallIdAppend represents the operation identifier for appending data.
// CallIdDelete represents the operation identifier for deleting data.
// CallIdSplice represents the operation identifier for splicing data.
// CallIdPanic represents the operation identifier for triggering a panic.
// CallIdRecover represents the operation identifier for recovering from a panic.
// CallIdInt represents the operation identifier for integer operations.
// CallIdBool represents the operation identifier for boolean operations.
// CallIdFloat represents the operation identifier for float operations.
// CallIdChar represents the operation identifier for character operations.
// CallIdString represents the operation identifier for string operations.
// CallIdTime represents the operation identifier for time-related operations.
// CallIdPrintf represents the operation identifier for formatted output using Printf.
// CallIdSprintf represents the operation identifier for formatted string output using Sprintf.
// CallIdMake represents the operation identifier for creating slices, maps, or channels.
const (
	CallIdLen = CallId(iota)
	CallIdCopy
	CallIdAppend
	CallIdDelete
	CallIdSplice
	CallIdPanic
	CallIdRecover
	CallIdInt
	CallIdBool
	CallIdFloat
	CallIdChar
	CallIdString
	CallIdTime
	CallIdPrintf
	CallIdSprintf
	CallIdMake
)

// _callIdContainer maps string keys to their corresponding CallId constants for identifying various operations or types.
var _callIdContainer = map[string]CallId{
	"len":     CallIdLen,
	"copy":    CallIdCopy,
	"append":  CallIdAppend,
	"delete":  CallIdDelete,
	"splice":  CallIdSplice,
	"panic":   CallIdPanic,
	"recover": CallIdRecover,
	"printf":  CallIdPrintf,
	"sprintf": CallIdSprintf,
	"make":    CallIdMake,
	"int":     CallIdInt, "int8": CallIdInt, "int32": CallIdInt, "int64": CallIdInt,
	"uint": CallIdInt, "uint8": CallIdInt, "uint32": CallIdInt, "uint64": CallIdInt,
	"float": CallIdFloat, "float32": CallIdFloat, "float64": CallIdFloat,
	"bool":   CallIdBool,
	"char":   CallIdChar,
	"byte":   CallIdChar,
	"string": CallIdString,
	"time":   CallIdTime,
}

// init initializes the package by registering the FuncInternal type with the gob encoder/decoder.
func init() {
	gob.Register(&FuncInternal{})
}

// CallIdFromString converts a string key to its associated CallId from a predefined container and returns its success status.
func CallIdFromString(in string) (CallId, bool) {
	v, ok := _callIdContainer[in]
	if !ok {
		return 0, false
	}
	return v, true
}

// FuncInternal is a type implementing IAllocator and storing a callable function with a unique identifier.
type FuncInternal struct {
	IAllocator
	id CallId
	fn func(frame int, args []IObject) (IObject, error)
}

// newFuncInternal initializes a new FuncInternal instance with the specified allocator and CallId, and prepares it.
// It returns the initialized IObject.
func newFuncInternal(allocator IAllocator, id CallId) IObject {
	fi := &FuncInternal{
		IAllocator: allocator,
		id:         id,
		fn:         nil,
	}
	fi.prepare(id)
	return fi
}

// setAllocator sets the allocator instance for the FuncInternal. It replaces the existing allocator with the provided one.
func (o *FuncInternal) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *FuncInternal) AsInterface() interface{} {
	return nil
}

// AsBool converts the internal state of the FuncInternal instance into a boolean representation and returns it.
func (o *FuncInternal) AsBool() bool {
	return false
}

// AsInt64 converts and returns the internal value of the receiver to an int64 type.
func (o *FuncInternal) AsInt64() int64 {
	return 0
}

// AsFloat64 converts the internal value of FuncInternal to a float64 representation and returns it.
func (o *FuncInternal) AsFloat64() float64 {
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *FuncInternal) AsBytes() []byte {
	return nil
}

// AsString returns the string representation of the FuncInternal instance.
func (o *FuncInternal) AsString() string {
	return "<FuncInternal>"
}

// AssignValue attempts to assign a value to the object and returns an error if assignment is not possible.
func (o *FuncInternal) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks whether the associated internal function is nil or not and returns false as a constant result.
func (o *FuncInternal) Nil() bool {
	return false
}

// LogicalOp performs a logical operation specified by the operator on the current object and the given right-hand side.
// Returns the result as an IObject or an error if the operation is invalid.
func (o *FuncInternal) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the provided operands and returns the resulting object or an error.
func (o *FuncInternal) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns a boolean value indicating false. It is a utility function that always evaluates to false.
func (o *FuncInternal) Falsy() bool {
	return false
}

// IndexGet retrieves an element based on an index from the object. Returns an error if the object is not indexable.
func (o *FuncInternal) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to set a value at a specified index but returns an error as indexing is not supported.
func (o *FuncInternal) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate is a method that returns an IIterator instance with the given integer input, intended for implementing iteration logic.
func (o *FuncInternal) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the current FuncInternal instance supports iteration and returns a boolean result.
func (o *FuncInternal) Iterable() bool {
	return false
}

// Length returns the integer representing the length or size as determined by the implementation.
func (o *FuncInternal) Length() int {
	return 0
}

// TypeName returns the type name of the FuncInternal instance as a string, including its identifier.
func (o *FuncInternal) TypeName() string {
	return "FuncInternal:" + strconv.Itoa(int(o.id))
}

// Copy creates a new instance of FuncInternal with the specified frame and retains the original id.
func (o *FuncInternal) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncInternal(frame, o.id)
}

// Equals determine if the provided IObject is equal to the current instance. Returns true if equal, otherwise false.
func (o *FuncInternal) Equals(_ IObject) bool {
	return false
}

// Call invokes the function with the given frame and arguments, returning a status code, the result, and an error if any.
func (o *FuncInternal) Call(frame int, args ...IObject) (uint, IObject, error) {
	v, err := o.fn(frame, args)
	return 1, v, err
}

// prepare initializes the function pointer `fn` based on the given CallId by mapping it to the appropriate method.
func (o *FuncInternal) prepare(id CallId) {
	switch id {
	case CallIdString:
		o.fn = o.string
	case CallIdInt:
		o.fn = o.int
	case CallIdFloat:
		o.fn = o.float
	case CallIdChar:
		o.fn = o.char
	case CallIdBool:
		o.fn = o.bool
	case CallIdTime:
		o.fn = o.time
	case CallIdLen:
		o.fn = o.len
	case CallIdPrintf:
		o.fn = o.printf
	case CallIdSprintf:
		o.fn = o.sprintf
	case CallIdCopy:
		o.fn = o.copy
	case CallIdAppend:
		o.fn = o.append
	case CallIdSplice:
		o.fn = o.splice
	case CallIdRecover:
		o.fn = o.recover
	case CallIdPanic:
		o.fn = o.panic
	case CallIdDelete:
		o.fn = o.delete
	case CallIdMake:
		o.fn = o.make
	default:
		o.fn = o.undefined
	}
}

// undefined returns an IObject representing an undefined value, leveraging the GateKeeper's UndefinedValue method.
func (o *FuncInternal) undefined(_ int, _ []IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), nil
}

// string processes a single argument, converts it to a string, and returns a new string object or an error.
func (o *FuncInternal) string(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return o.GateKeeper().NewString(frame, args[0].AsString()), nil
}

// int handles integer conversion and validation of input arguments, returning an integer object or an error if invalid.
func (o *FuncInternal) int(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return o.GateKeeper().NewInt(frame, args[0].AsInt64()), nil
}

// float converts the provided argument to a floating-point object or returns an error if the conversion fails.
func (o *FuncInternal) float(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return o.GateKeeper().NewFloat(frame, args[0].AsFloat64()), nil
}

// char converts an integer argument to a character and returns it as a new IObject.
func (o *FuncInternal) char(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return o.GateKeeper().NewChar(frame, rune(args[0].AsInt64())), nil
}

// bool evaluates the given argument and converts it to a boolean, returning the result as an IObject or an error.
func (o *FuncInternal) bool(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return o.GateKeeper().Boolean(args[0].AsBool()), nil
}

// time processes a frame and arguments to return a new time object or an error if arguments are invalid.
func (o *FuncInternal) time(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	if v, ok := args[0].(*Time); ok {
		return o.GateKeeper().NewTime(frame, v.Value()), nil
	}
	return o.GateKeeper().NewTime(frame, time.Unix(args[0].AsInt64(), 0)), nil
}

// copy creates a duplicate of the provided object, ensuring only one argument is passed for the operation.
func (o *FuncInternal) copy(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return args[0].Copy(frame, 0), nil
}

// len computes the length of the given argument and returns it as an integer object. Returns an error on failure.
func (o *FuncInternal) len(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return o.GateKeeper().NewInt(frame, int64(args[0].Length())), nil
}

// sprintf formats a string according to a format specifier and arguments, returning the formatted string or an error.
func (o *FuncInternal) sprintf(frame int, args []IObject) (IObject, error) {
	switch len(args) {
	case 0:
		return o.GateKeeper().UndefinedValue(), ErrInvalidArgumentsNumber
	case 1:
		return args[0], nil
	default:
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, v.AsInterface())
		}
		return o.GateKeeper().NewString(frame, fmt.Sprintf(args[0].AsString(), ar...)), nil
	}
}

// printf formats and prints the provided arguments according to the specified format.
// The first argument is treated as a format string, and the subsequent arguments are its values.
// Returns an UndefinedValue and error if the number of arguments is invalid.
func (o *FuncInternal) printf(_ int, args []IObject) (IObject, error) {
	switch len(args) {
	case 0:
		return o.GateKeeper().UndefinedValue(), ErrInvalidArgumentsNumber
	case 1:
		fmt.Printf(args[0].AsString())
		return o.GateKeeper().UndefinedValue(), nil
	default:
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, v.AsInterface())
		}
		fmt.Printf(args[0].AsString(), ar...)
		return o.GateKeeper().UndefinedValue(), nil
	}
}

// append adds elements to an array or creates a new array with the provided arguments.
// It requires at least two arguments: the target array and the elements to append.
// Returns the modified or new array and an error if invalid arguments are provided.
func (o *FuncInternal) append(frame int, args []IObject) (IObject, error) {
	if len(args) < 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Array:
		return o.GateKeeper().NewArray(frame, append(arg.Values(), args[1:]...)), nil
	default:
		return o.GateKeeper().NewArray(frame, args[1:]), nil
	}
}

// delete removes a key-value pair from a Map object based on the provided key and returns an undefined value or an error.
func (o *FuncInternal) delete(_ int, args []IObject) (IObject, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Map:
		arg.Delete(args[1].AsString())
		return o.GateKeeper().UndefinedValue(), nil
	default:
		return nil, NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// splice removes or replaces elements in an array, starting at a specified index, and optionally inserts new elements.
func (o *FuncInternal) splice(frame int, args []IObject) (IObject, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, ErrInvalidArgumentsNumber
	}
	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	arrayLen := array.Length()
	var startIdx int
	if argsLen > 1 {
		startIdx = int(args[1].AsInt64())
		if startIdx < 0 || startIdx > arrayLen {
			return nil, ErrIndexOutOfBounds
		}
	}
	delCount := array.Length()
	if argsLen > 2 {
		delCount = int(args[2].AsInt64())
		if delCount < 0 {
			return nil, ErrIndexOutOfBounds
		}
	}
	if startIdx+delCount > arrayLen {
		delCount = arrayLen - startIdx
	}
	endIdx := startIdx + delCount
	deleted := append([]IObject{}, array.Values()[startIdx:endIdx]...)
	head := array.Values()[:startIdx]
	var items []IObject
	if argsLen > 3 {
		items = make([]IObject, 0, argsLen-3)
		for i := 3; i < argsLen; i++ {
			items = append(items, args[i])
		}
	}
	items = append(items, array.Values()[endIdx:]...)
	array.Assign(append(head, items...))
	return o.GateKeeper().NewArray(frame, deleted), nil
}

// panic triggers a runtime panic by returning an error with the message from the first argument or a default "panic" message.
func (o *FuncInternal) panic(_ int, args []IObject) (IObject, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("%s", args[0].AsString())
	}
	return nil, fmt.Errorf("panic")
}

// recover attempts to gracefully handle potential panics or errors within the function and ensures a safe fallback response.
func (o *FuncInternal) recover(_ int, _ []IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), nil
}

// make creates and initializes objects such as slices, maps, or arrays based on the specified type and dimensions.
func (o *FuncInternal) make(frame int, args []IObject) (IObject, error) {
	var kind IObject
	count := 0
	reserve := 0
	switch len(args) {
	case 1:
		kind = args[0]
	case 2:
		kind = args[0]
		if count = int(args[1].AsInt64()); count < 0 {
			return nil, fmt.Errorf("make: len out of range: %d", count)
		}
		reserve = count
	case 3:
		kind = args[0]
		if count = int(args[1].AsInt64()); count < 0 {
			return nil, fmt.Errorf("make: len out of range: %d", count)
		}
		if reserve = int(args[2].AsInt64()); reserve < count {
			return nil, fmt.Errorf("make: cap out of range: %d", reserve)
		}
	default:
		return nil, ErrInvalidArgumentsNumber
	}
	switch kind.TypeName() {
	case BytesType:
		return o.GateKeeper().NewBytes(frame, make([]byte, count, reserve)), nil
	case ArrayType:
		return o.GateKeeper().NewArray(frame, make([]IObject, count, reserve)), nil
	case MapType:
		return o.GateKeeper().NewMap(frame, make(map[string]IObject, count)), nil
	default:
		return nil, fmt.Errorf("cannot make type %s", kind.TypeName())
	}
}

// Count returns an integer value typically representing the count or total computed by this method.
func (o *FuncInternal) Count() int {
	return 1
}

// GobEncode serializes the receiver into a byte slice using the gob encoding format and returns the encoded data and error.
func (o *FuncInternal) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.id); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.fn); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the FuncInternal fields using the gob package.
func (o *FuncInternal) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.id); err != nil {
		return err
	}
	if err := decoder.Decode(&o.fn); err != nil {
		return err
	}
	return nil
}
