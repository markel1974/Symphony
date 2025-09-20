package sdk

import (
	"math"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewMath)
}

// Math serves as a container for mathematical operations and modules, mapping module names to their respective objects.
type Math struct {
	*bytecode.Package
	nanObj objects.IObject
}

// NewMath initializes and returns a new instance of Math with predefined mathematical constants and function modules.
func NewMath(factory objects.IGateKeeper) bytecode.IPackage {
	const (
		defE         = "E"
		defPi        = "Pi"
		defPhi       = "Phi"
		defSqrt2     = "Sqrt2"
		defSqrtE     = "SqrtE"
		defSqrtPi    = "SqrtPi"
		defSqrtPhi   = "SqrtPhi"
		defLn2       = "Ln2"
		defLog2E     = "Log2E"
		defLn10      = "Ln10"
		defLog10E    = "Log10E"
		defNaN       = "NaN"
		defAbs       = "Abs"
		defAcos      = "Acos"
		defAcosh     = "Acosh"
		defAsin      = "Asin"
		defAsinh     = "Asinh"
		defAtan      = "Atan"
		defAtan2     = "Atan2"
		defAtanh     = "Atanh"
		defCbrt      = "Cbrt"
		defCeil      = "Ceil"
		defCopysign  = "Copysign"
		defCos       = "Cos"
		defCosh      = "Cosh"
		defDim       = "Dim"
		defErf       = "Erf"
		defErfc      = "Erfc"
		defExp       = "Exp"
		defExp2      = "Exp2"
		defExpm1     = "Expm1"
		defFloor     = "Floor"
		defGamma     = "Gamma"
		defHypot     = "Hypot"
		defIlogb     = "Ilogb"
		defInf       = "Inf"
		defIsInf     = "IsInf"
		defIsNaN     = "IsNaN"
		defJ0        = "J0"
		defJ1        = "J1"
		defJn        = "Jn"
		defLdexp     = "Ldexp"
		defLog       = "Log"
		defLog10     = "Log10"
		defLog1p     = "Log1p"
		defLog2      = "Log2"
		defLogb      = "Logb"
		defMax       = "Max"
		defMin       = "Min"
		defMod       = "Mod"
		defNextafter = "Nextafter"
		defPow       = "Pow"
		defPow10     = "Pow10"
		defRemainder = "Remainder"
		defSignbit   = "Signbit"
		defSin       = "Sin"
		defSinh      = "Sinh"
		defSqrt      = "Sqrt"
		defTan       = "Tan"
		defTanh      = "Tanh"
		defTrunc     = "Trunc"
		defY0        = "Y0"
		defY1        = "Y1"
		defYn        = "Yn"
	)

	m := &Math{Package: bytecode.NewPackage("math")}
	m.nanObj = factory.NewFloat(objects.FrameStatic, math.NaN())

	m.Add(defE, factory.NewFloat(objects.FrameStatic, math.E))
	m.Add(defPi, factory.NewFloat(objects.FrameStatic, math.Pi))
	m.Add(defPhi, factory.NewFloat(objects.FrameStatic, math.Phi))
	m.Add(defSqrt2, factory.NewFloat(objects.FrameStatic, math.Sqrt2))
	m.Add(defSqrtE, factory.NewFloat(objects.FrameStatic, math.SqrtE))
	m.Add(defSqrtPi, factory.NewFloat(objects.FrameStatic, math.SqrtPi))
	m.Add(defSqrtPhi, factory.NewFloat(objects.FrameStatic, math.SqrtPhi))
	m.Add(defLn2, factory.NewFloat(objects.FrameStatic, math.Ln2))
	m.Add(defLog2E, factory.NewFloat(objects.FrameStatic, math.Log2E))
	m.Add(defLn10, factory.NewFloat(objects.FrameStatic, math.Ln10))
	m.Add(defLog10E, factory.NewFloat(objects.FrameStatic, math.Log10E))
	m.Add(defNaN, factory.NewFuncImport(objects.FrameStatic, defNaN, 0, m.nan))
	m.Add(defAbs, factory.NewFuncImport(objects.FrameStatic, defAbs, 1, m.funcF64ToF64(math.Abs)))
	m.Add(defAcos, factory.NewFuncImport(objects.FrameStatic, defAcos, 1, m.funcF64ToF64(math.Acos)))
	m.Add(defAcosh, factory.NewFuncImport(objects.FrameStatic, defAcosh, 1, m.funcF64ToF64(math.Acosh)))
	m.Add(defAsin, factory.NewFuncImport(objects.FrameStatic, defAsin, 1, m.funcF64ToF64(math.Asin)))
	m.Add(defAsinh, factory.NewFuncImport(objects.FrameStatic, defAsinh, 1, m.funcF64ToF64(math.Asinh)))
	m.Add(defAtan, factory.NewFuncImport(objects.FrameStatic, defAtan, 1, m.funcF64ToF64(math.Atan)))
	m.Add(defAtan2, factory.NewFuncImport(objects.FrameStatic, defAtan2, 2, m.funcF64F64ToF64(math.Atan2)))
	m.Add(defAtanh, factory.NewFuncImport(objects.FrameStatic, defAtanh, 1, m.funcF64ToF64(math.Atanh)))
	m.Add(defCbrt, factory.NewFuncImport(objects.FrameStatic, defCbrt, 1, m.funcF64ToF64(math.Cbrt)))
	m.Add(defCeil, factory.NewFuncImport(objects.FrameStatic, defCeil, 1, m.funcF64ToF64(math.Ceil)))
	m.Add(defCopysign, factory.NewFuncImport(objects.FrameStatic, defCopysign, 2, m.funcF64F64ToF64(math.Copysign)))
	m.Add(defCos, factory.NewFuncImport(objects.FrameStatic, defCos, 1, m.funcF64ToF64(math.Cos)))
	m.Add(defCosh, factory.NewFuncImport(objects.FrameStatic, defCosh, 1, m.funcF64ToF64(math.Cosh)))
	m.Add(defDim, factory.NewFuncImport(objects.FrameStatic, defDim, 2, m.funcF64F64ToF64(math.Dim)))
	m.Add(defErf, factory.NewFuncImport(objects.FrameStatic, defErf, 1, m.funcF64ToF64(math.Erf)))
	m.Add(defErfc, factory.NewFuncImport(objects.FrameStatic, defErfc, 1, m.funcF64ToF64(math.Erfc)))
	m.Add(defExp, factory.NewFuncImport(objects.FrameStatic, defExp, 1, m.funcF64ToF64(math.Exp)))
	m.Add(defExp2, factory.NewFuncImport(objects.FrameStatic, defExp2, 1, m.funcF64ToF64(math.Exp2)))
	m.Add(defExpm1, factory.NewFuncImport(objects.FrameStatic, defExpm1, 1, m.funcF64ToF64(math.Expm1)))
	m.Add(defFloor, factory.NewFuncImport(objects.FrameStatic, defFloor, 1, m.funcF64ToF64(math.Floor)))
	m.Add(defGamma, factory.NewFuncImport(objects.FrameStatic, defGamma, 1, m.funcF64ToF64(math.Gamma)))
	m.Add(defHypot, factory.NewFuncImport(objects.FrameStatic, defHypot, 2, m.funcF64F64ToF64(math.Hypot)))
	m.Add(defIlogb, factory.NewFuncImport(objects.FrameStatic, defIlogb, 1, m.funcFloat64ToInt(math.Ilogb)))
	m.Add(defInf, factory.NewFuncImport(objects.FrameStatic, defInf, 1, m.funcIntToF64(math.Inf)))
	m.Add(defIsInf, factory.NewFuncImport(objects.FrameStatic, defIsInf, 1, m.funcF64IntToBool(math.IsInf)))
	m.Add(defIsNaN, factory.NewFuncImport(objects.FrameStatic, defIsNaN, 1, m.funcF64ToBool(math.IsNaN)))
	m.Add(defJ0, factory.NewFuncImport(objects.FrameStatic, defJ0, 1, m.funcF64ToF64(math.J0)))
	m.Add(defJ1, factory.NewFuncImport(objects.FrameStatic, defJ1, 1, m.funcF64ToF64(math.J1)))
	m.Add(defJn, factory.NewFuncImport(objects.FrameStatic, defJn, 1, m.funcIntF64ToF64(math.Jn)))
	m.Add(defLdexp, factory.NewFuncImport(objects.FrameStatic, defLdexp, 1, m.funcF64IntToF64(math.Ldexp)))
	m.Add(defLog, factory.NewFuncImport(objects.FrameStatic, defLog, 1, m.funcF64ToF64(math.Log)))
	m.Add(defLog10, factory.NewFuncImport(objects.FrameStatic, defLog10, 1, m.funcF64ToF64(math.Log10)))
	m.Add(defLog1p, factory.NewFuncImport(objects.FrameStatic, defLog1p, 1, m.funcF64ToF64(math.Log1p)))
	m.Add(defLog2, factory.NewFuncImport(objects.FrameStatic, defLog2, 1, m.funcF64ToF64(math.Log2)))
	m.Add(defLogb, factory.NewFuncImport(objects.FrameStatic, defLogb, 1, m.funcF64ToF64(math.Logb)))
	m.Add(defMax, factory.NewFuncImport(objects.FrameStatic, defMax, 2, m.funcF64F64ToF64(math.Max)))
	m.Add(defMin, factory.NewFuncImport(objects.FrameStatic, defMin, 2, m.funcF64F64ToF64(math.Min)))
	m.Add(defMod, factory.NewFuncImport(objects.FrameStatic, defMod, 2, m.funcF64F64ToF64(math.Mod)))
	m.Add(defNextafter, factory.NewFuncImport(objects.FrameStatic, defNextafter, 2, m.funcF64F64ToF64(math.Nextafter)))
	m.Add(defPow, factory.NewFuncImport(objects.FrameStatic, defPow, 2, m.funcF64F64ToF64(math.Pow)))
	m.Add(defPow10, factory.NewFuncImport(objects.FrameStatic, defPow10, 1, m.funcIntToF64(math.Pow10)))
	m.Add(defRemainder, factory.NewFuncImport(objects.FrameStatic, defRemainder, 2, m.funcF64F64ToF64(math.Remainder)))
	m.Add(defSignbit, factory.NewFuncImport(objects.FrameStatic, defSignbit, 1, m.funcF64ToBool(math.Signbit)))
	m.Add(defSin, factory.NewFuncImport(objects.FrameStatic, defSin, 1, m.funcF64ToF64(math.Sin)))
	m.Add(defSinh, factory.NewFuncImport(objects.FrameStatic, defSinh, 1, m.funcF64ToF64(math.Sinh)))
	m.Add(defSqrt, factory.NewFuncImport(objects.FrameStatic, defSqrt, 1, m.funcF64ToF64(math.Sqrt)))
	m.Add(defTan, factory.NewFuncImport(objects.FrameStatic, defTan, 1, m.funcF64ToF64(math.Tan)))
	m.Add(defTanh, factory.NewFuncImport(objects.FrameStatic, defTanh, 1, m.funcF64ToF64(math.Tanh)))
	m.Add(defTrunc, factory.NewFuncImport(objects.FrameStatic, defTrunc, 1, m.funcF64ToF64(math.Trunc)))
	m.Add(defY0, factory.NewFuncImport(objects.FrameStatic, defY0, 1, m.funcF64ToF64(math.Y0)))
	m.Add(defY1, factory.NewFuncImport(objects.FrameStatic, defY1, 1, m.funcF64ToF64(math.Y1)))
	m.Add(defYn, factory.NewFuncImport(objects.FrameStatic, defYn, 2, m.funcIntF64ToF64(math.Yn)))

	return m
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

// funcF64ToF64 creates a Invocable that applies a function transforming a float64 input to a float64 result.
// It validates a single argument, converts it to float64, applies the function, and returns the result as an IObject.
// Returns an error for invalid arguments or if the number of arguments is not exactly one.
func (m *Math) funcF64ToF64(fn func(float64) float64) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(f1)), nil
	}
}

