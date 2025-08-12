package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// FuncAR converts a zero-argument function into a CallableFunc that enforces no arguments and returns UndefinedValue.
// Returns an error if arguments are provided.
func FuncAR(fn func()) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		fn()
		return objects.UndefinedValue, nil
	}
}

// FuncARI wraps a zero-argument function returning an int into a CallableFunc compatible with objects.IObject behavior.
// It ensures the function is called with no arguments, else returns an error for wrong number of arguments.
// The return value is converted into an objects.IObject representation using objects.NewInt.
func FuncARI(fn func() int) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		return objects.NewInt(int64(fn())), nil
	}
}

// FuncARI64
func FuncARI64(fn func() int64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		return objects.NewInt(fn()), nil
	}
}

// FuncAI64RI64 wraps a function that takes an int64 and returns an int64 into a callable object function.
func FuncAI64RI64(fn func(int64) int64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}

		i1, ok := objects.ToInt64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		return objects.NewInt(fn(i1)), nil
	}
}

// FuncAI64R adapts a function taking a single int64 argument to an objects.CallableFunc type.
// It ensures the provided argument is an int64-compatible value and calls the underlying function.
// Returns an error if the argument type is invalid or the number of arguments is incorrect.
func FuncAI64R(fn func(int64)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}

		i1, ok := objects.ToInt64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		fn(i1)
		return objects.UndefinedValue, nil
	}
}

// FuncARB wraps a function that returns a bool into a CallableFunc compatible with the IObject interface.
// It converts the bool result into either TrueValue or FalseValue and validates that no arguments are provided.
// Returns ErrWrongNumArguments if arguments are passed.
func FuncARB(fn func() bool) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		if fn() {
			return objects.TrueValue, nil
		}
		return objects.FalseValue, nil
	}
}

// FuncARE wraps a function with no arguments and a single error return into a CallableFunc for use in the objects framework.
func FuncARE(fn func() error) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		return wrapError(fn()), nil
	}
}

// FuncARS wraps a function returning a string, converting it into a callable function returning an IObject or an error.
// It validates the number of arguments and limits the byte-length of the string based on MaxStringLen.
// Returns an error if argument count is non-zero or the string length exceeds the allowed limit.
func FuncARS(fn func() string) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		v, err := objects.NewString(fn())
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncARSE wraps a function returning a string and error into a CallableFunc ensuring no arguments are passed, and string limits obeyed.
func FuncARSE(fn func() (string, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return wrapError(err), nil
		}
		v, err := objects.NewString(res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncARYE wraps a function that returns a byte slice and an error, converting it into an objects.CallableFunc instance.
// The returned CallableFunc takes no arguments and returns a byte object or an error, enforcing argument count and size limits.
func FuncARYE(fn func() ([]byte, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return wrapError(err), nil
		}
		if len(res) > objects.MaxBytesLen {
			return nil, errors.ErrBytesLimit
		}
		return objects.NewBytes(res), nil
	}
}

// FuncARF wraps a function returning a float64 into a CallableFunc that accepts no arguments and returns an IObject.
// Returns an error if any arguments are provided.
func FuncARF(fn func() float64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		return objects.NewFloat(fn()), nil
	}
}

// FuncARSs creates a CallableFunc that invokes the provided function and wraps its string results in an Array of objects.
// Returns an error if arguments are provided, or if any string exceeds the maximum length allowed.
func FuncARSs(fn func() []string) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		arr := &objects.Array{}
		for _, elem := range fn() {
			v, err := objects.NewString(elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncARIsE transforms a no-argument function returning a slice of ints and error into a CallableFunc returning an IObject.
// Returns an error if arguments are provided, or wraps any error returned by the provided function.
// Converts the int slice output into an Array of IObject, appending elements one by one.
func FuncARIsE(fn func() ([]int, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, errors.ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return wrapError(err), nil
		}
		arr := objects.NewArray(nil)
		for _, v := range res {
			arr.Append(objects.NewInt(int64(v)))
		}
		return arr, nil
	}
}

// FuncAIRIs maps an external function operating on integers to a callable function that works with IObject arguments.
// The resulting function accepts one integer-compatible argument and returns an Array of integers as IObjects.
func FuncAIRIs(fn func(int) []int) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		res := fn(i1)
		arr := objects.NewArray(nil)
		for _, v := range res {
			arr.Append(objects.NewInt(int64(v)))
		}
		return arr, nil
	}
}

