package sdk

import (
	"math"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewMath)
}

// Math serves as a container for mathematical operations and modules, mapping module names to their respective objects.
type Math struct {
	container map[string]objects.IObject
	nanObj    objects.IObject
}

// NewMath initializes and returns a new instance of Math with predefined mathematical constants and function modules.
func NewMath(factory objects.IGateKeeper) IPackage {
	m := &Math{}
	m.nanObj = factory.NewFloat(objects.FrameStatic, math.NaN())
	constants := map[string]objects.IObject{
		"E":       factory.NewFloat(objects.FrameStatic, math.E),
		"Pi":      factory.NewFloat(objects.FrameStatic, math.Pi),
		"Phi":     factory.NewFloat(objects.FrameStatic, math.Phi),
		"Sqrt2":   factory.NewFloat(objects.FrameStatic, math.Sqrt2),
		"SqrtE":   factory.NewFloat(objects.FrameStatic, math.SqrtE),
		"SqrtPi":  factory.NewFloat(objects.FrameStatic, math.SqrtPi),
		"SqrtPhi": factory.NewFloat(objects.FrameStatic, math.SqrtPhi),
		"Ln2":     factory.NewFloat(objects.FrameStatic, math.Ln2),
		"Log2E":   factory.NewFloat(objects.FrameStatic, math.Log2E),
		"Ln10":    factory.NewFloat(objects.FrameStatic, math.Ln10),
		"Log10E":  factory.NewFloat(objects.FrameStatic, math.Log10E),
	}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "NaN", 0, m.nan),
		factory.NewFuncImport(objects.FrameStatic, "Abs", 1, m.funcF64ToF64(math.Abs)),
		factory.NewFuncImport(objects.FrameStatic, "Acos", 1, m.funcF64ToF64(math.Acos)),
		factory.NewFuncImport(objects.FrameStatic, "Acosh", 1, m.funcF64ToF64(math.Acosh)),
		factory.NewFuncImport(objects.FrameStatic, "Asin", 1, m.funcF64ToF64(math.Asin)),
		factory.NewFuncImport(objects.FrameStatic, "Asinh", 1, m.funcF64ToF64(math.Asinh)),
		factory.NewFuncImport(objects.FrameStatic, "Atan", 1, m.funcF64ToF64(math.Atan)),
		factory.NewFuncImport(objects.FrameStatic, "Atan2", 2, m.funcF64F64ToF64(math.Atan2)),
		factory.NewFuncImport(objects.FrameStatic, "Atanh", 1, m.funcF64ToF64(math.Atanh)),
		factory.NewFuncImport(objects.FrameStatic, "Cbrt", 1, m.funcF64ToF64(math.Cbrt)),
		factory.NewFuncImport(objects.FrameStatic, "Ceil", 1, m.funcF64ToF64(math.Ceil)),
		factory.NewFuncImport(objects.FrameStatic, "Copysign", 2, m.funcF64F64ToF64(math.Copysign)),
		factory.NewFuncImport(objects.FrameStatic, "Cos", 1, m.funcF64ToF64(math.Cos)),
		factory.NewFuncImport(objects.FrameStatic, "Cosh", 1, m.funcF64ToF64(math.Cosh)),
		factory.NewFuncImport(objects.FrameStatic, "Dim", 2, m.funcF64F64ToF64(math.Dim)),
		factory.NewFuncImport(objects.FrameStatic, "Erf", 1, m.funcF64ToF64(math.Erf)),
		factory.NewFuncImport(objects.FrameStatic, "Erfc", 1, m.funcF64ToF64(math.Erfc)),
		factory.NewFuncImport(objects.FrameStatic, "Exp", 1, m.funcF64ToF64(math.Exp)),
		factory.NewFuncImport(objects.FrameStatic, "Exp2", 1, m.funcF64ToF64(math.Exp2)),
		factory.NewFuncImport(objects.FrameStatic, "Expm1", 1, m.funcF64ToF64(math.Expm1)),
		factory.NewFuncImport(objects.FrameStatic, "Floor", 1, m.funcF64ToF64(math.Floor)),
		factory.NewFuncImport(objects.FrameStatic, "Gamma", 1, m.funcF64ToF64(math.Gamma)),
		factory.NewFuncImport(objects.FrameStatic, "Hypot", 2, m.funcF64F64ToF64(math.Hypot)),
		factory.NewFuncImport(objects.FrameStatic, "Ilogb", 1, m.funcFloat64ToInt(math.Ilogb)),
		factory.NewFuncImport(objects.FrameStatic, "Inf", 1, m.funcIntToF64(math.Inf)),
		factory.NewFuncImport(objects.FrameStatic, "IsInf", 1, m.funcF64IntToBool(math.IsInf)),
		factory.NewFuncImport(objects.FrameStatic, "IsNaN", 1, m.funcF64ToBool(math.IsNaN)),
		factory.NewFuncImport(objects.FrameStatic, "J0", 1, m.funcF64ToF64(math.J0)),
		factory.NewFuncImport(objects.FrameStatic, "J1", 1, m.funcF64ToF64(math.J1)),
		factory.NewFuncImport(objects.FrameStatic, "Jn", 1, m.funcIntF64ToF64(math.Jn)),
		factory.NewFuncImport(objects.FrameStatic, "Ldexp", 1, m.funcF64IntToF64(math.Ldexp)),
		factory.NewFuncImport(objects.FrameStatic, "Log", 1, m.funcF64ToF64(math.Log)),
		factory.NewFuncImport(objects.FrameStatic, "Log10", 1, m.funcF64ToF64(math.Log10)),
		factory.NewFuncImport(objects.FrameStatic, "Log1p", 1, m.funcF64ToF64(math.Log1p)),
		factory.NewFuncImport(objects.FrameStatic, "Log2", 1, m.funcF64ToF64(math.Log2)),
		factory.NewFuncImport(objects.FrameStatic, "Logb", 1, m.funcF64ToF64(math.Logb)),
		factory.NewFuncImport(objects.FrameStatic, "Max", 2, m.funcF64F64ToF64(math.Max)),
		factory.NewFuncImport(objects.FrameStatic, "Min", 2, m.funcF64F64ToF64(math.Min)),
		factory.NewFuncImport(objects.FrameStatic, "Mod", 2, m.funcF64F64ToF64(math.Mod)),
		factory.NewFuncImport(objects.FrameStatic, "Nextafter", 2, m.funcF64F64ToF64(math.Nextafter)),
		factory.NewFuncImport(objects.FrameStatic, "Pow", 2, m.funcF64F64ToF64(math.Pow)),
		factory.NewFuncImport(objects.FrameStatic, "Pow10", 1, m.funcIntToF64(math.Pow10)),
		factory.NewFuncImport(objects.FrameStatic, "Remainder", 2, m.funcF64F64ToF64(math.Remainder)),
		factory.NewFuncImport(objects.FrameStatic, "Signbit", 1, m.funcF64ToBool(math.Signbit)),
		factory.NewFuncImport(objects.FrameStatic, "Sin", 1, m.funcF64ToF64(math.Sin)),
		factory.NewFuncImport(objects.FrameStatic, "Sinh", 1, m.funcF64ToF64(math.Sinh)),
		factory.NewFuncImport(objects.FrameStatic, "Sqrt", 1, m.funcF64ToF64(math.Sqrt)),
		factory.NewFuncImport(objects.FrameStatic, "Tan", 1, m.funcF64ToF64(math.Tan)),
		factory.NewFuncImport(objects.FrameStatic, "Tanh", 1, m.funcF64ToF64(math.Tanh)),
		factory.NewFuncImport(objects.FrameStatic, "Trunc", 1, m.funcF64ToF64(math.Trunc)),
		factory.NewFuncImport(objects.FrameStatic, "Y0", 1, m.funcF64ToF64(math.Y0)),
		factory.NewFuncImport(objects.FrameStatic, "Y1", 1, m.funcF64ToF64(math.Y1)),
		factory.NewFuncImport(objects.FrameStatic, "Yn", 2, m.funcIntF64ToF64(math.Yn)),
	}
	m.container = BuildContainer(container, constants)
	return m
}

