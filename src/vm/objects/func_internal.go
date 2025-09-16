package objects

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)

// CallId is a custom integer type used to represent unique identifiers for specific function calls or operations.
type CallId int

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
	CallIdTypeName
	CallIdIsInt
	CallIdIsFloat
	CallIdIsString
	CallIdIsBool
	CallIdIsChar
	CallIdIsBytes
	CallIdIsArray
	CallIdIsMap
	CallIdIsIterable
	CallIdIsTime
	CallIdIsError
	CallIdIsUndefined
	CallIdIsFunction
	CallIdIsCallable
	CallIdPrintf
	CallIdSprintf
	CallIdMake
)

var _callIdContainer = map[string]CallId{
	"len":         CallIdLen,
	"copy":        CallIdCopy,
	"append":      CallIdAppend,
	"delete":      CallIdDelete,
	"splice":      CallIdSplice,
	"panic":       CallIdPanic,
	"recover":     CallIdRecover,
	"int":         CallIdInt,
	"int8":        CallIdInt,
	"int32":       CallIdInt,
	"int64":       CallIdInt,
	"uint8":       CallIdInt,
	"uint32":      CallIdInt,
	"uint64":      CallIdInt,
	"bool":        CallIdBool,
	"float":       CallIdFloat,
	"float32":     CallIdFloat,
	"float64":     CallIdFloat,
	"char":        CallIdChar,
	"byte":        CallIdChar,
	"string":      CallIdString,
	"time":        CallIdTime,
	"name":        CallIdTypeName,
	"isInt":       CallIdIsInt,
	"isFloat":     CallIdIsFloat,
	"isString":    CallIdIsString,
	"isBool":      CallIdIsBool,
	"isChar":      CallIdIsChar,
	"isBytes":     CallIdIsBytes,
	"isArray":     CallIdIsArray,
	"isMap":       CallIdIsMap,
	"isIterable":  CallIdIsIterable,
	"isTime":      CallIdIsTime,
	"isError":     CallIdIsError,
	"isUndefined": CallIdIsUndefined,
	"isFunction":  CallIdIsFunction,
	"isCallable":  CallIdIsCallable,
	"printf":      CallIdPrintf,
	"sprintf":     CallIdSprintf,
	"make":        CallIdMake,
}

func init() {
	gob.Register(&FuncInternal{})
}

func CallIdFromString(in string) (CallId, bool) {
	v, ok := _callIdContainer[in]
	if !ok {
		return 0, false
	}
	return v, true
}

// FuncInternal is a callable object type that encapsulates a function and provides execution context information.
type FuncInternal struct {
	IAllocator
	id CallId
	fn func(frame int, args []IObject) (IObject, error)
}

