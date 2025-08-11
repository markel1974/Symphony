package modules

import (
	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/formats"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var builtinFuncs = []*objects.BuiltinFunction{
	{
		Name:  "len",
		Value: builtinLen,
	},
	{
		Name:  "copy",
		Value: builtinCopy,
	},
	{
		Name:  "append",
		Value: builtinAppend,
	},
	{
		Name:  "delete",
		Value: builtinDelete,
	},
	{
		Name:  "splice",
		Value: builtinSplice,
	},
	{
		Name:  "string",
		Value: builtinString,
	},
	{
		Name:  "int",
		Value: builtinInt,
	},
	{
		Name:  "bool",
		Value: builtinBool,
	},
	{
		Name:  "float",
		Value: builtinFloat,
	},
	{
		Name:  "char",
		Value: builtinChar,
	},
	{
		Name:  "bytes",
		Value: builtinBytes,
	},
	{
		Name:  "time",
		Value: builtinTime,
	},
	{
		Name:  "is_int",
		Value: builtinIsInt,
	},
	{
		Name:  "is_float",
		Value: builtinIsFloat,
	},
	{
		Name:  "is_string",
		Value: builtinIsString,
	},
	{
		Name:  "is_bool",
		Value: builtinIsBool,
	},
	{
		Name:  "is_char",
		Value: builtinIsChar,
	},
	{
		Name:  "is_bytes",
		Value: builtinIsBytes,
	},
	{
		Name:  "is_array",
		Value: builtinIsArray,
	},
	{
		Name:  "is_immutable_array",
		Value: builtinIsImmutableArray,
	},
	{
		Name:  "is_map",
		Value: builtinIsMap,
	},
	{
		Name:  "is_immutable_map",
		Value: builtinIsImmutableMap,
	},
	{
		Name:  "is_iterable",
		Value: builtinIsIterable,
	},
	{
		Name:  "is_time",
		Value: builtinIsTime,
	},
	{
		Name:  "is_error",
		Value: builtinIsError,
	},
	{
		Name:  "is_undefined",
		Value: builtinIsUndefined,
	},
	{
		Name:  "is_function",
		Value: builtinIsFunction,
	},
	{
		Name:  "is_callable",
		Value: builtinIsCallable,
	},
	{
		Name:  "type_name",
		Value: builtinTypeName,
	},
	{
		Name:  "format",
		Value: builtinFormat,
	},
	{
		Name:  "range",
		Value: builtinRange,
	},
}

// GetAllBuiltinFunctions returns all builtin function objects.
func GetAllBuiltinFunctions() []*objects.BuiltinFunction {
	return append([]*objects.BuiltinFunction{}, builtinFuncs...)
}

func builtinTypeName(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	return &objects.String{Value: args[0].TypeName()}, nil
}

