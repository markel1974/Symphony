package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// _builtinFunctions is a slice of predefined builtin functions used for various standard operations within the system.
var _builtinFunctions = []*objects.FunctionBuiltin{
	objects.NewFunctionBuiltin("len", builtinLen),
	objects.NewFunctionBuiltin("copy", builtinCopy),
	objects.NewFunctionBuiltin("append", builtinAppend),
	objects.NewFunctionBuiltin("delete", builtinDelete),
	objects.NewFunctionBuiltin("splice", builtinSplice),
	objects.NewFunctionBuiltin("format", builtinFormat),
	objects.NewFunctionBuiltin("range", builtinRange),
	objects.NewFunctionBuiltin("string", builtinString),
	objects.NewFunctionBuiltin("int", builtinInt),
	objects.NewFunctionBuiltin("bool", builtinBool),
	objects.NewFunctionBuiltin("float", builtinFloat),
	objects.NewFunctionBuiltin("char", builtinChar),
	objects.NewFunctionBuiltin("bytes", builtinBytes),
	objects.NewFunctionBuiltin("time", builtinTime),
	objects.NewFunctionBuiltin("type_name", builtinTypeName),
	objects.NewFunctionBuiltin("is_int", builtinIsInt),
	objects.NewFunctionBuiltin("is_float", builtinIsFloat),
	objects.NewFunctionBuiltin("is_string", builtinIsString),
	objects.NewFunctionBuiltin("is_bool", builtinIsBool),
	objects.NewFunctionBuiltin("is_char", builtinIsChar),
	objects.NewFunctionBuiltin("is_bytes", builtinIsBytes),
	objects.NewFunctionBuiltin("is_array", builtinIsArray),
	objects.NewFunctionBuiltin("is_immutable_array", builtinIsImmutableArray),
	objects.NewFunctionBuiltin("is_map", builtinIsMap),
	objects.NewFunctionBuiltin("is_immutable_map", builtinIsImmutableMap),
	objects.NewFunctionBuiltin("is_iterable", builtinIsIterable),
	objects.NewFunctionBuiltin("is_time", builtinIsTime),
	objects.NewFunctionBuiltin("is_error", builtinIsError),
	objects.NewFunctionBuiltin("is_undefined", builtinIsUndefined),
	objects.NewFunctionBuiltin("is_function", builtinIsFunction),
	objects.NewFunctionBuiltin("is_callable", builtinIsCallable),
}

// builtinTypeName returns the type name of the given object argument as a string, or an error if the argument count is invalid.
func builtinTypeName(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	return objects.NewString(args[0].TypeName())
}

