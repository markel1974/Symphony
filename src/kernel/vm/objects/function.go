package objects

import (
	"fmt"
)

// FuncAR converts a no-argument, no-return Go function into a CallableFunc type that can be called with zero arguments.
// Returns ErrWrongNumArguments if any arguments are passed.
// Invokes the provided function and returns UndefinedValue upon successful execution.
func FuncAR(fn func()) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		fn()
		return UndefinedValue, nil
	}
}

// FuncARI wraps a no-argument integer-returning function into a callable functional interface of type CallableFunc.
// Returns an error if arguments are provided. Converts the integer result into an IObject using NewInt.
func FuncARI(fn func() int) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewInt(int64(fn())), nil
	}
}

// FuncARI64 wraps a function returning int64 into a CallableFunc with no arguments.
// Returns ErrWrongNumArguments if arguments are passed.
// Converts the result to an IObject using NewInt before returning.
func FuncARI64(fn func() int64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewInt(fn()), nil
	}
}

// FuncAI64RI64 wraps a function that takes int64 and returns int64, into a CallableFunc compatible with IObject interface.
func FuncAI64RI64(fn func(int64) int64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}

		i1, ok := ToInt64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		return NewInt(fn(i1)), nil
	}
}

// FuncAI64R wraps a function that accepts a single int64 argument into a CallableFunc that works with IObject arguments.
func FuncAI64R(fn func(int64)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}

		i1, ok := ToInt64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		fn(i1)
		return UndefinedValue, nil
	}
}

// FuncARB wraps a zero-argument boolean function into a CallableFunc that returns TrueValue or FalseValue.
func FuncARB(fn func() bool) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		if fn() {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncARE creates a CallableFunc wrapper around a zero-argument function that returns an error.
// Returns ErrWrongNumArguments if arguments are provided.
// Wraps the error returned by the given function into an IObject-compatible error object.
func FuncARE(fn func() error) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewObjectError(fn()), nil
	}
}

// FuncARS wraps a function that returns a string, creating a CallableFunc with IObject arguments and results.
// If called with arguments, it returns ErrWrongNumArguments. Otherwise, it returns a string-wrapped IObject result.
func FuncARS(fn func() string) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		v, err := NewString(fn())
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncARSE wraps a function that returns a string and error into a CallableFunc that accepts no arguments.
// Returns an error if arguments are provided or if the wrapped function encounters an error.
func FuncARSE(fn func() (string, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return NewObjectError(err), nil
		}
		v, err := NewString(res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncARYE converts a function returning ([]byte, error) into a CallableFunc that adheres to IObject function standards.
// It ensures the argument count is zero, wraps errors into IObject-compatible errors, and enforces byte slice size limits.
func FuncARYE(fn func() ([]byte, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return NewObjectError(err), nil
		}
		if len(res) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return NewBytes(res), nil
	}
}

// FuncARF wraps a zero-argument function that returns a float64 into a CallableFunc returning an IObject and an error.
// Returns ErrWrongNumArguments if called with arguments.
// Converts the float64 output of the provided function into an IObject using NewFloat.
func FuncARF(fn func() float64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewFloat(fn()), nil
	}
}

// FuncARSs takes a function that returns a slice of strings and wraps it into a CallableFunc returning an Array of strings.
// The CallableFunc expects zero arguments; passing others results in ErrWrongNumArguments.
// Converts each string from the slice into a String object and appends it to the Array.
func FuncARSs(fn func() []string) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		arr := &Array{}
		for _, elem := range fn() {
			v, err := NewString(elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncARIsE wraps a function that returns a slice of integers and an error into a CallableFunc compatible function.
// It validates zero arguments, invokes the wrapped function, wraps any error, and converts the slice to an array of IObject.
// Returns an IObject array containing the integers or a wrapped error if the wrapped function fails.
func FuncARIsE(fn func() ([]int, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return NewObjectError(err), nil
		}
		arr := NewArray(nil)
		for _, v := range res {
			arr.Append(NewInt(int64(v)))
		}
		return arr, nil
	}
}

// FuncAIRIs takes a function that converts an integer to a slice of integers and returns it as a callable function.
func FuncAIRIs(fn func(int) []int) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		res := fn(i1)
		arr := NewArray(nil)
		for _, v := range res {
			arr.Append(NewInt(int64(v)))
		}
		return arr, nil
	}
}

