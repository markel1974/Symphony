package sdk

import (
	"math"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Math serves as a container for mathematical operations and modules, mapping module names to their respective objects.
type Math struct {
	*Package
}

// NewMath initializes and returns a new instance of Math with predefined mathematical constants and function modules.
func NewMath() *Math {
	m := &Math{}
	constants := map[string]objects.IObject{
		"E":       objects.NewFloat(math.E),
		"Pi":      objects.NewFloat(math.Pi),
		"Phi":     objects.NewFloat(math.Phi),
		"Sqrt2":   objects.NewFloat(math.Sqrt2),
		"SqrtE":   objects.NewFloat(math.SqrtE),
		"SqrtPi":  objects.NewFloat(math.SqrtPi),
		"SqrtPhi": objects.NewFloat(math.SqrtPhi),
		"Ln2":     objects.NewFloat(math.Ln2),
		"Log2E":   objects.NewFloat(math.Log2E),
		"Ln10":    objects.NewFloat(math.Ln10),
		"Log10E":  objects.NewFloat(math.Log10E),
	}
	container := []*objects.FuncPackage{
		objects.NewFuncPackage(objects.FuncPackageDef, "Abs", objects.FuncIf64Of64(math.Abs)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Acos", objects.FuncIf64Of64(math.Acos)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Acosh", objects.FuncIf64Of64(math.Acosh)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Asin", objects.FuncIf64Of64(math.Asin)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Asinh", objects.FuncIf64Of64(math.Asinh)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Atan", objects.FuncIf64Of64(math.Atan)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Atan2", objects.FuncIf64f64Of64(math.Atan2)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Atanh", objects.FuncIf64Of64(math.Atanh)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Cbrt", objects.FuncIf64Of64(math.Cbrt)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Ceil", objects.FuncIf64Of64(math.Ceil)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Copysign", objects.FuncIf64f64Of64(math.Copysign)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Cos", objects.FuncIf64Of64(math.Cos)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Cosh", objects.FuncIf64Of64(math.Cosh)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Dim", objects.FuncIf64f64Of64(math.Dim)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Erf", objects.FuncIf64Of64(math.Erf)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Erfc", objects.FuncIf64Of64(math.Erfc)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Exp", objects.FuncIf64Of64(math.Exp)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Exp2", objects.FuncIf64Of64(math.Exp2)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Expm1", objects.FuncIf64Of64(math.Expm1)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Floor", objects.FuncIf64Of64(math.Floor)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Gamma", objects.FuncIf64Of64(math.Gamma)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Hypot", objects.FuncIf64f64Of64(math.Hypot)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Ilogb", objects.FuncIf64Oi(math.Ilogb)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Inf", objects.FuncIiOf64(math.Inf)),
		objects.NewFuncPackage(objects.FuncPackageDef, "IsInf", objects.FuncIf64iOb(math.IsInf)),
		objects.NewFuncPackage(objects.FuncPackageDef, "IsNaN", objects.FuncIf64Ob(math.IsNaN)),
		objects.NewFuncPackage(objects.FuncPackageDef, "J0", objects.FuncIf64Of64(math.J0)),
		objects.NewFuncPackage(objects.FuncPackageDef, "J1", objects.FuncIf64Of64(math.J1)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Jn", objects.FuncIif64Of64(math.Jn)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Ldexp", objects.FuncIf64iOf64(math.Ldexp)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Log", objects.FuncIf64Of64(math.Log)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Log10", objects.FuncIf64Of64(math.Log10)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Log1p", objects.FuncIf64Of64(math.Log1p)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Log2", objects.FuncIf64Of64(math.Log2)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Logb", objects.FuncIf64Of64(math.Logb)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Max", objects.FuncIf64f64Of64(math.Max)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Min", objects.FuncIf64f64Of64(math.Min)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Mod", objects.FuncIf64f64Of64(math.Mod)),
		objects.NewFuncPackage(objects.FuncPackageDef, "NaN", objects.FuncInOf64(math.NaN)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Nextafter", objects.FuncIf64f64Of64(math.Nextafter)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Pow", objects.FuncIf64f64Of64(math.Pow)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Pow10", objects.FuncIiOf64(math.Pow10)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Remainder", objects.FuncIf64f64Of64(math.Remainder)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Signbit", objects.FuncIf64Ob(math.Signbit)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Sin", objects.FuncIf64Of64(math.Sin)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Sinh", objects.FuncIf64Of64(math.Sinh)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Sqrt", objects.FuncIf64Of64(math.Sqrt)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Tan", objects.FuncIf64Of64(math.Tan)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Tanh", objects.FuncIf64Of64(math.Tanh)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Trunc", objects.FuncIf64Of64(math.Trunc)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Y0", objects.FuncIf64Of64(math.Y0)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Y1", objects.FuncIf64Of64(math.Y1)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Yn", objects.FuncIif64Of64(math.Yn)),
	}
	m.Package = NewPackage("math", container, constants)
	return m
}