// FuncAFRF creates a CallableFunc that applies a provided float64 function to the first argument and returns the result.
func FuncAFRF(fn func(float64) float64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		f1, ok := objects.ToFloat64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "float(compatible)", args[0].TypeName())
		}
		return objects.NewFloat(fn(f1)), nil
	}
}

// FuncAIR wraps a function taking an int argument, turning it into a CallableFunc that processes IObject arguments.
// Returns an error if the number of arguments is not 1 or if the argument is not convertible to an int.
func FuncAIR(fn func(int)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		fn(i1)
		return objects.UndefinedValue, nil
	}
}

// FuncAIRF converts a function from int to float64 into a CallableFunc for use with the objects.IObject interface.
// The function takes one argument, validates its type, and applies the provided function to return a Float object.
// Returns an error if the argument count differs from 1 or if the argument type is not int-compatible.
func FuncAIRF(fn func(int) float64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		return objects.NewFloat(fn(i1)), nil
	}
}

// FuncAFRI converts a Go function with a float64 input and int output into an objects.CallableFunc.
// It validates the arguments, enforcing a single float-compatible input, and wraps the output as an objects.Int.
func FuncAFRI(fn func(float64) int) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		f1, ok := objects.ToFloat64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "float(compatible)", args[0].TypeName())
		}
		return objects.NewInt(int64(fn(f1))), nil
	}
}

// FuncAFFRF creates a callable function that takes two float-compatible arguments and applies a provided binary function.
// Returns an error if arguments are not float-compatible or if the number of arguments is invalid.
func FuncAFFRF(fn func(float64, float64) float64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		f1, ok := objects.ToFloat64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "float(compatible)", args[0].TypeName())
		}
		f2, ok := objects.ToFloat64(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "float(compatible)", args[1].TypeName())
		}
		return objects.NewFloat(fn(f1, f2)), nil
	}
}

// FuncAIFRF converts a function with signature func(int, float64) float64 into an objects.CallableFunc.
// It validates arguments, converts them to appropriate types, and executes the provided function.
func FuncAIFRF(fn func(int, float64) float64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		f2, ok := objects.ToFloat64(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "float(compatible)", args[1].TypeName())
		}
		return objects.NewFloat(fn(i1, f2)), nil
	}
}

// FuncAFIRF wraps a function with signature func(float64, int) float64 into an objects.CallableFunc compatible callable.
// The wrapped function takes two arguments: a float-compatible object and an int-compatible object.
// Returns a callable function that invokes the underlying wrapped function and provides error handling for invalid arguments.
func FuncAFIRF(fn func(float64, int) float64) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		f1, ok := objects.ToFloat64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "float(compatible)", args[0].TypeName())
		}
		i2, ok := objects.ToInt(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
		return objects.NewFloat(fn(f1, i2)), nil
	}
}

// FuncAFIRB creates a CallableFunc that applies a given function fn to a float64 and int extracted from IObject arguments.
// The returned CallableFunc expects exactly two arguments: the first convertible to float64, the second to int.
// If the arguments are valid and fn evaluates to true, it returns objects.TrueValue; otherwise, objects.FalseValue.
// Returns an error if there are an incorrect number of arguments or type conversion fails.
func FuncAFIRB(fn func(float64, int) bool) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		f1, ok := objects.ToFloat64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "float(compatible)", args[0].TypeName())
		}
		i2, ok := objects.ToInt(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
		if fn(f1, i2) {
			return objects.TrueValue, nil
		}
		return objects.FalseValue, nil
	}
}

// FuncAFRB takes a function that accepts a float64 and returns a bool, and converts it into a CallableFunc.
// It validates the argument count, checks the type, and invokes the provided function, returning a truthy or falsy value.
func FuncAFRB(fn func(float64) bool) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		f1, ok := objects.ToFloat64(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "float(compatible)", args[0].TypeName())
		}
		if fn(f1) {
			return objects.TrueValue, nil
		}
		return objects.FalseValue, nil
	}
}

