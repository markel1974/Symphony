package sdk

import (
	"math"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewMath)
}

// Math serves as a container for mathematical operations and modules, mapping module names to their respective objects.
type Math struct {
	container map[string]objects.IObject
}

// NewMath initializes and returns a new instance of Math with predefined mathematical constants and function modules.
func NewMath(factory objects.IGateKeeper) IPackage {
	m := &Math{}
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
		factory.NewFuncPackage("Abs", factory.FuncIf64Of64(math.Abs)),
		factory.NewFuncPackage("Acos", factory.FuncIf64Of64(math.Acos)),
		factory.NewFuncPackage("Acosh", factory.FuncIf64Of64(math.Acosh)),
		factory.NewFuncPackage("Asin", factory.FuncIf64Of64(math.Asin)),
		factory.NewFuncPackage("Asinh", factory.FuncIf64Of64(math.Asinh)),
		factory.NewFuncPackage("Atan", factory.FuncIf64Of64(math.Atan)),
		factory.NewFuncPackage("Atan2", factory.FuncIf64f64Of64(math.Atan2)),
		factory.NewFuncPackage("Atanh", factory.FuncIf64Of64(math.Atanh)),
		factory.NewFuncPackage("Cbrt", factory.FuncIf64Of64(math.Cbrt)),
		factory.NewFuncPackage("Ceil", factory.FuncIf64Of64(math.Ceil)),
		factory.NewFuncPackage("Copysign", factory.FuncIf64f64Of64(math.Copysign)),
		factory.NewFuncPackage("Cos", factory.FuncIf64Of64(math.Cos)),
		factory.NewFuncPackage("Cosh", factory.FuncIf64Of64(math.Cosh)),
		factory.NewFuncPackage("Dim", factory.FuncIf64f64Of64(math.Dim)),
		factory.NewFuncPackage("Erf", factory.FuncIf64Of64(math.Erf)),
		factory.NewFuncPackage("Erfc", factory.FuncIf64Of64(math.Erfc)),
		factory.NewFuncPackage("Exp", factory.FuncIf64Of64(math.Exp)),
		factory.NewFuncPackage("Exp2", factory.FuncIf64Of64(math.Exp2)),
		factory.NewFuncPackage("Expm1", factory.FuncIf64Of64(math.Expm1)),
		factory.NewFuncPackage("Floor", factory.FuncIf64Of64(math.Floor)),
		factory.NewFuncPackage("Gamma", factory.FuncIf64Of64(math.Gamma)),
		factory.NewFuncPackage("Hypot", factory.FuncIf64f64Of64(math.Hypot)),
		factory.NewFuncPackage("Ilogb", factory.FuncIf64Oi(math.Ilogb)),
		factory.NewFuncPackage("Inf", factory.FuncIiOf64(math.Inf)),
		factory.NewFuncPackage("IsInf", factory.FuncIf64iOb(math.IsInf)),
		factory.NewFuncPackage("IsNaN", factory.FuncIf64Ob(math.IsNaN)),
		factory.NewFuncPackage("J0", factory.FuncIf64Of64(math.J0)),
		factory.NewFuncPackage("J1", factory.FuncIf64Of64(math.J1)),
		factory.NewFuncPackage("Jn", factory.FuncIif64Of64(math.Jn)),
		factory.NewFuncPackage("Ldexp", factory.FuncIf64iOf64(math.Ldexp)),
		factory.NewFuncPackage("Log", factory.FuncIf64Of64(math.Log)),
		factory.NewFuncPackage("Log10", factory.FuncIf64Of64(math.Log10)),
		factory.NewFuncPackage("Log1p", factory.FuncIf64Of64(math.Log1p)),
		factory.NewFuncPackage("Log2", factory.FuncIf64Of64(math.Log2)),
		factory.NewFuncPackage("Logb", factory.FuncIf64Of64(math.Logb)),
		factory.NewFuncPackage("Max", factory.FuncIf64f64Of64(math.Max)),
		factory.NewFuncPackage("Min", factory.FuncIf64f64Of64(math.Min)),
		factory.NewFuncPackage("Mod", factory.FuncIf64f64Of64(math.Mod)),
		factory.NewFuncPackage("NaN", m.funcInOf64(math.NaN)),
		factory.NewFuncPackage("Nextafter", factory.FuncIf64f64Of64(math.Nextafter)),
		factory.NewFuncPackage("Pow", factory.FuncIf64f64Of64(math.Pow)),
		factory.NewFuncPackage("Pow10", factory.FuncIiOf64(math.Pow10)),
		factory.NewFuncPackage("Remainder", factory.FuncIf64f64Of64(math.Remainder)),
		factory.NewFuncPackage("Signbit", factory.FuncIf64Ob(math.Signbit)),
		factory.NewFuncPackage("Sin", factory.FuncIf64Of64(math.Sin)),
		factory.NewFuncPackage("Sinh", factory.FuncIf64Of64(math.Sinh)),
		factory.NewFuncPackage("Sqrt", factory.FuncIf64Of64(math.Sqrt)),
		factory.NewFuncPackage("Tan", factory.FuncIf64Of64(math.Tan)),
		factory.NewFuncPackage("Tanh", factory.FuncIf64Of64(math.Tanh)),
		factory.NewFuncPackage("Trunc", factory.FuncIf64Of64(math.Trunc)),
		factory.NewFuncPackage("Y0", factory.FuncIf64Of64(math.Y0)),
		factory.NewFuncPackage("Y1", factory.FuncIf64Of64(math.Y1)),
		factory.NewFuncPackage("Yn", factory.FuncIif64Of64(math.Yn)),
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

// funcInOf64 wraps a no-argument function returning float64 into a FuncCallable to integrate with the GateAdapter system.
// It enforces no arguments and converts the float64 result to an IObject, returning ErrInvalidArgumentsNumber for invalid input.
func (m *Math) funcInOf64(fn func() float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, objects.ErrInvalidArgumentsNumber
		}
		return gk.NewFloat(frame, fn()), nil
	}
}
