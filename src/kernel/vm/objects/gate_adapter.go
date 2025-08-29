package objects

import (
	"fmt"
)

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

// GateAdapter is a type that wraps a GateKeeper and provides functional adapters to map Go functions to FuncCallable.
type GateAdapter struct {
	gk *GateKeeper
}

// NewGateAdapter creates and returns a new instance of GateAdapter, initialized with the provided GateKeeper.
func NewGateAdapter(gk *GateKeeper) *GateAdapter {
	return &GateAdapter{gk: gk}
}

// FuncIi64On creates a FuncCallable that invokes a function accepting an int64 argument, handling argument conversion.
func (ga *GateAdapter) FuncIi64On(fn func(int64)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		i1, err := ga.gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(i1)
		return ga.gk.UndefinedValue(), nil
	}
}

// FuncIiOiS wraps a function from int to a slice of int into a FuncCallable type for use with IObject arguments.
func (ga *GateAdapter) FuncIiOiS(fn func(int) []int) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		i1, err := ga.gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(int(i1))
		obj := ga.gk.NewArray(frame, nil)
		arr, ok := obj.(*Array)
		if !ok {
			return nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, v := range res {
			arr.Append(ga.gk.NewInt(frame, int64(v)))
		}
		return arr, nil
	}
}

// FuncIf64Of64 creates a FuncCallable that applies a function transforming a float64 input to a float64 result.
// It validates a single argument, converts it to float64, applies the function, and returns the result as an IObject.
// Returns an error for invalid arguments or if the number of arguments is not exactly one.
func (ga *GateAdapter) FuncIf64Of64(fn func(float64) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		f1, err := ga.gk.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewFloat(frame, fn(f1)), nil
	}
}

// FuncIiOf64 converts a function from int to float64 into a FuncCallable that operates on IObject arguments.
func (ga *GateAdapter) FuncIiOf64(fn func(int) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		i1, err := ga.gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewFloat(frame, fn(int(i1))), nil
	}
}

// FuncIf64Oi converts a function accepting a float64 and returning an int into a callable function for GateAdapter.
func (ga *GateAdapter) FuncIf64Oi(fn func(float64) int) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		f1, err := ga.gk.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewInt(frame, int64(fn(f1))), nil
	}
}

// FuncIf64f64Of64 adapts a function taking two float64 arguments and returning a float64 to a FuncCallable.
// It converts IObject arguments to float64, applies the input function, and returns the result as an IObject.
// Returns an error if the argument count is not 2 or if type conversion fails.
func (ga *GateAdapter) FuncIf64f64Of64(fn func(float64, float64) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		f1, err := ga.gk.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := ga.gk.ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewFloat(frame, fn(f1, f2)), nil
	}
}

// FuncIif64Of64 adapts a function accepting an int and float64 as inputs, returning a float64, into a FuncCallable.
func (ga *GateAdapter) FuncIif64Of64(fn func(int, float64) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		i1, err := ga.gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := ga.gk.ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewFloat(frame, fn(int(i1), f2)), nil
	}
}

// FuncIf64iOf64 converts a function accepting (float64, int) and returning float64 into a FuncCallable type.
// It requires exactly two arguments: the first must be convertible to float64 and the second to int.
// Returns an IObject representing the result of applying the function or an error if argument conversion fails.
func (ga *GateAdapter) FuncIf64iOf64(fn func(float64, int) float64) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		f1, err := ga.gk.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ga.gk.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewFloat(frame, fn(f1, int(i2))), nil
	}
}

// FuncIf64iOb wraps a function of type func(float64, int) bool into a FuncCallable compatible with IObject arguments.
// It converts the first argument to float64 and the second argument to int, validates their types, and applies the function.
// Returns the system-defined TrueValue or FalseValue based on the function's result, or an error on failure.
// Returns ErrInvalidArgumentsNumber if the number of arguments passed is not exactly two.
func (ga *GateAdapter) FuncIf64iOb(fn func(float64, int) bool) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		f1, err := ga.gk.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := ga.gk.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(f1, int(i2)) {
			return ga.gk.TrueValue(), nil
		}
		return ga.gk.FalseValue(), nil
	}
}

