package objects

import (
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
	"typeName":    CallIdTypeName,
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
	Allocator
	id CallId
	fn func(frame int, args []IObject) (IObject, error)
}

// NewFuncImport creates a new FuncImport instance with the specified Id and callable function.
func newFuncInternal(factory IGateKeeper, frame int, id CallId) IObject {
	fi := &FuncInternal{
		Allocator: Allocator{gk: factory, frame: frame},
		id:        id,
		fn:        nil,
	}
	fi.prepare(id)
	return fi
}

// AsBool returns a boolean representation of the object, always returning false for FuncJit.
func (h *FuncInternal) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (h *FuncInternal) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (h *FuncInternal) AsFloat64() float64 {
	return 0
}

// AsString returns the string representation of a FuncImport object.
func (h *FuncInternal) AsString() string {
	return "<FuncInternal>"
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (h *FuncInternal) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (h *FuncInternal) Nil() bool {
	return false
}

// LogicalOp performs a logical operation between the current object and a provided IObject using the specified operator.
// Always returns nil and ErrInvalidOperator as logical operations are not supported for this object type.
func (h *FuncInternal) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(h.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and operand, returning the result or an error.
func (h *FuncInternal) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (h *FuncInternal) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (h *FuncInternal) IndexGet(_ int, _ IObject) (IObject, error) {
	return h.gk.UndefinedValue(), ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (h *FuncInternal) IndexSet(_, _ IObject) error {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (h *FuncInternal) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (h *FuncInternal) Iterable() bool {
	return false
}

// Length returns the length of the Int object.
func (h *FuncInternal) Length() int {
	return 0
}

// TypeName returns the type name of the FuncImport as a string.
func (h *FuncInternal) TypeName() string {
	return "FuncInternal:" + strconv.Itoa(int(h.id))
}

// Copy creates and returns a new FuncImport instance with the same Value field as the original object.
func (h *FuncInternal) Copy(frame int, _ int) IObject {
	return h.GateKeeper().NewFuncInternal(frame, h.id)
}

// Equals checks whether the current FuncImport is equal to another object of type IObject. Always returns false.
func (h *FuncInternal) Equals(_ IObject) bool {
	return false
}

// Call executes the GateCall with the provided GateKeeper, frame, and arguments, returning an IObject or an error.
func (h *FuncInternal) Call(frame int, args ...IObject) (uint, IObject, error) {
	v, err := h.fn(frame, args)
	return 1, v, err
}

// prepare initializes the function pointer `fn` in GateCall based on the provided CallId.
// It maps CallId enums to specific methods or operations within the GateCall struct.
// Returns an error if an undefined CallId is provided.
func (h *FuncInternal) prepare(id CallId) {
	switch id {
	case CallIdIsString:
		h.fn = h.isString
	case CallIdIsInt:
		h.fn = h.isInt
	case CallIdIsFloat:
		h.fn = h.isFloat
	case CallIdIsBool:
		h.fn = h.isBool
	case CallIdIsChar:
		h.fn = h.isChar
	case CallIdIsBytes:
		h.fn = h.isBytes
	case CallIdIsArray:
		h.fn = h.isArray
	case CallIdIsMap:
		h.fn = h.isMap
	case CallIdIsTime:
		h.fn = h.isTime
	case CallIdIsError:
		h.fn = h.isError
	case CallIdIsUndefined:
		h.fn = h.isUndefined
	case CallIdIsFunction:
		h.fn = h.isFunction
	case CallIdIsIterable:
		h.fn = h.isIterable
	case CallIdTypeName:
		h.fn = h.callIdTypeName
	case CallIdString:
		h.fn = h.callIdString
	case CallIdInt:
		h.fn = h.int
	case CallIdFloat:
		h.fn = h.float
	case CallIdChar:
		h.fn = h.char
	case CallIdBool:
		h.fn = h.bool
	case CallIdTime:
		h.fn = h.time
	case CallIdLen:
		h.fn = h.len
	case CallIdPrintf:
		h.fn = h.printf
	case CallIdSprintf:
		h.fn = h.sprintf
	case CallIdCopy:
		h.fn = h.copy
	case CallIdAppend:
		h.fn = h.append
	case CallIdSplice:
		h.fn = h.splice
	case CallIdRecover:
		h.fn = h.recover
	case CallIdPanic:
		h.fn = h.panic
	case CallIdDelete:
		h.fn = h.delete
	case CallIdMake:
		h.fn = h.make
	default:
		h.fn = h.undefined
	}
}

// undefined is a placeholder function for handling undefined arguments.
func (h *FuncInternal) undefined(_ int, _ []IObject) (IObject, error) {
	return h.gk.UndefinedValue(), nil
}

// isString verifies if the first argument is of type String and returns a Boolean indicating the result.
func (h *FuncInternal) isString(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*String)
	return h.gk.Boolean(ok), nil
}

// isInt checks if the first argument in the args slice is of type Int and returns a boolean IObject indicating the result.
func (h *FuncInternal) isInt(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Int)
	return h.gk.Boolean(ok), nil
}

// isFloat checks if the first argument is of type Float and returns a Boolean result.
func (h *FuncInternal) isFloat(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Float)
	return h.gk.Boolean(ok), nil
}

// isBool checks if the first element in the args slice is of type Bool and returns a Boolean object accordingly.
func (h *FuncInternal) isBool(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Bool)
	return h.gk.Boolean(ok), nil
}

