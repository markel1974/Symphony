package sdk

import (
	"math"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Math serves as a container for mathematical operations and modules, mapping module names to their respective objects.
type Math struct {
	factory *objects.Factory
	*Package
}

// NewMath initializes and returns a new instance of Math with predefined mathematical constants and function modules.
func NewMath(factory *objects.Factory) *Math {
	m := &Math{
		factory: factory,
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
	container := []*objects.FuncPackage{
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Abs", factory.FuncIf64Of64(math.Abs)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Acos", factory.FuncIf64Of64(math.Acos)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Acosh", factory.FuncIf64Of64(math.Acosh)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Asin", factory.FuncIf64Of64(math.Asin)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Asinh", factory.FuncIf64Of64(math.Asinh)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Atan", factory.FuncIf64Of64(math.Atan)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Atan2", factory.FuncIf64f64Of64(math.Atan2)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Atanh", factory.FuncIf64Of64(math.Atanh)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Cbrt", factory.FuncIf64Of64(math.Cbrt)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Ceil", factory.FuncIf64Of64(math.Ceil)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Copysign", factory.FuncIf64f64Of64(math.Copysign)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Cos", factory.FuncIf64Of64(math.Cos)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Cosh", factory.FuncIf64Of64(math.Cosh)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Dim", factory.FuncIf64f64Of64(math.Dim)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Erf", factory.FuncIf64Of64(math.Erf)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Erfc", factory.FuncIf64Of64(math.Erfc)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Exp", factory.FuncIf64Of64(math.Exp)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Exp2", factory.FuncIf64Of64(math.Exp2)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Expm1", factory.FuncIf64Of64(math.Expm1)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Floor", factory.FuncIf64Of64(math.Floor)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Gamma", factory.FuncIf64Of64(math.Gamma)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Hypot", factory.FuncIf64f64Of64(math.Hypot)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Ilogb", factory.FuncIf64Oi(math.Ilogb)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Inf", factory.FuncIiOf64(math.Inf)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "IsInf", factory.FuncIf64iOb(math.IsInf)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "IsNaN", factory.FuncIf64Ob(math.IsNaN)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "J0", factory.FuncIf64Of64(math.J0)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "J1", factory.FuncIf64Of64(math.J1)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Jn", factory.FuncIif64Of64(math.Jn)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Ldexp", factory.FuncIf64iOf64(math.Ldexp)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Log", factory.FuncIf64Of64(math.Log)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Log10", factory.FuncIf64Of64(math.Log10)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Log1p", factory.FuncIf64Of64(math.Log1p)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Log2", factory.FuncIf64Of64(math.Log2)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Logb", factory.FuncIf64Of64(math.Logb)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Max", factory.FuncIf64f64Of64(math.Max)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Min", factory.FuncIf64f64Of64(math.Min)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Mod", factory.FuncIf64f64Of64(math.Mod)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "NaN", factory.FuncInOf64(math.NaN)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Nextafter", factory.FuncIf64f64Of64(math.Nextafter)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Pow", factory.FuncIf64f64Of64(math.Pow)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Pow10", factory.FuncIiOf64(math.Pow10)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Remainder", factory.FuncIf64f64Of64(math.Remainder)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Signbit", factory.FuncIf64Ob(math.Signbit)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Sin", factory.FuncIf64Of64(math.Sin)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Sinh", factory.FuncIf64Of64(math.Sinh)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Sqrt", factory.FuncIf64Of64(math.Sqrt)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Tan", factory.FuncIf64Of64(math.Tan)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Tanh", factory.FuncIf64Of64(math.Tanh)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Trunc", factory.FuncIf64Of64(math.Trunc)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Y0", factory.FuncIf64Of64(math.Y0)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Y1", factory.FuncIf64Of64(math.Y1)),
		factory.NewFuncPackage(objects.FrameStatic, objects.FuncPackageDef, "Yn", factory.FuncIif64Of64(math.Yn)),
	}
	m.Package = NewPackage("math", container, constants)
	return m
}