// FuncIf64Ob wraps a given float64 predicate function and returns a FuncCallable to evaluate the predicate using IObject.
// If the argument count is incorrect or conversion to float64 fails, it returns an error.
// The result of the predicate determines the boolean IObject returned: TrueValue or FalseValue.
func (ga *GateAdapter) FuncIf64Ob(fn func(float64) bool) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		f1, err := ga.gk.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		if fn(f1) {
			return ga.gk.TrueValue(), nil
		}
		return ga.gk.FalseValue(), nil
	}
}

// FuncIsOs creates a callable function using the provided transformation function that operates on string arguments.
// Returns a function that takes a frame and varargs, validates the input, applies the transformation, and returns the result.
func (ga *GateAdapter) FuncIsOs(fn func(string) string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v := ga.gk.NewString(frame, fn(s1))
		return v, nil
	}
}

// FuncIsOsS transforms a string-to-string-slice function into a FuncCallable that operates on an IObject and returns an array.
func (ga *GateAdapter) FuncIsOsS(fn func(string) []string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(s1)
		obj := ga.gk.NewArray(frame, nil)
		arr, ok := obj.(*Array)
		if !ok {
			return nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, elem := range res {
			v := ga.gk.NewString(frame, elem)
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIsOse converts a function with a string input and output into a FuncCallable for use in the current framework.
// It validates arguments, converts inputs/outputs to/from IObject, and handles errors gracefully.
func (ga *GateAdapter) FuncIsOse(fn func(string) (string, error)) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return ga.gk.NewError(frame, err.Error()), nil
		}
		v := ga.gk.NewString(frame, res)
		return v, nil
	}
}

// FuncIssOsS creates a callable function wrapping a string-transforming function and returns results as an array object.
func (ga *GateAdapter) FuncIssOsS(fn func(string, string) []string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ga.gk.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		obj := ga.gk.NewArray(frame, nil)
		arr, ok := obj.(*Array)
		if !ok {
			return nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, res := range fn(s1, s2) {
			v := ga.gk.NewString(frame, res)
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIssiOsS wraps a function of type func(string, string, int) []string into a FuncCallable compatible function.
func (ga *GateAdapter) FuncIssiOsS(fn func(string, string, int) []string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 3 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ga.gk.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := ga.gk.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		obj := ga.gk.NewArray(frame, nil)
		arr, ok := obj.(*Array)
		if !ok {
			return nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, res := range fn(s1, s2, int(i3)) {
			v := ga.gk.NewString(frame, res)
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIssOi wraps a function accepting two string arguments and returning an int, into a FuncCallable type function.
func (ga *GateAdapter) FuncIssOi(fn func(string, string) int) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ga.gk.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		return ga.gk.NewInt(frame, int64(fn(s1, s2))), nil
	}
}

// FuncIssOs wraps a function that takes two strings and returns a string, converting it into a FuncCallable handler.
func (ga *GateAdapter) FuncIssOs(fn func(string, string) string) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ga.gk.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v := ga.gk.NewString(frame, fn(s1, s2))
		return v, nil
	}
}

// FuncIssOb creates a FuncCallable that evaluates a function on two string arguments extracted from IObject inputs.
// It returns an IObject representing true or false based on the function's result or an error if arguments are invalid.
func (ga *GateAdapter) FuncIssOb(fn func(string, string) bool) FuncCallable {
	return func(frame int, args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := ga.gk.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(s1, s2) {
			return ga.gk.TrueValue(), nil
		}
		return ga.gk.FalseValue(), nil
	}
}

// FuncIbSOs wraps a given function to process a byte slice argument and returns a callable function for the system.
func (ga *GateAdapter) FuncIbSOs(fn func([]byte) string) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		bs1, err := ga.gk.ToBytesArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v := ga.gk.NewString(frame, fn(bs1))
		return v, nil
	}
}

// FuncIsOie wraps a function matching func(string) (int, error) into a FuncCallable, enabling its usage with IObject arguments.
// It validates the argument count, converts the first argument to a string, applies the given function, and returns the result.
// If an error occurs during argument conversion or function execution, it returns appropriate error objects or messages.
func (ga *GateAdapter) FuncIsOie(fn func(string) (int, error)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return ga.gk.NewError(frame, err.Error()), nil
		}
		return ga.gk.NewInt(frame, int64(res)), nil
	}
}

