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
		factory.NewFuncExternal(objects.FrameStatic, "Abs", factory.FuncIf64Of64(math.Abs)),
		factory.NewFuncExternal(objects.FrameStatic, "Acos", factory.FuncIf64Of64(math.Acos)),
		factory.NewFuncExternal(objects.FrameStatic, "Acosh", factory.FuncIf64Of64(math.Acosh)),
		factory.NewFuncExternal(objects.FrameStatic, "Asin", factory.FuncIf64Of64(math.Asin)),
		factory.NewFuncExternal(objects.FrameStatic, "Asinh", factory.FuncIf64Of64(math.Asinh)),
		factory.NewFuncExternal(objects.FrameStatic, "Atan", factory.FuncIf64Of64(math.Atan)),
		factory.NewFuncExternal(objects.FrameStatic, "Atan2", factory.FuncIf64f64Of64(math.Atan2)),
		factory.NewFuncExternal(objects.FrameStatic, "Atanh", factory.FuncIf64Of64(math.Atanh)),
		factory.NewFuncExternal(objects.FrameStatic, "Cbrt", factory.FuncIf64Of64(math.Cbrt)),
		factory.NewFuncExternal(objects.FrameStatic, "Ceil", factory.FuncIf64Of64(math.Ceil)),
		factory.NewFuncExternal(objects.FrameStatic, "Copysign", factory.FuncIf64f64Of64(math.Copysign)),
		factory.NewFuncExternal(objects.FrameStatic, "Cos", factory.FuncIf64Of64(math.Cos)),
		factory.NewFuncExternal(objects.FrameStatic, "Cosh", factory.FuncIf64Of64(math.Cosh)),
		factory.NewFuncExternal(objects.FrameStatic, "Dim", factory.FuncIf64f64Of64(math.Dim)),
		factory.NewFuncExternal(objects.FrameStatic, "Erf", factory.FuncIf64Of64(math.Erf)),
		factory.NewFuncExternal(objects.FrameStatic, "Erfc", factory.FuncIf64Of64(math.Erfc)),
		factory.NewFuncExternal(objects.FrameStatic, "Exp", factory.FuncIf64Of64(math.Exp)),
		factory.NewFuncExternal(objects.FrameStatic, "Exp2", factory.FuncIf64Of64(math.Exp2)),
		factory.NewFuncExternal(objects.FrameStatic, "Expm1", factory.FuncIf64Of64(math.Expm1)),
		factory.NewFuncExternal(objects.FrameStatic, "Floor", factory.FuncIf64Of64(math.Floor)),
		factory.NewFuncExternal(objects.FrameStatic, "Gamma", factory.FuncIf64Of64(math.Gamma)),
		factory.NewFuncExternal(objects.FrameStatic, "Hypot", factory.FuncIf64f64Of64(math.Hypot)),
		factory.NewFuncExternal(objects.FrameStatic, "Ilogb", factory.FuncIf64Oi(math.Ilogb)),
		factory.NewFuncExternal(objects.FrameStatic, "Inf", factory.FuncIiOf64(math.Inf)),
		factory.NewFuncExternal(objects.FrameStatic, "IsInf", factory.FuncIf64iOb(math.IsInf)),
		factory.NewFuncExternal(objects.FrameStatic, "IsNaN", factory.FuncIf64Ob(math.IsNaN)),
		factory.NewFuncExternal(objects.FrameStatic, "J0", factory.FuncIf64Of64(math.J0)),
		factory.NewFuncExternal(objects.FrameStatic, "J1", factory.FuncIf64Of64(math.J1)),
		factory.NewFuncExternal(objects.FrameStatic, "Jn", factory.FuncIif64Of64(math.Jn)),
		factory.NewFuncExternal(objects.FrameStatic, "Ldexp", factory.FuncIf64iOf64(math.Ldexp)),
		factory.NewFuncExternal(objects.FrameStatic, "Log", factory.FuncIf64Of64(math.Log)),
		factory.NewFuncExternal(objects.FrameStatic, "Log10", factory.FuncIf64Of64(math.Log10)),
		factory.NewFuncExternal(objects.FrameStatic, "Log1p", factory.FuncIf64Of64(math.Log1p)),
		factory.NewFuncExternal(objects.FrameStatic, "Log2", factory.FuncIf64Of64(math.Log2)),
		factory.NewFuncExternal(objects.FrameStatic, "Logb", factory.FuncIf64Of64(math.Logb)),
		factory.NewFuncExternal(objects.FrameStatic, "Max", factory.FuncIf64f64Of64(math.Max)),
		factory.NewFuncExternal(objects.FrameStatic, "Min", factory.FuncIf64f64Of64(math.Min)),
		factory.NewFuncExternal(objects.FrameStatic, "Mod", factory.FuncIf64f64Of64(math.Mod)),
		factory.NewFuncExternal(objects.FrameStatic, "NaN", m.funcInOf64(math.NaN)),
		factory.NewFuncExternal(objects.FrameStatic, "Nextafter", factory.FuncIf64f64Of64(math.Nextafter)),
		factory.NewFuncExternal(objects.FrameStatic, "Pow", factory.FuncIf64f64Of64(math.Pow)),
		factory.NewFuncExternal(objects.FrameStatic, "Pow10", factory.FuncIiOf64(math.Pow10)),
		factory.NewFuncExternal(objects.FrameStatic, "Remainder", factory.FuncIf64f64Of64(math.Remainder)),
		factory.NewFuncExternal(objects.FrameStatic, "Signbit", factory.FuncIf64Ob(math.Signbit)),
		factory.NewFuncExternal(objects.FrameStatic, "Sin", factory.FuncIf64Of64(math.Sin)),
		factory.NewFuncExternal(objects.FrameStatic, "Sinh", factory.FuncIf64Of64(math.Sinh)),
		factory.NewFuncExternal(objects.FrameStatic, "Sqrt", factory.FuncIf64Of64(math.Sqrt)),
		factory.NewFuncExternal(objects.FrameStatic, "Tan", factory.FuncIf64Of64(math.Tan)),
		factory.NewFuncExternal(objects.FrameStatic, "Tanh", factory.FuncIf64Of64(math.Tanh)),
		factory.NewFuncExternal(objects.FrameStatic, "Trunc", factory.FuncIf64Of64(math.Trunc)),
		factory.NewFuncExternal(objects.FrameStatic, "Y0", factory.FuncIf64Of64(math.Y0)),
		factory.NewFuncExternal(objects.FrameStatic, "Y1", factory.FuncIf64Of64(math.Y1)),
		factory.NewFuncExternal(objects.FrameStatic, "Yn", factory.FuncIif64Of64(math.Yn)),
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