// FuncAFRF converts a single-argument float64 function into a CallableFunc compatible with IObject arguments.
// It validates the input argument as a float-compatible type.
// Returns a new IObject representing the result or an appropriate error if validation fails.
func FuncAFRF(fn func(float64) float64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, ok := ToFloat64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "float(compatible)", args[0].TypeName())
		}
		return NewFloat(fn(f1)), nil
	}
}

// FuncAIR wraps a function with an int parameter to conform to the CallableFunc signature for custom runtime calls.
// It validates the argument count and type, invoking the provided function with the argument as an integer.
// Returns UndefinedValue on success or an error if the argument is invalid.
func FuncAIR(fn func(int)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		fn(i1)
		return UndefinedValue, nil
	}
}

// FuncAIRF wraps a function of type func(int) float64 as a CallableFunc, enabling its use within the IObject interface ecosystem.
// It validates that exactly one argument is provided and converts it to an int before calling the wrapped function.
// If the argument type is incompatible or the wrong number of arguments are passed, an appropriate error is returned.
func FuncAIRF(fn func(int) float64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		return NewFloat(fn(i1)), nil
	}
}

// FuncAFRI wraps a function transforming a float64 to an int, making it callable with IObject arguments.
// Returns an error if incorrect number or type of arguments are provided.
func FuncAFRI(fn func(float64) int) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, ok := ToFloat64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "float(compatible)", args[0].TypeName())
		}
		return NewInt(int64(fn(f1))), nil
	}
}

// FuncAFFRF creates a CallableFunc that applies the given binary float64 function to two converted IObject arguments.
// Returns an error if arguments are not exactly two or cannot be converted to float64.
func FuncAFFRF(fn func(float64, float64) float64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, ok := ToFloat64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "float(compatible)", args[0].TypeName())
		}
		f2, ok := ToFloat64(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "float(compatible)", args[1].TypeName())
		}
		return NewFloat(fn(f1, f2)), nil
	}
}

// FuncAIFRF wraps a provided function accepting an int and float64, returning it as a CallableFunc compatible with IObject arguments.
// It enforces argument type validation and handles potential type mismatches with descriptive errors.
func FuncAIFRF(fn func(int, float64) float64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		f2, ok := ToFloat64(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "float(compatible)", args[1].TypeName())
		}
		return NewFloat(fn(i1, f2)), nil
	}
}

// FuncAFIRF creates a CallableFunc wrapping a function that takes a float64 and int and returns a float64.
// It validates input argument types and converts them to the appropriate types expected by the wrapped function.
// Returns an IObject representing the result of the wrapped function or an error if argument validation fails.
func FuncAFIRF(fn func(float64, int) float64) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, ok := ToFloat64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "float(compatible)", args[0].TypeName())
		}
		i2, ok := ToInt(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
		}
		return NewFloat(fn(f1, i2)), nil
	}
}

// FuncAFIRB wraps a function that processes a float64 and an int, exposing it as a CallableFunc compatible with the IObject interface.
// It converts the first argument to a float64 and the second to an int, then applies the provided function.
// Returns TrueValue if the function evaluates to true; otherwise, returns FalseValue.
// Returns ErrWrongNumArguments if the argument count is not 2 or NewInvalidArgumentError on type conversion failures.
func FuncAFIRB(fn func(float64, int) bool) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, ok := ToFloat64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "float(compatible)", args[0].TypeName())
		}
		i2, ok := ToInt(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
		}
		if fn(f1, i2) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncAFRB wraps a function accepting a float64 and returning a boolean into a CallableFunc compatible with the IObject interface.
func FuncAFRB(fn func(float64) bool) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, ok := ToFloat64(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "float(compatible)", args[0].TypeName())
		}
		if fn(f1) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncASRS creates a CallableFunc that applies a provided string-to-string function to the first argument and returns the result.
