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

// funcInOn converts a no-argument void Go function into a FuncCallable for integration with the IObject system.
// Returns ErrWrongNumArguments if arguments are passed during invocation.
// Executes the provided function and returns the UndefinedValue from GateKeeper upon completion.
func funcInOn(f *GateKeeper, fn func()) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		fn()
		return f.UndefinedValue(), nil
	}
}

// funcInOi wraps a no-argument integer-returning function into a FuncCallable interface, returning an IObject representation of the result.
// It requires a GateKeeper for object creation and fails with ErrWrongNumArguments if any arguments are provided.
func funcInOi(f *GateKeeper, fn func() int) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return f.NewInt(frame, int64(fn())), nil
	}
}

// funcInOi64 wraps a function returning int64 into a FuncCallable accepting no arguments.
// Returns ErrWrongNumArguments if arguments are provided to the invocation.
// Converts the return value of the function to an IObject using GateKeeper's NewInt method.
func funcInOi64(f *GateKeeper, fn func() int64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return f.NewInt(frame, fn()), nil
	}
}

// funcIi64Oi64 wraps a function taking int64 and returning int64 into a FuncCallable compatible with IObject interface.
// It validates a single IObject argument, converts it to an int64, applies the wrapped function, and returns the result.
// Returns ErrWrongNumArguments if provided argument count is not 1 or an error on argument conversion.
func funcIi64Oi64(f *GateKeeper, fn func(int64) int64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewInt(frame, fn(i1)), nil
	}
}

// funcIi64On wraps a function taking an int64 as a FuncCallable that operates on IObject arguments.
// It validates the argument count and type before invoking the provided function.
// Returns an IObject output or an error.
func funcIi64On(f *GateKeeper, fn func(int64)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(i1)
		return f.UndefinedValue(), nil
	}
}

// funcInOb transforms a zero-argument boolean function into a FuncCallable returning TrueValue or FalseValue based on its result.
// Returns an error if extra arguments are passed during invocation.
func funcInOb(f *GateKeeper, fn func() bool) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		if fn() {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// funcInOe wraps a zero-argument function returning an error into a FuncCallable compatible function.
// Returns an error object if the wrapped function fails.
// Returns ErrWrongNumArguments if called with any arguments.
func funcInOe(f *GateKeeper, fn func() error) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		err := fn()
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// funcInOs creates a FuncCallable that wraps a function returning a string, converting it to an IObject within a given frame.
// Returns ErrWrongNumArguments if called with any arguments, otherwise wraps and returns the function's string output as an IObject.
func funcInOs(f *GateKeeper, fn func() string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		v, err := f.NewString(frame, fn())
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcInOse wraps a function returning a string and error into a FuncCallable that takes no arguments and handles errors.
func funcInOse(f *GateKeeper, fn func() (string, error)) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		v, err := f.NewString(frame, res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcInObSe converts a function returning ([]byte, error) into a FuncCallable that wraps logic for IObject compatibility.
// Ensures no arguments are passed, wraps errors into IObject errors, and enforces a byte slice size limit.
func funcInObSe(f *GateKeeper, fn func() ([]byte, error)) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		if len(res) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return f.NewBytes(frame, res), nil
	}
}

// funcInOf64 wraps a zero-argument function returning float64 into a FuncCallable returning an IObject and an error.
// Returns ErrWrongNumArguments if called with arguments.
// Converts the float64 result of fn into an IObject using the provided GateKeeper instance.
func funcInOf64(f *GateKeeper, fn func() float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return f.NewFloat(frame, fn()), nil
	}
}

