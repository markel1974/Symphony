package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// BuiltinFunctions is a structure that encapsulates and provides access to predefined built-in function modules.
type BuiltinFunctions struct {
	module []*objects.FunctionModule
}

// NewBuiltinFunctions initializes a new BuiltinFunctions instance with predefined standard functions.
// It returns a pointer to the newly created BuiltinFunctions object.
func NewBuiltinFunctions() *BuiltinFunctions {
	b := &BuiltinFunctions{}
	b.module = []*objects.FunctionModule{
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "len", b.Len),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "copy", b.Copy),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "append", b.Append),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "delete", b.Delete),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "splice", b.Splice),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "format", b.Format),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "string", b.String),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "int", b.Int),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "bool", b.Bool),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "float", b.Float),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "char", b.Char),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "bytes", b.Bytes),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "time", b.Time),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "type_name", b.TypeName),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_int", b.IsInt),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_float", b.IsFloat),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_string", b.IsString),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_bool", b.IsBool),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_char", b.IsChar),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_bytes", b.IsBytes),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_array", b.IsArray),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_immutable_array", b.IsImmutableArray),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_map", b.IsMap),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_immutable_map", b.IsImmutableMap),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_iterable", b.IsIterable),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_time", b.IsTime),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_error", b.IsError),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_undefined", b.IsUndefined),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_function", b.IsFunction),
		objects.NewFunctionModule(objects.FunctionBuiltinDef, "is_callable", b.IsCallable),
	}

	return b
}

// Name returns the name of the builtin function as a string.
func (h *BuiltinFunctions) Name() string {
	return "builtin"
}

// Module returns a slice of pointers to FunctionModule objects representing the built-in function modules.
func (h *BuiltinFunctions) Module() []*objects.FunctionModule {
	return h.module
}

// TypeName returns the type name of the provided object as a string. It expects a single argument and returns an error if the argument count is invalid.
func (h *BuiltinFunctions) TypeName(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	return objects.NewString(args[0].TypeName())
}

// IsString checks if the given argument is of type String and returns TrueValue if it is, otherwise FalseValue.
func (h *BuiltinFunctions) IsString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsInt checks if the provided argument is of type Int. Returns true if it is, otherwise false. Returns an error if the number of arguments is not exactly one.
func (h *BuiltinFunctions) IsInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsFloat checks if the provided argument is of type Float and returns a boolean result wrapped in an IObject.
func (h *BuiltinFunctions) IsFloat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsBool checks if the given argument is of type Bool. It returns TrueValue if true, and FalseValue otherwise.
func (h *BuiltinFunctions) IsBool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsChar checks if the given argument is of type Char. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) IsChar(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Char); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsBytes checks if the given argument is of type *Bytes. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) IsBytes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bytes); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsArray determines whether the first argument is of type Array. Returns TrueValue if true and FalseValue otherwise.
func (h *BuiltinFunctions) IsArray(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Array); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsImmutableArray checks if the provided argument is of type ArrayImmutable. Returns TrueValue if true, otherwise FalseValue.
func (h *BuiltinFunctions) IsImmutableArray(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.ArrayImmutable); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsMap checks if the provided argument is of type Map and returns TrueValue if true; otherwise, FalseValue.
func (h *BuiltinFunctions) IsMap(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Map); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsImmutableMap checks if the first argument is of type MapImmutable and returns TrueValue if true, otherwise FalseValue.
// Returns an error if the number of arguments provided is not exactly one.
func (h *BuiltinFunctions) IsImmutableMap(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.MapImmutable); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsTime determines if the provided argument is of type Time and returns TrueValue if so, otherwise FalseValue.
func (h *BuiltinFunctions) IsTime(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsError checks if the provided argument is of type Error. Returns TrueValue if it is, else FalseValue.
func (h *BuiltinFunctions) IsError(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Error); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsUndefined checks if the provided argument is the undefined value. Returns true if undefined, otherwise false.
func (h *BuiltinFunctions) IsUndefined(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0] == objects.UndefinedValue {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsFunction checks if the provided argument is a compiled function and returns true or false accordingly.
func (h *BuiltinFunctions) IsFunction(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch args[0].(type) {
	case *objects.FunctionCompiled:
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsCallable determines if the provided argument is callable and returns true or false accordingly. Returns an error if the number of arguments is not exactly one.
func (h *BuiltinFunctions) IsCallable(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0].CanCall() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// IsIterable checks if the provided argument is iterable and returns a Boolean object indicating the result.
// Returns objects.ErrWrongNumArguments if the number of arguments is not exactly one.
func (h *BuiltinFunctions) IsIterable(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0].CanIterate() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// Len determines the length of an array, string, bytes, or map. Returns an error if the argument type is unsupported.
func (h *BuiltinFunctions) Len(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return objects.NewInt(int64(arg.Length())), nil
	case *objects.ArrayImmutable:
		return objects.NewInt(int64(arg.Length())), nil
	case *objects.String:
		return objects.NewInt(int64(arg.Length())), nil
	case *objects.Bytes:
		return objects.NewInt(int64(arg.Length())), nil
	case *objects.Map:
		return objects.NewInt(int64(arg.Length())), nil
	case *objects.MapImmutable:
		return objects.NewInt(int64(arg.Length())), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "container", arg.TypeName())
	}
}

// Range generates a range of integers, requiring 2-3 arguments: start, stop, and an optional step (> 0). Returns an IObject or error.
func (h *BuiltinFunctions) Range(args ...objects.IObject) (objects.IObject, error) {
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
		step = objects.NewInt(int64(1))
	}
	return h.rangeCreate(start.Value(), stop.Value(), step.Value()), nil
}

// rangeCreate generates an array with elements starting from `start`, ending before `stop`, incremented by `step`.
func (h *BuiltinFunctions) rangeCreate(start int64, stop int64, step int64) *objects.Array {
	array := objects.NewArray(nil)
	if start <= stop {
		for i := start; i < stop; i += step {
			array.Append(objects.NewInt(i))
		}
	} else {
		for i := start; i > stop; i -= step {
			array.Append(objects.NewInt(i))
		}
	}
	return array
}

// Format applies a format string to a variable number of arguments, returning the formatted result as a String object.
func (h *BuiltinFunctions) Format(args ...objects.IObject) (objects.IObject, error) {
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
		ar = append(ar, objects.ToInterface(v))
	}
	return objects.NewString(fmt.Sprintf(formatString.Value(), ar...))
}