func FuncASRS(fn func(string) string) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		v, err := NewString(fn(s1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASRSs converts a string-to-string-array function into a CallableFunc that operates on IObject arguments.
// It takes one string-compatible argument, applies the provided function, and returns the result as an Array of strings.
// If argument count or type is invalid, it returns an error.
func FuncASRSs(fn func(string) []string) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		res := fn(s1)
		arr := NewArray(nil)
		for _, elem := range res {
			v, err := NewString(elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncASRSE wraps a string transformation function and adapts it to a CallableFunc with argument validation logic.
func FuncASRSE(fn func(string) (string, error)) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		res, err := fn(s1)
		if err != nil {
			return NewObjectError(err), nil
		}
		v, err := NewString(res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASRE converts a string-to-error function into a CallableFunc that operates on IObject arguments.
// It expects exactly one argument convertible to a string and returns an IObject error or result.
// Returns ErrWrongNumArguments if called with an incorrect number of arguments.
// Returns an invalid argument error if the first argument is not string-compatible.
func FuncASRE(fn func(string) error) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		return NewObjectError(fn(s1)), nil
	}
}

// FuncASSRE wraps a function accepting two strings and returning an error into a CallableFunc compatible with the IObject interface.
// It ensures the function is called with exactly two string arguments and returns an appropriate error for incorrect usage.
func FuncASSRE(fn func(string, string) error) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "string(compatible)", args[1].TypeName())
		}
		return NewObjectError(fn(s1, s2)), nil
	}
}

// FuncASSRSs converts a function that takes two strings and returns a slice of strings into a CallableFunc.
// The returned CallableFunc validates its arguments, invokes the provided function, and returns the results as an array.
func FuncASSRSs(fn func(string, string) []string) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[1].TypeName())

		}
		arr := NewArray(nil)
		for _, res := range fn(s1, s2) {
			v, err := NewString(res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncASSIRSs converts a function with parameters (string, string, int) -> []string into a CallableFunc.
// It validates arguments, applies the function, and wraps the output in an IObject-compatible Array.
// Returns an error if argument validation fails or function results cannot be converted to a String.
func FuncASSIRSs(fn func(string, string, int) []string) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "string(compatible)", args[1].TypeName())
		}
		i3, ok := ToInt(args[2])
		if !ok {
			return nil, NewInvalidArgumentError("third", "int(compatible)", args[2].TypeName())
		}
		arr := NewArray(nil)
		for _, res := range fn(s1, s2, i3) {
			v, err := NewString(res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncASSRI converts a function with two string inputs and an int output into a CallableFunc type.
// The returned CallableFunc validates that exactly two arguments are passed and they are string-compatible.
// If arguments are valid, the wrapped function is invoked, and its integer result is wrapped in an IObject.
// Returns an error if the number of arguments is incorrect or conversion to strings fails.
func FuncASSRI(fn func(string, string) int) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "string(compatible)", args[0].TypeName())
		}
		return NewInt(int64(fn(s1, s2))), nil
	}
}

// FuncASSRS wraps a function that takes two strings and returns a string into a CallableFunc accepting IObject arguments.
// It validates argument types and ensures exactly two arguments are passed or returns an appropriate error.
// The wrapped function's result is converted to an IObject before being returned.
func FuncASSRS(fn func(string, string) string) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "string(compatible)", args[1].TypeName())
		}
		v, err := NewString(fn(s1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil

	}
}

// FuncASSRB wraps a binary comparison function for strings as a callable function in the IObject system.
// The returned CallableFunc validates arguments, applies the provided function, and returns TrueValue or FalseValue.
// It expects the function to take two string arguments and return a boolean indicating the comparison result.
// Returns an error if the number of arguments is incorrect or arguments are not string-compatible.
func FuncASSRB(fn func(string, string) bool) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "string(compatible)", args[1].TypeName())
		}
		if fn(s1, s2) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncASsSRS creates a CallableFunc that processes a string slice and a string, applying the given transformation function.
func FuncASsSRS(fn func([]string, string) string) CallableFunc {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		var ss1 []string
		switch arg0 := args[0].(type) {
		case *Array:
			for idx, a := range arg0.Values() {
				as, ok := ToString(a)
				if !ok {
					return nil, NewInvalidArgumentError(fmt.Sprintf("first[%d]", idx), "string(compatible)", a.TypeName())
				}
				ss1 = append(ss1, as)
			}
		case *ArrayImmutable:
			for idx, a := range arg0.Values() {
				as, ok := ToString(a)
				if !ok {
					return nil, NewInvalidArgumentError(fmt.Sprintf("first[%d]", idx), "string(compatible)", a.TypeName())
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, NewInvalidArgumentError("first", "array", args[0].TypeName())
		}
		s2, ok := ToString(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "string(compatible)", args[1].TypeName())
		}
		v, err := NewString(fn(ss1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASI64RE transforms a function accepting a string and int64 into a CallableFunc that operates on IObject arguments.
// Takes exactly two arguments; the first must be string-compatible, the second int64-compatible, or errors are returned.
// Wraps the result of the provided function into an IObject or returns an appropriate error if validation fails.
func FuncASI64RE(fn func(string, int64) error) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		i2, ok := ToInt64(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
		}
		return NewObjectError(fn(s1, i2)), nil
	}
}

// FuncAIIRE wraps a function taking two integers and returning an error into a CallableFunc accepting two IObject arguments.
func FuncAIIRE(fn func(int, int) error) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		i2, ok := ToInt(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
		}
		return NewObjectError(fn(i1, i2)), nil
	}
}