// funcInOsS wraps a function that returns a slice of strings into a FuncCallable, returning an Array object of the strings.
// The resulting FuncCallable accepts no arguments. Passing arguments results in ErrWrongNumArguments.
// Each string is converted to an IObject representation and added to the Array.
func funcInOsS(f *GateKeeper, fn func() []string) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		arr := f.NewArray(frame, nil)
		for _, elem := range fn() {
			v, err := f.NewString(frame, elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// funcInOiSe wraps a function returning a slice of integers and an error into a FuncCallable-compatible function.
// Validates no arguments, calls the function, handles errors, and converts the result into an IObject array.
// Returns an IObject array of integers or a wrapped error if the function execution fails.
func funcInOiSe(f *GateKeeper, fn func() ([]int, error)) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		arr := f.NewArray(frame, nil)
		for _, v := range res {
			arr.Append(f.NewInt(frame, int64(v)))
		}
		return arr, nil
	}
}

// funcIiOiS converts a function from int to a slice of int into a FuncCallable for use in the system's callable context.
func funcIiOiS(f *GateKeeper, fn func(int) []int) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(int(i1))
		arr := f.NewArray(frame, nil)
		for _, v := range res {
			arr.Append(f.NewInt(frame, int64(v)))
		}
		return arr, nil
	}
}

// funcIf64Of64 wraps a float64-to-float64 function as a FuncCallable compatible with IObject arguments.
// Converts the first IObject argument to float64, applies the function, and returns the result as an IObject.
// Returns ErrWrongNumArguments if the arguments length is not exactly 1 or an error on conversion failure.
func funcIf64Of64(f *GateKeeper, fn func(float64) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(frame, fn(f1)), nil
	}
}

// funcIiOn wraps a function with an int parameter into a FuncCallable, ensuring proper argument validation and conversion.
func funcIiOn(f *GateKeeper, fn func(int)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(int(i1))
		return f.UndefinedValue(), nil
	}
}

// funcIiOf64 wraps a func(int) float64 as a FuncCallable, validating input and converting arguments to match the function's signature.
// It ensures exactly one argument is passed, converts it to int64, and invokes the provided function, returning a float result.
// Returns an error if argument count or type is invalid.
func funcIiOf64(f *GateKeeper, fn func(int) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(frame, fn(int(i1))), nil
	}
}

// funcIf64Oi converts a float64-to-int function into a FuncCallable compatible with the IObject interface.
// Returns ErrWrongNumArguments if provided arguments don't match the expected count or type.
// Leverages GateKeeper for argument parsing and construction of the return value.
func funcIf64Oi(f *GateKeeper, fn func(float64) int) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewInt(frame, int64(fn(f1))), nil
	}
}

// funcIf64f64Of64 creates a FuncCallable that applies a binary float64 function to two IObject arguments after conversion.
// Returns an error if the argument count is not exactly two or if conversion to float64 fails.
func funcIf64f64Of64(f *GateKeeper, fn func(float64, float64) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := f.ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(frame, fn(f1, f2)), nil
	}
}

// funcIif64Of64 wraps a function accepting an int and float64, exposing it as a FuncCallable compatible with IObject arguments.
// It validates the argument types, ensuring the first is convertible to int64 and the second to float64, or returns appropriate errors.
func funcIif64Of64(f *GateKeeper, fn func(int, float64) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := f.ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(frame, fn(int(i1), f2)), nil
	}
}

// funcIf64iOf64 creates a FuncCallable wrapping a function that takes float64, int arguments and returns a float64.
// Validates and converts input arguments to expected types, returning the result or an error on validation failure.
func funcIf64iOf64(f *GateKeeper, fn func(float64, int) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(frame, fn(f1, int(i2))), nil
	}
}