func builtinIsString(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsInt(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsFloat(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsBool(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bool); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsChar(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Char); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsBytes(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Bytes); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsArray(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Array); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsImmutableArray(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.ImmutableArray); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsMap(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Map); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsImmutableMap(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.ImmutableMap); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsTime(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsError(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Error); ok {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsUndefined(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if args[0] == objects.UndefinedValue {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsFunction(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	switch args[0].(type) {
	case *objects.CompiledFunction:
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsCallable(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if args[0].CanCall() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

func builtinIsIterable(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	if args[0].CanIterate() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// len(obj object) => int
func builtinLen(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return &objects.Int{Value: int64(len(arg.Value))}, nil
	case *objects.ImmutableArray:
		return &objects.Int{Value: int64(len(arg.Value))}, nil
	case *objects.String:
		return &objects.Int{Value: int64(len(arg.Value))}, nil
	case *objects.Bytes:
		return &objects.Int{Value: int64(len(arg.Value))}, nil
	case *objects.Map:
		return &objects.Int{Value: int64(len(arg.Value))}, nil
	case *objects.ImmutableMap:
		return &objects.Int{Value: int64(len(arg.Value))}, nil
	default:
		return nil, errors.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array/string/bytes/map",
			Found:    arg.TypeName(),
		}
	}
}

// range(start, stop[, step])
func builtinRange(args ...objects.Object) (objects.Object, error) {
	numArgs := len(args)
	if numArgs < 2 || numArgs > 3 {
		return nil, errors.ErrWrongNumArguments
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

			return nil, errors.ErrInvalidArgumentType{
				Name:     name,
				Expected: "int",
				Found:    arg.TypeName(),
			}
		}
		if i == 2 && v.Value <= 0 {
			return nil, errors.ErrInvalidRangeStep
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

	if step == nil {
		step = &objects.Int{Value: int64(1)}
	}

	return buildRange(start.Value, stop.Value, step.Value), nil
}

func buildRange(start, stop, step int64) *objects.Array {
	array := &objects.Array{}
	if start <= stop {
		for i := start; i < stop; i += step {
			array.Value = append(array.Value, &objects.Int{
				Value: i,
			})
		}
	} else {
		for i := start; i > stop; i -= step {
			array.Value = append(array.Value, &objects.Int{
				Value: i,
			})
		}
	}
	return array
}

func builtinFormat(args ...objects.Object) (objects.Object, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, errors.ErrWrongNumArguments
	}
	format, ok := args[0].(*objects.String)
	if !ok {
		return nil, errors.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := formats.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	return &objects.String{Value: s}, nil
}

func builtinCopy(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	return args[0].Copy(), nil
}

func builtinString(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.String); ok {
		return args[0], nil
	}
	v, ok := objects.ToString(args[0])
	if ok {
		if len(v) > objects.MaxStringLen {
			return nil, errors.ErrStringLimit
		}
		return &objects.String{Value: v}, nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

func builtinInt(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Int); ok {
		return args[0], nil
	}
	v, ok := objects.ToInt64(args[0])
	if ok {
		return &objects.Int{Value: v}, nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

func builtinFloat(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Float); ok {
		return args[0], nil
	}
	v, ok := objects.ToFloat64(args[0])
	if ok {
		return &objects.Float{Value: v}, nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

func builtinBool(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
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

func builtinChar(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Char); ok {
		return args[0], nil
	}
	v, ok := objects.ToRune(args[0])
	if ok {
		return &objects.Char{Value: v}, nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

func builtinBytes(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, errors.ErrWrongNumArguments
	}

	// bytes(N) => create a new bytes with given size N
	if n, ok := args[0].(*objects.Int); ok {
		if n.Value > int64(objects.MaxBytesLen) {
			return nil, errors.ErrBytesLimit
		}
		return &objects.Bytes{Value: make([]byte, int(n.Value))}, nil
	}
	v, ok := objects.ToByteSlice(args[0])
	if ok {
		if len(v) > objects.MaxBytesLen {
			return nil, errors.ErrBytesLimit
		}
		return &objects.Bytes{Value: v}, nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

func builtinTime(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, errors.ErrWrongNumArguments
	}
	if _, ok := args[0].(*objects.Time); ok {
		return args[0], nil
	}
	v, ok := objects.ToTime(args[0])
	if ok {
		return &objects.Time{Value: v}, nil
	}
	if argsLen == 2 {
		return args[1], nil
	}
	return objects.UndefinedValue, nil
}

// append(arr, items...)
func builtinAppend(args ...objects.Object) (objects.Object, error) {
	if len(args) < 2 {
		return nil, errors.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Array:
		return &objects.Array{Value: append(arg.Value, args[1:]...)}, nil
	case *objects.ImmutableArray:
		return &objects.Array{Value: append(arg.Value, args[1:]...)}, nil
	default:
		return nil, errors.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    arg.TypeName(),
		}
	}
}

// builtinDelete deletes Map keys
// usage: delete(map, "key")
// key must be a string
func builtinDelete(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if argsLen != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	switch arg := args[0].(type) {
	case *objects.Map:
		if key, ok := args[1].(*objects.String); ok {
			delete(arg.Value, key.Value)
			return objects.UndefinedValue, nil
		}
		return nil, errors.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	default:
		return nil, errors.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "map",
			Found:    arg.TypeName(),
		}
	}
}

// builtinSplice deletes and changes given Array, returns deleted items.
// usage:
// deleted_items := splice(array[,start[,delete_count[,item1[,item2[,...]]]])
func builtinSplice(args ...objects.Object) (objects.Object, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, errors.ErrWrongNumArguments
	}

	array, ok := args[0].(*objects.Array)
	if !ok {
		return nil, errors.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    args[0].TypeName(),
		}
	}
	arrayLen := len(array.Value)

	var startIdx int
	if argsLen > 1 {
		arg1, ok := args[1].(*objects.Int)
		if !ok {
			return nil, errors.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "int",
				Found:    args[1].TypeName(),
			}
		}
		startIdx = int(arg1.Value)
		if startIdx < 0 || startIdx > arrayLen {
			return nil, errors.ErrIndexOutOfBounds
		}
	}

	delCount := len(array.Value)
	if argsLen > 2 {
		arg2, ok := args[2].(*objects.Int)
		if !ok {
			return nil, errors.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "int",
				Found:    args[2].TypeName(),
			}
		}
		delCount = int(arg2.Value)
		if delCount < 0 {
			return nil, errors.ErrIndexOutOfBounds
		}
	}
	// if count of to be deleted items is bigger than expected, truncate it
	if startIdx+delCount > arrayLen {
		delCount = arrayLen - startIdx
	}
	// delete items
	endIdx := startIdx + delCount
	deleted := append([]objects.Object{}, array.Value[startIdx:endIdx]...)

	head := array.Value[:startIdx]
	var items []objects.Object
	if argsLen > 3 {
		items = make([]objects.Object, 0, argsLen-3)
		for i := 3; i < argsLen; i++ {
			items = append(items, args[i])
		}
	}
	items = append(items, array.Value[endIdx:]...)
	array.Value = append(head, items...)

	// return deleted items
	return &objects.Array{Value: deleted}, nil
}

func GetBuiltin(idx int) *objects.BuiltinFunction {
	return builtinFuncs[idx]
}
