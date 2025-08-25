package sdk

import (
	"math"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Math serves as a container for mathematical operations and modules, mapping module names to their respective objects.
type Math struct {
	gk objects.IGateKeeper
	*Package
}

// NewMath initializes and returns a new instance of Math with predefined mathematical constants and function modules.
func NewMath(factory objects.IGateKeeper) *Math {
	m := &Math{
		gk: factory,
	}
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
		factory.NewFuncPackage(objects.FuncPackageDef, "Abs", factory.FuncIf64Of64(math.Abs)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Acos", factory.FuncIf64Of64(math.Acos)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Acosh", factory.FuncIf64Of64(math.Acosh)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Asin", factory.FuncIf64Of64(math.Asin)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Asinh", factory.FuncIf64Of64(math.Asinh)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Atan", factory.FuncIf64Of64(math.Atan)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Atan2", factory.FuncIf64f64Of64(math.Atan2)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Atanh", factory.FuncIf64Of64(math.Atanh)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Cbrt", factory.FuncIf64Of64(math.Cbrt)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Ceil", factory.FuncIf64Of64(math.Ceil)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Copysign", factory.FuncIf64f64Of64(math.Copysign)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Cos", factory.FuncIf64Of64(math.Cos)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Cosh", factory.FuncIf64Of64(math.Cosh)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Dim", factory.FuncIf64f64Of64(math.Dim)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Erf", factory.FuncIf64Of64(math.Erf)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Erfc", factory.FuncIf64Of64(math.Erfc)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Exp", factory.FuncIf64Of64(math.Exp)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Exp2", factory.FuncIf64Of64(math.Exp2)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Expm1", factory.FuncIf64Of64(math.Expm1)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Floor", factory.FuncIf64Of64(math.Floor)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Gamma", factory.FuncIf64Of64(math.Gamma)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Hypot", factory.FuncIf64f64Of64(math.Hypot)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Ilogb", factory.FuncIf64Oi(math.Ilogb)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Inf", factory.FuncIiOf64(math.Inf)),
		factory.NewFuncPackage(objects.FuncPackageDef, "IsInf", factory.FuncIf64iOb(math.IsInf)),
		factory.NewFuncPackage(objects.FuncPackageDef, "IsNaN", factory.FuncIf64Ob(math.IsNaN)),
		factory.NewFuncPackage(objects.FuncPackageDef, "J0", factory.FuncIf64Of64(math.J0)),
		factory.NewFuncPackage(objects.FuncPackageDef, "J1", factory.FuncIf64Of64(math.J1)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Jn", factory.FuncIif64Of64(math.Jn)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Ldexp", factory.FuncIf64iOf64(math.Ldexp)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Log", factory.FuncIf64Of64(math.Log)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Log10", factory.FuncIf64Of64(math.Log10)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Log1p", factory.FuncIf64Of64(math.Log1p)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Log2", factory.FuncIf64Of64(math.Log2)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Logb", factory.FuncIf64Of64(math.Logb)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Max", factory.FuncIf64f64Of64(math.Max)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Min", factory.FuncIf64f64Of64(math.Min)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Mod", factory.FuncIf64f64Of64(math.Mod)),
		factory.NewFuncPackage(objects.FuncPackageDef, "NaN", m.funcInOf64(math.NaN)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Nextafter", factory.FuncIf64f64Of64(math.Nextafter)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Pow", factory.FuncIf64f64Of64(math.Pow)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Pow10", factory.FuncIiOf64(math.Pow10)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Remainder", factory.FuncIf64f64Of64(math.Remainder)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Signbit", factory.FuncIf64Ob(math.Signbit)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sin", factory.FuncIf64Of64(math.Sin)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sinh", factory.FuncIf64Of64(math.Sinh)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sqrt", factory.FuncIf64Of64(math.Sqrt)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Tan", factory.FuncIf64Of64(math.Tan)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Tanh", factory.FuncIf64Of64(math.Tanh)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Trunc", factory.FuncIf64Of64(math.Trunc)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Y0", factory.FuncIf64Of64(math.Y0)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Y1", factory.FuncIf64Of64(math.Y1)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Yn", factory.FuncIif64Of64(math.Yn)),
	}
	m.Package = NewPackage("math", container, constants)
	return m
}

// funcInOf64 wraps a no-argument function returning float64 into a FuncCallable to integrate with the GateAdapter system.
// It enforces no arguments and converts the float64 result to an IObject, returning ErrWrongNumArguments for invalid input.
func (z *Math) funcInOf64(fn func() float64) objects.FuncCallable {
	return func(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, objects.ErrWrongNumArguments
		}
		return z.gk.NewFloat(frame, fn()), nil
	}
}
