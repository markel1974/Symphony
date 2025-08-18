package sdk

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Rand is a struct that encapsulates a module mapping of string keys to objects implementing the IObject interface.
type Rand struct {
	*Module
}

// NewRand creates a new instance of Rand with a pre-defined set of random number generation functions.
func NewRand() *Rand {
	z := &Rand{
		Module: NewModule(),
	}
	z.attrs = map[string]objects.IObject{
		"Int63":       objects.NewFunctionModule(objects.FunctionModuleDef, "Int63", objects.FuncInOi64(rand.Int63)),
		"Float64":     objects.NewFunctionModule(objects.FunctionModuleDef, "Float64", objects.FuncInOf64(rand.Float64)),
		"Int63n":      objects.NewFunctionModule(objects.FunctionModuleDef, "Int63n", objects.FuncIi64Oi64(rand.Int63n)),
		"ExpFloat64":  objects.NewFunctionModule(objects.FunctionModuleDef, "ExpFloat64", objects.FuncInOf64(rand.ExpFloat64)),
		"NormFloat64": objects.NewFunctionModule(objects.FunctionModuleDef, "NormFloat64", objects.FuncInOf64(rand.NormFloat64)),
		"Perm":        objects.NewFunctionModule(objects.FunctionModuleDef, "Perm", objects.FuncIiOiS(rand.Perm)),
		"Seed":        objects.NewFunctionModule(objects.FunctionModuleDef, "Seed", objects.FuncIi64On(rand.Seed)),
		"Read":        objects.NewFunctionModule(objects.FunctionModuleDef, "Read", z.Read),
		"Rand":        objects.NewFunctionModule(objects.FunctionModuleDef, "Rand", z.Rand),
	}
	return z
}

// Name returns the name of Rand module.
func (z *Rand) Name() string {
	return "rand"
}

// Read reads random data into a byte slice and returns the number of bytes written as an integer or an error if it occurs.
func (z *Rand) Read(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := objects.ToByteSliceArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := rand.Read(bs1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(int64(res)), nil
}

// Rand generates a new random number generator using the provided seed argument and returns its options as a map.
func (z *Rand) Rand(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	src := rand.NewSource(i1)
	return z.RandOptions(rand.New(src)), nil
}

// RandOptions returns an immutable map of utility functions for the provided *rand.Rand object, enabling random operations.
func (z *Rand) RandOptions(r *rand.Rand) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"Int63":       objects.NewFunctionModule(objects.FunctionModuleDef, "Int63", objects.FuncInOi64(r.Int63)),
			"Float64":     objects.NewFunctionModule(objects.FunctionModuleDef, "Float64", objects.FuncInOf64(r.Float64)),
			"Int63n":      objects.NewFunctionModule(objects.FunctionModuleDef, "Int63n", objects.FuncIi64Oi64(r.Int63n)),
			"ExpFloat64":  objects.NewFunctionModule(objects.FunctionModuleDef, "ExpFloat64", objects.FuncInOf64(r.ExpFloat64)),
			"NormFloat64": objects.NewFunctionModule(objects.FunctionModuleDef, "NormFloat64", objects.FuncInOf64(r.NormFloat64)),
			"Perm":        objects.NewFunctionModule(objects.FunctionModuleDef, "Perm", objects.FuncIiOiS(r.Perm)),
			"Seed":        objects.NewFunctionModule(objects.FunctionModuleDef, "Seed", objects.FuncIi64On(r.Seed)),
			"Read":        objects.NewFunctionModule(objects.FunctionModuleDef, "Read", func(args ...objects.IObject) (objects.IObject, error) { return z.RandOptionsRead(r, args...) }),
		})
}

// RandOptionsRead reads random data into a provided byte slice, returning the number of bytes read or an error.
func (z *Rand) RandOptionsRead(r *rand.Rand, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := objects.ToByteSliceArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := r.Read(bs1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(int64(res)), nil
}
