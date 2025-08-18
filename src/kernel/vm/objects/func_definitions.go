package objects

// Code		Meaning		Go-Type			Description
//	I		Input			-			Prefix indicating input parameters.
//	O		Output			-			Prefix indicating return values.
//	n		None			()			Indicates absence of parameters or return values.
//	i		int				int			An integer.
//	b		bool			bool		A boolean value.
//	e		error			error		An error object.
//	i64		int64			int64		A 64-bit integer.
//	f64		float64			float64		A 64-bit floating point number.
//	s		string			string		A string.
//	sS		string-Slice 	[]string	A slice (array) of strings.
//	bS		bytes-Slice		[]byte		A slice (array) of bytes.
//  iS      int-Slice		[]int       A slice (array) of int.

// FuncInOn converts a no-argument, no-return Go function into a FuncCallable type that can be called with zero arguments.
// Returns ErrWrongNumArguments if any arguments are passed.
// Invokes the provided function and returns UndefinedValue upon successful execution.
func FuncInOn(fn func()) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		fn()
		return UndefinedValue, nil
	}
}

// FuncInOi wraps a no-argument integer-returning function into a callable functional interface of type FuncCallable.
// Returns an error if arguments are provided. Converts the integer result into an IObject using NewInt.
func FuncInOi(fn func() int) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewInt(int64(fn())), nil
	}
}

// FuncInOi64 wraps a function returning int64 into a FuncCallable with no arguments.
// Returns ErrWrongNumArguments if arguments are passed.
// Converts the result to an IObject using NewInt before returning.
func FuncInOi64(fn func() int64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewInt(fn()), nil
	}
}

// FuncIi64Oi64 wraps a function that takes int64 and returns int64, into a FuncCallable compatible with IObject interface.
func FuncIi64Oi64(fn func(int64) int64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return NewInt(fn(i1)), nil
	}
}

// FuncIi64On wraps a function that accepts a single int64 argument into a FuncCallable that works with IObject arguments.
func FuncIi64On(fn func(int64)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(i1)
		return UndefinedValue, nil
	}
}

// FuncInOb wraps a zero-argument boolean function into a FuncCallable that returns TrueValue or FalseValue.
func FuncInOb(fn func() bool) FuncCallable {
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

// FuncInOe creates a FuncCallable wrapper around a zero-argument function that returns an error.
// Returns ErrWrongNumArguments if arguments are provided.
// Wraps the error returned by the given function into an IObject-compatible error object.
func FuncInOe(fn func() error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewObjectError(fn()), nil
	}
}

// FuncInOs wraps a function that returns a string, creating a FuncCallable with IObject arguments and results.
// If called with arguments, it returns ErrWrongNumArguments. Otherwise, it returns a string-wrapped IObject result.
func FuncInOs(fn func() string) FuncCallable {
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

// FuncInOse wraps a function that returns a string and error into a FuncCallable that accepts no arguments.
// Returns an error if arguments are provided or if the wrapped function encounters an error.
func FuncInOse(fn func() (string, error)) FuncCallable {
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

// FuncInObSe converts a function returning ([]byte, error) into a FuncCallable that adheres to IObject function standards.
// It ensures the argument count is zero, wraps errors into IObject-compatible errors, and enforces byte slice size limits.
func FuncInObSe(fn func() ([]byte, error)) FuncCallable {
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

// FuncInOf64 wraps a zero-argument function that returns a float64 into a FuncCallable returning an IObject and an error.
// Returns ErrWrongNumArguments if called with arguments.
// Converts the float64 output of the provided function into an IObject using NewFloat.
func FuncInOf64(fn func() float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return NewFloat(fn()), nil
	}
}

// FuncInOsS takes a function that returns a slice of strings and wraps it into a FuncCallable returning an Array of strings.
// The FuncCallable expects zero arguments; passing others results in ErrWrongNumArguments.
// Converts each string from the slice into a String object and appends it to the Array.
func FuncInOsS(fn func() []string) FuncCallable {
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

// FuncInOiSe wraps a function that returns a slice of integers and an error into a FuncCallable compatible function.
// It validates zero arguments, invokes the wrapped function, wraps any error, and converts the slice to an array of IObject.
// Returns an IObject array containing the integers or a wrapped error if the wrapped function fails.
func FuncInOiSe(fn func() ([]int, error)) FuncCallable {
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

// FuncIiOiS takes a function that converts an integer to a slice of integers and returns it as a callable function.
func FuncIiOiS(fn func(int) []int) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(int(i1))
		arr := NewArray(nil)
		for _, v := range res {
			arr.Append(NewInt(int64(v)))
		}
		return arr, nil
	}
}

// FuncIf64Of64 converts a single-argument float64 function into a FuncCallable compatible with IObject arguments.
// It validates the input argument as a float-compatible type.
// Returns a new IObject representing the result or an appropriate error if validation fails.
func FuncIf64Of64(fn func(float64) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return NewFloat(fn(f1)), nil
	}
}

// FuncIiOn wraps a function with an int parameter to conform to the FuncCallable signature for custom runtime calls.
// It validates the argument count and type, invoking the provided function with the argument as an integer.
// Returns UndefinedValue on success or an error if the argument is invalid.
func FuncIiOn(fn func(int)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(int(i1))
		return UndefinedValue, nil
	}
}

