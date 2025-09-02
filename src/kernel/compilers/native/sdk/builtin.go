package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterBuiltin(NewBuiltinFunctions)
}

// BuiltinFunctions is a structure that encapsulates and provides access to predefined built-in function modules.
type BuiltinFunctions struct {
	pkg []objects.IObject
}

// NewBuiltinFunctions initializes a new BuiltinFunctions instance with predefined standard functions.
// It returns a pointer to the newly created BuiltinFunctions object.
func NewBuiltinFunctions(gk objects.IGateKeeper) IBuiltin {
	b := &BuiltinFunctions{}
	b.pkg = []objects.IObject{
		//functions
		gk.NewFuncImport(objects.FrameStatic, "len", b.len),
		gk.NewFuncImport(objects.FrameStatic, "copy", b.copy),
		gk.NewFuncImport(objects.FrameStatic, "append", b.append),
		gk.NewFuncImport(objects.FrameStatic, "delete", b.delete),
		gk.NewFuncImport(objects.FrameStatic, "splice", b.splice),
		gk.NewFuncImport(objects.FrameStatic, "format", b.format),
		gk.NewFuncImport(objects.FrameStatic, "panic", b.panic),
		gk.NewFuncImport(objects.FrameStatic, "recover", b.recover),
		//cast
		gk.NewFuncImport(objects.FrameStatic, "int", b.int),
		gk.NewFuncImport(objects.FrameStatic, "bool", b.bool),
		gk.NewFuncImport(objects.FrameStatic, "float", b.float),
		gk.NewFuncImport(objects.FrameStatic, "char", b.char),
		gk.NewFuncImport(objects.FrameStatic, "bytes", b.bytes),
		gk.NewFuncImport(objects.FrameStatic, "string", b.string),
		gk.NewFuncImport(objects.FrameStatic, "time", b.time),
		//
		gk.NewFuncImport(objects.FrameStatic, "typeName", b.typeName),
		//type check
		gk.NewFuncImport(objects.FrameStatic, "isInt", b.isInt),
		gk.NewFuncImport(objects.FrameStatic, "isFloat", b.isFloat),
		gk.NewFuncImport(objects.FrameStatic, "isString", b.isString),
		gk.NewFuncImport(objects.FrameStatic, "isBool", b.isBool),
		gk.NewFuncImport(objects.FrameStatic, "isChar", b.isChar),
		gk.NewFuncImport(objects.FrameStatic, "isBytes", b.isBytes),
		gk.NewFuncImport(objects.FrameStatic, "isArray", b.isArray),
		gk.NewFuncImport(objects.FrameStatic, "isMap", b.isMap),
		gk.NewFuncImport(objects.FrameStatic, "isIterable", b.isIterable),
		gk.NewFuncImport(objects.FrameStatic, "isTime", b.isTime),
		gk.NewFuncImport(objects.FrameStatic, "isError", b.isError),
		gk.NewFuncImport(objects.FrameStatic, "isUndefined", b.isUndefined),
		gk.NewFuncImport(objects.FrameStatic, "isFunction", b.isFunction),
		gk.NewFuncImport(objects.FrameStatic, "isCallable", b.isCallable),
	}
	return b
}

// Container retrieves the package-scoped slice of IObject instances associated with the BuiltinFunctions.
func (h *BuiltinFunctions) Container() []objects.IObject {
	return h.pkg
}

// Name returns the name of the builtin function as a string.
func (h *BuiltinFunctions) Name() string {
	return "builtin"
}

// Package returns a slice of pointers to FuncImport objects representing the built-in function modules.
func (h *BuiltinFunctions) Package() []objects.IObject {
	return h.pkg
}

// TypeName returns the type name of the provided object as a string. It expects a single argument and returns an error if the argument count is invalid.
func (h *BuiltinFunctions) typeName(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	return gk.NewString(frame, args[0].TypeName()), nil
}