// FuncASRS wraps a string transformation function into a CallableFunc compatible with the IObject interface.
// It ensures the function receives a single string argument, validates its type, applies the transformation,
// checks the resulting string length against a defined limit, and returns the transformed value as a String object.
func FuncASRS(fn func(string) string) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		v, err := objects.NewString(fn(s1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASRSs converts a function that takes a string and returns a slice of strings into an objects.CallableFunc.
// It checks the argument count and type, applies the function, and wraps the result in an objects.Array.
// Returns an error if the arguments are invalid or string size exceeds the limit.
func FuncASRSs(fn func(string) []string) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		res := fn(s1)
		arr := objects.NewArray(nil)
		for _, elem := range res {
			v, err := objects.NewString(elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncASRSE wraps a function accepting a string and returning a string and error into an objects.CallableFunc.
// Returns an error if the argument is not a string or the function exceeds the maximum string length.
func FuncASRSE(fn func(string) (string, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		res, err := fn(s1)
		if err != nil {
			return wrapError(err), nil
		}
		v, err := objects.NewString(res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASRE wraps a single-argument string function to conform to the CallableFunc type, with error handling for argument checks.
func FuncASRE(fn func(string) error) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		return wrapError(fn(s1)), nil
	}
}

// FuncASSRE converts a function accepting two string arguments and returning an error into an objects.CallableFunc.
// It validates the provided arguments to ensure they are string-compatible and wraps the error, if any, into IObject.
// Returns ErrWrongNumArguments for invalid argument count or NewInvalidArgumentType for type mismatches.
func FuncASSRE(fn func(string, string) error) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		}
		return wrapError(fn(s1, s2)), nil
	}
}

// FuncASSRSs wraps a function accepting two strings and returning a slice of strings into a CallableFunc.
// It ensures the right number of arguments (2 strings) and validates their types before calling the provided function.
// The result is converted into an Array
func FuncASSRSs(fn func(string, string) []string) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[1].TypeName())

		}
		arr := objects.NewArray(nil)
		for _, res := range fn(s1, s2) {
			v, err := objects.NewString(res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncASSIRSs wraps a given function with specific parameter types and converts its behavior into an objects.CallableFunc.
// It ensures type validation for three arguments: two strings and one integer, and applies the provided function.
// Returns an objects.Array containing the results or an appropriate error if validation fails.
func FuncASSIRSs(fn func(string, string, int) []string) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 3 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		}
		i3, ok := objects.ToInt(args[2])
		if !ok {
			return nil, errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
		}
		arr := objects.NewArray(nil)
		for _, res := range fn(s1, s2, i3) {
			v, err := objects.NewString(res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncASSRI wraps a function that accepts two string arguments and returns an int, making it compatible with IObject inputs.
func FuncASSRI(fn func(string, string) int) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[0].TypeName())
		}
		return objects.NewInt(int64(fn(s1, s2))), nil
	}
}

// FuncASSRS creates a CallableFunc that invokes a function taking two strings and returning a string.
// It validates the arguments as string-compatible, applies the provided function, and ensures the result respects string limits.
func FuncASSRS(fn func(string, string) string) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		}
		v, err := objects.NewString(fn(s1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil

	}
}

// FuncASSRB converts a function that takes two strings and returns a bool into a CallableFunc, processing IObject arguments.
// It enforces that the input is exactly two string-compatible arguments, returning errors otherwise.
// Based on the function result, it returns TrueValue for true or FalseValue for false as IObject.
func FuncASSRB(fn func(string, string) bool) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		}
		if fn(s1, s2) {
			return objects.TrueValue, nil
		}
		return objects.FalseValue, nil
	}
}