// isChar checks if the first argument is of type *Char and returns a boolean result wrapped in an IObject.
func (h *FuncInternal) isChar(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Char)
	return h.gk.Boolean(ok), nil
}

// isBytes checks if the first argument is of type Bytes and returns a boolean result wrapped in IObject.
func (h *FuncInternal) isBytes(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Bytes)
	return h.gk.Boolean(ok), nil
}

// isArray checks if the first argument is of type Array and returns a boolean IObject accordingly.
func (h *FuncInternal) isArray(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Array)
	return h.gk.Boolean(ok), nil
}

// isMap checks if the first argument is of type *Map and returns a boolean IObject and an error if applicable.
func (h *FuncInternal) isMap(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Map)
	return h.gk.Boolean(ok), nil
}

// isTime verifies if the provided argument is of type Time and returns a boolean result or an error.
func (h *FuncInternal) isTime(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Time)
	return h.gk.Boolean(ok), nil
}

// isError determines if the given argument is of type *Error and returns a boolean result with any error encountered.
func (h *FuncInternal) isError(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Error)
	return h.gk.Boolean(ok), nil
}

// isUndefined checks if the first argument is of type Undefined and returns a boolean result or an error if invalid input.
func (h *FuncInternal) isUndefined(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*Undefined)
	return h.gk.Boolean(ok), nil
}

// isFunction checks if the first argument is of type *FuncCompiled and returns a boolean IObject accordingly.
func (h *FuncInternal) isFunction(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	_, ok := args[0].(*FuncCompiled)
	return h.gk.Boolean(ok), nil
}

// isIterable determines if the argument is iterable and returns a boolean representation of the result.
func (h *FuncInternal) isIterable(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.Boolean(args[0].Iterable()), nil
}

// callIdTypeName generates a string representation of the type name for the provided argument.
func (h *FuncInternal) callIdTypeName(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.NewString(frame, args[0].TypeName()), nil
}

// callIdString generates a new string object using the provided frame and argument, ensuring exactly one argument is passed.
func (h *FuncInternal) callIdString(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.NewString(frame, args[0].AsString()), nil
}

// callIdInt generates a new integer object using the provided frame and the integer value extracted from the input argument.
func (h *FuncInternal) int(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.NewInt(frame, args[0].AsInt64()), nil
}

// callIdFloat converts the first argument to a float and creates a new float object for the given frame.
func (h *FuncInternal) float(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.NewFloat(frame, args[0].AsFloat64()), nil
}

// callIdChar converts an integer argument to a character and returns the new character object, or an error if input is invalid.
func (h *FuncInternal) char(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.NewChar(frame, rune(args[0].AsInt64())), nil
}

// callIdBool processes a single argument and returns a Boolean object based on the argument's boolean value or an error.
func (h *FuncInternal) bool(_ int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.Boolean(args[0].AsBool()), nil
}