// Name returns the name of the Math module as a string.
func (m *Math) Name() string {
	return "math"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (m *Math) Get(name string) (objects.IObject, bool) {
	v, ok := m.container[name]
	return v, ok
}

// nan returns the Not-A-Number value as a float64.
func (m *Math) nan(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	ret := gk.UndefinedValue()
	err := objects.ErrInvalidArgumentsNumber
	if len(args) == 0 {
		ret = m.nanObj
		err = nil
	}
	return 1, ret, err
}

// funcF64ToF64 creates a FuncCallable that applies a function transforming a float64 input to a float64 result.
// It validates a single argument, converts it to float64, applies the function, and returns the result as an IObject.
// Returns an error for invalid arguments or if the number of arguments is not exactly one.
func (m *Math) funcF64ToF64(fn func(float64) float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(f1)), nil
	}
}

// funcIntToF64 converts a function from int to float64 into a FuncCallable that operates on IObject arguments.
func (m *Math) funcIntToF64(fn func(int) float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(int(i1))), nil
	}
}

// funcF64F64ToF64 adapts a function taking two float64 arguments and returning a float64 to a FuncCallable.
// It converts IObject arguments to float64, applies the input function, and returns the result as an IObject.
// Returns an error if the argument count is not 2 or if type conversion fails.
func (m *Math) funcF64F64ToF64(fn func(float64, float64) float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		f2, err := gk.ToFloat64Arg(1, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(f1, f2)), nil
	}
}