// FuncIiOf64 wraps a function of type func(int) float64 as a FuncCallable, enabling its use within the IObject interface ecosystem.
// It validates that exactly one argument is provided and converts it to an int before calling the wrapped function.
// If the argument type is incompatible or the wrong number of arguments are passed, an appropriate error is returned.
func FuncIiOf64(fn func(int) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return NewFloat(fn(int(i1))), nil
	}
}

// FuncIf64Oi wraps a function transforming a float64 to an int, making it callable with IObject arguments.
// Returns an error if incorrect number or type of arguments are provided.
func FuncIf64Oi(fn func(float64) int) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return NewInt(int64(fn(f1))), nil
	}
}

// FuncIf64f64Of64 creates a FuncCallable that applies the given binary float64 function to two converted IObject arguments.
// Returns an error if arguments are not exactly two or cannot be converted to float64.
func FuncIf64f64Of64(fn func(float64, float64) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewFloat(fn(f1, f2)), nil
	}
}

// FuncIif64Of64 wraps a provided function accepting an int and float64, returning it as a FuncCallable compatible with IObject arguments.
// It enforces argument type validation and handles potential type mismatches with descriptive errors.
func FuncIif64Of64(fn func(int, float64) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewFloat(fn(int(i1), f2)), nil
	}
}

// FuncIf64iOf64 creates a FuncCallable wrapping a function that takes a float64 and int and returns a float64.
// It validates input argument types and converts them to the appropriate types expected by the wrapped function.
// Returns an IObject representing the result of the wrapped function or an error if argument validation fails.
func FuncIf64iOf64(fn func(float64, int) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewFloat(fn(f1, int(i2))), nil
	}
}

// FuncIf64iOb wraps a function that processes a float64 and an int, exposing it as a FuncCallable compatible with the IObject interface.
// It converts the first argument to a float64 and the second to an int, then applies the provided function.
// Returns TrueValue if the function evaluates to true; otherwise, returns FalseValue.
// Returns ErrWrongNumArguments if the argument count is not 2 or NewInvalidArgumentError on type conversion failures.
func FuncIf64iOb(fn func(float64, int) bool) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(f1, int(i2)) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncIf64Ob wraps a function accepting a float64 and returning a boolean into a FuncCallable compatible with the IObject interface.
func FuncIf64Ob(fn func(float64) bool) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		if fn(f1) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncIsOs creates a FuncCallable that applies a provided string-to-string function to the first argument and returns the result.
func FuncIsOs(fn func(string) string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v, err := NewString(fn(s1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsOsS converts a string-to-string-array function into a FuncCallable that operates on IObject arguments.
// It takes one string-compatible argument, applies the provided function, and returns the result as an Array of strings.
// If argument count or type is invalid, it returns an error.
func FuncIsOsS(fn func(string) []string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
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

// FuncIsOse wraps a string transformation function and adapts it to a FuncCallable with argument validation logic.
func FuncIsOse(fn func(string) (string, error)) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
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

// FuncIsOe converts a string-to-error function into a FuncCallable that operates on IObject arguments.
// It expects exactly one argument convertible to a string and returns an IObject error or result.
// Returns ErrWrongNumArguments if called with an incorrect number of arguments.
// Returns an invalid argument error if the first argument is not string-compatible.
func FuncIsOe(fn func(string) error) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		return NewObjectError(fn(s1)), nil
	}
}

// FuncIssOe wraps a function accepting two strings and returning an error into a FuncCallable compatible with the IObject interface.
// It ensures the function is called with exactly two string arguments and returns an appropriate error for incorrect usage.
func FuncIssOe(fn func(string, string) error) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewObjectError(fn(s1, s2)), nil
	}
}