// Copy returns a copy of the provided object. It accepts exactly one argument and raises an error on invalid input.
func (h *BuiltinFunctions) Copy(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	return args[0].Copy(), nil
}

// String converts an object to a string representation or returns a default value if conversion is not possible.
func (h *BuiltinFunctions) String(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return args[0], nil
	}
	v, ok := objects.ToString(args[0])
	if ok {
		return objects.NewString(v)
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// Int converts the first argument to an integer type, or returns the second argument or undefined if the conversion fails.
func (h *BuiltinFunctions) Int(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return args[0], nil
	}
	v, ok := objects.ToInt64(args[0])
	if ok {
		return objects.NewInt(v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// Float converts the provided argument to a Float object if possible, returning an error for invalid or unsupported input.
func (h *BuiltinFunctions) Float(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return args[0], nil
	}
	v, ok := objects.ToFloat64(args[0])
	if ok {
		return objects.NewFloat(v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// Bool converts the provided argument to a boolean type if possible, or returns an error if the argument count is invalid.
func (h *BuiltinFunctions) Bool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return args[0], nil
	}
	v, ok := objects.ToBool(args[0])
	if ok {
		if v {
			return objects.TrueValue, nil
		}
		return objects.FalseValue, nil
	}
	return objects.UndefinedValue, nil
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
	v, ok := objects.ToRune(args[0])
	if ok {
		return objects.NewChar(v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// Bytes creates a new byte slice object from the given input or returns an error if the arguments are invalid.
// Accepts one or two arguments: a size as an *objects.Int or convertible byte data; optional default value.
func (h *BuiltinFunctions) Bytes(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if n, ok := args[0].(*objects.Int); ok {
		if n.Value() > int64(objects.MaxBytesLen) {
			return nil, objects.ErrBytesLimit
		}
		return objects.NewBytes(make([]byte, int(n.Value()))), nil
	}
	v, ok := objects.ToByteSlice(args[0])
	if ok {
		if len(v) > objects.MaxBytesLen {
			return nil, objects.ErrBytesLimit
		}
		return objects.NewBytes(v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// Time converts the input argument(s) to a Time object if compatible, else returns a fallback or undefined values.
func (h *BuiltinFunctions) Time(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return args[0], nil
	}
	v, ok := objects.ToTime(args[0])
	if ok {
		return objects.NewTime(v), nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// Append adds one or more elements to an array or immutable array and returns the updated array or an error.
func (h *BuiltinFunctions) Append(args ...objects.IObject) (objects.IObject, error) {
	if len(args) < 2 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return objects.NewArray(append(arg.Values(), args[1:]...)), nil
	case *objects.ArrayImmutable:
		return objects.NewArray(append(arg.Values(), args[1:]...)), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "array", arg.TypeName())
	}
}

// Delete removes a key-value pair from a map. Requires a map as the first argument and a string key as the second argument.
func (h *BuiltinFunctions) Delete(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Map:
		if key, ok := args[1].(*objects.String); ok {
			arg.Delete(key.Value())
			return objects.UndefinedValue, nil
		}
		return nil, objects.NewInvalidArgumentError(1, "string", args[1].TypeName())
	default:
		return nil, objects.NewInvalidArgumentError(0, "map", arg.TypeName())
	}
}

// Splice removes or replaces existing elements and/or adds new elements in an array, returning a new array of deleted elements.
func (h *BuiltinFunctions) Splice(args ...objects.IObject) (objects.IObject, error) {
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
	return objects.NewArray(deleted), nil
}