// NewFuncImport creates a new FuncImport instance with the specified Id and callable function.
func newFuncInternal(allocator IAllocator, id CallId) IObject {
	fi := &FuncInternal{
		IAllocator: allocator,
		id:         id,
		fn:         nil,
	}
	fi.prepare(id)
	return fi
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (i *FuncInternal) setAllocator(allocator IAllocator) {
	i.IAllocator = allocator
}

// AsBool returns a boolean representation of the object, always returning false for FuncJit.
func (i *FuncInternal) AsBool() bool {
	return false
}

// AsInt64 returns the len of the array as an int64 Code.
func (i *FuncInternal) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the len of the array as an int64 Code.
func (i *FuncInternal) AsFloat64() float64 {
	return 0
}

// AsString returns the string representation of a FuncImport object.
func (i *FuncInternal) AsString() string {
	return "<FuncInternal>"
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (i *FuncInternal) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (i *FuncInternal) Nil() bool {
	return false
}

// LogicalOp performs a logical operation between the current object and a provided IObject using the specified operator.
// Always returns nil and ErrInvalidOperator as logical operations are not supported for this object type.
func (i *FuncInternal) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(i.GateKeeper(), op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and operand, returning the result or an error.
func (i *FuncInternal) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (i *FuncInternal) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (i *FuncInternal) IndexGet(_ int, _ IObject) (IObject, error) {
	return i.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (i *FuncInternal) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (i *FuncInternal) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (i *FuncInternal) Iterable() bool {
	return false
}

// Length returns the len of the Int object.
func (i *FuncInternal) Length() int {
	return 0
}

// TypeName returns the type name of the FuncImport as a string.
func (i *FuncInternal) TypeName() string {
	return "FuncInternal:" + strconv.Itoa(int(i.id))
}

// Copy creates and returns a new FuncImport instance with the same Value field as the original object.
func (i *FuncInternal) Copy(frame int, _ int) IObject {
	return i.GateKeeper().NewFuncInternal(frame, i.id)
}

// Equals checks whether the current FuncImport is equal to another object of type IObject. Always returns false.
func (i *FuncInternal) Equals(_ IObject) bool {
	return false
}

// Call executes the GateCall with the provided GateKeeper, frame, and arguments, returning an IObject or an error.
func (i *FuncInternal) Call(frame int, args ...IObject) (uint, IObject, error) {
	v, err := i.fn(frame, args)
	return 1, v, err
}

// prepare initializes the function pointer `fn` in GateCall based on the provided CallId.
// It maps CallId enums to specific methods or operations within the GateCall struct.
// Returns an error if an undefined CallId is provided.
func (i *FuncInternal) prepare(id CallId) {
	switch id {
	case CallIdIsString:
		i.fn = i.isString
	case CallIdIsInt:
		i.fn = i.isInt
	case CallIdIsFloat:
		i.fn = i.isFloat
	case CallIdIsBool:
		i.fn = i.isBool
	case CallIdIsChar:
		i.fn = i.isChar
	case CallIdIsBytes:
		i.fn = i.isBytes
	case CallIdIsArray:
		i.fn = i.isArray
	case CallIdIsMap:
		i.fn = i.isMap
	case CallIdIsTime:
		i.fn = i.isTime
	case CallIdIsError:
		i.fn = i.isError
	case CallIdIsUndefined:
		i.fn = i.isUndefined
	case CallIdIsFunction:
		i.fn = i.isFunction
	case CallIdIsIterable:
		i.fn = i.isIterable
	case CallIdTypeName:
		i.fn = i.callIdTypeName
	case CallIdString:
		i.fn = i.callIdString
	case CallIdInt:
		i.fn = i.int
	case CallIdFloat:
		i.fn = i.float
	case CallIdChar:
		i.fn = i.char
	case CallIdBool:
		i.fn = i.bool
	case CallIdTime:
		i.fn = i.time
	case CallIdLen:
		i.fn = i.len
	case CallIdPrintf:
		i.fn = i.printf
	case CallIdSprintf:
		i.fn = i.sprintf
	case CallIdCopy:
		i.fn = i.copy
	case CallIdAppend:
		i.fn = i.append
	case CallIdSplice:
		i.fn = i.splice
	case CallIdRecover:
		i.fn = i.recover
	case CallIdPanic:
		i.fn = i.panic
	case CallIdDelete:
		i.fn = i.delete
	case CallIdMake:
		i.fn = i.make
	default:
		i.fn = i.undefined
	}
}

// undefined is a placeholder function for handling undefined arguments.
func (i *FuncInternal) undefined(_ int, _ []IObject) (IObject, error) {
	return i.GateKeeper().UndefinedValue(), nil
}

// isString verifies if the first argument is of type String and returns a Boolean indicating the result.
func (i *FuncInternal) isString(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*String)
	return i.GateKeeper().Boolean(ok), nil
}

// isInt checks if the first argument in the Args slice is of type Int and returns a boolean IObject indicating the result.
func (i *FuncInternal) isInt(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Int)
	return i.GateKeeper().Boolean(ok), nil
}

// isFloat checks if the first argument is of type Float and returns a Boolean result.
func (i *FuncInternal) isFloat(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Float)
	return i.GateKeeper().Boolean(ok), nil
}

// isBool checks if the first element in the Args slice is of type Bool and returns a Boolean object accordingly.
func (i *FuncInternal) isBool(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Bool)
	return i.GateKeeper().Boolean(ok), nil
}

