package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// BuiltinFunctions is a structure that encapsulates and provides access to predefined built-in function modules.
type BuiltinFunctions struct {
	factory *objects.Factory
	pkg     []*objects.FuncPackage
}

// NewBuiltinFunctions initializes a new BuiltinFunctions instance with predefined standard functions.
// It returns a pointer to the newly created BuiltinFunctions object.
func NewBuiltinFunctions(factory *objects.Factory) *BuiltinFunctions {
	b := &BuiltinFunctions{
		factory: factory,
	}
	b.pkg = []*objects.FuncPackage{
		factory.NewFuncPackage(objects.FuncBuiltinDef, "len", b.len),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "copy", b.copy),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "append", b.append),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "delete", b.delete),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "splice", b.splice),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "format", b.format),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "string", b.string),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "int", b.int),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "bool", b.bool),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "float", b.float),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "char", b.Char),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "bytes", b.bytes),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "time", b.time),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "type_name", b.TypeName),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_int", b.isInt),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_float", b.isFloat),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_string", b.isString),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_bool", b.isBool),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_char", b.isChar),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_bytes", b.isBytes),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_array", b.isArray),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_immutable_array", b.isImmutableArray),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_map", b.isMap),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_immutable_map", b.isImmutableMap),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_iterable", b.isIterable),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_time", b.isTime),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_error", b.isError),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_undefined", b.isUndefined),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_function", b.isFunction),
		factory.NewFuncPackage(objects.FuncBuiltinDef, "is_callable", b.isCallable),
	}
	return b
}

// Name returns the name of the builtin function as a string.
func (h *BuiltinFunctions) Name() string {
	return "builtin"
}

// Package returns a slice of pointers to FuncPackage objects representing the built-in function modules.
func (h *BuiltinFunctions) Package() []*objects.FuncPackage {
	return h.pkg
}

// TypeName returns the type name of the provided object as a string. It expects a single argument and returns an error if the argument count is invalid.
func (h *BuiltinFunctions) TypeName(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	return h.factory.NewString(objects.FrameUndefined, args[0].TypeName())
}