// FuncIssOsS converts a function that takes two strings and returns a slice of strings into a FuncCallable.
// The returned FuncCallable validates its arguments, invokes the provided function, and returns the results as an array.
func FuncIssOsS(fn func(string, string) []string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
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

// FuncIssiOsS converts a function with parameters (string, string, int) -> []string into a FuncCallable.
// It validates arguments, applies the function, and wraps the output in an IObject-compatible Array.
// Returns an error if argument validation fails or function results cannot be converted to a String.
func FuncIssiOsS(fn func(string, string, int) []string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		arr := NewArray(nil)
		for _, res := range fn(s1, s2, int(i3)) {
			v, err := NewString(res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIssOi converts a function with two string inputs and an int output into a FuncCallable type.
// The returned FuncCallable validates that exactly two arguments are passed and they are string-compatible.
// If arguments are valid, the wrapped function is invoked, and its integer result is wrapped in an IObject.
// Returns an error if the number of arguments is incorrect or conversion to strings fails.
func FuncIssOi(fn func(string, string) int) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewInt(int64(fn(s1, s2))), nil
	}
}

// FuncIssOs wraps a function that takes two strings and returns a string into a FuncCallable accepting IObject arguments.
// It validates argument types and ensures exactly two arguments are passed or returns an appropriate error.
// The wrapped function's result is converted to an IObject before being returned.
func FuncIssOs(fn func(string, string) string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := NewString(fn(s1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil

	}
}

// FuncIssOb wraps a binary comparison function for strings as a callable function in the IObject system.
// The returned FuncCallable validates arguments, applies the provided function, and returns TrueValue or FalseValue.
// It expects the function to take two string arguments and return a boolean indicating the comparison result.
// Returns an error if the number of arguments is incorrect or arguments are not string-compatible.
func FuncIssOb(fn func(string, string) bool) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(s1, s2) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
}

// FuncIsSsOs creates a FuncCallable that processes a string slice and a string, applying the given transformation function.
func FuncIsSsOs(fn func([]string, string) string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		var ss1 []string
		switch arg0 := args[0].(type) {
		case *Array:
			for idx, a := range arg0.Values() {
				as, err := ToStringArg(idx, a)
				if err != nil {
					return nil, err
				}
				ss1 = append(ss1, as)
			}
		case *ArrayImmutable:
			for idx, a := range arg0.Values() {
				as, err := ToStringArg(idx, a)
				if err != nil {
					return nil, err
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, NewInvalidArgumentError(0, "array", args[0].TypeName())
		}
		s2, err := ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := NewString(fn(ss1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsi64Oe transforms a function accepting a string and int64 into a FuncCallable that operates on IObject arguments.
// Takes exactly two arguments; the first must be string-compatible, the second int64-compatible, or errors are returned.
// Wraps the result of the provided function into an IObject or returns an appropriate error if validation fails.
func FuncIsi64Oe(fn func(string, int64) error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewObjectError(fn(s1, i2)), nil
	}
}

// FuncIiiOe wraps a function taking two integers and returning an error into a FuncCallable accepting two IObject arguments.
func FuncIiiOe(fn func(int, int) error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return NewObjectError(fn(int(i1), int(i2))), nil
	}
}

// FuncIsiOs wraps a function that takes a string and int as inputs and returns a string, converting it to a FuncCallable.
// It validates the arguments, calls the wrapped function, and converts the result to an IObject.
// Returns an error if arguments are of invalid types or wrong number of arguments is supplied.
func FuncIsiOs(fn func(string, int) string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := NewString(fn(s1, int(i2)))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsiiOe converts a function with string, int, int inputs, and an error return into a FuncCallable with variadic IObject arguments.
func FuncIsiiOe(fn func(string, int, int) error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		return NewObjectError(fn(s1, int(i2), int(i3))), nil
	}
}

// FuncIbSOie wraps a function that takes a byte slice and returns an int and error into a FuncCallable for IObject use.
// It ensures the input argument is a single byte-compatible IObject and converts its result to IObject format.
// Returns ErrWrongNumArguments if called with more or less than one argument.
// Returns NewInvalidArgumentError if the input argument isn't byte-compatible.
// Converts the function's error output into an appropriate IObject error.
func FuncIbSOie(fn func([]byte) (int, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		bs1, err := ToByteSliceArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(bs1)
		if err != nil {
			return NewObjectError(err), nil
		}
		return NewInt(int64(res)), nil
	}
}

// FuncIbSOs wraps a function that converts a byte slice to a string, returning it as a FuncCallable in the custom object system.
// It ensures the input is a single argument of type bytes-compatible, and returns an error for invalid or unsupported types.
// The resulting FuncCallable checks argument validity, applies the provided function, and returns a new String object.
func FuncIbSOs(fn func([]byte) string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		bs1, err := ToByteSliceArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v, err := NewString(fn(bs1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsOie wraps a string-to-int function into a FuncCallable compatible with IObject interface arguments and error handling.
func FuncIsOie(fn func(string) (int, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return NewObjectError(err), nil
		}
		return NewInt(int64(res)), nil
	}
}

// FuncIsObSe returns a FuncCallable that wraps a function converting a string to a byte slice and error output.
// It validates input, reports invalid arguments, enforces byte length limits, and converts output to IObject format.
// Uses ErrWrongNumArguments, NewInvalidArgumentError, and ErrBytesLimit for error handling.
func FuncIsObSe(fn func(string) ([]byte, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := ToStringArg(0, args[0])
		if err != nil {
			return nil, err
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

// FuncIiOsSe converts a function mapping an integer to a slice of strings and an error into a FuncCallable.
func FuncIiOsSe(fn func(int) ([]string, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(int(i1))
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

// FuncIiOs wraps a function of type `func(int) string` into a FuncCallable compatible with the IObject interface system.
// It validates argument count and type, executes the provided function, and converts the result into an IObject.
func FuncIiOs(fn func(int) string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		s := fn(int(i1))
		v, err := NewString(s)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}