// IsString checks if the given argument is of type AsString and returns TrueValue if it is, otherwise FalseValue.
func (h *BuiltinFunctions) isString(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.String); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsInt checks if the provided argument is of type Int. Returns true if it is, otherwise false. Returns an error if the number of arguments is not exactly one.
func (h *BuiltinFunctions) isInt(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Int); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsFloat checks if the provided argument is of type Float and returns a boolean result wrapped in an IObject.
func (h *BuiltinFunctions) isFloat(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Float); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsBool checks if the given argument is of type Bool. It returns TrueValue if true, and FalseValue otherwise.
func (h *BuiltinFunctions) isBool(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsChar checks if the given argument is of type Char. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) isChar(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Char); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsBytes checks if the given argument is of type *Bytes. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) isBytes(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Bytes); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsArray determines whether the first argument is of type Array. Returns TrueValue if true and FalseValue otherwise.
func (h *BuiltinFunctions) isArray(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Array); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsMap checks if the provided argument is of type Map and returns TrueValue if true; otherwise, FalseValue.
func (h *BuiltinFunctions) isMap(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Map); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsTime determines if the provided argument is of type Time and returns TrueValue if so, otherwise FalseValue.
func (h *BuiltinFunctions) isTime(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Time); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsError checks if the provided argument is of type Error. Returns TrueValue if it is, else FalseValue.
func (h *BuiltinFunctions) isError(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Error); ok {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsUndefined checks if the provided argument is the undefined value. Returns true if undefined, otherwise false.
func (h *BuiltinFunctions) isUndefined(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if args[0] == gk.UndefinedValue() {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsFunction checks if the provided argument is a compiled function and returns true or false accordingly.
func (h *BuiltinFunctions) isFunction(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	switch args[0].(type) {
	case *objects.FuncCompiled:
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsCallable determines if the provided argument is callable and returns true or false accordingly. Returns an error if the number of arguments is not exactly one.
func (h *BuiltinFunctions) isCallable(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if args[0].CanCall() {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// IsIterable checks if the provided argument is iterable and returns a Falsy object indicating the result.
// Returns objects.ErrInvalidArgumentsNumber if the number of arguments is not exactly one.
func (h *BuiltinFunctions) isIterable(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if args[0].CanIterate() {
		return gk.TrueValue(), nil
	}
	return gk.FalseValue(), nil
}

// len determines the length of an array, string, bytes, or map. Returns an error if the argument type is unsupported.
func (h *BuiltinFunctions) len(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return gk.NewInt(frame, int64(arg.Length())), nil
	case *objects.String:
		return gk.NewInt(frame, int64(arg.Length())), nil
	case *objects.Bytes:
		return gk.NewInt(frame, int64(arg.Length())), nil
	case *objects.Map:
		return gk.NewInt(frame, int64(arg.Length())), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "container", arg.TypeName())
	}
}

// Range generates a range of integers, requiring 2-3 arguments: start, stop, and an optional step (> 0). Returns an IObject or error.
func (h *BuiltinFunctions) rangeInit(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs < 2 || numArgs > 3 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	var start, stop, step *objects.Int
	for i, arg := range args {
		v, ok := args[i].(*objects.Int)
		if !ok {
			return nil, objects.NewInvalidArgumentError(i, "int", arg.TypeName())
		}
		if i == 2 && v.Value() <= 0 {
			return nil, objects.ErrInvalidRangeStep
		}
		switch i {
		case 0:
			start = v
		case 1:
			stop = v
		case 2:
			step = v
		}
	}
	if start == nil || stop == nil {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if step == nil {
		ok := false
		st := gk.NewInt(frame, int64(1))
		if step, ok = st.(*objects.Int); !ok {
			return nil, objects.ErrLimitExceed
		}
		//step = gk.NewInt(frame, int64(1))
	}
	startV := start.Value()
	stopV := stop.Value()
	stepV := step.Value()
	obj := gk.NewArray(frame, nil)
	arr, ok := obj.(*objects.Array)
	if !ok {
		return nil, fmt.Errorf("expected Array, got %T", obj)
	}
	if startV <= stopV {
		for i := startV; i < stopV; i += stepV {
			arr.Append(gk.NewInt(frame, i))
		}
	} else {
		for i := startV; i > stopV; i -= stepV {
			arr.Append(gk.NewInt(frame, i))
		}
	}
	return arr, nil
}

// Format applies a format string to a variable number of arguments, returning the formatted result as a AsString object.
func (h *BuiltinFunctions) format(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	formatString, ok := args[0].(*objects.String)
	if !ok {
		return nil, objects.NewInvalidArgumentError(0, "string", args[0].TypeName())
	}
	if numArgs == 1 {
		return formatString, nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, gk.ToInterface(v))
	}
	return gk.NewString(frame, fmt.Sprintf(formatString.Value(), ar...)), nil
}

// Copy returns a copy of the provided object. It accepts exactly one argument and raises an error on invalid input.
func (h *BuiltinFunctions) copy(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	return args[0].Copy(frame, 0), nil
}

// AsString converts an object to a string representation or returns a default value if conversion is not possible.
func (h *BuiltinFunctions) string(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.String); ok {
		return args[0], nil
	}
	v, ok := gk.ToString(args[0])
	if ok {
		return gk.NewString(frame, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return gk.UndefinedValue(), nil
}

// Int converts the first argument to an integer type, or returns the second argument or undefined if the conversion fails.
func (h *BuiltinFunctions) int(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Int); ok {
		return args[0], nil
	}
	v, ok := gk.ToInt64(args[0])
	if ok {
		return gk.NewInt(frame, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return gk.UndefinedValue(), nil
}

// Float converts the provided argument to a Float object if possible, returning an error for invalid or unsupported input.
func (h *BuiltinFunctions) float(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Float); ok {
		return args[0], nil
	}
	v, ok := gk.ToFloat64(args[0])
	if ok {
		return gk.NewFloat(frame, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return gk.UndefinedValue(), nil
}

// Bool converts the provided argument to a boolean type if possible, or returns an error if the argument count is invalid.
func (h *BuiltinFunctions) bool(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return args[0], nil
	}
	v, ok := gk.ToBool(args[0])
	if ok {
		if v {
			return gk.TrueValue(), nil
		}
		return gk.FalseValue(), nil
	}
	return gk.UndefinedValue(), nil
}

// Char converts the given argument to a Char object if possible, returning the result, or a default fallback if provided.
func (h *BuiltinFunctions) char(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Char); ok {
		return args[0], nil
	}
	v, ok := gk.ToRune(args[0])
	if ok {
		return gk.NewChar(frame, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return gk.UndefinedValue(), nil
}

// Bytes creates a new byte slice object from the given input or returns an error if the arguments are invalid.
// Accepts one or two arguments: a size as an *objects.Int or convertible byte data; optional default value.
func (h *BuiltinFunctions) bytes(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if n, ok := args[0].(*objects.Int); ok {
		return gk.NewBytes(frame, make([]byte, int(n.Value()))), nil
	}
	v, ok := gk.ToBytes(args[0])
	if ok {
		return gk.NewBytes(frame, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return gk.UndefinedValue(), nil
}

// Time converts the input argument(s) to a Time object if compatible, else returns a fallback or undefined values.
func (h *BuiltinFunctions) time(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	if _, ok := args[0].(*objects.Time); ok {
		return args[0], nil
	}
	v, ok := gk.ToTime(args[0])
	if ok {
		return gk.NewTime(frame, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return gk.UndefinedValue(), nil
}

// Append adds one or more elements to an array or immutable array and returns the updated array or an error.
func (h *BuiltinFunctions) append(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) < 2 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return gk.NewArray(frame, append(arg.Values(), args[1:]...)), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "array", arg.TypeName())
	}
}

// Delete removes a key-value pair from a map. Requires a map as the first argument and a string key as the second argument.
func (h *BuiltinFunctions) delete(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	switch arg := args[0].(type) {
	case *objects.Map:
		if key, ok := args[1].(*objects.String); ok {
			arg.Delete(key.Value())
			return gk.UndefinedValue(), nil
		}
		return nil, objects.NewInvalidArgumentError(1, "string", args[1].TypeName())
	default:
		return nil, objects.NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// Splice removes or replaces existing elements and/or adds new elements in an array, returning a new array of deleted elements.
func (h *BuiltinFunctions) splice(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}

	array, ok := args[0].(*objects.Array)
	if !ok {
		return nil, objects.NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	arrayLen := array.Length()

	var startIdx int
	if argsLen > 1 {
		arg1, ok := args[1].(*objects.Int)
		if !ok {
			return nil, objects.NewInvalidArgumentError(1, "int", args[1].TypeName())
		}
		startIdx = int(arg1.Value())
		if startIdx < 0 || startIdx > arrayLen {
			return nil, objects.ErrIndexOutOfBounds
		}
	}

	delCount := array.Length()
	if argsLen > 2 {
		arg2, ok := args[2].(*objects.Int)
		if !ok {
			return nil, objects.NewInvalidArgumentError(2, "int", args[2].TypeName())
		}
		delCount = int(arg2.Value())
		if delCount < 0 {
			return nil, objects.ErrIndexOutOfBounds
		}
	}
	if startIdx+delCount > arrayLen {
		delCount = arrayLen - startIdx
	}
	endIdx := startIdx + delCount
	deleted := append([]objects.IObject{}, array.Values()[startIdx:endIdx]...)

	head := array.Values()[:startIdx]
	var items []objects.IObject
	if argsLen > 3 {
		items = make([]objects.IObject, 0, argsLen-3)
		for i := 3; i < argsLen; i++ {
			items = append(items, args[i])
		}
	}
	items = append(items, array.Values()[endIdx:]...)
	array.Assign(append(head, items...))
	return gk.NewArray(frame, deleted), nil
}

// panic raises an error with the provided message, halting execution if the argument count is exactly one.
func (h *BuiltinFunctions) panic(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	// Create an error with the panic message
	err := fmt.Errorf("%s", args[0].AsString())
	// Signal the error to the VM. This is the key point!
	// The VM already has a mechanism to stop execution in case of error.
	// Calling SetError will stop the main execution loop.
	// v.SetError(err) // This call should be made through an interface exposed to builtin

	// In practice, the VM panic becomes an error that stops the script
	return nil, err
}

// panic raises an error with the provided message, halting execution if the argument count is exactly one.
func (h *BuiltinFunctions) recover(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	return gk.UndefinedValue(), nil
}