// callIdTime processes a method call to generate a new Time object based on the provided frame and argument.
// It expects exactly one argument, either a Time object or a value convertible to a timestamp.
// Returns the newly created Time object or an error if the input is invalid.
func (h *FuncInternal) time(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	if v, ok := args[0].(*Time); ok {
		return h.gk.NewTime(frame, v.Value()), nil
	}
	return h.gk.NewTime(frame, time.Unix(args[0].AsInt64(), 0)), nil
}

// copy creates a copy of the given IObject argument and returns it, ensuring only a single argument is provided.
func (h *FuncInternal) copy(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return args[0].Copy(frame, 0), nil
}

// len computes the length of the provided container argument (Array, String, Bytes, or Map) and returns it as an integer.
func (h *FuncInternal) len(frame int, args []IObject) (IObject, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgumentsNumber
	}
	return h.gk.NewInt(frame, int64(args[0].Length())), nil
}

// sprintf formats a string using the format specified in args[0] and additional arguments provided.
// Returns a new formatted string or an error if the argument count is invalid.
// If args contains only one element, it directly returns that element.
func (h *FuncInternal) sprintf(frame int, args []IObject) (IObject, error) {
	switch len(args) {
	case 0:
		return h.gk.UndefinedValue(), ErrInvalidArgumentsNumber
	case 1:
		return args[0], nil
	default:
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, h.gk.ToInterface(v))
		}
		return h.gk.NewString(frame, fmt.Sprintf(args[0].AsString(), ar...)), nil
	}
}

// printf formats and prints the given arguments based on the format string provided in the first argument.
func (h *FuncInternal) printf(_ int, args []IObject) (IObject, error) {
	switch len(args) {
	case 0:
		return h.gk.UndefinedValue(), ErrInvalidArgumentsNumber
	case 1:
		fmt.Printf(args[0].AsString())
		return h.gk.UndefinedValue(), nil
	default:
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, h.gk.ToInterface(v))
		}
		fmt.Printf(args[0].AsString(), ar...)
		return h.gk.UndefinedValue(), nil
	}
}

// append appends provided arguments to an array and returns a new array or an error if arguments are invalid.
func (h *FuncInternal) append(frame int, args []IObject) (IObject, error) {
	if len(args) < 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Array:
		return h.gk.NewArray(frame, append(arg.Values(), args[1:]...)), nil
	default:
		return nil, NewInvalidArgumentError(0, "array", arg.TypeName())
	}
}

// delete removes a key-value pair from a Map object if the key exists and returns an undefined value or an error if invalid.
func (h *FuncInternal) delete(_ int, args []IObject) (IObject, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Map:
		arg.Delete(args[1].AsString())
		return h.gk.UndefinedValue(), nil
	default:
		return nil, NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// splice modifies the content of an array by removing and/or adding elements at the specified index.
// frame specifies the execution context, args contains arguments including the array, start index, delete count, and new elements.
// Returns a new array containing the removed elements or an error if arguments are invalid.
func (h *FuncInternal) splice(frame int, args []IObject) (IObject, error) {
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
	return h.gk.NewArray(frame, deleted), nil
}

// panic triggers a runtime panic with the provided argument or a default message if no argument is given.
func (h *FuncInternal) panic(_ int, args []IObject) (IObject, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("%s", args[0].AsString())
	}
	return nil, fmt.Errorf("panic")
}

// recover is a method that handles logic for recovering from specific scenarios and returns an IObject and an error.
// It takes an integer and a slice of IObject as arguments and returns a default undefined value and nil error.
func (h *FuncInternal) recover(_ int, _ []IObject) (IObject, error) {
	return h.gk.UndefinedValue(), nil
}

// make creates a new object of the specified type, with optional length and capacity, or returns an error on failure.
func (h *FuncInternal) make(frame int, args []IObject) (IObject, error) {
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
		return h.gk.NewBytes(frame, make([]byte, count, reserve)), nil
	case ArrayType:
		return h.gk.NewArray(frame, make([]IObject, count, reserve)), nil
	case MapType:
		return h.gk.NewMap(frame, make(map[string]IObject, count)), nil
	default:
		return nil, fmt.Errorf("cannot make type %s", kind.TypeName())
	}
}

// Count returns the total number of elements in the instance and its sub-elements.
func (h *FuncInternal) Count() int {
	return 1
}