// funcIntF64ToF64 adapts a function accepting an int and float64 as inputs, returning a float64, into a FuncCallable.
func (m *Math) funcIntF64ToF64(fn func(int, float64) float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		f2, err := gk.ToFloat64Arg(1, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(int(i1), f2)), nil
	}
}

// funcF64IntToF64 converts a function accepting (float64, int) and returning float64 into a FuncCallable type.
// It requires exactly two arguments: the first must be convertible to float64 and the second to int.
// Returns an IObject representing the result of applying the function or an error if argument conversion fails.
func (m *Math) funcF64IntToF64(fn func(float64, int) float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		i2, err := gk.ToInt64Arg(1, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(f1, int(i2))), nil
	}
}

// FuncIf64Oi converts a function accepting a float64 and returning an int into a callable function for GateAdapter.
func (m *Math) funcFloat64ToInt(fn func(float64) int) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewInt(frame, int64(fn(f1))), nil
	}
}

// funcF64IntToBool wraps a function of type func(float64, int) bool into a FuncCallable compatible with IObject arguments.
// It converts the first argument to float64 and the second argument to int, validates their types, and applies the function.
// Returns the system-defined TrueValue or FalseValue based on the function's result, or an error on failure.
// Returns ErrInvalidArgumentsNumber if the number of arguments passed is not exactly two.
func (m *Math) funcF64IntToBool(fn func(float64, int) bool) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		i2, err := gk.ToInt64Arg(1, args)
		if err != nil {
			return 0, nil, err
		}
		if fn(f1, int(i2)) {
			return 1, gk.TrueValue(), nil
		}
		return 1, gk.FalseValue(), nil
	}
}

// funcF64ToBool wraps a given float64 predicate function and returns a FuncCallable to evaluate the predicate using IObject.
// If the argument count is incorrect or conversion to float64 fails, it returns an error.
// The result of the predicate determines the boolean IObject returned: TrueValue or FalseValue.
func (m *Math) funcF64ToBool(fn func(float64) bool) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		if fn(f1) {
			return 1, gk.TrueValue(), nil
		}
		return 1, gk.FalseValue(), nil
	}
}