// funcIntToF64 converts a function from int to float64 into a Invocable that operates on IObject arguments.
func (m *Math) funcIntToF64(fn func(int) float64) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewFloat(frame, fn(int(i1))), nil
	}
}

// funcF64F64ToF64 adapts a function taking two float64 arguments and returning a float64 to a Invocable.
// It converts IObject arguments to float64, applies the input function, and returns the result as an IObject.
// Returns an error if the argument count is not 2 or if type conversion fails.
func (m *Math) funcF64F64ToF64(fn func(float64, float64) float64) objects.Invocable {
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

// funcIntF64ToF64 adapts a function accepting an int and float64 as inputs, returning a float64, into a Invocable.
func (m *Math) funcIntF64ToF64(fn func(int, float64) float64) objects.Invocable {
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

// funcF64IntToF64 converts a function accepting (float64, int) and returning float64 into a Invocable type.
// It requires exactly two arguments: the first must be convertible to float64 and the second to int.
// Returns an IObject representing the result of applying the function or an error if argument conversion fails.
func (m *Math) funcF64IntToF64(fn func(float64, int) float64) objects.Invocable {
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
func (m *Math) funcFloat64ToInt(fn func(float64) int) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		f1, err := gk.ToFloat64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewInt(frame, int64(fn(f1))), nil
	}
}

// funcF64IntToBool wraps a function of type func(float64, int) bool into a Invocable compatible with IObject arguments.
// It converts the first argument to float64 and the second argument to int, validates their types, and applies the function.
// Returns the system-defined TrueValue or FalseValue based on the function's result, or an error on failure.
// Returns ErrInvalidArgumentsNumber if the number of arguments passed is not exactly two.
func (m *Math) funcF64IntToBool(fn func(float64, int) bool) objects.Invocable {
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

// funcF64ToBool wraps a given float64 predicate function and returns a Invocable to evaluate the predicate using IObject.
// If the argument count is incorrect or conversion to float64 fails, it returns an error.
// The result of the predicate determines the boolean IObject returned: TrueValue or FalseValue.
func (m *Math) funcF64ToBool(fn func(float64) bool) objects.Invocable {
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