// isChar checks if the first argument is of type *Char and returns a boolean result wrapped in an IObject.
func (i *FuncInternal) isChar(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Char)
	return i.GateKeeper().Boolean(ok), nil
}

// isBytes checks if the first argument is of type Bytes and returns a boolean result wrapped in IObject.
func (i *FuncInternal) isBytes(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Bytes)
	return i.GateKeeper().Boolean(ok), nil
}

// isArray checks if the first argument is of type Array and returns a boolean IObject accordingly.
func (i *FuncInternal) isArray(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Array)
	return i.GateKeeper().Boolean(ok), nil
}

// isMap checks if the first argument is of type *Map and returns a boolean IObject and an error if applicable.
func (i *FuncInternal) isMap(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Map)
	return i.GateKeeper().Boolean(ok), nil
}

// isTime verifies if the provided argument is of type Time and returns a boolean result or an error.
func (i *FuncInternal) isTime(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Time)
	return i.GateKeeper().Boolean(ok), nil
}

// isError determines if the given argument is of type *Error and returns a boolean result with any error encountered.
func (i *FuncInternal) isError(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Error)
	return i.GateKeeper().Boolean(ok), nil
}

// isUndefined checks if the first argument is of type Undefined and returns a boolean result or an error if invalid input.
func (i *FuncInternal) isUndefined(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Undefined)
	return i.GateKeeper().Boolean(ok), nil
}

// isFunction checks if the first argument is of type *Func and returns a boolean IObject accordingly.
func (i *FuncInternal) isFunction(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Func)
	return i.GateKeeper().Boolean(ok), nil
}

// isIterable determines if the argument is iterable and returns a boolean representation of the result.
func (i *FuncInternal) isIterable(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().Boolean(args[0].Iterable()), nil
}

// callIdTypeName generates a string representation of the type name for the provided argument.
func (i *FuncInternal) callIdTypeName(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().NewString(frame, args[0].TypeName()), nil
}

// callIdString generates a new string object using the provided frame and argument, ensuring exactly one argument is passed.
func (i *FuncInternal) callIdString(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().NewString(frame, args[0].AsString()), nil
}

// callIdInt generates a new integer object using the provided frame and the integer Code extracted from the input argument.
func (i *FuncInternal) int(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().NewInt(frame, args[0].AsInt64()), nil
}

// callIdFloat converts the first argument to a float and creates a new float object for the given frame.
func (i *FuncInternal) float(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().NewFloat(frame, args[0].AsFloat64()), nil
}

// callIdChar converts an integer argument to a character and returns the new character object, or an error if input is invalid.
func (i *FuncInternal) char(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().NewChar(frame, rune(args[0].AsInt64())), nil
}

// callIdBool processes a single argument and returns a Boolean object based on the argument's boolean Code or an error.
func (i *FuncInternal) bool(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().Boolean(args[0].AsBool()), nil
}

// callIdTime processes a method call to generate a new Time object based on the provided frame and argument.
// It expects exactly one argument, either a Time object or a Code convertible to a timestamp.
// Returns the newly created Time object or an error if the input is invalid.
func (i *FuncInternal) time(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	if v, ok := args[0].(*Time); ok {
		return i.GateKeeper().NewTime(frame, v.Value()), nil
	}
	return i.GateKeeper().NewTime(frame, time.Unix(args[0].AsInt64(), 0)), nil
}

// copy creates a copy of the given IObject argument and returns it, ensuring only a single argument is provided.
func (i *FuncInternal) copy(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return args[0].Copy(frame, 0), nil
}

// len computes the len of the provided container argument (Array, String, Bytes, or Map) and returns it as an integer.
func (i *FuncInternal) len(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return i.GateKeeper().NewInt(frame, int64(args[0].Length())), nil
}

