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
		factory.NewFuncImport(objects.FrameStatic, "NaN", m.nan),
		factory.NewFuncImport(objects.FrameStatic, "Abs", m.abs),
		factory.NewFuncImport(objects.FrameStatic, "Acos", m.acos),
		factory.NewFuncImport(objects.FrameStatic, "Acosh", factory.FuncIf64Of64(math.Acosh)),
		factory.NewFuncImport(objects.FrameStatic, "Asin", factory.FuncIf64Of64(math.Asin)),
		factory.NewFuncImport(objects.FrameStatic, "Asinh", factory.FuncIf64Of64(math.Asinh)),
		factory.NewFuncImport(objects.FrameStatic, "Atan", factory.FuncIf64Of64(math.Atan)),
		factory.NewFuncImport(objects.FrameStatic, "Atan2", factory.FuncIf64f64Of64(math.Atan2)),
		factory.NewFuncImport(objects.FrameStatic, "Atanh", factory.FuncIf64Of64(math.Atanh)),
		factory.NewFuncImport(objects.FrameStatic, "Cbrt", factory.FuncIf64Of64(math.Cbrt)),
		factory.NewFuncImport(objects.FrameStatic, "Ceil", factory.FuncIf64Of64(math.Ceil)),
		factory.NewFuncImport(objects.FrameStatic, "Copysign", factory.FuncIf64f64Of64(math.Copysign)),
		factory.NewFuncImport(objects.FrameStatic, "Cos", factory.FuncIf64Of64(math.Cos)),
		factory.NewFuncImport(objects.FrameStatic, "Cosh", factory.FuncIf64Of64(math.Cosh)),
		factory.NewFuncImport(objects.FrameStatic, "Dim", factory.FuncIf64f64Of64(math.Dim)),
		factory.NewFuncImport(objects.FrameStatic, "Erf", factory.FuncIf64Of64(math.Erf)),
		factory.NewFuncImport(objects.FrameStatic, "Erfc", factory.FuncIf64Of64(math.Erfc)),
		factory.NewFuncImport(objects.FrameStatic, "Exp", factory.FuncIf64Of64(math.Exp)),
		factory.NewFuncImport(objects.FrameStatic, "Exp2", factory.FuncIf64Of64(math.Exp2)),
		factory.NewFuncImport(objects.FrameStatic, "Expm1", factory.FuncIf64Of64(math.Expm1)),
		factory.NewFuncImport(objects.FrameStatic, "Floor", factory.FuncIf64Of64(math.Floor)),
		factory.NewFuncImport(objects.FrameStatic, "Gamma", factory.FuncIf64Of64(math.Gamma)),
		factory.NewFuncImport(objects.FrameStatic, "Hypot", factory.FuncIf64f64Of64(math.Hypot)),
		factory.NewFuncImport(objects.FrameStatic, "Ilogb", factory.FuncIf64Oi(math.Ilogb)),
		factory.NewFuncImport(objects.FrameStatic, "Inf", factory.FuncIiOf64(math.Inf)),
		factory.NewFuncImport(objects.FrameStatic, "IsInf", factory.FuncIf64iOb(math.IsInf)),
		factory.NewFuncImport(objects.FrameStatic, "IsNaN", factory.FuncIf64Ob(math.IsNaN)),
		factory.NewFuncImport(objects.FrameStatic, "J0", factory.FuncIf64Of64(math.J0)),
		factory.NewFuncImport(objects.FrameStatic, "J1", factory.FuncIf64Of64(math.J1)),
		factory.NewFuncImport(objects.FrameStatic, "Jn", factory.FuncIif64Of64(math.Jn)),
		factory.NewFuncImport(objects.FrameStatic, "Ldexp", factory.FuncIf64iOf64(math.Ldexp)),
		factory.NewFuncImport(objects.FrameStatic, "Log", factory.FuncIf64Of64(math.Log)),
		factory.NewFuncImport(objects.FrameStatic, "Log10", factory.FuncIf64Of64(math.Log10)),
		factory.NewFuncImport(objects.FrameStatic, "Log1p", factory.FuncIf64Of64(math.Log1p)),
		factory.NewFuncImport(objects.FrameStatic, "Log2", factory.FuncIf64Of64(math.Log2)),
		factory.NewFuncImport(objects.FrameStatic, "Logb", factory.FuncIf64Of64(math.Logb)),
		factory.NewFuncImport(objects.FrameStatic, "Max", factory.FuncIf64f64Of64(math.Max)),
		factory.NewFuncImport(objects.FrameStatic, "Min", factory.FuncIf64f64Of64(math.Min)),
		factory.NewFuncImport(objects.FrameStatic, "Mod", factory.FuncIf64f64Of64(math.Mod)),
		factory.NewFuncImport(objects.FrameStatic, "Nextafter", factory.FuncIf64f64Of64(math.Nextafter)),
		factory.NewFuncImport(objects.FrameStatic, "Pow", factory.FuncIf64f64Of64(math.Pow)),
		factory.NewFuncImport(objects.FrameStatic, "Pow10", factory.FuncIiOf64(math.Pow10)),
		factory.NewFuncImport(objects.FrameStatic, "Remainder", factory.FuncIf64f64Of64(math.Remainder)),
		factory.NewFuncImport(objects.FrameStatic, "Signbit", factory.FuncIf64Ob(math.Signbit)),
		factory.NewFuncImport(objects.FrameStatic, "Sin", factory.FuncIf64Of64(math.Sin)),
		factory.NewFuncImport(objects.FrameStatic, "Sinh", factory.FuncIf64Of64(math.Sinh)),
		factory.NewFuncImport(objects.FrameStatic, "Sqrt", factory.FuncIf64Of64(math.Sqrt)),
		factory.NewFuncImport(objects.FrameStatic, "Tan", factory.FuncIf64Of64(math.Tan)),
		factory.NewFuncImport(objects.FrameStatic, "Tanh", factory.FuncIf64Of64(math.Tanh)),
		factory.NewFuncImport(objects.FrameStatic, "Trunc", factory.FuncIf64Of64(math.Trunc)),
		factory.NewFuncImport(objects.FrameStatic, "Y0", factory.FuncIf64Of64(math.Y0)),
		factory.NewFuncImport(objects.FrameStatic, "Y1", factory.FuncIf64Of64(math.Y1)),
		factory.NewFuncImport(objects.FrameStatic, "Yn", factory.FuncIif64Of64(math.Yn)),
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

func (m *Math) abs(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	ret := gk.UndefinedValue()
	err := objects.ErrInvalidArgumentsNumber
	if len(args) == 1 {
		var f1 float64
		if f1, err = gk.ToFloat64Arg(0, args[0]); err == nil {
			ret = gk.NewFloat(frame, math.Abs(f1))
		}
	}
	return 1, ret, err
}

func (m *Math) acos(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	ret := gk.UndefinedValue()
	err := objects.ErrInvalidArgumentsNumber
	if len(args) == 1 {
		var f1 float64
		if f1, err = gk.ToFloat64Arg(0, args[0]); err == nil {
			ret = gk.NewFloat(frame, math.Acos(f1))
		}
	}
	return 1, ret, err
}