// funcIf64iOb converts the first argument to float64 and the second to int, then applies the provided boolean function.
// Returns TrueValue/FalseValue based on the function's result or an error if arguments are invalid or conversion fails.
// Returns ErrWrongNumArguments if the number of arguments is not 2.
func funcIf64iOb(f *GateKeeper, fn func(float64, int) bool) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(f1, int(i2)) {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// funcIf64Ob converts a function that takes a float64 and returns a boolean into a FuncCallable-compatible wrapper.
// Returns `true` or `false` IObject based on the function's result.
// Ensures the input arguments are valid and converts the first argument to float64 or raises an error.
// Errors if an incorrect number of arguments is provided.
func funcIf64Ob(f *GateKeeper, fn func(float64) bool) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		if fn(f1) {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// funcIsOs creates a FuncCallable that applies a string-to-string function to the first argument and returns the result.
// It ensures the correct number of arguments and converts the input to a string before processing.
// Returns the transformed string as an IObject or an error if the input is invalid.
func funcIsOs(f *GateKeeper, fn func(string) string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(frame, fn(s1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcIsOsS creates a FuncCallable that converts IObject arguments to strings, applies a string-to-string-array function, and returns the result as an Array.
// It validates the argument count and type, returning an error for invalid inputs.
func funcIsOsS(f *GateKeeper, fn func(string) []string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(s1)
		arr := f.NewArray(frame, nil)
		for _, elem := range res {
			v, err := f.NewString(frame, elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// funcIsOse adapts a string transformation function to a FuncCallable, validating arguments and managing errors.
func funcIsOse(f *GateKeeper, fn func(string) (string, error)) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		v, err := f.NewString(frame, res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcIsOe wraps a string-to-error function into a FuncCallable compatible with the system's IObject interface.
// It validates the argument count, ensures the first argument is string-convertible, and propagates any errors.
// Returns ErrWrongNumArguments if the number of arguments is invalid or an IObject error for argument issues.
// On success, it returns f.TrueValue() indicating successful execution of the given function.
func funcIsOe(f *GateKeeper, fn func(string) error) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		err = fn(s1)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// funcIssOe wraps a function taking two strings and returning an error into a FuncCallable compatible with IObject.
// Ensures the function is invoked with exactly two string arguments, returning an error for invalid input.
// Returns a TrueValue IObject upon successful execution or an error message as a new error object.
func funcIssOe(f *GateKeeper, fn func(string, string) error) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		err = fn(s1, s2)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// funcIssOsS converts a function accepting two strings and returning a slice of strings into a FuncCallable.
// It validates arguments, invokes the function, and returns the results as an IObject Array.
func funcIssOsS(f *GateKeeper, fn func(string, string) []string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		arr := f.NewArray(frame, nil)
		for _, res := range fn(s1, s2) {
			v, err := f.NewString(frame, res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// funcIssiOsS converts a function with parameters (string, string, int) -> []string into a FuncCallable with type validation.
// It processes arguments, executes the provided function, and wraps the output into an IObject-compatible array.
// Returns an error if argument count or types are invalid, or if conversion to IObject fails.
func funcIssiOsS(f *GateKeeper, fn func(string, string, int) []string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := f.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		arr := f.NewArray(frame, nil)
		for _, res := range fn(s1, s2, int(i3)) {
			v, err := f.NewString(frame, res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// funcIssOi converts a function with two string arguments and an int return type into a FuncCallable.
// It validates that exactly two arguments are passed, converts them to strings, and wraps the result in an IObject.
// Returns an error if the argument count is incorrect or conversion to strings fails.
func funcIssOi(f *GateKeeper, fn func(string, string) int) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewInt(frame, int64(fn(s1, s2))), nil
	}
}

// funcIssOs wraps a two-string-argument function into a FuncCallable for arguments of type IObject.
// It ensures exactly two arguments are provided and converts them into strings or returns an error if conversion fails.
// The wrapped function's result is converted back into an IObject before being returned.
func funcIssOs(f *GateKeeper, fn func(string, string) string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(frame, fn(s1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil

	}
}

// funcIssOb wraps a binary string comparison function as a FuncCallable for the IObject system with validation and error handling.
// It ensures two arguments are provided, converts them to strings, applies the comparison function, and returns TrueValue or FalseValue.
// An error is returned if argument count is invalid or type conversion fails.
func funcIssOb(f *GateKeeper, fn func(string, string) bool) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(s1, s2) {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// funcIsSsOs constructs a FuncCallable that processes arguments as an array of strings and a single string, returning a new string.
// It validates input types, transforms them, applies a provided function, and creates a new result object or error.
func funcIsSsOs(f *GateKeeper, fn func([]string, string) string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		var ss1 []string
		switch arg0 := args[0].(type) {
		case *Array:
			for idx, a := range arg0.Values() {
				as, err := f.ToStringArg(idx, a)
				if err != nil {
					return nil, err
				}
				ss1 = append(ss1, as)
			}
		case *ArrayImmutable:
			for idx, a := range arg0.Values() {
				as, err := f.ToStringArg(idx, a)
				if err != nil {
					return nil, err
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, NewInvalidArgumentError(0, "array", args[0].TypeName())
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(frame, fn(ss1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcIsi64Oe converts a function accepting string and int64 into a FuncCallable operating on IObject arguments.
// Validates arguments count, ensures the first is string-compatible and second int64-compatible, then invokes the function.
// Returns the result as IObject or an error in case of validation or function execution failure.
func funcIsi64Oe(f *GateKeeper, fn func(string, int64) error) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		err = fn(s1, i2)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// funcIiiOe converts a function with two integer inputs and an error return to a FuncCallable that operates on IObject inputs.
func funcIiiOe(f *GateKeeper, fn func(int, int) error) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		err = fn(int(i1), int(i2))
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// funcIsiOs converts a function with a string and int as inputs to a FuncCallable, ensuring argument validation and type conversion.
func funcIsiOs(f *GateKeeper, fn func(string, int) string) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(frame, fn(s1, int(i2)))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcIsiiOe converts a function with string, int, int inputs and an error return into a FuncCallable with IObject arguments.
// It validates the number of arguments and converts them to the appropriate types before invoking the provided function.
func funcIsiiOe(f *GateKeeper, fn func(string, int, int) error) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := f.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		err = fn(s1, int(i2), int(i3))
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// funcIbSOie wraps a function into a FuncCallable compatible with IObject, handling argument validation and conversion.
// Ensures the function accepts a single byte-compatible argument and converts the result to an IObject type.
// Returns ErrWrongNumArguments if the wrong number of arguments are provided.
// Converts the wrapped function's error into an IObject error using the GateKeeper's error handling.
func funcIbSOie(f *GateKeeper, fn func([]byte) (int, error)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		bs1, err := f.ToByteSliceArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(bs1)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.NewInt(frame, int64(res)), nil
	}
}

// funcIbSOs creates a FuncCallable that transforms a byte slice into a string using a provided function and handles argument validation.
func funcIbSOs(f *GateKeeper, fn func([]byte) string) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		bs1, err := f.ToByteSliceArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(frame, fn(bs1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// funcIsOie converts a string-to-int function into a FuncCallable that handles IObject arguments and errors.
// It validates argument length, converts the first argument to a string, and processes the function result.
// Returns an IObject-based integer result or an error encapsulated in IObject.
func funcIsOie(f *GateKeeper, fn func(string) (int, error)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		return f.NewInt(frame, int64(res)), nil
	}
}

// funcIsObSe constructs a FuncCallable that wraps a string-to-byte transformation while enforcing argument and byte limits.
func funcIsObSe(f *GateKeeper, fn func(string) ([]byte, error)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		if len(res) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return f.NewBytes(frame, res), nil
	}
}

// funcIiOsSe wraps a function of type func(int) ([]string, error) into a callable function with error handling and array conversion.
func funcIiOsSe(f *GateKeeper, fn func(int) ([]string, error)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(int(i1))
		if err != nil {
			return f.NewError(frame, err.Error()), nil
		}
		arr := f.NewArray(frame, nil)
		for _, r := range res {
			if len(r) > MaxStringLen {
				return nil, ErrStringLimit
			}
			v, err := f.NewString(frame, r)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// funcIiOs wraps a `func(int) string` as a FuncCallable compatible with IObject, validating arguments and managing errors.
func funcIiOs(f *GateKeeper, fn func(int) string) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		s := fn(int(i1))
		v, err := f.NewString(frame, s)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}
