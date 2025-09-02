package objects

import (
	"fmt"
	"time"
)

const (
	CallIdLen = iota
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

type GateCall struct {
	gk *GateKeeper
}

func NewGateCall(gk *GateKeeper) *GateCall {
	return &GateCall{
		gk: gk,
	}
}

func (h *GateCall) Call(id int, frame int, args ...IObject) (IObject, error) {
	var arg0 IObject
	if len(args) > 0 {
		arg0 = args[0]
	}
	switch id {
	case CallIdIsString:
		_, ok := arg0.(*String)
		return h.gk.Boolean(ok), nil
	case CallIdIsInt:
		_, ok := arg0.(*Int)
		return h.gk.Boolean(ok), nil
	case CallIdIsFloat:
		_, ok := arg0.(*Float)
		return h.gk.Boolean(ok), nil
	case CallIdIsBool:
		_, ok := arg0.(*Bool)
		return h.gk.Boolean(ok), nil
	case CallIdIsChar:
		_, ok := arg0.(*Char)
		return h.gk.Boolean(ok), nil
	case CallIdIsBytes:
		_, ok := arg0.(*Bytes)
		return h.gk.Boolean(ok), nil
	case CallIdIsArray:
		_, ok := arg0.(*Array)
		return h.gk.Boolean(ok), nil
	case CallIdIsMap:
		_, ok := arg0.(*Map)
		return h.gk.Boolean(ok), nil
	case CallIdIsTime:
		_, ok := arg0.(*Time)
		return h.gk.Boolean(ok), nil
	case CallIdIsError:
		_, ok := arg0.(*Error)
		return h.gk.Boolean(ok), nil
	case CallIdIsUndefined:
		_, ok := arg0.(*Undefined)
		return h.gk.Boolean(ok), nil
	case CallIdIsFunction:
		_, ok := arg0.(*FuncCompiled)
		return h.gk.Boolean(ok), nil
	case CallIdIsCallable:
		return h.gk.Boolean(arg0.CanCall()), nil
	case CallIdIsIterable:
		return h.gk.Boolean(arg0.CanIterate()), nil
	case CallIdTypeName:
		return h.gk.NewString(frame, arg0.TypeName()), nil
	case CallIdString:
		return h.gk.NewString(frame, arg0.AsString()), nil
	case CallIdInt:
		return h.gk.NewInt(frame, arg0.AsInt64()), nil
	case CallIdFloat:
		return h.gk.NewFloat(frame, arg0.AsFloat64()), nil
	case CallIdChar:
		return h.gk.NewChar(frame, rune(arg0.AsInt64())), nil
	case CallIdBool:
		return h.gk.Boolean(arg0.AsBool()), nil
	case CallIdTime:
		if v, ok := arg0.(*Time); ok {
			return h.gk.NewTime(frame, v.Value()), nil
		}
		return h.gk.NewTime(frame, time.Unix(arg0.AsInt64(), 0)), nil
	case CallIdLen:
		return h.len(frame, arg0)
	case CallIdPrintf:
		return h.printf(args...)
	case CallIdSprintf:
		return h.sprintf(frame, args...)
	case CallIdCopy:
		return arg0.Copy(frame, 0), nil
	case CallIdAppend:
		return h.append(h.gk, frame, args...)
	case CallIdSplice:
		return h.splice(frame, args...)
	case CallIdRecover:
		return h.recover()
	case CallIdPanic:
		return h.panic(arg0)
	case CallIdDelete:
		return h.delete(h.gk, args...)
	case CallIdMake:
		return h.make(h.gk, frame, args...)
	default:
		return h.gk.UndefinedValue(), fmt.Errorf("undefined function %d", id)
	}
}

// len determines the length of an array, string, bytes, or map. Returns an error if the argument type is unsupported.
func (h *GateCall) len(frame int, argIn IObject) (IObject, error) {
	switch arg := argIn.(type) {
	case *Array:
		return h.gk.NewInt(frame, int64(arg.Length())), nil
	case *String:
		return h.gk.NewInt(frame, int64(arg.Length())), nil
	case *Bytes:
		return h.gk.NewInt(frame, int64(arg.Length())), nil
	case *Map:
		return h.gk.NewInt(frame, int64(arg.Length())), nil
	default:
		return nil, NewInvalidArgumentError(0, "container", arg.TypeName())
	}
}

// Format applies a format string to a variable number of arguments, returning the formatted result as a AsString object.
func (h *GateCall) sprintf(frame int, args ...IObject) (IObject, error) {
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

// Format applies a format string to a variable number of arguments, returning the formatted result as a AsString object.
func (h *GateCall) printf(args ...IObject) (IObject, error) {
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

// Append adds one or more elements to an array or immutable array and returns the updated array or an error.
func (h *GateCall) append(gk IGateKeeper, frame int, args ...IObject) (IObject, error) {
	if len(args) < 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Array:
		return gk.NewArray(frame, append(arg.Values(), args[1:]...)), nil
	default:
		return nil, NewInvalidArgumentError(0, "array", arg.TypeName())
	}
}

// Delete removes a key-value pair from a map. Requires a map as the first argument and a string key as the second argument.
func (h *GateCall) delete(gk IGateKeeper, args ...IObject) (IObject, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *Map:
		arg.Delete(args[1].AsString())
		return gk.UndefinedValue(), nil
	default:
		return nil, NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// Splice removes or replaces existing elements and/or adds new elements in an array, returning a new array of deleted elements.
func (h *GateCall) splice(frame int, args ...IObject) (IObject, error) {
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

// panic raises an error with the provided message, halting execution if the argument count is exactly one.
func (h *GateCall) panic(arg IObject) (IObject, error) {
	// Create an error with the panic message
	err := fmt.Errorf("%s", arg.AsString())
	// Signal the error to the VM. This is the key point!
	// The VM already has a mechanism to stop execution in case of error.
	// Calling SetError will stop the main execution loop.
	// v.SetError(err) // This call should be made through an interface exposed to builtin

	// In practice, the VM panic becomes an error that stops the script
	return nil, err
}

// panic raises an error with the provided message, halting execution if the argument count is exactly one.
func (h *GateCall) recover() (IObject, error) {
	return h.gk.UndefinedValue(), nil
}

// make creates and returns a new object based on the specified arguments, supporting byte slices or arrays.
// The kind of object is determined by the first argument, while subsequent arguments define size and capacity.
// Returns an undefined value if the kind is unsupported.
func (h *GateCall) make(gk IGateKeeper, frame int, args ...IObject) (IObject, error) {
	var kind IObject
	count := 0
	reserve := 0
	switch len(args) {
	case 2:
		kind = args[0]
		count = int(args[1].AsInt64())
	case 3:
		kind = args[0]
		count = int(args[1].AsInt64())
		reserve = int(args[2].AsInt64())
	default:
		return nil, ErrInvalidArgumentsNumber
	}
	if count < 0 {
		return nil, fmt.Errorf("make: len out of range: %d", count)
	}
	if reserve < count {
		return nil, fmt.Errorf("make: cap out of range: %d", reserve)
	}
	switch kind.TypeName() {
	case BytesType:
		return gk.NewBytes(frame, make([]byte, count, reserve)), nil
	case ArrayType:
		return gk.NewArray(frame, make([]IObject, count, reserve)), nil
	case MapType:
		return gk.NewMap(frame, make(map[string]IObject, count)), nil
	default:
		return nil, fmt.Errorf("cannot make type %s", kind.TypeName())
	}
}