// IsString checks if the given argument is of type String and returns TrueValue if it is, otherwise FalseValue.
func (h *BuiltinFunctions) isString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsInt checks if the provided argument is of type Int. Returns true if it is, otherwise false. Returns an error if the number of arguments is not exactly one.
func (h *BuiltinFunctions) isInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsFloat checks if the provided argument is of type Float and returns a boolean result wrapped in an IObject.
func (h *BuiltinFunctions) isFloat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsBool checks if the given argument is of type Bool. It returns TrueValue if true, and FalseValue otherwise.
func (h *BuiltinFunctions) isBool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsChar checks if the given argument is of type Char. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) isChar(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Char); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsBytes checks if the given argument is of type *Bytes. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) isBytes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bytes); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsArray determines whether the first argument is of type Array. Returns TrueValue if true and FalseValue otherwise.
func (h *BuiltinFunctions) isArray(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Array); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsImmutableArray checks if the provided argument is of type ArrayImmutable. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) isImmutableArray(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.ArrayImmutable); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsMap checks if the provided argument is of type Map and returns TrueValue if true; otherwise, FalseValue.
func (h *BuiltinFunctions) isMap(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Map); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsImmutableMap checks if the first argument is of type MapImmutable and returns TrueValue if true, otherwise FalseValue.
// Returns an error if the number of arguments provided is not exactly one.
func (h *BuiltinFunctions) isImmutableMap(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.MapImmutable); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsTime determines if the provided argument is of type Time and returns TrueValue if so, otherwise FalseValue.
func (h *BuiltinFunctions) isTime(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsError checks if the provided argument is of type Error. Returns TrueValue if it is, else FalseValue.
func (h *BuiltinFunctions) isError(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Error); ok {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsUndefined checks if the provided argument is the undefined value. Returns true if undefined, otherwise false.
func (h *BuiltinFunctions) isUndefined(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0] == h.factory.UndefinedValue() {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsFunction checks if the provided argument is a compiled function and returns true or false accordingly.
func (h *BuiltinFunctions) isFunction(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch args[0].(type) {
	case *objects.FuncCompiled:
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsCallable determines if the provided argument is callable and returns true or false accordingly. Returns an error if the number of arguments is not exactly one.
func (h *BuiltinFunctions) isCallable(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0].CanCall() {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// IsIterable checks if the provided argument is iterable and returns a Boolean object indicating the result.
// Returns objects.ErrWrongNumArguments if the number of arguments is not exactly one.
func (h *BuiltinFunctions) isIterable(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0].CanIterate() {
		return h.factory.TrueValue(), nil
	}
	return h.factory.FalseValue(), nil
}

// len determines the length of an array, string, bytes, or map. Returns an error if the argument type is unsupported.
func (h *BuiltinFunctions) len(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return h.factory.NewInt(objects.FrameUndefined, int64(arg.Length())), nil
	case *objects.ArrayImmutable:
		return h.factory.NewInt(objects.FrameUndefined, int64(arg.Length())), nil
	case *objects.String:
		return h.factory.NewInt(objects.FrameUndefined, int64(arg.Length())), nil
	case *objects.Bytes:
		return h.factory.NewInt(objects.FrameUndefined, int64(arg.Length())), nil
	case *objects.Map:
		return h.factory.NewInt(objects.FrameUndefined, int64(arg.Length())), nil
	case *objects.MapImmutable:
		return h.factory.NewInt(objects.FrameUndefined, int64(arg.Length())), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "container", arg.TypeName())
	}
}

// Range generates a range of integers, requiring 2-3 arguments: start, stop, and an optional step (> 0). Returns an IObject or error.
func (h *BuiltinFunctions) rangeInit(args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs < 2 || numArgs > 3 {
		return nil, objects.ErrWrongNumArguments
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
		return nil, objects.ErrWrongNumArguments
	}
	if step == nil {
		step = h.factory.NewInt(objects.FrameUndefined, int64(1))
	}
	startV := start.Value()
	stopV := stop.Value()
	stepV := step.Value()
	ret := h.factory.NewArray(objects.FrameUndefined, nil)
	if startV <= stopV {
		for i := startV; i < stopV; i += stepV {
			ret.Append(h.factory.NewInt(objects.FrameUndefined, i))
		}
	} else {
		for i := startV; i > stopV; i -= stepV {
			ret.Append(h.factory.NewInt(objects.FrameUndefined, i))
		}
	}
	return ret, nil
}

// Format applies a format string to a variable number of arguments, returning the formatted result as a String object.
func (h *BuiltinFunctions) format(args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, objects.ErrWrongNumArguments
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
		ar = append(ar, h.factory.ToInterface(v))
	}
	return h.factory.NewString(objects.FrameUndefined, fmt.Sprintf(formatString.Value(), ar...))
}

// Copy returns a copy of the provided object. It accepts exactly one argument and raises an error on invalid input.
func (h *BuiltinFunctions) copy(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	return args[0].Copy(objects.FrameUndefined), nil
}

// String converts an object to a string representation or returns a default value if conversion is not possible.
func (h *BuiltinFunctions) string(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return args[0], nil
	}
	v, ok := h.factory.ToString(args[0])
	if ok {
		return h.factory.NewString(objects.FrameUndefined, v)
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return h.factory.UndefinedValue(), nil
}

// Int converts the first argument to an integer type, or returns the second argument or undefined if the conversion fails.
func (h *BuiltinFunctions) int(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return args[0], nil
	}
	v, ok := h.factory.ToInt64(args[0])
	if ok {
		return h.factory.NewInt(objects.FrameUndefined, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return h.factory.UndefinedValue(), nil
}

// Float converts the provided argument to a Float object if possible, returning an error for invalid or unsupported input.
func (h *BuiltinFunctions) float(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return args[0], nil
	}
	v, ok := h.factory.ToFloat64(args[0])
	if ok {
		return h.factory.NewFloat(objects.FrameUndefined, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return h.factory.UndefinedValue(), nil
}

// Bool converts the provided argument to a boolean type if possible, or returns an error if the argument count is invalid.
func (h *BuiltinFunctions) bool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return args[0], nil
	}
	v, ok := h.factory.ToBool(args[0])
	if ok {
		if v {
			return h.factory.TrueValue(), nil
		}
		return h.factory.FalseValue(), nil
	}
	return h.factory.UndefinedValue(), nil
}

// Char converts the given argument to a Char object if possible, returning the result, or a default fallback if provided.
func (h *BuiltinFunctions) Char(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Char); ok {
		return args[0], nil
	}
	v, ok := h.factory.ToRune(args[0])
	if ok {
		return h.factory.NewChar(objects.FrameUndefined, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return h.factory.UndefinedValue(), nil
}

// Bytes creates a new byte slice object from the given input or returns an error if the arguments are invalid.
// Accepts one or two arguments: a size as an *objects.Int or convertible byte data; optional default value.
func (h *BuiltinFunctions) bytes(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if n, ok := args[0].(*objects.Int); ok {
		if n.Value() > int64(objects.MaxBytesLen) {
			return nil, objects.ErrBytesLimit
		}
		return h.factory.NewBytes(objects.FrameUndefined, make([]byte, int(n.Value()))), nil
	}
	v, ok := h.factory.ToByteSlice(args[0])
	if ok {
		if len(v) > objects.MaxBytesLen {
			return nil, objects.ErrBytesLimit
		}
		return h.factory.NewBytes(objects.FrameUndefined, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return h.factory.UndefinedValue(), nil
}

// Time converts the input argument(s) to a Time object if compatible, else returns a fallback or undefined values.
func (h *BuiltinFunctions) time(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return args[0], nil
	}
	v, ok := h.factory.ToTime(args[0])
	if ok {
		return h.factory.NewTime(objects.FrameUndefined, v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return h.factory.UndefinedValue(), nil
}

// Append adds one or more elements to an array or immutable array and returns the updated array or an error.
func (h *BuiltinFunctions) append(args ...objects.IObject) (objects.IObject, error) {
	if len(args) < 2 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return h.factory.NewArray(objects.FrameUndefined, append(arg.Values(), args[1:]...)), nil
	case *objects.ArrayImmutable:
		return h.factory.NewArray(objects.FrameUndefined, append(arg.Values(), args[1:]...)), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "array", arg.TypeName())
	}
}

// Delete removes a key-value pair from a map. Requires a map as the first argument and a string key as the second argument.
func (h *BuiltinFunctions) delete(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Map:
		if key, ok := args[1].(*objects.String); ok {
			arg.Delete(key.Value())
			return h.factory.UndefinedValue(), nil
		}
		return nil, objects.NewInvalidArgumentError(1, "string", args[1].TypeName())
	default:
		return nil, objects.NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// Splice removes or replaces existing elements and/or adds new elements in an array, returning a new array of deleted elements.
func (h *BuiltinFunctions) splice(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
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
	return h.factory.NewArray(objects.FrameUndefined, deleted), nil
}