// builtinIsString checks if the given argument is of type String.
// Returns TrueValue if the argument is a String, otherwise FalseValue.
// Returns an error if the number of arguments is not exactly 1.
func builtinIsString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsInt checks if the provided argument is of type Int.
// Returns TrueValue if the argument is an Int; otherwise, FalseValue.
// Returns ErrWrongNumArguments if the number of arguments is not exactly one.
func builtinIsInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsFloat checks if the given argument is of type Float and returns TrueValue if it is, otherwise FalseValue.
// It expects exactly one argument and returns an error if the number of arguments is incorrect.
func builtinIsFloat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsBool checks if the provided argument is of boolean type and returns true or false accordingly.
// Returns an error if the number of arguments is not equal to 1.
func builtinIsBool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsChar checks if the provided argument is of type Char, returning TrueValue if it is, otherwise FalseValue.
func builtinIsChar(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Char); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsBytes checks if the given argument is of type Bytes. Returns TrueValue if yes, otherwise FalseValue.
func builtinIsBytes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bytes); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsArray checks if the provided argument is an Array. Returns TrueValue if it is, otherwise returns FalseValue.
// Returns ErrWrongNumArguments if the number of arguments is not exactly 1.
func builtinIsArray(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Array); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsImmutableArray checks if the provided argument is of type ArrayImmutable and returns TrueValue or FalseValue.
// Returns an error if the number of arguments is not exactly one.
func builtinIsImmutableArray(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.ArrayImmutable); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsMap checks if the provided argument is of type map and returns true if it is, otherwise false.
// Returns an error if the number of arguments is not exactly one.
func builtinIsMap(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Map); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsImmutableMap checks if the provided argument is of type MapImmutable and returns TrueValue or FalseValue.
func builtinIsImmutableMap(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.MapImmutable); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsTime determines if the given argument is of type *objects.Time and returns true or false accordingly.
// Returns an error if the number of arguments provided is not exactly one.
func builtinIsTime(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsError checks if the provided argument is an error object. Returns true if it is, otherwise false.
func builtinIsError(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Error); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsUndefined checks if the given argument is the special Undefined value and returns a boolean result.
func builtinIsUndefined(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0] == objects.UndefinedValue {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsFunction checks if the given argument is a FunctionCompiled. Returns TrueValue for true, else FalseValue.
// Returns an error if the number of arguments is not 1.
func builtinIsFunction(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch args[0].(type) {
	case *objects.FunctionCompiled:
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsCallable checks if the given object is callable and returns a true or false value accordingly.
// Returns an error if the number of arguments is not exactly one.
func builtinIsCallable(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0].CanCall() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinIsIterable checks if the given argument is iterable.
// Takes a single argument and returns true if the object can iterate, otherwise false.
// Returns an error if the number of arguments is not exactly one.
func builtinIsIterable(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	if args[0].CanIterate() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// builtinLen calculates the length of a supported object such as array, string, bytes, or map.
// Returns an Int object representing the length or an error if the argument is invalid.
func builtinLen(args ...objects.IObject) (objects.IObject, error) {
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
		return nil, objects.NewInvalidArgumentError("first", "array/string/bytes/map", arg.TypeName())
	}
}

// builtinRange creates a range of integers based on the start, stop, and optional step parameters.
// Returns an error if arguments are missing, of incorrect types, or if step is less than or equal to zero.
func builtinRange(args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs < 2 || numArgs > 3 {
		return nil, objects.ErrWrongNumArguments
	}
	var start, stop, step *objects.Int

	for i, arg := range args {
		v, ok := args[i].(*objects.Int)
		if !ok {
			var name string
			switch i {
			case 0:
				name = "start"
			case 1:
				name = "stop"
			case 2:
				name = "step"
			}
			return nil, objects.NewInvalidArgumentError(name, "int", arg.TypeName())
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
	return buildRange(start.Value(), stop.Value(), step.Value()), nil
}

// buildRange generates an array of integers from start to stop with the specified step, supporting both ascending and descending ranges.
func buildRange(start int64, stop int64, step int64) *objects.Array {
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

// builtinFormat formats a string using the given format string and additional arguments, returning a new string object.
// Returns an error if no arguments are provided or if the first argument is not a string.
func builtinFormat(args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	formatString, ok := args[0].(*objects.String)
	if !ok {
		return nil, objects.NewInvalidArgumentError("format", "string", args[0].TypeName())
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

// builtinCopy creates and returns a copy of the provided object. Only one argument is expected; additional arguments result in an error.
func builtinCopy(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	return args[0].Copy(), nil
}

// builtinString converts the provided argument to a string or returns the default value if specified and supported.
// It returns an error if the argument count is invalid or if the string size exceeds the defined limit.
func builtinString(args ...objects.IObject) (objects.IObject, error) {
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

// builtinInt converts an object to an integer if possible; returns the second argument or undefined for invalid inputs.
func builtinInt(args ...objects.IObject) (objects.IObject, error) {
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

// builtinFloat converts the first argument to a Float object if possible, or returns the second argument or UndefinedValue.
// It takes 1 or 2 arguments, with the second argument serving as a default return value if conversion fails.
// Returns ErrWrongNumArguments if the number of arguments is not 1 or 2.
// Returns a Float object or an error based on the conversion success.
func builtinFloat(args ...objects.IObject) (objects.IObject, error) {
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

// builtinBool converts an object to its boolean representation or returns the object itself if already a boolean.
// It returns an error if the number of arguments is not exactly one.
func builtinBool(args ...objects.IObject) (objects.IObject, error) {
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

// builtinChar converts the given object to a character object. Returns the character or a default value if provided.
func builtinChar(args ...objects.IObject) (objects.IObject, error) {
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

// builtinBytes creates a new byte array of a given size N or converts argument(s) to a byte slice.
// Returns an error if the byte size exceeds the maximum limit or if arguments are invalid.
func builtinBytes(args ...objects.IObject) (objects.IObject, error) {
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

// builtinTime converts the input into a Time object if possible or returns undefined/default value if input is invalid.
func builtinTime(args ...objects.IObject) (objects.IObject, error) {
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

// builtinAppend appends elements to an array or immutable array, returning a new array. Errors if the first argument is not an array.
func builtinAppend(args ...objects.IObject) (objects.IObject, error) {
	if len(args) < 2 {
		return nil, objects.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return objects.NewArray(append(arg.Values(), args[1:]...)), nil
	case *objects.ArrayImmutable:
		return objects.NewArray(append(arg.Values(), args[1:]...)), nil
	default:
		return nil, objects.NewInvalidArgumentError("first", "array", arg.TypeName())
	}
}

// builtinDelete removes a key from a map object and returns UndefinedValue if successful; otherwise, returns an error.
func builtinDelete(args ...objects.IObject) (objects.IObject, error) {
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
		return nil, objects.NewInvalidArgumentError("second", "string", args[1].TypeName())
	default:
		return nil, objects.NewInvalidArgumentError("first", "map", arg.TypeName())
	}
}

// builtinSplice removes specified elements from an array and optionally inserts new elements at the same position.
// It expects an array as the first argument, followed by the start index, delete count, and optional items to insert.
// Returns a new array of the removed elements or an error for invalid types, out-of-bounds indices, or argument issues.
func builtinSplice(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}

	array, ok := args[0].(*objects.Array)
	if !ok {
		return nil, objects.NewInvalidArgumentError("first", "array", args[0].TypeName())
	}
	arrayLen := array.Length()

	var startIdx int
	if argsLen > 1 {
		arg1, ok := args[1].(*objects.Int)
		if !ok {
			return nil, objects.NewInvalidArgumentError("second", "int", args[1].TypeName())
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
			return nil, objects.NewInvalidArgumentError("third", "int", args[2].TypeName())
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