// FuncIsObSe wraps a function accepting a string and returning []byte and error to conform to the FuncCallable type.
// It validates the number of arguments, ensures the first argument is a string, and transforms outputs to IObject types.
func (ga *GateAdapter) FuncIsObSe(fn func(string) ([]byte, error)) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		s1, err := ga.gk.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return ga.gk.NewError(frame, err.Error()), nil
		}
		return ga.gk.NewBytes(frame, res), nil
	}
}

// FuncIiOs wraps a function from int to string into a FuncCallable, ensuring compatibility with the IObject interface.
func (ga *GateAdapter) FuncIiOs(fn func(int) string) FuncCallable {
	return func(frame int, args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgumentsNumber
		}
		i1, err := ga.gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		s := fn(int(i1))
		v := ga.gk.NewString(frame, s)
		return v, nil
	}
}

// BinaryOpInt64 performs a binary operation on two integer values and returns the result.
func (ga *GateAdapter) BinaryOpInt64(op Operator, lhs int64, rhs int64) (int64, error) {
	switch op {
	case OperatorAdd:
		return lhs + rhs, nil
	case OperatorSub:
		return lhs - rhs, nil
	case OperatorMul:
		return lhs * rhs, nil
	case OperatorQuo:
		if rhs == 0 {
			return 0, ErrDivisionByZero
		}
		return lhs / rhs, nil
	case OperatorRem:
		if rhs == 0 {
			return 0, ErrDivisionByZero
		}
		return lhs % rhs, nil
	case OperatorAnd:
		return lhs & rhs, nil
	case OperatorOr:
		return lhs | rhs, nil
	case OperatorXor:
		return lhs ^ rhs, nil
	case OperatorAndNot:
		return lhs &^ rhs, nil
	case OperatorShl:
		return lhs << uint64(rhs), nil
	case OperatorShr:
		return lhs >> uint64(rhs), nil
	case OperatorLess:
		if lhs < rhs {
			return 1, nil
		}
		return 0, nil
	case OperatorGreater:
		if lhs > rhs {
			return 1, nil
		}
		return 0, nil
	case OperatorLessEq:
		if lhs <= rhs {
			return 1, nil
		}
		return 0, nil
	case OperatorGreaterEq:
		if lhs >= rhs {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrInvalidOperator
	}
}

// BoundsCheck validates and adjusts slice bounds using provided low and high indices, ensuring they are within valid range.
func (ga *GateAdapter) BoundsCheck(lowStack IObject, highStack IObject, numElements int64) (int64, int64, error) {
	var lowIdx int64
	if lowStack != ga.gk.UndefinedValue() {
		if low, ok := lowStack.(*Int); ok {
			lowIdx = low.Value()
		} else {
			return 0, 0, fmt.Errorf("invalid slice index type: %s", low.TypeName())
		}
	}
	var highIdx int64
	if highStack == ga.gk.UndefinedValue() {
		highIdx = numElements
	} else if high, ok := highStack.(*Int); ok {
		highIdx = high.Value()
	} else {
		return 0, 0, fmt.Errorf("invalid slice index type: %s", high.TypeName())
	}
	if lowIdx > highIdx {
		return 0, 0, fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
	}
	if lowIdx < 0 {
		lowIdx = 0
	} else if lowIdx > numElements {
		lowIdx = numElements
	}
	if highIdx < 0 {
		highIdx = 0
	} else if highIdx > numElements {
		highIdx = numElements
	}
	return lowIdx, highIdx, nil
}

// IndexAssign assigns a value to a nested structure, using selectors to determine the target location.
// It navigates through the provided selectors and performs an assignment on the target object at the final index.
// Returns an error if any selector is invalid, the object is not indexable, or the assignment fails.
func (ga *GateAdapter) IndexAssign(frame int, dst IObject, src IObject, selectors []IObject) error {
	if len(selectors) == 0 {
		return ErrSelectorNotProvided
	}
	numSel := len(selectors)
	for sIdx := numSel - 1; sIdx > 0; sIdx-- {
		next, err := dst.IndexGet(frame, selectors[sIdx])
		if err != nil {
			return ComputeIndexGetError(err, dst.TypeName(), selectors[sIdx].TypeName())
		}
		dst = next
	}
	if err := dst.IndexSet(selectors[0], src); err != nil {
		return ComputeIndexSetError(err, dst.TypeName(), src.TypeName())
	}
	return nil
}