// FuncASIRS wraps a function that takes a string and int as inputs and returns a string, converting it to a CallableFunc.
// It validates the arguments, calls the wrapped function, and converts the result to an IObject.
// Returns an error if arguments are of invalid types or wrong number of arguments is supplied.
func FuncASIRS(fn func(string, int) string) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		i2, ok := ToInt(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
		}
		v, err := NewString(fn(s1, i2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASIIRE converts a function with string, int, int inputs, and an error return into a CallableFunc with variadic IObject arguments.
func FuncASIIRE(fn func(string, int, int) error) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		i2, ok := ToInt(args[1])
		if !ok {
			return nil, NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
		}
		i3, ok := ToInt(args[2])
		if !ok {
			return nil, NewInvalidArgumentError("third", "int(compatible)", args[2].TypeName())
		}
		return NewObjectError(fn(s1, i2, i3)), nil
	}
}

// FuncAYRIE wraps a function that takes a byte slice and returns an int and error into a CallableFunc for IObject use.
// It ensures the input argument is a single byte-compatible IObject and converts its result to IObject format.
// Returns ErrWrongNumArguments if called with more or less than one argument.
// Returns NewInvalidArgumentError if the input argument isn't byte-compatible.
// Converts the function's error output into an appropriate IObject error.
func FuncAYRIE(fn func([]byte) (int, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		y1, ok := ToByteSlice(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "bytes(compatible)", args[0].TypeName())
		}
		res, err := fn(y1)
		if err != nil {
			return NewObjectError(err), nil
		}
		return NewInt(int64(res)), nil
	}
}

// FuncAYRS wraps a function that converts a byte slice to a string, returning it as a CallableFunc in the custom object system.
// It ensures the input is a single argument of type bytes-compatible, and returns an error for invalid or unsupported types.
// The resulting CallableFunc checks argument validity, applies the provided function, and returns a new String object.
func FuncAYRS(fn func([]byte) string) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		y1, ok := ToByteSlice(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "bytes(compatible)", args[0].TypeName())
		}
		v, err := NewString(fn(y1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncASRIE wraps a string-to-int function into a CallableFunc compatible with IObject interface arguments and error handling.
func FuncASRIE(fn func(string) (int, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		res, err := fn(s1)
		if err != nil {
			return NewObjectError(err), nil
		}
		return NewInt(int64(res)), nil
	}
}

// FuncASRYE returns a CallableFunc that wraps a function converting a string to a byte slice and error output.
// It validates input, reports invalid arguments, enforces byte length limits, and converts output to IObject format.
// Uses ErrWrongNumArguments, NewInvalidArgumentError, and ErrBytesLimit for error handling.
func FuncASRYE(fn func(string) ([]byte, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, ok := ToString(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
		}
		res, err := fn(s1)
		if err != nil {
			return NewObjectError(err), nil
		}
		if len(res) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return NewBytes(res), nil
	}
}

// FuncAIRSsE converts a function mapping an integer to a slice of strings and an error into a CallableFunc.
func FuncAIRSsE(fn func(int) ([]string, error)) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		res, err := fn(i1)
		if err != nil {
			return NewObjectError(err), nil
		}
		arr := NewArray(nil)
		for _, r := range res {
			if len(r) > MaxStringLen {
				return nil, ErrStringLimit
			}
			v, err := NewString(r)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncAIRS wraps a function of type `func(int) string` into a CallableFunc compatible with the IObject interface system.
// It validates argument count and type, executes the provided function, and converts the result into an IObject.
func FuncAIRS(fn func(int) string) CallableFunc {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, ok := ToInt(args[0])
		if !ok {
			return nil, NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
		}
		s := fn(i1)
		v, err := NewString(s)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}