// FuncASsSRS converts a function that processes a slice of strings
func FuncASsSRS(fn func([]string, string) string) objects.CallableFunc {
	return func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		var ss1 []string
		switch arg0 := args[0].(type) {
		case *objects.Array:
			for idx, a := range arg0.Values() {
				as, ok := objects.ToString(a)
				if !ok {
					return nil, errors.NewInvalidArgumentType(fmt.Sprintf("first[%d]", idx), "string(compatible)", a.TypeName())
				}
				ss1 = append(ss1, as)
			}
		case *objects.ImmutableArray:
			for idx, a := range arg0.Values() {
				as, ok := objects.ToString(a)
				if !ok {
					return nil, errors.NewInvalidArgumentType(fmt.Sprintf("first[%d]", idx), "string(compatible)", a.TypeName())
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, errors.NewInvalidArgumentType("first", "array", args[0].TypeName())
		}
		s2, ok := objects.ToString(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		}
		v, err := objects.NewString(fn(ss1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASI64RE wraps a function taking a string and int64, returning a
func FuncASI64RE(fn func(string, int64) error) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		i2, ok := objects.ToInt64(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
		return wrapError(fn(s1, i2)), nil
	}
}

// FuncAIIRE creates a CallableFunc that validates two int-compatible arguments and invokes a given function returning error.
func FuncAIIRE(fn func(int, int) error) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		i2, ok := objects.ToInt(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
		return wrapError(fn(i1, i2)), nil
	}
}

// FuncASIRS converts a Go function to a callable runtime object that takes two arguments: a string and an int.
// It validates argument types, applies the function, and checks the result for the string size constraint.
func FuncASIRS(fn func(string, int) string) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 2 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		i2, ok := objects.ToInt(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
		v, err := objects.NewString(fn(s1, i2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASIIRE wraps a given function with string and two integer arguments into a CallableFunc returning IObject and error.
// It validates input argument count and types, returning errors for invalid values or wrapping the function's execution.
func FuncASIIRE(fn func(string, int, int) error) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 3 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		i2, ok := objects.ToInt(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
		i3, ok := objects.ToInt(args[2])
		if !ok {
			return nil, errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
		}
		return wrapError(fn(s1, i2, i3)), nil
	}
}

// FuncAYRIE converts a function of type func([]byte) (int, error) into an objects.CallableFunc compatible function.
// It enforces argument count and type, invoking the wrapped function with a byte slice and returning its result or error.
func FuncAYRIE(fn func([]byte) (int, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		y1, ok := objects.ToByteSlice(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "bytes(compatible)", args[0].TypeName())
		}
		res, err := fn(y1)
		if err != nil {
			return wrapError(err), nil
		}
		return objects.NewInt(int64(res)), nil
	}
}

// FuncAYRS wraps a function transforming a byte slice into a string as a callable that expects a single IObject argument.
// It converts the input argument into a byte slice, applies the provided function, and returns the result as a String object.
func FuncAYRS(fn func([]byte) string) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		y1, ok := objects.ToByteSlice(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "bytes(compatible)", args[0].TypeName())
		}
		v, err := objects.NewString(fn(y1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASRIE wraps a function that converts a string to an int and error into a CallableFunc implementation.
func FuncASRIE(fn func(string) (int, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		res, err := fn(s1)
		if err != nil {
			return wrapError(err), nil
		}
		return objects.NewInt(int64(res)), nil
	}
}

// FuncASRYE wraps a function taking a string input and returning bytes and an error into a CallableFunc interface.
// It validates arguments, ensures results do not exceed MaxBytesLen, and handles errors, returning them as
func FuncASRYE(fn func(string) ([]byte, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		s1, ok := objects.ToString(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		}
		res, err := fn(s1)
		if err != nil {
			return wrapError(err), nil
		}
		if len(res) > objects.MaxBytesLen {
			return nil, errors.ErrBytesLimit
		}
		return objects.NewBytes(res), nil
	}
}

// FuncAIRSsE wraps a function that takes an integer and returns a slice of strings and an error into a CallableFunc type.
func FuncAIRSsE(fn func(int) ([]string, error)) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		res, err := fn(i1)
		if err != nil {
			return wrapError(err), nil
		}
		arr := objects.NewArray(nil)
		for _, r := range res {
			if len(r) > objects.MaxStringLen {
				return nil, errors.ErrStringLimit
			}
			v, err := objects.NewString(r)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncAIRS wraps a function accepting an integer and returning a string, exposing it as an objects.CallableFunc.
// It validates argument count and types, ensuring correct usage and error handling for runtime invocation.
func FuncAIRS(fn func(int) string) objects.CallableFunc {
	return func(args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, errors.ErrWrongNumArguments
		}
		i1, ok := objects.ToInt(args[0])
		if !ok {
			return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
		}
		s := fn(i1)
		v, err := objects.NewString(s)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}