// sprintf formats a string using the format specified in Args[0] and additional arguments provided.
// Returns a new formatted string or an error if the argument count is invalid.
// If Args contains only one element, it directly returns that element.
func (i *FuncInternal) sprintf(frame int, args []IObject) (IObject, error) {
	switch len(args) {
	case 0:
		return i.GateKeeper().UndefinedValue(), ErrInvalidArgumentsNumber
	case 1:
		return args[0], nil
	default:
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, i.GateKeeper().ToInterface(v))
		}
		return i.GateKeeper().NewString(frame, fmt.Sprintf(args[0].AsString(), ar...)), nil
	}
}

// printf formats and prints the given arguments based on the format string provided in the first argument.
func (i *FuncInternal) printf(_ int, args []IObject) (IObject, error) {
	switch len(args) {
	case 0:
		return i.GateKeeper().UndefinedValue(), ErrInvalidArgumentsNumber
	case 1:
		fmt.Printf(args[0].AsString())
		return i.GateKeeper().UndefinedValue(), nil
	default:
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, i.GateKeeper().ToInterface(v))
		}
		fmt.Printf(args[0].AsString(), ar...)
		return i.GateKeeper().UndefinedValue(), nil
	}
}

// append appends provided arguments to an array and returns a new array or an error if arguments are invalid.
func (i *FuncInternal) append(frame int, args []IObject) (IObject, error) {
	if len(args) < 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Array:
		return i.GateKeeper().NewArray(frame, append(arg.Values(), args[1:]...)), nil
	default:
		return i.GateKeeper().NewArray(frame, args[1:]), nil
	}
}

// delete removes a key-Code pair from a Map object if the key exists and returns an undefined Code or an error if invalid.
func (i *FuncInternal) delete(_ int, args []IObject) (IObject, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Map:
		arg.Delete(args[1].AsString())
		return i.GateKeeper().UndefinedValue(), nil
	default:
		return nil, NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// splice modifies the content of an array by removing and/or adding elements at the specified index.
// frame specifies the execution context, Args contains arguments including the array, start index, delete count, and new elements.
// Returns a new array containing the removed elements or an error if arguments are invalid.
func (i *FuncInternal) splice(frame int, args []IObject) (IObject, error) {
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
	return i.GateKeeper().NewArray(frame, deleted), nil
}

// panic triggers a runtime panic with the provided argument or a default message if no argument is given.
func (i *FuncInternal) panic(_ int, args []IObject) (IObject, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("%s", args[0].AsString())
	}
	return nil, fmt.Errorf("panic")
}

// recover is a method that handles logic for recovering from specific scenarios and returns an IObject and an error.
// It takes an integer and a slice of IObject as arguments and returns a default undefined Code and nil error.
func (i *FuncInternal) recover(_ int, _ []IObject) (IObject, error) {
	return i.GateKeeper().UndefinedValue(), nil
}

// make creates a new object of the specified type, with optional len and capacity, or returns an error on failure.
func (i *FuncInternal) make(frame int, args []IObject) (IObject, error) {
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
		return i.GateKeeper().NewBytes(frame, make([]byte, count, reserve)), nil
	case ArrayType:
		return i.GateKeeper().NewArray(frame, make([]IObject, count, reserve)), nil
	case MapType:
		return i.GateKeeper().NewMap(frame, make(map[string]IObject, count)), nil
	default:
		return nil, fmt.Errorf("cannot make type %s", kind.TypeName())
	}
}

// Count returns the total number of elements in the instance and its sub-elements.
func (i *FuncInternal) Count() int {
	return 1
}

// GobEncode serializes the FuncInternal's data into a byte slice using gob encoding and returns the result or an error.
func (i *FuncInternal) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(i.id); err != nil {
		return nil, err
	}
	if err := encoder.Encode(i.fn); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the FuncInternal's data field using the gob package.
func (i *FuncInternal) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&i.id); err != nil {
		return err
	}
	if err := decoder.Decode(&i.fn); err != nil {
		return err
	}
	return nil
}
